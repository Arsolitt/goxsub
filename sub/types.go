package sub

import "encoding/json"

// Subscription represents a single subscription entry from a JSON subscription feed.
type Subscription struct {
	DNS       json.RawMessage `json:"dns"`
	Inbounds  json.RawMessage `json:"inbounds"`
	Log       json.RawMessage `json:"log"`
	Outbounds []Outbound      `json:"outbounds"`
	Remarks   string          `json:"remarks"`
	Routing   json.RawMessage `json:"routing"`
}

// Outbound represents a proxy outbound configuration.
type Outbound struct {
	StreamSettings StreamSettings   `json:"streamSettings"`
	Protocol       string           `json:"protocol"`
	Tag            string           `json:"tag"`
	Settings       OutboundSettings `json:"settings"`
}

// OutboundSettings holds the settings for an outbound connection.
type OutboundSettings struct {
	Vnext []VNext `json:"vnext,omitempty"`
}

// VNext represents a vnext server configuration.
type VNext struct {
	Address string `json:"address"`
	Users   []User `json:"users"`
	Port    int    `json:"port"`
}

// User represents a user account in a vnext entry.
type User struct {
	Encryption string `json:"encryption"`
	Flow       string `json:"flow,omitempty"`
	ID         string `json:"id"`
	Level      int    `json:"level,omitempty"`
}

// StreamSettings holds the transport and security configuration for an outbound.
type StreamSettings struct {
	RealitySettings *RealitySettings `json:"realitySettings,omitempty"`
	TLSSettings     *TLSSettings     `json:"tlsSettings,omitempty"`
	TCPSettings     *TCPSettings     `json:"tcpSettings,omitempty"`
	WSSettings      *WSSettings      `json:"wsSettings,omitempty"`
	GRPCSettings    *GRPCSettings    `json:"grpcSettings,omitempty"`
	HTTPSettings    *HTTPSettings    `json:"httpSettings,omitempty"`
	KCPSettings     *KCPSettings     `json:"kcpSettings,omitempty"`
	Network         string           `json:"network"`
	Security        string           `json:"security"`
}

// RealitySettings holds REALITY TLS settings.
type RealitySettings struct {
	PublicKey   string `json:"publicKey"`
	Fingerprint string `json:"fingerprint"`
	ServerName  string `json:"serverName"`
	ShortID     string `json:"shortId"`
	SpiderX     string `json:"spiderX,omitempty"`
}

// TLSSettings holds standard TLS settings.
type TLSSettings struct {
	ServerName  string   `json:"serverName,omitempty"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	ALPN        []string `json:"alpn,omitempty"`
}

// TCPSettings holds TCP transport settings.
type TCPSettings struct {
	Header *TCPHeader `json:"header,omitempty"`
}

// TCPHeader holds the TCP header type.
type TCPHeader struct {
	Type string `json:"type,omitempty"`
}

// WSSettings holds WebSocket transport settings.
type WSSettings struct {
	Path string `json:"path,omitempty"`
	Host string `json:"host,omitempty"`
}

// GRPCSettings holds gRPC transport settings.
type GRPCSettings struct {
	ServiceName string `json:"serviceName,omitempty"`
}

// HTTPSettings holds HTTP/2 transport settings.
type HTTPSettings struct {
	Path string   `json:"path,omitempty"`
	Host []string `json:"host,omitempty"`
}

// KCPSettings holds mKCP transport settings.
type KCPSettings struct {
	HeaderType *KCPHeader `json:"header,omitempty"`
	Seed       string     `json:"seed,omitempty"`
}

// KCPHeader holds the mKCP header type.
type KCPHeader struct {
	Type string `json:"type,omitempty"`
}
