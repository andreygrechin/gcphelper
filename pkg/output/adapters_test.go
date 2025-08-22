package output_test

import (
	"testing"
	"time"

	"github.com/andreygrechin/gcphelper/pkg/folders"
	"github.com/andreygrechin/gcphelper/pkg/organizations"
	"github.com/andreygrechin/gcphelper/pkg/output"
	"github.com/stretchr/testify/assert"
)

func TestFoldersToResources(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	folderList := []*folders.Folder{
		{
			ID:          "123",
			DisplayName: "Test Folder 1",
			Parent:      "organizations/456",
			State:       "ACTIVE",
			CreateTime:  baseTime,
			UpdateTime:  baseTime.Add(time.Hour),
		},
		{
			ID:          "789",
			DisplayName: "Test Folder 2",
			Parent:      "folders/123",
			State:       "ACTIVE",
			CreateTime:  baseTime.Add(time.Minute),
			UpdateTime:  baseTime.Add(2 * time.Hour),
		},
	}

	resources := output.FoldersToResources(folderList)

	assert.Len(t, resources, 2)
	assert.Equal(t, "123", resources[0].GetID())
	assert.Equal(t, "Test Folder 1", resources[0].GetDisplayName())
	assert.Equal(t, "ACTIVE", resources[0].GetState())
	assert.Equal(t, baseTime, resources[0].GetCreateTime())
	assert.Equal(t, baseTime.Add(time.Hour), resources[0].GetUpdateTime())

	assert.Equal(t, "789", resources[1].GetID())
	assert.Equal(t, "Test Folder 2", resources[1].GetDisplayName())
}

func TestOrganizationsToResources(t *testing.T) {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	orgList := []*organizations.Organization{
		{
			ID:          "123456789",
			DisplayName: "Test Organization",
			State:       "ACTIVE",
			CreateTime:  baseTime,
			UpdateTime:  baseTime.Add(time.Hour),
		},
	}

	resources := output.OrganizationsToResources(orgList)

	assert.Len(t, resources, 1)
	assert.Equal(t, "123456789", resources[0].GetID())
	assert.Equal(t, "Test Organization", resources[0].GetDisplayName())
	assert.Equal(t, "ACTIVE", resources[0].GetState())
	assert.Equal(t, baseTime, resources[0].GetCreateTime())
	assert.Equal(t, baseTime.Add(time.Hour), resources[0].GetUpdateTime())
}

func TestFolderHeaders(t *testing.T) {
	headers := output.FolderHeaders()
	expected := []string{"ID", "Display Name", "Parent", "State", "Create Time", "Update Time"}
	assert.Equal(t, expected, headers)
}

func TestOrganizationHeaders(t *testing.T) {
	headers := output.OrganizationHeaders()
	expected := []string{"ID", "Display Name", "State", "Create Time", "Update Time"}
	assert.Equal(t, expected, headers)
}
