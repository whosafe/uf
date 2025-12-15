package ubind

import (
	"strconv"
	"strings"
)

// parseJSON 解析 JSON 数据
func parseJSON(data []byte) *Value {
	p := &jsonParser{
		data:     data,
		pos:      0,
		depth:    0,
		maxDepth: 100, // 【安全修复】限制最大递归深度为 100,防止栈溢出 DoS
	}
	return p.parseValue()
}

// jsonParser JSON 解析器
type jsonParser struct {
	data     []byte
	pos      int
	depth    int // 当前递归深度
	maxDepth int // 最大递归深度限制
}

// parseValue 解析值
func (p *jsonParser) parseValue() *Value {
	p.skipWhitespace()

	if p.pos >= len(p.data) {
		return nil
	}

	ch := p.data[p.pos]

	switch ch {
	case 'n':
		return p.parseNull()
	case 't', 'f':
		return p.parseBool()
	case '"':
		return p.parseString()
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return p.parseNumber()
	default:
		return nil
	}
}

// parseNull 解析 null
func (p *jsonParser) parseNull() *Value {
	if p.pos+4 <= len(p.data) && string(p.data[p.pos:p.pos+4]) == "null" {
		p.pos += 4
		return &Value{Type: TypeNull}
	}
	return nil
}

// parseBool 解析布尔值
func (p *jsonParser) parseBool() *Value {
	if p.pos+4 <= len(p.data) && string(p.data[p.pos:p.pos+4]) == "true" {
		p.pos += 4
		return &Value{Type: TypeBool, Bool: true}
	}
	if p.pos+5 <= len(p.data) && string(p.data[p.pos:p.pos+5]) == "false" {
		p.pos += 5
		return &Value{Type: TypeBool, Bool: false}
	}
	return nil
}

// parseNumber 解析数字
func (p *jsonParser) parseNumber() *Value {
	start := p.pos

	// 负号
	if p.pos < len(p.data) && p.data[p.pos] == '-' {
		p.pos++
	}

	// 整数部分
	for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
		p.pos++
	}

	// 小数部分
	if p.pos < len(p.data) && p.data[p.pos] == '.' {
		p.pos++
		for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.pos++
		}
	}

	// 指数部分
	if p.pos < len(p.data) && (p.data[p.pos] == 'e' || p.data[p.pos] == 'E') {
		p.pos++
		if p.pos < len(p.data) && (p.data[p.pos] == '+' || p.data[p.pos] == '-') {
			p.pos++
		}
		for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.pos++
		}
	}

	numStr := string(p.data[start:p.pos])
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return nil
	}

	return &Value{Type: TypeNumber, Number: num}
}

// parseString 解析字符串
func (p *jsonParser) parseString() *Value {
	if p.data[p.pos] != '"' {
		return nil
	}
	p.pos++ // 跳过开始的 "

	var buf strings.Builder
	for p.pos < len(p.data) {
		ch := p.data[p.pos]

		if ch == '"' {
			p.pos++ // 跳过结束的 "
			return &Value{Type: TypeString, String: buf.String()}
		}

		if ch == '\\' && p.pos+1 < len(p.data) {
			p.pos++
			next := p.data[p.pos]
			switch next {
			case '"', '\\', '/':
				buf.WriteByte(next)
			case 'b':
				buf.WriteByte('\b')
			case 'f':
				buf.WriteByte('\f')
			case 'n':
				buf.WriteByte('\n')
			case 'r':
				buf.WriteByte('\r')
			case 't':
				buf.WriteByte('\t')
			default:
				buf.WriteByte(next)
			}
			p.pos++
		} else {
			buf.WriteByte(ch)
			p.pos++
		}
	}

	return nil
}

// parseObject 解析对象
func (p *jsonParser) parseObject() *Value {
	// 【安全修复】检查递归深度,防止栈溢出
	if p.depth >= p.maxDepth {
		return nil
	}
	p.depth++
	defer func() { p.depth-- }()

	if p.data[p.pos] != '{' {
		return nil
	}
	p.pos++ // 跳过 {

	obj := make(map[string]*Value)

	p.skipWhitespace()
	if p.pos < len(p.data) && p.data[p.pos] == '}' {
		p.pos++ // 空对象
		return &Value{Type: TypeObject, Object: obj}
	}

	for {
		p.skipWhitespace()

		// 解析 key
		keyVal := p.parseString()
		if keyVal == nil {
			return nil
		}
		key := keyVal.String

		p.skipWhitespace()
		if p.pos >= len(p.data) || p.data[p.pos] != ':' {
			return nil
		}
		p.pos++ // 跳过 :

		// 解析 value
		value := p.parseValue()
		if value == nil {
			return nil
		}

		obj[key] = value

		p.skipWhitespace()
		if p.pos >= len(p.data) {
			return nil
		}

		if p.data[p.pos] == '}' {
			p.pos++
			break
		}

		if p.data[p.pos] == ',' {
			p.pos++
			continue
		}

		return nil
	}

	return &Value{Type: TypeObject, Object: obj}
}

// parseArray 解析数组
func (p *jsonParser) parseArray() *Value {
	// 【安全修复】检查递归深度,防止栈溢出
	if p.depth >= p.maxDepth {
		return nil
	}
	p.depth++
	defer func() { p.depth-- }()

	if p.data[p.pos] != '[' {
		return nil
	}
	p.pos++ // 跳过 [

	var arr []*Value

	p.skipWhitespace()
	if p.pos < len(p.data) && p.data[p.pos] == ']' {
		p.pos++ // 空数组
		return &Value{Type: TypeArray, Array: arr}
	}

	for {
		value := p.parseValue()
		if value == nil {
			return nil
		}

		arr = append(arr, value)

		p.skipWhitespace()
		if p.pos >= len(p.data) {
			return nil
		}

		if p.data[p.pos] == ']' {
			p.pos++
			break
		}

		if p.data[p.pos] == ',' {
			p.pos++
			continue
		}

		return nil
	}

	return &Value{Type: TypeArray, Array: arr}
}

// skipWhitespace 跳过空白字符
func (p *jsonParser) skipWhitespace() {
	for p.pos < len(p.data) {
		ch := p.data[p.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			p.pos++
		} else {
			break
		}
	}
}
