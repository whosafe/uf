package unet

import (
	"context"
	"net"
)

// Server 服务器接口 (所有协议都要实现)
type Server interface {
	// Start 启动服务器
	Start(addr string) error

	// Stop 停止服务器
	Stop(ctx context.Context) error

	// Serve 处理连接 (阻塞)
	Serve(listener net.Listener) error

	// Use 注册中间件
	Use(middleware ...MiddlewareFunc)

	// Handle 注册处理器
	Handle(pattern string, handler HandlerFunc)
}
