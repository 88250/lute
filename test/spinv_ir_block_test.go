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

var spinVditorIRBlockDOMTests = []*parseTest{

	{"5", "<p data-block=\"0\">foo ((bar)) <wbr></p>", "<p data-block=\"0\" data-node-id=\"\">foo <span data-type=\"block-ref\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__link\">bar</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span> <wbr></p>"},
	{"4", "<p data-block=\"0\" data-node-id=\"1596459249782\">((foo \"text\")<wbr></p>\n", "<p data-block=\"0\" data-node-id=\"\">((foo &quot;text&quot;)<wbr></p>"},
	{"3", "<p data-block=\"0\" data-node-id=\"1596459249782\">((foo \"text\"))<wbr></p>\n", "<p data-block=\"0\" data-node-id=\"\"><span data-type=\"block-ref\" class=\"vditor-ir__node vditor-ir__node--expand\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__link\">foo</span>text<span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span><wbr></p>"},
	{"2", "<p data-block=\"0\" data-node-id=\"1596459249782\">((foo))<wbr></p>\n", "<p data-block=\"0\" data-node-id=\"\"><span data-type=\"block-ref\" class=\"vditor-ir__node vditor-ir__node--expand\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__link\">foo</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span><wbr></p>"},
	{"1", "<p data-block=\"0\" data-node-id=\"1\">foo</p><p data-block=\"0\"><wbr><br></p>", "<p data-block=\"0\" data-node-id=\"\">foo</p><p data-block=\"0\" data-node-id=\"\"><wbr></p>"},
	{"0", "<p data-block=\"0\" data-node-id=\"1\">foo</p><p data-block=\"0\"><wbr><br></p>", "<p data-block=\"0\" data-node-id=\"\">foo</p><p data-block=\"0\" data-node-id=\"\"><wbr></p>"},
}

func TestSpinVditorIRBlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.BlockRef = true

	for _, test := range spinVditorIRBlockDOMTests {
		html := luteEngine.SpinVditorIRBlockDOM(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
