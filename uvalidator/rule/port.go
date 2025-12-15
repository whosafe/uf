package rule

import (
	"iutime.com/utime/uf/uvalidator"
	"iutime.com/utime/uf/uvalidator/i18n"
)

// Port 端口号验证规�?(1-65535)
type Port struct{}

// Validate 执行验证
func (p *Port) Validate(value any) bool {
	switch v := value.(type) {
	case int:
		return v >= 1 && v <= 65535
	case int64:
		return v >= 1 && v <= 65535
	case string:
		if v == "" {
			return true
		}
		// 尝试转换为整�?
		var port int
		for _, r := range v {
			if r < '0' || r > '9' {
				return false
			}
			port = port*10 + int(r-'0')
		}
		return port >= 1 && port <= 65535
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (p *Port) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("port", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (p *Port) Name() string {
	return "port"
}

// NewPort 创建端口号验证规�?
func NewPort() *Port {
	return &Port{}
}
