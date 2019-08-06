// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import (
	"unsafe"
)

// TODO: 词法分析部分还有性能优化空间：不从原文本中解析出 rune 来，而是通过一些下标来标记并操作原文本 byte 数组。
//       这样做可以节省解析 rune 的时间并大大减少内存分配（rune 和 slice 增长），可以在很大程度上提升性能（30%?）。

// lexer 描述了词法分析器结构。
type lexer struct {
	input   items  // 输入的文本字符数组
	length  int    // 输入的文本字符数组的长度
	offset  int    // 当前读取位置
	line    []byte // 当前行
	lineNum int    // 当前行号
	lineLen int    // 当前行长度

	width int // TODO 最新一个 token 的宽度（字节数）
}

// nextLine 返回下一行。
func (l *lexer) nextLine() (line items) {
	if l.offset >= l.length {
		return
	}

	var b item
	i := l.offset
	for ; i < l.length; i++ {
		b = l.input[i]
		if '\n' == b {
			i++
			break
		}
	}
	line = l.input[l.offset:i]
	l.offset = i
	return
}

// lex 创建一个词法分析器。
func lex(input string) (ret *lexer) {
	ret = &lexer{}
	ret.load(input)
	ret.length = len(ret.input)

	return
}

func (l *lexer) load(str string) {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	l.input = *(*items)(unsafe.Pointer(&h))
}
