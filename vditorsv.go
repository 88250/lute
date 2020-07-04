// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

import (
	"bytes"
	"github.com/88250/lute/ast"
	"strings"

	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

// SpinVditorSVDOM 自旋 Vditor Split-View DOM，用于分屏预览模式下的编辑。
func (lute *Lute) SpinVditorSVDOM(markdown string) (ovHTML string) {
	lute.VditorSV = true
	lute.VditorWYSIWYG = true

	tree := parse.Parse("", []byte(markdown), lute.Options)
	lute.adjustVditorSVTree(tree)

	renderer := render.NewVditorSVRenderer(tree)
	output := renderer.Render()
	if renderer.Option.Footnotes && 0 < len(renderer.Tree.Context.FootnotesDefs) {
		output = renderer.RenderFootnotesDefs(renderer.Tree.Context)
	}
	// 替换插入符
	ovHTML = strings.ReplaceAll(string(output), util.Caret, "<wbr>")
	return
}

// HTML2VditorSVDOM 将 HTML 转换为 Vditor Split-View DOM，用于分屏预览模式下粘贴。
func (lute *Lute) HTML2VditorSVDOM(sHTML string) (vHTML string) {
	lute.VditorSV = true
	lute.VditorWYSIWYG = true

	markdown, err := lute.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), lute.Options)
	renderer := render.NewVditorSVRenderer(tree)
	for nodeType, rendererFunc := range lute.HTML2VditorSVDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	if renderer.Option.Footnotes && 0 < len(renderer.Tree.Context.FootnotesDefs) {
		output = renderer.RenderFootnotesDefs(renderer.Tree.Context)
	}
	vHTML = string(output)
	return
}

func (lute *Lute) adjustVditorSVTree(tree *parse.Tree) {
	lute.continueListItem(tree.Root)
}

// continueListItem 用于延续列表项，当指定的 node 结构如下时
//
//   * foo
//   ‸
//
//  将生成延续列表项
//
//   * foo
//   * ‸
func (lute *Lute) continueListItem(node *ast.Node) {
	caretNode := lute.findCaretNode(node)
	if nil == caretNode {
		return
	}

	if !bytes.Equal(caretNode.Tokens, util.CaretTokens) {
		return
	}

	if nil == caretNode.Previous || ast.NodeSoftBreak != caretNode.Previous.Type || (ast.NodeParagraph != caretNode.Parent.Type && ast.NodeBlockquote != caretNode.Parent.Type) {
		return
	}

	if nil == caretNode.Parent.Parent || ast.NodeListItem != caretNode.Parent.Parent.Type {
		return
	}

	prevLi := caretNode.Parent.Parent
	li := &ast.Node{Type: ast.NodeListItem, Tokens: prevLi.Tokens, ListData: &ast.ListData{
		Typ:          prevLi.ListData.Typ,
		Tight:        prevLi.ListData.Tight,
		BulletChar:   prevLi.ListData.BulletChar,
		Start:        prevLi.ListData.Start,
		Delimiter:    prevLi.ListData.Delimiter,
		Padding:      prevLi.ListData.Padding,
		MarkerOffset: prevLi.ListData.MarkerOffset,
		Checked:      prevLi.ListData.Checked,
		Marker:       prevLi.ListData.Marker,
		Num:          prevLi.ListData.Num + 1,
	}}
	li.AppendChild(&ast.Node{Type: ast.NodeParagraph})
	li.AppendChild(caretNode)
	prevLi.InsertAfter(li)
}

func (lute *Lute) findCaretNode(node *ast.Node) (ret *ast.Node) {
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			if ast.NodeText == n.Type {
				if bytes.Contains(n.Tokens, util.CaretTokens) {
					ret = n
					return ast.WalkStop
				}
			}
		}
		return ast.WalkContinue
	})
	return
}
