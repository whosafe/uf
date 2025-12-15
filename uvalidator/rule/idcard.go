package rule

import (
	"regexp"
	"strconv"

	"iutime.com/utime/uf/uvalidator"

	"iutime.com/utime/uf/uvalidator/i18n"
)

// 身份证号正则
var idCard15Regex = regexp.MustCompile(`^[1-9]\d{7}((0\d)|(1[0-2]))(([0|1|2]\d)|3[0-1])\d{3}$`)
var idCard18Regex = regexp.MustCompile(`^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`)

// IDCard 身份证号验证规则 (15位或18�?
type IDCard struct{}

// Validate 执行验证
func (i *IDCard) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	length := len(str)
	if length == 15 {
		return idCard15Regex.MatchString(str)
	} else if length == 18 {
		if !idCard18Regex.MatchString(str) {
			return false
		}
		// 验证校验�?
		return validateIDCard18CheckCode(str)
	}

	return false
}

// validateIDCard18CheckCode 验证18位身份证校验�?
func validateIDCard18CheckCode(idCard string) bool {
	// 加权因子
	factor := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	// 校验码对照表
	checkCode := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

	sum := 0
	for i := 0; i < 17; i++ {
		num, _ := strconv.Atoi(string(idCard[i]))
		sum += num * factor[i]
	}

	mod := sum % 11
	expectedCheckCode := checkCode[mod]
	actualCheckCode := idCard[17]

	// 处理小写x
	if actualCheckCode == 'x' {
		actualCheckCode = 'X'
	}

	return actualCheckCode == expectedCheckCode
}

// GetMessage 获取错误消息
func (i *IDCard) GetMessage(field string, params map[string]string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("idcard", lang...)
	return replaceAll(template, "{field}", field)
}

// Name 规则名称
func (i *IDCard) Name() string {
	return "idcard"
}

// NewIDCard 创建身份证号验证规则
func NewIDCard() *IDCard {
	return &IDCard{}
}
