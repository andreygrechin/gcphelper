package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/andreygrechin/gcphelper/internal/logger"
	"github.com/andreygrechin/gcphelper/pkg/folders"
	"github.com/andreygrechin/gcphelper/pkg/output"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrMutuallyExclusiveFlags is returned when both parent flags are specified.
var ErrMutuallyExclusiveFlags = errors.New("cannot specify both --parent-folder and --parent-organization")

// NewFoldersCommand creates and returns the folders command.
func NewFoldersCommand(log logger.Logger) *cobra.Command {
	var (
		parentFolder       string
		parentOrganization string
	)

	cmd := &cobra.Command{
		Use:     "folders",
		Aliases: []string{"folder"},
		Short:   "List Google Cloud folders",
		Long: `List Google Cloud folders using the SearchFolders API to discover all accessible folders.

This command uses the SearchFolders API which efficiently finds all folders you have
access to regardless of organizational hierarchy. This can discover folders even
when you don't have permissions on intermediate parent resources.

You can filter results by specifying a parent folder or organization.

Examples:
  # List all accessible folders
  gcphelper folders

  # List folders from a specific organization
  gcphelper folders --parent-organization 123456789

  # List folders under a specific parent folder
  gcphelper folders --parent-folder 987654321

  # List folders in JSON format
  gcphelper --format json folders

  # List only folder IDs for scripting
  gcphelper --format id folders

  # Pipe folder IDs to other commands
  gcphelper -f id folders | xargs -I {} gcloud resource-manager folders describe {}

  # List folders with verbose output
  gcphelper --verbose folders`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runFoldersCommand(parentFolder, parentOrganization, globalFormat, globalVerbose, log)
		},
	}

	cmd.Flags().StringVarP(&parentFolder, "parent-folder", "p", "", "Parent folder ID to filter folders by")
	cmd.Flags().StringVarP(&parentOrganization, "parent-organization", "o", "",
		"Parent organization ID to filter folders by")

	return cmd
}

func runFoldersCommand(parentFolder, parentOrganization, format string, verbose bool, log logger.Logger) error {
	ctx := context.Background()

	// create folders service
	service, err := folders.NewServiceFromContextWithLogger(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create folders service: %w", err)
	}
	defer func() {
		if closeErr := service.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close service: %v\n", closeErr)
		}
	}()

	// validate mutually exclusive flags
	if parentFolder != "" && parentOrganization != "" {
		return ErrMutuallyExclusiveFlags
	}

	// configure fetch options
	opts := folders.NewFetchOptions()
	if parentFolder != "" {
		opts.Parent = "folders/" + parentFolder
	} else if parentOrganization != "" {
		opts.Parent = "organizations/" + parentOrganization
	}

	// fetch folders using SearchFolders API
	folderList, err := service.ListFolders(ctx, opts)
	if err != nil {
		return HandleFoldersError(err, opts.Parent)
	}

	// output results
	return OutputFolders(folderList, format, verbose)
}

func OutputFolders(folderList []*folders.Folder, format string, verbose bool) error {
	formatter := output.NewFormatterWithType(os.Stdout, verbose, "folders")
	resources := output.FoldersToResources(folderList)
	headers := output.FolderHeaders()

	if err := formatter.Format(resources, output.Format(format), headers); err != nil {
		return fmt.Errorf("failed to format folders output: %w", err)
	}

	return nil
}

// HandleFoldersError provides enhanced error handling with helpful messages.
func HandleFoldersError(err error, parent string) error {
	// check if this is a permission denied error
	if st, ok := status.FromError(err); ok && st.Code() == codes.PermissionDenied {
		if parent == "" {
			return fmt.Errorf(`permission denied: insufficient permissions to search folders.

Ensure you have the required IAM permissions:
  - resourcemanager.folders.list (to access folders)
  - Or specify a parent with: gcphelper folders --parent-organization YOUR_ORG_ID

Original error: %w`, err)
		}

		return fmt.Errorf(`permission denied: insufficient permissions to list folders under parent %s.

Ensure you have the 'resourcemanager.folders.list' permission for this parent resource.

Original error: %w`, parent, err)
	}

	// return the original error for other types of errors
	return err
}
