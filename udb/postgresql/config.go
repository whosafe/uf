package postgresql

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/whosafe/uf/uerror"
)

// Config PostgreSQL 配置
type Config struct {
	// 连接配置
	Host     string
	Port     int
	Username string
	Password string
	Database string

	// SSL 配置
	SSLMode string // disable, require, verify-ca, verify-full

	// 连接池配置
	Pool *PoolConfig

	// 查询配置
	Query *QueryConfig

	// 日志配置
	Log *LogConfig
}

// PoolConfig 连接池配置
type PoolConfig struct {
	MaxConns          int32         // 最大连接数
	MinConns          int32         // 最小连接数
	MaxConnLifetime   time.Duration // 连接最大生命周期
	MaxConnIdleTime   time.Duration // 连接最大空闲时间
	HealthCheckPeriod time.Duration // 健康检查周期
}

// QueryConfig 查询配置
type QueryConfig struct {
	DefaultTimeout     time.Duration // 默认查询超时
	SlowQueryThreshold time.Duration // 慢查询阈值
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
	SlowQuery  bool       // 是否记录慢查询
	LogParams  bool       // 是否记录查询参数
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "",
		Database: "postgres",
		SSLMode:  "disable",
		Pool: &PoolConfig{
			MaxConns:          25,
			MinConns:          5,
			MaxConnLifetime:   1 * time.Hour,
			MaxConnIdleTime:   30 * time.Minute,
			HealthCheckPeriod: 1 * time.Minute,
		},
		Query: &QueryConfig{
			DefaultTimeout:     30 * time.Second,
			SlowQueryThreshold: 1 * time.Second,
		},
		Log: &LogConfig{
			Enabled:    true,
			Level:      slog.LevelInfo,
			Format:     "text",
			Output:     "stdout",
			FilePath:   "./logs/db.log",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   false,
			SlowQuery:  true,
			LogParams:  false,
		},
	}
}

// DSN 生成连接字符串
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
	)
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Host == "" {
		return uerror.New("host is required")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return uerror.New(fmt.Sprintf("invalid port: %d", c.Port))
	}
	if c.Username == "" {
		return uerror.New("username is required")
	}
	if c.Database == "" {
		return uerror.New("database is required")
	}
	if c.Pool != nil && c.Pool.MaxConns < c.Pool.MinConns {
		return uerror.New("max_conns must be >= min_conns")
	}
	return nil
}
