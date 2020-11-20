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

var superBlockTests = []parseTest{

	{"6", "start\n{{{\nfoo\n\n> bar\n>\n> baz\n\n{{{\n* list\n  * para\n\nbazz\n}}}\n\npara\n\n{{{\n# foo\n\nbar\n}}}\n\n}}}\nend\n", "<p>start</p>\n<p>foo</p>\n<blockquote>\n<p>bar</p>\n<p>baz</p>\n</blockquote>\n<ul>\n<li>list\n<ul>\n<li>para</li>\n</ul>\n</li>\n</ul>\n<p>bazz</p>\n<p>para</p>\n<h1 id=\"foo\">foo</h1>\n<p>bar</p>\n<p>end</p>\n"},
	{"5", "{{{\n# foo\n\n{{{\nbar\n\nbaz\n}}}\n}}}\n", "<h1 id=\"foo\">foo</h1>\n<p>bar</p>\n<p>baz</p>\n"},
	{"4", "{{{\nfoo\n\n{{{\nbar\n}}}\n\nbaz\n}}}", "<p>foo</p>\n<p>bar</p>\n<p>baz</p>\n"},
	{"3", "{{{\nfoo\n\n* bar\n\n  baz\n}}}", "<p>foo</p>\n<ul>\n<li>\n<p>bar</p>\n<p>baz</p>\n</li>\n</ul>\n"},
	{"2", "{{{\nfoo\n\n* bar\n  * baz\n}}}", "<p>foo</p>\n<ul>\n<li>bar\n<ul>\n<li>baz</li>\n</ul>\n</li>\n</ul>\n"},
	{"1", "{{{\nfoo\n\nbar\n\n}}}", "<p>foo</p>\n<p>bar</p>\n"},
	{"0", "{{{\nfoo\n}}}", "<p>foo</p>\n"},
}

func TestSuperBlock(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SuperBlock = true
	for _, test := range superBlockTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
