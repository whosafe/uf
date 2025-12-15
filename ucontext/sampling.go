package ucontext

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

// 采样配置
var (
	samplingRate float64 = 1.0 // 默认 100% 采样
	samplingMu   sync.RWMutex
	rng          = rand.New(rand.NewSource(time.Now().UnixNano()))
	rngMu        sync.Mutex
)

// SetSamplingRate 设置采样率
// rate 范围: 0.0 - 1.0
// 0.0 表示不采样，1.0 表示全部采样
func SetSamplingRate(rate float64) {
	if rate < 0 {
		rate = 0
	}
	if rate > 1 {
		rate = 1
	}

	samplingMu.Lock()
	samplingRate = rate
	samplingMu.Unlock()
}

// GetSamplingRate 获取当前采样率
func GetSamplingRate() float64 {
	samplingMu.RLock()
	defer samplingMu.RUnlock()
	return samplingRate
}

// shouldSample 判断是否应该采样
func shouldSample() bool {
	samplingMu.RLock()
	rate := samplingRate
	samplingMu.RUnlock()

	if rate >= 1.0 {
		return true
	}
	if rate <= 0.0 {
		return false
	}

	rngMu.Lock()
	sample := rng.Float64() < rate
	rngMu.Unlock()

	return sample
}

// ForceSample 强制采样（忽略采样率）
func ForceSample(ctx context.Context) context.Context {
	tc := FromContext(ctx)
	if tc == nil {
		tc = NewTraceContext()
	}
	tc.Sampled = true
	return WithContext(ctx, tc)
}

// IsSampled 检查是否被采样
func IsSampled(ctx context.Context) bool {
	tc := FromContext(ctx)
	if tc == nil {
		return false
	}
	return tc.Sampled
}
