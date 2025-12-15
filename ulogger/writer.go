package ulogger

// logWriter 实现 io.Writer 接口，用于将写入转换为日志输出
type logWriter struct {
	logger *Logger
}

// Write 实现 io.Writer 接口
func (lw *logWriter) Write(p []byte) (n int, err error) {
	lw.logger.Info(string(p))
	return len(p), nil
}
