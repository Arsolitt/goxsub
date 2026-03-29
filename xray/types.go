package xray

import "encoding/json"

type Subscription struct {
	DNS       json.RawMessage `json:"dns"`
	Inbounds  json.RawMessage `json:"inbounds"`
	Log       json.RawMessage `json:"log"`
	Outbounds []Outbound      `json:"outbounds"`
	Remarks   string          `json:"remarks"`
	Routing   json.RawMessage `json:"routing"`
}

type Outbound struct {
	StreamSettings StreamSettings   `json:"streamSettings"`
	Protocol       string           `json:"protocol"`
	Tag            string           `json:"tag"`
	Settings       OutboundSettings `json:"settings"`
}

type OutboundSettings struct {
	Vnext []VNext `json:"vnext,omitempty"`
}

type VNext struct {
	Address string `json:"address"`
	Users   []User `json:"users"`
	Port    int    `json:"port"`
}

type User struct {
	Encryption string `json:"encryption"`
	Flow       string `json:"flow,omitempty"`
	ID         string `json:"id"`
	Level      int    `json:"level,omitempty"`
}

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

type RealitySettings struct {
	PublicKey   string `json:"publicKey"`
	Fingerprint string `json:"fingerprint"`
	ServerName  string `json:"serverName"`
	ShortID     string `json:"shortId"`
	SpiderX     string `json:"spiderX,omitempty"`
}

type TLSSettings struct {
	ServerName  string   `json:"serverName,omitempty"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	ALPN        []string `json:"alpn,omitempty"`
}

type TCPSettings struct {
	Header *TCPHeader `json:"header,omitempty"`
}

type TCPHeader struct {
	Type string `json:"type,omitempty"`
}

type WSSettings struct {
	Path string `json:"path,omitempty"`
	Host string `json:"host,omitempty"`
}

type GRPCSettings struct {
	ServiceName string `json:"serviceName,omitempty"`
}

type HTTPSettings struct {
	Path string   `json:"path,omitempty"`
	Host []string `json:"host,omitempty"`
}

type KCPSettings struct {
	HeaderType *KCPHeader `json:"header,omitempty"`
	Seed       string     `json:"seed,omitempty"`
}

type KCPHeader struct {
	Type string `json:"type,omitempty"`
}

type VLESSProxy struct {
	Remarks  string
	Outbound Outbound
}
