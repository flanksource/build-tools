package cmd

import (
	"github.com/spf13/cobra"
)

var PullRequestComment = &cobra.Command{
	Use:   "comment",
	Short: "Posts a comment to a PR",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return client.Comment(args[0])
	},
}

func init() {
	PullRequests.AddCommand(PullRequestComment)

}
