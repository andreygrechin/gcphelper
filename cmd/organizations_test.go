package cmd_test

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"

	"github.com/andreygrechin/gcphelper/cmd"
	"github.com/andreygrechin/gcphelper/internal/logger"
	"github.com/andreygrechin/gcphelper/pkg/organizations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutputOrganizationsInvalidFormat(t *testing.T) {
	err := cmd.OutputOrganizations(nil, "invalid", false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported output format")
}

func TestOrganizationsCommandFlags(t *testing.T) {
	// create a fresh command instance for testing
	testLogger, _ := logger.NewDevelopmentLogger()
	organizationsCmd := cmd.NewOrganizationsCommand(testLogger)

	// test that global flags don't exist on the command itself (they're on the root)
	assert.Nil(t, organizationsCmd.Flags().Lookup("format"), "format flag should not exist on organizations command")
	assert.Nil(t, organizationsCmd.Flags().Lookup("verbose"), "verbose flag should not exist on organizations command")

	// test command aliases
	assert.Contains(t, organizationsCmd.Aliases, "org", "command should have 'org' alias")
}

func TestOutputOrganizationsIDFormat(t *testing.T) {
	testCases := map[string]struct {
		organizations []*organizations.Organization
		verbose       bool
		wantOut       string
		wantErr       string
	}{
		"empty list non-verbose": {
			organizations: []*organizations.Organization{},
			verbose:       false,
			wantOut:       "",
			wantErr:       "",
		},
		"empty list verbose": {
			organizations: []*organizations.Organization{},
			verbose:       true,
			wantOut:       "",
			wantErr:       "No organizations found.\n",
		},
		"single organization": {
			organizations: []*organizations.Organization{
				{
					ID:          "123456789",
					Name:        "organizations/123456789",
					DisplayName: "Test Organization",
					State:       "ACTIVE",
					CreateTime:  time.Now(),
					UpdateTime:  time.Now(),
				},
			},
			verbose: false,
			wantOut: "123456789\n",
			wantErr: "",
		},
		"multiple organizations": {
			organizations: []*organizations.Organization{
				{
					ID:          "123456789",
					Name:        "organizations/123456789",
					DisplayName: "Test Organization 1",
					State:       "ACTIVE",
					CreateTime:  time.Now(),
					UpdateTime:  time.Now(),
				},
				{
					ID:          "987654321",
					Name:        "organizations/987654321",
					DisplayName: "Test Organization 2",
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
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			os.Stdout = wOut
			os.Stderr = wErr

			// run the function
			err := cmd.OutputOrganizations(tc.organizations, "id", tc.verbose)

			// restore stdout/stderr and read output
			_ = wOut.Close()
			_ = wErr.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			outBytes, _ := io.ReadAll(rOut)
			errBytes, _ := io.ReadAll(rErr)

			// validate results
			require.NoError(t, err)
			assert.Equal(t, tc.wantOut, string(outBytes))
			assert.Equal(t, tc.wantErr, string(errBytes))
		})
	}
}

func TestOutputOrganizationsJSONFormat(t *testing.T) {
	orgList := []*organizations.Organization{
		{
			ID:          "123456789",
			Name:        "organizations/123456789",
			DisplayName: "Test Organization",
			State:       "ACTIVE",
			CreateTime:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:  time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	// capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// run the function
	err := cmd.OutputOrganizations(orgList, "json", false)

	// restore stdout and read output
	_ = w.Close()
	os.Stdout = oldStdout
	output, _ := io.ReadAll(r)

	// validate results
	require.NoError(t, err)
	assert.Contains(t, string(output), `"id": "123456789"`)
	assert.Contains(t, string(output), `"display_name": "Test Organization"`)
	assert.Contains(t, string(output), `"state": "ACTIVE"`)
}

func TestOutputOrganizationsTableFormat(t *testing.T) {
	testCases := map[string]struct {
		organizations []*organizations.Organization
		verbose       bool
		wantContains  []string
		wantStderr    string
	}{
		"empty list non-verbose": {
			organizations: []*organizations.Organization{},
			verbose:       false,
			wantContains:  []string{},
			wantStderr:    "",
		},
		"empty list verbose": {
			organizations: []*organizations.Organization{},
			verbose:       true,
			wantContains:  []string{},
			wantStderr:    "No organizations found.\n",
		},
		"single organization non-verbose": {
			organizations: []*organizations.Organization{
				{
					ID:          "123456789",
					Name:        "organizations/123456789",
					DisplayName: "Test Organization",
					State:       "ACTIVE",
					CreateTime:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdateTime:  time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			verbose: false,
			wantContains: []string{
				"123456789",
				"Test Organization",
				"ACTIVE",
				"2023-01-01 00:00:00",
				"2023-06-01 00:00:00",
			},
			wantStderr: "",
		},
		"single organization verbose": {
			organizations: []*organizations.Organization{
				{
					ID:          "123456789",
					Name:        "organizations/123456789",
					DisplayName: "Test Organization",
					State:       "ACTIVE",
					CreateTime:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdateTime:  time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			verbose: true,
			wantContains: []string{
				"123456789",
				"Test Organization",
				"ACTIVE",
				"Total organizations: 1",
			},
			wantStderr: "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// capture stdout and stderr
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			rOut, wOut, _ := os.Pipe()
			rErr, wErr, _ := os.Pipe()
			os.Stdout = wOut
			os.Stderr = wErr

			// run the function
			err := cmd.OutputOrganizations(tc.organizations, "table", tc.verbose)

			// restore stdout/stderr and read output
			_ = wOut.Close()
			_ = wErr.Close()
			os.Stdout = oldStdout
			os.Stderr = oldStderr

			outBytes, _ := io.ReadAll(rOut)
			errBytes, _ := io.ReadAll(rErr)

			// validate results
			require.NoError(t, err)
			output := string(outBytes)
			for _, want := range tc.wantContains {
				assert.Contains(t, output, want)
			}
			assert.Equal(t, tc.wantStderr, string(errBytes))
		})
	}
}

func TestOutputOrganizationsCSVFormat(t *testing.T) {
	orgList := []*organizations.Organization{
		{
			ID:          "123456789",
			Name:        "organizations/123456789",
			DisplayName: "Test Organization",
			State:       "ACTIVE",
			CreateTime:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdateTime:  time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	// capture stdout
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// run the function
	err := cmd.OutputOrganizations(orgList, "csv", false)

	// restore stdout and read output
	_ = w.Close()
	os.Stdout = oldStdout
	_, _ = io.Copy(&buf, r)
	output := buf.String()

	// validate results
	require.NoError(t, err)
	assert.Contains(t, output, "ID,Display Name,State,Create Time,Update Time")
	assert.Contains(t, output, "123456789,Test Organization,ACTIVE")
}
