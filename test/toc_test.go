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

var tocTests = []parseTest{

	{"0", "[toc]\n\n# 1\n\n## 1.1\n\n# 2\n", "<div class=\"toc-div\"><span class=\"toc-h1\"><a class=\"toc-a\" href=\"#1\">1</a></span><br>&emsp;&emsp;<span class=\"toc-h2\"><a class=\"toc-a\" href=\"#11\">1.1</a></span><br><span class=\"toc-h1\"><a class=\"toc-a\" href=\"#2\">2</a></span><br></div>\n\n<h1 id=\"1\">1</h1>\n<h2 id=\"11\">1.1</h2>\n<h1 id=\"2\">2</h1>\n"},
}

func TestToC(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ToC = true

	for _, test := range tocTests {
		html, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
