package umarshal

import (
	"testing"
)

// TestWriter 测试 Writer 基础功能
func TestWriter(t *testing.T) {
	w := AcquireWriter()
	defer ReleaseWriter(w)

	w.WriteObjectStart()
	w.WriteObjectField("name")
	w.WriteString("Alice")
	w.WriteComma()
	w.WriteObjectField("age")
	w.WriteInt(25)
	w.WriteObjectEnd()

	expected := `{"name":"Alice","age":25}`
	if string(w.Bytes()) != expected {
		t.Errorf("Expected %s, got %s", expected, string(w.Bytes()))
	}
}

// TestStringEscape 测试字符串转义
func TestStringEscape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", `"hello"`},
		{"hello\"world", `"hello\"world"`},
		{"hello\nworld", `"hello\nworld"`},
		{"hello\tworld", `"hello\tworld"`},
		{"hello\\world", `"hello\\world"`},
		{"hello\rworld", `"hello\rworld"`},
	}

	for _, tt := range tests {
		w := AcquireWriter()
		w.WriteString(tt.input)
		result := string(w.Bytes())
		ReleaseWriter(w)

		if result != tt.expected {
			t.Errorf("Input: %q, Expected: %s, Got: %s", tt.input, tt.expected, result)
		}
	}
}

// TestMarshal 测试 Marshal 函数
func TestMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"nil", nil, "null"},
		{"string", "hello", `"hello"`},
		{"int", 123, "123"},
		{"bool", true, "true"},
		{"float", 3.14, "3.14"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input)
			if err != nil {
				t.Errorf("Marshal error: %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result))
			}
		})
	}
}

// 自定义类型测试
type User struct {
	ID       int
	Username string
	Email    string
	Age      int
}

// umarshal 实现
func (u *User) Marshal(w *Writer) error {
	w.WriteObjectStart()
	w.WriteObjectField("id")
	w.WriteInt(u.ID)
	w.WriteComma()
	w.WriteObjectField("username")
	w.WriteString(u.Username)
	w.WriteComma()
	w.WriteObjectField("email")
	w.WriteString(u.Email)
	w.WriteComma()
	w.WriteObjectField("age")
	w.WriteInt(u.Age)
	w.WriteObjectEnd()
	return nil
}

func TestCustomMarshaler(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "Alice",
		Email:    "alice@example.com",
		Age:      25,
	}

	result, err := Marshal(user)
	if err != nil {
		t.Errorf("Marshal error: %v", err)
	}

	expected := `{"id":1,"username":"Alice","email":"alice@example.com","age":25}`
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

// BenchmarkMarshal 性能测试
func BenchmarkMarshal(b *testing.B) {
	user := &User{
		ID:       1,
		Username: "Alice",
		Email:    "alice@example.com",
		Age:      25,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Marshal(user)
	}
}

// BenchmarkWriter 性能测试
func BenchmarkWriter(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := AcquireWriter()
		w.WriteObjectStart()
		w.WriteObjectField("id")
		w.WriteInt(1)
		w.WriteComma()
		w.WriteObjectField("name")
		w.WriteString("Alice")
		w.WriteComma()
		w.WriteObjectField("age")
		w.WriteInt(25)
		w.WriteObjectEnd()
		ReleaseWriter(w)
	}
}
