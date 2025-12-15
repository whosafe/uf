package uvalidator

import (
	"strings"
	"sync"
)

// Language 语言类型
type Language string

const (
	LanguageEN Language = "en" // 英语
	LanguageZH Language = "zh" // 中文
)

var (
	currentLanguage = LanguageEN
	languageMu      sync.RWMutex
)

// SetLanguage 设置当前语言
func SetLanguage(lang Language) {
	languageMu.Lock()
	defer languageMu.Unlock()
	currentLanguage = lang
}

// GetLanguage 获取当前语言
func GetLanguage() Language {
	languageMu.RLock()
	defer languageMu.RUnlock()
	return currentLanguage
}

// ParseAcceptLanguage 从 Accept-Language 请求头解析语言
// 例如: "zh-CN,zh;q=0.9,en;q=0.8" -> LanguageZH
func ParseAcceptLanguage(acceptLang string) Language {
	if acceptLang == "" {
		return GetLanguage()
	}

	// 简单解析 Accept-Language,取第一个语言代码
	parts := strings.Split(acceptLang, ",")
	if len(parts) > 0 {
		// 取第一个语言,去除权重等参数
		firstLang := strings.TrimSpace(parts[0])
		// 去除 ;q=0.9 等权重参数
		if idx := strings.Index(firstLang, ";"); idx > 0 {
			firstLang = firstLang[:idx]
		}

		// 提取语言代码(前两位)
		if len(firstLang) >= 2 {
			lang := strings.ToLower(firstLang[:2])
			switch lang {
			case "zh":
				return LanguageZH
			case "en":
				return LanguageEN
			}
		}
	}

	return GetLanguage()
}
