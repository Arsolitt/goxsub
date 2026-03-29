package xray

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

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

// ToVLESSURI converts an Outbound with vless protocol to a vless:// URI string.
//
//nolint:gocognit,gocyclo,cyclop,funlen // handles multiple transport and security types
func ToVLESSURI(o Outbound, remarks string) (string, error) {
	if len(o.Settings.Vnext) == 0 {
		return "", errors.New("no vnext in outbound settings")
	}
	vnext := o.Settings.Vnext[0]
	if vnext.Address == "" {
		return "", errors.New("missing address")
	}
	if vnext.Port == 0 {
		return "", errors.New("missing port")
	}
	if len(vnext.Users) == 0 || vnext.Users[0].ID == "" {
		return "", errors.New("missing user id")
	}

	user := vnext.Users[0]
	ss := o.StreamSettings

	var params []string
	params = append(params, "type="+ss.Network)
	params = append(params, "encryption="+user.Encryption)
	params = append(params, "security="+ss.Security)

	switch ss.Security {
	case "reality":
		if ss.RealitySettings != nil {
			params = append(params, "pbk="+ss.RealitySettings.PublicKey)
			params = append(params, "fp="+ss.RealitySettings.Fingerprint)
			params = append(params, "sni="+ss.RealitySettings.ServerName)
			params = append(params, "sid="+ss.RealitySettings.ShortID)
			if ss.RealitySettings.SpiderX != "" {
				params = append(params, "spx="+url.QueryEscape(ss.RealitySettings.SpiderX))
			}
		}
	case "tls":
		if ss.TLSSettings != nil {
			if ss.TLSSettings.ServerName != "" {
				params = append(params, "sni="+ss.TLSSettings.ServerName)
			}
			if len(ss.TLSSettings.ALPN) > 0 {
				params = append(params, "alpn="+url.QueryEscape(strings.Join(ss.TLSSettings.ALPN, ",")))
			}
			if ss.TLSSettings.Fingerprint != "" {
				params = append(params, "fp="+ss.TLSSettings.Fingerprint)
			}
		}
	}

	switch ss.Network {
	case "ws":
		if ss.WSSettings != nil {
			if ss.WSSettings.Path != "" {
				params = append(params, "path="+url.QueryEscape(ss.WSSettings.Path))
			}
			if ss.WSSettings.Host != "" {
				params = append(params, "host="+ss.WSSettings.Host)
			}
		}
	case "grpc":
		if ss.GRPCSettings != nil && ss.GRPCSettings.ServiceName != "" {
			params = append(params, "serviceName="+ss.GRPCSettings.ServiceName)
		}
	case "h2":
		if ss.HTTPSettings != nil {
			if ss.HTTPSettings.Path != "" {
				params = append(params, "path="+url.QueryEscape(ss.HTTPSettings.Path))
			}
			if len(ss.HTTPSettings.Host) > 0 {
				params = append(params, "host="+url.QueryEscape(strings.Join(ss.HTTPSettings.Host, ",")))
			}
		}
	case "kcp":
		if ss.KCPSettings != nil {
			if ss.KCPSettings.Seed != "" {
				params = append(params, "type="+url.QueryEscape(ss.KCPSettings.Seed))
			}
			if ss.KCPSettings.HeaderType != nil && ss.KCPSettings.HeaderType.Type != "" {
				params = append(params, "headerType="+ss.KCPSettings.HeaderType.Type)
			}
		}
	}

	if user.Flow != "" {
		params = append(params, "flow="+user.Flow)
	}

	fragment := ""
	if remarks != "" {
		fragment = "#" + remarks
	}

	return fmt.Sprintf("vless://%s@%s?%s%s",
		user.ID, net.JoinHostPort(vnext.Address, strconv.Itoa(vnext.Port)),
		strings.Join(params, "&"), fragment), nil
}
