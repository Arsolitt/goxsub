package xray

import "testing"

func makeProxy(remarks string) VLESSProxy {
	return VLESSProxy{
		Remarks: remarks,
		Outbound: Outbound{
			Protocol: "vless",
			Settings: OutboundSettings{Vnext: []VNext{{Address: "a.com", Port: 443, Users: []User{{ID: "u"}}}}},
		},
	}
}

func TestFilterByRemark_EmptyPatterns(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("Russia Server"), makeProxy("NL Server")}
	result := FilterByRemark(proxies, nil)
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestFilterByRemark_SinglePattern(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("Russia Server"), makeProxy("NL Server")}
	result := FilterByRemark(proxies, []string{"*Russia*"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks != "NL Server" {
		t.Errorf("expected 'NL Server', got %q", result[0].Remarks)
	}
}

func TestFilterByRemark_MultiplePatterns(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("Russia Server"), makeProxy("China Node"), makeProxy("NL Server")}
	result := FilterByRemark(proxies, []string{"*Russia*", "*China*"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks != "NL Server" {
		t.Errorf("expected 'NL Server', got %q", result[0].Remarks)
	}
}

func TestFilterByRemark_CaseInsensitive(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("RUSSIA"), makeProxy("russia"), makeProxy("Russia"), makeProxy("NL")}
	result := FilterByRemark(proxies, []string{"*russia*"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks != "NL" {
		t.Errorf("expected 'NL', got %q", result[0].Remarks)
	}
}

func TestFilterByRemark_GlobSpecials(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("A1"), makeProxy("A2"), makeProxy("B1")}
	result := FilterByRemark(proxies, []string{"A?"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks != "B1" {
		t.Errorf("expected 'B1', got %q", result[0].Remarks)
	}
}

func TestFilterByRemark_GlobCharClass(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("A1"), makeProxy("B1"), makeProxy("C1")}
	result := FilterByRemark(proxies, []string{"[AB]1"})
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
	if result[0].Remarks != "C1" {
		t.Errorf("expected 'C1', got %q", result[0].Remarks)
	}
}

func TestFilterByRemark_NoMatches(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("NL Server"), makeProxy("DE Server")}
	result := FilterByRemark(proxies, []string{"*JP*"})
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestFilterByRemark_AllExcluded(t *testing.T) {
	proxies := []VLESSProxy{makeProxy("RUSSIA"), makeProxy("Russia")}
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
