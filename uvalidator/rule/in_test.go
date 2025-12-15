package rule

import (
	"testing"

	"github.com/whosafe/uf/uvalidator"
)

// TestIn 测试字符串 In 规则
func TestIn(t *testing.T) {
	tests := []struct {
		name     string
		allowed  []string
		value    any
		expected bool
	}{
		{
			name:     "有效值",
			allowed:  []string{"apple", "banana", "orange"},
			value:    "apple",
			expected: true,
		},
		{
			name:     "有效值 - banana",
			allowed:  []string{"apple", "banana", "orange"},
			value:    "banana",
			expected: true,
		},
		{
			name:     "无效值",
			allowed:  []string{"apple", "banana", "orange"},
			value:    "grape",
			expected: false,
		},
		{
			name:     "空字符串",
			allowed:  []string{"apple", "banana", "orange"},
			value:    "",
			expected: true, // 空字符串通过验证
		},
		{
			name:     "非字符串类型",
			allowed:  []string{"apple", "banana", "orange"},
			value:    123,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewIn(tt.allowed...)
			result := rule.Validate(tt.value)
			if result != tt.expected {
				t.Errorf("In.Validate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestInInt 测试整数 In 规则
func TestInInt(t *testing.T) {
	tests := []struct {
		name     string
		allowed  []int
		value    any
		expected bool
	}{
		{
			name:     "有效值",
			allowed:  []int{1, 2, 3, 4, 5},
			value:    3,
			expected: true,
		},
		{
			name:     "有效值 - 边界值",
			allowed:  []int{1, 2, 3, 4, 5},
			value:    1,
			expected: true,
		},
		{
			name:     "无效值",
			allowed:  []int{1, 2, 3, 4, 5},
			value:    6,
			expected: false,
		},
		{
			name:     "负数有效值",
			allowed:  []int{-1, 0, 1},
			value:    -1,
			expected: true,
		},
		{
			name:     "非整数类型",
			allowed:  []int{1, 2, 3},
			value:    "1",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewInInt(tt.allowed...)
			result := rule.Validate(tt.value)
			if result != tt.expected {
				t.Errorf("InInt.Validate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestInFloat 测试浮点数 In 规则
func TestInFloat(t *testing.T) {
	tests := []struct {
		name     string
		allowed  []float64
		value    any
		expected bool
	}{
		{
			name:     "有效值",
			allowed:  []float64{1.5, 2.5, 3.5},
			value:    2.5,
			expected: true,
		},
		{
			name:     "有效值 - 整数浮点数",
			allowed:  []float64{1.0, 2.0, 3.0},
			value:    1.0,
			expected: true,
		},
		{
			name:     "无效值",
			allowed:  []float64{1.5, 2.5, 3.5},
			value:    4.5,
			expected: false,
		},
		{
			name:     "负数有效值",
			allowed:  []float64{-1.5, 0.0, 1.5},
			value:    -1.5,
			expected: true,
		},
		{
			name:     "非浮点数类型",
			allowed:  []float64{1.5, 2.5, 3.5},
			value:    "2.5",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewInFloat(tt.allowed...)
			result := rule.Validate(tt.value)
			if result != tt.expected {
				t.Errorf("InFloat.Validate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestInGetMessage 测试错误消息生成
func TestInGetMessage(t *testing.T) {
	// 测试字符串 In 规则
	rule := NewIn("apple", "banana", "orange")
	msg := rule.GetMessage("Fruit")
	if msg == "" {
		t.Error("In.GetMessage() returned empty string")
	}

	// 测试整数 In 规则
	ruleInt := NewInInt(1, 2, 3)
	msgInt := ruleInt.GetMessage("Number")
	if msgInt == "" {
		t.Error("InInt.GetMessage() returned empty string")
	}

	// 测试浮点数 In 规则
	ruleFloat := NewInFloat(1.5, 2.5, 3.5)
	msgFloat := ruleFloat.GetMessage("Price")
	if msgFloat == "" {
		t.Error("InFloat.GetMessage() returned empty string")
	}
}

// TestInGetMessageWithLanguage 测试不同语言的错误消息
func TestInGetMessageWithLanguage(t *testing.T) {
	rule := NewIn("apple", "banana", "orange")

	// 英文消息
	msgEN := rule.GetMessage("Fruit", uvalidator.LanguageEN)
	if msgEN == "" {
		t.Error("In.GetMessage() with EN returned empty string")
	}

	// 中文消息
	msgZH := rule.GetMessage("Fruit", uvalidator.LanguageZH)
	if msgZH == "" {
		t.Error("In.GetMessage() with ZH returned empty string")
	}

	// 消息应该不同
	if msgEN == msgZH {
		t.Error("EN and ZH messages should be different")
	}
}

// TestInName 测试规则名称
func TestInName(t *testing.T) {
	rule := NewIn("apple", "banana")
	if rule.Name() != "in" {
		t.Errorf("In.Name() = %v, expected 'in'", rule.Name())
	}

	ruleInt := NewInInt(1, 2, 3)
	if ruleInt.Name() != "in" {
		t.Errorf("InInt.Name() = %v, expected 'in'", ruleInt.Name())
	}

	ruleFloat := NewInFloat(1.5, 2.5)
	if ruleFloat.Name() != "in" {
		t.Errorf("InFloat.Name() = %v, expected 'in'", ruleFloat.Name())
	}
}
