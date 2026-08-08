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
	"strings"
	"testing"

	"github.com/88250/lute"
)

var strikethroughTests = []parseTest{

	{"2", "*~foo~*bar\n", "<p><em><del>foo</del></em>bar</p>\n"},
	{"1", "~~foo~~", "<p><del>foo</del></p>\n"},
	{"0", "~foo~", "<p><del>foo</del></p>\n"},
}

func TestStrikethrough(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.GFMStrikethrough = true

	for _, test := range strikethroughTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var strikethroughDisabledTests = []parseTest{

	{"2", "~foo~**bar**", "<p>~foo~<strong>bar</strong></p>\n"},
	{"1", "~~foo~~", "<p>~~foo~~</p>\n"},
	{"0", "~foo~", "<p>~foo~</p>\n"},
}

func TestStrikethroughDisabled(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.GFMStrikethrough = false

	for _, test := range strikethroughDisabledTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

func TestFullWidthStrikethrough(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetFullWidthStrikethrough(true)

	tests := []parseTest{
		{"full-width", "～～foo～～", "<p><del>foo</del></p>\n"},
		{"suffix", "～～foo～～111", "<p><del>foo</del>111</p>\n"},
		{"ascii", "~~foo~~111", "<p><del>foo</del>111</p>\n"},
		{"mixed-open", "～～foo~~", "<p><del>foo</del></p>\n"},
		{"mixed-close", "~~foo～～", "<p><del>foo</del></p>\n"},
		{"single", "～foo～", "<p>～foo～</p>\n"},
		{"triple", "～～～foo～～～", "<p>～～～foo～～～</p>\n"},
		{"code", "`～～foo～～`", "<p><code>～～foo～～</code></p>\n"},
	}
	for _, test := range tests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}

	formatted := luteEngine.FormatStr("full-width-format", "～～foo～～")
	if "~~foo~~\n" != formatted {
		t.Fatalf("full-width strikethrough markers were not normalized, got %q", formatted)
	}
}

func TestFullWidthStrikethroughDisabled(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetFullWidthStrikethrough(false)

	html := luteEngine.MarkdownStr("full-width-disabled", "～～foo～～")
	if "<p>～～foo～～</p>\n" != html {
		t.Fatalf("disabled full-width strikethrough was parsed, got %q", html)
	}
}

func TestSpinBlockDOMFullWidthStrikethrough(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetTextMark(true)
	luteEngine.SetSpin(true)
	luteEngine.SetKramdownIAL(true)
	luteEngine.SetFullWidthStrikethrough(true)

	dom := `<div data-node-id="20260808120000-abcdefg" data-node-index="1" data-type="NodeParagraph" class="p"><div contenteditable="true" spellcheck="false">～～foo～～<wbr>111</div><div class="protyle-attr" contenteditable="false">​</div></div>`
	spun := luteEngine.SpinBlockDOM(dom)
	if !strings.Contains(spun, `<span data-type="s">foo</span><wbr>111`) {
		t.Fatalf("full-width strikethrough was not parsed around the caret, got %q", spun)
	}

	luteEngine.SetFullWidthStrikethrough(false)
	spun = luteEngine.SpinBlockDOM(dom)
	if !strings.Contains(spun, `～～foo～～<wbr>111`) {
		t.Fatalf("disabled full-width strikethrough changed the source text, got %q", spun)
	}
}
