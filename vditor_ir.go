// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"strconv"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

// SpinVditorIRDOM 自旋 Vditor Instant-Rendering DOM，用于即时渲染模式下的编辑。
func (lute *Lute) SpinVditorIRDOM(ivHTML string) (ovHTML string) {
	// 替换插入符
	ivHTML = strings.ReplaceAll(ivHTML, "<wbr>", editor.Caret)
	markdown := lute.vditorIRDOM2Md(ivHTML)
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorIRRenderer(tree, lute.RenderOptions)
	output := renderer.Render()
	// 替换插入符
	ovHTML = strings.ReplaceAll(string(output), editor.Caret, "<wbr>")
	return
}

// HTML2VditorIRDOM 将 HTML 转换为 Vditor Instant-Rendering DOM，用于即时渲染模式下粘贴。
func (lute *Lute) HTML2VditorIRDOM(sHTML string) (vHTML string) {
	markdown, err := lute.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorIRRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.HTML2VditorIRDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

// VditorIRDOM2HTML 将 Vditor Instant-Rendering DOM 转换为 HTML，用于 Vditor.getHTML() 接口。
func (lute *Lute) VditorIRDOM2HTML(vhtml string) (sHTML string) {
	markdown := lute.vditorIRDOM2Md(vhtml)
	sHTML = lute.Md2HTML(markdown)
	return
}

// Md2VditorIRDOM 将 markdown 转换为 Vditor Instant-Rendering DOM，用于从源码模式切换至即时渲染模式。
func (lute *Lute) Md2VditorIRDOM(markdown string) (vHTML string) {
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorIRRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2VditorIRDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

// VditorIRDOM2Md 将 Vditor Instant-Rendering DOM 转换为 markdown，用于从即时渲染模式切换至源码模式。
func (lute *Lute) VditorIRDOM2Md(htmlStr string) (markdown string) {
	htmlStr = strings.ReplaceAll(htmlStr, editor.Zwsp, "")
	markdown = lute.vditorIRDOM2Md(htmlStr)
	markdown = strings.ReplaceAll(markdown, editor.Zwsp, "")
	return
}

func (lute *Lute) vditorIRDOM2Md(htmlStr string) (markdown string) {
	// 删掉插入符
	htmlStr = strings.ReplaceAll(htmlStr, "<wbr>", "")

	// 替换结尾空白，否则 HTML 解析会产生冗余节点导致生成空的代码块
	htmlStr = strings.ReplaceAll(htmlStr, "\t\n", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "    \n", "  \n")

	// 将字符串解析为 DOM 树
	htmlRoot := lute.parseHTML(htmlStr)
	if nil == htmlRoot {
		return
	}

	// 调整 DOM 结构
	lute.adjustVditorDOM(htmlRoot)

	// 将 HTML 树转换为 Markdown AST
	tree := &parse.Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument}, Context: &parse.Context{ParseOption: lute.ParseOptions}}
	tree.Context.Tip = tree.Root
	for c := htmlRoot.FirstChild; nil != c; c = c.NextSibling {
		lute.genASTByVditorIRDOM(c, tree)
	}

	// 调整树结构
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n.Type {
			case ast.NodeInlineHTML, ast.NodeCodeSpan, ast.NodeInlineMath, ast.NodeHTMLBlock, ast.NodeCodeBlockCode, ast.NodeMathBlockContent:
				n.Tokens = html.UnescapeHTML(n.Tokens)
				if nil != n.Next && ast.NodeCodeSpan == n.Next.Type && n.CodeMarkerLen == n.Next.CodeMarkerLen {
					// 合并代码节点 https://github.com/Vanessa219/vditor/issues/167
					n.FirstChild.Next.Tokens = append(n.FirstChild.Next.Tokens, n.Next.FirstChild.Next.Tokens...)
					n.Next.Unlink()
				}
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
	options := render.NewOptions()
	options.AutoSpace = false
	options.FixTermTypo = false
	renderer := render.NewFormatRenderer(tree, options)
	formatted := renderer.Render()
	markdown = string(formatted)
	return
}

// genASTByVditorIRDOM 根据指定的 Vditor IR DOM 节点 n 进行深度优先遍历并逐步生成 Markdown 语法树 tree。
func (lute *Lute) genASTByVditorIRDOM(n *html.Node, tree *parse.Tree) {
	dataRender := util.DomAttrValue(n, "data-render")
	if "1" == dataRender || "2" == dataRender { // 1：浮动工具栏，2：preview 代码块、数学公式块或者不解析的节点
		return
	}

	dataType := util.DomAttrValue(n, "data-type")

	if atom.Div == n.DataAtom {
		if "code-block" == dataType || "html-block" == dataType || "math-block" == dataType || "yaml-front-matter" == dataType {
			if ("code-block" == dataType || "math-block" == dataType) &&
				!strings.Contains(util.DomAttrValue(n.FirstChild, "data-type"), "-block-open-marker") {
				// 处理在结尾 ``` 或者 $$ 后换行的情况
				// TODO: 插入符现在已经不可能出现在该位置，确认后移除该段代码
				p := &ast.Node{Type: ast.NodeParagraph}
				text := &ast.Node{Type: ast.NodeText, Tokens: []byte(util.DomText(n.FirstChild))}
				p.AppendChild(text)
				tree.Context.Tip.AppendChild(p)
				tree.Context.Tip = p
				return
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				lute.genASTByVditorIRDOM(c, tree)
			}
		} else if "link-ref-defs-block" == dataType {
			text := util.DomText(n)
			node := &ast.Node{Type: ast.NodeText, Tokens: []byte(text)}
			tree.Context.Tip.AppendChild(node)
		} else if "footnotes-def" == dataType {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				lute.genASTByVditorIRDOM(c, tree)
			}
		} else if "footnotes-block" == dataType {
			for def := n.FirstChild; nil != def; def = def.NextSibling {
				originalHTML := &bytes.Buffer{}
				if err := html.Render(originalHTML, def); nil == err {
					md := lute.vditorIRDOM2Md(originalHTML.String())
					lines := strings.Split(md, "\n")
					md = ""
					for i, line := range lines {
						if 0 < i {
							md += "    " + line
						} else {
							md = line
						}
						md += "\n"
					}
					node := &ast.Node{Type: ast.NodeText, Tokens: []byte(md)}
					tree.Context.Tip.AppendChild(node)
				}
			}
		} else if "toc-block" == dataType {
			node := &ast.Node{Type: ast.NodeToC}
			tree.Context.Tip.AppendChild(node)
		} else {
			text := util.DomText(n)
			if editor.Caret+"\n" == text { // 处理 FireFox 某些情况下产生的分段
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(editor.Caret + "\n")})
			}
		}
		return
	}

	class := util.DomAttrValue(n, "class")
	content := strings.ReplaceAll(n.Data, editor.Zwsp, "")
	node := &ast.Node{Type: ast.NodeText, Tokens: []byte(content)}
	switch n.DataAtom {
	case 0:
		if "" == content {
			return
		}

		checkIndentCodeBlock := strings.ReplaceAll(content, editor.Caret, "")
		checkIndentCodeBlock = strings.ReplaceAll(checkIndentCodeBlock, "\t", "    ")
		if (!lute.isInline(n.PrevSibling)) && strings.HasPrefix(checkIndentCodeBlock, "    ") {
			node.Type = ast.NodeCodeBlock
			node.IsFencedCodeBlock = true
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: []byte("```"), CodeBlockFenceLen: 3})
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker})
			startCaret := strings.HasPrefix(content, editor.Caret)
			if startCaret {
				content = strings.ReplaceAll(content, editor.Caret, "")
			}
			content = strings.TrimSpace(content)
			if startCaret {
				content = editor.Caret + content
			}
			content := &ast.Node{Type: ast.NodeCodeBlockCode, Tokens: []byte(content)}
			node.AppendChild(content)
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: []byte("```"), CodeBlockFenceLen: 3})
			tree.Context.Tip.AppendChild(node)
			return
		}
		if nil != n.Parent && atom.A == n.Parent.DataAtom {
			node.Type = ast.NodeLinkText
		}
		tree.Context.Tip.AppendChild(node)
	case atom.P, atom.Div:
		node.Type = ast.NodeParagraph
		text := util.DomText(n)
		if "\n" == text && ast.NodeBlockquote == tree.Context.Tip.Type && nil == tree.Context.Tip.FirstChild.Next {
			// 不允许在 bq 第一个节点前换行
			return
		} else {
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		}
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		text := util.DomText(n)
		if "" == strings.TrimSpace(text) {
			return
		}
		node.Type = ast.NodeHeading
		marker := util.DomAttrValue(n, "data-marker")
		node.HeadingSetext = "=" == marker || "-" == marker
		if !node.HeadingSetext {
			marker := util.DomText(n.FirstChild)
			node.HeadingLevel = bytes.Count([]byte(marker), []byte("#"))
		} else {
			if "" == strings.TrimSpace(strings.ReplaceAll(util.DomText(n.LastChild), editor.Caret, "")) {
				node.Type = ast.NodeText
				node.Tokens = []byte(text)
				tree.Context.Tip.AppendChild(node)
				tree.Context.Tip = node
				break
			}

			if "=" == marker {
				node.HeadingLevel = 1
			} else {
				node.HeadingLevel = 2
			}
			if nil != n.LastChild.PrevSibling {
				n.LastChild.PrevSibling.Data = strings.TrimSuffix(n.LastChild.PrevSibling.Data, "\n")
			}
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Hr:
		node.Type = ast.NodeThematicBreak
		tree.Context.Tip.AppendChild(node)
	case atom.Blockquote:
		content := strings.TrimSpace(util.DomText(n))
		if "" == content || "&gt;" == content {
			return
		}
		if editor.Caret == content {
			node.Type = ast.NodeText
			node.Tokens = []byte(content)
			tree.Context.Tip.AppendChild(node)
		}

		node.Type = ast.NodeBlockquote
		node.AppendChild(&ast.Node{Type: ast.NodeBlockquoteMarker, Tokens: []byte(">")})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Ol, atom.Ul:
		if nil == n.FirstChild {
			return
		}

		node.Type = ast.NodeList
		node.ListData = &ast.ListData{}
		if atom.Ol == n.DataAtom {
			node.ListData.Typ = 1
		}
		tight := util.DomAttrValue(n, "data-tight")
		if "true" == tight || "" == tight {
			node.ListData.Tight = true
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Li:
		if p := n.FirstChild; nil != p && atom.P == p.DataAtom && nil != p.NextSibling && atom.P == p.NextSibling.DataAtom {
			tree.Context.Tip.ListData.Tight = false
		}

		node.Type = ast.NodeListItem
		marker := util.DomAttrValue(n, "data-marker")
		var bullet byte
		if "" == marker {
			if nil != n.Parent && atom.Ol == n.Parent.DataAtom {
				firstLiMarker := util.DomAttrValue(n.Parent.FirstChild, "data-marker")
				if startAttr := util.DomAttrValue(n.Parent, "start"); "" == startAttr {
					marker = "1"
				} else {
					marker = startAttr
				}
				if "" != firstLiMarker {
					marker += firstLiMarker[len(firstLiMarker)-1:]
				} else {
					marker += "."
				}
			} else {
				marker = util.DomAttrValue(n.Parent, "data-marker")
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
					if "1." != marker && "1)" != marker && nil != n.PrevSibling && atom.Li != n.PrevSibling.DataAtom &&
						nil != n.Parent.Parent && (atom.Ol == n.Parent.Parent.DataAtom || atom.Ul == n.Parent.Parent.DataAtom) {
						// 子有序列表第一项必须从 1 开始
						marker = "1."
					}
					if "1." != marker && "1)" != marker && atom.Ol == n.Parent.DataAtom && n.Parent.FirstChild == n && "" == util.DomAttrValue(n.Parent, "start") {
						marker = "1."
					}
				} else {
					if "*" != marker && "-" != marker && "+" != marker {
						marker = "*"
					}
					bullet = marker[0]
				}
			} else {
				marker = util.DomAttrValue(n, "data-marker")
				if "" == marker {
					marker = "*"
				}
				bullet = marker[0]
			}
		}
		node.ListData = &ast.ListData{Marker: []byte(marker), BulletChar: bullet}
		if 0 == bullet {
			node.ListData.Num, _ = strconv.Atoi(marker[:len(marker)-1])
			node.ListData.Delimiter = marker[len(marker)-1]
		}

		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Pre:
		if atom.Code == n.FirstChild.DataAtom {
			var codeTokens []byte
			if nil != n.FirstChild.FirstChild {
				codeTokens = []byte(n.FirstChild.FirstChild.Data)
			}

			divDataType := util.DomAttrValue(n.Parent, "data-type")
			switch divDataType {
			case "math-block":
				node.Type = ast.NodeMathBlockContent
				node.Tokens = codeTokens
				tree.Context.Tip.AppendChild(node)
			case "html-block":
				node.Type = ast.NodeHTMLBlock
				node.Tokens = codeTokens
				tree.Context.Tip.AppendChild(node)
			case "yaml-front-matter":
				node.Type = ast.NodeYamlFrontMatter
				node.AppendChild(&ast.Node{Type: ast.NodeYamlFrontMatterContent, Tokens: codeTokens})
				tree.Context.Tip.AppendChild(node)
			default:
				node.Type = ast.NodeCodeBlockCode
				node.Tokens = codeTokens
				tree.Context.Tip.AppendChild(node)
			}
		}
		return
	case atom.Em, atom.I:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, editor.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(editor.Zwsp + "\n")})
		}
		text := strings.TrimSpace(util.DomText(n))
		if lute.isEmptyText(n) {
			return
		}
		if editor.Caret == text {
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeEmphasis
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Strong, atom.B:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, editor.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(editor.Zwsp + "\n")})
		}
		text := strings.TrimSpace(util.DomText(n))
		if lute.isEmptyText(n) {
			return
		}
		if editor.Caret == text {
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeStrong
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Del, atom.S, atom.Strike:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, editor.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(editor.Zwsp + "\n")})
		}
		text := strings.TrimSpace(util.DomText(n))
		if lute.isEmptyText(n) {
			return
		}
		if editor.Caret == text {
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeStrikethrough
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Mark:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, editor.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(editor.Zwsp + "\n")})
		}
		text := strings.TrimSpace(util.DomText(n))
		if lute.isEmptyText(n) {
			return
		}
		if editor.Caret == text {
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeMark
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Code:
		if nil == n.FirstChild {
			return
		}
		contentStr := strings.ReplaceAll(n.FirstChild.Data, editor.Zwsp, "")
		if editor.Caret == contentStr {
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		if "" == contentStr {
			return
		}
		codeTokens := []byte(contentStr)
		content := &ast.Node{Type: ast.NodeCodeSpanContent, Tokens: codeTokens}
		node.Type = ast.NodeCodeSpan
		node.AppendChild(content)
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Br:
		if nil != n.Parent {
			if lute.parentIs(n, atom.Td, atom.Th) {
				if (nil == n.PrevSibling || editor.Caret == n.PrevSibling.Data) && (nil == n.NextSibling || editor.Caret == n.NextSibling.Data) {
					return
				}
				if nil == n.NextSibling {
					return // 删掉表格中结尾的 br
				}

				node.Type = ast.NodeInlineHTML
				node.Tokens = []byte("<br />")
				tree.Context.Tip.AppendChild(node)
				return
			}
			if atom.P == n.Parent.DataAtom {
				if nil != n.Parent.NextSibling && (atom.Ul == n.Parent.NextSibling.DataAtom || atom.Ol == n.Parent.NextSibling.DataAtom || atom.Blockquote == n.Parent.NextSibling.DataAtom) {
					tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(editor.Zwsp)})
					return
				}
			}
		}

		node.Type = ast.NodeHardBreak
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.A:
		node.Type = ast.NodeLink
		node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Img:
		imgClass := class
		imgAlt := util.DomAttrValue(n, "alt")
		if "emoji" == imgClass {
			node.Type = ast.NodeEmoji
			emojiImg := &ast.Node{Type: ast.NodeEmojiImg, Tokens: tree.EmojiImgTokens(imgAlt, util.DomAttrValue(n, "src"))}
			emojiImg.AppendChild(&ast.Node{Type: ast.NodeEmojiAlias, Tokens: []byte(":" + imgAlt + ":")})
			node.AppendChild(emojiImg)
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		} else {
			return
		}
	case atom.Input:
		if nil == n.Parent || nil == n.Parent.Parent || (atom.P != n.Parent.DataAtom && atom.Li != n.Parent.DataAtom) {
			// 仅允许 input 出现在任务列表中
			return
		}
		if nil != n.NextSibling && atom.Span == n.NextSibling.DataAtom {
			// 在任务列表前退格
			n.NextSibling.FirstChild.Data = strings.TrimSpace(n.NextSibling.FirstChild.Data)
			break
		}
		node.Type = ast.NodeTaskListItemMarker
		node.TaskListItemChecked = lute.hasAttr(n, "checked")
		tree.Context.Tip.AppendChild(node)
		if nil != node.Parent.Parent.Parent && nil != node.Parent.Parent.Parent.ListData { // ul.li.p.input
			node.Parent.Parent.Parent.ListData.Typ = 3
		}
	case atom.Table:
		node.Type = ast.NodeTable
		var tableAligns []int
		if nil == n.FirstChild || nil == n.FirstChild.FirstChild || nil == n.FirstChild.FirstChild.FirstChild {
			return
		}

		for th := n.FirstChild.FirstChild.FirstChild; nil != th; th = th.NextSibling {
			align := util.DomAttrValue(th, "align")
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
		node.Tokens = nil
		tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeParagraph}) // 表格开头输入会导致解析问题，所以插入一个空段落进行分隔
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Thead:
		node.Type = ast.NodeTableHead
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Tbody:
	case atom.Tr:
		node.Type = ast.NodeTableRow
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Th, atom.Td:
		node.Type = ast.NodeTableCell
		align := util.DomAttrValue(n, "align")
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
		node.Tokens = nil
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Sup:
		if nil == n.FirstChild {
			break
		}
		if "footnotes-ref" == dataType {
			node.Type = ast.NodeText
			node.Tokens = []byte(util.DomText(n))
			tree.Context.Tip.AppendChild(node)
		}
		return
	case atom.Span:
		switch dataType {
		case "inline-node", "em", "strong", "s", "a", "link-ref", "img", "code", "heading-id", "html-inline", "inline-math", "html-entity":
			node.Type = ast.NodeText
			node.Tokens = []byte(util.DomText(n))
			tree.Context.Tip.AppendChild(node)
			return
		case "math-block-close-marker":
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker, Tokens: parse.MathBlockMarker})
			defer tree.Context.ParentTip()
			return
		case "math-block-open-marker":
			node.Type = ast.NodeMathBlock
			node.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker, Tokens: parse.MathBlockMarker})
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			return
		case "yaml-front-matter-close-marker":
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeYamlFrontMatterCloseMarker, Tokens: parse.YamlFrontMatterMarker})
			defer tree.Context.ParentTip()
			return
		case "yaml-front-matter-open-marker":
			node.Type = ast.NodeYamlFrontMatter
			node.AppendChild(&ast.Node{Type: ast.NodeYamlFrontMatterOpenMarker, Tokens: parse.YamlFrontMatterMarker})
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			return
		case "code-block-open-marker":
			if atom.Pre == n.NextSibling.DataAtom { // DOM 后缺少 info span 节点
				n.InsertAfter(&html.Node{DataAtom: atom.Span, Attr: []*html.Attribute{{Key: "data-type", Val: "code-block-info"}}})
			}
			marker := []byte(util.DomText(n))
			lastBacktick := bytes.LastIndex(marker, []byte("`")) + 1
			if 0 < lastBacktick {
				// 把 ` 后面的字符调整到 info 节点
				n.NextSibling.AppendChild(&html.Node{Data: string(marker[lastBacktick:])})
				marker = marker[:lastBacktick]
			}
			node.Type = ast.NodeCodeBlock
			node.IsFencedCodeBlock = true
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: marker, CodeBlockFenceLen: len(marker)})
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			return
		case "code-block-info":
			info := []byte(util.DomText(n))
			info = bytes.ReplaceAll(info, []byte(editor.Zwsp), nil)
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker, CodeBlockInfo: info})
			return
		case "code-block-close-marker":
			marker := []byte(util.DomText(n))
			lastBacktick := bytes.LastIndex(marker, []byte("`")) + 1
			if 0 < lastBacktick {
				marker = marker[:lastBacktick]
			}
			if 0 == len(marker) {
				marker = []byte("```")
			}
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: marker, CodeBlockFenceLen: len(marker)})
			defer tree.Context.ParentTip()
			return
		case "heading-marker":
			text := util.DomText(n)
			if caretInMarker := strings.Contains(text, editor.Caret); caretInMarker {
				caret := &html.Node{Type: html.TextNode, Data: editor.Caret}
				n.InsertAfter(caret)
				text = strings.ReplaceAll(text, "#", "")
				text = strings.ReplaceAll(text, editor.Caret, "")
				text = strings.TrimSpace(text)
				if 0 < len(text) {
					caret.Data = text + caret.Data
				}
			}
			return
		}

		if nil == n.FirstChild {
			break
		}

		var codeTokens []byte
		if editor.Zwsp == n.FirstChild.Data && "" == util.DomAttrValue(n, "style") && nil != n.FirstChild.NextSibling {
			codeTokens = []byte(n.FirstChild.NextSibling.FirstChild.Data)
		} else if atom.Code == n.FirstChild.DataAtom {
			codeTokens = []byte(n.FirstChild.FirstChild.Data)
			if editor.Zwsp == string(codeTokens) {
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
			tree.Context.Tip.AppendChild(node)
		} else if "html-inline" == dataType {
			node.Type = ast.NodeInlineHTML
			node.Tokens = codeTokens
			tree.Context.Tip.AppendChild(node)
		} else if "code-inline" == dataType {
			node.Tokens = codeTokens
			tree.Context.Tip.AppendChild(node)
		} else if "html-entity" == dataType {
			node.Type = ast.NodeText
			node.Tokens = codeTokens
			tree.Context.Tip.AppendChild(node)
		}
		return
	case atom.Font:
		node.Type = ast.NodeText
		node.Tokens = []byte(util.DomText(n))
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Details:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = util.DomHTML(n)
		node.Tokens = bytes.SplitAfter(node.Tokens, []byte("</summary>"))[0]
		tree.Context.Tip.AppendChild(node)
	case atom.Kbd:
		// kbd 标签由 code 标签构成节点
	case atom.Summary:
		return
	default:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = util.DomHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTByVditorIRDOM(c, tree)
	}

	switch n.DataAtom {
	case atom.A:
		node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
		node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
		href := util.DomAttrValue(n, "href")
		if "" != lute.RenderOptions.LinkBase {
			href = strings.ReplaceAll(href, lute.RenderOptions.LinkBase, "")
		}
		if "" != lute.RenderOptions.LinkPrefix {
			href = strings.ReplaceAll(href, lute.RenderOptions.LinkPrefix, "")
		}
		node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: []byte(href)})
		linkTitle := util.DomAttrValue(n, "title")
		if "" != linkTitle {
			node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
			node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: []byte(linkTitle)})
		}
		node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
	case atom.Details:
		tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeHTMLBlock, Tokens: []byte("</details>")})
	}
}
