package uconfig

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/whosafe/uf/uerror"
)

// parser 简单的 YAML 解析器上下文
type parser struct {
	lines   []string
	current int
}

// Parse 解析 YAML 字节流 (极简实现)
func Parse(data []byte) (*Node, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	p := &parser{
		lines:   lines,
		current: 0,
	}

	root := &Node{
		Kind:     MappingNode,
		Children: make(map[string]*Node),
	}

	// 解析顶层
	if err := p.parseBlock(root, -1); err != nil {
		return nil, err
	}

	return root, nil
}

// getIndent 计算行缩进空格数
func getIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

// parseBlock 解析一个缩进块
// parent: 父节点
// parentIndent: 父节点的缩进级别
func (p *parser) parseBlock(parent *Node, parentIndent int) error {
	// 委托给 parseRecursive 处理实际的解析逻辑
	return parseRecursive(p, parent, parentIndent)
}

func parseRecursive(p *parser, parent *Node, minIndent int) error {
	for p.current < len(p.lines) {
		line := p.lines[p.current]
		trimLine := strings.TrimSpace(line)

		if trimLine == "" || strings.HasPrefix(trimLine, "#") {
			p.current++
			continue
		}

		indent := getIndent(line)
		if indent <= minIndent && minIndent != -1 {
			// 缩进小于等于当前块的最小缩进要求（实际上 minIndent 应该是 parent ident）
			// 如果 indent <= parent's indent, return.
			// 这里的 minIndent 传入的是 "block indent".
			// 初始 -1. 顶层 indent 0.
			// 如果是子块，indent 必须 > parentIndent.
			// 我们动态探测当前块的 indent.
			return nil
		}

		// 列表处理
		if strings.HasPrefix(trimLine, "- ") {
			if parent.Kind != SequenceNode && len(parent.Children) == 0 && len(parent.List) == 0 {
				parent.Kind = SequenceNode
			}
			if parent.Kind != SequenceNode {
				return uerror.NewWithCode(1, fmt.Sprintf("line %d: mixed mapping and sequence", p.current+1))
			}

			// 列表项内容
			valStr := strings.TrimPrefix(trimLine, "- ")
			valStr = strings.TrimSpace(valStr)

			child := &Node{}
			if valStr != "" {
				// "- value" -> Scalar item
				child.Kind = ScalarNode
				child.Value = removeQuotes(valStr) // 处理引号
			} else {
				// "- " -> Block item (nested object/list)
				// 预读取下一行看缩进
				p.current++ // 消费当前行 "- "
				// 递归解析子节点作为 child
				// child 应该是一个 Map (通常列表里是对象)
				child.Kind = MappingNode // 默认为 Map，除非下文发现是 List
				child.Children = make(map[string]*Node)
				if err := parseRecursive(p, child, indent); err != nil {
					return err
				}
				parent.List = append(parent.List, child)
				continue
			}
			parent.List = append(parent.List, child)
			p.current++
			continue
		}

		// Map 处理 "key: value" or "key:"
		colonIdx := strings.Index(trimLine, ":")
		if colonIdx == -1 {
			// 既不是列表项 "- " 也没有 separator ": " -> 可能是纯 Scalar 或者是错误?
			// 在 Config 文件中，顶层通常是 Key: Value
			return uerror.NewWithCode(1, fmt.Sprintf("line %d: expected key: value", p.current+1))
		}

		key := strings.TrimSpace(trimLine[:colonIdx])
		valStr := strings.TrimSpace(trimLine[colonIdx+1:])

		if parent.Kind != MappingNode {
			// 如果之前已经被标记为 Sequence，这里出错了
			if parent.Kind == SequenceNode {
				return uerror.NewWithCode(1, fmt.Sprintf("line %d: mixed mapping and sequence", p.current+1))
			}
			parent.Kind = MappingNode
			if parent.Children == nil {
				parent.Children = make(map[string]*Node)
			}
		}

		child := &Node{}

		if valStr != "" {
			// Check if value is potentially a comment? e.g. "key: value # comment"
			if idx := strings.Index(valStr, " #"); idx != -1 {
				valStr = strings.TrimSpace(valStr[:idx])
			}

			// 简单的 Scalar
			child.Kind = ScalarNode
			child.Value = removeQuotes(valStr)
			parent.Children[key] = child
			p.current++
		} else {
			// "key:" -> Block
			// 下一行开始是子内容
			p.current++ // 消费 "key:"

			// 探测下一行缩进
			if p.current < len(p.lines) {
				nextLine := p.lines[p.current]
				nextIndent := getIndent(nextLine)
				if nextIndent <= indent {
					// 下一行缩进没有增加，说明 key 的值是空的 (Scalar empty)
					child.Kind = ScalarNode
					child.Value = ""
				} else {
					// 递归解析子块
					// 子块的 minIndent 是当前 indent (key 的 indent)
					// parseRecursive 会读取直到 indent <= current indent

					// 默认子节点初始化为 Map，由内部逻辑决定是否转 List
					child.Kind = MappingNode
					child.Children = make(map[string]*Node)

					if err := parseRecursive(p, child, indent); err != nil {
						return err
					}
				}
			} else {
				// 文件结束，Value 为空
				child.Kind = ScalarNode
				child.Value = ""
			}
			parent.Children[key] = child
		}
	}
	return nil
}

func removeQuotes(s string) string {
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}
