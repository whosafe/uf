package rule

import (
	"fmt"
	"strings"

	"github.com/whosafe/uf/uvalidator"
	"github.com/whosafe/uf/uvalidator/i18n"
)

// In 值在列表中验证规则（字符串版本）
type In struct {
	Allowed []string // 允许的值列表
}

// Validate 执行验证
func (in *In) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	for _, allowed := range in.Allowed {
		if str == allowed {
			return true
		}
	}
	return false
}

// GetMessage 获取错误消息
func (in *In) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("in", lang...)

	// 替换占位符
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", strings.Join(in.Allowed, ", "))

	return msg
}

// Name 规则名称
func (in *In) Name() string {
	return "in"
}

// NewIn 创建 In 验证规则
func NewIn(allowed ...string) *In {
	return &In{Allowed: allowed}
}

// InInt 值在列表中验证规则（整数版本）
type InInt struct {
	Allowed []int
}

// Validate 执行验证
func (in *InInt) Validate(value any) bool {
	num, ok := value.(int)
	if !ok {
		return false
	}

	for _, allowed := range in.Allowed {
		if num == allowed {
			return true
		}
	}
	return false
}

// GetMessage 获取错误消息
func (in *InInt) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("in", lang...)

	// 转换为字符串
	var strs []string
	for _, v := range in.Allowed {
		strs = append(strs, fmt.Sprintf("%d", v))
	}

	// 替换占位符
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", strings.Join(strs, ", "))

	return msg
}

// Name 规则名称
func (in *InInt) Name() string {
	return "in"
}

// NewInInt 创建整数 In 验证规则
func NewInInt(allowed ...int) *InInt {
	return &InInt{Allowed: allowed}
}

// InFloat 值在列表中验证规则（浮点数版本）
type InFloat struct {
	Allowed []float64
}

// Validate 执行验证
func (in *InFloat) Validate(value any) bool {
	num, ok := value.(float64)
	if !ok {
		return false
	}

	for _, allowed := range in.Allowed {
		if num == allowed {
			return true
		}
	}
	return false
}

// GetMessage 获取错误消息
func (in *InFloat) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("in", lang...)

	// 转换为字符串
	var strs []string
	for _, v := range in.Allowed {
		strs = append(strs, fmt.Sprintf("%.2f", v))
	}

	// 替换占位符
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", strings.Join(strs, ", "))

	return msg
}

// Name 规则名称
func (in *InFloat) Name() string {
	return "in"
}

// NewInFloat 创建浮点数 In 验证规则
func NewInFloat(allowed ...float64) *InFloat {
	return &InFloat{Allowed: allowed}
}
