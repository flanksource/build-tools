/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package cmd

import (
	"errors"
	"github.com/philipstaffordwood/build-tools/pkg/gh/pr"
	"github.com/philipstaffordwood/build-tools/pkg/junit"
	"github.com/philipstaffordwood/build-tools/util"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

const AuthTokenFlag = "auth-token"
const SilentSuccessFlag = "silent-success"
const SuccessMessageFlag = "success-message"
const FailureMessageFlag = "failure-message"

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
	cmd.Flags().String(SuccessMessageFlag,  "", "This message will be added to the top of the PR comment if no failed or skipped tests are found.")
	cmd.Flags().String(FailureMessageFlag,  "", "This message will be added to the top of the PR comment if failed or skipped tests are found.")
}

func validateReportJunitArguments(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		return errors.New("three arguments needed.")
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




func parseReportJunitFlagsAndArguments(cmd *cobra.Command)  (_pr pr.PR, junitFiles []string, silentSuccess bool, successMessage string, failureMessage string, err error ) {
	args := cmd.Flags().Args()
	if len(args)<3 {
		return pr.PR{}, []string{}, false, "","",errors.New(`Not enough arguments: 
At least give:
owner/repo pr-number junit-file`)
	}

	ownerRepoArg := args[0]
	prNumarg,_ := strconv.Atoi(args[1])
	args = args[2:]

	ownerRepoSplit := strings.Split(ownerRepoArg,"/")


	token,_ := cmd.Flags().GetString(AuthTokenFlag)
	silentSuccess,_ = cmd.Flags().GetBool(SilentSuccessFlag)
	successMessage,_ = cmd.Flags().GetString(SuccessMessageFlag)
	failureMessage,_ = cmd.Flags().GetString(FailureMessageFlag)

	_pr = pr.PR{
		APIToken: token,
		Owner: ownerRepoSplit[0],
		Repo: ownerRepoSplit[1],
		Num: prNumarg,
	}

	return _pr, args, silentSuccess, successMessage,failureMessage,nil
}

func runReportJUnitCmd (cmd *cobra.Command, args []string) error {
	_pr, files, silentSuccess, successMessage, failureMessage,err := parseReportJunitFlagsAndArguments(cmd)
	if err != nil {
		return err
	}
	rpts, err := util.GetFileString(files)
	md, hadFailures, err := junit.GenerateMarkdownReport(rpts, silentSuccess)
	if err != nil {
		return err
	}
	if successMessage != "" && !hadFailures {
		md = successMessage + "\n"+  md
	}
	if failureMessage != "" && hadFailures {
		md = failureMessage + "\n"+ md
	}
	if md != "" {
		err := _pr.Post(md)
		if err != nil {
			return err
		}
	}
	return nil
}





