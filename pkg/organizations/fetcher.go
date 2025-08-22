package organizations

import (
	"context"
	"errors"
	"fmt"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"google.golang.org/api/iterator"
)

// Fetcher defines the interface for fetching organizations from Google Cloud.
type Fetcher interface {
	// SearchOrganizations searches for organizations accessible to the caller.
	SearchOrganizations(ctx context.Context) ([]*Organization, error)

	// Close releases any resources held by the fetcher.
	Close() error
}

// Client implements the Fetcher interface using the Google Cloud Resource Manager API.
type Client struct {
	client *resourcemanager.OrganizationsClient
}

// NewClientFromContext creates a new organizations client using application default credentials.
func NewClientFromContext(ctx context.Context) (*Client, error) {
	c, err := resourcemanager.NewOrganizationsClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create organizations client: %w", err)
	}

	return &Client{
		client: c,
	}, nil
}

// SearchOrganizations searches for organizations accessible to the caller.
func (c *Client) SearchOrganizations(ctx context.Context) ([]*Organization, error) {
	req := &resourcemanagerpb.SearchOrganizationsRequest{}

	it := c.client.SearchOrganizations(ctx, req)

	var organizations []*Organization
	for {
		org, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate organizations: %w", err)
		}

		organizations = append(organizations, OrganizationFromProto(org))
	}

	return organizations, nil
}

// Close releases any resources held by the fetcher.
func (c *Client) Close() error {
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("failed to close organizations client: %w", err)
	}

	return nil
}
