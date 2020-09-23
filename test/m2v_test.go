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

var md2VditorTests = []parseTest{

	{"19", "foo\n{: id=\"fooid\"}\nbar\n{: id=\"barid\"}", "<p data-block=\"0\" id=\"fooid\">foo</p><p data-block=\"0\" id=\"barid\">bar</p>"},
	{"18", "![][foo]\n\n[foo]: bar", "<p data-block=\"0\">\u200b<img src=\"bar\" alt=\"\" data-type=\"link-ref\" data-link-label=\"foo\" /></p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[foo]: bar\n</div>"},
	{"17", "![text][foo]\n\n[foo]: bar", "<p data-block=\"0\">\u200b<img src=\"bar\" alt=\"text\" data-type=\"link-ref\" data-link-label=\"foo\" /></p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[foo]: bar\n</div>"},
	{"16", "# heading {#custom-id}\n", "<h1 data-block=\"0\" data-id=\"#custom-id\" id=\"wysiwyg-#custom-id\" data-marker=\"#\">heading</h1>"},
	{"15", "foo\n\n[^1]: 111\n\n[2]: 222\n", "<p data-block=\"0\">foo</p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[2]: 222\n</div><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"><p data-block=\"0\">111</p></li></ol></div>"},
	{"14", "[^1]\n\n[^1]:\n", "<p data-block=\"0\">\u200b<sup data-type=\"footnotes-ref\" data-footnotes-label=\"^1\" class=\"b3-tooltips b3-tooltips__s\" aria-label=\"\">1</sup>\u200b</p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"></li></ol></div>"},
	{"13", "[toc]\n\n# foo", "<div class=\"vditor-toc\" data-block=\"0\" data-type=\"toc-block\" contenteditable=\"false\"><span data-type=\"toc-h\">foo</span><br></div><p data-block=\"0\"></p><h1 data-block=\"0\" id=\"wysiwyg-foo\" data-marker=\"#\">foo</h1>"},
	{"12", "foo[^1]\n[^1]:bar\n    * baz", "<p data-block=\"0\">foo<sup data-type=\"footnotes-ref\" data-footnotes-label=\"^1\" class=\"b3-tooltips b3-tooltips__s\" aria-label=\"barbaz\">1</sup>\u200b</p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"><p data-block=\"0\">bar</p><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">baz</li></ul></li></ol></div>"},
	{"11", "[foo][1]\n\n[1]: /bar\n", "<p data-block=\"0\">\u200b<span data-type=\"link-ref\" data-link-label=\"1\">foo</span>\u200b</p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[1]: /bar\n</div>"},
	{"10", "Foo\n    ---\n", "<p data-block=\"0\">Foo\n---</p>"},
	{"9", "    ***\n     ***\n\n-     -      -      -", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"***\n ***\n\"><pre class=\"vditor-wysiwyg__pre\" style=\"display: none\"><code>***\n ***\n</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><code>***\n ***\n</code></pre></div><hr data-block=\"0\" />"},
	{"8", "    ***\n", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"***\n\"><pre class=\"vditor-wysiwyg__pre\" style=\"display: none\"><code>***\n</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><code>***\n</code></pre></div>"},
	{"7", "* a\n  * b", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">a<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">b</li></ul></li></ul>"},
	{"6", "[]()", "<p data-block=\"0\">[]()</p>"},
	{"5", "* [ ]", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> </li></ul>"},
	{"4", "*", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">\u200b</li></ul>"},
	{"3", "foo'%'bar", "<p data-block=\"0\">foo'%'bar</p>"},
	{"2", "<p align=\"center\">\nfoo</p>\n\nbar", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;p align=&quot;center&quot;&gt;\nfoo&lt;/p&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><p align=\"center\">\nfoo</p></pre></div><p data-block=\"0\">bar</p>"},
	{"1", `foo\<aa>bar`, "<p data-block=\"0\">foo<span data-type=\"backslash\"><span>\\</span>&lt;</span>aa&gt;bar</p>"},
	{"0", `<details>
<summary>foo</summary>

* bar

</details>`, "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;details&gt;\n&lt;summary&gt;foo&lt;/summary&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"><details>\n<summary>foo</summary></pre></div><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;/details&gt;</code></pre><pre class=\"vditor-wysiwyg__preview\" data-render=\"2\"></details></pre></div>"},
}

func TestMd2Vditor(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ToC = true
	luteEngine.KramdownIAL = true

	for _, test := range md2VditorTests {
		md := luteEngine.Md2VditorDOM(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}

var md2VditorIRTests = []parseTest{

	{"13", "foo\n{: id=\"fooid\"}\nbar\n{: id=\"barid\"}", "<p data-block=\"0\">foo</p><span data-type=\"kramdown-ial\">{: id=\"fooid\"}\n</span><p data-block=\"0\">bar</p><span data-type=\"kramdown-ial\">{: id=\"barid\"}\n</span>"},
	{"12", "![][foo]\n\n[foo]: bar", "<p data-block=\"0\"><span class=\"vditor-ir__node\" data-type=\"img\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">[foo]</span><img src=\"bar\" /></span></p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[foo]: bar\n</div>"},
	{"11", "![text][foo]\n\n[foo]: bar", "<p data-block=\"0\"><span class=\"vditor-ir__node\" data-type=\"img\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">text</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">[foo]</span><img src=\"bar\" alt=\"text\" /></span></p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[foo]: bar\n</div>"},
	{"10", "[^foo]\n\n[^foo]:", "<p data-block=\"0\">\u200b<sup data-type=\"footnotes-ref\" class=\"vditor-ir__node b3-tooltips b3-tooltips__s\" aria-label=\"\" data-footnotes-label=\"^foo\"><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">^foo</span><span class=\"vditor-ir__marker--hide\" data-render=\"1\">1</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span></sup>\u200b</p><div data-block=\"0\" data-type=\"footnotes-block\"><div data-type=\"footnotes-def\">[^foo]: </div></div>"},
	{"9", "# foo {id}", "<h1 data-block=\"0\" class=\"vditor-ir__node\" id=\"ir-id\" data-marker=\"#\"><span class=\"vditor-ir__marker vditor-ir__marker--heading\" data-type=\"heading-marker\"># </span>foo<span data-type=\"heading-id\" class=\"vditor-ir__marker\"> {id}</span></h1>"},
	{"8", "`foo`", "<p data-block=\"0\"><span data-type=\"code\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">`</span><code data-newline=\"1\">foo</code><span class=\"vditor-ir__marker\">`</span></span></p>"},
	{"7", "$$\nfoo\n$$", "<div data-block=\"0\" data-type=\"math-block\" class=\"vditor-ir__node\"><span data-type=\"math-block-open-marker\">$$</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code data-type=\"math-block\" class=\"language-math\">foo</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><code data-type=\"math-block\" class=\"language-math\">foo</code></pre><span data-type=\"math-block-close-marker\">$$</span></div>"},
	{"6", "foo<bar>baz", "<p data-block=\"0\">foo<span data-type=\"html-inline\" class=\"vditor-ir__node\"><code class=\"vditor-ir__marker\">&lt;bar&gt;</code></span>baz</p>"},
	{"5", "$foo$", "<p data-block=\"0\"><span data-type=\"inline-node\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">$</span><code data-newline=\"1\" class=\"vditor-ir__marker vditor-ir__marker--pre\" data-type=\"math-inline\">foo</code><span class=\"vditor-ir__preview\" data-render=\"2\"><code class=\"language-math\">foo</code></span><span class=\"vditor-ir__marker\">$</span></span></p>"},
	{"4", "```foo\nbar\n```", "<div data-block=\"0\" data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">\u200bfoo</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code class=\"language-foo\">bar\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><code class=\"language-foo\">bar\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div>"},
	{"3", "__foo__", "<p data-block=\"0\"><span data-type=\"strong\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bi\">__</span><strong data-newline=\"1\">foo</strong><span class=\"vditor-ir__marker vditor-ir__marker--bi\">__</span></span></p>"},
	{"2", "* foo\n  * bar", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">foo<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></li></ul>"},
	{"1", "# foo", "<h1 data-block=\"0\" class=\"vditor-ir__node\" id=\"ir-foo\" data-marker=\"#\"><span class=\"vditor-ir__marker vditor-ir__marker--heading\" data-type=\"heading-marker\"># </span>foo</h1>"},
	{"0", "*foo*", "<p data-block=\"0\"><span data-type=\"em\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span><em data-newline=\"1\">foo</em><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span></span></p>"},
}

func TestMd2VditorIR(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ToC = true
	luteEngine.KramdownIAL = true

	for _, test := range md2VditorIRTests {
		md := luteEngine.Md2VditorIRDOM(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}

var md2VditorIRBlockTests = []parseTest{

	{"10", "foo#bar#baz", "<p data-block=\"0\" data-node-id=\"\" data-type=\"p\">foo<span data-type=\"tag\" class=\"vditor-ir__tag\">#bar#</span>baz</p>"},
	{"9", "foo\n((id \"text\"))bar", "<p data-block=\"0\" data-node-id=\"\" data-type=\"p\">foo\n<span data-type=\"block-ref\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">id</span><span class=\"vditor-ir__marker\"> </span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__blockref\">text</span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span>bar</p>"},
	{"8", "foo\n!((id \"text\"))\nbar", "<p data-block=\"0\" data-node-id=\"\" data-type=\"p\">foo</p><div data-block=\"0\" data-node-id=\"\" data-type=\"block-ref-embed\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">id</span><span class=\"vditor-ir__marker\"> </span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__blockref\">text</span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><div data-block-def-id=\"id\" data-render=\"2\" data-type=\"block-render\"></div></div><p data-block=\"0\" data-node-id=\"\" data-type=\"p\">bar</p>"},
	{"7", "foo!((id \"text\"))", "<p data-block=\"0\" data-node-id=\"\" data-type=\"p\">foo!<span data-type=\"block-ref\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">id</span><span class=\"vditor-ir__marker\"> </span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__blockref\">text</span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span></p>"},
	{"7", "!((id))", "<div data-block=\"0\" data-node-id=\"\" data-type=\"block-ref-embed\" data-text=\"0\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">id</span><span class=\"vditor-ir__marker\"> </span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__blockref\"></span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><div data-block-def-id=\"id\" data-render=\"2\" data-type=\"block-render\"></div></div>"},
	{"6", "!((id \"text\"))", "<div data-block=\"0\" data-node-id=\"\" data-type=\"block-ref-embed\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">id</span><span class=\"vditor-ir__marker\"> </span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__blockref\">text</span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><div data-block-def-id=\"id\" data-render=\"2\" data-type=\"block-render\"></div></div>"},
	{"5", "<!--foo-->\nbar", "<div data-block=\"0\" data-node-id=\"\" data-type=\"html-block\"><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code data-type=\"html-block\">&lt;!--foo--&gt;</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><!--foo--></pre></div><p data-block=\"0\" data-node-id=\"\" data-type=\"p\">bar</p>"},
	{"4", "# foo\n{: id=\"test\" bookmark=\"bookmark\"}", "<h1 data-block=\"0\" class=\"vditor-ir__node\" data-node-id=\"test\" bookmark=\"bookmark\" data-type=\"h\" id=\"ir-foo\" data-marker=\"#\"><span class=\"vditor-ir__marker vditor-ir__marker--heading\" data-type=\"heading-marker\"># </span>foo</h1>"},
	{"3", "```\nfoo\n```\n{: id=\"test\" bookmark=\"bookmark\"}", "<div data-block=\"0\" data-node-id=\"test\" bookmark=\"bookmark\" data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">\u200b</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code>foo\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><code>foo\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div>"},
	{"2", "```\nfoo\n```\n{: id=\"test\"}", "<div data-block=\"0\" data-node-id=\"test\" data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">\u200b</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code>foo\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><code>foo\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div>"},
	{"1", "foo\n{: id=\"fooid\"}\nbar\n{: id=\"barid\"}", "<p data-block=\"0\" data-node-id=\"fooid\" data-type=\"p\">foo</p><p data-block=\"0\" data-node-id=\"barid\" data-type=\"p\">bar</p>"},
	{"0", "", ""},
}

func TestMd2VditorIRBlock(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ToC = true
	luteEngine.BlockRef = true
	luteEngine.KramdownIAL = true
	luteEngine.Tag = true

	for _, test := range md2VditorIRBlockTests {
		md := luteEngine.Md2VditorIRBlockDOM(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}
