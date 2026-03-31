# Sing-box Output Format Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add sing-box outbound JSON output format to goxsub CLI and library, refactor formatters to return raw string slices.

**Architecture:** New `format/singbox.go` with `SingboxConfig` struct and `Singbox` function that converts proxies to sing-box JSON objects. CLI layer handles output presentation (trailing commas for sing-box). `Podkop` refactored to return `[]string` instead of single joined string. `VLESSProxy` gains `SetRemarks` method for remark override.

**Tech Stack:** Go 1.26, standard library only (`encoding/json`).

---

### Task 1: Add SetRemarks to VLESSProxy

**Files:**
- Modify: `proxy/proxy.go:20`
- Test: `proxy/proxy_test.go` (create)

- [ ] **Step 1: Write the failing test**

Create `proxy/proxy_test.go`:

```go
package proxy

import (
	"testing"

	"github.com/Arsolitt/goxsub/sub"
)

func TestSetRemarks(t *testing.T) {
	p := &VLESSProxy{
		Outbound: sub.Outbound{Protocol: "vless"},
		remark:   "original",
	}
	if p.Remarks() != "original" {
		t.Fatalf("expected original remark, got %q", p.Remarks())
	}
	p.SetRemarks("changed")
	if p.Remarks() != "changed" {
		t.Fatalf("expected changed remark, got %q", p.Remarks())
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./proxy/ -run TestSetRemarks -v`
Expected: FAIL — `p.SetRemarks` undefined

- [ ] **Step 3: Add SetRemarks method to VLESSProxy**

Add to `proxy/proxy.go` after line 20:

```go
func (v *VLESSProxy) SetRemarks(r string) { v.remark = r }
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./proxy/ -run TestSetRemarks -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add proxy/proxy.go proxy/proxy_test.go
git commit -m "feat: add SetRemarks method to VLESSProxy"
```

---

### Task 2: Refactor Podkop to return []string

**Files:**
- Modify: `format/podkop.go`
- Modify: `format/podkop_test.go`
- Modify: `api.go`

- [ ] **Step 1: Write the failing tests**

Replace all of `format/podkop_test.go` with updated tests that expect `[]string`:

```go
package format

import (
	"strings"
	"testing"

	"github.com/Arsolitt/goxsub/proxy"
	"github.com/Arsolitt/goxsub/sub"
)

func TestPodkop(t *testing.T) {
	subs := []sub.Subscription{
		{
			Outbounds: []sub.Outbound{
				{
					Protocol: "vless",
					Tag:      "proxy",
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
				{
					Protocol: "vless",
					Tag:      "proxy2",
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
			Remarks: "Server A",
		},
	}
	proxies := proxy.ExtractProxies(subs)

	result, err := Podkop(proxies, "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(result))
	}
	if result[0] != "uci del podkop.main.urltest_proxy_links" {
		t.Errorf("first line mismatch:\ngot:      %s\nexpected: uci del podkop.main.urltest_proxy_links", result[0])
	}
	if !strings.HasPrefix(result[1], "uci add_list podkop.main.urltest_proxy_links='vless://uuid-a@") {
		t.Errorf("second line mismatch:\ngot: %s", result[1])
	}
	if !strings.HasPrefix(result[2], "uci add_list podkop.main.urltest_proxy_links='vless://uuid-b@") {
		t.Errorf("third line mismatch:\ngot: %s", result[2])
	}
	if result[3] != "service podkop restart" {
		t.Errorf("last line mismatch:\ngot:      %s\nexpected: service podkop restart", result[3])
	}
}

func TestPodkop_CustomSection(t *testing.T) {
	subs := []sub.Subscription{
		{
			Outbounds: []sub.Outbound{
				{
					Protocol: "vless",
					Tag:      "proxy",
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
			Remarks: "S",
		},
	}
	proxies := proxy.ExtractProxies(subs)

	result, err := Podkop(proxies, "backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(result))
	}
	expected := "uci add_list podkop.backup.urltest_proxy_links='vless://uuid-c@c.example.com:443?type=tcp&encryption=none&security=none#S'"
	if result[1] != expected {
		t.Errorf("second line mismatch:\ngot:      %s\nexpected: %s", result[1], expected)
	}
}

func TestPodkop_Empty(t *testing.T) {
	result, err := Podkop(nil, "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(result))
	}
	if result[0] != "uci del podkop.main.urltest_proxy_links" {
		t.Errorf("first line mismatch:\ngot:      %s\nexpected: uci del podkop.main.urltest_proxy_links", result[0])
	}
	if result[1] != "service podkop restart" {
		t.Errorf("last line mismatch:\ngot:      %s\nexpected: service podkop restart", result[1])
	}
}

func TestPodkop_InvalidProxy(t *testing.T) {
	subs := []sub.Subscription{
		{
			Outbounds: []sub.Outbound{
				{Protocol: "vless", Tag: "proxy", Settings: sub.OutboundSettings{}},
			},
			Remarks: "bad",
		},
	}
	proxies := proxy.ExtractProxies(subs)
	_, err := Podkop(proxies, "main")
	if err == nil {
		t.Fatal("expected error for invalid proxy")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./format/ -v`
