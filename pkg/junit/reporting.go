/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package junit

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flanksource/build-tools/pkg/tesults"
	"github.com/google/go-github/v32/github"
	"github.com/joshdk/go-junit"
)

var (
	AnnotationWarning = "warning"
	AnnotationNotice  = "notice"
	AnnotationFailure = "failure"
)

const mdTableHeader = `| Result | Class | Message |
|--------|-------|--------|
`

func (results TestResults) GenerateMarkdown() string {
	mdResult := ""
	mdResult += fmt.Sprintf("<details><summary>%d test suites - Totals:  %d tests, %d failed, %d skipped, %d passed</summary>\n\n", len(results.Suites), results.Total, results.Failed, results.Skipped, results.Passed)

	for _, suite := range results.Suites {
		//44 tests, 21 passed, 0 failed, 23 skipped
		mdResult += fmt.Sprintf("<details><summary>%s:  %d tests, %d failed, %d skipped, %d passed</summary>\n\n", suite.Name, suite.Totals.Tests, suite.Totals.Failed, suite.Totals.Skipped, suite.Totals.Passed)
		mdResult += mdTableHeader
		for _, test := range suite.Tests {
			switch test.Status {
			case junit.StatusFailed:
				mdResult += fmt.Sprintf("| :x: | **%s** | `%s` |\n", test.Classname, test.Name)
				// no default:
				// any other status will be ignored
			}
		}
		for _, test := range suite.Tests {
			switch test.Status {
			case junit.StatusSkipped:
				mdResult += fmt.Sprintf("| :white_circle: | **%s** | `%s` |\n", test.Classname, test.Name)
				// no default:
				// any other status will be ignored
			}
		}
		mdResult += "\n</details>\n"
	}
	mdResult += "\n</details>\n"
	return mdResult
}

func toPtr(s string) *string {
	return &s
}

const MAX_ANNOTATIONS = 50
const GH_DEBUG = "::debug %s::%s\n"
const GH_ERROR = "::error %s::%s\n"

const GH_WARNING = "::warning %s %s::%s\n"
const GH_SKIP = "::debug %s::SKIPPED%s\n"

func TestName(t junit.Test) string {
	name := strings.Split(t.Name, ":")[0]
	if t.Classname != name {
		return t.Classname + "." + name
	}
	return name
}

func GithubWorkflowName(t junit.Test, suite junit.Suite) string {
	parts := strings.Split(t.Name, ":")
	name := TestName(t)
	if name != suite.Name {
		name = suite.Name + "." + name
	}
	if len(parts) == 3 {
		return fmt.Sprintf("file=%s,line=%d,col=%d", name, parts[1], parts[2])
	}
	return name
}

func (results TestResults) GenerateGithubWorkflowCommands() string {
	result := ""
	result += fmt.Sprintf(GH_WARNING, "", "", fmt.Sprintf("%d suites %d tests, %d failed, %d skipped, %d passed", len(results.Suites), results.Total, results.Failed, results.Skipped, results.Passed))

	for _, suite := range results.Suites {
		for _, test := range suite.Tests {

			msg := ""
			if test.Error != nil {
				msg = strings.Split(strings.TrimSpace(test.Error.Error()), "\n")[0]
			}
			switch test.Status {
			case junit.StatusFailed:
				result += fmt.Sprintf(GH_ERROR, GithubWorkflowName(test, suite), msg)
			case junit.StatusSkipped:
				result += fmt.Sprintf(GH_SKIP, GithubWorkflowName(test, suite), msg)
			}
		}
	}
	return result
}

func (results TestResults) GetGithubAnnotations() []*github.CheckRunAnnotation {
	list := []*github.CheckRunAnnotation{}
	count := 0
	for _, suite := range results.Suites {
		count++
		if count > MAX_ANNOTATIONS {
			return list
		}
		for _, test := range suite.Tests {
			var msg string
			if test.Error != nil {
				msg = test.Error.Error()
			}
			if msg == "" {
				msg = test.Name
			}
			title := test.Classname
			path := test.Name
			zero := 0
			annotation := &github.CheckRunAnnotation{
				Title:     &title,
				StartLine: &zero,
				EndLine:   &zero,
				Path:      &path,
				Message:   &msg,
			}

			switch test.Status {
			case junit.StatusFailed:
				annotation.AnnotationLevel = &AnnotationFailure
			default:
				continue
			}

			list = append(list, annotation)
		}
	}

	return list
}

var RESULT_MAP = map[string]string{
	"passed":       "pass",
	"skipped-fail": "fail",
	"skipped-pass": "pass",
	"failed":       "fail",
	"error":        "fail",
}

func (results TestResults) UploadToTesults(token string, failOnSkip bool) error {
	testResults := make([]interface{}, 0)
	for _, suite := range results.Suites {
		for _, test := range suite.Tests {
			resultString := string(test.Status)
			if test.Status == junit.StatusSkipped {
				if failOnSkip {
					resultString += "-fail"
				} else {
					resultString += "-pass"
				}
			}
			result := map[string]interface{}{
				"name":        test.Classname,
				"suite":       suite.Name,
				"description": test.Name,
				"result":      RESULT_MAP[resultString],
				"reason":      "",
				"duration":    test.Duration.Milliseconds(),
				"_stdout":     test.SystemOut,
				"_stderr":     test.SystemErr,
			}
			if test.Error != nil {
				result["reason"] = test.Error.Error()
			}
			if failOnSkip && test.Status == junit.StatusSkipped {
				result["reason"] = test.Name
			}
			testResults = append(testResults, result)
		}
	}
	data := map[string]interface{}{
		"target": token,
		"results": map[string]interface{}{
			"cases": testResults,
		},
	}
	result := tesults.Results(data)
	if !result["success"].(bool) {
		return errors.New(result["message"].(string))
	}
	return nil
}
