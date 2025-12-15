package rule

import (
	"time"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

const (
	// DefaultDateFormat 默认日期格式
	DefaultDateFormat = "2006-01-02"
	// DefaultDateTimeFormat 默认日期时间格式
	DefaultDateTimeFormat = "2006-01-02 15:04:05"
)

// Date 日期格式验证规则
type Date struct {
	Format string // 日期格式,为空则使用默认格式
}

// Validate 执行验证
func (d *Date) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	format := d.Format
	if format == "" {
		format = DefaultDateFormat
	}

	_, err := time.Parse(format, str)
	return err == nil
}

// GetMessage 获取错误消息
func (d *Date) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("date", lang...)
	msg := replaceAll(template, "{field}", field)
	format := d.Format
	if format == "" {
		format = DefaultDateFormat
	}
	msg = replaceAll(msg, "{param}", format)
	return msg
}

// Name 规则名称
func (d *Date) Name() string {
	return "date"
}

// NewDate 创建日期验证规则
func NewDate(format ...string) *Date {
	d := &Date{}
	if len(format) > 0 {
		d.Format = format[0]
	}
	return d
}

// DateTime 日期时间格式验证规则
type DateTime struct {
	Format string // 日期时间格式,为空则使用默认格式
}

// Validate 执行验证
func (dt *DateTime) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	format := dt.Format
	if format == "" {
		format = DefaultDateTimeFormat
	}

	_, err := time.Parse(format, str)
	return err == nil
}

// GetMessage 获取错误消息
func (dt *DateTime) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("datetime", lang...)
	msg := replaceAll(template, "{field}", field)
	format := dt.Format
	if format == "" {
		format = DefaultDateTimeFormat
	}
	msg = replaceAll(msg, "{param}", format)
	return msg
}

// Name 规则名称
func (dt *DateTime) Name() string {
	return "datetime"
}

// NewDateTime 创建日期时间验证规则
func NewDateTime(format ...string) *DateTime {
	dt := &DateTime{}
	if len(format) > 0 {
		dt.Format = format[0]
	}
	return dt
}
