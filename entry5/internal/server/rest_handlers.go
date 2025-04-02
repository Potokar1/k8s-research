package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Potokar1/k8s-research/entry5/internal/worker"
)

type Server struct {
	client *http.Client

	worker *worker.Worker
}

func NewServer(w *worker.Worker) *Server {
	return &Server{
		client: http.DefaultClient,
		worker: w,
	}
}

func (s *Server) InitializeREST(ctx context.Context, mux *http.ServeMux) {
	// live and ready checks
	mux.HandleFunc("/live", s.restLive)
	mux.HandleFunc("/ready", s.restReady)

	// worker endpoints
	mux.HandleFunc("/sell", s.restSell)
	mux.HandleFunc("/inventory", s.restInventory)
}

// restLive implements the REST API for the live check
func (s *Server) restLive(w http.ResponseWriter, r *http.Request) {
	// Respond okay
	w.WriteHeader(http.StatusOK)
}

// restReady implements the REST API for the ready check
// The worker is ready if it has an inventory greater than a set min.
func (s *Server) restReady(w http.ResponseWriter, r *http.Request) {
	if s.worker.AboveMinimum() {
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "Not ready", http.StatusServiceUnavailable)
}

// restSell implement the REST API for selling items from the worker
func (s *Server) restSell(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Handle the buy request
	buyRequest, err := worker.DecodeBuyRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.DebugContext(ctx, "received sell request", "item", buyRequest.Item, "quantity", buyRequest.Quantity)

	// Sell the item(s)
	if sold := s.worker.Sell(ctx, buyRequest.Item, buyRequest.Quantity); !sold {
		http.Error(w, "Not enough inventory", http.StatusConflict)
		return

	}

	// Respond okay
	w.WriteHeader(http.StatusOK)
}

// restInventory implements the REST API for getting the inventory of the worker
func (s *Server) restInventory(w http.ResponseWriter, r *http.Request) {
	// Get the inventory
	invList := s.worker.InventoryList()

	// Respond with the inventory as JSON
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(invList); err != nil {
		slog.Debug("error encoding inventory", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
