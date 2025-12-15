package main

import (
	"fmt"
	"os"

	"github.com/whosafe/uf/cmd/uf/cli"
)

const (
	// Version 脚手架版本
	Version = "0.1.0"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		if err := cli.HandleInit(os.Args[2:]); err != nil {
			fmt.Printf("错误: %v\n", err)
			os.Exit(1)
		}
	case "gen", "generate":
		handleGen()
	case "build":
		handleBuild()
	case "run":
		handleRun()
	case "up", "update":
		handleUp()
	case "version", "-v", "--version":
		handleVersion()
	case "help", "-h", "--help":
		printHelp()
	default:
		fmt.Printf("未知命令: %s\n\n", command)
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`UF 框架脚手架工具

用法:
  uf <command> [arguments]

命令:
  init        初始化新项目
  gen         生成代码（handler/model/middleware/validator）
  build       编译项目
  run         运行开发服务器
  up          更新框架依赖
  version     显示版本信息
  help        显示帮助信息

使用 "uf <command> --help" 查看命令的详细帮助。

示例:
  uf init                    # 交互式创建项目
  uf init --protocol http    # 创建 HTTP 项目
  uf gen handler User        # 生成 User Handler
  uf build                   # 编译项目
  uf run                     # 运行项目
`)
}

func handleGen() {
	fmt.Println("uf gen - 功能开发中...")
}

func handleBuild() {
	fmt.Println("uf build - 功能开发中...")
}

func handleRun() {
	fmt.Println("uf run - 功能开发中...")
}

func handleUp() {
	fmt.Println("uf up - 功能开发中...")
}

func handleVersion() {
	fmt.Printf("UF 脚手架工具 v%s\n", Version)
}
