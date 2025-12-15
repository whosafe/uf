package ulogger

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// getFileSize 获取文件大小
func getFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// ensureDir 确保目录存在
func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// formatFileName 格式化文件名，添加时间戳
func formatFileName(baseName string, t time.Time) string {
	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)
	timestamp := t.Format("20060102-150405")
	return nameWithoutExt + "." + timestamp + ext
}

// listBackupFiles 列出所有备份文件
func listBackupFiles(dir, baseName string) ([]string, error) {
	if dir == "" {
		return nil, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)
	prefix := nameWithoutExt + "."

	var backups []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// 匹配备份文件：原名.时间戳.log 或 原名.时间戳.log.gz
		if strings.HasPrefix(name, prefix) && name != baseName {
			backups = append(backups, filepath.Join(dir, name))
		}
	}

	// 按文件名排序（时间戳排序）
	sort.Strings(backups)
	return backups, nil
}

// getFileModTime 获取文件修改时间
func getFileModTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// shortFilePath 获取短文件路径（只保留文件名）
func shortFilePath(fullPath string) string {
	return filepath.Base(fullPath)
}
