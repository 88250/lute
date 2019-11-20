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
	"github.com/b3log/lute/html"
	"github.com/b3log/lute/html/atom"
	"strings"
)

// 插入符 \u2038
const caret = "‸"

// RenderVditorDOM 用于渲染 Vditor DOM。
func (lute *Lute) RenderVditorDOM(htmlStr string) (html string, err error) {
	lute.VditorWYSIWYG = true

	var md string
	md, err = lute.VditorDOM2Md(htmlStr)
	if nil != err {
		return
	}

	var tree *Tree
	tree, err = lute.parse("", []byte(md))
	if nil != err {
		return
	}

	renderer := lute.newVditorRenderer(tree)
	var output []byte
	output, err = renderer.Render()
	// 替换插入符
	html = strings.ReplaceAll(string(output), caret, "<wbr>")
	return
}

// VditorDOM2Md 将 Vditor DOM 转换为 Markdown 文本。
func (lute *Lute) VditorDOM2Md(htmlStr string) (md string, err error) {
	// 替换插入符
	htmlStr = strings.ReplaceAll(htmlStr, "<wbr>", caret)

	// 将字符串解析为 DOM 树

	reader := strings.NewReader(htmlStr)
	htmlRoot := &html.Node{Type: html.ElementNode}
	htmlNodes, err := html.ParseFragment(reader, htmlRoot)
	if nil != err {
		return
	}

	// 将 HTML 树转换为 Markdown AST

	tree := &Tree{Name: "", Root: &Node{typ: NodeDocument}, context: &Context{option: lute.options}}
	tree.context.tip = tree.Root
	for _, htmlNode := range htmlNodes {
		lute.genASTByVditorDOM(htmlNode, tree)
	}

	// 将 AST 进行 Markdown 格式化渲染

	var formatted []byte
	renderer := lute.newFormatRenderer(tree)
	formatted, err = renderer.Render()
	if nil != err {
		return
	}
	md = bytesToStr(formatted)
	return
}

