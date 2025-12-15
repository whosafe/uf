package rule

import (
	"fmt"
	"time"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// DateBefore 日期早于指定日期验证规则
type DateBefore struct {
	CompareDate string // 比较日期
	Format      string // 日期格式
}

// Validate 执行验证
func (db *DateBefore) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	format := db.Format
	if format == "" {
		format = DefaultDateFormat
	}

	date, err := time.Parse(format, str)
	if err != nil {
		return false
	}

	compareDate, err := time.Parse(format, db.CompareDate)
	if err != nil {
		return false
	}

	return date.Before(compareDate)
}

// GetMessage 获取错误消息
func (db *DateBefore) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("date_before", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", db.CompareDate)
	return msg
}

// Name 规则名称
func (db *DateBefore) Name() string {
	return "date_before"
}

// NewDateBefore 创建日期早于验证规则
func NewDateBefore(compareDate string, format ...string) *DateBefore {
	db := &DateBefore{CompareDate: compareDate}
	if len(format) > 0 {
		db.Format = format[0]
	}
	return db
}

// DateAfter 日期晚于指定日期验证规则
type DateAfter struct {
	CompareDate string // 比较日期
	Format      string // 日期格式
}

// Validate 执行验证
func (da *DateAfter) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	format := da.Format
	if format == "" {
		format = DefaultDateFormat
	}

	date, err := time.Parse(format, str)
	if err != nil {
		return false
	}

	compareDate, err := time.Parse(format, da.CompareDate)
	if err != nil {
		return false
	}

	return date.After(compareDate)
}

// GetMessage 获取错误消息
func (da *DateAfter) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("date_after", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", da.CompareDate)
	return msg
}

// Name 规则名称
func (da *DateAfter) Name() string {
	return "date_after"
}

// NewDateAfter 创建日期晚于验证规则
func NewDateAfter(compareDate string, format ...string) *DateAfter {
	da := &DateAfter{CompareDate: compareDate}
	if len(format) > 0 {
		da.Format = format[0]
	}
	return da
}

// DateBetween 日期在指定范围内验证规则
type DateBetween struct {
	StartDate string // 开始日�?
	EndDate   string // 结束日期
	Format    string // 日期格式
}

// Validate 执行验证
func (db *DateBetween) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	format := db.Format
	if format == "" {
		format = DefaultDateFormat
	}

	date, err := time.Parse(format, str)
	if err != nil {
		return false
	}

	startDate, err := time.Parse(format, db.StartDate)
	if err != nil {
		return false
	}

	endDate, err := time.Parse(format, db.EndDate)
	if err != nil {
		return false
	}

	return (date.After(startDate) || date.Equal(startDate)) &&
		(date.Before(endDate) || date.Equal(endDate))
}

// GetMessage 获取错误消息
func (db *DateBetween) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("date_between", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{min}", db.StartDate)
	msg = replaceAll(msg, "{max}", fmt.Sprintf("%s", db.EndDate))
	return msg
}

// Name 规则名称
func (db *DateBetween) Name() string {
	return "date_between"
}

// NewDateBetween 创建日期范围验证规则
func NewDateBetween(startDate, endDate string, format ...string) *DateBetween {
	db := &DateBetween{StartDate: startDate, EndDate: endDate}
	if len(format) > 0 {
		db.Format = format[0]
	}
	return db
}
