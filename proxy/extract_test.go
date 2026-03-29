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
