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
	"strings"
)

type ListItem struct {
	*Node
	int
	items
	t *Tree

	Checked bool
	Tight   bool

	Spaces int
}

func (n *ListItem) String() string {
	return fmt.Sprintf("%s", n.Children())
}

func (n *ListItem) HTML() string {
	var content string
	for _, c := range n.Children() {
		content += c.HTML()
	}

	if strings.Contains(content, "<ul>") {
		return fmt.Sprintf("<li>%s</li>\n", content)
	}

	if 1 < len(n.Children()) || strings.Contains(content, "<pre><code") {
		return fmt.Sprintf("<li>\n%s</li>\n", content)
	}

	return fmt.Sprintf("<li>%s</li>\n", content)
}

func newListItem(indentSpaces int, t *Tree, token *item) (ret *Node, listItem *ListItem) {
	ret = &Node{NodeType: NodeListItem}
	listItem = &ListItem{
		ret, token.pos, items{}, t,
		false,
		true,
		indentSpaces,
	}
	t.context.CurNode = ret

	return
}

func (t *Tree) parseListItem(line items) (ret *Node, listItem *ListItem) {
	indentSpaces := t.context.IndentSpaces
	ret, listItem = newListItem(indentSpaces, t, line[0])
	blankLineBetweenBlocks := false
	for {
		c := t.parseBlock(line)
		if nil == c {
			continue
		}

		blankLines := t.skipBlankLines()
		if 1 <= blankLines && !blankLineBetweenBlocks {
			blankLineBetweenBlocks = true
		}

		line = t.nextLine()
		if line.isEOF() {
			break
		}

		spaces, tabs, _, _ := t.nonWhitespace(line)
		totalSpaces := spaces + tabs*4
		if totalSpaces < indentSpaces {
			t.backupLine(line)

			break
		}

		line = indentOffset(line, indentSpaces, t)
	}

	if 1 < len(ret.Children()) && blankLineBetweenBlocks {
		listItem.Tight = false
	}

	return
}
