package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CreateDir 创建目录
func CreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// WriteFile 写入文件
func WriteFile(path string, content string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := CreateDir(dir); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return WriteFile(dst, string(data))
}

// RenderTemplate 渲染模板
func RenderTemplate(template string, vars map[string]string) string {
	result := template
	for key, value := range vars {
		// 支持 {{.Key}} 和 {{Key}} 两种格式
		placeholder1 := "{{." + key + "}}"
		placeholder2 := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder1, value)
		result = strings.ReplaceAll(result, placeholder2, value)
	}
	return result
}

// RunCommand 执行命令
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommandInDir 在指定目录执行命令
func RunCommandInDir(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetCurrentDir 获取当前目录
func GetCurrentDir() (string, error) {
	return os.Getwd()
}

// ToSnakeCase 转换为蛇形命名
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ToPascalCase 转换为帕斯卡命名
func ToPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})

	var result strings.Builder
	for _, word := range words {
		if len(word) > 0 {
			result.WriteString(strings.ToUpper(word[:1]))
			if len(word) > 1 {
				result.WriteString(strings.ToLower(word[1:]))
			}
		}
	}
	return result.String()
}

// PrintStep 打印步骤信息
func PrintStep(step int, total int, msg string) {
	fmt.Printf("[%d/%d] %s\n", step, total, msg)
}
