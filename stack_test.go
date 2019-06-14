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

import (
	"testing"
)

func TestStack(t *testing.T) {
	e1 := mkItem(itemBacktick, "`")
	e2 := mkItem(itemStr, "lute")
	e3 := mkItem(itemBacktick, "`")

	s := &stack{}
	s.push(&e1)
	s.push(&e2)
	s.push(&e3)

	if "`" != s.pop().(*item).val {
		t.Log("unexpected stack item")
	}

	if "lute" != s.pop().(*item).val {
		t.Log("unexpected stack item")
	}

	if "`" != s.peek().(*item).val {
		t.Log("unexpected stack item")
	}

	if "`" != s.pop().(*item).val {
		t.Log("unexpected stack item")
	}
}
