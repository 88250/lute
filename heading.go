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
	Pos
	*Tree
	Children

	Depth int
}

func (n *Heading) String() string {
	return fmt.Sprintf("# %s", n.Children)
}

func (n *Heading) HTML() string {
	content := html(n.Children)

	return fmt.Sprintf("<h%d>%s</h%d>\n", n.Depth, content, n.Depth)
}

func (n *Heading) append(c Node) {
	n.Children = append(n.Children, c)
}

func (t *Tree) parseHeading() Node {
	token := t.next()
	t.next() // consume spaces

	ret := &Heading{
		NodeHeading, token.pos, t, Children{},
		len(token.val),
	}

	c := t.parsePhrasingContent()
	ret.append(c)

	return ret
}
