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
	t1 := mkItem(itemBacktick, "`")
	t2 := mkItem(itemStr, "lute")
	t3 := mkItem(itemBacktick, "`")

	s := &stack{}
	s.push(&t1)
	s.push(&t2)
	s.push(&t3)

	if "`" != s.pop().val {
		t.Error("unexpected stack item")
	}

	if "lute" != s.pop().val {
		t.Error("unexpected stack item")
	}

	if "`" != s.peek().val {
		t.Error("unexpected stack item")
	}

	if "`" != s.pop().val {
		t.Error("unexpected stack item")
	}

	s.push(&t1)
	s.push(&t2)

	tokens := s.popMatch(t1)
	if &t1 != tokens[0] || &t2 != tokens[1] {
		t.Error("unexpected stack item")
	}
}
