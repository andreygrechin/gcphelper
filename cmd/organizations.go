package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/andreygrechin/gcphelper/internal/logger"
	"github.com/andreygrechin/gcphelper/pkg/organizations"
	"github.com/andreygrechin/gcphelper/pkg/output"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewOrganizationsCommand creates and returns the organizations command.
func NewOrganizationsCommand(log logger.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "organizations",
		Aliases: []string{"organization", "org"},
		Short:   "List Google Cloud organizations",
		Long: `List Google Cloud organizations accessible to the caller.

This command searches for organizations accessible to your credentials and displays
information about them. This requires the following IAM permissions:
- resourcemanager.organizations.get (to search organizations)

Examples:
  # List all accessible organizations
  gcphelper organizations

  # List organizations in JSON format
  gcphelper --format json organizations

  # List only organization IDs for scripting
  gcphelper --format id organizations

  # Pipe organization IDs to other commands
  gcphelper -f id organizations | xargs -I {} gcloud resource-manager organizations describe {}

  # List organizations with verbose output
  gcphelper --verbose organizations

  # Use the short alias
  gcphelper org`,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runOrganizationsCommand(globalFormat, globalVerbose, log)
		},
	}

	return cmd
}

func runOrganizationsCommand(format string, verbose bool, log logger.Logger) error {
	ctx := context.Background()

	// create organizations service
	service, err := organizations.NewServiceFromContextWithLogger(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create organizations service: %w", err)
	}
	defer func() {
		if closeErr := service.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close service: %v\n", closeErr)
		}
	}()

	// search for organizations
	organizationList, err := service.SearchOrganizations(ctx)
	if err != nil {
		return HandleOrganizationsError(err)
	}

	// output results
	return OutputOrganizations(organizationList, format, verbose)
}

func OutputOrganizations(organizationList []*organizations.Organization, format string, verbose bool) error {
	formatter := output.NewFormatterWithType(os.Stdout, verbose, "organizations")
	resources := output.OrganizationsToResources(organizationList)
	headers := output.OrganizationHeaders()

	if err := formatter.Format(resources, output.Format(format), headers); err != nil {
		return fmt.Errorf("failed to format organizations output: %w", err)
	}

	return nil
}

// HandleOrganizationsError provides enhanced error handling with helpful messages.
func HandleOrganizationsError(err error) error {
	// check if this is a permission denied error
	if st, ok := status.FromError(err); ok && st.Code() == codes.PermissionDenied {
		return fmt.Errorf(`permission denied: insufficient permissions to search organizations.

Ensure you have the 'resourcemanager.organizations.get' permission.

Original error: %w`, err)
	}

	// return the original error for other types of errors
	return err
}
