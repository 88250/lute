// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

// +build javascript

package test

import (
	"testing"

	"github.com/88250/lute"
)

var spinVditorDOMTests = []*parseTest{

	{"114", "<p data-block=\"0\">```<wbr>a b\nc\n</p>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code class=\"language-a\"><wbr>c\n</code></pre></div>"},
	{"113", "<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\"><p>[ <wbr>]</p></li></ul>", "<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\" class=\"vditor-task\"><input type=\"checkbox\" /> <wbr></li></ul>"},
	{"112", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"></li><li data-marker=\"*\"><p>f<wbr></p></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">\u200b</li><li data-marker=\"*\">f<wbr></li></ul>"},
	{"111", "<h1 data-block=\"0\" data-marker=\"#\">foo {#custom-id}<wbr></h1>", "<h1 data-block=\"0\" data-id=\"#custom-id\" data-marker=\"#\">foo <wbr></h1>"},
	{"110", "<div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^111\"><p data-block=\"0\">1<wbr>\n</p></li></ol></div>", "<div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^111\"><p data-block=\"0\">1<wbr>\n</p></li></ol></div>"},
	// 109：重复的脚注定义 marker 会被去重，重现步骤：在脚注定义中换行
	{"109", "<div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"><p data-block=\"0\">foo</p></li><li data-type=\"footnotes-li\" data-marker=\"^1\"><p data-block=\"0\"><wbr>bar\n</p></li></ol></div>", "<div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"><p data-block=\"0\">foo\n</p></li></ol></div>"},
	{"108", "<p data-block=\"0\"><wbr>## heading\n</p>", "<h2 data-block=\"0\" data-marker=\"#\"><wbr>heading</h2>"},
	{"107", "<div class=\"toc-div\" data-type=\"toc-block\"><span class=\"toc-h1\"><a class=\"toc-a\" href=\"#foo\">foo</a></span><br></div>\n\n<h1 data-block=\"0\" data-marker=\"#\">foo</h1>", "<div class=\"vditor-toc\" data-block=\"0\" data-type=\"toc-block\" contenteditable=\"false\"><span data-type=\"toc-h\">foo</span><br></div><p data-block=\"0\"></p><h1 data-block=\"0\" data-marker=\"#\">foo</h1>"},
	{"106", "<p data-block=\"0\"><sup data-type=\"footnotes-ref\" data-footnotes-label=\"^1\">1</sup>\n</p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"></li></ol></div>", "<p data-block=\"0\">\u200b<sup data-type=\"footnotes-ref\" data-footnotes-label=\"^1\">1</sup>\u200b\n</p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"></li></ol></div>"},
	{"105", "<p data-block=\"0\"><span data-type=\"link-ref\" data-link-text=\"1\" data-link-label=\"1\">1</span>\n</p><p data-block=\"0\" data-type=\"link-ref-defs\">[1]: f<wbr>\n</p>", "<p data-block=\"0\">\u200b<span data-type=\"link-ref\" data-link-label=\"1\">1</span>\u200b\n</p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[1]: f<wbr>\n</div>"},
	{"104", "<a href=\"\" title=\"baz\">foo</a>", "<p data-block=\"0\"><a href=\"\"baz\"\">foo</a>\n</p>"},
	{"103", "<p data-block=\"0\"><strong data-marker=\"**\">foo\n<em data-marker=\"*\">ba<wbr></em></strong>\n</p>", "<p data-block=\"0\"><strong data-marker=\"**\">foo\n<em data-marker=\"*\">ba<wbr></em></strong>\n</p>"},
	{"102", "<p data-block=\"0\"><strong data-marker=\"**\">foo<em>\u200b\nb<wbr></em></strong>\n</p>", "<p data-block=\"0\"><strong data-marker=\"**\">foo\n<em data-marker=\"*\">b<wbr></em></strong>\n</p>"},
	{"101", `<ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*"><ul data-tight="true" data-marker="-" data-block="0"><li data-marker="-"><p>- -<wbr></p></li></ul></li></ul>`, "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"><hr data-block=\"0\" /><p data-block=\"0\"><wbr>\n</p></li></ul>"},
	{"100", "<em data-marker=\"*\">foo\nbar</em>\t\n---", "<h2 data-block=\"0\" data-marker=\"-\"><em data-marker=\"*\">foo\nbar</em></h2>"},
	{"99", `<ol data-block="0"><li><p>f<wbr></p></li></ol>`, "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\">f<wbr></li></ol>"},
	{"98", "<p data-block=\"0\">\u200b<code data-marker=\"`\">\u200bcode</code>​    <wbr><code data-marker=\"`\">\u200bcode</code>\u200b\n</p>", "<p data-block=\"0\">\u200b<code data-marker=\"`\">\u200bcode</code>\u200b    <wbr><code data-marker=\"`\">\u200bcode</code>\u200b\n</p>"},
	{"97", "<p data-block=\"0\"><strong data-marker=\"**\">foo\n</strong>b<wbr></p>", "<p data-block=\"0\"><strong data-marker=\"**\">foo</strong>\nb<wbr>\n</p>"},
	{"96", "<p data-block=\"0\">​<code data-marker=\"`\">\u200bcode<wbr></code><code data-marker=\"`\">\u200bspan</code><span>\u200b</span></p>", "<p data-block=\"0\">\u200b<code data-marker=\"`\">\u200bcode<wbr>span</code>\u200b\n</p>"},
	{"95", "<p data-block=\"0\"><strong><em><wbr>\u200b</em></strong></p>", ""},
	{"94", "<p data-block=\"0\">\u200b<code data-marker=\"`\">\u200bcode\nspan<wbr></code>\u200b\n</p>", "<p data-block=\"0\">\u200b<code data-marker=\"`\">\u200bcode span<wbr></code>\u200b\n</p>"},
	{"93", "<p data-block=\"0\"><strong data-marker=\"**\"><wbr></strong>\n</p>", "<p data-block=\"0\"><wbr>\n</p>"},
	{"92", "<p data-block=\"0\"><strong><em>\u200b<wbr></em></strong></p>", ""},
	{"91", "<p data-block=\"0\"><wbr>    ***\n</p>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code><wbr>***\n</code></pre></div>"},
	{"90", "<p data-block=\"0\">    ***<wbr>\n</p>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code>***<wbr>\n</code></pre></div>"},
	{"89", `<ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*" class="vditor-task"><p><wbr><input type="checkbox"> foo</p></li></ul>`, "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> <wbr>foo</li></ul>"},
	{"88", `<ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*"><p>test</p></li></ul><ul data-tight="true" data-marker="-" data-block="0"><li data-marker="-"><p>--<wbr></p></li></ul>`, "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">test</li></ul><hr data-block=\"0\" /><p data-block=\"0\"><wbr>\n</p>"},
	{"87", `<ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*" class="vditor-task"><p><input type="checkbox"> foo</p><ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*" class="vditor-task"><p><input type="checkbox"> bar</p></li></ul><p data-block="0">b<wbr></p></li></ul>`, "<ul data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><p data-block=\"0\"><input type=\"checkbox\" /> foo\n</p><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> bar</li></ul><p data-block=\"0\">b<wbr>\n</p></li></ul>"},
	{"86", "<p data-block=\"0\"><strong><s>foo</s></strong>bar<wbr></p>", "<p data-block=\"0\"><strong data-marker=\"**\"><s data-marker=\"~~\">foo</s></strong>bar<wbr>\n</p>"},
	{"85", "<p data-block=\"0\"><span data-marker=\"**\"><b>foo </b><span>b<wbr></span></span>", "<p data-block=\"0\"><strong data-marker=\"**\">foo</strong> b<wbr>\n</p>"},
	{"84", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1)\"><p>f<wbr></p></li></ol>", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1)\">f<wbr></li></ol>"},
	{"83", "<p data-block=\"0\">1) f<wbr></p>", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1)\">f<wbr></li></ol>"},
	{"82", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"2.\"><p>bar<wbr></p></li></ol>", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\">bar<wbr></li></ol>"},
	{"81", "<p data-block=\"0\"><strong>\u200b<em>\u200b<s>\u200b1</s></em></strong><wbr></p>", "<p data-block=\"0\"><em data-marker=\"*\"><strong data-marker=\"**\"><s data-marker=\"~~\">1</s></strong></em><wbr>\n</p>"},
	{"80", "<s><em>\u200b</em></s>", ""},
	{"79", "<p data-block=\"0\"><b>&#8203;</b></p>", ""},
	{"78", "<p data-block=\"0\">``foo``<wbr></p>", "<p data-block=\"0\">\u200b<code data-marker=\"``\">\u200bfoo</code>\u200b<wbr>\n</p>"},
	{"77", `<p data-block="0">1<wbr><span style="background-color: var(--textarea-focus-background-color); color: var(--textarea-text-color);">​</span><span class="vditor-wysiwyg__block" data-type="math-inline" style="background-color: var(--textarea-focus-background-color); color: var(--textarea-text-color);"><code data-type="math-inline">foo</code><span class="vditor-wysiwyg__preview" data-render="false"><span class="vditor-math" data-math="foo"><span class="katex"><span class="katex-html" aria-hidden="true"><span class="base"><span class="strut" style="height:0.85396em;vertical-align:-0.19444em;"></span><span class="mord mathdefault" style="margin-right:0.05724em;">f</span><span class="mord mathdefault" style="margin-right:0.05724em;">o</span><span class="mord mathdefault" style="margin-right:0.05724em;">o</span></span></span></span></span></span></span><span style="background-color: var(--textarea-focus-background-color); color: var(--textarea-text-color);">​</span></p>`, "<p data-block=\"0\">1<wbr><span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\">\u200bfoo</code></span>\u200b\n</p>"},
	{"76", "<ul><li data-marker=\"1.\"><p>12<wbr></p></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">12<wbr></li></ul>"},
	{"75", "<ol><li data-marker=\"*\"><p>foo<wbr></p></li></ol>", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\">foo<wbr></li></ol>"},
	{"74", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\" class=\"vditor-task\"><p><input checked=\"\" type=\"checkbox\"> f<wbr></p></li></ol>", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> f<wbr></li></ol>"},
	{"73", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\"><p>[x] foo<wbr></p></li></ol>", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> foo<wbr></li></ol>"},
	{"72", "<p data-block=\"0\">foo\n-<wbr></p>", "<h2 data-block=\"0\" data-marker=\"-\">foo<wbr></h2>"},
	{"71", "<p data-block=\"0\">foo<span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\">\u200bbar</code></span> \n</p>", "<p data-block=\"0\">foo<span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\">\u200bbar</code></span>\u200b\n</p>"},
	{"70", "<p data-block=\"0\"> <span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\">\u200bfoo</code></span> \n</p>", "<p data-block=\"0\">\u200b<span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\">\u200bfoo</code></span>\u200b\n</p>"},
	{"69", "<p data-block=\"0\"><a href=\"/bar\"><code>foo</code></a><wbr>\n</p>", "<p data-block=\"0\"><a href=\"/bar\">\u200b<code data-marker=\"`\">\u200bfoo</code></a><wbr>\n</p>"},
	{"68", `<p data-block="0">|foo|bar|<wbr></p>`, "<p data-block=\"0\">|foo|bar|<wbr>\n</p>"},
	{"67", `<ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*"><p>[ ]<wbr></p></li></ul>`, "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> <wbr></li></ul>"},
	{"66", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><p><input type=\"checkbox\" checked=\"checked\"><wbr> foo</p></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> <wbr> foo</li></ul>"},
	{"65", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"><p>foo<em data-marker=\"*\">bar</em></p></li><li data-marker=\"*\"><p><em data-marker=\"*\"><wbr><br></em></p></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">foo<em data-marker=\"*\">bar</em></li><li data-marker=\"*\"><wbr></li></ul>"},
	{"64", "<p data-block=\"0\">[foo<wbr>](/bar)", "<p data-block=\"0\"><a href=\"/bar\">foo<wbr></a>\n</p>"},
	{"63", "<p data-block=\"0\">![foo<wbr>](/bar)", "<p data-block=\"0\"><img src=\"/bar\" alt=\"foo\" />\n</p>"},
	{"62", "<p data-block=\"0\"><strong data-marker=\"__\"><wbr><br></strong></p>", "<p data-block=\"0\"><wbr>\n</p>"},
	{"61", "<p data-block=\"0\">_foo_<wbr></p>", "<p data-block=\"0\"><em data-marker=\"_\">foo</em><wbr>\n</p>"},
	{"60", "<p data-block=\"0\">foo\n=<wbr></p>", "<h1 data-block=\"0\" data-marker=\"=\">foo<wbr></h1>"},
	{"59", "<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\"><p>foo</p><ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\"><p>bar<wbr><p></li></ul></li></ul>", "<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\">foo<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\">bar<wbr></li></ul></li></ul>"},
	{"58", "<p data-block=\"0\">![](/bar)<wbr>\n</p>", "<p data-block=\"0\"><img src=\"/bar\" alt=\"\" /><wbr>\n</p>"},
	{"57", "<p data-block=\"0\">/<span data-type=\"backslash\"><span>\\</span>_</span><span data-type=\"backslash\"><span>\\</span>_</span>foo__.\n</p>", "<p data-block=\"0\">/<span data-type=\"backslash\"><span>\\</span>_</span><span data-type=\"backslash\"><span>\\</span>_</span>foo__.\n</p>"},
	{"56", "<p data-block=\"0\">[foo](bar<wbr>)\n</p>", "<p data-block=\"0\"><a href=\"bar\">foo<wbr></a>\n</p>"},
	{"55", "<p data-block=\"0\">[foo<wbr>](bar)\n</p>", "<p data-block=\"0\"><a href=\"bar\">foo<wbr></a>\n</p>"},
	{"54", "<table data-block=\"0\"><thead><tr><th>col1</th></tr></thead><tbody><tr><td><wbr>\n</td></tr></tbody></table>", "<table data-block=\"0\"><thead><tr><th>col1</th></tr></thead><tbody><tr><td><wbr> </td></tr><tr><td> </td></tr></tbody></table>"},
	{"53", "<table data-block=\"0\"><thead><tr><th>col1</th></tr></thead><tbody><tr><td><wbr><br></td></tr></tbody></table>", "<table data-block=\"0\"><thead><tr><th>col1</th></tr></thead><tbody><tr><td><wbr> </td></tr></tbody></table>"},
	{"52", "<p data-block=\"0\">---<wbr>\n</p>", "<hr data-block=\"0\" /><p data-block=\"0\"><wbr>\n</p>"},
	{"51", "<p data-block=\"0\">### <wbr>\n</p>", "<p data-block=\"0\">### <wbr>\n</p>"},
	{"50", "<details open=\"\">\n<summary>foo</summary><ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></details>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;details open=&quot;&quot;&gt;\n&lt;summary&gt;foo&lt;/summary&gt;</code></pre></div><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;/details&gt;</code></pre></div>"},
	{"49", "<details>\n<summary>foo</summary><ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></details>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;details&gt;\n&lt;summary&gt;foo&lt;/summary&gt;</code></pre></div><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;/details&gt;</code></pre></div>"},
	{"48", "<p data-block=\"0\"><a href=\"中文\">link</a><wbr>\n</p>", "<p data-block=\"0\"><a href=\"中文\">link</a><wbr>\n</p>"},
	{"47", "<p data-block=\"0\">`1<wbr>`\n</p>", "<p data-block=\"0\">\u200b<code data-marker=\"`\">\u200b1<wbr></code>\u200b\n</p>"},
	{"46", "<p>- [x] f<wbr>\n</p>", "<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\" class=\"vditor-task\"><input checked=\"\" type=\"checkbox\" /> f<wbr></li></ul>"},
	{"45", "<ul data-tight=\"true\"><li data-marker=\"-\" class=\"vditor-task\"><input type=\"checkbox\"> foo</li></ul><p>- [ ] b<wbr>\n</p>", "<ul data-marker=\"-\" data-block=\"0\"><li data-marker=\"-\" class=\"vditor-task\"><p data-block=\"0\"><input type=\"checkbox\" /> foo\n</p></li><li data-marker=\"-\" class=\"vditor-task\"><p data-block=\"0\"><input type=\"checkbox\" /> b<wbr>\n</p></li></ul>"},
	{"44", "<p data-block=\"0\">* [ ]<wbr>\n</p>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> <wbr></li></ul>"},
	{"43", "<p>* [ <wbr>\n</p>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">[ <wbr></li></ul>"},
	{"42", "<p>* [<wbr>\n</p>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">[<wbr></li></ul>"},
	{"40", "<h3>隐藏细节</h3><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\"><pre><code>&lt;details&gt;\n&lt;summary&gt;\n这里是摘要部分。\n&lt;/summary&gt;\n这里是细节部分。\n&lt;/details&gt;<br></code></pre><div class=\"vditor-wysiwyg__preview\" contenteditable=\"false\" data-render=\"false\"></div></div><p>1<wbr></p>", "<h3 data-block=\"0\" data-marker=\"#\">隐藏细节</h3><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;details&gt;\n&lt;summary&gt;\n这里是摘要部分。\n&lt;/summary&gt;\n这里是细节部分。\n&lt;/details&gt;</code></pre></div><p data-block=\"0\">1<wbr>\n</p>"},
	{"39", "<p>*foo<wbr>*bar\n</p>", "<p data-block=\"0\"><em data-marker=\"*\">foo<wbr></em>bar\n</p>"},
	{"38", "<p>[foo](b<wbr>)\n</p>", "<p data-block=\"0\"><a href=\"b\">foo<wbr></a>\n</p>"},
	{"37", "<blockquote><p><wbr>\n</p></blockquote>", ""},
	{"36", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-marker=\"```\"><pre><code class=\"language-go\">foo</code></pre></div>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code class=\"language-go\">foo\n</code></pre></div>"},
	{"35", "<p><em data-marker=\"*\">foo</em></p><p><em data-marker=\"*\"><wbr><br></em></p>", "<p data-block=\"0\"><em data-marker=\"*\">foo</em>\n</p><p data-block=\"0\"><wbr>\n</p>"},
	{"34", "<p> <span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\">foo</code></span> \n</p>", "<p data-block=\"0\">\u200b<span class=\"vditor-wysiwyg__block\" data-type=\"math-inline\"><code data-type=\"math-inline\">\u200bfoo</code></span>\u200b\n</p>"},
	{"33", "<p><code>foo</code><wbr>\n</p>", "<p data-block=\"0\">\u200b<code data-marker=\"`\">\u200bfoo</code>\u200b<wbr>\n</p>"},
	{"32", "<p>```<wbr></p>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code><wbr>\n</code></pre></div>"},
	{"31", "<p>\u200bfoo<wbr></p>", "<p data-block=\"0\">foo<wbr>\n</p>"},
	{"30", "<p>1. Node.js</p><p>2. Go<wbr></p>", "<ol data-block=\"0\"><li data-marker=\"1.\"><p data-block=\"0\">Node.js\n</p></li><li data-marker=\"2.\"><p data-block=\"0\">Go<wbr>\n</p></li></ol>"},
	{"29", "<p><wbr><br></p>", "<p data-block=\"0\"><wbr>\n</p>"},
	{"28", "<p data-block=\"0\">❤️<wbr>\n</p>", "<p data-block=\"0\">❤️<wbr>\n</p>"},
	{"27", "<p><wbr></p>", "<p data-block=\"0\"><wbr>\n</p>"},
	{"26", "<p>![alt](src \"title\")</p>", "<p data-block=\"0\"><img src=\"src\" alt=\"alt\" title=\"title\" />\n</p>"},
	{"25", "<pre><code class=\"language-java\"><wbr>\n</code></pre>", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code class=\"language-java\"><wbr>\n</code></pre></div>"},
	{"24", "<ul data-tight=\"true\"><li data-marker=\"*\"><wbr><br></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\"><wbr></li></ul>"},
	{"23", "<ol><li data-marker=\"1.\">foo</li></ol>", "<ol data-tight=\"true\" data-block=\"0\"><li data-marker=\"1.\">foo</li></ol>"},
	{"22", "<ul><li data-marker=\"*\">foo</li><li data-marker=\"*\"><ul><li data-marker=\"*\">bar</li></ul></li></ul>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">foo</li><li data-marker=\"*\"><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></li></ul>"},
	{"21", "<p>[foo](/bar \"baz\")</p>", "<p data-block=\"0\"><a href=\"/bar\" title=\"baz\">foo</a>\n</p>"},
	{"20", "<p>[foo](/bar)</p>", "<p data-block=\"0\"><a href=\"/bar\">foo</a>\n</p>"},
	{"19", "<p>[foo]()</p>", "<p data-block=\"0\"><a href=\"\">foo</a>\n</p>"},
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
	luteEngine.ToC = true

	for _, test := range spinVditorDOMTests {
		html := luteEngine.SpinVditorDOM(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var spinVditorIRDOMTests = []*parseTest{

	{"4", "<blockquote data-block=\"0\"><p data-block=\"0\"><wbr>\n</p></blockquote>", "<p data-block=\"0\">&gt; <wbr>\n</p>"},
	{"3", "<blockquote data-block=\"0\"><p data-block=\"0\">fo<wbr>\n</p></blockquote>", "<blockquote data-block=\"0\"><p data-block=\"0\">fo<wbr>\n</p></blockquote>"},
	{"2", "<p data-block=\"0\"><span data-type=\"inline-node\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">*</span><em data-newline=\"1\">foo</em><span class=\"vditor-ir__marker\">*</span></span>\n</p>", "<p data-block=\"0\"><span data-type=\"inline-node\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">*</span><em data-newline=\"1\">foo</em><span class=\"vditor-ir__marker\">*</span></span>\n</p>"},
	{"1", "<p data-block=\"0\">f<wbr></p><p data-block=\"0\">bar\n</p>", "<p data-block=\"0\">f<wbr>\n</p><p data-block=\"0\">bar\n</p>"},
	{"0", "<p data-block=\"0\">foo\n</p><p data-block=\"0\"><wbr><br></p>", "<p data-block=\"0\">foo\n</p><p data-block=\"0\"><wbr>\n</p>"},
}

func TestSpinVditorIRDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ToC = true

	for _, test := range spinVditorIRDOMTests {
		html := luteEngine.SpinVditorIRDOM(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
