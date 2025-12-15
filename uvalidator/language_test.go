package uvalidator_test

import (
	"testing"

	"github.com/whosafe/uf/uvalidator"
	"github.com/whosafe/uf/uvalidator/rule"
)

// TestGlobalLanguage 测试全局语言设置
func TestGlobalLanguage(t *testing.T) {
	// 设置全局语言为中文
	uvalidator.SetLanguage(uvalidator.LanguageZH)

	requiredRule := rule.NewRequired()

	// 不传语言参数,使用全局设置
	msg := requiredRule.GetMessage("Username")
	expected := "Username不能为空"
	if msg != expected {
		t.Errorf("Expected '%s', got '%s'", expected, msg)
	}

	// 切换全局语言为英文
	uvalidator.SetLanguage(uvalidator.LanguageEN)
	msg = requiredRule.GetMessage("Username")
	expected = "Username is required"
	if msg != expected {
		t.Errorf("Expected '%s', got '%s'", expected, msg)
	}
}

// TestRequestLevelLanguage 测试请求级语言选择
func TestRequestLevelLanguage(t *testing.T) {
	// 设置全局语言为英文
	uvalidator.SetLanguage(uvalidator.LanguageEN)

	requiredRule := rule.NewRequired()

	// 测试1: 显式传入中文
	msg := requiredRule.GetMessage("Username", uvalidator.LanguageZH)
	expected := "Username不能为空"
	if msg != expected {
		t.Errorf("Expected '%s', got '%s'", expected, msg)
	}

	// 测试2: 显式传入英文
	msg = requiredRule.GetMessage("Email", uvalidator.LanguageEN)
	expected = "Email is required"
	if msg != expected {
		t.Errorf("Expected '%s', got '%s'", expected, msg)
	}

	// 测试3: 不传语言参数,应使用全局设置(英文)
	msg = requiredRule.GetMessage("Password")
	expected = "Password is required"
	if msg != expected {
		t.Errorf("Expected '%s', got '%s'", expected, msg)
	}
}

// TestParseAcceptLanguage 测试 Accept-Language 解析
func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		acceptLang string
		expected   uvalidator.Language
	}{
		{"zh-CN,zh;q=0.9,en;q=0.8", uvalidator.LanguageZH},
		{"en-US,en;q=0.9", uvalidator.LanguageEN},
		{"zh", uvalidator.LanguageZH},
		{"en", uvalidator.LanguageEN},
		{"", uvalidator.LanguageEN},               // 空字符串使用全局设置
		{"fr-FR,fr;q=0.9", uvalidator.LanguageEN}, // 不支持的语言回退到全局设置
	}

	// 设置全局语言为英文
	uvalidator.SetLanguage(uvalidator.LanguageEN)

	for _, tt := range tests {
		result := uvalidator.ParseAcceptLanguage(tt.acceptLang)
		if result != tt.expected {
			t.Errorf("ParseAcceptLanguage(%q) = %v, expected %v", tt.acceptLang, result, tt.expected)
		}
	}
}

// TestMultipleRulesWithLanguage 测试多个规则的语言选择
func TestMultipleRulesWithLanguage(t *testing.T) {
	// 测试不同规则使用不同语言
	emailRule := rule.NewEmail()
	minRule := rule.NewMinLength(3)
	phoneRule := rule.NewPhone()

	// 使用中文
	lang := uvalidator.LanguageZH

	emailMsg := emailRule.GetMessage("Email", lang)
	if emailMsg != "Email必须是有效的邮箱地址" {
		t.Errorf("Email rule message incorrect: %s", emailMsg)
	}

	minMsg := minRule.GetMessage("Username", lang)
	expected := "Username长度不能少于3个字符"
	if minMsg != expected {
		t.Errorf("Min rule message incorrect: expected '%s', got '%s'", expected, minMsg)
	}

	phoneMsg := phoneRule.GetMessage("Phone", lang)
	if phoneMsg != "Phone必须是有效的手机号" {
		t.Errorf("Phone rule message incorrect: %s", phoneMsg)
	}
}

// TestConcurrentLanguageSelection 测试并发场景下的语言选择
func TestConcurrentLanguageSelection(t *testing.T) {
	requiredRule := rule.NewRequired()

	// 模拟并发请求,每个请求使用不同的语言
	done := make(chan bool, 2)

	// 协程1: 使用中文
	go func() {
		for i := 0; i < 100; i++ {
			msg := requiredRule.GetMessage("Field", uvalidator.LanguageZH)
			if msg != "Field不能为空" {
				t.Errorf("Concurrent test failed: expected Chinese message, got %s", msg)
			}
		}
		done <- true
	}()

	// 协程2: 使用英文
	go func() {
		for i := 0; i < 100; i++ {
			msg := requiredRule.GetMessage("Field", uvalidator.LanguageEN)
			if msg != "Field is required" {
				t.Errorf("Concurrent test failed: expected English message, got %s", msg)
			}
		}
		done <- true
	}()

	// 等待两个协程完成
	<-done
	<-done
}

// TestNewRulesWithLanguage 测试新增规则的语言选择
func TestNewRulesWithLanguage(t *testing.T) {
	tests := []struct {
		name     string
		rule     uvalidator.Rule
		field    string
		lang     uvalidator.Language
		expected string
	}{
		{
			name:     "UUID - Chinese",
			rule:     rule.NewUUID(),
			field:    "ID",
			lang:     uvalidator.LanguageZH,
			expected: "ID必须是有效的UUID",
		},
		{
			name:     "UUID - English",
			rule:     rule.NewUUID(),
			field:    "ID",
			lang:     uvalidator.LanguageEN,
			expected: "ID must be a valid UUID",
		},
		{
			name:     "Between - Chinese",
			rule:     rule.NewBetween(10, 100),
			field:    "Age",
			lang:     uvalidator.LanguageZH,
			expected: "Age必须在10和100之间",
		},
		{
			name:     "Between - English",
			rule:     rule.NewBetween(10, 100),
			field:    "Age",
			lang:     uvalidator.LanguageEN,
			expected: "Age must be between 10 and 100",
		},
		{
			name:     "IDCard - Chinese",
			rule:     rule.NewIDCard(),
			field:    "IDCard",
			lang:     uvalidator.LanguageZH,
			expected: "IDCard必须是有效的身份证号",
		},
		{
			name:     "StrongPassword - Chinese",
			rule:     rule.NewStrongPassword(),
			field:    "Password",
			lang:     uvalidator.LanguageZH,
			expected: "Password必须是强密码(至少8位,包含大写、小写、数字和特殊字符)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.rule.GetMessage(tt.field, tt.lang)
			if msg != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, msg)
			}
		})
	}
}
