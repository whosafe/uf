package uhttp

import (
	"context"
	"net"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/ulogger"
	"github.com/whosafe/uf/uprotocol/unet"
)

// Server HTTP 服务器
type Server struct {
	config         *Config
	router         *Router
	middlewares    []unet.MiddlewareFunc
	httpServer     *http.Server
	accessLogger   *ulogger.Logger // 访问日志
	errorLogger    *ulogger.Logger // 错误日志
	sessionManager *SessionManager // Session 管理器
	mu             sync.RWMutex
}

// New 创建新的 HTTP 服务器
func New() *Server {
	return NewWithConfig(GetConfig())
}

// NewWithConfig 使用配置创建 HTTP 服务器
func NewWithConfig(cfg *Config) *Server {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	s := &Server{
		config:      cfg,
		router:      NewRouter(),
		middlewares: make([]unet.MiddlewareFunc, 0),
	}

	// 创建访问日志 Logger
	if cfg.AccessLog != nil && cfg.AccessLog.Enabled {
		accessLogger, err := createLogger(cfg.AccessLog)
		if err != nil {
			ulogger.Warn("创建访问日志失败,使用默认 logger", "error", err)
			s.accessLogger = ulogger.Default()
		} else {
			s.accessLogger = accessLogger
		}
	} else {
		s.accessLogger = ulogger.Default()
	}

	// 创建错误日志 Logger
	if cfg.ErrorLog != nil && cfg.ErrorLog.Enabled {
		errorLogger, err := createLogger(cfg.ErrorLog)
		if err != nil {
			ulogger.Warn("创建错误日志失败,使用默认 logger", "error", err)
			s.errorLogger = ulogger.Default()
		} else {
			s.errorLogger = errorLogger
		}
	} else {
		s.errorLogger = ulogger.Default()
	}

	// 创建 HTTP Server
	s.httpServer = &http.Server{
		Addr:           cfg.Address,
		Handler:        s,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.IdleTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	// 自动启用静态文件服务
	if cfg.Static != nil && cfg.Static.Enabled {
		s.StaticWithConfig(&StaticConfig{
			Root:   cfg.Static.Root,
			Prefix: cfg.Static.Prefix,
			Index:  cfg.Static.Index,
			Browse: cfg.Static.Browse,
		})
	}

	// 自动启用 Session
	if cfg.Session != nil && cfg.Session.Enabled {
		var store unet.SessionStore
		switch cfg.Session.Provider {
		case "redis":
			// Redis 存储需要外部提供 RedisClient,这里使用内存存储
			store = NewMemoryStore()
		case "file":
			// 文件存储暂未实现,使用内存存储
			store = NewMemoryStore()
		default: // memory
			store = NewMemoryStore()
		}

		cookieName := cfg.Session.CookieName
		if cookieName == "" {
			cookieName = "session_id"
		}
		maxAge := cfg.Session.MaxAge
		if maxAge == 0 {
			maxAge = 3600
		}

		s.sessionManager = NewSessionManager(store, cookieName, maxAge)
	}

	return s
}

// createLogger 根据配置创建 Logger
func createLogger(cfg *LogConfig) (*ulogger.Logger, error) {
	loggerCfg := &ulogger.Config{
		Level:  cfg.Level,
		Stdout: cfg.Output == "stdout" || cfg.Output == "stderr",
	}

	// 设置格式
	if cfg.Format == "json" {
		loggerCfg.Format = "json"
	} else {
		loggerCfg.Format = "text"
	}

	// 如果输出到文件
	if cfg.Output == "file" && cfg.FilePath != "" {
		loggerCfg.Path = filepath.Dir(cfg.FilePath)
		loggerCfg.File = filepath.Base(cfg.FilePath)
		loggerCfg.RotateSize = cfg.MaxSize * 1024 * 1024 // MB to bytes
		loggerCfg.RotateBackupLimit = cfg.MaxBackups
		loggerCfg.RotateBackupExpire = cfg.MaxAge * 24 * 3600 // days to seconds
		if cfg.Compress {
			loggerCfg.RotateBackupCompress = 6
		}
	}

	return ulogger.New(loggerCfg)
}

// SetAccessLogger 设置访问日志 Logger
func (s *Server) SetAccessLogger(logger *ulogger.Logger) {
	s.accessLogger = logger
}

// AccessLogger 获取访问日志 Logger
func (s *Server) AccessLogger() *ulogger.Logger {
	return s.accessLogger
}

// SetErrorLogger 设置错误日志 Logger
func (s *Server) SetErrorLogger(logger *ulogger.Logger) {
	s.errorLogger = logger
}

// ErrorLogger 获取错误日志 Logger
func (s *Server) ErrorLogger() *ulogger.Logger {
	return s.errorLogger
}

// Start 启动服务器
func (s *Server) Start(addr ...string) error {
	if len(addr) > 0 {
		s.httpServer.Addr = addr[0]
	}

	s.accessLogger.Info("HTTP 服务器启动", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Stop 停止服务器
func (s *Server) Stop(ctx context.Context) error {
	s.accessLogger.Info("正在关闭 HTTP 服务器...")
	return s.httpServer.Shutdown(ctx)
}

// Serve 处理连接 (阻塞)
func (s *Server) Serve(listener net.Listener) error {
	return s.httpServer.Serve(listener)
}

// MiddlewareConfig 获取中间件配置
func (s *Server) MiddlewareConfig() *MiddlewareConfig {
	// 如果没有配置中间件,使用默认配置(全部启用)
	if s.config.Middleware == nil {
		return &MiddlewareConfig{
			EnableTrace:    true,
			EnableLogger:   true,
			EnableRecovery: true,
		}
	}
	return s.config.Middleware
}

// Use 注册全局中间件
func (s *Server) Use(middleware ...unet.MiddlewareFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.middlewares = append(s.middlewares, middleware...)
}

// Handle 注册处理器
func (s *Server) Handle(pattern string, handler unet.HandlerFunc) {
	s.router.addRoute("ANY", pattern, handler)
}

// ServeHTTP 实现 http.Handler 接口
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 创建 Request 和 Response (从对象池获取)
	req := newRequest(r, w, s)

	// 请求处理完成后释放回对象池
	defer func() {
		req.response.release()
		req.release()
	}()

	// 查找路由
	handler, params := s.router.getValue(r.Method, r.URL.Path)
	if handler == nil {
		// 404
		req.response.NotFound("路由不存在")
		return
	}

	// 设置路径参数
	req.params = params

	// 应用中间件
	finalHandler := applyMiddlewares(handler, s.middlewares)

	// 创建追踪上下文
	ctx := ucontext.NewWithContext(r.Context())

	// 执行处理器
	if err := finalHandler(ctx, req); err != nil {
		s.errorLogger.ErrorCtx(ctx.Context(), "处理请求失败", "error", err)
		// 只有在响应未写入时才返回错误响应
		if !req.response.IsWritten() {
			req.response.InternalError("服务器内部错误")
		}
	}
}

// GET 注册 GET 请求处理器
func (s *Server) GET(path string, handler unet.HandlerFunc) {
	s.router.addRoute(http.MethodGet, path, handler)
}

// POST 注册 POST 请求处理器
func (s *Server) POST(path string, handler unet.HandlerFunc) {
	s.router.addRoute(http.MethodPost, path, handler)
}

// PUT 注册 PUT 请求处理器
func (s *Server) PUT(path string, handler unet.HandlerFunc) {
	s.router.addRoute(http.MethodPut, path, handler)
}

// DELETE 注册 DELETE 请求处理器
func (s *Server) DELETE(path string, handler unet.HandlerFunc) {
	s.router.addRoute(http.MethodDelete, path, handler)
}

// PATCH 注册 PATCH 请求处理器
func (s *Server) PATCH(path string, handler unet.HandlerFunc) {
	s.router.addRoute(http.MethodPatch, path, handler)
}

// HEAD 注册 HEAD 请求处理器
func (s *Server) HEAD(path string, handler unet.HandlerFunc) {
	s.router.addRoute(http.MethodHead, path, handler)
}

// OPTIONS 注册 OPTIONS 请求处理器
func (s *Server) OPTIONS(path string, handler unet.HandlerFunc) {
	s.router.addRoute(http.MethodOptions, path, handler)
}

// Group 创建路由组
func (s *Server) Group(prefix string) *Group {
	return &Group{
		prefix: prefix,
		server: s,
	}
}

// applyMiddlewares 应用中间件链
func applyMiddlewares(handler unet.HandlerFunc, middlewares []unet.MiddlewareFunc) unet.HandlerFunc {
	// 从后往前应用中间件
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
