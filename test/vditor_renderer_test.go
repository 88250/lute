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

	{"15", "**foo**\n\n<br />\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></strong><span class=\"marker\">**</span></span></p><br />"},
	{"14", "**foo**\n\nbar\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></strong><span class=\"marker\">**</span></span></p><p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">bar</span></p>"},
	{"13", "**foo** _bar_\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></strong><span class=\"marker\">**</span></span><span data-ntype=\"10\" data-mtype=\"2\"> </span><span class=\"node\"><span class=\"marker\">_</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">bar</span></em><span class=\"marker\">_</span></span></p>"},
	{"12", "[Lute](https://github.com/b3log/lute)", "<p data-ntype=\"1\" data-mtype=\"0\"><span><span class=\"marker\">[</span><a href=\"https://github.com/b3log/lute\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></a><span class=\"marker\">]</span><span class=\"marker\">(</span><span>https://github.com/b3log/lute</span><span class=\"marker\">)</span></span></p>"},
	{"11", "Lu\nte\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lu</span><span data-ntype=\"15\" data-mtype=\"2\" /></span><span data-ntype=\"10\" data-mtype=\"2\">te</span></p>"},
	{"10", "Lu  \nte\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lu</span><span data-ntype=\"14\" data-mtype=\"2\"></span><span data-ntype=\"10\" data-mtype=\"2\">te</span></p>"},
	{"9", "Lu\\\nte\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lu</span><span data-ntype=\"14\" data-mtype=\"2\"></span><span data-ntype=\"10\" data-mtype=\"2\">te</span></p>"},
	{"8", "`Lute`\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span><span class=\"marker\">`</span><code data-ntype=\"13\" data-mtype=\"2\">Lute</code><span class=\"marker\">`</span></p>"},
	{"7", "**Lute**\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">**</span><strong data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></strong><span class=\"marker\">**</span></span></p>"},
	{"6", "*Lute*\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">*</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></em><span class=\"marker\">*</span></span></p>"},
	{"5", "_Lute_\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\"><span class=\"marker\">_</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></em><span class=\"marker\">_</span></span></p>"},
	{"4", "* Lute\n", "<ul data-ntype=\"5\" data-mtype=\"1\"><li data-ntype=\"6\" data-mtype=\"1\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></li></ul>"},
	{"3", "> Lute\n", "<span class=\"node\"><span class=\"marker\">&gt;</span><blockquote data-ntype=\"4\" data-mtype=\"1\"><p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></p></blockquote></span>"},
	{"2", "---\n", "<hr data-ntype=\"3\" data-mtype=\"0\" />"},
	{"1", "## Lute\n", "<h2 data-ntype=\"2\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></h2>"},
	{"0", "Lute\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"10\" data-mtype=\"2\">Lute</span></p>"},
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
