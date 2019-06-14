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

type parseTest struct {
	name   string
	input  string
	result string
}

var parseTests = []parseTest{
	// commonmark spec cases
	//{"spec12", "- `one\n- two`\n", "<ul>\n<li>`one</li>\n<li>two`</li>\n</ul>\n"},
	{"spec11", "*\t*\t*\t\n", "<hr />\n"},
	{"spec10", "#\tFoo\n", "<h1>Foo</h1>\n"},
	{"spec9", " - foo\n   - bar\n\t - baz\n", "<ul>\n<li>foo\n<ul>\n<li>bar\n<ul>\n<li>baz</li>\n</ul>\n</li>\n</ul>\n</li>\n</ul>\n"},
	{"spec8", "    foo\n\tbar\n", "<pre><code>foo\nbar\n</code></pre>\n"},
	{"spec7", "-\t\tfoo\n", "<ul>\n<li>\n<pre><code>  foo\n</code></pre>\n</li>\n</ul>\n"},
	{"spec6", ">\t\tfoo\n", "<blockquote>\n<pre><code>  foo\n</code></pre>\n</blockquote>\n"},
	{"spec5", "- foo\n\n\t\tbar\n", "<ul>\n<li>\n<p>foo</p>\n<pre><code>  bar\n</code></pre>\n</li>\n</ul>\n"},
	{"spec4.1", "   - foo\n\n\tbar\n", "<ul>\n<li>foo</li>\n</ul>\n<pre><code>bar\n</code></pre>\n"},
	{"spec4", "  - foo\n\n\tbar\n", "<ul>\n<li>\n<p>foo</p>\n<p>bar</p>\n</li>\n</ul>\n"},
	{"spec3", "    a\ta\n    ὐ\ta\n", "<pre><code>a\ta\nὐ\ta\n</code></pre>\n"},
	{"spec2", "  \tfoo\tbaz\t\tbim\n", "<pre><code>foo\tbaz\t\tbim\n</code></pre>\n"},
	{"spce1", "\tfoo\tbaz\t\tbim\n", "<pre><code>foo\tbaz\t\tbim\n</code></pre>\n"},

	// some simple cases
	{"simple11", "`lu\nte`", "<p><code>lu te</code></p>\n"},
	{"simple10", "lu\n\nte", "<p>lu</p>\n<p>te</p>\n"},
	{"simple9", "* lute", "<ul>\n<li>lute</li>\n</ul>\n"},
	{"simple8", "# lute", "<h1>lute</h1>\n"},
	{"simple7", "> lute", "<blockquote>\n<p>lute</p>\n</blockquote>\n"},
	{"simple6", "l**u**te", "<p>l<strong>u</strong>te</p>\n"},
	{"simple5", "l*u*te", "<p>l<em>u</em>te</p>\n"},
	{"simple4", "    lute\n", "<pre><code>lute\n</code></pre>\n"},
	{"simple3", "\tlute\n", "<pre><code>lute\n</code></pre>\n"},
	{"simple2", "l`u`te", "<p>l<code>u</code>te</p>\n"},
	{"simple1", "lute", "<p>lute</p>\n"},
	{"simple0", "", ""},
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
