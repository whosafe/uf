package unet

import "github.com/whosafe/uf/ucontext"

// HandlerFunc 处理器函数
// ctx: 封装的上下文 (包含 context.Context 和 TraceContext)
// req: 请求对象
type HandlerFunc func(ctx *ucontext.Context, req Request) error

// MiddlewareFunc 中间件函数
// 中间件包装一个 HandlerFunc,返回新的 HandlerFunc
type MiddlewareFunc func(HandlerFunc) HandlerFunc
