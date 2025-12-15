package rule

import (
	"regexp"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// 中国手机号正则（支持 13x, 14x, 15x, 16x, 17x, 18x, 19x
var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// Phone 手机号验证规则（中国手机号）
type Phone struct{}

// Validate 执行验证
func (p *Phone) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true // 空值由 required 验证
	}

	return phoneRegex.MatchString(str)
}

// GetMessage 获取错误消息
func (p *Phone) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("phone", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (p *Phone) Name() string {
	return "phone"
}

// NewPhone 创建手机号验证规则
func NewPhone() *Phone {
	return &Phone{}
}
