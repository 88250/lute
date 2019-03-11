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
	"runtime"
	"strings"
)

// Tree is the representation of the markdown ast.
type Tree struct {
	Root      *Root
	name      string // the name of the input; used only for error reports
	text      string
	lex       *lexer
	token     [3]item
	peekCount int
}

func (t *Tree) HTML() string {
	return t.Root.HTML()
}

func Parse(name, text string) (*Tree, error) {
	t := &Tree{
		name: name,
		text: text,
	}
	_, err := t.Parse(text)

	return t, err
}

// next returns the next token.
func (t *Tree) next() item {
	if t.peekCount > 0 {
		t.peekCount--
	} else {
		t.token[0] = t.lex.nextItem()
	}

	return t.token[t.peekCount]
}

// backup backs the input stream up one token.
func (t *Tree) backup() {
	t.peekCount++
}

// backup2 backs the input stream up two tokens.
// The zeroth token is already there.
func (t *Tree) backup2(t1 item) {
	t.token[1] = t1
	t.peekCount = 2
}

// backup3 backs the input stream up three tokens
// The zeroth token is already there.
func (t *Tree) backup3(t2, t1 item) {
	// Reverse order: we're pushing back.
	t.token[1] = t1
	t.token[2] = t2
	t.peekCount = 3
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

// nextNonSpace returns the next non-space token.
func (t *Tree) nextNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemSpace {
			break
		}
	}

	return token
}

// peekNonSpace returns but does not consume the next non-space token.
func (t *Tree) peekNonSpace() (token item) {
	for {
		token = t.next()
		if token.typ != itemSpace {
			break
		}
	}
	t.backup()

	return token
}

// Parsing.

// ErrorContext returns a textual representation of the location of the node in the input text.
// The receiver is only used when the node does not have a pointer to the tree inside,
// which can occur in old code.
func (t *Tree) ErrorContext(n Node) (location string) {
	pos := int(n.Position())

	text := t.text[:pos]
	byteNum := strings.LastIndex(text, "\n")
	if byteNum == -1 {
		byteNum = pos // On first line.
	} else {
		byteNum++ // After the newline.
		byteNum = pos - byteNum
	}
	lineNum := 1 + strings.Count(text, "\n")

	return fmt.Sprintf("%s:%d:%d", t.name, lineNum, byteNum)
}

// errorf formats the error and terminates processing.
func (t *Tree) errorf(format string, args ...interface{}) {
	t.Root = nil
	format = fmt.Sprintf("tree: %s:%d: %s", t.name, t.token[0].line, format)
	panic(fmt.Errorf(format, args...))
}

// error terminates processing.
func (t *Tree) error(err error) {
	t.errorf("%s", err)
}

// expect consumes the next token and guarantees it has the required type.
func (t *Tree) expect(expected itemType, context string) item {
	token := t.nextNonSpace()
	if token.typ != expected {
		t.unexpected(token, context)
	}

	return token
}

// expectOneOf consumes the next token and guarantees it has one of the required types.
func (t *Tree) expectOneOf(expected1, expected2 itemType, context string) item {
	token := t.nextNonSpace()
	if token.typ != expected1 && token.typ != expected2 {
		t.unexpected(token, context)
	}

	return token
}

// unexpected complains about the token and terminates processing.
func (t *Tree) unexpected(token item, context string) {
	t.errorf("unexpected %s in %s", token, context)
}

// recover is the handler that turns panics into returns from the top level of Parse.
func (t *Tree) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		if t != nil {
			t.lex.drain()
			t.stopParse()
		}
		*errp = e.(error)
	}
}

// startParse initializes the parser, using the lexer.
func (t *Tree) startParse(lex *lexer) {
	t.Root = nil
	t.lex = lex
}

// stopParse terminates parsing.
func (t *Tree) stopParse() {
	t.lex = nil
}

func (t *Tree) Parse(text string) (tree *Tree, err error) {
	defer t.recover(&err)
	t.startParse(lex(t.name, text))
	t.text = text
	t.parseContent()
	t.stopParse()

	return t, nil
}

func (t *Tree) parseContent() {
	t.Root = &Root{Parent{NodeType: NodeRoot, Pos: 0}}

	for token := t.peek(); itemEOF != token.typ; token = t.peek() {
		var c Node
		switch token.typ {
		case itemSpace:
			for {
				token := t.next()
				if itemSpace != token.typ {
					t.backup()

					break
				}
			}
			fallthrough
		case itemStr, itemHeading, itemThematicBreak, itemQuote /* List, Table, HTML */, itemCode, // BlockContent
			itemTab:
			c = t.parseTopLevelContent()
		default:
			c = t.parsePhrasingContent()
		}

		t.Root.append(c)
	}
}

