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

var headingIDTests = []parseTest{

	{"3", "# Same\n\n ## Same\n\n # Same-\n\n ## Test\n\n ## Test-\n\n ## Test\n\n ## Same\n\n ## Test\n", "<h1 id=\"Same\">Same</h1>\n<h2 id=\"Same-\">Same</h2>\n<h1 id=\"Same--\">Same-</h1>\n<h2 id=\"Test\">Test</h2>\n<h2 id=\"Test-\">Test-</h2>\n<h2 id=\"Test--\">Test</h2>\n<h2 id=\"Same---\">Same</h2>\n<h2 id=\"Test---\">Test</h2>\n"},
	{"2", "# Same\n\n ## Same\n\n # Same-\n", "<h1 id=\"Same\">Same</h1>\n<h2 id=\"Same-\">Same</h2>\n<h1 id=\"Same--\">Same-</h1>\n"},
	{"1", "# Same\n\n ## Same\n", "<h1 id=\"Same\">Same</h1>\n<h2 id=\"Same-\">Same</h2>\n"},
	{"0", "### Heading3 {#custom-id}\n", "<h3 id=\"custom-id\">Heading3</h3>\n"},
}

func TestHeadingID(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.HeadingID = true
	luteEngine.RenderOptions.HeadingID = true

	for _, test := range headingIDTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var headingAnchorTests = []parseTest{

	// 为标题添加锚点 https://github.com/b3log/lute/issues/35
	{"40", "# [标题一](https://ld246.com) 二**三**\n", "<h1 id=\"标题一-二三\"><a href=\"https://ld246.com\">标题一</a> 二<strong>三</strong><a id=\"vditorAnchor-标题一-二三\" class=\"vditor-anchor\" href=\"#标题一-二三\"><svg viewBox=\"0 0 16 16\" version=\"1.1\" width=\"16\" height=\"16\"><path fill-rule=\"evenodd\" d=\"M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z\"></path></svg></a></h1>\n"},
	{"39", "# 标题一\n", "<h1 id=\"标题一\">标题一<a id=\"vditorAnchor-标题一\" class=\"vditor-anchor\" href=\"#标题一\"><svg viewBox=\"0 0 16 16\" version=\"1.1\" width=\"16\" height=\"16\"><path fill-rule=\"evenodd\" d=\"M4 9h1v1H4c-1.5 0-3-1.69-3-3.5S2.55 3 4 3h4c1.45 0 3 1.69 3 3.5 0 1.41-.91 2.72-2 3.25V8.59c.58-.45 1-1.27 1-2.09C10 5.22 8.98 4 8 4H4c-.98 0-2 1.22-2 2.5S3 9 4 9zm9-3h-1v1h1c1 0 2 1.22 2 2.5S13.98 12 13 12H9c-.98 0-2-1.22-2-2.5 0-.83.42-1.64 1-2.09V6.25c-1.09.53-2 1.84-2 3.25C6 11.31 7.55 13 9 13h4c1.45 0 3-1.69 3-3.5S14.5 6 13 6z\"></path></svg></a></h1>\n"},
}

func TestHeadingAnchor(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetHeadingID(true)
	luteEngine.RenderOptions.HeadingAnchor = true
	for _, test := range headingAnchorTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
