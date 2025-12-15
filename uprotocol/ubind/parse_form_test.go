package ubind_test

import (
	"testing"

	"github.com/whosafe/uf/uprotocol/ubind"
)

func TestParseForm(t *testing.T) {
	formData := []byte("name=Alice&age=25&city=Beijing")

	val := ubind.Parse(formData)
	if val == nil {
		t.Fatal("Parse returned nil")
	}

	if !val.IsObject() {
		t.Fatal("Expected object type")
	}

	// 测试字段
	if name := val.Get("name"); name == nil || name.Str() != "Alice" {
		t.Errorf("Expected name=Alice, got %v", name)
	}

	if age := val.Get("age"); age == nil || age.Str() != "25" {
		t.Errorf("Expected age=25, got %v", age)
	}

	if city := val.Get("city"); city == nil || city.Str() != "Beijing" {
		t.Errorf("Expected city=Beijing, got %v", city)
	}
}

func TestParseFormWithURLEncoding(t *testing.T) {
	// URL 编码: "Hello World" -> "Hello+World"
	formData := []byte("message=Hello+World&email=test%40example.com")

	val := ubind.Parse(formData)
	if val == nil {
		t.Fatal("Parse returned nil")
	}

	// 测试空格解码
	if msg := val.Get("message"); msg == nil || msg.Str() != "Hello World" {
		t.Errorf("Expected message='Hello World', got '%s'", msg.Str())
	}

	// 测试 %XX 解码
	if email := val.Get("email"); email == nil || email.Str() != "test@example.com" {
		t.Errorf("Expected email='test@example.com', got '%s'", email.Str())
	}
}

func TestParseFormEmpty(t *testing.T) {
	formData := []byte("")

	val := ubind.Parse(formData)
	if val == nil {
		t.Fatal("Parse returned nil")
	}

	if !val.IsObject() {
		t.Fatal("Expected object type")
	}
}

func TestBindForm(t *testing.T) {
	type FormData struct {
		Name  string
		Email string
		Age   string
	}

	var fd FormData

	// 实现 Bind 方法
	bindFunc := func(key string, value *ubind.Value) error {
		switch key {
		case "name":
			fd.Name = value.Str()
		case "email":
			fd.Email = value.Str()
		case "age":
			fd.Age = value.Str()
		}
		return nil
	}

	// 创建匿名结构体实现 Binder
	type formBinder struct {
		bindFn func(string, *ubind.Value) error
	}
	fb := &formBinder{bindFn: bindFunc}

	formData := []byte("name=Alice&email=alice%40example.com&age=25")
	val := ubind.Parse(formData)

	// 手动调用 bind
	for key, value := range val.Object {
		fb.bindFn(key, value)
	}

	if fd.Name != "Alice" {
		t.Errorf("Expected Name=Alice, got %s", fd.Name)
	}
	if fd.Email != "alice@example.com" {
		t.Errorf("Expected Email=alice@example.com, got %s", fd.Email)
	}
	if fd.Age != "25" {
		t.Errorf("Expected Age=25, got %s", fd.Age)
	}
}
