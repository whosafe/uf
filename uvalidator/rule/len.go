package rule

import (
	"fmt"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// Len 固定长度验证规则
type Len struct {
	Length int // 固定长度
}

// Validate 执行验证
func (l *Len) Validate(value any) bool {
	switch v := value.(type) {
	case string:
		return len(v) == l.Length
	case []any:
		return len(v) == l.Length
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (l *Len) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("len", lang...)

	// 替换占位符
	msg := template
	msg = replaceAll(msg, "{field}", field)
	msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", l.Length))

	return msg
}

// Name 规则名称
func (l *Len) Name() string {
	return "len"
}

// NewLen 创建固定长度验证规则
func NewLen(length int) *Len {
	return &Len{Length: length}
}
