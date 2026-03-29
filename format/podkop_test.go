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
