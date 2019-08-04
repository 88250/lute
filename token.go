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
	"fmt"
	"strings"
	"unicode"
)

// item 描述了一个读取进来的 Token.
type item struct {
	typ           itemType // Token 类型
	input         *string  // 指向整个输入文本
	valueStartPos int      // 词素（lexeme）起始位置
	valueEndPos   int      // 词素结束位置
}

func (i *item) Value() string {
	val := *i.input
	if "" == val {
		return ""
	}

	return (val)[i.valueStartPos:i.valueEndPos]
}

func (i *item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case len(string(i.Value())) > 10:
		return fmt.Sprintf("%.10q...", i.Value())
	}

	return fmt.Sprintf("%q", i.Value())
}

func (i *item) isWhitespace() bool {
	return itemSpace == i.typ || itemTab == i.typ || itemNewline == i.typ || "\u000A" == i.Value() || "\u000C" == i.Value() || "\u000D" == i.Value()
}

func (i *item) isUnicodeWhitespace() bool {
	length := len(i.Value())
	if 1 != length && 2 != length {
		return false
	}

	r := rune(i.Value()[0])
	if 2 == length {
		r = rune(i.Value()[1])
	}

	return unicode.Is(unicode.Zs, r) || itemTab == i.typ || "\u000D" == i.Value() || itemNewline == i.typ || "\u000C" == i.Value()
}

func (i *item) isSpaceOrTab() bool {
	return i.isSpace() || i.isTab()
}

func (i *item) isSpace() bool {
	return itemSpace == i.typ
}

func (i *item) isTab() bool {
	return itemTab == i.typ
}

func (i *item) isNewline() bool {
	return itemNewline == i.typ
}

func (i *item) isNumInt() bool {
	for _, c := range i.Value() {
		if '0' > c || '9' < c {
			return false
		}
	}

	return true
}

func (i *item) isControl() bool {
	return itemControl == i.typ
}

func (i *item) isPunct() bool {
	return i.isASCIIPunct() || (1 == len(i.Value()) && unicode.IsPunct(rune(i.Value()[0])))
}

func (i *item) isASCIIPunct() bool {
	if 1 != len(i.Value()) {
		return false
	}

	c := i.Value()[0]
	return (0x21 <= c && 0x2F >= c) || (0x3A <= c && 0x40 >= c) || (0x5B <= c && 0x60 >= c) || (0x7B <= c && 0x7E >= c)
}

func (i *item) isASCIILetter() bool {
	for _, c := range i.Value() {
		if !('A' <= c && 'Z' >= c) && !('a' <= c && 'z' >= c) {
			return false
		}
	}

	return true
}

func (i *item) isASCIILetterNumHyphen() bool {
	for _, c := range i.Value() {
		if !('A' <= c && 'Z' >= c) && !('a' <= c && 'z' >= c) && !('0' <= c && '9' >= c) && '-' != c {
			return false
		}
	}

	return true
}

func (i *item) isEOF() bool {
	return itemEOF == i.typ
}

// itemType identifies the type of lex items.
type itemType int

// Make the types pretty print.
var itemName = map[itemType]string{
	itemEOF:          "EOF",
	itemStr:          "str",
	itemBacktick:     "`",
	itemTilde:        "~",
	itemBang:         "!",
	itemCrosshatch:   "#",
	itemAsterisk:     "*",
	itemOpenParen:    "(",
	itemCloseParen:   ")",
	itemHyphen:       "-",
	itemUnderscore:   "_",
	itemPlus:         "+",
	itemTab:          "tab",
	itemOpenBracket:  "[",
	itemCloseBracket: "]",
	itemDoublequote:  "\"",
	itemSinglequote:  "'",
	itemLess:         "<",
	itemGreater:      ">",
	itemSpace:        "space",
	itemNewline:      "newline",
	itemDot:          ".",
	itemColon:        ":",
	itemQuestion:     "?",
	itemAmpersand:    "&",
	itemSemicolon:    ";",
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}

	return s
}

