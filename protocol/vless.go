package protocol

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/Arsolitt/goxsub/proxy"
)

// VLESSURI converts a VLESSProxy to a vless:// URI string.

//nolint:gocognit,gocyclo,cyclop,funlen // handles multiple transport and security types
func VLESSURI(p *proxy.VLESSProxy) (string, error) {
	o := p.Outbound
	remarks := p.Remarks()
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
	params = append(params, "type="+url.QueryEscape(ss.Network))
	params = append(params, "encryption="+url.QueryEscape(user.Encryption))
	params = append(params, "security="+url.QueryEscape(ss.Security))

	switch ss.Security {
	case "reality":
		if ss.RealitySettings != nil {
			params = append(params, "pbk="+url.QueryEscape(ss.RealitySettings.PublicKey))
			params = append(params, "fp="+url.QueryEscape(ss.RealitySettings.Fingerprint))
			params = append(params, "sni="+url.QueryEscape(ss.RealitySettings.ServerName))
			params = append(params, "sid="+url.QueryEscape(ss.RealitySettings.ShortID))
			if ss.RealitySettings.SpiderX != "" {
				params = append(params, "spx="+url.QueryEscape(ss.RealitySettings.SpiderX))
			}
		}
	case "tls":
		if ss.TLSSettings != nil {
			if ss.TLSSettings.ServerName != "" {
				params = append(params, "sni="+url.QueryEscape(ss.TLSSettings.ServerName))
			}
			if len(ss.TLSSettings.ALPN) > 0 {
				params = append(params, "alpn="+url.QueryEscape(strings.Join(ss.TLSSettings.ALPN, ",")))
			}
			if ss.TLSSettings.Fingerprint != "" {
				params = append(params, "fp="+url.QueryEscape(ss.TLSSettings.Fingerprint))
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
				params = append(params, "host="+url.QueryEscape(ss.WSSettings.Host))
			}
		}
	case "grpc":
		if ss.GRPCSettings != nil && ss.GRPCSettings.ServiceName != "" {
			params = append(params, "serviceName="+url.QueryEscape(ss.GRPCSettings.ServiceName))
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
				params = append(params, "headerType="+url.QueryEscape(ss.KCPSettings.HeaderType.Type))
			}
		}
	}

	if user.Flow != "" {
		params = append(params, "flow="+url.QueryEscape(user.Flow))
	}

	fragment := ""
	if remarks != "" {
		fragment = "#" + url.PathEscape(remarks)
	}

	return fmt.Sprintf("vless://%s@%s?%s%s",
		user.ID, net.JoinHostPort(vnext.Address, strconv.Itoa(vnext.Port)),
		strings.Join(params, "&"), fragment), nil
}

// ToVLESSURI converts a VLESSProxy (as a Proxy interface) to a vless:// URI string.
// Returns an error if the proxy is not a VLESSProxy or conversion fails.
func ToVLESSURI(p proxy.Proxy) (string, error) {
	vp, ok := p.(*proxy.VLESSProxy)
	if !ok {
		return "", fmt.Errorf("expected VLESSProxy, got %T", p)
	}
	return VLESSURI(vp)
}

// ToURI converts a proxy to its protocol URI string.
// Currently supports vless. Returns an error for unsupported protocols.
func ToURI(p proxy.Proxy) (string, error) {
	switch p.Protocol() {
	case "vless":
		return ToVLESSURI(p)
	default:
		return "", fmt.Errorf("unsupported protocol: %s", p.Protocol())
	}
}
