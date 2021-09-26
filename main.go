package main

import (
	"os"

	"github.com/joaom00/gh-discussions/cmd"
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use: "discussion",
	}
}

func main() {
	rc := rootCmd()
	rc.AddCommand(cmd.NewListCmd())

	if err := rc.Execute(); err != nil {
		os.Exit(1)
	}
}
