package ulogger

import (
	"compress/gzip"
	"io"
	"os"
	"time"
)

// backupManager 备份文件管理器
type backupManager struct {
	config *Config
}

// newBackupManager 创建备份管理器
func newBackupManager(config *Config) *backupManager {
	return &backupManager{
		config: config,
	}
}

// cleanup 清理备份文件
func (bm *backupManager) cleanup() error {
	if !bm.config.IsBackupEnabled() {
		return nil
	}

	baseName := time.Now().Format(bm.config.File)
	backups, err := listBackupFiles(bm.config.Path, baseName)
	if err != nil {
		return err
	}

	// 按数量限制清理
	if bm.config.RotateBackupLimit > 0 {
		if err := bm.cleanupByLimit(backups); err != nil {
			return err
		}
	}

	// 按过期时间清理
	if bm.config.RotateBackupExpire > 0 {
		if err := bm.cleanupByExpire(backups); err != nil {
			return err
		}
	}

	return nil
}

// cleanupByLimit 按数量限制清理备份文件
func (bm *backupManager) cleanupByLimit(backups []string) error {
	if len(backups) <= bm.config.RotateBackupLimit {
		return nil
	}

	// 删除最旧的文件（保留最新的 N 个）
	toDelete := len(backups) - bm.config.RotateBackupLimit
	for i := 0; i < toDelete; i++ {
		if err := os.Remove(backups[i]); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

// cleanupByExpire 按过期时间清理备份文件
func (bm *backupManager) cleanupByExpire(backups []string) error {
	expireDuration := time.Duration(bm.config.RotateBackupExpire) * time.Second
	now := time.Now()

	for _, backup := range backups {
		modTime, err := getFileModTime(backup)
		if err != nil {
			continue
		}

		if now.Sub(modTime) > expireDuration {
			if err := os.Remove(backup); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}

// compressFile 压缩文件
func (bm *backupManager) compressFile(srcPath string) error {
	if bm.config.RotateBackupCompress == 0 {
		return nil
	}

	// 打开源文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建压缩文件
	dstPath := srcPath + ".gz"
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 创建 gzip writer
	gzWriter, err := gzip.NewWriterLevel(dstFile, int(bm.config.RotateBackupCompress))
	if err != nil {
		return err
	}
	defer gzWriter.Close()

	// 复制数据
	if _, err := io.Copy(gzWriter, srcFile); err != nil {
		return err
	}

	// 关闭 gzip writer 以刷新缓冲区
	if err := gzWriter.Close(); err != nil {
		return err
	}

	// 删除原文件
	return os.Remove(srcPath)
}
