package unet

import (
	"net"
	"net/http"
	"net/url"

	"github.com/whosafe/uf/uprotocol/ubind"
)

// Request 请求接口 (所有协议都要实现)
type Request interface {
	// Protocol 获取协议类型
	Protocol() Protocol

	// RemoteAddr 获取远程地址
	RemoteAddr() net.Addr

	// LocalAddr 获取本地地址
	LocalAddr() net.Addr

	// Get/Set 键值对存储
	Get(key string) (any, bool)
	Set(key string, value any)

	// Bind 绑定请求数据
	Bind(obj ubind.Binder) error

	// Response 获取响应接口
	Response() Response

	// Param 获取路径参数
	Param(key string) string

	// Query 获取查询参数
	Query(key string) string

	// QueryDefault 获取查询参数，如果不存在则返回默认值
	QueryDefault(key, def string) string

	// Header 获取请求头
	Header(key string) string

	// Method 获取请求方法
	Method() string

	// Path 获取请求路径
	Path() string

	// Body 获取请求体
	Body() ([]byte, error)

	// Cookie 获取 Cookie
	Cookie(name string) (*http.Cookie, error)

	// BindJSON 绑定 JSON 数据
	BindJSON(obj ubind.Binder) error

	// BindForm 绑定 Form 数据
	BindForm(obj ubind.Binder) error

	// BindQuery 绑定查询参数
	BindQuery(obj ubind.Binder) error

	// URL 获取完整 URL
	URL() *url.URL

	// Session 获取 Session
	Session() (Session, error)
}
