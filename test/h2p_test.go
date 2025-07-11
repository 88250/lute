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

var html2BlockDOMTests = []parseTest{

	{"0", "<table><tr><td><span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"foo\" contenteditable=\"false\" class=\"render-node\" data-render=\"true\">​<span class=\"katex\"><span class=\"katex-html\" aria-hidden=\"true\"><span class=\"base\"><span class=\"strut\" style=\"height:0.8889em;vertical-align:-0.1944em;\"></span><span class=\"mord mathnormal\" style=\"margin-right:0.10764em;\">f</span><span class=\"mord mathnormal\">oo</span></span></span></span></span>​</td></tr><tr><td></td></tr></table>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeTable\" class=\"table\"><div contenteditable=\"false\"><table contenteditable=\"true\" spellcheck=\"false\"><colgroup><col /></colgroup><thead><tr><th><span data-type=\"inline-math\" data-subtype=\"math\" data-content=\"foo\" contenteditable=\"false\" class=\"render-node\"></span>\u200b</th></tr></thead><tbody><tr><td></td></tr></tbody></table><div class=\"protyle-action__table\"><div class=\"table__resize\"></div><div class=\"table__select\"></div></div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
}

func TestHTML2BlockDOM(t *testing.T) {
	luteEngine := lute.New()

	ast.Testing = true
	for _, test := range html2BlockDOMTests {
		result := luteEngine.HTML2BlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}
