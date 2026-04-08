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
	"github.com/88250/lute/ast"
)

// 测试 HTML 渲染器开启 DataTask 后在 checkbox 上输出 data-task 属性

var dataTaskHTMLTests = []parseTest{
	// 未勾选任务项
	{"0", "- [ ] foo\n", "<ul>\n<li class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" data-task=\" \" /> foo</li>\n</ul>\n"},
	// 已勾选任务项 (x)
	{"1", "- [x] bar\n", "<ul>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" data-task=\"X\" /> bar</li>\n</ul>\n"},
	// 已勾选任务项 (X)
	{"2", "- [X] baz\n", "<ul>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" data-task=\"X\" /> baz</li>\n</ul>\n"},
	// 混合列表
	{"3", "- [ ] foo\n- [x] bar\n", "<ul>\n<li class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" data-task=\" \" /> foo</li>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" data-task=\"X\" /> bar</li>\n</ul>\n"},
	// 嵌套任务列表
	{"4", "- [x] foo\n  - [ ] bar\n  - [x] baz\n- [ ] bim\n", "<ul>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" data-task=\"X\" /> foo\n<ul>\n<li class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" data-task=\" \" /> bar</li>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" data-task=\"X\" /> baz</li>\n</ul>\n</li>\n<li class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" data-task=\" \" /> bim</li>\n</ul>\n"},
	// 有序任务列表
	{"5", "1. [x] ordered\n2. [ ] item\n", "<ol>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" data-task=\"X\" /> ordered</li>\n<li class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" data-task=\" \" /> item</li>\n</ol>\n"},
}

func TestDataTaskHTML(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetDataTask(true)

	for _, test := range dataTaskHTMLTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

// 测试 HTML 渲染器关闭 DataTask 时不输出 data-task 属性（默认行为）

var dataTaskDisabledHTMLTests = []parseTest{
	{"0", "- [ ] foo\n", "<ul>\n<li class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" /> foo</li>\n</ul>\n"},
	{"1", "- [x] bar\n", "<ul>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" /> bar</li>\n</ul>\n"},
}

func TestDataTaskDisabledHTML(t *testing.T) {
	luteEngine := lute.New()
	// DataTask 默认为 false，无需显式设置

	for _, test := range dataTaskDisabledHTMLTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

// 测试 Protyle 渲染器（Md2BlockDOM）开启 DataTask 后输出 data-task 属性

var dataTaskMd2BlockDOMTests = []parseTest{
	// 未勾选
	{"0", "* [ ] foo", "<div data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"t\" data-task=\" \" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\" draggable=\"true\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	// 已勾选
	{"1", "* [x] bar", "<div data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"t\" data-task=\" \" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li protyle-task--done\"><div class=\"protyle-action protyle-action--task\" draggable=\"true\"><svg><use xlink:href=\"#iconCheck\"></use></svg></div><div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
}

func TestDataTaskMd2BlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetDataTask(true)

	ast.Testing = true
	for _, test := range dataTaskMd2BlockDOMTests {
		result := luteEngine.Md2BlockDOM(test.from, true)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, result, test.from)
		}
	}
	ast.Testing = false
}

// 测试关闭 GFMTaskListItemClass 后 DataTask 仍然可以正常输出

var dataTaskNoClassHTMLTests = []parseTest{
	{"0", "- [ ] foo\n", "<ul>\n<li><input disabled=\"\" type=\"checkbox\" data-task=\" \" /> foo</li>\n</ul>\n"},
	{"1", "- [x] bar\n", "<ul>\n<li><input checked=\"\" disabled=\"\" type=\"checkbox\" data-task=\"X\" /> bar</li>\n</ul>\n"},
}

func TestDataTaskNoClassHTML(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetDataTask(true)
	luteEngine.RenderOptions.GFMTaskListItemClass = ""

	for _, test := range dataTaskNoClassHTMLTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

// 测试自定义任务标记字符的 data-task 属性保留

var dataTaskCustomMarkerTests = []parseTest{
	// "/" 标记 (自定义进行中状态)
	{"0", "- [/] in progress\n", "<ul>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" data-task=\"/\" /> in progress</li>\n</ul>\n"},
	// "?" 标记 (自定义待确认状态)
	{"1", "- [?] maybe\n", "<ul>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" data-task=\"?\" /> maybe</li>\n</ul>\n"},
	// "!" 标记 (自定义重要状态)
	{"2", "- [!] important\n", "<ul>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" data-task=\"!\" /> important</li>\n</ul>\n"},
}

func TestDataTaskCustomMarker(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetArbitraryTaskListItemMarker(true)
	luteEngine.SetDataTask(true)

	for _, test := range dataTaskCustomMarkerTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
