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
	"bytes"
	"strings"

	"github.com/b3log/lute/html"
	"github.com/b3log/lute/html/atom"
)

// genASTByDOM 根据指定的 DOM 节点 n 进行深度优先遍历并逐步生成 Markdown 语法树 tree。
func (lute *Lute) genASTByDOM(n *html.Node, tree *Tree) {
	skipChildren := false
	node := &Node{typ: -1, tokens: strToItems(n.Data)}
	switch n.DataAtom {
	case 0:
		node.typ = NodeText
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		node.typ = NodeHeading
		node.headingLevel = int(node.tokens[1].term() - byte('0'))
		node.AppendChild(&Node{typ:NodeHeadingC8hMarker, tokens: strToItems(strings.Repeat("#", node.headingLevel))})
	case atom.Br:
		node.typ = NodeInlineHTML
		node.tokens = strToItems("<br />")
	case atom.Em:
		node.typ = NodeEmphasis
		node.AppendChild(&Node{typ: NodeEmU8eOpenMarker, tokens: strToItems("_")})
	case atom.Strong:
		node.typ = NodeStrong
		node.AppendChild(&Node{typ: NodeStrongA6kOpenMarker, tokens: strToItems("**")})
	case atom.P:
		node.typ = NodeParagraph
	case atom.Blockquote:
		node.typ = NodeBlockquote
	case atom.Hr:
		node.typ = NodeThematicBreak
	case atom.Span:
		mtype := lute.domAttrValue(n, "data-mtype")
		if "2" == mtype { // 行级元素可以直接用其 text
			node.typ = NodeText
			node.tokens = strToItems(lute.domText(n))
			break
		}

		class := lute.domAttrValue(n, "class")
		skipChildren = !strings.Contains(class, "node")
		if skipChildren {
			ntype := lute.domAttrValue(n, "data-ntype")
			skipChildren = "10" != ntype
		}
	}

	if -1 != node.typ {
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
	}
	if !skipChildren {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lute.genASTByDOM(c, tree)

			switch n.DataAtom {
			case atom.Em:
				node.AppendChild(&Node{typ: NodeEmU8eCloseMarker, tokens: strToItems("_")})
			case atom.Strong:
				node.AppendChild(&Node{typ: NodeStrongA6kCloseMarker, tokens: strToItems("**")})
			}
		}
	}
	if -1 != node.typ {
		tree.context.tip = tree.context.tip.parent
	}
}

func (lute *Lute) domAttrValue(n *html.Node, attrName string) string {
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

func (lute *Lute) domText(n *html.Node) string {
	buf := &bytes.Buffer{}
	lute.domText0(n, buf)
	return buf.String()
}

func (lute *Lute) domText0(n *html.Node, buffer *bytes.Buffer) {
	if nil == n {
		return
	}
	if 0 == n.DataAtom {
		buffer.WriteString(n.Data)
	}
	for child := n.FirstChild; nil != child; child = child.NextSibling {
		lute.domText0(child, buffer)
	}
}
