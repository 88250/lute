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
	"regexp"
)

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
	nullRegexp.ReplaceAllString(ret, "\uFFFD") // https://spec.commonmark.org/0.29/#insecure-characters

	return
}

// Context use to store common data in parsing.
type Context struct {
	// Blocks parsing

	CurLine      items
	CurNode      Node
	IndentSpaces int

	// Inlines parsing

	Delimiters *delimiter
	Pos        int
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

func (t *Tree) Render(renderer *Renderer) (output string, err error) {
	err = renderer.Render(t.Root)
	if nil != err {
		return "", err
	}
	output = renderer.writer.String()

	return
}

func (t *Tree) nonWhitespace(line items) (spaces, tabs int, tokens items, firstNonWhitespace *item) {
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


func (t *Tree) trimLeft(tokens items) (ret items) {
	ret = tokens

	size := len(tokens)
	if 1 > size {
		return
	}

	i := 0
	for ; i < size; i++ {
		if !tokens[i].isWhitespace() {
			break
		}
	}

	ret = tokens[i:]

	return
}

func (t *Tree) trimRight(tokens items) (ret items) {
	ret = tokens

	size := len(tokens)
	if 1 > size {
		return
	}

	i := size - 1
	for ; 0 <= size; i-- {
		if !tokens[i].isWhitespace() {
			break
		}
	}

	ret = tokens[:i+1]

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

func (t *Tree) firstNonSpace(line items) (index int, token *item) {
	for index, token = range line {
		if itemSpace != token.typ {
			return
		}
	}

	return
}

func (t *Tree) accept(line items, itemType itemType) (pos int) {
	for ; pos < len(line); pos++ {
		if itemType != line[pos].typ {
			break
		}
	}

	return
}

func (t *Tree) isBlankLine(line items) bool {
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

func (t *Tree) removeSpacesTabs(line items) (tokens items) {
	for _, token := range line {
		if itemSpace != token.typ && itemTab != token.typ {
			tokens = append(tokens, token)
		}
	}

	return
}

func indentOffset(tokens items, indentSpaces int, t *Tree) (ret items) {
	var nonWhitespaces items
	compSpaces := 0
	i := 0
	for ; i < len(tokens); i++ {
		typ := tokens[i].typ
		if itemSpace == typ {
			compSpaces++
		} else if itemTab == typ {
			compSpaces += 4
		} else {
			nonWhitespaces = append(nonWhitespaces, tokens[i:]...)
			break
		}
	}

	remains := compSpaces - indentSpaces
	if 0 >= remains {
		return tokens
	}

	for j := 0; j < remains/4; j++ {
		ret = append(ret, &item{itemTab, 0, "\t", 0})
	}
	for j := 0; j < remains%4; j++ {
		ret = append(ret, &item{itemSpace, 0, " ", 0})
	}
	ret = append(ret, nonWhitespaces...)

	return
}

func (t *Tree) nextLine() (line items) {
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

func (t *Tree) backupLine(line items) {
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
	t.Root = &Root{&BaseNode{typ: NodeRoot}}
	t.context.CurNode = t.Root
	t.parseBlocks()
	t.parseInlines()
	t.lex = nil

	return nil
}
