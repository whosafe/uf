package rule

import (
	"fmt"
	"strings"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// Confirmed 确认字段验证规则 (用于密码确认等场�?
type Confirmed struct {
	ConfirmationValue any // 确认�?
}

// Validate 执行验证
func (c *Confirmed) Validate(value any) bool {
	return value == c.ConfirmationValue
}

// GetMessage 获取错误消息
func (c *Confirmed) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("confirmed", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (c *Confirmed) Name() string {
	return "confirmed"
}

// NewConfirmed 创建确认字段验证规则
func NewConfirmed(confirmationValue any) *Confirmed {
	return &Confirmed{ConfirmationValue: confirmationValue}
}

// Distinct 不同于指定值验证规�?
type Distinct struct {
	DisallowedValue any
}

// Validate 执行验证
func (d *Distinct) Validate(value any) bool {
	return value != d.DisallowedValue
}

// GetMessage 获取错误消息
func (d *Distinct) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("distinct", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%v", d.DisallowedValue))
	return msg
}

// Name 规则名称
func (d *Distinct) Name() string {
	return "distinct"
}

// NewDistinct 创建不同于指定值验证规�?
func NewDistinct(disallowedValue any) *Distinct {
	return &Distinct{DisallowedValue: disallowedValue}
}

// NotIn 不在指定列表中验证规�?
type NotIn struct {
	Disallowed []string
}

// Validate 执行验证
func (n *NotIn) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	for _, disallowed := range n.Disallowed {
		if str == disallowed {
			return false
		}
	}
	return true
}

// GetMessage 获取错误消息
func (n *NotIn) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("not_in", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", strings.Join(n.Disallowed, ", "))
	return msg
}

// Name 规则名称
func (n *NotIn) Name() string {
	return "not_in"
}

// NewNotIn 创建不在列表中验证规�?
func NewNotIn(disallowed ...string) *NotIn {
	return &NotIn{Disallowed: disallowed}
}

// Nullable 允许为null验证规则
type Nullable struct{}

// Validate 执行验证
func (n *Nullable) Validate(value any) bool {
	return true // 总是通过,允许任何值包括nil
}

// GetMessage 获取错误消息
func (n *Nullable) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("nullable", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (n *Nullable) Name() string {
	return "nullable"
}

// NewNullable 创建允许null验证规则
func NewNullable() *Nullable {
	return &Nullable{}
}
