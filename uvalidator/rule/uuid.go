package rule

import (
	"regexp"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// UUID v4 格式正则
var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// UUID UUID格式验证规则 (v4)
type UUID struct{}

// Validate 执行验证
func (u *UUID) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return uuidRegex.MatchString(str)
}

// GetMessage 获取错误消息
func (u *UUID) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("uuid", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (u *UUID) Name() string {
	return "uuid"
}

// NewUUID 创建UUID验证规则
func NewUUID() *UUID {
	return &UUID{}
}
