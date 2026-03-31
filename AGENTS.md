## Project Overview

goxsub — Go library and CLI for parsing xray core JSON subscriptions and converting vless outbounds to vless:// URIs, sing-box outbounds, or podkop UCI commands.

Standard library only. No third-party dependencies. Go 1.26.

## Architecture

```
api.go         Public API surface: re-exports types and functions from subpackages
cmd/goxsub/    CLI binary: fetches subscription URL, prints vless:// URIs, sing-box JSON, or podkop UCI to stdout
sub/           JSON subscription parsing and types
proxy/         Proxy extraction and filtering
protocol/      URI conversion (vless:// and others)
format/        Output formatters (podkop uci commands, sing-box outbound JSON)
```

- `sub/types.go` — JSON models (Subscription, Outbound, StreamSettings, RealitySettings, etc.)
- `sub/parse.go` — `ParseSubscription([]byte) ([]Subscription, error)`
- `proxy/proxy.go` — `Proxy` interface, `VLESSProxy` type
- `proxy/extract.go` — `ExtractProxies([]Subscription) []Proxy`
- `proxy/filter.go` — `FilterByRemark([]Proxy, []string) []Proxy`
- `protocol/vless.go` — `VLESSURI(*VLESSProxy) (string, error)`, `ToVLESSURI(Proxy) (string, error)`, `ToURI(Proxy) (string, error)`
- `format/podkop.go` — `Podkop([]Proxy, string) ([]string, error)`
- `format/singbox.go` — `Singbox([]Proxy, SingboxConfig) ([]string, error)`, `SingboxConfig` struct

`api.go` re-exports all public types and functions for convenience (`import goxsub "github.com/Arsolitt/goxsub"`).

## File Generation

All generated files (build artifacts, test output, coverage reports, etc.) must go to `build/` or `/tmp`. Never write generated files into the repository source tree.

## Build/Test Commands

- Build: `go build -o build/goxsub ./cmd/goxsub/`
- Install: `go install github.com/Arsolitt/goxsub/cmd/goxsub@latest`
- Run all tests: `go test ./...`
- Run single test: `go test -run TestFunctionName ./package/`
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
- No `sync.Mutex`/`sync.RWMutex` as embedded struct fields (enforced by `embeddedstructfieldcheck`)

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

### Interfaces

- `Proxy` interface in `proxy/proxy.go` defines protocol-agnostic proxy contract
- Concrete types like `VLESSProxy` implement the interface; protocol packages use type assertions

### Logging

- Use `log/slog` — `log` (standard) is forbidden in non-main files (enforced by `depguard`)
- Context-aware slog methods required when context is in scope (enforced by `sloglint`)

### Random

- Use `math/rand/v2` — `math/rand` forbidden in non-test files (enforced by `depguard`)

## Testing

- Tests in same package (`package sub`, not `package sub_test`)
- Standard library assertions only (`t.Fatal`, `t.Errorf`, `t.Fatalf`) — no testify
- Test data is inline in test functions; no `testdata/` directory currently
- Use `t.TempDir()` instead of `os.TempDir()` in tests
- Many linters are relaxed in `_test.go` files (bodyclose, dupl, errcheck, funlen, goconst, gosec, noctx, gocognit, nestif)
- No `t.Helper()` required for test helpers

## Linting

Config: `.golangci.yaml` (maratori golden config for golangci-lint v2.6.2, 50+ linters enabled).

Key thresholds: funlen 100 lines/50 statements, gocognit 20, cyclop 30.

`nolint` directives must specify the linter name and provide a reason (except for `funlen`, `gocognit`, `golines` which need name but not reason).

## Commit Messages

Conventional commits: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`
