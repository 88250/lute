// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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

	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
)

// HTML2Markdown 将 HTML 转换为 Markdown。
func (lute *Lute) HTML2Markdown(htmlStr string) (markdown string, err error) {
	// 将字符串解析为 DOM 树

	reader := strings.NewReader(htmlStr)
	htmlRoot := &html.Node{Type: html.ElementNode}
	htmlNodes, err := html.ParseFragment(reader, htmlRoot)
	if nil != err {
		return
	}

	// 将 HTML 树转换为 Markdown AST

	tree := &Tree{Name: "", Root: &Node{Typ: NodeDocument}, context: &Context{option: lute.options}}
	tree.context.tip = tree.Root
	for _, htmlNode := range htmlNodes {
		lute.genASTByDOM(htmlNode, tree)
	}

	// 调整树结构
	// TODO: 列表项依赖入参带有 p 节点，需要在此调整为自动插入 p 节点

	Walk(tree.Root, func(n *Node, entering bool) (status WalkStatus, e error) {
		if entering {
			if NodeList == n.Typ {
				// ul.ul => ul.li.ul
				if nil != n.Parent && NodeList == n.Parent.Typ {
					previousLi := n.Previous
					if nil != previousLi {
						n.Unlink()
						previousLi.AppendChild(n)
					}
				}
			}
		}
		return WalkContinue, nil
	})

	// 将 AST 进行 Markdown 格式化渲染

	var formatted []byte
	renderer := lute.newFormatRenderer(tree)
	for nodeType, render := range lute.HTML2MdRendererFuncs {
		renderer.(*FormatRenderer).extRendererFuncs[nodeType] = render
	}
	formatted, err = renderer.Render()
	if nil != err {
		return
	}
	markdown = bytesToStr(formatted)
	return
}

