package cmd

import (
	"github.com/flanksource/build-tools/cmd/test"
	"github.com/flanksource/build-tools/pkg/gh/pr"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestGetReportJunitCmd(t *testing.T) {
	cmd := GetReportJUnitCommand()
	assert.NotNil(t,cmd, "We must have a report-junit command")
}

func TestReportJunit_HasTokenField(t *testing.T) {
	cmd := GetReportJUnitCommand()
	token := "SOME_TOKEN_HERE"
	args := []string{"--auth-token",token}
	test.ParsesStringFlag(t,cmd,"auth-token",token, args,"We need to be able to parse an auth token to access Github")
}

func TestReportJunit_HasReportJunitSubcommand(t *testing.T) {
	cmd := GetReportJUnitCommand()
	test.HasSubcommand(t,cmd,"report-junit","We want a report-junit subcommand")
}

func Test_parseReportJunitFlagsAndArguments(t *testing.T) {

	tests := []struct {
		name string
		args []string
		shouldFail bool
		wantedPR *pr.PR
		wantedFiles []string
	}{
		{
			name: "happy path",
			args: []string {"flanksource/platform-cli", "1","junit.xml", "--auth-token", "SOME_TOKEN"},
			shouldFail: false,
			wantedPR: &pr.PR{
				APIToken: "SOME_TOKEN",
				Owner:    "flanksource",
				Repo:     "platform-cli",
				Num:      1,
			},
			wantedFiles: []string{"junit.xml"},
		},
		{
			name: "multiple files",
			args: []string {"flanksource/platform-cli", "1","junit1.xml", "junit2.xml", "--auth-token", "SOME_TOKEN"},
			shouldFail: false,
			wantedPR: &pr.PR{
				APIToken: "SOME_TOKEN",
				Owner:    "flanksource",
				Repo:     "platform-cli",
				Num:      1,
			},
			wantedFiles: []string{"junit1.xml", "junit2.xml" },
		},
		{
			name: "Invalid PR number fails",
			args: []string {"flanksource/platform-cli", "NaN","junit1.xml", "--auth-token", "SOME_TOKEN"},
			shouldFail: true,
			wantedPR: nil,
			wantedFiles: []string{},
		},
		{
			name: "Invalid owner/repo fails",
			args: []string {"flanksource","platform-cli", "1","junit1.xml", "--auth-token", "SOME_TOKEN"},
			shouldFail: true,
			wantedPR: nil,
			wantedFiles: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd:= GetReportJUnitCommand()
			cmd.ParseFlags(tt.args)
			err := validateReportJunitArguments(cmd, cmd.Flags().Args())

			if !tt.shouldFail {
				assert.NoError(t,err,"Validation failed with error: %v", err)
			} else if tt.shouldFail {
				t.Logf("Testcase %v should fail, but and did with error %v", tt.name, err)
				assert.Error(t,err,"Validation failed with error: %v", err)
				return // not testing further
			}
			gotPR, gotFiles, err := parseReportJunitFlagsAndArguments(cmd)
			assert.NoError(t,err,"Parsing failed with error: %v", err)
			assert.ElementsMatch(t,tt.wantedFiles,gotFiles,"We wanted different files.")
			t.Logf("%v -> %v, %v",tt.args, gotPR, gotFiles)
			assert.Equal(t,*(tt.wantedPR),gotPR,"We needed a PR to post a comment to.")
			//assert.Equal(t,tt.wantedPR,got,"We needed a PR to post a comment to.")
		})
	}
}