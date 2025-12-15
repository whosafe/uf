package uhttp

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// CSRFConfig CSRF 保护配置
type CSRFConfig struct {
	// TokenLength Token 长度(字节)
	TokenLength int

	// CookieName Cookie 名称
	CookieName string

	// HeaderName 请求头名称
	HeaderName string

	// FormFieldName 表单字段名称
	FormFieldName string

	// CookiePath Cookie 路径
	CookiePath string

	// CookieMaxAge Cookie 有效期(秒)
	CookieMaxAge int

	// SkipMethods 跳过验证的 HTTP 方法
	SkipMethods []string

	// Store Token 存储器 (可选,默认使用内存存储)
	Store CSRFTokenStore
}

// DefaultCSRFConfig 默认 CSRF 配置
func DefaultCSRFConfig() *CSRFConfig {
	return &CSRFConfig{
		TokenLength:   32,
		CookieName:    "csrf_token",
		HeaderName:    "X-CSRF-Token",
		FormFieldName: "csrf_token",
		CookiePath:    "/",
		CookieMaxAge:  3600, // 1小时
		SkipMethods:   []string{"GET", "HEAD", "OPTIONS", "TRACE"},
		Store:         NewMemoryCSRFStore(), // 默认使用内存存储
	}
}

// generateCSRFToken 生成 CSRF Token
func generateCSRFToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// MiddlewareCSRF CSRF 保护中间件(使用默认配置)
func MiddlewareCSRF() unet.MiddlewareFunc {
	return MiddlewareCSRFWithConfig(*DefaultCSRFConfig())
}

// MiddlewareCSRFWithConfig 使用自定义配置的 CSRF 中间件
func MiddlewareCSRFWithConfig(config CSRFConfig) unet.MiddlewareFunc {
	// 确保 Store 不为 nil
	if config.Store == nil {
		config.Store = NewMemoryCSRFStore()
	}

	// 启动定期清理过期 Token
	go func() {
		ticker := time.NewTicker(time.Duration(config.CookieMaxAge) * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			config.Store.Cleanup()
		}
	}()

	return func(next unet.HandlerFunc) unet.HandlerFunc {
		return func(ctx *ucontext.Context, req unet.Request) error {
			httpReq := req.(*Request)
			httpResp := req.Response().(*Response)

			// 检查是否需要验证
			if shouldSkipCSRF(httpReq.Method(), config.SkipMethods) {
				// 对于安全方法,生成并设置新 Token
				token, err := generateCSRFToken(config.TokenLength)
				if err != nil {
					return err
				}

				// 保存 Token
				if err := config.Store.Store(token, config.CookieMaxAge); err != nil {
					return err
				}

				// 设置 Cookie
				httpResp.SetSessionCookie(config.CookieName, token, config.CookiePath, "", config.CookieMaxAge, true, true, http.SameSiteStrictMode)

				// 将 Token 存储到请求上下文,供模板使用
				httpReq.Set("csrf_token", token)

				return next(ctx, req)
			}

			// 对于状态改变的方法,验证 Token
			// 1. 从 Cookie 获取 Token
			cookieToken, err := httpReq.GetCookie(config.CookieName)
			if err != nil || cookieToken == "" {
				return httpResp.Forbidden("CSRF token missing in cookie")
			}

			// 2. 从请求头或表单获取 Token
			requestToken := httpReq.Header(config.HeaderName)
			if requestToken == "" {
				// 尝试从表单获取
				if err := httpReq.raw.ParseForm(); err == nil {
					requestToken = httpReq.raw.FormValue(config.FormFieldName)
				}
			}

			if requestToken == "" {
				return httpResp.Forbidden("CSRF token missing in request")
			}

			// 3. 验证 Token
			if !config.Store.Validate(cookieToken) || cookieToken != requestToken {
				return httpResp.Forbidden("CSRF token mismatch")
			}

			// 验证通过,继续处理
			return next(ctx, req)
		}
	}
}

// CSRFProtection CSRF 保护中间件(兼容旧接口)
// 推荐使用 MiddlewareCSRF() 或 MiddlewareCSRFWithConfig()
func CSRFProtection(config *CSRFConfig) func(unet.HandlerFunc) unet.HandlerFunc {
	if config == nil {
		return MiddlewareCSRF()
	}
	return MiddlewareCSRFWithConfig(*config)
}

// shouldSkipCSRF 检查是否应该跳过 CSRF 验证
func shouldSkipCSRF(method string, skipMethods []string) bool {
	for _, m := range skipMethods {
		if method == m {
			return true
		}
	}
	return false
}

// GetCSRFToken 从请求中获取 CSRF Token (供模板使用)
func (r *Request) GetCSRFToken() string {
	if token, ok := r.Get("csrf_token"); ok {
		if str, ok := token.(string); ok {
			return str
		}
	}
	return ""
}
