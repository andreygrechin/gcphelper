package organizations

import (
	"strings"
	"time"

	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"github.com/jedib0t/go-pretty/v6/table"
)

const orgPrefix = "organizations/"

// Organization represents a Google Cloud organization.
type Organization struct {
	ID          string    `json:"id"`           // ID is the organization's numeric ID ("123456789")
	Name        string    `json:"name"`         // Name is the organization's resource name ("organizations/123456789")
	DisplayName string    `json:"display_name"` // DisplayName is the organization's human-readable name
	State       string    `json:"state"`        // State indicates the organization's lifecycle state
	CreateTime  time.Time `json:"create_time"`  // CreateTime is when the organization was created
	UpdateTime  time.Time `json:"update_time"`  // UpdateTime is when the organization was last updated
}

// OrganizationFromProto converts a protobuf organization to our internal type.
func OrganizationFromProto(pb *resourcemanagerpb.Organization) *Organization {
	if pb == nil {
		return nil
	}

	org := &Organization{
		ID:          strings.TrimPrefix(pb.GetName(), orgPrefix),
		Name:        pb.GetName(),
		DisplayName: pb.GetDisplayName(),
		State:       pb.GetState().String(),
	}

	if pb.GetCreateTime() != nil {
		org.CreateTime = pb.GetCreateTime().AsTime()
	}

	if pb.GetUpdateTime() != nil {
		org.UpdateTime = pb.GetUpdateTime().AsTime()
	}

	return org
}

// GetID returns the organization's ID.
func (o *Organization) GetID() string {
	return o.ID
}

// GetDisplayName returns the organization's display name.
func (o *Organization) GetDisplayName() string {
	return o.DisplayName
}

// GetState returns the organization's state.
func (o *Organization) GetState() string {
	return o.State
}

// GetCreateTime returns the organization's creation time.
func (o *Organization) GetCreateTime() time.Time {
	return o.CreateTime
}

// GetUpdateTime returns the organization's last update time.
func (o *Organization) GetUpdateTime() time.Time {
	return o.UpdateTime
}

// TableRow returns the organization data as a table row.
func (o *Organization) TableRow() []interface{} {
	return table.Row{
		o.ID,
		o.DisplayName,
		o.State,
		o.CreateTime.Format("2006-01-02 15:04:05"),
		o.UpdateTime.Format("2006-01-02 15:04:05"),
	}
}
