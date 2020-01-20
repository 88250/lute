// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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

package test

import (
	"testing"

	"github.com/88250/lute"
)

var html2MdTests = []parseTest{

	{"11", "<li>foo</li><li>bar</li>", "* foo\n* bar\n"},
	{"10", `<p data-block="0">foo'%'bar</p>`, "foo'%'bar\n"},
	{"9", `<code class="language-text">&gt;</code>`, "`>`\n"},
	{"8", `<div><a href="/bar">foo</a></div>`, "[foo](/bar)\n"},
	{"7", `<ul><li><p>Java</p><ul><li><p>Spring</p></li></ul></li></ul>`, "* Java\n  * Spring\n"},
	{"6", `<!--StartFragment--><p>这是一篇讲解如何正确使用<span>&nbsp;</span><strong>Markdown</strong><span>&nbsp;</span>的排版示例，学会这个很有必要，能让你的文章有更佳清晰的排版。</p><!--EndFragment-->`, "这是一篇讲解如何正确使用 **Markdown** 的排版示例，学会这个很有必要，能让你的文章有更佳清晰的排版。\n"},
	{"5", `<!--StartFragment--><ul><li><input checked="" disabled="" type="checkbox"><span>&nbsp;</span>发布 Solo</li></ul><!--EndFragment-->`, "* [X] 发布 Solo\n"},
	{"4", "<span>&nbsp;</span>发布 Solo", "发布 Solo\n"},
	{"3", "<pre><ul><li>foo</li></ul></pre>", "<pre><ul><li>foo</li></ul></pre>\n"},
	{"2", "<pre><span>//&#32;Lute&#32;-&#32;A&#32;structured&#32;markdown&#32;engine.<br></span><span>//&#32;Copyright&#32;(c)&#32;2019-present,&#32;b3log.org</span></pre>", "<pre><span>// Lute - A structured markdown engine.<br/></span><span>// Copyright (c) 2019-present, b3log.org</span></pre>\n"},
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
