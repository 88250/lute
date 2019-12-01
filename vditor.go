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

	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
)

// 插入符 \u2038
const caret = "‸"

// Md2HTML 将 markdown 转换为标准 HTML，用于源码模式预览。
func (lute *Lute) Md2HTML(markdown string) (html string) {
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
	return
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
	markdown := lute.vditorDOM2Md(vhtml)
	html = lute.Md2HTML(markdown)
	return
}

// Md2VditorDOM 将 markdown 转换为 Vditor DOM，用于从源码模式切换至所见即所得模式。
func (lute *Lute) Md2VditorDOM(markdown string) (html string) {
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
	return lute.vditorDOM2Md(htmlStr)
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
					// .li => .ul.li
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
	if "code-block" == dataType {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lute.genASTByVditorDOM(c, tree)
		}
		return
	}

	if atom.Div == n.DataAtom && ("html-block" == dataType || "html-inline" == dataType || "math-block" == dataType || "math-inline" == dataType) {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lute.genASTByVditorDOM(c, tree)
		}
		return
	}

	class := lute.domAttrValue(n, "class")
	node := &Node{typ: NodeText, tokens: []byte(n.Data)}
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
		node.AppendChild(&Node{typ: NodeHeadingC8hMarker, tokens: []byte(strings.Repeat("#", node.headingLevel))})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Hr:
		node.typ = NodeThematicBreak
		tree.context.tip.AppendChild(node)
	case atom.Blockquote:
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
		node.listData = &listData{marker: []byte(marker)}
		if lute.firstChildIsText(n) {
			tree.context.tip.AppendChild(node)
			tree.context.tip = node
			node = &Node{typ: NodeParagraph}
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Pre:
		if atom.Code == n.FirstChild.DataAtom {
			switch dataType {
			case "math-block":
				node.typ = NodeMathBlock
				node.AppendChild(&Node{typ: NodeMathBlockOpenMarker})
				node.AppendChild(&Node{typ: NodeMathBlockContent, tokens: []byte(n.FirstChild.Data)})
				node.AppendChild(&Node{typ: NodeMathBlockCloseMarker})
				tree.context.tip.AppendChild(node)
			case "html-block":
				node.typ = NodeHTMLBlock
				node.tokens = []byte(n.FirstChild.Data)
				tree.context.tip.AppendChild(node)
			default:
				node.typ = NodeCodeBlock
				node.isFencedCodeBlock = true
				node.AppendChild(&Node{typ: NodeCodeBlockFenceOpenMarker, tokens: []byte("```"), codeBlockFenceLen: 3})
				node.AppendChild(&Node{typ: NodeCodeBlockFenceInfoMarker})
				class := lute.domAttrValue(n.FirstChild, "class")
				if strings.Contains(class, "language-") {
					language := class[len("language-"):]
					node.lastChild.codeBlockInfo = []byte(language)
				}
				content := &Node{typ: NodeCodeBlockCode, tokens: []byte(n.FirstChild.FirstChild.Data)}
				node.AppendChild(content)
				node.AppendChild(&Node{typ: NodeCodeBlockFenceCloseMarker, tokens: []byte("```"), codeBlockFenceLen: 3})
				tree.context.tip.AppendChild(node)
			}
		}
		return
	case atom.Em, atom.I:
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
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Strong, atom.B:
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
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Code:
		content := &Node{typ: NodeCodeSpanContent, tokens: []byte(n.FirstChild.Data)}
		node.typ = NodeCodeSpan
		node.AppendChild(&Node{typ: NodeCodeSpanOpenMarker, tokens: []byte("`")})
		node.AppendChild(content)
		node.AppendChild(&Node{typ: NodeCodeSpanCloseMarker, tokens: []byte("`")})
		tree.context.tip.AppendChild(node)
		return
	case atom.Br:
		node.typ = NodeHardBreak
		node.tokens = []byte("\n")
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
			node.AppendChild(&Node{typ: NodeLinkDest, tokens: []byte(lute.domAttrValue(n, "src"))})
			node.AppendChild(&Node{typ: NodeCloseParen})
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Input:
		node.typ = NodeTaskListItemMarker
		if lute.hasAttr(n, "checked") {
			node.taskListItemChecked = true
			node.tokens = []byte("[X]")
		} else {
			node.tokens = []byte("[ ]")
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
		if nil != node.parent.parent {
			node.parent.parent.listData.typ = 3
		}
	case atom.Del, atom.S, atom.Strike:
		node.typ = NodeStrikethrough
		marker := lute.domAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&Node{typ: NodeStrikethrough1OpenMarker, tokens: []byte(marker)})
		} else {
			node.AppendChild(&Node{typ: NodeStrikethrough2OpenMarker, tokens: []byte(marker)})
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
	case atom.Span:
		if "math-inline" == dataType {
			node.typ = NodeInlineMath
			node.AppendChild(&Node{typ: NodeInlineMathOpenMarker})
			node.AppendChild(&Node{typ: NodeInlineMathContent, tokens: []byte(n.FirstChild.FirstChild.Data)})
			node.AppendChild(&Node{typ: NodeInlineMathCloseMarker})
			tree.context.tip.AppendChild(node)
		} else if "html-inline" == dataType {
			node.typ = NodeInlineHTML
			node.tokens = []byte(n.FirstChild.FirstChild.Data)
			tree.context.tip.AppendChild(node)
		} else {
			node.tokens = []byte(lute.domText(n))
			tree.context.tip.AppendChild(node)
		}
		return
	default:
		node.typ = NodeHTMLBlock
		buf := &bytes.Buffer{}
		html.Render(buf, n)
		node.tokens = buf.Bytes()
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
		node.AppendChild(&Node{typ: NodeLinkDest, tokens: []byte(lute.domAttrValue(n, "href"))})
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

func (lute *Lute) domHTML(n *html.Node) string {
	buf := &bytes.Buffer{}
	lute.domHTML0(n, buf)
	return buf.String()
}

func (lute *Lute) domHTML0(n *html.Node, buffer *bytes.Buffer) {
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
		lute.domHTML0(child, buffer)
	}
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
