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

// +build javascript

package lute

import (
	"unicode/utf8"
)

// item 描述了词法分析的一个 token。
type item struct {
	node     *Node // 所属节点
	termByte byte  // 源码字节值
	ln       int   // 源码行号，从 1 开始
	col      int   // 源码列号，从 1 开始
	offset   int   // 源码偏移位置，从 0 开始
}

// items 定义了 token 数组。
type items []item

// nilItem 返回一个空值 token。
func nilItem() item {
	return item{termByte: 0}
}

// isNilItem 判断 item 是否为空值。
func isNilItem(item item) bool {
	return 0 == item.termByte
}

// newItem 构造一个 token。
func newItem(term byte, ln, col, offset int) item {
	return item{termByte: term, ln: ln, col: col, offset: offset}
}

// term 返回 item 的词素。
func (item item) term() byte {
	return item.termByte
}

// Offset 返回 item 的 offset。
func (item item) Offset() int {
	return item.offset
}

// setTerm 用于设置 tokens 中第 i 个 token 的词素。
func setTerm(tokens *items, i int, term byte) {
	(*tokens)[i].termByte = term
}

// strToItems 将 str 转为 items。
func strToItems(str string) (ret items) {
	ret = make(items, 0, len(str))
	length := len(str)
	for i := 0; i < length; i++ {
		ret = append(ret, item{termByte: str[i]})
	}
	return
}

// itemsToStr 将 items 转为 string。
func itemsToStr(items items) string {
	return string(itemsToBytes(items))
}

// itemsToBytes 将 items 转为 []byte。
func itemsToBytes(items items) (ret []byte) {
	length := len(items)
	for i := 0; i < length; i++ {
		ret = append(ret, items[i].termByte)
	}
	return
}

// bytesToItems 将 bytes 转为 items。
func bytesToItems(bytes []byte) (ret items) {
	ret = make(items, 0, len(bytes))
	length := len(bytes)
	for i := 0; i < length; i++ {
		ret = append(ret, item{termByte: bytes[i]})
	}
	return
}

// bytesToStr 快速转换 []byte 为 string。
func bytesToStr(bytes []byte) string {
	return string(bytes)
}

// strToBytes 快速转换 string 为 []byte。
func strToBytes(str string) []byte {
	return []byte(str)
}

// nextLine 返回下一行。
func (l *lexer) nextLine() (ret items) {
	if l.offset >= l.length {
		return
	}
	ret = make(items, 0, 256)

	l.ln++
	l.col = 0

	var b, nb byte
	i := l.offset
	offset := 0
	for ; i < l.length; i += l.width {
		b = l.input[i]
		l.col++
		if itemNewline == b {
			i++
			ret = append(ret, newItem(b, l.ln, l.col, l.offset+offset))
			break
		} else if itemCarriageReturn == b {
			// 按照规范定义的 line ending (https://spec.commonmark.org/0.29/#line-ending) 处理 \r
			ret = append(ret, newItem(b, l.ln, l.col, l.offset+offset))
			if i < l.length-1 {
				nb = l.input[i+1]
				if itemNewline == nb {
					l.input = append(l.input[:i], l.input[i+1:]...) // 移除 \r，依靠下一个的 \n 切行
					l.length--                                      // 重新计算总长
					ret = ret[:len(ret)-1]
					ret = append(ret, newItem(nb, l.ln, l.col, l.offset+offset))
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
			ret = append(ret, newItem(l.input[i], l.ln, l.col, l.offset))
			ret = append(ret, newItem(l.input[i+1], l.ln, l.col, l.offset+1))
			ret = append(ret, newItem(l.input[i+2], l.ln, l.col, l.offset+2))
			continue
		}

		if utf8.RuneSelf <= b { // 说明占用多个字节
			_, l.width = utf8.DecodeRune(l.input[i:])
			for j := 0; j < l.width; j++ {
				ret = append(ret, newItem(l.input[i+j], l.ln, l.col, l.offset+offset+j))
			}
		} else {
			l.width = 1
			ret = append(ret, newItem(b, l.ln, l.col, l.offset+offset))
		}
		offset+=l.width
	}
	l.offset = i
	return
}
