package junit

import (
	"fmt"
	"github.com/philipstaffordwood/build-tools/util"
	"github.com/flanksource/commons/files"
	log "github.com/sirupsen/logrus"
	loghooks "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"strings"
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

func TestGenerateMarkdownFailures(t *testing.T) {
	tests := []struct {
		name string
		silentSuccess bool
		simpleReport string
	}{
		{
			name: "Malformed XML fails",
			silentSuccess: false,
			simpleReport: files.SafeRead("../../fixtures/junit/malformed-xml.xml"),
		},

	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := GenerateMarkdown(test.simpleReport, test.silentSuccess)
			assert.Error(t, err, "We expect ts test to fail")


		})
	}
}

func ignoreErrorHelper(contents []string, err error) []string {
	return contents
}

func TestGenerateMarkdownReport(t *testing.T) {
	boolAssignHelper := func (_bool bool) *bool {
		return &_bool
	}
	type args struct {
		reports       []string
		silentSuccess bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantFoundFailures *bool //nil means we don't care about the result
		wantErr bool
	}{
		{
			name: "no reports gives blank output",
			args: args{
				reports:       []string{},
				silentSuccess: true,
			},
			want: "",
			wantFoundFailures: boolAssignHelper(false),
			wantErr: false,
		},
		{
			name: "only success reports with silent set gives blank output",
			args: args{
				reports:       ignoreErrorHelper(util.GetFileString([]string{"../../fixtures/junit/junit-single-success.xml"})),
				silentSuccess: true,
			},
			want: "",
			wantFoundFailures: boolAssignHelper(false),
			wantErr: false,
		},
		{
			name: "only success reports with silent not set gives single success message",
			args: args{
				reports:       ignoreErrorHelper(util.GetFileString([]string{"../../fixtures/junit/junit-single-success.xml"})),
				silentSuccess: false,
			},
			want: SuccessMessage,
			wantFoundFailures: boolAssignHelper(false),
			wantErr: false,
		},
		{
			name: "single report with single failure produces result",
			args: args{
				reports:       ignoreErrorHelper(util.GetFileString([]string{"../../fixtures/junit/junit-single-failure.xml"})),
				silentSuccess: false,
			},
			want: tableHeader+singleLoneFailureMessageRow,
			wantFoundFailures: boolAssignHelper(true),
			wantErr: false,
		},
		{
			name: "single report with single skip produces result",
			args: args{
				reports:       ignoreErrorHelper(util.GetFileString([]string{"../../fixtures/junit/junit-single-skipped.xml"})),
				silentSuccess: false,
			},
			want: tableHeader+singleLoneSkippedMessageRow,
			wantFoundFailures: boolAssignHelper(true),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, foundFailures, err := GenerateMarkdownReport(tt.args.reports, tt.args.silentSuccess)
			if tt.wantErr{
				assert.Error(t,err)
				return
			}
			assert.Equal(t, tt.want, got)
			if tt.wantFoundFailures != nil {
				assert.Equal(t, *(tt.wantFoundFailures), foundFailures )
			}

		})
	}
}

func TestGenerateMarkdownReportLogging(t *testing.T) {
	stringRefHelper := func (_string string) *string {
		return &_string
	}
	type args struct {
		reports       []string
		silentSuccess bool
	}
	tests := []struct {
		name    string
		args    args
		wantLog    string
		wantLevel log.Level
		wantErr bool
		wantRpt    *string
	}{
		{
			name: "an empty report results in a log warning",
			args: args{
				reports:       []string{""},
				silentSuccess: false,
			},
			wantLog: "Empty report.",
			wantLevel: log.WarnLevel,
			wantErr: false,
		},
		{
			name: "a report causing failure results in a log error, but not failure",
			args: args{
				reports:  ignoreErrorHelper(
					util.GetFileString([]string{
						"../../fixtures/junit/junit-single-failure.xml",
						"../../fixtures/junit/malformed-xml.xml",
					})),
				silentSuccess: false,
			},
			wantLog: "Error generating report:",
			wantLevel: log.ErrorLevel,
			wantErr: false,
			wantRpt: stringRefHelper(tableHeader+ singleLoneFailureMessageRow),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loghooks := loghooks.NewGlobal()
			gotRpt, _, err := GenerateMarkdownReport(tt.args.reports, tt.args.silentSuccess)
			if tt.wantErr{
				assert.Error(t,err)
				return
			}
			ContainsLogEntryWithMessagePrefix(t, loghooks.AllEntries(), tt.wantLog, tt.wantLevel)
			if tt.wantRpt != nil {
				assert.Equal(t, *(tt.wantRpt), gotRpt)
			}
		})
	}
}


func ContainsLogEntryWithMessagePrefix(t *testing.T, entries []*log.Entry, wantedMessage string, wantedLevel log.Level) bool {
	if entries == nil {
		return assert.Fail(t, fmt.Sprintf("Wanted log entry '%v', but got a nil collection of entries",wantedMessage))
	}
	for _,entry := range entries {
		if strings.HasPrefix(entry.Message, wantedMessage) {
			return assert.Equalf(t, wantedLevel, entry.Level,"Found wanted log entry '%v', but it had loglevel of '%v' and we wanted '%v' ",wantedMessage, entry.Level, wantedLevel)
		}
	}
	return assert.Fail(t, fmt.Sprintf("Wanted log entry '%v', but didn't find it ",wantedMessage))

}



