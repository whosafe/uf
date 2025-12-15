package uhttp

import (
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// CORSConfig CORS 配置·
type CORSConfig struct {
	AllowOrigins     string // 允许的源,默认 "*"
	AllowMethods     string // 允许的方法,默认 "GET,POST,PUT,DELETE,PATCH,HEAD,OPTIONS"
	AllowHeaders     string // 允许的请求头,默认 "*"
	AllowCredentials bool   // 是否允许携带凭证
	ExposeHeaders    string // 暴露的响应头
	MaxAge           int    // 预检请求缓存时间(秒)
}

// DefaultCORSConfig 默认 CORS 配置
var DefaultCORSConfig = CORSConfig{
	AllowOrigins:     "*",
	AllowMethods:     "GET,POST,PUT,DELETE,PATCH,HEAD,OPTIONS",
	AllowHeaders:     "*",
	AllowCredentials: false,
	ExposeHeaders:    "",
	MaxAge:           3600,
}

// MiddlewareCORS 跨域支持中间件(使用默认配置)
func MiddlewareCORS() unet.MiddlewareFunc {
	return MiddlewareCORSWithConfig(DefaultCORSConfig)
}

// MiddlewareCORSWithConfig 使用自定义配置的 CORS 中间件
func MiddlewareCORSWithConfig(config CORSConfig) unet.MiddlewareFunc {
	return func(next unet.HandlerFunc) unet.HandlerFunc {
		return func(ctx *ucontext.Context, req unet.Request) error {
			httpReq := req.(*Request)
			resp := req.Response().(*Response)

			// 设置 CORS 头
			resp.SetHeader("Access-Control-Allow-Origin", config.AllowOrigins)
			resp.SetHeader("Access-Control-Allow-Methods", config.AllowMethods)
			resp.SetHeader("Access-Control-Allow-Headers", config.AllowHeaders)

			if config.AllowCredentials {
				resp.SetHeader("Access-Control-Allow-Credentials", "true")
			}

			if config.ExposeHeaders != "" {
				resp.SetHeader("Access-Control-Expose-Headers", config.ExposeHeaders)
			}

			if config.MaxAge > 0 {
				resp.SetHeader("Access-Control-Max-Age", string(rune(config.MaxAge)))
			}

			// 处理 OPTIONS 预检请求
			if httpReq.Method() == "OPTIONS" {
				return resp.String(204, "")
			}

			return next(ctx, req)
		}
	}
}
