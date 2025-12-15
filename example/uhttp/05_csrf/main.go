package main

import (
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uprotocol/uhttp"
	"github.com/whosafe/uf/uprotocol/unet"
)

func main() {
	// 创建服务器
	server := uhttp.New()

	// 【安全修复】启用 CSRF 保护中间件
	// 方式1: 使用默认配置
	server.Use(uhttp.MiddlewareCSRF())

	// 方式2: 使用自定义配置
	// csrfConfig := uhttp.CSRFConfig{
	// 	TokenLength:   32,
	// 	CookieName:    "csrf_token",
	// 	HeaderName:    "X-CSRF-Token",
	// 	FormFieldName: "csrf_token",
	// 	CookieMaxAge:  3600,
	// 	SkipMethods:   []string{"GET", "HEAD", "OPTIONS"},
	// }
	// server.Use(uhttp.MiddlewareCSRFWithConfig(csrfConfig))

	// GET 请求 - 自动生成 CSRF Token
	server.GET("/form", func(ctx *ucontext.Context, req unet.Request) error {
		httpReq := req.(*uhttp.Request)
		httpResp := req.Response().(*uhttp.Response)

		// 获取 CSRF Token
		csrfToken := httpReq.GetCSRFToken()

		// 返回包含 CSRF Token 的表单
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>CSRF Protection Demo</title></head>
		<body>
			<h1>CSRF Protection Demo</h1>
			
			<!-- 方式1: 使用隐藏字段 -->
			<h2>Form with Hidden Field</h2>
			<form action="/submit" method="POST">
				<input type="hidden" name="csrf_token" value="` + csrfToken + `">
				<input type="text" name="username" placeholder="Username">
				<button type="submit">Submit</button>
			</form>

			<!-- 方式2: 使用 JavaScript 和请求头 -->
			<h2>AJAX with Header</h2>
			<button onclick="submitWithAjax()">Submit via AJAX</button>

			<script>
			function submitWithAjax() {
				fetch('/submit', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						'X-CSRF-Token': '` + csrfToken + `'
					},
					body: JSON.stringify({ username: 'test' })
				})
				.then(response => response.json())
				.then(data => alert('Success: ' + JSON.stringify(data)))
				.catch(error => alert('Error: ' + error));
			}
			</script>
		</body>
		</html>
		`

		return httpResp.HTML(200, html)
	})

	// POST 请求 - 需要验证 CSRF Token
	server.POST("/submit", func(ctx *ucontext.Context, req unet.Request) error {
		httpResp := req.Response().(*uhttp.Response)

		// 如果到达这里,说明 CSRF 验证已通过
		return httpResp.JSON(200, map[string]any{
			"success": true,
			"message": "CSRF validation passed!",
		})
	})

	// 启动服务器
	server.Start(":8080")
}
