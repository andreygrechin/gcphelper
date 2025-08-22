package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/andreygrechin/gcphelper/pkg/output"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockResource implements the Resource interface for testing.
type mockResource struct {
	id          string
	displayName string
	state       string
	createTime  time.Time
	updateTime  time.Time
}

func (m *mockResource) GetID() string            { return m.id }
func (m *mockResource) GetDisplayName() string   { return m.displayName }
func (m *mockResource) GetState() string         { return m.state }
func (m *mockResource) GetCreateTime() time.Time { return m.createTime }
func (m *mockResource) GetUpdateTime() time.Time { return m.updateTime }
func (m *mockResource) TableRow() []interface{} {
	return table.Row{
		m.id,
		m.displayName,
		m.state,
		m.createTime.Format("2006-01-02 15:04:05"),
		m.updateTime.Format("2006-01-02 15:04:05"),
	}
}

func createTestResources() []output.Resource {
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	return []output.Resource{
		&mockResource{
			id:          "123",
			displayName: "Test Resource 1",
			state:       "ACTIVE",
			createTime:  baseTime,
			updateTime:  baseTime.Add(time.Hour),
		},
		&mockResource{
			id:          "456",
			displayName: "Test Resource 2",
			state:       "INACTIVE",
			createTime:  baseTime.Add(time.Minute),
			updateTime:  baseTime.Add(2 * time.Hour),
		},
	}
}

func TestFormatter_FormatJSON(t *testing.T) {
	tests := map[string]struct {
		resources []output.Resource
		verbose   bool
		want      string
	}{
		"formats resources as JSON": {
			resources: createTestResources(),
			verbose:   false,
			want:      "formatted JSON output",
		},
		"formats empty resources": {
			resources: []output.Resource{},
			verbose:   true,
			want:      "[]",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			formatter := output.NewFormatterWithType(&buf, tt.verbose, "")
			headers := []string{"ID", "Name", "State", "Created", "Updated"}

			err := formatter.Format(tt.resources, output.FormatJSON, headers)
			require.NoError(t, err)

			// verify valid JSON was produced
			var result []interface{}
			err = json.Unmarshal(buf.Bytes(), &result)
			require.NoError(t, err)

			assert.Len(t, result, len(tt.resources))
		})
	}
}

func TestFormatter_FormatTable(t *testing.T) {
	tests := map[string]struct {
		resources []output.Resource
		verbose   bool
		headers   []string
	}{
		"formats resources as table": {
			resources: createTestResources(),
			verbose:   false,
			headers:   []string{"ID", "Name", "State", "Created", "Updated"},
		},
		"formats empty resources with verbose": {
			resources: []output.Resource{},
			verbose:   true,
			headers:   []string{"ID", "Name"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			formatter := output.NewFormatterWithType(&buf, tt.verbose, "")

			err := formatter.Format(tt.resources, output.FormatTable, tt.headers)
			require.NoError(t, err)

			result := buf.String()
			if len(tt.resources) > 0 {
				// verify table headers are present (they are uppercase in the output)
				for _, header := range tt.headers {
					assert.Contains(t, result, strings.ToUpper(header))
				}
				// verify resource data is present
				for _, resource := range tt.resources {
					assert.Contains(t, result, resource.GetID())
					assert.Contains(t, result, resource.GetDisplayName())
				}
			}
		})
	}
}

func TestFormatter_FormatCSV(t *testing.T) {
	tests := map[string]struct {
		resources []output.Resource
		headers   []string
	}{
		"formats resources as CSV": {
			resources: createTestResources(),
			headers:   []string{"ID", "Name", "State", "Created", "Updated"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			formatter := output.NewFormatterWithType(&buf, false, "")

			err := formatter.Format(tt.resources, output.FormatCSV, tt.headers)
			require.NoError(t, err)

			result := buf.String()
			// verify CSV headers
			for _, header := range tt.headers {
				assert.Contains(t, result, header)
			}
			// verify CSV content
			lines := strings.Split(strings.TrimSpace(result), "\n")
			assert.GreaterOrEqual(t, len(lines), len(tt.resources)+1) // +1 for header
		})
	}
}

func TestFormatter_FormatID(t *testing.T) {
	tests := map[string]struct {
		resources []output.Resource
		verbose   bool
		expected  []string
	}{
		"formats resource IDs": {
			resources: createTestResources(),
			verbose:   false,
			expected:  []string{"123", "456"},
		},
		"handles empty resources with verbose": {
			resources: []output.Resource{},
			verbose:   true,
			expected:  []string{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			formatter := output.NewFormatterWithType(&buf, tt.verbose, "")

			err := formatter.Format(tt.resources, output.FormatID, nil)
			require.NoError(t, err)

			result := strings.TrimSpace(buf.String())
			if len(tt.expected) == 0 {
				assert.Empty(t, result)
			} else {
				lines := strings.Split(result, "\n")
				assert.Equal(t, tt.expected, lines)
			}
		})
	}
}

func TestFormatter_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	formatter := output.NewFormatterWithType(&buf, false, "")
	resources := createTestResources()

	err := formatter.Format(resources, "invalid", []string{"ID"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported output format")
}