// genASTByDOM 根据指定的 DOM 节点 n 进行深度优先遍历并逐步生成 Markdown 语法树 tree。
func (lute *Lute) genASTByDOM(n *html.Node, tree *Tree) {
	if html.CommentNode == n.Type || atom.Meta == n.DataAtom {
		return
	}

	dataRender := lute.domAttrValue(n, "data-render")
	if "false" == dataRender {
		return
	}

	node := &Node{Typ: NodeText, Tokens: strToBytes(n.Data)}
	switch n.DataAtom {
	case 0:
		if nil != n.Parent && atom.A == n.Parent.DataAtom {
			node.Typ = NodeLinkText
		}
		node.Tokens = bytes.ReplaceAll(node.Tokens, []byte{194, 160}, []byte{' '}) // 将 &nbsp; 转换为空格
		tree.context.tip.AppendChild(node)
	case atom.P, atom.Div:
		node.Typ = NodeParagraph
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		node.Typ = NodeHeading
		node.headingLevel = int(node.Tokens[1] - byte('0'))
		node.AppendChild(&Node{Typ: NodeHeadingC8hMarker, Tokens: strToBytes(strings.Repeat("#", node.headingLevel))})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Hr:
		node.Typ = NodeThematicBreak
		tree.context.tip.AppendChild(node)
	case atom.Blockquote:
		node.Typ = NodeBlockquote
		node.AppendChild(&Node{Typ: NodeBlockquoteMarker, Tokens: strToBytes(">")})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Ol, atom.Ul:
		node.Typ = NodeList
		node.listData = &listData{}
		if atom.Ol == n.DataAtom {
			node.listData.typ = 1
		}
		node.tight = true
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Li:
		node.Typ = NodeListItem
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			if nil != n.Parent && atom.Ol == n.Parent.DataAtom {
				start := lute.domAttrValue(n.Parent, "start")
				if "" == start {
					marker = "1."
				} else {
					marker = start + "."
				}
			} else {
				marker = "*"
			}
		} else {
			if nil != n.Parent && "1." != marker && atom.Ol == n.Parent.DataAtom && nil != n.Parent.Parent && (atom.Ol == n.Parent.Parent.DataAtom || atom.Ul == n.Parent.Parent.DataAtom) {
				// 子有序列表必须从 1 开始
				marker = "1."
			}
		}
		node.listData = &listData{marker: []byte(marker)}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Pre:
		firstc := n.FirstChild
		if nil != firstc {
			if atom.Code == firstc.DataAtom {
				node.Typ = NodeCodeBlock
				node.isFencedCodeBlock = true
				node.AppendChild(&Node{Typ: NodeCodeBlockFenceOpenMarker, Tokens: strToBytes("```"), codeBlockFenceLen: 3})
				node.AppendChild(&Node{Typ: NodeCodeBlockFenceInfoMarker})
				buf := &bytes.Buffer{}
				class := lute.domAttrValue(firstc, "class")
				if strings.Contains(class, "language-") {
					language := class[len("language-"):]
					node.LastChild.codeBlockInfo = []byte(language)
				}
				firstc = firstc.FirstChild
				for c := firstc; nil != c; c = c.NextSibling {
					buf.WriteString(lute.domText(c))
				}
				content := &Node{Typ: NodeCodeBlockCode, Tokens: buf.Bytes()}
				node.AppendChild(content)
				node.AppendChild(&Node{Typ: NodeCodeBlockFenceCloseMarker, Tokens: strToBytes("```"), codeBlockFenceLen: 3})
				tree.context.tip.AppendChild(node)
			} else {
				node.Typ = NodeHTMLBlock
				node.Tokens = lute.domHTML(n)
				tree.context.tip.AppendChild(node)
			}
		}
		return
	case atom.Em, atom.I:
		node.Typ = NodeEmphasis
		marker := "*"
		node.AppendChild(&Node{Typ: NodeEmA6kOpenMarker, Tokens: strToBytes(marker)})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Strong, atom.B:
		node.Typ = NodeStrong
		marker := "**"
		node.AppendChild(&Node{Typ: NodeStrongA6kOpenMarker, Tokens: strToBytes(marker)})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Code:
		if nil == n.FirstChild {
			return
		}

		code := lute.domHTML(n.FirstChild)
		unescaped := html.UnescapeString(string(code))
		code = []byte(unescaped)
		content := &Node{Typ: NodeCodeSpanContent, Tokens: code}
		node.Typ = NodeCodeSpan
		node.AppendChild(&Node{Typ: NodeCodeSpanOpenMarker, Tokens: []byte("`")})
		node.AppendChild(content)
		node.AppendChild(&Node{Typ: NodeCodeSpanCloseMarker, Tokens: []byte("`")})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
		return
	case atom.Br:
		node.Typ = NodeHardBreak
		node.Tokens = strToBytes("\n")
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.A:
		node.Typ = NodeLink
		node.AppendChild(&Node{Typ: NodeOpenBracket})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Img:
		imgClass := lute.domAttrValue(n, "class")
		imgAlt := lute.domAttrValue(n, "alt")
		if "emoji" == imgClass {
			node.Typ = NodeEmoji
			emojiImg := &Node{Typ: NodeEmojiImg, Tokens: tree.emojiImgTokens(imgAlt, lute.domAttrValue(n, "src"))}
			emojiImg.AppendChild(&Node{Typ: NodeEmojiAlias, Tokens: strToBytes(":" + imgAlt + ":")})
			node.AppendChild(emojiImg)
		} else {
			node.Typ = NodeImage
			node.AppendChild(&Node{Typ: NodeBang})
			node.AppendChild(&Node{Typ: NodeOpenBracket})
			if "" != imgAlt {
				node.AppendChild(&Node{Typ: NodeLinkText, Tokens: strToBytes(imgAlt)})
			}
			node.AppendChild(&Node{Typ: NodeCloseBracket})
			node.AppendChild(&Node{Typ: NodeOpenParen})
			node.AppendChild(&Node{Typ: NodeLinkDest, Tokens: strToBytes(lute.domAttrValue(n, "src"))})
			linkTitle := lute.domAttrValue(n, "title")
			if "" != linkTitle {
				node.AppendChild(&Node{Typ: NodeLinkSpace})
				node.AppendChild(&Node{Typ: NodeLinkTitle, Tokens: []byte(linkTitle)})
			}
			node.AppendChild(&Node{Typ: NodeCloseParen})
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Input:
		node.Typ = NodeTaskListItemMarker
		if lute.hasAttr(n, "checked") {
			node.taskListItemChecked = true
		}
		tree.context.tip.AppendChild(node)
		if nil != node.Parent.Parent {
			node.Parent.Parent.listData.typ = 3
		}
	case atom.Del, atom.S, atom.Strike:
		node.Typ = NodeStrikethrough
		marker := "~"
		node.AppendChild(&Node{Typ: NodeStrikethrough1OpenMarker, Tokens: strToBytes(marker)})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Table:
		node.Typ = NodeTable
		var tableAligns []int
		for th := n.FirstChild.FirstChild.FirstChild; nil != th; th = th.NextSibling {
			align := lute.domAttrValue(th, "align")
			switch align {
			case "left":
				tableAligns = append(tableAligns, 1)
			case "center":
				tableAligns = append(tableAligns, 2)
			case "right":
				tableAligns = append(tableAligns, 3)
			default:
				tableAligns = append(tableAligns, 0)
			}
		}
		node.tableAligns = tableAligns
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Thead:
		node.Typ = NodeTableHead
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Tbody:
	case atom.Tr:
		node.Typ = NodeTableRow
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Th, atom.Td:
		node.Typ = NodeTableCell
		align := lute.domAttrValue(n, "align")
		var tableAlign int
		switch align {
		case "left":
			tableAlign = 1
		case "center":
			tableAlign = 2
		case "right":
			tableAlign = 3
		default:
			tableAlign = 0
		}
		node.tableCellAlign = tableAlign
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Span:
		if nil == n.FirstChild {
			return
		}
	case atom.Font:
		return
	case atom.Details:
		node.Typ = NodeHTMLBlock
		node.Tokens = lute.domHTML(n)
		node.Tokens = bytes.SplitAfter(node.Tokens, []byte("</summary>"))[0]
		tree.context.tip.AppendChild(node)
	case atom.Summary:
		return
	default:
		node.Typ = NodeHTMLBlock
		tokens := lute.domHTML(n)
		node.Tokens = tokens
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTByDOM(c, tree)
	}

	switch n.DataAtom {
	case atom.Em, atom.I:
		marker := "*"
		node.AppendChild(&Node{Typ: NodeEmA6kCloseMarker, Tokens: strToBytes(marker)})
	case atom.Strong, atom.B:
		marker := "**"
		node.AppendChild(&Node{Typ: NodeStrongA6kCloseMarker, Tokens: strToBytes(marker)})
	case atom.A:
		node.AppendChild(&Node{Typ: NodeCloseBracket})
		node.AppendChild(&Node{Typ: NodeOpenParen})
		node.AppendChild(&Node{Typ: NodeLinkDest, Tokens: strToBytes(lute.domAttrValue(n, "href"))})
		linkTitle := lute.domAttrValue(n, "title")
		if "" != linkTitle {
			node.AppendChild(&Node{Typ: NodeLinkSpace})
			node.AppendChild(&Node{Typ: NodeLinkTitle, Tokens: strToBytes(linkTitle)})
		}
		node.AppendChild(&Node{Typ: NodeCloseParen})
	case atom.Del, atom.S, atom.Strike:
		marker := "~"
		node.AppendChild(&Node{Typ: NodeStrikethrough1CloseMarker, Tokens: strToBytes(marker)})
	case atom.Details:
		tree.context.tip.AppendChild(&Node{Typ: NodeHTMLBlock, Tokens: []byte("</details>")})
	}
}
