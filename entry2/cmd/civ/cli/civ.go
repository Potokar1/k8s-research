package cli

import (
	"context"
	"strings"

	"github.com/Potokar1/k8s-research/entry2/internal/k8s"
	"github.com/spf13/cobra"
)

// NewCivCmd creates the base civ command
func NewCivCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "civ",
		Short: "civ is a CLI tool for managing k8s resources",
	}

	cmd.AddCommand(NewKingdomsCmd())
	cmd.AddCommand(NewTownsCmd())
	cmd.AddCommand(NewWorkersCmd())
	cmd.AddCommand(NewLogsCmd())

	return cmd
}

func KingdomsValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	ctx := context.Background()
	kingdoms, err := k8s.ListNamespaces(ctx)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var comps []string
	for _, kingdom := range kingdoms {
		if strings.HasPrefix(kingdom, "kingdom-of-") && strings.HasPrefix(kingdom, toComplete) {
			comps = append(comps, kingdom)
		}
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

func TownsValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	kingdom, _ := cmd.Flags().GetString("kingdom")
	if kingdom == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ctx := context.Background()
	towns, err := k8s.ListPods(ctx, kingdom)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var comps []string
	for _, town := range towns {
		if strings.HasPrefix(town, toComplete) {
			comps = append(comps, town)
		}
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}

func WorkersValidArgsFunction(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	kingdom, _ := cmd.Flags().GetString("kingdom")
	town, _ := cmd.Flags().GetString("town")
	if kingdom == "" || town == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	ctx := context.Background()
	workers, err := k8s.ListContainers(ctx, kingdom, town)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	var comps []string
	for _, worker := range workers {
		if strings.HasPrefix(worker, toComplete) {
			comps = append(comps, worker)
		}
	}
	return comps, cobra.ShellCompDirectiveNoFileComp
}
