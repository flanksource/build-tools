package cmd

import (
	"fmt"

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
}
