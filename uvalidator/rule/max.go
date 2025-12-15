package rule

import (
	"fmt"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// Max 最大值验证规�?
type Max struct {
	MaxValue int // 最大�?
}

// Validate 执行验证
func (m *Max) Validate(value any) bool {
	switch v := value.(type) {
	case string:
		return len(v) <= m.MaxValue
	case int:
		return v <= m.MaxValue
	case int64:
		return int(v) <= m.MaxValue
	case float64:
		return int(v) <= m.MaxValue
	case []any:
		return len(v) <= m.MaxValue
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (m *Max) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	// 判断是字符串长度还是数�?
	key := "max"
	if params["type"] == "string" {
		key = "max_length"
	}

	template := i18n.GetMessage(key)

	// 替换占位�?
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", m.MaxValue))

	return msg
}

// Name 规则名称
func (m *Max) Name() string {
	return "max"
}

// NewMax 创建最大值验证规�?
func NewMax(maxValue int) *Max {
	return &Max{MaxValue: maxValue}
}
