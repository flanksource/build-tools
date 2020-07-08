package cmd

import (
	"fmt"
	"os"

	"github.com/flanksource/build-tools/pkg/junit"
	"github.com/spf13/cobra"
)

var Junit = &cobra.Command{
	Use: "junit",
}

func init() {
	Junit.AddCommand(&cobra.Command{
		Use:     "markdown",
		Aliases: []string{"md"},
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			results, err := junit.ParseJunitResultFiles(args...)
			if err != nil {
				return err
			}
			fmt.Println(results.GenerateMarkdown())

			return nil
		},
	})
	Junit.AddCommand(&cobra.Command{
		Use:     "passfail",
		Aliases: []string{"pf"},
		Short: "Print result summary for JUnit test and set exit code to 1 if there are failed tests.",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			results, err := junit.ParseJunitResultFiles(args...)
			if err != nil {
				return err
			}
			fmt.Println(results.String())
			if results.Failed > 0 {
				os.Exit(1)
			}

			return nil
		},
	})
}
