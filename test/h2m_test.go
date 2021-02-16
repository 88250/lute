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

	{"40", "<!--StartFragment--><strong>foo.</strong><span>bar</span><!--EndFragment-->", "**foo.** bar\n"},
	{"39", "<!--StartFragment--><p><strong>Js版</strong></p><pre>&lt;script&gt;\n&nbsp;&nbsp;&nbsp;&nbsp; test = \"你好abc\"\n&nbsp;&nbsp;&nbsp;&nbsp; str = \"\"\n&nbsp;&nbsp;&nbsp;&nbsp; for( i=0;&nbsp;&nbsp;&nbsp; i&lt;test.length; i++ )\n&nbsp;&nbsp;&nbsp;&nbsp; {\n&nbsp;&nbsp;&nbsp;&nbsp;  temp = test.charCodeAt(i).toString(16);\n&nbsp;&nbsp;&nbsp;&nbsp;  str&nbsp;&nbsp;&nbsp; += \"\\\\u\"+ new Array(5-String(temp).length).join(\"0\") +temp;\n&nbsp;&nbsp;&nbsp;&nbsp; }\n&nbsp;&nbsp;&nbsp;&nbsp; document.write (str)\n&lt;/script&gt;</pre><br><!--EndFragment-->", "**Js 版**\n\n```\n<script>\n\u00a0\u00a0\u00a0\u00a0 test = \"你好abc\"\n\u00a0\u00a0\u00a0\u00a0 str = \"\"\n\u00a0\u00a0\u00a0\u00a0 for( i=0;\u00a0\u00a0\u00a0 i<test.length; i++ )\n\u00a0\u00a0\u00a0\u00a0 {\n\u00a0\u00a0\u00a0\u00a0  temp = test.charCodeAt(i).toString(16);\n\u00a0\u00a0\u00a0\u00a0  str\u00a0\u00a0\u00a0 += \"\\\\u\"+ new Array(5-String(temp).length).join(\"0\") +temp;\n\u00a0\u00a0\u00a0\u00a0 }\n\u00a0\u00a0\u00a0\u00a0 document.write (str)\n</script>\n```\n"},
	{"38", "<!--StartFragment--><table width=\"778\"><tbody><tr><td class=\"key\">ú</td><td>&amp;uacute;</td><td>&amp;#250;</td><td class=\"key\">û</td><td>&amp;ucirc;</td><td>&amp;#251;</td><td class=\"key\">ü</td><td>&amp;uuml;</td><td>&amp;#252;</td><td class=\"key\">ý</td><td>&amp;yacute;</td><td>&amp;#253;</td><td class=\"key\">þ</td><td>&amp;thorn;</td><td>&amp;#254;</td></tr><tr><td class=\"key\">ÿ</td><td>&amp;yuml;</td></tr></tbody></table><!--EndFragment-->\n", "| ú | &uacute; | &#250; | û | &ucirc; | &#251; | ü | &uuml; | &#252; | ý | &yacute; | &#253; | þ | &thorn; | &#254; |\n| ---- | ---------- | -------- | ---- | --------- | -------- | ---- | -------- | -------- | ---- | ---------- | -------- | ---- | --------- | -------- |\n| ÿ | &yuml;   |\n"},
	{"37", "<!--StartFragment--><table width=\"400\"><tbody><tr><th>显示</th><th>说明</th><th>实体名称</th><th>实体编号</th></tr><tr><td class=\"key\"></td><td>半方大的空白</td><td>&amp;ensp;</td><td>&amp;#8194;</td></tr><tr></tr><tr><td class=\"key\"></td><td>全方大的空白</td><td>&amp;emsp;</td><td>&amp;#8195;</td></tr><tr></tr><tr><td class=\"key\"></td><td>不断行的空白格</td><td>&amp;nbsp;</td><td>&amp;#160;</td></tr><tr><td class=\"key\">&lt;</td><td>小于</td><td>&amp;lt;</td><td>&amp;#60;</td></tr><tr><td class=\"key\">&gt;</td><td>大于</td><td>&amp;gt;</td><td>&amp;#62;</td></tr><tr><td class=\"key\">&amp;</td><td>&amp;符号</td><td>&amp;amp;</td><td>&amp;#38;</td></tr><tr><td class=\"key\">\"</td><td>双引号</td><td>&amp;quot;</td><td>&amp;#34;</td></tr><tr><td class=\"key\">©</td><td>版权</td><td>&amp;copy;</td><td>&amp;#169;</td></tr><tr><td class=\"key\">®</td><td>已注册商标</td><td>&amp;reg;</td><td>&amp;#174;</td></tr><tr><td class=\"key\">™</td><td>商标（美国）</td><td>™</td><td>&amp;#8482;</td></tr><tr></tr><tr><td class=\"key\">×</td><td>乘号</td><td>&amp;times;</td><td>&amp;#215;</td></tr><tr><td class=\"key\">÷</td><td>除号</td><td>&amp;divide;</td><td>&amp;#247;</td></tr></tbody></table><!--EndFragment-->\n", "| 显示 | 说明           | 实体名称 | 实体编号 |\n| ------ | ---------------- | ---------- | ---------- |\n|      | 半方大的空白   | &ensp;   | &#8194;  |\n|      | 全方大的空白   | &emsp;   | &#8195;  |\n|      | 不断行的空白格 | &nbsp;   | &#160;   |\n| <    | 小于           | &lt;     | &#60;    |\n| >    | 大于           | &gt;     | &#62;    |\n| &    | &符号          | &amp;    | &#38;    |\n| \"    | 双引号         | &quot;   | &#34;    |\n| ©   | 版权           | &copy;   | &#169;   |\n| ®   | 已注册商标     | &reg;    | &#174;   |\n| ™   | 商标（美国）   | ™       | &#8482;  |\n| ×   | 乘号           | &times;  | &#215;   |\n| ÷   | 除号           | &divide; | &#247;   |\n"},
	{"36", "<!--StartFragment--><h1><b><span>foo</span></b><b><span><o:p></o:p></span></b></h1><p class=\"MsoNormal\"><img width=\"554\" height=\"337\" src=\"file:///C:\\WINDOWS\\TEMP\\ksohtml15220\\wps4.jpg\"><span><o:p>&nbsp;</o:p></span></p><!--EndFragment-->", "# **foo**\n\n![](file:///C:\\WINDOWS\\TEMP\\ksohtml15220\\wps4.jpg)\n"},
	{"35", "<a href=\"bar\">&lt;foo&gt;</a>", "[&lt;foo&gt;](bar)\n"},
	{"34", "<div class=\"gatsby-highlight\" data-language=\"js\"><pre class=\"blog-code language-js\"><span class=\"token keyword\">const</span></pre></div>", "```js\nconst\n```\n"},
	{"33", "<table><tr><td><p>事件编号</p></td><td><p>事件类别(category)</p></td><td><p>事件操作(action)</p></td><td><p>事件标签(label)</p></td><td><p>事件值(value)</p></td></tr><tr><td><p>1</p></td><td><p>合作行业标签</p></td><td><p>点击</p></td><td><p>选择条件</p></td><td></td></tr><tr><td><p>2</p></td><td><p>营销目的标签</p></td><td><p>点击</p></td><td><p>选择条件</p></td><td></td></tr><tr><td><p>3</p></td><td><p>合作资源标签</p></td><td><p>点击</p></td><td><p>选择条件</p></td><td></td></tr><tr><td><p>4</p></td><td><p>合作平台标签</p></td><td><p>点击</p></td><td><p>选择条件</p></td><td></td></tr><tr><td><p>5</p></td><td><p>卡片</p></td><td><p>查看详情</p></td><td><p>案例名称</p></td><td></td></tr><tr><td><p>6</p></td><td><p>卡片</p></td><td><p>点赞</p></td><td><p>案例名称</p></td><td><p>点赞数</p></td></tr><tr><td><p>7</p></td><td><p>卡片</p></td><td><p>取消点赞</p></td><td><p>案例名称</p></td><td><p>点赞数</p></td></tr><tr><td><p>8</p></td><td><p>卡片</p></td><td><p>下载分享图</p></td><td><p>案例名称</p></td><td></td></tr></table>", "| 事件编号 | 事件类别(category) | 事件操作(action) | 事件标签(label) | 事件值(value) |\n| ---------- | -------------------- | ------------------ | ----------------- | --------------- |\n| 1        | 合作行业标签       | 点击             | 选择条件        |               |\n| 2        | 营销目的标签       | 点击             | 选择条件        |               |\n| 3        | 合作资源标签       | 点击             | 选择条件        |               |\n| 4        | 合作平台标签       | 点击             | 选择条件        |               |\n| 5        | 卡片               | 查看详情         | 案例名称        |               |\n| 6        | 卡片               | 点赞             | 案例名称        | 点赞数        |\n| 7        | 卡片               | 取消点赞         | 案例名称        | 点赞数        |\n| 8        | 卡片               | 下载分享图       | 案例名称        |               |\n"},
	{"32", "<ul>\n  <li>咖啡</li>\n  <li>茶\n    <ul>\n    <li>红茶</li>\n    <li>绿茶</li>\n    </ul>\n  </li>\n  <li>牛奶</li>\n</ul>", "* 咖啡\n* 茶\n  * 红茶\n  * 绿茶\n* 牛奶\n"},
	{"32", "<ul>\n<li>foo</li>\n<li>bar\n<ul>\n<li>baz</li>\n<li>baz</li>\n</ul>\n</li>\n<li>bar</li>\n</ul>", "* foo\n* bar\n  * baz\n  * baz\n* bar\n"},
	{"31", "<ul>\n<li>foo\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>", "* foo\n  * bar\n"},
	{"30", "<ul><li>foo</li><li>bar<ul><li>baz</li><li>baz</li></ul></li><li>bar</li></ul>", "* foo\n* bar\n  * baz\n  * baz\n* bar\n"},
	{"29", "<p>测试 <code>name</code> 属性。调用 <code>getValue</code> 时 <code>Method</code></p>", "测试 `name` 属性。调用 `getValue` 时 `Method`\n"},
	{"28", `<html>
<body>
	<table>
		<tr>
      	<th>Month</th>
          <th>Savings</th>
		</tr>
		<tr>
      	<td>January</td>
          <td>$100</td>
		</tr>
		<tr>
      	<td>February</td>
          <td>$80</td>
		</tr>
	</table>
</body>
</html>`, "| Month    | Savings |\n| ---------- | --------- |\n| January  | $100    |\n| February | $80     |\n"},
	{"27", `<html>
<body>
    <table>
            <thead>
                    <tr>
                            <th>Month</th>
                            <th>Savings</th>
                    </tr>
				</thead>
            <tbody>
                    <tr>
                            <td>January</td>
                            <td>$100</td>
                    </tr>
                    <tr>
                            <td>February</td>
                            <td>$80</td>
                    </tr>
            </tbody>
    </table>
</body>
</html>`, "| Month    | Savings |\n| ---------- | --------- |\n| January  | $100    |\n| February | $80     |\n"},
	{"26", "<table class=\"markdown-reference\"><thead><tr><th>Type</th><th class=\"second-example\">Or</th><th>… to Get</th></tr></thead><tbody><tr><td class=\"preformatted\">*Italic*</td><td class=\"preformatted second-example\">_Italic_</td><td><em>Italic</em></td></tr><tr><td class=\"preformatted\">**Bold**</td><td class=\"preformatted second-example\">__Bold__</td><td><strong>Bold</strong></td></tr><tr><td class=\"preformatted\"># Heading 1</td><td class=\"preformatted second-example\">Heading 1<br>=========</td><td><h1 class=\"smaller-h1\">Heading 1</h1></td></tr><tr><td class=\"preformatted\">## Heading 2</td><td class=\"preformatted second-example\">Heading 2<br>---------</td><td><h2 class=\"smaller-h2\">Heading 2</h2></td></tr><tr><td class=\"preformatted\">[Link](http://a.com)</td><td class=\"preformatted second-example\">[Link][1]<br>⋮<br>[1]: http://b.org</td><td><a href=\"https://commonmark.org/\">Link</a></td></tr><tr><td class=\"preformatted\">![Image](http://url/a.png)</td><td class=\"preformatted second-example\">![Image][1]<br>⋮<br>[1]: http://url/b.jpg</td><td><img src=\"https://commonmark.org/help/images/favicon.png\" width=\"36\" height=\"36\" alt=\"Markdown\"></td></tr><tr><td class=\"preformatted\">&gt; Blockquote</td><td class=\"preformatted second-example\">&nbsp;</td><td><blockquote>Blockquote</blockquote></td></tr><tr><td class=\"preformatted\"><p>* List<br>* List<br>* List</p></td><td class=\"preformatted second-example\"><p>- List<br>- List<br>- List<br></p></td><td><ul><li>List</li><li>List</li><li>List</li></ul></td></tr></tbody></table>", "| Type                       | Or                                   | … to Get                                                 |\n| ---------------------------- | -------------------------------------- | ----------------------------------------------------------- |\n| *Italic*                   | _Italic_                             | *Italic*                                                |\n| **Bold**                   | __Bold__                             | **Bold**                                            |\n| # Heading 1                | Heading 1<br/>=========                  | # Heading 1                                              |\n| ## Heading 2               | Heading 2<br/>---------                  | ## Heading 2                                             |\n| [Link](http://a.com)       | [Link][1]<br/>⋮<br/>[1]: http://b.org       | [Link](https://commonmark.org/)                              |\n| ![Image](http://url/a.png) | ![Image][1]<br/>⋮<br/>[1]: http://url/b.jpg | ![Markdown](https://commonmark.org/help/images/favicon.png) |\n| > Blockquote               |                                      | > Blockquote                                     |\n| * List<br/>* List<br/>* List       | - List<br/>- List<br/>- List<br/>                | * List* List* List                                      |\n"},
	{"25", "<table class=\"table table-bordered\"><thead class=\"thead-light\"><tr><th>Element</th><th>Markdown Syntax</th></tr></thead><tbody><tr><td><a href=\"https://www.markdownguide.org/extended-syntax/#tables\">Table</a></td><td><code>| Syntax | Description |<br>| ----------- | ----------- |<br>| Header | Title |<br>| Paragraph | Text |</code></td></tr><tr><td><a href=\"https://www.markdownguide.org/extended-syntax/#fenced-code-blocks\">Fenced Code Block</a></td><td><code>```<br>{<br>&nbsp;&nbsp;\"firstName\": \"John\",<br>&nbsp;&nbsp;\"lastName\": \"Smith\",<br>&nbsp;&nbsp;\"age\": 25<br>}<br>```</code></td></tr></tbody></table>", "| Element                                                                             | Markdown Syntax                                                                                                  |\n| ------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------ |\n| [Table](https://www.markdownguide.org/extended-syntax/#tables)                         | `\\| Syntax \\| Description \\|\\| ----------- \\| ----------- \\|\\| Header \\| Title \\|\\| Paragraph \\| Text \\|` |\n| [Fenced Code Block](https://www.markdownguide.org/extended-syntax/#fenced-code-blocks) | ````{\u00a0\u00a0\"firstName\": \"John\",\u00a0\u00a0\"lastName\": \"Smith\",\u00a0\u00a0\"age\": 25}````        |\n"},
	{"24", "<table><thead><tr><th>Element</th><th>Markdown Syntax</th></tr></thead><tbody><tr><td>Table</td><td><code>| Syntax | Description |<br>| ----------- | ----------- |<br>| Header | Title |<br>| Paragraph | Text |</code></td></tr></tbody></table>", "| Element | Markdown Syntax                                                                                                  |\n| --------- | ------------------------------------------------------------------------------------------------------------------ |\n| Table   | `\\| Syntax \\| Description \\|\\| ----------- \\| ----------- \\|\\| Header \\| Title \\|\\| Paragraph \\| Text \\|` |\n"},
	{"23", "<h2 style=\"box-sizing: border-box; margin-top: 24px; margin-bottom: 16px; font-weight: 600; font-size: 1.5em; line-height: 1.25; padding-bottom: 0.3em; border-bottom: 1px solid rgb(234, 236, 239); color: rgb(36, 41, 46); font-family: -apple-system, BlinkMacSystemFont, &quot;Segoe UI&quot;, Helvetica, Arial, sans-serif, &quot;Apple Color Emoji&quot;, &quot;Segoe UI Emoji&quot;; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; white-space: normal; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; background-color: rgb(255, 255, 255); text-decoration-style: initial; text-decoration-color: initial;\"><g-emoji class=\"g-emoji\" alias=\"m\" fallback-src=\"https://github.githubassets.com/images/icons/emoji/unicode/24c2.png\" style=\"box-sizing: border-box; font-family: &quot;Apple Color Emoji&quot;, &quot;Segoe UI&quot;, &quot;Segoe UI Emoji&quot;, &quot;Segoe UI Symbol&quot;; font-size: 1.2em; font-weight: 400; line-height: 20px; vertical-align: middle; font-style: normal !important;\">Ⓜ️</g-emoji><span> </span>Markdown User Guide</h2>", "## Ⓜ️ Markdown User Guide\n"},
	{"22", "<div class=\"highlight highlight-source-shell\"><pre>npm install vditor --save</pre></div>", "```shell\nnpm install vditor --save\n```\n"},
	{"21", "<h4><a id=\"user-content-id\" class=\"anchor\" aria-hidden=\"true\" href=\"https://github.com/Vanessa219/vditor/blob/master/README.md#id\"><svg class=\"octicon octicon-link\" viewBox=\"0 0 16 16\" version=\"1.1\" width=\"16\" height=\"16\" aria-hidden=\"true\"><path fill-rule=\"evenodd\" d=\"M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z\"></path></svg></a>id</h4>", "#### id\n"},
	{"20", "<h2 id=\"whats-markdown\">What’s Markdown?<a class=\"anchorjs-link \" aria-label=\"Anchor\" data-anchorjs-icon=\"\uE9CB\" href=\"https://www.markdownguide.org/getting-started/#whats-markdown\"></a></h2>", "## What’s Markdown?\n"},
	{"19", "<pre><span>`foo`</span></pre>", "```\n`foo`\n```\n"},
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
	{"3", "<pre><ul><li>foo</li></ul></pre>", "<pre><ul><li><p>foo</p></li></ul></pre>\n"},
	{"2", "<pre><span>//&#32;Lute&#32;-&#32;A&#32;structured&#32;markdown&#32;engine.<br></span><span>//&#32;Copyright&#32;(c)&#32;2019-present,&#32;b3log.org</span></pre>", "```\n// Lute - A structured markdown engine.\n// Copyright (c) 2019-present, b3log.org\n```\n"},
	{"1", "<meta charset='utf-8'><span>foo</span>", "foo\n"},
	{"0", "<html><body><!--StartFragment--><p>foo</p><!--EndFragment--></body></html>", "foo\n"},
}

func TestHTML2Md(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetAutoSpace(true)
	for _, test := range html2MdTests {
		md := luteEngine.HTML2Md(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}
