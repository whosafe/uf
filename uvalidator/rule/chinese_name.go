package rule

import (
	"regexp"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// 中文姓名正则 (2-20个中文字字符支持·)
var chineseNameRegex = regexp.MustCompile(`^[\p{Han}·]{2,20}$`)

// ChineseName 中文姓名验证规则
type ChineseName struct{}

// Validate 执行验证
func (c *ChineseName) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return chineseNameRegex.MatchString(str)
}

// GetMessage 获取错误消息
func (c *ChineseName) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("chinese_name", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (c *ChineseName) Name() string {
	return "chinese_name"
}

// NewChineseName 创建中文姓名验证规则
func NewChineseName() *ChineseName {
	return &ChineseName{}
}
