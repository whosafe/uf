package uhttp

// SessionManager 获取 Session 管理器
func (s *Server) SessionManager() *SessionManager {
	return s.sessionManager
}

// SetSessionManager 设置 Session 管理器
func (s *Server) SetSessionManager(manager *SessionManager) {
	s.sessionManager = manager
}
