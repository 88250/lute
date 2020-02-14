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
	"strconv"
	"strings"

	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
)

// caret 插入符 \u2038。
const caret = "‸"

// zwsp 零宽空格。
const zwsp = "\u200b"

// Md2HTML 将 markdown 转换为标准 HTML，用于源码模式预览。
func (lute *Lute) Md2HTML(markdown string) (html string) {
	lute.VditorWYSIWYG = false
	html, err := lute.MarkdownStr("", markdown)
	if nil != err {
		html = err.Error()
	}
	return
}

// FormatMd 将 markdown 进行格式化输出 formatted，用于源码模式格式化。
func (lute *Lute) FormatMd(markdown string) (formatted string) {
	formatted, err := lute.FormatStr("", markdown)
	if nil != err {
		formatted = err.Error()
	}
	return
}

// SpinVditorDOM 自旋 Vditor DOM，用于所见即所得模式下的编辑。
func (lute *Lute) SpinVditorDOM(htmlStr string) (html string) {
	lute.VditorWYSIWYG = true

	// 替换插入符
	htmlStr = strings.ReplaceAll(htmlStr, "<wbr>", caret)

	markdown := lute.vditorDOM2Md(htmlStr)

	tree, err := lute.parse("", []byte(markdown))
	if nil != err {
		html = err.Error()
		return
	}

	renderer := lute.newVditorRenderer(tree)
	var output []byte
	output, err = renderer.Render()
	if nil != err {
		html = err.Error()
		return
	}

	// 替换插入符
	html = strings.ReplaceAll(string(output), caret, "<wbr>")
	return html
}

// HTML2VditorDOM 将 HTML 转换为 Vditor DOM，用于所见即所得模式下粘贴。
func (lute *Lute) HTML2VditorDOM(htmlStr string) (html string) {
	lute.VditorWYSIWYG = true

	markdown, err := lute.HTML2Markdown(htmlStr)
	if nil != err {
		html = err.Error()
		return
	}

	var tree *Tree
	tree, err = lute.parse("", []byte(markdown))
	if nil != err {
		html = err.Error()
		return
	}

	renderer := lute.newVditorRenderer(tree)
	for nodeType, render := range lute.HTML2VditorDOMRendererFuncs {
		renderer.extRendererFuncs[nodeType] = render
	}
	var output []byte
	output, err = renderer.Render()
	if nil != err {
		html = err.Error()
	}
	html = string(output)
	return
}

// VditorDOM2HTML 将 Vditor DOM 转换为 HTML，用于 Vditor.getHTML() 接口。
func (lute *Lute) VditorDOM2HTML(vhtml string) (html string) {
	lute.VditorWYSIWYG = true

	markdown := lute.vditorDOM2Md(vhtml)
	html = lute.Md2HTML(markdown)
	return
}

// Md2VditorDOM 将 markdown 转换为 Vditor DOM，用于从源码模式切换至所见即所得模式。
func (lute *Lute) Md2VditorDOM(markdown string) (html string) {
	lute.VditorWYSIWYG = true

	tree, err := lute.parse("", []byte(markdown))
	if nil != err {
		html = err.Error()
		return
	}

	renderer := lute.newVditorRenderer(tree)
	var output []byte
	output, err = renderer.Render()
	if nil != err {
		html = err.Error()
	}
	html = string(output)
	return
}

// VditorDOM2Md 将 Vditor DOM 转换为 markdown，用于从所见即所得模式切换至源码模式。
func (lute *Lute) VditorDOM2Md(htmlStr string) (markdown string) {
	lute.VditorWYSIWYG = true

	md := lute.vditorDOM2Md(htmlStr)
	md = lute.FormatMd(md) // 再格式化一次处理表格对齐
	return strings.ReplaceAll(md, zwsp, "")
}

// RenderEChartsJSON 用于渲染 ECharts JSON 格式数据。
func (lute *Lute) RenderEChartsJSON(markdown string) (json string) {
	tree, err := lute.parse("", []byte(markdown))
	if nil != err {
		json = err.Error()
		return
	}

	renderer := lute.newEChartsJSONRenderer(tree)
	var output []byte
	output, err = renderer.Render()
	if nil != err {
		json = err.Error()
		return
	}
	json = string(output)
	return
}

