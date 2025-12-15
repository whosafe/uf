package postgresql

import (
	"log/slog"
	"time"

	"github.com/whosafe/uf/uconfig"
	"github.com/whosafe/uf/uconv"
)

var globalConfig *Config

func init() {
	// 注册配置解析器
	uconfig.Register("database.postgres", parseConfig)
}

// parseConfig 解析 PostgreSQL 配置
func parseConfig(key string, value *uconfig.Node) error {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}

	switch key {
	case "host":
		globalConfig.Host = value.String()
	case "port":
		globalConfig.Port = uconv.ToIntDef(value, 5432)
	case "username":
		globalConfig.Username = value.String()
	case "password":
		globalConfig.Password = value.String()
	case "database":
		globalConfig.Database = value.String()
	case "ssl_mode":
		globalConfig.SSLMode = value.String()
	case "pool":
		return value.Decode(globalConfig.Pool)
	case "query":
		return value.Decode(globalConfig.Query)
	case "log":
		return value.Decode(globalConfig.Log)
	}

	return nil
}

// UnmarshalYAML 实现 PoolConfig 的 YAML 解析
func (p *PoolConfig) UnmarshalYAML(key string, value *uconfig.Node) error {
	switch key {
	case "max_conns":
		p.MaxConns = int32(uconv.ToIntDef(value, 25))
	case "min_conns":
		p.MinConns = int32(uconv.ToIntDef(value, 5))
	case "max_conn_lifetime":
		if d, err := time.ParseDuration(value.String()); err == nil {
			p.MaxConnLifetime = d
		}
	case "max_conn_idle_time":
		if d, err := time.ParseDuration(value.String()); err == nil {
			p.MaxConnIdleTime = d
		}
	case "health_check_period":
		if d, err := time.ParseDuration(value.String()); err == nil {
			p.HealthCheckPeriod = d
		}
	}
	return nil
}

// UnmarshalYAML 实现 QueryConfig 的 YAML 解析
func (q *QueryConfig) UnmarshalYAML(key string, value *uconfig.Node) error {
	switch key {
	case "default_timeout":
		if d, err := time.ParseDuration(value.String()); err == nil {
			q.DefaultTimeout = d
		}
	case "slow_query_threshold":
		if d, err := time.ParseDuration(value.String()); err == nil {
			q.SlowQueryThreshold = d
		}
	}
	return nil
}

// UnmarshalYAML 实现 LogConfig 的 YAML 解析
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
	case "max_size":
		l.MaxSize = uconv.ToIntDef(value, 100)
	case "max_backups":
		l.MaxBackups = uconv.ToIntDef(value, 10)
	case "max_age":
		l.MaxAge = uconv.ToIntDef(value, 30)
	case "compress":
		l.Compress = uconv.ToBoolDef(value, false)
	case "slow_query":
		l.SlowQuery = uconv.ToBoolDef(value, true)
	case "log_params":
		l.LogParams = uconv.ToBoolDef(value, false)
	}
	return nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if globalConfig == nil {
		return DefaultConfig()
	}
	return globalConfig
}
