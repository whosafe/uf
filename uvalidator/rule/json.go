package rule

import (
	"encoding/json"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// JSON JSON格式验证规则
type JSON struct{}

// Validate 执行验证
func (j *JSON) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	return json.Valid([]byte(str))
}

// GetMessage 获取错误消息
func (j *JSON) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("json", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (j *JSON) Name() string {
	return "json"
}

// NewJSON 创建JSON验证规则
func NewJSON() *JSON {
	return &JSON{}
}
