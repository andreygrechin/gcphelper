package organizations

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

// Service provides high-level operations for working with Google Cloud organizations.
type Service struct {
	fetcher Fetcher
	logger  logger.Logger
}

// NewServiceWithLogger creates a new organizations service with the provided fetcher and logger.
func NewServiceWithLogger(fetcher Fetcher, log logger.Logger) *Service {
	return &Service{
		fetcher: fetcher,
		logger:  log,
	}
}

// NewServiceFromContextWithLogger creates a new organizations service using
// application default credentials with logger.
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

// SearchOrganizations searches for organizations accessible to the caller.
func (s *Service) SearchOrganizations(ctx context.Context) ([]*Organization, error) {
	if s.logger != nil {
		s.logger.Info("searching for accessible organizations")
	}

	// show progress indicator for potentially long-running operations
	spin := spinner.New(spinner.CharSets[spinnerStyle], spinnerSpeed)
	spin.Suffix = " Searching for accessible organizations..."
	spin.Start()
	defer spin.Stop()

	organizations, err := s.fetcher.SearchOrganizations(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.Debug("failed to search organizations", zap.Error(err))
		}

		return nil, fmt.Errorf("failed to search organizations: %w", err)
	}

	if s.logger != nil {
		s.logger.Debug("successfully found organizations", zap.Int("count", len(organizations)))
	}

	return organizations, nil
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
