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

import "regexp"

func Parse(name, text string) (*Tree, error) {
	text = sanitize(text)

	t := &Tree{name: name, text: text}
	err := t.parse()

	return t, err
}

var newlinesRe = regexp.MustCompile("\r[\n\u0085]?|[\u2424\u2028\u0085]")
var nullRe = regexp.MustCompile("\u0000")

func sanitize(text string) (ret string) {
	ret = newlinesRe.ReplaceAllString(text, "\n")
	nullRe.ReplaceAllString(ret, "\uFFFD") // https://github.github.com/gfm/#insecure-characters

	return
}

// Tree is the representation of the markdown ast.
type Tree struct {
	Root      *Root
	CurNode   Node
	name      string // the name of the input; used only for error reports
	text      string
	lex       *lexer
	token     [64]item
	peekCount int
}

func (t *Tree) HTML() string {
	return t.Root.HTML()
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

func (t *Tree) backups(tokens []item) {
	i := 0
	l := len(tokens)
	for ; i < l; i++ {
		t.token[l-1-i] = tokens[i] // push back
	}
	t.peekCount = i
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

func (t *Tree) nextNonWhitespace() (spaces, tabs int, tokens []item) {
	for {
		token := t.next()
		tokens = append(tokens, token)
		switch token.typ {
		case itemTab:
			tabs++
			continue
		case itemSpace:
			spaces++
			continue
		default:
			return
		}
	}
}

// Parsing.

// recover is the handler that turns panics into returns from the top level of Parse.
func (t *Tree) recover(errp *error) {
	e := recover()
	if e != nil {
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

func (t *Tree) parse() (err error) {
	defer t.recover(&err)
	t.startParse(lex(t.name, t.text))
	t.parseContent()
	t.stopParse()

	return nil
}

func (t *Tree) parseContent() {
	t.Root = &Root{NodeType: NodeRoot, Pos: 0}

	for token := t.peek(); itemEOF != token.typ; token = t.peek() {
		var c Node
		switch token.typ {
		case itemAsterisk, itemHyphen:
			c = t.parseList()
		case itemStr, itemCrosshatch, itemGreater, itemTab:
			c = t.parseTopLevelContent()
		case itemSpace:
			spaces, tabs, tokens := t.nextNonWhitespace()
			if 1 > tabs && 4 > spaces {
				last := tokens[len(tokens) - 1]
				if itemAsterisk == last.typ || itemHyphen == last.typ {
					t.backups(tokens)
					c = t.parseList()
				}
				t.Root.append(c)
				continue
			}

			t.backups(tokens)
			c = t.parseIndentCode()
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

func (t *Tree) parseBlockContent() Node {
	switch token := t.peek(); token.typ {
	case itemStr:
		return t.parseParagraph()
	case itemCrosshatch:
		return t.parseHeading()
	case itemGreater:
		return t.parseBlockquote()
	case itemBackquote:
		return t.parseInlineCode()
	case itemTab:
		return t.parseCode()
	case itemAsterisk:
		return t.parseList()
	default:
		return nil
	}
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
	case itemStr, itemTab:
		return t.parseText()
	case itemAsterisk:
		ret = t.parseEmOrStrong()
	case itemBackquote:
		ret = t.parseInlineCode()
	case itemNewline:
		ret = t.parseBreak()
	}

	return
}

func (t *Tree) parseParagraph() Node {
	token := t.peek()

	ret := &Paragraph{NodeParagraph, token.pos, t, Children{}, "<p>", "</p>"}

	for {
		c := t.parsePhrasingContent()
		if nil == c {
			ret.trim()

			break
		}
		ret.append(c)

		token = t.peek()
		if itemNewline == token.typ {
			token = t.next()
			if itemNewline == token.typ {
				t.next()
				break
			}

			t.backup()
		}
	}

	return ret
}

func (t *Tree) parseHeading() (ret Node) {
	token := t.next()
	t.next() // consume spaces

	ret = &Heading{
		NodeHeading, token.pos, t, Children{t.parsePhrasingContent()},
		len(token.val),
	}

	return
}

func (t *Tree) parseThematicBreak() (ret Node) {
	token := t.next()
	ret = &ThematicBreak{NodeThematicBreak, token.pos}

	return
}

func (t *Tree) parseBlockquote() (ret Node) {
	token := t.next() // >

	indentSpaces := 2

	spaces, tabs, tokens := t.nextNonWhitespace()
	totalSpaces := spaces + tabs*4
	if totalSpaces <= indentSpaces {
		t.backup()
		ret = &Blockquote{NodeParagraph, token.pos, Children{t.parseBlockContent()}}

		return
	}

	indentOffset(tokens, indentSpaces, t)

	ret = &Blockquote{NodeParagraph, token.pos, Children{t.parseBlockContent()}}

	return
}

func (t *Tree) parseText() Node {
	token := t.next()

	return &Text{NodeText, token.pos, t, token.val}
}

func (t *Tree) parseEmOrStrong() (ret Node) {
	t.next() // consume open *
	token := t.peek()
	if itemAsterisk == token.typ {
		ret = t.parseStrong()
	} else {
		ret = &Emphasis{NodeEmphasis, token.pos, t, Children{t.parsePhrasingContent()}}
	}
	t.next() // consume close *

	return
}

func (t *Tree) parseStrong() (ret Node) {
	t.next() // consume open *
	token := t.peek()
	ret = &Strong{NodeStrong, token.pos, t, Children{t.parsePhrasingContent()}}
	t.next() // consume close *

	return
}

func (t *Tree) parseDelete() (ret Node) {
	t.next() // consume open ~~
	token := t.peek()
	ret = &Delete{NodeDelete, token.pos, t, Children{t.parsePhrasingContent()}}
	t.next() // consume close ~~

	return
}

func (t *Tree) parseHTML() (ret Node) {
	return nil
}

func (t *Tree) parseBreak() (ret Node) {
	token := t.next()
	ret = &Break{NodeBreak, token.pos, t}

	return
}

func (t *Tree) parseInlineCode() (ret Node) {
	t.next() // consume open `

	code := t.next()
	ret = &InlineCode{NodeInlineCode, code.pos, t, code.val}

	t.next() // consume close `

	return
}

func (t *Tree) parseIndentCode() (ret Node) {
	var code string

Loop:
	for {
		for i := 0; i < 4; {
			token := t.next()
			switch token.typ {
			case itemSpace:
				i++
			case itemTab:
				i += 4
			default:
				break
			}
		}

		token := t.next()
		for ; itemBackquote != token.typ && itemEOF != token.typ; token = t.next() {
			code += token.val
			if itemNewline == token.typ {
				spaces, tabs, tokens := t.nextNonWhitespace()
				if 1 > tabs && 4 > spaces {
					t.backup()
					break Loop
				} else {
					t.backups(tokens)
					continue Loop
				}
			}
		}
	}

	ret = &Code{NodeCode, 0, t, code, "", ""}

	return
}

func (t *Tree) parseCode() (ret Node) {
	t.next() // consume open ```

	token := t.next()
	pos := token.pos
	var code string
	for ; itemBackquote != token.typ && itemEOF != token.typ; token = t.next() {
		code += token.val
		if itemNewline == token.typ {
			if itemBackquote == t.peek().typ {
				break
			}
		}
	}

	ret = &Code{NodeCode, pos, t, code, "", ""}

	if itemEOF == t.peek().typ {
		return
	}

	t.next() // consume close ```

	return
}

func indentOffset(tokens []item, indentSpaces int, t *Tree) {
	var restoreTokens, nonWhitespaces []item
	compSpaces := 0
	i := 0
	for ; i < len(tokens); i++ {
		if itemSpace == tokens[i].typ {
			compSpaces++
		} else if itemTab == tokens[i].typ {
			compSpaces += 4
		} else {
			nonWhitespaces = append(nonWhitespaces, tokens[i])
		}
	}

	remains := compSpaces - indentSpaces
	if 4 <= remains {
		for j := 0; j < remains/4; j++ {
			restoreTokens = append(restoreTokens, item{itemTab, 0, "\t", 0})
		}
		for j := 0; j < remains%4; j++ {
			restoreTokens = append(restoreTokens, item{itemSpace, 0, " ", 0})
		}
		restoreTokens = append(restoreTokens, nonWhitespaces...)
		t.backups(restoreTokens)
	} else {
		t.backup()
	}
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
