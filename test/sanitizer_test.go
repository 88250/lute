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
	"github.com/88250/lute/render"
)

var sanitizerTests = []parseTest{

	{"16", "<form><input formaction=javascript:alert('xss') type=submit value='click me'></input></form>", "<form><input type=\"submit\" value=\"click me\"></input></form>\n"},
	{"15", "<meta http-equiv=\"refresh\" content=\"0\" />", "<meta content=\"0\"/>\n"},
	{"14", "<div><p /><a href=\"javascript:alert('xss')\">click me</a></div>", "<div><p/><a>click me</a></div>\n"},
	{"13", "<iframe src=\"data&NewLine;:text/html,%3Cscript%3Ealert('xss')%3C%2Fscript%3E\"></iframe>", "<iframe></iframe>\n"},
	{"12", "<iframe src=\"data&Tab;:text/html,%3Cscript%3Ealert('xss')%3C%2Fscript%3E\"></iframe>", "<iframe></iframe>\n"},
	{"11", "<iframe src=\"java&NewLine;script:alert('xss')\"></iframe>", "<iframe></iframe>\n"},
	{"10", "<iframe src=\"java&Tab;script:alert('xss')\"></iframe>", "<iframe></iframe>\n"},
	{"9", "<iframe srcdoc=\"&lt;img src&equals;x onerror&equals;alert&lpar;'xss'&rpar;&gt;\" />", "<iframe/>\n"},
	{"8", "[xss](javascript:alert(document.domain))", "<p><a href=\"\">xss</a></p>\n"},
	{"7", "![a](\"<img src=xss onerror=alert(1)>)\n", "<p>![a](&quot;&lt;img src=xss onerror=alert(1)&gt;)</p>\n"},
	{"6", "<img src=\"foo\" onload=\"alert(1)\" onerror=\"alert(2)\"/>", "<img src=\"foo\" />\n"},
	{"5", "<iframe src='javascript:parent.require(\"child_process\").exec(\"open -a Calculator\")'></iframe>", "<iframe></iframe>\n"},
	{"4", "![Escape SRC - onerror](\"onerror=\"alert('ImageOnError'))", "<p><img src=\"%22onerror=%22alert(&#39;ImageOnError&#39;)\" alt=\"Escape SRC - onerror\" /></p>\n"},
	{"3", "<EMBED SRC=\"data:image/svg+xml;base64,mock payload\" type=\"image/svg+xml\" AllowScriptAccess=\"always\"></EMBED>", "<p><embed type=\"image/svg+xml\" allowscriptaccess=\"always\"></embed></p>\n"},
	{"2", "<FOo>bar", "<p><foo>bar</p>\n"},
	{"1", "<img onerror=\"alert(1)\" src=\"bar.png\" />", "<img src=\"bar.png\" />\n"},
	{"0", "foo<script>alert(1)</script>bar", "<p>foo alert(1) bar</p>\n"},
}

func TestSanitizer(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.RenderOptions.Sanitize = true

	for _, test := range sanitizerTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var sanitizerVditorTests = []parseTest{

	{"8", "<img OnError=\"alert(1)\" src=\"bar.png\" />", "<p data-block=\"0\"><img src=\"bar.png\" alt=\"\" /></p>"},
	{"7", "<iframe src=\"//player.bilibili.com/player.html?aid=test&page=1\" scrolling=\"no\" border=\"0\" frameborder=\"no\" framespacing=\"0\" allowfullscreen=\"true\"> </iframe>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;iframe src=&quot;https://player.bilibili.com/player.html?aid=test&amp;page=1&quot; scrolling=&quot;no&quot; border=&quot;0&quot; frameborder=&quot;no&quot; framespacing=&quot;0&quot; allowfullscreen=&quot;true&quot;&gt;&lt;/iframe&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><iframe src=\"https://player.bilibili.com/player.html?aid=test&amp;page=1\" scrolling=\"no\" border=\"0\" frameborder=\"no\" framespacing=\"0\" allowfullscreen=\"true\"></iframe></pre></div>"},
	{"6", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;img src=\"test1<wbr>\" onerror=\"alert('XSS')\"&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"1\"><img src=\"test\" onerror=\"alert('XSS')\"></pre></div>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;img src=&quot;test1<wbr>&quot; onerror=&quot;alert('XSS')&quot;&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><img src=\"test1\"></pre></div>"},
	{"5", "<iframe src=\"javascript:parent.require('child_process').exec('open -a Calculator')\"></iframe>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;iframe src=&quot;javascript:parent.require('child_process').exec('open -a Calculator')&quot;&gt;&lt;/iframe&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><iframe></iframe></pre></div>"},
	{"4", "![Escape SRC - onerror](\"onerror=\"alert('ImageOnError'))", "<p data-block=\"0\"><img src=\"\" alt=\"Escape SRC - onerror\" /></p>"},
	{"3", "<EMBED SRC=\"data:image/svg+xml;base64,mock payload\" type=\"image/svg+xml\" AllowScriptAccess=\"always\"></EMBED>", "<p data-block=\"0\">\u200b<code data-type=\"html-inline\">\u200b&lt;embed src=&quot;data:image/svg+xml;base64,mock payload&quot; type=&quot;image/svg+xml&quot; allowscriptaccess=&quot;always&quot;/&gt;</code></p>"},
	{"2", "<FOo>bar", "<p data-block=\"0\">foobar</p>"},
	{"1", "<img onerror=\"alert(1)\" src=\"bar.png\" />", "<p data-block=\"0\"><img src=\"bar.png\" alt=\"\" /></p>"},
	{"0", "foo<script>alert(1)</script>bar", "<p data-block=\"0\">foo</p><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;script&gt;alert(1)&lt;/script&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\">  </pre></div><p data-block=\"0\">bar</p>"},
}

func TestSanitizerVditor(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetVditorWYSIWYG(true)
	luteEngine.RenderOptions.Sanitize = true

	for _, test := range sanitizerVditorTests {
		html := luteEngine.SpinVditorDOM(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

func TestSanitize(t *testing.T) {
	output := render.Sanitize("<img src=\"foo\" onload=\"alert(1)\" onerror=\"alert(2)\"/>")
	if "<img src=\"foo\" />" != output {
		t.Fatalf("sanitize failed")
	}

	output = render.Sanitize("![a](&quot;<img src=xss onerror=alert(1)>)\n")
	if "![a](&#34;<img src=\"xss\">)\n" != output {
		t.Fatalf("sanitize failed")
	}
}
