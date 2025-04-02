package cli

import (
	"strings"

	"github.com/Potokar1/k8s-research/entry4/internal/k8s"
	"github.com/spf13/cobra"
)

// NewKingdomsCmd creates the kingdoms command
func NewKingdomsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kingdoms",
		Short: "kingdoms is a CLI tool for managing k8s kingdoms",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// list kingdoms
			kingdoms, err := k8s.ListNamespaces(cmd.Context())
			if err != nil {
				return err
			}

			// print kingdoms that start with "kingdom-of-"
			for _, kingdom := range kingdoms {
				if strings.HasPrefix(kingdom, "kingdom-of-") {
					cmd.Println(kingdom)
				}
			}

			return nil
		},
	}

	return cmd
}
