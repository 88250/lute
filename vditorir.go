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
	"strconv"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/util"

	"github.com/88250/lute/html"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

// SpinVditorIRDOM 自旋 Vditor Instant-Rendering DOM，用于即时渲染模式下的编辑。
func (lute *Lute) SpinVditorIRDOM(ivHTML string) (ovHTML string) {
	lute.VditorIR = true
	lute.VditorWYSIWYG = true

	// 替换插入符
	ivHTML = strings.ReplaceAll(ivHTML, "<wbr>", parse.Caret)
	markdown := lute.vditorIRDOM2Md(ivHTML)
	tree := parse.Parse("", []byte(markdown), lute.Options)
	renderer := render.NewVditorIRRenderer(tree)
	output := renderer.Render()
	if renderer.Option.Footnotes && 0 < len(renderer.Tree.Context.FootnotesDefs) {
		output = renderer.RenderFootnotesDefs(renderer.Tree.Context)
	}
	// 替换插入符
	ovHTML = strings.ReplaceAll(string(output), parse.Caret, "<wbr>")
	return
}

// HTML2VditorIRDOM 将 HTML 转换为 Vditor Instant-Rendering DOM，用于即时渲染模式下粘贴。
func (lute *Lute) HTML2VditorIRDOM(sHTML string) (vHTML string) {
	lute.VditorIR = true
	lute.VditorWYSIWYG = true

	markdown, err := lute.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), lute.Options)
	renderer := render.NewVditorIRRenderer(tree)
	for nodeType, rendererFunc := range lute.HTML2VditorIRDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	if renderer.Option.Footnotes && 0 < len(renderer.Tree.Context.FootnotesDefs) {
		output = renderer.RenderFootnotesDefs(renderer.Tree.Context)
	}
	vHTML = string(output)
	return
}

// VditorIRDOM2HTML 将 Vditor Instant-Rendering DOM 转换为 HTML，用于 Vditor.getHTML() 接口。
func (lute *Lute) VditorIRDOM2HTML(vhtml string) (sHTML string) {
	lute.VditorIR = true
	lute.VditorWYSIWYG = true

	markdown := lute.vditorIRDOM2Md(vhtml)
	sHTML = lute.Md2HTML(markdown)
	return
}

// Md2VditorIRDOM 将 markdown 转换为 Vditor Instant-Rendering DOM，用于从源码模式切换至所见即所得模式。
func (lute *Lute) Md2VditorIRDOM(markdown string) (vHTML string) {
	lute.VditorIR = true
	lute.VditorWYSIWYG = true

	tree := parse.Parse("", []byte(markdown), lute.Options)
	renderer := render.NewVditorIRRenderer(tree)
	output := renderer.Render()
	if renderer.Option.Footnotes && 0 < len(renderer.Tree.Context.FootnotesDefs) {
		output = renderer.RenderFootnotesDefs(renderer.Tree.Context)
	}
	vHTML = string(output)
	return
}

// VditorIRDOM2Md 将 Vditor Instant-Rendering DOM 转换为 markdown，用于从所见即所得模式切换至源码模式。
func (lute *Lute) VditorIRDOM2Md(htmlStr string) (markdown string) {
	lute.VditorIR = true
	lute.VditorWYSIWYG = true

	htmlStr = strings.ReplaceAll(htmlStr, parse.Zwsp, "")
	markdown = lute.vditorIRDOM2Md(htmlStr)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	return
}

func (lute *Lute) vditorIRDOM2Md(htmlStr string) (markdown string) {
	// 删掉插入符
	htmlStr = strings.ReplaceAll(htmlStr, "<wbr>", "")

	// 替换结尾空白，否则 HTML 解析会产生冗余节点导致生成空的代码块
	htmlStr = strings.ReplaceAll(htmlStr, "\t\n", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "    \n", "  \n")

	// 将字符串解析为 DOM 树

	reader := strings.NewReader(htmlStr)
	htmlRoot := &html.Node{Type: html.ElementNode}
	htmlNodes, err := html.ParseFragment(reader, htmlRoot)
	if nil != err {
		markdown = err.Error()
		return
	}

	// 调整 DOM 结构
	lute.adjustVditorDOM(htmlNodes)

	// 将 HTML 树转换为 Markdown AST

	tree := &parse.Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument}, Context: &parse.Context{Option: lute.Options}}
	tree.Context.Tip = tree.Root
	for _, htmlNode := range htmlNodes {
		lute.genASTByVditorIRDOM(htmlNode, tree)
	}

	// 调整树结构

	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n.Type {
			case ast.NodeInlineHTML, ast.NodeCodeSpan, ast.NodeInlineMath, ast.NodeHTMLBlock, ast.NodeCodeBlockCode, ast.NodeMathBlockContent:
				n.Tokens = util.UnescapeHTML(n.Tokens)
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

	renderer := render.NewFormatRenderer(tree)
	formatted := renderer.Render()
	markdown = string(formatted)
	return
}

