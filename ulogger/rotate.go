package ulogger

import (
	"io"
	"os"
	"sync"
	"time"
)

// rotateWriter 支持文件轮转的 Writer
type rotateWriter struct {
	config        *Config
	file          *os.File
	currentPath   string
	currentSize   int64
	lastRotate    time.Time
	mu            sync.Mutex
	backupManager *backupManager
	stopChan      chan struct{}
	wg            sync.WaitGroup
}

// newRotateWriter 创建轮转 Writer
func newRotateWriter(config *Config) (*rotateWriter, error) {
	rw := &rotateWriter{
		config:        config,
		lastRotate:    time.Now(),
		backupManager: newBackupManager(config),
		stopChan:      make(chan struct{}),
	}

	// 打开初始文件
	if err := rw.openFile(); err != nil {
		return nil, err
	}

	// 如果启用了时间轮转，启动定时检查
	if config.RotateExpire > 0 {
		rw.wg.Add(1)
		go rw.rotateByTimeLoop()
	}

	return rw, nil
}

// Write 实现 io.Writer 接口
func (rw *rotateWriter) Write(p []byte) (n int, err error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	// 检查是否需要按时间轮转
	if rw.config.RotateExpire > 0 && time.Since(rw.lastRotate) >= time.Duration(rw.config.RotateExpire)*time.Second {
		if err := rw.rotate(); err != nil {
			return 0, err
		}
	}

	// 检查是否需要按大小轮转
	if rw.config.RotateSize > 0 && rw.currentSize+int64(len(p)) > int64(rw.config.RotateSize) {
		if err := rw.rotate(); err != nil {
			return 0, err
		}
	}

	// 写入数据
	n, err = rw.file.Write(p)
	rw.currentSize += int64(n)
	return n, err
}

// Close 关闭 Writer
func (rw *rotateWriter) Close() error {
	close(rw.stopChan)
	rw.wg.Wait()

	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.file != nil {
		return rw.file.Close()
	}
	return nil
}

// openFile 打开日志文件
func (rw *rotateWriter) openFile() error {
	path := rw.config.GetLogFilePath()
	if path == "" {
		return nil
	}

	// 确保目录存在
	if err := ensureDir(rw.config.Path); err != nil {
		return err
	}

	// 打开或创建文件
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// 获取当前文件大小
	info, err := file.Stat()
	if err != nil {
		file.Close()
		return err
	}

	rw.file = file
	rw.currentPath = path
	rw.currentSize = info.Size()
	return nil
}

// rotate 执行文件轮转
func (rw *rotateWriter) rotate() error {
	// 关闭当前文件
	if rw.file != nil {
		if err := rw.file.Close(); err != nil {
			return err
		}
	}

	// 重命名当前文件为备份文件
	if rw.currentPath != "" {
		backupPath := formatFileName(rw.currentPath, time.Now())
		if err := os.Rename(rw.currentPath, backupPath); err != nil && !os.IsNotExist(err) {
			return err
		}

		// 压缩备份文件
		if rw.config.RotateBackupCompress > 0 {
			go rw.backupManager.compressFile(backupPath)
		}
	}

	// 打开新文件
	if err := rw.openFile(); err != nil {
		return err
	}

	// 清理旧备份
	go rw.backupManager.cleanup()

	rw.lastRotate = time.Now()
	return nil
}

// rotateByTimeLoop 定时检查是否需要轮转
func (rw *rotateWriter) rotateByTimeLoop() {
	defer rw.wg.Done()

	ticker := time.NewTicker(time.Second * 10) // 每 10 秒检查一次
	defer ticker.Stop()

	for {
		select {
		case <-rw.stopChan:
			return
		case <-ticker.C:
			rw.mu.Lock()
			if time.Since(rw.lastRotate) >= time.Duration(rw.config.RotateExpire)*time.Second {
				rw.rotate()
			}
			rw.mu.Unlock()
		}
	}
}

// Sync 同步文件
func (rw *rotateWriter) Sync() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.file != nil {
		return rw.file.Sync()
	}
	return nil
}

// multiWriter 多输出 Writer
type multiWriter struct {
	writers []io.Writer
}

// newMultiWriter 创建多输出 Writer
func newMultiWriter(writers ...io.Writer) *multiWriter {
	return &multiWriter{
		writers: writers,
	}
}

// Write 实现 io.Writer 接口
func (mw *multiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writers {
		n, err = w.Write(p)
		if err != nil {
			return
		}
	}
	return len(p), nil
}
