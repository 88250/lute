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

var md2VditorDOMTests = []parseTest{

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

var md2VditorIRBlockDOMTests = []parseTest{

	{"27", "* {: id=\"20210221200613-7vpmc8h\"}foo\n  {: id=\"20210221195351-x5tgalq\" updated=\"20210221201411\"}\n{: id=\"20210221195349-czsad7f\" updated=\"20210221195351\"}\n\n\n{: id=\"20210215183533-l36k5mo\" type=\"doc\"}", "<ul data-marker=\"*\" data-block-index=\"1\" data-node-id=\"20210221195349-czsad7f\" updated=\"20210221195351\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20210221200613-7vpmc8h\"><p data-node-id=\"20210221195351-x5tgalq\" data-type=\"p\" updated=\"20210221201411\">foo</p></li></ul>"},
	{"26", "<<<<<<< HEAD\n{: id=\"20210121085548-9vnyjk4\"}\n=======\n{: id=\"20210120223935-11oegu7\"}\n>>>>>>> parent of fe4124a (Revert \"commit data\")", "<div data-node-id=\"20060102150405-1a2b3c4\" data-block-index=\"1\" data-type=\"git-conflict\" class=\"vditor-ir__node\"><span data-type=\"git-conflict-open-marker\" class=\"vditor-ir__marker\"><<<<<<< HEAD</span><pre class=\"vditor-ir__marker--pre\"><code>{: id=&quot;20210121085548-9vnyjk4&quot;}\n=======\n{: id=&quot;20210120223935-11oegu7&quot;}\n</code></pre><span data-type=\"git-conflict-close-marker\" class=\"vditor-ir__marker\">>>>>>>> parent of fe4124a (Revert \"commit data\")</span></div>"},
	{"25", "<img src=\"https://foo?bar=baz\"/>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-block-index=\"1\" data-type=\"html-block\" class=\"vditor-ir__node\"><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code data-type=\"html-block\">&lt;img src=&quot;https://foo?bar=baz&quot;/&gt;</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><img src=\"https://foo?bar=baz\"/></pre></div>"},
	{"24", "* foo `{: id=\"\" type=\"doc\"}` bar", "<ul data-marker=\"*\" data-block-index=\"1\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"20060102150405-1a2b3c4\"><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">foo <span data-type=\"code\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">`</span><code data-newline=\"1\">{: id=&quot;&quot; type=&quot;doc&quot;}</code><span class=\"vditor-ir__marker\">`</span></span> bar</p></li></ul>"},
	{"23", "![foo](bar \"*baz* $bazz$\")", "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"1\"><span class=\"vditor-ir__node\" data-type=\"img\" style=\"display: block; text-align: center; white-space: initial;\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">foo</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">bar</span> <span class=\"vditor-ir__marker vditor-ir__marker--title\">\"*baz* $bazz$\"</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><img src=\"bar\" alt=\"foo\" /><span class=\"vditor-ir__preview\" data-render=\"2\"><em>baz</em> <span class=\"language-math\">bazz</span></span></span></p>"},
	{"22", "{{{\nfoo\n{: id=\"20201124225532-6mzfxc5\"}\n\n{{{\nbar\n{: id=\"20201124202133-sujebkr\"}\n\n{: id=\"20201124225918-sg44ucv\"}\n\n}}}\n{: id=\"20201124225743-1quklut\"}\n\nbaz\n{: id=\"20201124225821-xqjdfct\"}\n\n}}}‸\n{: id=\"20201124225824-fy4jie3\"}", "<div data-node-id=\"20201124225824-fy4jie3\" data-block-index=\"1\" data-type=\"super-block\" class=\"vditor-ir__node vditor-ir__node--expand\"><span data-type=\"super-block-open-marker\" class=\"vditor-ir__marker\">{{{</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"super-block-layout\"></span><p data-node-id=\"20201124225532-6mzfxc5\" data-type=\"p\">foo</p><div data-node-id=\"20201124225743-1quklut\" data-type=\"super-block\" class=\"vditor-ir__node\"><span data-type=\"super-block-open-marker\" class=\"vditor-ir__marker\">{{{</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"super-block-layout\"></span><p data-node-id=\"20201124202133-sujebkr\" data-type=\"p\">bar</p><span data-type=\"super-block-close-marker\" class=\"vditor-ir__marker\">}}}</span></div><p data-node-id=\"20201124225821-xqjdfct\" data-type=\"p\">baz‸</p><span data-type=\"super-block-close-marker\" class=\"vditor-ir__marker\">}}}</span></div>"},
	{"21", "{{{\nfoo\n{: id=\"20201124112436-amupplg\"}\n\n}}}\n{: id=\"20201124113702-drj95ri\"}\n\nbar\n{: id=\"20201124114838-5qung2s\"}", "<div data-node-id=\"20201124113702-drj95ri\" data-block-index=\"1\" data-type=\"super-block\" class=\"vditor-ir__node\"><span data-type=\"super-block-open-marker\" class=\"vditor-ir__marker\">{{{</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"super-block-layout\"></span><p data-node-id=\"20201124112436-amupplg\" data-type=\"p\">foo</p><span data-type=\"super-block-close-marker\" class=\"vditor-ir__marker\">}}}</span></div><p data-node-id=\"20201124114838-5qung2s\" data-type=\"p\" data-block-index=\"2\">bar</p>"},
	{"20", "![foo](bar \"baz\")", "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"1\"><span class=\"vditor-ir__node\" data-type=\"img\" style=\"display: block; text-align: center; white-space: initial;\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">foo</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">bar</span> <span class=\"vditor-ir__marker vditor-ir__marker--title\">\"baz\"</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><img src=\"bar\" alt=\"foo\" /><span class=\"vditor-ir__preview\" data-render=\"2\">baz</span></span></p>"},
	{"19", "<audio controls=\"controls\" src=\"assets/20201118233326-ulhglhc-record1605713606242.wav\"></audio>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-block-index=\"1\" data-type=\"html-block\" class=\"vditor-ir__node\"><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code data-type=\"html-block\">&lt;audio controls=&quot;controls&quot; src=&quot;assets/20201118233326-ulhglhc-record1605713606242.wav&quot;&gt;&lt;/audio&gt;</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><audio controls=\"controls\" src=\"assets/20201118233326-ulhglhc-record1605713606242.wav\"></audio></pre></div>"},
	{"18", "* {: id=\"20201118232739-leknj8q\"}[ ] foo\n  {: id=\"20201118233118-tdrr6q0\"}\n\n   bar\n  {: id=\"20201118233118-y62y9hq\"}\n{: id=\"20201118232736-8ljalfq\"}", "<ul data-marker=\"*\" data-block-index=\"1\" data-node-id=\"20201118232736-8ljalfq\" data-type=\"ul\"><li data-marker=\"*\" class=\"vditor-task\" data-node-id=\"20201118232739-leknj8q\"><p data-node-id=\"20201118233118-tdrr6q0\" data-type=\"p\"><input type=\"checkbox\" /> foo</p><p data-node-id=\"20201118233118-y62y9hq\" data-type=\"p\">bar</p></li></ul>"},
	{"17", "* {: id=\"20201112104625-ix4w985\"}[ ] ‸\n* {: id=\"20201112110516-nrtkabp\"}[ ] bar\n{: id=\"20201112104625-zoon9oa\"}", "<ul data-marker=\"*\" data-block-index=\"1\" data-node-id=\"20201112104625-zoon9oa\" data-type=\"ul\"><li data-marker=\"*\" class=\"vditor-task\" data-node-id=\"20201112104625-ix4w985\"><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\"><input type=\"checkbox\" /> ‸</p></li><li data-marker=\"*\" class=\"vditor-task\" data-node-id=\"20201112110516-nrtkabp\"><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\"><input type=\"checkbox\" /> bar</p></li></ul>"},
	{"16", "foo\n`bar`\n", "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"1\">foo\n<span data-type=\"code\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">`</span><code data-newline=\"1\">bar</code><span class=\"vditor-ir__marker\">`</span></span></p>"},
	{"15", "![](assets/中文/foo.png)\n", "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"1\"><span class=\"vditor-ir__node\" data-type=\"img\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">[</span><span class=\"vditor-ir__marker vditor-ir__marker--bracket\">]</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">assets/中文/foo.png</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><img src=\"assets/中文/foo.png\" /></span></p>"},
	{"14", "| foo |\n| - |\n| bar *baz* |\n", "<table data-type=\"table\" data-node-id=\"20060102150405-1a2b3c4\" data-block-index=\"1\"><thead><tr><th>foo</th></tr></thead><tbody><tr><td>bar <span data-type=\"em\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--em\">*</span><em data-newline=\"1\">baz</em><span class=\"vditor-ir__marker vditor-ir__marker--em\">*</span></span></td></tr></tbody></table>"},
	{"13", "* {: id=\"fooid\"}foo\n  * {: id=\"barid\"}bar\n  {: id=\"ul2\"}\n{: id=\"ul1\"}\n", "<ul data-marker=\"*\" data-block-index=\"1\" data-node-id=\"ul1\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"fooid\"><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-node-id=\"ul2\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"barid\"><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">bar</p></li></ul></li></ul>"},
	{"12", "* {: id=\"fooid\"}foo\n  * {: id=\"barid\"}bar\n  {: id=\"ul2\"}\n{: id=\"ul1\"}", "<ul data-marker=\"*\" data-block-index=\"1\" data-node-id=\"ul1\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"fooid\"><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">foo</p><ul data-marker=\"*\" data-node-id=\"ul2\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"barid\"><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">bar</p></li></ul></li></ul>"},
	{"11", "* {: id=\"fooid\"}foo\n{: id=\"id\"}", "<ul data-marker=\"*\" data-block-index=\"1\" data-node-id=\"id\" data-type=\"ul\"><li data-marker=\"*\" data-node-id=\"fooid\"><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\">foo</p></li></ul>"},
	{"10", "foo#bar#baz", "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"1\">foo<span data-type=\"tag\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--tag\">#</span><span class=\"vditor-ir__tag\">bar</span><span class=\"vditor-ir__marker vditor-ir__marker--tag\">#</span></span>baz</p>"},
	{"9", "foo\n((id \"text\"))bar", "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"1\">foo\n<span data-type=\"block-ref\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">id</span><span class=\"vditor-ir__marker\"> </span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__marker vditor-ir__marker--info\">text</span><span class=\"vditor-ir__marker\">\"</span><span data-type=\"ref-text-tpl-render-result\" class=\"vditor-ir__blockref\"></span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span>bar</p>"},
	{"8", "foo\n!((id \"text\"))\nbar", "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"1\">foo</p><div data-node-id=\"20060102150405-1a2b3c4\" data-block-index=\"2\" data-type=\"block-ref-embed\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">id</span><span class=\"vditor-ir__marker\"> </span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__marker vditor-ir__marker--info\">text</span><span class=\"vditor-ir__marker\">\"</span><span data-type=\"ref-text-tpl-render-result\" class=\"vditor-ir__blockref\"></span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><div data-block-def-id=\"id\" data-render=\"2\" data-type=\"block-render\"></div></div><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"3\">bar</p>"},
	{"7.1", "foo!((id \"text\"))", "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"1\">foo!<span data-type=\"block-ref\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">id</span><span class=\"vditor-ir__marker\"> </span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__marker vditor-ir__marker--info\">text</span><span class=\"vditor-ir__marker\">\"</span><span data-type=\"ref-text-tpl-render-result\" class=\"vditor-ir__blockref\"></span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span></p>"},
	{"7", "!((id))", "<div data-node-id=\"20060102150405-1a2b3c4\" data-block-index=\"1\" data-type=\"block-ref-embed\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">id</span><span class=\"vditor-ir__marker\"> </span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__marker vditor-ir__marker--info\"></span><span class=\"vditor-ir__marker\">\"</span><span data-type=\"ref-text-tpl-render-result\" class=\"vditor-ir__blockref\"></span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><div data-block-def-id=\"id\" data-render=\"2\" data-type=\"block-render\"></div></div>"},
	{"6", "!((id \"text\"))", "<div data-node-id=\"20060102150405-1a2b3c4\" data-block-index=\"1\" data-type=\"block-ref-embed\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker\">!</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">id</span><span class=\"vditor-ir__marker\"> </span><span class=\"vditor-ir__marker\">\"</span><span class=\"vditor-ir__marker vditor-ir__marker--info\">text</span><span class=\"vditor-ir__marker\">\"</span><span data-type=\"ref-text-tpl-render-result\" class=\"vditor-ir__blockref\"></span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><div data-block-def-id=\"id\" data-render=\"2\" data-type=\"block-render\"></div></div>"},
	{"5", "<!--foo-->\nbar", "<div data-node-id=\"20060102150405-1a2b3c4\" data-block-index=\"1\" data-type=\"html-block\" class=\"vditor-ir__node\"><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code data-type=\"html-block\">&lt;!--foo--&gt;</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><!--foo--></pre></div><p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"2\">bar</p>"},
	{"4", "# foo\n{: id=\"test\" bookmark=\"bookmark\"}", "<h1 data-block-index=\"1\" class=\"vditor-ir__node\" data-node-id=\"test\" data-type=\"h\" bookmark=\"bookmark\" tip=\"bookmark\" data-marker=\"#\"><span class=\"vditor-ir__marker vditor-ir__marker--heading\" data-type=\"heading-marker\"># </span>foo</h1>"},
	{"3", "```\nfoo\n```\n{: id=\"test\" bookmark=\"bookmark\"}", "<div data-node-id=\"test\" data-block-index=\"1\" bookmark=\"bookmark\" tip=\"bookmark\" data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">\u200b</span><pre><code>foo\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div>"},
	{"2", "```\nfoo\n```\n{: id=\"test\"}", "<div data-node-id=\"test\" data-block-index=\"1\" data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">\u200b</span><pre><code>foo\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div>"},
	{"1", "foo\n{: id=\"fooid\"}\nbar\n{: id=\"barid\"}", "<p data-node-id=\"fooid\" data-type=\"p\" data-block-index=\"1\">foo</p><p data-node-id=\"barid\" data-type=\"p\" data-block-index=\"2\">bar</p>"},
	{"0", "", ""},
}

func TestMd2VditorIRBlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetVditorIR(true)
	luteEngine.SetToC(true)
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.ParseOptions.KramdownBlockIAL = true
	luteEngine.RenderOptions.KramdownBlockIAL = true
	luteEngine.ParseOptions.Tag = true
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.ParseOptions.GitConflict = true

	ast.Testing = true
	for _, test := range md2VditorIRBlockDOMTests {
		result := luteEngine.Md2VditorIRBlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
	ast.Testing = false
}

var inlineMd2VditorIRBlockDOMTests = []parseTest{

	{"0", "1. foo", "<p data-node-id=\"20060102150405-1a2b3c4\" data-type=\"p\" data-block-index=\"1\">1. foo</p>"},
}

func TestInlineMd2VditorIRBlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetVditorIR(true)
	luteEngine.SetToC(true)
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.ParseOptions.KramdownBlockIAL = true
	luteEngine.RenderOptions.KramdownBlockIAL = true
	luteEngine.ParseOptions.Tag = true
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.ParseOptions.GitConflict = true

	ast.Testing = true
	for _, test := range inlineMd2VditorIRBlockDOMTests {
		result := luteEngine.InlineMd2VditorIRBlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
	ast.Testing = false
}

var md2BlockDOMTests = []parseTest{

	{"12", "<p>foo</p>\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">&lt;p&gt;foo&lt;/p&gt;</div><div class=\"protyle-attr\"></div></div>"},
	{"11", "foo((20210121085548-9vnyjk4 \"{{.text}}\"))**bar**\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo<span data-type=\"block-ref\" data-id=\"20210121085548-9vnyjk4\" data-anchor=\"{{.text}}\"></span><strong>bar</strong></div><div class=\"protyle-attr\"></div></div>"},
	{"10", "$$\nfoo\n$$\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeMathBlock\" class=\"render-node\" data-content=\"foo\" data-subtype=\"math\"><div spin=\"1\"></div><div class=\"protyle-attr\"></div></div>"},
	{"9", "```abc\nfoo\n```\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"render-node\" data-content=\"foo\" data-subtype=\"abc\"><div spin=\"1\"></div><div class=\"protyle-attr\"></div></div>"},
	{"8", "```\nfoo\n```\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"code-block\"><div class=\"protyle-code\"><div class=\"protyle-code__language\"></div><div class=\"protyle-code__copy\"></div></div><div contenteditable=\"true\" spellcheck=\"false\">foo\n</div><div class=\"protyle-attr\"></div></div>"},
	{"7", "foo\n{: id=\"20210408204847-qyy54hz\"}\n---\n{: id=\"20210408204848-qyy54ha\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div data-node-id=\"20210408204848-qyy54ha\" data-node-index=\"2\" data-type=\"NodeThematicBreak\" class=\"hr\"><div></div></div>"},
	{"6", "<<<<<<< HEAD\n{: id=\"20210121085548-9vnyjk4\"}\n=======\n{: id=\"20210120223935-11oegu7\"}\n>>>>>>> parent of fe4124a (Revert \"commit data\")", "<div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeType(497)\" class=\"git-conflict\"><div contenteditable=\"true\" spellcheck=\"false\">{: id=&quot;20210121085548-9vnyjk4&quot;}\n=======\n{: id=&quot;20210120223935-11oegu7&quot;}</div><div class=\"protyle-attr\"></div></div>"},
	{"5", "* {: id=\"20210415082227-m67yq1v\"}foo\n  {: id=\"20210415082227-z9mgkh5\"}\n* {: id=\"20210415091213-c387rm0\"}bar\n  {: id=\"20210415091222-knbamrt\"}\n{: id=\"20210414223654-vfqydjh\"}\n\n\n{: id=\"20210413000727-8ua3vhv\" type=\"doc\"}\n", "<div data-subtype=\"u\" data-node-id=\"20210414223654-vfqydjh\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210415082227-m67yq1v\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-bullet\"></div><div data-node-id=\"20210415082227-z9mgkh5\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210415091213-c387rm0\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-bullet\"></div><div data-node-id=\"20210415091222-knbamrt\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>"},
	{"4", "{{{col\nfoo\n{: id=\"20210413003741-qbi5h4h\"}\n\nbar\n{: id=\"20210413001027-1mo28cc\"}\n}}}\n{: id=\"20210413005325-jljhnw5\"}\n\n\n{: id=\"20210413000727-8ua3vhv\" type=\"doc\"}\n", "<div data-node-id=\"20210413005325-jljhnw5\" data-node-index=\"1\" data-type=\"NodeSuperBlock\" class=\"sb\" data-sb-layout=\"col\"><div data-node-id=\"20210413003741-qbi5h4h\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div data-node-id=\"20210413001027-1mo28cc\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>"},
	{"3", "# foo\n{: id=\"20210408153138-td774lp\"}", "<div data-subtype=\"h1\" data-node-id=\"20210408153138-td774lp\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div>"},
	{"2", "> foo\n> {: id=\"20210408153138-td774lp\"}\n{: id=\"20210408153137-zds0o4x\"}", "<div data-node-id=\"20210408153137-zds0o4x\" data-node-index=\"1\" data-type=\"NodeBlockquote\" class=\"bq\"><div data-node-id=\"20210408153138-td774lp\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>"},
	{"1", "foo\n{: id=\"20210408204847-qyy54hz\" bookmark=\"bm\"}\nbar\n{: id=\"20210408204848-qyy54ha\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" bookmark=\"bm\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"><div class=\"protyle-attr--bookmark\">bm</div></div></div><div data-node-id=\"20210408204848-qyy54ha\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\"></div></div>"},
	{"0", "", ""},
}

func TestMd2BlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetToC(true)
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.ParseOptions.KramdownBlockIAL = true
	luteEngine.RenderOptions.KramdownBlockIAL = true
	luteEngine.ParseOptions.Tag = true
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.ParseOptions.GitConflict = true

	ast.Testing = true
	for _, test := range md2BlockDOMTests {
		result := luteEngine.Md2BlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
	ast.Testing = false
}
