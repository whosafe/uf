package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/ulogger"
)

func main() {
	println("=== ucontext 链路追踪演示 ===\n")

	// 初始化雪花算法（worker ID = 1）
	ucontext.InitSnowflake(1)

	// 设置采样率为 100%
	ucontext.SetSamplingRate(1.0)

	// 1. 基本使用
	println("1. 基本使用")
	println("-------------------")
	basicUsage()
	println()

	// 2. 嵌套 Span
	println("2. 嵌套 Span（父子关系）")
	println("-------------------")
	nestedSpan()
	println()

	// 3. Logger 集成
	println("3. Logger 集成")
	println("-------------------")
	loggerIntegration()
	println()

	// 4. HTTP 传播
	println("4. HTTP 传播")
	println("-------------------")
	httpPropagation()
	println()

	// 5. 采样控制
	println("5. 采样控制")
	println("-------------------")
	samplingControl()
}

// basicUsage 基本使用示例
func basicUsage() {
	// 创建追踪上下文
	ctx := ucontext.NewContext(context.Background())
	tc := ucontext.FromContext(ctx)

	fmt.Printf("Trace ID:   %s\n", tc.TraceID)
	fmt.Printf("Span ID:    %s\n", tc.SpanID)
	fmt.Printf("Request ID: %s\n", tc.RequestID)
	fmt.Printf("Sampled:    %v\n", tc.Sampled)

	// 设置元数据
	tc.SetMetadata("user", "alice")
	tc.SetMetadata("action", "login")
	fmt.Printf("Metadata:   user=%s, action=%s\n",
		tc.GetMetadata("user"),
		tc.GetMetadata("action"))
}

// nestedSpan 嵌套 Span 示例
func nestedSpan() {
	// 父操作
	parentCtx := ucontext.NewContext(context.Background())
	parentTC := ucontext.FromContext(parentCtx)
	fmt.Printf("父 Span - Trace ID: %s, Span ID: %s\n",
		parentTC.TraceID, parentTC.SpanID)

	// 子操作 1
	child1Ctx := ucontext.NewSpan(parentCtx)
	child1TC := ucontext.FromContext(child1Ctx)
	fmt.Printf("子 Span 1 - Trace ID: %s, Span ID: %s, Parent: %s\n",
		child1TC.TraceID, child1TC.SpanID, child1TC.ParentSpanID)

	// 子操作 2
	child2Ctx := ucontext.NewSpan(parentCtx)
	child2TC := ucontext.FromContext(child2Ctx)
	fmt.Printf("子 Span 2 - Trace ID: %s, Span ID: %s, Parent: %s\n",
		child2TC.TraceID, child2TC.SpanID, child2TC.ParentSpanID)

	// 验证 Trace ID 相同
	if child1TC.TraceID == child2TC.TraceID && child1TC.TraceID == parentTC.TraceID {
		fmt.Println("✓ 所有 Span 的 Trace ID 相同")
	}
}

// loggerIntegration Logger 集成示例
func loggerIntegration() {
	// 创建带追踪信息的 context
	ctx := ucontext.NewContext(context.Background())

	// 方式 2: 使用全局 Logger 的 Context 方法
	ulogger.InfoCtx(ctx, "使用 InfoCtx 记录日志", "key", "value")
	ulogger.DebugCtx(ctx, "使用 DebugCtx 记录日志", "key", "value")
	ulogger.WarnCtx(ctx, "使用 WarnCtx 记录日志", "key", "value")
	ulogger.ErrorCtx(ctx, "使用 ErrorCtx 记录日志", "key", "value")
}

// httpPropagation HTTP 传播示例
func httpPropagation() {
	// 创建追踪上下文
	tc := ucontext.NewTraceContext()
	fmt.Printf("原始 Trace ID: %s\n", tc.TraceID)

	// 模拟注入到 HTTP Header
	header := http.Header{}
	ucontext.InjectHTTPHeaders(header, tc)
	fmt.Printf("注入 Header: X-Trace-ID=%s\n", header.Get("X-Trace-ID"))

	// 模拟从 HTTP Header 提取
	extractedTC := ucontext.ExtractHTTPHeaders(header)
	fmt.Printf("提取 Trace ID: %s\n", extractedTC.TraceID)
	fmt.Printf("新 Span ID: %s (自动生成)\n", extractedTC.SpanID)
	fmt.Printf("Parent Span ID: %s (来自上游)\n", extractedTC.ParentSpanID)

	if extractedTC.TraceID == tc.TraceID {
		fmt.Println("✓ Trace ID 传播成功")
	}
}

// samplingControl 采样控制示例
func samplingControl() {
	// 设置 50% 采样率
	ucontext.SetSamplingRate(0.5)
	fmt.Printf("设置采样率: %.0f%%\n", ucontext.GetSamplingRate()*100)

	// 创建多个上下文，统计采样情况
	sampled := 0
	total := 100
	for i := 0; i < total; i++ {
		tc := ucontext.NewTraceContext()
		if tc.Sampled {
			sampled++
		}
	}
	fmt.Printf("采样结果: %d/%d (%.0f%%)\n", sampled, total, float64(sampled)/float64(total)*100)

	// 强制采样
	ucontext.SetSamplingRate(0.0) // 设置为 0%
	ctx := ucontext.NewContext(context.Background())
	fmt.Printf("0%% 采样率，是否采样: %v\n", ucontext.IsSampled(ctx))

	ctx = ucontext.ForceSample(ctx)
	fmt.Printf("强制采样后，是否采样: %v\n", ucontext.IsSampled(ctx))

	// 恢复 100% 采样
	ucontext.SetSamplingRate(1.0)
}
