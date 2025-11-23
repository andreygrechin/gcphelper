# REVIEW

## Dependencies

- github.com/briandowns/spinner v1.23.2: Fine for UX, but gate by TTY/no‑progress to avoid garbage output in pipes/CI.
- google.golang.org/api, google.golang.org/grpc, google.golang.org/protobuf are correctly listed as direct dependencies since they are directly imported in the codebase.
- Recommendations:
  - Run go mod tidy periodically; use govulncheck (already in Makefile) before releases.
  - Consider pinning via -u=patch upgrades periodically and lock with CI to catch breaking changes.

## Naming & Data Model

- Interface/type names:
  - NewServiceFromContextWithLogger: verbose. NewService(ctx, log) communicates the same within the package namespace.

## CLI & UX

- Global flags:
  - globalFormat and globalVerbose as package globals are common with Cobra but harder to test/extend. Prefer a config struct carried via cmd.Context()/dependency injection, or bind flags to rootCmd.PersistentFlags() and read via cmd.Flags().
- Context:
  - Use cmd.Context() in RunE and propagate for cancellation (SIGINT/SIGTERM). Right now context.Background() prevents graceful cancel.
- Spinners:
  - Gate spinners on TTY and/or a --no-progress flag; avoid during non‑interactive runs and for JSON/CSV output. Also ensure they write to stderr not stdout to keep data streams clean. The current code does not set the spinner writer explicitly.

## Docs & Consistency

- ✅ README and ARCHITECTURE.md have been updated and now accurately reflect the current implementation:
  - Correctly documents organizations and folders commands
  - Accurately documents flags: --parent-folder, --parent-organization, --format, --verbose
  - Lists TTY detection as a future enhancement rather than current feature

## Suggested Changes

- CLI polish:
  - Use cmd.Context() in RunE paths; propagate context throughout services/clients.
  - Add --no-progress and only show spinner on TTY and in table output, writing to stderr.
