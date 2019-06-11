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

type lexTest struct {
	name  string
	input string
	items []item
}

func mkItem(typ itemType, text string) item {
	return item{
		typ: typ,
		val: text,
	}
}

var (
	tEOF = mkItem(itemEOF, "")
)

var lexTests = []lexTest{

	{"spec7", "-\t\tfoo\n", []item{mkItem(itemHyphen, "-"), mkItem(itemTab, "\t"), mkItem(itemTab, "\t"), mkItem(itemStr, "foo"), mkItem(itemNewline, "\n"), tEOF}},

	{"crosshatch", "# lute", []item{mkItem(itemCrosshatch, "#"), mkItem(itemSpace, " "), mkItem(itemStr, "lute"), tEOF}},
	{"greater", "> lute", []item{mkItem(itemGreater, ">"), mkItem(itemSpace, " "), mkItem(itemStr, "lute"), tEOF}},
	{"asterisk", "*lute*", []item{mkItem(itemAsterisk, "*"), mkItem(itemStr, "lute"), mkItem(itemAsterisk, "*"), tEOF}},
	{"backtick", "`lute`", []item{mkItem(itemBacktick, "`"), mkItem(itemStr, "lute"), mkItem(itemBacktick, "`"), tEOF}},
	{"tab", "\tlute", []item{mkItem(itemTab, "\t"), mkItem(itemStr, "lute"), tEOF}},
	{"str", "lute", []item{mkItem(itemStr, "lute"), tEOF}},
	{"newline", " \n", []item{mkItem(itemSpace, " "), mkItem(itemNewline, "\n"), tEOF}},
	{"space", " ", []item{mkItem(itemSpace, " "), tEOF}},
	{"empty", "", []item{tEOF}},
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		l := lex(test.name, test.input)
		var items []item
		for _, line := range l.items {
			for _, item := range line {
				items = append(items, item)
			}
		}
		if !equal(items, test.items, false) {
			t.Fatalf("%s:\nexpected\n\t%v\ngot\n\t%v\n", test.name, test.items, items)
		}
	}
}

func equal(i1, i2 []item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}

	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}

	return true
}
