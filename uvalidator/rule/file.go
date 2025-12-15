package rule

import (
	"fmt"
	"strings"

	"github.com/whosafe/uf/uvalidator"

	"github.com/whosafe/uf/uvalidator/i18n"
)

// FileExtension 文件扩展名验证规则
type FileExtension struct {
	AllowedExtensions []string // 允许的扩展名列表
}

// Validate 执行验证
func (f *FileExtension) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	// 获取文件扩展名
	lastDot := -1
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] == '.' {
			lastDot = i
			break
		}
	}

	if lastDot == -1 {
		return false
	}

	ext := strings.ToLower(str[lastDot+1:])

	// 检查是否在允许列表中
	for _, allowed := range f.AllowedExtensions {
		if strings.ToLower(allowed) == ext {
			return true
		}
	}

	return false
}

// GetMessage 获取错误消息
func (f *FileExtension) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("file_extension", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", strings.Join(f.AllowedExtensions, ", "))
	return msg
}

// Name 规则名称
func (f *FileExtension) Name() string {
	return "file_extension"
}

// NewFileExtension 创建文件扩展名验证规则
func NewFileExtension(allowedExtensions ...string) *FileExtension {
	return &FileExtension{AllowedExtensions: allowedExtensions}
}

// MimeType MIME类型验证规则
type MimeType struct {
	AllowedMimeTypes []string // 允许的MIME类型列表
}

// Validate 执行验证
func (m *MimeType) Validate(value any) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	if str == "" {
		return true
	}

	// 检查是否在允许列表中
	for _, allowed := range m.AllowedMimeTypes {
		if str == allowed {
			return true
		}
	}

	return false
}

// GetMessage 获取错误消息
func (m *MimeType) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("mime_type", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{param}", strings.Join(m.AllowedMimeTypes, ", "))
	return msg
}

// Name 规则名称
func (m *MimeType) Name() string {
	return "mime_type"
}

// NewMimeType 创建MIME类型验证规则
func NewMimeType(allowedMimeTypes ...string) *MimeType {
	return &MimeType{AllowedMimeTypes: allowedMimeTypes}
}

// FileSize 文件大小验证规则
type FileSize struct {
	MinSize int64 // 最小字节),0表示不限制
	MaxSize int64 // 最大字节),0表示不限制
}

// Validate 执行验证
func (f *FileSize) Validate(value any) bool {
	var size int64

	switch v := value.(type) {
	case int:
		size = int64(v)
	case int64:
		size = v
	default:
		return false
	}

	if f.MinSize > 0 && size < f.MinSize {
		return false
	}

	if f.MaxSize > 0 && size > f.MaxSize {
		return false
	}

	return true
}

// GetMessage 获取错误消息
func (f *FileSize) GetMessage(field string, lang ...uvalidator.Language) string {
	template := i18n.GetMessage("file_size", lang...)
	msg := replaceAll(template, "{field}", field)
	msg = replaceAll(msg, "{min}", fmt.Sprintf("%d", f.MinSize))
	msg = replaceAll(msg, "{max}", fmt.Sprintf("%d", f.MaxSize))
	return msg
}

// Name 规则名称
func (f *FileSize) Name() string {
	return "file_size"
}

// NewFileSize 创建文件大小验证规则
func NewFileSize(minSize, maxSize int64) *FileSize {
	return &FileSize{MinSize: minSize, MaxSize: maxSize}
}
