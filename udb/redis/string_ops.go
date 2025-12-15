package redis

import (
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uconv"
	"github.com/whosafe/uf/uerror"
)

// ==================== 字符串操作 ====================

// Get 获取字符串值
func (c *Connection) Get(ctx *ucontext.Context, key string) (string, error) {
	startTime := time.Now()
	result, err := c.client.Get(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "GET", []any{key}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return "", ErrNil
		}
		return "", uerror.Wrap(err, "GET 失败")
	}

	return result, nil
}

// Set 设置字符串值
func (c *Connection) Set(ctx *ucontext.Context, key string, value any, expiration time.Duration) error {
	startTime := time.Now()
	err := c.client.Set(ctx, key, value, expiration).Err()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SET", []any{key, value, expiration}, duration, err)

	if err != nil {
		return uerror.Wrap(err, "SET 失败")
	}

	return nil
}

// SetNX 仅在键不存在时设置值
func (c *Connection) SetNX(ctx *ucontext.Context, key string, value any, expiration time.Duration) (bool, error) {
	startTime := time.Now()
	result, err := c.client.SetNX(ctx, key, value, expiration).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SETNX", []any{key, value, expiration}, duration, err)

	if err != nil {
		return false, uerror.Wrap(err, "SETNX 失败")
	}

	return result, nil
}

// GetSet 设置新值并返回旧值
func (c *Connection) GetSet(ctx *ucontext.Context, key string, value any) (string, error) {
	startTime := time.Now()
	result, err := c.client.GetSet(ctx, key, value).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "GETSET", []any{key, value}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return "", ErrNil
		}
		return "", uerror.Wrap(err, "GETSET 失败")
	}

	return result, nil
}

// MGet 批量获取多个键的值
func (c *Connection) MGet(ctx *ucontext.Context, keys ...string) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.MGet(ctx, keys...).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "MGET", []any{keys}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "MGET 失败")
	}

	// 转换结果
	values := make([]string, len(result))
	for i, v := range result {
		values[i] = uconv.ToString(v)
	}

	return values, nil
}

// MSet 批量设置多个键值对
func (c *Connection) MSet(ctx *ucontext.Context, pairs ...any) error {
	startTime := time.Now()
	err := c.client.MSet(ctx, pairs...).Err()
	duration := time.Since(startTime)

	c.logCommand(ctx, "MSET", pairs, duration, err)

	if err != nil {
		return uerror.Wrap(err, "MSET 失败")
	}

	return nil
}

// Incr 将键的整数值加 1
func (c *Connection) Incr(ctx *ucontext.Context, key string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.Incr(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "INCR", []any{key}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "INCR 失败")
	}

	return result, nil
}

// IncrBy 将键的整数值增加指定数值
func (c *Connection) IncrBy(ctx *ucontext.Context, key string, value int64) (int64, error) {
	startTime := time.Now()
	result, err := c.client.IncrBy(ctx, key, value).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "INCRBY", []any{key, value}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "INCRBY 失败")
	}

	return result, nil
}

// Decr 将键的整数值减 1
func (c *Connection) Decr(ctx *ucontext.Context, key string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.Decr(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "DECR", []any{key}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "DECR 失败")
	}

	return result, nil
}

// DecrBy 将键的整数值减少指定数值
func (c *Connection) DecrBy(ctx *ucontext.Context, key string, value int64) (int64, error) {
	startTime := time.Now()
	result, err := c.client.DecrBy(ctx, key, value).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "DECRBY", []any{key, value}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "DECRBY 失败")
	}

	return result, nil
}

// Append 追加字符串到键的值
func (c *Connection) Append(ctx *ucontext.Context, key, value string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.Append(ctx, key, value).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "APPEND", []any{key, value}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "APPEND 失败")
	}

	return result, nil
}

// StrLen 获取字符串长度
func (c *Connection) StrLen(ctx *ucontext.Context, key string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.StrLen(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "STRLEN", []any{key}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "STRLEN 失败")
	}

	return result, nil
}
