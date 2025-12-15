package rule

import (
	"fmt"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// Gt 大于验证规则
type Gt struct {
	Value int
}

// Validate 执行验证
func (g *Gt) Validate(value any) bool {
	switch v := value.(type) {
	case int:
		return v > g.Value
	case int64:
		return int(v) > g.Value
	case float64:
		return int(v) > g.Value
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (g *Gt) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("gt", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", g.Value))
	return msg
}

// Name 规则名称
func (g *Gt) Name() string {
	return "gt"
}

// NewGt 创建大于验证规则
func NewGt(value int) *Gt {
	return &Gt{Value: value}
}

// Gte 大于等于验证规则
type Gte struct {
	Value int
}

// Validate 执行验证
func (g *Gte) Validate(value any) bool {
	switch v := value.(type) {
	case int:
		return v >= g.Value
	case int64:
		return int(v) >= g.Value
	case float64:
		return int(v) >= g.Value
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (g *Gte) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("gte", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", g.Value))
	return msg
}

// Name 规则名称
func (g *Gte) Name() string {
	return "gte"
}

// NewGte 创建大于等于验证规则
func NewGte(value int) *Gte {
	return &Gte{Value: value}
}

// Lt 小于验证规则
type Lt struct {
	Value int
}

// Validate 执行验证
func (l *Lt) Validate(value any) bool {
	switch v := value.(type) {
	case int:
		return v < l.Value
	case int64:
		return int(v) < l.Value
	case float64:
		return int(v) < l.Value
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (l *Lt) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("lt", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", l.Value))
	return msg
}

// Name 规则名称
func (l *Lt) Name() string {
	return "lt"
}

// NewLt 创建小于验证规则
func NewLt(value int) *Lt {
	return &Lt{Value: value}
}

// Lte 小于等于验证规则
type Lte struct {
	Value int
}

// Validate 执行验证
func (l *Lte) Validate(value any) bool {
	switch v := value.(type) {
	case int:
		return v <= l.Value
	case int64:
		return int(v) <= l.Value
	case float64:
		return int(v) <= l.Value
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (l *Lte) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("lte", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", l.Value))
	return msg
}

// Name 规则名称
func (l *Lte) Name() string {
	return "lte"
}

// NewLte 创建小于等于验证规则
func NewLte(value int) *Lte {
	return &Lte{Value: value}
}
