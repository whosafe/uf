package ulogger

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

// TestDefaultConfig 测试默认配置
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.File != "2006-01-02.log" {
		t.Errorf("Expected File to be '2006-01-02.log', got '%s'", config.File)
	}

	if config.Level != slog.LevelInfo {
		t.Errorf("Expected Level to be Info, got %v", config.Level)
	}

	if !config.UseStandardLogFormat {
		t.Error("Expected UseStandardLogFormat to be true")
	}

	if !config.Stdout {
		t.Error("Expected Stdout to be true")
	}
}

// TestBasicLogging 测试基本日志输出
func TestBasicLogging(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	config := &Config{
		Path:                 tmpDir,
		File:                 "test.log",
		Level:                slog.LevelDebug,
		UseStandardLogFormat: true,
		Stdout:               false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// 写入日志
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// 同步到磁盘
	if err := logger.Sync(); err != nil {
		t.Errorf("Failed to sync: %v", err)
	}

	// 检查文件是否存在
	logPath := filepath.Join(tmpDir, "test.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

// TestLogLevel 测试日志级别过滤
func TestLogLevel(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		Path:   tmpDir,
		File:   "level.log",
		Level:  slog.LevelWarn, // 只记录 Warn 及以上
		Stdout: false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Debug("should not appear")
	logger.Info("should not appear")
	logger.Warn("should appear")
	logger.Error("should appear")

	logger.Sync()

	// 读取日志文件
	logPath := filepath.Join(tmpDir, "level.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	logContent := string(content)
	if len(logContent) == 0 {
		t.Error("Log file is empty")
	}
}

// TestWithAttributes 测试带属性的日志
func TestWithAttributes(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		Path:   tmpDir,
		File:   "attrs.log",
		Stdout: false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("user login", "user", "alice", "ip", "192.168.1.1")
	logger.Error("database error", "error", "connection timeout", "retry", 3)

	logger.Sync()

	logPath := filepath.Join(tmpDir, "attrs.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Log file is empty")
	}
}

// TestPrefix 测试日志前缀
func TestPrefix(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		Path:   tmpDir,
		File:   "prefix.log",
		Prefix: "[APP]",
		Stdout: false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("test message")
	logger.Sync()

	logPath := filepath.Join(tmpDir, "prefix.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Log file is empty")
	}
}

// TestStdoutOnly 测试仅输出到终端
func TestStdoutOnly(t *testing.T) {
	config := &Config{
		Path:   "", // 不写文件
		Stdout: true,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("stdout only message")
}

// TestGlobalLogger 测试全局 Logger
func TestGlobalLogger(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		Path:   tmpDir,
		File:   "global.log",
		Stdout: false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	SetDefault(logger)

	Info("global info message")
	Warn("global warn message")

	logger.Sync()

	logPath := filepath.Join(tmpDir, "global.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Global log file was not created")
	}
}
