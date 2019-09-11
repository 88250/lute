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

	{"11", "[Lute](https://github.com/b3log/lute)", "<p><span><span class=\"marker\">[</span><a href=\"https://github.com/b3log/lute\"><span>Lute</span></a><span class=\"marker\">]</span><span class=\"marker\">(</span><span>https://github.com/b3log/lute</span><span class=\"marker\">)</span></span></p>\n"},
	{"10", "Lu\nte\n", "<p><span>Lu</span><span /></span>\n<span>te</span></p>\n"},
	{"9", "Lu  \nte\n", "<p><span>Lu</span><span></span>\n<span>te</span></p>\n"},
	{"8", "Lu\\\nte\n", "<p><span>Lu</span><span></span>\n<span>te</span></p>\n"},
	{"7", "`Lute`\n", "<p><span><span class=\"marker\">`</span><code>Lute</code><span class=\"marker\">`</span></p>\n"},
	{"6", "**Lute**\n", "<p><span class=\"node\"><span class=\"marker\">**</span><strong><span>Lute</span></strong><span class=\"marker\">**</span></p>\n"},
	{"5", "*Lute*\n", "<p><span><span class=\"marker\">*</span><em><span>Lute</span></em><span class=\"marker\">*</span></p>\n"},
	{"4", "* Lute\n", "<ul>\n<li><span>Lute</span></li>\n</ul>\n"},
	{"3", "> Lute\n", "<blockquote>\n<p><span>Lute</span></p>\n</blockquote>\n"},
	{"2", "---\n", "<hr />\n"},
	{"1", "## Lute\n", "<h2><span>Lute</span></h2>\n"},
	{"0", "Lute\n", "<p><span>Lute</span></p>\n"},
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
