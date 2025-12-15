package rule

import (
	"strings"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// NotBlank 非空白验证规�?排除纯空�?
type NotBlank struct{}

// Validate 执行验证
func (n *NotBlank) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	return strings.TrimSpace(str) != ""
}

// GetMessage 获取错误消息
func (n *NotBlank) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("not_blank", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (n *NotBlank) Name() string {
	return "not_blank"
}

// NewNotBlank 创建非空白验证规�?
func NewNotBlank() *NotBlank {
	return &NotBlank{}
}
