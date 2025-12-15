package ulogger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

// TestTextFormatter 测试文本格式化器
func TestTextFormatter(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		Path:                 tmpDir,
		File:                 "text.log",
		Format:               "text",
		UseStandardLogFormat: true,
		ShortFile:            true,
		Level:                slog.LevelDebug,
		Stdout:               false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("test message", "key", "value")
	logger.Sync()

	// 验证文件存在
	logPath := filepath.Join(tmpDir, "text.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file was not created")
	}
}

// TestJSONFormatter 测试 JSON 格式化器
func TestJSONFormatter(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		Path:      tmpDir,
		File:      "json.log",
		Format:    "json",
		ShortFile: true,
		Level:     slog.LevelDebug,
		Stdout:    false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("test message", "key", "value", "number", 42)
	logger.Sync()

	// 验证文件存在
	logPath := filepath.Join(tmpDir, "json.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// 验证是 JSON 格式（包含大括号）
	if len(content) == 0 {
		t.Error("Log file is empty")
	}
	if content[0] != '{' {
		t.Error("Log output is not JSON format")
	}
}

// TestCustomFormatter 测试自定义格式化器
func TestCustomFormatter(t *testing.T) {
	tmpDir := t.TempDir()

	// 自定义格式化器
	customFormatter := &testFormatter{}

	config := &Config{
		Path:      tmpDir,
		File:      "custom.log",
		Format:    "custom",
		Formatter: customFormatter,
		Level:     slog.LevelDebug,
		Stdout:    false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("test message")
	logger.Sync()

	// 验证文件存在
	logPath := filepath.Join(tmpDir, "custom.log")
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	// 验证使用了自定义格式
	expected := "CUSTOM: test message"
	if string(content[:len(expected)]) != expected {
		t.Errorf("Expected custom format, got: %s", string(content))
	}
}

// testFormatter 测试用的自定义格式化器
type testFormatter struct{}

func (f *testFormatter) Format(r slog.Record, config *Config) ([]byte, error) {
	return []byte(fmt.Sprintf("CUSTOM: %s", r.Message)), nil
}

// TestFormatSwitch 测试格式切换
func TestFormatSwitch(t *testing.T) {
	tmpDir := t.TempDir()

	// 测试文本格式
	textConfig := &Config{
		Path:   tmpDir,
		File:   "switch-text.log",
		Format: "text",
		Stdout: false,
	}
	textLogger, _ := New(textConfig)
	textLogger.Info("text format")
	textLogger.Close()

	// 测试 JSON 格式
	jsonConfig := &Config{
		Path:   tmpDir,
		File:   "switch-json.log",
		Format: "json",
		Stdout: false,
	}
	jsonLogger, _ := New(jsonConfig)
	jsonLogger.Info("json format")
	jsonLogger.Close()

	// 验证两个文件都存在
	textPath := filepath.Join(tmpDir, "switch-text.log")
	jsonPath := filepath.Join(tmpDir, "switch-json.log")

	if _, err := os.Stat(textPath); os.IsNotExist(err) {
		t.Error("Text log file was not created")
	}

	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		t.Error("JSON log file was not created")
	}
}
