/*
Copyright © 2020 Flanksource
This file is part of Flanksource build tools
*/
package util

import "github.com/flanksource/commons/files"

func GetFileString(filenames []string) (contents []string, err error) {
	contents = make([]string, 0, len(filenames))

	for _, file := range filenames {
		contents = append(contents, files.SafeRead(file))
	}

	return contents, nil

}
