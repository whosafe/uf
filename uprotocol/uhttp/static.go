package uhttp

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/unet"
)

// StaticConfig 静态文件配置
type StaticConfig struct {
	Root   string   // 静态文件根目录
	Prefix string   // URL 前缀
	Index  []string // 索引文件列表
	Browse bool     // 是否允许目录浏览
}

// Static 注册静态文件服务
func (s *Server) Static(prefix, root string) {
	s.StaticWithConfig(&StaticConfig{
		Root:   root,
		Prefix: prefix,
		Index:  []string{"index.html", "index.htm"},
		Browse: false,
	})
}

// StaticWithConfig 使用配置注册静态文件服务
func (s *Server) StaticWithConfig(cfg *StaticConfig) {
	// 确保前缀以 / 开头
	if !strings.HasPrefix(cfg.Prefix, "/") {
		cfg.Prefix = "/" + cfg.Prefix
	}

	// 创建文件服务器
	fileServer := http.FileServer(http.Dir(cfg.Root))

	// 注册路由
	pattern := cfg.Prefix + "/*filepath"
	s.GET(pattern, func(ctx *ucontext.Context, req unet.Request) error {
		httpReq := req.(*Request)
		httpResp := req.Response().(*Response)

		// 获取文件路径
		path := httpReq.Param("filepath")
		if path == "" {
			path = "/"
		}

		// 【安全修复】清理路径,防止路径遍历攻击
		path = filepath.Clean(path)

		// 【安全修复】拒绝包含 .. 的路径
		if strings.Contains(path, "..") {
			return httpResp.Forbidden("非法路径")
		}

		// 构建完整路径
		fullPath := filepath.Join(cfg.Root, path)

		// 【安全修复】验证最终路径是否在根目录内
		absRoot, err := filepath.Abs(cfg.Root)
		if err != nil {
			return httpResp.InternalError("服务器配置错误")
		}
		absPath, err := filepath.Abs(fullPath)
		if err != nil {
			return httpResp.InternalError("路径解析失败")
		}
		// 确保访问路径在根目录内
		if !strings.HasPrefix(absPath, absRoot) {
			return httpResp.Forbidden("非法路径")
		}

		// 检查文件是否存在
		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) {
				return httpResp.NotFound("文件不存在")
			}
			return err
		}

		// 如果是目录
		if info.IsDir() {
			// 尝试查找索引文件
			for _, index := range cfg.Index {
				indexPath := filepath.Join(fullPath, index)
				if _, err := os.Stat(indexPath); err == nil {
					path = filepath.Join(path, index)
					break
				}
			}

			// 如果不允许目录浏览且没有找到索引文件
			if !cfg.Browse {
				return httpResp.Forbidden("不允许目录浏览")
			}
		}

		// 设置正确的 URL 路径
		httpReq.raw.URL.Path = path

		// 使用文件服务器处理
		fileServer.ServeHTTP(httpReq.writer, httpReq.raw)

		return nil
	})
}

// File 发送文件
func (s *Server) File(path, filepath string) {
	s.GET(path, func(ctx *ucontext.Context, req unet.Request) error {
		httpReq := req.(*Request)
		http.ServeFile(httpReq.writer, httpReq.raw, filepath)
		return nil
	})
}
