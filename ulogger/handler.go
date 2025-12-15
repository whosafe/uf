package ulogger

import (
	"context"
	"io"
	"log/slog"
)

// customHandler 自定义 slog Handler
type customHandler struct {
	config    *Config
	writer    io.Writer
	formatter Formatter
	attrs     []slog.Attr
	groups    []string
	level     slog.Leveler
}

// newCustomHandler 创建自定义 Handler
func newCustomHandler(config *Config, writer io.Writer) *customHandler {
	// 根据配置选择格式化器
	var formatter Formatter
	switch config.Format {
	case "json":
		formatter = NewJSONFormatter()
	case "custom":
		if config.Formatter != nil {
			formatter = config.Formatter
		} else {
			// 如果没有提供自定义格式化器，回退到文本格式
			formatter = NewTextFormatter()
		}
	default: // "text" 或其他
		formatter = NewTextFormatter()
	}

	return &customHandler{
		config:    config,
		writer:    writer,
		formatter: formatter,
		level:     config.Level,
		attrs:     make([]slog.Attr, 0),
		groups:    make([]string, 0),
	}
}

// Enabled 实现 slog.Handler 接口
func (h *customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

// Handle 实现 slog.Handler 接口
func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	// 使用格式化器格式化日志
	data, err := h.formatter.Format(r, h.config)
	if err != nil {
		return err
	}

	// 添加换行符
	data = append(data, '\n')

	// 写入
	_, err = h.writer.Write(data)
	return err
}

// WithAttrs 实现 slog.Handler 接口
func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := *h
	newHandler.attrs = append(make([]slog.Attr, 0, len(h.attrs)+len(attrs)), h.attrs...)
	newHandler.attrs = append(newHandler.attrs, attrs...)
	return &newHandler
}

// WithGroup 实现 slog.Handler 接口
func (h *customHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newHandler := *h
	newHandler.groups = append(make([]string, 0, len(h.groups)+1), h.groups...)
	newHandler.groups = append(newHandler.groups, name)
	return &newHandler
}

// multiHandler 多输出 Handler
type multiHandler struct {
	handlers []slog.Handler
}

// newMultiHandler 创建多输出 Handler
func newMultiHandler(handlers ...slog.Handler) *multiHandler {
	return &multiHandler{
		handlers: handlers,
	}
}

// Enabled 实现 slog.Handler 接口
func (mh *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range mh.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle 实现 slog.Handler 接口
func (mh *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range mh.handlers {
		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// WithAttrs 实现 slog.Handler 接口
func (mh *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

// WithGroup 实现 slog.Handler 接口
func (mh *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(mh.handlers))
	for i, h := range mh.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}
