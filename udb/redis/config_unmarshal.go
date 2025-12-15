package redis

import (
	"log/slog"
	"time"

	"github.com/whosafe/uf/uconfig"
	"github.com/whosafe/uf/uconv"
)

var globalConfig = DefaultConfig()

func init() {
	// 注册配置解析器
	uconfig.Register("database.redis", globalConfig.UnmarshalYAML)
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return globalConfig
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (c *Config) UnmarshalYAML(key string, value *uconfig.Node) error {
	switch key {
	case "host":
		c.Host = value.String()
	case "port":
		c.Port = uconv.ToIntDef(value, 6379)
	case "password":
		c.Password = value.String()
	case "db":
		c.DB = uconv.ToIntDef(value, 0)
	case "pool":
		if c.Pool == nil {
			c.Pool = &PoolConfig{}
		}
		return value.Decode(c.Pool)
	case "query":
		if c.Query == nil {
			c.Query = &QueryConfig{}
		}
		return value.Decode(c.Query)
	case "log":
		if c.Log == nil {
			c.Log = &LogConfig{}
		}
		return value.Decode(c.Log)
	}
	return nil
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (p *PoolConfig) UnmarshalYAML(key string, value *uconfig.Node) error {
	switch key {
	case "pool_size":
		p.PoolSize = uconv.ToIntDef(value, 10)
	case "min_idle_conn":
		p.MinIdleConn = uconv.ToIntDef(value, 5)
	case "max_idle":
		p.MaxIdle = uconv.ToIntDef(value, 10)
	case "max_active":
		p.MaxActive = uconv.ToIntDef(value, 100)
	case "idle_timeout":
		p.IdleTimeout = uconv.ToDurationDef(value, 5*time.Minute)
	case "max_lifetime":
		p.MaxLifetime = uconv.ToDurationDef(value, 1*time.Hour)
	}
	return nil
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (q *QueryConfig) UnmarshalYAML(key string, value *uconfig.Node) error {
	switch key {
	case "default_timeout":
		q.DefaultTimeout = uconv.ToDurationDef(value, 30*time.Second)
	case "slow_query_threshold":
		q.SlowQueryThreshold = uconv.ToDurationDef(value, 100*time.Millisecond)
	}
	return nil
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (l *LogConfig) UnmarshalYAML(key string, value *uconfig.Node) error {
	switch key {
	case "enabled":
		l.Enabled = uconv.ToBoolDef(value, true)
	case "level":
		levelStr := value.String()
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
		l.Format = value.String()
	case "output":
		l.Output = value.String()
	case "file_path":
		l.FilePath = value.String()
	case "slow_query":
		l.SlowQuery = uconv.ToBoolDef(value, true)
	case "log_params":
		l.LogParams = uconv.ToBoolDef(value, false)
	}
	return nil
}
