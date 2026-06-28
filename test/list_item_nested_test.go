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

// 开关开启：空列表项下创建子列表前补一个空段落（HTML 渲染会省略空段落，效果与默认一致）
var disableNestedListTests = []parseTest{
	{"空列表项-子列表", "-\n  - bar\n",
		"<ul>\n<li>\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>\n"},
	{"一行-连写", "- - bar\n",
		"<ul>\n<li>\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>\n"},
	{"一行*连写", "* * bar\n",
		"<ul>\n<li>\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>\n"},
	{"一行有序连写", "1. 1. bar\n",
		"<ol>\n<li>\n<ol>\n<li>bar</li>\n</ol>\n</li>\n</ol>\n"},
	{"有内容列表项-子列表-不受影响", "- foo\n  - bar\n",
		"<ul>\n<li>foo\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>\n"},
	{"同级列表-不受影响", "- foo\n- bar\n",
		"<ul>\n<li>foo</li>\n<li>bar</li>\n</ul>\n"},
}

func TestEnsureListItemParagraph(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetEnsureListItemParagraph(true)
	luteEngine.SetHeadingID(false)
	luteEngine.SetKramdownIAL(false)

	for _, test := range disableNestedListTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal\n\t%q",
				test.name, test.to, html, test.from)
		}
	}
}

// 开关关闭（默认）：空列表项下直接挂子列表
var defaultNestedListTests = []parseTest{
	{"空列表项-子列表", "-\n  - bar\n",
		"<ul>\n<li>\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>\n"},
}

func TestDefaultNestedList(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetHeadingID(false)
	luteEngine.SetKramdownIAL(false)

	for _, test := range defaultNestedListTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal\n\t%q",
				test.name, test.to, html, test.from)
		}
	}
}

// containsText 判断节点子树中是否存在带文本的 NodeText 节点
func containsText(n *ast.Node) bool {
	if n.Type == ast.NodeText && 0 < len(n.Tokens) {
		return true
	}
	for c := n.FirstChild; nil != c; c = c.Next {
		if containsText(c) {
			return true
		}
	}
	return false
}

// emptyParagraphCount 统计列表项下不含任何文本的空段落数量
func emptyParagraphCount(listItem *ast.Node) (count int) {
	for c := listItem.FirstChild; nil != c; c = c.Next {
		if ast.NodeParagraph == c.Type && !containsText(c) {
			count++
		}
	}
	return
}

var issue17890Tests = []string{
	"- 1\n\n  -\n- 3\n- 4\n", // 顶层有内容的列表项下挂空嵌套项，其后跟同级非空项
	"-\n  -\n- bar\n",         // 顶层空列表项下挂空嵌套项，其后跟同级非空项
	"-\n  -\n  - baz\n",       // 空列表项下连续两个嵌套项，前一个为空
}

// 空列表项下创建子列表前补空段落后，同一个空列表项里不应出现重复的空段落
// https://github.com/siyuan-note/siyuan/issues/17890
func TestIssue17890EmptyListItemNoDuplicateParagraph(t *testing.T) {
	for _, markdown := range issue17890Tests {
		// 与前端 Protyle 编辑器一致的解析选项：ProtyleWYSIWYG + KramdownBlockIAL 同时开启时，
		// ListStart 与 listFinalize 两条路径都会为空列表项补段落，需保证不重复
		luteEngine := lute.New()
		luteEngine.SetProtyleWYSIWYG(true)
		luteEngine.SetKramdownIAL(true)
		luteEngine.SetSpin(true)
		luteEngine.SetEnsureListItemParagraph(true)
		_, tree := luteEngine.Md2BlockDOMTree(markdown, false)

		var offenders []string
		ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
			if !entering || ast.NodeListItem != n.Type {
				return ast.WalkContinue
			}
			if 1 < emptyParagraphCount(n) {
				offenders = append(offenders, n.ID)
			}
			return ast.WalkContinue
		})
		if 0 < len(offenders) {
			t.Fatalf("markdown %q 中列表项 %v 出现重复空段落\nexpected 每个空列表项至多一个空段落",
				markdown, offenders)
		}
	}
}
