package xray

// ExtractVLESSOutbounds filters outbounds by tag "proxy" or protocol "vless" and returns them as VLESSProxy values.
func ExtractVLESSOutbounds(subs []Subscription) []VLESSProxy {
	var proxies []VLESSProxy
	for _, sub := range subs {
		for _, ob := range sub.Outbounds {
			if ob.Tag == "proxy" || ob.Protocol == "vless" {
				proxies = append(proxies, VLESSProxy{Outbound: ob, Remarks: sub.Remarks})
			}
		}
	}
	return proxies
}
