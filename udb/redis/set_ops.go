package redis

import (
	"time"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uerror"
)

// ==================== 集合操作 ====================

// SAdd 向集合添加一个或多个成员
func (c *Connection) SAdd(ctx *ucontext.Context, key string, members ...any) (int64, error) {
	startTime := time.Now()
	result, err := c.client.SAdd(ctx, key, members...).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SADD", append([]any{key}, members...), duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "SADD 失败")
	}

	return result, nil
}

// SRem 移除集合中一个或多个成员
func (c *Connection) SRem(ctx *ucontext.Context, key string, members ...any) (int64, error) {
	startTime := time.Now()
	result, err := c.client.SRem(ctx, key, members...).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SREM", append([]any{key}, members...), duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "SREM 失败")
	}

	return result, nil
}

// SMembers 获取集合中所有成员
func (c *Connection) SMembers(ctx *ucontext.Context, key string) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.SMembers(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SMEMBERS", []any{key}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "SMEMBERS 失败")
	}

	return result, nil
}

// SIsMember 判断成员是否在集合中
func (c *Connection) SIsMember(ctx *ucontext.Context, key string, member any) (bool, error) {
	startTime := time.Now()
	result, err := c.client.SIsMember(ctx, key, member).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SISMEMBER", []any{key, member}, duration, err)

	if err != nil {
		return false, uerror.Wrap(err, "SISMEMBER 失败")
	}

	return result, nil
}

// SCard 获取集合的成员数
func (c *Connection) SCard(ctx *ucontext.Context, key string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.SCard(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SCARD", []any{key}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "SCARD 失败")
	}

	return result, nil
}

// SPop 移除并返回集合中的一个随机元素
func (c *Connection) SPop(ctx *ucontext.Context, key string) (string, error) {
	startTime := time.Now()
	result, err := c.client.SPop(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SPOP", []any{key}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return "", ErrNil
		}
		return "", uerror.Wrap(err, "SPOP 失败")
	}

	return result, nil
}

// SPopN 移除并返回集合中的 N 个随机元素
func (c *Connection) SPopN(ctx *ucontext.Context, key string, count int64) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.SPopN(ctx, key, count).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SPOP", []any{key, count}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "SPOP 失败")
	}

	return result, nil
}

// SRandMember 返回集合中的一个随机元素
func (c *Connection) SRandMember(ctx *ucontext.Context, key string) (string, error) {
	startTime := time.Now()
	result, err := c.client.SRandMember(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SRANDMEMBER", []any{key}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return "", ErrNil
		}
		return "", uerror.Wrap(err, "SRANDMEMBER 失败")
	}

	return result, nil
}

// SRandMemberN 返回集合中的 N 个随机元素
func (c *Connection) SRandMemberN(ctx *ucontext.Context, key string, count int64) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.SRandMemberN(ctx, key, count).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "SRANDMEMBER", []any{key, count}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "SRANDMEMBER 失败")
	}

	return result, nil
}

// SUnion 返回多个集合的并集
func (c *Connection) SUnion(ctx *ucontext.Context, keys ...string) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.SUnion(ctx, keys...).Result()
	duration := time.Since(startTime)

	args := make([]any, len(keys))
	for i, k := range keys {
		args[i] = k
	}
	c.logCommand(ctx, "SUNION", args, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "SUNION 失败")
	}

	return result, nil
}

// SInter 返回多个集合的交集
func (c *Connection) SInter(ctx *ucontext.Context, keys ...string) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.SInter(ctx, keys...).Result()
	duration := time.Since(startTime)

	args := make([]any, len(keys))
	for i, k := range keys {
		args[i] = k
	}
	c.logCommand(ctx, "SINTER", args, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "SINTER 失败")
	}

	return result, nil
}

// SDiff 返回多个集合的差集
func (c *Connection) SDiff(ctx *ucontext.Context, keys ...string) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.SDiff(ctx, keys...).Result()
	duration := time.Since(startTime)

	args := make([]any, len(keys))
	for i, k := range keys {
		args[i] = k
	}
	c.logCommand(ctx, "SDIFF", args, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "SDIFF 失败")
	}

	return result, nil
}
