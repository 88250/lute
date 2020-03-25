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

	{"16", "# heading {#custom-id}\n", "<h1 data-block=\"0\" data-id=\"#custom-id\" data-marker=\"#\">heading</h1>"},
	{"15", "foo\n\n[^1]: 111\n\n[2]: 222\n", "<p data-block=\"0\">foo\n</p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[2]: 222\n</div><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"><p data-block=\"0\">111\n</p></li></ol></div>"},
	{"14", "[^1]\n\n[^1]:\n", "<p data-block=\"0\">\u200b<sup data-type=\"footnotes-ref\" data-footnotes-label=\"^1\">1</sup>\u200b\n</p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"></li></ol></div>"},
	{"13", "[toc]\n\n# foo", "<div class=\"vditor-toc\" data-block=\"0\" data-type=\"toc-block\" contenteditable=\"false\"><span data-type=\"toc-h\">foo</span><br></div><p data-block=\"0\"></p><h1 data-block=\"0\" data-marker=\"#\">foo</h1>"},
	{"12", "foo[^1]\n[^1]:bar\n    * baz", "<p data-block=\"0\">foo<sup data-type=\"footnotes-ref\" data-footnotes-label=\"^1\">1</sup>\u200b\n</p><div data-block=\"0\" data-type=\"footnotes-block\"><ol data-type=\"footnotes-defs-ol\"><li data-type=\"footnotes-li\" data-marker=\"^1\"><p data-block=\"0\">bar\n</p><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">baz</li></ul></li></ol></div>"},
	{"11", "[foo][1]\n\n[1]: /bar\n", "<p data-block=\"0\">\u200b<span data-type=\"link-ref\" data-link-label=\"1\">foo</span>\u200b\n</p><div data-block=\"0\" data-type=\"link-ref-defs-block\">[1]: /bar\n</div>"},
	{"10", "Foo\n    ---\n", "<p data-block=\"0\">Foo\n---\n</p>"},
	{"9", "    ***\n     ***\n\n-     -      -      -", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"***\n ***\n\"><pre><code>***\n ***\n</code></pre></div><hr data-block=\"0\" />"},
	{"8", "    ***\n", "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"***\n\"><pre><code>***\n</code></pre></div>"},
	{"7", "* a\n  * b", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">a<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">b</li></ul></li></ul>"},
	{"6", "[]()", "<p data-block=\"0\">[]()\n</p>"},
	{"5", "* [ ]", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\" class=\"vditor-task\"><input type=\"checkbox\" /> </li></ul>"},
	{"4", "*", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">\u200b</li></ul>"},
	{"3", "foo'%'bar", "<p data-block=\"0\">foo'%'bar\n</p>"},
	{"2", "<p align=\"center\">\nfoo\n</p>\n\nbar", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;p align=&quot;center&quot;&gt;\nfoo\n&lt;/p&gt;</code></pre></div><p data-block=\"0\">bar\n</p>"},
	{"1", `foo\<aa>bar`, "<p data-block=\"0\">foo<span data-type=\"backslash\"><span>\\</span>&lt;</span>aa&gt;bar\n</p>"},
	{"0", `<details>
<summary>foo</summary>

* bar

</details>`, "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;details&gt;\n&lt;summary&gt;foo&lt;/summary&gt;</code></pre></div><ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul><div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;/details&gt;</code></pre></div>"},
}

func TestMd2Vditor(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ToC = true

	for _, test := range md2VditorTests {
		md := luteEngine.Md2VditorDOM(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}

var md2VditorIRTests = []parseTest{

	{"7", "$$\nfoo\n$$", "<div data-block=\"0\" data-type=\"math-block\" class=\"vditor-ir__node\"><span data-type=\"math-block-open-marker\">$$</span><pre class=\"vditor-ir__marker vditor-ir__marker--pre\"><code data-type=\"math-block\" class=\"language-math\">foo</code></pre><pre class=\"vditor-ir__preview\" data-render=\"false\"><code data-type=\"math-block\" class=\"language-math\">foo</code></pre><span data-type=\"math-block-close-marker\">$$</span></div>"},
	{"6", "foo<bar>baz", "<p data-block=\"0\">foo<span data-type=\"inline-node\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker--link\">&lt;bar&gt;</span></span>baz\n</p>"},
	{"5", "$foo$", "<p data-block=\"0\"><span data-type=\"inline-node\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">$</span><code data-newline=\"1\" class=\"vditor-ir__marker vditor-ir__marker--pre\">foo</code><span class=\"vditor-ir__preview\" data-render=\"false\"><code class=\"language-math\">foo</code></span><span class=\"vditor-ir__marker\">$</span></span>\n</p>"},
	{"4", "```foo\nbar\n```", "<div data-block=\"0\" data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">\u200bfoo</span><pre class=\"vditor-ir__marker vditor-ir__marker--pre\"><code class=\"language-foo\">bar\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"false\"><code class=\"language-foo\">bar\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div>"},
	{"3", "__foo__", "<p data-block=\"0\"><span data-type=\"strong\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bi\">__</span><strong data-newline=\"1\">foo</strong><span class=\"vditor-ir__marker vditor-ir__marker--bi\">__</span></span>\n</p>"},
	{"2", "* foo\n  * bar", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">foo<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></li></ul>"},
	{"1", "# foo", "<h1 data-block=\"0\" class=\"vditor-ir__node\" data-marker=\"#\"><span class=\"vditor-ir__marker vditor-ir__marker--heading\"># </span>foo</h1>"},
	{"0", "*foo*", "<p data-block=\"0\"><span data-type=\"em\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span><em data-newline=\"1\">foo</em><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span></span>\n</p>"},
}

func TestMd2VditorIR(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ToC = true

	for _, test := range md2VditorIRTests {
		md := luteEngine.Md2VditorIRDOM(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}
