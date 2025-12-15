package uhttp

import (
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// MiddlewareLogger 请求日志中间件
func MiddlewareLogger() unet.MiddlewareFunc {
	return func(next unet.HandlerFunc) unet.HandlerFunc {
		return func(ctx *ucontext.Context, req unet.Request) error {
			start := time.Now()
			httpReq := req.(*Request)

			// 执行处理器
			err := next(ctx, req)

			// 获取访问日志 logger
			logger := httpReq.Server().AccessLogger()

			// 记录日志
			duration := time.Since(start)
			resp := req.Response().(*Response)

			logger.InfoCtx(ctx.Context(), "HTTP Request",
				"method", httpReq.Method(),
				"path", httpReq.Path(),
				"status", resp.StatusCode(),
				"duration_ms", duration.Milliseconds(),
				"client_ip", httpReq.RemoteAddr().String(),
			)

			return err
		}
	}
}
