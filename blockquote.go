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

import "fmt"

type Blockquote struct {
	NodeType
	int
	RawText
	items
	t        *Tree
	Parent   Node
	Subnodes Children
}

func (n *Blockquote) String() string {
	return fmt.Sprintf("%s", n.Subnodes)
}

func (n *Blockquote) HTML() string {
	content := html(n.Subnodes)

	return fmt.Sprintf("<blockquote>\n%s</blockquote>\n", content)
}

func (n *Blockquote) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *Blockquote) Children() Children {
	return n.Subnodes
}

func newBlockquote(t *Tree, token *item) *Blockquote {
	ret := &Blockquote{
		NodeBlockquote, token.pos, "", items{}, t, t.context.CurNode, Children{}}
	t.context.CurNode = ret

	return ret
}

func (t *Tree) parseBlockquote(line items) Node {
	token := line[0]
	indentSpaces := t.context.IndentSpaces + 2

	ret := newBlockquote(t, token)
	line = indentOffset(line[1:], indentSpaces, t)
	for {
		c := t.parseBlock(line)
		if nil == c {
			break
		}

		line = t.nextLine()
		if line.isEOF() {
			break
		}

		//spaces, tabs, tokens, _ := t.nonWhitespace(line)
		//
		//totalSpaces := spaces + tabs*4
		//if totalSpaces < indentSpaces {
		//	t.backups(tokens)
		//	break
		//} else if totalSpaces == indentSpaces {
		//	t.backup()
		//	continue
		//}
		//
		//indentOffset(tokens, indentSpaces, t)
	}

	return ret
}

// https://spec.commonmark.org/0.29/#block-quotes
func (t *Tree) isBlockquote(line items) bool {
	if 2 > len(line) { // at least > and newline
		return false
	}

	_, marker := t.firstNonSpace(line)
	if ">" != marker.val {
		return false
	}

	return true
}
