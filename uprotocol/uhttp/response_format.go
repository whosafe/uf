package uhttp

import (
	"github.com/whosafe/uf/uprotocol/umarshal"
)

// APIResponse 统一的 API 响应格式
type APIResponse struct {
	Code    int    `json:"code"`              // 业务状态码
	Message string `json:"message,omitempty"` // 提示信息
	Data    any    `json:"data,omitempty"`    // 响应数据
	Error   string `json:"error,omitempty"`   // 错误信息
}

// Marshal 实现 umarshal.IMarshaler 接口
func (a *APIResponse) Marshal(w *umarshal.Writer) error {
	w.WriteObjectStart()
	w.WriteObjectField("code")
	w.WriteInt(a.Code)

	if a.Message != "" {
		w.WriteComma()
		w.WriteObjectField("message")
		w.WriteString(a.Message)
	}

	if a.Data != nil {
		w.WriteComma()
		w.WriteObjectField("data")
		// 如果 Data 实现了 IMarshaler,使用自定义序列化
		if m, ok := a.Data.(umarshal.IMarshaler); ok {
			m.Marshal(w)
		} else {
			// 否则使用标准序列化
			data, _ := umarshal.Marshal(a.Data)
			w.WriteRawString(string(data))
		}
	}

	if a.Error != "" {
		w.WriteComma()
		w.WriteObjectField("error")
		w.WriteString(a.Error)
	}

	w.WriteObjectEnd()
	return nil
}

// Success 成功响应
func Success(data any) *APIResponse {
	return &APIResponse{
		Code: 0,
		Data: data,
	}
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(message string, data any) *APIResponse {
	return &APIResponse{
		Code:    0,
		Message: message,
		Data:    data,
	}
}

// Error 错误响应
func Error(code int, message string) *APIResponse {
	return &APIResponse{
		Code:  code,
		Error: message,
	}
}

// ErrorWithData 带数据的错误响应
func ErrorWithData(code int, message string, data any) *APIResponse {
	return &APIResponse{
		Code:  code,
		Error: message,
		Data:  data,
	}
}

// 常用业务错误码
const (
	CodeSuccess           = 0     // 成功
	CodeInvalidParams     = 10001 // 参数错误
	CodeNotFound          = 10002 // 资源不存在
	CodeUnauthorized      = 10003 // 未授权
	CodeForbidden         = 10004 // 禁止访问
	CodeInternalError     = 10005 // 内部错误
	CodeDatabaseError     = 10006 // 数据库错误
	CodeValidationError   = 10007 // 验证失败
	CodeDuplicateError    = 10008 // 重复数据
	CodeRateLimitExceeded = 10009 // 超过限流
)

// Response 辅助方法 - 成功响应
func (r *Response) Success(data any) error {
	return r.JSON(200, Success(data))
}

// SuccessWithMessage 辅助方法 - 带消息的成功响应
func (r *Response) SuccessWithMessage(message string, data any) error {
	return r.JSON(200, SuccessWithMessage(message, data))
}

// Error 辅助方法 - 错误响应
func (r *Response) Error(httpCode int, bizCode int, message string) error {
	return r.JSON(httpCode, Error(bizCode, message))
}

// BadRequest 400 错误
func (r *Response) BadRequest(message string) error {
	return r.JSON(400, Error(CodeInvalidParams, message))
}

// NotFound 404 错误
func (r *Response) NotFound(message string) error {
	return r.JSON(404, Error(CodeNotFound, message))
}

// Unauthorized 401 错误
func (r *Response) Unauthorized(message string) error {
	return r.JSON(401, Error(CodeUnauthorized, message))
}

// Forbidden 403 错误
func (r *Response) Forbidden(message string) error {
	return r.JSON(403, Error(CodeForbidden, message))
}

// InternalError 500 错误
func (r *Response) InternalError(message string) error {
	return r.JSON(500, Error(CodeInternalError, message))
}
