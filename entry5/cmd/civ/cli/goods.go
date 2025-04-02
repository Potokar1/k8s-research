package cli

import (
	"github.com/Potokar1/k8s-research/entry5/internal/k8s"
	"github.com/spf13/cobra"
)

// NewGoodsCmd creates the goods command
func NewGoodsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "goods",
		Short: "View total goods give query",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get flag values for kingdom, town, and shop
			kingdom, err := cmd.Flags().GetString("kingdom")
			if err != nil {
				return err
			}
			townName, err := cmd.Flags().GetString("town")
			if err != nil {
				return err
			}
			shopName, err := cmd.Flags().GetString("shop")
			if err != nil {
				return err
			}
			return k8s.GetInventory(cmd.Context(), kingdom, townName, shopName)
		},
	}

	// flags for kingdom, town, and shop, all optional
	cmd.Flags().String("kingdom", "", "Filter by kingdom")
	cmd.Flags().String("town", "", "Filter by town")
	cmd.Flags().String("shop", "", "Filter by shop")
	cmd.MarkFlagRequired("kingdom")
	cmd.RegisterFlagCompletionFunc("kingdom", KingdomsValidArgsFunction)
	cmd.RegisterFlagCompletionFunc("town", TownsValidArgsFunction)
	cmd.RegisterFlagCompletionFunc("shop", ShopsValidArgsFunction)

	return cmd
}
