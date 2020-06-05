package gh

import (
	"context"
	"fmt"

	"github.com/flanksource/build-tools/pkg/junit"
	"github.com/google/go-github/v31/github"
	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2"
)

type Client struct {
	context.Context
	*github.Client
	*github.PullRequestEvent
	Repo, Owner     string
	Token           string
	Workflow, Build string
	PR              int
	RunID           int64
	SHA, Ref        string
}

func (gh *Client) Init() {
	gh.Context = context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh.Token},
	)
	tc := oauth2.NewClient(gh.Context, ts)

	gh.Client = github.NewClient(tc)
}

type CheckRun struct {
	*Client
	*github.CheckRun
}

func (gh *Client) GetActionRun() (*CheckRun, error) {
	gh.Init()

	if gh.RunID > 0 {
		run, _, err := gh.Checks.GetCheckRun(gh, gh.Owner, gh.Repo, gh.RunID)
		if err != nil {
			return nil, err
		}
		return &CheckRun{
			Client:   gh,
			CheckRun: run,
		}, nil
	}
	if gh.SHA == "" || gh.Build == "" {
		return nil, fmt.Errorf("must specify either --run-id or --sha and --build")
	}
	results, _, err := gh.Checks.ListCheckRunsForRef(gh, gh.Owner, gh.Repo, *gh.PullRequestEvent.PullRequest.Head.SHA, &github.ListCheckRunsOptions{})
	if err != nil {
		return nil, stacktrace.Propagate(err, "Error get action run %s/%s: %s", gh.Owner, gh.Repo, gh.SHA)
	}
	builds := []string{}
	for _, checkRun := range results.CheckRuns {
		fmt.Println(*checkRun.Name)
		builds = append(builds, *checkRun.Name)
		if *checkRun.Name == gh.Build {
			return &CheckRun{
				Client:   gh,
				CheckRun: checkRun,
			}, nil
		}
	}

	return nil, fmt.Errorf("check run '%s' not found for  %s, valid: %v", gh.Build, gh.SHA, builds)
}

func (gh *Client) Comment(comment string) error {
	gh.Init()
	if gh.PR == 0 {
		return fmt.Errorf("--pr not specified")
	}

	_, _, err := gh.Issues.CreateComment(gh, gh.Owner, gh.Repo, gh.PR, &github.IssueComment{
		Body: &comment,
	})
	return err
}

func (run *CheckRun) Annotate(results junit.TestResults) error {
	title := "Test Results"
	summary := results.String()

	// _, _, err := run.Checks.CreateCheckRun(run.Context, run.Owner, run.Repo, github.CreateCheckRunOptions{
	// 	Name:    run.Build + "3",
	// 	HeadSHA: run.SHA,
	// 	Output: &github.CheckRunOutput{
	// 		Title:       &title,
	// 		Summary:     &summary,
	// 		Annotations: results.GetGithubAnnotations(),
	// 	},
	// })
	_, _, err := run.Checks.UpdateCheckRun(run, run.Owner, run.Repo, *run.ID, github.UpdateCheckRunOptions{
		Output: &github.CheckRunOutput{
			Title:       &title,
			Summary:     &summary,
			Annotations: results.GetGithubAnnotations(),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
