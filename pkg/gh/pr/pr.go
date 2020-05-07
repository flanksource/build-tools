package pr

import (
	"fmt"
	"github.com/google/go-github/v31/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"net/http"
)



type PR struct {
	APIToken string
	Owner string
	Repo string
	Num int
	client *github.Client
}

func (p *PR) Post(msg string ) error {
	ctx := context.Background()
	_ = p.initClient(ctx)

	prComment := github.IssueComment{
		Body: &msg,
	}

	if _, _, err := p.client.Issues.CreateComment(ctx, p.Owner, p.Repo, p.Num, &prComment); err != nil {
		return fmt.Errorf("Failed to post comment to PR with error %v", err)
	}
	return nil
}


func (p *PR) initClient(ctx context.Context  ) error {
	if p.client != nil {
		return nil
	}

	ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: p.APIToken},
		)
	tc := oauth2.NewClient(ctx, ts)
		p.client = github.NewClient(tc)
	return nil
}




// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func (p *PR) setTestClient(fn RoundTripFunc)  {
	ctx := context.Background()


	tstClnt := &http.Client{
		Transport: RoundTripFunc(fn),
	}
	tstCtx := context.WithValue(ctx, oauth2.HTTPClient, tstClnt )
	_ = p.initClient(tstCtx)
	return
}


