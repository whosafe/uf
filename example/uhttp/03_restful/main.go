package main

import (
	"github.com/whosafe/uf/ucontext"
	"github.com/whosafe/uf/uconv"
	"github.com/whosafe/uf/uprotocol/ubind"
	"github.com/whosafe/uf/uprotocol/uhttp"
	"github.com/whosafe/uf/uprotocol/unet"
)

// User 用户模型
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Bind 实现 ubind.Binder 接口
func (u *User) Bind(key string, value *ubind.Value) error {
	switch key {
	case "id":
		u.ID = value.Int()
	case "name":
		u.Name = value.Str()
	case "age":
		u.Age = value.Int()
	}
	return nil
}

// 模拟数据库
var users = []User{
	{ID: 1, Name: "Alice", Age: 25},
	{ID: 2, Name: "Bob", Age: 30},
}

func main() {
	server := uhttp.New()

	// RESTful API
	api := server.Group("/api")
	{
		// 用户相关
		users := api.Group("/users")
		{
			users.GET("", listUsers)         // GET /api/users
			users.GET("/:id", getUser)       // GET /api/users/1
			users.POST("", createUser)       // POST /api/users
			users.PUT("/:id", updateUser)    // PUT /api/users/1
			users.DELETE("/:id", deleteUser) // DELETE /api/users/1
		}
	}

	server.Start(":8080")
}

// 列出所有用户
func listUsers(ctx *ucontext.Context, req unet.Request) error {
	resp := req.Response().(*uhttp.Response)
	return resp.Success(users)
}

// 获取单个用户
func getUser(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)
	httpResp := req.Response().(*uhttp.Response)

	idStr := httpReq.Param("id")
	id := uconv.ToIntDef(idStr, 0)

	for _, user := range users {
		if user.ID == id {
			return httpResp.Success(user)
		}
	}

	return httpResp.NotFound("用户不存在")
}

// 创建用户
func createUser(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)
	httpResp := req.Response().(*uhttp.Response)

	var user User
	if err := httpReq.BindJSON(&user); err != nil {
		return httpResp.BadRequest("请求参数错误")
	}

	// 验证
	if user.Name == "" {
		return httpResp.BadRequest("用户名不能为空")
	}
	if user.Age <= 0 {
		return httpResp.BadRequest("年龄必须大于0")
	}

	user.ID = len(users) + 1
	users = append(users, user)

	return httpResp.SuccessWithMessage("创建成功", user)
}

// 更新用户
func updateUser(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)
	httpResp := req.Response().(*uhttp.Response)

	idStr := httpReq.Param("id")
	id := uconv.ToIntDef(idStr, 0)

	var updatedUser User
	if err := httpReq.BindJSON(&updatedUser); err != nil {
		return httpResp.BadRequest("请求参数错误")
	}

	for i, user := range users {
		if user.ID == id {
			updatedUser.ID = id
			users[i] = updatedUser
			return httpResp.SuccessWithMessage("更新成功", updatedUser)
		}
	}

	return httpResp.NotFound("用户不存在")
}

// 删除用户
func deleteUser(ctx *ucontext.Context, req unet.Request) error {
	httpReq := req.(*uhttp.Request)
	httpResp := req.Response().(*uhttp.Response)

	idStr := httpReq.Param("id")
	id := uconv.ToIntDef(idStr, 0)

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			return httpResp.SuccessWithMessage("删除成功", nil)
		}
	}

	return httpResp.NotFound("用户不存在")
}
