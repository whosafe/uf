package main

import (
	"github.com/whosafe/uf/uconfig"
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/uhttp"
	"github.com/whosafe/uf/uprotocol/unet"
)

func main() {
	// 加载配置文件
	uconfig.Load("example/uhttp/04_advanced/config.yaml")

	// 创建服务器
	server := uhttp.New()

	// 自动应用核心中间件 (Trace, Logger, Recovery)
	// 根据配置文件中的 middleware 配置决定是否启用
	uhttp.ApplyDefaultMiddlewares(server)

	// 路由
	server.GET("/", func(ctx *ucontext.Context, req unet.Request) error {
		httpResp := req.Response().(*uhttp.Response)
		return httpResp.Success(map[string]string{
			"message": "Server with config file",
		})
	})

	// Session 示例
	server.GET("/login", loginHandler)
	server.GET("/profile", profileHandler)

	// 文件上传示例
	server.POST("/upload", uploadHandler)

	// 启动服务器 (地址来自配置文件)
	if err := server.Start(""); err != nil {
		panic(err)
	}
}

func loginHandler(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)
	httpResp := req.Response().(*uhttp.Response)

	// 获取 Session 管理器
	sessionMgr := httpReq.Server().SessionManager()
	if sessionMgr != nil {
		session, _ := sessionMgr.Start(httpReq)
		session.Set("user_id", 123)
		session.Set("username", "alice")
		session.Save()
		return httpResp.SuccessWithMessage("登录成功", map[string]any{
			"user_id":  123,
			"username": "alice",
		})
	}

	// Session 未启用
	return httpResp.SuccessWithMessage("登录成功(Session未启用)", map[string]any{
		"user_id":  123,
		"username": "alice",
	})
}

func profileHandler(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)
	httpResp := req.Response().(*uhttp.Response)

	sessionMgr := httpReq.Server().SessionManager()
	if sessionMgr != nil {
		session, _ := sessionMgr.Start(httpReq)
		userID, _ := session.Get("user_id")
		username, _ := session.Get("username")

		return httpResp.Success(map[string]any{
			"user_id":  userID,
			"username": username,
		})
	}

	return httpResp.Unauthorized("未登录")
}

func uploadHandler(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)
	httpResp := req.Response().(*uhttp.Response)

	// 获取上传文件
	file, err := httpReq.FormFile("file")
	if err != nil {
		return httpResp.BadRequest("未上传文件")
	}

	// 保存文件
	path, err := httpReq.SaveUploadedFileWithConfig(file, &uhttp.FileUploadConfig{
		MaxSize:     10 << 20, // 10MB
		AllowedExts: []string{".jpg", ".png", ".gif"},
		UploadDir:   "./uploads",
	})

	if err != nil {
		return httpResp.BadRequest(err.Error())
	}

	return httpResp.SuccessWithMessage("上传成功", map[string]string{
		"path": path,
	})
}
