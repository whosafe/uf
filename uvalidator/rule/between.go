package rule

import (
	"fmt"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// Between 数值范围验证规�?(min <= value <= max)
type Between struct {
	Min int
	Max int
}

// Validate 执行验证
func (b *Between) Validate(value any) bool {
	switch v := value.(type) {
	case int:
		return v >= b.Min && v <= b.Max
	case int64:
		return int(v) >= b.Min && int(v) <= b.Max
	case float64:
		return int(v) >= b.Min && int(v) <= b.Max
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (b *Between) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("between", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{min}", fmt.Sprintf("%d", b.Min))
	msg = replaceAll(msg, "{max}", fmt.Sprintf("%d", b.Max))
	return msg
}

// Name 规则名称
func (b *Between) Name() string {
	return "between"
}

// NewBetween 创建数值范围验证规�?
func NewBetween(min, max int) *Between {
	return &Between{Min: min, Max: max}
}
