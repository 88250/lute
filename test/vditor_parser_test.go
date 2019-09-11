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

var vditorParserTests = []parseTest{

	{"5", "<p><span class=\"node\"><span class=\"marker\">**</span><strong><span>foo</span></strong><span class=\"marker\">**</span></span><span> </span><span class=\"node\"><span class=\"marker\">_</span><em><span>bar</span></em><span class=\"marker\">_</span></span></p>\n", "**foo** _bar_\n"},
	{"4", "<h2><span>Lute</span></h2>\n", "## Lute\n"},
	{"3", "<p><span class=\"node\"><span class=\"marker\">**</span><strong><span>Lute</span></strong><span class=\"marker\">**</span></p>\n", "**Lute**\n"},
	{"2", "<p><span><span class=\"marker\">*</span><em><span>Lute</span></em><span class=\"marker\">*</span></p>\n", "*Lute*\n"},
	{"1", "<p><span class=\"node\"><span class=\"marker\">_</span><em><span>Lute</span></em><span class=\"marker\">_</span></span></p>\n", "_Lute_\n"},
	{"0", "<p><span>Lute</span></p>\n", "Lute\n"},
}

func TestVditorParser(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorParserTests {
		html, err := luteEngine.VditorDOMMarkdown(test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
