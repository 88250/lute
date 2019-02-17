// Lute - A structural markdown engine.
// Copyright (C) 2019, b3log.org
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
	"strings"
	"unicode"
	"unicode/utf8"
)

// Pos represents a byte position in the original input text.
type Pos int

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
	case i.typ == itemError:
		return i.val
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}

	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

// Make the types pretty print.
var itemName = map[itemType]string{
	itemError:         "error",
	itemEOF:           "EOF",
	itemStr:           "str",
	itemHeader:        "#",
	itemQuote:         ">",
	itemListItem:      "-",
	itemCode:          "`",
	itemStrong:        "**",
	itemEm:            "*",
	itemDel:           "~~",
	itemOpenLinkText:  "[",
	itemCloseLinkText: "]",
	itemOpenLinkHref:  "(",
	itemCloseLinkHref: ")",
	itemImg:           "!",
	itemSpace:         "space",
	itemNewline:       "newline",
}

func (i itemType) String() string {
	s := itemName[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}

	return s
}

const (
	itemError         itemType = iota // error occurred; value is text of error
	itemEOF                           // EOF
	itemStr                           // plain text
	itemHeader                        // #
	itemQuote                         // >
	itemListItem                      // -
	itemCode                          // `
	itemStrong                        // **
	itemEm                            // *
	itemDel                           // ~~
	itemOpenLinkText                  // [
	itemCloseLinkText                 // ]
	itemOpenLinkHref                  // (
	itemCloseLinkHref                 // )
	itemImg                           // !
	itemSpace                         // space
	itemNewline                       // newline
)

const (
	eof = -1
)

// stateFn represents the state of the scanner as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name    string    // the name of the input; used only for error reports
	input   string    // the string being scanned
	state   stateFn   // the next lexing function to enter
	pos     Pos       // current position in the input
	start   Pos       // start position of this item
	width   Pos       // width of last rune read from input
	lastPos Pos       // position of most recent item returned by nextItem
	items   chan item // channel of scanned items
	line    int       // 1+number of newlines seen
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
	l.items <- item{t, l.start, l.input[l.start:l.pos], l.line}

	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()

	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.nextItem.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{itemError, l.start, fmt.Sprintf(format, args...), l.line}

	return nil
}

// nextItem returns the next item from the input.
// Called by the parser, not in the lexing goroutine.
func (l *lexer) nextItem() item {
	item := <-l.items
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
		items: make(chan item),
		line:  1,
	}

	go l.run()

	return l
}

// run runs the state machine for the lexer.
func (l *lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}

	close(l.items)
}

// State functions.

// lexText scans until a mark character.
func lexText(l *lexer) stateFn {
	r := l.next()
	switch {
	case '#' == r:
		return lexHeader
	case '>' == r:
		l.emit(itemQuote)

		return lexQuote
	case '*' == r:
		return lexEmorStrong
	case '`' == r:
		return lexCode
	case ' ' == r, '\t' == r:
		l.emit(itemSpace)

		return lexText
	case '\n' == r:
		l.emit(itemNewline)

		return lexText
	case '[' == r:
		l.emit(itemOpenLinkText)

		return lexStr
	case ']' == r:
		l.emit(itemCloseLinkText)
		l.next()
		l.emit(itemOpenLinkHref)

		return lexStr
	case ')' == r:
		l.emit(itemCloseLinkHref)

		return lexText
	case '!' == r:
		if '[' == l.peek() {
			l.emit(itemImg)

			return lexText
		} else {
			l.emit(itemStr)

			return lexText
		}
	case eof == r:
		l.emit(itemEOF)
	default:
		return lexStr
	}

	return nil
}

// lexHeader scans '#'.
func lexHeader(l *lexer) stateFn {
	l.acceptRun("#")
	l.emit(itemHeader)

	r := l.next()
	switch {
	case ' ' == r:
		l.emit(itemSpace)

		return lexStr
	default:
		return l.errorf("# must be followed by a space")
	}
}

// lexQuote scans '>'.
func lexQuote(l *lexer) stateFn {
	r := l.next()
	switch {
	case ' ' == r:
		l.emit(itemSpace)
	}

	return lexText
}

// lexEmOrStrong scans '*' or '**'.
func lexEmorStrong(l *lexer) stateFn {
	r := l.next()
	switch {
	case '*' == r:
		l.emit(itemStrong)

		return lexText
	default:
		l.backup()
		l.emit(itemEm)

		return lexText
	}
}

// lexCode scans '`' or '```'.
func lexCode(l *lexer) stateFn {
	l.acceptRun("`")
	l.emit(itemCode)

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

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n'
}
