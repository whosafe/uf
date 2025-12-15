package uconfig

import (
	"fmt"
	"os"
	"sync"

	"github.com/whosafe/uf/uerror"
)

var (
	// 存储解析后的根节点
	rootNode *Node
	rootMu   sync.RWMutex
)

// invokeCallback 执行回调逻辑
// 如果 node 是 Map，遍历其子节点调用 cb (模拟 v1 行为)
// 否则直接调用 cb
func invokeCallback(key string, node *Node, cb ICallback) error {
	if node.Kind == MappingNode {
		for childKey, childVal := range node.Children {
			if err := cb(childKey, childVal); err != nil {
				return uerror.Wrap(err, fmt.Sprintf("解析配置项 '%s.%s' 失败", key, childKey))
			}
		}
	} else {
		if err := cb(key, node); err != nil {
			return uerror.Wrap(err, fmt.Sprintf("解析配置项 '%s' 失败", key))
		}
	}
	return nil
}

// Callback 手动触发配置回调
func Callback(key string, cb ICallback) error {
	rootMu.RLock()
	defer rootMu.RUnlock()

	if rootNode == nil {
		return uerror.New("config not loaded")
	}

	if rootNode.Kind != MappingNode {
		return uerror.New("config root must be a mapping")
	}

	if child, ok := rootNode.Children[key]; ok {
		return invokeCallback(key, child, cb)
	}
	return nil
}

// ParseConfig 解析配置内容并分发回调
func ParseConfig(data []byte) error {
	// 使用 V2 自制 Parser
	node, err := Parse(data)
	if err != nil {
		return uerror.Wrap(err, "解析 YAML 失败")
	}

	rootMu.Lock()
	rootNode = node
	rootMu.Unlock()

	registryMu.RLock()
	defer registryMu.RUnlock()

	processedKeys := make(map[string]bool)

	// 处理注册的 Key
	// rootNode 应该是 MappingNode
	if rootNode.Kind != MappingNode {
		// 如果根节点不是 Map，无法通过 Key 路由，这通常不符合 Config 文件的习惯
		return uerror.New("config root must be a mapping")
	}

	for key, cb := range registry {
		processedKeys[key] = true
		if child, ok := rootNode.Children[key]; ok {
			if err := invokeCallback(key, child, cb); err != nil {
				return err
			}
		}
	}

	// 处理未知 Key
	if unknownCb != nil {
		for key, child := range rootNode.Children {
			if !processedKeys[key] {
				if err := unknownCb(key, child); err != nil {
					return uerror.Wrap(err, fmt.Sprintf("解析未知配置项 '%s' 失败", key))
				}
			}
		}
	}

	return nil
}

// Load 加载配置文件
func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return uerror.Wrap(err, "读取配置文件失败")
	}
	return ParseConfig(data)
}
