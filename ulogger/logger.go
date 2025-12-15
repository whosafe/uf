package ulogger

import (
	"context"
	"io"
	"log/slog"
	"os"

	"iutime.com/utime/uf/ucontext"
)

// Logger 日志器
type Logger struct {
	config       *Config
	slogger      *slog.Logger
	rotateWriter *rotateWriter
}

// New 创建新的日志器
func New(config *Config) (*Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, err
	}

	logger := &Logger{
		config: config,
	}

	// 创建 handlers
	var handlers []slog.Handler

	// 文件输出
	if config.Path != "" {
		rw, err := newRotateWriter(config)
		if err != nil {
			return nil, err
		}
		logger.rotateWriter = rw
		fileHandler := newCustomHandler(config, rw)
		handlers = append(handlers, fileHandler)
	}

	// 终端输出
	if config.Stdout {
		stdoutHandler := newCustomHandler(config, os.Stdout)
		handlers = append(handlers, stdoutHandler)
	}

	// 如果没有任何输出，默认输出到终端
	if len(handlers) == 0 {
		stdoutHandler := newCustomHandler(config, os.Stdout)
		handlers = append(handlers, stdoutHandler)
	}

	// 创建 handler
	var handler slog.Handler
	if len(handlers) == 1 {
		handler = handlers[0]
	} else {
		handler = newMultiHandler(handlers...)
	}

	// 创建 slog.Logger
	logger.slogger = slog.New(handler)

	return logger, nil
}

// Close 关闭日志器
func (l *Logger) Close() error {
	if l.rotateWriter != nil {
		return l.rotateWriter.Close()
	}
	return nil
}

// Sync 同步日志到磁盘
func (l *Logger) Sync() error {
	if l.rotateWriter != nil {
		return l.rotateWriter.Sync()
	}
	return nil
}

// Logger 获取底层 slog.Logger
func (l *Logger) Logger() *slog.Logger {
	return l.slogger
}

// DebugCtx 输出 Debug 级别日志（支持 context）
func (l *Logger) DebugCtx(ctx context.Context, msg string, args ...any) {
	if ctx != nil {
		tc := ucontext.FromContext(ctx)
		if tc != nil {
			args = append(args, "trace_id", tc.TraceID, "span_id", tc.SpanID)
		}
	}
	l.Debug(msg, args...)
}

// InfoCtx 输出 Info 级别日志（支持 context）
func (l *Logger) InfoCtx(ctx context.Context, msg string, args ...any) {
	if ctx != nil {
		tc := ucontext.FromContext(ctx)
		if tc != nil {
			args = append(args, "trace_id", tc.TraceID, "span_id", tc.SpanID)
		}
	}
	l.Info(msg, args...)
}

// WarnCtx 输出 Warn 级别日志（支持 context）
func (l *Logger) WarnCtx(ctx context.Context, msg string, args ...any) {
	if ctx != nil {
		tc := ucontext.FromContext(ctx)
		if tc != nil {
			args = append(args, "trace_id", tc.TraceID, "span_id", tc.SpanID)
		}
	}
	l.Warn(msg, args...)
}

// ErrorCtx 输出 Error 级别日志（支持 context）
func (l *Logger) ErrorCtx(ctx context.Context, msg string, args ...any) {
	if ctx != nil {
		tc := ucontext.FromContext(ctx)
		if tc != nil {
			args = append(args, "trace_id", tc.TraceID, "span_id", tc.SpanID)
		}
	}
	l.Error(msg, args...)
}

// Debug 输出 Debug 级别日志（兼容旧接口）
func (l *Logger) Debug(msg string, args ...any) {
	l.slogger.Debug(msg, args...)
}

// Info 输出 Info 级别日志（兼容旧接口）
func (l *Logger) Info(msg string, args ...any) {
	l.slogger.Info(msg, args...)
}

// Warn 输出 Warn 级别日志（兼容旧接口）
func (l *Logger) Warn(msg string, args ...any) {
	l.slogger.Warn(msg, args...)
}

// Error 输出 Error 级别日志（兼容旧接口）
func (l *Logger) Error(msg string, args ...any) {
	l.slogger.Error(msg, args...)
}

// Log 输出指定级别的日志
func (l *Logger) Log(level slog.Level, msg string, args ...any) {
	l.slogger.Log(nil, level, msg, args...)
}

// LogContext 输出指定级别的日志（支持 context）
func (l *Logger) LogContext(ctx context.Context, level slog.Level, msg string, args ...any) {
	l.slogger.Log(ctx, level, msg, args...)
}

// With 创建带有额外属性的子 Logger
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		config:       l.config,
		slogger:      l.slogger.With(args...),
		rotateWriter: l.rotateWriter,
	}
}

// WithGroup 创建带有分组的子 Logger
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		config:       l.config,
		slogger:      l.slogger.WithGroup(name),
		rotateWriter: l.rotateWriter,
	}
}

// Writer 返回一个 io.Writer，可以用于其他需要 io.Writer 的场景
func (l *Logger) Writer() io.Writer {
	return &logWriter{logger: l}
}
