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

// Md2HTML 将 markdown 转换为标准 HTML，用于源码模式预览。
func (lute *Lute) Md2HTML(markdown string) (sHTML string) {
	sHTML = lute.MarkdownStr("", markdown)
	return
}

// SpinVditorDOM 自旋 Vditor DOM，用于所见即所得模式下的编辑。
func (lute *Lute) SpinVditorDOM(ivHTML string) (ovHTML string) {
	ivHTML = strings.ReplaceAll(ivHTML, editor.FrontEndCaret, editor.Caret)
	markdown := lute.vditorDOM2Md(ivHTML)
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorRenderer(tree, lute.RenderOptions)
	output := renderer.Render()
	ovHTML = strings.ReplaceAll(string(output), editor.Caret, editor.FrontEndCaret)
	return
}

// HTML2VditorDOM 将 HTML 转换为 Vditor DOM，用于所见即所得模式下粘贴。
func (lute *Lute) HTML2VditorDOM(sHTML string) (vHTML string) {
	markdown, err := lute.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.HTML2VditorDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

// VditorDOM2HTML 将 Vditor DOM 转换为 HTML，用于 Vditor.getHTML() 接口。
func (lute *Lute) VditorDOM2HTML(vhtml string) (sHTML string) {
	markdown := lute.vditorDOM2Md(vhtml)
	sHTML = lute.Md2HTML(markdown)
	return
}

// Md2VditorDOM 将 markdown 转换为 Vditor DOM，用于从源码模式切换至所见即所得模式。
func (lute *Lute) Md2VditorDOM(markdown string) (vHTML string) {
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2VditorDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

// VditorDOM2Md 将 Vditor DOM 转换为 markdown，用于从所见即所得模式切换至源码模式。
func (lute *Lute) VditorDOM2Md(htmlStr string) (markdown string) {
	htmlStr = strings.ReplaceAll(htmlStr, editor.Zwsp, "")
	markdown = lute.vditorDOM2Md(htmlStr)
	markdown = strings.ReplaceAll(markdown, editor.Zwsp, "")
	return
}

// RenderEChartsJSON 用于渲染 ECharts JSON 格式数据。
func (lute *Lute) RenderEChartsJSON(markdown string) (json string) {
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewEChartsJSONRenderer(tree, lute.RenderOptions)
	output := renderer.Render()
	json = string(output)
	return
}

// RenderKityMinderJSON 用于渲染 KityMinder JSON 格式数据。
func (lute *Lute) RenderKityMinderJSON(markdown string) (json string) {
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewKityMinderJSONRenderer(tree, lute.RenderOptions)
	output := renderer.Render()
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
	htmlStr = strings.ReplaceAll(htmlStr, editor.FrontEndCaret, "")

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
		lute.genASTByVditorDOM(c, tree)
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

func (lute *Lute) parseHTML(htmlStr string) *html.Node {
	reader := strings.NewReader(htmlStr)
	doc, err := html.Parse(reader)
	if nil != err {
		return nil
	}
	if "html" != doc.FirstChild.Data {
		return doc
	}
	return doc.FirstChild.LastChild // doc.html.body
}

func (lute *Lute) adjustVditorDOM(root *html.Node) {
	lute.removeEmptyNodes(root)
	lute.removeHighlightJSSpans(root)

	for c := root.FirstChild; nil != c; c = c.NextSibling {
		lute.mergeVditorDOMList0(c)
	}

	for c := root.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustVditorDOMListTight0(c)
	}

	for c := root.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustVditorDOMListList(c)
	}

	for c := root.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustVditorDOMListItemInP(c)
	}

	for c := root.FirstChild; nil != c; {
		next := c.NextSibling
		lute.removeCodeCode(c)
		c = next
	}

	for c := root.FirstChild; nil != c; {
		next := c.NextSibling
		lute.adjustVditorDOMCodeA(c)
		c = next
	}

	for c := root.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustCustomTag(c)
	}

	for c := root.FirstChild; nil != c; c = c.NextSibling {
		lute.mergeSameStrong(c)
	}

	for c := root.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustTableCode(c)
	}

	for c := root.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustMath(c)
	}

	for c := root.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustNoscriptImg(c)
	}
}

func (lute *Lute) adjustCustomTag(n *html.Node) {
	// 将某些自定义标签转换为标准标签

	if html.ElementNode == n.Type && 0 == n.DataAtom {
		if "ucapcontent" == n.Data {
			n.DataAtom = atom.Div
		} else if "ucaptitle" == n.Data {
			n.DataAtom = atom.H2
			n.Data = "h2"
		} else if "markerow8" == n.Data {
			n.DataAtom = atom.Span
		} else if "app-document-text" == n.Data {
			n.DataAtom = atom.Div
		}
	}

	for c := n.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustCustomTag(c)
	}
}

