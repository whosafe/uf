package ulogger

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestRotateBySize 测试按大小轮转
func TestRotateBySize(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		Path:       tmpDir,
		File:       "rotate.log",
		RotateSize: 100, // 100 字节
		Stdout:     false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// 写入足够多的日志触发轮转
	for i := 0; i < 20; i++ {
		logger.Info("This is a test message that will trigger rotation")
	}

	logger.Sync()
	time.Sleep(100 * time.Millisecond) // 等待轮转完成

	// 检查是否生成了备份文件
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(files) < 2 {
		t.Errorf("Expected at least 2 files (current + backup), got %d", len(files))
	}
}

// TestRotateByTime 测试按时间轮转
func TestRotateByTime(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		Path:         tmpDir,
		File:         "time-rotate.log",
		RotateExpire: 2, // 2 秒
		Stdout:       false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	logger.Info("First message")
	logger.Sync()

	// 等待超过轮转时间
	time.Sleep(3 * time.Second)

	logger.Info("Second message after rotation")
	logger.Sync()

	time.Sleep(100 * time.Millisecond)

	// 检查是否生成了备份文件
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	if len(files) < 2 {
		t.Errorf("Expected at least 2 files, got %d", len(files))
	}
}

// TestBackupLimit 测试备份数量限制
func TestBackupLimit(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		Path:              tmpDir,
		File:              "backup.log",
		RotateSize:        50,
		RotateBackupLimit: 2, // 只保留 2 个备份
		Stdout:            false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// 写入大量日志触发多次轮转
	for i := 0; i < 50; i++ {
		logger.Info("Test message for backup limit testing")
		time.Sleep(10 * time.Millisecond)
	}

	logger.Sync()
	time.Sleep(500 * time.Millisecond) // 等待清理完成

	// 检查文件数量
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	// 应该有当前文件 + 最多 2 个备份 = 最多 3 个文件
	if len(files) > 3 {
		t.Errorf("Expected at most 3 files, got %d", len(files))
	}
}

// TestCompression 测试压缩功能
func TestCompression(t *testing.T) {
	tmpDir := t.TempDir()

	config := &Config{
		Path:                 tmpDir,
		File:                 "compress.log",
		RotateSize:           100,
		RotateBackupCompress: 6, // gzip 压缩级别 6
		Stdout:               false,
	}

	logger, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Close()

	// 写入足够多的日志触发轮转
	for i := 0; i < 30; i++ {
		logger.Info("This is a test message for compression testing")
	}

	logger.Sync()
	time.Sleep(500 * time.Millisecond) // 等待压缩完成

	// 检查是否有 .gz 文件
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	hasGzFile := false
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".gz" {
			hasGzFile = true
			break
		}
	}

	if !hasGzFile {
		t.Error("Expected to find .gz compressed file")
	}
}
