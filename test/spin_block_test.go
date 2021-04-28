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

	{"7", "<div data-node-id=\"20210428163425-gd63njj\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210428163507\"><div contenteditable=\"true\" spellcheck=\"false\">{{foo}}<wbr></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-content=\"foo\" data-node-id=\"20210428163425-gd63njj\" data-node-index=\"1\" data-type=\"NodeBlockQueryEmbed\" class=\"render-node\" updated=\"20210428163507\"></div>"},
	{"6", "<div data-node-id=\"20210428155259-1j2zqx0\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210428155312\"><div contenteditable=\"true\" spellcheck=\"false\">foo\n\nb<wbr></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-node-id=\"20210428155259-1j2zqx0\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210428155312\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">b<wbr></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"5", "<div data-node-id=\"20210428094047-w9di4p3\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210428094048\"><div contenteditable=\"true\" spellcheck=\"false\"># <wbr></div><div class=\"protyle-attr\"></div></div>", "<div data-subtype=\"h1\" data-node-id=\"20210428094047-w9di4p3\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20210428094048\"><div contenteditable=\"true\" spellcheck=\"false\"><wbr></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"4", "<div data-node-id=\"20210428094047-w9di4p3\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210428094048\"><div contenteditable=\"true\" spellcheck=\"false\">#<wbr></div><div class=\"protyle-attr\"></div></div>", "<div data-node-id=\"20210428094047-w9di4p3\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210428094048\"><div contenteditable=\"true\" spellcheck=\"false\">#<wbr></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"3", "<div data-content=\"name:foo\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeBlockQueryEmbed\" class=\"render-node\"></div>", "<div data-content=\"name:foo\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeBlockQueryEmbed\" class=\"render-node\"></div>"},
	{"2", "<div data-subtype=\"t\" data-node-id=\"20210426154414-99vddn3\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210426154416\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210426154418-c3qfhdw\" data-type=\"NodeListItem\" class=\"li\"><input type=\"checkbox\"><div data-node-id=\"20210426154418-08q0qvv\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210426154423\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-subtype=\"t\" data-node-id=\"20210426154414-99vddn3\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210426154416\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210426154418-c3qfhdw\" data-type=\"NodeListItem\" class=\"li\"><input type=\"checkbox\" /><div data-node-id=\"20210426154418-08q0qvv\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210426154423\"><div contenteditable=\"true\" spellcheck=\"false\"> foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"1", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><kbd>foo</kbd></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><kbd>foo</kbd></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"0", "<div data-node-id=\"20210426094859-uataalw\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210426101601\"><div contenteditable=\"true\" spellcheck=\"false\">[[<wbr><wbr><span data-type=\"block-ref\" data-id=\"20210426091959-npvs57l\" data-anchor=\"\" contenteditable=\"false\"></span><span data-type=\"block-ref\" data-id=\"20210426091959-npvs57l\" data-anchor=\"\" contenteditable=\"false\"></span><span data-type=\"block-ref\" data-id=\"20210426091959-npvs57l\" data-anchor=\"\" contenteditable=\"false\"></span><span data-type=\"block-ref\" data-id=\"20210426091959-npvs57l\" data-anchor=\"\" contenteditable=\"false\"></span><span data-type=\"block-ref\" data-id=\"20210426091959-npvs57l\" data-anchor=\"\" contenteditable=\"false\"></span></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-node-id=\"20210426094859-uataalw\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210426101601\"><div contenteditable=\"true\" spellcheck=\"false\">[[<wbr><wbr><span data-type=\"block-ref\" data-id=\"20210426091959-npvs57l\" data-anchor=\"\" contenteditable=\"false\"></span><span data-type=\"block-ref\" data-id=\"20210426091959-npvs57l\" data-anchor=\"\" contenteditable=\"false\"></span><span data-type=\"block-ref\" data-id=\"20210426091959-npvs57l\" data-anchor=\"\" contenteditable=\"false\"></span><span data-type=\"block-ref\" data-id=\"20210426091959-npvs57l\" data-anchor=\"\" contenteditable=\"false\"></span><span data-type=\"block-ref\" data-id=\"20210426091959-npvs57l\" data-anchor=\"\" contenteditable=\"false\"></span></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestSpinBlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
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
