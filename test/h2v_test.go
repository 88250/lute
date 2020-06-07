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

var html2VditorDOMTests = []parseTest{

	{"4", `<iframe src="foo" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true"> </iframe>`, "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;iframe src=&quot;foo&quot; scrolling=&quot;no&quot; border=&quot;0&quot; frameborder=&quot;no&quot; framespacing=&quot;0&quot; allowfullscreen=&quot;true&quot;&gt; &lt;/iframe&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><iframe src=\"foo\" scrolling=\"no\" border=\"0\" frameborder=\"no\" framespacing=\"0\" allowfullscreen=\"true\"> </iframe></pre></div>"},
	{"3", `<!--StartFragment--><a href="https://hacpai.com/article/1553314676872?r=Vanessa">每天 30 秒系列</a><!--EndFragment-->`, "<p data-block=\"0\"><a href=\"https://hacpai.com/article/1553314676872?r=Vanessa\">每天 30 秒系列</a>\n</p>"},
	{"2", `<!--StartFragment--><span>Use<span>&nbsp;</span></span><code class="language-text">new Date()</code><span><span>&nbsp;</span>and</span><!--EndFragment-->`, "<p data-block=\"0\">Use <code data-marker=\"`\">\u200bnew Date()</code>\u200b and\n</p>"},
	{"1", `<pre><code class="language-text">&gt;</code></pre>`, "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre class=\"vditor-wysiwyg__pre\"><code class=\"language-text\">&gt;\n</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><code class=\"language-text\">&gt;\n</code></pre></div>"},
	{"0", `<code class="language-text">&gt;</code>`, "<p data-block=\"0\">\u200b<code data-marker=\"`\">\u200b&gt;</code>\u200b\n</p>"},
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
