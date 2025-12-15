package uhttp

import (
	"net/http"
)

// SetCookie 设置 Cookie (便捷方法)
func (r *Response) SetCookieValue(name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	r.SetCookie(cookie)
}

// GetCookie 获取 Cookie 值
func (r *Request) GetCookie(name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// DeleteCookie 删除 Cookie
func (r *Response) DeleteCookie(name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	r.SetCookie(cookie)
}

// SetSecureCookie 设置安全 Cookie
func (r *Response) SetSecureCookie(name, value string, maxAge int, domain string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	r.SetCookie(cookie)
}

func (r *Response) SetSessionCookie(name string, id string, path string, domain string, age int, secure bool, only bool, site http.SameSite) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    id,
		Path:     path,
		Domain:   domain,
		MaxAge:   age,
		Secure:   secure,
		HttpOnly: only,
		SameSite: site,
	}
	r.SetCookie(cookie)
}
