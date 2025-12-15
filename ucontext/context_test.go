package ucontext

import (
	"context"
	"testing"
	"time"
)

// TestSnowflake 测试雪花算法
func TestSnowflake(t *testing.T) {
	sf, err := NewSnowflake(1)
	if err != nil {
		t.Fatalf("Failed to create snowflake: %v", err)
	}

	// 生成多个 ID，检查唯一性
	ids := make(map[int64]bool)
	for i := 0; i < 1000; i++ {
		id := sf.Generate()
		if ids[id] {
			t.Errorf("Duplicate ID generated: %d", id)
		}
		ids[id] = true
	}
}

// TestTraceContext 测试追踪上下文
func TestTraceContext(t *testing.T) {
	// 创建追踪上下文
	tc := NewTraceContext()

	if tc.TraceID == "" {
		t.Error("TraceID should not be empty")
	}
	if tc.SpanID == "" {
		t.Error("SpanID should not be empty")
	}
	if tc.RequestID == "" {
		t.Error("RequestID should not be empty")
	}
}

// TestContextPropagation 测试上下文传播
func TestContextPropagation(t *testing.T) {
	// 创建带追踪信息的 context
	ctx := NewContext(context.Background())

	// 提取追踪信息
	tc := FromContext(ctx)
	if tc == nil {
		t.Fatal("TraceContext should not be nil")
	}

	if tc.TraceID == "" {
		t.Error("TraceID should not be empty")
	}
}

// TestNewSpan 测试创建子 Span
func TestNewSpan(t *testing.T) {
	// 创建父 context
	parentCtx := NewContext(context.Background())
	parentTC := FromContext(parentCtx)

	// 创建子 Span
	childCtx := NewSpan(parentCtx)
	childTC := FromContext(childCtx)

	if childTC == nil {
		t.Fatal("Child TraceContext should not be nil")
	}

	// 验证 Trace ID 相同
	if childTC.TraceID != parentTC.TraceID {
		t.Error("Child should have same TraceID as parent")
	}

	// 验证 Span ID 不同
	if childTC.SpanID == parentTC.SpanID {
		t.Error("Child should have different SpanID from parent")
	}

	// 验证父子关系
	if childTC.ParentSpanID != parentTC.SpanID {
		t.Error("Child's ParentSpanID should be parent's SpanID")
	}
}

// TestSampling 测试采样
func TestSampling(t *testing.T) {
	// 设置 50% 采样率
	SetSamplingRate(0.5)

	if GetSamplingRate() != 0.5 {
		t.Errorf("Expected sampling rate 0.5, got %f", GetSamplingRate())
	}

	// 生成多个上下文，检查采样率
	sampled := 0
	total := 1000
	for i := 0; i < total; i++ {
		tc := NewTraceContext()
		if tc.Sampled {
			sampled++
		}
	}

	// 允许 10% 的误差
	expectedMin := int(float64(total) * 0.4)
	expectedMax := int(float64(total) * 0.6)

	if sampled < expectedMin || sampled > expectedMax {
		t.Errorf("Sampling rate out of expected range: got %d/%d", sampled, total)
	}

	// 恢复 100% 采样
	SetSamplingRate(1.0)
}

// TestForceSample 测试强制采样
func TestForceSample(t *testing.T) {
	// 设置 0% 采样率
	SetSamplingRate(0.0)

	// 创建上下文（应该不采样）
	ctx := NewContext(context.Background())
	if IsSampled(ctx) {
		t.Error("Should not be sampled with 0% rate")
	}

	// 强制采样
	ctx = ForceSample(ctx)
	if !IsSampled(ctx) {
		t.Error("Should be sampled after ForceSample")
	}

	// 恢复 100% 采样
	SetSamplingRate(1.0)
}

// TestMetadata 测试元数据
func TestMetadata(t *testing.T) {
	tc := NewTraceContext()

	tc.SetMetadata("user", "alice")
	tc.SetMetadata("action", "login")

	if tc.GetMetadata("user") != "alice" {
		t.Error("Metadata 'user' should be 'alice'")
	}

	if tc.GetMetadata("action") != "login" {
		t.Error("Metadata 'action' should be 'login'")
	}
}

// TestDuration 测试时长计算
func TestDuration(t *testing.T) {
	tc := NewTraceContext()

	time.Sleep(10 * time.Millisecond)

	duration := tc.Duration()
	if duration < 10*time.Millisecond {
		t.Errorf("Duration should be at least 10ms, got %v", duration)
	}
}