Expected: FAIL — cannot use `result` (type `string`) as type `[]string`

- [ ] **Step 3: Refactor Podkop to return []string**

Replace all of `format/podkop.go` with:

```go
package format

import (
	"fmt"

	"github.com/Arsolitt/goxsub/protocol"
	"github.com/Arsolitt/goxsub/proxy"
)

// Podkop converts proxies to uci shell commands for podkop OpenWrt package.
// Returns a slice of command lines.
func Podkop(proxies []proxy.Proxy, section string) ([]string, error) {
	lines := []string{
		fmt.Sprintf("uci del podkop.%s.urltest_proxy_links", section),
	}
	for _, p := range proxies {
		uri, err := protocol.ToURI(p)
		if err != nil {
			return nil, fmt.Errorf("convert proxy to URI: %w", err)
		}
		lines = append(lines, fmt.Sprintf("uci add_list podkop.%s.urltest_proxy_links='%s'", section, uri))
	}
	lines = append(lines, "service podkop restart")
	return lines, nil
}
```

- [ ] **Step 4: Update api.go Podkop re-export**

No change needed — `var Podkop = format.Podkop` still works since the variable type is inferred.

- [ ] **Step 5: Update CLI to join Podkop lines**

In `cmd/goxsub/main.go`, replace the podkop case (lines 100-106):

Old:
```go
	case "podkop":
		output, err := goxsub.Podkop(proxies, *podkopSection)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		fmt.Println(output)
```

New:
```go
	case "podkop":
		lines, err := goxsub.Podkop(proxies, *podkopSection)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		for _, line := range lines {
			fmt.Println(line)
		}
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `go test ./... -v`
Expected: ALL PASS

- [ ] **Step 7: Lint**

Run: `golangci-lint run --fix`
Expected: No errors

- [ ] **Step 8: Commit**

```bash
git add format/podkop.go format/podkop_test.go cmd/goxsub/main.go
git commit -m "feat: change Podkop return type from string to []string"
```

---

### Task 3: Add Singbox formatter

**Files:**
- Create: `format/singbox.go`
- Create: `format/singbox_test.go`

- [ ] **Step 1: Write the failing tests**

Create `format/singbox_test.go`:

```go
package format

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Arsolitt/goxsub/proxy"
	"github.com/Arsolitt/goxsub/sub"
)

