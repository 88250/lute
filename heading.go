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

func (t *Tree) parseHeading() Node {
	token := t.next()

	ret := &Heading{
		NodeHeading, token.pos, "", items{}, t, Children{},
		len(token.val),
	}

	t.nextNonWhitespace()
	t.backup()
	for {
		token = t.next()
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
