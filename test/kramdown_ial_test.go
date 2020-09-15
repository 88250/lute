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

var kramBlockIALTests = []parseTest{

	{"19", "* foo\n\n  bar\n  {: id=\"barid\"}\n\n  > baz\n  {: id=\"bazid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>\n<p>foo</p>\n<p id=\"barid\">bar</p>\n<blockquote id=\"bazid\">\n<p>baz</p>\n</blockquote>\n</li>\n</ul>\n"},
	{"18", "* foo\n\n  bar\n  {: id=\"barid\"}\n\n  baz\n{: id=\"id\"}", "<ul id=\"id\">\n<li>\n<p>foo</p>\n<p id=\"barid\">bar</p>\n<p>baz</p>\n</li>\n</ul>\n"},
	{"17", "> * foo\n>   * bar\n>     * baz\n>\n>       bazz\n>       {: id=\"bazzid\"}\n>     {: id=\"bazid\"}\n>   {: id=\"barid\"}\n> {: id=\"fooid\"}\n{: id=\"id\"}", "<blockquote id=\"id\">\n<ul id=\"fooid\">\n<li>foo\n<ul id=\"barid\">\n<li>bar\n<ul id=\"bazid\">\n<li>\n<p>baz</p>\n<p id=\"bazzid\">bazz</p>\n</li>\n</ul>\n</li>\n</ul>\n</li>\n</ul>\n</blockquote>\n"},
	{"16", "> foo\n> {: id=\"fooid\"}\n> * bar\n> {: id=\"barid\"}\n{: id=\"id\"}", "<blockquote id=\"id\">\n<p id=\"fooid\">foo</p>\n<ul id=\"barid\">\n<li>bar</li>\n</ul>\n</blockquote>\n"},
	{"15", "> foo\n>\n> * bar\n> {: id=\"barid\"}\n{: id=\"id\"}", "<blockquote id=\"id\">\n<p>foo</p>\n<ul id=\"barid\">\n<li>bar</li>\n</ul>\n</blockquote>\n"},
	{"14", "foo\n{: id=\"fooid\"}\nbar\n{: id=\"barid\"}", "<p id=\"fooid\">foo</p>\n<p id=\"barid\">bar</p>\n"},
	{"13", "foo\n{: id=\"fooid\"}\nbar", "<p id=\"fooid\">foo</p>\n<p>bar</p>\n"},
	{"12", "* foo\n\n  > bar\n  {: id=\"bqid\"}\n  > baz\n  > {: id=\"bazid\"}\n* baz\n  {: id=\"bazid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>\n<p>foo</p>\n<blockquote id=\"bqid\">\n<p>bar</p>\n</blockquote>\n<blockquote>\n<p id=\"bazid\">baz</p>\n</blockquote>\n</li>\n<li>\n<p id=\"bazid\">baz</p>\n</li>\n</ul>\n"},
	{"11", "* foo\n  * bar\n  * baz\n  {: id=\"subid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>foo\n<ul id=\"subid\">\n<li>bar</li>\n<li>baz</li>\n</ul>\n</li>\n</ul>\n"},
	{"10", "* foo\n  * bar\n  * baz\n  {: id=\"subid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>foo\n<ul id=\"subid\">\n<li>bar</li>\n<li>baz</li>\n</ul>\n</li>\n</ul>\n"},
	{"9", "* foo\n\n  > bar\n  > {: id=\"barid\"}\n  {: id=\"bqid\"}\n\n  baz\n{: id=\"id\"}", "<ul id=\"id\">\n<li>\n<p>foo</p>\n<blockquote id=\"bqid\">\n<p id=\"barid\">bar</p>\n</blockquote>\n<p>baz</p>\n</li>\n</ul>\n"},
	{"8", "* foo\n\n  > bar\n  {: id=\"bqid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>\n<p>foo</p>\n<blockquote id=\"bqid\">\n<p>bar</p>\n</blockquote>\n</li>\n</ul>\n"},
	{"7", "* > foo\n  {: id=\"bqid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>\n<blockquote id=\"bqid\">\n<p>foo</p>\n</blockquote>\n</li>\n</ul>\n"},
	{"6", "* foo\n\n* bar\n  {: id=\"barid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>\n<p>foo</p>\n</li>\n<li>\n<p id=\"barid\">bar</p>\n</li>\n</ul>\n"},
	{"5", "* foo\n  {: id=\"fooid\"}\n{: id=\"id\"}", "<ul id=\"id\">\n<li>foo</li>\n</ul>\n"},
	{"4", "* foo\n{: id=\"fooid\"}", "<ul id=\"fooid\">\n<li>foo</li>\n</ul>\n"},
	{"3", "> foo\n> {: id=\"fooid\"}\n>\n> baz\n> {: id=\"bazid\"}\n>\n{: id=\"bqid\"}", "<blockquote id=\"bqid\">\n<p id=\"fooid\">foo</p>\n<p id=\"bazid\">baz</p>\n</blockquote>\n"},
	{"2", "> foo\n> {: id=\"fooid\"}\n{: id=\"bqid\"}", "<blockquote id=\"bqid\">\n<p id=\"fooid\">foo</p>\n</blockquote>\n"},
	{"1", "> foo\n> {: id=\"fooid\" name=\"bar\"}", "<blockquote>\n<p id=\"fooid\" name=\"bar\">foo</p>\n</blockquote>\n"},
	{"0", "foo\n{: id=\"fooid\" class=\"bar\"}", "<p id=\"fooid\" class=\"bar\">foo</p>\n"},
}

func TestKramBlockIALs(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.KramdownIAL = true

	for _, test := range kramBlockIALTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
