/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package cmd

import (
	"github.com/flanksource/build-tools/pkg/junit"
	"github.com/spf13/cobra"
)

var successMessage, failureMessage string

var PullRequestReport = &cobra.Command{
	Use:   "report-junit",
	Short: "Posts a comment detailing JUnit test results to a PR",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		results, err := junit.ParseJunitResults(args...)
		if err != nil {
			return err
		}

		md := results.GenerateMarkdown()

		if failureMessage != "" && !results.Success() {
			md = failureMessage + "\n" + md
		}
		if successMessage != "" && results.Success() {
			md = successMessage + "\n" + md
		}
		md += results.GenerateMarkdown()

		if md != "" {
			return client.Comment(md)
		}
		return nil
	},
}

func init() {
	PullRequests.AddCommand(PullRequestReport)
	Actions.AddCommand(ActionReport)
	PullRequestReport.Flags().StringVar(&successMessage, "success-message", "", "This message will be added to the top of the PR comment if no failed or skipped tests are found.")
	PullRequestReport.Flags().StringVar(&failureMessage, "failure-message", "", "This message will be added to the top of the PR comment if failed or skipped tests are found.")
}

var ActionReport = &cobra.Command{
	Use:   "report-junit",
	Short: "Update a check result with annotations",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		results, err := junit.ParseJunitResults(args...)
		if err != nil {
			return err
		}
		run, err := client.GetActionRun()
		if err != nil {
			return err
		}
		return run.Annotate(*results)
	},
}
