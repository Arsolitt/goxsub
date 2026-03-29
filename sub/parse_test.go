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
