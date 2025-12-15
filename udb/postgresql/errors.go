package postgresql

import "errors"

// 错误定义
var (
	// ErrNoRows 没有找到记录
	ErrNoRows = errors.New("no rows in result set")

	// ErrInvalidConfig 无效的配置
	ErrInvalidConfig = errors.New("invalid config")

	// ErrConnectionFailed 连接失败
	ErrConnectionFailed = errors.New("connection failed")

	// ErrQueryTimeout 查询超时
	ErrQueryTimeout = errors.New("query timeout")
)
