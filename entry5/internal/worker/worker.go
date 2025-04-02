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
)

type Worker struct {
	mu         sync.Mutex
	Inventory  map[string]int
	directions []Direction
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

func NewWorker(directions []Direction) *Worker {
	return &Worker{
		Inventory:  make(map[string]int),
		directions: directions,
	}
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
		w.Inventory[item.Product] += item.Amount
		slog.InfoContext(ctx, "Purchased", "product", item.Product, "amount", item.Amount, "Inventory", w.Inventory)
		return true
	default:
		// any other non-200 status code is treated as an error
		slog.DebugContext(ctx, "received non-200 status code from store", "status_code", resp.StatusCode, "store", item.Store)
		return false
	}
}

func (w *Worker) Sell(ctx context.Context, item string, quantity int) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Inventory[item] < quantity {
		return false
	}
	w.Inventory[item] -= quantity
	slog.InfoContext(ctx, "Sold", "item", item, "amount", quantity, "Inventory", w.Inventory)
	return true
}

func (w *Worker) AboveMinimum() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, direction := range w.directions {
		if w.Inventory[direction.Product] < direction.Minimum {
			return false
		}
	}
	return true
}

func (w *Worker) InventoryCount(item string) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Inventory[item]
}

func (w *Worker) InventoryList() map[string]int {
	w.mu.Lock()
	defer w.mu.Unlock()
	invList := make(map[string]int, len(w.Inventory))
	maps.Copy(invList, w.Inventory)
	return invList
}

// produce increments the inventory of a product by a set amount
func (w *Worker) produce(ctx context.Context, direction Direction) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// check that we have enough inventory to produce the product
	for _, input := range direction.ProductInputList {
		if w.Inventory[input.Product] < input.Amount {
			// attempt to buy the missing inputs
			slog.DebugContext(ctx, "buying product input from store", "product", direction.Product, "input", input.Product, "store", input.Store, "amount", input.Amount)
			if bought := w.buy(ctx, input); !bought {
				slog.DebugContext(ctx, "failed to buy input", "product", input.Product, "store", input.Store, "amount", input.Amount)
			}
		}
	}

	// use inputs to make the product
	for _, input := range direction.ProductInputList {
		// decrement the inventory of the input product
		if w.Inventory[input.Product] < input.Amount {
			slog.DebugContext(ctx, "not enough inventory to produce product", "product", direction.Product, "missing_input", input.Product, "required_amount", input.Amount)
			return
		}
		// decrement the inventory of the input product
		w.Inventory[input.Product] -= input.Amount
	}

	w.Inventory[direction.Product] += direction.Amount
	slog.InfoContext(ctx, "Produced", "product", direction.Product, "amount", direction.Amount, "inventory", w.Inventory)
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
			time.Sleep(1 * time.Second)
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
