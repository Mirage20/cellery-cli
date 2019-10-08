package main

import (
	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "celleryadm",
		Short: "Manage cellery runtime",
		Long:  `A cli tool to manage your cellery runtime`,
	}
	cmd.AddCommand(
		newInstallCmd(),
	)
	return cmd
}
