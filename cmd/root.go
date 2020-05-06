/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// GetRootCommand returns a command that represents the base command when called without any subcommands,
// adds all child commands to the root command and sets flags appropriately.
func GetRootCommand() *cobra.Command {

	var rootCmd = &cobra.Command{
		Use:   "build-tools",
		Short: "build-tools : A swiss-army knife of CI/CI related commands",
		Long: ``,
	}
	initRootCommand(rootCmd)
	return rootCmd
}

// initRootCommand defines the flags, persistent flags and configuration settings
// for the root command and adds all sub commands.
func initRootCommand(cmd *cobra.Command) {

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.build-tools.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//cmd.PFlags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.AddCommand(GetGhCommand())

}
