package rule

import (
	"unicode"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// Alphanum 只包含字母和数字验证规则
type Alphanum struct{}

// Validate 执行验证
func (a *Alphanum) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// GetMessage 获取错误消息
func (a *Alphanum) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("alphanum", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (a *Alphanum) Name() string {
	return "alphanum"
}

// NewAlphanum 创建字母数字验证规则
func NewAlphanum() *Alphanum {
	return &Alphanum{}
}
