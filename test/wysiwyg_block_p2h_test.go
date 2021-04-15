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

var p2hTests = []*parseTest{

	{"0", "<div data-node-id=\"20210415115149-65rm92t\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210415190606\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"vditor-attr\"></div></div>", "<div data-subtype=\"h1\" data-node-id=\"20210415115149-65rm92t\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20210415190606\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"vditor-attr\"></div></div>"},
}

func TestP2H(t *testing.T) {
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
	for _, test := range p2hTests {
		ovHTML := luteEngine.P2H(test.from, "1")
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}
