// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package test

import (
	"testing"

	"github.com/88250/lute"
)

var fnTests = []parseTest{

	{"0", "foo[^1]\n[^1]:bar\n    * baz", "<p>foo<sup class=\"footnotes-ref\" id=\"footnotes-ref-1\"><a href=\"#footnotes-def-1\">1</a></sup></p>\n<div class=\"footnotes-defs-div\">\n<ol class=\"footnotes-defs-ol\"><li id=\"footnotes-def-1\"><p>:bar</p>\n<ul>\n<li>baz</li>\n</ul>\n</li>\n</ol>\n</div>"},
}

func TestFootnotes(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range fnTests {
		html, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
