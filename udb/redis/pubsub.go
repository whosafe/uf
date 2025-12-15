package redis

import (
	"github.com/redis/go-redis/v9"
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uerror"
)

// ==================== Pub/Sub 操作 ====================

// Publish 发布消息到指定频道
func (c *Connection) Publish(ctx *ucontext.Context, channel string, message any) (int64, error) {
	result, err := c.client.Publish(ctx, channel, message).Result()
	if err != nil {
		return 0, uerror.Wrap(err, "发布消息失败")
	}

	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.InfoCtx(ctx, "发布消息",
			"channel", channel,
			"subscribers", result)
	}

	return result, nil
}

// Subscribe 订阅一个或多个频道
func (c *Connection) Subscribe(ctx *ucontext.Context, channels ...string) *redis.PubSub {
	pubsub := c.client.Subscribe(ctx, channels...)

	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.InfoCtx(ctx, "订阅频道",
			"channels", channels)
	}

	return pubsub
}

// PSubscribe 订阅一个或多个模式
func (c *Connection) PSubscribe(ctx *ucontext.Context, patterns ...string) *redis.PubSub {
	pubsub := c.client.PSubscribe(ctx, patterns...)

	if c.config.Log != nil && c.config.Log.Enabled {
		c.logger.InfoCtx(ctx, "订阅模式",
			"patterns", patterns)
	}

	return pubsub
}
