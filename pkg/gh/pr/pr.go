package pr

import (
	"fmt"
	"github.com/flanksource/build-tools/pkg/junit"
	"github.com/google/go-github/v31/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type PR struct {
	APIToken string
	Owner string
	Repo string
	Num int
}

func (p *PR) Post(msg string ) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.APIToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	prComment := github.IssueComment{
		Body: &msg,
	}

	if _, _, err := client.Issues.CreateComment(ctx, p.Owner, p.Repo, p.Num, &prComment); err != nil {
		return fmt.Errorf("Failed to post comment to PR with error %v", err)
	}
	return nil
}

func (p *PR) PostJUnitResults(junitFiles []string ) error {
	msg, err := junit.GenerateMarkdownReport(junitFiles)
	if err != nil {
		return fmt.Errorf("Failed to generate JUnit Markdown report with error %v", err)
	}
	p.Post(msg)
	return nil
}