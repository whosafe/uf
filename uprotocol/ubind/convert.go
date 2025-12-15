package ubind

import "net/url"

// FromURLValues 从 url.Values 创建 Value
// url.Values 是 map[string][]string 类型,用于表单和查询参数
func FromURLValues(values url.Values) *Value {
	if values == nil {
		return &Value{Type: TypeNull}
	}

	obj := make(map[string]*Value)
	for key, vals := range values {
		if len(vals) == 0 {
			continue
		}
		// 如果只有一个值,直接使用字符串
		if len(vals) == 1 {
			obj[key] = &Value{Type: TypeString, String: vals[0]}
		} else {
			// 如果有多个值,创建数组
			arr := make([]*Value, len(vals))
			for i, val := range vals {
				arr[i] = &Value{Type: TypeString, String: val}
			}
			obj[key] = &Value{Type: TypeArray, Array: arr}
		}
	}

	return &Value{
		Type:   TypeObject,
		Object: obj,
	}
}

// FromMap 从 map[string]any 创建 Value
func FromMap(data map[string]any) *Value {
	if data == nil {
		return &Value{Type: TypeNull}
	}

	obj := make(map[string]*Value)
	for key, val := range data {
		obj[key] = FromAny(val)
	}

	return &Value{
		Type:   TypeObject,
		Object: obj,
	}
}

// FromAny 从 any 创建 Value
func FromAny(data any) *Value {
	if data == nil {
		return &Value{Type: TypeNull}
	}

	switch v := data.(type) {
	case bool:
		return &Value{Type: TypeBool, Bool: v}
	case int:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case int8:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case int16:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case int32:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case int64:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case uint:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case uint8:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case uint16:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case uint32:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case uint64:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case float32:
		return &Value{Type: TypeNumber, Number: float64(v)}
	case float64:
		return &Value{Type: TypeNumber, Number: v}
	case string:
		return &Value{Type: TypeString, String: v}
	case []any:
		arr := make([]*Value, len(v))
		for i, item := range v {
			arr[i] = FromAny(item)
		}
		return &Value{Type: TypeArray, Array: arr}
	case map[string]any:
		return FromMap(v)
	default:
		// 其他类型转为字符串
		return &Value{Type: TypeString, String: ""}
	}
}