func (lute *Lute) adjustNoscriptImg(n *html.Node) {
	if nil != n.Parent && atom.Figure == n.Parent.DataAtom &&
		atom.Noscript == n.DataAtom && nil != n.FirstChild && strings.HasPrefix(n.FirstChild.Data, "<img ") {
		img := n.FirstChild
		img.Unlink()
		img.DataAtom = atom.Img
		fragment, err := html.ParseFragment(strings.NewReader(img.Data), &html.Node{Type: html.ElementNode})
		if nil != err || 1 > len(fragment) {
			return
		}
		img = fragment[0]
		n.InsertBefore(img)
		var unlinks []*html.Node
		for c := n; nil != c; c = c.NextSibling {
			if atom.Figcaption == c.DataAtom {
				continue
			}
			unlinks = append(unlinks, c)
		}
		for _, unlink := range unlinks {
			unlink.Unlink()
		}
	}

	for c := n.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustNoscriptImg(c)
	}
}

func (lute *Lute) adjustMath(n *html.Node) {
	class := util.DomAttrValue(n, "class")
	if ((atom.Span == n.DataAtom || atom.Div == n.DataAtom) && strings.Contains(class, "mwe-math-element")) || (strings.Contains(class, "tex") && !strings.Contains(class, "text")) {
		if annos := util.DomChildrenByType(n, atom.Annotation); 0 < len(annos) {
			anno := annos[0]
			if "application/x-tex" == util.DomAttrValue(anno, "encoding") {
				util.SetDomAttrValue(n, "data-tex", util.DomText(anno))
				return
			}
		}

		if imgs := util.DomChildrenByType(n, atom.Img); 0 < len(imgs) {
			img := imgs[0]
			if mathContent := util.DomAttrValue(img, "alt"); "" != mathContent {
				util.SetDomAttrValue(n, "data-tex", mathContent)
				return
			}
		}
	}

	if atom.Img == n.DataAtom && strings.Contains(class, "ma-tex-img") {
		if mathContent := util.DomAttrValue(n, "alt"); "" != mathContent {
			n.DataAtom = atom.Span
			util.SetDomAttrValue(n, "data-tex", mathContent)
		}
		return
	}

	formula := util.DomAttrValue(n, "data-formula")
	if "" != formula {
		if html.ElementNode == n.Type && "mjx-container" == n.Data {
			n.DataAtom = atom.Span
		}

		util.SetDomAttrValue(n, "data-tex", formula)
		return
	}

	if strings.Contains(class, "texhtml") {
		if mathContent := util.DomTexhtml(n); "" != mathContent {
			util.SetDomAttrValue(n, "data-tex", mathContent)
			return
		}
	}

	if strings.Contains(class, "math") {
		scripts := util.DomChildrenByType(n, atom.Script)
		if 0 < len(scripts) {
			script := scripts[0]
			if "math/tex" == util.DomAttrValue(script, "type") {
				if mathContent := util.DomText(script); "" != mathContent {
					util.SetDomAttrValue(n, "data-tex", mathContent)
				}
			}
		}
	}

	for c := n.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustMath(c)
	}
}

func (lute *Lute) adjustTableCode(n *html.Node) {
	if atom.Table == n.DataAtom {
		// 表格类型的代码块进行预处理 https://github.com/siyuan-note/siyuan/issues/11540
		// 移除 <td class="gutter">
		// td class="code"> 下的 <div class="container"> 改为 <pre>

		tds := util.DomChildrenByType(n, atom.Td)
		var unlinks []*html.Node
		for _, td := range tds {
			tdClass := util.DomAttrValue(td, "class")
			if strings.Contains(tdClass, "gutter") {
				unlinks = append(unlinks, td)
				continue
			}

			// 移除 <span class="lnt"> 格式的行号 https://github.com/siyuan-note/siyuan/issues/13242
			spans := util.DomChildrenByType(td, atom.Span)
			removed := false
			for _, span := range spans {
				if "lnt" == util.DomAttrValue(span, "class") {
					removed = true
					break
				}
			}
			if removed {
				unlinks = append(unlinks, td)
				continue
			}

			if strings.Contains(tdClass, "code") {
				if c := td.FirstChild; nil != c && atom.Div == c.DataAtom {
					c.DataAtom = atom.Pre
					c.Data = "pre"
					lang := util.DomAttrValue(n, "class")
					lang = strings.ReplaceAll(lang, "syntaxhighlighter", "")
					lang = strings.TrimSpace(lang)
					if "" != lang {
						util.SetDomAttrValue(c, "class", lang)
					}
					for div := c.FirstChild; nil != div; div = div.NextSibling {
						div.DataAtom = atom.Code
						div.Data = "code"
					}
				}
			}
		}

		for _, unlink := range unlinks {
			unlink.Unlink()
		}
	}

	for c := n.FirstChild; nil != c; c = c.NextSibling {
		lute.adjustTableCode(c)
	}
}

