package ubind

// parseBinary 解析 Binary 数据
// TODO: 实现自定义二进制格式解析
// 可能的格式:
// - Length-prefixed: [4 bytes length][data]
// - Type-Length-Value (TLV)
// - Protocol Buffers like
func parseBinary(data []byte) *Value {
	// TODO: 实现 Binary 解析
	// 暂时返回空对象
	return &Value{Type: TypeObject, Object: make(map[string]*Value)}
}
