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

var imgPathSpaceTests = []parseTest{

	{"1", "![foo](bar baz.png 'title')", "<p><img src=\"bar baz.png\" alt=\"foo\" title=\"title\" /></p>\n"},
	{"0", "![foo](bar baz.png \"title\")", "<p><img src=\"bar baz.png\" alt=\"foo\" title=\"title\" /></p>\n"},
}

func TestImgPathSpace(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.ImgPathAllowSpace = true
	for _, test := range imgPathSpaceTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