const (
	itemEOF          itemType = iota // EOF
	itemStr                          // plain text
	itemBacktick                     // `
	itemTilde                        // ~
	itemBang                         // !
	itemCrosshatch                   // #
	itemAsterisk                     // *
	itemOpenParen                    // (
	itemCloseParen                   // )
	itemHyphen                       // -
	itemUnderscore                   // _
	itemPlus                         // +
	itemEqual                        // =
	itemTab                          // \t
	itemOpenBracket                  // [
	itemCloseBracket                 // ]
	itemDoublequote                  // "
	itemSinglequote                  // '
	itemLess                         // <
	itemGreater                      // >
	itemSpace                        // space
	itemNewline                      // \n
	itemBackslash                    // \
	itemSlash                        // /
	itemDot                          // .
	itemColon                        // :
	itemQuestion                     // ?
	itemAmpersand                    // &
	itemSemicolon                    // ;
	itemControl
)

var (
	tEOF             = makeItem(itemEOF, "")
	tSpace           = makeItem(itemSpace, " ")
	tNewLine         = makeItem(itemNewline, "\n")
	tTab             = makeItem(itemTab, "\t")
	tBacktick        = makeItem(itemBacktick, "`")
	tAsterisk        = makeItem(itemAsterisk, "*")
	tHypen           = makeItem(itemHyphen, "-")
	tUnderscore      = makeItem(itemUnderscore, "_")
	tPlus            = makeItem(itemPlus, "+")
	tBangOpenBracket = makeItem(itemBang, "!")
	tOpenBracket     = makeItem(itemOpenBracket, "[")
	tCloseBracket    = makeItem(itemCloseBracket, "]")
	tOpenParen       = makeItem(itemOpenParen, "(")
	tCloseParan      = makeItem(itemCloseParen, ")")
	tBackslash       = makeItem(itemBackslash, "\\")
	tSlash           = makeItem(itemSlash, "/")
	tCrosshatch      = makeItem(itemCrosshatch, "#")
	tLess            = makeItem(itemLess, "<")
	tGreater         = makeItem(itemGreater, ">")
	tEqual           = makeItem(itemEqual, "=")
	tDoublequote     = makeItem(itemDoublequote, "\"")
	tDot             = makeItem(itemDot, ".")
	tAmpersand       = makeItem(itemAmpersand, "&")
)

func makeItem(typ itemType, text string) *item {
	return &item{typ, &text, 0, 1}
}

const (
	end = -1
)

type items []*item

func (tokens items) Tokens() items {
	return tokens
}

func (tokens items) rawText() (ret string) {
	b := &strings.Builder{}
	for i := 0; i < len(tokens); i++ {
		b.WriteString(tokens[i].Value())
	}
	ret = b.String()

	return
}

func (tokens items) trimLeftSpace() (spaces int, remains items) {
	size := len(tokens)
	if 1 > size {
		return 0, tokens
	}

	i := 0
	for ; i < size; i++ {
		if tokens[i].isSpace() {
			spaces++
		} else if tokens[i].isTab() {
			spaces += 4
		} else {
			break
		}
	}

	remains = tokens[i:]

	return
}

func (tokens items) trim() (ret items) {
	_, ret = tokens.trimLeft()
	ret = ret.trimRight()

	return
}

func (tokens items) trimLeft() (whitespaces, remains items) {
	size := len(tokens)
	if 1 > size {
		return nil, tokens
	}

	i := 0
	for ; i < size; i++ {
		if !tokens[i].isWhitespace() {
			break
		} else {
			whitespaces = append(whitespaces, tokens[i])
		}
	}

	return whitespaces, tokens[i:]
}

