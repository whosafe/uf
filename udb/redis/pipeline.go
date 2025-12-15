package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uerror"
)

// ==================== Pipeline 操作 ====================

// Pipeline 创建一个 Pipeline
func (c *Connection) Pipeline() redis.Pipeliner {
	return c.client.Pipeline()
}

// TxPipeline 创建一个事务 Pipeline
func (c *Connection) TxPipeline() redis.Pipeliner {
	return c.client.TxPipeline()
}

// Pipelined 在 Pipeline 中执行函数
func (c *Connection) Pipelined(ctx *ucontext.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	cmds, err := c.client.Pipelined(ctx, fn)
	if err != nil {
		return nil, uerror.Wrap(err, "Pipeline 执行失败")
	}
	return cmds, nil
}

// TxPipelined 在事务 Pipeline 中执行函数
func (c *Connection) TxPipelined(ctx *ucontext.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error) {
	cmds, err := c.client.TxPipelined(ctx, fn)
	if err != nil {
		return nil, uerror.Wrap(err, "事务 Pipeline 执行失败")
	}
	return cmds, nil
}
