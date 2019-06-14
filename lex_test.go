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
	items items
}

func mkItem(typ itemType, text string) *item {
	return &item{
		typ: typ,
		val: text,
	}
}

var (
	tEOF      = mkItem(itemEOF, "")
	tSpace    = mkItem(itemSpace, " ")
	tNewLine  = mkItem(itemNewline, "\n")
	tTab      = mkItem(itemTab, "\t")
	tBacktick = mkItem(itemBacktick, "`")
	tAsterisk = mkItem(itemAsterisk, "*")
)

var lexTests = []lexTest{

	//{"spec7", "-\t\tfoo\n", []item{mkItem(itemHyphen, "-"), tTab, tTab, mkItem(itemStr, "foo"), tNewLine, tEOF}},

	{"simple11", "`lu\nte`", items{tBacktick, mkItem(itemStr, "lu"), tNewLine, mkItem(itemStr, "te"), tBacktick, tEOF}},
	{"simple10", "# lute", items{mkItem(itemCrosshatch, "#"), tSpace, mkItem(itemStr, "lute"), tEOF}},
	{"simple9", "> lute", items{mkItem(itemGreater, ">"), tSpace, mkItem(itemStr, "lute"), tEOF}},
	{"simple8", "*lute*", items{tAsterisk, mkItem(itemStr, "lute"), tAsterisk, tEOF}},
	{"simple7", "`lute`", items{tBacktick, mkItem(itemStr, "lute"), tBacktick, tEOF}},
	{"simple6", "\tlute", items{tTab, mkItem(itemStr, "lute"), tEOF}},
	{"simple5", "lute", items{mkItem(itemStr, "lute"), tEOF}},
	{"simple4", "1\n\n2", items{mkItem(itemStr, "1"), tNewLine, tNewLine, mkItem(itemStr, "2"), tEOF}},
	{"simple3", "\n\n", items{tNewLine, tNewLine, tEOF}},
	{"simple2", " \n", items{tSpace, tNewLine, tEOF}},
	{"simple1", " ", items{tSpace, tEOF}},
	{"simple0", "", items{tEOF}},
}

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		l := lex(test.name, test.input)
		var items items
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

func equal(i1, i2 items, checkPos bool) bool {
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