func (lute *Lute) mergeSameStrong(n *html.Node) {
	for c := n.FirstChild; nil != c; {
		next := c.NextSibling
		if nil != next && atom.Strong == c.DataAtom && atom.Strong == next.DataAtom {
			for cc := next.FirstChild; nil != cc; {
				nextChild := cc.NextSibling
				cc.Unlink()
				c.AppendChild(cc)
				cc = nextChild
			}
			next.Unlink()
			next = c.NextSibling
		}
		c = next
	}

	for c := n.FirstChild; nil != c; c = c.NextSibling {
		lute.mergeSameStrong(c)
	}
}

// adjustVditorDOMListList 用于将 ul.ul 调整为 ul.li.ul。
func (lute *Lute) adjustVditorDOMListList(n *html.Node) {
	if atom.Ul != n.DataAtom && atom.Ol != n.DataAtom && atom.Li != n.DataAtom {
		return
	}

	if atom.Li == n.DataAtom {
		if nil != n.FirstChild && atom.Br == n.FirstChild.DataAtom {
			// 规范化换行时 li 的结构，对调 ZWSP 和 <br> 的位置
			n.FirstChild.DataAtom = 0
			n.FirstChild.Data = editor.Zwsp
			if nextLi := n.NextSibling; nil != n.NextSibling && atom.Li == n.NextSibling.DataAtom {
				if caret := nextLi.FirstChild; nil != caret && editor.Caret+editor.Zwsp == caret.Data {
					caret.Data = editor.Caret + "\n"
				}
			}
		}
	} else {
		if nil != n.Parent && (atom.Ul == n.Parent.DataAtom || atom.Ol == n.Parent.DataAtom) {
			if prevLi := n.PrevSibling; nil != prevLi {
				n.Unlink()
				prevLi.AppendChild(n)
			}
		}
	}

	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		lute.adjustVditorDOMListList(c)
		c = next
	}
}

func (lute *Lute) removeHighlightJSSpans(node *html.Node) {
	var spans []*html.Node
	for c := node; nil != c; c = c.NextSibling {
		lute.hljsSpans(c, &spans)
	}
	for _, span := range spans {
		span.Unlink()
	}
}

func (lute *Lute) hljsSpans(n *html.Node, spans *[]*html.Node) {
	if atom.Span == n.DataAtom && strings.HasPrefix(util.DomAttrValue(n, "class"), "hljs-") {
		*spans = append(*spans, n)
		text := util.DomText(n)
		n.InsertBefore(&html.Node{Type: html.TextNode, Data: text})
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.hljsSpans(c, spans)
	}
}

func (lute *Lute) removeEmptyNodes(node *html.Node) {
	var emptyNodes []*html.Node
	for c := node; nil != c; c = c.NextSibling {
		lute.searchEmptyNodes(c, &emptyNodes)
	}
	for _, emptyNode := range emptyNodes {
		emptyNode.Unlink()
	}
}

func (lute *Lute) searchEmptyNodes(n *html.Node, emptyNodes *[]*html.Node) {
	switch n.DataAtom {
	case 0:
		if lute.isInline(n.PrevSibling) || lute.isInline(n.NextSibling) || lute.isInline(n.Parent) {
			// 前节点或者后节点是行级节点的话保留该空白
			break
		}

		if html.TextNode == n.Type {
			data := strings.TrimLeft(n.Data, " ")
			data = strings.TrimRight(data, " ")
			for strings.Contains(data, "\n\n") {
				// 浏览器剪藏扩展列表下方段落缩进成为子块 https://github.com/siyuan-note/siyuan/issues/6289
				data = strings.ReplaceAll(data, "\n\n", "")
			}
			if "" == data {
				*emptyNodes = append(*emptyNodes, n)
				return
			}
		}

		parent := n.Parent
		if nil != parent && (atom.Ol == parent.DataAtom || atom.Ul == parent.DataAtom || atom.Li == parent.DataAtom) {
			if nil == n.NextSibling || (html.TextNode == n.NextSibling.Type || atom.Ul == n.NextSibling.DataAtom) || "" == strings.TrimSpace(n.Data) {
				n.Data = strings.TrimRight(n.Data, "\n\t ")
			}
		}
		if nil != parent && (atom.Table == parent.DataAtom || atom.Thead == parent.DataAtom || atom.Tbody == parent.DataAtom || atom.Tr == parent.DataAtom) {
			n.Data = strings.TrimSpace(n.Data)
		}

		if "" == n.Data {
			*emptyNodes = append(*emptyNodes, n)
		}

		if html.CommentNode == n.Type {
			*emptyNodes = append(*emptyNodes, n)
		}
	case atom.Span:
		if lc := n.LastChild; nil != lc && atom.Br == lc.DataAtom {
			// 如果行级标记节点最后一个子节点是 <br>，则将该 <br> 移动到该行级标记节点的后面
			// 表格内多个连续的超链接无法换行显示 https://github.com/siyuan-note/siyuan/issues/5966
			n.InsertAfter(lc)
		}

		if util.IsTempMarkSpan(n) {
			// 将嵌套在临时标记中的节点提升到临时标记节点之前
			*emptyNodes = append(*emptyNodes, n)
			var children []*html.Node
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				children = append(children, c)
			}
			for _, c := range children {
				n.InsertBefore(c)
			}
			return
		}
	case atom.Strong, atom.B, atom.Em, atom.I, atom.Del, atom.S, atom.Strike, atom.Mark:
		if nil != n.FirstChild {
			if atom.Br == n.FirstChild.DataAtom {
				*emptyNodes = append(*emptyNodes, n.FirstChild)
				n.InsertBefore(&html.Node{Type: html.ElementNode, DataAtom: atom.Br, Data: "br"})
			}
			if html.TextNode == n.FirstChild.Type {
				text := n.FirstChild.Data
				spaces := lute.prefixSpaces(text)
				if "" != spaces {
					n.FirstChild.Data = editor.Zwsp + n.FirstChild.Data
				}
			}
		}
		if nil != n.LastChild {
			if atom.Br == n.LastChild.DataAtom {
				*emptyNodes = append(*emptyNodes, n.LastChild)
				n.InsertAfter(&html.Node{Type: html.ElementNode, DataAtom: atom.Br, Data: "br"})
			}
			if html.TextNode == n.LastChild.Type {
				text := n.LastChild.Data
				spaces := lute.suffixSpaces(text)
				if "" != spaces {
					n.FirstChild.Data = n.FirstChild.Data + editor.Zwsp
				}
			}
		}
	default:
		if "katex" == util.DomAttrValue(n, "class") {
			*emptyNodes = append(*emptyNodes, n)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.searchEmptyNodes(c, emptyNodes)
	}

	switch n.DataAtom {
	case atom.Ol, atom.Ul:
		if dataType := util.DomAttrValue(n, "data-type"); "footnotes-defs-ol" == dataType {
			return
		}
		if nil != n.FirstChild && nil != n.FirstChild.FirstChild && atom.Input != n.FirstChild.FirstChild.DataAtom {
			return
		}

		text := util.DomText(n)
		if "" == text {
			*emptyNodes = append(*emptyNodes, n)
		}
	}
}

