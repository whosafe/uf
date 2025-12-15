package rule

import (
	"regexp"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// Regex 正则表达式验证规则
type Regex struct {
	Pattern string
	regex   *regexp.Regexp
}

// Validate 执行验证
func (r *Regex) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	// 延迟编译正则表达
	if r.regex == nil {
		var err error
		r.regex, err = regexp.Compile(r.Pattern)
		if err != nil {
			return false
		}
	}

	return r.regex.MatchString(str)
}

// GetMessage 获取错误消息
func (r *Regex) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("regex", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", r.Pattern)
	return msg
}

// Name 规则名称
func (r *Regex) Name() string {
	return "regex"
}

// NewRegex 创建正则表达式验证规则
func NewRegex(pattern string) *Regex {
	return &Regex{Pattern: pattern}
}
