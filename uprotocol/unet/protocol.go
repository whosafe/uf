package unet

// Protocol 协议类型
type Protocol string

const (
	// ProtocolHTTP HTTP 协议
	ProtocolHTTP Protocol = "http"
	// ProtocolHTTPS HTTPS 协议
	ProtocolHTTPS Protocol = "https"
	// ProtocolTCP TCP 协议
	ProtocolTCP Protocol = "tcp"
	// ProtocolQUIC QUIC 协议
	ProtocolQUIC Protocol = "quic"
)

// String 返回协议字符串
func (p Protocol) String() string {
	return string(p)
}
