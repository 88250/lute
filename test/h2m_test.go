// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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

var html2MdTests = []parseTest{

	{"20", "<h2 id=\"whats-markdown\">What’s Markdown?<a class=\"anchorjs-link \" aria-label=\"Anchor\" data-anchorjs-icon=\"\uE9CB\" href=\"https://www.markdownguide.org/getting-started/#whats-markdown\"></a></h2>", "## What’s Markdown?\n"},
	{"19", "<pre><span>`foo`</span></pre>", "`foo`\n"},
	{"18", "<del>foo</del>", "~foo~\n"},
	{"17", "<img src=\"bar.png\" alt=\"foo\">", "![foo](bar.png)\n"},
	{"16", "foo<br>bar", "foo\nbar\n"},
	{"15", "<em>foo</em>", "*foo*\n"},
	{"14", "<hr>", "---\n"},
	{"13", "<blockquote>foo</blockquote>", "> foo\n"},
	{"12", "<h1>foo</h1>", "# foo\n"},
	{"11", "<li>foo</li><li>bar</li>", "* foo\n* bar\n"},
	{"10", `<p data-block="0">foo'%'bar</p>`, "foo'%'bar\n"},
	{"9", `<code class="language-text">&gt;</code>`, "`>`\n"},
	{"8", `<div><a href="/bar">foo</a></div>`, "[foo](/bar)\n"},
	{"7", `<ul><li><p>Java</p><ul><li><p>Spring</p></li></ul></li></ul>`, "* Java\n  * Spring\n"},
	{"6", `<!--StartFragment--><p>这是一篇讲解如何正确使用<span>&nbsp;</span><strong>Markdown</strong><span>&nbsp;</span>的排版示例，学会这个很有必要，能让你的文章有更佳清晰的排版。</p><!--EndFragment-->`, "这是一篇讲解如何正确使用 **Markdown** 的排版示例，学会这个很有必要，能让你的文章有更佳清晰的排版。\n"},
	{"5", `<!--StartFragment--><ul><li><input checked="" disabled="" type="checkbox"><span>&nbsp;</span>发布 Solo</li></ul><!--EndFragment-->`, "* [X] 发布 Solo\n"},
	{"4", "<span>&nbsp;</span>发布 Solo", "发布 Solo\n"},
	{"3", "<pre><ul><li>foo</li></ul></pre>", "<pre><ul><li>foo</li></ul></pre>\n"},
	{"2", "<pre><span>//&#32;Lute&#32;-&#32;A&#32;structured&#32;markdown&#32;engine.<br></span><span>//&#32;Copyright&#32;(c)&#32;2019-present,&#32;b3log.org</span></pre>", "// Lute - A structured Markdown engine.\n// Copyright (c) 2019-present, b3log.org\n"},
	{"1", "<meta charset='utf-8'><span>foo</span>", "foo\n"},
	{"0", "<html><body><!--StartFragment--><p>foo</p><!--EndFragment--></body></html>", "foo\n"},
}

func TestHTML2Md(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range html2MdTests {
		md := luteEngine.HTML2Md(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}
