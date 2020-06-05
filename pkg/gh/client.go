package gh

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/flanksource/build-tools/pkg/junit"
	"github.com/google/go-github/v31/github"
	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2"
)

type Client struct {
	context.Context
	*github.Client
	*github.PullRequestEvent
	Repo, Owner          string
	Token                string
	EventPath, EventType string
	Workflow, Build      string
	PR                   int
	RunID                int64
	SHA, Ref             string
}

func (gh *Client) Init() error {
	gh.Context = context.TODO()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: gh.Token},
	)
	tc := oauth2.NewClient(gh.Context, ts)

	gh.Client = github.NewClient(tc)
	if gh.EventPath != "" {
		var data []byte
		var err error
		if data, err = ioutil.ReadFile(gh.EventPath); err != nil {
			return fmt.Errorf("error reading %s: %v", gh.EventPath, err)
		}
		if gh.EventType == "pull_request" {
			gh.PullRequestEvent = &github.PullRequestEvent{}
			if err := json.Unmarshal(data, gh.PullRequestEvent); err != nil {
				return fmt.Errorf("error unmarshal %s: %v", gh.EventPath, err)
			}
		}
	}
	return nil
}

type CheckRun struct {
	*Client
	*github.CheckRun
}

func (gh *Client) GetActionRun() (*CheckRun, error) {
	if err := gh.Init(); err != nil {
		return nil, err
	}

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
	zero := 0
	list := []*github.CheckRunAnnotation{&github.CheckRunAnnotation{
		AnnotationLevel: &junit.AnnotationNotice,
		StartLine:       &zero,
		EndLine:         &zero,
		Path:            &run.Build,
		Message:         &summary,
	}}
	list = append(list, results.GetGithubAnnotations()...)
	_, _, err := run.Checks.UpdateCheckRun(run, run.Owner, run.Repo, *run.ID, github.UpdateCheckRunOptions{
		Output: &github.CheckRunOutput{
			Title:       &title,
			Summary:     &summary,
			Annotations: list,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
