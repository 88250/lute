// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

// +build javascript

package test

import (
	"testing"

	"github.com/88250/lute"
)

var spinVditorDOMTests = []*parseTest{

	{"68", `<p data-block="0">|foo|bar|<wbr></p>`, "<p data-block=\"0\">|foo|bar|<wbr>\n</p>"},
	{"67", `<ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*"><p>[ ]<wbr></p></li></ul>`, "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> <wbr></li></ul>"},
	{"66", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><p><input type=\"checkbox\" checked=\"checked\"><wbr> foo</p></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> <wbr> foo</li></ul>"},
	{"65", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"><p>foo<em data-marker=\"*\">bar</em></p></li><li data-marker=\"*\"><p><em data-marker=\"*\"><wbr><br></em></p></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">foo<em data-marker=\"*\">bar</em></li><li data-marker=\"*\"><wbr></li></ul>"},
	{"64", "<p data-block=\"0\">[foo<wbr>](/bar)", "<p data-block=\"0\"><a href=\"/bar\">foo<wbr></a>\n</p>"},
	{"63", "<p data-block=\"0\">![foo<wbr>](/bar)", "<p data-block=\"0\"><img src=\"/bar\" alt=\"foo\" />\n</p>"},
	{"62", "<p data-block=\"0\"><strong data-marker=\"__\"><wbr><br></strong></p>", "<p data-block=\"0\"><wbr>\n</p>"},
	{"61", "<p data-block=\"0\">_foo_<wbr></p>", "<p data-block=\"0\"><em data-marker=\"_\">foo</em><wbr>\n</p>"},
	{"60", "<p data-block=\"0\">foo\n=<wbr></p>", "<h1 data-block=\"0\">foo</h1>"},
	{"59", "<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\"><p>foo</p><ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\"><p>bar<wbr><p></li></ul></li></ul>", "<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\">foo<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\">bar<wbr></li></ul></li></ul>"},
	{"58", "<p data-block=\"0\">![](/bar)<wbr>\n</p>", "<p data-block=\"0\"><img src=\"/bar\" alt=\"\" /><wbr>\n</p>"},
	{"57", "<p data-block=\"0\">/<span data-type=\"backslash\"><span>\\</span>_</span><span data-type=\"backslash\"><span>\\</span>_</span>foo__.\n</p>", "<p data-block=\"0\">/<span data-type=\"backslash\"><span>\\</span>_</span><span data-type=\"backslash\"><span>\\</span>_</span>foo__.\n</p>"},
	{"56", "<p data-block=\"0\">[foo](bar<wbr>)\n</p>", "<p data-block=\"0\"><a href=\"bar\">foo<wbr></a>\n</p>"},
	{"55", "<p data-block=\"0\">[foo<wbr>](bar)\n</p>", "<p data-block=\"0\"><a href=\"bar\">foo<wbr></a>\n</p>"},
	{"54", "<table data-block=\"0\"><thead><tr><th>col1</th></tr></thead><tbody><tr><td><wbr>\n</td></tr></tbody></table>", "<table data-block=\"0\"><thead><tr><th>col1</th></tr></thead><tbody><tr><td><wbr></td></tr><tr><td></td></tr></tbody></table>"},
	{"53", "<table data-block=\"0\"><thead><tr><th>col1</th></tr></thead><tbody><tr><td><wbr><br></td></tr></tbody></table>", "<table data-block=\"0\"><thead><tr><th>col1</th></tr></thead><tbody><tr><td><wbr></td></tr></tbody></table>"},
	{"52", "<p data-block=\"0\">---<wbr>\n</p>", "<hr data-block=\"0\" />"},
	{"51", "<p data-block=\"0\">### <wbr>\n</p>", "<p data-block=\"0\">### <wbr>\n</p>"},
	{"50", "<details open=\"\">\n<summary>foo</summary><ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></details>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;details open=&quot;&quot;&gt;\n&lt;summary&gt;foo&lt;/summary&gt;</code></pre></div><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;/details&gt;</code></pre></div>"},
	{"49", "<details>\n<summary>foo</summary><ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></details>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;details&gt;\n&lt;summary&gt;foo&lt;/summary&gt;</code></pre></div><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;/details&gt;</code></pre></div>"},
	{"49", "<p data-block=\"0\"><a href=\"/bar\"><code>foo</code></a><wbr>\n</p>", "<p data-block=\"0\"><a href=\"/bar\"><code>foo</code></a><wbr>\n</p>"},
	{"48", "<p data-block=\"0\"><a href=\"中文\">link</a><wbr>\n</p>", "<p data-block=\"0\"><a href=\"中文\">link</a><wbr>\n</p>"},
	{"47", "<p data-block=\"0\">`1<wbr>`\n</p>", "<p data-block=\"0\"> <code>1<wbr></code> \n</p>"},
	{"46", "<p>- [x] f<wbr>\n</p>", "<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> f<wbr></li></ul>"},
	{"45", "<ul data-tight=\"true\"><li data-marker=\"-\" class=\"vditor-task\"><input type=\"checkbox\"> foo</li></ul><p>- [ ] b<wbr>\n</p>", "<ul data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\" class=\"vditor-task\"><p data-block=\"0\"><input type=\"checkbox\" /> foo\n</p></li><li data-marker=\"-\" class=\"vditor-task\"><p data-block=\"0\"><input type=\"checkbox\" /> b<wbr>\n</p></li></ul>"},
	{"44", "<p data-block=\"0\">* [ ]<wbr>\n</p>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> <wbr></li></ul>"},
	{"43", "<p>* [ <wbr>\n</p>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">[ <wbr></li></ul>"},
	{"42", "<p>* [<wbr>\n</p>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">[<wbr></li></ul>"},
	{"40", "<h3>隐藏细节</h3><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\"><pre><code>&lt;details&gt;\n&lt;summary&gt;\n这里是摘要部分。\n&lt;/summary&gt;\n这里是细节部分。\n&lt;/details&gt;<br></code></pre><div class=\"vditor-wysiwyg__preview\" contenteditable=\"false\" data-render=\"false\"></div></div><p>1<wbr></p>", "<h3 data-block=\"0\">隐藏细节</h3><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;details&gt;\n&lt;summary&gt;\n这里是摘要部分。\n&lt;/summary&gt;\n这里是细节部分。\n&lt;/details&gt;</code></pre></div><p data-block=\"0\">1<wbr>\n</p>"},
	{"39", "<p>*foo<wbr>*bar\n</p>", "<p data-block=\"0\"><em data-marker=\"*\">foo<wbr></em>bar\n</p>"},
	{"38", "<p>[foo](b<wbr>)\n</p>", "<p data-block=\"0\"><a href=\"b\">foo<wbr></a>\n</p>"},
	{"37", "<blockquote><p><wbr>\n</p></blockquote>", ""},
	{"36", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-marker=\"```\"><pre><code class=\"language-go\">foo</code></pre></div>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code class=\"language-go\">foo\n</code></pre></div>"},
	{"35", "<p><em data-marker=\"*\">foo</em></p><p><em data-marker=\"*\"><wbr><br></em></p>", "<p data-block=\"0\"><em data-marker=\"*\">foo</em>\n</p><p data-block=\"0\"><wbr>\n</p>"},
	{"34", "<p> <span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\">foo</code></span> \n</p>", "<p data-block=\"0\"> <span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\">foo</code></span> \n</p>"},
	{"33", "<p><code>foo</code><wbr>\n</p>", "<p data-block=\"0\"> <code>foo</code> <wbr>\n</p>"},
	{"32", "<p>```<wbr></p>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code><wbr>\n</code></pre></div>"},

	{"30", "<p>1. Node.js</p><p>2. Go<wbr></p>", "<ol data-block=\"0\"><li data-marker=\"1.\"><p data-block=\"0\">Node.js\n</p></li><li data-marker=\"2.\"><p data-block=\"0\">Go<wbr>\n</p></li></ol>"},
	{"29", "<p><wbr><br></p>", "<p data-block=\"0\"><wbr>\n</p>"},

	{"27", "<p><wbr></p>", "<p data-block=\"0\"><wbr>\n</p>"},
	{"26", "<p>![alt](src \"title\")</p>", "<p data-block=\"0\"><img src=\"src\" alt=\"alt\" title=\"title\" />\n</p>"},
	{"25", "<pre><code class=\"language-java\"><wbr>\n</code></pre>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code class=\"language-java\"><wbr>\n</code></pre></div>"},
	{"24", "<ul data-tight=\"true\"><li data-marker=\"*\"><wbr><br></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"><wbr></li></ul>"},
	{"23", "<ol><li data-marker=\"1.\">foo</li></ol>", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\">foo</li></ol>"},
	{"22", "<ul><li data-marker=\"*\">foo</li><li data-marker=\"*\"><ul><li data-marker=\"*\">bar</li></ul></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">foo</li><li data-marker=\"*\"><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></li></ul>"},
	{"21", "<p>[foo](/bar \"baz\")</p>", "<p data-block=\"0\"><a href=\"/bar\" title=\"baz\">foo</a>\n</p>"},
	{"20", "<p>[foo](/bar)</p>", "<p data-block=\"0\"><a href=\"/bar\">foo</a>\n</p>"},
	{"19", "<p>[foo]()</p>", "<p data-block=\"0\">[foo]()\n</p>"},
	{"18", "<p>[](/bar)</p>", "<p data-block=\"0\">[](/bar)\n</p>"},
	{"17", "<p>[]()</p>", "<p data-block=\"0\">[]()\n</p>"},
	{"16", "<p>[](</p>", "<p data-block=\"0\">[](\n</p>"},
	{"15", "<p><img alt=\"octocat\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/octocat.png\" title=\"octocat\" /></p>", "<p data-block=\"0\"><img alt=\"octocat\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/octocat.png\" title=\"octocat\" />\n</p>"},
	{"14", ":octocat:", "<p data-block=\"0\"><img alt=\"octocat\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/octocat.png\" title=\"octocat\" />\n</p>"},
	{"13", "<p>1、foo</p>", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\">foo</li></ol>"},
	{"12", "<p><s data-marker=\"~~\">Hi</s> Hello, world!</p>", "<p data-block=\"0\"><s data-marker=\"~~\">Hi</s> Hello, world!\n</p>"},
	{"11", "<p><del data-marker=\"~\">Hi</del> Hello, world!</p>", "<p data-block=\"0\"><s data-marker=\"~\">Hi</s> Hello, world!\n</p>"},
	{"10", "<ul data-tight=\"true\"><li data-marker=\"*\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> foo<wbr></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> foo<wbr></li></ul>"},
	{"9", "<ul data-tight=\"true\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> foo<wbr></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> foo<wbr></li></ul>"},
	{"8", "> <wbr>", "<p data-block=\"0\">&gt; <wbr>\n</p>"},
	{"7", "><wbr>", "<p data-block=\"0\">&gt;<wbr>\n</p>"},
	{"6", "<p>> foo<wbr></p>", "<blockquote data-block=\"0\"><p data-block=\"0\">foo<wbr>\n</p></blockquote>"},
	{"5", "<p>foo</p><p><wbr><br></p>", "<p data-block=\"0\">foo\n</p><p data-block=\"0\"><wbr>\n</p>"},
	{"4", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<wbr></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">foo<wbr></li></ul>"},
	{"3", "<p><em data-marker=\"*\">foo<wbr></em></p>", "<p data-block=\"0\"><em data-marker=\"*\">foo<wbr></em>\n</p>"},
	{"2", "<p>foo<wbr></p>", "<p data-block=\"0\">foo<wbr>\n</p>"},
	{"1", "<p><strong data-marker=\"**\">foo</strong></p>", "<p data-block=\"0\"><strong data-marker=\"**\">foo</strong>\n</p>"},
	{"0", "<p>foo</p>", "<p data-block=\"0\">foo\n</p>"},
}

func TestSpinVditorDOM(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range spinVditorDOMTests {
		html := luteEngine.SpinVditorDOM(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
