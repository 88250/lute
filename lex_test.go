// Lute - A structural markdown engine.
// Copyright (C) 2019, b3log.org
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

	// combined cases
	{"comb1", "### foo\n## foo", []item{mkItem(itemHeader, "###"), mkItem(itemSpace, " "), mkItem(itemStr, "foo"), mkItem(itemSpace, "\n"), mkItem(itemHeader, "##"), mkItem(itemSpace, " "), mkItem(itemStr, "foo"), tEOF}},
	{"comb0", "# h\n\nl**u**te", []item{mkItem(itemHeader, "#"), mkItem(itemSpace, " "), mkItem(itemStr, "h"), mkItem(itemSpace, "\n"), mkItem(itemSpace, "\n"), mkItem(itemStr, "l"), mkItem(itemStrong, "**"), mkItem(itemStr, "u"), mkItem(itemStrong, "**"), mkItem(itemStr, "te"), tEOF}},

	// simple cases
	{"paragraph", "p1\n\np2", []item{mkItem(itemStr, "p1"), mkItem(itemSpace, "\n"), mkItem(itemSpace, "\n"), mkItem(itemStr, "p2"), tEOF}},
	{"img", `![alt](/uri "title")`, []item{mkItem(itemImg, "!"), mkItem(itemOpenLinkText, "["), mkItem(itemStr, "alt"),
		mkItem(itemCloseLinkText, "]"), mkItem(itemOpenLinkHref, "("), mkItem(itemStr, "/uri"), mkItem(itemSpace, " "), mkItem(itemStr, `"title"`), mkItem(itemCloseLinkHref, ")"), tEOF}},
	{"link", `[link](/uri "title")`, []item{mkItem(itemOpenLinkText, "["), mkItem(itemStr, "link"),
		mkItem(itemCloseLinkText, "]"), mkItem(itemOpenLinkHref, "("), mkItem(itemStr, "/uri"), mkItem(itemSpace, " "), mkItem(itemStr, `"title"`), mkItem(itemCloseLinkHref, ")"), tEOF}},
	{"header", "# lute", []item{mkItem(itemHeader, "#"), mkItem(itemSpace, " "), mkItem(itemStr, "lute"), tEOF}},
	{"quote", "> lute", []item{mkItem(itemQuote, ">"), mkItem(itemSpace, " "), mkItem(itemStr, "lute"), tEOF}},
	{"strong", "l**u**te", []item{mkItem(itemStr, "l"), mkItem(itemStrong, "**"), mkItem(itemStr, "u"), mkItem(itemStrong, "**"), mkItem(itemStr, "te"), tEOF}},
	{"em", "l*u*te", []item{mkItem(itemStr, "l"), mkItem(itemEm, "*"), mkItem(itemStr, "u"), mkItem(itemEm, "*"), mkItem(itemStr, "te"), tEOF}},
	{"code", "l`u`te", []item{mkItem(itemStr, "l"), mkItem(itemCode, "`"), mkItem(itemStr, "u"), mkItem(itemCode, "`"), mkItem(itemStr, "te"), tEOF}},
	{"str", "lute", []item{mkItem(itemStr, "lute"), tEOF}},
	{"spaces", " \t\n", []item{mkItem(itemSpace, " "), mkItem(itemSpace, "\t"), mkItem(itemSpace, "\n"), tEOF}},
	{"empty", "", []item{tEOF}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest) (items []item) {
	l := lex(t.name, t.input)
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}

	return
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

func TestLex(t *testing.T) {
	for _, test := range lexTests {
		items := collect(&test)
		if !equal(items, test.items, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", test.name, items, test.items)
		}
	}
}
