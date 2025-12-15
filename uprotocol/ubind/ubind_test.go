package ubind_test

import (
	"testing"

	"github.com/whosafe/uf/uprotocol/ubind"
)

// User 测试用户结构
type User struct {
	ID   int
	Name string
	Age  int
}

func (u *User) Bind(key string, value *ubind.Value) error {
	switch key {
	case "id":
		u.ID = value.Int()
	case "name":
		u.Name = value.Str()
	case "age":
		u.Age = value.Int()
	}
	return nil
}

func TestParseJSON(t *testing.T) {
	jsonData := []byte(`{"id":1,"name":"Alice","age":25}`)

	val := ubind.Parse(jsonData)
	if val == nil {
		t.Fatal("Parse returned nil")
	}

	if !val.IsObject() {
		t.Fatal("Expected object type")
	}

	// 测试字段
	if id := val.Get("id"); id == nil || id.Int() != 1 {
		t.Errorf("Expected id=1, got %v", id)
	}

	if name := val.Get("name"); name == nil || name.Str() != "Alice" {
		t.Errorf("Expected name=Alice, got %v", name)
	}

	if age := val.Get("age"); age == nil || age.Int() != 25 {
		t.Errorf("Expected age=25, got %v", age)
	}
}

func TestBind(t *testing.T) {
	jsonData := []byte(`{"id":1,"name":"Alice","age":25}`)

	val := ubind.Parse(jsonData)
	var user User
	err := ubind.Bind(val, &user)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("Expected ID=1, got %d", user.ID)
	}
	if user.Name != "Alice" {
		t.Errorf("Expected Name=Alice, got %s", user.Name)
	}
	if user.Age != 25 {
		t.Errorf("Expected Age=25, got %d", user.Age)
	}
}

// Address 测试地址结构
type Address struct {
	City   string
	Street string
}

func (a *Address) Bind(key string, value *ubind.Value) error {
	switch key {
	case "city":
		a.City = value.Str()
	case "street":
		a.Street = value.Str()
	}
	return nil
}

// UserWithAddress 带地址的用户
type UserWithAddress struct {
	ID      int
	Name    string
	Address Address
}

func (u *UserWithAddress) Bind(key string, value *ubind.Value) error {
	switch key {
	case "id":
		u.ID = value.Int()
	case "name":
		u.Name = value.Str()
	case "address":
		if value.IsObject() {
			return ubind.Bind(value, &u.Address)
		}
	}
	return nil
}

func TestNestedObject(t *testing.T) {
	jsonData := []byte(`{"id":1,"name":"Alice","address":{"city":"Beijing","street":"Main St"}}`)

	val := ubind.Parse(jsonData)
	var user UserWithAddress
	err := ubind.Bind(val, &user)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("Expected ID=1, got %d", user.ID)
	}
	if user.Name != "Alice" {
		t.Errorf("Expected Name=Alice, got %s", user.Name)
	}
	if user.Address.City != "Beijing" {
		t.Errorf("Expected City=Beijing, got %s", user.Address.City)
	}
	if user.Address.Street != "Main St" {
		t.Errorf("Expected Street=Main St, got %s", user.Address.Street)
	}
}

// Team 团队结构
type Team struct {
	Name    string
	Members []User
}

func (t *Team) Bind(key string, value *ubind.Value) error {
	switch key {
	case "name":
		t.Name = value.Str()
	case "members":
		if value.IsArray() {
			t.Members = make([]User, value.Len())
			for i := 0; i < value.Len(); i++ {
				elem := value.Index(i)
				if elem.IsObject() {
					ubind.Bind(elem, &t.Members[i])
				}
			}
		}
	}
	return nil
}

func TestArray(t *testing.T) {
	jsonData := []byte(`{"name":"Dev Team","members":[{"id":1,"name":"Alice","age":25},{"id":2,"name":"Bob","age":30}]}`)

	val := ubind.Parse(jsonData)
	var team Team
	err := ubind.Bind(val, &team)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	if team.Name != "Dev Team" {
		t.Errorf("Expected Name=Dev Team, got %s", team.Name)
	}
	if len(team.Members) != 2 {
		t.Fatalf("Expected 2 members, got %d", len(team.Members))
	}
	if team.Members[0].Name != "Alice" {
		t.Errorf("Expected first member=Alice, got %s", team.Members[0].Name)
	}
	if team.Members[1].Name != "Bob" {
		t.Errorf("Expected second member=Bob, got %s", team.Members[1].Name)
	}
}
