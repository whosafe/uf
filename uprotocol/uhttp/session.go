package uhttp

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/whosafe/uf/uprotocol/unet"
)

// SessionManager 会话管理器
type SessionManager struct {
	store       unet.SessionStore
	cookieName  string
	maxAge      int
	maxLifetime time.Duration // 会话最大存活时间
}

// NewSessionManager 创建会话管理器
func NewSessionManager(store unet.SessionStore, cookieName string, maxAge int) *SessionManager {
	return NewSessionManagerWithLifetime(store, cookieName, maxAge, 30*time.Minute)
}

// NewSessionManagerWithLifetime 创建会话管理器并指定存活周期
func NewSessionManagerWithLifetime(store unet.SessionStore, cookieName string, maxAge int, maxLifetime time.Duration) *SessionManager {
	return &SessionManager{
		store:       store,
		cookieName:  cookieName,
		maxAge:      maxAge,
		maxLifetime: maxLifetime,
	}
}

// Start 启动会话
func (m *SessionManager) Start(req *Request) (unet.Session, error) {
	// 尝试从 Cookie 获取会话 ID
	sessionID, err := req.GetCookie(m.cookieName)
	if err == nil && sessionID != "" {
		// 尝试获取现有会话
		session, err := m.store.Get(sessionID)
		if err == nil && session != nil {
			return session, nil
		}
	}

	// 创建新会话
	session, err := m.store.Create()
	if err != nil {
		return nil, err
	}

	// 设置 Cookie (使用配置)
	m.setSessionCookie(req, req.Response(), session.ID())

	return session, nil
}

// setSessionCookie 设置 Session Cookie
func (m *SessionManager) setSessionCookie(req *Request, resp unet.Response, sessionID string) {
	// 获取 Cookie 配置
	cookieCfg := req.Server().config.Cookie

	if cookieCfg != nil {
		// 使用配置的 Cookie 设置
		resp.SetSessionCookie(m.cookieName, sessionID, cookieCfg.Path, cookieCfg.Domain, m.maxAge, cookieCfg.Secure, cookieCfg.HttpOnly, parseSameSite(cookieCfg.SameSite))
	} else {
		// 【安全修复】使用安全的默认设置
		// HttpOnly: 防止 XSS 攻击
		// SameSite=Lax: 防止 CSRF 攻击
		// Secure: 在 HTTPS 下传输(通过检测协议)
		isHTTPS := req.Protocol() == "https" || req.Header("X-Forwarded-Proto") == "https"
		resp.SetSessionCookie(m.cookieName, sessionID, "/", "", m.maxAge, isHTTPS, true, http.SameSiteLaxMode)
	}
}

// Destroy 销毁会话
func (m *SessionManager) Destroy(req *Request, resp *Response) error {
	sessionID, err := req.GetCookie(m.cookieName)
	if err != nil {
		return err
	}

	// 删除会话
	if err := m.store.Destroy(sessionID); err != nil {
		return err
	}

	// 删除 Cookie
	resp.DeleteCookie(m.cookieName)

	return nil
}

// generateSessionID 生成会话 ID
// 【安全修复】检查随机数生成错误,防止生成不安全的 Session ID
func generateSessionID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
