package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Potokar1/k8s-research/entry4/internal/server"
	"github.com/Potokar1/k8s-research/entry4/internal/worker"
	"github.com/spf13/cobra"
)

// NewServeCmd creates the serve command
func NewServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "serve <directions-file>",
		Short:  "Start the server",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancelFunc := context.WithCancel(cmd.Context())
			// go routine that listens for signals
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, os.Interrupt)
			go func() {
				<-sigs
				slog.InfoContext(ctx, "received interrupt signal, shutting down")
				cancelFunc()
			}()

			// create worker
			directions, err := worker.ParseDirectionsFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse directions file: %w", err)
			}
			slog.InfoContext(ctx, "creating worker with directions", "count", len(directions))
			for i, d := range directions {
				slog.InfoContext(ctx, "direction", "index", i, "product", d.Product, "amount", d.Amount, "interval", d.Interval,
					"inputs", d.ProductInputList)
			}
			worker := worker.NewWorker(directions)
			go worker.Work(ctx)

			// create the server
			s := server.NewServer(worker)
			mux := http.DefaultServeMux
			s.InitializeREST(ctx, mux)

			srv := &http.Server{
				Handler: mux,
				Addr:    ":8080",
			}

			go func() {
				slog.InfoContext(ctx, "Listening", "addr", srv.Addr)
				if err := srv.ListenAndServe(); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						slog.ErrorContext(ctx, "serve failed", "error", err)
						panic(err)
					}
				}
			}()

			<-ctx.Done()
			slog.InfoContext(ctx, "Shutting down server")

			// create a deadline to wait for
			timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			slog.InfoContext(ctx, "Waiting for server to shut down")
			if err := srv.Shutdown(timeoutCtx); err != nil {
				return fmt.Errorf("server shutdown failed: %w", err)
			}

			return nil
		},
	}

	return cmd
}
