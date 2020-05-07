package main

import (
	"github.com/philipstaffordwood/build-tools/cmd"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
	command := exec.Command("go", "run", "main.go")
	out, err := command.CombinedOutput()
	assert.NoError(t,err,"main CLI entrypoint with no args should run without failure.")
	sout := string(out) // because out is []byte
	assert.True(t,strings.HasPrefix(sout, cmd.GetRootCommand().Short),"main CLI entrypoint with no args output should start with usage usage.")
}

func Test_main_failure(t *testing.T) {
	command := exec.Command("go", "run", "main.go","this-command-does-not-exist")
	out, err := command.CombinedOutput()
	assert.Error(t,err,"main CLI entrypoint with a non-existent command should fail.")
	sout := string(out) // because out is []byte
	t.Logf("%v",sout)

}

