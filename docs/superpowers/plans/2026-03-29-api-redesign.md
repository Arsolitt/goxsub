# API Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move the public API to the module root, split internal logic into `sub/`, `proxy/`, `protocol/`, `format/` packages, and delete the `xray/` package.

**Architecture:** Subscription parsing and JSON models in `sub/`. Proxy interface and concrete types in `proxy/`. Protocol URI converters in `protocol/`. Output formatters in `format/`. Root `api.go` re-exports everything. CLI imports only `goxsub`.

**Tech Stack:** Go 1.26, standard library only.

---

### Task 1: Create `sub/` package — types and parsing

**Files:**
- Create: `sub/types.go`
- Create: `sub/parse.go`
- Create: `sub/parse_test.go`
- Delete: `xray/types.go`
- Delete: `xray/parse.go`
- Delete: `xray/parse_test.go`

- [ ] **Step 1: Create `sub/types.go`**

Copy all types from `xray/types.go` except `VLESSProxy`. Change `package xray` to `package sub`. Add godoc to `Subscription`, `Outbound`, `OutboundSettings`, `VNext`, `User`, `StreamSettings`.

```go
package sub

import "encoding/json"

// Subscription represents a single subscription entry from a JSON subscription feed.
type Subscription struct {
	DNS       json.RawMessage `json:"dns"`
	Inbounds  json.RawMessage `json:"inbounds"`
	Log       json.RawMessage `json:"log"`
	Outbounds []Outbound      `json:"outbounds"`
	Remarks   string          `json:"remarks"`
	Routing   json.RawMessage `json:"routing"`
}

// Outbound represents a proxy outbound configuration.
type Outbound struct {
	StreamSettings StreamSettings   `json:"streamSettings"`
	Protocol       string           `json:"protocol"`
	Tag            string           `json:"tag"`
	Settings       OutboundSettings `json:"settings"`
}

// OutboundSettings holds the settings for an outbound connection.
type OutboundSettings struct {
	Vnext []VNext `json:"vnext,omitempty"`
}

// VNext represents a vnext server configuration.
type VNext struct {
	Address string `json:"address"`
	Users   []User `json:"users"`
	Port    int    `json:"port"`
}

// User represents a user account in a vnext entry.
type User struct {
	Encryption string `json:"encryption"`
	Flow       string `json:"flow,omitempty"`
	ID         string `json:"id"`
	Level      int    `json:"level,omitempty"`
}

// StreamSettings holds the transport and security configuration for an outbound.
type StreamSettings struct {
	RealitySettings *RealitySettings `json:"realitySettings,omitempty"`
	TLSSettings     *TLSSettings     `json:"tlsSettings,omitempty"`
	TCPSettings     *TCPSettings     `json:"tcpSettings,omitempty"`
	WSSettings      *WSSettings      `json:"wsSettings,omitempty"`
	GRPCSettings    *GRPCSettings    `json:"grpcSettings,omitempty"`
	HTTPSettings    *HTTPSettings    `json:"httpSettings,omitempty"`
	KCPSettings     *KCPSettings     `json:"kcpSettings,omitempty"`
	Network         string           `json:"network"`
	Security        string           `json:"security"`
}

// RealitySettings holds REALITY TLS settings.
type RealitySettings struct {
	PublicKey   string `json:"publicKey"`
	Fingerprint string `json:"fingerprint"`
	ServerName  string `json:"serverName"`
	ShortID     string `json:"shortId"`
	SpiderX     string `json:"spiderX,omitempty"`
}

// TLSSettings holds standard TLS settings.
type TLSSettings struct {
	ServerName  string   `json:"serverName,omitempty"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	ALPN        []string `json:"alpn,omitempty"`
}

// TCPSettings holds TCP transport settings.
type TCPSettings struct {
	Header *TCPHeader `json:"header,omitempty"`
}

// TCPHeader holds the TCP header type.
type TCPHeader struct {
	Type string `json:"type,omitempty"`
}

// WSSettings holds WebSocket transport settings.
type WSSettings struct {
	Path string `json:"path,omitempty"`
	Host string `json:"host,omitempty"`
}

// GRPCSettings holds gRPC transport settings.
type GRPCSettings struct {
	ServiceName string `json:"serviceName,omitempty"`
}

// HTTPSettings holds HTTP/2 transport settings.
type HTTPSettings struct {
	Path string   `json:"path,omitempty"`
	Host []string `json:"host,omitempty"`
}

// KCPSettings holds mKCP transport settings.
type KCPSettings struct {
	HeaderType *KCPHeader `json:"header,omitempty"`
	Seed       string     `json:"seed,omitempty"`
}

// KCPHeader holds the mKCP header type.
type KCPHeader struct {
	Type string `json:"type,omitempty"`
}
```

- [ ] **Step 2: Create `sub/parse.go`**

Copy `ParseSubscription` from `xray/parse.go`. Change `package xray` to `package sub`. Keep the function identical.

```go
package sub

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ParseSubscription parses a JSON array of subscription objects from raw byte data.
func ParseSubscription(data []byte) ([]Subscription, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	t, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("parse subscription: %w", err)
	}
	delim, ok := t.(json.Delim)
	if !ok || delim != '[' {
		return nil, fmt.Errorf("parse subscription: expected JSON array, got %v", t)
	}
	subs := make([]Subscription, 0)
	for dec.More() {
		var sub Subscription
		if err := dec.Decode(&sub); err != nil {
			return nil, fmt.Errorf("parse subscription: decode element: %w", err)
		}
		subs = append(subs, sub)
	}
	return subs, nil
}
```

- [ ] **Step 3: Create `sub/parse_test.go`**

Copy tests from `xray/parse_test.go`. Change `package xray` to `package sub`. Tests are identical since they only use `Subscription` and `ParseSubscription`.

```go
package sub

