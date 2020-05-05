package pr

import (
	"fmt"
	"github.com/google/go-github/v31/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type PRIssueComment struct {
	PR
	Msg string
}

func (c *PRIssueComment) Post() error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "258abf1195a9cc5d3e9d0e700a2a02174e6f8002"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	prComment := github.IssueComment{
		Body: &c.Msg,
	}

	if _, _, err := client.Issues.CreateComment(ctx, c.Owner, c.Repo, c.Num, &prComment); err != nil {
		return fmt.Errorf("Failed to post comment to PR with error %v", err)
	}
	return nil
}

func (c *PRIssueComment) PostJUnit(report string ) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "258abf1195a9cc5d3e9d0e700a2a02174e6f8002"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	prComment := github.IssueComment{
		Body: &c.Msg,
	}

	if _, _, err := client.Issues.CreateComment(ctx, c.Owner, c.Repo, c.Num, &prComment); err != nil {
		return fmt.Errorf("Failed to post comment to PR with error %v", err)
	}
	return nil
}