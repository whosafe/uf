package ulogger

import "strings"

// 敏感关键词列表
var sensitiveKeywords = []string{
	"password", "passwd", "pwd",
	"token", "secret", "key",
	"auth", "credential", "api_key",
	"private", "session", "cookie",
}

// isSensitiveKey 判断key是否为敏感字段
func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

// sanitizeString 脱敏字符串值
// 如果字符串中包含敏感关键词,返回脱敏后的字符串
func sanitizeString(s string) string {
	lower := strings.ToLower(s)
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lower, keyword) {
			return "***REDACTED***"
		}
	}
	return s
}

// sanitizeAttrValue 脱敏 slog.Value 类型的值
func sanitizeAttrValue(v interface{}) interface{} {
	// 如果是字符串,进行脱敏检查
	if str, ok := v.(string); ok {
		return sanitizeString(str)
	}
	// 其他类型直接返回
	return v
}
