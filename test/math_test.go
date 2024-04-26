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

var mathTests = []parseTest{

	{"14", "$ $\n", "<p>$ $</p>\n"},

	// 解析问题 10$ https://github.com/88250/lute/issues/1
	{"解析问题 10$", "10$\n", "<p>10$</p>\n"},

	{"13", "$$a^2 + b^2 = \\color{red}c^2$$\n## foo\n", "<div class=\"language-math\">a^2 + b^2 = \\color{red}c^2</div>\n<h2 id=\"foo\">foo</h2>\n"},
	{"12", "$$\n", "<div class=\"language-math\"></div>\n"},
	{"11", "$\n", "<p>$</p>\n"},
	{"10", "lu$$a^2 + b^2 = \\color{red}c^2$$te\n", "<p>lu\n<div class=\"language-math\">a^2 + b^2 = \\color{red}c^2</div>\nte</p>\n"},
	{"9", "b$\\color{red}a^2$a\n", "<p>b<span class=\"language-math\">\\color{red}a^2</span>a</p>\n"},
	{"8", "lu$a^2 + b^2 = \\color{red}c^2$1te\n", "<p>lu$a^2 + b^2 = \\color{red}c^2$1te</p>\n"},
	{"7", "lu$1a^2 + b^2 = \\color{red}c^2$te\n", "<p>lu$1a^2 + b^2 = \\color{red}c^2$te</p>\n"},
	{"6", "lu$a^2 + b^2 = \\color{red}c^2$te$a^2$m\n", "<p>lu<span class=\"language-math\">a^2 + b^2 = \\color{red}c^2</span>te<span class=\"language-math\">a^2</span>m</p>\n"},
	{"5", "lu$a^2 + b^2 = \\color{red}c^2$te\n", "<p>lu<span class=\"language-math\">a^2 + b^2 = \\color{red}c^2</span>te</p>\n"},
	{"4", "lu$$a^2 + b^2 = \\color{red}c^2$$te\n", "<p>lu\n<div class=\"language-math\">a^2 + b^2 = \\color{red}c^2</div>\nte</p>\n"},
	{"3", "$$\na^2 + b^2 = \\color{red}c^2\n$$\n", "<div class=\"language-math\">a^2 + b^2 = \\color{red}c^2</div>\n"},
	{"2", "| $a^2 + b^2 = \\color{red}c^2$ | bar |\n| --- | --- |\n| baz | bim |\n", "<table>\n<thead>\n<tr>\n<th><span class=\"language-math\">a^2 + b^2 = \\color{red}c^2</span></th>\n<th>bar</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>baz</td>\n<td>bim</td>\n</tr>\n</tbody>\n</table>\n"},
	{"1", "$a^2 + b^2 = \\color{red}c^2$\n", "<p><span class=\"language-math\">a^2 + b^2 = \\color{red}c^2</span></p>\n"},
	{"0", "$$a^2 + b^2 = \\color{red}c^2$$\n", "<div class=\"language-math\">a^2 + b^2 = \\color{red}c^2</div>\n"},
}

func TestMath(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetHeadingID(true)
	for _, test := range mathTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var inlineMathTests = []parseTest{
	{"0", "$foo$", "<p><span class=\"language-math\">foo</span></p>\n"},
}

func TestInlineMath(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetInlineMath(true)

	for _, test := range inlineMathTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var inlineMathDisabledTests = []parseTest{
	{"0", "$foo$", "<p>$foo$</p>\n"},
}

func TestInlineMathDisabled(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetInlineMath(false)

	for _, test := range inlineMathDisabledTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var inlineMathDigitAllowTests = []parseTest{

	{"1", "$1$2", "<p><span class=\"language-math\">1</span>2</p>\n"},
	{"0", "$1$", "<p><span class=\"language-math\">1</span></p>\n"},
}

func TestInlineMathDigit(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetInlineMathAllowDigitAfterOpenMarker(true)

	for _, test := range inlineMathDigitAllowTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
