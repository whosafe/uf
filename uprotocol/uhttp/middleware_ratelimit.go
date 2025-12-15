package uhttp

import (
	"sync"
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	MaxRequests int                       // 最大请求数
	Window      time.Duration             // 时间窗口
	KeyFunc     func(req *Request) string // 获取限流 key 的函数
}

// DefaultRateLimitConfig 默认限流配置
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		MaxRequests: 100,
		Window:      1 * time.Minute,
		KeyFunc: func(req *Request) string {
			// 默认使用 IP 地址作为 key
			return req.RemoteAddr().String()
		},
	}
}

// rateLimiter 限流器
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	config   *RateLimitConfig
}

// visitor 访问者
type visitor struct {
	count     int
	lastReset time.Time
	mu        sync.Mutex
}

// newRateLimiter 创建限流器
func newRateLimiter(config *RateLimitConfig) *rateLimiter {
	limiter := &rateLimiter{
		visitors: make(map[string]*visitor),
		config:   config,
	}

	// 启动清理 goroutine
	go limiter.cleanup()

	return limiter
}

// cleanup 定期清理过期的访问者
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.Window)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for key, v := range rl.visitors {
			v.mu.Lock()
			if time.Since(v.lastReset) > rl.config.Window*2 {
				delete(rl.visitors, key)
			}
			v.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// getVisitor 获取或创建访问者
func (rl *rateLimiter) getVisitor(key string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[key]
	if !exists {
		v = &visitor{
			count:     0,
			lastReset: time.Now(),
		}
		rl.visitors[key] = v
	}

	return v
}

// allow 检查是否允许请求
func (rl *rateLimiter) allow(key string) bool {
	v := rl.getVisitor(key)

	v.mu.Lock()
	defer v.mu.Unlock()

	now := time.Now()

	// 检查是否需要重置
	if now.Sub(v.lastReset) > rl.config.Window {
		v.count = 0
		v.lastReset = now
	}

	// 检查是否超过限制
	if v.count >= rl.config.MaxRequests {
		return false
	}

	v.count++
	return true
}

// MiddlewareRateLimit 限流中间件 (使用默认配置)
func MiddlewareRateLimit() unet.MiddlewareFunc {
	return MiddlewareRateLimitWithConfig(DefaultRateLimitConfig())
}

// MiddlewareRateLimitWithConfig 使用自定义配置的限流中间件
func MiddlewareRateLimitWithConfig(config *RateLimitConfig) unet.MiddlewareFunc {
	limiter := newRateLimiter(config)

	return func(next unet.HandlerFunc) unet.HandlerFunc {
		return func(ctx *ucontext.Context, req unet.Request) error {
			httpReq := req.(*Request)

			// 获取限流 key
			key := config.KeyFunc(httpReq)

			// 检查是否允许
			if !limiter.allow(key) {
				// 超过限制,返回 429
				httpResp := req.Response().(*Response)
				return httpResp.Error(429, CodeRateLimitExceeded, "请求过于频繁,请稍后再试")
			}

			return next(ctx, req)
		}
	}
}

// MiddlewareRateLimitByIP 基于 IP 的限流中间件
func MiddlewareRateLimitByIP(maxRequests int, window time.Duration) unet.MiddlewareFunc {
	return MiddlewareRateLimitWithConfig(&RateLimitConfig{
		MaxRequests: maxRequests,
		Window:      window,
		KeyFunc: func(req *Request) string {
			return req.RemoteAddr().String()
		},
	})
}

// MiddlewareRateLimitByPath 基于路径的限流中间件
func MiddlewareRateLimitByPath(maxRequests int, window time.Duration) unet.MiddlewareFunc {
	return MiddlewareRateLimitWithConfig(&RateLimitConfig{
		MaxRequests: maxRequests,
		Window:      window,
		KeyFunc: func(req *Request) string {
			return req.Path()
		},
	})
}

// MiddlewareRateLimitByIPAndPath 基于 IP 和路径的限流中间件
func MiddlewareRateLimitByIPAndPath(maxRequests int, window time.Duration) unet.MiddlewareFunc {
	return MiddlewareRateLimitWithConfig(&RateLimitConfig{
		MaxRequests: maxRequests,
		Window:      window,
		KeyFunc: func(req *Request) string {
			return req.RemoteAddr().String() + ":" + req.Path()
		},
	})
}
