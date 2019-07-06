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
	"testing"
)

type lexTest struct {
	name  string
	input string
	items items
}

var lexTests = []lexTest{

	{"spec14", "+++\n", items{tPlus, tPlus, tPlus, tNewLine, tEOF}},
	{"spec13", "***\n---\n___\n", items{tAsterisk, tAsterisk, tAsterisk, tNewLine, tHypen, tHypen, tHypen, tNewLine, tUnderscore, tUnderscore, tUnderscore, tNewLine, tEOF}},
	{"spec7", "-\t\tfoo\n", items{makeItem(itemHyphen, "-"), tTab, tTab, makeItem(itemStr, "foo"), tNewLine, tEOF}},

	{"simple13", "![lute]()", items{tBangOpenBracket, makeItem(itemStr, "lute"), tCloseBracket, tOpenParen, tCloseParan, tEOF}},
	{"simple12", "[lute]()", items{tOpenBracket, makeItem(itemStr, "lute"), tCloseBracket, tOpenParen, tCloseParan, tEOF}},
	{"simple11", "`lu\nte`", items{tBacktick, makeItem(itemStr, "lu"), tNewLine, makeItem(itemStr, "te"), tBacktick, tEOF}},
	{"simple10", "# lute", items{makeItem(itemCrosshatch, "#"), tSpace, makeItem(itemStr, "lute"), tEOF}},
	{"simple9", "> lute", items{makeItem(itemGreater, ">"), tSpace, makeItem(itemStr, "lute"), tEOF}},
	{"simple8", "*lute*", items{tAsterisk, makeItem(itemStr, "lute"), tAsterisk, tEOF}},
	{"simple7", "`lute`", items{tBacktick, makeItem(itemStr, "lute"), tBacktick, tEOF}},
	{"simple6", "\tlute", items{tTab, makeItem(itemStr, "lute"), tEOF}},
	{"simple5", "lute", items{makeItem(itemStr, "lute"), tEOF}},
	{"simple4", "1\n\n2", items{makeItem(itemStr, "1"), tNewLine, tNewLine, makeItem(itemStr, "2"), tEOF}},
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
