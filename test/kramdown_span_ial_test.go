// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"testing"

	"github.com/88250/lute"
)

var kramdownSpanIALTests = []parseTest{

	{"5", "![foo](bar){: width=\"80%\" height=\"80%\"}bar", "<p><img src=\"bar\" alt=\"foo\" width=\"80%\" height=\"80%\" />bar</p>\n"},
	{"4", "![foo](bar){: style=\"zoom:80%;\"}bar", "<p><img src=\"bar\" alt=\"foo\" style=\"zoom:80%;\" />bar</p>\n"},
	{"3", "__foo__{: style=\"color: red\"}bar", "<p><strong style=\"color: red\">foo</strong>bar</p>\n"},
	{"2", "**foo**{: style=\"color: red\"}bar", "<p><strong style=\"color: red\">foo</strong>bar</p>\n"},
	{"1", "_foo_{: style=\"color: red\"}bar", "<p><em style=\"color: red\">foo</em>bar</p>\n"},
	{"0", "*foo*{: style=\"color: red\"}bar", "<p><em style=\"color: red\">foo</em>bar</p>\n"},
}

func TestKramdownSpanIALs(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.KramdownIAL = true

	for _, test := range kramdownSpanIALTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