func TestSingbox_Reality(t *testing.T) {
	subs := []sub.Subscription{
		{
			Outbounds: []sub.Outbound{
				{
					Protocol: "vless",
					Tag:      "proxy",
					Settings: sub.OutboundSettings{
						Vnext: []sub.VNext{{
							Address: "lt-cherry-01.com",
							Port:    8443,
							Users: []sub.User{{
								ID:         "18a77cca-d7d8-41a3-bf31-8144854623b5",
								Encryption: "none",
								Flow:       "xtls-rprx-vision",
							}},
						}},
					},
					StreamSettings: sub.StreamSettings{
						Network:  "tcp",
						Security: "reality",
						RealitySettings: &sub.RealitySettings{
							PublicKey:   "L3X1eh1Jq_6PKJ6LlwjgiWq0XNaDOqCVKgIElJ5nkVA",
							Fingerprint: "chrome",
							ServerName:  "rbc.ru",
							ShortID:     "e0ef3d5c0aacb615",
						},
					},
				},
			},
			Remarks: "lt-1",
		},
	}
	proxies := proxy.ExtractProxies(subs)
	cfg := SingboxConfig{
		OutboundPrefix: "proxy-",
		KeepRemark:     true,
		DNSResolver:    "dns-local",
	}

	result, err := Singbox(proxies, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 line, got %d", len(result))
	}

	var out singboxOutbound
	if err := json.Unmarshal([]byte(result[0]), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if out.Type != "vless" {
		t.Errorf("type mismatch: got %q, want %q", out.Type, "vless")
	}
	if out.Tag != "proxy-lt-1" {
		t.Errorf("tag mismatch: got %q, want %q", out.Tag, "proxy-lt-1")
	}
	if out.Server != "lt-cherry-01.com" {
		t.Errorf("server mismatch: got %q, want %q", out.Server, "lt-cherry-01.com")
	}
	if out.ServerPort != 8443 {
		t.Errorf("server_port mismatch: got %d, want %d", out.ServerPort, 8443)
	}
	if out.UUID != "18a77cca-d7d8-41a3-bf31-8144854623b5" {
		t.Errorf("uuid mismatch: got %q", out.UUID)
	}
	if out.Flow != "xtls-rprx-vision" {
		t.Errorf("flow mismatch: got %q", out.Flow)
	}
	if out.TLS == nil {
		t.Fatal("expected tls block")
	}
	if !out.TLS.Enabled {
		t.Error("expected tls enabled")
	}
	if out.TLS.ServerName != "rbc.ru" {
		t.Errorf("tls server_name mismatch: got %q", out.TLS.ServerName)
	}
	if out.TLS.UTLS == nil || !out.TLS.UTLS.Enabled || out.TLS.UTLS.Fingerprint != "chrome" {
		t.Errorf("tls utls mismatch: %+v", out.TLS.UTLS)
	}
	if out.TLS.Reality == nil || !out.TLS.Reality.Enabled {
		t.Error("expected reality enabled")
	}
	if out.TLS.Reality.PublicKey != "L3X1eh1Jq_6PKJ6LlwjgiWq0XNaDOqCVKgIElJ5nkVA" {
		t.Errorf("reality public_key mismatch: got %q", out.TLS.Reality.PublicKey)
	}
	if out.TLS.Reality.ShortID != "e0ef3d5c0aacb615" {
		t.Errorf("reality short_id mismatch: got %q", out.TLS.Reality.ShortID)
	}
	if out.DomainResolver != "dns-local" {
		t.Errorf("domain_resolver mismatch: got %q", out.DomainResolver)
	}
}

func TestSingbox_NoRemark(t *testing.T) {
	subs := []sub.Subscription{
		{
			Outbounds: []sub.Outbound{
				{
					Protocol: "vless",
					Tag:      "proxy",
					Settings: sub.OutboundSettings{
						Vnext: []sub.VNext{{
							Address: "a.com",
							Port:    443,
							Users:   []sub.User{{ID: "uuid-1", Encryption: "none"}},
						}},
					},
					StreamSettings: sub.StreamSettings{
						Network:  "tcp",
						Security: "none",
					},
				},
				{
					Protocol: "vless",
					Tag:      "proxy2",
					Settings: sub.OutboundSettings{
						Vnext: []sub.VNext{{
							Address: "b.com",
							Port:    8443,
							Users:   []sub.User{{ID: "uuid-2", Encryption: "none"}},
						}},
					},
					StreamSettings: sub.StreamSettings{
						Network:  "tcp",
						Security: "none",
					},
				},
			},
			Remarks: "Server A",
		},
	}
	proxies := proxy.ExtractProxies(subs)
	cfg := SingboxConfig{
		OutboundPrefix: "p-",
		OutboundSuffix: "-x",
		KeepRemark:     false,
		DNSResolver:    "dns-local",
	}

	result, err := Singbox(proxies, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(result))
	}

	var out1, out2 singboxOutbound
	if err := json.Unmarshal([]byte(result[0]), &out1); err != nil {
		t.Fatalf("invalid JSON line 0: %v", err)
	}
	if err := json.Unmarshal([]byte(result[1]), &out2); err != nil {
		t.Fatalf("invalid JSON line 1: %v", err)
	}

	if out1.Tag != "p-1-x" {
		t.Errorf("tag mismatch line 0: got %q, want %q", out1.Tag, "p-1-x")
	}
	if out2.Tag != "p-2-x" {
		t.Errorf("tag mismatch line 1: got %q, want %q", out2.Tag, "p-2-x")
	}
}

