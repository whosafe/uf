package udb

import (
	"fmt"
	"strings"
)

// 敏感关键词列表
var sensitiveKeywords = []string{
	"password", "passwd", "pwd",
	"token", "secret", "key",
	"auth", "credential", "api_key",
	"private", "session", "cookie",
}

// SanitizeArgs 脱敏参数列表
// 【安全修复】防止敏感信息泄露到日志中
func SanitizeArgs(args []any) []any {
	if len(args) == 0 {
		return args
	}

	sanitized := make([]any, len(args))
	for i, arg := range args {
		sanitized[i] = sanitizeValue(arg)
	}
	return sanitized
}

// sanitizeValue 脱敏单个值
func sanitizeValue(value any) any {
	switch v := value.(type) {
	case string:
		// 检查字符串是否包含敏感关键词
		if containsSensitiveKeyword(v) {
			return "***REDACTED***"
		}
		return v
	case map[string]any:
		// 递归处理 map
		return sanitizeMap(v)
	case []any:
		// 递归处理 slice
		return SanitizeArgs(v)
	default:
		return v
	}
}

// sanitizeMap 脱敏 map
func sanitizeMap(m map[string]any) map[string]any {
	sanitized := make(map[string]any)
	for k, v := range m {
		// 检查 key 是否是敏感字段
		if isSensitiveKey(k) {
			sanitized[k] = "***REDACTED***"
		} else {
			sanitized[k] = sanitizeValue(v)
		}
	}
	return sanitized
}

// containsSensitiveKeyword 检查字符串是否包含敏感关键词
func containsSensitiveKeyword(s string) bool {
	lower := strings.ToLower(s)
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

// isSensitiveKey 检查 key 是否是敏感字段
func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

// SanitizeCommand 脱敏命令字符串(用于 Redis 等)
func SanitizeCommand(cmd string, args []any) string {
	// 如果命令本身包含敏感操作,脱敏所有参数
	cmdLower := strings.ToLower(cmd)
	if strings.Contains(cmdLower, "auth") ||
		strings.Contains(cmdLower, "password") ||
		strings.Contains(cmdLower, "secret") {
		return fmt.Sprintf("%s ***REDACTED***", cmd)
	}

	// 否则只脱敏敏感参数
	sanitizedArgs := SanitizeArgs(args)
	return fmt.Sprintf("%s %v", cmd, sanitizedArgs)
}
