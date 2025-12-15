package rule

import (
	"regexp"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

var urlRegex = regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)

// URL URL 验证规则
type URL struct{}

// Validate 执行验证
func (u *URL) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true // 空值由 required 验证
	}

	return urlRegex.MatchString(str)
}

// GetMessage 获取错误消息
func (u *URL) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("url", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (u *URL) Name() string {
	return "url"
}

// NewURL 创建 URL 验证规则
func NewURL() *URL {
	return &URL{}
}
