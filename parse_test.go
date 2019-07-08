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
	{"spec96", "`````\n\n```\naaa\n", "<pre><code>\n```\naaa\n</code></pre>\n"},
	{"spec96", "```\n", "<pre><code></code></pre>\n"},
	{"spec91", "``\nfoo\n``\n", "<p><code>foo</code></p>\n"},
	{"spec90", "~~~\n<\n >\n~~~\n", "<pre><code>&lt;\n &gt;\n</code></pre>\n"},
	{"spec87", "\n    \n    foo\n    \n\n", "<pre><code>foo\n</code></pre>\n"},
	{"spec81", "    chunk1\n\n    chunk2\n  \n \n \n    chunk3\n", "<pre><code>chunk1\n\nchunk2\n\n\n\nchunk3\n</code></pre>\n"},
	{"spec79", "1.  foo\n\n    - bar\n", "<ol>\n<li>\n<p>foo</p>\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ol>\n"},
	{"spec62", "> Foo\n---\n", "<blockquote>\n<p>Foo</p>\n</blockquote>\n<hr />\n"},
	{"spec61", "`Foo\n----\n`\n\n<a title=\"a lot\n---\nof dashes\"/>\n", "<h2>`Foo</h2>\n<p>`</p>\n<h2>&lt;a title=&quot;a lot</h2>\n<p>of dashes&quot;/&gt;</p>\n"},
	{"spec60", "Foo\\\n----\n", "<h2>Foo\\</h2>\n"},
	{"spec56", "Foo\n   ----      \n", "<h2>Foo</h2>\n"},
	{"spec55", "    Foo\n    ---\n\n    Foo\n---\n", "<pre><code>Foo\n---\n\nFoo\n</code></pre>\n<hr />\n"},
	{"spec51", "Foo *bar\nbaz*\n====\n", "<h1>Foo <em>bar\nbaz</em></h1>\n"},
	{"spec41", "## \n#\n### ###\n", "<h2></h2>\n<h1></h1>\n<h3></h3>\n"},
	{"spec41", "## foo ##\n  ###   bar    ###\n", "<h2>foo</h2>\n<h3>bar</h3>\n"},
	{"spec38", " ### foo\n  ## foo\n   # foo\n", "<h3>foo</h3>\n<h2>foo</h2>\n<h1>foo</h1>\n"},
	{"spec37", "#                  foo                     \n", "<h1>foo</h1>\n"},
	{"spec36", "# foo *bar* \\*baz\\*\n", "<h1>foo <em>bar</em> *baz*</h1>\n"},
	{"spec35", "\\## foo\n", "<p>## foo</p>\n"},
	{"spec34", "#5 bolt\n\n#hashtag\n", "<p>#5 bolt</p>\n<p>#hashtag</p>\n"},
	{"spec32", "# foo\n## foo\n### foo\n#### foo\n##### foo\n###### foo\n", "<h1>foo</h1>\n<h2>foo</h2>\n<h3>foo</h3>\n<h4>foo</h4>\n<h5>foo</h5>\n<h6>foo</h6>\n"},
	{"spec30", "* Foo\n* * *\n* Bar\n", "<ul>\n<li>Foo</li>\n</ul>\n<hr />\n<ul>\n<li>Bar</li>\n</ul>\n"},
	{"spec29", "Foo\n---\nbar\n", "<h2>Foo</h2>\n<p>bar</p>\n"},
	{"spec26", " *-*\n", "<p><em>-</em></p>\n"},
	{"spec18", "Foo\n    ***\n", "<p>Foo\n***</p>\n"},
	{"spec18", "    ***\n", "<pre><code>***\n</code></pre>\n"},
	{"spec16", "--\n**\n__\n", "<p>--\n**\n__</p>\n"},
	{"spec14", "+++\n", "<p>+++</p>\n"},
	{"spec13", "***\n---\n___\n", "<hr />\n<hr />\n<hr />\n"},
	{"spec12", "- `one\n- two`\n", "<ul>\n<li>`one</li>\n<li>two`</li>\n</ul>\n"},
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
	{"simple13", "- lu\n  - te", "<ul>\n<li>lu\n<ul>\n<li>te</li>\n</ul>\n</li>\n</ul>\n"},
	{"simple12", "`l*ut*e", "<p>`l<em>ut</em>e</p>\n"},
	{"simple11", "`lu\nte`", "<p><code>lu te</code></p>\n"},
	{"simple10", "lu\n\nte", "<p>lu</p>\n<p>te</p>\n"},
	{"simple9", "* lute", "<ul>\n<li>lute</li>\n</ul>\n"},
	{"simple8", "# lute", "<h1>lute</h1>\n"},
	{"simple7", "> lute", "<blockquote>\n<p>lute</p>\n</blockquote>\n"},
	{"simple6", "l**ut**e", "<p>l<strong>ut</strong>e</p>\n"},
	{"simple5", "l*ut*e", "<p>l<em>ut</em>e</p>\n"},
	{"simple4", "    lute\n", "<pre><code>lute\n</code></pre>\n"},
	{"simple3", "\tlute\n", "<pre><code>lute\n</code></pre>\n"},
	{"simple2", "l`ut`e", "<p>l<code>ut</code>e</p>\n"},
	{"simple1", "lute", "<p>lute</p>\n"},
	{"simple0", "", ""},
}

func TestParse(t *testing.T) {
	for _, test := range parseTests {
		fmt.Println("Test [" + test.name + "]")
		tree, err := Parse(test.name, test.input)
		if nil != err {
			t.Fatalf("%q: unexpected error: %v", test.name, err)
		}

		renderer := NewHTMLRenderer()
		html, err := tree.Render(renderer)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.result != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", tree.name, test.result, html, test.input)
		}
	}
}
