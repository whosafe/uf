package redis

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/whosafe/uf/uerror"
)

// Config Redis 配置
type Config struct {
	// 连接配置
	Host     string
	Port     int
	Password string
	DB       int

	// 连接池配置
	Pool *PoolConfig

	// 查询配置
	Query *QueryConfig

	// 日志配置
	Log *LogConfig
}

// PoolConfig 连接池配置
type PoolConfig struct {
	MaxIdle     int           // 最大空闲连接数
	MaxActive   int           // 最大活跃连接数
	IdleTimeout time.Duration // 空闲连接超时时间
	MaxLifetime time.Duration // 连接最大生命周期
	PoolSize    int           // 连接池大小
	MinIdleConn int           // 最小空闲连接数
}

// QueryConfig 查询配置
type QueryConfig struct {
	DefaultTimeout     time.Duration // 默认查询超时
	SlowQueryThreshold time.Duration // 慢查询阈值
}

// LogConfig 日志配置
type LogConfig struct {
	Enabled   bool       // 是否启用
	Level     slog.Level // 日志级别
	Format    string     // 格式: json, text
	Output    string     // 输出: stdout, stderr, file
	FilePath  string     // 文件路径
	SlowQuery bool       // 是否记录慢查询
	LogParams bool       // 是否记录查询参数
}

// DefaultConfig 返回默认配置
// 适用于中小型应用，QPS < 1000
func DefaultConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
		Pool: &PoolConfig{
			PoolSize:    10,  // 默认连接池大小
			MinIdleConn: 5,   // 保持 5 个空闲连接
			MaxIdle:     10,  // 最大空闲连接
			MaxActive:   100, // 最大活跃连接
			IdleTimeout: 5 * time.Minute,
			MaxLifetime: 1 * time.Hour,
		},
		Query: &QueryConfig{
			DefaultTimeout:     30 * time.Second,
			SlowQueryThreshold: 100 * time.Millisecond, // 100ms 为慢查询
		},
		Log: &LogConfig{
			Enabled:   true,
			Level:     slog.LevelInfo,
			Format:    "text",
			Output:    "stdout",
			FilePath:  "./logs/redis.log",
			SlowQuery: true,  // 记录慢查询
			LogParams: false, // 生产环境建议关闭参数记录
		},
	}
}

// HighPerformanceConfig 返回高性能配置
// 适用于高并发场景，QPS > 5000
func HighPerformanceConfig() *Config {
	config := DefaultConfig()
	config.Pool.PoolSize = 50    // 增加连接池大小
	config.Pool.MinIdleConn = 20 // 保持更多空闲连接
	config.Pool.MaxActive = 500  // 支持更高并发
	config.Query.DefaultTimeout = 10 * time.Second
	config.Query.SlowQueryThreshold = 50 * time.Millisecond // 更严格的慢查询阈值
	config.Log.LogParams = false                            // 高性能场景不记录参数
	return config
}

// LowLatencyConfig 返回低延迟配置
// 适用于对延迟敏感的场景
func LowLatencyConfig() *Config {
	config := DefaultConfig()
	config.Pool.PoolSize = 30
	config.Pool.MinIdleConn = 15 // 保持足够的预热连接
	config.Query.DefaultTimeout = 5 * time.Second
	config.Query.SlowQueryThreshold = 20 * time.Millisecond // 20ms 即为慢查询
	config.Log.SlowQuery = true
	config.Log.LogParams = true // 低延迟场景需要详细日志排查问题
	return config
}

// Addr 生成连接地址
func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Host == "" {
		return uerror.New("主机地址不能为空")
	}
	if c.Port <= 0 || c.Port > 65535 {
		return uerror.New(fmt.Sprintf("无效的端口号: %d", c.Port))
	}
	if c.DB < 0 {
		return uerror.New(fmt.Sprintf("无效的数据库索引: %d", c.DB))
	}
	if c.Pool != nil {
		if c.Pool.PoolSize < 0 {
			return uerror.New("连接池大小必须 >= 0")
		}
		if c.Pool.MinIdleConn < 0 {
			return uerror.New("最小空闲连接数必须 >= 0")
		}
		if c.Pool.MinIdleConn > c.Pool.PoolSize {
			return uerror.New("最小空闲连接数不能大于连接池大小")
		}
	}
	return nil
}
