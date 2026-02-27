package sqlite

import (
	"log/slog"
	"time"

	"github.com/whosafe/uf/uerror"
)

// Config SQLite 配置
type Config struct {
	// 连接配置
	Path string // 数据库文件路径, 例如: ./data.db, :memory:
	Mode string // ro, rw, rwc, memory

	// 连接池配置
	Pool *PoolConfig

	// 查询配置
	Query *QueryConfig

	// 日志配置
	Log *LogConfig
}

// PoolConfig 连接池配置
type PoolConfig struct {
	MaxConns        int           // 最大连接数
	MaxConnLifetime time.Duration // 连接最大生命周期
	MaxConnIdleTime time.Duration // 连接最大空闲时间
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
		Path: "data.db",
		Mode: "rwc",
		Pool: &PoolConfig{
			MaxConns:        10,
			MaxConnLifetime: 1 * time.Hour,
			MaxConnIdleTime: 30 * time.Minute,
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
	dsn := c.Path
	if c.Mode != "" {
		dsn += "?mode=" + c.Mode
	}
	return dsn
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Path == "" {
		return uerror.New("db path is required")
	}
	if c.Pool != nil && c.Pool.MaxConns <= 0 {
		return uerror.New("max_conns must be > 0")
	}
	return nil
}
