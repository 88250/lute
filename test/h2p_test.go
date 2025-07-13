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

	{"4", "<pre><code class=\"language-json\">{\n    &quot;propName&quot;: &quot;propValue&quot;\n}\n</code></pre>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"code-block\"><div class=\"protyle-action\"><span class=\"protyle-action--first protyle-action__language\" contenteditable=\"false\">json</span><span class=\"fn__flex-1\"></span><span class=\"b3-tooltips__nw b3-tooltips protyle-icon protyle-icon--first protyle-action__copy\"><svg><use xlink:href=\"#iconCopy\"></use></svg></span><span class=\"b3-tooltips__nw b3-tooltips protyle-icon protyle-icon--last protyle-action__menu\"><svg><use xlink:href=\"#iconMore\"></use></svg></span></div><div class=\"hljs\"><div></div><div contenteditable=\"true\" style=\"flex: 1\" spellcheck=\"false\">{\n    &quot;propName&quot;: &quot;propValue&quot;\n}\n</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"3", "<div class=\"language-plantuml\">@startuml component\nactor client\nnode app\ndatabase db\ndb -&gt; app\napp -&gt; client\n@enduml\n</div>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"render-node\" data-content=\"@startuml component\nactor client\nnode app\ndatabase db\ndb -&gt; app\napp -&gt; client\n@enduml\" data-subtype=\"plantuml\"><div spin=\"1\"></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"2", "<div class=\"language-math\" id=\"20250713111927-m1adqgm\" name=\"公式1\" updated=\"20250713111932\">123</div>\n", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeMathBlock\" class=\"render-node\" name=\"公式1\" data-content=\"123\" data-subtype=\"math\"><div spin=\"1\"></div><div class=\"protyle-attr\" contenteditable=\"false\"><div class=\"protyle-attr--name\"><svg><use xlink:href=\"#iconN\"></use></svg>公式1</div>\u200b</div></div>"},
	{"1", "<h2 id=\"---我们的家\" name=\"社区\" updated=\"20210601233355\">🏘️ 我们的家</h2>", "<div data-subtype=\"h2\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h2\" name=\"社区\"><div contenteditable=\"true\" spellcheck=\"false\">🏘️ 我们的家</div><div class=\"protyle-attr\" contenteditable=\"false\"><div class=\"protyle-attr--name\"><svg><use xlink:href=\"#iconN\"></use></svg>社区</div>\u200b</div></div>"},
	{"0", "<table><tr><td><span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"foo\" contenteditable=\"false\" class=\"render-node\" data-render=\"true\">​<span class=\"katex\"><span class=\"katex-html\" aria-hidden=\"true\"><span class=\"base\"><span class=\"strut\" style=\"height:0.8889em;vertical-align:-0.1944em;\"></span><span class=\"mord mathnormal\" style=\"margin-right:0.10764em;\">f</span><span class=\"mord mathnormal\">oo</span></span></span></span></span>​</td></tr><tr><td></td></tr></table>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeTable\" class=\"table\"><div contenteditable=\"false\"><table contenteditable=\"true\" spellcheck=\"false\"><colgroup><col /></colgroup><thead><tr><th><span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"foo\" contenteditable=\"false\" class=\"render-node\"></span>\u200b</th></tr></thead><tbody><tr><td></td></tr></tbody></table><div class=\"protyle-action__table\"><div class=\"table__resize\"></div><div class=\"table__select\"></div></div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
}

func TestHTML2BlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetKramdownBlockIAL(true)

	ast.Testing = true
	for _, test := range html2BlockDOMTests {
		result := luteEngine.HTML2BlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}
