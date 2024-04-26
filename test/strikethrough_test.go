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
	"github.com/88250/lute"
	"testing"
)

var strikethroughTests = []parseTest{

	{"1", "~~foo~~", "<p><del>foo</del></p>\n"},
	{"0", "~foo~", "<p><del>foo</del></p>\n"},
}

func TestStrikethrough(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.GFMStrikethrough = true

	for _, test := range strikethroughTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var strikethroughDisabledTests = []parseTest{
	{"1", "~~foo~~", "<p>~~foo~~</p>\n"},
	{"0", "~foo~", "<p>~foo~</p>\n"},
}

func TestStrikethroughDisabled(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.GFMStrikethrough = false

	for _, test := range strikethroughDisabledTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
