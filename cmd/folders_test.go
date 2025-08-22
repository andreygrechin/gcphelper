package cmd_test

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/andreygrechin/gcphelper/cmd"
	"github.com/andreygrechin/gcphelper/internal/logger"
	"github.com/andreygrechin/gcphelper/pkg/folders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestOutputFoldersInvalidFormat(t *testing.T) {
	err := cmd.OutputFolders(nil, "invalid", false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported output format")
}

func TestFoldersCommandFlags(t *testing.T) {
	// create a fresh command instance for testing
	testLogger, _ := logger.NewDevelopmentLogger()
	foldersCmd := cmd.NewFoldersCommand(testLogger)

	// test that the command has its specific flags
	assert.NotNil(t, foldersCmd.Flags().Lookup("parent-folder"), "parent-folder flag should exist")
	assert.NotNil(t, foldersCmd.Flags().Lookup("parent-organization"), "parent-organization flag should exist")

	// test that global flags don't exist on the command itself (they're on the root)
	assert.Nil(t, foldersCmd.Flags().Lookup("format"), "format flag should not exist on folders command")
	assert.Nil(t, foldersCmd.Flags().Lookup("verbose"), "verbose flag should not exist on folders command")

	// test default values for command-specific flags
	parentFolderFlag := foldersCmd.Flags().Lookup("parent-folder")
	assert.Empty(t, parentFolderFlag.DefValue, "parent-folder flag default should be empty")

	parentOrgFlag := foldersCmd.Flags().Lookup("parent-organization")
	assert.Empty(t, parentOrgFlag.DefValue, "parent-organization flag default should be empty")
}

func TestOutputFoldersIDFormat(t *testing.T) {
	testCases := map[string]struct {
		folders []*folders.Folder
		verbose bool
		wantOut string
		wantErr string
	}{
		"empty list non-verbose": {
			folders: []*folders.Folder{},
			verbose: false,
			wantOut: "",
			wantErr: "",
		},
		"empty list verbose": {
			folders: []*folders.Folder{},
			verbose: true,
			wantOut: "",
			wantErr: "No folders found.\n",
		},
		"single folder": {
			folders: []*folders.Folder{
				{
					ID:          "123456789",
					DisplayName: "Test Folder",
					Parent:      "organizations/987654321",
					State:       "ACTIVE",
					CreateTime:  time.Now(),
					UpdateTime:  time.Now(),
				},
			},
			verbose: false,
			wantOut: "123456789\n",
			wantErr: "",
		},
		"multiple folders": {
			folders: []*folders.Folder{
				{
					ID:          "123456789",
					DisplayName: "Test Folder 1",
					Parent:      "organizations/987654321",
					State:       "ACTIVE",
					CreateTime:  time.Now(),
					UpdateTime:  time.Now(),
				},
				{
					ID:          "987654321",
					DisplayName: "Test Folder 2",
					Parent:      "organizations/987654321",
					State:       "ACTIVE",
					CreateTime:  time.Now(),
					UpdateTime:  time.Now(),
				},
			},
			verbose: false,
			wantOut: "123456789\n987654321\n",
			wantErr: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// capture stdout and stderr
			originalStdout := os.Stdout
			originalStderr := os.Stderr
			defer func() {
				os.Stdout = originalStdout
				os.Stderr = originalStderr
			}()

			stdoutReader, stdoutWriter, _ := os.Pipe()
			stderrReader, stderrWriter, _ := os.Pipe()
			os.Stdout = stdoutWriter
			os.Stderr = stderrWriter

			// run the function
			err := cmd.OutputFolders(tc.folders, "id", tc.verbose)
			require.NoError(t, err)

			// close writers and read output
			_ = stdoutWriter.Close()
			_ = stderrWriter.Close()

			stdoutBuf := new(bytes.Buffer)
			stderrBuf := new(bytes.Buffer)
			_, _ = io.Copy(stdoutBuf, stdoutReader)
			_, _ = io.Copy(stderrBuf, stderrReader)

			// verify output
			assert.Equal(t, tc.wantOut, stdoutBuf.String(), "stdout should match expected")
			assert.Equal(t, tc.wantErr, stderrBuf.String(), "stderr should match expected")
		})
	}
}

var errTestNetwork = errors.New("network error")

func TestHandleFoldersError(t *testing.T) {
	testCases := map[string]struct {
		err     error
		parent  string
		wantMsg string
	}{
		"permission denied without parent": {
			err:     status.Error(codes.PermissionDenied, "insufficient permissions"),
			parent:  "",
			wantMsg: "permission denied: insufficient permissions to search folders",
		},
		"permission denied with parent": {
			err:     status.Error(codes.PermissionDenied, "insufficient permissions"),
			parent:  "organizations/123456789",
			wantMsg: "permission denied: insufficient permissions to list folders under parent organizations/123456789",
		},
		"other error": {
			err:     errTestNetwork,
			parent:  "",
			wantMsg: "network error",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			result := cmd.HandleFoldersError(tc.err, tc.parent)
			assert.Contains(t, result.Error(), tc.wantMsg)
		})
	}
}