// genASTByVditorDOM 根据指定的 Vditor DOM 节点 n 进行深度优先遍历并逐步生成 Markdown 语法树 tree。
func (lute *Lute) genASTByVditorDOM(n *html.Node, tree *Tree) {
	node := &Node{typ: NodeText, tokens: strToBytes(n.Data)}
	switch n.DataAtom {
	case 0:
		if nil != n.Parent {
			switch n.Parent.DataAtom {
			case atom.Code:
				if nil == n.Parent.Parent {
					node.typ = NodeCodeSpanContent
				} else {
					node.typ = NodeCodeBlockCode
					class := lute.domAttrValue(n.Parent, "class")
					if strings.Contains(class, "language-") {
						language := class[len("language-"):]
						tree.context.tip.lastChild.codeBlockInfo = strToBytes(language)
					}
				}
			case atom.A:
				node.typ = NodeLinkText
			}
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Wbr:
		node.typ = NodeText
		node.tokens = strToBytes("<wbr>")
		lastc := tree.context.tip.lastDeepestChild()
		if NodeParagraph == lastc.typ {
			lastc.AppendChild(node)
		} else {
			lastc.AppendTokens(node.tokens)
		}
		return
	case atom.Div:
		class := lute.domAttrValue(n, "class")
		if strings.Contains(class, "vditor-panel") {
			return
		}
		fallthrough
	case atom.P:
		node.typ = NodeParagraph
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		node.typ = NodeHeading
		node.headingLevel = int(node.tokens[1] - byte('0'))
		node.AppendChild(&Node{typ: NodeHeadingC8hMarker, tokens: strToBytes(strings.Repeat("#", node.headingLevel))})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Hr:
		node.typ = NodeThematicBreak
		tree.context.tip.AppendChild(node)
	case atom.Blockquote:
		node.typ = NodeBlockquote
		node.AppendChild(&Node{typ: NodeBlockquoteMarker, tokens: strToBytes(">")})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Ol, atom.Ul:
		node.typ = NodeList
		node.listData = &listData{}
		if atom.Ol == n.DataAtom {
			node.listData.typ = 1
		}
		tight := lute.domAttrValue(n, "data-tight")
		if "true" == tight || "" == tight {
			node.tight = true
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Li:
		node.typ = NodeListItem
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			if atom.Ol == n.Parent.DataAtom {
				start := lute.domAttrValue(n.Parent, "start")
				if "" == start {
					marker = "1."
				} else {
					marker = start + "."
				}
			} else {
				marker = "*"
			}
		}
		node.listData = &listData{marker: strToBytes(marker)}
		if lute.firstChildIsText(n) {
			tree.context.tip.AppendChild(node)
			tree.context.tip = node
			node = &Node{typ: NodeParagraph}
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Pre:
		node.typ = NodeCodeBlock
		node.isFencedCodeBlock = true
		node.AppendChild(&Node{typ: NodeCodeBlockFenceOpenMarker, tokens: strToBytes("```"), codeBlockFenceLen: 3})
		node.AppendChild(&Node{typ: NodeCodeBlockFenceInfoMarker})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Em, atom.I:
		node.typ = NodeEmphasis
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "*"
		}
		if "_" == marker {
			node.AppendChild(&Node{typ: NodeEmU8eOpenMarker, tokens: strToBytes(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeEmA6kOpenMarker, tokens: strToBytes(marker)})
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Strong, atom.B:
		node.typ = NodeStrong
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "**"
		}
		if "__" == marker {
			node.AppendChild(&Node{typ: NodeStrongU8eOpenMarker, tokens: strToBytes(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeStrongA6kOpenMarker, tokens: strToBytes(marker)})
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Code:
		if nil == n.Parent || atom.Pre != n.Parent.DataAtom {
			node.typ = NodeCodeSpan
			node.AppendChild(&Node{typ: NodeCodeSpanOpenMarker, tokens: strToBytes("`")})
			tree.context.tip.AppendChild(node)
			tree.context.tip = node
			defer tree.context.parentTip(n)
		}
	case atom.Br:
		node.typ = NodeInlineHTML
		node.tokens = strToBytes("<br />")
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.A:
		node.typ = NodeLink
		node.AppendChild(&Node{typ: NodeOpenBracket})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Img:
		imgClass := lute.domAttrValue(n, "class")
		imgAlt := lute.domAttrValue(n, "alt")
		if "emoji" == imgClass {
			node.typ = NodeEmoji
			emojiImg := &Node{typ: NodeEmojiImg, tokens: tree.emojiImgTokens(imgAlt, lute.domAttrValue(n, "src"))}
			emojiImg.AppendChild(&Node{typ: NodeEmojiAlias, tokens: strToBytes(":" + imgAlt + ":")})
			node.AppendChild(emojiImg)
		} else {
			node.typ = NodeImage
			node.AppendChild(&Node{typ: NodeBang})
			node.AppendChild(&Node{typ: NodeOpenBracket})
			if "" != imgAlt {
				node.AppendChild(&Node{typ: NodeLinkText, tokens: strToBytes(imgAlt)})
			}
			node.AppendChild(&Node{typ: NodeCloseBracket})
			node.AppendChild(&Node{typ: NodeOpenParen})
			node.AppendChild(&Node{typ: NodeLinkDest, tokens: strToBytes(lute.domAttrValue(n, "src"))})
			node.AppendChild(&Node{typ: NodeCloseParen})
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Input:
		node.typ = NodeTaskListItemMarker
		if lute.hasAttr(n, "checked") {
			node.taskListItemChecked = true
			node.tokens = strToBytes("[X]")
		} else {
			node.tokens = strToBytes("[ ]")
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
		if nil != node.parent.parent {
			node.parent.parent.listData.typ = 3
		}
	case atom.Del, atom.S:
		node.typ = NodeStrikethrough
		marker := lute.domAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&Node{typ: NodeStrikethrough1OpenMarker, tokens: strToBytes(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeStrikethrough2OpenMarker, tokens: strToBytes(marker)})
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Table:
		node.typ = NodeTable
		var tableAligns []int
		for c := n.FirstChild.FirstChild.FirstChild; nil != c; c = c.NextSibling {
			tableAligns = append(tableAligns, 0)
		}
		node.tableAligns = tableAligns
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Thead:
		node.typ = NodeTableHead
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Tr:
		node.typ = NodeTableRow
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Th, atom.Td:
		node.typ = NodeTableCell
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Span:
		marker := lute.domAttrValue(n, "data-marker")
		if "" != marker {
			switch marker {
			case "*", "_":
				node.typ = NodeEmphasis
				if "_" == marker {
					node.AppendChild(&Node{typ: NodeEmU8eOpenMarker, tokens: strToBytes(marker)})
				} else {
					node.AppendChild(&Node{typ: NodeEmA6kOpenMarker, tokens: strToBytes(marker)})
				}
			case "**", "__":
				if "__" == marker {
					node.AppendChild(&Node{typ: NodeStrongU8eOpenMarker, tokens: strToBytes(marker)})
				} else {
					node.AppendChild(&Node{typ: NodeStrongA6kOpenMarker, tokens: strToBytes(marker)})
				}

			}
		} else {
			node.tokens = strToBytes("")
		}

		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTByVditorDOM(c, tree)
	}

	switch n.DataAtom {
	case atom.Em, atom.I:
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "*"
		}
		if "_" == marker {
			node.AppendChild(&Node{typ: NodeEmU8eCloseMarker, tokens: strToBytes(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeEmA6kCloseMarker, tokens: strToBytes(marker)})
		}
	case atom.Strong, atom.B:
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "**"
		}
		if "__" == marker {
			node.AppendChild(&Node{typ: NodeStrongU8eCloseMarker, tokens: strToBytes(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeStrongA6kCloseMarker, tokens: strToBytes(marker)})
		}
	case atom.Pre:
		node.AppendChild(&Node{typ: NodeCodeBlockFenceCloseMarker, tokens: strToBytes("```"), codeBlockFenceLen: 3})
	case atom.Code:
		if nil == n.Parent || atom.Pre != n.Parent.DataAtom {
			node.AppendChild(&Node{typ: NodeCodeSpanCloseMarker, tokens: strToBytes("`")})
		}
	case atom.A:
		node.AppendChild(&Node{typ: NodeCloseBracket})
		node.AppendChild(&Node{typ: NodeOpenParen})
		node.AppendChild(&Node{typ: NodeLinkDest, tokens: strToBytes(lute.domAttrValue(n, "href"))})
		linkTitle := lute.domAttrValue(n, "title")
		if "" != linkTitle {
			node.AppendChild(&Node{typ: NodeLinkSpace})
			node.AppendChild(&Node{typ: NodeLinkTitle, tokens: strToBytes(linkTitle)})
		}
		node.AppendChild(&Node{typ: NodeCloseParen})
	case atom.Del, atom.S:
		marker := lute.domAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&Node{typ: NodeStrikethrough1CloseMarker, tokens: strToBytes(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeStrikethrough2CloseMarker, tokens: strToBytes(marker)})
		}
	case atom.Span:
		marker := lute.domAttrValue(n, "data-marker")
		if "" != marker {
			switch marker {
			case "*", "_":
				node.typ = NodeEmphasis
				if "_" == marker {
					node.AppendChild(&Node{typ: NodeEmU8eCloseMarker, tokens: strToBytes(marker)})
				} else {
					node.AppendChild(&Node{typ: NodeEmA6kCloseMarker, tokens: strToBytes(marker)})
				}
			case "**", "__":
				if "__" == marker {
					node.AppendChild(&Node{typ: NodeStrongU8eCloseMarker, tokens: strToBytes(marker)})
				} else {
					node.AppendChild(&Node{typ: NodeStrongA6kCloseMarker, tokens: strToBytes(marker)})
				}
			}
		}
	}
}

func (context *Context) parentTip(n *html.Node) {
	for tip := context.tip.parent; nil != tip; tip = tip.parent {
		if NodeParagraph == tip.typ {
			if nil == n.NextSibling {
				continue
			}
			nextType := n.NextSibling.DataAtom
			if atom.Ul == nextType ||
				atom.Ol == nextType {
				continue
			}
		}
		context.tip = tip
		break
	}
}

// firstChildIsText 用于判断 n 的第一个子节点是否是文本节点。
func (lute *Lute) firstChildIsText(n *html.Node) bool {
	for c := n.FirstChild; nil != c; c = c.NextSibling {
		if caret == c.Data {
			continue // 不考虑插入符
		}
		return 0 == c.DataAtom || atom.Em == c.DataAtom
	}
	return false
}

func (lute *Lute) hasAttr(n *html.Node, attrName string) bool {
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return true
		}
	}
	return false
}

func (lute *Lute) domAttrValue(n *html.Node, attrName string) string {
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
