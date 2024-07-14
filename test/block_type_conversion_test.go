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

var cancelSuperBlockTests = []*parseTest{
	{"0", "<div data-node-id=\"20060102150405-1a2b3c4\" data-node-index=\"1\" data-type=\"NodeSuperBlock\" class=\"sb\" data-sb-layout=\"col\"><div data-node-id=\"20210518185646-hjjhl5p\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20210518185646-pcdktyp\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "<div data-node-id=\"20210518185646-hjjhl5p\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20210518185646-pcdktyp\" data-node-index=\"3\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
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

var blocks2psTests = []*parseTest{
	{"2", "<div data-subtype=\"u\" data-node-id=\"20231219120217-oa6yl6b\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20231219120506\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20231219120506-qjotc6n\" data-type=\"NodeListItem\" class=\"li\" updated=\"20231219120506\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20231219120506-dckai7w\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219120506\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div data-node-id=\"20231219120509-fjj9w2d\" data-type=\"NodeList\" class=\"list\" data-subtype=\"u\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20231219120509-fvl0qrv\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20231219120509-4odpk5z\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219120510\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div contenteditable=\"false\" class=\"protyle-attr\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div data-subtype=\"u\" data-node-id=\"20231219120223-7xqch02\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20231219120511\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20231219120511-ybzcnet\" data-type=\"NodeListItem\" class=\"li\" updated=\"20231219120511\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20231219120511-581gs9e\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219120511\"><div contenteditable=\"true\" spellcheck=\"false\">foo1</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div data-node-id=\"20231219120513-4fyeovb\" data-type=\"NodeList\" class=\"list\" data-subtype=\"u\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20231219120513-tb4gg96\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20231219120513-r5ao3jn\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219120514\"><div contenteditable=\"true\" spellcheck=\"false\">bar1</div><div contenteditable=\"false\" class=\"protyle-attr\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "<div data-node-id=\"20231219120506-dckai7w\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219120506\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-subtype=\"u\" data-node-id=\"20231219120509-fjj9w2d\" data-node-index=\"2\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20231219120509-fvl0qrv\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20231219120509-4odpk5z\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219120510\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20231219120511-581gs9e\" data-node-index=\"3\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219120511\"><div contenteditable=\"true\" spellcheck=\"false\">foo1</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-subtype=\"u\" data-node-id=\"20231219120513-4fyeovb\" data-node-index=\"4\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20231219120513-tb4gg96\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20231219120513-r5ao3jn\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219120514\"><div contenteditable=\"true\" spellcheck=\"false\">bar1</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"1", "<div data-subtype=\"u\" data-node-id=\"20231219110700-cjhm3pd\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20231219113459\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20231219110702-z9t33x1\" data-type=\"NodeListItem\" class=\"li\" updated=\"20231219113459\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20231219113527-focntxe\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219113530\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div data-subtype=\"u\" data-node-id=\"20231219114338-lyf6uer\" data-node-index=\"2\" data-type=\"NodeList\" class=\"list\" updated=\"20231219114340\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20231219114340-d4yq4f9\" data-type=\"NodeListItem\" class=\"li\" updated=\"20231219114340\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20231219114340-ij6r9wa\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219114341\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "<div data-node-id=\"20231219113527-focntxe\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219113530\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20231219114340-ij6r9wa\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20231219114341\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"0", "<div data-node-id=\"20220426231411-0xk2nas\" data-node-index=\"0\" data-type=\"NodeBlockquote\" class=\"bq\" updated=\"20220426233927\"><div data-node-id=\"20220426233908-2ocolr9\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20220426233927\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div data-node-id=\"20220426231415-p90jao5\" data-node-index=\"1\" data-type=\"NodeBlockquote\" class=\"bq\" updated=\"20220426233910\"><div data-node-id=\"20220426233910-fzxzuf7\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20220426233910\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div data-node-id=\"20220426233911-bfxns1v\" data-node-index=\"2\" data-type=\"NodeBlockquote\" class=\"bq\" updated=\"20220426233913\"><div data-node-id=\"20220426233912-bvcbp40\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20220426233913\"><div contenteditable=\"true\" spellcheck=\"false\">baz</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "<div data-node-id=\"20220426233908-2ocolr9\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20220426233927\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20220426233910-fzxzuf7\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20220426233910\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20220426233912-bvcbp40\" data-node-index=\"3\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20220426233913\"><div contenteditable=\"true\" spellcheck=\"false\">baz</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
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

	{"3", "<div data-node-id=\"20240714101759-nc7r7j3\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20240714102018\"><div contenteditable=\"true\" spellcheck=\"false\">\t123</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "<div data-subtype=\"h1\" data-node-id=\"20240714101759-nc7r7j3\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20240714102018\"><div contenteditable=\"true\" spellcheck=\"false\">123</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"2", "<div data-node-id=\"20240702163651-4ww2c7k\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20240702163654\"><div contenteditable=\"true\" spellcheck=\"false\">foo\nbar</div><div class=\"protyle-attr\" contenteditable=\"false\">​</div></div>", "<div data-subtype=\"h1\" data-node-id=\"20240702163651-4ww2c7k\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20240702163654\"><div contenteditable=\"true\" spellcheck=\"false\">foobar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"1", "<div data-subtype=\"h2\" data-node-id=\"20210610154900-5r20m2j\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h2\" updated=\"20210610154902\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-subtype=\"h2\" data-node-id=\"20210610154909-wemyt8x\" data-node-index=\"2\" data-type=\"NodeHeading\" class=\"h2\" updated=\"20210610154910\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "<div data-subtype=\"h1\" data-node-id=\"20210610154900-5r20m2j\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20210610154902\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-subtype=\"h1\" data-node-id=\"20210610154909-wemyt8x\" data-node-index=\"2\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20210610154910\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"0", "<div data-node-id=\"20210518185646-hjjhl5p\" data-node-index=\"0\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20210518191256-m1ij6pn\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20210518185646-pcdktyp\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "<div data-subtype=\"h1\" data-node-id=\"20210518185646-hjjhl5p\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\" updated=\"20210518191256\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-subtype=\"h1\" data-node-id=\"20210518191256-m1ij6pn\" data-node-index=\"2\" data-type=\"NodeHeading\" class=\"h1\" fold=\"0\" updated=\"20210518191257\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-subtype=\"h1\" data-node-id=\"20210518185646-pcdktyp\" data-node-index=\"3\" data-type=\"NodeHeading\" class=\"h1\"><div contenteditable=\"true\" spellcheck=\"false\"></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
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
	{"0", "<div data-subtype=\"u\" data-node-id=\"20210526172440-zgfgiiv\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210526172454\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210526172454-nryi4wk\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210526172454-me3ybh0\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210526172455\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210526172455-8oc3ap6\" data-type=\"NodeListItem\" class=\"li protyle-wysiwyg--select\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210526172455-8lbabpq\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210526172456\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "<div data-subtype=\"t\" data-node-id=\"20210526172440-zgfgiiv\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210526172454\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210526172454-nryi4wk\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\" draggable=\"true\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210526172454-me3ybh0\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210526172455\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210526172455-8oc3ap6\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\" draggable=\"true\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210526172455-8lbabpq\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210526172456\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
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
	{"0", "<div data-subtype=\"t\" data-node-id=\"20210621152530-fu33pvt\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210621152825\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210621152826-p1o0uwy\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210621152826-n4o2vir\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152827\"><div contenteditable=\"true\" spellcheck=\"false\">1</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210621152827-wvqo7lo\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210621152827-ttmqnj4\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152827\"><div contenteditable=\"true\" spellcheck=\"false\">2</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\"></div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210621152828-j4hgdt5\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210621152828-pgwkd1s\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152828\"><div contenteditable=\"true\" spellcheck=\"false\">3</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "<div data-subtype=\"o\" data-node-id=\"20210621152530-fu33pvt\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210621152825\"><div data-marker=\"1.\" data-subtype=\"o\" data-node-id=\"20210621152826-p1o0uwy\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">1.</div><div data-node-id=\"20210621152826-n4o2vir\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152827\"><div contenteditable=\"true\" spellcheck=\"false\">1</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-marker=\"2.\" data-subtype=\"o\" data-node-id=\"20210621152827-wvqo7lo\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">2.</div><div data-node-id=\"20210621152827-ttmqnj4\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152827\"><div contenteditable=\"true\" spellcheck=\"false\">2</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-marker=\"3.\" data-subtype=\"o\" data-node-id=\"20210621152828-j4hgdt5\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--order\" contenteditable=\"false\" draggable=\"true\">3.</div><div data-node-id=\"20210621152828-pgwkd1s\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210621152828\"><div contenteditable=\"true\" spellcheck=\"false\">3</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
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
	{"1", "<div data-subtype=\"t\" data-node-id=\"20210610155819-dvxb0ws\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210610155823\"><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210610155824-ibbpfdr\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210610155824-6pcn10x\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210610155825\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-marker=\"*\" data-subtype=\"t\" data-node-id=\"20210610155825-trhucha\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action protyle-action--task\"><svg><use xlink:href=\"#iconUncheck\"></use></svg></div><div data-node-id=\"20210610155825-8v3wz67\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210610155825\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "<div data-node-id=\"20210610155824-6pcn10x\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210610155825\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20210610155825-8v3wz67\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210610155825\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
	{"0", "<div data-subtype=\"u\" data-node-id=\"20210509095907-rixmte6\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\" updated=\"20210509103643\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210509103643-es01df0\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210509103643-xr4nn64\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103643\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20210509103701-dnmzej3\" data-type=\"NodeList\" class=\"list\" data-subtype=\"u\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210509103644-zfqz75l\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210509103644-uednsdt\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103651\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "<div data-node-id=\"20210509103643-xr4nn64\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103643\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-subtype=\"u\" data-node-id=\"20210509103701-dnmzej3\" data-node-index=\"2\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210509103644-zfqz75l\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\" draggable=\"true\"><svg><use xlink:href=\"#iconDot\"></use></svg></div><div data-node-id=\"20210509103644-uednsdt\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20210509103651\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
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

var cancelBlockquoteTests = []*parseTest{
	{"0", "<div data-node-id=\"20211107212413-7y4w6vm\" data-node-index=\"1\" data-type=\"NodeBlockquote\" class=\"bq\" updated=\"20211107235603\"><div data-node-id=\"20211107235603-js5hmfn\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20211107235606\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20211107235606-7ceugn5\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20211107235607\"><div contenteditable=\"true\" spellcheck=\"false\" data-tip=\"键入文字或 '/' 选择\">bar</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>", "<div data-node-id=\"20211107235603-js5hmfn\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20211107235606\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div><div data-node-id=\"20211107235606-7ceugn5\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\" updated=\"20211107235607\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\" contenteditable=\"false\">\u200b</div></div>"},
}

func TestCancelBlockquote(t *testing.T) {
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
	for _, test := range cancelBlockquoteTests {
		ovHTML := luteEngine.CancelBlockquote(test.from)
		if test.to != ovHTML {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, ovHTML, test.from)
		}
	}
	ast.Testing = false
}

var ul2olTests = []*parseTest{}

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

var ol2ulTests = []*parseTest{}

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
