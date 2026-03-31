package goxsub

import (
	"github.com/Arsolitt/goxsub/format"
	"github.com/Arsolitt/goxsub/protocol"
	"github.com/Arsolitt/goxsub/proxy"
	"github.com/Arsolitt/goxsub/sub"
)

type Subscription = sub.Subscription
type Outbound = sub.Outbound
type OutboundSettings = sub.OutboundSettings
type VNext = sub.VNext
type User = sub.User
type StreamSettings = sub.StreamSettings
type RealitySettings = sub.RealitySettings
type TLSSettings = sub.TLSSettings
type TCPSettings = sub.TCPSettings
type TCPHeader = sub.TCPHeader
type WSSettings = sub.WSSettings
type GRPCSettings = sub.GRPCSettings
type HTTPSettings = sub.HTTPSettings
type KCPSettings = sub.KCPSettings
type KCPHeader = sub.KCPHeader

type Proxy = proxy.Proxy
type VLESSProxy = proxy.VLESSProxy

var ParseSubscription = sub.ParseSubscription
var ExtractProxies = proxy.ExtractProxies
var FilterByRemark = proxy.FilterByRemark
var ToURI = protocol.ToURI
var ToVLESSURI = protocol.ToVLESSURI
var VLESSURI = protocol.VLESSURI
var Podkop = format.Podkop
var Singbox = format.Singbox

type SingboxConfig = format.SingboxConfig
