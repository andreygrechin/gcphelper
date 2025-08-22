package folders

import (
	"context"
	"errors"
	"fmt"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"google.golang.org/api/iterator"
)

// Fetcher defines the interface for fetching folders from Google Cloud.
type Fetcher interface {
	// ListFolders lists all accessible folders.
	ListFolders(ctx context.Context, opts *FetchOptions) ([]*Folder, error)

	// ListFoldersFromParent lists folders under a specific parent resource.
	ListFoldersFromParent(ctx context.Context, parent string, opts *FetchOptions) ([]*Folder, error)

	// Close releases any resources held by the fetcher.
	Close() error
}

// Client implements the Fetcher interface using the Google Cloud Resource Manager API.
type Client struct {
	foldersClient *resourcemanager.FoldersClient
}

// NewClientFromContext creates a new folders client using application default credentials.
func NewClientFromContext(ctx context.Context) (*Client, error) {
	c, err := resourcemanager.NewFoldersClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create folders client: %w", err)
	}

	return &Client{
		foldersClient: c,
	}, nil
}

// ListFolders lists all accessible folders, optionally filtered by parent.
func (c *Client) ListFolders(ctx context.Context, opts *FetchOptions) ([]*Folder, error) {
	if opts == nil {
		opts = NewFetchOptions()
	}

	// use SearchFolders API with parent filtering
	return c.searchAllAccessibleFolders(ctx, opts)
}

// ListFoldersFromParent lists folders under a specific parent resource.
// This method uses the direct ListFolders API for the specified parent.
func (c *Client) ListFoldersFromParent(ctx context.Context, parent string, _ *FetchOptions) ([]*Folder, error) {
	return c.searchAllAccessibleFolders(ctx, &FetchOptions{Parent: parent})
}

// Close releases any resources held by the fetcher.
func (c *Client) Close() error {
	if err := c.foldersClient.Close(); err != nil {
		return fmt.Errorf("failed to close folders client: %w", err)
	}

	return nil
}

// searchAllAccessibleFolders lists folders accessible to the current user, optionally filtered by parent.
func (c *Client) searchAllAccessibleFolders(ctx context.Context, opts *FetchOptions) ([]*Folder, error) {
	req := &resourcemanagerpb.SearchFoldersRequest{
		Query: "state:ACTIVE",
	}
	if opts.Parent != "" {
		req.Query += " AND parent:" + opts.Parent
	}

	it := c.foldersClient.SearchFolders(ctx, req)

	var folders []*Folder
	for {
		folder, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate all folders: %w", err)
		}

		folders = append(folders, FolderFromProto(folder))
	}

	return folders, nil
}
