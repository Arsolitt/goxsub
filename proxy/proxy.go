package proxy

import "github.com/Arsolitt/goxsub/sub"

// Proxy is the interface implemented by all protocol-specific proxy types.
type Proxy interface {
	Protocol() string
	Tag() string
	Remarks() string
}

// VLESSProxy represents a VLESS protocol proxy extracted from a subscription.
type VLESSProxy struct {
	remark   string
	Outbound sub.Outbound
}

func (v *VLESSProxy) Protocol() string { return vlessProtocol }
func (v *VLESSProxy) Tag() string      { return v.Outbound.Tag }
func (v *VLESSProxy) Remarks() string  { return v.remark }
