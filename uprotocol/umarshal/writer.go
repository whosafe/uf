package umarshal

import (
	"strconv"
	"sync"
)

// Writer JSON 写入器
type Writer struct {
	buf []byte
}

// writerPool Writer 对象池
var writerPool = sync.Pool{
	New: func() any {
		return &Writer{
			buf: make([]byte, 0, 512),
		}
	},
}

// AcquireWriter 从对象池获取 Writer
func AcquireWriter() *Writer {
	return writerPool.Get().(*Writer)
}

// ReleaseWriter 释放 Writer 到对象池
func ReleaseWriter(w *Writer) {
	w.buf = w.buf[:0]
	writerPool.Put(w)
}

// Bytes 返回生成的 JSON 字节
func (w *Writer) Bytes() []byte {
	return w.buf
}

// Reset 重置 Writer
func (w *Writer) Reset() {
	w.buf = w.buf[:0]
}

// WriteByte 写入单个字节
func (w *Writer) WriteByte(b byte) {
	w.buf = append(w.buf, b)
}

// WriteBytes 写入字节数组
func (w *Writer) WriteBytes(b []byte) {
	w.buf = append(w.buf, b...)
}

// WriteString 写入字符串 (不转义)
func (w *Writer) WriteRawString(s string) {
	w.buf = append(w.buf, s...)
}

// WriteObjectStart 写入对象开始 {
func (w *Writer) WriteObjectStart() {
	w.buf = append(w.buf, '{')
}

// WriteObjectEnd 写入对象结束 }
func (w *Writer) WriteObjectEnd() {
	w.buf = append(w.buf, '}')
}

// WriteArrayStart 写入数组开始 [
func (w *Writer) WriteArrayStart() {
	w.buf = append(w.buf, '[')
}

// WriteArrayEnd 写入数组结束 ]
func (w *Writer) WriteArrayEnd() {
	w.buf = append(w.buf, ']')
}

// WriteComma 写入逗号
func (w *Writer) WriteComma() {
	w.buf = append(w.buf, ',')
}

// WriteColon 写入冒号
func (w *Writer) WriteColon() {
	w.buf = append(w.buf, ':')
}

// WriteNull 写入 null
func (w *Writer) WriteNull() {
	w.buf = append(w.buf, "null"...)
}

// WriteBool 写入布尔值
func (w *Writer) WriteBool(b bool) {
	if b {
		w.buf = append(w.buf, "true"...)
	} else {
		w.buf = append(w.buf, "false"...)
	}
}

// WriteInt 写入整数
func (w *Writer) WriteInt(n int) {
	w.buf = strconv.AppendInt(w.buf, int64(n), 10)
}

// WriteInt64 写入 int64
func (w *Writer) WriteInt64(n int64) {
	w.buf = strconv.AppendInt(w.buf, n, 10)
}

// WriteUint 写入无符号整数
func (w *Writer) WriteUint(n uint) {
	w.buf = strconv.AppendUint(w.buf, uint64(n), 10)
}

// WriteUint64 写入 uint64
func (w *Writer) WriteUint64(n uint64) {
	w.buf = strconv.AppendUint(w.buf, n, 10)
}

// WriteFloat32 写入 float32
func (w *Writer) WriteFloat32(f float32) {
	w.buf = strconv.AppendFloat(w.buf, float64(f), 'f', -1, 32)
}

// WriteFloat64 写入 float64
func (w *Writer) WriteFloat64(f float64) {
	w.buf = strconv.AppendFloat(w.buf, f, 'f', -1, 64)
}

// WriteObjectField 写入对象字段名 "key":
func (w *Writer) WriteObjectField(key string) {
	w.WriteString(key)
	w.WriteColon()
}
