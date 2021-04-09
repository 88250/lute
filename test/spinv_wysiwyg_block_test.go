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

	{"0", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"p\" class=\"p\">foo</div><div data-node-id=\"20210408204848-qyy54ha\" data-node-index=\"2\" data-type=\"p\" class=\"p\">bar</div>", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"p\" class=\"p\">foo</div><div data-node-id=\"20210408204848-qyy54ha\" data-node-index=\"2\" data-type=\"p\" class=\"p\">bar</div>"},
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