// genASTByVditorIRDOM 根据指定的 Vditor IR DOM 节点 n 进行深度优先遍历并逐步生成 Markdown 语法树 tree。
func (lute *Lute) genASTByVditorIRDOM(n *html.Node, tree *parse.Tree) {
	dataRender := lute.domAttrValue(n, "data-render")
	if "1" == dataRender || "2" == dataRender { // 1：浮动工具栏，2：preview 代码块、数学公式块
		return
	}

	dataType := lute.domAttrValue(n, "data-type")

	if atom.Div == n.DataAtom {
		if "code-block" == dataType || "html-block" == dataType || "math-block" == dataType {
			if ("code-block" == dataType || "math-block" == dataType) &&
				!strings.Contains(lute.domAttrValue(n.FirstChild, "data-type"), "-block-open-marker") {
				// 处理在结尾 ``` 或者 $$ 后换行的情况
				p := &ast.Node{Type: ast.NodeParagraph}
				text := &ast.Node{Type: ast.NodeText, Tokens: []byte(lute.domText(n.FirstChild))}
				p.AppendChild(text)
				tree.Context.Tip.AppendChild(p)
				tree.Context.Tip = p
				return
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				lute.genASTByVditorIRDOM(c, tree)
			}
		} else if "link-ref-defs-block" == dataType {
			text := lute.domText(n)
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
			node := &ast.Node{Type: ast.NodeText, Tokens: []byte("[toc]\n\n")}
			tree.Context.Tip.AppendChild(node)
		} else {
			text := lute.domText(n)
			if parse.Caret+"\n" == text { // 处理 FireFox 某些情况下产生的分段
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Caret + "\n")})
			}
		}
		return
	}

	class := lute.domAttrValue(n, "class")
	content := strings.ReplaceAll(n.Data, parse.Zwsp, "")
	node := &ast.Node{Type: ast.NodeText, Tokens: []byte(content)}
	switch n.DataAtom {
	case 0:
		if "" == content {
			return
		}

		checkIndentCodeBlock := strings.ReplaceAll(content, parse.Caret, "")
		checkIndentCodeBlock = strings.ReplaceAll(checkIndentCodeBlock, "\t", "    ")
		if (!lute.isInline(n.PrevSibling)) && strings.HasPrefix(checkIndentCodeBlock, "    ") {
			node.Type = ast.NodeCodeBlock
			node.IsFencedCodeBlock = true
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: []byte("```"), CodeBlockFenceLen: 3})
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker})
			startCaret := strings.HasPrefix(content, parse.Caret)
			if startCaret {
				content = strings.ReplaceAll(content, parse.Caret, "")
			}
			content = strings.TrimSpace(content)
			if startCaret {
				content = parse.Caret + content
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
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		if "" == strings.TrimSpace(lute.domText(n)) {
			return
		}
		node.Type = ast.NodeHeading
		marker := lute.domAttrValue(n, "data-marker")
		id := lute.domAttrValue(n, "data-id")
		if "" != id {
			node.HeadingID = []byte(id)
		}
		node.HeadingSetext = "=" == marker || "-" == marker
		if !node.HeadingSetext {
			marker := lute.domText(n.FirstChild)
			level := bytes.Count([]byte(marker), []byte("#"))
			node.HeadingLevel = level
			//headingC8hMarker := &ast.Node{Type: ast.NodeHeadingC8hMarker}
			//headingC8hMarker.Tokens = []byte(marker)
			//node.AppendChild(headingC8hMarker)
		} else {
			// 将 Setext 强制转为 ATX
			node.HeadingSetext = false
			if "=" == marker {
				node.HeadingLevel = 1
			} else {
				node.HeadingLevel = 2
			}
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Hr:
		node.Type = ast.NodeThematicBreak
		tree.Context.Tip.AppendChild(node)
	case atom.Blockquote:
		content := strings.TrimSpace(lute.domText(n))
		if "" == content || "&gt;" == content || parse.Caret == content {
			return
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
		tight := lute.domAttrValue(n, "data-tight")
		if "true" == tight || "" == tight {
			node.Tight = true
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Li:
		if p := n.FirstChild; nil != p && atom.P == p.DataAtom && nil != p.NextSibling && atom.P == p.NextSibling.DataAtom {
			tree.Context.Tip.Tight = false
		}

		node.Type = ast.NodeListItem
		marker := lute.domAttrValue(n, "data-marker")
		var bullet byte
		if "" == marker {
			if nil != n.Parent && atom.Ol == n.Parent.DataAtom {
				firstLiMarker := lute.domAttrValue(n.Parent.FirstChild, "data-marker")
				if startAttr := lute.domAttrValue(n.Parent, "start"); "" == startAttr {
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
					if "1." != marker && "1)" != marker && nil != n.PrevSibling && atom.Li != n.PrevSibling.DataAtom &&
						nil != n.Parent.Parent && (atom.Ol == n.Parent.Parent.DataAtom || atom.Ul == n.Parent.Parent.DataAtom) {
						// 子有序列表第一项必须从 1 开始
						marker = "1."
					}
					if "1." != marker && "1)" != marker && atom.Ol == n.Parent.DataAtom && n.Parent.FirstChild == n && "" == lute.domAttrValue(n.Parent, "start") {
						marker = "1."
					}
				} else {
					if "*" != marker && "-" != marker && "+" != marker {
						marker = "*"
					}
					bullet = marker[0]
				}
			} else {
				marker = lute.domAttrValue(n, "data-marker")
				if "" == marker {
					marker = "*"
				}
				bullet = marker[0]
			}
		}
		node.ListData = &ast.ListData{Marker: []byte(marker), BulletChar: bullet}
		if 0 == bullet {
			node.ListData.Num, _ = strconv.Atoi(string(marker[0]))
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

			divDataType := lute.domAttrValue(n.Parent, "data-type")
			switch divDataType {
			case "math-block":
				node.Type = ast.NodeMathBlock
				node.AppendChild(&ast.Node{Type: ast.NodeMathBlockContent, Tokens: codeTokens})
				tree.Context.Tip.AppendChild(node)
			case "html-block":
				node.Type = ast.NodeHTMLBlock
				node.Tokens = codeTokens
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
		if lute.starstWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if parse.Caret == text {
			node.Tokens = []byte(parse.Caret)
			tree.Context.Tip.AppendChild(node)
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
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && parse.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Caret)})
			if "_" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeEmU8eCloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, parse.Zwsp, "")

		// 开头结尾空格后会形成 * foo * 导致强调、加粗删除线标记失效，这里将空格移到右标记符前后 _*foo*_
		if strings.HasPrefix(n.FirstChild.Data, " ") && nil == n.FirstChild.PrevSibling {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
		}
		if strings.HasSuffix(n.FirstChild.Data, " ") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
		}
		if strings.HasSuffix(n.FirstChild.Data, "\n") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, "\n")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: "\n"})
		}

		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Strong, atom.B:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.starstWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if parse.Caret == text {
			node.Tokens = []byte(parse.Caret)
			tree.Context.Tip.AppendChild(node)
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
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && parse.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Caret)})
			if "__" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeStrongU8eCloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kCloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, parse.Zwsp, "")
		if strings.HasPrefix(n.FirstChild.Data, " ") && nil == n.FirstChild.PrevSibling {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
		}
		if strings.HasSuffix(n.FirstChild.Data, " ") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
		}
		if strings.HasSuffix(n.FirstChild.Data, "\n") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, "\n")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: "\n"})
		}

		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Del, atom.S, atom.Strike:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.starstWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if parse.Caret == text {
			node.Tokens = []byte(parse.Caret)
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeStrikethrough
		marker := lute.domAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1OpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2OpenMarker, Tokens: []byte(marker)})
		}
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && parse.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Caret)})
			if "~" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1CloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2CloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, parse.Zwsp, "")
		if strings.HasPrefix(n.FirstChild.Data, " ") && nil == n.FirstChild.PrevSibling {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
			node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
		}
		if strings.HasSuffix(n.FirstChild.Data, " ") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
		}
		if strings.HasSuffix(n.FirstChild.Data, "\n") && nil == n.FirstChild.NextSibling {
			n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, "\n")
			n.InsertAfter(&html.Node{Type: html.TextNode, Data: "\n"})
		}

		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Code:
		if nil == n.FirstChild {
			return
		}
		contentStr := strings.ReplaceAll(n.FirstChild.Data, parse.Zwsp, "")
		if parse.Caret == contentStr {
			node.Tokens = []byte(parse.Caret)
			tree.Context.Tip.AppendChild(node)
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
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Br:
		if nil != n.Parent {
			if lute.parentIs(n, atom.Td, atom.Th) {
				if (nil == n.PrevSibling || parse.Caret == n.PrevSibling.Data) && (nil == n.NextSibling || parse.Caret == n.NextSibling.Data) {
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
					tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp)})
					return
				}
				if nil != n.Parent.Parent && nil != n.Parent.Parent.NextSibling && atom.Li == n.Parent.Parent.NextSibling.DataAtom {
					tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp)})
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
		imgAlt := lute.domAttrValue(n, "alt")
		if "emoji" == imgClass {
			node.Type = ast.NodeEmoji
			emojiImg := &ast.Node{Type: ast.NodeEmojiImg, Tokens: tree.EmojiImgTokens(imgAlt, lute.domAttrValue(n, "src"))}
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
		if lute.hasAttr(n, "checked") {
			node.TaskListItemChecked = true
		}
		tree.Context.Tip.AppendChild(node)
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
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Sup:
		if nil == n.FirstChild {
			break
		}
		if "footnotes-ref" == dataType {
			node.Type = ast.NodeText
			node.Tokens = []byte(lute.domText(n))
			tree.Context.Tip.AppendChild(node)
		}
		return
	case atom.Span:
		switch dataType {
		case "inline-node", "em", "strong", "s", "a", "link-ref", "img", "code":
			node.Type = ast.NodeText
			node.Tokens = []byte(lute.domText(n))
			tree.Context.Tip.AppendChild(node)
			return
		case "math-block-close-marker":
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker, Tokens: []byte("$$")})
			defer tree.Context.ParentTip()
			return
		case "math-block-open-marker":
			node.Type = ast.NodeMathBlock
			node.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker, Tokens: []byte("$$")})
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			return
		case "code-block-open-marker":
			if atom.Pre == n.NextSibling.DataAtom { // DOM 后缺少 info span 节点
				n.InsertAfter(&html.Node{DataAtom: atom.Span, Attr: []html.Attribute{{Key: "data-type", Val: "code-block-info"}}})
			}
			marker := []byte(lute.domText(n))
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
			info := []byte(lute.domText(n))
			info = bytes.ReplaceAll(info, []byte(parse.Zwsp), nil)
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker, CodeBlockInfo: info})
			return
		case "code-block-close-marker":
			marker := []byte(lute.domText(n))
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
			if caretInMarker := strings.Contains(lute.domText(n), parse.Caret); caretInMarker {
				n.InsertAfter(&html.Node{Type:html.TextNode, Data:parse.Caret})
			}
			return
		}

		if nil == n.FirstChild {
			break
		}

		var codeTokens []byte
		if parse.Zwsp == n.FirstChild.Data && "" == lute.domAttrValue(n, "style") && nil != n.FirstChild.NextSibling {
			codeTokens = []byte(n.FirstChild.NextSibling.FirstChild.Data)
		} else if atom.Code == n.FirstChild.DataAtom {
			codeTokens = []byte(n.FirstChild.FirstChild.Data)
			if parse.Zwsp == string(codeTokens) {
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
		}
		return
	case atom.Font:
		return
	case atom.Details:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = lute.domHTML(n)
		node.Tokens = bytes.SplitAfter(node.Tokens, []byte("</summary>"))[0]
		tree.Context.Tip.AppendChild(node)
	case atom.Kbd:
		node.Type = ast.NodeInlineHTML
		node.Tokens = lute.domHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Summary:
		return
	default:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = lute.domHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTByVditorIRDOM(c, tree)
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
		tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeHTMLBlock, Tokens: []byte("</details>")})
	}
}
