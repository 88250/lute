// Lute - A structural markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lute

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

// Pos represents a byte position in the original input text.
type Pos int

func (p Pos) Position() Pos {
	return p
}

// item represents a token returned from the scanner.
type item struct {
	typ  itemType // the type of this item
	pos  Pos      // the starting position, in bytes, of this item in the input string
	val  string   // the value of this item, aka lexeme
	line int      // the line number at the start of this item
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}

	return fmt.Sprintf("%q", i.val)
}

// https://github.github.com/gfm/#whitespace-character
func (i item) isWhitespace() bool {
	return itemSpace == i.typ || itemTab == i.typ || itemNewline == i.typ // TODO(D): line tabulation (U+000B), form feed (U+000C), or carriage return (U+000D)
}

// itemType identifies the type of lex items.
type itemType int

// Make the types pretty print.
var itemName = map[itemType]string{
	itemEOF:          "EOF",
	itemStr:          "str",
	itemBacktick:     "`",
	itemTilde:        "~",
	itemExclamation:  "!",
	itemCrosshatch:   "#",
	itemAsterisk:     "*",
	itemOpenParen:    "(",
	itemCloseParen:   ")",
	itemHyphen:       "-",
	itemPlus:         "+",
	itemTab:          "tab",
	itemOpenBracket:  "[",
	itemCloseBracket: "]",
	itemDoublequote:  "\"",
	itemSinglequote:  "'",
	itemGreater:      ">",
	itemSpace:        "space",
	itemNewline:      "newline",
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
	itemExclamation                  // !
	itemCrosshatch                   // #
	itemAsterisk                     // *
	itemOpenParen                    // (
	itemCloseParen                   // )
	itemHyphen                       // -
	itemPlus                         // +
	itemTab                          // \t
	itemOpenBracket                  // [
	itemCloseBracket                 // ]
	itemDoublequote                  // "
	itemSinglequote                  // '
	itemGreater                      // >
	itemSpace                        // space
	itemNewline                      // \n
)

const (
	eof = -1
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name    string  // the name of the input; used only for error reports
	input   string  // the string being scanned
	state   stateFn // the next lexing function to enter
	pos     Pos     // current position in the input
	start   Pos     // start position of this item
	width   Pos     // width of last rune read from input
	lastPos Pos     // position of most recent item returned by nextItem
	items   []item  // scanned items
	line    int     // 1+number of newlines seen
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0

		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	if r == '\n' {
		l.line++
	}

	return r
}

// peek returns but does not consume the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()

	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
	// Correct newline count.
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// emit passes an item back to the parser.
func (l *lexer) emit(t itemType) {
	l.items = append(l.items, item{t, l.start, l.input[l.start:l.pos], l.line})
	l.start = l.pos
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	item := l.items[l.lastPos]
	l.lastPos = item.pos

	return item
}

// drain drains the output so the lexing goroutine will exit.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) drain() {
	for range l.items {
	}
}

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: []item{},
		line:  1,
	}

	l.run()

	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexText; nil != l.state; {
		l.state = l.state(l)
	}
}

// State functions.

// lexText scans until a mark character.
func lexText(l *lexer) stateFn {
	r := l.next()
	switch {
	case '`' == r:
		l.emit(itemBacktick)
	case '!' == r:
		l.emit(itemExclamation)
	case '#' == r:
		l.emit(itemCrosshatch)
	case '*' == r:
		l.emit(itemAsterisk)
	case '(' == r:
		l.emit(itemOpenParen)
	case ')' == r:
		l.emit(itemCloseParen)
	case '-' == r:
		l.emit(itemHyphen)
	case '+' == r:
		l.emit(itemPlus)
	case '\t' == r:
		l.emit(itemTab)
	case '[' == r:
		l.emit(itemOpenBracket)
	case ']' == r:
		l.emit(itemCloseBracket)
	case '"' == r:
		l.emit(itemDoublequote)
	case '\'' == r:
		l.emit(itemSinglequote)
	case '>' == r:
		l.emit(itemGreater)
	case ' ' == r:
		l.emit(itemSpace)
	case '\n' == r:
		l.emit(itemNewline)
	case eof == r:
		l.emit(itemEOF)

		return nil
	default:
		return lexStr
	}

	return lexText
}

// lexStr scans a str.
func lexStr(l *lexer) stateFn {
	for {
		r := l.next()
		switch {
		case unicode.IsLetter(r), '/' == r, '"' == r, unicode.IsNumber(r):
		// absorb
		default:
			l.backup()
			l.emit(itemStr)

			return lexText
		}
	}
}
