package junit

import (
	"fmt"
	"github.com/flanksource/commons/files"
	"github.com/joshdk/go-junit"
)

const mdTableHeader = `| Class| Message | Result |
|------|---------|--------|
`

func GenerateMarkdownReport(junitFiles []string, silentSuccess bool) (string, error) {
	var hasFailures = false
	result := ""
	for _, file := range junitFiles {
		rpt := files.SafeRead(file)
		if rpt == "" {
			//log warning, but continue
		} else {
			failures, md, err := GenerateMarkdown(rpt,silentSuccess)
			if err != nil {
				//log error, but continue
			}
			if failures {
				hasFailures = true
			}
			result += md
		}
	}
	if !hasFailures {
		if !silentSuccess {
			return ":thumbsup: All good - no test failures.", nil

		} else {
			return "", nil
		}
	}
	return result, nil
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
				mdTable += fmt.Sprintf("| **%s** | `%s` | :x: |\n", test.Classname, test.Name)
			case junit.StatusSkipped:
				mdTable += fmt.Sprintf("| **%s** | `%s` | :white_circle: |\n", test.Classname, test.Name)
			case junit.StatusPassed:
				break
				default:
				return hasFailures, "", fmt.Errorf("Not implemented")
			}
		}
	}

	if !hasFailures {
		if !silentSuccess {
			return hasFailures,":thumbsup: All good - no test failures.", nil

		} else {
			return hasFailures, "", nil
		}
	}

	return hasFailures, mdTable, nil


}