import "testing"

func TestParseSubscription_ValidArray(t *testing.T) {
	data := []byte(`[
		{
			"dns": {},
			"inbounds": [],
			"log": {},
			"outbounds": [
				{"protocol": "vless", "tag": "proxy", "settings": {"vnext": []}, "streamSettings": {"network": "tcp", "security": "none"}}
			],
			"remarks": "test",
			"routing": {}
		}
	]`)

	subs, err := ParseSubscription(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(subs) != 1 {
		t.Fatalf("expected 1 subscription, got %d", len(subs))
	}
	if subs[0].Remarks != "test" {
		t.Errorf("expected remarks 'test', got %q", subs[0].Remarks)
	}
	if len(subs[0].Outbounds) != 1 {
		t.Fatalf("expected 1 outbound, got %d", len(subs[0].Outbounds))
	}
	if subs[0].Outbounds[0].Protocol != "vless" {
		t.Errorf("expected protocol 'vless', got %q", subs[0].Outbounds[0].Protocol)
	}
}

func TestParseSubscription_InvalidJSON(t *testing.T) {
	_, err := ParseSubscription([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestParseSubscription_NonArray(t *testing.T) {
	_, err := ParseSubscription([]byte(`{"key": "value"}`))
	if err == nil {
		t.Fatal("expected error for non-array JSON")
	}
}

func TestParseSubscription_EmptyArray(t *testing.T) {
	subs, err := ParseSubscription([]byte(`[]`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(subs) != 0 {
		t.Errorf("expected 0 subscriptions, got %d", len(subs))
	}
}

func TestParseSubscription_MultipleSubscriptions(t *testing.T) {
	data := []byte(`[
		{"outbounds": [], "remarks": "first"},
		{"outbounds": [], "remarks": "second"}
	]`)

	subs, err := ParseSubscription(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(subs) != 2 {
		t.Fatalf("expected 2 subscriptions, got %d", len(subs))
	}
	if subs[0].Remarks != "first" || subs[1].Remarks != "second" {
		t.Errorf("remarks mismatch: %q, %q", subs[0].Remarks, subs[1].Remarks)
	}
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./sub/...`
Expected: PASS (all 5 tests)

- [ ] **Step 5: Run lint**

Run: `golangci-lint run --fix`
Expected: no errors

- [ ] **Step 6: Commit**

```bash
git add sub/ && git commit -m "refactor: extract sub package with subscription types and parsing"
```

---

### Task 2: Create `proxy/` package — interface, types, extraction, filtering

**Files:**
- Create: `proxy/proxy.go`
- Create: `proxy/extract.go`
- Create: `proxy/extract_test.go`
- Create: `proxy/filter.go`
- Create: `proxy/filter_test.go`

- [ ] **Step 1: Create `proxy/proxy.go`**

Define the `Proxy` interface and `VLESSProxy` concrete type. `VLESSProxy` wraps `sub.Outbound`.

```go
package proxy

import "github.com/Arsolitt/goxsub/sub"

// Proxy is the interface implemented by all protocol-specific proxy types.
type Proxy interface {
	Protocol() string
	Tag() string
	Remarks() string
}

// VLESSProxy represents a VLESS protocol proxy extracted from a subscription.
type VLESSProxy struct {
	Outbound sub.Outbound
	Remarks  string
}

func (v *VLESSProxy) Protocol() string { return "vless" }
func (v *VLESSProxy) Tag() string      { return v.Outbound.Tag }
func (v *VLESSProxy) Remarks() string   { return v.Remarks }
```

- [ ] **Step 2: Create `proxy/extract.go`**

```go
package proxy

import "github.com/Arsolitt/goxsub/sub"

// ExtractProxies extracts proxy outbounds from subscriptions by protocol and tag.
// Currently supports "vless" protocol and "proxy" tag.
func ExtractProxies(subs []sub.Subscription) []Proxy {
	var proxies []Proxy
	for _, s := range subs {
		for _, ob := range s.Outbounds {
			if ob.Protocol == "vless" || ob.Tag == "proxy" {
				proxies = append(proxies, &VLESSProxy{Outbound: ob, Remarks: s.Remarks})
			}
		}
	}
	return proxies
}
```

- [ ] **Step 3: Create `proxy/extract_test.go`**

Port `TestExtractVLESSOutbounds` and `TestExtractVLESSOutbounds_NoMatches` from `xray/vless_test.go`. Adapt to use `proxy` package types.

```go
package proxy

import (
	"testing"

	"github.com/Arsolitt/goxsub/sub"
)

func TestExtractProxies(t *testing.T) {
	subs := []sub.Subscription{
		{
			Outbounds: []sub.Outbound{
				{Protocol: "vless", Tag: "proxy", Settings: sub.OutboundSettings{Vnext: []sub.VNext{{}}}},
				{Protocol: "socks", Tag: "upstream", Settings: sub.OutboundSettings{}},
				{Protocol: "freedom", Tag: "direct", Settings: sub.OutboundSettings{}},
			},
			Remarks: "Server A",
		},
		{
			Outbounds: []sub.Outbound{
				{Protocol: "vless", Tag: "alt-proxy", Settings: sub.OutboundSettings{Vnext: []sub.VNext{{}}}},
				{Protocol: "blackhole", Tag: "block", Settings: sub.OutboundSettings{}},
			},
			Remarks: "Server B",
		},
	}

	proxies := ExtractProxies(subs)
	if len(proxies) != 2 {
		t.Fatalf("expected 2 proxies, got %d", len(proxies))
	}
	if proxies[0].Remarks() != "Server A" {
		t.Errorf("expected remarks 'Server A', got %q", proxies[0].Remarks())
	}
	if proxies[0].Tag() != "proxy" {
		t.Errorf("expected tag 'proxy', got %q", proxies[0].Tag())
	}
	if proxies[0].Protocol() != "vless" {
		t.Errorf("expected protocol 'vless', got %q", proxies[0].Protocol())
	}
	if proxies[1].Remarks() != "Server B" {
		t.Errorf("expected remarks 'Server B', got %q", proxies[1].Remarks())
	}
	if proxies[1].Tag() != "alt-proxy" {
		t.Errorf("expected tag 'alt-proxy', got %q", proxies[1].Tag())
	}
}

func TestExtractProxies_NoMatches(t *testing.T) {
	subs := []sub.Subscription{
		{
			Outbounds: []sub.Outbound{
				{Protocol: "socks", Tag: "upstream", Settings: sub.OutboundSettings{}},
				{Protocol: "freedom", Tag: "direct", Settings: sub.OutboundSettings{}},
			},
			Remarks: "No VLESS",
		},
	}
	proxies := ExtractProxies(subs)
	if len(proxies) != 0 {
		t.Errorf("expected 0 proxies, got %d", len(proxies))
	}
}
```

- [ ] **Step 4: Create `proxy/filter.go`**

Port `FilterByRemark` from `xray/filter.go`. Change signature from `[]VLESSProxy` to `[]Proxy`.

```go
package proxy

import (
	"path/filepath"
	"strings"
)

// FilterByRemark returns proxies whose Remarks field does not match any of the given glob patterns.
// Matching is case-insensitive. If patterns is empty, all proxies are returned unchanged.
func FilterByRemark(proxies []Proxy, patterns []string) []Proxy {
	if len(patterns) == 0 || len(proxies) == 0 {
		return proxies
	}

	var result []Proxy
	for _, p := range proxies {
		excluded := false
		remark := strings.ToLower(p.Remarks())
		for _, pattern := range patterns {
			matched, _ := filepath.Match(strings.ToLower(pattern), remark)
			if matched {
				excluded = true
				break
			}
		}
		if !excluded {
			result = append(result, p)
		}
	}
	return result
}
```

- [ ] **Step 5: Create `proxy/filter_test.go`**

Port filter tests from `xray/filter_test.go`. Use `proxy.Proxy` interface and `makeProxy` helper.

```go
package proxy

import (
	"testing"

	"github.com/Arsolitt/goxsub/sub"
)

func makeProxy(remarks string) Proxy {
	return &VLESSProxy{
		Remarks: remarks,
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{Vnext: []sub.VNext{{Address: "a.com", Port: 443, Users: []sub.User{{ID: "u"}}}}},
		},
	}
}

func TestFilterByRemark_EmptyPatterns(t *testing.T) {
	proxies := []Proxy{makeProxy("Russia Server"), makeProxy("NL Server")}
	result := FilterByRemark(proxies, nil)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestFilterByRemark_SinglePattern(t *testing.T) {
	proxies := []Proxy{makeProxy("Russia Server"), makeProxy("NL Server")}
	result := FilterByRemark(proxies, []string{"*Russia*"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks() != "NL Server" {
		t.Errorf("expected 'NL Server', got %q", result[0].Remarks())
	}
}

func TestFilterByRemark_MultiplePatterns(t *testing.T) {
	proxies := []Proxy{makeProxy("Russia Server"), makeProxy("China Node"), makeProxy("NL Server")}
	result := FilterByRemark(proxies, []string{"*Russia*", "*China*"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks() != "NL Server" {
		t.Errorf("expected 'NL Server', got %q", result[0].Remarks())
	}
}

func TestFilterByRemark_CaseInsensitive(t *testing.T) {
	proxies := []Proxy{makeProxy("RUSSIA"), makeProxy("russia"), makeProxy("Russia"), makeProxy("NL")}
	result := FilterByRemark(proxies, []string{"*russia*"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks() != "NL" {
		t.Errorf("expected 'NL', got %q", result[0].Remarks())
	}
}

func TestFilterByRemark_GlobSpecials(t *testing.T) {
	proxies := []Proxy{makeProxy("A1"), makeProxy("A2"), makeProxy("B1")}
	result := FilterByRemark(proxies, []string{"A?"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks() != "B1" {
		t.Errorf("expected 'B1', got %q", result[0].Remarks())
	}
}

func TestFilterByRemark_GlobCharClass(t *testing.T) {
	proxies := []Proxy{makeProxy("A1"), makeProxy("B1"), makeProxy("C1")}
	result := FilterByRemark(proxies, []string{"[AB]1"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks() != "C1" {
		t.Errorf("expected 'C1', got %q", result[0].Remarks())
	}
}

func TestFilterByRemark_NoMatches(t *testing.T) {
	proxies := []Proxy{makeProxy("NL Server"), makeProxy("DE Server")}
	result := FilterByRemark(proxies, []string{"*JP*"})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestFilterByRemark_AllExcluded(t *testing.T) {
	proxies := []Proxy{makeProxy("RUSSIA"), makeProxy("Russia")}
	result := FilterByRemark(proxies, []string{"*Russia*"})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}

func TestFilterByRemark_NilInput(t *testing.T) {
	result := FilterByRemark(nil, []string{"*Russia*"})
	if len(result) != 0 {
		t.Fatalf("expected 0, got %d", len(result))
	}
}
```

- [ ] **Step 6: Run tests**

Run: `go test ./proxy/...`
Expected: PASS (2 extract + 9 filter tests)

- [ ] **Step 7: Run lint**

Run: `golangci-lint run --fix`
Expected: no errors

- [ ] **Step 8: Commit**

```bash
git add proxy/ && git commit -m "refactor: add proxy package with interface, VLESSProxy, extraction, and filtering"
```

---

### Task 3: Create `protocol/` package — VLESS URI converter

**Files:**
- Create: `protocol/vless.go`
- Create: `protocol/vless_test.go`

- [ ] **Step 1: Create `protocol/vless.go`**

Port `ToVLESSURI` from `xray/vless.go`. Change signature to accept `*proxy.VLESSProxy` instead of `(Outbound, string)`. Import `github.com/Arsolitt/goxsub/proxy` and `github.com/Arsolitt/goxsub/sub` for types. Keep all the URI-building logic identical — just read fields from `p.Outbound` and `p.Remarks` instead of function parameters.

```go
package protocol

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/Arsolitt/goxsub/proxy"
	"github.com/Arsolitt/goxsub/sub"
)

// VLESSURI converts a VLESSProxy to a vless:// URI string.

//nolint:gocognit,gocyclo,cyclop,funlen // handles multiple transport and security types
func VLESSURI(p *proxy.VLESSProxy) (string, error) {
	o := p.Outbound
	remarks := p.Remarks
	if len(o.Settings.Vnext) == 0 {
		return "", errors.New("no vnext in outbound settings")
	}
	vnext := o.Settings.Vnext[0]
	if vnext.Address == "" {
		return "", errors.New("missing address")
	}
	if vnext.Port == 0 {
		return "", errors.New("missing port")
	}
	if len(vnext.Users) == 0 || vnext.Users[0].ID == "" {
		return "", errors.New("missing user id")
	}

	user := vnext.Users[0]
	ss := o.StreamSettings

	var params []string
	params = append(params, "type="+url.QueryEscape(ss.Network))
	params = append(params, "encryption="+url.QueryEscape(user.Encryption))
	params = append(params, "security="+url.QueryEscape(ss.Security))

	switch ss.Security {
	case "reality":
		if ss.RealitySettings != nil {
			params = append(params, "pbk="+url.QueryEscape(ss.RealitySettings.PublicKey))
			params = append(params, "fp="+url.QueryEscape(ss.RealitySettings.Fingerprint))
			params = append(params, "sni="+url.QueryEscape(ss.RealitySettings.ServerName))
			params = append(params, "sid="+url.QueryEscape(ss.RealitySettings.ShortID))
			if ss.RealitySettings.SpiderX != "" {
				params = append(params, "spx="+url.QueryEscape(ss.RealitySettings.SpiderX))
			}
		}
	case "tls":
		if ss.TLSSettings != nil {
			if ss.TLSSettings.ServerName != "" {
				params = append(params, "sni="+url.QueryEscape(ss.TLSSettings.ServerName))
			}
			if len(ss.TLSSettings.ALPN) > 0 {
				params = append(params, "alpn="+url.QueryEscape(strings.Join(ss.TLSSettings.ALPN, ",")))
			}
			if ss.TLSSettings.Fingerprint != "" {
				params = append(params, "fp="+url.QueryEscape(ss.TLSSettings.Fingerprint))
			}
		}
	}

	switch ss.Network {
	case "ws":
		if ss.WSSettings != nil {
			if ss.WSSettings.Path != "" {
				params = append(params, "path="+url.QueryEscape(ss.WSSettings.Path))
			}
			if ss.WSSettings.Host != "" {
				params = append(params, "host="+url.QueryEscape(ss.WSSettings.Host))
			}
		}
	case "grpc":
		if ss.GRPCSettings != nil && ss.GRPCSettings.ServiceName != "" {
			params = append(params, "serviceName="+url.QueryEscape(ss.GRPCSettings.ServiceName))
		}
	case "h2":
		if ss.HTTPSettings != nil {
			if ss.HTTPSettings.Path != "" {
				params = append(params, "path="+url.QueryEscape(ss.HTTPSettings.Path))
			}
			if len(ss.HTTPSettings.Host) > 0 {
				params = append(params, "host="+url.QueryEscape(strings.Join(ss.HTTPSettings.Host, ",")))
			}
		}
	case "kcp":
		if ss.KCPSettings != nil {
			if ss.KCPSettings.Seed != "" {
				params = append(params, "type="+url.QueryEscape(ss.KCPSettings.Seed))
			}
			if ss.KCPSettings.HeaderType != nil && ss.KCPSettings.HeaderType.Type != "" {
				params = append(params, "headerType="+url.QueryEscape(ss.KCPSettings.HeaderType.Type))
			}
		}
	}

	if user.Flow != "" {
		params = append(params, "flow="+url.QueryEscape(user.Flow))
	}

	fragment := ""
	if remarks != "" {
		fragment = "#" + url.PathEscape(remarks)
	}

	return fmt.Sprintf("vless://%s@%s?%s%s",
		user.ID, net.JoinHostPort(vnext.Address, strconv.Itoa(vnext.Port)),
		strings.Join(params, "&"), fragment), nil
}

// ToVLESSURI converts a VLESSProxy (as a Proxy interface) to a vless:// URI string.
// Returns an error if the proxy is not a VLESSProxy or conversion fails.
func ToVLESSURI(p proxy.Proxy) (string, error) {
	vp, ok := p.(*proxy.VLESSProxy)
	if !ok {
		return "", fmt.Errorf("expected VLESSProxy, got %T", p)
	}
	return VLESSURI(vp)
}

// ToURI converts a proxy to its protocol URI string.
// Currently supports vless. Returns an error for unsupported protocols.
func ToURI(p proxy.Proxy) (string, error) {
	switch p.Protocol() {
	case "vless":
		return ToVLESSURI(p)
	default:
		return "", fmt.Errorf("unsupported protocol: %s", p.Protocol())
	}
}
```

- [ ] **Step 2: Create `protocol/vless_test.go`**

Port all `TestToVLESSURI_*` tests from `xray/vless_test.go`. Each test constructs a `*proxy.VLESSProxy` and calls `VLESSURI()`. Expected values are identical.

```go
package protocol

import (
	"testing"

	"github.com/Arsolitt/goxsub/proxy"
	"github.com/Arsolitt/goxsub/sub"
)

func TestVLESSURI_TCPReality(t *testing.T) {
	p := &proxy.VLESSProxy{
		Remarks: "Test TCP Reality",
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{
					Address: "example.com",
					Port:    443,
					Users:   []sub.User{{ID: "test-uuid-1234", Encryption: "none", Flow: "xtls-rprx-vision"}},
				}},
			},
			StreamSettings: sub.StreamSettings{
				Network:  "tcp",
				Security: "reality",
				RealitySettings: &sub.RealitySettings{
					PublicKey:   "pub_key_value",
					Fingerprint: "firefox",
					ServerName:  "sni.example.com",
					ShortID:     "abcd1234",
					SpiderX:     "/",
				},
			},
		},
	}
	uri, err := VLESSURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://test-uuid-1234@example.com:443?type=tcp&encryption=none&security=reality&pbk=pub_key_value&fp=firefox&sni=sni.example.com&sid=abcd1234&spx=%2F&flow=xtls-rprx-vision#Test%20TCP%20Reality"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_TCPTLS(t *testing.T) {
	p := &proxy.VLESSProxy{
		Remarks: "Test TCP TLS",
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{
					Address: "tls.example.com",
					Port:    8443,
					Users:   []sub.User{{ID: "uuid-tls", Encryption: "none"}},
				}},
			},
			StreamSettings: sub.StreamSettings{
				Network:  "tcp",
				Security: "tls",
				TLSSettings: &sub.TLSSettings{
					ServerName:  "tls.example.com",
					ALPN:        []string{"h2", "http/1.1"},
					Fingerprint: "chrome",
				},
			},
		},
	}
	uri, err := VLESSURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-tls@tls.example.com:8443?type=tcp&encryption=none&security=tls&sni=tls.example.com&alpn=h2%2Chttp%2F1.1&fp=chrome#Test%20TCP%20TLS"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_WSTLS(t *testing.T) {
	p := &proxy.VLESSProxy{
		Remarks: "Test WS TLS",
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{
					Address: "ws.example.com",
					Port:    443,
					Users:   []sub.User{{ID: "uuid-ws", Encryption: "none"}},
				}},
			},
			StreamSettings: sub.StreamSettings{
				Network:  "ws",
				Security: "tls",
				TLSSettings: &sub.TLSSettings{
					ServerName: "ws.example.com",
				},
				WSSettings: &sub.WSSettings{
					Path: "/websocket-path",
					Host: "ws-host.example.com",
				},
			},
		},
	}
	uri, err := VLESSURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-ws@ws.example.com:443?type=ws&encryption=none&security=tls&sni=ws.example.com&path=%2Fwebsocket-path&host=ws-host.example.com#Test%20WS%20TLS"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_GrpctLS(t *testing.T) {
	p := &proxy.VLESSProxy{
		Remarks: "Test GRPC TLS",
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{
					Address: "grpc.example.com",
					Port:    443,
					Users:   []sub.User{{ID: "uuid-grpc", Encryption: "none"}},
				}},
			},
			StreamSettings: sub.StreamSettings{
				Network:  "grpc",
				Security: "tls",
				TLSSettings: &sub.TLSSettings{
					ServerName: "grpc.example.com",
				},
				GRPCSettings: &sub.GRPCSettings{
					ServiceName: "grpc-service",
				},
			},
		},
	}
	uri, err := VLESSURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-grpc@grpc.example.com:443?type=grpc&encryption=none&security=tls&sni=grpc.example.com&serviceName=grpc-service#Test%20GRPC%20TLS"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_H2TLS(t *testing.T) {
	p := &proxy.VLESSProxy{
		Remarks: "Test H2 TLS",
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{
					Address: "h2.example.com",
					Port:    443,
					Users:   []sub.User{{ID: "uuid-h2", Encryption: "none"}},
				}},
			},
			StreamSettings: sub.StreamSettings{
				Network:  "h2",
				Security: "tls",
				TLSSettings: &sub.TLSSettings{
					ServerName: "h2.example.com",
				},
				HTTPSettings: &sub.HTTPSettings{
					Path: "/h2-path",
					Host: []string{"h2a.example.com", "h2b.example.com"},
				},
			},
		},
	}
	uri, err := VLESSURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-h2@h2.example.com:443?type=h2&encryption=none&security=tls&sni=h2.example.com&path=%2Fh2-path&host=h2a.example.com%2Ch2b.example.com#Test%20H2%20TLS"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_KCPNone(t *testing.T) {
	p := &proxy.VLESSProxy{
		Remarks: "Test KCP None",
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{
					Address: "kcp.example.com",
					Port:    443,
					Users:   []sub.User{{ID: "uuid-kcp", Encryption: "none"}},
				}},
			},
			StreamSettings: sub.StreamSettings{
				Network:  "kcp",
				Security: "none",
				KCPSettings: &sub.KCPSettings{
					Seed:       "mykcpseed",
					HeaderType: &sub.KCPHeader{Type: "none"},
				},
			},
		},
	}
	uri, err := VLESSURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-kcp@kcp.example.com:443?type=kcp&encryption=none&security=none&type=mykcpseed&headerType=none#Test%20KCP%20None"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_TCPNone(t *testing.T) {
	p := &proxy.VLESSProxy{
		Remarks: "Test Plain",
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{
					Address: "plain.example.com",
					Port:    80,
					Users:   []sub.User{{ID: "uuid-plain", Encryption: "none"}},
				}},
			},
			StreamSettings: sub.StreamSettings{
				Network:  "tcp",
				Security: "none",
			},
		},
	}
	uri, err := VLESSURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-plain@plain.example.com:80?type=tcp&encryption=none&security=none#Test%20Plain"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_NoVnext(t *testing.T) {
	p := &proxy.VLESSProxy{
		Outbound: sub.Outbound{Protocol: "vless", Settings: sub.OutboundSettings{}},
	}
	_, err := VLESSURI(p)
	if err == nil {
		t.Fatal("expected error for no vnext")
	}
}

