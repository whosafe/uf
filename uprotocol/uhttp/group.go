package uhttp

import "github.com/whosafe/uf/uprotocol/unet"

// Group 路由组
type Group struct {
	prefix      string
	server      *Server
	middlewares []unet.MiddlewareFunc
}

// Group 创建子路由组
func (g *Group) Group(prefix string) *Group {
	return &Group{
		prefix:      g.prefix + prefix,
		server:      g.server,
		middlewares: append([]unet.MiddlewareFunc{}, g.middlewares...),
	}
}

// Use 注册组级中间件
func (g *Group) Use(middleware ...unet.MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middleware...)
}

// GET 注册 GET 请求处理器
func (g *Group) GET(path string, handler unet.HandlerFunc) {
	g.handle("GET", path, handler)
}

// POST 注册 POST 请求处理器
func (g *Group) POST(path string, handler unet.HandlerFunc) {
	g.handle("POST", path, handler)
}

// PUT 注册 PUT 请求处理器
func (g *Group) PUT(path string, handler unet.HandlerFunc) {
	g.handle("PUT", path, handler)
}

// DELETE 注册 DELETE 请求处理器
func (g *Group) DELETE(path string, handler unet.HandlerFunc) {
	g.handle("DELETE", path, handler)
}

// PATCH 注册 PATCH 请求处理器
func (g *Group) PATCH(path string, handler unet.HandlerFunc) {
	g.handle("PATCH", path, handler)
}

// HEAD 注册 HEAD 请求处理器
func (g *Group) HEAD(path string, handler unet.HandlerFunc) {
	g.handle("HEAD", path, handler)
}

// OPTIONS 注册 OPTIONS 请求处理器
func (g *Group) OPTIONS(path string, handler unet.HandlerFunc) {
	g.handle("OPTIONS", path, handler)
}

// handle 处理路由注册
func (g *Group) handle(method, path string, handler unet.HandlerFunc) {
	// 应用组级中间件
	finalHandler := applyMiddlewares(handler, g.middlewares)

	// 添加路由
	fullPath := g.prefix + path
	g.server.router.addRoute(method, fullPath, finalHandler)
}
