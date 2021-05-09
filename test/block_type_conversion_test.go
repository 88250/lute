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

var cancelListTests = []*parseTest{
	{"0", "<div data-subtype=\"u\" data-node-id=\"20210509095907-rixmte6\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210509103643\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210509103643-es01df0\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210509103643-xr4nn64\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103643\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210509103701-dnmzej3\" data-type=\"NodeList\" class=\"list\" data-subtype=\"u\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210509103644-zfqz75l\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210509103644-uednsdt\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103651\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-node-id=\"20210509103643-xr4nn64\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103643\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-subtype=\"u\" data-node-id=\"20210509103701-dnmzej3\" data-node-index=\"2\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210509103644-zfqz75l\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210509103644-uednsdt\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103651\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestCancelList(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
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
	for _, test := range cancelListTests {
		ovHTML := luteEngine.CancelList(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var ul2olTests = []*parseTest{

}

func TestUL2OL(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
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

var ol2ulTests = []*parseTest{

}

func TestOL2UL(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
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
	for _, test := range ol2ulTests {
		ovHTML := luteEngine.OL2UL(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var p2hTests = []*parseTest{

}

func TestP2H(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
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
