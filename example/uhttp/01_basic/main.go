package main

import (
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/uhttp"
	"github.com/whosafe/uf/uprotocol/unet"
)

func main() {
	// 创建服务器
	server := uhttp.New()

	// 注册路由
	server.GET("/", func(ctx *ucontext.Context, req unet.Request) error {
		return req.Response().JSON(200, map[string]string{
			"message": "Hello, World!",
		})
	})

	server.GET("/ping", func(ctx *ucontext.Context, req unet.Request) error {
		return req.Response().JSON(200, map[string]string{
			"status": "pong",
		})
	})

	// 启动服务器
	server.Start(":8080")
}
