package cli

import "github.com/spf13/cobra"

// NewCivCmd creates the base civ command
func NewCivCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "civ",
		Short: "civ is a CLI tool for managing k8s resources",
	}

	cmd.AddCommand(NewKingdomsCmd())

	return cmd
}
