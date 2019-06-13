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
	CurLine      line
	CurNode      Node
	IndentSpaces int
}

// Tree is the representation of the markdown ast.
type Tree struct {
	Root      *Root
	name      string
	text      string
	lex       *lexer
	peekCount int
	context   *Context
}

func (t *Tree) HTML() string {
	return t.Root.HTML()
}

func (t *Tree) nonWhitespace(line line) (spaces, tabs int, tokens line, firstNonWhitespace item) {
	for i := 0; i < len(line); i++ {
		token := line[i]
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

func (t *Tree) skipWhitespaces(line line) (tokens line) {
	for _, token := range line {
		if !token.isWhitespace() {
			tokens = append(tokens, token)
		}
	}

	return
}

func (t *Tree) skipBlankLines() (count int) {
	for {
		line := t.nextLine()
		if line.isEOF() {
			return
		}
		if !t.isBlankLine(line) {
			t.backupLine(line)
			return
		}
		count++
	}
}

func (t *Tree) firstNonSpace(line line) (index int, token item) {
	for index, token = range line {
		if itemSpace != token.typ {
			return
		}
	}

	return
}

// https://spec.commonmark.org/0.29/#blank-line
func (t *Tree) isBlankLine(line line) bool {
	if line.isEOF() {
		return true
	}

	for _, token := range line {
		typ := token.typ
		if itemSpace != typ && itemTab != typ && itemNewline != typ {
			return false
		}
	}

	return true
}

func (t *Tree) removeSpaces(line line) (tokens line) {
	for _, token := range line {
		if itemSpace != token.typ {
			tokens = append(tokens, token)
		}
	}

	return
}

func indentOffset(tokens line, indentSpaces int, t *Tree) (ret line) {
	var nonWhitespaces line
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
	if 0 >= remains {
		return tokens
	}

	for j := 0; j < remains/4; j++ {
		ret = append(ret, item{itemTab, 0, "\t", 0})
	}
	for j := 0; j < remains%4; j++ {
		ret = append(ret, item{itemSpace, 0, " ", 0})
	}
	ret = append(ret, nonWhitespaces...)

	return
}

type line []item

func (line *line) isEOF() bool {
	return 1 == len(*line) && (*line)[0].isEOF()
}

func (line *line) rawText() (ret RawText) {
	for i := 0; i < len(*line); i++ {
		ret += RawText((*line)[i].val)
	}

	return
}

func (t *Tree) nextLine() (line line) {
	if nil != t.context.CurLine {
		line = t.context.CurLine
		t.context.CurLine = nil

		return
	}

	for {
		token := t.lex.nextItem()
		line = append(line, token)
		if token.isLineEnding() || token.isEOF() {
			return
		}
	}
}

func (t *Tree) backupLine(line line) {
	t.context.CurLine = line
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
