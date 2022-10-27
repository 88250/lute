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

var blockRefTests = []parseTest{

	{"4", "((20221026202632-wqhfhhb 'foo'))", "<p>'foo'</p>\n"},
	{"3", "((id \"foo<a>bar\"))", "<p>((id &quot;foo<a>bar&quot;))</p>\n"},
	{"2", "((20201105103725-dd01qas \"foo<a>bar\"))", "<p>\"foo&lt;a&gt;bar\"</p>\n"},
	{"1", "((20201105103725-dd01qas \"$foo$\"))", "<p>\"$foo$\"</p>\n"},
	{"0", "((20201105103725-dd01qas \"思源笔记\"))", "<p>\"思源笔记\"</p>\n"},
}

func TestBlockRef(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.BlockRef = true
	for _, test := range blockRefTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
