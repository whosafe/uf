package rule

import (
	"github.com/whosafe/uf/uvalidator"
	"github.com/whosafe/uf/uvalidator/i18n"
)

// Positive 正数验证规则 (value > 0)
type Positive struct{}

// Validate 执行验证
func (p *Positive) Validate(value any) bool {
	switch v := value.(type) {
	case int:
		return v > 0
	case int64:
		return v > 0
	case float64:
		return v > 0
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (p *Positive) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("positive", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (p *Positive) Name() string {
	return "positive"
}

// NewPositive 创建正数验证规则
func NewPositive() *Positive {
	return &Positive{}
}

// Negative 负数验证规则 (value < 0)
type Negative struct{}

// Validate 执行验证
func (n *Negative) Validate(value any) bool {
	switch v := value.(type) {
	case int:
		return v < 0
	case int64:
		return v < 0
	case float64:
		return v < 0
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (n *Negative) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("negative", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (n *Negative) Name() string {
	return "negative"
}

// NewNegative 创建负数验证规则
func NewNegative() *Negative {
	return &Negative{}
}
