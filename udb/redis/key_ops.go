package redis

import (
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uerror"
)

// ==================== 键管理操作 ====================

// Del 删除一个或多个键
func (c *Connection) Del(ctx *ucontext.Context, keys ...string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.Del(ctx, keys...).Result()
	duration := time.Since(startTime)

	args := make([]any, len(keys))
	for i, k := range keys {
		args[i] = k
	}
	c.logCommand(ctx, "DEL", args, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "DEL 失败")
	}

	return result, nil
}

// Exists 检查一个或多个键是否存在
func (c *Connection) Exists(ctx *ucontext.Context, keys ...string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.Exists(ctx, keys...).Result()
	duration := time.Since(startTime)

	args := make([]any, len(keys))
	for i, k := range keys {
		args[i] = k
	}
	c.logCommand(ctx, "EXISTS", args, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "EXISTS 失败")
	}

	return result, nil
}

// Expire 设置键的过期时间（秒）
func (c *Connection) Expire(ctx *ucontext.Context, key string, expiration time.Duration) (bool, error) {
	startTime := time.Now()
	result, err := c.client.Expire(ctx, key, expiration).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "EXPIRE", []any{key, expiration}, duration, err)

	if err != nil {
		return false, uerror.Wrap(err, "EXPIRE 失败")
	}

	return result, nil
}

// ExpireAt 设置键在指定时间点过期
func (c *Connection) ExpireAt(ctx *ucontext.Context, key string, tm time.Time) (bool, error) {
	startTime := time.Now()
	result, err := c.client.ExpireAt(ctx, key, tm).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "EXPIREAT", []any{key, tm}, duration, err)

	if err != nil {
		return false, uerror.Wrap(err, "EXPIREAT 失败")
	}

	return result, nil
}

// TTL 获取键的剩余生存时间（秒）
func (c *Connection) TTL(ctx *ucontext.Context, key string) (time.Duration, error) {
	startTime := time.Now()
	result, err := c.client.TTL(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "TTL", []any{key}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "TTL 失败")
	}

	return result, nil
}

// Persist 移除键的过期时间
func (c *Connection) Persist(ctx *ucontext.Context, key string) (bool, error) {
	startTime := time.Now()
	result, err := c.client.Persist(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "PERSIST", []any{key}, duration, err)

	if err != nil {
		return false, uerror.Wrap(err, "PERSIST 失败")
	}

	return result, nil
}

// Keys 查找所有符合给定模式的键
func (c *Connection) Keys(ctx *ucontext.Context, pattern string) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.Keys(ctx, pattern).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "KEYS", []any{pattern}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "KEYS 失败")
	}

	return result, nil
}

// Rename 重命名键
func (c *Connection) Rename(ctx *ucontext.Context, key, newKey string) error {
	startTime := time.Now()
	err := c.client.Rename(ctx, key, newKey).Err()
	duration := time.Since(startTime)

	c.logCommand(ctx, "RENAME", []any{key, newKey}, duration, err)

	if err != nil {
		return uerror.Wrap(err, "RENAME 失败")
	}

	return nil
}

// RenameNX 仅在新键不存在时重命名键
func (c *Connection) RenameNX(ctx *ucontext.Context, key, newKey string) (bool, error) {
	startTime := time.Now()
	result, err := c.client.RenameNX(ctx, key, newKey).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "RENAMENX", []any{key, newKey}, duration, err)

	if err != nil {
		return false, uerror.Wrap(err, "RENAMENX 失败")
	}

	return result, nil
}

// Type 返回键所存储值的类型
func (c *Connection) Type(ctx *ucontext.Context, key string) (string, error) {
	startTime := time.Now()
	result, err := c.client.Type(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "TYPE", []any{key}, duration, err)

	if err != nil {
		return "", uerror.Wrap(err, "TYPE 失败")
	}

	return result, nil
}

// Scan 迭代数据库中的键
func (c *Connection) Scan(ctx *ucontext.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	startTime := time.Now()
	result, nextCursor, err := c.client.Scan(ctx, cursor, match, count).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SCAN", []any{cursor, match, count}, duration, err)

	if err != nil {
		return nil, 0, uerror.Wrap(err, "SCAN 失败")
	}

	return result, nextCursor, nil
}
