package rule

import (
	"encoding/base64"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// Base64 Base64编码验证规则
type Base64 struct{}

// Validate 执行验证
func (b *Base64) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	_, err := base64.StdEncoding.DecodeString(str)
	return err == nil
}

// GetMessage 获取错误消息
func (b *Base64) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("base64", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (b *Base64) Name() string {
	return "base64"
}

// NewBase64 创建Base64验证规则
func NewBase64() *Base64 {
	return &Base64{}
}
