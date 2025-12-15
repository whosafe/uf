package ubind

// ValueType JSON 值类型
type ValueType int

const (
	// TypeNull null 类型
	TypeNull ValueType = iota
	// TypeBool 布尔类型
	TypeBool
	// TypeNumber 数字类型
	TypeNumber
	// TypeString 字符串类型
	TypeString
	// TypeArray 数组类型
	TypeArray
	// TypeObject 对象类型
	TypeObject
)

// Value JSON 值
type Value struct {
	Type   ValueType
	Bool   bool
	Number float64
	String string
	Array  []*Value
	Object map[string]*Value
}

// IsNull 是否为 null
func (v *Value) IsNull() bool {
	return v.Type == TypeNull
}

// IsBool 是否为布尔值
func (v *Value) IsBool() bool {
	return v.Type == TypeBool
}

// IsNumber 是否为数字
func (v *Value) IsNumber() bool {
	return v.Type == TypeNumber
}

// IsString 是否为字符串
func (v *Value) IsString() bool {
	return v.Type == TypeString
}

// IsArray 是否为数组
func (v *Value) IsArray() bool {
	return v.Type == TypeArray
}

// IsObject 是否为对象
func (v *Value) IsObject() bool {
	return v.Type == TypeObject
}

// Int 获取整数值
func (v *Value) Int() int {
	return int(v.Number)
}

// Int64 获取 int64 值
func (v *Value) Int64() int64 {
	return int64(v.Number)
}

// Float 获取浮点数值
func (v *Value) Float() float64 {
	return v.Number
}

// Str 获取字符串值
func (v *Value) Str() string {
	return v.String
}

// Get 获取对象字段
func (v *Value) Get(key string) *Value {
	if v.Type == TypeObject && v.Object != nil {
		return v.Object[key]
	}
	return nil
}

// Index 获取数组元素
func (v *Value) Index(i int) *Value {
	if v.Type == TypeArray && i >= 0 && i < len(v.Array) {
		return v.Array[i]
	}
	return nil
}

// Len 获取数组长度
func (v *Value) Len() int {
	if v.Type == TypeArray {
		return len(v.Array)
	}
	return 0
}
