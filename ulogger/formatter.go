package ulogger

import (
	"log/slog"
	"runtime"
)

// Formatter 日志格式化器接口
type Formatter interface {
	// Format 格式化日志记录
	// 返回格式化后的字节数组，不包含换行符
	Format(r slog.Record, config *Config) ([]byte, error)
}

// FormatType 格式类型
type FormatType string

const (
	FormatText   FormatType = "text"   // 文本格式
	FormatJSON   FormatType = "json"   // JSON 格式
	FormatCustom FormatType = "custom" // 自定义格式
)

// getSourceLocation 获取源代码位置信息
func getSourceLocation(pc uintptr, shortFile bool) (file string, line int) {
	if pc == 0 {
		return "", 0
	}

	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	file = f.File
	if shortFile {
		file = shortFilePath(file)
	}
	line = f.Line
	return
}

// levelString 获取级别字符串
func levelString(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARN"
	case slog.LevelError:
		return "ERROR"
	default:
		return "LEVEL(" + level.String() + ")"
	}
}
