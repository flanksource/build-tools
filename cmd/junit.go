package cmd

import (
	"errors"
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
		Short:   "Print result summary for JUnit test and set exit code to 1 if there are failed tests.",
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
	tesultsCommand := &cobra.Command{
		Use:     "upload-tesults",
		Aliases: []string{"ut"},
		Short:   "Upload test results for a JUnit test to Tesults. Requires a Tesults token passed in or as TESULTS_TOKEN environment variable.",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if tesultsToken, _ := cmd.Flags().GetString("token"); tesultsToken != "" {
				results, err := junit.ParseJunitResultFiles(args...)
				if err != nil {
					return err
				}
				return results.UploadToTesults(tesultsToken)
			}
			return errors.New("No Tesults token supplied")
		},
	}
	tesultsCommand.Flags().StringP("token", "t", os.Getenv("TESULTS_TOKEN"),
		"The tesults token to use for the upload. Defaults to the TESULTS_TOKEN environment variable.")
	Junit.AddCommand(tesultsCommand)
}
