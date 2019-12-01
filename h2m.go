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
	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"strings"
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

	tree := &Tree{Name: "", Root: &Node{typ: NodeDocument}, context: &Context{option: lute.options}}
	tree.context.tip = tree.Root
	for _, htmlNode := range htmlNodes {
		lute.genASTByDOM(htmlNode, tree)
	}

	// 调整树结构

	Walk(tree.Root, func(n *Node, entering bool) (status WalkStatus, e error) {
		if entering {
			if NodeList == n.typ {
				// ul.ul => ul.li.ul
				if nil != n.parent && NodeList == n.parent.typ {
					previousLi := n.previous
					if nil != previousLi {
						n.Unlink()
						previousLi.AppendChild(n)
					}
				}
			} else if NodeListItem == n.typ {
				if nil != n.parent && NodeList != n.parent.typ {
					// doc.li => doc.ul.li
					previousList := n.previous
					if nil != previousList {
						n.Unlink()
						previousList.AppendChild(n)
					}
				}
			}
		}
		return WalkContinue, nil
	})

	// 将 AST 进行 Markdown 格式化渲染

	var formatted []byte
	renderer := lute.newFormatRenderer(tree)
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

	node := &Node{typ: NodeText, tokens: strToBytes(n.Data)}
	switch n.DataAtom {
	case 0:
		if nil != n.Parent && atom.A == n.Parent.DataAtom {
			node.typ = NodeLinkText
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
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
		node.tight = true
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Li:
		node.typ = NodeListItem
		marker := "*"
		if atom.Ol == n.Parent.DataAtom {
			start := lute.domAttrValue(n.Parent, "start")
			if "" == start {
				marker = "1."
			} else {
				marker = start + "."
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
		buf := &bytes.Buffer{}
		firstc := n.FirstChild
		if nil != firstc {
			if atom.Code == firstc.DataAtom {
				class := lute.domAttrValue(firstc, "class")
				if strings.Contains(class, "language-") {
					language := class[len("language-"):]
					node.lastChild.codeBlockInfo = []byte(language)
				}
				firstc = firstc.FirstChild
			}

			for c := firstc; nil != c; c = c.NextSibling {
				buf.WriteString(lute.domText(c))
			}
		}
		content := &Node{typ: NodeCodeBlockCode, tokens: buf.Bytes()}
		node.AppendChild(content)
		node.AppendChild(&Node{typ: NodeCodeBlockFenceCloseMarker, tokens: strToBytes("```"), codeBlockFenceLen: 3})
		tree.context.tip.AppendChild(node)
		return
	case atom.Em, atom.I:
		node.typ = NodeEmphasis
		marker := "*"
		node.AppendChild(&Node{typ: NodeEmA6kOpenMarker, tokens: strToBytes(marker)})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Strong, atom.B:
		node.typ = NodeStrong
		marker := "**"
		node.AppendChild(&Node{typ: NodeStrongA6kOpenMarker, tokens: strToBytes(marker)})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Code:
		buf := &bytes.Buffer{}
		for c := n.FirstChild; nil != c; c = c.NextSibling {
			html.Render(buf, c)
		}
		content := &Node{typ: NodeCodeSpanContent, tokens: buf.Bytes()}
		node.typ = NodeCodeSpan
		node.AppendChild(&Node{typ: NodeCodeSpanOpenMarker, tokens: []byte("`")})
		node.AppendChild(content)
		node.AppendChild(&Node{typ: NodeCodeSpanCloseMarker, tokens: []byte("`")})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
		return
	case atom.Br:
		node.typ = NodeHardBreak
		node.tokens = strToBytes("\n")
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
	case atom.Del, atom.S, atom.Strike:
		node.typ = NodeStrikethrough
		marker := "~"
		node.AppendChild(&Node{typ: NodeStrikethrough1OpenMarker, tokens: strToBytes(marker)})
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
	case atom.Tbody:
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
	default:
		node.typ = NodeHTMLBlock
		buf := &bytes.Buffer{}
		html.Render(buf, n)
		tokens := buf.Bytes()
		node.tokens = tokens
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
		node.AppendChild(&Node{typ: NodeEmA6kCloseMarker, tokens: strToBytes(marker)})
	case atom.Strong, atom.B:
		marker := "**"
		node.AppendChild(&Node{typ: NodeStrongA6kCloseMarker, tokens: strToBytes(marker)})
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
	case atom.Del, atom.S, atom.Strike:
		marker := "~"
		node.AppendChild(&Node{typ: NodeStrikethrough1CloseMarker, tokens: strToBytes(marker)})
	}
}