func TestSingbox_Empty(t *testing.T) {
	cfg := SingboxConfig{DNSResolver: "dns-local"}
	result, err := Singbox(nil, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(result))
	}
}

func TestSingbox_PrefixSuffix(t *testing.T) {
	subs := []sub.Subscription{
		{
			Outbounds: []sub.Outbound{
				{
					Protocol: "vless",
					Tag:      "proxy",
					Settings: sub.OutboundSettings{
						Vnext: []sub.VNext{{
							Address: "a.com",
							Port:    443,
							Users:   []sub.User{{ID: "uuid", Encryption: "none"}},
						}},
					},
					StreamSettings: sub.StreamSettings{
						Network:  "tcp",
						Security: "none",
					},
				},
			},
			Remarks: "myserver",
		},
	}
	proxies := proxy.ExtractProxies(subs)
	cfg := SingboxConfig{
		OutboundPrefix: "[",
		OutboundSuffix: "]",
		KeepRemark:     true,
		DNSResolver:    "dns-local",
	}

	result, err := Singbox(proxies, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result[0], `"tag":"[myserver]"`) {
		t.Errorf("expected tag with prefix/suffix, got:\n%s", result[0])
	}
}

func TestSingbox_NoTLS(t *testing.T) {
	subs := []sub.Subscription{
		{
			Outbounds: []sub.Outbound{
				{
					Protocol: "vless",
					Tag:      "proxy",
					Settings: sub.OutboundSettings{
						Vnext: []sub.VNext{{
							Address: "a.com",
							Port:    443,
							Users:   []sub.User{{ID: "uuid", Encryption: "none"}},
						}},
					},
					StreamSettings: sub.StreamSettings{
						Network:  "tcp",
						Security: "none",
					},
				},
			},
			Remarks: "plain",
		},
	}
	proxies := proxy.ExtractProxies(subs)
	cfg := SingboxConfig{KeepRemark: true, DNSResolver: "dns-local"}

	result, err := Singbox(proxies, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(result[0]), &raw); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := raw["tls"]; ok {
		t.Error("expected no tls field when security is none")
	}
}

