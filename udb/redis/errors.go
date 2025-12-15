package redis

import "github.com/whosafe/uf/uerror"

var (
	// ErrNil Redis 返回 nil 错误
	ErrNil = uerror.New("Redis 键不存在")

	// ErrClosed 连接已关闭
	ErrClosed = uerror.New("Redis 客户端已关闭")
)
