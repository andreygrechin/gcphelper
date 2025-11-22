# Architecture Documentation

This document describes the architecture and design patterns used in the gcphelper project.

## Overview

gcphelper is a CLI tool for fetching Google Cloud Platform resource information. It provides commands to list organizations and folders using the Google Cloud Resource Manager API.

## Project Structure

```text
gcphelper/
├── cmd/                        # CLI command definitions
│   ├── root.go                # Root command and global flags
│   ├── organizations.go       # Organizations command
│   └── folders.go             # Folders command
├── pkg/
│   ├── folders/               # Folder fetching logic
│   │   ├── fetcher.go        # API client and Fetcher interface
│   │   ├── service.go        # High-level service with UX features
│   │   └── types.go          # Data types and conversions
│   ├── organizations/         # Organization fetching logic
│   │   ├── fetcher.go        # API client and Fetcher interface
│   │   ├── service.go        # High-level service with UX features
│   │   └── types.go          # Data types and conversions
│   └── output/                # Output formatting
│       ├── formatter.go      # Format handling (table, JSON, CSV, ID)
│       └── adapters.go       # Resource conversion for output
└── internal/
    └── logger/                # Logging utilities
```

## Architecture Layers

The application follows a layered architecture separating CLI, service, and API concerns:

```text
┌────────────────────────────────────┐
│         CLI Layer (cmd/)           │  ← Command definitions, flags, output
├────────────────────────────────────┤
│      Service Layer (pkg/*/service) │  ← UX features (spinners), orchestration
├────────────────────────────────────┤
│    Fetcher Interface (pkg/*/fetcher)│ ← Contract for resource operations
├────────────────────────────────────┤
│     Client Implementation          │  ← Google Cloud API integration
├────────────────────────────────────┤
│   Data Types & Conversion          │  ← Domain models, protobuf conversion
└────────────────────────────────────┘
```

## Folder Fetching System (`pkg/folders/`)

### Components

#### 1. Fetcher Interface (`fetcher.go:13-23`)

Defines the contract for folder operations:

```go
type Fetcher interface {
    ListFolders(ctx context.Context, opts *FetchOptions) ([]*Folder, error)
    ListFoldersFromParent(ctx context.Context, parent string, opts *FetchOptions) ([]*Folder, error)
    Close() error
}
```

**Purpose:**
- Enables dependency injection and testability
- Allows mocking for unit tests
- Clean separation between interface and implementation

#### 2. Client Implementation (`fetcher.go:25-92`)

Implements the Fetcher interface using Google Cloud Resource Manager API:

```go
type Client struct {
    foldersClient *resourcemanager.FoldersClient
}
```

**Key Method:** `searchAllAccessibleFolders`

Uses the SearchFolders API with query-based filtering:

```go
req := &resourcemanagerpb.SearchFoldersRequest{
    Query: "state:ACTIVE",
}
if opts.Parent != "" {
    req.Query += " AND parent:" + opts.Parent
}
```

**Fetching Strategy:**
- Single API call using SearchFolders
- Simple sequential iteration through results
- Parent filtering via query parameter
- Returns all matching folders in one request

**Location:** `fetcher.go:68-92`

#### 3. Service Layer (`service.go`)

Provides high-level operations with user experience features:

```go
type Service struct {
    fetcher Fetcher
    logger  logger.Logger
}
```

**Key Features:**
- Progress indicators using spinner (pkg/briandowns/spinner)
- Debug logging for operations
- Resource lifecycle management
- Simplified API for CLI commands

**Spinner Integration:**

```go
spin := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
spin.Suffix = " Fetching folders..."
spin.Start()
defer spin.Stop()
```

**Location:** `service.go:45-79`

#### 4. Data Types (`types.go`)

**Folder Struct:**

```go
type Folder struct {
    ID          string    // Folder ID ("123456789")
    Name        string    // Resource name ("folders/123456789")
    DisplayName string    // Human-readable name
    Parent      string    // Parent resource
    State       string    // Lifecycle state
    CreateTime  time.Time
    UpdateTime  time.Time
}
```

**FetchOptions:**

```go
type FetchOptions struct {
    Parent string // Filter by parent (e.g., "organizations/123", "folders/456")
}
```

**Conversion:** `FolderFromProto` converts protobuf Folder to domain model

**Location:** `types.go:11-90`

### Data Flow

```text
User Command
    ↓
cmd/folders.go (runFoldersCommand)
    ↓
Validate flags (--parent-organization XOR --parent-folder)
    ↓
Create FetchOptions with Parent field
    ↓
service.ListFolders(ctx, opts)
    ↓
[Spinner starts]
    ↓
client.ListFolders(ctx, opts)
    ↓
client.searchAllAccessibleFolders(ctx, opts)
    ↓
Build SearchFoldersRequest with query
    ↓
Call GCP SearchFolders API
    ↓
Iterate results sequentially
    ↓
Convert each protobuf Folder to domain Folder
    ↓
Return []*Folder
    ↓
[Spinner stops]
    ↓
Format output (table/JSON/CSV/ID)
    ↓
Display to user
```

