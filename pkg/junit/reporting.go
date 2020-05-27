/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package junit

import (
	"fmt"
	log "github.com/flanksource/commons/logger"
	"github.com/joshdk/go-junit"
)

const mdTableHeader = `| Result | Class | Message |
|--------|-------|--------|
`

const SuccessMessage = ":thumbsup: All good - no test failures."

func GenerateMarkdownReport(reports []string, silentSuccess bool) (string, bool, error) {
	var hasFailures = false
	var hadError error = nil
	result := ""
	for _, rpt := range reports {
		if rpt == "" {
			//log warning, but continue
			log.Warnf("Empty report.")
		} else {
			failures, md, err := GenerateMarkdown(rpt, silentSuccess)
			if err != nil {
				//log error, but continue
				log.Errorf("Error generating report: %v", err)
				hadError = err
			}
			if failures {
				hasFailures = true
			}
			result += md
		}
	}
	if !hasFailures {
		if !silentSuccess {
			return SuccessMessage, hasFailures, hadError

		} else {
			return "", hasFailures, hadError
		}
	}
	return result, hasFailures, hadError
}

func GenerateMarkdown(reportXml string, silentSuccess bool) (hasFailures bool, md string, err error) {
	hasFailures = false

	xml := []byte(reportXml)

	suites, err := junit.Ingest(xml)
	if err != nil {
		return false, "", err
	}

	mdTable := mdTableHeader

	for _, suite := range suites {
		for _, test := range suite.Tests {
			if test.Status != junit.StatusPassed {
				hasFailures = true
			}
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

	if !hasFailures {
		if !silentSuccess {
			return hasFailures, ":thumbsup: All good - no test failures.", nil

		} else {
			return hasFailures, "", nil
		}
	}

	return hasFailures, mdTable, nil

}
