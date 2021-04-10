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
	"github.com/88250/lute/ast"
	"testing"

	"github.com/88250/lute"
)

var spinVditorBlockDOMTests = []*parseTest{

	{"1", "<div data-marker=\"*\" data-node-id=\"20210410181814-3ar1jd9\" data-node-index=\"1\" data-type=\"list\" class=\"list\" updated=\"20210410223029\"><div data-marker=\"*\" data-node-id=\"20210410223029-ay4v26n\" data-type=\"listitem\" class=\"li\"><div class=\"vditor-bullet\"></div><div data-node-id=\"20210410223029-f91ozr5\" data-type=\"p\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><wbr>foo</div><div class=\"vditor-attr\"></div></div><div class=\"vditor-attr\"></div></div><div class=\"vditor-attr\"></div></div>", "<div data-marker=\"*\" data-node-id=\"20210410181814-3ar1jd9\" data-node-index=\"1\" data-type=\"list\" class=\"list\" updated=\"20210410223029\"><div data-marker=\"*\" data-node-id=\"20210410223029-ay4v26n\" data-type=\"listitem\" class=\"li\"><div class=\"vditor-bullet\"></div><div data-node-id=\"20210410223029-f91ozr5\" data-type=\"p\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><wbr>foo</div><div class=\"vditor-attr\"></div></div><div class=\"vditor-attr\"></div></div><div class=\"vditor-attr\"></div></div>"},
	{"0", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"p\" class=\"p\" bookmark=\"bm\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"vditor-attr\"><div class=\"vditor-attr--bookmark\">bm</div></div></div><div data-node-id=\"20210408204848-qyy54ha\" data-node-index=\"2\" data-type=\"p\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"vditor-attr\"></div></div>", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"p\" class=\"p\" bookmark=\"bm\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"vditor-attr\"><div class=\"vditor-attr--bookmark\">bm</div></div></div><div data-node-id=\"20210408204848-qyy54ha\" data-node-index=\"2\" data-type=\"p\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"vditor-attr\"></div></div>"},
}

func TestSpinVditorBlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetVditorWYSIWYG(true)
	luteEngine.ParseOptions.Mark = true
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.SetKramdownIAL(true)
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.SetLinkBase("/siyuan/0/测试笔记/")
	luteEngine.SetAutoSpace(false)
	luteEngine.SetSub(true)
	luteEngine.SetSup(true)
	luteEngine.SetGitConflict(true)

	ast.Testing = true
	for _, test := range spinVditorBlockDOMTests {
		html := luteEngine.SpinVditorBlockDOM(test.from)

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
	ast.Testing = false
}
