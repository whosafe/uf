package ulogger

import (
	"context"
	"log/slog"

	"iutime.com/utime/uf/uconfig"
	"iutime.com/utime/uf/ucontext"
)

// 全局默认 Logger
var defaultLogger *Logger

// 全局配置（用于从 uconfig 加载）
var globalConfig *Config

// init 初始化默认 Logger
func init() {
	// 初始化全局配置为默认值
	globalConfig = DefaultConfig()
	// 注册到 uconfig
	Register()
	// 创建默认 Logger
	defaultLogger, _ = New(globalConfig)
}

// SetDefault 设置全局默认 Logger
func SetDefault(logger *Logger) {
	defaultLogger = logger
	slog.SetDefault(logger.slogger)
}

// Default 获取全局默认 Logger
func Default() *Logger {
	return defaultLogger
}

// DebugCtx 使用默认 Logger 输出 Debug 日志（支持 context）
func DebugCtx(ctx context.Context, msg string, args ...any) {
	if ctx != nil {
		tc := ucontext.FromContext(ctx)
		if tc != nil {
			args = append(args, "trace_id", tc.TraceID, "span_id", tc.SpanID)
		}
	}
	Debug(msg, args...)
}

// InfoCtx 使用默认 Logger 输出 Info 日志（支持 context）
func InfoCtx(ctx context.Context, msg string, args ...any) {
	if ctx != nil {
		tc := ucontext.FromContext(ctx)
		if tc != nil {
			args = append(args, "trace_id", tc.TraceID, "span_id", tc.SpanID)
		}
	}
	Info(msg, args...)
}

// WarnCtx 使用默认 Logger 输出 Warn 日志（支持 context）
func WarnCtx(ctx context.Context, msg string, args ...any) {
	if ctx != nil {
		tc := ucontext.FromContext(ctx)
		if tc != nil {
			args = append(args, "trace_id", tc.TraceID, "span_id", tc.SpanID)
		}
	}
	Warn(msg, args...)
}

// ErrorCtx 使用默认 Logger 输出 Error 日志（支持 context）
func ErrorCtx(ctx context.Context, msg string, args ...any) {
	if ctx != nil {
		tc := ucontext.FromContext(ctx)
		if tc != nil {
			args = append(args, "trace_id", tc.TraceID, "span_id", tc.SpanID)
		}
	}
	Error(msg, args...)
}

// Debug 使用默认 Logger 输出 Debug 日志（兼容旧接口）
func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

// Info 使用默认 Logger 输出 Info 日志（兼容旧接口）
func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

// Warn 使用默认 Logger 输出 Warn 日志（兼容旧接口）
func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

// Error 使用默认 Logger 输出 Error 日志（兼容旧接口）
func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

// Register 注册 logger 配置到 uconfig
// 使用方式：在 main 函数中调用 ulogger.Register()，然后调用 uconfig.Load()
// 这样配置文件中的 logger 配置会自动加载并重新初始化默认 logger
//
// 示例：
//
//	ulogger.Register()
//	uconfig.Load("config.yaml")
//	ulogger.Info("使用配置文件中的 logger 设置")
func Register() {
	// 注册到 uconfig，使用 "logger" 作为配置键
	uconfig.Register("logger", func(key string, value *uconfig.Node) error {
		// 解析配置到全局配置对象
		if err := globalConfig.UnmarshalYAML(key, value); err != nil {
			return err
		}

		// 验证配置
		if err := globalConfig.Validate(); err != nil {
			return err
		}

		// 重新创建默认 logger
		newLogger, err := New(globalConfig)
		if err != nil {
			return err
		}

		// 关闭旧的 logger
		if defaultLogger != nil {
			defaultLogger.Close()
		}

		// 更新默认 logger
		defaultLogger = newLogger

		// 同时更新 slog 的默认 logger
		slog.SetDefault(defaultLogger.slogger)

		return nil
	})
}
