package rule

import (
	"fmt"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// Min 最小值验证规�?
type Min struct {
	MinValue int // 最小�?
}

// Validate 执行验证
func (m *Min) Validate(value any) bool {
	switch v := value.(type) {
	case string:
		return len(v) >= m.MinValue
	case int:
		return v >= m.MinValue
	case int64:
		return int(v) >= m.MinValue
	case float64:
		return int(v) >= m.MinValue
	case []any:
		return len(v) >= m.MinValue
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (m *Min) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	// 判断是字符串长度还是数�?
	key := "min"
	if params["type"] == "string" {
		key = "min_length"
	}

	template := i18n.GetMessage(key, lang...)

	// 替换占位�?
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", m.MinValue))

	return msg
}

// Name 规则名称
func (m *Min) Name() string {
	return "min"
}

// NewMin 创建最小值验证规�?
func NewMin(minValue int) *Min {
	return &Min{MinValue: minValue}
}

// replaceAll 简单的字符串替�?
func replaceAll(s, old, new string) string {
	result := ""
	for {
		i := indexOf(s, old)
		if i == -1 {
			result += s
			break
		}
		result += s[:i] + new
		s = s[i+len(old):]
	}
	return result
}

// indexOf 查找子串位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
