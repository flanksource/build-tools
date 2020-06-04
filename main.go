/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package main

import (
	"fmt"
	"os"

	"github.com/flanksource/build-tools/cmd"
	"github.com/spf13/cobra"
)

func main() {

	root := &cobra.Command{
		Use:   "build-tools",
		Short: "build-tools : A swiss-army knife of CI/CI related commands",
	}

	root.AddCommand(cmd.Github, cmd.Junit)

	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
