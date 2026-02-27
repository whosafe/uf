package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/whosafe/uf/uerror"
	"github.com/whosafe/uf/ulogger"
	_ "modernc.org/sqlite"
)

// Connection SQLite 客户连接
type Connection struct {
	db     *sql.DB
	config *Config
	logger *ulogger.Logger
}

// New 创建 SQLite 客户连接
func New(config *Config) (*Connection, error) {
	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, uerror.Wrap(err, "invalid config")
	}

	// 创建连接
	db, err := sql.Open("sqlite", config.DSN())
	if err != nil {
		return nil, uerror.Wrap(err, "failed to create sqlite connection")
	}

	// 设置连接池参数
	if config.Pool != nil {
		if config.Pool.MaxConns > 0 {
			db.SetMaxOpenConns(config.Pool.MaxConns)
			db.SetMaxIdleConns(config.Pool.MaxConns) // SQLite 里最大开和空闲一般一致就好
		}
		db.SetConnMaxLifetime(config.Pool.MaxConnLifetime)
		db.SetConnMaxIdleTime(config.Pool.MaxConnIdleTime)
	}

	// 检查连接
	if err := db.Ping(); err != nil {
		return nil, uerror.Wrap(err, "failed to ping sqlite database")
	}

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
			loggerConfig.File = "db.log"
			loggerConfig.RotateSize = config.Log.MaxSize
			loggerConfig.RotateBackupLimit = config.Log.MaxBackups
			loggerConfig.RotateBackupExpire = config.Log.MaxAge * 24 * 3600
			if config.Log.Compress {
				loggerConfig.RotateBackupCompress = 6
			}
		}

		logger, _ = ulogger.New(loggerConfig)
	} else {
		// 使用默认 logger
		logger, _ = ulogger.New(ulogger.DefaultConfig())
	}

	conn := &Connection{
		db:     db,
		config: config,
		logger: logger,
	}

	// 记录连接创建成功
	if config.Log != nil && config.Log.Enabled {
		logger.Info("数据库连接创建成功",
			"path", config.Path,
			"mode", config.Mode,
			"max_conns", config.Pool.MaxConns)
	}

	return conn, nil
}

// Query 创建查询构建器
func (c *Connection) Query(ctx context.Context) *QueryBuilder {
	return &QueryBuilder{
		ctx:        ctx,
		Connection: c,
	}
}

// Insert 创建插入构建器
func (c *Connection) Insert(ctx context.Context) *InsertBuilder {
	return &InsertBuilder{
		ctx:        ctx,
		Connection: c,
	}
}

// Update 创建更新构建器
func (c *Connection) Update(ctx context.Context) *UpdateBuilder {
	return &UpdateBuilder{
		ctx:        ctx,
		Connection: c,
	}
}

// Delete 创建删除构建器
func (c *Connection) Delete(ctx context.Context) *DeleteBuilder {
	return &DeleteBuilder{
		ctx:        ctx,
		Connection: c,
	}
}

// Exec 执行 SQL
func (c *Connection) Exec(ctx context.Context, sqlQuery string, args ...any) (int64, error) {
	// 记录开始时间
	startTime := time.Now()

	// 记录 Debug 日志
	if c.config.Log != nil && c.config.Log.Enabled {
		if c.config.Log.LogParams {
			c.logger.DebugCtx(ctx, "执行 SQL", "sql", sqlQuery, "args", args)
		} else {
			c.logger.DebugCtx(ctx, "执行 SQL", "sql", sqlQuery)
		}
	}

	// 执行
	result, err := c.db.ExecContext(ctx, sqlQuery, args...)
	if err != nil {
		// 记录错误日志
		if c.config.Log != nil && c.config.Log.Enabled {
			c.logger.ErrorCtx(ctx, "SQL 执行失败",
				"sql", sqlQuery,
				"error", err.Error())
		}
		return 0, uerror.Wrap(err, "执行 SQL 失败")
	}

	rowsAffected, _ := result.RowsAffected()
	duration := time.Since(startTime)

	// 记录成功日志
	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.InfoCtx(ctx, "SQL 执行成功",
			"affected_rows", rowsAffected,
			"duration", duration)

		// 记录慢查询
		if c.config.Log.SlowQuery && c.config.Query != nil {
			if duration > c.config.Query.SlowQueryThreshold {
				c.logger.WarnCtx(ctx, "慢 SQL",
					"sql", sqlQuery,
					"duration", duration,
					"threshold", c.config.Query.SlowQueryThreshold,
					"affected_rows", rowsAffected)
			}
		}
	}

	return rowsAffected, nil
}

// Begin 开始事务
func (c *Connection) Begin(ctx context.Context) (*Transaction, error) {
	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.InfoCtx(ctx, "开始事务")
	}

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		if c.config.Log != nil && c.config.Log.Enabled {
			c.logger.ErrorCtx(ctx, "开始事务失败", "error", err.Error())
		}
		return nil, uerror.Wrap(err, "开始事务失败")
	}

	return &Transaction{
		tx:         tx,
		ctx:        ctx,
		Connection: c,
	}, nil
}

// Close 关闭连接池
func (c *Connection) Close() {
	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.Info("关闭数据库连接")
	}
	c.db.Close()
}

// Ping 健康检查
func (c *Connection) Ping(ctx context.Context) error {
	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.DebugCtx(ctx, "执行健康检查")
	}

	err := c.db.PingContext(ctx)
	if err != nil {
		if c.config.Log != nil && c.config.Log.Enabled {
			c.logger.ErrorCtx(ctx, "健康检查失败", "error", err.Error())
		}
		return uerror.Wrap(err, "健康检查失败")
	}

	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.InfoCtx(ctx, "健康检查成功")
	}
	return nil
}
