package rule

import (
	"fmt"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// Min 最小值验证规则
type Min struct {
	MinValue int // 最小值
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
func (m *Min) GetMessage(field string, lang ...uvalidator.Language) string {
	// 默认使用数值消息
	template := i18n.GetMessage("min", lang...)

	// 替换占位符
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", m.MinValue))

	return msg
}

// Name 规则名称
func (m *Min) Name() string {
	return "min"
}

// NewMin 创建最小值验证规则
func NewMin(minValue int) *Min {
	return &Min{MinValue: minValue}
}

// replaceAll 简单的字符串替换
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
