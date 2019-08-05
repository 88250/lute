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

// 词法分析部分还有性能优化空间：可通过减少

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

// lex creates a new lexer for the input string.
func lex(input string) *lexer {
	ret := &lexer{items: make([]items, 0, 64)}
	if "" == input {
		ret.items = append(ret.items, items{})
		ret.items[ret.line] = append(ret.items[ret.line], itemEOF)

		return ret
	}

	lineScanner := bufio.NewScanner(strings.NewReader(input))
	var line string
	var itemScanners []*scanner
	for lineScanner.Scan() {
		line = lineScanner.Text() + "\n"
		itemScanner := &scanner{input: line, items: make([]item, 0, 16)}
		itemScanners = append(itemScanners, itemScanner)
		itemScanner.run()
	}

	for _, itemScanner := range itemScanners {
		ret.items = append(ret.items, itemScanner.items)
	}
	ret.length = len(ret.items)

	return ret
}

type scanner struct {
	input string // the string being scanned
	pos   int    // current position in the input
	width int    // width of last rune read from input
	items items  // scanned tokens
}

func (s *scanner) run() {
	for {
		r := s.next()
		switch {
		case itemEOF == r:
			return
		default:
			s.items = append(s.items, r)
		}
	}
}

// next returns the next token in the input.
func (s *scanner) next() (ret item) {
	if s.pos >= len(s.input) {
		s.width = 0
		return itemEOF
	}

	var r rune
	r, s.width = utf8.DecodeRuneInString(s.input[s.pos:])
	s.pos += s.width
	ret = item(r)

	return ret
}
