package redis

import (
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uerror"
)

// ==================== 哈希操作 ====================

// HGet 获取哈希表中指定字段的值
func (c *Connection) HGet(ctx *ucontext.Context, key, field string) (string, error) {
	startTime := time.Now()
	result, err := c.client.HGet(ctx, key, field).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "HGET", []any{key, field}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return "", ErrNil
		}
		return "", uerror.Wrap(err, "HGET 失败")
	}

	return result, nil
}

// HSet 设置哈希表中字段的值
func (c *Connection) HSet(ctx *ucontext.Context, key string, values ...any) (int64, error) {
	startTime := time.Now()
	result, err := c.client.HSet(ctx, key, values...).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "HSET", append([]any{key}, values...), duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "HSET 失败")
	}

	return result, nil
}

// HSetNX 仅在字段不存在时设置值
func (c *Connection) HSetNX(ctx *ucontext.Context, key, field string, value any) (bool, error) {
	startTime := time.Now()
	result, err := c.client.HSetNX(ctx, key, field, value).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "HSETNX", []any{key, field, value}, duration, err)

	if err != nil {
		return false, uerror.Wrap(err, "HSETNX 失败")
	}

	return result, nil
}

// HGetAll 获取哈希表中所有字段和值
func (c *Connection) HGetAll(ctx *ucontext.Context, key string) (map[string]string, error) {
	startTime := time.Now()
	result, err := c.client.HGetAll(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "HGETALL", []any{key}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "HGETALL 失败")
	}

	return result, nil
}

// HDel 删除哈希表中的字段
func (c *Connection) HDel(ctx *ucontext.Context, key string, fields ...string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.HDel(ctx, key, fields...).Result()
	duration := time.Since(startTime)

	args := make([]any, len(fields)+1)
	args[0] = key
	for i, f := range fields {
		args[i+1] = f
	}
	c.logCommand(ctx, "HDEL", args, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "HDEL 失败")
	}

	return result, nil
}

// HExists 检查哈希表中字段是否存在
func (c *Connection) HExists(ctx *ucontext.Context, key, field string) (bool, error) {
	startTime := time.Now()
	result, err := c.client.HExists(ctx, key, field).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "HEXISTS", []any{key, field}, duration, err)

	if err != nil {
		return false, uerror.Wrap(err, "HEXISTS 失败")
	}

	return result, nil
}

// HLen 获取哈希表中字段的数量
func (c *Connection) HLen(ctx *ucontext.Context, key string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.HLen(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "HLEN", []any{key}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "HLEN 失败")
	}

	return result, nil
}

// HKeys 获取哈希表中所有字段名
func (c *Connection) HKeys(ctx *ucontext.Context, key string) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.HKeys(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "HKEYS", []any{key}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "HKEYS 失败")
	}

	return result, nil
}

// HVals 获取哈希表中所有值
func (c *Connection) HVals(ctx *ucontext.Context, key string) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.HVals(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "HVALS", []any{key}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "HVALS 失败")
	}

	return result, nil
}

// HIncrBy 为哈希表中字段的整数值增加指定数值
func (c *Connection) HIncrBy(ctx *ucontext.Context, key, field string, incr int64) (int64, error) {
	startTime := time.Now()
	result, err := c.client.HIncrBy(ctx, key, field, incr).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "HINCRBY", []any{key, field, incr}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "HINCRBY 失败")
	}

	return result, nil
}

// HIncrByFloat 为哈希表中字段的浮点值增加指定数值
func (c *Connection) HIncrByFloat(ctx *ucontext.Context, key, field string, incr float64) (float64, error) {
	startTime := time.Now()
	result, err := c.client.HIncrByFloat(ctx, key, field, incr).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "HINCRBYFLOAT", []any{key, field, incr}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "HINCRBYFLOAT 失败")
	}

	return result, nil
}

// HMGet 批量获取哈希表中多个字段的值
func (c *Connection) HMGet(ctx *ucontext.Context, key string, fields ...string) ([]any, error) {
	startTime := time.Now()
	result, err := c.client.HMGet(ctx, key, fields...).Result()
	duration := time.Since(startTime)

	args := make([]any, len(fields)+1)
	args[0] = key
	for i, f := range fields {
		args[i+1] = f
	}
	c.logCommand(ctx, "HMGET", args, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "HMGET 失败")
	}

	return result, nil
}
