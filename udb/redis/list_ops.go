package redis

import (
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uerror"
)

// ==================== 列表操作 ====================

// LPush 将一个或多个值插入到列表头部
func (c *Connection) LPush(ctx *ucontext.Context, key string, values ...any) (int64, error) {
	startTime := time.Now()
	result, err := c.client.LPush(ctx, key, values...).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "LPUSH", append([]any{key}, values...), duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "LPUSH 失败")
	}

	return result, nil
}

// RPush 将一个或多个值插入到列表尾部
func (c *Connection) RPush(ctx *ucontext.Context, key string, values ...any) (int64, error) {
	startTime := time.Now()
	result, err := c.client.RPush(ctx, key, values...).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "RPUSH", append([]any{key}, values...), duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "RPUSH 失败")
	}

	return result, nil
}

// LPop 移除并返回列表头部元素
func (c *Connection) LPop(ctx *ucontext.Context, key string) (string, error) {
	startTime := time.Now()
	result, err := c.client.LPop(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "LPOP", []any{key}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return "", ErrNil
		}
		return "", uerror.Wrap(err, "LPOP 失败")
	}

	return result, nil
}

// RPop 移除并返回列表尾部元素
func (c *Connection) RPop(ctx *ucontext.Context, key string) (string, error) {
	startTime := time.Now()
	result, err := c.client.RPop(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "RPOP", []any{key}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return "", ErrNil
		}
		return "", uerror.Wrap(err, "RPOP 失败")
	}

	return result, nil
}

// LRange 获取列表指定范围内的元素
func (c *Connection) LRange(ctx *ucontext.Context, key string, start, stop int64) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.LRange(ctx, key, start, stop).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "LRANGE", []any{key, start, stop}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "LRANGE 失败")
	}

	return result, nil
}

// LLen 获取列表长度
func (c *Connection) LLen(ctx *ucontext.Context, key string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.LLen(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "LLEN", []any{key}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "LLEN 失败")
	}

	return result, nil
}

// LIndex 获取列表中指定索引的元素
func (c *Connection) LIndex(ctx *ucontext.Context, key string, index int64) (string, error) {
	startTime := time.Now()
	result, err := c.client.LIndex(ctx, key, index).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "LINDEX", []any{key, index}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return "", ErrNil
		}
		return "", uerror.Wrap(err, "LINDEX 失败")
	}

	return result, nil
}

// LSet 设置列表中指定索引的元素值
func (c *Connection) LSet(ctx *ucontext.Context, key string, index int64, value any) error {
	startTime := time.Now()
	err := c.client.LSet(ctx, key, index, value).Err()
	duration := time.Since(startTime)

	c.logCommand(ctx, "LSET", []any{key, index, value}, duration, err)

	if err != nil {
		return uerror.Wrap(err, "LSET 失败")
	}

	return nil
}

// LRem 移除列表中与参数值相等的元素
func (c *Connection) LRem(ctx *ucontext.Context, key string, count int64, value any) (int64, error) {
	startTime := time.Now()
	result, err := c.client.LRem(ctx, key, count, value).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "LREM", []any{key, count, value}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "LREM 失败")
	}

	return result, nil
}

// LTrim 修剪列表，只保留指定范围内的元素
func (c *Connection) LTrim(ctx *ucontext.Context, key string, start, stop int64) error {
	startTime := time.Now()
	err := c.client.LTrim(ctx, key, start, stop).Err()
	duration := time.Since(startTime)

	c.logCommand(ctx, "LTRIM", []any{key, start, stop}, duration, err)

	if err != nil {
		return uerror.Wrap(err, "LTRIM 失败")
	}

	return nil
}

// BLPop 阻塞式移除并返回列表头部元素
func (c *Connection) BLPop(ctx *ucontext.Context, timeout time.Duration, keys ...string) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.BLPop(ctx, timeout, keys...).Result()
	duration := time.Since(startTime)

	args := make([]any, len(keys)+1)
	for i, k := range keys {
		args[i] = k
	}
	args[len(keys)] = timeout
	c.logCommand(ctx, "BLPOP", args, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, ErrNil
		}
		return nil, uerror.Wrap(err, "BLPOP 失败")
	}

	return result, nil
}

// BRPop 阻塞式移除并返回列表尾部元素
func (c *Connection) BRPop(ctx *ucontext.Context, timeout time.Duration, keys ...string) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.BRPop(ctx, timeout, keys...).Result()
	duration := time.Since(startTime)

	args := make([]any, len(keys)+1)
	for i, k := range keys {
		args[i] = k
	}
	args[len(keys)] = timeout
	c.logCommand(ctx, "BRPOP", args, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, ErrNil
		}
		return nil, uerror.Wrap(err, "BRPOP 失败")
	}

	return result, nil
}
