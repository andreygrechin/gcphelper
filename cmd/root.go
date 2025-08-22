package cmd

import (
	"fmt"

	"github.com/andreygrechin/gcphelper/internal/logger"
	"github.com/spf13/cobra"
)

type VersionInfo struct {
	Version   string
	BuildTime string
	Commit    string
}

// Global flags accessible to all subcommands.
var (
	globalFormat  string
	globalVerbose bool
)

// NewRootCommand creates and returns the root command.
func NewRootCommand(v VersionInfo, log logger.Logger) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gcphelper",
		Short: "A CLI tool to fetch information from Google Cloud",
		Long: `gcphelper is a command-line tool that helps you fetch and manage
information from Google Cloud Platform resources.

Features:
- List all projects in an organization
- List all folders in an organization
- List all accessible organizations

The tool uses Application Default Credentials for authentication.
Make sure you have authenticated with Google Cloud using:
  gcloud auth application-default login`,
	}

	rootCmd.AddCommand(NewFoldersCommand(log))
	rootCmd.AddCommand(NewOrganizationsCommand(log))

	rootCmd.Version = fmt.Sprintf("\n  Version: %s\n  Commit: %s\n  Built: %s", v.Version, v.Commit, v.BuildTime)

	// Add global persistent flags
	rootCmd.PersistentFlags().StringVarP(&globalFormat, "format", "f", "table", "Output format (table, json, csv, id)")
	rootCmd.PersistentFlags().BoolVarP(&globalVerbose, "verbose", "v", false,
		"Show additional output like counts and status messages")

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// Returns an exit code: 0 for success, 1 for error.
func Execute(v VersionInfo) int {
	log, err := logger.NewDevelopmentLogger()
	if err != nil {
		_ = fmt.Errorf("error creating logger: %w", err)

		return 1
	}
	defer func() {
		err := log.Close()
		if err != nil {
			_ = fmt.Errorf("error syncing logger: %w", err)
		}
	}()

	rootCmd := NewRootCommand(v, log)
	if err := rootCmd.Execute(); err != nil {
		_ = fmt.Errorf("error executing root command: %w", err)

		return 1
	}

	return 0
}
