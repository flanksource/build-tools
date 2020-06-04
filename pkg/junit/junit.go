package junit

import (
	"fmt"
	"io/ioutil"

	"github.com/joshdk/go-junit"
)

type TestResults struct {
	Suites  []junit.Suite
	Total   int
	Failed  int
	Skipped int
	Passed  int
}

func ParseJunitResults(files ...string) (*TestResults, error) {
	results := TestResults{}

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		suites, err := junit.Ingest(data)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %v", file, err)
		}
		results.Suites = append(results.Suites, suites...)
	}

	for _, suite := range results.Suites {
		for _, test := range suite.Tests {
			results.Total++
			switch test.Status {
			case "skipped":
				results.Skipped++
			case "failed", "error":
				results.Failed++
			case "passed":
				results.Passed++
			}
		}
	}

	return &results, nil
}

func (results TestResults) Success() bool {
	return results.Passed > 0 && results.Failed == 0
}

func (results TestResults) String() string {
	return fmt.Sprintf("%d tests, %d passed, %d failed, %d skipped", results.Total, results.Passed, results.Failed, results.Skipped)
}
