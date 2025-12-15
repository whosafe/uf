package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uerror"
)

// ==================== 事务操作 ====================

// Watch 监视一个或多个键
func (c *Connection) Watch(ctx *ucontext.Context, fn func(*redis.Tx) error, keys ...string) error {
	err := c.client.Watch(ctx, fn, keys...)
	if err != nil {
		return uerror.Wrap(err, "监视事务失败")
	}
	return nil
}

// TxPipelineClient 返回事务 Pipeline
// 使用示例：
//
//	pipe := conn.TxPipelineClient()
//	pipe.Set(ctx, "key1", "value1", 0)
//	pipe.Set(ctx, "key2", "value2", 0)
//	_, err := pipe.Exec(ctx)
func (c *Connection) TxPipelineClient() redis.Pipeliner {
	return c.client.TxPipeline()
}
