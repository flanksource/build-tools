/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// GetGhCommand returns the github (gh) command, adds all child commands and sets flags appropriately.
func GetGhCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "gh",
		Short: "github related actions",
		Long: ``,
		//RunE: func(cmd *cobra.Command, args []string) error { return fmt.Errorf("test")},
	}
	initGhCommand(cmd)
	return cmd
}

// initGhCommand defines the flags, persistent flags and configuration settings
// for the gh command and adds all sub commands.
func initGhCommand(cmd *cobra.Command) {

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.build-tools.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	cmd.AddCommand(GetReportJUnitCommand())

}
