package uhttp

import (
	"net/http"
	"sync"

	"github.com/whosafe/uf/uprotocol/umarshal"
)

// responsePool Response 对象池
var responsePool = sync.Pool{
	New: func() any {
		return &Response{}
	},
}

// Response HTTP 响应 (实现 unet.Response 接口)
type Response struct {
	writer     http.ResponseWriter
	statusCode int
	written    bool
}

// newResponse 创建新的 Response (从对象池获取)
func newResponse(w http.ResponseWriter) *Response {
	resp := responsePool.Get().(*Response)
	resp.writer = w
	resp.statusCode = http.StatusOK
	resp.written = false
	return resp
}

// release 释放 Response 回对象池
func (r *Response) release() {
	r.writer = nil
	r.statusCode = 0
	r.written = false
	responsePool.Put(r)
}

// JSON 返回 JSON 响应
func (r *Response) JSON(code int, data any) error {
	r.Status(code)
	r.SetHeader("Content-Type", "application/json") // Changed from r.Header to r.SetHeader for syntactic correctness

	// 使用 umarshal 序列化
	jsonData, err := umarshal.Marshal(data)
	if err != nil {
		return r.JSON(code, Error(CodeInternalError, err.Error()))
	}

	_, err = r.writer.Write(jsonData)
	return err
}

// String 返回字符串响应
func (r *Response) String(code int, text string) error {
	r.SetHeader("Content-Type", "text/plain; charset=utf-8")
	r.Status(code)

	_, err := r.writer.Write([]byte(text))
	return err
}

// Bytes 返回字节响应
func (r *Response) Bytes(code int, data []byte) error {
	r.Status(code)
	_, err := r.writer.Write(data)
	return err
}

// HTML 返回 HTML 响应
func (r *Response) HTML(code int, html string) error {
	r.SetHeader("Content-Type", "text/html; charset=utf-8")
	r.Status(code)

	_, err := r.writer.Write([]byte(html))
	return err
}

// Redirect 重定向
func (r *Response) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		code = http.StatusFound
	}
	r.SetHeader("Location", url)
	r.Status(code)
	return nil
}

// SetHeader 设置响应头
func (r *Response) SetHeader(key, value string) {
	r.writer.Header().Set(key, value)
}

// AddHeader 添加响应头
func (r *Response) AddHeader(key, value string) {
	r.writer.Header().Add(key, value)
}

// SetCookie 设置 Cookie
func (r *Response) SetCookie(cookie *http.Cookie) {
	http.SetCookie(r.writer, cookie)
}

// Status 设置状态码
func (r *Response) Status(code int) {
	if !r.written {
		r.statusCode = code
		r.writer.WriteHeader(code)
		r.written = true
	}
}

// Write 实现 io.Writer 接口
func (r *Response) Write(data []byte) (int, error) {
	if !r.written {
		r.Status(http.StatusOK)
	}
	return r.writer.Write(data)
}

// Writer 获取原始 http.ResponseWriter
func (r *Response) Writer() http.ResponseWriter {
	return r.writer
}

// StatusCode 获取状态码
func (r *Response) StatusCode() int {
	return r.statusCode
}

// IsWritten 检查响应是否已写入
func (r *Response) IsWritten() bool {
	return r.written
}
