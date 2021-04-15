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

var ul2olTests = []*parseTest{

	{"0", "<div data-subtype=\"u\" data-node-id=\"20210414223654-vfqydjh\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210415082227-m67yq1v\" data-type=\"NodeListItem\" class=\"li\"><div class=\"vditor-bullet\"></div><div data-node-id=\"20210415082227-z9mgkh5\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"vditor-attr\"></div></div><div class=\"vditor-attr\"></div></div><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210415091213-c387rm0\" data-type=\"NodeListItem\" class=\"li\"><div class=\"vditor-bullet\"></div><div data-node-id=\"20210415091222-knbamrt\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"vditor-attr\"></div></div><div class=\"vditor-attr\"></div></div><div class=\"vditor-attr\"></div></div>", "<div data-subtype=\"o\" data-node-id=\"20210414223654-vfqydjh\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"1.\" data-subtype=\"o\" data-node-id=\"20210415082227-m67yq1v\" data-type=\"NodeListItem\" class=\"li\"><div class=\"vditor-bullet vditor-bullet--order\">1.</div><div data-node-id=\"20210415082227-z9mgkh5\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"vditor-attr\"></div></div><div class=\"vditor-attr\"></div></div><div data-marker=\"2.\" data-subtype=\"o\" data-node-id=\"20210415091213-c387rm0\" data-type=\"NodeListItem\" class=\"li\"><div class=\"vditor-bullet vditor-bullet--order\">2.</div><div data-node-id=\"20210415091222-knbamrt\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"vditor-attr\"></div></div><div class=\"vditor-attr\"></div></div><div class=\"vditor-attr\"></div></div>"},
}

func TestUL2OL(t *testing.T) {
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
	for _, test := range ul2olTests {
		ovHTML := luteEngine.UL2OL(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}
