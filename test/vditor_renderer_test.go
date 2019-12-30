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

var vditorDOM2MdTests = []parseTest{

	{"65", `<table><thead><tr><th align="center">col1</th></tr></thead><tbody><tr><td align="center">12</td></tr><tr><td align="center">34<wbr></td></tr></tbody></table>`, "|col1|\n|:---:|\n|12|\n|34|\n"},
	{"64", `<ul data-tight="true"><li data-marker="*">a<ul data-tight="true"><li data-marker="*">a1</li></ul></li><li data-marker="*">b</li><li data-marker="*">c<wbr></li></ul>`, "* a\n  * a1\n* b\n* c\n"},
	{"63", "<ul data-tight=\"true\"><li data-marker=\"*\">foo</li></ul><p>b<wbr>\n</p>", "* foo\n\nb\n"},
	{"62", "<ul><li data-marker=\"*\"><p>foo\n</p></li></ul><p>b<wbr>\n</p>", "* foo\n\nb\n"},
	{"61", "<ul data-tight=\"true\"><li data-marker=\"*\">foo</li><li data-marker=\"*\"><wbr><br></li></ul>", "* foo\n*\n"},
	{"60", "<ul><li data-marker=\"-\" class=\"vditor-task\"><p><input type=\"checkbox\"> foo\n</p></li><li data-marker=\"-\" class=\"vditor-task\"><p><input type=\"checkbox\"> b<wbr>\n</p></li></ul>", "- [ ] foo\n- [ ] b\n"},
	{"59", "<ul><li data-marker=\"-\" class=\"vditor-task\"><p><input type=\"checkbox\" /> foo\n</p></li><li data-marker=\"-\" class=\"vditor-task\"><p><input type=\"checkbox\" /> b<wbr>\n</p></li></ul>", "- [ ] foo\n- [ ] b\n"},
	{"58", "<p><em data-marker=\"*\">foo </em>bar<wbr>\n</p>", "*foo*bar\n"},
	{"57", "<h3>隐藏细节</h3><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\"><pre><code data-code=\"%3Cdetails%3E%0A%3Csummary%3E%E8%BF%99%E9%87%8C%E6%98%AF%E6%91%98%E8%A6%81%E9%83%A8%E5%88%86%E3%80%82%3C%2Fsummary%3E%0A%E8%BF%99%E9%87%8C%E6%98%AF%E7%BB%86%E8%8A%82%E9%83%A8%E5%88%86%E3%80%82%0A%3C%2Fdetails%3E%0A\">&lt;details&gt;&lt;summary&gt;这里是摘要部分。&lt;/summary&gt;这里是细节部分。&lt;/details&gt;<br></code></pre><div class=\"vditor-wysiwyg__preview\" contenteditable=\"false\" data-render=\"false\"></div></div><p>1<wbr></p>", "### 隐藏细节\n\n<details>\n<summary>这里是摘要部分。</summary>\n这里是细节部分。\n</details>\n\n1\n"},
	{"56", "<p>~删除线~</p>", "~删除线~\n"},
	{"55", "<ul data-tight=\"true\"><li data-marker=\"*\">foo</li><li data-marker=\"*\"><br></li><li data-marker=\"*\"><wbr>bar</li></ul>", "* foo\n* bar\n"}, // 在 bar 前面换行剔除空的列表项节点
	{"54", "<p>f<code data-code=\"o\"></code><wbr>o\n</p>", "f`o`o\n"},
	{"53", "<blockquote><p><br></p><p><wbr>foo\n</p></blockquote>", "> foo\n"}, // 在块引用第一个字符前换行
	{"52", "<blockquote><p>foo\n</p><blockquote><p>bar<wbr>\n</p></blockquote></blockquote>", "> foo\n>\n> > bar\n> >\n"},
	{"51", "<blockquote><blockquote><p><wbr>\n</p></blockquote></blockquote>", "\n"},
	{"50", "<blockquote><p><wbr>\n</p></blockquote>", "\n"},
	{"49", "<blockquote><p>f<wbr>\n</p></blockquote>", "> f\n"},
	{"48", "<div class=\"vditor-wysiwyg__block\" data-type=\"math-block\"><pre><code data-code=\"foo\"></code></pre></div>", "$$\nfoo\n$$\n"},
	{"47", "<p><em data-marker=\"*\"><br></em></p><p><em data-marker=\"*\"><wbr>foo</em></p>", "*foo*\n"},
	{"46", "<p><em data-marker=\"*\">foo<wbr></em></p><p><em data-marker=\"*\"></em></p>", "*foo*\n"},
	{"45", "<p><em data-marker=\"*\">foo</em></p><p><em data-marker=\"*\"><wbr><br></em></p>", "*foo*\n"},
	{"44", "<ol><li data-marker=\"1.\"><p>Node.js</p></li><li data-marker=\"2.\"><p>Go<wbr></p></li></ol>", "1. Node.js\n2. Go\n"},
	{"43", "<p>f<i>o</i>o<wbr></p>", "f*o*o\n"},
	{"42", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<br></li><ul><li data-marker=\"*\">b<wbr></li></ul></ul>", "* foo\n  * b\n"},
	{"41", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-marker=\"```\"><pre><code data-code=\"foo\" class=\"language-go\">foo<br></code></pre></div>", "```go\nfoo\n```\n"},
	{"40", "<p>f<span data-marker=\"*\">o</span>ob<wbr></p>", "foob\n"},
	{"39", "<p><b>foo<wbr></b></p>", "**foo**\n"},
	{"38", "<p>```java</p><p><wbr><br></p>", "```java\n"},
	{"37", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<wbr></li><li data-marker=\"*\"></li><li data-marker=\"*\"><br></li></ul>", "* foo\n*\n"}, // 剔除中间的空项
	{"36", "<ul data-tight=\"true\"><li data-marker=\"*\">1<em data-marker=\"*\">2</em></li><li data-marker=\"*\"><em data-marker=\"*\"><wbr><br></em></li></ul>", "* 1*2*\n*\n"},
	{"35", "<ul data-tight=\"true\"><li data-marker=\"*\"><wbr><br></li></ul>", "*\n"},
	{"34", "<p>中<wbr>文</p>", "中文\n"},
	{"33", "<ol data-tight=\"true\"><li data-marker=\"1.\">foo</li></ul>", "1. foo\n"},
	{"32", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<wbr></li></ul>", "* foo\n"},
	{"31", "<ul><li data-marker=\"*\">foo<ul><li data-marker=\"*\">bar</li></ul></li></ul>", "* foo\n  * bar\n"},
	{"30", "<ul><li data-marker=\"*\">foo</li><li data-marker=\"*\"><ul><li data-marker=\"*\"><br /></li></ul></li></ul>", "* foo\n* *\n"},
	{"29", "<p><s>del</s></p>", "~~del~~\n"},
	{"29", "<p>[]()</p>", "[]()\n"},
	{"28", ":octocat:", ":octocat:\n"},
	{"27", "<table><thead><tr><th>abc</th><th>def</th></tr></thead></table>\n", "|abc|def|\n|---|---|\n"},
	{"26", "<p><del data-marker=\"~~\">Hi</del> Hello, world!</p>", "~~Hi~~ Hello, world!\n"},
	{"25", "<p><del data-marker=\"~\">Hi</del> Hello, world!</p>", "~Hi~ Hello, world!\n"},
	{"24", "<ul><li data-marker=\"*\" class=\"vditor-task\"><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo<wbr></li></ul>", "* [X] foo\n"},
	{"23", "<ul><li data-marker=\"*\" class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" /> foo<wbr></li></ul>", "* [ ] foo\n"},
	{"22", "><wbr>", ">\n"},
	{"21", "<p>> foo<wbr></p>", "> foo\n"},
	{"20", "<p>foo</p><p><wbr><br></p>", "foo\n"},
	{"19", "<ul><li data-marker=\"*\">foo</li></ul>", "* foo\n"},
	{"18", "<p><em data-marker=\"*\">foo<wbr></em></p>", "*foo*\n"},
	{"17", "foo bar", "foo bar\n"},
	{"16", "<p><em><strong>foo</strong></em></p>", "***foo***\n"},
	{"15", "<p><strong data-marker=\"__\">foo</strong></p>", "__foo__\n"},
	{"14", "<p><strong data-marker=\"**\">foo</strong></p>", "**foo**\n"},
	{"13", "<h2>foo</h2><p>para<em>em</em></p>", "## foo\n\npara*em*\n"},
	{"12", "<a href=\"/bar\" title=\"baz\">foo</a>", "[foo](/bar \"baz\")\n"},
	{"11", "<img src=\"/bar\" alt=\"foo\" />", "![foo](/bar)\n"},
	{"10", "<img src=\"/bar\" />", "![](/bar)\n"},
	{"9", "<a href=\"/bar\">foo</a>", "[foo](/bar)\n"},
	{"8", "foo<br />bar", "foo\nbar\n"},
	{"7", "<p><code data-code=\"foo\"></code><wbr>\n</p>", "`foo`\n"},
	{"6", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-marker=\"```\"><pre><code data-code=\"foo\">foo<br></code></pre></div>", "```\nfoo\n```\n"},
	{"5", "<ul><li data-marker=\"*\">foo</li></ul>", "* foo\n"},
	{"4", "<blockquote>foo</blockquote>", "> foo\n"},
	{"3", "<h2>foo</h2>", "## foo\n"},
	{"2", "<p><strong><em>foo</em></strong></p>", "***foo***\n"},
	{"1", "<p><strong>foo</strong></p>", "**foo**\n"},
	{"0", "<p>foo</p>", "foo\n"},
}

func TestVditorDOM2Md(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorDOM2MdTests {
		md := luteEngine.VditorDOM2Md(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}

var spinVditorDOMTests = []*parseTest{

	{"47", "<p data-block=\"0\">`1<wbr>`\n</p>", "<p data-block=\"0\"> <code data-code=\"1\">1<wbr></code> \n</p>"},
	{"46", "<p>- [x] f<wbr>\n</p>", "<ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"-\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> f<wbr></li></ul>"},
	{"45", "<ul data-tight=\"true\"><li data-marker=\"-\" class=\"vditor-task\"><input type=\"checkbox\"> foo</li></ul><p>- [ ] b<wbr>\n</p>", "<ul data-block=\"0\"><li data-marker=\"-\" class=\"vditor-task\"><p data-block=\"0\"><input type=\"checkbox\" /> foo\n</p></li><li data-marker=\"-\" class=\"vditor-task\"><p data-block=\"0\"><input type=\"checkbox\" /> b<wbr>\n</p></li></ul>"},
	{"44", "<p>* [ ]<wbr>\n</p>", "<p data-block=\"0\">* [ ]<wbr>\n</p>"},
	{"43", "<p>* [ <wbr>\n</p>", "<p data-block=\"0\">* [ <wbr>\n</p>"},
	{"42", "<p>* [<wbr>\n</p>", "<p data-block=\"0\">* [<wbr>\n</p>"},
	{"40", "<h3>隐藏细节</h3><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\"><pre><code data-code=\"%3Cdetails%3E%0A%3Csummary%3E%E8%BF%99%E9%87%8C%E6%98%AF%E6%91%98%E8%A6%81%E9%83%A8%E5%88%86%E3%80%82%3C%2Fsummary%3E%0A%E8%BF%99%E9%87%8C%E6%98%AF%E7%BB%86%E8%8A%82%E9%83%A8%E5%88%86%E3%80%82%0A%3C%2Fdetails%3E%0A\">&lt;details&gt;&lt;summary&gt;这里是摘要部分。&lt;/summary&gt;这里是细节部分。&lt;/details&gt;<br></code></pre><div class=\"vditor-wysiwyg__preview\" contenteditable=\"false\" data-render=\"false\"></div></div><p>1<wbr></p>", "<h3 data-block=\"0\">隐藏细节</h3><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code data-code=\"%3Cdetails%3E%0A%3Csummary%3E%E8%BF%99%E9%87%8C%E6%98%AF%E6%91%98%E8%A6%81%E9%83%A8%E5%88%86%E3%80%82%3C%2Fsummary%3E%0A%E8%BF%99%E9%87%8C%E6%98%AF%E7%BB%86%E8%8A%82%E9%83%A8%E5%88%86%E3%80%82%0A%3C%2Fdetails%3E\"></code></pre></div><p data-block=\"0\">1<wbr>\n</p>"},
	{"39", "<p>*foo<wbr>*bar\n</p>", "<p data-block=\"0\"><em data-marker=\"*\">foo<wbr></em>bar\n</p>"},
	{"38", "<p>[foo](h<wbr>)\n</p>", "<p data-block=\"0\"><a href=\"h\">foo<wbr></a>\n</p>"},
	{"37", "<blockquote><p><wbr>\n</p></blockquote>", ""},
	{"36", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-marker=\"```\"><pre><code data-code=\"foo\" class=\"language-go\">foo</code></pre></div>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code data-code=\"foo\" class=\"language-go\"></code></pre></div>"},
	{"35", "<p><em data-marker=\"*\">foo</em></p><p><em data-marker=\"*\"><wbr><br></em></p>", "<p data-block=\"0\"><em data-marker=\"*\">foo</em>\n</p><p data-block=\"0\"><em data-marker=\"*\">\n<wbr></em>\n</p>"},
	{"34", "<p> <span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\" data-code=\"a1a\"></code></span> \n</p>", "<p data-block=\"0\"> <span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\" data-code=\"a1a\">a1a</code></span> \n</p>"},
	{"33", "<p><code data-code=\"foo\"></code><wbr>\n</p>", "<p data-block=\"0\"> <code data-code=\"foo\">foo</code> <wbr>\n</p>"},
	{"32", "<p>```<wbr></p>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code data-code=\"\"><wbr>\n</code></pre></div>"},

	{"30", "<p>1. Node.js</p><p>2. Go<wbr></p>", "<ol data-block=\"0\"><li data-marker=\"1.\"><p data-block=\"0\">Node.js\n</p></li><li data-marker=\"2.\"><p data-block=\"0\">Go<wbr>\n</p></li></ol>"},
	{"29", "<p><wbr><br></p>", "<p data-block=\"0\"><wbr>\n</p>"},

	{"27", "<p><wbr></p>", "<p data-block=\"0\"><wbr>\n</p>"},
	{"26", "<p>![alt](src \"title\")</p>", "<p data-block=\"0\"><img src=\"src\" alt=\"alt\" title=\"title\" />\n</p>"},
	{"25", "<pre><code class=\"language-java\"><wbr>\n</code></pre>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code data-code=\"\" class=\"language-java\">\n</code></pre></div>"},
	{"24", "<ul data-tight=\"true\"><li data-marker=\"*\"><wbr><br></li></ul>", "<ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\"><wbr></li></ul>"},
	{"23", "<ol><li data-marker=\"1.\">foo</li></ol>", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\">foo</li></ol>"},
	{"22", "<ul><li data-marker=\"*\">foo</li><li data-marker=\"*\"><ul><li data-marker=\"*\">bar</li></ul></li></ul>", "<ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\">foo</li><li data-marker=\"*\"><ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></li></ul>"},
	{"21", "<p>[foo](/bar \"baz\")</p>", "<p data-block=\"0\"><a href=\"/bar\" title=\"baz\">foo</a>\n</p>"},
	{"20", "<p>[foo](/bar)</p>", "<p data-block=\"0\"><a href=\"/bar\">foo</a>\n</p>"},
	{"19", "<p>[foo]()</p>", "<p data-block=\"0\">[foo]()\n</p>"},
	{"18", "<p>[](/bar)</p>", "<p data-block=\"0\">[](/bar)\n</p>"},
	{"17", "<p>[]()</p>", "<p data-block=\"0\">[]()\n</p>"},
	{"16", "<p>[](</p>", "<p data-block=\"0\">[](\n</p>"},
	{"15", "<p><img alt=\"octocat\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/octocat.png\" title=\"octocat\" /></p>", "<p data-block=\"0\"><img alt=\"octocat\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/octocat.png\" title=\"octocat\" />\n</p>"},
	{"14", ":octocat:", "<p data-block=\"0\"><img alt=\"octocat\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/octocat.png\" title=\"octocat\" />\n</p>"},

	{"12", "<p><s data-marker=\"~~\">Hi</s> Hello, world!</p>", "<p data-block=\"0\"><s data-marker=\"~~\">Hi</s> Hello, world!\n</p>"},
	{"11", "<p><del data-marker=\"~\">Hi</del> Hello, world!</p>", "<p data-block=\"0\"><s data-marker=\"~\">Hi</s> Hello, world!\n</p>"},
	{"10", "<ul><li data-marker=\"*\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> foo<wbr></li></ul>", "<ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> foo<wbr></li></ul>"},
	{"9", "<ul><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> foo<wbr></li></ul>", "<ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> foo<wbr></li></ul>"},
	{"8", "> <wbr>", "<p data-block=\"0\">> <wbr>\n</p>"},
	{"7", "><wbr>", "<p data-block=\"0\">><wbr>\n</p>"},
	{"6", "<p>> foo<wbr></p>", "<blockquote data-block=\"0\"><p data-block=\"0\">foo<wbr>\n</p></blockquote>"},
	{"5", "<p>foo</p><p><wbr><br></p>", "<p data-block=\"0\">foo\n</p><p data-block=\"0\"><wbr>\n</p>"},
	{"4", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<wbr></li></ul>", "<ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\">foo<wbr></li></ul>"},
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
