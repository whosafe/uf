package uhttp

import (
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// MiddlewareTrace 链路追踪中间件
func MiddlewareTrace() unet.MiddlewareFunc {
	return func(next unet.HandlerFunc) unet.HandlerFunc {
		return func(ctx *ucontext.Context, req unet.Request) error {
			httpReq := req.(*Request)

			// 从 HTTP Header 提取或创建追踪上下文
			tc := ucontext.ExtractHTTPHeaders(httpReq.Raw().Header)
			if tc == nil {
				tc = ucontext.NewTraceContext()
			}

			// 创建新的 ucontext.Context 并注入追踪信息
			stdCtx := ucontext.WithContext(ctx.Context(), tc)
			newCtx := ucontext.NewWithContext(stdCtx)

			// 设置响应头
			resp := req.Response().(*Response)
			resp.SetHeader("X-Trace-ID", tc.TraceID)
			resp.SetHeader("X-Span-ID", tc.SpanID)

			return next(newCtx, req)
		}
	}
}
