package proxy

import "github.com/Arsolitt/goxsub/sub"

const vlessProtocol = "vless"

// ExtractProxies extracts proxy outbounds from subscriptions by protocol and tag.
// Currently supports "vless" protocol and "proxy" tag.
func ExtractProxies(subs []sub.Subscription) []Proxy {
	var proxies []Proxy
	for _, s := range subs {
		for _, ob := range s.Outbounds {
			if ob.Protocol == vlessProtocol || ob.Tag == "proxy" {
				proxies = append(proxies, &VLESSProxy{Outbound: ob, remark: s.Remarks})
			}
		}
	}
	return proxies
}
