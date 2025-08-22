package output

import (
	"github.com/andreygrechin/gcphelper/pkg/folders"
	"github.com/andreygrechin/gcphelper/pkg/organizations"
)

// FoldersToResources converts a slice of folders to a slice of resources.
func FoldersToResources(folderList []*folders.Folder) []Resource {
	resources := make([]Resource, len(folderList))
	for i, folder := range folderList {
		resources[i] = folder
	}

	return resources
}

// OrganizationsToResources converts a slice of organizations to a slice of resources.
func OrganizationsToResources(organizationList []*organizations.Organization) []Resource {
	resources := make([]Resource, len(organizationList))
	for i, org := range organizationList {
		resources[i] = org
	}

	return resources
}

// FolderHeaders returns the table headers for folder output.
func FolderHeaders() []string {
	return []string{"ID", "Display Name", "Parent", "State", "Create Time", "Update Time"}
}

// OrganizationHeaders returns the table headers for organization output.
func OrganizationHeaders() []string {
	return []string{"ID", "Display Name", "State", "Create Time", "Update Time"}
}
