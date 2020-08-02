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

var blockRefTests = []parseTest{

	{"6", "(( 12345678  \"text\" ))\n", "<p><a href=\"12345678\">text</a></p>\n"},
	{"5", "(( 12345678  ))\n", "<p><a href=\"12345678\">placeholder</a></p>\n"},
	{"4", "((12345678 text))\n", "<p>((12345678 text))</p>\n"},
	{"3", "((12345678\"text\"))\n", "<p>((12345678&quot;text&quot;))</p>\n"},
	{"2", "((12345678 \"text\"))\n", "<p><a href=\"12345678\">text</a></p>\n"},
	{"1", "((12345678))\n", "<p><a href=\"12345678\">placeholder</a></p>\n"},
	{"0", "((foo))\n", "<p>((foo))</p>\n"},
}

func TestBlockRef(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range blockRefTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
