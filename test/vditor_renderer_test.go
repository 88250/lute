// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

// +build javascript

package test

import (
	"testing"

	"github.com/b3log/lute"
)

var vditorRendererTests = []parseTest{

	//{"21", "---\n", "<hr data-ntype=\"3\" data-mtype=\"0\" data-caret=\"start\" data-caretoffset=\"0\" />"},
	//{"20", "> *foo\n> bar*\n", "<blockquote data-ntype=\"4\" data-mtype=\"1\"><span class=\"node\"><span class=\"marker\" data-ntype=\"5\" data-mtype=\"2\" data-caret=\"start\" data-caretoffset=\"0\">&gt;</span></span><p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\" data-ntype=\"12\" data-mtype=\"2\"><span class=\"marker\" data-ntype=\"13\" data-mtype=\"2\">*</span><em data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"11\" data-mtype=\"2\">foo</span><span data-ntype=\"24\" data-mtype=\"2\"><br><span class=\"newline\" data-ntype=\"24\" data-mtype=\"2\" />\n</span><span data-ntype=\"11\" data-mtype=\"2\">bar</span></em><span class=\"marker\" data-ntype=\"14\" data-mtype=\"2\">*</span></span></p><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span></blockquote>"},
	{"6", "`foo`\n", "<span><span class=\"marker\">`</span><code data-ntype=\"13\" data-mtype=\"2\">Lute</code><span class=\"marker\">`</span>"},
	{"5", "**foo**\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\" data-ntype=\"19\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\">**</span><strong data-ntype=\"18\" data-mtype=\"2\"><span data-ntype=\"12\" data-mtype=\"2\">foo</span></strong><span class=\"marker\" data-ntype=\"20\" data-mtype=\"2\">**</span></span></p><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span>"},
	{"4", "_foo_\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"13\" data-mtype=\"2\"><span class=\"marker\" data-ntype=\"16\" data-mtype=\"2\">_</span><em data-ntype=\"13\" data-mtype=\"2\"><span data-ntype=\"12\" data-mtype=\"2\" data-cso=\"1\" data-ceo=\"1\">foo</span></em><span class=\"marker\" data-ntype=\"17\" data-mtype=\"2\">_</span></span></p><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span>"},
	{"3", "> foo\n", "<blockquote data-ntype=\"5\" data-mtype=\"1\"><span class=\"node\"><span class=\"marker\" data-ntype=\"6\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\">&gt;</span></span><p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"12\" data-mtype=\"2\">foo</span></p><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span></blockquote>"},
	{"2", "## foo\n", "<h2 data-ntype=\"2\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\" data-ntype=\"3\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\">##</span></span><span data-ntype=\"12\" data-mtype=\"2\">foo</span></h2>"},
	{"1", "*foo*\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"13\" data-mtype=\"2\"><span class=\"marker\" data-ntype=\"14\" data-mtype=\"2\">*</span><em data-ntype=\"13\" data-mtype=\"2\"><span data-ntype=\"12\" data-mtype=\"2\" data-cso=\"1\" data-ceo=\"1\">foo</span></em><span class=\"marker\" data-ntype=\"15\" data-mtype=\"2\">*</span></span></p><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span>"},
	{"0", "foo\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"12\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\">foo</span></p><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span>"},
}

func TestVditorRenderer(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorRendererTests {
		html, err := luteEngine.RenderVditorDOM(test.from, 2, 2)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
