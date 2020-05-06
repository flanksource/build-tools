/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package cmd

import (
	"errors"
	"github.com/flanksource/build-tools/pkg/gh/pr"
	"github.com/flanksource/build-tools/pkg/junit"
	"github.com/flanksource/build-tools/util"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

const AuthTokenFlag = "auth-token"
const SilentSuccessFlag = "silent-success"

// GetReportJUnitCommand returns the report-junit command, adds all child commands and sets flags appropriately.
func GetReportJUnitCommand() *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "report-junit",
		Short: "Posts a comment detailing JUnit test results to a PR",
		Long: ``,
		Args:  validateReportJunitArguments,
		RunE:  runReportJUnitCmd,

	}
	initReportJUnitCommand(cmd)
	return cmd
}

// initReportJUnitCommand defines the flags, persistent flags and configuration settings
// for the report-junit command and adds all sub commands.
func initReportJUnitCommand(cmd *cobra.Command) {

	cmd.Flags().StringP(AuthTokenFlag, "t", "", "The Github API key to be used to access Github.")
	cmd.Flags().Bool(SilentSuccessFlag,  false, "If set to 'true' posts no PR comment when JUnit test contains no failures or skips.")
}

func validateReportJunitArguments(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		return errors.New("four arguments needed.")
	}

	if len(strings.Split(args[0],"/")) < 2 {
		return errors.New("an owner and repo in the form owner/repo is required.")
	}

	if _,err := strconv.Atoi(args[1]); err != nil {
		return errors.New("valid PR number required.")
	}

	if token, err :=cmd.Flags().GetString(AuthTokenFlag); err!= nil || token == "" {
		return errors.New("a github API Token is required.")
	}
	return nil
}


func runReportJUnitCmd (cmd *cobra.Command, args []string) error {
	_pr, files, silentSuccess, err := parseReportJunitFlagsAndArguments(cmd)
	if err != nil {
		return err
	}
	rpts, err := util.GetFileString(files)
	md, err := junit.GenerateMarkdownReport(rpts, silentSuccess)
	if err != nil {
		return err
	}
	if md != "" {
		err := _pr.Post(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseReportJunitFlagsAndArguments(cmd *cobra.Command)  (_pr pr.PR, junitFiles []string, silentSuccess bool, err error ) {
	args := cmd.Flags().Args()
	if len(args)<3 {
		return pr.PR{}, []string{}, false, errors.New(`Not enough arguments: 
At least give:
owner/repo pr-number junit-file`)
	}

	ownerRepoArg := args[0]
	prNumarg,_ := strconv.Atoi(args[1])
	args = args[2:]

	ownerRepoSplit := strings.Split(ownerRepoArg,"/")


	token,_ := cmd.Flags().GetString(AuthTokenFlag)
	silentSuccess,_ = cmd.Flags().GetBool(SilentSuccessFlag)


	_pr = pr.PR{
		APIToken: token,
		Owner: ownerRepoSplit[0],
		Repo: ownerRepoSplit[1],
		Num: prNumarg,
	}

	return _pr, args, silentSuccess, nil
}





