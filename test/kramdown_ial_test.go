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

var kramIALTests = []parseTest{

	//{"8", "* foo\n\n  > bar\n  {: id=\"barid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>\n<p>foo</p>\n</li>\n<li>\n<p id=\"fooid\">foo</p>\n</li>\n</ul>\n"},
	{"7", "* > foo\n  {: id=\"fooid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>\n<p>foo</p>\n</li>\n<li>\n<p id=\"fooid\">foo</p>\n</li>\n</ul>\n"},
	{"6", "* foo\n\n* foo\n  {: id=\"fooid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>\n<p>foo</p>\n</li>\n<li>\n<p id=\"fooid\">foo</p>\n</li>\n</ul>\n"},
	{"5", "* foo\n  {: id=\"fooid\"}\n{: id=\"id\"}\n", "<ul id=\"id\">\n<li>foo</li>\n</ul>\n"},
	{"4", "* foo\n{: id=\"fooid\"}\n", "<ul id=\"fooid\">\n<li>foo</li>\n</ul>\n"},
	{"3", "> foo\n> {: id=\"fooid\"}\n>\n> baz\n> {: id=\"bazid\"}\n>\n{: id=\"bqid\"}\n", "<blockquote id=\"bqid\">\n<p id=\"fooid\">foo</p>\n<p id=\"bazid\">baz</p>\n</blockquote>\n"},
	{"2", "> foo\n> {: id=\"fooid\"}\n{: id=\"bqid\"}\n", "<blockquote id=\"bqid\">\n<p id=\"fooid\">foo</p>\n</blockquote>\n"},
	{"1", "> foo\n> {: id=\"fooid\" name=\"bar\"}\n", "<blockquote>\n<p id=\"fooid\" name=\"bar\">foo</p>\n</blockquote>\n"},
	{"0", "foo\n{: id=\"fooid\" class=\"bar\"}\n", "<p id=\"fooid\" class=\"bar\">foo</p>\n"},
}

func TestKramIALs(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.KramdownIAL = true

	for _, test := range kramIALTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
