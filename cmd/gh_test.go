/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package cmd

import (
	"github.com/philipstaffordwood/build-tools/cmd/test"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestGetGhCmd(t *testing.T) {
	cmd := GetGhCommand()
	assert.NotNil(t,cmd, "We must have a gh command")
}

func TestGh_HasReportJunitSubcommand(t *testing.T) {
	cmd := GetGhCommand()
	test.HasSubcommand(t,cmd,"report-junit","We want a report-junit subcommand")
}


