package ubind

// Binder 数据绑定接口
// 所有需要从请求中绑定数据的结构体都应该实现这个接口
type Binder interface {
	Bind(key string, value *Value) error
}
