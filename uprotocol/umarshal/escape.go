package umarshal

// 转义表: 标记哪些字符需要转义
var escapeTable = [256]bool{
	'"':  true, // 双引号
	'\\': true, // 反斜杠
	'\b': true, // 退格
	'\f': true, // 换页
	'\n': true, // 换行
	'\r': true, // 回车
	'\t': true, // 制表符
}

// 转义字符映射
var escapeChar = [256]byte{
	'"':  '"',
	'\\': '\\',
	'\b': 'b',
	'\f': 'f',
	'\n': 'n',
	'\r': 'r',
	'\t': 't',
}

// WriteString 写入字符串 (带转义)
func (w *Writer) WriteString(s string) {
	w.WriteByte('"')

	// 快速路径: 检查是否需要转义
	needEscape := false
	for i := 0; i < len(s); i++ {
		if s[i] < 0x20 || escapeTable[s[i]] {
			needEscape = true
			break
		}
	}

	if !needEscape {
		// 无需转义,直接写入
		w.WriteRawString(s)
	} else {
		// 需要转义,逐字符处理
		w.writeEscapedString(s)
	}

	w.WriteByte('"')
}

// writeEscapedString 写入转义字符串
func (w *Writer) writeEscapedString(s string) {
	start := 0
	for i := 0; i < len(s); i++ {
		c := s[i]

		// 控制字符 (0x00-0x1F) 需要 \uXXXX 转义
		if c < 0x20 {
			// 写入之前的正常字符
			if i > start {
				w.WriteRawString(s[start:i])
			}

			// 特殊控制字符
			if escapeTable[c] {
				w.WriteByte('\\')
				w.WriteByte(escapeChar[c])
			} else {
				// 其他控制字符用 \uXXXX
				w.WriteRawString("\\u00")
				w.WriteByte(hexChar(c >> 4))
				w.WriteByte(hexChar(c & 0x0F))
			}
			start = i + 1
			continue
		}

		// 特殊字符转义
		if escapeTable[c] {
			// 写入之前的正常字符
			if i > start {
				w.WriteRawString(s[start:i])
			}

			w.WriteByte('\\')
			w.WriteByte(escapeChar[c])
			start = i + 1
		}
	}

	// 写入剩余的正常字符
	if start < len(s) {
		w.WriteRawString(s[start:])
	}
}

// hexChar 返回十六进制字符
func hexChar(n byte) byte {
	if n < 10 {
		return '0' + n
	}
	return 'a' + (n - 10)
}
