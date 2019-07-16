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

type ListItem struct {
	*BaseNode

	Checked bool
	Tight   bool
}

func newListItem(t *Tree, token *item) (ret Node) {
	baseNode := &BaseNode{typ: NodeListItem, tokens:items{}}
	ret = &ListItem{
		baseNode,
		false,
		true,
	}
	t.context.CurNode = ret

	return
}

func (t *Tree) parseListItem(line items) (ret Node) {
	ret = newListItem(t, line[0])
	blankLineBetweenBlocks := false
	for {
		n := t.parseBlock(line)
		if nil == n {
			line = t.nextLine()
			continue
		}
		ret.AppendChild(ret, n)

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
		if totalSpaces < t.context.IndentSpaces {
			t.backupLine(line)

			break
		}

		line = t.indentOffset(line, t.context.IndentSpaces)
	}

	if 1 < len(ret.Children()) && blankLineBetweenBlocks {
		ret.(*ListItem).Tight = false
	}

	return
}
