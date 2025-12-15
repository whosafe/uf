package uhttp

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/whosafe/uf/uprotocol/ubind"
	"github.com/whosafe/uf/uprotocol/unet"
)

// 错误定义
var (
	// ErrRequestBodyTooLarge 请求体过大
	ErrRequestBodyTooLarge = errors.New("request body too large")
)

// requestPool Request 对象池
var requestPool = sync.Pool{
	New: func() any {
		return &Request{
			store: make(map[string]any),
		}
	},
}

// Request HTTP 请求 (实现 unet.Request 接口)
type Request struct {
	raw      *http.Request
	writer   http.ResponseWriter
	params   Params
	query    url.Values
	store    map[string]any
	response *Response
	server   *Server
}

// Params 路径参数
type Params map[string]string

// newRequest 创建新的 Request (从对象池获取)
func newRequest(r *http.Request, w http.ResponseWriter, server *Server) *Request {
	req := requestPool.Get().(*Request)
	req.raw = r
	req.writer = w
	req.params = make(Params)
	req.query = nil
	// 清空 store
	for k := range req.store {
		delete(req.store, k)
	}
	req.response = newResponse(w)
	req.server = server
	return req
}

// release 释放 Request 回对象池
func (r *Request) release() {
	r.raw = nil
	r.writer = nil
	r.params = nil
	r.query = nil
	r.response = nil
	r.server = nil
	requestPool.Put(r)
}

// Protocol 获取协议类型
func (r *Request) Protocol() unet.Protocol {
	if r.raw.TLS != nil {
		return unet.ProtocolHTTPS
	}
	return unet.ProtocolHTTP
}

// RemoteAddr 获取远程地址
func (r *Request) RemoteAddr() net.Addr {
	addr, _ := net.ResolveTCPAddr("tcp", r.raw.RemoteAddr)
	return addr
}

// LocalAddr 获取本地地址
func (r *Request) LocalAddr() net.Addr {
	addr, _ := net.ResolveTCPAddr("tcp", r.raw.Host)
	return addr
}

// Get 获取存储的值
func (r *Request) Get(key string) (any, bool) {
	val, ok := r.store[key]
	return val, ok
}

// setServer 设置 Server
func (r *Request) setServer(server *Server) {
	r.server = server
}

// server 获取 Server
func (r *Request) Server() *Server {
	return r.server
}

// Set 设置存储的值
func (r *Request) Set(key string, value any) {
	r.store[key] = value
}

// Bind 绑定请求数据
func (r *Request) Bind(obj ubind.Binder) error {
	// 【安全修复】限制请求体大小,防止 DoS 攻击
	maxBodySize := r.getMaxBodySize()
	limitedReader := io.LimitReader(r.raw.Body, maxBodySize)

	// 读取请求体
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return err
	}
	defer r.raw.Body.Close()

	// 检查是否超过限制
	if int64(len(body)) >= maxBodySize {
		return ErrRequestBodyTooLarge
	}

	// 解析数据
	val := ubind.Parse(body)
	return ubind.Bind(val, obj)
}

// Response 获取响应接口
func (r *Request) Response() unet.Response {
	return r.response
}

// Param 获取路径参数
func (r *Request) Param(key string) string {
	return r.params[key]
}

// Query 获取查询参数
func (r *Request) Query(key string) string {
	if r.query == nil {
		r.query = r.raw.URL.Query()
	}
	return r.query.Get(key)
}

// QueryDefault 获取查询参数,如果不存在则返回默认值
func (r *Request) QueryDefault(key, def string) string {
	val := r.Query(key)
	if val == "" {
		return def
	}
	return val
}

// Header 获取请求头
func (r *Request) Header(key string) string {
	return r.raw.Header.Get(key)
}

// Cookie 获取 Cookie
func (r *Request) Cookie(name string) (*http.Cookie, error) {
	return r.raw.Cookie(name)
}

// Body 获取请求体
func (r *Request) Body() ([]byte, error) {
	// 【安全修复】限制请求体大小,防止 DoS 攻击
	maxBodySize := r.getMaxBodySize()
	limitedReader := io.LimitReader(r.raw.Body, maxBodySize)

	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, err
	}

	// 检查是否超过限制
	if int64(len(body)) >= maxBodySize {
		return nil, ErrRequestBodyTooLarge
	}

	return body, nil
}

// BindJSON 绑定 JSON 数据
func (r *Request) BindJSON(obj ubind.Binder) error {
	// 【安全修复】限制请求体大小,防止 DoS 攻击
	maxBodySize := r.getMaxBodySize()
	limitedReader := io.LimitReader(r.raw.Body, maxBodySize)

	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return err
	}
	defer r.raw.Body.Close()

	// 检查是否超过限制
	if int64(len(body)) >= maxBodySize {
		return ErrRequestBodyTooLarge
	}

	val := ubind.ParseJSON(body)
	return ubind.Bind(val, obj)
}

// BindForm 绑定 Form 数据
func (r *Request) BindForm(obj ubind.Binder) error {
	if err := r.raw.ParseForm(); err != nil {
		return err
	}

	val := ubind.ParseForm([]byte(r.raw.Form.Encode()))
	return ubind.Bind(val, obj)
}

// BindQuery 绑定查询参数
func (r *Request) BindQuery(obj ubind.Binder) error {
	val := ubind.ParseForm([]byte(r.raw.URL.RawQuery))
	return ubind.Bind(val, obj)
}

// Method 获取请求方法
func (r *Request) Method() string {
	return r.raw.Method
}

// Path 获取请求路径
func (r *Request) Path() string {
	return r.raw.URL.Path
}

// URL 获取完整 URL
func (r *Request) URL() *url.URL {
	return r.raw.URL
}

// Session 获取 Session
func (r *Request) Session() (unet.Session, error) {
	// 获取 Session 管理器
	sessionMgr := r.Server().SessionManager()
	if sessionMgr != nil {
		session, _ := sessionMgr.Start(r)
		return session, nil
	}
	return nil, errors.New("session manager not found")
}

// Raw 获取原始 *http.Request
func (r *Request) Raw() *http.Request {
	return r.raw
}

// getMaxBodySize 获取最大请求体大小
func (r *Request) getMaxBodySize() int64 {
	// 从服务器配置获取
	if r.server != nil && r.server.config != nil && r.server.config.MaxBodyBytes > 0 {
		return r.server.config.MaxBodyBytes
	}
	// 默认 10MB
	return 10 << 20
}
