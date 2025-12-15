package uvalidator

// Validator 验证器接口
type Validator interface {
	Validate() error
}

// Rule 验证规则接口
type Rule interface {
	// Validate 执行验证
	Validate(value interface{}) bool

	// GetMessage 获取错误消息
	// lang 参数可选,如果为空则使用全局语言设置
	GetMessage(field string, lang ...Language) string

	// Name 规则名称
	Name() string
}
