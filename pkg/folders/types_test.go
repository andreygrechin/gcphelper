package folders_test

import (
	"testing"
	"time"

	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"github.com/andreygrechin/gcphelper/pkg/folders"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestNewFetchOptions(t *testing.T) {
	tests := map[string]struct {
		want *folders.FetchOptions
	}{
		"returns default options": {
			want: &folders.FetchOptions{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := folders.NewFetchOptions()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFolderFromProto(t *testing.T) {
	createTime := time.Now().UTC()
	updateTime := createTime.Add(time.Hour)

	tests := map[string]struct {
		input *resourcemanagerpb.Folder
		want  *folders.Folder
	}{
		"converts full folder": {
			input: &resourcemanagerpb.Folder{
				Name:        "folders/123456789",
				Parent:      "organizations/987654321",
				DisplayName: "Test Folder",
				State:       resourcemanagerpb.Folder_ACTIVE,
				CreateTime:  timestamppb.New(createTime),
				UpdateTime:  timestamppb.New(updateTime),
			},
			want: &folders.Folder{
				ID:          "123456789",
				Name:        "folders/123456789",
				DisplayName: "Test Folder",
				Parent:      "organizations/987654321",
				State:       "ACTIVE",
				CreateTime:  createTime,
				UpdateTime:  updateTime,
			},
		},
		"converts folder with minimal fields": {
			input: &resourcemanagerpb.Folder{
				Name:        "folders/123456789",
				Parent:      "organizations/987654321",
				DisplayName: "Test Folder",
				State:       resourcemanagerpb.Folder_DELETE_REQUESTED,
			},
			want: &folders.Folder{
				ID:          "123456789",
				Name:        "folders/123456789",
				DisplayName: "Test Folder",
				Parent:      "organizations/987654321",
				State:       "DELETE_REQUESTED",
				CreateTime:  time.Time{},
				UpdateTime:  time.Time{},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := folders.FolderFromProto(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
