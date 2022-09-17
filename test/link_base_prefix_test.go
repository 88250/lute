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

var linkBaseTests = []parseTest{
	{"3", "[foo][^label]\n[^label]: bar\n", "<p><sup class=\"footnotes-ref\" id=\"footnotes-ref-1\"><a href=\"http://domain.com/path/#footnotes-def-1\">1</a></sup></p>\n<div class=\"footnotes-defs-div\"><hr class=\"footnotes-defs-hr\" />\n<ol class=\"footnotes-defs-ol\"><li id=\"footnotes-def-1\"><p>bar <a href=\"#footnotes-ref-1\" class=\"vditor-footnotes__goto-ref\">↩</a></p>\n</li>\n</ol></div>"},
	{"2", "![foo](D:\\bar.png)\n", "<p><img src=\"D:%5Cbar.png\" alt=\"foo\" /></p>\n"},
	{"1", "![foo](bar.png)\n", "<p><img src=\"http://domain.com/path/bar.png\" alt=\"foo\" /></p>\n"},
	{"0", "[foo](bar.png)\n", "<p><a href=\"http://domain.com/path/bar.png\">foo</a></p>\n"},
}

func TestLinkBase(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetSup(true)
	luteEngine.RenderOptions.LinkBase = "http://domain.com/path/"

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
	luteEngine.RenderOptions.LinkBase = "http://domain.com/path/"
	luteEngine.RenderOptions.LinkPrefix = "prefix:"

	for _, test := range linkBasePrefixTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
