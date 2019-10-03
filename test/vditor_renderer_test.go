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
	//{"20", "> *foo\n> bar*\n", "<blockquote data-ntype=\"4\" data-mtype=\"1\"><span class=\"node\"><span class=\"marker\" data-ntype=\"5\" data-mtype=\"2\" data-caret=\"start\" data-caretoffset=\"0\">&gt;</span></span><p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\" data-ntype=\"12\" data-mtype=\"2\"><span class=\"marker\" data-ntype=\"13\" data-mtype=\"2\">*</span><em data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"11\" data-mtype=\"2\">foo</span><span data-ntype=\"24\" data-mtype=\"2\"><br><span class=\"newline\" data-ntype=\"24\" data-mtype=\"2\" />\n</span><span data-ntype=\"11\" data-mtype=\"2\">bar</span></em><span class=\"marker\" data-ntype=\"14\" data-mtype=\"2\">*</span></span></p><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span></blockquote>"},
	{"8", "-- ---\n", "<hr data-ntype=\"4\" data-mtype=\"0\" data-cso=\"2\" data-ceo=\"2\" />"},
	{"7", "__foo__\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\" data-ntype=\"21\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\">__</span><strong data-ntype=\"18\" data-mtype=\"2\"><span data-ntype=\"12\" data-mtype=\"2\">foo</span></strong><span class=\"marker\" data-ntype=\"22\" data-mtype=\"2\">__</span></span><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"},
	{"6", "`foo`\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"23\" data-mtype=\"2\"><span class=\"marker\" data-ntype=\"24\" data-mtype=\"2\">`</span><code data-ntype=\"23\" data-mtype=\"2\"><span data-ntype=\"25\" data-mtype=\"2\" data-cso=\"1\" data-ceo=\"1\">foo</span></code><span class=\"marker\" data-ntype=\"26\" data-mtype=\"2\">`</span></span><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"},
	{"5", "**foo**\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\" data-ntype=\"19\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\">**</span><strong data-ntype=\"18\" data-mtype=\"2\"><span data-ntype=\"12\" data-mtype=\"2\">foo</span></strong><span class=\"marker\" data-ntype=\"20\" data-mtype=\"2\">**</span></span><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"},
	{"4", "_foo_\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"13\" data-mtype=\"2\"><span class=\"marker\" data-ntype=\"16\" data-mtype=\"2\">_</span><em data-ntype=\"13\" data-mtype=\"2\"><span data-ntype=\"12\" data-mtype=\"2\" data-cso=\"1\" data-ceo=\"1\">foo</span></em><span class=\"marker\" data-ntype=\"17\" data-mtype=\"2\">_</span></span><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"},
	{"3", "> foo\n", "<div class=\"node\"><span class=\"marker\" data-ntype=\"6\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\">&gt; </span><blockquote data-ntype=\"6\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\"><p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"12\" data-mtype=\"2\">foo</span><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p></blockquote></div>"},
	{"2", "## foo\n", "<span class=\"node\"><span class=\"marker\" data-ntype=\"3\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\">## </span><h2 data-ntype=\"3\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\"><span data-ntype=\"12\" data-mtype=\"2\">foo</span></h2></span>"},
	{"1", "*foo*\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"13\" data-mtype=\"2\"><span class=\"marker\" data-ntype=\"14\" data-mtype=\"2\">*</span><em data-ntype=\"13\" data-mtype=\"2\"><span data-ntype=\"12\" data-mtype=\"2\" data-cso=\"1\" data-ceo=\"1\">foo</span></em><span class=\"marker\" data-ntype=\"15\" data-mtype=\"2\">*</span></span><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"},
	{"0", "foo\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"12\" data-mtype=\"2\" data-cso=\"2\" data-ceo=\"2\">foo</span><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"},
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

func TestVditorNewline(t *testing.T) {
	luteEngine := lute.New()

	html, err := luteEngine.VditorNewline(1, nil)
	if nil != err {
		t.Fatalf("unexpected: %s", err)
	}
	expected := "<p data-ntype=\"1\" data-mtype=\"0\"><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"
	if expected != html {
		t.Fatalf("vditor newline failed\nexpected\n\t%q\ngot\n\t%q", expected, html)
	}
}
