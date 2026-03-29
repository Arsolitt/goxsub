package protocol

import (
	"testing"

	"github.com/Arsolitt/goxsub/proxy"
	"github.com/Arsolitt/goxsub/sub"
)

func TestVLESSURI_TCPReality(t *testing.T) {
	p := &proxy.VLESSProxy{
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
	expected := "vless://test-uuid-1234@example.com:443?type=tcp&encryption=none&security=reality&pbk=pub_key_value&fp=firefox&sni=sni.example.com&sid=abcd1234&spx=%2F&flow=xtls-rprx-vision"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_TCPTLS(t *testing.T) {
	p := &proxy.VLESSProxy{
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
	expected := "vless://uuid-tls@tls.example.com:8443?type=tcp&encryption=none&security=tls&sni=tls.example.com&alpn=h2%2Chttp%2F1.1&fp=chrome"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_WSTLS(t *testing.T) {
	p := &proxy.VLESSProxy{
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
	expected := "vless://uuid-ws@ws.example.com:443?type=ws&encryption=none&security=tls&sni=ws.example.com&path=%2Fwebsocket-path&host=ws-host.example.com"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_GrpctLS(t *testing.T) {
	p := &proxy.VLESSProxy{
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
	expected := "vless://uuid-grpc@grpc.example.com:443?type=grpc&encryption=none&security=tls&sni=grpc.example.com&serviceName=grpc-service"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_H2TLS(t *testing.T) {
	p := &proxy.VLESSProxy{
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
	expected := "vless://uuid-h2@h2.example.com:443?type=h2&encryption=none&security=tls&sni=h2.example.com&path=%2Fh2-path&host=h2a.example.com%2Ch2b.example.com"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_KCPNone(t *testing.T) {
	p := &proxy.VLESSProxy{
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
	expected := "vless://uuid-kcp@kcp.example.com:443?type=kcp&encryption=none&security=none&type=mykcpseed&headerType=none"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_TCPNone(t *testing.T) {
	p := &proxy.VLESSProxy{
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
	expected := "vless://uuid-plain@plain.example.com:80?type=tcp&encryption=none&security=none"
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
	expected := "vless://uuid@example.com:443?type=tcp&encryption=none&security=reality&pbk=pubkey&fp=chrome&sni=sni.com&sid=1234"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_RealityEncodedPublicKey(t *testing.T) {
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
	expected := "vless://uuid@example.com:443?type=tcp&encryption=none&security=reality&pbk=abc123%2Bdef%2Fghi%3D&fp=chrome&sni=sni.com&sid=abcd"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestVLESSURI_WithRemarks(t *testing.T) {
	subs := []sub.Subscription{{
		Remarks: "My Proxy",
		Outbounds: []sub.Outbound{{
			Protocol: "vless",
			Tag:      "proxy",
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
		}},
	}}
	proxies := proxy.ExtractProxies(subs)
	if len(proxies) == 0 {
		t.Fatal("expected at least one proxy")
	}
	vp, ok := proxies[0].(*proxy.VLESSProxy)
	if !ok {
		t.Fatal("expected VLESSProxy")
	}
	uri, err := VLESSURI(vp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid@example.com:443?type=tcp&encryption=none&security=none#My%20Proxy"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

func TestToVLESSURI_WrongType(t *testing.T) {
	p := &mockProxy{proto: "vless"}
	_, err := ToVLESSURI(p)
	if err == nil {
		t.Fatal("expected error for non-VLESSProxy type")
	}
}

func TestToURI_Unsupported(t *testing.T) {
	p := &mockProxy{proto: "trojan"}
	_, err := ToURI(p)
	if err == nil {
		t.Fatal("expected error for unsupported protocol")
	}
}

func TestToURI_VLESS(t *testing.T) {
	subs := []sub.Subscription{{
		Remarks: "test",
		Outbounds: []sub.Outbound{{
			Protocol: "vless",
			Tag:      "proxy",
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
		}},
	}}
	proxies := proxy.ExtractProxies(subs)
	if len(proxies) == 0 {
		t.Fatal("expected at least one proxy")
	}
	uri, err := ToURI(proxies[0])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "vless://uuid@example.com:443?type=tcp&encryption=none&security=none#test"
	if uri != expected {
		t.Errorf("URI mismatch:\ngot:      %s\nexpected: %s", uri, expected)
	}
}

type mockProxy struct {
	proto string
}

func (m *mockProxy) Protocol() string { return m.proto }
func (m *mockProxy) Tag() string      { return "" }
func (m *mockProxy) Remarks() string  { return "" }