func (tokens items) trimRight() items {
	size := len(tokens)
	if 1 > size {
		return tokens
	}

	i := size - 1
	for ; 0 <= i; i-- {
		if !tokens[i].isWhitespace() && !tokens[i].isEOF() {
			break
		}
	}

	return tokens[:i+1]
}

func (tokens items) firstNonSpace() (index int, token *item) {
	for index, token = range tokens {
		if itemSpace != token.typ {
			return
		}
	}

	return
}

func (tokens items) accept(itemType itemType) (pos int) {
	for ; pos < len(tokens); pos++ {
		if itemType != tokens[pos].typ {
			break
		}
	}

	return
}

func (tokens items) index(itemType itemType) (pos int) {
	for ; pos < len(tokens); pos++ {
		if itemType == tokens[pos].typ {
			return
		}
	}

	return -1
}

func (tokens items) contain(itemTypes ...itemType) bool {
	for _, token := range tokens {
		for _, it := range itemTypes {
			if token.typ == it {
				return true
			}
		}
	}

	return false
}

func (tokens items) containWhitespace() bool {
	for _, token := range tokens {
		if token.isWhitespace() {
			return true
		}
	}

	return false
}

func (tokens items) isBlankLine() bool {
	for _, token := range tokens {
		typ := token.typ
		if itemSpace != typ && itemTab != typ && itemNewline != typ {
			return false
		}
	}

	return true
}

func (tokens items) leftSpaces() (count int) {
	for _, token := range tokens {
		if itemSpace == token.typ {
			count++
		} else if itemTab == token.typ {
			count += 4
		} else {
			break
		}
	}

	return
}

func (tokens items) splitWhitespace() (ret []items) {
	ret = []items{}
	i := 0
	ret = append(ret, items{})
	lastIsWhitespace := false
	for _, token := range tokens {
		if token.isWhitespace() {
			if !lastIsWhitespace {
				i++
				ret = append(ret, items{})
			}
			lastIsWhitespace = true
		} else {
			ret[i] = append(ret[i], token)
			lastIsWhitespace = false
		}
	}

	return
}

func (tokens items) split(itemType itemType) (ret []items) {
	ret = []items{}
	i := 0
	ret = append(ret, items{})
	for j, token := range tokens {
		if itemType == token.typ {
			ret = append(ret, items{})
			ret[i+1] = append(ret[i+1], tokens[j+1:]...)
			return
		} else {
			ret[i] = append(ret[i], token)
		}
	}

	return
}

func (tokens items) isASCIILetterNumHyphen() bool {
	for _, token := range tokens {
		if !token.isASCIILetterNumHyphen() {
			return false
		}
	}

	return true
}

func (tokens items) startWith(itemType itemType) bool {
	if 1 > len(tokens) {
		return false
	}

	return itemType == tokens[0].typ
}

func (tokens items) endWith(itemType itemType) bool {
	length := len(tokens)
	if 1 > length {
		return false
	}

	return itemType == tokens[length-1].typ
}

func (tokens items) isBackslashEscape(pos int) bool {
	if !tokens[pos].isASCIIPunct() {
		return false
	}

	backslashes := 0
	for i := pos - 1; 0 <= i; i-- {
		if itemBackslash != tokens[i].typ {
			break
		}

		backslashes++
	}

	return 0 != backslashes%2
}

func (tokens items) statWhitespace() (newlines, spaces, tabs int) {
	for _, token := range tokens {
		if itemNewline == token.typ {
			newlines++
		} else if itemSpace == token.typ {
			spaces++
		} else if itemTab == token.typ {
			tabs++
		}
	}

	return
}

func (tokens items) spnl() (ret bool, passed, remains items) {
	passed, remains = tokens.trimLeft()
	newlines, _, _ := passed.statWhitespace()
	if 1 < newlines {
		return false, nil, tokens
	}
	ret = true

	return
}

func (tokens items) peek(pos int) *item {
	if pos < len(tokens) {
		return tokens[pos]
	}

	return nil
}
