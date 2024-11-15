package cli

import (
	"strings"

	"github.com/Potokar1/k8s-research/entry1/internal/k8s"
	"github.com/spf13/cobra"
)

// NewKingdomsCmd creates the kingdoms command
func NewKingdomsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kingdoms",
		Short: "kingdoms is a CLI tool for managing k8s namespaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			// list namespaces
			namespaces, err := k8s.ListNamespaces(cmd.Context())
			if err != nil {
				return err
			}

			// print namespaces that start with "kingdom-of-"
			for _, namespace := range namespaces {
				if strings.HasPrefix(namespace, "kingdom-of-") {
					cmd.Println(namespace)
				}
			}

			return nil
		},
	}

	return cmd
}
