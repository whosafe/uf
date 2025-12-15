package rule

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// Integer 整数验证规则
type Integer struct{}

// Validate 执行验证
func (i *Integer) Validate(value any) bool {
	switch v := value.(type) {
	case int, int64:
		return true
	case float64:
		return v == math.Floor(v)
	case string:
		if v == "" {
			return true
		}
		_, err := strconv.Atoi(v)
		return err == nil
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (i *Integer) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("integer", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (i *Integer) Name() string {
	return "integer"
}

// NewInteger 创建整数验证规则
func NewInteger() *Integer {
	return &Integer{}
}

// Decimal 小数验证规则
type Decimal struct {
	DecimalPlaces int // 小数位数,0表示不限至
}

// Validate 执行验证
func (d *Decimal) Validate(value any) bool {
	switch v := value.(type) {
	case float64:
		if d.DecimalPlaces == 0 {
			return true
		}
		// 检查小数位
		str := fmt.Sprintf("%f", v)
		parts := strings.Split(str, ".")
		if len(parts) != 2 {
			return false
		}
		// 去除尾部
		decimal := strings.TrimRight(parts[1], "0")
		return len(decimal) <= d.DecimalPlaces
	case string:
		if v == "" {
			return true
		}
		_, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return false
		}
		if d.DecimalPlaces == 0 {
			return true
		}
		parts := strings.Split(v, ".")
		if len(parts) != 2 {
			return false
		}
		decimal := strings.TrimRight(parts[1], "0")
		return len(decimal) <= d.DecimalPlaces
	default:
		return false
	}
}

// GetMessage 获取错误消息
func (d *Decimal) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("decimal", lang...)
	msg := replaceAll(template, "{field}", field)
	if d.DecimalPlaces > 0 {
		msg = replaceAll(msg, "{param}", fmt.Sprintf("%d", d.DecimalPlaces))
	}
	return msg
}

// Name 规则名称
func (d *Decimal) Name() string {
	return "decimal"
}

// NewDecimal 创建小数验证规则
func NewDecimal(decimalPlaces int) *Decimal {
	return &Decimal{DecimalPlaces: decimalPlaces}
}
