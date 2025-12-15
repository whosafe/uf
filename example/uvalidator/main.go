package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"iutime.com/utime/uf/uvalidator"
	"iutime.com/utime/uf/uvalidator/rule"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Phone    string `json:"phone"`
}

// Validate 验证请求数据
func (r *CreateUserRequest) Validate(lang ...uvalidator.Language) error {
	var errs uvalidator.ValidationErrors

	// Username 验证
	requiredRule := rule.NewRequired()
	if !requiredRule.Validate(r.Username) {
		errs = append(errs, uvalidator.NewFieldError(
			"Username",
			requiredRule.Name(),
			r.Username,
			requiredRule.GetMessage("Username", nil, lang...),
		))
	}

	minRule := rule.NewMin(3)
	if !minRule.Validate(r.Username) {
		errs = append(errs, uvalidator.NewFieldError(
			"Username",
			minRule.Name(),
			r.Username,
			minRule.GetMessage("Username", map[string]string{"type": "string"}, lang...),
		))
	}

	maxRule := rule.NewMax(20)
	if !maxRule.Validate(r.Username) {
		errs = append(errs, uvalidator.NewFieldError(
			"Username",
			maxRule.Name(),
			r.Username,
			maxRule.GetMessage("Username", map[string]string{"type": "string"}, lang...),
		))
	}

	// Email 验证
	if !requiredRule.Validate(r.Email) {
		errs = append(errs, uvalidator.NewFieldError(
			"Email",
			requiredRule.Name(),
			r.Email,
			requiredRule.GetMessage("Email", nil, lang...),
		))
	}

	emailRule := rule.NewEmail()
	if !emailRule.Validate(r.Email) {
		errs = append(errs, uvalidator.NewFieldError(
			"Email",
			emailRule.Name(),
			r.Email,
			emailRule.GetMessage("Email", nil, lang...),
		))
	}

	// Password 验证
	if !requiredRule.Validate(r.Password) {
		errs = append(errs, uvalidator.NewFieldError(
			"Password",
			requiredRule.Name(),
			r.Password,
			requiredRule.GetMessage("Password", nil, lang...),
		))
	}

	strongPwdRule := rule.NewStrongPassword(8)
	if !strongPwdRule.Validate(r.Password) {
		errs = append(errs, uvalidator.NewFieldError(
			"Password",
			strongPwdRule.Name(),
			r.Password,
			strongPwdRule.GetMessage("Password", nil, lang...),
		))
	}

	// Age 验证
	betweenRule := rule.NewBetween(18, 100)
	if !betweenRule.Validate(r.Age) {
		errs = append(errs, uvalidator.NewFieldError(
			"Age",
			betweenRule.Name(),
			r.Age,
			betweenRule.GetMessage("Age", nil, lang...),
		))
	}

	// Phone 验证 (可选)
	if r.Phone != "" {
		phoneRule := rule.NewPhone()
		if !phoneRule.Validate(r.Phone) {
			errs = append(errs, uvalidator.NewFieldError(
				"Phone",
				phoneRule.Name(),
				r.Phone,
				phoneRule.GetMessage("Phone", nil, lang...),
			))
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// CreateUserHandler HTTP 处理器示例
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	// 从请求头解析语言
	lang := uvalidator.ParseAcceptLanguage(r.Header.Get("Accept-Language"))

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 使用请求级语言验证
	if err := req.Validate(lang); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error": err.Error(),
		})
		return
	}

	// 处理成功
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "User created successfully",
		"user":    req,
	})
}

