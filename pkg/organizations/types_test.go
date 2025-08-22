package organizations_test

import (
	"testing"
	"time"

	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"github.com/andreygrechin/gcphelper/pkg/organizations"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestOrgFromProto(t *testing.T) {
	testCases := map[string]struct {
		input *resourcemanagerpb.Organization
		want  *organizations.Organization
	}{
		"valid organization": {
			input: &resourcemanagerpb.Organization{
				Name:        "organizations/123456789",
				DisplayName: "Test Organization",
				State:       resourcemanagerpb.Organization_ACTIVE,
				CreateTime:  timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdateTime:  timestamppb.New(time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)),
			},
			want: &organizations.Organization{
				ID:          "123456789",
				Name:        "organizations/123456789",
				DisplayName: "Test Organization",
				State:       "ACTIVE",
				CreateTime:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdateTime:  time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		"organization without timestamps": {
			input: &resourcemanagerpb.Organization{
				Name:        "organizations/987654321",
				DisplayName: "Test Organization 2",
				State:       resourcemanagerpb.Organization_DELETE_REQUESTED,
			},
			want: &organizations.Organization{
				ID:          "987654321",
				Name:        "organizations/987654321",
				DisplayName: "Test Organization 2",
				State:       "DELETE_REQUESTED",
				CreateTime:  time.Time{},
				UpdateTime:  time.Time{},
			},
		},
		"nil organization": {
			input: nil,
			want:  nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := organizations.OrganizationFromProto(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
