// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"github.com/88250/lute"
	"testing"
)

var relativePathLinkTests = []parseTest{

	{"1", "![foo](bar.png)\n", "<p><img src=\"http://domain.com/path/bar.png\" alt=\"foo\" /></p>\n"},
	{"0", "[foo](bar.png)\n", "<p><a href=\"http://domain.com/path/bar.png\">foo</a></p>\n"},
}

func TestRelativePathLink(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.LinkBase = "http://domain.com/path/"

	for _, test := range relativePathLinkTests {
		html, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("test case [%s] unexpected: %s", test.name, err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
