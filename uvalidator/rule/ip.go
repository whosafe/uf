package rule

import (
	"net"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// IP IP地址验证规则 (IPv4 与 IPv6)
type IP struct{}

// Validate 执行验证
func (i *IP) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return net.ParseIP(str) != nil
}

// GetMessage 获取错误消息
func (i *IP) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("ip", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (i *IP) Name() string {
	return "ip"
}

// NewIP 创建IP地址验证规则
func NewIP() *IP {
	return &IP{}
}

// IPv4 IPv4地址验证规则
type IPv4 struct{}

// Validate 执行验证
func (i *IPv4) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	ip := net.ParseIP(str)
	if ip == nil {
		return false
	}

	return ip.To4() != nil
}

// GetMessage 获取错误消息
func (i *IPv4) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("ipv4", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (i *IPv4) Name() string {
	return "ipv4"
}

// NewIPv4 创建IPv4地址验证规则
func NewIPv4() *IPv4 {
	return &IPv4{}
}

// IPv6 IPv6地址验证规则
type IPv6 struct{}

// Validate 执行验证
func (i *IPv6) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	ip := net.ParseIP(str)
	if ip == nil {
		return false
	}

	return ip.To4() == nil
}

// GetMessage 获取错误消息
func (i *IPv6) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("ipv6", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (i *IPv6) Name() string {
	return "ipv6"
}

// NewIPv6 创建IPv6地址验证规则
func NewIPv6() *IPv6 {
	return &IPv6{}
}
