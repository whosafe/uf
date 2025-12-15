package ubind

import (
	"net/url"
	"testing"
)

func TestFromURLValues(t *testing.T) {
	// 测试单值
	values := url.Values{
		"name": []string{"Alice"},
		"age":  []string{"25"},
	}

	val := FromURLValues(values)
	if val.Type != TypeObject {
		t.Errorf("期望类型为 TypeObject, 实际为 %v", val.Type)
	}

	nameVal := val.Get("name")
	if nameVal == nil || nameVal.Type != TypeString || nameVal.String != "Alice" {
		t.Errorf("name 字段解析错误: %v", nameVal)
	}

	ageVal := val.Get("age")
	if ageVal == nil || ageVal.Type != TypeString || ageVal.String != "25" {
		t.Errorf("age 字段解析错误: %v", ageVal)
	}
}

func TestFromURLValuesMultiple(t *testing.T) {
	// 测试多值(数组)
	values := url.Values{
		"tags": []string{"go", "web", "framework"},
	}

	val := FromURLValues(values)
	tagsVal := val.Get("tags")
	if tagsVal == nil || tagsVal.Type != TypeArray {
		t.Errorf("tags 应该是数组类型")
	}

	if tagsVal.Len() != 3 {
		t.Errorf("期望数组长度为 3, 实际为 %d", tagsVal.Len())
	}

	if tagsVal.Index(0).String != "go" {
		t.Errorf("tags[0] 应该是 'go', 实际为 '%s'", tagsVal.Index(0).String)
	}
}

func TestFromURLValuesEmpty(t *testing.T) {
	// 测试空值
	values := url.Values{}
	val := FromURLValues(values)

	if val.Type != TypeObject {
		t.Errorf("期望类型为 TypeObject, 实际为 %v", val.Type)
	}

	if len(val.Object) != 0 {
		t.Errorf("期望空对象, 实际有 %d 个字段", len(val.Object))
	}
}

func TestFromURLValuesNil(t *testing.T) {
	// 测试 nil
	val := FromURLValues(nil)

	if val.Type != TypeNull {
		t.Errorf("期望类型为 TypeNull, 实际为 %v", val.Type)
	}
}

func TestFromURLValuesMixed(t *testing.T) {
	// 测试混合单值和多值
	values := url.Values{
		"name":  []string{"Alice"},
		"tags":  []string{"go", "web"},
		"email": []string{"test@example.com"},
	}

	val := FromURLValues(values)

	// 检查单值
	nameVal := val.Get("name")
	if nameVal == nil || nameVal.Type != TypeString {
		t.Errorf("name 应该是字符串类型")
	}

	// 检查多值
	tagsVal := val.Get("tags")
	if tagsVal == nil || tagsVal.Type != TypeArray {
		t.Errorf("tags 应该是数组类型")
	}

	// 检查单值
	emailVal := val.Get("email")
	if emailVal == nil || emailVal.Type != TypeString {
		t.Errorf("email 应该是字符串类型")
	}
}

func TestFromURLValuesBind(t *testing.T) {
	// 测试与 Bind 配合使用
	type User struct {
		Name string
		Age  string // 改为字符串,因为 FromURLValues 返回的都是字符串
	}

	values := url.Values{
		"name": []string{"Bob"},
		"age":  []string{"30"},
	}

	val := FromURLValues(values)
	user := &User{}

	// 简单的 Bind 实现
	err := Bind(val, &testBinder{
		bindFunc: func(key string, value *Value) error {
			switch key {
			case "name":
				user.Name = value.String
			case "age":
				user.Age = value.String
			}
			return nil
		},
	})

	if err != nil {
		t.Errorf("绑定失败: %v", err)
	}

	if user.Name != "Bob" {
		t.Errorf("期望 Name 为 'Bob', 实际为 '%s'", user.Name)
	}

	if user.Age != "30" {
		t.Errorf("期望 Age 为 '30', 实际为 '%s'", user.Age)
	}
}

// testBinder 用于测试的 Binder 实现
type testBinder struct {
	bindFunc func(key string, value *Value) error
}

func (b *testBinder) Bind(key string, value *Value) error {
	return b.bindFunc(key, value)
}