// HTML2Md 用于将 HTML 转换为 markdown。
func (lute *Lute) HTML2Md(html string) (markdown string) {
	markdown, err := lute.HTML2Markdown(html)
	if nil != err {
		markdown = err.Error()
		return
	}
	return
}

func (lute *Lute) vditorDOM2Md(htmlStr string) (markdown string) {
	// 删掉插入符
	htmlStr = strings.ReplaceAll(htmlStr, "<wbr>", "")

	// 将字符串解析为 DOM 树

	reader := strings.NewReader(htmlStr)
	htmlRoot := &html.Node{Type: html.ElementNode}
	htmlNodes, err := html.ParseFragment(reader, htmlRoot)
	if nil != err {
		markdown = err.Error()
		return
	}

	// 将 HTML 树转换为 Markdown AST

	tree := &Tree{Name: "", Root: &Node{typ: NodeDocument}, context: &Context{option: lute.options}}
	tree.context.tip = tree.Root
	for _, htmlNode := range htmlNodes {
		lute.genASTByVditorDOM(htmlNode, tree)
	}

	// 调整树结构

	Walk(tree.Root, func(n *Node, entering bool) WalkStatus {
		if entering {
			switch n.typ {
			case NodeInlineHTML, NodeCodeSpan, NodeInlineMath, NodeHTMLBlock, NodeCodeBlockCode, NodeMathBlockContent:
				n.tokens = unescapeHTML(n.tokens)
			case NodeList:
				// 浏览器生成的子列表是 ul.ul 形式，需要将其调整为 ul.li.ul
				if nil != n.parent && NodeList == n.parent.typ {
					if previousLi := n.previous; nil != previousLi {
						previousLi.AppendChild(n)
					}
				}
			}
		}
		return WalkContinue
	})

	// 将 AST 进行 Markdown 格式化渲染

	var formatted []byte
	renderer := lute.newFormatRenderer(tree)
	formatted, err = renderer.Render()
	if nil != err {
		markdown = err.Error()
		return
	}
	markdown = string(formatted)
	return
}

