package i18n

import "iutime.com/utime/uf/uvalidator"

// GetMessage 获取错误消息
// lang 参数可选,如果不传则使用全局语言设置
// 并发安全：en 和 zh map 为只读，不会在运行时修改
func GetMessage(key string, lang ...uvalidator.Language) string {
	var currentLang uvalidator.Language
	if len(lang) > 0 {
		currentLang = lang[0]
	} else {
		currentLang = uvalidator.GetLanguage()
	}

	var messages map[string]string
	switch currentLang {
	case uvalidator.LanguageZH:
		messages = zh
	default:
		messages = en
	}

	if msg, ok := messages[key]; ok {
		return msg
	}

	// 回退到英文
	if msg, ok := en[key]; ok {
		return msg
	}

	return "{field} validation failed"
}
