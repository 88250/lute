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
	"github.com/88250/lute/ast"
	"testing"

	"github.com/88250/lute"
)

var html2VditorDOMTests = []parseTest{

	{"5", `<!--StartFragment--><p class="MsoNormal"><span><font face="宋体">测试</font>WPS<font face="宋体">粘贴</font></span><span><o:p></o:p></span></p><!--EndFragment-->`, "<p data-block=\"0\">测试WPS粘贴</p>"},
	{"4", `<iframe src="foo" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"> </iframe>`, "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;iframe src=&quot;foo&quot; scrolling=&quot;no&quot; border=&quot;0&quot; frameborder=&quot;no&quot; framespacing=&quot;0&quot; allowfullscreen=&quot;true&quot;&gt; &lt;/iframe&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><iframe src=\"foo\" scrolling=\"no\" border=\"0\" frameborder=\"no\" framespacing=\"0\" allowfullscreen=\"true\"> </iframe></pre></div>"},
	{"3", `<!--StartFragment--><a href="https://ld246.com/article/1553314676872?r=Vanessa">每天 30 秒系列</a><!--EndFragment-->`, "<p data-block=\"0\"><a href=\"https://ld246.com/article/1553314676872?r=Vanessa\">每天 30 秒系列</a></p>"},
	{"2", `<!--StartFragment--><span>Use<span>&nbsp;</span></span><code class="language-text">new Date()</code><span><span>&nbsp;</span>and</span><!--EndFragment-->`, "<p data-block=\"0\">Use <code data-marker=\"`\">\u200bnew Date()</code>\u200b and</p>"},
	{"1", `<pre><code class="language-text">&gt;</code></pre>`, "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre class=\"vditor-wysiwyg__pre\" style=\"display: none\"><code class=\"language-text\">&gt;\n</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><code class=\"language-text\">&gt;\n</code></pre></div>"},
	{"0", `<code class="language-text">&gt;</code>`, "<p data-block=\"0\">\u200b<code data-marker=\"`\">\u200b&gt;</code>\u200b</p>"},
}

func TestHTML2VditorDOM(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range html2VditorDOMTests {
		result := luteEngine.HTML2VditorDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}

var html2VditorIRBlockDOMTests = []parseTest{
	{"4", "<!--StartFragment--><span class=\"tiao\">第一条</span><span>　为</span><!--EndFragment-->", "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\"><span data-type=\"strong\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--strong\">**</span><strong data-newline=\"1\">第一条</strong><span class=\"vditor-ir__marker vditor-ir__marker--strong\">**</span></span>\u3000为</p>"},
	{"3", `<!--StartFragment-->

<p class="MsoNormal"><span lang="EN-US">Para1<o:p></o:p></span></p>

<table class="MsoTableGrid" border="1" cellspacing="0" cellpadding="0" align="left">
 <tbody><tr>
  <td width="189" valign="top">
  <p class="MsoNormal"><span lang="EN-US">Foo<o:p></o:p></span></p>
  </td>
  <td width="189" valign="top">
  <p class="MsoNormal"><span lang="EN-US">Bar<o:p></o:p></span></p>
  </td>
  <td width="189" valign="top">
  <p class="MsoNormal"><span lang="EN-US">baz<o:p></o:p></span></p>
  </td>
 </tr>
 <tr>
  <td width="189" valign="top">
  <p class="MsoNormal"><span lang="EN-US">Foo2<o:p></o:p></span></p>
  <p class="MsoNormal"><span lang="EN-US">Foo3<o:p></o:p></span></p>
  </td>
  <td width="189" valign="top">
  <p class="MsoNormal"><span lang="EN-US">Bar2<o:p></o:p></span></p>
  </td>
  <td width="189" valign="top">
  <p class="MsoNormal"><span lang="EN-US">Baz2<o:p></o:p></span></p>
  </td>
 </tr>
</tbody></table>

<p class="MsoNormal"><span lang="EN-US"><o:p>&nbsp;</o:p></span></p>

<p class="MsoNormal"><span lang="EN-US">Para2<o:p></o:p></span></p>

<!--EndFragment-->`, "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">Para1</p><table data-type=\"table\" data-node-id=\"20060102150405-1a2b3c4\"><thead><tr><th>Foo</th><th>Bar</th><th>baz</th></tr></thead><tbody><tr><td>Foo2 Foo3</td><td>Bar2</td><td>Baz2</td></tr></tbody></table><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">Para2</p>"},
	{"2", `<table border="0" cellpadding="0" cellspacing="0" width="144">
<!--StartFragment-->
 <colgroup><col width="72" span="2">
 </colgroup><tbody><tr height="19">
  <td height="19" width="72">foo</td>
  <td width="72">bar</td>
 </tr>
 <tr height="19">
  <td height="19">baz</td>
  <td>bazz</td>
 </tr>
<!--EndFragment-->
</tbody></table>`, "<table data-type=\"table\" data-node-id=\"20060102150405-1a2b3c4\"><thead><tr><th>foo</th><th>bar</th></tr></thead><tbody><tr><td>baz</td><td>bazz</td></tr></tbody></table>"},
	{"1", `<!--StartFragment--><span>第9号</span><!--EndFragment-->`, "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">第9号</p>"},
	{"0", `<!--StartFragment-->
<table border="0" cellpadding="0" cellspacing="0" width="144" height="18">
 <colgroup><col width="72" span="2">
 </colgroup><tbody><tr height="18">
  <td height="18" width="72" x:str="">foo</td>
  <td width="72" x:str="">bar</td>
 </tr>
</tbody></table>
<!--EndFragment-->`, "<table data-type=\"table\" data-node-id=\"20060102150405-1a2b3c4\"><thead><tr><th>foo</th><th>bar</th></tr></thead></table>"},
	{"0", `<!--StartFragment--><a class="d-inline-block" data-hovercard-type="user" data-hovercard-url="/users/88250/hovercard" data-octo-click="hovercard-link-click" data-octo-dimensions="link_type:self" href="https://github.com/88250"><img class="avatar avatar-user" height="20" width="20" alt="@88250" src="https://avatars2.githubusercontent.com/u/873584?s=60&amp;u=f7f95251dd56b576aefd20094b6695a2db23a927&amp;v=4"></a><span><span>&nbsp;</span></span><!--EndFragment-->`, "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\"><span data-type=\"a\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__node\" data-type=\"img\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">@88250</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">https://avatars2.githubusercontent.com/u/873584?s=60&amp;u=f7f95251dd56b576aefd20094b6695a2db23a927&amp;v=4</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><img src=\"https://avatars2.githubusercontent.com/u/873584?s=60&u=f7f95251dd56b576aefd20094b6695a2db23a927&v=4\" alt=\"@88250\" /></span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">https://github.com/88250</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span></p>"},
}

func TestHTML2VditorIRBlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetAutoSpace(false)

	ast.Testing = true
	for _, test := range html2VditorIRBlockDOMTests {
		result := luteEngine.HTML2VditorIRBlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}
