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

	{"15", "**foo**\n\n<br />\n", "<p><span class=\"node\"><span class=\"marker\">**</span><strong><span>foo</span></strong><span class=\"marker\">**</span></span></p><br />"},
	{"14", "**foo**\n\nbar\n", "<p><span class=\"node\"><span class=\"marker\">**</span><strong><span>foo</span></strong><span class=\"marker\">**</span></span></p><p><span>bar</span></p>"},
	{"13", "**foo** _bar_\n", "<p><span class=\"node\"><span class=\"marker\">**</span><strong><span>foo</span></strong><span class=\"marker\">**</span></span><span> </span><span class=\"node\"><span class=\"marker\">_</span><em><span>bar</span></em><span class=\"marker\">_</span></span></p>"},
	{"12", "[Lute](https://github.com/b3log/lute)", "<p><span><span class=\"marker\">[</span><a href=\"https://github.com/b3log/lute\"><span>Lute</span></a><span class=\"marker\">]</span><span class=\"marker\">(</span><span>https://github.com/b3log/lute</span><span class=\"marker\">)</span></span></p>"},
	{"11", "Lu\nte\n", "<p><span>Lu</span><span /></span><span>te</span></p>"},
	{"10", "Lu  \nte\n", "<p><span>Lu</span><span></span><span>te</span></p>"},
	{"9", "Lu\\\nte\n", "<p><span>Lu</span><span></span><span>te</span></p>"},
	{"8", "`Lute`\n", "<p><span><span class=\"marker\">`</span><code>Lute</code><span class=\"marker\">`</span></p>"},
	{"7", "**Lute**\n", "<p><span class=\"node\"><span class=\"marker\">**</span><strong><span>Lute</span></strong><span class=\"marker\">**</span></span></p>"},
	{"6", "*Lute*\n", "<p><span class=\"node\"><span class=\"marker\">*</span><em><span>Lute</span></em><span class=\"marker\">*</span></span></p>"},
	{"5", "_Lute_\n", "<p><span class=\"node\"><span class=\"marker\">_</span><em><span>Lute</span></em><span class=\"marker\">_</span></span></p>"},
	{"4", "* Lute\n", "<ul><li><span>Lute</span></li></ul>"},
	{"3", "> Lute\n", "<blockquote><p><span>Lute</span></p></blockquote>"},
	{"2", "---\n", "<hr />"},
	{"1", "## Lute\n", "<h2><span>Lute</span></h2>"},
	{"0", "Lute\n", "<p><span>Lute</span></p>"},
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
