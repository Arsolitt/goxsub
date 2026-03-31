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
