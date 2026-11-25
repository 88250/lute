// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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

// An empty list item cannot interrupt a paragraph (CommonMark 5.2, example 285),
// and a marker followed only by spaces is still an empty item (example 282).
// Only a marker directly followed by a newline was recognized as blank; trailing
// spaces or tabs after the marker let the item interrupt the paragraph.
var blankListItemParagraphTests = []parseTest{
	{"bullet-space", "foo\n* \nbar\n", "<p>foo\n*\nbar</p>\n"},
	{"bullet-tab", "foo\n+\t\nbar\n", "<p>foo\n+\t\nbar</p>\n"},
	{"bullet-indented", "foo\n  * \nbar\n", "<p>foo\n*\nbar</p>\n"},
	{"ordered-space", "foo\n1. \nbar\n", "<p>foo\n1.\nbar</p>\n"},
	{"bullet-newline", "foo\n*\n\nfoo\n1.\n", "<p>foo\n*</p>\n<p>foo\n1.</p>\n"},
	{"non-blank-interrupts", "foo\n* bar\n", "<p>foo</p>\n<ul>\n<li>bar</li>\n</ul>\n"},
	{"blank-item-starts-list", "* \nbar\n", "<ul>\n<li></li>\n</ul>\n<p>bar</p>\n"},
}

func TestBlankListItemParagraph(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.RenderOptions.SoftBreak2HardBreak = false

	for _, test := range blankListItemParagraphTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