## Organization Fetching System (`pkg/organizations/`)

Similar layered architecture to folders:

1. **Fetcher Interface** - Defines SearchOrganizations contract
2. **Client Implementation** - Uses SearchOrganizations API
3. **Service Layer** - Adds spinner and logging
4. **Data Types** - Organization struct and conversion

**Key Difference:** No parent filtering - searches all accessible organizations.

## Output System (`pkg/output/`)

### Formatter

Supports multiple output formats:

```go
type Format string

const (
    FormatTable Format = "table"
    FormatJSON  Format = "json"
    FormatCSV   Format = "csv"
    FormatID    Format = "id"
)
```

**Format Implementations:**
- **Table**: Uses `github.com/jedib0t/go-pretty/v6` for formatted tables
- **JSON**: Standard `encoding/json` with indentation
- **CSV**: Standard `encoding/csv`
- **ID**: Outputs only resource IDs, one per line

### Resource Adapters

Convert domain types to generic Resource interface for formatting:

```go
type Resource interface {
    GetID() string
    GetDisplayName() string
    GetState() string
    GetCreateTime() time.Time
    GetUpdateTime() time.Time
}
```

**Functions:**
- `FoldersToResources()` - Converts []*folders.Folder
- `OrganizationsToResources()` - Converts []*organizations.Organization

## CLI Layer (`cmd/`)

### Root Command

Defines global flags available to all commands:

```go
var (
    globalFormat  string  // --format, -f (table|json|csv|id)
    globalVerbose bool    // --verbose, -v
)
```

**Initialization:**
- Creates logger
- Registers subcommands (folders, organizations)
- Sets up persistent flags

### Command Pattern

Each command follows this structure:

1. **Define flags** - Command-specific flags
2. **RunE handler** - Validation and execution
3. **Service creation** - Initialize service from context
4. **Fetch resources** - Call service methods
5. **Format output** - Use output formatter
6. **Error handling** - Enhanced error messages

**Example:** `cmd/folders.go`

- Flags: `--parent-organization`, `--parent-folder`
- Validation: Mutually exclusive parent flags
- Enhanced errors: Permission denied with helpful messages

## Design Patterns

### 1. Interface Segregation

```go
// Small, focused interfaces
type Fetcher interface {
    ListFolders(ctx context.Context, opts *FetchOptions) ([]*Folder, error)
    ListFoldersFromParent(ctx context.Context, parent string, opts *FetchOptions) ([]*Folder, error)
    Close() error
}
```

Benefits:
- Easy to mock for testing
- Clear contracts
- Minimal coupling

### 2. Dependency Injection

```go
// Service accepts interface, not concrete type
func NewServiceWithLogger(fetcher Fetcher, log logger.Logger) *Service
```

Benefits:
- Testability with mocks
- Flexible implementations
- Clear dependencies

### 3. Resource Pattern

```go
// Defer cleanup in service creation
service, err := folders.NewServiceFromContextWithLogger(ctx, log)
defer service.Close()
```

Benefits:
- Guaranteed cleanup
- No resource leaks
- Clear ownership

### 4. Adapter Pattern

```go
// Convert domain types to common interface
resources := output.FoldersToResources(folderList)
formatter.Format(resources, format, headers)
```

Benefits:
- Uniform output handling
- Format-agnostic domain types
- Reusable formatters

## Testing Strategy

### Mock-Based Testing

Generated mocks using `github.com/vektra/mockery/v3`:

```go
//go:generate mockery --name=Fetcher --output=mocks --outpkg=mocks
```

**Mock Usage:**
- Unit tests for Service layer
- Controlled API responses
- Error scenario testing

### Test Structure

Table-driven tests with subtests:

```go
tests := map[string]struct {
    setupMock func(*mocks.Fetcher)
    wantErr   bool
}{
    "success": {...},
    "api_error": {...},
}

for name, tt := range tests {
    t.Run(name, func(t *testing.T) {
        // Test implementation
    })
}
```

## Error Handling

### Error Wrapping

All errors are wrapped with context:

```go
return nil, fmt.Errorf("failed to list folders: %w", err)
```

### Enhanced Error Messages

CLI commands provide helpful error messages:

```go
if st.Code() == codes.PermissionDenied {
    return fmt.Errorf(`permission denied: insufficient permissions.

Ensure you have the 'resourcemanager.folders.list' permission.

Original error: %w`, err)
}
```

## Authentication

Uses Google Cloud Application Default Credentials:

```go
client, err := resourcemanager.NewFoldersClient(ctx)
// Automatically uses ADC from gcloud or GOOGLE_APPLICATION_CREDENTIALS
```

## Future Enhancements

1. **Caching**: Add optional response caching for repeated queries
2. **Filtering**: Client-side filtering by display name, state, etc.
3. **Sorting**: Configurable sort order for results
4. **Pagination**: Better handling of very large result sets
5. **Rate Limiting**: Adaptive throttling for API quotas
6. **Progress Bars**: Show progress for large operations
7. **TTY Detection**: Disable spinners in non-interactive mode
