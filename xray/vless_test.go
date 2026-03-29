package xray

import (
	"strings"
	"testing"
)

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

func TestToVLESSURI_TCPReality(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{
				Address: "example.com",
				Port:    443,
				Users:   []User{{ID: "test-uuid-1234", Encryption: "none", Flow: "xtls-rprx-vision"}},
			}},
		},
		StreamSettings: StreamSettings{
			Network:  "tcp",
			Security: "reality",
			RealitySettings: &RealitySettings{
				PublicKey:   "pub_key_value",
				Fingerprint: "firefox",
				ServerName:  "sni.example.com",
				ShortID:     "abcd1234",
				SpiderX:     "/",
			},
		},
	}

	uri, err := ToVLESSURI(ob, "Test TCP Reality")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://test-uuid-1234@example.com:443?type=tcp&encryption=none&security=reality&pbk=pub_key_value&fp=firefox&sni=sni.example.com&sid=abcd1234&spx=%2F&flow=xtls-rprx-vision#Test%20TCP%20Reality"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestToVLESSURI_TCPTLS(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{
				Address: "tls.example.com",
				Port:    8443,
				Users:   []User{{ID: "uuid-tls", Encryption: "none"}},
			}},
		},
		StreamSettings: StreamSettings{
			Network:  "tcp",
			Security: "tls",
			TLSSettings: &TLSSettings{
				ServerName:  "tls.example.com",
				ALPN:        []string{"h2", "http/1.1"},
				Fingerprint: "chrome",
			},
		},
	}
	uri, err := ToVLESSURI(ob, "Test TCP TLS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-tls@tls.example.com:8443?type=tcp&encryption=none&security=tls&sni=tls.example.com&alpn=h2%2Chttp%2F1.1&fp=chrome#Test%20TCP%20TLS"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestToVLESSURI_WSTLS(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{
				Address: "ws.example.com",
				Port:    443,
				Users:   []User{{ID: "uuid-ws", Encryption: "none"}},
			}},
		},
		StreamSettings: StreamSettings{
			Network:  "ws",
			Security: "tls",
			TLSSettings: &TLSSettings{
				ServerName: "ws.example.com",
			},
			WSSettings: &WSSettings{
				Path: "/websocket-path",
				Host: "ws-host.example.com",
			},
		},
	}
	uri, err := ToVLESSURI(ob, "Test WS TLS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-ws@ws.example.com:443?type=ws&encryption=none&security=tls&sni=ws.example.com&path=%2Fwebsocket-path&host=ws-host.example.com#Test%20WS%20TLS"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestToVLESSURI_GrpctLS(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{
				Address: "grpc.example.com",
				Port:    443,
				Users:   []User{{ID: "uuid-grpc", Encryption: "none"}},
			}},
		},
		StreamSettings: StreamSettings{
			Network:  "grpc",
			Security: "tls",
			TLSSettings: &TLSSettings{
				ServerName: "grpc.example.com",
			},
			GRPCSettings: &GRPCSettings{
				ServiceName: "grpc-service",
			},
		},
	}
	uri, err := ToVLESSURI(ob, "Test GRPC TLS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-grpc@grpc.example.com:443?type=grpc&encryption=none&security=tls&sni=grpc.example.com&serviceName=grpc-service#Test%20GRPC%20TLS"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestToVLESSURI_H2TLS(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{
				Address: "h2.example.com",
				Port:    443,
				Users:   []User{{ID: "uuid-h2", Encryption: "none"}},
			}},
		},
		StreamSettings: StreamSettings{
			Network:  "h2",
			Security: "tls",
			TLSSettings: &TLSSettings{
				ServerName: "h2.example.com",
			},
			HTTPSettings: &HTTPSettings{
				Path: "/h2-path",
				Host: []string{"h2a.example.com", "h2b.example.com"},
			},
		},
	}
	uri, err := ToVLESSURI(ob, "Test H2 TLS")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-h2@h2.example.com:443?type=h2&encryption=none&security=tls&sni=h2.example.com&path=%2Fh2-path&host=h2a.example.com%2Ch2b.example.com#Test%20H2%20TLS"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestToVLESSURI_KCPNone(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{
				Address: "kcp.example.com",
				Port:    443,
				Users:   []User{{ID: "uuid-kcp", Encryption: "none"}},
			}},
		},
		StreamSettings: StreamSettings{
			Network:  "kcp",
			Security: "none",
			KCPSettings: &KCPSettings{
				Seed:       "mykcpseed",
				HeaderType: &KCPHeader{Type: "none"},
			},
		},
	}
	uri, err := ToVLESSURI(ob, "Test KCP None")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-kcp@kcp.example.com:443?type=kcp&encryption=none&security=none&type=mykcpseed&headerType=none#Test%20KCP%20None"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestToVLESSURI_TCPNone(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{
				Address: "plain.example.com",
				Port:    80,
				Users:   []User{{ID: "uuid-plain", Encryption: "none"}},
			}},
		},
		StreamSettings: StreamSettings{
			Network:  "tcp",
			Security: "none",
		},
	}
	uri, err := ToVLESSURI(ob, "Test Plain")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid-plain@plain.example.com:80?type=tcp&encryption=none&security=none#Test%20Plain"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestToVLESSURI_NoVnext(t *testing.T) {
	ob := Outbound{Protocol: "vless", Settings: OutboundSettings{}}
	_, err := ToVLESSURI(ob, "test")
	if err == nil {
		t.Fatal("expected error for no vnext")
	}
}

