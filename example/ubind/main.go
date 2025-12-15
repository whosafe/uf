package main

import (
	"fmt"

	"github.com/whosafe/uf/uprotocol/ubind"
)

// User 用户结构
type User struct {
	ID   int
	Name string
	Age  int
}

func (u *User) Bind(key string, value *ubind.Value) error {
	switch key {
	case "id":
		if value.IsNumber() {
			u.ID = value.Int()
		} else if value.IsString() {
			// Form 数据是字符串,需要手动转换
			fmt.Sscanf(value.Str(), "%d", &u.ID)
		}
	case "name":
		u.Name = value.Str()
	case "age":
		if value.IsNumber() {
			u.Age = value.Int()
		} else if value.IsString() {
			// Form 数据是字符串,需要手动转换
			fmt.Sscanf(value.Str(), "%d", &u.Age)
		}
	}
	return nil
}

// Address 地址结构
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
				ubind.Bind(value.Index(i), &t.Members[i])
			}
		}
	}
	return nil
}

func main() {
	fmt.Println("=== ubind 数据绑定示例 ===")

	// 示例 1: JSON 简单对象
	fmt.Println("--- 示例 1: JSON 简单对象 ---")
	jsonData := []byte(`{"id":1,"name":"Alice","age":25}`)
	fmt.Printf("JSON: %s\n", string(jsonData))

	val := ubind.Parse(jsonData)
	var user User
	ubind.Bind(val, &user)

	fmt.Printf("解析结果: ID=%d, Name=%s, Age=%d\n", user.ID, user.Name, user.Age)
	fmt.Println()

	// 示例 2: JSON 嵌套对象
	fmt.Println("--- 示例 2: JSON 嵌套对象 ---")
	jsonNested := []byte(`{"id":2,"name":"Bob","address":{"city":"Beijing","street":"Main St"}}`)
	fmt.Printf("JSON: %s\n", string(jsonNested))

	val2 := ubind.Parse(jsonNested)
	var userWithAddr UserWithAddress
	ubind.Bind(val2, &userWithAddr)

	fmt.Printf("解析结果: ID=%d, Name=%s, City=%s, Street=%s\n",
		userWithAddr.ID, userWithAddr.Name, userWithAddr.Address.City, userWithAddr.Address.Street)
	fmt.Println()

	// 示例 3: JSON 数组
	fmt.Println("--- 示例 3: JSON 数组 ---")
	jsonArray := []byte(`{"name":"Dev Team","members":[{"id":1,"name":"Alice","age":25},{"id":2,"name":"Bob","age":30}]}`)
	fmt.Printf("JSON: %s\n", string(jsonArray))

	val3 := ubind.Parse(jsonArray)
	var team Team
	ubind.Bind(val3, &team)

	fmt.Printf("解析结果: Team=%s, Members=%d\n", team.Name, len(team.Members))
	for i, member := range team.Members {
		fmt.Printf("  Member %d: ID=%d, Name=%s, Age=%d\n", i+1, member.ID, member.Name, member.Age)
	}
	fmt.Println()

	// 示例 4: Form 数据
	fmt.Println("--- 示例 4: Form 数据 ---")
	formData := []byte("name=Charlie&age=28&id=3")
	fmt.Printf("Form: %s\n", string(formData))

	val4 := ubind.Parse(formData)
	var user2 User
	ubind.Bind(val4, &user2)

	fmt.Printf("解析结果: ID=%d, Name=%s, Age=%d\n", user2.ID, user2.Name, user2.Age)
	fmt.Println()

	// 示例 5: Form URL 编码
	fmt.Println("--- 示例 5: Form URL 编码 ---")
	formEncoded := []byte("name=Hello+World&id=4&age=35")
	fmt.Printf("Form: %s\n", string(formEncoded))

	val5 := ubind.Parse(formEncoded)
	var user3 User
	ubind.Bind(val5, &user3)

	fmt.Printf("解析结果: ID=%d, Name=%s, Age=%d\n", user3.ID, user3.Name, user3.Age)
	fmt.Println("(注意: 'Hello+World' 被解码为 'Hello World')")
	fmt.Println()

	// 示例 6: 自动格式识别
	fmt.Println("--- 示例 6: 自动格式识别 ---")
	fmt.Println("Parse() 会自动识别 JSON 或 Form 格式:")
	fmt.Println("  - JSON: 以 { 或 [ 开头")
	fmt.Println("  - Form: 包含 = 字符")
	fmt.Println()

	fmt.Println("所有示例完成! ✓")
}
