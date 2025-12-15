package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Input 获取用户输入
func Input(prompt string, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" && defaultValue != "" {
		return defaultValue
	}

	return input
}

// Select 选择菜单
func Select(prompt string, options []string, defaultIndex int) int {
	fmt.Println(prompt)
	for i, opt := range options {
		if i == defaultIndex {
			fmt.Printf("  > %s\n", opt)
		} else {
			fmt.Printf("    %s\n", opt)
		}
	}

	fmt.Printf("\n选择 [1-%d] (默认: %d): ", len(options), defaultIndex+1)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultIndex
	}

	var choice int
	fmt.Sscanf(input, "%d", &choice)
	choice-- // 转换为 0-based index

	if choice < 0 || choice >= len(options) {
		return defaultIndex
	}

	return choice
}

// Confirm 确认对话框
func Confirm(prompt string, defaultYes bool) bool {
	var suffix string
	if defaultYes {
		suffix = "[Y/n]"
	} else {
		suffix = "[y/N]"
	}

	fmt.Printf("%s %s: ", prompt, suffix)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))

	if input == "" {
		return defaultYes
	}

	return input == "y" || input == "yes"
}

// Success 成功消息
func Success(msg string) {
	fmt.Printf("✓ %s\n", msg)
}

// Error 错误消息
func Error(msg string) {
	fmt.Printf("✗ %s\n", msg)
}

// Info 信息消息
func Info(msg string) {
	fmt.Printf("ℹ %s\n", msg)
}

// Warning 警告消息
func Warning(msg string) {
	fmt.Printf("⚠ %s\n", msg)
}
