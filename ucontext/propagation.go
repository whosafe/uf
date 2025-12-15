package ucontext

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTP Header 常量
const (
	HeaderTraceID      = "X-Trace-ID"
	HeaderSpanID       = "X-Span-ID"
	HeaderParentSpanID = "X-Parent-Span-ID"
	HeaderRequestID    = "X-Request-ID"
	HeaderSampled      = "X-Sampled"
)

// InjectHTTPHeaders 将追踪信息注入 HTTP Header
func InjectHTTPHeaders(header http.Header, tc *TraceContext) {
	if tc == nil {
		return
	}

	header.Set(HeaderTraceID, tc.TraceID)
	header.Set(HeaderSpanID, tc.SpanID)
	if tc.ParentSpanID != "" {
		header.Set(HeaderParentSpanID, tc.ParentSpanID)
	}
	header.Set(HeaderRequestID, tc.RequestID)
	if tc.Sampled {
		header.Set(HeaderSampled, "1")
	} else {
		header.Set(HeaderSampled, "0")
	}
}

// ExtractHTTPHeaders 从 HTTP Header 提取追踪信息
func ExtractHTTPHeaders(header http.Header) *TraceContext {
	traceID := header.Get(HeaderTraceID)
	if traceID == "" {
		// 如果没有 Trace ID，创建新的追踪上下文
		return NewTraceContext()
	}

	tc := &TraceContext{
		TraceID:      traceID,
		SpanID:       GenerateID(),             // 创建新的 Span ID
		ParentSpanID: header.Get(HeaderSpanID), // 上游的 Span ID 成为当前的 Parent Span ID
		RequestID:    header.Get(HeaderRequestID),
		StartTime:    time.Now(),
		Sampled:      parseSampled(header.Get(HeaderSampled)),
		Metadata:     make(map[string]string),
	}

	// 如果没有 Request ID，生成新的
	if tc.RequestID == "" {
		tc.RequestID = GenerateID()
	}

	return tc
}

// HTTPMiddleware HTTP 中间件，自动处理追踪上下文
func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从 Header 提取或创建追踪上下文
		tc := ExtractHTTPHeaders(r.Header)

		// 注入到 request context
		ctx := WithContext(r.Context(), tc)
		r = r.WithContext(ctx)

		// 将追踪信息添加到响应 Header
		InjectHTTPHeaders(w.Header(), tc)

		next.ServeHTTP(w, r)
	})
}

// parseSampled 解析采样标志
func parseSampled(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return shouldSample()
	}

	sampled, err := strconv.ParseBool(s)
	if err != nil {
		// 尝试解析为数字
		if s == "1" {
			return true
		}
		return false
	}
	return sampled
}
