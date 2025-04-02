package cli

import (
	"context"
	"fmt"

	"github.com/Potokar1/k8s-research/entry4/internal/k8s"
	"github.com/spf13/cobra"
)

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
			_, err = cmd.Flags().GetString("town") // TODO: We don't need town name for this anymore, but the completion function still needs it...
			if err != nil {
				return err
			}
			shopName, err := cmd.Flags().GetString("shop")
			if err != nil {
				return err
			}
			logs, err := k8s.GetContainerLogs(context.Background(), kingdom, shopName)
			if err != nil {
				return err
			}
			fmt.Println(logs)
			return nil
		},
	}

	cmd.Flags().String("kingdom", "", "Kingdom of the town")
	cmd.Flags().String("town", "", "Name of the town")
	cmd.Flags().String("shop", "", "Name of the shop")
	cmd.MarkFlagRequired("kingdom")
	cmd.MarkFlagRequired("town")
	cmd.MarkFlagRequired("shop")
	cmd.RegisterFlagCompletionFunc("kingdom", KingdomsValidArgsFunction)
	cmd.RegisterFlagCompletionFunc("town", TownsValidArgsFunction)
	cmd.RegisterFlagCompletionFunc("shop", ShopsValidArgsFunction)

	return cmd
}
