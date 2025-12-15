package uhttp

import (
	"github.com/whosafe/uf/uprotocol/unet"
	"log/slog"
	"time"
)

// Config HTTP 服务器配置
type Config struct {
	// 基础配置
	Name     string // 服务名称
	Protocol string // 协议类型: http, https
	Address  string // 监听地址

	// 超时配置
	ReadTimeout  time.Duration // 读取超时
	WriteTimeout time.Duration // 写入超时
	IdleTimeout  time.Duration // 空闲超时

	// 大小限制
	MaxHeaderBytes int   // 最大请求头字节数
	MaxBodyBytes   int64 // 最大请求体字节数
	MaxFormBytes   int64 // 最大表单字节数

	// 连接配置
	KeepAlive   bool   // 是否启用 Keep-Alive
	ServerAgent string // Server 头

	// 静态文件配置
	Static *StaticFileConfig // 静态文件配置

	// Cookie 配置
	Cookie *CookieFileConfig // Cookie 配置

	// Session 配置
	Session *SessionFileConfig // Session 配置

	// 中间件配置
	Middleware *MiddlewareConfig // 中间件配置

	// 日志配置
	AccessLog *LogConfig // 访问日志配置
	ErrorLog  *LogConfig // 错误日志配置
}

// MiddlewareConfig 中间件配置
type MiddlewareConfig struct {
	// 核心中间件
	EnableTrace    bool // 是否启用追踪中间件
	EnableLogger   bool // 是否启用日志中间件
	EnableRecovery bool // 是否启用恢复中间件

	// CORS 中间件
	EnableCORS bool        // 是否启用 CORS 中间件
	CORS       *CORSConfig // CORS 配置

	// CSRF 中间件
	EnableCSRF bool        // 是否启用 CSRF 中间件
	CSRF       *CSRFConfig // CSRF 配置

	// 超时中间件
	EnableTimeout bool   // 是否启用超时中间件
	Timeout       string // 超时时间,如 "30s"

	// 限流中间件
	EnableRateLimit bool             // 是否启用限流中间件
	RateLimit       *RateLimitConfig // 限流配置
}

// CookieFileConfig Cookie 文件配置
type CookieFileConfig struct {
	Domain   string // Cookie 域
	Path     string // Cookie 路径
	MaxAge   int    // 最大存活时间(秒)
	Secure   bool   // 是否只在 HTTPS 下传输
	HttpOnly bool   // 是否禁止 JavaScript 访问
	SameSite string // SameSite 策略: strict, lax, none
}

// SessionFileConfig Session 文件配置
type SessionFileConfig struct {
	Enabled    bool   // 是否启用
	Provider   string // 提供者: memory, redis, file
	CookieName string // Cookie 名称
	MaxAge     int    // Cookie 最大存活时间(秒)
}

// StaticFileConfig 静态文件配置
type StaticFileConfig struct {
	Enabled bool     // 是否启用
	Root    string   // 静态文件根目录
	Prefix  string   // URL 前缀
	Index   []string // 索引文件列表
	Browse  bool     // 是否允许目录浏览
}

// LogConfig 日志配置
type LogConfig struct {
	Enabled    bool       // 是否启用
	Level      slog.Level // 日志级别
	Format     string     // 格式: json, text
	Output     string     // 输出: stdout, stderr, file
	FilePath   string     // 文件路径
	MaxSize    int        // 最大文件大小 (MB)
	MaxBackups int        // 最大备份数量
	MaxAge     int        // 最大保留天数
	Compress   bool       // 是否压缩
}

// SessionConfig 会话配置
type SessionConfig struct {
	Store       unet.SessionStore // 存储实现
	CookieName  string            // Cookie 名称
	MaxAge      int               // Cookie 最大存活时间(秒)
	MaxLifetime time.Duration     // 会话最大存活时间
}

// DefaultSessionConfig 默认会话配置
func DefaultSessionConfig() *SessionConfig {
	return &SessionConfig{
		Store:       NewMemoryStore(),
		CookieName:  "session_id",
		MaxAge:      3600,             // 1小时
		MaxLifetime: 30 * time.Minute, // 30分钟
	}
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Name:           "uhttp-server",
		Protocol:       "http",
		Address:        ":8080",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,  // 1MB
		MaxBodyBytes:   10 << 20, // 10MB
		MaxFormBytes:   10 << 20, // 10MB
		KeepAlive:      true,
		ServerAgent:    "UF/1.0",
		Middleware: &MiddlewareConfig{
			// 核心中间件默认启用
			EnableTrace:    true,
			EnableLogger:   true,
			EnableRecovery: true,
			// 其他中间件默认禁用
			EnableCORS:      false,
			EnableTimeout:   false,
			EnableRateLimit: false,
		},
		AccessLog: &LogConfig{
			Enabled: true,
			Level:   slog.LevelInfo,
			Format:  "text",
			Output:  "stdout",
		},
		ErrorLog: &LogConfig{
			Enabled: true,
			Level:   slog.LevelError,
			Format:  "text",
			Output:  "stderr",
		},
	}
}
