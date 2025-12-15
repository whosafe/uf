package ulogger

import (
	"encoding/json"
	"log/slog"
	"strconv"
)

// JSONFormatter JSON 格式化器
type JSONFormatter struct{}

// NewJSONFormatter 创建 JSON 格式化器
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format 实现 Formatter 接口
func (f *JSONFormatter) Format(r slog.Record, config *Config) ([]byte, error) {
	// 构建 JSON 对象
	logEntry := make(map[string]interface{})

	// 时间
	logEntry["time"] = r.Time.Format("2006-01-02T15:04:05.000Z07:00")

	// 级别
	logEntry["level"] = levelString(r.Level)

	// 消息
	logEntry["msg"] = r.Message

	// 文件和行号
	file, line := getSourceLocation(r.PC, config.ShortFile)
	if file != "" {
		logEntry["source"] = file + ":" + strconv.Itoa(line)
	}

	// 属性
	attrs := make(map[string]interface{})
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})
	if len(attrs) > 0 {
		logEntry["attrs"] = attrs
	}

	// 序列化为 JSON
	data, err := json.Marshal(logEntry)
	if err != nil {
		return nil, err
	}

	return data, nil
}
