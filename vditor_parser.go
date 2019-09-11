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
	case atom.Br:
		node.typ = NodeInlineHTML
		node.tokens = toItems("<br />")
	case atom.Em:
		node.typ = NodeEmphasis
		node.strongEmDelMarker = n.PrevSibling.FirstChild.Data[0]
		node.strongEmDelMarkenLen = 1
	case atom.Strong:
		node.typ = NodeStrong
		node.strongEmDelMarker = n.PrevSibling.FirstChild.Data[0]
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
		tree.context.tip.AppendChild(tree.context.tip, node)
		tree.context.tip = node
	}
	if !skipChildren {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lute.genASTByVditorDOM(c, tree)
		}
	}
	if 1 <= node.typ {
		tree.context.tip = tree.context.tip.parent
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
