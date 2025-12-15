package uvalidator

import (
	"fmt"
	"strings"
)

// FieldError 字段验证错误
type FieldError struct {
	Field   string      // 字段名
	Rule    string      // 规则名称
	Value   interface{} // 字段值
	Message string      // 错误消息
}

// Error 实现 error 接口
func (e *FieldError) Error() string {
	return e.Message
}

// ValidationErrors 多个验证错误
type ValidationErrors []FieldError

// Error 实现 error 接口
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}

	var buf strings.Builder
	for i, err := range ve {
		if i > 0 {
			buf.WriteString("; ")
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}

// HasErrors 是否有错误
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}

// First 获取第一个错误
func (ve ValidationErrors) First() *FieldError {
	if len(ve) > 0 {
		return &ve[0]
	}
	return nil
}

// ByField 按字段名获取错误
func (ve ValidationErrors) ByField(field string) []FieldError {
	var errs []FieldError
	for _, err := range ve {
		if err.Field == field {
			errs = append(errs, err)
		}
	}
	return errs
}

// NewFieldError 创建字段错误
func NewFieldError(field, rule string, value interface{}, message string) FieldError {
	return FieldError{
		Field:   field,
		Rule:    rule,
		Value:   value,
		Message: message,
	}
}

// Errorf 格式化错误消息
func Errorf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
