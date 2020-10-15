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
	"github.com/88250/lute"
	"testing"
)

var linkBaseTests = []parseTest{
	{"2", "![foo](D:\\bar.png)\n", "<p><img src=\"D:\\bar.png\" alt=\"foo\" /></p>\n"},
	{"1", "![foo](bar.png)\n", "<p><img src=\"http://domain.com/path/bar.png\" alt=\"foo\" /></p>\n"},
	{"0", "[foo](bar.png)\n", "<p><a href=\"http://domain.com/path/bar.png\">foo</a></p>\n"},
}

func TestLinkBase(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.LinkBase = "http://domain.com/path/"

	for _, test := range linkBaseTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var linkBasePrefixTests = []parseTest{
	{"0", "[foo](bar.png)\n", "<p><a href=\"prefix:http://domain.com/path/bar.png\">foo</a></p>\n"},
}

func TestLinkBasePrefix(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.LinkBase = "http://domain.com/path/"
	luteEngine.LinkPrefix = "prefix:"

	for _, test := range linkBasePrefixTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}