func TestVLESSURI_MissingAddress(t *testing.T) {
	p := &proxy.VLESSProxy{
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{Port: 443, Users: []sub.User{{ID: "uuid"}}}},
			},
		},
	}
	_, err := VLESSURI(p)
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestVLESSURI_MissingPort(t *testing.T) {
	p := &proxy.VLESSProxy{
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{Address: "example.com", Users: []sub.User{{ID: "uuid"}}}},
			},
		},
	}
	_, err := VLESSURI(p)
	if err == nil {
		t.Fatal("expected error for missing port")
	}
}

func TestVLESSURI_MissingUserID(t *testing.T) {
	p := &proxy.VLESSProxy{
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{Address: "example.com", Port: 443, Users: []sub.User{}}},
			},
		},
	}
	_, err := VLESSURI(p)
	if err == nil {
		t.Fatal("expected error for missing user id")
	}
}

func TestVLESSURI_NoRemarks(t *testing.T) {
	p := &proxy.VLESSProxy{
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{
					Address: "example.com",
					Port:    443,
					Users:   []sub.User{{ID: "uuid", Encryption: "none"}},
				}},
			},
			StreamSettings: sub.StreamSettings{
				Network:  "tcp",
				Security: "none",
			},
		},
	}
	uri, err := VLESSURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid@example.com:443?type=tcp&encryption=none&security=none"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_RealityNoSpiderX(t *testing.T) {
	p := &proxy.VLESSProxy{
		Remarks: "NoSpiderX",
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{
					Address: "example.com",
					Port:    443,
					Users:   []sub.User{{ID: "uuid", Encryption: "none"}},
				}},
			},
			StreamSettings: sub.StreamSettings{
				Network:  "tcp",
				Security: "reality",
				RealitySettings: &sub.RealitySettings{
					PublicKey:   "pubkey",
					Fingerprint: "chrome",
					ServerName:  "sni.com",
					ShortID:     "1234",
				},
			},
		},
	}
	uri, err := VLESSURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid@example.com:443?type=tcp&encryption=none&security=reality&pbk=pubkey&fp=chrome&sni=sni.com&sid=1234#NoSpiderX"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_RealityEncodedPublicKey(t *testing.T) {
	p := &proxy.VLESSProxy{
		Remarks: "Encoded PBK",
		Outbound: sub.Outbound{
			Protocol: "vless",
			Settings: sub.OutboundSettings{
				Vnext: []sub.VNext{{
					Address: "example.com",
					Port:    443,
					Users:   []sub.User{{ID: "uuid", Encryption: "none"}},
				}},
			},
			StreamSettings: sub.StreamSettings{
				Network:  "tcp",
				Security: "reality",
				RealitySettings: &sub.RealitySettings{
					PublicKey:   "abc123+def/ghi=",
					Fingerprint: "chrome",
					ServerName:  "sni.com",
					ShortID:     "abcd",
				},
			},
		},
	}
	uri, err := VLESSURI(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid@example.com:443?type=tcp&encryption=none&security=reality&pbk=abc123%2Bdef%2Fghi%3D&fp=chrome&sni=sni.com&sid=abcd#Encoded%20PBK"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestToURI_Unsupported(t *testing.T) {
	p := &mockProxy{proto: "trojan"}
	_, err := ToURI(p)
	if err == nil {
		t.Fatal("expected error for unsupported protocol")
	}
}

