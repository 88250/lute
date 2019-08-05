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
//       这样可以减少解析时间和内存分配，可以在很大程度上提升性能

// lexer 描述了词法分析器结构。
type lexer struct {
	items  []items // 所有文本行
	length int     // 总行数
	line   int     // 当前行号
}

// nextLine 返回下一行。
func (lexer *lexer) nextLine() (line items) {
	if lexer.line >= lexer.length {
		return
	}

	line = lexer.items[lexer.line]
	lexer.line++
	return
}

// newLexer 创建一个词法分析器。
func newLexer(input string) *lexer {
	ret := &lexer{items: make([]items, 0, 64)}
	if "" == input {
		ret.items = append(ret.items, items{})
		ret.items[ret.line] = append(ret.items[ret.line], itemEnd)

		return ret
	}

	lineScanner := bufio.NewScanner(strings.NewReader(input))
	var line string
	var itemScanners []*lineLexer
	for lineScanner.Scan() {
		line = lineScanner.Text() + "\n"
		itemScanner := &lineLexer{input: line, items: make([]item, 0, 16)}
		itemScanners = append(itemScanners, itemScanner)
		itemScanner.run()
	}

	for _, itemScanner := range itemScanners {
		ret.items = append(ret.items, itemScanner.items)
	}
	ret.length = len(ret.items)

	return ret
}

// lineLexer 描述了文本行词法分析器。
type lineLexer struct {
	input string // 原文本行
	pos   int    // 当前读取位置
	width int    // 最新一个 token 的宽度（字节数）
	items items  // 分析好的 tokens
}

// run 执行词法分析。
func (s *lineLexer) run() {
	for {
		if r := s.next(); itemEnd == r {
			return
		} else {
			s.items = append(s.items, r)
		}
	}
}

// next 从原文本中返回最新的一个 token，如果读取结束则返回 itemEnd。
func (s *lineLexer) next() (ret item) {
	if s.pos >= len(s.input) {
		s.width = 0
		return itemEnd
	}

	var r rune
	r, s.width = utf8.DecodeRuneInString(s.input[s.pos:])
	s.pos += s.width
	ret = item(r)

	return ret
}
