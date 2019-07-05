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
	"fmt"
	"testing"
)

type parseTest struct {
	name   string
	input  string
	result string
}

var parseTests = []parseTest{
	// commonmark spec cases
	//{"spec16", "--\n**\n__\n", "<p>--\n**\n__</p>\n"},
	//{"spec14", "+++\n", "<p>+++</p>\n"},
	//{"spec13", "***\n---\n___\n", "<hr />\n<hr />\n<hr />\n"},
	//{"spec12", "- `one\n- two`\n", "<ul>\n<li>`one</li>\n<li>two`</li>\n</ul>\n"},
	//{"spec11", "*\t*\t*\t\n", "<hr />\n"},
	//{"spec10", "#\tFoo\n", "<h1>Foo</h1>\n"},
	//{"spec9", " - foo\n   - bar\n\t - baz\n", "<ul>\n<li>foo\n<ul>\n<li>bar\n<ul>\n<li>baz</li>\n</ul>\n</li>\n</ul>\n</li>\n</ul>\n"},
	//{"spec8", "    foo\n\tbar\n", "<pre><code>foo\nbar\n</code></pre>\n"},
	//{"spec7", "-\t\tfoo\n", "<ul>\n<li>\n<pre><code>  foo\n</code></pre>\n</li>\n</ul>\n"},
	//{"spec6", ">\t\tfoo\n", "<blockquote>\n<pre><code>  foo\n</code></pre>\n</blockquote>\n"},
	//{"spec5", "- foo\n\n\t\tbar\n", "<ul>\n<li>\n<p>foo</p>\n<pre><code>  bar\n</code></pre>\n</li>\n</ul>\n"},
	//{"spec4.1", "   - foo\n\n\tbar\n", "<ul>\n<li>foo</li>\n</ul>\n<pre><code>bar\n</code></pre>\n"},
	//{"spec4", "  - foo\n\n\tbar\n", "<ul>\n<li>\n<p>foo</p>\n<p>bar</p>\n</li>\n</ul>\n"},
	//{"spec3", "    a\ta\n    ὐ\ta\n", "<pre><code>a\ta\nὐ\ta\n</code></pre>\n"},
	//{"spec2", "  \tfoo\tbaz\t\tbim\n", "<pre><code>foo\tbaz\t\tbim\n</code></pre>\n"},
	//{"spce1", "\tfoo\tbaz\t\tbim\n", "<pre><code>foo\tbaz\t\tbim\n</code></pre>\n"},

	// some simple cases
	//{"simple12", "`l*ut*e", "<p>`l<em>ut</em>e</p>\n"},
	//{"simple11", "`lu\nte`", "<p><code>lu te</code></p>\n"},
	//{"simple10", "lu\n\nte", "<p>lu</p>\n<p>te</p>\n"},
	//{"simple9", "* lute", "<ul>\n<li>lute</li>\n</ul>\n"},
	//{"simple8", "# lute", "<h1>lute</h1>\n"},
	//{"simple7", "> lute", "<blockquote>\n<p>lute</p>\n</blockquote>\n"},
	//{"simple6", "l**ut**e", "<p>l<strong>ut</strong>e</p>\n"},
	//{"simple5", "l*ut*e", "<p>l<em>ut</em>e</p>\n"},
	//{"simple4", "    lute\n", "<pre><code>lute\n</code></pre>\n"},
	//{"simple3", "\tlute\n", "<pre><code>lute\n</code></pre>\n"},
	//{"simple2", "l`ut`e", "<p>l<code>ut</code>e</p>\n"},
	{"simple1", "lute", "<p>lute</p>\n"},
	{"simple0", "", ""},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		fmt.Println("Test [" + test.name + "]")
		tree, err := Parse(test.name, test.input)
		if nil != err {
			t.Errorf("%q: unexpected error: %v", test.name, err)
		}

		html := tree.Render()
		if test.result != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", tree.name, test.result, html, test.input)
		}
	}
}
