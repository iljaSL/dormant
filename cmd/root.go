package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dormant",
		Short: "A tool to analyze go.mod files for inactive Go packages",
		Long: `Dormant is a CLI tool for analyzing go.mod files for inactive Go packages.
This tool should help Go developers to analyze their go.mod files and report packages
that are not being actively maintained anymore.
		`,
		Run:     root,
		Version: "0.1",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}
