package rule

import (
	"regexp"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// 统一社会信用代码正则 (18为)
var socialCreditRegex = regexp.MustCompile(`^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$`)

// UnifiedSocialCreditCode 统一社会信用代码验证规则
type UnifiedSocialCreditCode struct{}

// Validate 执行验证
func (u *UnifiedSocialCreditCode) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	if len(str) != 18 {
		return false
	}

	return socialCreditRegex.MatchString(str)
}

// GetMessage 获取错误消息
func (u *UnifiedSocialCreditCode) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("social_credit", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (u *UnifiedSocialCreditCode) Name() string {
	return "social_credit"
}

// NewUnifiedSocialCreditCode 创建统一社会信用代码验证规则
func NewUnifiedSocialCreditCode() *UnifiedSocialCreditCode {
	return &UnifiedSocialCreditCode{}
}
