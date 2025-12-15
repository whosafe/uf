package rule

import (
	"regexp"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// 邮政编码正则 (6位数�?
var postalCodeRegex = regexp.MustCompile(`^\d{6}$`)

// PostalCode 邮政编码验证规则 (中国6�?
type PostalCode struct{}

// Validate 执行验证
func (p *PostalCode) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return postalCodeRegex.MatchString(str)
}

// GetMessage 获取错误消息
func (p *PostalCode) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("postalcode", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (p *PostalCode) Name() string {
	return "postalcode"
}

// NewPostalCode 创建邮政编码验证规则
func NewPostalCode() *PostalCode {
	return &PostalCode{}
}
