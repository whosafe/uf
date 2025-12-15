package uhttp

import (
	"strings"

	"github.com/whosafe/uf/uprotocol/unet"
)

// Router 路由器
type Router struct {
	trees map[string]*node // 每个 HTTP 方法一棵树
}

// NewRouter 创建新的路由器
func NewRouter() *Router {
	return &Router{
		trees: make(map[string]*node),
	}
}

// addRoute 添加路由
func (r *Router) addRoute(method, path string, handler unet.HandlerFunc) {
	if path == "" {
		panic("path cannot be empty")
	}
	if path[0] != '/' {
		panic("path must begin with '/'")
	}
	if handler == nil {
		panic("handler cannot be nil")
	}

	root := r.trees[method]
	if root == nil {
		root = &node{path: "/"}
		r.trees[method] = root
	}

	root.addRoute(path, handler)
}

// getValue 获取路由处理器和参数
func (r *Router) getValue(method, path string) (unet.HandlerFunc, Params) {
	root := r.trees[method]
	if root == nil {
		// 尝试 ANY 方法
		root = r.trees["ANY"]
		if root == nil {
			return nil, nil
		}
	}

	handler, params := root.getValue(path)
	return handler, params
}

// node 路由树节点
type node struct {
	path     string           // 节点路径
	handler  unet.HandlerFunc // 处理器
	children []*node          // 子节点
	isWild   bool             // 是否是通配符节点 (:param 或 *path)
	param    string           // 参数名称
}

// addRoute 添加路由到节点
func (n *node) addRoute(path string, handler unet.HandlerFunc) {
	// 移除前导斜杠
	if path == "/" {
		n.handler = handler
		return
	}

	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "/")

	current := n
	for i, part := range parts {
		child := current.matchChild(part)
		if child == nil {
			// 创建新节点
			child = &node{
				path: part,
			}

			// 检查是否是参数节点
			if strings.HasPrefix(part, ":") {
				child.isWild = true
				child.param = part[1:]
			} else if strings.HasPrefix(part, "*") {
				child.isWild = true
				child.param = part[1:]
			}

			current.children = append(current.children, child)
		}

		current = child

		// 最后一个部分,设置处理器
		if i == len(parts)-1 {
			current.handler = handler
		}
	}
}

// getValue 获取处理器和参数
func (n *node) getValue(path string) (unet.HandlerFunc, Params) {
	// 根路径
	if path == "/" {
		return n.handler, nil
	}

	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "/")
	params := make(Params)

	current := n
	for i, part := range parts {
		child := current.matchChild(part)
		if child == nil {
			return nil, nil
		}

		// 如果是参数节点,保存参数值
		if child.isWild && child.param != "" {
			params[child.param] = part
		}

		current = child

		// 最后一个部分,返回处理器
		if i == len(parts)-1 {
			return current.handler, params
		}
	}

	return nil, nil
}

// matchChild 匹配子节点
func (n *node) matchChild(part string) *node {
	// 精确匹配
	for _, child := range n.children {
		if child.path == part {
			return child
		}
	}

	// 通配符匹配
	for _, child := range n.children {
		if child.isWild {
			return child
		}
	}

	return nil
}
