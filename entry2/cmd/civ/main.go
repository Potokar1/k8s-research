package main

import (
	"os"

	"github.com/Potokar1/k8s-research/entry2/cmd/civ/cli"
)

func main() {
	root := cli.NewCivCmd()

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
