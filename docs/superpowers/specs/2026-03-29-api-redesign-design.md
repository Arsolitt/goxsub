# API Redesign: Public Root Interface with Internal Packages

## Problem

The library's public API lives in the `xray/` subpackage, requiring users to import `goxsub/xray`. The `xray` name ties subscription parsing to a specific implementation. Output is limited to VLESS URI and podkop format, with no path for adding new protocols or output formats.

## Decision

Move the public API to the module root. Users import `goxsub` and call `goxsub.ParseSubscription()`, `goxsub.ExtractProxies()`, etc. Internal logic is split into focused packages. Protocol-specific types (VLESSProxy, SSProxy) share a common interface. Output formatters and protocol URI converters are separate packages.

## Package Structure

```
goxsub/
├── sub/           Subscription parsing and JSON models
│   ├── types.go       Subscription, Outbound, StreamSettings, User, etc.
│   └── parse.go       ParseSubscription([]byte) ([]Subscription, error)
├── proxy/         Proxy interface, concrete types, extraction, filtering
│   ├── proxy.go       Proxy interface, VLESSProxy
│   ├── extract.go     ExtractProxies([]sub.Subscription) []Proxy
│   └── filter.go      FilterByRemark([]Proxy, []string) []Proxy
├── protocol/      Protocol to URI conversion
│   └── vless.go       VLESSURI(*VLESSProxy) (string, error)
├── format/        Output formatters
│   └── podkop.go      Podkop([]proxy.Proxy, string) (string, error)
├── api.go         Re-exports all public types and functions
├── go.mod
└── cmd/goxsub/    CLI (imports goxsub, not internal packages)
```

The `xray/` package is removed entirely.

## Proxy Interface

Defined in `proxy/proxy.go`:

```go
type Proxy interface {
    Protocol() string
    Tag() string
    Remarks() string
}
```

Concrete types implement this interface. `VLESSProxy` is the first:

```go
type VLESSProxy struct {
    Outbound sub.Outbound
    Remarks  string
}
```

`Protocol()` returns `"vless"`, `Tag()` returns `Outbound.Tag`, `Remarks()` returns the Remarks field.

Adding a new protocol (trojan, ss) means adding a new type implementing `Proxy` in `proxy/`. `ExtractProxies` switches on `Outbound.Protocol` to create the correct type.

## Protocol Converters

`protocol/` contains pure functions that accept concrete proxy types:

```go
func VLESSURI(p *VLESSProxy) (string, error)
```

Adding a new protocol adds a new file: `protocol/trojan.go` with `TrojanURI(*TrojanProxy) (string, error)`.

## Output Formatters

`format/` contains output formatters. Each works with `[]proxy.Proxy` and switches on `Protocol()` to call the appropriate protocol converter:

```go
func Podkop(proxies []proxy.Proxy, section string) (string, error) {
    for _, p := range proxies {
        switch p.Protocol() {
        case "vless":
            uri, err := protocol.VLESSURI(p.(*VLESSProxy))
        }
    }
}
```

Future formatters: `format/singbox.go` (singbox JSON), `format/xray.go` (xray JSON config). Each new protocol adds a case to existing formatters.

## Root Re-exports

`api.go` re-exports all public types and functions from internal packages. Users import only `goxsub`:

```go
import "github.com/Arsolitt/goxsub"

subs, err := goxsub.ParseSubscription(data)
proxies := goxsub.ExtractProxies(subs)
proxies = goxsub.FilterByRemark(proxies, patterns)
```

## Migration

- Delete `xray/` package entirely
- Update `cmd/goxsub/main.go` to import `goxsub` instead of `goxsub/xray`
- External users change import path from `goxsub/xray` to `goxsub`
