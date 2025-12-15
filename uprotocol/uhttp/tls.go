package uhttp

import (
	"crypto/tls"

	"github.com/whosafe/uf/uerror"
)

// TLSConfig TLS 配置
type TLSConfig struct {
	CertFile   string // 证书文件路径
	KeyFile    string // 密钥文件路径
	MinVersion uint16 // 最小 TLS 版本
	MaxVersion uint16 // 最大 TLS 版本
}

// DefaultTLSConfig 默认 TLS 配置
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}
}

// StartTLS 启动 HTTPS 服务器
func (s *Server) StartTLS(addr, certFile, keyFile string) error {
	s.config.Address = addr
	s.config.Protocol = "https"

	// 配置 TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}

	s.httpServer.Addr = addr
	s.httpServer.TLSConfig = tlsConfig

	// 启动 HTTPS 服务器
	if err := s.httpServer.ListenAndServeTLS(certFile, keyFile); err != nil {
		s.errorLogger.Error("HTTPS 服务器启动失败", "error", err)
		return uerror.Wrap(err, "HTTPS 服务器启动失败")
	}

	return nil
}

// StartTLSWithConfig 使用配置启动 HTTPS 服务器
func (s *Server) StartTLSWithConfig(addr string, tlsCfg *TLSConfig) error {
	if tlsCfg.CertFile == "" || tlsCfg.KeyFile == "" {
		return uerror.New("证书文件和密钥文件不能为空")
	}

	s.config.Address = addr
	s.config.Protocol = "https"

	// 配置 TLS
	tlsConfig := &tls.Config{
		MinVersion: tlsCfg.MinVersion,
		MaxVersion: tlsCfg.MaxVersion,
	}

	s.httpServer.Addr = addr
	s.httpServer.TLSConfig = tlsConfig

	// 启动 HTTPS 服务器
	if err := s.httpServer.ListenAndServeTLS(tlsCfg.CertFile, tlsCfg.KeyFile); err != nil {
		s.errorLogger.Error("HTTPS 服务器启动失败", "error", err)
		return uerror.Wrap(err, "HTTPS 服务器启动失败")
	}

	return nil
}
