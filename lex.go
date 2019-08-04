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

type lexer struct {
	items  []items // 总行
	length int     // 总行数
	line   int     // 当前行号
}

type scanner struct {
	input string // the string being scanned
	pos   int    // current position in the input
	width int    // width of last rune read from input
	items items  // scanned tokens
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

func (s *scanner) run() {
	for {
		r := s.next()
		switch {
		case itemEOF == r:
			return
		case r.isNewline():
			s.newItem(itemNewline)
		default:
			s.newItem(r)
		}
	}
}

// next returns the next rune in the input.
func (s *scanner) next() item {
	if s.pos >= len(s.input) {
		s.width = 0
		return itemEOF
	}

	r, w := utf8.DecodeRuneInString(s.input[s.pos:])
	s.width = w
	s.pos += s.width

	return item(r)
}

// backup steps back one rune. Can only be called once per call of next.
func (s *scanner) backup() {
	s.pos -= s.width
}

// newItem creates an item with the specified item type.
func (s *scanner) newItem(r item) {
	s.items = append(s.items, r)
}
