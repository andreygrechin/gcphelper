package folders_test

import (
	"errors"
	"testing"
	"time"

	"github.com/andreygrechin/gcphelper/internal/logger"
	"github.com/andreygrechin/gcphelper/pkg/folders"
	"github.com/andreygrechin/gcphelper/pkg/folders/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Test error variables for err113 compliance.
var (
	errServiceTestAPIError    = errors.New("API error")
	errServiceTestOrgNotFound = errors.New("org not found")
)

func TestService_ListFolders(t *testing.T) {
	expectedFolders := []*folders.Folder{
		{
			ID:          "folders/123456789",
			Name:        "folders/123456789",
			DisplayName: "Test Folder 1",
			Parent:      "organizations/987654321",
			State:       "ACTIVE",
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		},
		{
			ID:          "folders/987654321",
			Name:        "folders/987654321",
			DisplayName: "Test Folder 2",
			Parent:      "organizations/987654321",
			State:       "ACTIVE",
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		},
	}

	tests := map[string]struct {
		opts        *folders.FetchOptions
		setupMock   func(*mocks.MockFetcher)
		want        []*folders.Folder
		wantErr     bool
		errContains string
	}{
		"successful fetch with default options": {
			opts: nil,
			setupMock: func(m *mocks.MockFetcher) {
				m.On("ListFolders", mock.Anything, mock.MatchedBy(func(opts *folders.FetchOptions) bool {
					return opts.Parent == ""
				})).Return(expectedFolders, nil)
			},
			want:    expectedFolders,
			wantErr: false,
		},
		"successful fetch with custom options": {
			opts: &folders.FetchOptions{
				Parent: "organizations/123456789",
			},
			setupMock: func(m *mocks.MockFetcher) {
				m.On("ListFolders", mock.Anything, mock.MatchedBy(func(opts *folders.FetchOptions) bool {
					return opts.Parent == "organizations/123456789"
				})).Return(expectedFolders, nil)
			},
			want:    expectedFolders,
			wantErr: false,
		},
		"fetcher returns error": {
			opts: nil,
			setupMock: func(m *mocks.MockFetcher) {
				m.On("ListFolders", mock.Anything, mock.Anything).Return(nil, errServiceTestAPIError)
			},
			want:        nil,
			wantErr:     true,
			errContains: "failed to list folders",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockFetcher := mocks.NewMockFetcher(t)
			testLogger, _ := logger.NewDevelopmentLogger()
			service := folders.NewServiceWithLogger(mockFetcher, testLogger)

			tt.setupMock(mockFetcher)

			got, err := service.ListFolders(t.Context(), tt.opts)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockFetcher.AssertExpectations(t)
		})
	}
}

func TestService_ListFoldersWithParent(t *testing.T) {
	expectedFolders := []*folders.Folder{
		{
			ID:          "folders/123456789",
			Name:        "folders/123456789",
			DisplayName: "Test Folder",
			Parent:      "organizations/987654321",
			State:       "ACTIVE",
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		},
	}

	tests := map[string]struct {
		parent      string
		setupMock   func(*mocks.MockFetcher)
		want        []*folders.Folder
		wantErr     bool
		errContains string
	}{
		"successful fetch with organization parent": {
			parent: "organizations/987654321",
			setupMock: func(m *mocks.MockFetcher) {
				m.On("ListFolders", mock.Anything, mock.MatchedBy(func(opts *folders.FetchOptions) bool {
					return opts.Parent == "organizations/987654321"
				})).Return(expectedFolders, nil)
			},
			want:    expectedFolders,
			wantErr: false,
		},
		"successful fetch with folder parent": {
			parent: "folders/123456789",
			setupMock: func(m *mocks.MockFetcher) {
				m.On("ListFolders", mock.Anything, mock.MatchedBy(func(opts *folders.FetchOptions) bool {
					return opts.Parent == "folders/123456789"
				})).Return(expectedFolders, nil)
			},
			want:    expectedFolders,
			wantErr: false,
		},
		"fetcher returns error": {
			parent: "organizations/987654321",
			setupMock: func(m *mocks.MockFetcher) {
				m.On("ListFolders", mock.Anything, mock.Anything).Return(nil, errServiceTestOrgNotFound)
			},
			want:        nil,
			wantErr:     true,
			errContains: "failed to list folders",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			mockFetcher := mocks.NewMockFetcher(t)
			testLogger, _ := logger.NewDevelopmentLogger()
			service := folders.NewServiceWithLogger(mockFetcher, testLogger)

			tt.setupMock(mockFetcher)

			opts := &folders.FetchOptions{
				Parent: tt.parent,
			}
			got, err := service.ListFolders(t.Context(), opts)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockFetcher.AssertExpectations(t)
		})
	}
}
