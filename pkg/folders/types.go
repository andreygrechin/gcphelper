package folders

import (
	"strings"
	"time"

	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"github.com/jedib0t/go-pretty/v6/table"
)

// Folder represents a Google Cloud folder resource.
type Folder struct {
	ID          string    `json:"id"`           // ID is the folder's unique identifier ("123456789")
	Name        string    `json:"name"`         // Name is the folder's resource name ("folders/123456789")
	DisplayName string    `json:"display_name"` // DisplayName is the folder's human-readable name
	Parent      string    `json:"parent"`       // Parent is the parent resource (organization or folder)
	State       string    `json:"state"`        // State indicates the folder's lifecycle state
	CreateTime  time.Time `json:"create_time"`  // CreateTime is when the folder was created
	UpdateTime  time.Time `json:"update_time"`  // UpdateTime is when the folder was last updated
}

// FetchOptions configures how folders are fetched.
type FetchOptions struct {
	Parent string // Parent specifies the parent resource to filter folders by (e.g., "folders/123", "organizations/456").
}

// NewFetchOptions creates a new FetchOptions with default values.
func NewFetchOptions() *FetchOptions {
	return &FetchOptions{}
}

const folderPrefix = "folders/"

// FolderFromProto converts a protobuf Folder to our Folder type.
func FolderFromProto(pb *resourcemanagerpb.Folder) *Folder {
	folder := &Folder{
		ID:          strings.TrimPrefix(pb.GetName(), folderPrefix),
		Name:        pb.GetName(),
		DisplayName: pb.GetDisplayName(),
		Parent:      pb.GetParent(),
		State:       pb.GetState().String(),
	}

	if pb.GetCreateTime() != nil {
		folder.CreateTime = pb.GetCreateTime().AsTime()
	}

	if pb.GetUpdateTime() != nil {
		folder.UpdateTime = pb.GetUpdateTime().AsTime()
	}

	return folder
}

// GetID returns the folder's ID.
func (f *Folder) GetID() string {
	return f.ID
}

// GetDisplayName returns the folder's display name.
func (f *Folder) GetDisplayName() string {
	return f.DisplayName
}

// GetState returns the folder's state.
func (f *Folder) GetState() string {
	return f.State
}

// GetCreateTime returns the folder's creation time.
func (f *Folder) GetCreateTime() time.Time {
	return f.CreateTime
}

// GetUpdateTime returns the folder's last update time.
func (f *Folder) GetUpdateTime() time.Time {
	return f.UpdateTime
}

// TableRow returns the folder data as a table row.
func (f *Folder) TableRow() []interface{} {
	return table.Row{
		f.ID,
		f.DisplayName,
		f.Parent,
		f.State,
		f.CreateTime.Format("2006-01-02 15:04:05"),
		f.UpdateTime.Format("2006-01-02 15:04:05"),
	}
}