func (lute *Lute) mergeVditorDOMList0(n *html.Node) {
	switch n.DataAtom {
	case atom.Ul, atom.Ol:
		if nil != n.NextSibling && n.DataAtom == n.NextSibling.DataAtom && 1 == len(n.NextSibling.Attr) {
			for c := n.NextSibling.FirstChild; nil != c; {
				next := c.NextSibling
				c.Unlink()
				n.AppendChild(c)
				c = next
			}
			n.NextSibling.Unlink()
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.mergeVditorDOMList0(c)
	}
}

func (lute *Lute) adjustVditorDOMListTight0(n *html.Node) {
	switch n.DataAtom {
	case atom.Ul:
		if !lute.parentIs(n, atom.Pre) {
			lute.setDOMAttrValue(n, "data-tight", lute.isTightList(n))
		}
	case atom.Ol:
		if !lute.parentIs(n, atom.Pre) {
			lute.setDOMAttrValue(n, "data-tight", lute.isTightList(n))
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.adjustVditorDOMListTight0(c)
	}
}

func (lute *Lute) adjustVditorDOMListItemInP(n *html.Node) {
	switch n.DataAtom {
	case atom.Li:
		// li 换行时 id 重复需要重新生成
		if nil != n.PrevSibling && util.DomAttrValue(n.PrevSibling, "data-node-id") == util.DomAttrValue(n, "data-node-id") {
			lute.setDOMAttrValue(n, "data-node-id", ast.NewNodeID())
		}
		// 松散 li 换行时和上一个 li.last id 重复
		if nil != n.PrevSibling && nil != n.FirstChild {
			id := util.DomAttrValue(n.FirstChild, "data-node-id") // id 为空的话是行级节点，列表项行级排版自动换行问题 https://github.com/siyuan-note/siyuan/issues/379
			if "" != id && nil != n.PrevSibling.LastChild && util.DomAttrValue(n.PrevSibling.LastChild, "data-node-id") == id {
				lute.setDOMAttrValue(n.FirstChild, "data-node-id", ast.NewNodeID())
			}
		}

		// 在 li 下的每个非容器块节点用 p 包裹
		for c := n.FirstChild; nil != c; c = c.NextSibling {
			if lute.listItemEnter(n) {
				p := &html.Node{Type: html.ElementNode, Data: "p", DataAtom: atom.P}
				p.AppendChild(&html.Node{Type: html.TextNode, Data: editor.Caret})
				p.AppendChild(&html.Node{Type: html.ElementNode, Data: "br", DataAtom: atom.Br})
				n.FirstChild.Unlink()
				n.FirstChild.Unlink()
				n.AppendChild(p)
				c = p
				continue
			}

			if atom.P != c.DataAtom && atom.Blockquote != c.DataAtom && atom.Ul != c.DataAtom && atom.Ol != c.DataAtom && atom.Div != c.DataAtom {
				spans, nextBlock := lute.forwardNextBlock(c)
				p := &html.Node{Type: html.ElementNode, Data: "p", DataAtom: atom.P}
				c.InsertBefore(p)
				for _, span := range spans {
					span.Unlink()
					p.AppendChild(span)
				}
				if c = nextBlock; nil == c {
					break
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.adjustVditorDOMListItemInP(c)
	}
}

func (lute *Lute) removeCodeCode(n *html.Node) {
	if atom.Code == n.DataAtom && nil != n.FirstChild && atom.Code == n.FirstChild.DataAtom {
		// code.code 重复嵌套，则不处理外层 code
		for c := n.FirstChild; nil != c; {
			next := c.NextSibling
			c.Unlink()
			n.InsertBefore(c)
			c = next
		}
		n.Unlink()
		return
	}

	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		lute.removeCodeCode(c)
		c = next
	}
}
func (lute *Lute) adjustVditorDOMCodeA(n *html.Node) {
	// https://github.com/siyuan-note/siyuan/issues/11370
	if atom.Code == n.DataAtom && nil != n.FirstChild && atom.A == n.FirstChild.DataAtom && n.FirstChild == n.LastChild {
		// code.a 的情况将 a 移到 code 外层，即 a.code
		prev := n.PrevSibling
		next := n.NextSibling
		parent := n.Parent
		a := n.FirstChild
		a.Unlink()
		n.Unlink()

		var anchorTexts []*html.Node
		for c := a.FirstChild; nil != c; c = c.NextSibling {
			anchorTexts = append(anchorTexts, c)
			c.Unlink()
		}
		for _, anchorText := range anchorTexts {
			n.AppendChild(anchorText)
		}
		a.AppendChild(n)

		if nil != prev {
			prev.InsertAfter(a)
		} else if nil != next {
			next.InsertBefore(a)
		} else if nil != parent {
			parent.AppendChild(a)
		}
		return
	}

	for c := n.FirstChild; c != nil; {
		next := c.NextSibling
		lute.adjustVditorDOMCodeA(c)
		c = next
	}
}

// forwardNextBlock 向前移动至下一个块级节点，即跳过行级节点。
func (lute *Lute) forwardNextBlock(spanNode *html.Node) (spans []*html.Node, nextBlock *html.Node) {
	for next := spanNode; nil != next; next = next.NextSibling {
		switch next.DataAtom {
		case atom.Ol, atom.Ul, atom.Div, atom.Blockquote:
			return
		}
		spans = append(spans, next)
	}
	return
}

func (lute *Lute) listItemEnter(li *html.Node) bool {
	if nil == li.FirstChild {
		return false
	}
	if editor.Caret == li.FirstChild.Data && "br" == li.LastChild.Data {
		return true
	}
	return false
}

func (lute *Lute) isTightList(list *html.Node) string {
	for li := list.FirstChild; nil != li; li = li.NextSibling {
		var subLists, subDivs, subBlockquotes, subParagraphs int
		for c := li.FirstChild; nil != c; c = c.NextSibling {
			switch c.DataAtom {
			case atom.Ul, atom.Ol:
				subLists++
			case atom.Div:
				subDivs++
			case atom.Blockquote:
				subBlockquotes++
			case atom.P:
				subParagraphs++
			}
		}
		if 1 < subParagraphs || 1 < subBlockquotes || 1 < subDivs || 1 < subLists {
			return "false"
		}

		if 1 < subParagraphs+subDivs || 1 < subParagraphs+subBlockquotes || 1 < subParagraphs+subLists {
			return "false"
		}

	}
	return "true"
}

// genASTByVditorDOM 根据指定的 Vditor DOM 节点 n 进行深度优先遍历并逐步生成 Markdown 语法树 tree。
func (lute *Lute) genASTByVditorDOM(n *html.Node, tree *parse.Tree) {
	dataRender := util.DomAttrValue(n, "data-render")
	if "1" == dataRender || "2" == dataRender { // 1：浮动工具栏，2：preview 代码块、数学公式块
		return
	}

	dataType := util.DomAttrValue(n, "data-type")

	if atom.Div == n.DataAtom {
		if "code-block" == dataType || "html-block" == dataType || "math-block" == dataType || "yaml-front-matter" == dataType {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				lute.genASTByVditorDOM(c, tree)
			}
		} else if "link-ref-defs-block" == dataType {
			text := util.DomText(n)
			node := &ast.Node{Type: ast.NodeText, Tokens: []byte(text)}
			tree.Context.Tip.AppendChild(node)
		} else if "footnotes-block" == dataType {
			ol := n.FirstChild
			if atom.Ol != ol.DataAtom {
				return
			}

			for li := ol.FirstChild; nil != li; li = li.NextSibling {
				if "\n" == li.Data {
					continue
				}

				originalHTML := &bytes.Buffer{}
				if err := html.Render(originalHTML, li); nil == err {
					md := lute.vditorDOM2Md("<ol data-type=\"footnotes-defs-ol\">" + originalHTML.String() + "</ol>")
					label := util.DomAttrValue(li, "data-marker")
					md = md[3:] // 去掉列表项标记符 1.
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
					md = "[" + label + "]: " + md
					node := &ast.Node{Type: ast.NodeText, Tokens: []byte(md)}
					tree.Context.Tip.AppendChild(node)
				} else {
					panic(err)
				}
			}
		} else if "toc-block" == dataType {
			node := &ast.Node{Type: ast.NodeToC}
			tree.Context.Tip.AppendChild(node)
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
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		if "" == strings.TrimSpace(util.DomText(n)) {
			return
		}
		node.Type = ast.NodeHeading
		node.HeadingLevel = int(node.Tokens[1] - byte('0'))
		marker := util.DomAttrValue(n, "data-marker")
		if id := util.DomAttrValue(n, "data-id"); "" != id {
			n.LastChild.InsertAfter(&html.Node{Type: html.TextNode, Data: " {" + id + "}"})
		}
		node.HeadingSetext = "=" == marker || "-" == marker
		if !node.HeadingSetext {
			headingC8hMarker := &ast.Node{Type: ast.NodeHeadingC8hMarker}
			headingC8hMarker.Tokens = []byte(strings.Repeat("#", node.HeadingLevel))
			node.AppendChild(headingC8hMarker)
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Hr:
		node.Type = ast.NodeThematicBreak
		tree.Context.Tip.AppendChild(node)
	case atom.Blockquote:
		content := strings.TrimSpace(util.DomText(n))
		if "" == content || "&gt;" == content || editor.Caret == content {
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
		tight := util.DomAttrValue(n, "data-tight")
		if "true" == tight || "" == tight {
			node.ListData.Tight = true
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Li:
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
			node.ListData.Num, _ = strconv.Atoi(string(marker[0]))
			node.ListData.Delimiter = marker[len(marker)-1]
		}

		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Pre:
		if atom.Code == n.FirstChild.DataAtom {
			marker := util.DomAttrValue(n.Parent, "data-marker")
			if "" == marker {
				marker = "```"
			}

			var codeTokens []byte
			if nil != n.FirstChild.FirstChild {
				codeTokens = []byte(n.FirstChild.FirstChild.Data)
			}

			divDataType := util.DomAttrValue(n.Parent, "data-type")
			switch divDataType {
			case "math-block":
				node.Type = ast.NodeMathBlock
				node.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker})
				node.AppendChild(&ast.Node{Type: ast.NodeMathBlockContent, Tokens: codeTokens})
				node.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker})
				tree.Context.Tip.AppendChild(node)
			case "yaml-front-matter":
				node.Type = ast.NodeYamlFrontMatter
				node.AppendChild(&ast.Node{Type: ast.NodeYamlFrontMatterOpenMarker})
				node.AppendChild(&ast.Node{Type: ast.NodeYamlFrontMatterContent, Tokens: codeTokens})
				node.AppendChild(&ast.Node{Type: ast.NodeYamlFrontMatterCloseMarker})
				tree.Context.Tip.AppendChild(node)
			case "html-block":
				node.Type = ast.NodeHTMLBlock
				node.Tokens = codeTokens
				tree.Context.Tip.AppendChild(node)
			default:
				node.Type = ast.NodeCodeBlock
				node.IsFencedCodeBlock = true
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: []byte(marker), CodeBlockFenceLen: len(marker)})
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker})
				class := util.DomAttrValue(n.FirstChild, "class")
				if strings.Contains(class, "language-") {
					language := class[len("language-"):]
					node.LastChild.CodeBlockInfo = []byte(language)
				}

				content := &ast.Node{Type: ast.NodeCodeBlockCode, Tokens: codeTokens}
				node.AppendChild(content)
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: []byte(marker), CodeBlockFenceLen: len(marker)})
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
		marker := util.DomAttrValue(n, "data-marker")
		if "" == marker {
			marker = "*"
		}
		if "_" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeEmU8eOpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeEmA6kOpenMarker, Tokens: []byte(marker)})
		}
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && editor.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: editor.CaretTokens})
			if "_" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeEmU8eCloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, editor.Zwsp, "")

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
		marker := util.DomAttrValue(n, "data-marker")
		if "" == marker {
			marker = "**"
		}
		if "__" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongU8eOpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kOpenMarker, Tokens: []byte(marker)})
		}
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && editor.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: editor.CaretTokens})
			if "__" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeStrongU8eCloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kCloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, editor.Zwsp, "")
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
		marker := util.DomAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1OpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2OpenMarker, Tokens: []byte(marker)})
		}
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && editor.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: editor.CaretTokens})
			if "~" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1CloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2CloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, editor.Zwsp, "")
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
		marker := util.DomAttrValue(n, "data-marker")
		if "=" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeMark1OpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeMark2OpenMarker, Tokens: []byte(marker)})
		}
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && editor.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: editor.CaretTokens})
			if "=" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeMark1CloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeMark2CloseMarker, Tokens: []byte(marker)})
			}
			return
		}

		n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, editor.Zwsp, "")
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
		if "html-inline" == dataType {
			// 所见即所得行级 HTML 解析 https://github.com/Vanessa219/vditor/issues/1156
			node.Type = ast.NodeInlineHTML
			node.Tokens = codeTokens
			tree.Context.Tip.AppendChild(node)
			return
		}

		marker := util.DomAttrValue(n, "data-marker")
		if "" == marker {
			marker = "`"
		}
		if bytes.HasPrefix(codeTokens, []byte("`")) {
			codeTokens = append([]byte(" "), codeTokens...)
			codeTokens = append(codeTokens, ' ')
		}
		node.Type = ast.NodeCodeSpan
		node.CodeMarkerLen = len(marker)
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanOpenMarker})
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanContent, Tokens: codeTokens})
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanCloseMarker})
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
				if nil != n.Parent.Parent && nil != n.Parent.Parent.NextSibling && atom.Li == n.Parent.Parent.NextSibling.DataAtom {
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
		if n.FirstChild == nil || n.FirstChild.Type == html.TextNode {
			text := util.DomText(n)
			if "" == text || editor.Zwsp == text {
				return
			}
		}

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
		} else {
			if "link-ref" == dataType {
				node.Type = ast.NodeText
				content := "![" + util.DomAttrValue(n, "alt") + "][" + util.DomAttrValue(n, "data-link-label") + "]"
				node.Tokens = []byte(content)
				tree.Context.Tip.AppendChild(node)
				return
			}

			node.Type = ast.NodeImage
			node.AppendChild(&ast.Node{Type: ast.NodeBang})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
			if "" != imgAlt {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkText, Tokens: []byte(imgAlt)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			src := util.DomAttrValue(n, "src")
			if "" != lute.RenderOptions.LinkBase {
				src = strings.ReplaceAll(src, lute.RenderOptions.LinkBase, "")
			}
			if "" != lute.RenderOptions.LinkPrefix {
				src = strings.ReplaceAll(src, lute.RenderOptions.LinkPrefix, "")
			}
			node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: []byte(src)})
			linkTitle := util.DomAttrValue(n, "title")
			if "" != linkTitle {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
				node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: []byte(linkTitle)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
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
		if nil != node.Parent.Parent && nil != node.Parent.Parent.ListData { // ul.li.input
			node.Parent.Parent.ListData.Typ = 3
		}
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
			node.Tokens = []byte("[" + util.DomAttrValue(n, "data-footnotes-label") + "]")
			if strings.Contains(n.FirstChild.Data, editor.Caret) {
				node.Tokens = append(node.Tokens, editor.CaretTokens...)
			}
			tree.Context.Tip.AppendChild(node)
		}
		return
	case atom.Span:
		if nil == n.FirstChild {
			break
		}

		if strings.Contains(class, "vditor-comment") {
			node.Type = ast.NodeInlineHTML
			buf := bytes.Buffer{}
			buf.WriteString("<span ")
			for i, attr := range n.Attr {
				buf.WriteString(attr.Key)
				if "" != attr.Val {
					buf.WriteString("=\"")
					buf.WriteString(attr.Val)
					buf.WriteString("\"")
				}
				if i < len(n.Attr)-1 {
					buf.WriteString(" ")
				}
			}
			buf.WriteString(">")
			node.Tokens = buf.Bytes()
			tree.Context.Tip.AppendChild(node)
			break
		}

		if "link-ref" == dataType {
			node.Type = ast.NodeText
			content := "[" + n.FirstChild.Data + "][" + util.DomAttrValue(n, "data-link-label") + "]"
			if nil != n.NextSibling && "2" == util.DomAttrValue(n.NextSibling, "data-render") {
				// 图片引用风格 ![text][label]
				content = "!" + content
			}
			node.Tokens = []byte(content)
			tree.Context.Tip.AppendChild(node)
			return
		}

		var codeTokens []byte
		if editor.Zwsp == n.FirstChild.Data && "" == util.DomAttrValue(n, "style") && nil != n.FirstChild.NextSibling {
			codeTokens = []byte(n.FirstChild.NextSibling.FirstChild.Data)
		} else if atom.Code == n.FirstChild.DataAtom {
			codeTokens = []byte(n.FirstChild.FirstChild.Data)
			if editor.Zwsp == string(codeTokens) {
				break
			}
		}
		if "math-inline" == dataType {
			node.Type = ast.NodeInlineMath
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathOpenMarker})
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathContent, Tokens: codeTokens})
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathCloseMarker})
			tree.Context.Tip.AppendChild(node)
			return
		} else if "html-inline" == dataType {
			node.Type = ast.NodeInlineHTML
			node.Tokens = codeTokens
			tree.Context.Tip.AppendChild(node)
			return
		} else if "code-inline" == dataType {
			node.Tokens = codeTokens
			tree.Context.Tip.AppendChild(node)
			return
		} else if "html-entity" == dataType {
			node.Type = ast.NodeText
			node.Tokens = codeTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		break
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
		node.Type = ast.NodeInlineHTML
		node.Tokens = util.DomHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Summary:
		return
	default:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = util.DomHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTByVditorDOM(c, tree)
	}

	switch n.DataAtom {
	case atom.Span:
		if strings.Contains(class, "vditor-comment") {
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeInlineHTML, Tokens: []byte("</span>")})
		}
	case atom.Em, atom.I:
		marker := util.DomAttrValue(n, "data-marker")
		if "" == marker {
			marker = "*"
		}
		if "_" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeEmU8eCloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker, Tokens: []byte(marker)})
		}
	case atom.Strong, atom.B:
		marker := util.DomAttrValue(n, "data-marker")
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
	case atom.Del, atom.S, atom.Strike:
		marker := util.DomAttrValue(n, "data-marker")
		if "~" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1CloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2CloseMarker, Tokens: []byte(marker)})
		}
	case atom.Mark:
		marker := util.DomAttrValue(n, "data-marker")
		if "=" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeMark1CloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeMark2CloseMarker, Tokens: []byte(marker)})
		}
	case atom.Details:
		tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeHTMLBlock, Tokens: []byte("</details>")})
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

