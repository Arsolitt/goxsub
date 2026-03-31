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

	var raw map[string]any
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
