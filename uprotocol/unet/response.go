package unet

import (
	"net/http"
)

// Response 响应接口 (所有协议都要实现)
type Response interface {
	// JSON 返回 JSON 响应
	JSON(code int, data any) error

	// String 返回字符串响应
	String(code int, text string) error

	// Bytes 返回字节响应
	Bytes(code int, data []byte) error

	// HTML 返回 HTML 响应
	HTML(code int, html string) error

	// Redirect 重定向
	Redirect(code int, url string) error

	// SetHeader 设置响应头
	SetHeader(key, value string)

	// AddHeader 添加响应头
	AddHeader(key, value string)

	// SetCookie 设置 Cookie
	SetCookie(cookie *http.Cookie)

	// Status 设置状态码
	Status(code int)

	// Write 写入响应数据（实现 io.Writer）
	Write(data []byte) (int, error)

	// StatusCode 获取状态码
	StatusCode() int

	// IsWritten 检查响应是否已写入
	IsWritten() bool

	// SetSessionCookie 设置session Cookie
	SetSessionCookie(name string, id string, path string, domain string, age int, secure bool, only bool, site http.SameSite)
}
