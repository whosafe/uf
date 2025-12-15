package ulogger

import (
	"log/slog"
	"strconv"
	"strings"
)

// TextFormatter 文本格式化器
type TextFormatter struct{}

// NewTextFormatter 创建文本格式化器
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{}
}

// Format 实现 Formatter 接口
func (f *TextFormatter) Format(r slog.Record, config *Config) ([]byte, error) {
	var buf strings.Builder

	// 添加前缀
	if config.Prefix != "" {
		buf.WriteString(config.Prefix)
		buf.WriteString(" ")
	}

	if config.UseStandardLogFormat {
		f.formatStandard(&buf, r, config)
	} else {
		f.formatSimple(&buf, r, config)
	}

	return []byte(buf.String()), nil
}

// formatStandard 标准格式化
func (f *TextFormatter) formatStandard(buf *strings.Builder, r slog.Record, config *Config) {
	// 时间
	buf.WriteString(r.Time.Format("2006-01-02 15:04:05"))
	buf.WriteString(" ")

	// 级别
	buf.WriteString("[")
	buf.WriteString(levelString(r.Level))
	buf.WriteString("]")
	buf.WriteString(" ")

	// 文件和行号
	file, line := getSourceLocation(r.PC, config.ShortFile)
	if file != "" {
		buf.WriteString(file)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(line))
		buf.WriteString(" ")
	}

	// 消息
	buf.WriteString(r.Message)

	// 属性
	r.Attrs(func(a slog.Attr) bool {
		buf.WriteString(" ")
		buf.WriteString(a.Key)
		buf.WriteString("=")
		buf.WriteString(a.Value.String())
		return true
	})
}

// formatSimple 简洁格式化
func (f *TextFormatter) formatSimple(buf *strings.Builder, r slog.Record, config *Config) {
	// 时间
	buf.WriteString(r.Time.Format("15:04:05"))
	buf.WriteString(" ")

	// 消息
	buf.WriteString(r.Message)

	// 属性
	r.Attrs(func(a slog.Attr) bool {
		buf.WriteString(" ")
		buf.WriteString(a.Key)
		buf.WriteString("=")
		buf.WriteString(a.Value.String())
		return true
	})
}
