package template

// internal/middleware/session.go 模板
const MiddlewareSessionTemplate = `package middleware

import (
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// RequireAuth 认证中间件
// 要求用户必须登录才能访问
func RequireAuth() unet.MiddlewareFunc {
	return func(next unet.HandlerFunc) unet.HandlerFunc {
		return func(ctx *ucontext.Context, req unet.Request) error {
			// 获取 Session
			session, err := req.Session()
			if err != nil {
				return req.Response().JSON(401, map[string]string{
					"error": "未授权",
				})
			}

			// 检查用户是否已登录
			userID, ok := session.Get("user_id")
			if !ok || userID == "" {
				return req.Response().JSON(401, map[string]string{
					"error": "请先登录",
				})
			}

			// 将用户 ID 存储到请求上下文中，供后续处理器使用
			req.Set("user_id", userID)

			// 继续处理请求
			return next(ctx, req)
		}
	}
}
`