type mockProxy struct {
	proto string
}

func (m *mockProxy) Protocol() string { return m.proto }
func (m *mockProxy) Tag() string      { return "" }
func (m *mockProxy) Remarks() string   { return "" }
```

- [ ] **Step 3: Run tests**

Run: `go test ./protocol/...`
Expected: PASS (13 VLESSURI tests + 1 ToURI test)

- [ ] **Step 4: Run lint**

Run: `golangci-lint run --fix`
Expected: no errors

- [ ] **Step 5: Commit**

```bash
git add protocol/ && git commit -m "refactor: add protocol package with VLESS URI converter"
```

---

### Task 4: Create `format/` package — Podkop formatter

**Files:**
- Create: `format/podkop.go`
- Create: `format/podkop_test.go`

- [ ] **Step 1: Create `format/podkop.go`**

Port `FormatPodkop` from `xray/vless.go`. Change signature to accept `[]proxy.Proxy`. Switch on `p.Protocol()` and type-assert to `*proxy.VLESSProxy`.

```go
package format

import (
	"fmt"
	"strings"

	"github.com/Arsolitt/goxsub/protocol"
	"github.com/Arsolitt/goxsub/proxy"
)

// Podkop converts proxies to uci shell commands for podkop OpenWrt package.
func Podkop(proxies []proxy.Proxy, section string) (string, error) {
	var lines []string
	lines = append(lines, fmt.Sprintf("uci del podkop.%s.urltest_proxy_links", section))
	for _, p := range proxies {
		uri, err := protocol.ToURI(p)
		if err != nil {
			return "", fmt.Errorf("convert proxy to URI: %w", err)
		}
		lines = append(lines, fmt.Sprintf("uci add_list podkop.%s.urltest_proxy_links='%s'", section, uri))
	}
	lines = append(lines, "service podkop restart")
	return strings.Join(lines, "\n"), nil
}
```

- [ ] **Step 2: Create `format/podkop_test.go`**

Port `TestFormatPodkop*` tests from `xray/vless_test.go`. Use `proxy.Proxy` interface.

```go
package format

