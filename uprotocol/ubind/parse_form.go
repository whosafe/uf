package ubind

import "strings"

// parseForm 解析 Form 数据 (application/x-www-form-urlencoded)
// 格式: key1=value1&key2=value2
func parseForm(data []byte) *Value {
	obj := make(map[string]*Value)

	if len(data) == 0 {
		return &Value{Type: TypeObject, Object: obj}
	}

	// 按 & 分割
	pairs := splitBytes(data, '&')
	for _, pair := range pairs {
		if len(pair) == 0 {
			continue
		}

		// 按 = 分割
		parts := splitBytes(pair, '=')
		if len(parts) != 2 {
			continue
		}

		key := urlDecode(string(parts[0]))
		value := urlDecode(string(parts[1]))

		obj[key] = &Value{Type: TypeString, String: value}
	}

	return &Value{Type: TypeObject, Object: obj}
}

// splitBytes 按分隔符分割字节数组
func splitBytes(data []byte, sep byte) [][]byte {
	var result [][]byte
	start := 0

	for i := 0; i < len(data); i++ {
		if data[i] == sep {
			result = append(result, data[start:i])
			start = i + 1
		}
	}

	// 添加最后一部分
	if start < len(data) {
		result = append(result, data[start:])
	}

	return result
}

// urlDecode URL 解码
func urlDecode(s string) string {
	var buf strings.Builder

	for i := 0; i < len(s); i++ {
		ch := s[i]

		if ch == '+' {
			buf.WriteByte(' ')
		} else if ch == '%' && i+2 < len(s) {
			// 解码 %XX
			hex := s[i+1 : i+3]
			if b := hexToByte(hex); b >= 0 {
				buf.WriteByte(byte(b))
				i += 2
			} else {
				buf.WriteByte(ch)
			}
		} else {
			buf.WriteByte(ch)
		}
	}

	return buf.String()
}

// hexToByte 将十六进制字符串转换为字节
func hexToByte(s string) int {
	if len(s) != 2 {
		return -1
	}

	high := hexDigit(s[0])
	low := hexDigit(s[1])

	if high < 0 || low < 0 {
		return -1
	}

	return high*16 + low
}

// hexDigit 将十六进制字符转换为数字
func hexDigit(c byte) int {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0')
	case 'a' <= c && c <= 'f':
		return int(c - 'a' + 10)
	case 'A' <= c && c <= 'F':
		return int(c - 'A' + 10)
	default:
		return -1
	}
}
