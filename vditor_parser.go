// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

import (
	"strings"

	"github.com/b3log/lute/html"
	"github.com/b3log/lute/html/atom"
)

// Vditor DOM Parser

// parseVditorDOM 解析 Vditor DOM 生成 Markdown 文本。
func (lute *Lute) parseVditorDOM(htmlStr string) (tree *Tree, err error) {
	defer recoverPanic(&err)

	// 将字符串解析为 HTML 树

	reader := strings.NewReader(htmlStr)
	htmlRoot := &html.Node{Type: html.ElementNode}
	htmlNodes, err := html.ParseFragment(reader, htmlRoot)
	if nil != err {
		return
	}

	// 将 HTML 树转换为 Markdown 语法树

	tree = &Tree{Name: "", Root: &Node{typ: NodeDocument}, context: &Context{option: lute.options}}
	for _, htmlNode := range htmlNodes {
		tree.context.tip = tree.Root
		lute.genASTByVditorDOM(htmlNode, tree)
	}

	return
}

// genASTByVditorDOM 根据指定的 Vditor DOM 节点 n 进行深度优先遍历并逐步生成 Markdown 语法树 tree。
func (lute *Lute) genASTByVditorDOM(n *html.Node, tree *Tree) {
	skipChildren := false
	node := &Node{typ: -1, tokens: toItems(n.Data)}
	switch n.DataAtom {
	case 0:
		node.typ = NodeText
	case atom.Em:
		node.typ = NodeEmphasis
		node.strongEmDelMarker = itemAsterisk
		node.strongEmDelMarkenLen = 1
	case atom.Strong:
		node.typ = NodeStrong
		node.strongEmDelMarker = itemAsterisk
		node.strongEmDelMarkenLen = 2
	case atom.P:
		node.typ = NodeParagraph
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		node.typ = NodeHeading
		node.headingLevel = int(node.tokens[1] - byte('0'))
	case atom.Span:
		if nil != n.Attr {
			class := lute.domAttrValue(n, "class")
			skipChildren = !strings.Contains(class, "node")
		}
	}

	if 1 <= node.typ {
		tree.context.vditorAddChild(node)
	}
	if !skipChildren {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lute.genASTByVditorDOM(c, tree)
		}
	}
}

func (lute *Lute) domAttrValue(n *html.Node, attrName string) (string) {
	if 1 > len(n.Attr) {
		return ""
	}
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

// vditorAddChild 将 child 作为子节点添加到 context.tip 上。如果 tip 节点不能接受子节点（非块级容器不能添加子节点），则最终化该 tip
// 节点并向父节点方向尝试，直到找到一个能接受 child 的节点为止。
func (context *Context) vditorAddChild(child *Node) {
	for !context.tip.CanContain(child.typ) {
		context.vditorFinalize(context.tip) // 注意调用 vditorFinalize 会向父节点方向进行迭代
	}

	context.tip.AppendChild(context.tip, child)
	context.tip = child
}

// vditorFinalize 执行 block 的最终化处理。调用该方法会将 context.tip 置为 block 的父节点。
func (context *Context) vditorFinalize(block *Node) {
	var parent = block.parent
	block.close = true
	block.Finalize(context)
	context.tip = parent
}
