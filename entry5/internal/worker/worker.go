package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Potokar1/k8s-research/entry5/internal/k8s"
)

type Worker struct {
	kingdom string // Kingdom is the namespace the worker belongs to
	name    string // Name is the name of the pod

	directions []Direction

	inventoryLock sync.RWMutex
	inventory     map[string]int
}

type ProductInput struct {
	Product string // Product is the name of the product to buy
	Store   string // Store is the URL of the store to buy from
	Amount  int    // Amount is the quantity of the product to buy
}

// Direction is a struct that represents what a worker can do and how often
type Direction struct {
	Product          string         // Product is the name of the product to produce
	ProductInputList []ProductInput // ProductInputList is a list of inputs required to produce the product
	Amount           int            // Amount is the amount of product to produce at each interval
	Minimum          int            // Minimum is the minimum amount of product to keep in inventory
	Interval         int            // Interval is the rate in seconds at which the product should be produced
}

func NewWorker(kingdom, name string, directions []Direction) *Worker {
	return &Worker{
		kingdom:    kingdom,
		name:       name,
		inventory:  make(map[string]int),
		directions: directions,
	}
}

func (w *Worker) UpdateStoreLog(ctx context.Context) error {
	// Get inventory list
	invList := w.InventoryList()

	// Patch Pod
	return k8s.PatchPod(ctx, w.kingdom, w.name, invList)
}

func (w *Worker) addInventory(ctx context.Context, item string, amount int) {
	w.inventoryLock.Lock()
	if _, exists := w.inventory[item]; !exists {
		w.inventory[item] = amount
	}
	w.inventory[item] += amount
	w.inventoryLock.Unlock()

	// patch the pod with the new inventory
	if err := w.UpdateStoreLog(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to update store log", "error", err)
	}
}

func (w *Worker) ifEnoughInventory(ctx context.Context, item string, amount int) bool {
	w.inventoryLock.RLock()
	defer w.inventoryLock.RUnlock()
	if _, exists := w.inventory[item]; !exists {
		slog.DebugContext(ctx, "Attempted to check inventory for item that does not exist", "item", item)
		return false
	}
	if w.inventory[item] < amount {
		slog.DebugContext(ctx, "Not enough inventory for item", "item", item, "requested_amount", amount, "available_amount", w.inventory[item])
		return false
	}
	return true
}

// returns true if item was removed from inventory, false if failure
func (w *Worker) removeInventory(ctx context.Context, item string, amount int) bool {
	if !w.ifEnoughInventory(ctx, item, amount) {
		return false
	}
	w.inventoryLock.Lock()
	w.inventory[item] -= amount
	slog.DebugContext(ctx, "Removed inventory", "item", item, "amount", amount, "remaining_inventory", w.inventory[item])
	w.inventoryLock.Unlock()

	// patch the pod with the new inventory
	if err := w.UpdateStoreLog(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to update store log", "error", err)
	}
	return true
}

// BuyRequest is the data payload received by another service to buy an item
type BuyRequest struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

// DecodeBuyRequest json decodes the request body into a BuyRequest struct
func DecodeBuyRequest(r *http.Request) (*BuyRequest, error) {
	var req BuyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

// buy allows the worker to buy a product from a store given the ProductInput
// buy is expected to only be called by a locked worker.
func (w *Worker) buy(ctx context.Context, item ProductInput) bool {
	// create a buy request for the item
	BuyRequest := BuyRequest{
		Item:     item.Product,
		Quantity: item.Amount,
	}

	// Marshal the BuyRequest into JSON
	payload, err := json.Marshal(BuyRequest)
	if err != nil {
		slog.ErrorContext(ctx, "failed to marshal BuyRequest", "error", err)
		return false
	}

	// create a new reader for the payload
	reader := bytes.NewReader(payload)

	// Make a request to the store to buy the product
	req, err := http.NewRequest("POST", item.Store+"/sell", reader)
	if err != nil {
		slog.ErrorContext(ctx, "failed to create HTTP request", "error", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.WarnContext(ctx, "failed to send HTTP request", "error", err)
		return false
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.ErrorContext(ctx, "failed to close response body", "error", err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusConflict:
		// Conflict means the store could not fulfill the request due to insufficient inventory
		slog.DebugContext(ctx, "store could not fulfill buy request due to insufficient inventory")
		return false
	case http.StatusOK:
		// if the request was successful, we assume the item was bought
		w.addInventory(ctx, item.Product, item.Amount)
		slog.InfoContext(ctx, "Purchased", "product", item.Product, "amount", item.Amount)
		return true
	default:
		// any other non-200 status code is treated as an error
		slog.DebugContext(ctx, "received non-200 status code from store", "status_code", resp.StatusCode, "store", item.Store)
		return false
	}
}

func (w *Worker) Sell(ctx context.Context, item string, quantity int) bool {
	return w.removeInventory(ctx, item, quantity)
}

func (w *Worker) AboveMinimum() bool {
	w.inventoryLock.RLock()
	defer w.inventoryLock.RUnlock()
	for _, direction := range w.directions {
		if w.inventory[direction.Product] < direction.Minimum {
			return false
		}
	}
	return true
}

func (w *Worker) InventoryList() map[string]string {
	w.inventoryLock.RLock()
	invList := make(map[string]int, len(w.inventory))
	maps.Copy(invList, w.inventory)
	w.inventoryLock.RUnlock()

	invListStr := make(map[string]string, len(invList))
	for item, amount := range invList {
		invListStr[item] = fmt.Sprintf("%d", amount)
	}

	return invListStr
}

// produce increments the inventory of a product by a set amount
func (w *Worker) produce(ctx context.Context, direction Direction) {
	// check that we have enough inventory to produce the product
	for _, input := range direction.ProductInputList {
		if !w.ifEnoughInventory(ctx, input.Product, input.Amount) {
			// attempt to buy the missing inputs
			if bought := w.buy(ctx, input); !bought {
				slog.DebugContext(ctx, "failed to buy input", "product", input.Product, "store", input.Store, "amount", input.Amount)
			}
			slog.DebugContext(ctx, "Bought product input from store", "product", direction.Product, "input", input.Product, "store", input.Store, "amount", input.Amount)
			// only let the workers do one action at a time, so return early
			return
		}
	}

	// use inputs to make the product
	for _, input := range direction.ProductInputList {
		if !w.removeInventory(ctx, input.Product, input.Amount) {
			slog.WarnContext(ctx, "not enough inventory to produce product", "product", direction.Product, "input", input.Product, "amount", input.Amount)
			return // if we can't remove the input, we can't produce the product
		}
	}

	// increment the inventory of the product. This is the worker producing the product
	w.addInventory(ctx, direction.Product, direction.Amount)
	slog.InfoContext(ctx, "Produced product", "product", direction.Product, "amount", direction.Amount)
}

// Work is the loop that will run the worker until the context is canceled
func (w *Worker) Work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for _, direction := range w.directions {
				select {
				case <-time.After(time.Duration(direction.Interval) * time.Second):
					w.produce(ctx, direction)
				case <-ctx.Done():
					return
				}
			}
			// sleep for a second to prevent a busy loop
			time.Sleep(1 * time.Second) // This could also be thought of as the time it takes for a worker to do a task (buy or produce)
		}
	}
}

// ParseDirectionsFile reads a json file and returns a slice of Directions
func ParseDirectionsFile(filename string) ([]Direction, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading worker directions file: %w", err)
	}

	var directions []Direction
	if err := json.Unmarshal(data, &directions); err != nil {
		return nil, fmt.Errorf("error unmarshalling worker directions: %w", err)
	}

	return directions, nil
}
