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
	"github.com/88250/lute/util"
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

	tree := &Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument}, context: &Context{option: lute.options}}
	tree.context.tip = tree.Root
	for _, htmlNode := range htmlNodes {
		lute.genASTByVditorDOM(htmlNode, tree)
	}

	// 调整树结构

	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n.Type {
			case ast.NodeInlineHTML, ast.NodeCodeSpan, ast.NodeInlineMath, ast.NodeHTMLBlock, ast.NodeCodeBlockCode, ast.NodeMathBlockContent:
				n.Tokens = util.UnescapeHTML(n.Tokens)
			case ast.NodeList:
				// 浏览器生成的子列表是 ul.ul 形式，需要将其调整为 ul.li.ul
				if nil != n.Parent && ast.NodeList == n.Parent.Type {
					if previousLi := n.Previous; nil != previousLi {
						previousLi.AppendChild(n)
					}
				}
			}
		}
		return ast.WalkContinue
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
	node := &ast.Node{Type: ast.NodeText, Tokens: []byte(content)}
	switch n.DataAtom {
	case 0:
		if "" == content {
			return
		}

		if nil != n.Parent && atom.A == n.Parent.DataAtom {
			node.Type = ast.NodeLinkText
		}
		tree.context.tip.AppendChild(node)
	case atom.P:
		if nil != n.Parent && atom.Blockquote == n.Parent.DataAtom && "" == strings.TrimSpace(lute.domText(n)) { // vditorDOM2MdTests case 53
			return
		}

		node.Type = ast.NodeParagraph
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		if "" == strings.TrimSpace(lute.domText(n)) {
			return
		}
		node.Type = ast.NodeHeading
		node.HeadingLevel = int(node.Tokens[1] - byte('0'))
		node.AppendChild(&ast.Node{Type: ast.NodeHeadingC8hMarker, Tokens: []byte(strings.Repeat("#", node.HeadingLevel))})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Hr:
		node.Type = ast.NodeThematicBreak
		tree.context.tip.AppendChild(node)
	case atom.Blockquote:
		content := strings.TrimSpace(lute.domText(n))
		if "" == content || "&gt;" == content || caret == content {
			return
		}

		node.Type = ast.NodeBlockquote
		node.AppendChild(&ast.Node{Type: ast.NodeBlockquoteMarker, Tokens: []byte(">")})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Ol, atom.Ul:
		node.Type = ast.NodeList
		node.ListData = &ast.ListData{}
		if atom.Ol == n.DataAtom {
			node.ListData.Typ = 1
		}
		tight := lute.domAttrValue(n, "data-tight")
		if "true" == tight || "" == tight {
			node.Tight = true
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Li:
		if nil != n.FirstChild && nil == n.FirstChild.NextSibling && (atom.Ul == n.FirstChild.DataAtom || atom.Ol == n.FirstChild.DataAtom) {
			break
		}

		if p := n.FirstChild; nil != p && atom.P == p.DataAtom && nil != p.NextSibling && atom.P == p.NextSibling.DataAtom {
			tree.context.tip.Tight = false
		}

		node.Type = ast.NodeListItem
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
		node.ListData = &ast.ListData{Marker: []byte(marker), BulletChar: bullet}
		if 0 == bullet {
			node.ListData.Num, _ = strconv.Atoi(string(marker[0]))
			node.ListData.Delimiter = marker[len(marker)-1]
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
				node.Type = ast.NodeMathBlock
				node.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker})
				node.AppendChild(&ast.Node{Type: ast.NodeMathBlockContent, Tokens: codeTokens})
				node.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker})
				tree.context.tip.AppendChild(node)
			case "html-block":
				node.Type = ast.NodeHTMLBlock
				node.Tokens = codeTokens
				tree.context.tip.AppendChild(node)
			default:
				node.Type = ast.NodeCodeBlock
				node.IsFencedCodeBlock = true
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: []byte(marker), CodeBlockFenceLen: len(marker)})
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker})
				class := lute.domAttrValue(n.FirstChild, "class")
				if strings.Contains(class, "language-") {
					language := class[len("language-"):]
					node.LastChild.CodeBlockInfo = []byte(language)
				}

				content := &ast.Node{Type: ast.NodeCodeBlockCode, Tokens: codeTokens}
				node.AppendChild(content)
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: []byte(marker), CodeBlockFenceLen: len(marker)})
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
			node.Tokens = []byte(caret)
			tree.context.tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeEmphasis
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "*"
		}
		if "_" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeEmU8eOpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeEmA6kOpenMarker, Tokens: []byte(marker)})
		}
		tree.context.tip.AppendChild(node)

		if nil != n.FirstChild && caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(caret)})
			if "_" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeEmU8eCloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, zwsp, "")

		// 开头结尾空格后会形成 * foo * 导致强调、加粗删除线标记失效，这里将空格移到右标记符前后 _*foo*_
		if strings.HasPrefix(n.FirstChild.Data, " ") {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
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
			node.Tokens = []byte(caret)
			tree.context.tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeStrong
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "**"
		}
		if "__" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongU8eOpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kOpenMarker, Tokens: []byte(marker)})
		}
		tree.context.tip.AppendChild(node)

		if nil != n.FirstChild && caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(caret)})
			if "__" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeStrongU8eCloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kCloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, zwsp, "")
		if strings.HasPrefix(n.FirstChild.Data, " ") {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
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
			node.Tokens = []byte(caret)
			tree.context.tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeStrikethrough
		marker := lute.domAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1OpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2OpenMarker, Tokens: []byte(marker)})
		}
		tree.context.tip.AppendChild(node)

		if nil != n.FirstChild && caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(caret)})
			if "~" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1CloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2CloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, zwsp, "")

		if strings.HasPrefix(n.FirstChild.Data, " ") {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
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
			node.Tokens = []byte(caret)
			tree.context.tip.AppendChild(node)
			return
		}
		if "" == contentStr {
			return
		}
		codeTokens := []byte(contentStr)
		content := &ast.Node{Type: ast.NodeCodeSpanContent, Tokens: codeTokens}
		marker := lute.domAttrValue(n, "marker")
		if "" == marker {
			marker = "`"
		}
		node.Type = ast.NodeCodeSpan
		node.CodeMarkerLen = len(marker)
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanOpenMarker})
		node.AppendChild(content)
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanCloseMarker})
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

				node.Type = ast.NodeInlineHTML
				node.Tokens = []byte("<br />")
				tree.context.tip.AppendChild(node)
				return
			}
			if atom.P == n.Parent.DataAtom {
				if nil != n.Parent.NextSibling && (atom.Ul == n.Parent.NextSibling.DataAtom || atom.Ol == n.Parent.NextSibling.DataAtom || atom.Blockquote == n.Parent.NextSibling.DataAtom) {
					tree.context.tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(zwsp)})
					return
				}
				if nil != n.Parent.Parent && nil != n.Parent.Parent.NextSibling && atom.Li == n.Parent.Parent.NextSibling.DataAtom {
					tree.context.tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(zwsp)})
					return
				}
			}
		}

		node.Type = ast.NodeHardBreak
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.A:
		node.Type = ast.NodeLink
		node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Img:
		imgClass := class
		imgAlt := lute.domAttrValue(n, "alt")
		if "emoji" == imgClass {
			node.Type = ast.NodeEmoji
			emojiImg := &ast.Node{Type: ast.NodeEmojiImg, Tokens: tree.emojiImgTokens(imgAlt, lute.domAttrValue(n, "src"))}
			emojiImg.AppendChild(&ast.Node{Type: ast.NodeEmojiAlias, Tokens: []byte(":" + imgAlt + ":")})
			node.AppendChild(emojiImg)
		} else {
			node.Type = ast.NodeImage
			node.AppendChild(&ast.Node{Type: ast.NodeBang})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
			if "" != imgAlt {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkText, Tokens: []byte(imgAlt)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			src := lute.domAttrValue(n, "src")
			if "" != lute.LinkBase {
				src = strings.ReplaceAll(src, lute.LinkBase, "")
			}
			node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: []byte(src)})
			linkTitle := lute.domAttrValue(n, "title")
			if "" != linkTitle {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
				node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: []byte(linkTitle)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
		}
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Input:
		if nil == n.Parent || nil == n.Parent.Parent || (atom.P != n.Parent.DataAtom && atom.Li != n.Parent.DataAtom) {
			// 仅允许 input 出现在任务列表中
			return
		}
		node.Type = ast.NodeTaskListItemMarker
		if lute.hasAttr(n, "checked") {
			node.TaskListItemChecked = true
		}
		tree.context.tip.AppendChild(node)
		if nil != node.Parent.Parent && nil != node.Parent.Parent.ListData { // ul.li.input
			node.Parent.Parent.ListData.Typ = 3
		}
		if nil != node.Parent.Parent.Parent && nil != node.Parent.Parent.Parent.ListData { // ul.li.p.input
			node.Parent.Parent.Parent.ListData.Typ = 3
		}
	case atom.Table:
		node.Type = ast.NodeTable
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
		node.TableAligns = tableAligns
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Thead:
		node.Type = ast.NodeTableHead
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Tbody:
	case atom.Tr:
		node.Type = ast.NodeTableRow
		tree.context.tip.AppendChild(node)
		tree.context.tip = node
		defer tree.context.parentTip(n)
	case atom.Th, atom.Td:
		node.Type = ast.NodeTableCell
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
		node.TableCellAlign = tableAlign
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
			node.Type = ast.NodeInlineMath
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathOpenMarker})
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathContent, Tokens: codeTokens})
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathCloseMarker})
			tree.context.tip.AppendChild(node)
		} else if "html-inline" == dataType {
			node.Type = ast.NodeInlineHTML
			node.Tokens = codeTokens
			tree.context.tip.AppendChild(node)
		} else if "code-inline" == dataType {
			node.Tokens = codeTokens
			tree.context.tip.AppendChild(node)
		}
		return
	case atom.Font:
		return
	case atom.Details:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = lute.domHTML(n)
		node.Tokens = bytes.SplitAfter(node.Tokens, []byte("</summary>"))[0]
		tree.context.tip.AppendChild(node)
	case atom.Kbd:
		node.Type = ast.NodeInlineHTML
		node.Tokens = lute.domHTML(n)
		tree.context.tip.AppendChild(node)
		return
	case atom.Summary:
		return
	default:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = lute.domHTML(n)
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
			node.AppendChild(&ast.Node{Type: ast.NodeEmU8eCloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker, Tokens: []byte(marker)})
		}
	case atom.Strong, atom.B:
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "**"
		}
		if "__" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongU8eCloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kCloseMarker, Tokens: []byte(marker)})
		}
	case atom.A:
		node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
		node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
		href := lute.domAttrValue(n, "href")
		if "" != lute.LinkBase {
			href = strings.ReplaceAll(href, lute.LinkBase, "")
		}
		node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: []byte(href)})
		linkTitle := lute.domAttrValue(n, "title")
		if "" != linkTitle {
			node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
			node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: []byte(linkTitle)})
		}
		node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
	case atom.Del, atom.S, atom.Strike:
		marker := lute.domAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1CloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2CloseMarker, Tokens: []byte(marker)})
		}
	case atom.Details:
		tree.context.tip.AppendChild(&ast.Node{Type: ast.NodeHTMLBlock, Tokens: []byte("</details>")})
	}
}

func (context *Context) parentTip(n *html.Node) {
	if tip := context.tip.Parent; nil != tip {
		context.tip = context.tip.Parent
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