// genASTByVditorDOM 根据指定的 Vditor DOM 节点 n 进行深度优先遍历并逐步生成 Markdown 语法树 tree。
func (lute *Lute) genASTByVditorDOM(n *html.Node, tree *Tree) {
	dataRender := lute.domAttrValue(n, "data-render")
	if "false" == dataRender {
		return
	}

	dataType := lute.domAttrValue(n, "data-type")

	if atom.Div == n.DataAtom && ("code-block" == dataType || "html-block" == dataType || "html-inline" == dataType || "math-block" == dataType || "math-inline" == dataType ||
		"backslash" == dataType) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lute.genASTByVditorDOM(c, tree)
		}
		return
	}

	class := lute.domAttrValue(n, "class")
	content := strings.ReplaceAll(n.Data, zwsp, "")
	node := &Node{typ: NodeText, tokens: []byte(content)}
	switch n.DataAtom {
	case 0:
		if "" == content {
			return
		}

		if nil != n.Parent && atom.A == n.Parent.DataAtom {
			node.typ = NodeLinkText
		}
		tree.context.tip.AppendChild(node)
	case atom.P:
		if nil != n.Parent && atom.Blockquote == n.Parent.DataAtom && "" == strings.TrimSpace(lute.domText(n)) { // vditorDOM2MdTests case 53
			return
		}

		node.typ = NodeParagraph
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		if "" == strings.TrimSpace(lute.domText(n)) {
			return
		}
		node.typ = NodeHeading
		node.headingLevel = int(node.tokens[1] - byte('0'))
		node.AppendChild(&Node{typ: NodeHeadingC8hMarker, tokens: []byte(strings.Repeat("#", node.headingLevel))})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Hr:
		node.typ = NodeThematicBreak
		tree.context.tip.AppendChild(node)
	case atom.Blockquote:
		content := strings.TrimSpace(lute.domText(n))
		if "" == content || "&gt;" == content || caret == content {
			return
		}

		node.typ = NodeBlockquote
		node.AppendChild(&Node{typ: NodeBlockquoteMarker, tokens: []byte(">")})
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
		if nil != n.FirstChild && nil == n.FirstChild.NextSibling && (atom.Ul == n.FirstChild.DataAtom || atom.Ol == n.FirstChild.DataAtom) {
			break
		}

		if p := n.FirstChild; nil != p && atom.P == p.DataAtom && nil != p.NextSibling && atom.P == p.NextSibling.DataAtom {
			tree.context.tip.tight = false
		}

		node.typ = NodeListItem
		marker := lute.domAttrValue(n, "data-marker")
		var bullet byte
		if "" == marker {
			if atom.Ol == n.Parent.DataAtom {
				if startAttr := lute.domAttrValue(n.Parent, "start"); "" == startAttr {
					marker = "1."
				} else {
					marker = startAttr + "."
				}
			} else {
				marker = lute.domAttrValue(n.Parent, "data-marker")
				if "" == marker {
					marker = "*"
				}
				bullet = marker[0]
			}
		} else {
			if nil != n.Parent {
				if atom.Ol == n.Parent.DataAtom {
					if "*" == marker || "-" == marker || "+" == marker {
						marker = "1."
					}
					if "1." != marker && "1)" != marker && nil != n.Parent.Parent && (atom.Ol == n.Parent.Parent.DataAtom || atom.Ul == n.Parent.Parent.DataAtom) {
						// 子有序列表必须从 1 开始
						marker = "1."
					}
					if "1." != marker && "1)" != marker && atom.Ol == n.Parent.DataAtom && n.Parent.FirstChild == n {
						marker = "1."
					}
				} else {
					if "*" != marker && "-" != marker && "+" != marker {
						marker = "*"
					}
					bullet = marker[0]
				}
			}
		}
		node.listData = &listData{marker: []byte(marker), bulletChar: bullet}
		if 0 == bullet {
			node.listData.num, _ = strconv.Atoi(string(marker[0]))
			node.listData.delimiter = marker[len(marker)-1]
		}

		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Pre:
		if atom.Code == n.FirstChild.DataAtom {
			marker := lute.domAttrValue(n.Parent, "data-marker")
			if "" == marker {
				marker = "```"
			}

			var codeTokens []byte
			if nil != n.FirstChild.FirstChild {
				codeTokens = []byte(n.FirstChild.FirstChild.Data)
			}

			divDataType := lute.domAttrValue(n.Parent, "data-type")
			switch divDataType {
			case "math-block":
				node.typ = NodeMathBlock
				node.AppendChild(&Node{typ: NodeMathBlockOpenMarker})
				node.AppendChild(&Node{typ: NodeMathBlockContent, tokens: codeTokens})
				node.AppendChild(&Node{typ: NodeMathBlockCloseMarker})
				tree.context.tip.AppendChild(node)
			case "html-block":
				node.typ = NodeHTMLBlock
				node.tokens = codeTokens
				tree.context.tip.AppendChild(node)
			default:
				node.typ = NodeCodeBlock
				node.isFencedCodeBlock = true
				node.AppendChild(&Node{typ: NodeCodeBlockFenceOpenMarker, tokens: []byte(marker), codeBlockFenceLen: len(marker)})
				node.AppendChild(&Node{typ: NodeCodeBlockFenceInfoMarker})
				class := lute.domAttrValue(n.FirstChild, "class")
				if strings.Contains(class, "language-") {
					language := class[len("language-"):]
					node.lastChild.codeBlockInfo = []byte(language)
				}

				content := &Node{typ: NodeCodeBlockCode, tokens: codeTokens}
				node.AppendChild(content)
				node.AppendChild(&Node{typ: NodeCodeBlockFenceCloseMarker, tokens: []byte(marker), codeBlockFenceLen: len(marker)})
				tree.context.tip.AppendChild(node)
			}
		}
		return
	case atom.Em, atom.I:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		text := strings.TrimSpace(lute.domText(n))
		if zwsp == text || "" == text {
			return
		}
		if caret == text {
			node.tokens = []byte(caret)
			tree.context.tip.AppendChild(node)
			return
		}

		node.typ = NodeEmphasis
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "*"
		}
		if "_" == marker {
			node.AppendChild(&Node{typ: NodeEmU8eOpenMarker, tokens: []byte(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeEmA6kOpenMarker, tokens: []byte(marker)})
		}
		tree.context.tip.AppendChild(node)

		if nil != n.FirstChild && caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&Node{typ: NodeText, tokens: []byte(caret)})
			if "_" == marker {
				node.AppendChild(&Node{typ: NodeEmU8eCloseMarker, tokens: []byte(marker)})
			} else {
				node.AppendChild(&Node{typ: NodeEmA6kCloseMarker, tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, zwsp, "")

		// 开头结尾空格后会形成 * foo * 导致强调、加粗删除线标记失效，这里将空格移到右标记符前后 _*foo*_
		if strings.HasPrefix(n.FirstChild.Data, " ") {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&Node{typ: NodeText, tokens: []byte(" ")})
		}
		if strings.HasSuffix(n.FirstChild.Data, " ") {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
		}

		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Strong, atom.B:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		text := strings.TrimSpace(lute.domText(n))
		if zwsp == text || "" == text {
			return
		}
		if caret == text {
			node.tokens = []byte(caret)
			tree.context.tip.AppendChild(node)
			return
		}

		node.typ = NodeStrong
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "**"
		}
		if "__" == marker {
			node.AppendChild(&Node{typ: NodeStrongU8eOpenMarker, tokens: []byte(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeStrongA6kOpenMarker, tokens: []byte(marker)})
		}
		tree.context.tip.AppendChild(node)

		if nil != n.FirstChild && caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&Node{typ: NodeText, tokens: []byte(caret)})
			if "__" == marker {
				node.AppendChild(&Node{typ: NodeStrongU8eCloseMarker, tokens: []byte(marker)})
			} else {
				node.AppendChild(&Node{typ: NodeStrongA6kCloseMarker, tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, zwsp, "")
		if strings.HasPrefix(n.FirstChild.Data, " ") {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&Node{typ: NodeText, tokens: []byte(" ")})
		}
		if strings.HasSuffix(n.FirstChild.Data, " ") {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
		}

		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Del, atom.S, atom.Strike:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		text := strings.TrimSpace(lute.domText(n))
		if zwsp == text || "" == text {
			return
		}
		if caret == text {
			node.tokens = []byte(caret)
			tree.context.tip.AppendChild(node)
			return
		}

		node.typ = NodeStrikethrough
		marker := lute.domAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&Node{typ: NodeStrikethrough1OpenMarker, tokens: []byte(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeStrikethrough2OpenMarker, tokens: []byte(marker)})
		}
		tree.context.tip.AppendChild(node)

		if nil != n.FirstChild && caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&Node{typ: NodeText, tokens: []byte(caret)})
			if "~" == marker {
				node.AppendChild(&Node{typ: NodeStrikethrough1CloseMarker, tokens: []byte(marker)})
			} else {
				node.AppendChild(&Node{typ: NodeStrikethrough2CloseMarker, tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, zwsp, "")

		if strings.HasPrefix(n.FirstChild.Data, " ") {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&Node{typ: NodeText, tokens: []byte(" ")})
		}
		if strings.HasSuffix(n.FirstChild.Data, " ") {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
		}

		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Code:
		if nil == n.FirstChild {
			return
		}
		contentStr := strings.ReplaceAll(n.FirstChild.Data, zwsp, "")
		if caret == contentStr {
			node.tokens = []byte(caret)
			tree.context.tip.AppendChild(node)
			return
		}
		if "" == contentStr {
			return
		}
		codeTokens := []byte(contentStr)
		content := &Node{typ: NodeCodeSpanContent, tokens: codeTokens}
		marker := lute.domAttrValue(n, "marker")
		if "" == marker {
			marker = "`"
		}
		node.typ = NodeCodeSpan
		node.codeMarkerLen = len(marker)
		node.AppendChild(&Node{typ: NodeCodeSpanOpenMarker})
		node.AppendChild(content)
		node.AppendChild(&Node{typ: NodeCodeSpanCloseMarker})
		tree.context.tip.AppendChild(node)
		return
	case atom.Br:
		if nil != n.Parent {
			if lute.parentIs(n, atom.Td, atom.Th) {
				if (nil == n.PrevSibling || caret == n.PrevSibling.Data) && (nil == n.NextSibling || caret == n.NextSibling.Data) {
					return
				}
				if nil == n.NextSibling {
					return // 删掉表格中结尾的 br
				}

				node.typ = NodeInlineHTML
				node.tokens = []byte("<br />")
				tree.context.tip.AppendChild(node)
				return
			}
			if atom.P == n.Parent.DataAtom {
				if nil != n.Parent.NextSibling && (atom.Ul == n.Parent.NextSibling.DataAtom || atom.Ol == n.Parent.NextSibling.DataAtom || atom.Blockquote == n.Parent.NextSibling.DataAtom) {
					tree.context.tip.AppendChild(&Node{typ: NodeText, tokens: []byte(zwsp)})
					return
				}
				if nil != n.Parent.Parent && nil != n.Parent.Parent.NextSibling && atom.Li == n.Parent.Parent.NextSibling.DataAtom {
					tree.context.tip.AppendChild(&Node{typ: NodeText, tokens: []byte(zwsp)})
					return
				}
			}
		}

		node.typ = NodeHardBreak
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
		imgClass := class
		imgAlt := lute.domAttrValue(n, "alt")
		if "emoji" == imgClass {
			node.typ = NodeEmoji
			emojiImg := &Node{typ: NodeEmojiImg, tokens: tree.emojiImgTokens(imgAlt, lute.domAttrValue(n, "src"))}
			emojiImg.AppendChild(&Node{typ: NodeEmojiAlias, tokens: []byte(":" + imgAlt + ":")})
			node.AppendChild(emojiImg)
		} else {
			node.typ = NodeImage
			node.AppendChild(&Node{typ: NodeBang})
			node.AppendChild(&Node{typ: NodeOpenBracket})
			if "" != imgAlt {
				node.AppendChild(&Node{typ: NodeLinkText, tokens: []byte(imgAlt)})
			}
			node.AppendChild(&Node{typ: NodeCloseBracket})
			node.AppendChild(&Node{typ: NodeOpenParen})
			src := lute.domAttrValue(n, "src")
			if "" != lute.LinkBase {
				src = strings.ReplaceAll(src, lute.LinkBase, "")
			}
			node.AppendChild(&Node{typ: NodeLinkDest, tokens: []byte(src)})
			linkTitle := lute.domAttrValue(n, "title")
			if "" != linkTitle {
				node.AppendChild(&Node{typ: NodeLinkSpace})
				node.AppendChild(&Node{typ: NodeLinkTitle, tokens: []byte(linkTitle)})
			}
			node.AppendChild(&Node{typ: NodeCloseParen})
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Input:
		if nil == n.Parent || nil == n.Parent.Parent || (atom.P != n.Parent.DataAtom && atom.Li != n.Parent.DataAtom) {
			// 仅允许 input 出现在任务列表中
			return
		}
		node.typ = NodeTaskListItemMarker
		if lute.hasAttr(n, "checked") {
			node.taskListItemChecked = true
		}
		tree.context.tip.AppendChild(node)
		if nil != node.parent.parent && nil != node.parent.parent.listData { // ul.li.input
			node.parent.parent.listData.typ = 3
		}
		if nil != node.parent.parent.parent && nil != node.parent.parent.parent.listData { // ul.li.p.input
			node.parent.parent.parent.listData.typ = 3
		}
	case atom.Table:
		node.typ = NodeTable
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
			break
		}
		var codeTokens []byte
		if zwsp == n.FirstChild.Data && "" == lute.domAttrValue(n, "style") && nil != n.FirstChild.NextSibling {
			codeTokens = []byte(n.FirstChild.NextSibling.FirstChild.Data)
		} else if atom.Code == n.FirstChild.DataAtom {
			codeTokens = []byte(n.FirstChild.FirstChild.Data)
			if zwsp == string(codeTokens) {
				break
			}
		} else {
			break
		}
		if "math-inline" == dataType {
			node.typ = NodeInlineMath
			node.AppendChild(&Node{typ: NodeInlineMathOpenMarker})
			node.AppendChild(&Node{typ: NodeInlineMathContent, tokens: codeTokens})
			node.AppendChild(&Node{typ: NodeInlineMathCloseMarker})
			tree.context.tip.AppendChild(node)
		} else if "html-inline" == dataType {
			node.typ = NodeInlineHTML
			node.tokens = codeTokens
			tree.context.tip.AppendChild(node)
		} else if "code-inline" == dataType {
			node.tokens = codeTokens
			tree.context.tip.AppendChild(node)
		}
		return
	case atom.Font:
		return
	case atom.Details:
		node.typ = NodeHTMLBlock
		node.tokens = lute.domHTML(n)
		node.tokens = bytes.SplitAfter(node.tokens, []byte("</summary>"))[0]
		tree.context.tip.AppendChild(node)
	case atom.Kbd:
		node.typ = NodeInlineHTML
		node.tokens = lute.domHTML(n)
		tree.context.tip.AppendChild(node)
		return
	case atom.Summary:
		return
	default:
		node.typ = NodeHTMLBlock
		node.tokens = lute.domHTML(n)
		tree.context.tip.AppendChild(node)
		return
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
			node.AppendChild(&Node{typ: NodeEmU8eCloseMarker, tokens: []byte(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeEmA6kCloseMarker, tokens: []byte(marker)})
		}
	case atom.Strong, atom.B:
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "**"
		}
		if "__" == marker {
			node.AppendChild(&Node{typ: NodeStrongU8eCloseMarker, tokens: []byte(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeStrongA6kCloseMarker, tokens: []byte(marker)})
		}
	case atom.A:
		node.AppendChild(&Node{typ: NodeCloseBracket})
		node.AppendChild(&Node{typ: NodeOpenParen})
		href := lute.domAttrValue(n, "href")
		if "" != lute.LinkBase {
			href = strings.ReplaceAll(href, lute.LinkBase, "")
		}
		node.AppendChild(&Node{typ: NodeLinkDest, tokens: []byte(href)})
		linkTitle := lute.domAttrValue(n, "title")
		if "" != linkTitle {
			node.AppendChild(&Node{typ: NodeLinkSpace})
			node.AppendChild(&Node{typ: NodeLinkTitle, tokens: []byte(linkTitle)})
		}
		node.AppendChild(&Node{typ: NodeCloseParen})
	case atom.Del, atom.S, atom.Strike:
		marker := lute.domAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&Node{typ: NodeStrikethrough1CloseMarker, tokens: []byte(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeStrikethrough2CloseMarker, tokens: []byte(marker)})
		}
	case atom.Details:
		tree.context.tip.AppendChild(&Node{typ: NodeHTMLBlock, tokens: []byte("</details>")})
	}
}

