package uhttp

import (
	"fmt"
	"runtime/debug"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uerror"
	"github.com/whosafe/uf/uprotocol/unet"
)

// MiddlewareRecovery 异常恢复中间件
func MiddlewareRecovery() unet.MiddlewareFunc {
	return func(next unet.HandlerFunc) unet.HandlerFunc {
		return func(ctx *ucontext.Context, req unet.Request) (err error) {
			defer func() {
				if r := recover(); r != nil {
					httpReq := req.(*Request)
					resp := req.Response().(*Response)

					// 获取错误日志 logger
					logger := httpReq.Server().ErrorLogger()

					// 记录错误日志
					logger.ErrorCtx(ctx.Context(), "Panic recovered",
						"error", r,
						"stack", string(debug.Stack()),
					)

					// 返回 500 错误
					resp.InternalError("服务器内部错误")

					// 设置错误
					err = uerror.New(fmt.Sprintf("panic recovered: %v", r))
				}
			}()

			return next(ctx, req)
		}
	}
}
