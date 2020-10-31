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

	{"5", `<!--StartFragment--><p class="MsoNormal"><span><font face="宋体">测试</font>WPS<font face="宋体">粘贴</font></span><span><o:p></o:p></span></p><!--EndFragment-->`, "<p data-block=\"0\">测试 WPS 粘贴</p>"},
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

	{"0", `<!--StartFragment--><a class="d-inline-block" data-hovercard-type="user" data-hovercard-url="/users/88250/hovercard" data-octo-click="hovercard-link-click" data-octo-dimensions="link_type:self" href="https://github.com/88250"><img class="avatar avatar-user" height="20" width="20" alt="@88250" src="https://avatars2.githubusercontent.com/u/873584?s=60&amp;u=f7f95251dd56b576aefd20094b6695a2db23a927&amp;v=4"></a><span><span>&nbsp;</span></span><!--EndFragment-->`, "<p data-block=\"0\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\"><span data-type=\"a\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__node\" data-type=\"img\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">@88250</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">https://avatars2.githubusercontent.com/u/873584?s=60&u=f7f95251dd56b576aefd20094b6695a2db23a927&v=4</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><img src=\"https://avatars2.githubusercontent.com/u/873584?s=60&u=f7f95251dd56b576aefd20094b6695a2db23a927&v=4\" alt=\"@88250\" /></span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">https://github.com/88250</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span></p>"},
}

func TestHTML2VditorIRBlockDOM(t *testing.T) {
	luteEngine := lute.New()

	ast.Testing = true
	for _, test := range html2VditorIRBlockDOMTests {
		result := luteEngine.HTML2VditorIRBlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}

