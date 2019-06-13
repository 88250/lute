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

import "regexp"

func Parse(name, text string) (*Tree, error) {
	text = sanitize(text)

	t := &Tree{name: name, text: text, context: &Context{}}
	err := t.parse()

	return t, err
}

var newlinesRegexp = regexp.MustCompile("\r[\n\u0085]?|[\u2424\u2028\u0085]")
var nullRegexp = regexp.MustCompile("\u0000")

func sanitize(text string) (ret string) {
	ret = newlinesRegexp.ReplaceAllString(text, "\n")
	nullRegexp.ReplaceAllString(ret, "\uFFFD") // https://github.github.com/0.29/#insecure-characters

	return
}

// Context use to store common data in parsing.
type Context struct {
	CurNode      Node
	IndentSpaces int
}

// Tree is the representation of the markdown ast.
type Tree struct {
	Root      *Root
	name      string
	text      string
	lex       *lexer
	token     [64]item
	peekCount int
	context   *Context
}

func (t *Tree) HTML() string {
	return t.Root.HTML()
}

// nextToken returns the next token.
func (t *Tree) nextToken() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}

	return t.token[t.peekCount]
}

func (t *Tree) nextNonWhitespace() (spaces, tabs int, tokens []item, firstNonWhitespace item) {
	for {
		token := t.nextToken()
		tokens = append(tokens, token)
		switch token.typ {
		case itemTab:
			tabs++
		case itemSpace:
			spaces++
		case itemNewline:
		default:
			firstNonWhitespace = token
			return
		}
	}
}

func (t *Tree) nonWhitespace(line []item) (spaces, tabs int, tokens []item, firstNonWhitespace item) {
	for i := 0; i < len(line); i++ {
		token := line[0]
		tokens = append(tokens, token)
		switch token.typ {
		case itemTab:
			tabs++
		case itemSpace:
			spaces++
		case itemNewline:
		default:
			firstNonWhitespace = token
			return
		}
	}

	return
}

func (t *Tree) skipWhitespace(line [] item) (tokens []item) {
	for _, token := range line{
		if !token.isWhitespace() {
			tokens = append(tokens, token)
		}
	}

	return
}

func (t *Tree) firstNonSpace(line []item) (index int, token item) {
	for index, token = range line {
		if itemSpace != token.typ {
			return
		}
	}

	return
}

// https://spec.commonmark.org/0.29/#blank-line
func (t *Tree) isBlankLine(line []item) bool {
	for _, token := range line {
		typ := token.typ
		if itemSpace != typ && itemTab != typ && itemNewline != typ {
			return false
		}
	}

	return true
}

func (t *Tree) removeSpaces(line []item) (tokens []item) {
	for _, token := range line {
		if itemSpace != token.typ {
			tokens = append(tokens, token)
		}
	}

	return
}

func (t *Tree) nextLineEnding() (tokens []item) {
	for {
		token := t.nextToken()
		tokens = append(tokens, token)
		if token.isLineEnding() || token.isEOF() {
			return
		}
	}
}

// peek returns but does not consume the next token.
func (t *Tree) peek() item {
	if t.peekCount > 0 {
		return t.token[t.peekCount-1]
	}

	t.peekCount = 1
	t.token[0] = t.lex.nextItem()

	return t.token[0]
}

// backup backs the input stream up one token.
func (t *Tree) backup() {
	t.peekCount++
}

func (t *Tree) backups(tokens []item) {
	i := 0
	l := len(tokens)
	for ; i < l; i++ {
		t.token[l-1-i] = tokens[i] // push back
	}
	t.peekCount = i
}

// Parsing.

// recover is the handler that turns panics into returns from the top level of Parse.
func (t *Tree) recover(err *error) {
	e := recover()
	if e != nil {
		*err = e.(error)
	}
}

func (t *Tree) parse() (err error) {
	defer t.recover(&err)

	t.lex = lex(t.name, t.text)
	t.Root = &Root{NodeType: NodeRoot, Pos: 0}
	t.context.CurNode = t.Root
	t.parseBlocks()
	t.parseInlines()
	t.lex = nil

	return nil
}

func (t *Tree) parseListContent() Node {

	return nil
}

func (t *Tree) parseTableContent() Node {

	return nil
}

func (t *Tree) parseRowContent() Node {

	return nil
}

func (t *Tree) parsePhrasingContent() (ret Node) {
	return
}

func (t *Tree) parseDelete() (ret Node) {
	t.nextToken() // consume open ~~
	token := t.peek()
	ret = &Delete{NodeDelete, token.pos, "", items{}, t, Children{t.parsePhrasingContent()}}
	t.nextToken() // consume close ~~

	return
}

func (t *Tree) parseHTML() (ret Node) {
	return nil
}

func (t *Tree) parseBreak() (ret Node) {
	token := t.nextToken()
	ret = &Break{NodeBreak, token.pos, "", items{}, t}

	return
}

func indentOffset(tokens []item, indentSpaces int, t *Tree) {
	var restoreTokens, nonWhitespaces []item
	compSpaces := 0
	i := 0
	for ; i < len(tokens); i++ {
		typ := tokens[i].typ
		if itemSpace == typ {
			compSpaces++
		} else if itemTab == typ {
			compSpaces += 4
		} else if itemNewline != typ {
			nonWhitespaces = append(nonWhitespaces, tokens[i])
		}
	}

	remains := compSpaces - indentSpaces
	for j := 0; j < remains/4; j++ {
		restoreTokens = append(restoreTokens, item{itemTab, 0, "\t", 0})
	}
	for j := 0; j < remains%4; j++ {
		restoreTokens = append(restoreTokens, item{itemSpace, 0, " ", 0})
	}
	restoreTokens = append(restoreTokens, nonWhitespaces...)
	t.backups(restoreTokens)
}

type stack struct {
	items []interface{}
	count int
}

func (s *stack) push(e interface{}) {
	s.items = append(s.items[:s.count], e)
	s.count++
}

func (s *stack) pop() interface{} {
	if s.count == 0 {
		return nil
	}

	s.count--

	return s.items[s.count]
}

func (s *stack) peek() interface{} {
	if s.count == 0 {
		return nil
	}

	return s.items[s.count-1]
}

func (s *stack) isEmpty() bool {
	return 0 == len(s.items)
}
