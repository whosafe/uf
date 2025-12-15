package umarshal

// Marshaler 自定义序列化接口
type IMarshaler interface {
	Marshal(w *Writer) error
}

// Marshal 序列化任意值
func Marshal(v any) ([]byte, error) {
	w := AcquireWriter()
	defer ReleaseWriter(w)

	if err := MarshalToWriter(w, v); err != nil {
		return nil, err
	}

	// 复制结果
	result := make([]byte, len(w.buf))
	copy(result, w.buf)
	return result, nil
}

// MarshalToWriter 序列化到 Writer
// 支持的类型:
// - 实现了 IMarshaler 接口的类型
// - 基础类型: string, int, int64, uint, uint64, float32, float64, bool, []byte
// - map[string]any, map[string]string, map[string]int
// - []any (需要元素类型也是支持的类型)
// - 其他类型将序列化为 null
func MarshalToWriter(w *Writer, v any) error {
	if v == nil {
		w.WriteNull()
		return nil
	}

	// 检查是否实现了 Marshaler 接口
	if m, ok := v.(IMarshaler); ok {
		return m.Marshal(w)
	}

	// 基础类型
	switch val := v.(type) {
	case string:
		w.WriteString(val)
	case int:
		w.WriteInt(val)
	case int64:
		w.WriteInt64(val)
	case uint:
		w.WriteUint(val)
	case uint64:
		w.WriteUint64(val)
	case float32:
		w.WriteFloat32(val)
	case float64:
		w.WriteFloat64(val)
	case bool:
		w.WriteBool(val)
	case []byte:
		w.WriteString(string(val))

	// map 类型
	case map[string]any:
		w.WriteObjectStart()
		first := true
		for k, v := range val {
			if !first {
				w.WriteComma()
			}
			first = false
			w.WriteObjectField(k)
			MarshalToWriter(w, v)
		}
		w.WriteObjectEnd()
	case map[string]string:
		w.WriteObjectStart()
		first := true
		for k, v := range val {
			if !first {
				w.WriteComma()
			}
			first = false
			w.WriteObjectField(k)
			w.WriteString(v)
		}
		w.WriteObjectEnd()
	case map[string]int:
		w.WriteObjectStart()
		first := true
		for k, v := range val {
			if !first {
				w.WriteComma()
			}
			first = false
			w.WriteObjectField(k)
			w.WriteInt(v)
		}
		w.WriteObjectEnd()

	// 数组/切片类型 - 只支持 []any
	case []any:
		w.WriteArrayStart()
		for i, item := range val {
			if i > 0 {
				w.WriteComma()
			}
			MarshalToWriter(w, item)
		}
		w.WriteArrayEnd()

	default:
		// 不支持的类型,返回 null
		// 如果需要序列化自定义类型,请实现 IMarshaler 接口
		w.WriteNull()
	}

	return nil
}