func main() {
	fmt.Println("=== uvalidator 使用示例 ===\n")

	// 示例 1: 基本验证 (使用全局语言)
	fmt.Println("--- 示例 1: 基本验证 (全局语言) ---")
	uvalidator.SetLanguage(uvalidator.LanguageZH)

	req1 := CreateUserRequest{
		Username: "ab",          // 太短
		Email:    "invalid",     // 无效邮箱
		Password: "weak",        // 弱密码
		Age:      15,            // 年龄不够
		Phone:    "12345678901", // 无效手机号
	}

	if err := req1.Validate(); err != nil {
		fmt.Printf("验证失败:\n%s\n\n", err)
	}

	// 示例 2: 有效数据
	fmt.Println("--- 示例 2: 有效数据 ---")
	req2 := CreateUserRequest{
		Username: "john_doe",
		Email:    "john@example.com",
		Password: "StrongP@ss123",
		Age:      25,
		Phone:    "13812345678",
	}

	if err := req2.Validate(); err != nil {
		fmt.Printf("验证失败: %s\n", err)
	} else {
		fmt.Println("验证通过! ✓")
	}
	fmt.Println()

	// 示例 3: 请求级语言选择
	fmt.Println("--- 示例 3: 请求级语言选择 ---")
	req3 := CreateUserRequest{
		Username: "ab",
		Email:    "invalid",
	}

	// 使用英文
	fmt.Println("使用英文:")
	if err := req3.Validate(uvalidator.LanguageEN); err != nil {
		fmt.Printf("%s\n\n", err)
	}

	// 使用中文
	fmt.Println("使用中文:")
	if err := req3.Validate(uvalidator.LanguageZH); err != nil {
		fmt.Printf("%s\n\n", err)
	}

	// 示例 4: 各种规则演示
	fmt.Println("--- 示例 4: 各种规则演示 ---")
	demonstrateRules()

	// 示例 5: HTTP 服务器
	fmt.Println("\n--- 示例 5: HTTP 服务器 ---")
	fmt.Println("启动 HTTP 服务器在 :8080")
	fmt.Println("测试命令:")
	fmt.Println("  curl -X POST http://localhost:8080/users \\")
	fmt.Println("    -H 'Content-Type: application/json' \\")
	fmt.Println("    -H 'Accept-Language: zh-CN' \\")
	fmt.Println("    -d '{\"username\":\"ab\",\"email\":\"invalid\"}'")
	fmt.Println()

	http.HandleFunc("/users", CreateUserHandler)
	// http.ListenAndServe(":8080", nil) // 取消注释以启动服务器
}

func demonstrateRules() {
	// 数字规则
	fmt.Println("数字规则:")
	betweenRule := rule.NewBetween(10, 100)
	fmt.Printf("  Between(50): %v\n", betweenRule.Validate(50))
	fmt.Printf("  Between(5): %v\n", betweenRule.Validate(5))

	positiveRule := rule.NewPositive()
	fmt.Printf("  Positive(10): %v\n", positiveRule.Validate(10))
	fmt.Printf("  Positive(-5): %v\n", positiveRule.Validate(-5))

	// 字符串规则
	fmt.Println("\n字符串规则:")
	uuidRule := rule.NewUUID()
	fmt.Printf("  UUID(valid): %v\n", uuidRule.Validate("550e8400-e29b-41d4-a716-446655440000"))
	fmt.Printf("  UUID(invalid): %v\n", uuidRule.Validate("not-a-uuid"))

	lowercaseRule := rule.NewLowercase()
	fmt.Printf("  Lowercase(hello): %v\n", lowercaseRule.Validate("hello"))
	fmt.Printf("  Lowercase(Hello): %v\n", lowercaseRule.Validate("Hello"))

	// 网络规则
	fmt.Println("\n网络规则:")
	ipRule := rule.NewIP()
	fmt.Printf("  IP(192.168.1.1): %v\n", ipRule.Validate("192.168.1.1"))
	fmt.Printf("  IP(invalid): %v\n", ipRule.Validate("999.999.999.999"))

	// 日期规则
	fmt.Println("\n日期规则:")
	dateRule := rule.NewDate()
	fmt.Printf("  Date(2024-01-01): %v\n", dateRule.Validate("2024-01-01"))
	fmt.Printf("  Date(2024/01/01): %v\n", dateRule.Validate("2024/01/01"))

	// 安全规则
	fmt.Println("\n安全规则:")
	strongPwdRule := rule.NewStrongPassword()
	fmt.Printf("  StrongPassword(Abc123!@#): %v\n", strongPwdRule.Validate("Abc123!@#"))
	fmt.Printf("  StrongPassword(weak): %v\n", strongPwdRule.Validate("weak"))

	noHTMLRule := rule.NewNoHTML()
	fmt.Printf("  NoHTML(plain text): %v\n", noHTMLRule.Validate("plain text"))
	fmt.Printf("  NoHTML(<script>): %v\n", noHTMLRule.Validate("<script>alert(1)</script>"))

	// 数组规则
	fmt.Println("\n数组规则:")
	uniqueRule := rule.NewUnique()
	fmt.Printf("  Unique([1,2,3]): %v\n", uniqueRule.Validate([]int{1, 2, 3}))
	fmt.Printf("  Unique([1,2,2]): %v\n", uniqueRule.Validate([]int{1, 2, 2}))
}
