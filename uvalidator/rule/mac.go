package rule

import (
	"net"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// MAC MAC地址验证规则
type MAC struct{}

// Validate 执行验证
func (m *MAC) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	_, err := net.ParseMAC(str)
	return err == nil
}

// GetMessage 获取错误消息
func (m *MAC) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("mac", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (m *MAC) Name() string {
	return "mac"
}

// NewMAC 创建MAC地址验证规则
func NewMAC() *MAC {
	return &MAC{}
}
