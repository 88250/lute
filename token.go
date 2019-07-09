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
	"unicode"
)

// item represents a token returned from the scanner.
type item struct {
	typ  itemType // the type of this item
	pos  int      // the starting position, in bytes, of this item in the input string
	val  string   // the value of this item, aka lexeme
	line int      // the line number at the start of this item
}

func (i *item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case len(string(i.val)) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}

	return fmt.Sprintf("%q", i.val)
}

func (i *item) isWhitespace() bool {
	return itemSpace == i.typ || itemTab == i.typ || itemNewline == i.typ // TODO(D): line tabulation (U+000B), form feed (U+000C), or carriage return (U+000D)
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
	for _, c := range i.val {
		if '0' > c || '9' < c {
			return false
		}
	}

	return true
}

func (i *item) isPunct() bool {
	return unicode.IsPunct(rune(i.val[0]))
}

func (i *item) isASCIIPunct() bool {
	c := i.val[0]

	return (0x21 <= c && 0x2F >= c) || (0x3A <= c && 0x40 >= c) || (0x5B <= c && 0x60 >= c) || (0x7B <= c && 0x7E >= c)
}

func (i *item) isASCIILetter() bool {
	for _, c := range i.val {
		if !('A' <= c && 'Z' >= c) && !('a' <= c && 'z' >= c) {
			return false
		}
	}

	return true
}

func (i *item) isASCIILetterNumHyphen() bool {
	for _, c := range i.val {
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
	itemEOF:             "EOF",
	itemStr:             "str",
	itemBacktick:        "`",
	itemTilde:           "~",
	itemBangOpenBracket: "![",
	itemCrosshatch:      "#",
	itemAsterisk:        "*",
	itemOpenParen:       "(",
	itemCloseParen:      ")",
	itemHyphen:          "-",
	itemUnderscore:      "_",
	itemPlus:            "+",
	itemTab:             "tab",
	itemOpenBracket:     "[",
	itemCloseBracket:    "]",
	itemDoublequote:     "\"",
	itemSinglequote:     "'",
	itemLess:            "<",
	itemGreater:         ">",
	itemSpace:           "space",
	itemNewline:         "newline",
	itemDot:             ".",
	itemColon:           ":",
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}

	return s
}

const (
	itemEOF             itemType = iota // EOF
	itemStr                             // plain text
	itemBacktick                        // `
	itemTilde                           // ~
	itemBangOpenBracket                 // ![
	itemCrosshatch                      // #
	itemAsterisk                        // *
	itemOpenParen                       // (
	itemCloseParen                      // )
	itemHyphen                          // -
	itemUnderscore                      // _
	itemPlus                            // +
	itemEqual                           // =
	itemTab                             // \t
	itemOpenBracket                     // [
	itemCloseBracket                    // ]
	itemDoublequote                     // "
	itemSinglequote                     // '
	itemLess                            // <
	itemGreater                         // >
	itemSpace                           // space
	itemNewline                         // \n
	itemBackslash                       // \
	itemSlash                           // /
	itemDot                             // .
	itemColon                           // :
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
	tBangOpenBracket = makeItem(itemBangOpenBracket, "![")
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
)

func makeItem(typ itemType, text string) *item {
	return &item{
		typ: typ,
		val: text,
	}
}

func makeItems(typ itemType, text string, num int) (ret items) {
	for i := 0; i < num; i++ {
		ret = append(ret, makeItem(typ, text))
	}

	return
}

const (
	end = -1
)

type items []*item

func (tokens items) Tokens() items {
	return tokens
}

func (tokens items) isEOF() bool {
	return 1 == len(tokens) && (tokens)[0].isEOF()
}

func (tokens items) rawText() (ret string) {
	for i := 0; i < len(tokens); i++ {
		ret += (tokens)[i].val
	}

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
	ret = tokens.trimLeft()
	ret = ret.trimRight()

	return
}

func (tokens items) trimLeft() items {
	size := len(tokens)
	if 1 > size {
		return tokens
	}

	i := 0
	for ; i < size; i++ {
		if !tokens[i].isWhitespace() {
			break
		}
	}

	return tokens[i:]
}

func (tokens items) trimRight() items {
	size := len(tokens)
	if 1 > size {
		return tokens
	}

	i := size - 1
	for ; 0 <= size; i-- {
		if !tokens[i].isWhitespace() {
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

func (tokens items) allAre(itemType itemType) bool {
	for _, token := range tokens {
		if token.typ != itemType {
			return false
		}
	}

	return true
}

func (tokens items) isBlankLine() bool {
	if tokens.isEOF() {
		return true
	}

	for _, token := range tokens {
		typ := token.typ
		if itemSpace != typ && itemTab != typ && itemNewline != typ {
			return false
		}
	}

	return true
}

func (tokens items) removeSpacesTabs() (ret items) {
	for _, token := range tokens {
		if itemSpace != token.typ && itemTab != token.typ {
			ret = append(ret, token)
		}
	}

	return
}

func (tokens items) whitespaceCountLeft() (count int) {
	for _, token := range tokens {
		if !token.isWhitespace() {
			break
		} else {
			count++
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
	for j, token := range tokens {
		if itemType == token.typ {
			ret[i+1] = append(ret[i+1], tokens[j:]...)
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
