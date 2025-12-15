package rule

import (
	"iutime.com/utime/uf/uvalidator"
	"iutime.com/utime/uf/uvalidator/i18n"
)

// BankCard 银行卡号验证规则 (使用Luhn算法)
type BankCard struct{}

// Validate 执行验证
func (b *BankCard) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	// 银行卡号通常�?3-19位数�?
	length := len(str)
	if length < 13 || length > 19 {
		return false
	}

	// 检查是否全是数�?
	for _, r := range str {
		if r < '0' || r > '9' {
			return false
		}
	}

	// Luhn算法验证
	return luhnCheck(str)
}

// luhnCheck Luhn算法校验
func luhnCheck(cardNumber string) bool {
	sum := 0
	isDouble := false

	// 从右往左遍�?
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')

		if isDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isDouble = !isDouble
	}

	return sum%10 == 0
}

// GetMessage 获取错误消息
func (b *BankCard) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("bankcard", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (b *BankCard) Name() string {
	return "bankcard"
}

// NewBankCard 创建银行卡号验证规则
func NewBankCard() *BankCard {
	return &BankCard{}
}
