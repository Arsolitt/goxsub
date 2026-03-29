package xray

import "testing"

func TestExtractVLESSOutbounds(t *testing.T) {
	subs := []Subscription{
		{
			Outbounds: []Outbound{
				{Protocol: "vless", Tag: "proxy", Settings: OutboundSettings{Vnext: []VNext{{}}}},
				{Protocol: "socks", Tag: "upstream", Settings: OutboundSettings{}},
				{Protocol: "freedom", Tag: "direct", Settings: OutboundSettings{}},
			},
			Remarks: "Server A",
		},
		{
			Outbounds: []Outbound{
				{Protocol: "vless", Tag: "alt-proxy", Settings: OutboundSettings{Vnext: []VNext{{}}}},
				{Protocol: "blackhole", Tag: "block", Settings: OutboundSettings{}},
			},
			Remarks: "Server B",
		},
	}

	proxies := ExtractVLESSOutbounds(subs)
	if len(proxies) != 2 {
		t.Fatalf("expected 2 proxies, got %d", len(proxies))
	}
	if proxies[0].Remarks != "Server A" {
		t.Errorf("expected remarks 'Server A', got %q", proxies[0].Remarks)
	}
	if proxies[0].Outbound.Tag != "proxy" {
		t.Errorf("expected tag 'proxy', got %q", proxies[0].Outbound.Tag)
	}
	if proxies[1].Remarks != "Server B" {
		t.Errorf("expected remarks 'Server B', got %q", proxies[1].Remarks)
	}
	if proxies[1].Outbound.Tag != "alt-proxy" {
		t.Errorf("expected tag 'alt-proxy', got %q", proxies[1].Outbound.Tag)
	}
}

func TestExtractVLESSOutbounds_NoMatches(t *testing.T) {
	subs := []Subscription{
		{
			Outbounds: []Outbound{
				{Protocol: "socks", Tag: "upstream", Settings: OutboundSettings{}},
				{Protocol: "freedom", Tag: "direct", Settings: OutboundSettings{}},
			},
			Remarks: "No VLESS",
		},
	}
	proxies := ExtractVLESSOutbounds(subs)
	if len(proxies) != 0 {
		t.Errorf("expected 0 proxies, got %d", len(proxies))
	}
}
