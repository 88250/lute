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
	"bufio"
	"strings"
	"unicode/utf8"
)

// TODO: 词法分析部分还有性能优化空间：不从原文本中解析出 rune 来，而是通过一些下标来标记并操作原文本 byte 数组。
//       这样做可以节省解析 rune 的时间并大大减少内存分配（rune 和 slice 增长），可以在很大程度上提升性能（30%?）。

// lexer 描述了词法分析器结构。
type lexer struct {
	items   []items // 分析好的所有文本行
	length  int     // 总行数
	lineNum int     // 当前行号
	line    string  // 当前行
	lineLen int     // 当前行长度
	pos     int     // 当前行读取位置
	width   int     // 最新一个 token 的宽度（字节数）
}

// nextLine 返回下一行。
func (l *lexer) nextLine() (line items) {
	if l.lineNum >= l.length {
		return
	}

	line = l.items[l.lineNum]
	l.lineNum++
	return
}

// lex 创建一个词法分析器并对 input 进行词法分析。
func lex(input string) *lexer {
	ret := &lexer{items: make([]items, 0, 64)}

	lineScanner := bufio.NewScanner(strings.NewReader(input))
	for lineScanner.Scan() {
		ret.items = append(ret.items, make([]item, 0, 16))
		ret.line = lineScanner.Text() + "\n"
		ret.lineLen = len(ret.line)
		ret.run()
	}
	ret.length = len(ret.items)
	ret.line = ""
	ret.lineNum = 0
	ret.pos = 0
	ret.lineLen = 0
	ret.width = 0
	return ret
}

// run 执行词法分析。
func (l *lexer) run() {
	for {
		if r := l.next(); itemEnd == r {
			return
		} else {
			l.items[l.lineNum] = append(l.items[l.lineNum], r)
		}
	}
}

// next 从原文本中返回最新的一个 token，如果读取结束则返回 itemEnd。
func (l *lexer) next() (ret item) {
	if l.pos >= l.lineLen {
		l.width = 0
		l.lineNum++
		l.pos = 0

		return itemEnd
	}

	var r rune
	r, l.width = utf8.DecodeRuneInString(l.line[l.pos:])
	l.pos += l.width
	ret = item(r)

	return ret
}
