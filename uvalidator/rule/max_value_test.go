package rule

import (
	"testing"

	"github.com/whosafe/uf/uvalidator"
)

// TestMaxValue 测试最大值规则
func TestMaxValue(t *testing.T) {
	tests := []struct {
		name     string
		maxValue int
		value    any
		expected bool
	}{
		{
			name:     "有效值 - int刚好等于最大值",
			maxValue: 100,
			value:    100,
			expected: true,
		},
		{
			name:     "有效值 - int小于最大值",
			maxValue: 100,
			value:    50,
			expected: true,
		},
		{
			name:     "无效值 - int大于最大值",
			maxValue: 100,
			value:    150,
			expected: false,
		},
		{
			name:     "有效值 - int64刚好等于最大值",
			maxValue: 1000,
			value:    int64(1000),
			expected: true,
		},
		{
			name:     "有效值 - int64小于最大值",
			maxValue: 1000,
			value:    int64(500),
			expected: true,
		},
		{
			name:     "无效值 - int64大于最大值",
			maxValue: 1000,
			value:    int64(1500),
			expected: false,
		},
		{
			name:     "有效值 - float64刚好等于最大值",
			maxValue: 100,
			value:    100.0,
			expected: true,
		},
		{
			name:     "有效值 - float64小于最大值",
			maxValue: 100,
			value:    50.5,
			expected: true,
		},
		{
			name:     "无效值 - float64大于最大值",
			maxValue: 100,
			value:    150.5,
			expected: false,
		},
		{
			name:     "边界值 - 负数",
			maxValue: 0,
			value:    -5,
			expected: true,
		},
		{
			name:     "边界值 - 零",
			maxValue: 0,
			value:    0,
			expected: true,
		},
		{
			name:     "边界值 - 零但值为正",
			maxValue: 0,
			value:    1,
			expected: false,
		},
		{
			name:     "无效类型 - 字符串",
			maxValue: 100,
			value:    "100",
			expected: false,
		},
		{
			name:     "无效类型 - nil",
			maxValue: 100,
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewMaxValue(tt.maxValue)
			result := rule.Validate(tt.value)
			if result != tt.expected {
				t.Errorf("MaxValue.Validate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestMaxValueGetMessage 测试错误消息生成
func TestMaxValueGetMessage(t *testing.T) {
	rule := NewMaxValue(100)

	// 测试默认语言
	msg := rule.GetMessage("Score")
	if msg == "" {
		t.Error("MaxValue.GetMessage() returned empty string")
	}

	// 测试英文
	msgEN := rule.GetMessage("Score", uvalidator.LanguageEN)
	if msgEN == "" {
		t.Error("MaxValue.GetMessage() with EN returned empty string")
	}

	// 测试中文
	msgZH := rule.GetMessage("Score", uvalidator.LanguageZH)
	if msgZH == "" {
		t.Error("MaxValue.GetMessage() with ZH returned empty string")
	}

	// 验证消息不同
	if msgEN == msgZH {
		t.Error("EN and ZH messages should be different")
	}
}

// TestMaxValueName 测试规则名称
func TestMaxValueName(t *testing.T) {
	rule := NewMaxValue(100)
	if rule.Name() != "max" {
		t.Errorf("MaxValue.Name() = %v, expected 'max'", rule.Name())
	}
}
