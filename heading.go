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

type Heading struct {
	NodeType
	int
	RawText
	items
	*Tree
	Subnodes Children

	Depth int
}

func (n *Heading) String() string {
	return fmt.Sprintf("# %s", n.Subnodes)
}

func (n *Heading) HTML() string {
	content := html(n.Subnodes)

	return fmt.Sprintf("<h%d>%s</h%d>\n", n.Depth, content, n.Depth)
}

func (n *Heading) Append(c Node) {
	n.Subnodes = append(n.Subnodes, c)
}

func (n *Heading) Children() Children {
	return n.Subnodes
}

func (t *Tree) parseHeading(line items) Node {
	marker := line[0]

	ret := &Heading{
		NodeHeading, marker.pos, "", items{}, t, Children{},
		len(marker.val),
	}

	tokens := t.skipWhitespaces(line[1:])
	for _, token := range tokens {
		if itemEOF == token.typ {
			break
		}
		if itemNewline == token.typ {
			break
		}

		ret.RawText += RawText(token.val)
		ret.items = append(ret.items, token)
	}

	return ret
}

// https://spec.commonmark.org/0.29/#atx-headings
func (t *Tree) isATXHeading(line items) bool {
	if 2 > len(line) { // at least # and newline
		return false
	}

	_, marker := t.firstNonSpace(line)
	// TODO: # 后面还需要空格才能确认是否是列表
	if "#" != marker.val {
		return false
	}

	return true
}
