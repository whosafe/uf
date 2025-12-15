package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/udb"
	"github.com/whosafe/uf/uerror"
	"github.com/whosafe/uf/ulogger"
)

// Connection Redis 客户端连接
type Connection struct {
	client *redis.Client
	config *Config
	logger *ulogger.Logger
}

// New 创建 Redis 客户端连接
func New(config *Config) (*Connection, error) {
	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, uerror.Wrap(err, "无效的配置")
	}

	// 创建 Redis 客户端选项
	opts := &redis.Options{
		Addr:         config.Addr(),
		Password:     config.Password,
		DB:           config.DB,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  config.Query.DefaultTimeout,
		WriteTimeout: config.Query.DefaultTimeout,
	}

	// 设置连接池参数
	if config.Pool != nil {
		opts.PoolSize = config.Pool.PoolSize
		opts.MinIdleConns = config.Pool.MinIdleConn
		opts.ConnMaxIdleTime = config.Pool.IdleTimeout
		opts.ConnMaxLifetime = config.Pool.MaxLifetime
	}

	// 创建客户端
	client := redis.NewClient(opts)

	// 创建 logger
	var logger *ulogger.Logger
	if config.Log != nil && config.Log.Enabled {
		loggerConfig := &ulogger.Config{
			Level:  config.Log.Level,
			Format: config.Log.Format,
			Stdout: config.Log.Output == "stdout" || config.Log.Output == "stderr",
		}

		// 如果输出到文件
		if config.Log.Output == "file" {
			loggerConfig.Path = "./logs"
			loggerConfig.File = "redis.log"
		}

		logger, _ = ulogger.New(loggerConfig)
	} else {
		// 使用默认 logger
		logger, _ = ulogger.New(ulogger.DefaultConfig())
	}

	conn := &Connection{
		client: client,
		config: config,
		logger: logger,
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, uerror.Wrap(err, "连接 Redis 失败")
	}

	// 记录连接创建成功
	if config.Log != nil && config.Log.Enabled {
		logger.Info("Redis 连接创建成功",
			"host", config.Host,
			"port", config.Port,
			"db", config.DB,
			"pool_size", config.Pool.PoolSize)
	}

	return conn, nil
}

// Client 获取原生 Redis 客户端
func (c *Connection) Client() *redis.Client {
	return c.client
}

// Close 关闭连接
func (c *Connection) Close() error {
	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.Info("关闭 Redis 连接")
	}
	return c.client.Close()
}

// Ping 健康检查
func (c *Connection) Ping(ctx *ucontext.Context) error {
	startTime := time.Now()

	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.DebugCtx(ctx, "执行 PING")
	}

	err := c.client.Ping(ctx).Err()
	duration := time.Since(startTime)

	if err != nil {
		if c.config.Log != nil && c.config.Log.Enabled {
			c.logger.ErrorCtx(ctx, "PING 失败", "error", err.Error())
		}
		return uerror.Wrap(err, "PING 失败")
	}

	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.InfoCtx(ctx, "PING 成功", "duration", duration)
	}

	return nil
}

// logCommand 记录命令执行日志
func (c *Connection) logCommand(ctx *ucontext.Context, cmd string, args []any, duration time.Duration, err error) {
	if c.config.Log == nil || !c.config.Log.Enabled {
		return
	}

	// 【安全修复】脱敏参数,防止敏感信息泄露
	sanitizedArgs := args
	if c.config.Log.LogParams {
		sanitizedArgs = udb.SanitizeArgs(args)
	}

	if err != nil {
		// 记录错误
		if c.config.Log.LogParams {
			c.logger.ErrorCtx(ctx, "Redis 命令执行失败",
				"cmd", cmd,
				"args", sanitizedArgs,
				"error", err.Error(),
				"duration", duration)
		} else {
			c.logger.ErrorCtx(ctx, "Redis 命令执行失败",
				"cmd", cmd,
				"error", err.Error(),
				"duration", duration)
		}
		return
	}

	// 记录慢查询
	if c.config.Log.SlowQuery && duration > c.config.Query.SlowQueryThreshold {
		if c.config.Log.LogParams {
			c.logger.WarnCtx(ctx, "Redis 慢查询",
				"cmd", cmd,
				"args", sanitizedArgs,
				"duration", duration,
				"threshold", c.config.Query.SlowQueryThreshold)
		} else {
			c.logger.WarnCtx(ctx, "Redis 慢查询",
				"cmd", cmd,
				"duration", duration,
				"threshold", c.config.Query.SlowQueryThreshold)
		}
		return
	}

	// 记录 Debug 日志
	if c.config.Log.LogParams {
		c.logger.DebugCtx(ctx, "Redis 命令执行",
			"cmd", cmd,
			"args", sanitizedArgs,
			"duration", duration)
	} else {
		c.logger.DebugCtx(ctx, "Redis 命令执行",
			"cmd", cmd,
			"duration", duration)
	}
}
