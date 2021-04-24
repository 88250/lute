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

var md2BlockDOMTests = []parseTest{

	{"17", "[foo]\n{: id=\"20210408204847-qyy54ha\"}\n[foo]: bar\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54ha\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">[foo]</div><div class=\"protyle-attr\"></div></div><div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">[foo]: bar</div><div class=\"protyle-attr\"></div></div>"},
	{"16", "<iframe src=\"foo\" scrolling=\"no\" border=\"0\" frameborder=\"no\" framespacing=\"0\" allowfullscreen=\"true\"> </iframe>\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeIFrame\" class=\"iframe\"><span class=\"protyle-action\"><svg class=\"svg\"><use xlink:href=\"#iconMore\"></use></svg><span class=\"protyle-action__drag\"></span></span><iframe src=\"foo\" scrolling=\"no\" border=\"0\" frameborder=\"no\" framespacing=\"0\" allowfullscreen=\"true\"> </iframe><div class=\"protyle-attr\"></div></div>"},
	{"15", "foo((20210121085548-9vnyjk4))**baz**\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo<span data-type=\"block-ref\" data-id=\"20210121085548-9vnyjk4\" data-anchor=\"\" contenteditable=\"false\"></span><strong>baz</strong></div><div class=\"protyle-attr\"></div></div>"},
	{"14", "![alt](bar)\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\"><span contenteditable=\"false\" data-type=\"img\" class=\"img\"><span class=\"protyle-action\"><svg class=\"svg\"><use xlink:href=\"#iconMore\"></use></svg></span><img src=\"/siyuan/0/测试笔记/bar\" data-src=\"bar\" alt=\"alt\" /><span class=\"protyle-action__drag\"></span><span class=\"protyle-action__title\"></span></span></div><div class=\"protyle-attr\"></div></div>"},
	{"13", "|foo|\n|-|\n|bar|\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeTable\" class=\"table\"><div class=\"protyle-action\"><svg class=\"svg\"><use xlink:href=\"#iconMore\"></use></svg></div><div contenteditable=\"true\" spellcheck=\"false\"><table><thead><tr><th>foo</th></tr></thead><tbody><tr><td>bar</td></tr></tbody></table></div><div class=\"protyle-attr\"></div></div>"},
	{"12", "<p>foo</p>\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">&lt;p&gt;foo&lt;/p&gt;</div><div class=\"protyle-attr\"></div></div>"},
	{"11", "foo((20210121085548-9vnyjk4 \"bar\"))**baz**\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo<span data-type=\"block-ref\" data-id=\"20210121085548-9vnyjk4\" data-anchor=\"bar\" contenteditable=\"false\"></span><strong>baz</strong></div><div class=\"protyle-attr\"></div></div>"},
	{"10", "$$\nfoo\n$$\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeMathBlock\" class=\"render-node\" data-content=\"foo\" data-subtype=\"math\"><div spin=\"1\"></div><div class=\"protyle-attr\"></div></div>"},
	{"9", "```abc\nfoo\n```\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"render-node\" data-content=\"foo\" data-subtype=\"abc\"><div spin=\"1\"></div><div class=\"protyle-attr\"></div></div>"},
	{"8", "```\nfoo\n```\n{: id=\"20210408204847-qyy54hz\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeCodeBlock\" class=\"code-block\"><div class=\"protyle-action\"><div class=\"protyle-action__language\"></div><div class=\"protyle-action__copy\"></div></div><div contenteditable=\"true\" spellcheck=\"false\">foo\n</div><div class=\"protyle-attr\"></div></div>"},
	{"7", "foo\n{: id=\"20210408204847-qyy54hz\"}\n---\n{: id=\"20210408204848-qyy54ha\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div data-node-id=\"20210408204848-qyy54ha\" data-node-index=\"2\" data-type=\"NodeThematicBreak\" class=\"hr\"><div></div></div>"},
	{"6", "<<<<<<< HEAD\n{: id=\"20210121085548-9vnyjk4\"}\n=======\n{: id=\"20210120223935-11oegu7\"}\n>>>>>>> parent of fe4124a (Revert \"commit data\")", "<div data-node-id=\"20060102150405-1a2b3c4\" data-type=\"NodeGitConflictContent\" class=\"git-conflict\"><div contenteditable=\"true\" spellcheck=\"false\">{: id=&quot;20210121085548-9vnyjk4&quot;}\n=======\n{: id=&quot;20210120223935-11oegu7&quot;}</div><div class=\"protyle-attr\"></div></div>"},
	{"5", "* {: id=\"20210415082227-m67yq1v\"}foo\n  {: id=\"20210415082227-z9mgkh5\"}\n* {: id=\"20210415091213-c387rm0\"}bar\n  {: id=\"20210415091222-knbamrt\"}\n{: id=\"20210414223654-vfqydjh\"}\n\n\n{: id=\"20210413000727-8ua3vhv\" type=\"doc\"}\n", "<div data-subtype=\"u\" data-node-id=\"20210414223654-vfqydjh\" data-node-index=\"1\" data-type=\"NodeList\" class=\"list\"><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210415082227-m67yq1v\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"></div><div data-node-id=\"20210415082227-z9mgkh5\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div data-marker=\"*\" data-subtype=\"u\" data-node-id=\"20210415091213-c387rm0\" data-type=\"NodeListItem\" class=\"li\"><div class=\"protyle-action\"></div><div data-node-id=\"20210415091222-knbamrt\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>"},
	{"4", "{{{col\nfoo\n{: id=\"20210413003741-qbi5h4h\"}\n\nbar\n{: id=\"20210413001027-1mo28cc\"}\n}}}\n{: id=\"20210413005325-jljhnw5\"}\n\n\n{: id=\"20210413000727-8ua3vhv\" type=\"doc\"}\n", "<div data-node-id=\"20210413005325-jljhnw5\" data-node-index=\"1\" data-type=\"NodeSuperBlock\" class=\"sb\" data-sb-layout=\"col\"><div data-node-id=\"20210413003741-qbi5h4h\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div data-node-id=\"20210413001027-1mo28cc\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>"},
	{"3", "# foo\n{: id=\"20210408153138-td774lp\"}", "<div data-subtype=\"h1\" data-node-id=\"20210408153138-td774lp\" data-node-index=\"1\" data-type=\"NodeHeading\" class=\"h1\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div>"},
	{"2", "> foo\n> {: id=\"20210408153138-td774lp\"}\n{: id=\"20210408153137-zds0o4x\"}", "<div data-node-id=\"20210408153137-zds0o4x\" data-node-index=\"1\" data-type=\"NodeBlockquote\" class=\"bq\"><div data-node-id=\"20210408153138-td774lp\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"></div></div><div class=\"protyle-attr\"></div></div>"},
	{"1", "foo\n{: id=\"20210408204847-qyy54hz\" bookmark=\"bm\"}\nbar\n{: id=\"20210408204848-qyy54ha\"}", "<div data-node-id=\"20210408204847-qyy54hz\" data-node-index=\"1\" data-type=\"NodeParagraph\" class=\"p\" bookmark=\"bm\"><div contenteditable=\"true\" spellcheck=\"false\">foo</div><div class=\"protyle-attr\"><div class=\"protyle-attr--bookmark\">bm</div></div></div><div data-node-id=\"20210408204848-qyy54ha\" data-node-index=\"2\" data-type=\"NodeParagraph\" class=\"p\"><div contenteditable=\"true\" spellcheck=\"false\">bar</div><div class=\"protyle-attr\"></div></div>"},
	{"0", "", ""},
}

func TestMd2BlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetToC(true)
	luteEngine.ParseOptions.BlockRef = true
	luteEngine.ParseOptions.KramdownBlockIAL = true
	luteEngine.RenderOptions.KramdownBlockIAL = true
	luteEngine.ParseOptions.Tag = true
	luteEngine.ParseOptions.SuperBlock = true
	luteEngine.SetLinkBase("/siyuan/0/测试笔记/")
	luteEngine.ParseOptions.GitConflict = true
	luteEngine.ParseOptions.LinkRef = false

	ast.Testing = true
	for _, test := range md2BlockDOMTests {
		result := luteEngine.Md2BlockDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
	ast.Testing = false
}
