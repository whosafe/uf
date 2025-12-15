package rule

import (
	"regexp"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// 域名正则
var domainRegex = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

// Domain 域名验证规则
type Domain struct{}

// Validate 执行验证
func (d *Domain) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return domainRegex.MatchString(str)
}

// GetMessage 获取错误消息
func (d *Domain) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("domain", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (d *Domain) Name() string {
	return "domain"
}

// NewDomain 创建域名验证规则
func NewDomain() *Domain {
	return &Domain{}
}
