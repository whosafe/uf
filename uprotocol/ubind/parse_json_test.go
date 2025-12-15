package ubind

import (
	"strings"
	"testing"
)

// TestParseJSON_DepthLimit 测试深度限制
func TestParseJSON_DepthLimit(t *testing.T) {
	// 生成深度为 150 的嵌套数组 (超过限制)
	nested := strings.Repeat("[", 150) + strings.Repeat("]", 150)
	result := parseJSON([]byte(nested))
	if result != nil {
		t.Error("应该返回 nil (超过深度限制)")
	}
}

// TestParseJSON_NormalDepth 测试正常深度
func TestParseJSON_NormalDepth(t *testing.T) {
	// 深度为 10 的正常嵌套
	nested := `[[[[[[[[[[1]]]]]]]]]]`
	result := parseJSON([]byte(nested))
	if result == nil {
		t.Error("正常深度应该解析成功")
	}
	if result.Type != TypeArray {
		t.Error("应该是数组类型")
	}
}

// TestParseJSON_DepthLimitObject 测试对象深度限制
func TestParseJSON_DepthLimitObject(t *testing.T) {
	// 生成深度为 150 的嵌套对象
	var nested strings.Builder
	for i := 0; i < 150; i++ {
		nested.WriteString(`{"a":`)
	}
	nested.WriteString("1")
	for i := 0; i < 150; i++ {
		nested.WriteString("}")
	}

	result := parseJSON([]byte(nested.String()))
	if result != nil {
		t.Error("应该返回 nil (超过深度限制)")
	}
}

// TestParseJSON_NormalDepthObject 测试正常对象深度
func TestParseJSON_NormalDepthObject(t *testing.T) {
	// 深度为 5 的正常嵌套对象
	nested := `{"a":{"b":{"c":{"d":{"e":1}}}}}`
	result := parseJSON([]byte(nested))
	if result == nil {
		t.Error("正常深度应该解析成功")
	}
	if result.Type != TypeObject {
		t.Error("应该是对象类型")
	}
}

// TestParseJSON_MixedDepth 测试混合嵌套深度
func TestParseJSON_MixedDepth(t *testing.T) {
	// 数组和对象混合嵌套,深度为 50
	var nested strings.Builder
	for i := 0; i < 25; i++ {
		nested.WriteString(`[{"a":`)
	}
	nested.WriteString("1")
	for i := 0; i < 25; i++ {
		nested.WriteString("}]")
	}

	result := parseJSON([]byte(nested.String()))
	if result == nil {
		t.Error("深度 50 应该解析成功")
	}
}

// TestParseJSON_ExactMaxDepth 测试刚好达到最大深度
func TestParseJSON_ExactMaxDepth(t *testing.T) {
	// 深度刚好为 100 (应该被允许，因为深度从0开始计数)
	nested := strings.Repeat("[", 100) + "1" + strings.Repeat("]", 100)
	result := parseJSON([]byte(nested))
	if result == nil {
		t.Error("深度 100 应该被允许")
	}
}

// TestParseJSON_JustBelowMaxDepth 测试略低于最大深度
func TestParseJSON_JustBelowMaxDepth(t *testing.T) {
	// 深度为 99
	nested := strings.Repeat("[", 99) + "1" + strings.Repeat("]", 99)
	result := parseJSON([]byte(nested))
	if result == nil {
		t.Error("深度 99 应该解析成功")
	}
}

// TestParseJSON_ExceedMaxDepth 测试超过最大深度
func TestParseJSON_ExceedMaxDepth(t *testing.T) {
	// 深度为 101 (应该被拒绝)
	nested := strings.Repeat("[", 101) + "1" + strings.Repeat("]", 101)
	result := parseJSON([]byte(nested))
	if result != nil {
		t.Error("深度 101 应该被拒绝")
	}
}
