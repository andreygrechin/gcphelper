package folders

import (
	"context"
	"fmt"
	"time"

	"github.com/andreygrechin/gcphelper/internal/logger"
	"github.com/briandowns/spinner"
	"go.uber.org/zap"
)

const (
	spinnerSpeed = 100 * time.Millisecond
	spinnerStyle = 11
)

// Service provides high-level operations for working with Google Cloud folders.
type Service struct {
	fetcher Fetcher
	logger  logger.Logger
}

// NewServiceWithLogger creates a new folders service with the provided fetcher and logger.
func NewServiceWithLogger(fetcher Fetcher, log logger.Logger) *Service {
	return &Service{
		fetcher: fetcher,
		logger:  log,
	}
}

// NewServiceFromContextWithLogger creates a new folders service using application default credentials with logger.
func NewServiceFromContextWithLogger(ctx context.Context, log logger.Logger) (*Service, error) {
	client, err := NewClientFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if log == nil {
		log = logger.NewNoOpLogger()
	}

	return NewServiceWithLogger(client, log), nil
}

// ListFolders lists all accessible folders.
func (s *Service) ListFolders(ctx context.Context, opts *FetchOptions) ([]*Folder, error) {
	if opts == nil {
		opts = NewFetchOptions()
	}

	if s.logger != nil {
		if opts.Parent != "" {
			s.logger.Debug("fetching folders from parent", zap.String("parent", opts.Parent))
		} else {
			s.logger.Debug("fetching all accessible folders")
		}
	}

	// show progress indicator for potentially long-running operations
	spin := spinner.New(spinner.CharSets[spinnerStyle], spinnerSpeed)
	if opts.Parent != "" {
		spin.Suffix = fmt.Sprintf(" Fetching folders from parent %s...", opts.Parent)
	} else {
		spin.Suffix = " Fetching folders..."
	}
	spin.Start()
	defer spin.Stop()

	folders, err := s.fetcher.ListFolders(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list folders: %w", err)
	}

	if s.logger != nil {
		s.logger.Debug("successfully fetched folders", zap.Int("count", len(folders)))
	}

	return folders, nil
}

// Close releases any resources held by the service.
func (s *Service) Close() error {
	if s.fetcher != nil {
		if err := s.fetcher.Close(); err != nil {
			return fmt.Errorf("failed to close service fetcher: %w", err)
		}
	}

	return nil
}
