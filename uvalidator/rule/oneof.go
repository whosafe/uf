package rule

import (
	"fmt"
	"strings"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// OneOf 枚举值验证规�?
type OneOf struct {
	Allowed []string // 允许的值列�?
}

// Validate 执行验证
func (o *OneOf) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	for _, allowed := range o.Allowed {
		if str == allowed {
			return true
		}
	}
	return false
}

// GetMessage 获取错误消息
func (o *OneOf) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("oneof", lang...)

	// 替换占位�?
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", strings.Join(o.Allowed, ", "))

	return msg
}

// Name 规则名称
func (o *OneOf) Name() string {
	return "oneof"
}

// NewOneOf 创建枚举值验证规�?
func NewOneOf(allowed ...string) *OneOf {
	return &OneOf{Allowed: allowed}
}

// NewOneOfInt 创建整数枚举值验证规�?
type OneOfInt struct {
	Allowed []int
}

// Validate 执行验证
func (o *OneOfInt) Validate(value any) bool {
	num, ok := value.(int)
	if !ok {
		return false
	}

	for _, allowed := range o.Allowed {
		if num == allowed {
			return true
		}
	}
	return false
}

// GetMessage 获取错误消息
func (o *OneOfInt) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("oneof", lang...)

	// 转换为字符串
	var strs []string
	for _, v := range o.Allowed {
		strs = append(strs, fmt.Sprintf("%d", v))
	}

	// 替换占位�?
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", strings.Join(strs, ", "))

	return msg
}

// Name 规则名称
func (o *OneOfInt) Name() string {
	return "oneof"
}

// NewOneOfInt 创建整数枚举值验证规�?
func NewOneOfInt(allowed ...int) *OneOfInt {
	return &OneOfInt{Allowed: allowed}
}
