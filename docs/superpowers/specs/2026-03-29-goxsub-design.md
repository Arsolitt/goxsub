# goxsub — xray core subscription parser

Go library and CLI for parsing xray core JSON subscriptions and converting vless outbounds to vless:// URIs.

## Scope

- Protocol: vless only
- Transport: tcp, ws, grpc, h2, kcp
- Security: none, tls, reality
- Input format: JSON array only (no base64, no single objects)
- CLI: fetch URL, print vless:// URIs to stdout (no filtering, sorting, or file output)

## Project structure

```
goxsub/
├── xray/           # Library package (importable)
│   ├── types.go    # JSON models: Subscription, Outbound, StreamSettings, etc.
│   ├── parse.go    # ParseSubscription([]byte) ([]Subscription, error)
│   └── vless.go    # ToVLESSURI(Outbound) (string, error)
├── cmd/goxsub/     # CLI binary
│   └── main.go     # goxsub <subscription-url>
├── go.mod
└── go.sum
```

No third-party dependencies. Standard library only (`encoding/json`, `net/http`, `fmt`, `url`).

## Data models (xray/types.go)

Typed structs reflecting the xray core subscription JSON format:

```go
type Subscription struct {
    DNS       DNS        `json:"dns"`
    Inbounds  []Inbound  `json:"inbounds"`
    Log       Log        `json:"log"`
    Outbounds []Outbound `json:"outbounds"`
    Remarks   string     `json:"remarks"`
    Routing   Routing    `json:"routing"`
}

type Outbound struct {
    Mux            *Mux             `json:"mux,omitempty"`
    Protocol       string           `json:"protocol"`
    Settings       OutboundSettings `json:"settings"`
    StreamSettings StreamSettings   `json:"streamSettings"`
    Tag            string           `json:"tag"`
}

type StreamSettings struct {
    Network         string           `json:"network"`
    Security        string           `json:"security"`
    RealitySettings *RealitySettings `json:"realitySettings,omitempty"`
    TLSSettings     *TLSSettings     `json:"tlsSettings,omitempty"`
    TCPSettings     *TCPSettings     `json:"tcpSettings,omitempty"`
    WSSettings      *WSSettings      `json:"wsSettings,omitempty"`
    GRPCSettings    *GRPCSettings    `json:"grpcSettings,omitempty"`
    HTTPSettings    *HTTPSettings    `json:"httpSettings,omitempty"`
    KCPSettings     *KCPSettings     `json:"kcpSettings,omitempty"`
    Sockopt         *Sockopt        `json:"sockopt,omitempty"`
}
```

All optional fields are pointers so `omitempty` works correctly and "absent" can be distinguished from "zero value". Sub-structures for each transport/security type (`RealitySettings`, `TLSSettings`, `TCPSettings`, `WSSettings`, `GRPCSettings`, `HTTPSettings`, `KCPSettings`, `Mux`, `Sockopt`) contain only fields needed for vless:// URI construction. DNS, Inbounds, Log, and Routing are parsed minimally to avoid JSON decode errors but are not used for conversion.

## Parsing (xray/parse.go)

```go
func ParseSubscription(data []byte) ([]Subscription, error)
```

Validates that input is a JSON array. Returns an error if JSON is invalid or root element is not an array. Decodes each element into a `Subscription`.

## Conversion (xray/vless.go)

```go
type VLESSProxy struct {
    Outbound Outbound
    Remarks  string
}

func ExtractVLESSOutbounds(subs []Subscription) []VLESSProxy
```

Filters outbounds across all subscriptions where `tag == "proxy" || protocol == "vless"`. Returns a flat slice of matching outbounds with their parent subscription's remarks.

```go
func ToVLESSURI(o Outbound) (string, error)
```

Builds a vless:// URI:

```
vless://<uuid>@<address>:<port>?<params>#<fragment>
```

Steps:
1. Extract `uuid` and `address:port` from `settings.vnext[0]`
2. Add query params from `streamSettings`:
   - `type` — network (tcp, ws, grpc, h2, kcp)
   - `encryption` — from `users[0].encryption`
   - `security` — security (none, tls, reality)
   - For reality: `pbk`, `fp`, `sni`, `sid`, `spx`
   - For tls: `sni`, `alpn`, `fp`
   - For ws: `path`, `host`
   - For grpc: `serviceName`
   - For h2: `path`, `host`
   - For kcp: `type` (seed), `headerType`
   - `flow` — from `users[0].flow`
3. Fragment (`#`) — from `remarks` parameter (the parent subscription's `remarks` field)
4. Returns error if required fields are missing (uuid, address, port)

Signature: `ToVLESSURI(o Outbound, remarks string) (string, error)`

All parameters are URL-encoded via `url.QueryEscape`. Parameter order is fixed for deterministic output.

## Error handling

- `ParseSubscription` — error on invalid JSON or non-array root
- `ToVLESSURI` — error if no `vnext` or empty, missing uuid/address/port
- Non-matching outbounds (wrong protocol and tag) are silently skipped

## Testing (xray/*_test.go)

- Golden tests: fixed JSON array in `testdata/` with expected vless:// URIs
- One test per transport type: tcp+reality, tcp+tls, ws+tls, grpc+tls, h2+tls, kcp+none, none+none
- Edge cases: empty array, outbound without vnext, outbound with protocol!=vless and tag!=proxy

## CLI (cmd/goxsub/main.go)

```
goxsub <subscription-url>
```

- GET request to URL, read body
- Pass to `ParseSubscription`
- Filter via `ExtractVLESSOutbounds`
- Print each vless:// URI on a separate line to stdout
- On error: message to stderr, exit code 1
