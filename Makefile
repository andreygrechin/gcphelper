.PHONY: all lint format fmt test build coverage-report coverage-report-html security-scan check-clean release-test release

# Build variables
VERSION    := $(shell git describe --tags --always --dirty)
COMMIT     := $(shell git rev-parse --short HEAD)
BUILDTIME  := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
MOD_PATH   := $(shell go list -m)
PACKAGES   := $(shell go list ./... | grep -vE '/mocks|/testhelpers')
APP_NAME   := gcphelper
GOCOVERDIR := ./covdatafiles

all: lint format test build

lint: fmt
	golangci-lint run --show-stats

format:
	gofumpt -l -w $(shell find . -type f -name "*.go" -not -path "./vendor/*" -not -name "mock_*.go")

fmt: format

test:
	mockery
	go test ./...

build:
	CGO_ENABLED=0 \
	go build \
		-ldflags \
		"-s \
		-w \
		-X main.Version=$(VERSION) \
		-X main.BuildTime=$(BUILDTIME) \
		-X main.Commit=$(COMMIT)" \
		-o bin/$(APP_NAME)

coverage-report:
	go test -race -coverprofile="${GOCOVERDIR}/coverage.out" ${PACKAGES} && go tool cover -func="${GOCOVERDIR}/coverage.out"

coverage-report-html:
	go test -race -coverprofile="${GOCOVERDIR}/coverage.out" ${PACKAGES} && go tool cover -func="${GOCOVERDIR}/coverage.out"
	go tool cover -html="${GOCOVERDIR}/coverage.out"

security-scan:
	gosec ./...
	govulncheck

check-clean:
	@if [ -n "$(shell git status --porcelain)" ]; then \
		echo "Error: Dirty working tree. Commit or stash changes before proceeding."; \
		exit 1; \
	fi

release-test: lint test security-scan
	goreleaser check
	goreleaser release --snapshot --clean

release: check-clean lint test security-scan
	goreleaser release --clean
