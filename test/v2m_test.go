// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"testing"

	"github.com/88250/lute"
)

var vditorDOM2MdTests = []parseTest{

	{"117", "<a><img></a>", "[![]()]()\n"},
	{"116", "<a href=\"https://example.com\"><img src=\"https://example.org\" alt=\"example\" title=\"example\"></a>", "[![example](https://example.org \"example\")](https://example.com)\n"},
	{"115", "<ul><li><img src=\"xxx.png\"></li></ul>", "* ![](xxx.png)\n"},
	{"114", "<p data-block=\"0\"><span class=\"vditor-comment vditor-comment--focus\" data-cmtids=\"20201105103606-fpmdc18\">foo<wbr></span></p>", "<span class=\"vditor-comment vditor-comment--focus\" data-cmtids=\"20201105103606-fpmdc18\">foo</span>\n"},
	{"113", "<div data-block=\"0\" data-type=\"footnotes-block\"><p><wbr><br></p></div>", "\n"},
	{"112", "<div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^bignote\"><p data-block=\"0\">foo</p><p data-block=\"0\">bar</p><div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre class=\"vditor-wysiwyg__pre\" style=\"display: none;\"><code>baz\n</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><code>codeblock\n</code></pre></div></li></ol></div>", "[^bignote]: foo\n    \n       bar\n    \n       ```\n       baz\n       ```\n"},
	{"111", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">1</li><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">2</li></ul><li><wbr><br></li></ul>", "* 1\n  * 2\n*\n"},
	{"110", "<ol data-tight=\"true\" start=\"3\" data-block=\"0\"><li data-marker=\"3.\"><p>foo<wbr></p></li></ol>", "3. foo\n"},
	{"109", "<li data-marker=\"*\" class=\"vditor-task\"><p><input checked=\"\" type=\"checkbox\"> foo</p></li><li data-marker=\"*\" class=\"vditor-task\"><p><input type=\"checkbox\"> bar</p></li>", "* [X] foo\n* [ ] bar\n"},
	{"108", "<p data-block=\"0\">foo[^1]</p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-marker=\"^1\"><p data-block=\"0\">bar</p><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">baz</li></ul></li></ol></div>", "foo[^1]\n\n[^1]: bar\n    \n       * baz\n"},
	{"107", "<p data-block=\"0\"><sup data-type=\"footnotes-ref\" data-footnotes-label=\"^1\">1</sup></p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"></li></ol></div>", "[^1]\n\n[^1]:\n"},
	{"106", "<a href=\"\" title=\"baz\">foo</a>", "[foo]( \"baz\")\n"},
	{"105", "<p data-block=\"0\"><strong data-marker=\"**\">foo<em>\u200b\nb<wbr></em></strong></p>", "**foo\n*b***\n"},
	{"104", "<p data-block=\"0\">\u200b<span class=\"vditor-wysiwyg__block\" data-type=\"html-inline\"><code data-type=\"html-inline\">\u200b&lt;foo&gt;</code><span class=\"vditor-wysiwyg__preview\" data-render=\"1\" data-html=\"&amp;lt;foo&amp;gt;\">....</span></span>\u200b	bar<wbr></p>", "<foo>\tbar\n"},
	{"103", `<ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*"><ul data-tight="true" data-marker="-" data-block="0"><li data-marker="-"><p>- -<wbr></p></li></ul></li></ul>`, "* - - -\n"},
	{"102", `<h2 data-block="0" data-marker="-">Setext 标题<wbr></h2>`, "Setext 标题\n-----------\n"},
	{"101", `<ol data-tight="true" data-block="0"><li data-marker="1)"><p>foo</p></li><ol data-tight="true" data-block="0"><li data-marker="1)"><p>bar</p></li><li data-marker="2)"><p>baz</p></li></ol><li><p><wbr><br></p></li></ol>`, "1) foo\n   1) bar\n   2) baz\n1)\n"},
	{"100", "<p data-block=\"0\"><strong data-marker=\"**\">foo <em data-marker=\"_\">bar</em></strong>​ b<wbr></p>", "**foo _bar_** b\n"},
	{"99", `<strong data-marker="**">test and <em data-marker="_">test</em></strong> test`, "**test and _test_** test\n"},
	{"98", `<ul data-tight="true" data-marker="+" data-block="0"><li data-marker="+"><p>foo</p></li><ul data-tight="true" data-marker="+" data-block="0"><li data-marker="+"><p>bar</p></li></ul><li><p><wbr><br></p></li></ul>`, "+ foo\n  + bar\n+\n"},
	{"97", `<ol data-tight="true" data-block="0"><li data-marker="1."><p>foo</p></li><ol data-tight="true" data-block="0"><li data-marker="1."><p>bar</p></li></ol><li><p><wbr><br></p></li></ol>`, "1. foo\n   1. bar\n1.\n"},
	{"96", "<p data-block=\"0\">\u200b<code marker=\"`\">\u200bcode\nspan<wbr></code>\u200b</p>", "`code\nspan`\n"},
	{"95", "<p data-block=\"0\"><strong data-marker=\"**\"><wbr></strong></p>", "\n"},
	{"94", "<blockquote data-block=\"0\"><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">foo<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre style=\"display: none;\"><code>code block\n</code></pre><div class=\"vditor-wysiwyg__preview\" data-render=\"1\"><pre><div class=\"vditor-copy\"><textarea></textarea><span aria-label=\"复制\" onmouseover=\"this.setAttribute('aria-label', '复制')\" class=\"b3-tooltips b3-tooltips__w\" onclick=\"this.previousElementSibling.select();document.execCommand('copy');this.setAttribute('aria-label', '已复制')\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 32 32\" width=\"32px\" height=\"32px\"> <path d=\"M28.681 11.159c-0.694-0.947-1.662-2.053-2.724-3.116s-2.169-2.030-3.116-2.724c-1.612-1.182-2.393-1.319-2.841-1.319h-11.5c-1.379 0-2.5 1.121-2.5 2.5v23c0 1.378 1.121 2.5 2.5 2.5h19c1.378 0 2.5-1.122 2.5-2.5v-15.5c0-0.448-0.137-1.23-1.319-2.841zM24.543 9.457c0.959 0.959 1.712 1.825 2.268 2.543h-4.811v-4.811c0.718 0.556 1.584 1.309 2.543 2.268v0zM28 29.5c0 0.271-0.229 0.5-0.5 0.5h-19c-0.271 0-0.5-0.229-0.5-0.5v-23c0-0.271 0.229-0.5 0.5-0.5 0 0 11.499-0 11.5 0v7c0 0.552 0.448 1 1 1h7v15.5zM18.841 1.319c-1.612-1.182-2.393-1.319-2.841-1.319h-11.5c-1.378 0-2.5 1.121-2.5 2.5v23c0 1.207 0.86 2.217 2 2.45v-25.45c0-0.271 0.229-0.5 0.5-0.5h15.215c-0.301-0.248-0.595-0.477-0.873-0.681z\"></path> </svg></span></div><code class=\"hljs properties\" style=\"max-height: 1000px;\"><span class=\"hljs-attr\">code</span> <span class=\"hljs-string\">block</span></code></pre></div></div></li></ul><p><wbr><br></p></blockquote>", "> * foo\n>   ```\n>   code block\n>   ```\n>\n>\n"},
	{"93", "<p data-block=\"0\"><wbr>    ***</p>", "```\n***\n```\n"},
	{"92", "<p data-block=\"0\">    ***<wbr></p>", "```\n***\n```\n"},
	{"91", "<p data-block=\"0\"><strong><em>\u200b</em></strong></p>", "\n"},
	{"90", "<p data-block=\"0\"><strong>\u200b</strong></p>", "\n"},
	{"90", `<ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*" class="vditor-task"><p><wbr><input type="checkbox"> foo</p></li><li data-marker="*" class="vditor-task"><p><input type="checkbox"> bar</p></li></ul>`, "* [ ] foo\n* [ ] bar\n"},
	{"89", `<ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*" class="vditor-task"><p><input type="checkbox"> foo<wbr><input type="checkbox"><span style="background-color: var(--textarea-background-color); color: var(--textarea-text-color);"> bar</span></p></li></ul>`, "* [ ] foobar\n"},
	{"88", "<p data-block=\"0\">foo<strong>​ b<wbr></strong></p>", "foo **b**\n"},
	{"87", "<p data-block=\"0\">foo<b> b<wbr></b>", "foo **b**\n"},
	{"86", "<strong><em></em>foo</strong>", "**foo**\n"},
	{"85", "<p data-block=\"0\">foo<strong>\u200b</strong></p><span>\u200b</span>", "foo\n"},
	{"84", `<table data-block="0"><thead><tr><th>col1</th></tr></thead><tbody><tr><td>foo<wbr><br></td></tr></tbody></table>`, "| col1 |\n| ---- |\n| foo  |\n"},
	{"83", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"><p>foo</p><p data-block=\"0\">b<wbr></p></li></ul>", "* foo\n\n  b\n"},
	{"82", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\"><p>[x] foo<wbr></p></li></ol>", "1. [x] foo\n"},
	{"81", "<p data-block=\"0\">f&#8203;b</p>", "fb\n"},
	{"80", "<p data-block=\"0\"><span class=\"vditor-wysiwyg__block\" data-type=\"html-inline\">\u200b<code data-type=\"html-inline\" style=\"display: none;\">&lt;foo&gt;</code></span>b<wbr></p>", "<foo>b\n"},
	{"79", "<p>\u200bfoo<wbr></p>", "foo\n"},
	{"78", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"><p>a​​​​</p><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"><p><br></p></li><li data-marker=\"*\"><p><wbr>b</p></li></ul></li></ul>", "* a\n\n  * \n  * b\n"},
	{"77", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"><p>a</p></li><li data-marker=\"*\"><p><wbr><br></p><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"><p>b</p></li></ul></li></ul>", "* a\n* \n\n  * b\n"},
	{"76", "<table data-block=\"0\"><thead><tr><th>col1<wbr></th></tr></thead></table>", "| col1 |\n| ---- |\n"},
	{"75", "<li data-marker=\"*\"><p>foo</p></li><li data-marker=\"*\"><p>bar</p></li>", "* foo\n* bar\n"},
	{"74", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre style=\"display: none;\"><code>foo'%'bar</code></pre></div>", "```\nfoo'%'bar\n```\n"},
	{"73", `<p data-block="0"><img src="/bar" alt="alt text" title="title">`, "![alt text](/bar \"title\")\n"},
	{"72", `<div class="vditor-wysiwyg__block" data-type="code-block"><pre data-block="0"><code></code></pre></div>`, "```\n```\n"},
	{"71", `<ul data-tight="true" data-block="0"><li data-marker="*"><p>123</p><ul data-tight="true" data-block="0"><li data-marker="*"><p>456</p><ul data-tight="true" data-block="0"><li data-marker="*"><p>789</p></li></ul></li></ul></li><li data-marker="*">1</li><li data-marker="*"><wbr><br></li></ul>`, "* 123\n\n  * 456\n\n    * 789\n* 1\n*\n"},
	{"70", "<p data-block=\"0\">/\\_\\_foo__.</p>", "/\\_\\_foo__.\n"},
	{"69", "<p data-block=\"0\">foo<kbd>code</kbd>bar</p>", "foo<kbd>code</kbd>bar\n"},
	{"68", "<p data-block=\"0\">1<wbr><span class=\"vditor-wysiwyg__block\" data-type=\"html-inline\"><code data-type=\"html-inline\" style=\"display:none\">&lt;br&gt;</code><span class=\"vditor-wysiwyg__preview\" data-render=\"1\"><br></span></span>2</p>", "1<br>2\n"},
	{"67", "<table data-block=\"0\"><thead><tr><th>col1</th></tr></thead><tbody><tr><td>1<br>2<wbr></td></tr></tbody></table>", "| col1     |\n| -------- |\n| 1<br />2 |\n"},
	{"66", "<table data-block=\"0\"><thead><tr><th>col1</th></tr></thead><tbody><tr><td><wbr><br></td></tr></tbody></table>", "| col1 |\n| ---- |\n|      |\n"},
	{"65", `<table><thead><tr><th align="center">col1</th></tr></thead><tbody><tr><td align="center">12</td></tr><tr><td align="center">34<wbr></td></tr></tbody></table>`, "| col1 |\n| :--: |\n|  12  |\n|  34  |\n"},
	{"64", `<ul data-tight="true"><li data-marker="*"><p>a</p><ul data-tight="true"><li data-marker="*"><p>a1</p></li></ul></li><li data-marker="*"><p>b<wbr></p></li></ul>`, "* a\n\n  * a1\n* b\n"},
	{"63", "<ul data-tight=\"true\"><li data-marker=\"*\">foo</li></ul><p>b<wbr></p>", "* foo\n\nb\n"},
	{"62", "<ul><li data-marker=\"*\"><p>foo</p></li></ul><p>b<wbr></p>", "* foo\n\nb\n"},
	{"61", "<ul data-tight=\"true\"><li data-marker=\"*\">foo</li><li data-marker=\"*\"><wbr><br></li></ul>", "* foo\n*\n"},
	{"60", "<ul><li data-marker=\"-\" class=\"vditor-task\"><p><input type=\"checkbox\"> foo</p></li><li data-marker=\"-\" class=\"vditor-task\"><p><input type=\"checkbox\"> b<wbr></p></li></ul>", "- [ ] foo\n- [ ] b\n"},
	{"59", "<ul><li data-marker=\"-\" class=\"vditor-task\"><p><input type=\"checkbox\" /> foo</p></li><li data-marker=\"-\" class=\"vditor-task\"><p><input type=\"checkbox\" /> b<wbr></p></li></ul>", "- [ ] foo\n- [ ] b\n"},
	{"58", "<p><em data-marker=\"*\">foo </em>bar<wbr></p>", "*foo* bar\n"},
	{"57", "<h3>隐藏细节</h3><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\"><pre><code>&lt;details&gt;\n&lt;summary&gt;\n这里是摘要部分。\n&lt;/summary&gt;\n这里是细节部分。&lt;/details&gt;<br></code></pre><div class=\"vditor-wysiwyg__preview\" contenteditable=\"false\" data-render=\"1\"></div></div><p>1<wbr></p>", "### 隐藏细节\n\n<details>\n<summary>\n这里是摘要部分。\n</summary>\n这里是细节部分。</details>\n\n1\n"},
	{"56", "<p>~删除线~</p>", "~删除线~\n"},
	{"55", "<ul data-tight=\"true\"><li data-marker=\"*\">foo</li><li data-marker=\"*\"><br></li><li data-marker=\"*\"><wbr>bar</li></ul>", "* foo\n*\n* bar\n"},
	{"54", "<p>f<code>o</code><wbr>o</p>", "f`o`o\n"},
	{"53", "<blockquote><p><br></p><p><wbr>foo</p></blockquote>", ">\n>\n> foo\n"}, // 在块引用第一个字符前换行
	{"52", "<blockquote><p>foo</p><blockquote><p>bar<wbr></p></blockquote></blockquote>", "> foo\n>\n>> bar\n>>\n"},
	{"51", "<blockquote><blockquote><p><wbr></p></blockquote></blockquote>", "\n"},
	{"50", "<blockquote><p><wbr></p></blockquote>", "\n"},
	{"49", "<blockquote><p>f<wbr></p></blockquote>", "> f\n"},
	{"48", "<div class=\"vditor-wysiwyg__block\" data-type=\"math-block\"><pre><code>foo</code></pre></div>", "$$\nfoo\n$$\n"},
	{"47", "<p><em data-marker=\"*\"><br></em></p><p><em data-marker=\"*\"><wbr>foo</em></p>", "*foo*\n"},
	{"46", "<p><em data-marker=\"*\">foo<wbr></em></p><p><em data-marker=\"*\"></em></p>", "*foo*\n"},
	{"45", "<p><em data-marker=\"*\">foo</em></p><p><em data-marker=\"*\"><wbr><br></em></p>", "*foo*\n"},
	{"44", "<ol><li data-marker=\"1.\"><p>Node.js</p></li><li data-marker=\"2.\"><p>Go<wbr></p></li></ol>", "1. Node.js\n2. Go\n"},
	{"43", "<p>f<i>o</i>o<wbr></p>", "f*o*o\n"},
	{"42", "<ul data-tight=\"true\"><li data-marker=\"*\"><p>foo</p></li><ul><li data-marker=\"*\"><p>b<wbr></p></li></ul></ul>", "* foo\n  * b\n"},
	{"41", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-marker=\"```\"><pre><code class=\"language-go\">foo<br></code></pre></div>", "```go\nfoo\n```\n"},
	{"40", "<p>f<span data-marker=\"*\">o</span>ob<wbr></p>", "foob\n"},
	{"39", "<p><b>foo<wbr></b></p>", "**foo**\n"},
	{"38", "<p>```java</p><p><wbr><br></p>", "```java\n"},
	{"37", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<wbr></li><li data-marker=\"*\"></li><li data-marker=\"*\"><br></li></ul>", "* foo\n*\n*\n"},
	{"36", "<ul data-tight=\"true\"><li data-marker=\"*\">1<em data-marker=\"*\">2</em></li><li data-marker=\"*\"><em data-marker=\"*\"><wbr><br></em></li></ul>", "* 1*2*\n*\n"},
	{"35", "<ul data-tight=\"true\"><li data-marker=\"*\"><wbr><br></li></ul>", "*\n"},
	{"34", "<p>中<wbr>文</p>", "中文\n"},
	{"33", "<ol data-tight=\"true\"><li data-marker=\"1.\">foo</li></ul>", "1. foo\n"},
	{"32", "<ul data-tight=\"true\"><li data-marker=\"*\">foo<wbr></li></ul>", "* foo\n"},
	{"31", "<ul data-tight=\"true\"><li data-marker=\"*\"><p>foo</p><ul data-tight=\"true\"><li data-marker=\"*\"><p>bar</p></li></ul></li></ul>", "* foo\n\n  * bar\n"},
	{"30", "<ul data-tight=\"true\"><li data-marker=\"*\"><p>foo</p></li><li data-marker=\"*\"><ul><li data-marker=\"*\"><p><br /></p></li></ul></li></ul>", "* foo\n* *\n"},
	{"29", "<p><s>del</s></p>", "~~del~~\n"},
	{"29", "<p>[]()</p>", "[]()\n"},
	{"28", ":octocat:", ":octocat:\n"},
	{"27", "<table><thead><tr><th>abc</th><th>def</th></tr></thead></table>\n", "| abc | def |\n| --- | --- |\n"},
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
	{"7", "<p><code>foo</code><wbr></p>", "`foo`\n"},
	{"6", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-marker=\"```\"><pre><code>foo<br></code></pre></div>", "```\nfoo\n```\n"},
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

var vditorIRDOM2MdTests = []parseTest{

	{"12", "<div data-block=\"0\" data-type=\"math-block\" class=\"vditor-ir__node\"><span data-type=\"math-block-open-marker\">$$</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code data-type=\"math-block\" class=\"language-math\">s_i = P(V | z_i) = f_s(z_i; \\theta^{s})</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><div data-type=\"math-block\" class=\"language-math\">s_i = P(V | z_i) = f_s(z_i; \\theta^{s})</div></pre><span data-type=\"math-block-close-marker\">$$</span></div>", "$$\ns_i = P(V | z_i) = f_s(z_i; \\theta^{s})\n$$\n"},
	{"11", "<table data-type=\"table\" data-node-id=\"20201028173434-2i1b4ah\"><thead><tr><th>foo</th></tr></thead><tbody><tr><td>bar *baz*<wbr></td></tr></tbody></table>", "| foo       |\n| --------- |\n| bar *baz* |\n"},
	{"10", "<ul data-tight=\"true\" data-marker=\"-\"><li data-marker=\"-\" class=\"vditor-task\" data-node-id=\"\"><input checked=\"\" type=\"checkbox\"></li></ul>", "\n"},
	{"9", "<ul data-tight=\"true\" data-marker=\"*\"><li data-marker=\"*\">foo<div data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">​</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code>bar\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"1\"><div class=\"vditor-copy\"><textarea></textarea><span aria-label=\"复制\" onmouseover=\"this.setAttribute('aria-label', '复制')\" class=\"b3-tooltips b3-tooltips__w\" onclick=\"this.previousElementSibling.select();document.execCommand('copy');this.setAttribute('aria-label', '已复制')\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 32 32\" width=\"32px\" height=\"32px\"> <path d=\"M28.681 11.159c-0.694-0.947-1.662-2.053-2.724-3.116s-2.169-2.030-3.116-2.724c-1.612-1.182-2.393-1.319-2.841-1.319h-11.5c-1.379 0-2.5 1.121-2.5 2.5v23c0 1.378 1.121 2.5 2.5 2.5h19c1.378 0 2.5-1.122 2.5-2.5v-15.5c0-0.448-0.137-1.23-1.319-2.841zM24.543 9.457c0.959 0.959 1.712 1.825 2.268 2.543h-4.811v-4.811c0.718 0.556 1.584 1.309 2.543 2.268v0zM28 29.5c0 0.271-0.229 0.5-0.5 0.5h-19c-0.271 0-0.5-0.229-0.5-0.5v-23c0-0.271 0.229-0.5 0.5-0.5 0 0 11.499-0 11.5 0v7c0 0.552 0.448 1 1 1h7v15.5zM18.841 1.319c-1.612-1.182-2.393-1.319-2.841-1.319h-11.5c-1.378 0-2.5 1.121-2.5 2.5v23c0 1.207 0.86 2.217 2 2.45v-25.45c0-0.271 0.229-0.5 0.5-0.5h15.215c-0.301-0.248-0.595-0.477-0.873-0.681z\"></path> </svg></span></div><code class=\"hljs nginx\"><span class=\"hljs-attribute\">bar</span>\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div>ba<wbr></li></ul>", "* foo\n  ```\n  bar\n  ```\n  ba\n"},
	{"8", "<ul data-marker=\"*\"><li data-marker=\"*\"><p>foo</p><div data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">​</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code>bar\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"1\"><div class=\"vditor-copy\"><textarea></textarea><span aria-label=\"复制\" onmouseover=\"this.setAttribute('aria-label', '复制')\" class=\"b3-tooltips b3-tooltips__w\" onclick=\"this.previousElementSibling.select();document.execCommand('copy');this.setAttribute('aria-label', '已复制')\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 32 32\" width=\"32px\" height=\"32px\"> <path d=\"M28.681 11.159c-0.694-0.947-1.662-2.053-2.724-3.116s-2.169-2.030-3.116-2.724c-1.612-1.182-2.393-1.319-2.841-1.319h-11.5c-1.379 0-2.5 1.121-2.5 2.5v23c0 1.378 1.121 2.5 2.5 2.5h19c1.378 0 2.5-1.122 2.5-2.5v-15.5c0-0.448-0.137-1.23-1.319-2.841zM24.543 9.457c0.959 0.959 1.712 1.825 2.268 2.543h-4.811v-4.811c0.718 0.556 1.584 1.309 2.543 2.268v0zM28 29.5c0 0.271-0.229 0.5-0.5 0.5h-19c-0.271 0-0.5-0.229-0.5-0.5v-23c0-0.271 0.229-0.5 0.5-0.5 0 0 11.499-0 11.5 0v7c0 0.552 0.448 1 1 1h7v15.5zM18.841 1.319c-1.612-1.182-2.393-1.319-2.841-1.319h-11.5c-1.378 0-2.5 1.121-2.5 2.5v23c0 1.207 0.86 2.217 2 2.45v-25.45c0-0.271 0.229-0.5 0.5-0.5h15.215c-0.301-0.248-0.595-0.477-0.873-0.681z\"></path> </svg></span></div><code class=\"hljs nginx\" style=\"max-height: 1000px;\"><span class=\"hljs-attribute\">bar</span>\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div><p>baz</p><div data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">​</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code>b<wbr></code></pre><pre class=\"vditor-ir__preview\" data-render=\"1\"><div class=\"vditor-copy\"><textarea></textarea><span aria-label=\"复制\" onmouseover=\"this.setAttribute('aria-label', '复制')\" class=\"b3-tooltips b3-tooltips__w\" onclick=\"this.previousElementSibling.select();document.execCommand('copy');this.setAttribute('aria-label', '已复制')\"><svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 32 32\" width=\"32px\" height=\"32px\"> <path d=\"M28.681 11.159c-0.694-0.947-1.662-2.053-2.724-3.116s-2.169-2.030-3.116-2.724c-1.612-1.182-2.393-1.319-2.841-1.319h-11.5c-1.379 0-2.5 1.121-2.5 2.5v23c0 1.378 1.121 2.5 2.5 2.5h19c1.378 0 2.5-1.122 2.5-2.5v-15.5c0-0.448-0.137-1.23-1.319-2.841zM24.543 9.457c0.959 0.959 1.712 1.825 2.268 2.543h-4.811v-4.811c0.718 0.556 1.584 1.309 2.543 2.268v0zM28 29.5c0 0.271-0.229 0.5-0.5 0.5h-19c-0.271 0-0.5-0.229-0.5-0.5v-23c0-0.271 0.229-0.5 0.5-0.5 0 0 11.499-0 11.5 0v7c0 0.552 0.448 1 1 1h7v15.5zM18.841 1.319c-1.612-1.182-2.393-1.319-2.841-1.319h-11.5c-1.378 0-2.5 1.121-2.5 2.5v23c0 1.207 0.86 2.217 2 2.45v-25.45c0-0.271 0.229-0.5 0.5-0.5h15.215c-0.301-0.248-0.595-0.477-0.873-0.681z\"></path> </svg></span></div><code class=\"hljs\"></code></pre><span data-type=\"code-block-close-marker\">```</span></div></li></ul>", "* foo\n\n  ```\n  bar\n  ```\n  baz\n\n  ```\n  b\n  ```\n"},
	{"7", "<span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span><em data-newline=\"1\">foo</em><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span>", "*foo*\n"},
	{"6", "<span class=\"vditor-ir__marker\">`</span><code data-newline=\"1\">foo</code><span class=\"vditor-ir__marker\">`</span>", "`foo`\n"},
	{"5", "<p><span data-type=\"code\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">`</span><code data-newline=\"1\">foo</code><span class=\"vditor-ir__marker\">`</span></span></p>", "`foo`\n"},
	{"4", "<ul data-tight=\"true\" data-marker=\"+\"><li data-marker=\"+\">foo</li></ul><ul data-tight=\"true\" data-marker=\"-\"><li data-marker=\"-\">bar<ul data-tight=\"true\" data-marker=\"-\"><li data-marker=\"-\">b<wbr></li></ul></li></ul>", "+ foo\n\n- bar\n  - b\n"},
	{"3", "<ul data-tight=\"true\" data-marker=\"*\"><li data-marker=\"*\"><span data-type=\"inline-node\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span><em data-newline=\"1\">foo</em><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span></span>ba<wbr></li></ul>", "* *foo*ba\n"},
	{"2", "<ul data-tight=\"true\" data-marker=\"*\"><li data-marker=\"*\">foo<ul data-tight=\"true\" data-marker=\"*\"><li data-marker=\"*\">bar</li></ul></li></ul>", "* foo\n  * bar\n"},
	{"1", "<h1 class=\"vditor-ir__node\" data-marker=\"#\"><span class=\"vditor-ir__marker\" data-type=\"heading-marker\"># </span>foo</h1>", "# foo\n"},
	{"0", "<p><span data-type=\"inline-node\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">*</span><em data-newline=\"1\">foo</em><span class=\"vditor-ir__marker\">*</span></span></p>", "*foo*\n"},
}

func TestVditorIRDOM2Md(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorIRDOM2MdTests {
		md := luteEngine.VditorIRDOM2Md(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}
