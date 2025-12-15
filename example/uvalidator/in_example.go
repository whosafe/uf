package main

import (
	"fmt"

	"github.com/whosafe/uf/uvalidator"
	"github.com/whosafe/uf/uvalidator/rule"
)

// UserRole 用户角色
type UserRole struct {
	Role string
}

// ProductStatus 产品状态
type ProductStatus struct {
	Status int
}

// PriceLevel 价格等级
type PriceLevel struct {
	Level float64
}

func InExample() {
	// 示例 1: 字符串 In 验证
	fmt.Println("=== 示例 1: 字符串 In 验证 ===")
	roleRule := rule.NewIn("admin", "user", "guest")

	testRoles := []string{"admin", "user", "moderator", ""}
	for _, role := range testRoles {
		valid := roleRule.Validate(role)
		if !valid {
			msg := roleRule.GetMessage("Role", uvalidator.LanguageZH)
			fmt.Printf("角色 '%s': %s\n", role, msg)
		} else {
			fmt.Printf("角色 '%s': 验证通过\n", role)
		}
	}

	// 示例 2: 整数 In 验证
	fmt.Println("\n=== 示例 2: 整数 In 验证 ===")
	statusRule := rule.NewInInt(0, 1, 2, 3) // 0=待审核, 1=已发布, 2=已下架, 3=已删除

	testStatuses := []int{0, 1, 2, 5}
	statusNames := map[int]string{0: "待审核", 1: "已发布", 2: "已下架", 3: "已删除", 5: "未知状态"}
	for _, status := range testStatuses {
		valid := statusRule.Validate(status)
		if !valid {
			msg := statusRule.GetMessage("Status", uvalidator.LanguageZH)
			fmt.Printf("状态 %d (%s): %s\n", status, statusNames[status], msg)
		} else {
			fmt.Printf("状态 %d (%s): 验证通过\n", status, statusNames[status])
		}
	}

	// 示例 3: 浮点数 In 验证
	fmt.Println("\n=== 示例 3: 浮点数 In 验证 ===")
	priceRule := rule.NewInFloat(9.9, 19.9, 29.9, 49.9, 99.9)

	testPrices := []float64{9.9, 19.9, 39.9, 99.9}
	for _, price := range testPrices {
		valid := priceRule.Validate(price)
		if !valid {
			msg := priceRule.GetMessage("Price", uvalidator.LanguageZH)
			fmt.Printf("价格 %.2f: %s\n", price, msg)
		} else {
			fmt.Printf("价格 %.2f: 验证通过\n", price)
		}
	}

	// 示例 4: 在结构体验证中使用
	fmt.Println("\n=== 示例 4: 在结构体验证中使用 ===")
	type CreateUserRequest struct {
		Username string
		Role     string
		Status   int
	}

	validateUser := func(req CreateUserRequest) error {
		var errs uvalidator.ValidationErrors

		// 验证 Username
		requiredRule := rule.NewRequired()
		if !requiredRule.Validate(req.Username) {
			errs = append(errs, uvalidator.NewFieldError(
				"Username",
				requiredRule.Name(),
				req.Username,
				requiredRule.GetMessage("Username", uvalidator.LanguageZH),
			))
		}

		// 验证 Role
		roleRule := rule.NewIn("admin", "user", "guest")
		if !roleRule.Validate(req.Role) {
			errs = append(errs, uvalidator.NewFieldError(
				"Role",
				roleRule.Name(),
				req.Role,
				roleRule.GetMessage("Role", uvalidator.LanguageZH),
			))
		}

		// 验证 Status
		statusRule := rule.NewInInt(0, 1, 2)
		if !statusRule.Validate(req.Status) {
			errs = append(errs, uvalidator.NewFieldError(
				"Status",
				statusRule.Name(),
				req.Status,
				statusRule.GetMessage("Status", uvalidator.LanguageZH),
			))
		}

		if errs.HasErrors() {
			return errs
		}
		return nil
	}

	// 测试有效请求
	validReq := CreateUserRequest{
		Username: "alice",
		Role:     "admin",
		Status:   1,
	}
	if err := validateUser(validReq); err != nil {
		fmt.Printf("验证失败: %v\n", err)
	} else {
		fmt.Printf("用户 %s (角色: %s, 状态: %d): 验证通过\n", validReq.Username, validReq.Role, validReq.Status)
	}

	// 测试无效请求
	invalidReq := CreateUserRequest{
		Username: "bob",
		Role:     "moderator", // 无效角色
		Status:   5,           // 无效状态
	}
	if err := validateUser(invalidReq); err != nil {
		fmt.Printf("用户 %s 验证失败:\n%v\n", invalidReq.Username, err)
	} else {
		fmt.Printf("用户 %s: 验证通过\n", invalidReq.Username)
	}
}
