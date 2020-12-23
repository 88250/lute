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

var markTests = []parseTest{

	{"11", "==foo=t= **bar**\n", "<p>==foo=t= <strong>bar</strong></p>\n"},
	{"10", "*[==foo==*](bar)\n", "<p>*<a href=\"bar\"><mark>foo</mark>*</a></p>\n"},
	{"9", "==[*foo*==](bar)\n", "<p>==<a href=\"bar\"><em>foo</em>==</a></p>\n"},
	{"8", "==[*foo*](bar)==\n", "<p><mark><a href=\"bar\"><em>foo</em></a></mark></p>\n"},
	{"7", "==*[foo](bar)*==\n", "<p><mark><em><a href=\"bar\">foo</a></em></mark></p>\n"},
	{"6", "==*[foo](bar)==*\n", "<p><mark>*<a href=\"bar\">foo</a></mark>*</p>\n"},
	{"5", "==[foo](bar)==\n", "<p><mark><a href=\"bar\">foo</a></mark></p>\n"},
	{"5", "[==foo==](bar)\n", "<p><a href=\"bar\"><mark>foo</mark></a></p>\n"},
	{"4", "=foo=\n", "<p>=foo=</p>\n"},
	{"3", "==**foo==**\n", "<p><mark>**foo</mark>**</p>\n"},
	{"2", "**==foo==**\n", "<p><strong><mark>foo</mark></strong></p>\n"},
	{"1", "==**foo**==\n", "<p><mark><strong>foo</strong></mark></p>\n"},
	{"0", "==foo==\n", "<p><mark>foo</mark></p>\n"},
}

func TestMark(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.Mark = true

	for _, test := range markTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var markDisableTests = []parseTest{

	{"7", "*[==foo==*](bar)\n", "<p>*<a href=\"bar\">==foo==*</a></p>\n"},
	{"6", "==[*foo*==](bar)\n", "<p>==<a href=\"bar\"><em>foo</em>==</a></p>\n"},
	{"5", "==[*foo*](bar)==\n", "<p>==<a href=\"bar\"><em>foo</em></a>==</p>\n"},
	{"4", "==*[foo](bar)*==\n", "<p>==<em><a href=\"bar\">foo</a></em>==</p>\n"},
	{"3", "[==foo==](bar)\n", "<p><a href=\"bar\">==foo==</a></p>\n"},
	{"2", "*=foo*=\n", "<p><em>=foo</em>=</p>\n"},
	{"1", "=foo=\n", "<p>=foo=</p>\n"},
	{"0", "==foo==\n", "<p>==foo==</p>\n"},
}

func TestMarkDisable(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.Mark = false

	for _, test := range markDisableTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
