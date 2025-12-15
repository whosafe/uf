package uconfig

import (
	"fmt"

	"github.com/whosafe/uf/uerror"
)

// NodeType 节点类型
type NodeType int

const (
	ScalarNode NodeType = iota
	MappingNode
	SequenceNode
)

// Node 配置节点
type Node struct {
	Kind     NodeType
	Value    string
	Children map[string]*Node // Mapping
	List     []*Node          // Sequence

	// 为了保持顺序便于测试或显示，可选
	Keys []string
}

// String 简单返回 Value
func (n *Node) String() string {
	return n.Value
}

// Iter 遍历 Sequence
func (n *Node) Iter(cb func(i int, v *Node) error) error {
	if n.Kind != SequenceNode {
		return uerror.New("node is not a list")
	}
	for i, v := range n.List {
		if err := cb(i, v); err != nil {
			return err
		}
	}
	return nil
}

// Unmarshaler 接口
type Unmarshaler interface {
	UnmarshalYAML(key string, value *Node) error
}

// Decode 解析到 Struct (Simplistic)
// 只支持实现了 Unmarshaler 的 Struct
func (n *Node) Decode(v any) error {
	// 如果 v 实现了 Unmarshaler
	// 我们需要 cast v 为 Unmarshaler interface
	if u, ok := v.(interface{ UnmarshalYAML(string, *Node) error }); ok {
		if n.Kind != MappingNode {
			return uerror.New("cannot decode non-map node to struct")
		}
		for k, child := range n.Children {
			if err := u.UnmarshalYAML(k, child); err != nil {
				return err
			}
		}
		return nil
	}

	return uerror.New(fmt.Sprintf("type %T does not implement UnmarshalYAML(key string, value *Node) error", v))
}