func (context *Context) parentTip(n *html.Node) {
	if tip := context.tip.parent; nil != tip {
		context.tip = context.tip.parent
	}
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
	if nil == n {
		return ""
	}

	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}
	return ""
}

func (lute *Lute) domCode(n *html.Node) string {
	buf := &bytes.Buffer{}
	lute.domCode0(n, buf)
	return buf.String()
}

func (lute *Lute) domCode0(n *html.Node, buffer *bytes.Buffer) {
	if nil == n {
		return
	}
	switch n.DataAtom {
	case 0:
		buffer.WriteString(n.Data)
	default:
		buffer.Write(lute.domHTML(n))
		return
	}

	for child := n.FirstChild; nil != child; child = child.NextSibling {
		lute.domCode0(child, buffer)
	}
}

func (lute *Lute) parentIs(n *html.Node, parentTypes ...atom.Atom) bool {
	for p := n.Parent; nil != p; p = p.Parent {
		for _, pt := range parentTypes {
			if pt == p.DataAtom {
				return true
			}
		}
	}
	return false
}

func (lute *Lute) domText(n *html.Node) string {
	buf := &bytes.Buffer{}
	for next := n; nil != next; next = next.NextSibling {
		lute.domText0(n, buf)
	}
	return buf.String()
}

func (lute *Lute) domText0(n *html.Node, buffer *bytes.Buffer) {
	if nil == n {
		return
	}
	switch n.DataAtom {
	case 0:
		buffer.WriteString(n.Data)
	case atom.Br:
		buffer.WriteString("\n")
	}

	for child := n.FirstChild; nil != child; child = child.NextSibling {
		lute.domText0(child, buffer)
	}
}

func (lute *Lute) domHTML(n *html.Node) []byte {
	buf := &bytes.Buffer{}
	html.Render(buf, n)
	return bytes.ReplaceAll(buf.Bytes(), []byte(zwsp), []byte(""))
}
