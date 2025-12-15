package uhttp

import (
	"log/slog"
	"time"

	"github.com/whosafe/uf/uconfig"
	"github.com/whosafe/uf/uconv"
	"github.com/whosafe/uf/uerror"
)

// globalConfig 全局配置
var globalConfig *Config

// init 自动注册配置回调
func init() {
	globalConfig = DefaultConfig()
	uconfig.Register("server", globalConfig.UnmarshalYAML)
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (c *Config) UnmarshalYAML(key string, node *uconfig.Node) error {

	switch key {
	case "name":
		c.Name = node.String()
	case "protocol":
		c.Protocol = node.String()
	case "address":
		c.Address = node.String()
	case "read_timeout":
		if d, err := time.ParseDuration(node.String()); err == nil {
			c.ReadTimeout = d
		}
	case "write_timeout":
		if d, err := time.ParseDuration(node.String()); err == nil {
			c.WriteTimeout = d
		}
	case "idle_timeout":
		if d, err := time.ParseDuration(node.String()); err == nil {
			c.IdleTimeout = d
		}
	case "max_header_bytes":
		c.MaxHeaderBytes = uconv.ToIntDef(node, 0)
	case "max_body_bytes":
		c.MaxBodyBytes = int64(uconv.ToIntDef(node, 0))
	case "max_form_bytes":
		c.MaxFormBytes = int64(uconv.ToIntDef(node, 0))
	case "keep_alive":
		c.KeepAlive = uconv.ToBoolDef(node, true)
	case "server_agent":
		c.ServerAgent = node.String()
	case "static":
		if c.Static == nil {
			c.Static = &StaticFileConfig{}
		}
		if err := node.Decode(c.Static); err != nil {
			return uerror.Wrap(err, "解析 static 失败")
		}
	case "cookie":
		if c.Cookie == nil {
			c.Cookie = &CookieFileConfig{}
		}
		if err := node.Decode(c.Cookie); err != nil {
			return uerror.Wrap(err, "解析 cookie 失败")
		}
	case "session":

		if c.Session == nil {
			c.Session = &SessionFileConfig{}
		}
		if err := node.Decode(c.Session); err != nil {
			return uerror.Wrap(err, "解析 session 失败")
		}

	case "access_log":
		if c.AccessLog == nil {
			c.AccessLog = &LogConfig{}
		}
		if err := node.Decode(c.AccessLog); err != nil {
			return uerror.Wrap(err, "解析 access_log 失败")
		}
	case "error_log":
		if c.ErrorLog == nil {
			c.ErrorLog = &LogConfig{}
		}
		if err := node.Decode(c.ErrorLog); err != nil {
			return uerror.Wrap(err, "解析 error_log 失败")
		}
	case "middleware":
		if c.Middleware == nil {
			c.Middleware = &MiddlewareConfig{}
		}
		if err := node.Decode(c.Middleware); err != nil {
			return uerror.Wrap(err, "解析 middleware 失败")
		}
	}
	return nil
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (l *LogConfig) UnmarshalYAML(key string, node *uconfig.Node) error {
	switch key {
	case "enabled":
		l.Enabled = uconv.ToBoolDef(node, true)
	case "level":
		levelStr := node.String()
		switch levelStr {
		case "debug":
			l.Level = slog.LevelDebug
		case "info":
			l.Level = slog.LevelInfo
		case "warn":
			l.Level = slog.LevelWarn
		case "error":
			l.Level = slog.LevelError
		default:
			l.Level = slog.LevelInfo
		}
	case "format":
		l.Format = node.String()
	case "output":
		l.Output = node.String()
	case "file_path":
		l.FilePath = node.String()
	case "max_size":
		l.MaxSize = uconv.ToIntDef(node, 100)
	case "max_backups":
		l.MaxBackups = uconv.ToIntDef(node, 10)
	case "max_age":
		l.MaxAge = uconv.ToIntDef(node, 30)
	case "compress":
		l.Compress = uconv.ToBoolDef(node, false)
	}
	return nil
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (s *StaticFileConfig) UnmarshalYAML(key string, node *uconfig.Node) error {

	switch key {
	case "enabled":
		s.Enabled = uconv.ToBoolDef(node, false)
	case "root":
		s.Root = node.String()
	case "prefix":
		s.Prefix = node.String()
	case "index":
		// 解析数组

		s.Index = make([]string, 0)
		if err := node.Iter(func(i int, v *uconfig.Node) error {
			s.Index = append(s.Index, v.String())
			return nil
		}); err != nil {
			return err
		}
	case "browse":
		s.Browse = uconv.ToBoolDef(node, false)
	}
	return nil
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (s *SessionFileConfig) UnmarshalYAML(key string, node *uconfig.Node) error {
	switch key {
	case "enabled":
		s.Enabled = uconv.ToBoolDef(node, false)
	case "provider":
		s.Provider = node.String()
	case "cookie_name":
		s.CookieName = node.String()
	case "max_age":
		s.MaxAge = uconv.ToIntDef(node, 3600)
	}
	return nil
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (c *CookieFileConfig) UnmarshalYAML(key string, node *uconfig.Node) error {
	switch key {
	case "domain":
		c.Domain = node.String()
	case "path":
		c.Path = node.String()
	case "max_age":
		c.MaxAge = uconv.ToIntDef(node, 86400)
	case "secure":
		c.Secure = uconv.ToBoolDef(node, false)
	case "http_only":
		c.HttpOnly = uconv.ToBoolDef(node, true)
	case "same_site":
		c.SameSite = node.String()
	}
	return nil
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (m *MiddlewareConfig) UnmarshalYAML(key string, node *uconfig.Node) error {
	switch key {
	case "enable_trace":
		m.EnableTrace = uconv.ToBoolDef(node, true)
	case "enable_logger":
		m.EnableLogger = uconv.ToBoolDef(node, true)
	case "enable_recovery":
		m.EnableRecovery = uconv.ToBoolDef(node, true)
	case "enable_cors":
		m.EnableCORS = uconv.ToBoolDef(node, false)
	case "cors":
		if m.CORS == nil {
			m.CORS = &CORSConfig{}
		}
		if err := node.Decode(m.CORS); err != nil {
			return uerror.Wrap(err, "解析 cors 失败")
		}
	case "enable_timeout":
		m.EnableTimeout = uconv.ToBoolDef(node, false)
	case "timeout":
		m.Timeout = node.String()
	case "enable_rate_limit":
		m.EnableRateLimit = uconv.ToBoolDef(node, false)
	case "rate_limit":
		if m.RateLimit == nil {
			m.RateLimit = &RateLimitConfig{}
		}
		if err := node.Decode(m.RateLimit); err != nil {
			return uerror.Wrap(err, "解析 rate_limit 失败")
		}
	}
	return nil
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (c *CORSConfig) UnmarshalYAML(key string, node *uconfig.Node) error {
	switch key {
	case "allow_origins":
		c.AllowOrigins = node.String()
	case "allow_methods":
		c.AllowMethods = node.String()
	case "allow_headers":
		c.AllowHeaders = node.String()
	case "allow_credentials":
		c.AllowCredentials = uconv.ToBoolDef(node, false)
	case "expose_headers":
		c.ExposeHeaders = node.String()
	case "max_age":
		c.MaxAge = uconv.ToIntDef(node, 3600)
	}
	return nil
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (r *RateLimitConfig) UnmarshalYAML(key string, node *uconfig.Node) error {
	switch key {
	case "max_requests":
		r.MaxRequests = uconv.ToIntDef(node, 100)
	case "window":
		windowStr := node.String()
		if d, err := time.ParseDuration(windowStr); err == nil {
			r.Window = d
		} else {
			r.Window = 1 * time.Minute // 默认 1 分钟
		}
	}
	return nil
}
