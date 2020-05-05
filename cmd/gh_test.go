package cmd

import (
	"github.com/flanksource/build-tools/cmd/test"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestGetGhCmd(t *testing.T) {
	cmd := GetGhCommand()
	assert.NotNil(t,cmd, "We must have a gh command")
}

//func TestGh_HasTokenField(t *testing.T) {
//	cmd := GetGhCommand()
//	token := "SOME_TOKEN_HERE"
//	args := []string{"--auth-token",token}
//	test.ParsesStringFlag(t,cmd,"auth-token",token, args,"We need to be able to parse an auth token for gh sub-commands")
//}

func TestGh_HasReportJunitSubcommand(t *testing.T) {
	cmd := GetGhCommand()
	test.HasSubcommand(t,cmd,"report-junit","We want a report-junit subcommand")
}


