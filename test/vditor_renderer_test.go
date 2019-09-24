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

	"github.com/b3log/lute"
)

var vditorRendererTests = []parseTest{

	{"20", "> *foo\n> bar*\n", "<blockquote>\n<p><em>foo<br />\nbar</em></p>\n</blockquote>\n"},
	//{"19", "*foo*\n", "<p data-ntype=\"1\" data-mtype=\"0\" data-pos-start=\"1:1\" data-pos-end=\"1:5\"><span class=\"node\" data-ntype=\"12\" data-mtype=\"2\" data-pos-start=\"1:1\" data-pos-end=\"1:5\"><span class=\"marker\" data-ntype=\"13\" data-mtype=\"2\" data-pos-start=\"1:1\" data-pos-end=\"1:1\">*</span><em data-ntype=\"12\" data-mtype=\"2\" data-pos-start=\"1:1\" data-pos-end=\"1:5\"><span data-ntype=\"11\" data-mtype=\"2\" data-pos-start=\"1:2\" data-pos-end=\"1:4\">foo</span></em><span class=\"marker\" data-ntype=\"14\" data-mtype=\"2\" data-pos-start=\"1:5\" data-pos-end=\"1:5\">*</span></span></p><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span>"},
	//{"18", "**789*<span></span>", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">*</span><span class=\"node\" data-ntype=\"11\" data-mtype=\"2\"><span class=\"marker\">*</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">789</span></em><span class=\"marker\">*</span></span><span></span><span class=\"newline\">\n\n</span></p>"},
	//{"16", "**foo*\n", "<span data-ntype=\"10\" data-mtype=\"2\">*</span><span class=\"node\" data-ntype=\"11\" data-mtype=\"2\"><span class=\"marker\">*</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></em><span class=\"marker\">*</span></span>"},
	//{"15", "**foo**\n\n<br />\n", "<span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></strong><span class=\"marker\">**</span></span><br />"},
	//{"14", "**foo**\n\nbar\n", "<span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></strong><span class=\"marker\">**</span></span><span data-ntype=\"10\" data-mtype=\"2\">bar</span>"},
	//{"13", "**foo** _bar_\n", "<span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></strong><span class=\"marker\">**</span></span><span data-ntype=\"10\" data-mtype=\"2\"> </span><span class=\"node\" data-ntype=\"11\" data-mtype=\"2\"><span class=\"marker\">_</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">bar</span></em><span class=\"marker\">_</span></span>"},
	//{"12", "[Lute](https://github.com/b3log/lute)", "<span><span class=\"marker\">[</span><a href=\"https://github.com/b3log/lute\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></a><span class=\"marker\">]</span><span class=\"marker\">(</span><span>https://github.com/b3log/lute</span><span class=\"marker\">)</span></span>"},
	//{"11", "Lu\nte\n", "<span data-ntype=\"10\" data-mtype=\"2\">Lu</span><span data-ntype=\"15\" data-mtype=\"2\" /></span><span data-ntype=\"10\" data-mtype=\"2\">te</span>"},
	//{"10", "Lu  \nte\n", "<span data-ntype=\"10\" data-mtype=\"2\">Lu</span><span data-ntype=\"14\" data-mtype=\"2\"></span><span data-ntype=\"10\" data-mtype=\"2\">te</span>"},
	//{"9", "Lu\\\nte\n", "<span data-ntype=\"10\" data-mtype=\"2\">Lu</span><span data-ntype=\"14\" data-mtype=\"2\"></span><span data-ntype=\"10\" data-mtype=\"2\">te</span>"},
	//{"8", "`Lute`\n", "<span><span class=\"marker\">`</span><code data-ntype=\"13\" data-mtype=\"2\">Lute</code><span class=\"marker\">`</span>"},
	//{"7", "**Lute**\n", "<span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></strong><span class=\"marker\">**</span></span>"},
	//{"6", "*Lute*\n", "<span class=\"node\" data-ntype=\"11\" data-mtype=\"2\"><span class=\"marker\">*</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></em><span class=\"marker\">*</span></span>"},
	//{"5", "_Lute_\n", "<span class=\"node\" data-ntype=\"11\" data-mtype=\"2\"><span class=\"marker\">_</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></em><span class=\"marker\">_</span></span>"},
	//{"4", "* Lute\n", "<ul data-ntype=\"5\" data-mtype=\"1\"><li data-ntype=\"6\" data-mtype=\"1\"><span class=\"node\"><span class=\"marker\">* </span></span><p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></p></li></ul>"},
	//{"3", "> Lute\n", "<blockquote><span class=\"node\"><span class=\"marker\" data-ntype=\"5\" data-mtype=\"2\" data-pos-start=\"1:1\" data-pos-end=\"1:2\">&gt;</span></span><p data-ntype=\"1\" data-mtype=\"0\" data-pos-start=\"1:3\" data-pos-end=\"1:7\"><span data-ntype=\"11\" data-mtype=\"2\" data-pos-start=\"1:3\" data-pos-end=\"1:7\">Lute</span><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span></blockquote>"},
	//{"2", "---\n", "<hr data-ntype=\"3\" data-mtype=\"0\" />"},
	//{"1", "## Lute\n", "<h2 data-ntype=\"2\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></h2>"},
	{"0", "Lute\n", "<p data-ntype=\"1\" data-mtype=\"0\" data-pos-start=\"1:1\" data-pos-end=\"1:5\"><span data-ntype=\"11\" data-mtype=\"2\" data-pos-start=\"1:1\" data-pos-end=\"1:5\">Lute</span></p><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span>"},
}

func TestVditorRenderer(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorRendererTests {
		html, err := luteEngine.RenderVditorDOM(test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
