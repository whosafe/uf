package rule

import (
	"fmt"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// Unique 数组元素唯一性验证规则
type Unique struct{}

// Validate 执行验证
func (u *Unique) Validate(value any) bool {
	switch v := value.(type) {
	case []string:
		seen := make(map[string]bool)
		for _, item := range v {
			if seen[item] {
				return false
			}
			seen[item] = true
		}
		return true
	case []int:
		seen := make(map[int]bool)
		for _, item := range v {
			if seen[item] {
				return false
			}
			seen[item] = true
		}
		return true
	case []any:
		seen := make(map[any]bool)
		for _, item := range v {
			if seen[item] {
				return false
			}
			seen[item] = true
		}
		return true
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (u *Unique) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("unique", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (u *Unique) Name() string {
	return "unique"
}

// NewUnique 创建数组唯一性验证规则
func NewUnique() *Unique {
	return &Unique{}
}

// ArrayMin 数组最小长度验证规则
type ArrayMin struct {
	MinLength int
}

// Validate 执行验证
func (a *ArrayMin) Validate(value any) bool {
	switch v := value.(type) {
	case []string:
		return len(v) >= a.MinLength
	case []int:
		return len(v) >= a.MinLength
	case []any:
		return len(v) >= a.MinLength
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (a *ArrayMin) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("array_min", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", a.MinLength))
	return msg
}

// Name 规则名称
func (a *ArrayMin) Name() string {
	return "array_min"
}

// NewArrayMin 创建数组最小长度验证规则
func NewArrayMin(minLength int) *ArrayMin {
	return &ArrayMin{MinLength: minLength}
}

// ArrayMax 数组最大长度验证规则
type ArrayMax struct {
	MaxLength int
}

// Validate 执行验证
func (a *ArrayMax) Validate(value any) bool {
	switch v := value.(type) {
	case []string:
		return len(v) <= a.MaxLength
	case []int:
		return len(v) <= a.MaxLength
	case []any:
		return len(v) <= a.MaxLength
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (a *ArrayMax) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("array_max", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", a.MaxLength))
	return msg
}

// Name 规则名称
func (a *ArrayMax) Name() string {
	return "array_max"
}

// NewArrayMax 创建数组最大长度验证规则
func NewArrayMax(maxLength int) *ArrayMax {
	return &ArrayMax{MaxLength: maxLength}
}

// ArrayContains 数组包含特定值验证规则
type ArrayContains struct {
	ContainsValue any
}

// Validate 执行验证
func (a *ArrayContains) Validate(value any) bool {
	switch v := value.(type) {
	case []string:
		target, ok := a.ContainsValue.(string)
		if !ok {
			return false
		}
		for _, item := range v {
			if item == target {
				return true
			}
		}
		return false
	case []int:
		target, ok := a.ContainsValue.(int)
		if !ok {
			return false
		}
		for _, item := range v {
			if item == target {
				return true
			}
		}
		return false
	case []any:
		for _, item := range v {
			if item == a.ContainsValue {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (a *ArrayContains) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("array_contains", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%v", a.ContainsValue))
	return msg
}

// Name 规则名称
func (a *ArrayContains) Name() string {
	return "array_contains"
}

// NewArrayContains 创建数组包含验证规则
func NewArrayContains(containsValue any) *ArrayContains {
	return &ArrayContains{ContainsValue: containsValue}
}
