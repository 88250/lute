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

var renderListStyleTests = []parseTest{
	{"4", "1) foo\n", "<ol data-style=\"1)\">\n<li>foo</li>\n</ol>\n"},
	{"3", "1. foo\n", "<ol data-style=\"1.\">\n<li>foo</li>\n</ol>\n"},
	{"2", "- foo\n", "<ul data-style=\"-\">\n<li>foo</li>\n</ul>\n"},
	{"1", "+ foo\n", "<ul data-style=\"+\">\n<li>foo</li>\n</ul>\n"},
	{"0", "* foo\n", "<ul data-style=\"*\">\n<li>foo</li>\n</ul>\n"},
}

func TestRenderListStyle(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.RenderListStyle = true
	for _, test := range renderListStyleTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var renderListStyleVditorTests = []parseTest{
	{"4", "1) foo\n", "<ol data-tight=\"true\" data-marker=\"1)\" data-block=\"0\" data-style=\"1)\"><li data-marker=\"1)\">foo</li></ol>"},
	{"3", "1. foo\n", "<ol data-tight=\"true\" data-marker=\"1.\" data-block=\"0\" data-style=\"1.\"><li data-marker=\"1.\">foo</li></ol>"},
	{"2", "- foo\n", "<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\" data-style=\"-\"><li data-marker=\"-\">foo</li></ul>"},
	{"1", "+ foo\n", "<ul data-tight=\"true\" data-marker=\"+\" data-block=\"0\" data-style=\"+\"><li data-marker=\"+\">foo</li></ul>"},
	{"0", "* foo\n", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\" data-style=\"*\"><li data-marker=\"*\">foo</li></ul>"},
}

func TestRenderListStyleVditor(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.RenderListStyle = true
	for _, test := range renderListStyleVditorTests {
		html := luteEngine.SpinVditorDOM(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var renderListStyleVditorIRTests = []parseTest{
	{"4", "1) foo\n", "<ol data-tight=\"true\" data-marker=\"1)\" data-block=\"0\" data-style=\"1)\"><li data-marker=\"1)\">foo</li></ol>"},
	{"3", "1. foo\n", "<ol data-tight=\"true\" data-marker=\"1.\" data-block=\"0\" data-style=\"1.\"><li data-marker=\"1.\">foo</li></ol>"},
	{"2", "- foo\n", "<ul data-tight=\"true\" data-marker=\"-\" data-block=\"0\" data-style=\"-\"><li data-marker=\"-\">foo</li></ul>"},
	{"1", "+ foo\n", "<ul data-tight=\"true\" data-marker=\"+\" data-block=\"0\" data-style=\"+\"><li data-marker=\"+\">foo</li></ul>"},
	{"0", "* foo\n", "<ul data-tight=\"true\" data-marker=\"*\" data-block=\"0\" data-style=\"*\"><li data-marker=\"*\">foo</li></ul>"},
}

func TestRenderListStyleVditorIR(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.RenderListStyle = true
	for _, test := range renderListStyleVditorIRTests {
		html := luteEngine.SpinVditorIRDOM(test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
