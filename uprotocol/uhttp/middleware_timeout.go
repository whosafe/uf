package uhttp

import (
	"context"
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// MiddlewareTimeout 超时控制中间件
func MiddlewareTimeout(timeout time.Duration) unet.MiddlewareFunc {
	return func(next unet.HandlerFunc) unet.HandlerFunc {
		return func(ctx *ucontext.Context, req unet.Request) error {
			// 创建带超时的 context
			timeoutCtx, cancel := context.WithTimeout(ctx.Context(), timeout)
			defer cancel()

			// 创建新的 ucontext.Context
			newCtx := ucontext.NewWithContext(timeoutCtx)

			// 使用通道接收结果
			done := make(chan error, 1)

			go func() {
				done <- next(newCtx, req)
			}()

			select {
			case err := <-done:
				return err
			case <-timeoutCtx.Done():
				httpResp := req.Response().(*Response)
				httpResp.Error(408, CodeInternalError, "请求超时")
				return timeoutCtx.Err()
			}
		}
	}
}
