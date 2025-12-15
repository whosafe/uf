package rule

import (
	"testing"

	"github.com/whosafe/uf/uvalidator"
)

// TestMaxLength 测试最大长度规则
func TestMaxLength(t *testing.T) {
	tests := []struct {
		name     string
		maxValue int
		value    any
		expected bool
	}{
		{
			name:     "有效值 - 刚好等于最大长度",
			maxValue: 5,
			value:    "hello",
			expected: true,
		},
		{
			name:     "有效值 - 小于最大长度",
			maxValue: 10,
			value:    "hello",
			expected: true,
		},
		{
			name:     "无效值 - 大于最大长度",
			maxValue: 3,
			value:    "hello",
			expected: false,
		},
		{
			name:     "边界值 - 空字符串",
			maxValue: 10,
			value:    "",
			expected: true,
		},
		{
			name:     "边界值 - 最大长度为0",
			maxValue: 0,
			value:    "",
			expected: true,
		},
		{
			name:     "边界值 - 最大长度为0但有内容",
			maxValue: 0,
			value:    "a",
			expected: false,
		},
		{
			name:     "无效类型 - 整数",
			maxValue: 5,
			value:    12345,
			expected: false,
		},
		{
			name:     "无效类型 - nil",
			maxValue: 5,
			value:    nil,
			expected: false,
		},
		{
			name:     "中文字符串",
			maxValue: 15, // "你好世界" 在 UTF-8 中是 12 字节
			value:    "你好世界",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewMaxLength(tt.maxValue)
			result := rule.Validate(tt.value)
			if result != tt.expected {
				t.Errorf("MaxLength.Validate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestMaxLengthGetMessage 测试错误消息生成
func TestMaxLengthGetMessage(t *testing.T) {
	rule := NewMaxLength(20)

	// 测试默认语言
	msg := rule.GetMessage("Username")
	if msg == "" {
		t.Error("MaxLength.GetMessage() returned empty string")
	}

	// 测试英文
	msgEN := rule.GetMessage("Username", uvalidator.LanguageEN)
	if msgEN == "" {
		t.Error("MaxLength.GetMessage() with EN returned empty string")
	}

	// 测试中文
	msgZH := rule.GetMessage("Username", uvalidator.LanguageZH)
	if msgZH == "" {
		t.Error("MaxLength.GetMessage() with ZH returned empty string")
	}

	// 验证消息不同
	if msgEN == msgZH {
		t.Error("EN and ZH messages should be different")
	}
}

// TestMaxLengthName 测试规则名称
func TestMaxLengthName(t *testing.T) {
	rule := NewMaxLength(20)
	if rule.Name() != "max_length" {
		t.Errorf("MaxLength.Name() = %v, expected 'max_length'", rule.Name())
	}
}
