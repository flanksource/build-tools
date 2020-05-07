package util

import (
	"github.com/flanksource/commons/files"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestGetFileString(t *testing.T) {
	tests := []struct {
		name         string
		files        []string
		wantContents []string
		wantErr      bool
	}{
		{
			name: "Single file",
			files: []string { "../fixtures/junit/junit-single-success.xml"},
			wantContents: []string {files.SafeRead("../fixtures/junit/junit-single-success.xml")},
			wantErr: false,
			},
		{
			name: "Multiple files",
			files: []string {
				"../fixtures/junit/junit-single-success.xml",
				"../fixtures/junit/junit-single-failure.xml",
				"../fixtures/junit/junit-multiple-each.xml",
				},
			wantContents: []string {
				files.SafeRead("../fixtures/junit/junit-single-success.xml"),
				files.SafeRead("../fixtures/junit/junit-single-failure.xml"),
				files.SafeRead("../fixtures/junit/junit-multiple-each.xml"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContents, err := GetFileString(tt.files)
			if tt.wantErr {
				assert.Error(t, err, "This test should fail.")
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantContents, gotContents)

		})
	}
}