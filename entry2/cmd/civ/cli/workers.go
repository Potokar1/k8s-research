package cli

import (
	"context"
	"fmt"

	"github.com/Potokar1/k8s-research/entry2/internal/k8s"
	"github.com/spf13/cobra"
)

// NewWorkersCmd creates the workers command
func NewWorkersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workers",
		Short: "List all workers in a town",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			kingdom, err := cmd.Flags().GetString("kingdom")
			if err != nil {
				return err
			}
			townName, err := cmd.Flags().GetString("town")
			if err != nil {
				return err
			}
			workers, err := k8s.ListContainers(context.Background(), kingdom, townName)
			if err != nil {
				return err
			}
			for _, worker := range workers {
				fmt.Println(worker)
			}
			return nil
		},
	}

	cmd.Flags().String("kingdom", "", "Kingdom of the town")
	cmd.Flags().String("town", "", "Name of the town")
	cmd.MarkFlagRequired("kingdom")
	cmd.MarkFlagRequired("town")
	cmd.RegisterFlagCompletionFunc("kingdom", KingdomsValidArgsFunction)
	cmd.RegisterFlagCompletionFunc("town", TownsValidArgsFunction)

	return cmd
}

// NewLogsCmd creates the logs command
func NewLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Get logs from a worker in a town",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			kingdom, err := cmd.Flags().GetString("kingdom")
			if err != nil {
				return err
			}
			townName, err := cmd.Flags().GetString("town")
			if err != nil {
				return err
			}
			workerName, err := cmd.Flags().GetString("worker")
			if err != nil {
				return err
			}
			logs, err := k8s.GetContainerLogs(context.Background(), kingdom, townName, workerName)
			if err != nil {
				return err
			}
			fmt.Println(logs)
			return nil
		},
	}

	cmd.Flags().String("kingdom", "", "Kingdom of the town")
	cmd.Flags().String("town", "", "Name of the town")
	cmd.Flags().String("worker", "", "Name of the worker")
	cmd.MarkFlagRequired("kingdom")
	cmd.MarkFlagRequired("town")
	cmd.MarkFlagRequired("worker")
	cmd.RegisterFlagCompletionFunc("kingdom", KingdomsValidArgsFunction)
	cmd.RegisterFlagCompletionFunc("town", TownsValidArgsFunction)
	cmd.RegisterFlagCompletionFunc("worker", WorkersValidArgsFunction)

	return cmd
}
