package uconfig

import "sync"

// ICallback 定义配置解析回调函数 (v2)
// key: 键名
// value: *Node
type ICallback func(key string, value *Node) error

var (
	registry   = make(map[string]ICallback)
	unknownCb  ICallback
	registryMu sync.RWMutex
)

// Register 注册配置回调
func Register(key string, cb ICallback) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[key] = cb
}

// RegisterUnknown 注册未知配置回调
func RegisterUnknown(cb ICallback) {
	registryMu.Lock()
	defer registryMu.Unlock()
	unknownCb = cb
}
