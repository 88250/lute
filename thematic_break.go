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

type ThematicBreak struct {
	NodeType
	int
	RawText
	items
}

func (n *ThematicBreak) String() string {
	return fmt.Sprintf("'***'")
}

func (n *ThematicBreak) HTML() string {
	return fmt.Sprintf("<hr />\n")
}

func (n *ThematicBreak) Append(c Node) {}

func (n *ThematicBreak) Children() Children {
	return nil
}

func (t *Tree) parseThematicBreak() (ret Node) {
	token := t.nextToken()
	ret = &ThematicBreak{NodeThematicBreak, token.pos, "", items{}}

	return
}
