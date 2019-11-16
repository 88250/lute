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

// +build javascript

package test

import (
	"testing"

	"github.com/b3log/lute"
)

var vditorRendererTests = []*parseTest{

	{"13", "<table><thead><tr><th>abc</th><th>def</th></tr></thead></table>", "<table><thead><tr><th>abc</th><th>def</th></tr></thead></table>"},
	{"12", "<p><del data-marker=\"~~\">Hi</del> Hello, world!</p>", "<p><del data-marker=\"~~\">Hi</del> Hello, world!</p>"},
	{"11", "<p><del data-marker=\"~\">Hi</del> Hello, world!</p>", "<p><del data-marker=\"~\">Hi</del> Hello, world!</p>"},
	{"10", "<ul><li class=\"vditor-task\"><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo<wbr></li></ul>", "<ul><li class=\"vditor-task\"><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo<wbr></li></ul>"},
	{"9", "<ul><li class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" /> foo<wbr></li></ul>", "<ul><li class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" /> foo<wbr></li></ul>"},
	{"8", "> <wbr>", "<blockquote><wbr></blockquote>"},
	{"7", "><wbr>", "<p>><wbr></p>"},
	{"6", "<p>> foo<wbr></p>", "<blockquote><p>foo<wbr></p></blockquote>"},
	{"5", "<p>foo</p><p><wbr><br></p>", "<p>foo</p><p><wbr><br /></p>"},
	{"4", "<ul><li>foo</li></ul><div><wbr><br></div>", "<ul><li>foo</li></ul><p><wbr><br /></p>"},
	{"3", "<p><em data-marker=\"*\">foo<wbr></em></p>", "<p><em data-marker=\"*\">foo<wbr></em></p>"},
	{"2", "<p>foo<wbr></p>", "<p>foo<wbr></p>"},
	{"1", "<p><strong data-marker=\"**\">foo</strong></p>", "<p><strong data-marker=\"**\">foo</strong></p>"},
	{"0", "<p>foo</p>", "<p>foo</p>"},
}

func TestVditorRenderer(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorRendererTests {
		html, err := luteEngine.RenderVditorDOM(test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
