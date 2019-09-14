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

package test

import (
	"testing"
)

var vditorParserTests = []parseTest{

	{"17", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lu</span><span data-caret></span><span data-ntype=\"10\" data-mtype=\"2\">te</span><span class=\"newline\">\n\n</span></p>", "Lu\x07te\n"},
	{"16", "<span data-ntype=\"10\" data-mtype=\"2\">*</span><span class=\"node\"><span class=\"marker\">*</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></em><span class=\"marker\">*</span></span>", "**foo*\n"},
	{"15", "<span class=\"node\"><span class=\"marker\">&gt;</span><blockquote data-ntype=\"4\" data-mtype=\"1\"><p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></p></blockquote></span>", "> foo\n"},
	{"14", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></strong><span class=\"marker\">**</span></span></p><p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">bar</span></p>", "**foo**\n\nbar\n"},
	{"13", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></strong><span class=\"marker\">**</span></span><span data-ntype=\"10\" data-mtype=\"2\"> </span><span class=\"node\"><span class=\"marker\">_</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">bar</span></em><span class=\"marker\">_</span></span></p>", "**foo** _bar_\n"},
	{"12", "<p data-ntype=\"1\" data-mtype=\"0\"><span><span class=\"marker\">[</span><a href=\"https://github.com/b3log/lute\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></a><span class=\"marker\">]</span><span class=\"marker\">(</span><span>https://github.com/b3log/lute</span><span class=\"marker\">)</span></span></p>", "[Lute](https://github.com/b3log/lute)"},
	{"11", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lu</span><span data-ntype=\"15\" data-mtype=\"2\" /></span><span data-ntype=\"10\" data-mtype=\"2\">te</span></p>","Lu\nte\n" },
	{"10", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lu</span><span data-ntype=\"14\" data-mtype=\"2\"></span><span data-ntype=\"10\" data-mtype=\"2\">te</span></p>", "Lu  \nte\n"},
	{"9", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lu</span><span data-ntype=\"14\" data-mtype=\"2\"></span><span data-ntype=\"10\" data-mtype=\"2\">te</span></p>", "Lu\\\nte\n"},
	{"8", "`Lute`\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span><span class=\"marker\">`</span><code data-ntype=\"13\" data-mtype=\"2\">Lute</code><span class=\"marker\">`</span></p>"},
	{"7", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></strong><span class=\"marker\">**</span></span></p>", "**Lute**\n"},
	{"6", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">*</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></em><span class=\"marker\">*</span></span></p>", "*Lute*\n"},
	{"5", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">_</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></em><span class=\"marker\">_</span></span></p>", "_Lute_\n"},
	{"4", "<ul data-ntype=\"5\" data-mtype=\"1\"><li data-ntype=\"6\" data-mtype=\"1\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></li></ul>", "* Lute\n"},
	{"3", "<span class=\"node\"><span class=\"marker\">&gt;</span><blockquote data-ntype=\"4\" data-mtype=\"1\"><p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></p></blockquote></span>", "> Lute\n"},
	{"2", "<hr data-ntype=\"3\" data-mtype=\"0\" />", "---\n"},
	{"1", "<h2 data-ntype=\"2\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></h2>", "## Lute\n"},
	{"0", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></p>", "Lute\n"},
}

func TestVditorParser(t *testing.T) {
	//luteEngine := lute.New()
	//
	//for _, test := range vditorParserTests {
	//	html, err := luteEngine.VditorDOMMarkdown(test.from)
	//	if nil != err {
	//		t.Fatalf("unexpected: %s", err)
	//	}
	//
	//	if test.to != html {
	//		t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
	//	}
	//}
}