import (
	"strings"
	"testing"

	"github.com/Arsolitt/goxsub/proxy"
	"github.com/Arsolitt/goxsub/sub"
)

func TestPodkop(t *testing.T) {
	proxies := []proxy.Proxy{
		&proxy.VLESSProxy{
			Remarks: "Server A",
			Outbound: sub.Outbound{
				Protocol: "vless",
				Settings: sub.OutboundSettings{
					Vnext: []sub.VNext{{
						Address: "a.example.com",
						Port:    443,
						Users:   []sub.User{{ID: "uuid-a", Encryption: "none"}},
					}},
				},
				StreamSettings: sub.StreamSettings{
					Network:  "tcp",
					Security: "none",
				},
			},
		},
		&proxy.VLESSProxy{
			Remarks: "Server B",
			Outbound: sub.Outbound{
				Protocol: "vless",
				Settings: sub.OutboundSettings{
					Vnext: []sub.VNext{{
						Address: "b.example.com",
						Port:    8443,
						Users:   []sub.User{{ID: "uuid-b", Encryption: "none"}},
					}},
				},
				StreamSettings: sub.StreamSettings{
					Network:  "tcp",
					Security: "none",
				},
			},
		},
	}

	result, err := Podkop(proxies, "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(result, "\n")
	if len(lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(lines))
	}
	if lines[0] != "uci del podkop.main.urltest_proxy_links" {
		t.Errorf("first line mismatch:\ngot:      %s\nexpected: uci del podkop.main.urltest_proxy_links", lines[0])
	}
	if !strings.HasPrefix(lines[1], "uci add_list podkop.main.urltest_proxy_links='vless://uuid-a@") {
		t.Errorf("second line mismatch:\ngot: %s", lines[1])
	}
	if !strings.HasPrefix(lines[2], "uci add_list podkop.main.urltest_proxy_links='vless://uuid-b@") {
		t.Errorf("third line mismatch:\ngot: %s", lines[2])
	}
	if lines[3] != "service podkop restart" {
		t.Errorf("last line mismatch:\ngot:      %s\nexpected: service podkop restart", lines[3])
	}
}

