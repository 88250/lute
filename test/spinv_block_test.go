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

var spinBlockDOMTests = []*parseTest{

	{"7", "<div data-node-id=\"20210415235639-52nicpn\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210417175151\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"a\" data-href=\"https://foo.bar\" data-title=\"foo\">foob<wbr></span></div><div class=\"protyle-attr\"></div></div>", "<div data-node-id=\"20210415235639-52nicpn\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210417175151\"><div contenteditable=\"true\" spellcheck=\"false\"><span data-type=\"a\" data-href=\"https://foo.bar\" data-title=\"foo\">foob<wbr></span></div><div class=\"protyle-attr\"></div></div>"},
	{"6", "<div data-node-id=\"20210417163402-wmhbf42\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"code-block\" updated=\"20210417163422\"><div class=\"protyle-code\"><div class=\"protyle-code__language\">foo</div><div class=\"protyle-code__copy\"></div></div><div contenteditable=\"true\" spellcheck=\"false\">f<wbr></div><div class=\"protyle-attr\"></div></div>", "<div data-node-id=\"20210417163402-wmhbf42\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"code-block\" updated=\"20210417163422\"><div class=\"protyle-code\"><div class=\"protyle-code__language\"></div><div class=\"protyle-code__copy\"></div></div><div contenteditable=\"true\" spellcheck=\"false\">f<wbr>\n</div><div class=\"protyle-attr\"></div></div>"},
	{"5", "<div data-node-id=\"20210417102544-ti6msnt\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210417162215\"><div contenteditable=\"true\" spellcheck=\"false\">```<wbr></div><div class=\"protyle-attr\"></div></div>", "<div data-node-id=\"20210417102544-ti6msnt\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"code-block\" updated=\"20210417162215\"><div class=\"protyle-code\"><div class=\"protyle-code__language\"></div><div class=\"protyle-code__copy\"></div></div><div contenteditable=\"true\" spellcheck=\"false\"><wbr></div><div class=\"protyle-attr\"></div></div>"},
	{"4", "<div data-node-id=\"20210416090625-wpvw0au\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210417153511\"><div contenteditable=\"true\" spellcheck=\"false\"><strong style=\"color: var(--b3-font-color1);\">foob<wbr></strong></div><div class=\"protyle-attr\"></div></div>", "<div data-node-id=\"20210416090625-wpvw0au\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210417153511\"><div contenteditable=\"true\" spellcheck=\"false\"><strong style=\"color: var(--b3-font-color1);\">foob<wbr></strong></div><div class=\"protyle-attr\"></div></div>"},
	{"3", "<div data-node-id=\"20210413005325-jljhnw5\" data-node-index=\"1\" data-type=\"NodeSuperBlock\" class=\"sb\" data-sb-layout=\"col\"><div data-node-id=\"20210413003741-qbi5h4h\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div data-node-id=\"20210413001027-1mo28cc\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>", "<div data-node-id=\"20210413005325-jljhnw5\" data-node-index=\"1\" data-type=\"NodeSuperBlock\" class=\"sb\" data-sb-layout=\"col\"><div data-node-id=\"20210413003741-qbi5h4h\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div data-node-id=\"20210413001027-1mo28cc\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>"},
	{"2", "<div data-subtype=\"h1\" data-node-id=\"20210408153138-td774lp\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div>", "<div data-subtype=\"h1\" data-node-id=\"20210408153138-td774lp\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div>"},
	{"1", "<div data-node-id=\"20210408153137-zds0o4x\" data-node-index=\"1\" data-type=\"NodeBlockquote\" class=\"bq\"><div data-node-id=\"20210408153138-td774lp\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div></div>", "<div data-node-id=\"20210408153137-zds0o4x\" data-node-index=\"1\" data-type=\"NodeBlockquote\" class=\"bq\"><div data-node-id=\"20210408153138-td774lp\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>"},
	{"0", "<div data-subtype=\"u\" data-node-id=\"20210414223654-vfqydjh\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210415082227-m67yq1v\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-bullet\"></div><div data-node-id=\"20210415082227-z9mgkh5\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210415091213-c387rm0\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-bullet\"></div><div data-node-id=\"20210415091222-knbamrt\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>", "<div data-subtype=\"u\" data-node-id=\"20210414223654-vfqydjh\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210415082227-m67yq1v\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-bullet\"></div><div data-node-id=\"20210415082227-z9mgkh5\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210415091213-c387rm0\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-bullet\"></div><div data-node-id=\"20210415091222-knbamrt\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>"},
}

func TestSpinBlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetVditorWYSIWYG(true)
	luteEngine.ParseOptions.Mark = true
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.SetKramdownIAL(true)
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.SetLinkBase("/siyuan/0/测试笔记/")
	luteEngine.SetAutoSpace(false)
	luteEngine.SetSub(true)
	luteEngine.SetSup(true)
	luteEngine.SetGitConflict(true)

	ast.Testing = true
	for _, test := range spinBlockDOMTests {
		html := luteEngine.SpinBlockDOM(test.from)

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
	ast.Testing = false
}
