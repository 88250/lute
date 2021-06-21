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
	"github.com/88250/lute/ast"
	"testing"

	"github.com/88250/lute"
)

var cancelSuperBlockTests = []*parseTest{
	{"0", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeSuperBlock\" class=\"sb\" data-sb-layout=\"col\"><div data-node-id=\"20210518185646-hjjhl5p\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518185646-pcdktyp\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-node-id=\"20210518185646-hjjhl5p\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518185646-pcdktyp\" data-node-index=\"3\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestCancelSuperBlock(t *testing.T) {
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
	for _, test := range cancelSuperBlockTests {
		ovHTML := luteEngine.CancelSuperBlock(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var blocksMergeSuperBlockTests = []*parseTest{
	{"0", "<div data-node-id=\"20210518185646-hjjhl5p\" data-node-index=\"0\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518185646-pcdktyp\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeSuperBlock\" class=\"sb\" data-sb-layout=\"col\"><div data-node-id=\"20210518185646-hjjhl5p\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518185646-pcdktyp\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestBlocksMergeSuperBlock(t *testing.T) {
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
	for _, test := range blocksMergeSuperBlockTests {
		ovHTML := luteEngine.BlocksMergeSuperBlock(test.from, "col")
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var blocks2tlsTests = []*parseTest{
	{"0", "<div data-node-id=\"20210518185646-hjjhl5p\" data-node-index=\"0\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518185646-pcdktyp\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210518185646-hjjhl5p\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210518191256-m1ij6pn\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210518185646-pcdktyp\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestBlocks2TLs(t *testing.T) {
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
	for _, test := range blocks2tlsTests {
		ovHTML := luteEngine.Blocks2TLs(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var blocks2olsTests = []*parseTest{
	{"0", "<div data-node-id=\"20210518185646-hjjhl5p\" data-node-index=\"0\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518185646-pcdktyp\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-subtype=\"o\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"1.\" data-subtype=\"o\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">1.</div><div data-node-id=\"20210518185646-hjjhl5p\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"2.\" data-subtype=\"o\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">2.</div><div data-node-id=\"20210518191256-m1ij6pn\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"3.\" data-subtype=\"o\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">3.</div><div data-node-id=\"20210518185646-pcdktyp\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestBlocks2OLs(t *testing.T) {
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
	for _, test := range blocks2olsTests {
		ovHTML := luteEngine.Blocks2OLs(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var blocks2ulsTests = []*parseTest{
	{"0", "<div data-node-id=\"20210518185646-hjjhl5p\" data-node-index=\"0\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518185646-pcdktyp\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-subtype=\"u\" data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210518185646-hjjhl5p\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210518191256-m1ij6pn\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210518185646-pcdktyp\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestBlocks2ULs(t *testing.T) {
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
	for _, test := range blocks2ulsTests {
		ovHTML := luteEngine.Blocks2ULs(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var blocks2blockquoteTests = []*parseTest{
	{"0", "<div data-node-id=\"20210604110012-rv2gfx9\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210604155021\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210604155021-047dhjy\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210604155021\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeBlockquote\" class=\"bq\"><div data-node-id=\"20210604110012-rv2gfx9\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210604155021\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210604155021-047dhjy\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210604155021\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestBlocks2Blockquote(t *testing.T) {
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
	for _, test := range blocks2blockquoteTests {
		ovHTML := luteEngine.Blocks2Blockquote(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var blocks2psTests = []*parseTest{
	{"0", "", ""},
}

func TestBlocks2Ps(t *testing.T) {
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
	for _, test := range blocks2psTests {
		ovHTML := luteEngine.Blocks2Ps(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var blocks2hsTests = []*parseTest{
	{"1", "<div data-subtype=\"h2\" data-node-id=\"20210610154900-5r20m2j\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h2\" updated=\"20210610154902\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-subtype=\"h2\" data-node-id=\"20210610154909-wemyt8x\" data-node-index=\"2\" data-type=\"NodeHeading\" class=\"h2\" updated=\"20210610154910\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-subtype=\"h1\" data-node-id=\"20210610154900-5r20m2j\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20210610154902\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-subtype=\"h1\" data-node-id=\"20210610154909-wemyt8x\" data-node-index=\"2\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20210610154910\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"0", "<div data-node-id=\"20210518185646-hjjhl5p\" data-node-index=\"0\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210518185646-pcdktyp\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-subtype=\"h1\" data-node-id=\"20210518185646-hjjhl5p\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-subtype=\"h1\" data-node-id=\"20210518191256-m1ij6pn\" data-node-index=\"2\" data-type=\"NodeHeading\" class=\"h1\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-subtype=\"h1\" data-node-id=\"20210518185646-pcdktyp\" data-node-index=\"3\" data-type=\"NodeHeading\" class=\"h1\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestBlocks2Hs(t *testing.T) {
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
	for _, test := range blocks2hsTests {
		ovHTML := luteEngine.Blocks2Hs(test.from, "1")
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var ul2tlTests = []*parseTest{
	{"0", "<div data-subtype=\"u\" data-node-id=\"20210526172440-zgfgiiv\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210526172454\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210526172454-nryi4wk\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210526172454-me3ybh0\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210526172455\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210526172455-8oc3ap6\" data-type=\"NodeListItem\" class=\"li protyle-wysiwyg--select\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210526172455-8lbabpq\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210526172456\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-subtype=\"t\" data-node-id=\"20210526172440-zgfgiiv\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210526172454\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210526172454-nryi4wk\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210526172454-me3ybh0\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210526172455\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210526172455-8oc3ap6\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210526172455-8lbabpq\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210526172456\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestUL2TL(t *testing.T) {
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
	for _, test := range ul2tlTests {
		ovHTML := luteEngine.UL2TL(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var tl2olTests = []*parseTest{
	{"0", "<div data-subtype=\"t\" data-node-id=\"20210621152530-fu33pvt\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210621152825\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210621152826-p1o0uwy\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210621152826-n4o2vir\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152827\"><div contenteditable=\"true\" spellcheck=\"false\">1</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210621152827-wvqo7lo\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210621152827-ttmqnj4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152827\"><div contenteditable=\"true\" spellcheck=\"false\">2</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\"></div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210621152828-j4hgdt5\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210621152828-pgwkd1s\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152828\"><div contenteditable=\"true\" spellcheck=\"false\">3</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-subtype=\"o\" data-node-id=\"20210621152530-fu33pvt\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210621152825\"><div data-marker=\"1.\" data-subtype=\"o\" data-node-id=\"20210621152826-p1o0uwy\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">1.</div><div data-node-id=\"20210621152826-n4o2vir\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152827\"><div contenteditable=\"true\" spellcheck=\"false\">1</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"2.\" data-subtype=\"o\" data-node-id=\"20210621152827-wvqo7lo\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">2.</div><div data-node-id=\"20210621152827-ttmqnj4\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152827\"><div contenteditable=\"true\" spellcheck=\"false\">2</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"3.\" data-subtype=\"o\" data-node-id=\"20210621152828-j4hgdt5\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">3.</div><div data-node-id=\"20210621152828-pgwkd1s\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152828\"><div contenteditable=\"true\" spellcheck=\"false\">3</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
}

func TestTL2OL(t *testing.T) {
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
	for _, test := range tl2olTests {
		ovHTML := luteEngine.TL2OL(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var cancelListTests = []*parseTest{
	{"1", "<div data-subtype=\"t\" data-node-id=\"20210610155819-dvxb0ws\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210610155823\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210610155824-ibbpfdr\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210610155824-6pcn10x\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210610155825\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210610155825-trhucha\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210610155825-8v3wz67\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210610155825\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-node-id=\"20210610155824-6pcn10x\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210610155825\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210610155825-8v3wz67\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210610155825\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
	{"0", "<div data-subtype=\"u\" data-node-id=\"20210509095907-rixmte6\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210509103643\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210509103643-es01df0\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210509103643-xr4nn64\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103643\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-node-id=\"20210509103701-dnmzej3\" data-type=\"NodeList\" class=\"list\" data-subtype=\"u\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210509103644-zfqz75l\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210509103644-uednsdt\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103651\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>", "<div data-node-id=\"20210509103643-xr4nn64\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103643\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div data-subtype=\"u\" data-node-id=\"20210509103701-dnmzej3\" data-node-index=\"2\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210509103644-zfqz75l\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210509103644-uednsdt\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103651\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div>"},
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