func (lute *Lute) domChild(n *html.Node, dataAtom atom.Atom) *html.Node {
	if nil == n {
		return nil
	}

	for c := n.FirstChild; nil != c; c = c.NextSibling {
		if ret := lute.domChild0(c, dataAtom); nil != ret {
			return ret
		}
	}
	return nil
}

func (lute *Lute) domChild0(n *html.Node, dataAtom atom.Atom) *html.Node {
	if n.DataAtom == dataAtom {
		return n
	}

	for c := n.FirstChild; nil != c; c = c.NextSibling {
		if ret := lute.domChild0(c, dataAtom); nil != ret {
			return ret
		}
	}
	return nil
}

func (lute *Lute) setDOMAttrValue(n *html.Node, attrName, attrVal string) {
	if nil == n {
		return
	}

	for _, attr := range n.Attr {
		if attr.Key == attrName {
			attr.Val = attrVal
			return
		}
	}

	n.Attr = append(n.Attr, &html.Attribute{Key: attrName, Val: attrVal})
}

func (lute *Lute) removeDOMAttr(n *html.Node, attrName string) {
	if nil == n {
		return
	}

	if 1 > len(n.Attr) {
		return
	}

	tmp := (n.Attr)[:0]
	for _, attr := range n.Attr {
		if attr.Key != attrName {
			tmp = append(tmp, attr)
		}
	}
	n.Attr = tmp
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
		buffer.Write(util.DomHTML(n))
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

func (lute *Lute) getParent(n *html.Node, parentType atom.Atom) *html.Node {
	for p := n.Parent; nil != p; p = p.Parent {
		if parentType == p.DataAtom {
			return p
		}
	}
	return nil
}

func (lute *Lute) isCaret(n *html.Node) (isCaret, isEmptyText bool) {
	text := util.DomText(n)
	trimSpaceText := strings.TrimSpace(text)
	if 1 > len(trimSpaceText) && 1 < len(text) && strings.Contains(text, editor.Caret) {
		return true, false
	}
	isCaret = editor.Caret == text || editor.Zwsp+editor.Caret == text || editor.Caret+editor.Zwsp == text
	isEmptyText = "" == text || editor.Zwsp == text
	if "" != util.DomAttrValue(n, "data-content") {
		isEmptyText = false
	}
	return
}

func (lute *Lute) isEmptyText(n *html.Node) bool {
	if nil != n.FirstChild && "block-ref" == util.DomAttrValue(n.FirstChild, "data-type") {
		return false
	}

	text := strings.TrimSpace(util.DomText(n))
	if "" == text || editor.Zwsp == text {
		return true
	}
	if editor.Zwsp+editor.Caret == text || editor.Caret+editor.Zwsp == text {
		return true
	}
	return false
}

func (lute *Lute) startsWithNewline(n *html.Node) bool {
	return strings.HasPrefix(n.Data, "\n") || strings.HasPrefix(n.Data, editor.Zwsp+"\n")
}

func (lute *Lute) isInline(n *html.Node) bool {
	if nil == n {
		return false
	}

	return 0 == n.DataAtom ||
		atom.Code == n.DataAtom ||
		atom.Strong == n.DataAtom || atom.B == n.DataAtom ||
		atom.Em == n.DataAtom || atom.I == n.DataAtom ||
		atom.Mark == n.DataAtom ||
		atom.Del == n.DataAtom || atom.S == n.DataAtom || atom.Strike == n.DataAtom ||
		atom.A == n.DataAtom ||
		atom.Img == n.DataAtom ||
		atom.U == n.DataAtom ||
		atom.Kbd == n.DataAtom ||
		atom.Span == n.DataAtom
}

func (lute *Lute) prefixSpaces(text string) (ret string) {
	for _, c := range text {
		if ' ' == c || 160 == c {
			ret += " "
		} else {
			return
		}
	}
	return
}

func (lute *Lute) suffixSpaces(text string) (ret string) {
	for i := len(text) - 1; i >= 0; i-- {
		if ' ' == text[i] || 160 == text[i] {
			ret += " "
		} else {
			return
		}
	}
	return
}
