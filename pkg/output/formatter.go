package output

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

// ErrUnsupportedOutputFormat is returned when an unsupported output format is requested.
var ErrUnsupportedOutputFormat = errors.New("unsupported output format")

// Constants for resource types.
const defaultResourceType = "resources"

// Format represents the output format type.
type Format string

// Output format constants.
const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
	FormatID    Format = "id"
)

// Resource represents a generic cloud resource with common fields.
type Resource interface {
	GetID() string
	GetDisplayName() string
	GetState() string
	GetCreateTime() time.Time
	GetUpdateTime() time.Time
	TableRow() []interface{}
}

// Formatter handles output formatting for resources.
type Formatter struct {
	writer       io.Writer
	verbose      bool
	resourceType string
}

// NewFormatterWithType creates a new formatter with a specific resource type for messages.
func NewFormatterWithType(writer io.Writer, verbose bool, resourceType string) *Formatter {
	if writer == nil {
		writer = os.Stdout
	}

	return &Formatter{
		writer:       writer,
		verbose:      verbose,
		resourceType: resourceType,
	}
}

// Format outputs the resources in the specified format.
func (f *Formatter) Format(resources []Resource, format Format, headers []string) error {
	switch format {
	case FormatJSON:
		return f.formatJSON(resources)
	case FormatCSV:
		return f.formatCSV(resources, headers)
	case FormatTable:
		return f.formatTable(resources, headers)
	case FormatID:
		return f.formatID(resources)
	default:
		return fmt.Errorf("%w: %s", ErrUnsupportedOutputFormat, format)
	}
}

func (f *Formatter) formatJSON(resources []Resource) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(resources); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func (f *Formatter) formatTable(resources []Resource, headers []string) error {
	if len(resources) == 0 {
		if f.verbose {
			resourceType := f.resourceType
			if resourceType == "" {
				resourceType = "resources"
			}
			fmt.Fprintf(os.Stderr, "No %s found.\n", resourceType)
		}

		return nil
	}

	t := table.NewWriter()
	t.SetOutputMirror(f.writer)
	t.SetStyle(table.StyleDefault)

	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	for _, resource := range resources {
		t.AppendRow(resource.TableRow())
	}

	if f.verbose {
		resourceType := f.resourceType
		if resourceType == "" {
			resourceType = defaultResourceType
		}
		t.SetCaption("Total %s: %d", resourceType, len(resources))
	}
	t.Render()

	return nil
}

func (f *Formatter) formatCSV(resources []Resource, headers []string) error {
	t := table.NewWriter()
	t.SetOutputMirror(f.writer)
	t.SetStyle(table.StyleDefault)

	headerRow := make(table.Row, len(headers))
	for i, h := range headers {
		headerRow[i] = h
	}
	t.AppendHeader(headerRow)

	for _, resource := range resources {
		t.AppendRow(resource.TableRow())
	}

	t.RenderCSV()

	return nil
}

func (f *Formatter) formatID(resources []Resource) error {
	if len(resources) == 0 && f.verbose {
		resourceType := f.resourceType
		if resourceType == "" {
			resourceType = defaultResourceType
		}
		fmt.Fprintf(os.Stderr, "No %s found.\n", resourceType)

		return nil
	}

	for _, resource := range resources {
		if _, err := fmt.Fprintln(f.writer, resource.GetID()); err != nil {
			return fmt.Errorf("failed to write resource ID: %w", err)
		}
	}

	return nil
}
