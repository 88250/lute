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

type vditorTest struct {
	*parseTest
	startOffset, endOffset int
}

var vditorRendererTests = []*vditorTest{

	//{&parseTest{"26", "foo\n\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">:</span><span data-hidden=\"❤️\"></span><span class=\"marker\">heart:</span></span><span data-cso=\"4\" data-ceo=\"4\"> foo</span><span class=\"newline\">\n\n</span></p>"}, 5, 5},
	{&parseTest{"25", ":heart: foo\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">:</span><span data-hidden=\"❤️\"></span><span class=\"marker\">heart:</span></span><span data-cso=\"4\" data-ceo=\"4\"> foo</span><span class=\"newline\">\n\n</span></p>"}, 11, 11},
	{&parseTest{"24", "foo:heart:bar\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span>foo</span><span class=\"node\"><span class=\"marker\">:</span><span data-hidden=\"❤️\"></span><span class=\"marker\" data-cso=\"1\" data-ceo=\"1\">heart:</span></span><span>bar</span><span class=\"newline\">\n\n</span></p>"}, 4, 4},
	{&parseTest{"23", "-\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"0\" data-ceo=\"0\">-</span><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"22", "https://hacpai.com\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"2\" data-ceo=\"2\">https://hacpai.com</span><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"21", "[foo](/bar )1\n", "<p data-ntype=\"1\" data-mtype=\"0\"><a class=\"node node--expand\" href=\"/bar\" data-ntype=\"29\" data-mtype=\"2\"><span class=\"marker\">[</span><span data-cso=\"1\" data-ceo=\"1\">foo</span><span class=\"marker\">]</span><span class=\"marker\">(</span><span class=\"marker\">/bar</span><span class=\"marker\"> </span><span class=\"marker\">)</span></a><span>1</span><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"20", "[foo](/bar)\n", "<p data-ntype=\"1\" data-mtype=\"0\"><a class=\"node node--expand\" href=\"/bar\" data-ntype=\"29\" data-mtype=\"2\"><span class=\"marker\">[</span><span data-cso=\"1\" data-ceo=\"1\">foo</span><span class=\"marker\">]</span><span class=\"marker\">(</span><span class=\"marker\">/bar</span><span class=\"marker\">)</span></a><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"19", "* foo\n\n  bar", "<ul data-ntype=\"7\" data-mtype=\"1\"><li class=\"node node--block node--expand\" data-ntype=\"8\" data-mtype=\"1\"><span class=\"marker\">* </span><p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"0\" data-ceo=\"0\">foo</span><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><span>bar</span><span class=\"newline\">\n\n</span></p></li></ul>"}, 2, 2},
	{&parseTest{"18", "* foo\n  bar", "<ul data-ntype=\"7\" data-mtype=\"1\"><li class=\"node node--block node--expand\" data-ntype=\"8\" data-mtype=\"1\"><span class=\"marker\">* </span><span data-cso=\"0\" data-ceo=\"0\">foo</span>\n<span>bar</span><span class=\"newline\">\n</span></li></ul>"}, 2, 2},
	{&parseTest{"17", "* foo\n", "<ul data-ntype=\"7\" data-mtype=\"1\"><li class=\"node node--block node--expand\" data-ntype=\"8\" data-mtype=\"1\"><span class=\"marker\">* </span><span data-cso=\"0\" data-ceo=\"0\">foo</span><span class=\"newline\">\n</span></li></ul>"}, 2, 2},
	{&parseTest{"16", "> 一  \n  \n* 二  \n  \n# 三\n", "<blockquote class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">> </span><p data-ntype=\"1\" data-mtype=\"0\"><span>一</span><span class=\"newline\">\n\n</span></p></blockquote><ul data-ntype=\"7\" data-mtype=\"1\"><li class=\"node node--block\" data-ntype=\"8\" data-mtype=\"1\"><span class=\"marker\">* </span><span>二</span><span class=\"newline\">\n</span></li></ul><h1 class=\"node\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\"># </span><span>三</span></h1>"}, 2, 2},
	{&parseTest{"15", "一\n\n二\n\n三\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span>一</span><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"0\" data-ceo=\"0\">二</span><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><span>三</span><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"14", "fo\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"2\" data-ceo=\"2\">fo</span><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"13", "", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"0\" data-ceo=\"0\"></span><span class=\"newline\">\n</span></p>"}, 2, 2},
	{&parseTest{"12", "foo\nbar\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"2\" data-ceo=\"2\">foo</span>\n<span>bar</span><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"11", "> # foo\n", "<blockquote class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">> </span><h1 class=\"node\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\"># </span><span>foo</span></h1></blockquote>"}, 2, 2},
	{&parseTest{"10", "> #\n", "<blockquote class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">> </span><h1 class=\"node\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\">#\n</span></h1></blockquote>"}, 2, 2},
	{&parseTest{"9", "> ## foo\n", "<blockquote class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">> </span><h2 class=\"node\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\">## </span><span>foo</span></h2></blockquote>"}, 2, 2},
	{&parseTest{"8", "-- ---\n", "<div class=\"node node--block node--hr\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">-- ---</span><hr /><span class=\"newline\">\n\n</span></div>"}, 2, 2},
	{&parseTest{"7", "__foo__\n", "<p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">__</span><span>foo</span><span class=\"marker\">__</span></strong><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"6", "`foo`\n", "<p data-ntype=\"1\" data-mtype=\"0\"><code class=\"node node--expand\" data-ntype=\"23\" data-mtype=\"2\"><span class=\"marker\">`</span>foo<span class=\"marker\">`</span></code><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"5", "**foo**\n", "<p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">**</span><span>foo</span><span class=\"marker\">**</span></strong><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"4", "_foo_\n", "<p data-ntype=\"1\" data-mtype=\"0\"><em class=\"node node--expand\" data-ntype=\"13\" data-mtype=\"2\"><span class=\"marker\">_</span><span data-cso=\"1\" data-ceo=\"1\">foo</span><span class=\"marker\">_</span></em><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"3", "> foo\n", "<blockquote class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">> </span><p data-ntype=\"1\" data-mtype=\"0\"><span>foo</span><span class=\"newline\">\n\n</span></p></blockquote>"}, 2, 2},
	{&parseTest{"2", "## foo\n", "<h2 class=\"node node--expand\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">## </span><span>foo</span></h2>"}, 2, 2},
	{&parseTest{"1", "*foo*\n", "<p data-ntype=\"1\" data-mtype=\"0\"><em class=\"node node--expand\" data-ntype=\"13\" data-mtype=\"2\"><span class=\"marker\">*</span><span data-cso=\"1\" data-ceo=\"1\">foo</span><span class=\"marker\">*</span></em><span class=\"newline\">\n\n</span></p>"}, 2, 2},
	{&parseTest{"0", "foo\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"2\" data-ceo=\"2\">foo</span><span class=\"newline\">\n\n</span></p>"}, 2, 2},
}

func TestVditorRenderer(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorRendererTests {
		html, err := luteEngine.RenderVditorDOM(test.from, test.startOffset, test.endOffset)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var vditorOperationTests = []*vditorTest{

	{&parseTest{"5", "**foo**\n", "<p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\">**</span><span>foo</span><span class=\"marker\">**</span></strong><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\" data-cso=\"0\" data-ceo=\"0\">**</span><span class=\"marker\">**</span></strong><span class=\"newline\">\n\n</span></p>"}, 7, 7},
	{&parseTest{"4", "***\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span>***</span><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"0\" data-ceo=\"0\">\n</span><span class=\"newline\">\n\n</span></p>"}, 3, 3},
	{&parseTest{"3", "***\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span>*</span><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"0\" data-ceo=\"0\">**</span><span class=\"newline\">\n\n</span></p>"}, 1, 1},
	{&parseTest{"2", "***foo\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span>***</span><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><span data-cso=\"0\" data-ceo=\"0\">foo</span><span class=\"newline\">\n\n</span></p>"}, 3, 3},
	{&parseTest{"1", "**foo**\n\n**bar**\n", "<p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\">**</span><span>foo</span><span class=\"marker\">**</span></strong><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\">**</span><span>b</span><span class=\"marker\">**</span></strong><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\" data-cso=\"0\" data-ceo=\"0\">**</span><span>ar</span><span class=\"marker\">**</span></strong><span class=\"newline\">\n\n</span></p>"}, 12, 12},
	{&parseTest{"0", "**foobar**\n", "<p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\">**</span><span>foo</span><span class=\"marker\">**</span></strong><span class=\"newline\">\n\n</span></p><p data-ntype=\"1\" data-mtype=\"0\"><strong class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\" data-cso=\"0\" data-ceo=\"0\">**</span><span>bar</span><span class=\"marker\">**</span></strong><span class=\"newline\">\n\n</span></p>"}, 5, 5},
}

func TestVditorOperation(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorOperationTests {
		html, err := luteEngine.VditorOperation(test.from, test.startOffset, test.endOffset, "newline")
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
