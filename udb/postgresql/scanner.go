package postgresql

// Scanner 扫描器接口（零反射）
// 用户需要实现此接口来扫描查询结果
type Scanner interface {
	Scan(key string, values any) error
}
