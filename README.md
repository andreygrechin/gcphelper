# gcphelper

gcphelper is the CLI tool to fetch information from Google Cloud.

[![license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/andreygrechin/gcphelper/blob/main/LICENSE)

## Features

- List all accessible organizations
- List all accessible folders
- Search folders by parent organization or folder
- Export information in multiple formats (table, JSON, CSV, ID)

## Installation

### go install

```shell
go install github.com/andreygrechin/gcphelper@latest
```

### Manually

Download the pre-compiled binaries from [the releases page](https://github.com/andreygrechin/gcphelper/releases/) and copy them to a desired location.

## Prerequisites

- Go 1.24 or later (for building from source)
- Google Cloud SDK (gcloud) installed and configured
- Valid Google Cloud credentials with appropriate permissions

## Authentication

This tool uses Google Cloud Application Default Credentials. Before using gcphelper, authenticate using:

```shell
gcloud auth application-default login
```

For service account authentication, set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable:

```shell
export GOOGLE_APPLICATION_CREDENTIALS="/path/to/service-account-key.json"
```

## Required Permissions

### For Organizations

To list organizations, your account needs:

- `resourcemanager.organizations.get` permission

### For Folders

To list folders, your account needs:

- `resourcemanager.folders.list` on the organization or parent folders
- `resourcemanager.folders.get` on individual folders

## Usage

### Available Commands

```shell
# View help for all commands
gcphelper --help

# Generate shell completion
gcphelper completion <shell>
```

### Global Flags

All commands support these global flags:

- `--format`, `-f`: Output format (table, json, csv, id) - default: table
- `--verbose`, `-v`: Show additional output like counts and status messages

### List Organizations

List all Google Cloud organizations accessible to your credentials.

```shell
# List all accessible organizations
gcphelper organizations

# List organizations in JSON format
gcphelper --format json organizations

# List only organization IDs for scripting
gcphelper --format id organizations

# Pipe organization IDs to other commands
gcphelper -f id organizations | xargs -I {} gcloud resource-manager organizations describe {}

# List organizations with verbose output
gcphelper --verbose organizations

# Use the short alias
gcphelper org
```

### List Folders

List Google Cloud folders using the SearchFolders API to discover all accessible folders.

```shell
# List all accessible folders
gcphelper folders

# List folders from a specific organization
gcphelper folders --parent-organization 123456789

# List folders under a specific parent folder
gcphelper folders --parent-folder 987654321

# List folders in JSON format
gcphelper --format json folders

# List folders in CSV format
gcphelper --format csv folders

# List only folder IDs for scripting
gcphelper --format id folders

# Pipe folder IDs to other commands
gcphelper -f id folders | xargs -I {} gcloud resource-manager folders describe {}

# List folders with verbose output
gcphelper --verbose folders

# Combine options
gcphelper --format json folders --parent-organization 123456789
```

#### Folder Command Flags

- `--parent-organization`, `-o`: Filter folders by parent organization ID
- `--parent-folder`, `-p`: Filter folders by parent folder ID

Note: You cannot specify both `--parent-organization` and `--parent-folder` at the same time.

## Output Formats

### Table (default)

Human-readable table format with columns for all resource attributes.

### JSON

Machine-readable JSON format for programmatic processing.

### CSV

Comma-separated values format for spreadsheet imports.

### ID

Outputs only resource IDs, one per line - useful for piping to other commands.

## License

This project is licensed under the [MIT License](LICENSE).

`SPDX-License-Identifier: MIT`
