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
	{"16", "> 一  \n  \n* 二  \n  \n# 三\n", "<blockquote class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\">&gt; </span><p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"0\" data-ceo=\"0\">一</span><span class=\"newline\">\n\n</span></p></blockquote><ul data-ntype=\"7\" data-mtype=\"1\"><li class=\"node node--block\" data-ntype=\"8\" data-mtype=\"1\"><span class=\"marker\">* </span><p data-ntype=\"1\" data-mtype=\"0\"><span>二</span><span class=\"newline\">\n\n</span></p></li></ul><h1 class=\"node\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\"># </span><span>三</span></h1>"},
	{"15", "一\n\n二\n\n三\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span>一</span><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"0\" data-ceo=\"0\">二</span><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><span>三</span><span class=\"newline\">\n\n</span></p>"},
	{"14", "fo\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"2\" data-ceo=\"2\">fo</span><span class=\"newline\">\n\n</span></p>"},
	{"13", "", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"1 data-mtype=\"2\" data-cso=\"0\" data-ceo=\"0\"></span><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"},
	{"12", "foo\nbar\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"2\" data-ceo=\"2\">foo</span>\n<span>bar</span><span class=\"newline\">\n\n</span></p>"},
	{"11", "> # foo\n", "<blockquote class=\"node node--block\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\">&gt; </span><h1 class=\"node node--expand\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\" data-cso=\"0\" data-ceo=\"0\"># </span><span>foo</span></h1></blockquote>"},
	{"10", "> #\n", "<blockquote class=\"node node--block\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\">&gt; </span><h1 class=\"node node--expand\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\" data-cso=\"0\" data-ceo=\"0\">#\n</span></h1></blockquote>"},
	{"9", "> ## foo\n", "<blockquote class=\"node node--block\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\">&gt; </span><h2 class=\"node node--expand\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\" data-cso=\"0\" data-ceo=\"0\">## </span><span>foo</span></h2></blockquote>"},
	{"8", "-- ---\n", "<hr data-ntype=\"4\" data-mtype=\"0\" data-cso=\"2\" data-ceo=\"2\" />"},
	{"7", "__foo__\n", "<p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\">__</span><span data-cso=\"0\" data-ceo=\"0\">foo</span><span class=\"marker\">__</span></strong><span class=\"newline\">\n\n</span></p>"},
	{"6", "`foo`\n", "<p data-ntype=\"1\" data-mtype=\"0\"><code class=\"node node--expand\" data-ntype=\"23\" data-mtype=\"2\"><span class=\"marker\">`</span>foo<span class=\"marker\">`</span></code><span class=\"newline\">\n\n</span></p>"},
	{"5", "**foo**\n", "<p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\">**</span><span data-cso=\"0\" data-ceo=\"0\">foo</span><span class=\"marker\">**</span></strong><span class=\"newline\">\n\n</span></p>"},
	{"4", "_foo_\n", "<p data-ntype=\"1\" data-mtype=\"0\"><em class=\"node node--expand\" data-ntype=\"13\" data-mtype=\"2\"><span class=\"marker\">_</span><span data-cso=\"1\" data-ceo=\"1\">foo</span><span class=\"marker\">_</span></em><span class=\"newline\">\n\n</span></p>"},
	{"3", "> foo\n", "<blockquote class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\">&gt; </span><p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"0\" data-ceo=\"0\">foo</span><span class=\"newline\">\n\n</span></p></blockquote>"},
	{"2", "## foo\n", "<h2 class=\"node node--expand\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">## </span><span>foo</span></h2>"},
	{"1", "*foo*\n", "<p data-ntype=\"1\" data-mtype=\"0\"><em class=\"node node--expand\" data-ntype=\"13\" data-mtype=\"2\"><span class=\"marker\">*</span><span data-cso=\"1\" data-ceo=\"1\">foo</span><span class=\"marker\">*</span></em><span class=\"newline\">\n\n</span></p>"},
	{"0", "foo\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"2\" data-ceo=\"2\">foo</span><span class=\"newline\">\n\n</span></p>"},
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
