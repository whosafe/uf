package rule

import (
	"strings"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// NotBlank 非空白验证规则排除纯空值
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
func (n *NotBlank) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("not_blank", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (n *NotBlank) Name() string {
	return "not_blank"
}

// NewNotBlank 创建非空白验证规则
func NewNotBlank() *NotBlank {
	return &NotBlank{}
}