func (t *Tree) parseTopLevelContent() (ret Node) {
	ret = t.parseBlockContent()

	return
}

func (t *Tree) parseBlockContent() (ret Node) {
	switch token := t.peek(); token.typ {
	case itemStr:
		return t.parseParagraph()
	case itemHeading:
		ret = t.parseHeading()
	case itemThematicBreak:
		ret = t.parseThematicBreak()
	case itemQuote:
		ret = t.parseBlockquote()
	case itemInlineCode:
		ret = t.parseInlineCode()
	case itemCode:
		ret = t.parseCode()
	case itemTab:
		ret = t.parseTabCode()
	default:
		t.unexpected(token, "input")
	}

	return
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
	ret = t.parseStaticPhrasingContent()

	return
}

func (t *Tree) parseStaticPhrasingContent() (ret Node) {
	switch token := t.peek(); token.typ {
	case itemStr:
		return t.parseText()
	case itemInlineCode:
		ret = t.parseInlineCode()
	case itemEm:
		ret = t.parseEm()
	case itemStrong:
		ret = t.parseStrong()
	}

	return
}

func (t *Tree) parseParagraph() Node {
	token := t.peek()

	ret := &Paragraph{
		Parent{NodeParagraph, token.pos, nil},
		[]Node{},
	}

	for {
		c := t.parsePhrasingContent()
		if nil == c {
			break
		}

		ret.append(c)
	}

	return ret
}

func (t *Tree) parseHeading() (ret Node) {
	token := t.next()
	t.next() // consume spaces

	ret = &Heading{
		Parent{NodeHeading, token.pos, nil},
		int8(len(token.val)),
		[]Node{t.parsePhrasingContent()},
	}

	return
}

func (t *Tree) parseThematicBreak() (ret Node) {
	token := t.next()
	ret = &ThematicBreak{NodeThematicBreak, token.pos}

	return
}

func (t *Tree) parseBlockquote() (ret Node) {
	token := t.next()
	t.next() // consume spaces

	ret = &Blockquote{
		Parent{NodeParagraph, token.pos, nil},
		[]Node{t.parseBlockContent()},
	}

	return
}

func (t *Tree) parseText() Node {
	token := t.next()

	return &Text{Literal{NodeText, token.pos, token.val}}
}

func (t *Tree) parseEm() (ret Node) {
	t.next() // consume open *
	token := t.peek()
	ret = &Emphasis{
		Parent{NodeEmphasis, token.pos, nil},
		[]Node{t.parsePhrasingContent()},
	}
	t.next() // consume close *

	return
}

func (t *Tree) parseStrong() (ret Node) {
	t.next() // consume open **
	token := t.peek()
	ret = &Strong{
		Parent{NodeStrong, token.pos, nil},
		[]Node{t.parsePhrasingContent()},
	}
	t.next() // consume close **

	return
}

func (t *Tree) parseDelete() (ret Node) {
	t.next() // consume open ~~
	token := t.peek()
	ret = &Delete{
		Parent{NodeDelete, token.pos, nil},
		[]Node{t.parsePhrasingContent()},
	}
	t.next() // consume close ~~

	return
}

func (t *Tree) parseHTML() (ret Node) {
	return nil
}

func (t *Tree) parseBreak() (ret Node) {
	token := t.next()
	ret = &Break{NodeBreak, token.pos}

	return
}

func (t *Tree) parseInlineCode() (ret Node) {
	t.next() // consume open `

	code := t.next()
	ret = &InlineCode{Literal{NodeInlineCode, code.pos, code.val}}

	t.next() // consume close `

	return
}

func (t *Tree) parseCode() (ret Node) {
	t.next() // consume open ```

	code := t.next()
	ret = &Code{Literal{NodeCode, code.pos, code.val}, "", ""}

	t.next() // consume close ```

	return
}

func (t *Tree) parseTabCode() (ret Node) {
	t.next() // consume \t
	var code string
	token := t.next()
	pos := token.pos
	for {
		code += token.val
		if itemNewline == token.typ {
			if itemTab != t.peek().typ {
				break
			}
		}
		token = t.next()
	}

	ret = &Code{Literal{NodeCode, pos, code}, "", ""}

	return
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
