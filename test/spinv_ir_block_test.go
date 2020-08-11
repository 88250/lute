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

	{"10", "<p data-block=\"0\" data-node-id=\"20200811140724-dxmm3jo\"><span data-type=\"em\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span><em data-newline=\"1\">foo</em><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*b<wbr></span></span></p>", "<p data-block=\"0\" data-node-id=\"20200811140724-dxmm3jo\"><span data-type=\"em\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span><em data-newline=\"1\">foo</em><span class=\"vditor-ir__marker vditor-ir__marker--bi\">*</span></span>b<wbr></p>"},
	{"9", "<p data-block=\"0\" data-node-id=\"20200810211034-9d34ae\"><span data-type=\"block-ref\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">20200810191413-6f5616</span> <span class=\"vditor-ir__blockref\">\"foo\"</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">) <wbr></span></span></p>", "<p data-block=\"0\" data-node-id=\"20200810211034-9d34ae\"><span data-type=\"block-ref\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">20200810191413-6f5616</span> <span class=\"vditor-ir__blockref\">&#34;foo&#34;</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span> <wbr></p>"},
	{"8", "<p data-block=\"0\" data-node-id=\"20200809211933-b81ed7\">* foo<wbr></p>", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\" data-node-id=\"20200809211933-b81ed7\"><li data-marker=\"*\" data-node-id=\"\">foo<wbr></li></ul>"},
	{"7", "<p data-block=\"0\" data-node-id=\"20200809184752-a537de\">&gt; foo<wbr></p>", "<blockquote data-block=\"0\" data-node-id=\"20200809184752-a537de\"><p data-block=\"0\" data-node-id=\"\">foo<wbr></p></blockquote>"},
	{"6", "<p data-block=\"0\" data-node-id=\"20200809093825-b06abb\"><span data-type=\"block-ref\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">foo</span> <span class=\"vditor-ir__blockref\">\"text<wbr>\"</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span></p>", "<p data-block=\"0\" data-node-id=\"20200809093825-b06abb\"><span data-type=\"block-ref\" class=\"vditor-ir__node vditor-ir__node--expand\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">foo</span> <span class=\"vditor-ir__blockref\">&#34;text<wbr>&#34;</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span></p>"},
	{"5", "<p data-block=\"0\" data-node-id=\"1596459249782\">foo ((bar)) <wbr></p>", "<p data-block=\"0\" data-node-id=\"1596459249782\">foo <span data-type=\"block-ref\" class=\"vditor-ir__node\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">bar</span> <span class=\"vditor-ir__blockref\">&#34;bar&#34;</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span> <wbr></p>"},
	{"4", "<p data-block=\"0\" data-node-id=\"1596459249782\">((foo \"text\")<wbr></p>\n", "<p data-block=\"0\" data-node-id=\"1596459249782\">((foo &#34;text&#34;)<wbr></p>"},
	{"3", "<p data-block=\"0\" data-node-id=\"1596459249782\">((foo \"text\"))<wbr></p>\n", "<p data-block=\"0\" data-node-id=\"1596459249782\"><span data-type=\"block-ref\" class=\"vditor-ir__node vditor-ir__node--expand\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">foo</span> <span class=\"vditor-ir__blockref\">&#34;text&#34;</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span><wbr></p>"},
	{"2", "<p data-block=\"0\" data-node-id=\"1596459249782\">((foo))<wbr></p>\n", "<p data-block=\"0\" data-node-id=\"1596459249782\"><span data-type=\"block-ref\" class=\"vditor-ir__node vditor-ir__node--expand\"><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">(</span><span class=\"vditor-ir__marker vditor-ir__marker--link\">foo</span> <span class=\"vditor-ir__blockref\">&#34;foo&#34;</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span><span class=\"vditor-ir__marker vditor-ir__marker--paren\">)</span></span><wbr></p>"},
	{"1", "<p data-block=\"0\" data-node-id=\"1\">foo</p><p data-block=\"0\" data-node-id=\"2\"><wbr><br></p>", "<p data-block=\"0\" data-node-id=\"1\">foo</p><p data-block=\"0\" data-node-id=\"2\"><wbr></p>"},
	{"0", "<p data-block=\"0\" data-node-id=\"1\">foo</p><p data-block=\"0\" data-node-id=\"20200811112006-322210\"><wbr><br></p>", "<p data-block=\"0\" data-node-id=\"1\">foo</p><p data-block=\"0\" data-node-id=\"20200811112006-322210\"><wbr></p>"},
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
