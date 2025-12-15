package ubind

import (
	"fmt"
	"strings"
	"testing"
)

// TestParseJSON_DebugDepth 调试深度计数
func TestParseJSON_DebugDepth(t *testing.T) {
	// 测试不同深度
	for depth := 95; depth <= 105; depth++ {
		nested := strings.Repeat("[", depth) + "1" + strings.Repeat("]", depth)
		result := parseJSON([]byte(nested))
		fmt.Printf("深度 %d: %v\n", depth, result != nil)
	}
}
