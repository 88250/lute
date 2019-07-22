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
	nullRegexp.ReplaceAllString(ret, "\uFFFD")

	return
}

type BlockContainer struct {
	nodes []Node
}

func (bc *BlockContainer) push(node Node) {
	bc.nodes = append(bc.nodes, node)
}

func (bc *BlockContainer) pop() (ret Node) {
	length := len(bc.nodes)
	if 1 > length {
		return nil
	}

	ret = bc.nodes[length-1]
	bc.nodes = bc.nodes[:length-1]

	return
}

func (bc *BlockContainer) peek() Node {
	length := len(bc.nodes)
	if 1 > length {
		return nil
	}

	return bc.nodes[length-1]
}

// Context use to store common data in parsing.
type Context struct {
	LinkRefDef map[string]*Link
	CurLines   []items

	// Inlines parsing

	Pos               int
	Delimiters        *delimiter
	Brackets          *delimiter
	previousDelimiter *delimiter
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

func (t *Tree) nonSpaceTab(tokens items) (spaces, tabs int, remains items) {
	i := 0
Loop:
	for ; i < len(tokens); i++ {
		token := tokens[i]
		switch token.typ {
		case itemTab:
			tabs++
		case itemSpace:
			spaces++
		default:
			break Loop
		}
	}

	remains = tokens[i:]

	return
}

func (t *Tree) skipBlankLines() (blankLines []items) {
	for {
		line := t.nextLine()
		if line.isEOF() {
			return
		}

		//tokens := t.removeStartBlockquoteMarker(line, t.context.BlockquoteLevel)
		if !line.isBlankLine() {
			t.backupLine(line)
			return
		}

		blankLines = append(blankLines, line)
	}
}

func (t *Tree) indentOffset(tokens items, indentSpaces int) (ret items) {
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
	if 0 > remains {
		return nonWhitespaces
	}

	for j := 0; j < remains; j++ {
		ret = append(ret, &item{itemSpace, 0, " ", 0})
	}
	ret = append(ret, nonWhitespaces...)

	return
}

func (t *Tree) nextLine() (line items) {
	length := len(t.context.CurLines)
	if 0 < length {
		line = t.context.CurLines[0]
		t.context.CurLines = t.context.CurLines[1:]

		return
	}

	for {
		token := t.lex.nextItem()
		line = append(line, token)
		if token.isNewline() || token.isEOF() {
			return
		}
	}
}

func (t *Tree) backupLine(line items) {
	if 0 < len(t.context.CurLines) {
		t.context.CurLines = append([]items{line}, t.context.CurLines...)
	} else {
		t.context.CurLines = append(t.context.CurLines, line)
	}
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
	t.context.LinkRefDef = map[string]*Link{}
	t.parseBlocks()
	t.parseInlines()
	t.lex = nil

	return nil
}
