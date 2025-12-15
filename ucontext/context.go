package ucontext

import (
	"context"
	"time"
)

// contextKey 用于在 context 中存储追踪信息的 key 类型
type contextKey string

const (
	traceContextKey contextKey = "trace_context"
)

// TraceContext 追踪上下文
type TraceContext struct {
	TraceID      string            // 追踪 ID
	SpanID       string            // 当前 Span ID
	ParentSpanID string            // 父 Span ID
	RequestID    string            // 请求 ID
	StartTime    time.Time         // 开始时间
	Sampled      bool              // 是否采样
	Metadata     map[string]string // 额外元数据
}

// NewTraceContext 创建新的追踪上下文
func NewTraceContext() *TraceContext {
	return &TraceContext{
		TraceID:   GenerateID(),
		SpanID:    GenerateID(),
		RequestID: GenerateID(),
		StartTime: time.Now(),
		Sampled:   shouldSample(),
		Metadata:  make(map[string]string),
	}
}

// NewSpanContext 创建子 Span 上下文
func NewSpanContext(parent *TraceContext) *TraceContext {
	if parent == nil {
		return NewTraceContext()
	}

	return &TraceContext{
		TraceID:      parent.TraceID,
		SpanID:       GenerateID(),
		ParentSpanID: parent.SpanID,
		RequestID:    parent.RequestID,
		StartTime:    time.Now(),
		Sampled:      parent.Sampled,
		Metadata:     copyMetadata(parent.Metadata),
	}
}

// WithContext 将追踪上下文注入 context.Context
func WithContext(ctx context.Context, tc *TraceContext) context.Context {
	if tc == nil {
		tc = NewTraceContext()
	}
	return context.WithValue(ctx, traceContextKey, tc)
}

// FromContext 从 context.Context 提取追踪上下文
func FromContext(ctx context.Context) *TraceContext {
	if ctx == nil {
		return nil
	}

	tc, ok := ctx.Value(traceContextKey).(*TraceContext)
	if !ok {
		return nil
	}
	return tc
}

// NewContext 创建带追踪信息的新 context
func NewContext(parent context.Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	return WithContext(parent, NewTraceContext())
}

// NewSpan 创建子 Span 的 context
func NewSpan(parent context.Context) context.Context {
	if parent == nil {
		return NewContext(context.Background())
	}

	parentTC := FromContext(parent)
	childTC := NewSpanContext(parentTC)
	return WithContext(parent, childTC)
}

// SetMetadata 设置元数据
func (tc *TraceContext) SetMetadata(key, value string) {
	if tc.Metadata == nil {
		tc.Metadata = make(map[string]string)
	}
	tc.Metadata[key] = value
}

// GetMetadata 获取元数据
func (tc *TraceContext) GetMetadata(key string) string {
	if tc.Metadata == nil {
		return ""
	}
	return tc.Metadata[key]
}

// Duration 获取从开始到现在的时长
func (tc *TraceContext) Duration() time.Duration {
	return time.Since(tc.StartTime)
}

// copyMetadata 复制元数据
func copyMetadata(src map[string]string) map[string]string {
	if src == nil {
		return make(map[string]string)
	}

	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
