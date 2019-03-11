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

type htmlTest struct {
	name   string
	input  string
	result string // what the user would see in an error message.
}

var htmlTests = []htmlTest{
	{"heading6", "###### lute", `<h6>lute</h6>`},
	{"heading5", "##### lute", `<h5>lute</h5>`},
	{"heading4", "#### lute", `<h4>lute</h4>`},
	{"heading3", "### lute", `<h3>lute</h3>`},
	{"heading2", "## lute", `<h2>lute</h2>`},
	{"heading1", "# lute", `<h1>lute</h1>`},
	{"quote", "> lute", `<blockquote><p>lute</p></blockquote>`},
	{"strong", "l**u**te", `<p>l<strong>u</strong>te</p>`},
	{"em", "l*u*te", `<p>l<em>u</em>te</p>`},
	{"inlineCode", "l`u`te", `<p>l<code>u</code>te</p>`},
	{"str", "lute", `<p>lute</p>`},
	{"empty", "", ``},
}

func TestHTML(t *testing.T) {
	for _, test := range htmlTests {
		tree, err := Parse(test.name, test.input)
		if nil != err {
			t.Fatalf("%q: unexpected error: %v", test.name, err)
		}

		html := tree.HTML()
		if test.result != html {
			t.Fatalf("%s: expected is %s, but actual is %s\n", tree.name, test.result, html)
		}
	}
}
