# Sing-box Output Format Design

## Summary

Add sing-box outbound JSON output format to the goxsub CLI and library. Refactor existing formatters to return raw string slices, leaving presentation (joining, trailing commas) to the CLI layer.

## New Files

- `format/singbox.go` — sing-box outbound conversion
- `format/singbox_test.go` — tests

## Modified Files

- `format/podkop.go` — change return type from `(string, error)` to `([]string, error)`
- `format/podkop_test.go` — update to match new return type
- `cmd/goxsub/main.go` — new flags, sing-box format case, CLI-level output formatting
- `api.go` — re-export `Singbox` and `SingboxConfig`

## Types: `format/singbox.go`

```go
type SingboxConfig struct {
    OutboundPrefix string
    OutboundSuffix string
    KeepRemark     bool
    DNSResolver    string
}
```

`SingboxConfig` holds options for sing-box outbound generation. All fields are exported.

Unexported JSON mapping types (local to `format` package):

```go
type singboxOutbound struct {
    Type           string        `json:"type"`
    Tag            string        `json:"tag"`
    Server         string        `json:"server"`
    ServerPort     int           `json:"server_port"`
    UUID           string        `json:"uuid"`
    Flow           string        `json:"flow,omitempty"`
    TLS            *singboxTLS   `json:"tls,omitempty"`
    DomainResolver string        `json:"domain_resolver"`
}

type singboxTLS struct {
    Enabled  bool            `json:"enabled"`
    ServerName string        `json:"server_name,omitempty"`
    UTLS     *singboxUTLS    `json:"utls,omitempty"`
    Reality  *singboxReality `json:"reality,omitempty"`
}

type singboxUTLS struct {
    Enabled     bool   `json:"enabled"`
    Fingerprint string `json:"fingerprint"`
}

type singboxReality struct {
    Enabled   bool   `json:"enabled"`
    PublicKey string `json:"public_key"`
    ShortID   string `json:"short_id"`
}
```

## Function: `format.Singbox`

```go
func Singbox(proxies []proxy.Proxy, cfg SingboxConfig) ([]string, error)
```

Converts each proxy to a sing-box outbound JSON object string. Returns a slice where each element is a single JSON object (no trailing comma, no array wrapping). Tag generation rules:

- `KeepRemark=true`: tag = `OutboundPrefix + proxy.Remarks() + OutboundSuffix`
- `KeepRemark=false`: tag = `OutboundPrefix + strconv.Itoa(i+1) + OutboundSuffix` where `i` is the 0-based index in the filtered proxy list.

TLS block is populated from `StreamSettings`:
- When `Security == "reality"`: TLS enabled with UTLS (fingerprint from `RealitySettings.Fingerprint`) and Reality (public_key, short_id, server_name).
- When `Security == "tls"`: TLS enabled with server_name and optional UTLS fingerprint.
- Otherwise: TLS block omitted.

`DomainResolver` is always set from `cfg.DNSResolver` (default `"dns-local"`).

Non-vless proxies are skipped (not included in output). This is defensive since `ExtractProxies` already filters to vless only.

## Refactor: `format.Podkop`

```go
func Podkop(proxies []proxy.Proxy, section string) ([]string, error)
```

Return type changes from `(string, error)` to `([]string, error)`. Each element is one UCI command line. No `strings.Join`. The header (`uci del`) and footer (`service podkop restart`) lines are included as first and last elements of the slice.

## CLI: `cmd/goxsub/main.go`

### New Flags

| Flag | Type | Default | Scope |
|---|---|---|---|
| `--singbox-dns-resolver` | string | `"dns-local"` | sing-box only |
| `--singbox-outbound-prefix` | string | `""` | sing-box only |
| `--singbox-outbound-suffix` | string | `""` | sing-box only |
| `--keep-remark` | bool | `true` | all formats |

### Validation

- `--singbox-*` flags used with `--format != singbox` → error, exit 1 (same pattern as `--podkop-section`).
- `--keep-remark` is global, no format restriction.

### Output Formatting (CLI layer)

- `singbox`: each line printed as `line + ","` (trailing comma).
- `podkop`: each line printed as-is (no trailing comma).
- `uri` (default): each URI printed as-is.

### Flow Change

1. Parse subscription, extract proxies, filter by remark.
2. If `--keep-remark=false`, replace remarks with sequential numbers before formatting.
3. Format and output.

Remark replacement is done in CLI before formatting by iterating proxies and setting sequential remark values via a wrapper or by adding a `SetRemark` mechanism to `VLESSProxy`. This ensures all formatters (singbox tag, URI fragment) see the updated remark.

## `api.go` Additions

```go
type SingboxConfig = format.SingboxConfig
var Singbox = format.Singbox
```

## Testing

- `format/singbox_test.go`: test with REALITY security proxy, TLS security proxy, remark numbering, prefix/suffix, dns-resolver, empty proxies.
- `format/podkop_test.go`: update existing tests for new return type.
- No CLI integration tests (CLI is not tested currently).
