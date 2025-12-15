package ulogger

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"iutime.com/utime/uf/uconfig"
	"iutime.com/utime/uf/uconv"
)

// Config 日志配置
type Config struct {
	Path                 string     // 日志文件路径。默认为空，表示关闭，仅输出到终端
	File                 string     // 日志文件格式。默认为"Y-m-d.log"
	Prefix               string     // 日志内容输出前缀。默认为空
	Level                slog.Level // 日志输出级别
	Format               string     // 日志格式："text", "json", "custom"。默认为"text"
	Formatter            Formatter  // 自定义格式化器（当 Format="custom" 时使用）
	UseStandardLogFormat bool       // 是否使用标准日志格式。默认true（仅对 text 格式有效）
	ShortFile            bool       // 日志文件是否只输出文件名。默认false
	Stdout               bool       // 日志是否同时输出到终端。默认true
	RotateSize           int        // 按照日志文件大小对文件进行滚动切分。默认为0，表示关闭滚动切分特性
	RotateExpire         int64      // 按照日志文件时间间隔对文件滚动切分。默认为0，表示关闭滚动切分特性
	RotateBackupLimit    int        // 按照切分的文件数量清理切分文件，当滚动切分特性开启时有效。默认为0，表示不备份，切分则删除
	RotateBackupExpire   int        // 按照切分的文件有效期清理切分文件，当滚动切分特性开启时有效。默认为0，表示不备份，切分则删除
	RotateBackupCompress uint16     // 滚动切分文件的压缩比（0-9）。默认为0，表示不压缩
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Path:                 "",
		File:                 "2006-01-02.log", // Go 时间格式
		Prefix:               "",
		Level:                slog.LevelInfo,
		Format:               "text",
		Formatter:            nil,
		UseStandardLogFormat: true,
		ShortFile:            false,
		Stdout:               true,
		RotateSize:           0,
		RotateExpire:         0,
		RotateBackupLimit:    0,
		RotateBackupExpire:   0,
		RotateBackupCompress: 0,
	}
}

// Validate 验证并规范化配置
func (c *Config) Validate() error {
	// 如果指定了路径，确保目录存在
	if c.Path != "" {
		if err := os.MkdirAll(c.Path, 0755); err != nil {
			return err
		}
	}

	// 如果文件名为空，使用默认值
	if c.File == "" {
		c.File = "2006-01-02.log"
	}

	// 压缩级别限制在 0-9
	if c.RotateBackupCompress > 9 {
		c.RotateBackupCompress = 9
	}

	return nil
}

// GetLogFilePath 获取当前日志文件的完整路径
func (c *Config) GetLogFilePath() string {
	if c.Path == "" {
		return ""
	}

	// 使用当前时间格式化文件名
	filename := time.Now().Format(c.File)
	return filepath.Join(c.Path, filename)
}

// IsRotateEnabled 检查是否启用了轮转功能
func (c *Config) IsRotateEnabled() bool {
	return c.RotateSize > 0 || c.RotateExpire > 0
}

// IsBackupEnabled 检查是否启用了备份功能
func (c *Config) IsBackupEnabled() bool {
	return c.RotateBackupLimit > 0 || c.RotateBackupExpire > 0
}

// UnmarshalYAML 实现 uconfig.Unmarshaler 接口
func (c *Config) UnmarshalYAML(key string, value *uconfig.Node) error {
	switch key {
	case "path":
		c.Path = value.String()
	case "file":
		c.File = value.String()
	case "prefix":
		c.Prefix = value.String()
	case "level":
		// 解析日志级别
		levelStr := value.String()
		switch levelStr {
		case "debug", "DEBUG":
			c.Level = slog.LevelDebug
		case "info", "INFO":
			c.Level = slog.LevelInfo
		case "warn", "WARN", "warning", "WARNING":
			c.Level = slog.LevelWarn
		case "error", "ERROR":
			c.Level = slog.LevelError
		default:
			// 尝试解析为数字
			if level, err := strconv.Atoi(levelStr); err == nil {
				c.Level = slog.Level(level)
			}
		}
	case "use_standard_log_format", "useStandardLogFormat":
		c.UseStandardLogFormat = uconv.ToBoolDef(value, true)
	case "short_file", "shortFile":
		c.ShortFile = uconv.ToBoolDef(value, false)
	case "stdout":
		c.Stdout = uconv.ToBoolDef(value, true)
	case "rotate_size", "rotateSize":
		c.RotateSize = uconv.ToIntDef(value, 0)
	case "rotate_expire", "rotateExpire":
		c.RotateExpire = int64(uconv.ToIntDef(value, 0))
	case "rotate_backup_limit", "rotateBackupLimit":
		c.RotateBackupLimit = uconv.ToIntDef(value, 0)
	case "rotate_backup_expire", "rotateBackupExpire":
		c.RotateBackupExpire = uconv.ToIntDef(value, 0)
	case "rotate_backup_compress", "rotateBackupCompress":
		c.RotateBackupCompress = uint16(uconv.ToIntDef(value, 0))
	case "format":
		c.Format = value.String()
	}
	return nil
}
