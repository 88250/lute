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

// https://spec.commonmark.org/0.29/#punctuation-character
func (i *item) isPunct() bool {
	return unicode.IsPunct(rune(i.val[0]))
}

func (i *item) isASCIIPunct() bool {
	c := i.val[0]

	return (0x21 <= c && 0x2F >= c) || (0x3A <= c && 0x40 >= c) || (0x5B <= c && 0x60 >= c) || (0x7B <= c && 0x7E >= c)
}

// https://spec.commonmark.org/0.29/#line-ending
func (i *item) isLineEnding() bool {
	return itemNewline == i.typ
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
	itemGreater:         ">",
	itemSpace:           "space",
	itemNewline:         "newline",
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
	itemGreater                         // >
	itemSpace                           // space
	itemNewline                         // \n
	itemBackslash                       // \
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
	tCrosshatch      = makeItem(itemCrosshatch, "#")
)

func makeItem(typ itemType, text string) *item {
	return &item{
		typ: typ,
		val: text,
	}
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

func (tokens items) isNewline() bool {
	if 1 != len(tokens) {
		return false
	}

	return itemNewline == tokens[0].typ
}

func (tokens items) removeSpacesTabs() (ret items) {
	for _, token := range tokens {
		if itemSpace != token.typ && itemTab != token.typ {
			ret = append(ret, token)
		}
	}

	return
}
