/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package cmd

import (
	"os"
	"strconv"
	"strings"

	"github.com/flanksource/build-tools/pkg/gh"
	"github.com/spf13/cobra"
)

var Github = &cobra.Command{
	Use:     "gh",
	Aliases: []string{"github"},
	Short:   "github related actions",
}

var Actions = &cobra.Command{
	Use: "actions",
}

var PullRequests = &cobra.Command{
	Use:     "pull-requests",
	Aliases: []string{"pr"},
}

var client = gh.Client{}

func init() {
	Github.AddCommand(Actions, PullRequests)
	if os.Getenv("GITHUB_REPOSITORY") != "" {
		Github.PersistentFlags().StringVar(&client.Repo, "repo", strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[1], "")
		Github.PersistentFlags().StringVar(&client.Owner, "owner", strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[0], "")
	} else {
		Github.PersistentFlags().StringVar(&client.Repo, "repo", os.Getenv("GITHUB_REPO"), "")
		Github.PersistentFlags().StringVar(&client.Owner, "owner", os.Getenv("GITHUB_OWNER"), "")
	}
	if os.Getenv("GITHUB_REF") != "" && strings.Contains(os.Getenv("GITHUB_REF"), "refs/pull/") {
		pr, _ := strconv.Atoi(strings.Split(os.Getenv("GITHUB_REF"), "/")[2])
		Github.PersistentFlags().IntVar(&client.PR, "pr", pr, "")
	} else {
		Github.PersistentFlags().IntVar(&client.PR, "pr", 0, "")
	}

	Github.PersistentFlags().StringVar(&client.SHA, "sha", os.Getenv("GITHUB_SHA"), "")
	// if os.Getenv("GITHUB_RUN_ID") != "" {
	// 	runId, _ := strconv.Atoi(os.Getenv("GITHUB_RUN_ID"))
	// 	Github.PersistentFlags().Int64Var(&client.RunID, "run-id", int64(runId), "")
	// } else {

	// }

	Github.PersistentFlags().Int64Var(&client.RunID, "run-id", 0, "")
	Github.PersistentFlags().StringVar(&client.Token, "token", os.Getenv("GITHUB_TOKEN"), "")
	Github.PersistentFlags().StringVar(&client.EventType, "event-type", os.Getenv("GITHUB_EVENT_NAME"), "")
	Github.PersistentFlags().StringVar(&client.EventPath, "event-path", os.Getenv("GITHUB_EVENT_PATH"), "")
	Github.PersistentFlags().StringVar(&client.Build, "build", "", "")
	Github.AddCommand(Actions, PullRequests)
}
