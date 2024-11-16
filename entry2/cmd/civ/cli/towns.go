package cli

import (
	"context"
	"fmt"

	"github.com/Potokar1/k8s-research/entry2/internal/k8s"
	"github.com/spf13/cobra"
)

// NewTownsCmd creates the towns command
func NewTownsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "towns",
		Short: "List all towns in a kingdom",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			kingdom, err := cmd.Flags().GetString("kingdom")
			if err != nil {
				return err
			}
			towns, err := k8s.ListPods(context.Background(), kingdom)
			if err != nil {
				return err
			}
			for _, town := range towns {
				fmt.Println(town)
			}
			return nil
		},
	}

	cmd.Flags().String("kingdom", "", "Select which kingdom to list towns from")
	cmd.MarkFlagRequired("kingdom")
	cmd.RegisterFlagCompletionFunc("kingdom", KingdomsValidArgsFunction)

	return cmd
}
