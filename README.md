# gcphelper

gcphelper is the CLI tool to fetch information from Google Cloud.

[![license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/andreygrechin/gcphelper/blob/main/LICENSE)

## Features

- List all projects in an organization
- List all folders in an organization with async fetching
- Export folder information in multiple formats (table, JSON, CSV)
- Configurable concurrency for API requests

## Installation

### go install

```shell
go install github.com/andreygrechin/gcphelper@latest
```

### Manually

Download the pre-compiled binaries from [the releases page](https://github.com/andreygrechin/gcphelper/releases/) and copy them to a desired location.

## Prerequisites

- Go 1.25 or later (for building from source)
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

To list folders, your account needs the following IAM permissions:

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

### List Projects

```shell
# Fetch all accessible projects
gcphelper projects

# Fetch projects recursively from a specific organization
gcphelper projects --org <org-id>
```

### List Folders

```shell
# Fetch all accessible folders
gcphelper folders

# Fetch folders recursively from a specific organization
gcphelper folders --org <org-id>

# Output folders in JSON format
gcphelper folders --format json

# Output folders in CSV format
gcphelper folders --format csv

# Fetch folders with custom concurrency (default: 10)
gcphelper folders --concurrency 20

# Combine options
gcphelper folders --org 123456789 --format json --concurrency 15
```

## License

This project is licensed under the [MIT License](LICENSE).

`SPDX-License-Identifier: MIT`
