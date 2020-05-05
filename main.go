/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package main

import (
	"fmt"
	"github.com/flanksource/build-tools/cmd"
	"os"
)

func main() {

	if err := cmd.GetRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
