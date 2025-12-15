package uhttp

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/whosafe/uf/uerror"
)

// FormFile 获取上传的文件
func (r *Request) FormFile(name string) (*multipart.FileHeader, error) {
	_, fileHeader, err := r.raw.FormFile(name)
	if err != nil {
		return nil, uerror.Wrap(err, "获取上传文件失败")
	}
	return fileHeader, nil
}

// MultipartForm 获取 multipart form
func (r *Request) MultipartForm() (*multipart.Form, error) {
	if err := r.raw.ParseMultipartForm(32 << 20); err != nil { // 32MB
		return nil, uerror.Wrap(err, "解析 multipart form 失败")
	}
	return r.raw.MultipartForm, nil
}

// SaveUploadedFile 保存上传的文件
func (r *Request) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return uerror.Wrap(err, "打开上传文件失败")
	}
	defer src.Close()

	// 创建目标目录
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return uerror.Wrap(err, "创建目录失败")
	}

	// 创建目标文件
	out, err := os.Create(dst)
	if err != nil {
		return uerror.Wrap(err, "创建文件失败")
	}
	defer out.Close()

	// 复制文件
	_, err = io.Copy(out, src)
	if err != nil {
		return uerror.Wrap(err, "复制文件失败")
	}
	return nil
}

// SaveUploadedFiles 保存多个上传的文件
func (r *Request) SaveUploadedFiles(formName, dstDir string) ([]string, error) {
	form, err := r.MultipartForm()
	if err != nil {
		return nil, uerror.Wrap(err, "获取 multipart form 失败")
	}

	files := form.File[formName]
	if len(files) == 0 {
		return nil, uerror.New("没有找到上传文件: " + formName)
	}

	savedPaths := make([]string, 0, len(files))
	for _, file := range files {
		// 【安全修复】强制使用 filepath.Base 清理文件名,防止路径遍历攻击
		dst := filepath.Join(dstDir, filepath.Base(file.Filename))
		if err := r.SaveUploadedFile(file, dst); err != nil {
			return savedPaths, err
		}
		savedPaths = append(savedPaths, dst)
	}

	return savedPaths, nil
}

// FileUploadConfig 文件上传配置
type FileUploadConfig struct {
	MaxSize      int64    // 最大文件大小 (字节)
	AllowedExts  []string // 允许的文件扩展名
	UploadDir    string   // 上传目录
	GenerateName bool     // 是否生成新文件名
}

// SaveUploadedFileWithConfig 使用配置保存上传的文件
func (r *Request) SaveUploadedFileWithConfig(file *multipart.FileHeader, cfg *FileUploadConfig) (string, error) {
	// 检查文件大小
	if cfg.MaxSize > 0 && file.Size > cfg.MaxSize {
		return "", uerror.New(fmt.Sprintf("文件大小 %d 超过最大限制 %d", file.Size, cfg.MaxSize))
	}

	// 检查文件扩展名
	if len(cfg.AllowedExts) > 0 {
		ext := filepath.Ext(file.Filename)
		allowed := false
		for _, allowedExt := range cfg.AllowedExts {
			if ext == allowedExt || ext == "."+allowedExt {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", uerror.New("文件扩展名不允许: " + ext)
		}
	}

	// 生成文件名
	// 【安全修复】强制使用 filepath.Base 清理文件名,防止路径遍历攻击
	filename := filepath.Base(file.Filename)
	if cfg.GenerateName {
		ext := filepath.Ext(file.Filename)
		filename = fmt.Sprintf("%d%s", time.Now().Unix(), ext) // 简单示例,实际应使用更好的方法
	}

	// 保存文件
	dst := filepath.Join(cfg.UploadDir, filename)
	if err := r.SaveUploadedFile(file, dst); err != nil {
		return "", uerror.Wrap(err, "保存文件失败")
	}

	return dst, nil
}

// GetFileSize 获取上传文件的大小
func GetFileSize(file *multipart.FileHeader) int64 {
	return file.Size
}

// GetFileExt 获取上传文件的扩展名
func GetFileExt(file *multipart.FileHeader) string {
	return filepath.Ext(file.Filename)
}
