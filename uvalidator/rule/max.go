package rule

import (
	"fmt"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// Max 最大值验证规则
type Max struct {
	MaxValue int // 最大值
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
func (m *Max) GetMessage(field string, lang ...uvalidator.Language) string {
	// 默认使用数值消息
	template := i18n.GetMessage("max", lang...)

	// 替换占位符
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", m.MaxValue))

	return msg
}

// Name 规则名称
func (m *Max) Name() string {
	return "max"
}

// NewMax 创建最大值验证规则
func NewMax(maxValue int) *Max {
	return &Max{MaxValue: maxValue}
}
