package ucontext

import (
	"context"
	"sync"
	"time"
)

// ============================================================================
// 类型定义
// ============================================================================

// contextKey 用于在 context 中存储追踪信息的 key 类型
type contextKey string

const (
	traceContextKey contextKey = "trace_context"
)

// Context 封装的上下文
// 包含标准 context.Context 和追踪信息 TraceContext
type Context struct {
	ctx   context.Context
	trace *TraceContext
}

// TraceContext 追踪上下文
type TraceContext struct {
	TraceID       string            // 追踪 ID
	SpanID        string            // 当前 Span ID
	ParentSpanID  string            // 父 Span ID
	RequestID     string            // 请求 ID
	StartTime     time.Time         // 开始时间
	Sampled       bool              // 是否采样
	Metadata      map[string]string // 额外元数据
	MetadataMutex sync.RWMutex
}

// ============================================================================
// Context 构造函数
// ============================================================================

// New 创建新的 Context
func New() *Context {
	return &Context{
		ctx:   context.Background(),
		trace: NewTraceContext(),
	}
}

// NewWithContext 从标准 context 创建
func NewWithContext(ctx context.Context) *Context {
	if ctx == nil {
		ctx = context.Background()
	}

	trace := FromContext(ctx)
	if trace == nil {
		trace = NewTraceContext()
	}

	return &Context{
		ctx:   ctx,
		trace: trace,
	}
}

// ============================================================================
// Context 方法
// ============================================================================

// Context 获取标准 context.Context
func (c *Context) Context() context.Context {
	return c.ctx
}

// Trace 获取追踪上下文
func (c *Context) Trace() *TraceContext {
	return c.trace
}

// WithValue 设置值
func (c *Context) WithValue(key, value any) *Context {
	return &Context{
		ctx:   context.WithValue(c.ctx, key, value),
		trace: c.trace,
	}
}

// Value 获取值
func (c *Context) Value(key any) any {
	return c.ctx.Value(key)
}

// WithTimeout 设置超时
func (c *Context) WithTimeout(timeout time.Duration) (*Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	return &Context{
		ctx:   ctx,
		trace: c.trace,
	}, cancel
}

// Done 返回 done channel
func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err 返回错误
func (c *Context) Err() error {
	return c.ctx.Err()
}

// Deadline 返回截止时间
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// ============================================================================
// Context 包级函数
// ============================================================================

// WithTimeout 为 Context 设置超时（包级函数）
func WithTimeout(ctx *Context, timeout time.Duration) (*Context, context.CancelFunc) {
	stdCtx, cancel := context.WithTimeout(ctx.Context(), timeout)
	return &Context{
		ctx:   stdCtx,
		trace: ctx.trace,
	}, cancel
}

// WithCancel 创建可取消的上下文
func WithCancel(ctx *Context) (*Context, context.CancelFunc) {
	stdCtx, cancel := context.WithCancel(ctx.Context())
	return &Context{
		ctx:   stdCtx,
		trace: ctx.trace,
	}, cancel
}

// ============================================================================
// TraceContext 构造函数
// ============================================================================

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

	parent.MetadataMutex.RLock()
	metadata := copyMetadata(parent.Metadata)
	parent.MetadataMutex.RUnlock()

	return &TraceContext{
		TraceID:      parent.TraceID,
		SpanID:       GenerateID(),
		ParentSpanID: parent.SpanID,
		RequestID:    parent.RequestID,
		StartTime:    time.Now(),
		Sampled:      parent.Sampled,
		Metadata:     metadata,
	}
}

// ============================================================================
// TraceContext 方法
// ============================================================================

// SetMetadata 设置元数据
func (tc *TraceContext) SetMetadata(key, value string) {
	tc.MetadataMutex.Lock()
	defer tc.MetadataMutex.Unlock()
	if tc.Metadata == nil {
		tc.Metadata = make(map[string]string)
	}
	tc.Metadata[key] = value
}

// GetMetadata 获取元数据
func (tc *TraceContext) GetMetadata(key string) string {
	tc.MetadataMutex.RLock()
	defer tc.MetadataMutex.RUnlock()
	if tc.Metadata == nil {
		return ""
	}
	return tc.Metadata[key]
}

// Duration 获取从开始到现在的时长
func (tc *TraceContext) Duration() time.Duration {
	return time.Since(tc.StartTime)
}

// ============================================================================
// 标准 context.Context 集成
// ============================================================================

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

// ============================================================================
// 辅助函数
// ============================================================================

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
