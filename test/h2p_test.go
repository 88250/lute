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
	"github.com/88250/lute/ast"
)

var html2BlockDOMTests = []parseTest{

	{"6", "<ol><li><ul><li><strong>对比</strong>：<table><thead><tr><th><strong>维度</strong></th><th><strong>深度工作</strong></th><th><strong>浮浅工作</strong></th></tr></thead><tbody><tr><td><strong>认知强度</strong></td><td>高（需全神贯注）</td><td>低（可多任务处理）</td></tr><tr><td><strong>产出价值</strong></td><td>创新性、高价值（如产品设计）</td><td>重复性、低价值（如填报表）</td></tr><tr><td><strong>稀缺性</strong></td><td>信息时代的核心竞争力</td><td>易被自动化替代</td></tr></tbody></table></li></ul></li><li><p><strong>时代价值</strong></p></li></ol>", "<div data-subtype=\"o\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"1.\" data-subtype=\"o\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">1.</div><div data-subtype=\"u\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">\u200b<span data-type=\"strong\">对比</span>\u200b：</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeTable\" class=\"table\"><div contenteditable=\"false\"><table contenteditable=\"true\" spellcheck=\"false\"><colgroup><col /><col /><col /></colgroup><thead><tr><th><span data-type=\"strong\">维度</span></th><th><span data-type=\"strong\">深度工作</span></th><th><span data-type=\"strong\">浮浅工作</span></th></tr></thead><tbody><tr><td><span data-type=\"strong\">认知强度</span></td><td>高（需全神贯注）</td><td>低（可多任务处理）</td></tr><tr><td><span data-type=\"strong\">产出价值</span></td><td>创新性、高价值（如产品设计）</td><td>重复性、低价值（如填报表）</td></tr><tr><td><span data-type=\"strong\">稀缺性</span></td><td>信息时代的核心竞争力</td><td>易被自动化替代</td></tr></tbody></table><div class=\"protyle-action__table\"><div class=\"table__resize\"></div><div class=\"table__select\"></div></div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-marker=\"2.\" data-subtype=\"o\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">2.</div><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"strong\">时代价值</span></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"5", "<strong>对比</strong><span>：</span><table><thead><tr><th><strong>维度</strong></th><th><strong>深度工作</strong></th><th><strong>浮浅工作</strong></th></tr></thead><tbody><tr><td><strong>认知强度</strong></td><td>高（需全神贯注）</td><td>低（可多任务处理）</td></tr></tbody></table>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">\u200b<span data-type=\"strong\">对比</span>\u200b：</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"2\" data-type=\"NodeTable\" class=\"table\"><div contenteditable=\"false\"><table contenteditable=\"true\" spellcheck=\"false\"><colgroup><col /><col /><col /></colgroup><thead><tr><th><span data-type=\"strong\">维度</span></th><th><span data-type=\"strong\">深度工作</span></th><th><span data-type=\"strong\">浮浅工作</span></th></tr></thead><tbody><tr><td><span data-type=\"strong\">认知强度</span></td><td>高（需全神贯注）</td><td>低（可多任务处理）</td></tr></tbody></table><div class=\"protyle-action__table\"><div class=\"table__resize\"></div><div class=\"table__select\"></div></div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"4", "<pre><code class=\"language-json\">{\n    &quot;propName&quot;: &quot;propValue&quot;\n}\n</code></pre>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"code-block\"><div class=\"protyle-action\"><span class=\"protyle-action--first protyle-action__language\" contenteditable=\"false\">json</span><span class=\"fn__flex-1\"></span><span class=\"ariaLabel protyle-icon protyle-icon--first protyle-action__copy\" data-position=\"4north\"><svg><use xlink:href=\"#iconCopy\"></use></svg></span><span class=\"ariaLabel protyle-icon protyle-icon--last protyle-action__menu\" data-position=\"4north\"><svg><use xlink:href=\"#iconMore\"></use></svg></span></div><div class=\"hljs\"><div></div><div contenteditable=\"true\" style=\"flex: 1\" spellcheck=\"false\">{\n    &quot;propName&quot;: &quot;propValue&quot;\n}\n</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"3", "<div class=\"language-plantuml\">@startuml component\nactor client\nnode app\ndatabase db\ndb -&gt; app\napp -&gt; client\n@enduml\n</div>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"render-node\" data-content=\"@startuml component\nactor client\nnode app\ndatabase db\ndb -&gt; app\napp -&gt; client\n@enduml\" data-subtype=\"plantuml\"><div spin=\"1\"></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"2", "<div class=\"language-math\" id=\"20250713111927-m1adqgm\" name=\"公式1\" updated=\"20250713111932\">123</div>\n", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeMathBlock\" class=\"render-node\" name=\"公式1\" data-content=\"123\" data-subtype=\"math\"><div spin=\"1\"></div><div class=\"protyle-attr\" contenteditable=\"false\"><div class=\"protyle-attr--name\"><svg><use xlink:href=\"#iconN\"></use></svg>公式1</div>\u200b</div></div>"},
	{"1", "<h2 id=\"---我们的家\" name=\"社区\" updated=\"20210601233355\">🏘️ 我们的家</h2>", "<div data-subtype=\"h2\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h2\" name=\"社区\"><div contenteditable=\"true\" spellcheck=\"false\">🏘️ 我们的家</div><div class=\"protyle-attr\" contenteditable=\"false\"><div class=\"protyle-attr--name\"><svg><use xlink:href=\"#iconN\"></use></svg>社区</div>\u200b</div></div>"},
	{"0", "<table><tr><td><span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"foo\" contenteditable=\"false\" class=\"render-node\" data-render=\"true\">​<span class=\"katex\"><span class=\"katex-html\" aria-hidden=\"true\"><span class=\"base\"><span class=\"strut\" style=\"height:0.8889em;vertical-align:-0.1944em;\"></span><span class=\"mord mathnormal\" style=\"margin-right:0.10764em;\">f</span><span class=\"mord mathnormal\">oo</span></span></span></span></span>​</td></tr><tr><td></td></tr></table>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeTable\" class=\"table\"><div contenteditable=\"false\"><table contenteditable=\"true\" spellcheck=\"false\"><colgroup><col /></colgroup><thead><tr><th><span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"foo\" contenteditable=\"false\" class=\"render-node\"></span>\u200b</th></tr></thead><tbody><tr><td></td></tr></tbody></table><div class=\"protyle-action__table\"><div class=\"table__resize\"></div><div class=\"table__select\"></div></div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
}

func TestHTML2BlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetKramdownBlockIAL(true)
	luteEngine.SetHTML2MarkdownAttrs([]string{"name", "custom-*"})

	ast.Testing = true
	for _, test := range html2BlockDOMTests {
		result := luteEngine.HTML2BlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}

func TestHTML2BlockDOMTableCaption(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetKramdownIAL(true)
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetHTML2MarkdownAttrs([]string{"custom-*"})

	ast.Testing = true
	from := "<table custom-caption-test=\"yes\"><caption style=\"color: red; caption-side: BOTTOM !important\"><strong>A &amp; B</strong></caption><tr><td colspan=\"2\">Wide</td><td>C</td></tr><tr><td>D</td><td>E</td><td>F</td></tr></table>"
	expected := "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeTable\" class=\"table\" caption=\"&amp;lt;caption contenteditable=&amp;quot;false&amp;quot; style=&amp;quot;caption-side: bottom;&amp;quot;&amp;gt;A &amp;amp;amp; B&amp;lt;/caption&amp;gt;\" custom-caption-test=\"yes\" updated=\"20060102150405\"><div contenteditable=\"false\"><table contenteditable=\"true\" spellcheck=\"false\"><caption contenteditable=\"false\" style=\"caption-side: bottom;\">A &amp; B</caption><colgroup><col /><col /><col /></colgroup><thead><tr><th colspan=\"2\">Wide</th><th class=\"fn__none\"></th><th>C</th></tr></thead><tbody><tr><td>D</td><td>E</td><td>F</td></tr></tbody></table><div class=\"protyle-action__table\"><div class=\"table__resize\"></div><div class=\"table__select\"></div></div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"
	if result := luteEngine.HTML2BlockDOM(from); expected != result {
		t.Fatalf("expected\n\t%q\ngot\n\t%q", expected, result)
	}
}
