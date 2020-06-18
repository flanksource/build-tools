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

func ParseJunitResultFiles(files ...string) (*TestResults, error) {
	results := make([][]byte, len(files))

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		results = append(results, data)
	}
	return ParseJunitResults(results...)
}

func ParseJunitResultStrings(resultStrings ...string) (*TestResults, error) {
	results := make([][]byte, len(resultStrings))
	for _, result := range resultStrings {
		data := []byte(result)
		results = append(results, data)
	}
	return ParseJunitResults(results...)
}

func ParseJunitResults(results ...[]byte) (*TestResults, error) {
	tr := TestResults{}

	for _, result := range results {
		suites, err := junit.Ingest(result)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %v", result, err)
		}
		tr.Suites = append(tr.Suites, suites...)
	}

	for _, suite := range tr.Suites {
		for _, test := range suite.Tests {
			tr.Total++
			switch test.Status {
			case "skipped":
				tr.Skipped++
			case "failed", "error":
				tr.Failed++
			case "passed":
				tr.Passed++
			}
		}
	}

	return &tr, nil
}

func (results TestResults) Success() bool {
	return results.Passed > 0 && results.Failed == 0
}

func (results TestResults) String() string {
	return fmt.Sprintf("%d tests, %d passed, %d failed, %d skipped", results.Total, results.Passed, results.Failed, results.Skipped)
}
