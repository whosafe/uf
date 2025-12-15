package uhttp

import (
	"time"

	"github.com/whosafe/uf/uprotocol/unet"
)

// ApplyDefaultMiddlewares 应用默认中间件到服务器
// 根据服务器配置自动应用中间件 (Trace, Logger, Recovery, CORS, CSRF, Timeout, RateLimit)
func ApplyDefaultMiddlewares(server *Server) {
	cfg := server.MiddlewareConfig()

	// 按照执行顺序添加中间件
	// 1. Trace - 最先执行,为请求添加追踪信息
	if cfg.EnableTrace {
		server.Use(MiddlewareTrace())
	}

	// 2. CORS - 跨域支持
	if cfg.EnableCORS {
		if cfg.CORS != nil {
			server.Use(MiddlewareCORSWithConfig(*cfg.CORS))
		} else {
			server.Use(MiddlewareCORS())
		}
	}

	// 3. CSRF - 跨站请求伪造保护
	if cfg.EnableCSRF {
		if cfg.CSRF != nil {
			server.Use(MiddlewareCSRFWithConfig(*cfg.CSRF))
		} else {
			server.Use(MiddlewareCSRF())
		}
	}

	// 4. Timeout - 超时控制
	if cfg.EnableTimeout {
		timeout := 30 * time.Second // 默认 30 秒
		if cfg.Timeout != "" {
			if d, err := time.ParseDuration(cfg.Timeout); err == nil {
				timeout = d
			}
		}
		server.Use(MiddlewareTimeout(timeout))
	}

	// 4. RateLimit - 限流
	if cfg.EnableRateLimit {
		if cfg.RateLimit != nil {
			// 如果没有设置 KeyFunc,使用默认的基于 IP 的限流
			if cfg.RateLimit.KeyFunc == nil {
				cfg.RateLimit.KeyFunc = func(req *Request) string {
					return req.RemoteAddr().String()
				}
			}
			server.Use(MiddlewareRateLimitWithConfig(cfg.RateLimit))
		} else {
			server.Use(MiddlewareRateLimit())
		}
	}

	// 5. Logger - 记录请求日志
	if cfg.EnableLogger {
		server.Use(MiddlewareLogger())
	}

	// 6. Recovery - 最后执行,捕获 panic
	if cfg.EnableRecovery {
		server.Use(MiddlewareRecovery())
	}
}

// DefaultMiddlewares 返回默认核心中间件列表
// 可以用于手动应用中间件: server.Use(uhttp.DefaultMiddlewares()...)
func DefaultMiddlewares() []unet.MiddlewareFunc {
	return []unet.MiddlewareFunc{
		MiddlewareTrace(),
		MiddlewareLogger(),
		MiddlewareRecovery(),
	}
}
