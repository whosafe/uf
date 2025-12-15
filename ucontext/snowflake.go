package ucontext

import (
	"fmt"
	"sync"
	"time"

	"github.com/whosafe/uf/uerror"
)

// Snowflake ID 生成器
// 64位ID结构: 1位符号位(0) + 41位时间戳 + 10位机器ID + 12位序列号

const (
	epoch          = int64(1640995200000) // 2022-01-01 00:00:00 UTC (毫秒)
	workerIDBits   = uint(10)             // 机器ID位数
	sequenceBits   = uint(12)             // 序列号位数
	workerIDShift  = sequenceBits
	timestampShift = sequenceBits + workerIDBits
	sequenceMask   = int64(-1) ^ (int64(-1) << sequenceBits)
	maxWorkerID    = int64(-1) ^ (int64(-1) << workerIDBits)
)

// Snowflake 雪花算法 ID 生成器
type Snowflake struct {
	mu            sync.Mutex
	workerID      int64
	sequence      int64
	lastTimestamp int64
}

// NewSnowflake 创建雪花算法生成器
func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, uerror.New(fmt.Sprintf("worker ID must be between 0 and %d", maxWorkerID))
	}

	return &Snowflake{
		workerID:      workerID,
		sequence:      0,
		lastTimestamp: 0,
	}, nil
}

// Generate 生成唯一 ID
func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamp := time.Now().UnixNano() / 1e6 // 转换为毫秒

	if timestamp < s.lastTimestamp {
		// 时钟回拨，等待
		timestamp = s.lastTimestamp
	}

	if timestamp == s.lastTimestamp {
		// 同一毫秒内，序列号递增
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for timestamp <= s.lastTimestamp {
				timestamp = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 新的毫秒，序列号重置
		s.sequence = 0
	}

	s.lastTimestamp = timestamp

	// 组装 ID
	id := ((timestamp - epoch) << timestampShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id
}

// GenerateString 生成字符串格式的 ID
func (s *Snowflake) GenerateString() string {
	return fmt.Sprintf("%d", s.Generate())
}

// 全局雪花算法生成器
var (
	globalSnowflake *Snowflake
	snowflakeOnce   sync.Once
)

// InitSnowflake 初始化全局雪花算法生成器
func InitSnowflake(workerID int64) error {
	var err error
	snowflakeOnce.Do(func() {
		globalSnowflake, err = NewSnowflake(workerID)
	})
	return err
}

// GenerateID 使用全局生成器生成 ID
func GenerateID() string {
	if globalSnowflake == nil {
		// 如果未初始化，使用默认 worker ID 0
		InitSnowflake(0)
	}
	return globalSnowflake.GenerateString()
}
