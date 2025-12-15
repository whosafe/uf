package main

import (
	"log/slog"
	"time"

	"iutime.com/utime/uf/ulogger"
)

func main() {
	// 示例 1: 基本用法 - 仅输出到终端
	basicExample()

	// 示例 2: 输出到文件
	fileExample()

	// 示例 3: 文件轮转（按大小）
	rotateBySizeExample()

	// 示例 4: 文件轮转（按时间）
	rotateByTimeExample()

	// 示例 5: 带压缩的备份
	compressedBackupExample()

	// 示例 6: 全局 Logger
	globalLoggerExample()
}

// 示例 1: 基本用法
func basicExample() {
	println("\n=== 示例 1: 基本用法 ===")

	config := &ulogger.Config{
		Prefix: "[APP]",
		Level:  slog.LevelDebug,
		Stdout: true,
	}

	logger, err := ulogger.New(config)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	logger.Debug("这是调试信息")
	logger.Info("这是普通信息")
	logger.Warn("这是警告信息")
	logger.Error("这是错误信息")

	// 带属性的日志
	logger.Info("用户登录", "user", "alice", "ip", "192.168.1.1")
}

// 示例 2: 输出到文件
func fileExample() {
	println("\n=== 示例 2: 输出到文件 ===")

	config := &ulogger.Config{
		Path:                 "./logs",
		File:                 "app.log",
		Level:                slog.LevelInfo,
		UseStandardLogFormat: true,
		ShortFile:            true,
		Stdout:               true, // 同时输出到终端
	}

	logger, err := ulogger.New(config)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	logger.Info("日志已写入文件", "path", "./logs/app.log")
	logger.Warn("这条日志同时输出到文件和终端")

	// 确保写入磁盘
	logger.Sync()
}

// 示例 3: 按大小轮转
func rotateBySizeExample() {
	println("\n=== 示例 3: 按大小轮转 ===")

	config := &ulogger.Config{
		Path:              "./logs",
		File:              "rotate-size.log",
		RotateSize:        1024, // 1KB
		RotateBackupLimit: 3,    // 保留 3 个备份
		Stdout:            false,
	}

	logger, err := ulogger.New(config)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// 写入大量日志触发轮转
	for i := 0; i < 50; i++ {
		logger.Info("这是一条测试日志，用于触发按大小轮转", "index", i)
	}

	logger.Sync()
	println("日志已写入，检查 ./logs 目录查看轮转文件")
}

// 示例 4: 按时间轮转
func rotateByTimeExample() {
	println("\n=== 示例 4: 按时间轮转 ===")

	config := &ulogger.Config{
		Path:         "./logs",
		File:         "rotate-time.log",
		RotateExpire: 5, // 5 秒轮转一次
		Stdout:       false,
	}

	logger, err := ulogger.New(config)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	logger.Info("第一批日志")

	println("等待 6 秒触发时间轮转...")
	time.Sleep(6 * time.Second)

	logger.Info("第二批日志（已轮转）")
	logger.Sync()

	println("日志已写入，检查 ./logs 目录查看轮转文件")
}

// 示例 5: 带压缩的备份
func compressedBackupExample() {
	println("\n=== 示例 5: 带压缩的备份 ===")

	config := &ulogger.Config{
		Path:                 "./logs",
		File:                 "compressed.log",
		RotateSize:           512,
		RotateBackupLimit:    5,
		RotateBackupCompress: 6, // gzip 压缩级别 6
		Stdout:               false,
	}

	logger, err := ulogger.New(config)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// 写入日志触发轮转和压缩
	for i := 0; i < 30; i++ {
		logger.Info("这条日志会被压缩保存", "index", i, "data", "some data here")
	}

	logger.Sync()
	time.Sleep(1 * time.Second) // 等待压缩完成

	println("日志已写入并压缩，检查 ./logs 目录查看 .gz 文件")
}

// 示例 6: 全局 Logger
func globalLoggerExample() {
	println("\n=== 示例 6: 全局 Logger ===")

	config := &ulogger.Config{
		Path:   "./logs",
		File:   "global.log",
		Prefix: "[GLOBAL]",
		Stdout: true,
	}

	logger, err := ulogger.New(config)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	// 设置为全局默认 Logger
	ulogger.SetDefault(logger)

	// 使用全局函数
	ulogger.Info("使用全局 Logger")
	ulogger.Warn("这是全局警告")

	// 也可以使用标准库的 slog
	slog.Info("标准库 slog 也会使用我们的 Logger")

	logger.Sync()
}