func TestToVLESSURI_MissingAddress(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{Port: 443, Users: []User{{ID: "uuid"}}}},
		},
	}
	_, err := ToVLESSURI(ob, "test")
	if err == nil {
		t.Fatal("expected error for missing address")
	}
}

func TestToVLESSURI_MissingPort(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{Address: "example.com", Users: []User{{ID: "uuid"}}}},
		},
	}
	_, err := ToVLESSURI(ob, "test")
	if err == nil {
		t.Fatal("expected error for missing port")
	}
}

func TestToVLESSURI_MissingUserID(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{Address: "example.com", Port: 443, Users: []User{}}},
		},
	}
	_, err := ToVLESSURI(ob, "test")
	if err == nil {
		t.Fatal("expected error for missing user id")
	}
}

func TestToVLESSURI_NoRemarks(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{
				Address: "example.com",
				Port:    443,
				Users:   []User{{ID: "uuid", Encryption: "none"}},
			}},
		},
		StreamSettings: StreamSettings{
			Network:  "tcp",
			Security: "none",
		},
	}
	uri, err := ToVLESSURI(ob, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid@example.com:443?type=tcp&encryption=none&security=none"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestToVLESSURI_RealityNoSpiderX(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{
				Address: "example.com",
				Port:    443,
				Users:   []User{{ID: "uuid", Encryption: "none"}},
			}},
		},
		StreamSettings: StreamSettings{
			Network:  "tcp",
			Security: "reality",
			RealitySettings: &RealitySettings{
				PublicKey:   "pubkey",
				Fingerprint: "chrome",
				ServerName:  "sni.com",
				ShortID:     "1234",
			},
		},
	}
	uri, err := ToVLESSURI(ob, "NoSpiderX")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid@example.com:443?type=tcp&encryption=none&security=reality&pbk=pubkey&fp=chrome&sni=sni.com&sid=1234#NoSpiderX"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
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

func TestFormatPodkop(t *testing.T) {
	proxies := []VLESSProxy{
		{
			Remarks: "Server A",
			Outbound: Outbound{
				Protocol: "vless",
				Settings: OutboundSettings{
					Vnext: []VNext{{
						Address: "a.example.com",
						Port:    443,
						Users:   []User{{ID: "uuid-a", Encryption: "none"}},
					}},
				},
				StreamSettings: StreamSettings{
					Network:  "tcp",
					Security: "none",
				},
			},
		},
		{
			Remarks: "Server B",
			Outbound: Outbound{
				Protocol: "vless",
				Settings: OutboundSettings{
					Vnext: []VNext{{
						Address: "b.example.com",
						Port:    8443,
						Users:   []User{{ID: "uuid-b", Encryption: "none"}},
					}},
				},
				StreamSettings: StreamSettings{
					Network:  "tcp",
					Security: "none",
				},
			},
		},
	}

	result, err := FormatPodkop(proxies, "main")
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

func TestFormatPodkop_CustomSection(t *testing.T) {
	proxies := []VLESSProxy{
		{
			Remarks: "S",
			Outbound: Outbound{
				Protocol: "vless",
				Settings: OutboundSettings{
					Vnext: []VNext{{
						Address: "c.example.com",
						Port:    443,
						Users:   []User{{ID: "uuid-c", Encryption: "none"}},
					}},
				},
				StreamSettings: StreamSettings{
					Network:  "tcp",
					Security: "none",
				},
			},
		},
	}

	result, err := FormatPodkop(proxies, "backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "uci del podkop.backup.urltest_proxy_links\nuci add_list podkop.backup.urltest_proxy_links='vless://uuid-c@c.example.com:443?type=tcp&encryption=none&security=none#S'\nservice podkop restart"
	if result != expected {
		t.Errorf("result mismatch:\ngot:      %s\nexpected: %s", result, expected)
	}
}

func TestFormatPodkop_Empty(t *testing.T) {
	result, err := FormatPodkop(nil, "main")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "uci del podkop.main.urltest_proxy_links\nservice podkop restart"
	if result != expected {
		t.Errorf("result mismatch:\ngot:      %s\nexpected: %s", result, expected)
	}
}

func TestFormatPodkop_InvalidProxy(t *testing.T) {
	proxies := []VLESSProxy{
		{Remarks: "bad", Outbound: Outbound{Protocol: "vless", Settings: OutboundSettings{}}},
	}
	_, err := FormatPodkop(proxies, "main")
	if err == nil {
		t.Fatal("expected error for invalid proxy")
	}
}

func TestToVLESSURI_RealityEncodedPublicKey(t *testing.T) {
	ob := Outbound{
		Protocol: "vless",
		Settings: OutboundSettings{
			Vnext: []VNext{{
				Address: "example.com",
				Port:    443,
				Users:   []User{{ID: "uuid", Encryption: "none"}},
			}},
		},
		StreamSettings: StreamSettings{
			Network:  "tcp",
			Security: "reality",
			RealitySettings: &RealitySettings{
				PublicKey:   "abc123+def/ghi=",
				Fingerprint: "chrome",
				ServerName:  "sni.com",
				ShortID:     "abcd",
			},
		},
	}
	uri, err := ToVLESSURI(ob, "Encoded PBK")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid@example.com:443?type=tcp&encryption=none&security=reality&pbk=abc123%2Bdef%2Fghi%3D&fp=chrome&sni=sni.com&sid=abcd#Encoded%20PBK"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}
