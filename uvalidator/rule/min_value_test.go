package rule

import (
	"testing"

	"github.com/whosafe/uf/uvalidator"
)

// TestMinValue 测试最小值规则
func TestMinValue(t *testing.T) {
	tests := []struct {
		name     string
		minValue int
		value    any
		expected bool
	}{
		{
			name:     "有效值 - int刚好等于最小值",
			minValue: 18,
			value:    18,
			expected: true,
		},
		{
			name:     "有效值 - int大于最小值",
			minValue: 18,
			value:    25,
			expected: true,
		},
		{
			name:     "无效值 - int小于最小值",
			minValue: 18,
			value:    15,
			expected: false,
		},
		{
			name:     "有效值 - int64刚好等于最小值",
			minValue: 100,
			value:    int64(100),
			expected: true,
		},
		{
			name:     "有效值 - int64大于最小值",
			minValue: 100,
			value:    int64(200),
			expected: true,
		},
		{
			name:     "无效值 - int64小于最小值",
			minValue: 100,
			value:    int64(50),
			expected: false,
		},
		{
			name:     "有效值 - float64刚好等于最小值",
			minValue: 10,
			value:    10.0,
			expected: true,
		},
		{
			name:     "有效值 - float64大于最小值",
			minValue: 10,
			value:    15.5,
			expected: true,
		},
		{
			name:     "无效值 - float64小于最小值",
			minValue: 10,
			value:    5.5,
			expected: false,
		},
		{
			name:     "边界值 - 负数",
			minValue: -10,
			value:    -5,
			expected: true,
		},
		{
			name:     "边界值 - 零",
			minValue: 0,
			value:    0,
			expected: true,
		},
		{
			name:     "无效类型 - 字符串",
			minValue: 10,
			value:    "10",
			expected: false,
		},
		{
			name:     "无效类型 - nil",
			minValue: 10,
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewMinValue(tt.minValue)
			result := rule.Validate(tt.value)
			if result != tt.expected {
				t.Errorf("MinValue.Validate() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestMinValueGetMessage 测试错误消息生成
func TestMinValueGetMessage(t *testing.T) {
	rule := NewMinValue(18)

	// 测试默认语言
	msg := rule.GetMessage("Age")
	if msg == "" {
		t.Error("MinValue.GetMessage() returned empty string")
	}

	// 测试英文
	msgEN := rule.GetMessage("Age", uvalidator.LanguageEN)
	if msgEN == "" {
		t.Error("MinValue.GetMessage() with EN returned empty string")
	}

	// 测试中文
	msgZH := rule.GetMessage("Age", uvalidator.LanguageZH)
	if msgZH == "" {
		t.Error("MinValue.GetMessage() with ZH returned empty string")
	}

	// 验证消息不同
	if msgEN == msgZH {
		t.Error("EN and ZH messages should be different")
	}
}

// TestMinValueName 测试规则名称
func TestMinValueName(t *testing.T) {
	rule := NewMinValue(18)
	if rule.Name() != "min" {
		t.Errorf("MinValue.Name() = %v, expected 'min'", rule.Name())
	}
}