func TestPodkop_CustomSection(t *testing.T) {
	proxies := []proxy.Proxy{
		&proxy.VLESSProxy{
			Remarks: "S",
			Outbound: sub.Outbound{
				Protocol: "vless",
				Settings: sub.OutboundSettings{
					Vnext: []sub.VNext{{
						Address: "c.example.com",
						Port:    443,
						Users:   []sub.User{{ID: "uuid-c", Encryption: "none"}},
					}},
				},
				StreamSettings: sub.StreamSettings{
					Network:  "tcp",
					Security: "none",
				},
			},
		},
	}

	result, err := Podkop(proxies, "backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "uci del podkop.backup.urltest_proxy_links\nuci add_list podkop.backup.urltest_proxy_links='vless://uuid-c@c.example.com:443?type=tcp&encryption=none&security=none#S'\nservice podkop restart"
	if result != expected {
		t.Errorf("result mismatch:\ngot:      %s\nexpected: %s", result, expected)
	}
}

func TestPodkop_Empty(t *testing.T) {
	result, err := Podkop(nil, "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "uci del podkop.main.urltest_proxy_links\nservice podkop restart"
	if result != expected {
		t.Errorf("result mismatch:\ngot:      %s\nexpected: %s", result, expected)
	}
}

func TestPodkop_InvalidProxy(t *testing.T) {
	proxies := []proxy.Proxy{
		&proxy.VLESSProxy{
			Remarks:  "bad",
			Outbound: sub.Outbound{Protocol: "vless", Settings: sub.OutboundSettings{}},
		},
	}
	_, err := Podkop(proxies, "main")
	if err == nil {
		t.Fatal("expected error for invalid proxy")
	}
}
```

- [ ] **Step 3: Run tests**

Run: `go test ./format/...`
Expected: PASS (4 tests)

- [ ] **Step 4: Run lint**

Run: `golangci-lint run --fix`
Expected: no errors

- [ ] **Step 5: Commit**

```bash
git add format/ && git commit -m "refactor: add format package with Podkop formatter"
```

---

### Task 5: Create root `api.go` — public re-exports

**Files:**
- Create: `api.go`

- [ ] **Step 1: Create `api.go`**

Re-export all public types and functions from internal packages. This is the file users interact with.

```go
package goxsub

import (
	"github.com/Arsolitt/goxsub/format"
	"github.com/Arsolitt/goxsub/proxy"
	"github.com/Arsolitt/goxsub/protocol"
	"github.com/Arsolitt/goxsub/sub"
)

type Subscription = sub.Subscription
type Outbound = sub.Outbound
type OutboundSettings = sub.OutboundSettings
type VNext = sub.VNext
type User = sub.User
type StreamSettings = sub.StreamSettings
type RealitySettings = sub.RealitySettings
type TLSSettings = sub.TLSSettings
type TCPSettings = sub.TCPSettings
type TCPHeader = sub.TCPHeader
type WSSettings = sub.WSSettings
type GRPCSettings = sub.GRPCSettings
type HTTPSettings = sub.HTTPSettings
type KCPSettings = sub.KCPSettings
type KCPHeader = sub.KCPHeader

type Proxy = proxy.Proxy
type VLESSProxy = proxy.VLESSProxy

var ParseSubscription = sub.ParseSubscription
var ExtractProxies = proxy.ExtractProxies
var FilterByRemark = proxy.FilterByRemark
var ToURI = protocol.ToURI
var ToVLESSURI = protocol.ToVLESSURI
var VLESSURI = protocol.VLESSURI
var Podkop = format.Podkop
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./...`
Expected: no errors

- [ ] **Step 3: Run lint**

Run: `golangci-lint run --fix`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add api.go && git commit -m "refactor: add root api.go with public re-exports"
```

---

### Task 6: Update CLI to use root import

**Files:**
- Modify: `cmd/goxsub/main.go`

- [ ] **Step 1: Update `cmd/goxsub/main.go`**

Replace `"github.com/Arsolitt/goxsub/xray"` with `"github.com/Arsolitt/goxsub"`. Replace all `xray.` calls with `goxsub.`. Replace `xray.ExtractVLESSOutbounds(subs)` with `goxsub.ExtractProxies(subs)`. Replace `xray.ToVLESSURI(p.Outbound, p.Remarks)` with `goxsub.ToURI(p)`. Replace `xray.FormatPodkop(proxies, *podkopSection)` with `goxsub.Podkop(proxies, *podkopSection)`.

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	goxsub "github.com/Arsolitt/goxsub"
)

type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ", ") }

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	os.Exit(run())
}

