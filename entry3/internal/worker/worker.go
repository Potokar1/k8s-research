package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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

// Direction is a struct that represents what a worker can do and how often
type Direction struct {
	Product  string
	Amount   int
	Minimum  int
	Interval int
}

func NewWorker(directions []Direction) *Worker {
	return &Worker{
		Inventory:  make(map[string]int),
		directions: directions,
	}
}

type BuyRequest struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

func DecodeBuyRequest(r *http.Request) (*BuyRequest, error) {
	var req BuyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return &req, nil
}

func (w *Worker) Buy(item string, quantity int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Inventory[item] += quantity
}

func (w *Worker) Sell(item string, quantity int) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Inventory[item] < quantity {
		return false
	}
	w.Inventory[item] -= quantity
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
	for k, v := range w.Inventory {
		invList[k] = v
	}
	return invList
}

// produce increments the inventory of a product by a set amount
func (w *Worker) produce(direction Direction) {
	w.mu.Lock()
	defer w.mu.Unlock()
	slog.Info("producing", "product", direction.Product, "amount", direction.Amount)
	w.Inventory[direction.Product] += direction.Amount
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
					w.produce(direction)
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
