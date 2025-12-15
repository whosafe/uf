package ubind

// Parse 自动识别格式并解析
func Parse(data []byte) *Value {
	if len(data) == 0 {
		return &Value{Type: TypeObject, Object: make(map[string]*Value)}
	}

	// 尝试 JSON (以 { 或 [ 开头)
	if data[0] == '{' || data[0] == '[' {
		if val := parseJSON(data); val != nil {
			return val
		}
	}

	// 尝试 Form
	if val := tryParseForm(data); val != nil {
		return val
	}

	// 尝试 Binary
	if val := tryParseBinary(data); val != nil {
		return val
	}

	// 默认返回空对象
	return &Value{Type: TypeObject, Object: make(map[string]*Value)}
}

// ParseJSON 强制使用 JSON 解析
func ParseJSON(data []byte) *Value {
	return parseJSON(data)
}

// ParseForm 强制使用 Form 解析
func ParseForm(data []byte) *Value {
	return parseForm(data)
}

// ParseBinary 强制使用 Binary 解析
func ParseBinary(data []byte) *Value {
	return parseBinary(data)
}

// tryParseForm 尝试解析 Form 数据
func tryParseForm(data []byte) *Value {
	// 检测是否包含 = 和 & (Form 格式特征)
	hasEqual := false
	for _, b := range data {
		if b == '=' {
			hasEqual = true
			break
		}
	}

	if hasEqual {
		return parseForm(data)
	}

	return nil
}

// tryParseBinary 尝试解析 Binary 数据
func tryParseBinary(data []byte) *Value {
	// TODO: 实现 Binary 解析
	return nil
}
