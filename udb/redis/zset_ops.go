package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uerror"
)

// ==================== 有序集合操作 ====================

// ZAdd 向有序集合添加一个或多个成员
func (c *Connection) ZAdd(ctx *ucontext.Context, key string, members ...redis.Z) (int64, error) {
	startTime := time.Now()
	result, err := c.client.ZAdd(ctx, key, members...).Result()
	duration := time.Since(startTime)

	args := make([]any, len(members)+1)
	args[0] = key
	for i, m := range members {
		args[i+1] = m
	}
	c.logCommand(ctx, "ZADD", args, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "ZADD 失败")
	}

	return result, nil
}

// ZRem 移除有序集合中的一个或多个成员
func (c *Connection) ZRem(ctx *ucontext.Context, key string, members ...any) (int64, error) {
	startTime := time.Now()
	result, err := c.client.ZRem(ctx, key, members...).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZREM", append([]any{key}, members...), duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "ZREM 失败")
	}

	return result, nil
}

// ZRange 返回有序集合中指定区间内的成员
func (c *Connection) ZRange(ctx *ucontext.Context, key string, start, stop int64) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.ZRange(ctx, key, start, stop).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZRANGE", []any{key, start, stop}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "ZRANGE 失败")
	}

	return result, nil
}

// ZRangeWithScores 返回有序集合中指定区间内的成员及其分数
func (c *Connection) ZRangeWithScores(ctx *ucontext.Context, key string, start, stop int64) ([]redis.Z, error) {
	startTime := time.Now()
	result, err := c.client.ZRangeWithScores(ctx, key, start, stop).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZRANGE", []any{key, start, stop, "WITHSCORES"}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "ZRANGE WITHSCORES 失败")
	}

	return result, nil
}

// ZRevRange 返回有序集合中指定区间内的成员（按分数从高到低）
func (c *Connection) ZRevRange(ctx *ucontext.Context, key string, start, stop int64) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.ZRevRange(ctx, key, start, stop).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZREVRANGE", []any{key, start, stop}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "ZREVRANGE 失败")
	}

	return result, nil
}

// ZRevRangeWithScores 返回有序集合中指定区间内的成员及其分数（按分数从高到低）
func (c *Connection) ZRevRangeWithScores(ctx *ucontext.Context, key string, start, stop int64) ([]redis.Z, error) {
	startTime := time.Now()
	result, err := c.client.ZRevRangeWithScores(ctx, key, start, stop).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZREVRANGE", []any{key, start, stop, "WITHSCORES"}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "ZREVRANGE WITHSCORES 失败")
	}

	return result, nil
}

// ZRangeByScore 返回有序集合中指定分数区间的成员
func (c *Connection) ZRangeByScore(ctx *ucontext.Context, key string, opt *redis.ZRangeBy) ([]string, error) {
	startTime := time.Now()
	result, err := c.client.ZRangeByScore(ctx, key, opt).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZRANGEBYSCORE", []any{key, opt.Min, opt.Max}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "ZRANGEBYSCORE 失败")
	}

	return result, nil
}

// ZRangeByScoreWithScores 返回有序集合中指定分数区间的成员及其分数
func (c *Connection) ZRangeByScoreWithScores(ctx *ucontext.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	startTime := time.Now()
	result, err := c.client.ZRangeByScoreWithScores(ctx, key, opt).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZRANGEBYSCORE", []any{key, opt.Min, opt.Max, "WITHSCORES"}, duration, err)

	if err != nil {
		return nil, uerror.Wrap(err, "ZRANGEBYSCORE WITHSCORES 失败")
	}

	return result, nil
}

// ZCard 获取有序集合的成员数
func (c *Connection) ZCard(ctx *ucontext.Context, key string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.ZCard(ctx, key).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZCARD", []any{key}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "ZCARD 失败")
	}

	return result, nil
}

// ZScore 获取有序集合中成员的分数
func (c *Connection) ZScore(ctx *ucontext.Context, key, member string) (float64, error) {
	startTime := time.Now()
	result, err := c.client.ZScore(ctx, key, member).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZSCORE", []any{key, member}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, ErrNil
		}
		return 0, uerror.Wrap(err, "ZSCORE 失败")
	}

	return result, nil
}

// ZRank 获取有序集合中成员的排名（从小到大）
func (c *Connection) ZRank(ctx *ucontext.Context, key, member string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.ZRank(ctx, key, member).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZRANK", []any{key, member}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, ErrNil
		}
		return 0, uerror.Wrap(err, "ZRANK 失败")
	}

	return result, nil
}

// ZRevRank 获取有序集合中成员的排名（从大到小）
func (c *Connection) ZRevRank(ctx *ucontext.Context, key, member string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.ZRevRank(ctx, key, member).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZREVRANK", []any{key, member}, duration, err)

	if err != nil {
		if err.Error() == "redis: nil" {
			return 0, ErrNil
		}
		return 0, uerror.Wrap(err, "ZREVRANK 失败")
	}

	return result, nil
}

// ZIncrBy 为有序集合中成员的分数增加指定值
func (c *Connection) ZIncrBy(ctx *ucontext.Context, key string, increment float64, member string) (float64, error) {
	startTime := time.Now()
	result, err := c.client.ZIncrBy(ctx, key, increment, member).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZINCRBY", []any{key, increment, member}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "ZINCRBY 失败")
	}

	return result, nil
}

// ZCount 计算有序集合中指定分数区间的成员数量
func (c *Connection) ZCount(ctx *ucontext.Context, key, min, max string) (int64, error) {
	startTime := time.Now()
	result, err := c.client.ZCount(ctx, key, min, max).Result()
	duration := time.Since(startTime)

	c.logCommand(ctx, "ZCOUNT", []any{key, min, max}, duration, err)

	if err != nil {
		return 0, uerror.Wrap(err, "ZCOUNT 失败")
	}

	return result, nil
}
