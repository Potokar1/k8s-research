package cli

import (
	"context"
	"fmt"

	"github.com/Potokar1/k8s-research/entry4/internal/k8s"
	"github.com/spf13/cobra"
)

// NewShopsCmd creates the shops command
func NewShopsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shops",
		Short: "List all shops in a town",
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
			shops, err := k8s.ListPods(context.Background(), kingdom, townName)
			if err != nil {
				return err
			}
			for _, shop := range shops {
				fmt.Println(shop)
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
