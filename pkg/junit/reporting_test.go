package junit

import (
	"github.com/flanksource/commons/files"
	"github.com/stretchr/testify/assert"
	"testing"
)

const successMessage string = ":thumbsup: All good - no test failures."

const tableHeader =
`| Class| Message | Result |
|------|---------|--------|
`

const singleLoneFailureMessageRow string =
	"| **test.harbor** | `[harbor] Expected pods but none running - did you deploy?` | :x: |\n"


const singleLoneSkippedMessageRow string =
"| **test.thanos** | `Must specify a thanos server under e2e.server in client mode` | :white_circle: |\n"

func TestGenerateMarkdown(t *testing.T) {
	tests := []struct {
		name string
		silentSuccess bool
		simpleReport string
		wantMarkdown string
	}{
		{
			name: "Single Success",
			silentSuccess: false,
			simpleReport: files.SafeRead("../../fixtures/junit/junit-single-success.xml"),
			wantMarkdown: successMessage,
		},
		{
			name: "Single Success - Silent",
			silentSuccess: true,
			simpleReport: files.SafeRead("../../fixtures/junit/junit-single-success.xml"),
			wantMarkdown: "",
		},
		{
			name: "Single Lone Failure",
			silentSuccess: false,
			simpleReport: files.SafeRead("../../fixtures/junit/junit-single-failure.xml"),
			wantMarkdown: tableHeader + singleLoneFailureMessageRow,
		},
		{
			name: "Single Lone Skipped",
			silentSuccess: false,
			simpleReport: files.SafeRead("../../fixtures/junit/junit-single-skipped.xml"),
			wantMarkdown: tableHeader + singleLoneSkippedMessageRow,
		},
		{
			name: "Simple Multiple - one of each",
			silentSuccess: false,
			simpleReport: files.SafeRead("../../fixtures/junit/junit-multiple-each.xml"),
			wantMarkdown: tableHeader + singleLoneSkippedMessageRow + singleLoneFailureMessageRow,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			failures, gotMd, err := GenerateMarkdown(test.simpleReport, test.silentSuccess)
			assert.NoError(t, err, "We expect no errors converting to a comment.")
			t.Logf("%v - \n%v", failures, gotMd)
			assert.Equal(t,test.wantMarkdown, gotMd, "We wanted another markdown result.")


		})
	}
}