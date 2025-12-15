package main

import (
	"github.com/whosafe/uf/uconfig"
	"github.com/whosafe/uf/ulogger"
)

func TestConfig() {
	println("=== 测试 uconfig 集成 ===\n")

	// 1. 注册 logger 到 uconfig
	println("1. 注册 logger 到 uconfig")
	ulogger.Register()

	// 2. 加载配置文件
	println("2. 加载配置文件: config.yaml")
	if err := uconfig.Load("config.yaml"); err != nil {
		println("加载配置失败:", err.Error())
		return
	}

	// 3. 使用默认 logger（应该已经使用配置文件中的设置）
	println("\n3. 使用默认 logger（应该使用配置文件中的设置）\n")
	ulogger.Debug("这是 Debug 日志")
	ulogger.Info("这是 Info 日志")
	ulogger.Warn("这是 Warn 日志", "key", "value")
	ulogger.Error("这是 Error 日志", "user", "alice")

	// 4. 验证配置是否生效
	println("\n4. 检查日志文件是否创建在 ./logs/app.log")
	println("   前缀应该是 [MyApp]")
	println("   日志级别应该是 debug")
	println("   应该同时输出到终端和文件")
}
