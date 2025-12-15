package rule

import (
	"iutime.com/utime/uf/uvalidator"
	"iutime.com/utime/uf/uvalidator/i18n"
)

// Required 必填验证规则
type Required struct{}

// Validate 执行验证
func (r *Required) Validate(value any) bool {
	switch v := value.(type) {
	case string:
		return v != ""
	case int:
		return v != 0
	case int64:
		return v != 0
	case float64:
		return v != 0
	default:
		return value != nil
	}
}

// GetMessage 获取错误消息
func (r *Required) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("required", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (r *Required) Name() string {
	return "required"
}

// NewRequired 创建必填验证规则
func NewRequired() *Required {
	return &Required{}
}
