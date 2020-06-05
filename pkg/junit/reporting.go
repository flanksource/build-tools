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
	mdTable := ""
	for _, suite := range results.Suites {
		for _, test := range suite.Tests {
			switch test.Status {
			case junit.StatusFailed:
				mdTable += fmt.Sprintf("| :x: | **%s** | `%s` |\n", test.Classname, test.Name)
			case junit.StatusSkipped:
				mdTable += fmt.Sprintf("| :white_circle: | **%s** | `%s` |\n", test.Classname, test.Name)
			case junit.StatusPassed:
				// we ignore successes - we comment only on failed and skipped results to cut down report size
				break
				// no default:
				// any other status will be ignored
			}
		}
	}

	return mdTable
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
			case junit.StatusSkipped:
				annotation.AnnotationLevel = &AnnotationWarning
			default:
				annotation.AnnotationLevel = &AnnotationNotice
			}

			list = append(list, annotation)
		}
	}

	return list
}
