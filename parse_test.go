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

type parseTest struct {
	name   string
	input  string
	result string
}

var parseTests = []parseTest{
	// commonmark spec cases
	//{"spec7", "-\t\tfoo\n", "<ul>\n<li>\n<pre><code>  foo\n</code></pre>\n</li>\n</ul>\n"},
	//{"spec6", ">\t\tfoo\n", "<blockquote>\n<pre><code>  foo\n</code></pre>\n</blockquote>\n"},
	{"spec5", "- foo\n\n\t\tbar\n", "<ul>\n<li>\n<p>foo</p>\n<pre><code>  bar\n</code></pre>\n</li>\n</ul>\n"},
	{"spec4", "  - foo\n\n\tbar\n", "<ul>\n<li>\n<p>foo</p>\n<p>bar</p>\n</li>\n</ul>\n"},
	{"spec4.1", "   - foo\n\n\tbar\n", "<ul>\n<li>foo</li>\n</ul>\n<pre><code>bar\n</code></pre>\n"},
	{"spec3", "    a\ta\n    ὐ\ta\n", "<pre><code>a\ta\nὐ\ta\n</code></pre>\n"},
	{"spec2", "  \tfoo\tbaz\t\tbim\n", "<pre><code>foo\tbaz\t\tbim\n</code></pre>\n"},

	// some simple cases
	{"paragraph2", "p1\n\np2", "<p>p1</p>\n<p>p2</p>\n"},
	{"paragraph", "p", "<p>p</p>\n"},
	{"list", "* lute", "<ul>\n<li>lute</li>\n</ul>\n"},
	{"heading", "# lute", "<h1>lute</h1>\n"},
	{"quote", "> lute", "<blockquote>\n<p>lute</p>\n</blockquote>\n"},
	{"strong", "l**u**te", "<p>l<strong>u</strong>te</p>\n"},
	{"em", "l*u*te", "<p>l<em>u</em>te</p>\n"},
	{"indent code block2", "    lute\n", "<pre><code>lute\n</code></pre>\n"},
	{"indent code block", "\tlute\n", "<pre><code>lute\n</code></pre>\n"},
	{"inline code", "l`u`te", "<p>l<code>u</code>te</p>\n"},
	{"str", "lute", "<p>lute</p>\n"},
	{"empty", "", ""},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		tree, err := Parse(test.name, test.input)
		if nil != err {
			t.Errorf("%q: unexpected error: %v", test.name, err)
		}

		html := tree.HTML()
		if test.result != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", tree.name, test.result, html, test.input)
		}
	}
}

func TestStack(t *testing.T) {
	e1 := mkItem(itemBackquote, "`")
	e2 := mkItem(itemStr, "lute")
	e3 := mkItem(itemBackquote, "`")

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
