/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package junit

import (
	"fmt"

	"github.com/google/go-github/v31/github"
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
	for _, suite := range results.Suites {
		mdResult += "<details><summary>"+suite.Name+"</summary>\n"
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

	return mdResult
}

func toPtr(s string) *string {
	return &s
}

const MAX_ANNOTATIONS = 50

func (results TestResults) GetGithubAnnotations() []*github.CheckRunAnnotation {
	list := []*github.CheckRunAnnotation{}
	count := 0
	for _, suite := range results.Suites {
		count++
		if count > MAX_ANNOTATIONS {
			return list
		}
		for _, test := range suite.Tests {
			msg := fmt.Sprintf("stdout:%s stderr:%s", test.SystemOut, test.SystemErr)
			zero := 0
			annotation := &github.CheckRunAnnotation{
				Title:     &test.Classname,
				StartLine: &zero,
				EndLine:   &zero,
				Path:      &test.Name,
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
