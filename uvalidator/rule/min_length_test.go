package rule

import (
	"testing"

	"github.com/whosafe/uf/uvalidator"
)

// TestMinLength 测试最小长度规则
func TestMinLength(t *testing.T) {
	tests := []struct {
		name     string
		minValue int
		value    any
		expected bool
	}{
		{
			name:     "有效值 - 刚好等于最小长度",
			minValue: 5,
			value:    "hello",
			expected: true,
		},
		{
			name:     "有效值 - 大于最小长度",
			minValue: 3,
			value:    "hello",
			expected: true,
		},
		{
			name:     "无效值 - 小于最小长度",
			minValue: 10,
			value:    "hello",
			expected: false,
		},
		{
			name:     "边界值 - 空字符串",
			minValue: 1,
			value:    "",
			expected: false,
		},
		{
			name:     "边界值 - 最小长度为0",
			minValue: 0,
			value:    "",
			expected: true,
		},
		{
			name:     "无效类型 - 整数",
			minValue: 5,
			value:    12345,
			expected: false,
		},
		{
			name:     "无效类型 - nil",
			minValue: 5,
			value:    nil,
			expected: false,
		},
		{
			name:     "中文字符串",
			minValue: 3,
			value:    "你好世界",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewMinLength(tt.minValue)
			result := rule.Validate(tt.value)
			if result != tt.expected {
				t.Errorf("MinLength.Validate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestMinLengthGetMessage 测试错误消息生成
func TestMinLengthGetMessage(t *testing.T) {
	rule := NewMinLength(5)

	// 测试默认语言
	msg := rule.GetMessage("Password")
	if msg == "" {
		t.Error("MinLength.GetMessage() returned empty string")
	}

	// 测试英文
	msgEN := rule.GetMessage("Password", uvalidator.LanguageEN)
	if msgEN == "" {
		t.Error("MinLength.GetMessage() with EN returned empty string")
	}

	// 测试中文
	msgZH := rule.GetMessage("Password", uvalidator.LanguageZH)
	if msgZH == "" {
		t.Error("MinLength.GetMessage() with ZH returned empty string")
	}

	// 验证消息不同
	if msgEN == msgZH {
		t.Error("EN and ZH messages should be different")
	}
}

// TestMinLengthName 测试规则名称
func TestMinLengthName(t *testing.T) {
	rule := NewMinLength(5)
	if rule.Name() != "min_length" {
		t.Errorf("MinLength.Name() = %v, expected 'min_length'", rule.Name())
	}
}
