// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

// +build !javascript

package lute

import (
	"unicode/utf8"
	"unsafe"
)

// item 描述了词法分析的一个 token。
type item byte

// items 定义了 token 数组。
type items []item

// nilItem 返回一个空值 token。
func nilItem() item {
	return item(0)
}

// isNilItem 判断 item 是否为空值。
func isNilItem(item item) bool {
	return 0 == item
}

// newItem 构造一个 token。
func newItem(term byte, ln, col, offset int) item {
	return item(term)
}

// term 返回 item 的词素。
func term(item item) byte {
	return byte(item)
}

// Offset 返回 item 的 offset。
func (item item) Offset() int {
	return 0
}

// TODO: 作为 item 的方法

// setTerm 用于设置 tokens 中第 i 个 token 的词素。
func setTerm(tokens *items, i int, term byte) {
	(*tokens)[i] = item(term)
}

// strToItems 将 str 转为 items。
func strToItems(str string) items {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*items)(unsafe.Pointer(&h))
}

// itemsToStr 将 items 转为 string。
func itemsToStr(items items) string {
	return *(*string)(unsafe.Pointer(&items))
}

// itemsToBytes 将 items 转为 []byte。
func itemsToBytes(items items) []byte {
	return *(*[]byte)(unsafe.Pointer(&items))
}

// bytesToItems 将 bytes 转为 items。
func bytesToItems(bytes []byte) items {
	return *(*items)(unsafe.Pointer(&bytes))
}

// bytesToStr 快速转换 []byte 为 string。
func bytesToStr(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// strToBytes 快速转换 string 为 []byte。
func strToBytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// nextLine 返回下一行。
func (l *lexer) nextLine() (ret items) {
	if l.offset >= l.length {
		return
	}

	l.ln++
	l.col = 0

	var b, nb byte
	i := l.offset
	for ; i < l.length; i += l.width {
		b = l.input[i]
		l.col++
		if itemNewline == b {
			i++
			break
		} else if itemCarriageReturn == b {
			// 按照规范定义的 line ending (https://spec.commonmark.org/0.29/#line-ending) 处理 \r
			if i < l.length-1 {
				nb = l.input[i+1]
				if itemNewline == nb {
					l.input = append(l.input[:i], l.input[i+1:]...) // 移除 \r，依靠下一个的 \n 切行
					l.length--                                      // 重新计算总长
				}
			}
			i++
			break
		} else if '\u0000' == b {
			// 将 \u0000 替换为 \uFFFD https://spec.commonmark.org/0.29/#insecure-characters
			l.input = append(l.input, 0, 0)
			copy(l.input[i+2:], l.input[i:])
			// \uFFFD 的 UTF-8 编码为 \xEF\xBF\xBD 共三个字节
			l.input[i] = '\xEF'
			l.input[i+1] = '\xBF'
			l.input[i+2] = '\xBD'
			l.length += 2 // 重新计算总长
			l.width = 3
			continue
		}

		if utf8.RuneSelf <= b { // 说明占用多个字节
			_, l.width = utf8.DecodeRune(l.input[i:])
		} else {
			l.width = 1
		}
	}
	ret = bytesToItems(l.input[l.offset:i])
	l.offset = i
	return
}
