package main

import (
	"fmt"
	"log/slog"

	"github.com/whosafe/uf/ulogger"
)

func FormatDemo() {
	println("=== 日志格式演示 ===\n")

	// 1. 文本格式（标准）
	println("1. 文本格式（标准）")
	println("-------------------")
	textStandardConfig := &ulogger.Config{
		Format:               "text",
		UseStandardLogFormat: true,
		ShortFile:            true,
		Prefix:               "[APP]",
		Level:                slog.LevelDebug,
		Stdout:               true,
	}
	logger1, _ := ulogger.New(textStandardConfig)
	logger1.Info("用户登录", "user", "alice", "ip", "192.168.1.1")
	logger1.Warn("内存使用率高", "usage", "85%")
	logger1.Close()

	println()

	// 2. 文本格式（简洁）
	println("2. 文本格式（简洁）")
	println("-------------------")
	textSimpleConfig := &ulogger.Config{
		Format:               "text",
		UseStandardLogFormat: false,
		Level:                slog.LevelDebug,
		Stdout:               true,
	}
	logger2, _ := ulogger.New(textSimpleConfig)
	logger2.Info("任务开始", "task", "backup")
	logger2.Debug("处理中", "progress", "50%")
	logger2.Close()

	println()

	// 3. JSON 格式
	println("3. JSON 格式")
	println("-------------------")
	jsonConfig := &ulogger.Config{
		Format:    "json",
		ShortFile: true,
		Level:     slog.LevelDebug,
		Stdout:    true,
	}
	logger3, _ := ulogger.New(jsonConfig)
	logger3.Info("API 请求", "method", "GET", "path", "/api/users", "status", 200)
	logger3.Error("数据库错误", "error", "connection timeout", "retry", 3)
	logger3.Close()

	println()

	// 4. 自定义格式
	println("4. 自定义格式")
	println("-------------------")
	customConfig := &ulogger.Config{
		Format:    "custom",
		Formatter: &CustomFormatter{},
		Level:     slog.LevelDebug,
		Stdout:    true,
	}
	logger4, _ := ulogger.New(customConfig)
	logger4.Info("自定义消息", "key1", "value1", "key2", "value2")
	logger4.Warn("这是警告")
	logger4.Close()

	println()

	// 5. JSON 格式写入文件
	println("5. JSON 格式写入文件")
	println("-------------------")
	jsonFileConfig := &ulogger.Config{
		Path:      "./logs",
		File:      "app.json",
		Format:    "json",
		ShortFile: true,
		Level:     slog.LevelInfo,
		Stdout:    false,
	}
	logger5, _ := ulogger.New(jsonFileConfig)
	logger5.Info("写入 JSON 日志", "format", "json", "file", "./logs/app.json")
	logger5.Sync()
	logger5.Close()
	println("JSON 日志已写入 ./logs/app.json")
}

// CustomFormatter 自定义格式化器示例
type CustomFormatter struct{}

func (f *CustomFormatter) Format(r slog.Record, config *ulogger.Config) ([]byte, error) {
	// 自定义格式：[级别] 消息 (属性...)
	output := fmt.Sprintf("[%s] %s", r.Level, r.Message)

	// 添加属性
	r.Attrs(func(a slog.Attr) bool {
		output += fmt.Sprintf(" (%s=%v)", a.Key, a.Value.Any())
		return true
	})

	return []byte(output), nil
}
