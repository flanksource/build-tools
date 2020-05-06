package pr

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)


const postURLTemplate = "https://api.github.com/repos/%s/%s/issues/%d/comments"
func expectedPostURL(owner string, repo string, prNum int) string {
	return fmt.Sprintf(postURLTemplate, owner, repo, prNum)
}


func TestPR_Post(t *testing.T) {
	tests := []struct {
		Name string
		Message string
		_PR  PR
		ShouldFail bool
		TestRequest RoundTripFunc

	}{{
		Name: "Simple Post",
		Message: "Hello",
		_PR: PR{
			APIToken: "SOME_TOKEN",
			Owner:    "flanksource",
			Repo:     "platform-cli",
			Num:      1,
		},
		ShouldFail: false,
		TestRequest: func(req *http.Request) *http.Response {
			// Test request parameters
			assert.Containsf(t, req.Header,"Authorization","We need to be using the API token, but the Auth header is missing.")
			assert.Equalf(t, "Bearer SOME_TOKEN", req.Header.Get("Authorization"), "We need to be using the API token, but the Auth header doesn't have the API key we set.")
			assert.Equalf(t, expectedPostURL("flanksource","platform-cli",1), req.URL.String(), "The owner, repo and PR ID needs to be mapped to the correct URL path locations.")
			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		},
	},
		{
			Name: "Test Error Handling - 401",
			Message: "Hello",
			_PR: PR{
				APIToken: "SOME_TOKEN",
				Owner:    "flanksource",
				Repo:     "platform-cli",
				Num:      1,
			},
			ShouldFail: true,
			TestRequest: func(req *http.Request) *http.Response {
				// Test request parameters
				assert.Containsf(t, req.Header,"Authorization","We need to be using the API token, but the Auth header is missing.")
				assert.Equalf(t, "Bearer SOME_TOKEN", req.Header.Get("Authorization"), "We need to be using the API token, but the Auth header doesn't have the API key we set.")
				assert.Equalf(t, expectedPostURL("flanksource","platform-cli",1), req.URL.String(), "The owner, repo and PR ID needs to be mapped to the correct URL path locations.")
				return &http.Response{
					StatusCode: 401,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{"message":"Bad credentials","documentation_url":"https://developer.github.com/v3"}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test._PR.setTestClient(test.TestRequest)
			err := test._PR.Post("Test")
			if !test.ShouldFail {
				assert.NoErrorf(t, err, "Normal calls should not return errors.")
			} else {
				assert.Errorf(t, err, "This test should have failed, but didn't")
			}
		})
	}
}

