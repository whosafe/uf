package uerror

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

// Error 是一个包含错误码和消息的自定义错误类型。
type Error struct {
	Code    int
	Message string
	Err     error
	stack   []uintptr
}

// New 创建一个默认错误码为 0 的新 Error。
func New(message string) *Error {
	return &Error{
		Code:    0,
		Message: message,
		stack:   callers(),
	}
}

// NewWithCode 创建一个指定错误码的新 Error。
func NewWithCode(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		stack:   callers(),
	}
}

// Wrap 使用消息包装现有错误，默认错误码为 0。
func Wrap(err error, message string) *Error {
	return &Error{
		Code:    0,
		Message: message,
		Err:     err,
		stack:   callers(),
	}
}

// WrapWithCode 使用指定错误码和消息包装现有错误。
func WrapWithCode(code int, err error, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
		stack:   callers(),
	}
}

// Error 返回错误的字符串表示形式。
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 返回底层错误。
func (e *Error) Unwrap() error {
	return e.Err
}

// Format 实现 fmt.Formatter 接口以允许自定义格式化。
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "Code: %d\nMessage: %s\n", e.Code, e.Message)
			if e.Err != nil {
				fmt.Fprintf(s, "Cause: %+v\n", e.Err)
			}
			e.formatStack(s)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

// formatStack 格式化堆栈信息到 fmt.State
func (e *Error) formatStack(s fmt.State) {
	if len(e.stack) == 0 {
		return
	}
	frames := runtime.CallersFrames(e.stack)
	for {
		frame, more := frames.Next()
		if !strings.HasPrefix(frame.Function, "runtime.") {
			fmt.Fprintf(s, "%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line)
		}
		if !more {
			break
		}
	}
}

// callers 获取调用堆栈
func callers() []uintptr {
	const depth = 64
	var pcs [depth]uintptr
	// 跳过 callers, New/Wrap, 及其调用者 (共3层)
	n := runtime.Callers(3, pcs[:])
	return pcs[0:n]
}
