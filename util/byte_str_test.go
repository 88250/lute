//go:build !javascript

package util

import (
	"bytes"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

func TestStrToBytes(t *testing.T) {
	for _, input := range []string{"", "text", "文本🙂", "a\x00b", strings.Repeat("text", 1024)} {
		got := StrToBytes(input)
		if string(got) != input || len(got) != len(input) || cap(got) != len(input) {
			t.Fatalf("conversion changed content or bounds: %q, len=%d, cap=%d", got, len(got), cap(got))
		}
		if len(input) > 0 && unsafe.SliceData(got) != unsafe.StringData(input) {
			t.Fatal("conversion must share the string's backing storage")
		}
	}
}

func TestStrToBytesLifetime(t *testing.T) {
	values := make([][]byte, 128)
	for i := range values {
		values[i] = stackStringBytes(byte(i))
	}
	runtime.GC()
	for i, value := range values {
		if !bytes.Equal(value, bytes.Repeat([]byte{byte(i)}, 16)) {
			t.Fatalf("value %d corrupted after the source frame returned: %x", i, value)
		}
	}
}

// stackStringBytes 让短字符串离开构造它的栈帧，检查返回切片是否保留底层数据。
//
//go:noinline
func stackStringBytes(seed byte) []byte {
	var data [16]byte
	for i := range data {
		data[i] = seed
	}
	return StrToBytes(string(data[:]))
}