//nolint:funlen
func run() int {
	formatFlag := flag.String("format", "uri", "output format: uri, podkop")
	podkopSection := flag.String("podkop-section", "main", "podkop uci section name")
	var excludePatterns stringSlice
	flag.Var(
		&excludePatterns,
		"exclude-by-remark",
		"exclude outbounds by remark glob pattern (case-insensitive, repeatable)",
	)
	flag.Parse()

	if *formatFlag != "podkop" && flag.Lookup("podkop-section").DefValue != *podkopSection {
		fmt.Fprintf(os.Stderr, "error: --podkop-section can only be used with --format podkop\n")
		return 1
	}

	for _, p := range excludePatterns {
		if _, err := filepath.Match(p, ""); err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid glob pattern %q: %v\n", p, err)
			return 1
		}
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: goxsub [flags] <subscription-url>\n")
		fmt.Fprintf(os.Stderr, "flags:\n")
		flag.PrintDefaults()
		return 1
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		args[0],
		nil,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "error: HTTP %d\n", resp.StatusCode)
		return 1
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	subs, err := goxsub.ParseSubscription(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	proxies := goxsub.ExtractProxies(subs)
	proxies = goxsub.FilterByRemark(proxies, excludePatterns)

	switch *formatFlag {
	case "podkop":
		output, err := goxsub.Podkop(proxies, *podkopSection)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		fmt.Println(output)
	default:
		for _, p := range proxies {
			uri, err := goxsub.ToURI(p)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return 1
			}
			fmt.Println(uri)
		}
	}

	return 0
}
```

- [ ] **Step 2: Build CLI**

Run: `go build -o build/goxsub ./cmd/goxsub/`
Expected: no errors

- [ ] **Step 3: Run lint**

Run: `golangci-lint run --fix`
Expected: no errors

- [ ] **Step 4: Commit**

```bash
git add cmd/goxsub/main.go && git commit -m "refactor: update CLI to use root goxsub import"
```

---

### Task 7: Delete `xray/` package and run full verification

**Files:**
- Delete: `xray/` (all files)

- [ ] **Step 1: Delete xray package**

```bash
rm -rf xray/
```

- [ ] **Step 2: Run all tests**

Run: `go test ./...`
Expected: PASS (all tests across sub/, proxy/, protocol/, format/)

- [ ] **Step 3: Run lint**

Run: `golangci-lint run --fix`
Expected: no errors

- [ ] **Step 4: Build**

Run: `go build -o build/goxsub ./cmd/goxsub/`
Expected: no errors

- [ ] **Step 5: Commit**

```bash
git add -A && git commit -m "refactor: remove xray package, complete API redesign"
```
