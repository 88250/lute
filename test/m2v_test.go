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

	"github.com/88250/lute/ast"

	"github.com/88250/lute"
)

var md2VditorDOMTests = []parseTest{

	{"25", "<input autofocus>\n<img src=https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png onmouseleave=alert('XSS')>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;input autofocus&gt;\n&lt;img src=https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png onmouseleave=alert('XSS')&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><input autofocus=\"\">\n<img src=\"https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png\"></pre></div>"},
	{"24", "<form ><iframe/src=\"data:text/html,<script>alert('xss');</script>\"></iframe>", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;form &gt;&lt;iframe/src=&quot;data:text/html,&lt;script&gt;alert('xss');&lt;/script&gt;&quot;&gt;&lt;/iframe&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><form><iframe></iframe></pre></div>"},
	{"23", "[**foo**][bar]\n\n[bar]:https://github.com", "<p data-block=\"0\">\u200b<span data-type=\"link-ref\" data-link-label=\"bar\">foo</span>\u200b</p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[bar]: https://github.com\n</div>"},
	{"22", "<span class=\"vditor-comment\" data-cmtids=\"20201105091940-wtpsc3a\">foo</span>", "<p data-block=\"0\"><span class=\"vditor-comment\" data-cmtids=\"20201105091940-wtpsc3a\">foo</span></p>"},
	{"21", "<span class=\"vditor-comment\" data-cmtids=\"20201105091940-wtpsc3a\">foo</span>b", "<p data-block=\"0\"><span class=\"vditor-comment\" data-cmtids=\"20201105091940-wtpsc3a\">foo</span>b</p>"},
	{"20", "\n      > foo", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre class=\"vditor-wysiwyg__pre\" style=\"display: none\"><code>  &gt; foo\n</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><code>  &gt; foo\n</code></pre></div>"},
	{"19", "foo\n{: id=\"fooid\"}\nbar\n{: id=\"barid\"}", "<p data-block=\"0\" id=\"fooid\">foo</p><p data-block=\"0\" id=\"barid\">bar</p>"},
	{"18", "![][foo]\n\n[foo]: bar", "<p data-block=\"0\">\u200b<img src=\"bar\" alt=\"\" data-type=\"link-ref\" data-link-label=\"foo\" /></p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[foo]: bar\n</div>"},
	{"17", "![text][foo]\n\n[foo]: bar", "<p data-block=\"0\">\u200b<img src=\"bar\" alt=\"text\" data-type=\"link-ref\" data-link-label=\"foo\" /></p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[foo]: bar\n</div>"},
	{"16", "# heading {#custom-id}\n", "<h1 data-block=\"0\" data-id=\"#custom-id\" id=\"wysiwyg-#custom-id\" data-marker=\"#\">heading</h1>"},
	{"15", "foo\n\n[^1]: 111\n\n[2]: 222\n", "<p data-block=\"0\">foo</p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"><p data-block=\"0\">111</p></li></ol></div><div data-block=\"0\" data-type=\"link-ref-defs-block\">[2]: 222\n</div>"},
	{"14", "[^1]\n\n[^1]:\n", "<p data-block=\"0\">\u200b<sup data-type=\"footnotes-ref\" data-footnotes-label=\"^1\" class=\"vditor-tooltipped vditor-tooltipped__s\" aria-label=\"\">1</sup>\u200b</p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"></li></ol></div>"},
	{"13", "[toc]\n\n# foo", "<div class=\"vditor-toc\" data-block=\"0\" data-type=\"toc-block\" contenteditable=\"false\"><ul><li><span data-target-id=\"wysiwyg-foo\">foo</span></li></ul></div><h1 data-block=\"0\" id=\"wysiwyg-foo\" data-marker=\"#\">foo</h1>"},
	{"12", "foo[^1]\n[^1]:bar\n    * baz", "<p data-block=\"0\">foo<sup data-type=\"footnotes-ref\" data-footnotes-label=\"^1\" class=\"vditor-tooltipped vditor-tooltipped__s\" aria-label=\"barbaz\">1</sup>\u200b</p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"><p data-block=\"0\">bar</p><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">baz</li></ul></li></ol></div>"},
	{"11", "[foo][1]\n\n[1]: /bar\n", "<p data-block=\"0\">\u200b<span data-type=\"link-ref\" data-link-label=\"1\">foo</span>\u200b</p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[1]: /bar\n</div>"},
	{"10", "Foo\n    ---\n", "<p data-block=\"0\">Foo\n---</p>"},
	{"9", "    ***\n     ***\n\n-     -      -      -", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre class=\"vditor-wysiwyg__pre\" style=\"display: none\"><code>***\n ***\n</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><code>***\n ***\n</code></pre></div><hr data-block=\"0\" />"},
	{"8", "    ***\n", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre class=\"vditor-wysiwyg__pre\" style=\"display: none\"><code>***\n</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><code>***\n</code></pre></div>"},
	{"7", "* a\n  * b", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">a<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">b</li></ul></li></ul>"},
	{"6", "[]()", "<p data-block=\"0\">[]()</p>"},
	{"5", "* [ ]", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">[ ]</li></ul>"},
	{"4", "*", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">\u200b</li></ul>"},
	{"3", "foo'%'bar", "<p data-block=\"0\">foo'%'bar</p>"},
	{"2", "<p align=\"center\">\nfoo</p>\n\nbar", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;p align=&quot;center&quot;&gt;\nfoo&lt;/p&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><p align=\"center\">\nfoo</p></pre></div><p data-block=\"0\">bar</p>"},
	{"1", `foo\<aa>bar`, "<p data-block=\"0\">foo<span data-type=\"backslash\"><span>\\</span>&lt;</span>aa&gt;bar</p>"},
	{"0", `<details>
<summary>foo</summary>

* bar

</details>`, "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;details&gt;\n&lt;summary&gt;foo&lt;/summary&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><details>\n<summary>foo</summary></pre></div><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;/details&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"></details></pre></div>"},
}

func TestMd2VditorDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetHeadingID(true)
	luteEngine.SetVditorWYSIWYG(true)
	luteEngine.ParseOptions.ToC = true
	luteEngine.RenderOptions.ToC = true
	luteEngine.ParseOptions.KramdownBlockIAL = true
	luteEngine.RenderOptions.KramdownBlockIAL = true
	luteEngine.RenderOptions.Sanitize = true

	for _, test := range md2VditorDOMTests {
		md := luteEngine.Md2VditorDOM(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}

var md2VditorIRDOMTests = []parseTest{

	{"16", "<img src=\"https://foo\"/>", "<div data-block=\"0\" data-type=\"html-block\" class=\"vditor-ir__node\"><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code data-type=\"html-block\">&lt;img src=&quot;https://foo&quot;/&gt;</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><img src=\"https://foo\"/></pre></div>"},
	{"15", "<img src=\"https://foo?bar=baz\"/>", "<div data-block=\"0\" data-type=\"html-block\" class=\"vditor-ir__node\"><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code data-type=\"html-block\">&lt;img src=&quot;https://foo?bar=baz&quot;/&gt;</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><img src=\"https://foo?bar=baz\"/></pre></div>"},
	{"14", "foo\n`bar`\n", "<p data-block=\"0\">foo\n<span data-type=\"code\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">`</span><code data-newline=\"1\">bar</code><span class=\"vditor-ir__marker\">`</span></span></p>"},
	{"13", "foo\n{: id=\"fooid\"}\nbar\n{: id=\"barid\"}", "<p data-block=\"0\">foo</p><p data-block=\"0\">bar</p>"},
	{"12", "![][foo]\n\n[foo]: bar", "<p data-block=\"0\"><span class=\"vditor-ir__node\" data-type=\"img\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">[foo]</span><img src=\"bar\" /></span></p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[foo]: bar\n</div>"},
	{"11", "![text][foo]\n\n[foo]: bar", "<p data-block=\"0\"><span class=\"vditor-ir__node\" data-type=\"img\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">text</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">[foo]</span><img src=\"bar\" alt=\"text\" /></span></p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[foo]: bar\n</div>"},
	{"10", "[^foo]\n\n[^foo]:", "<p data-block=\"0\">\u200b<sup data-type=\"footnotes-ref\" class=\"vditor-ir__node vditor-tooltipped vditor-tooltipped__s\" aria-label=\"\" data-footnotes-label=\"^foo\"><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">^foo</span><span class=\"vditor-ir__marker--hide\" data-render=\"1\">1</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span></sup>\u200b</p><div data-block=\"0\" data-type=\"footnotes-block\"><div data-type=\"footnotes-def\">[^foo]: </div></div>"},
	{"9", "# foo {id}", "<h1 data-block=\"0\" class=\"vditor-ir__node\" id=\"ir-id\" data-marker=\"#\"><span class=\"vditor-ir__marker vditor-ir__marker--heading\" data-type=\"heading-marker\"># </span>foo<span data-type=\"heading-id\" class=\"vditor-ir__marker\"> {id}</span></h1>"},
	{"8", "`foo`", "<p data-block=\"0\"><span data-type=\"code\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">`</span><code data-newline=\"1\">foo</code><span class=\"vditor-ir__marker\">`</span></span></p>"},
	{"7", "$$\nfoo\n$$", "<div data-block=\"0\" data-type=\"math-block\" class=\"vditor-ir__node\"><span data-type=\"math-block-open-marker\">$$</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code data-type=\"math-block\" class=\"language-math\">foo</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><div data-type=\"math-block\" class=\"language-math\">foo</div></pre><span data-type=\"math-block-close-marker\">$$</span></div>"},
	{"6", "foo<bar>baz", "<p data-block=\"0\">foo<span data-type=\"html-inline\" class=\"vditor-ir__node\"><code class=\"vditor-ir__marker\">&lt;bar&gt;</code></span>baz</p>"},
	{"5", "$foo$", "<p data-block=\"0\"><span data-type=\"inline-node\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">$</span><code data-newline=\"1\" class=\"vditor-ir__marker vditor-ir__marker--pre\" data-type=\"math-inline\">foo</code><span class=\"vditor-ir__preview\" data-render=\"2\"><span class=\"language-math\">foo</span></span><span class=\"vditor-ir__marker\">$</span></span></p>"},
	{"4", "```foo\nbar\n```", "<div data-block=\"0\" data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">\u200bfoo</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code class=\"language-foo\">bar\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><code class=\"language-foo\">bar\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div>"},
	{"3", "__foo__", "<p data-block=\"0\"><span data-type=\"strong\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bi\">__</span><strong data-newline=\"1\">foo</strong><span class=\"vditor-ir__marker vditor-ir__marker--bi\">__</span></span></p>"},
	{"2", "* foo\n  * bar", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">foo<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></li></ul>"},
	{"1", "# foo", "<h1 data-block=\"0\" class=\"vditor-ir__node\" id=\"ir-foo\" data-marker=\"#\"><span class=\"vditor-ir__marker vditor-ir__marker--heading\" data-type=\"heading-marker\"># </span>foo</h1>"},
	{"0", "*foo*", "<p data-block=\"0\"><span data-type=\"em\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span><em data-newline=\"1\">foo</em><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span></span></p>"},
}

func TestMd2VditorIRDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.ToC = true
	luteEngine.RenderOptions.ToC = true
	luteEngine.ParseOptions.KramdownBlockIAL = true
	luteEngine.RenderOptions.KramdownBlockIAL = true

	ast.Testing = true
	for _, test := range md2VditorIRDOMTests {
		md := luteEngine.Md2VditorIRDOM(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}