func TestSingbox_InvalidProxy(t *testing.T) {
	subs := []sub.Subscription{
		{
			Outbounds: []sub.Outbound{
				{Protocol: "vless", Tag: "proxy", Settings: sub.OutboundSettings{}},
			},
			Remarks: "bad",
		},
	}
	proxies := proxy.ExtractProxies(subs)
	cfg := SingboxConfig{KeepRemark: true, DNSResolver: "dns-local"}
	_, err := Singbox(proxies, cfg)
	if err == nil {
		t.Fatal("expected error for invalid proxy")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./format/ -run TestSingbox -v`
Expected: FAIL — `Singbox` undefined

- [ ] **Step 3: Implement Singbox formatter**

Create `format/singbox.go`:

```go
package format

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Arsolitt/goxsub/proxy"
)

// SingboxConfig holds options for sing-box outbound generation.
type SingboxConfig struct {
	OutboundPrefix string
	OutboundSuffix string
	KeepRemark     bool
	DNSResolver    string
}

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
	Enabled    bool            `json:"enabled"`
	ServerName string          `json:"server_name,omitempty"`
	UTLS       *singboxUTLS    `json:"utls,omitempty"`
	Reality    *singboxReality `json:"reality,omitempty"`
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

// Singbox converts proxies to sing-box outbound JSON object strings.
// Returns a slice where each element is a single JSON object (no trailing comma, no array wrapping).
func Singbox(proxies []proxy.Proxy, cfg SingboxConfig) ([]string, error) {
	var lines []string
	for i, p := range proxies {
		out, err := proxyToSingbox(p, i, cfg)
		if err != nil {
			return nil, fmt.Errorf("proxy %d: %w", i, err)
		}
		data, err := json.Marshal(out)
		if err != nil {
			return nil, fmt.Errorf("marshal proxy %d: %w", i, err)
		}
		lines = append(lines, string(data))
	}
	return lines, nil
}

//nolint:gocognit
func proxyToSingbox(p proxy.Proxy, index int, cfg SingboxConfig) (*singboxOutbound, error) {
	vp, ok := p.(*proxy.VLESSProxy)
	if !ok {
		return nil, fmt.Errorf("unsupported proxy type: %T", p)
	}

	o := vp.Outbound
	if len(o.Settings.Vnext) == 0 {
		return nil, fmt.Errorf("no vnext in outbound settings")
	}
	vnext := o.Settings.Vnext[0]
	if vnext.Address == "" {
		return nil, fmt.Errorf("missing address")
	}
	if vnext.Port == 0 {
		return nil, fmt.Errorf("missing port")
	}
	if len(vnext.Users) == 0 || vnext.Users[0].ID == "" {
		return nil, fmt.Errorf("missing user id")
	}

	tag := cfg.OutboundPrefix + strconv.Itoa(index+1) + cfg.OutboundSuffix
	if cfg.KeepRemark {
		tag = cfg.OutboundPrefix + p.Remarks() + cfg.OutboundSuffix
	}

	user := vnext.Users[0]
	out := &singboxOutbound{
		Type:           "vless",
		Tag:            tag,
		Server:         vnext.Address,
		ServerPort:     vnext.Port,
		UUID:           user.ID,
		Flow:           user.Flow,
		DomainResolver: cfg.DNSResolver,
	}

	ss := o.StreamSettings
	switch ss.Security {
	case "reality":
		if ss.RealitySettings != nil {
			out.TLS = &singboxTLS{
				Enabled:    true,
				ServerName: ss.RealitySettings.ServerName,
				UTLS: &singboxUTLS{
					Enabled:     true,
					Fingerprint: ss.RealitySettings.Fingerprint,
				},
				Reality: &singboxReality{
					Enabled:   true,
					PublicKey: ss.RealitySettings.PublicKey,
					ShortID:   ss.RealitySettings.ShortID,
				},
			}
		}
	case "tls":
		if ss.TLSSettings != nil {
			tls := &singboxTLS{Enabled: true, ServerName: ss.TLSSettings.ServerName}
			if ss.TLSSettings.Fingerprint != "" {
				tls.UTLS = &singboxUTLS{Enabled: true, Fingerprint: ss.TLSSettings.Fingerprint}
			}
			out.TLS = tls
		}
	}

	return out, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./format/ -run TestSingbox -v`
Expected: ALL PASS

- [ ] **Step 5: Lint**

Run: `golangci-lint run --fix`
Expected: No errors

- [ ] **Step 6: Commit**

```bash
git add format/singbox.go format/singbox_test.go
git commit -m "feat: add sing-box outbound formatter"
```

---

### Task 4: Update api.go

**Files:**
- Modify: `api.go`

- [ ] **Step 1: Add Singbox re-exports**

Add to `api.go` after the `Podkop` line:

```go
var Singbox = format.Singbox
type SingboxConfig = format.SingboxConfig
```

- [ ] **Step 2: Verify build**

Run: `go build ./...`
Expected: No errors

- [ ] **Step 3: Commit**

```bash
git add api.go
git commit -m "refactor: re-export Singbox and SingboxConfig from api"
```

---

### Task 5: Update CLI with new flags and sing-box output

**Files:**
- Modify: `cmd/goxsub/main.go`

- [ ] **Step 1: Add new flags and sing-box format case**

Replace the entire `run()` function in `cmd/goxsub/main.go` with:

```go
//nolint:funlen
func run() int {
	formatFlag := flag.String("format", "uri", "output format: uri, podkop, singbox")
	podkopSection := flag.String("podkop-section", "main", "podkop uci section name")
	singboxDNSResolver := flag.String("singbox-dns-resolver", "dns-local", "sing-box domain_resolver value")
	singboxOutboundPrefix := flag.String("singbox-outbound-prefix", "", "sing-box outbound tag prefix")
	singboxOutboundSuffix := flag.String("singbox-outbound-suffix", "", "sing-box outbound tag suffix")
	keepRemark := flag.Bool("keep-remark", true, "keep original remark or replace with sequential number")
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

	if *formatFlag != "singbox" && (flag.Lookup("singbox-dns-resolver").DefValue != *singboxDNSResolver ||
		*singboxOutboundPrefix != "" || *singboxOutboundSuffix != "") {
		fmt.Fprintf(os.Stderr, "error: --singbox-* flags can only be used with --format singbox\n")
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

	if !*keepRemark {
		for i, p := range proxies {
			if vp, ok := p.(*goxsub.VLESSProxy); ok {
				vp.SetRemarks(strconv.Itoa(i + 1))
			}
		}
	}

	switch *formatFlag {
	case "podkop":
		lines, err := goxsub.Podkop(proxies, *podkopSection)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		for _, line := range lines {
			fmt.Println(line)
		}
	case "singbox":
		cfg := goxsub.SingboxConfig{
			OutboundPrefix: *singboxOutboundPrefix,
			OutboundSuffix: *singboxOutboundSuffix,
			KeepRemark:     *keepRemark,
			DNSResolver:    *singboxDNSResolver,
		}
		lines, err := goxsub.Singbox(proxies, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		for _, line := range lines {
			fmt.Println(line + ",")
		}
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

Add `"strconv"` to the import block.

- [ ] **Step 2: Verify build**

Run: `go build -o build/goxsub ./cmd/goxsub/`
Expected: No errors

- [ ] **Step 3: Lint**

Run: `golangci-lint run --fix`
Expected: No errors

- [ ] **Step 4: Run all tests**

Run: `go test ./... -v`
Expected: ALL PASS

- [ ] **Step 5: Commit**

```bash
git add cmd/goxsub/main.go
git commit -m "feat: add sing-box format and --keep-remark, --singbox-* flags to CLI"
```

---

### Task 6: Final verification

- [ ] **Step 1: Run full test suite**

Run: `go test ./... -cover`
Expected: ALL PASS

- [ ] **Step 2: Lint**

Run: `golangci-lint run`
Expected: No errors

- [ ] **Step 3: Build**

Run: `go build -o build/goxsub ./cmd/goxsub/`
Expected: No errors
