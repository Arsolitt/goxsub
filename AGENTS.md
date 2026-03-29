## Project Overview

goxsub — Go library and CLI for parsing xray core JSON subscriptions and converting vless outbounds to vless:// URIs.

Standard library only. No third-party dependencies. Go 1.26.

## Architecture

```
cmd/goxsub/    CLI binary: fetches subscription URL, prints vless:// URIs to stdout
xray/          Library package (importable): types, parsing, conversion
```

- `xray/types.go` — JSON models (Subscription, Outbound, StreamSettings, etc.)
- `xray/parse.go` — `ParseSubscription([]byte) ([]Subscription, error)`
- `xray/vless.go` — `ExtractVLESSOutbounds([]Subscription) []VLESSProxy`, `ToVLESSURI(Outbound, string) (string, error)`

Design spec: `docs/superpowers/specs/2026-03-29-goxsub-design.md`
Implementation plan: `docs/superpowers/plans/2026-03-29-goxsub.md`

## File Generation

All generated files (build artifacts, test output, coverage reports, etc.) must go to `build/` or `/tmp`. Never write generated files into the repository source tree.

## Build/Test Commands

- Build: `go build -o build/goxsub ./cmd/goxsub/`
- Install: `go install github.com/Arsolitt/goxsub/cmd/goxsub@latest`
- Run all tests: `go test ./...`
- Run single test: `go test -run TestFunctionName ./internal/package`
- Run tests with coverage: `go test -cover ./...`
- Lint code: `golangci-lint run`
- Lint and auto-fix: `golangci-lint run --fix`
- Always run `golangci-lint run --fix` first to auto-resolve issues before fixing remaining ones manually

## Code Conventions

### General

- Tab indentation (Go standard)
- Max line length: 120 (enforced by `golines`)
- No comments unless asked
- Exported types and functions must have godoc comments (enforced by `revive`)
- Comments must end with a period (enforced by `godot`)

### Imports

- `goimports` formatter with local prefix `github.com/Arsolitt/goxsub`
- Group order: stdlib, external (none currently), local

### Error Handling

- Always check errors (enforced by `errcheck`)
- Use `fmt.Errorf("context: %w", err)` for wrapping (enforced by `errorlint`)
- Sentinel errors: prefix with `Err` (enforced by `errname`)
- Error types: suffix with `Error` (enforced by `errname`)
- Never return `(nil, nil)` (enforced by `nilnil`)
- Printf-like functions: suffix with `f` (enforced by `goprintffuncname`)

### Types

- Optional struct fields use pointers with `omitempty` to distinguish absent from zero-value
- Unused fields in JSON (DNS, Inbounds, Log, Routing) use `json.RawMessage`
- No meaningless package names (`utils`, `helpers`)

### Logging

- Use `log/slog` — `log` (standard) is forbidden in non-main files (enforced by `depguard`)
- Context-aware slog methods required when context is in scope (enforced by `sloglint`)

### Random

- Use `math/rand/v2` — `math/rand` forbidden in non-test files (enforced by `depguard`)

## Testing

- Tests in same package (`package xray`, not `package xray_test`)
- Standard library assertions only (no testify)
- Golden test data in `testdata/`
- Use `t.TempDir()` instead of `os.TempDir()` in tests
- Many linters are relaxed in `_test.go` files
- No `t.Helper()` required for test helpers

## Linting

Config: `.golangci.yaml` (maratori golden config for golangci-lint v2.6.2, 50+ linters enabled).

`nolint` directives must specify the linter name and provide a reason (except for `funlen`, `gocognit`, `golines`).

## Commit Messages

Conventional commits: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`
