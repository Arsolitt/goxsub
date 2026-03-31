package format

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Arsolitt/goxsub/proxy"
)

// SingboxConfig holds options for sing-box outbound generation.
type SingboxConfig struct {
	OutboundPrefix string
	OutboundSuffix string
	DNSResolver    string
	KeepRemark     bool
}

type singboxOutbound struct {
	TLS            *singboxTLS `json:"tls,omitempty"`
	Type           string      `json:"type"`
	Tag            string      `json:"tag"`
	Server         string      `json:"server"`
	UUID           string      `json:"uuid"`
	Flow           string      `json:"flow,omitempty"`
	DomainResolver string      `json:"domain_resolver"`
	ServerPort     int         `json:"server_port"`
}

type singboxTLS struct {
	UTLS       *singboxUTLS    `json:"utls,omitempty"`
	Reality    *singboxReality `json:"reality,omitempty"`
	ServerName string          `json:"server_name,omitempty"`
	Enabled    bool            `json:"enabled"`
}

type singboxUTLS struct {
	Fingerprint string `json:"fingerprint"`
	Enabled     bool   `json:"enabled"`
}

type singboxReality struct {
	PublicKey string `json:"public_key"`
	ShortID   string `json:"short_id"`
	Enabled   bool   `json:"enabled"`
}

// Singbox converts proxies to sing-box outbound JSON object strings.
// Returns a slice where each element is a single JSON object (no trailing comma, no array wrapping).
func Singbox(proxies []proxy.Proxy, cfg SingboxConfig) ([]string, error) {
	var lines []string
	for i, p := range proxies {
		out, err := proxyToSingbox(p, i, cfg)
		if err != nil {
			return nil, fmt.Errorf("proxy %d: %w", i, err)
		}
		data, err := json.Marshal(out)
		if err != nil {
			return nil, fmt.Errorf("marshal proxy %d: %w", i, err)
		}
		lines = append(lines, string(data))
	}
	return lines, nil
}

func proxyToSingbox(p proxy.Proxy, index int, cfg SingboxConfig) (*singboxOutbound, error) {
	vp, ok := p.(*proxy.VLESSProxy)
	if !ok {
		return nil, fmt.Errorf("unsupported proxy type: %T", p)
	}

	o := vp.Outbound
	if len(o.Settings.Vnext) == 0 {
		return nil, errors.New("no vnext in outbound settings")
	}
	vnext := o.Settings.Vnext[0]
	if vnext.Address == "" {
		return nil, errors.New("missing address")
	}
	if vnext.Port == 0 {
		return nil, errors.New("missing port")
	}
	if len(vnext.Users) == 0 || vnext.Users[0].ID == "" {
		return nil, errors.New("missing user id")
	}

	tag := cfg.OutboundPrefix + strconv.Itoa(index+1) + cfg.OutboundSuffix
	if cfg.KeepRemark {
		tag = cfg.OutboundPrefix + p.Remarks() + cfg.OutboundSuffix
	}

	user := vnext.Users[0]
	out := &singboxOutbound{
		Type:           "vless",
		Tag:            tag,
		Server:         vnext.Address,
		ServerPort:     vnext.Port,
		UUID:           user.ID,
		Flow:           user.Flow,
		DomainResolver: cfg.DNSResolver,
	}

	ss := o.StreamSettings
	switch ss.Security {
	case "reality":
		if ss.RealitySettings != nil {
			out.TLS = &singboxTLS{
				Enabled:    true,
				ServerName: ss.RealitySettings.ServerName,
				UTLS: &singboxUTLS{
					Enabled:     true,
					Fingerprint: ss.RealitySettings.Fingerprint,
				},
				Reality: &singboxReality{
					Enabled:   true,
					PublicKey: ss.RealitySettings.PublicKey,
					ShortID:   ss.RealitySettings.ShortID,
				},
			}
		}
	case "tls":
		if ss.TLSSettings != nil {
			tls := &singboxTLS{Enabled: true, ServerName: ss.TLSSettings.ServerName}
			if ss.TLSSettings.Fingerprint != "" {
				tls.UTLS = &singboxUTLS{Enabled: true, Fingerprint: ss.TLSSettings.Fingerprint}
			}
			out.TLS = tls
		}
	}

	return out, nil
}
