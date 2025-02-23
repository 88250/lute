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
	"fmt"
	"strconv"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

func (lute *Lute) SpinBlockDOM(ivHTML string) (ovHTML string) {
	//fmt.Println(ivHTML)
	markdown := lute.blockDOM2Md(ivHTML)
	markdown = strings.ReplaceAll(markdown, editor.Zwsp, "")
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)

	firstChild := tree.Root.FirstChild
	lastChildMaybeIAL := tree.Root.LastChild.Previous
	if ast.NodeParagraph == firstChild.Type && "" == firstChild.ID && nil != lastChildMaybeIAL && firstChild != lastChildMaybeIAL.Previous &&
		ast.NodeKramdownBlockIAL == lastChildMaybeIAL.Type {
		// 软换行后生成多个块，需要把老 ID 调整到第一个块上
		firstChild.ID, lastChildMaybeIAL.Previous.ID = lastChildMaybeIAL.Previous.ID, ""
		firstChild.KramdownIAL, lastChildMaybeIAL.Previous.KramdownIAL = lastChildMaybeIAL.Previous.KramdownIAL, nil
		firstChild.InsertAfter(lastChildMaybeIAL)
	}
	if ast.NodeKramdownBlockIAL == firstChild.Type && nil != firstChild.Next && ast.NodeKramdownBlockIAL == firstChild.Next.Type && util.IsDocIAL(firstChild.Next.Tokens) {
		// 空段落块还原
		ialArray := parse.Tokens2IAL(firstChild.Tokens)
		ial := parse.IAL2Map(ialArray)
		p := &ast.Node{Type: ast.NodeParagraph, ID: ial["id"], KramdownIAL: ialArray}
		firstChild.InsertBefore(p)
	}

	// 使用 Markdown 标记符嵌套行级元素后被还原为纯文本 https://github.com/siyuan-note/siyuan/issues/7637
	// 这里需要将混合嵌套（比如 <strong><span a></span></strong>）的行级元素拆分为多个平铺的行级元素（<span strong> 和 <span strong a>）
	parse.NestedInlines2FlattedSpansHybrid(tree, false)

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) HTML2BlockDOM(sHTML string) (vHTML string) {
	//fmt.Println(sHTML)
	markdown, err := lute.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewProtyleRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.HTML2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	return
}

func (lute *Lute) BlockDOM2HTML(vHTML string) (sHTML string) {
	markdown := lute.blockDOM2Md(vHTML)
	sHTML = lute.Md2HTML(markdown)
	return
}

func (lute *Lute) BlockDOM2InlineBlockDOM(vHTML string) (vIHTML string) {
	markdown := lute.blockDOM2Md(vHTML)
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	var inlines []*ast.Node
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		if ast.NodeTableCell == n.Type {
			n.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
			return ast.WalkContinue
		}

		if !n.IsBlock() {
			if ast.NodeCodeBlockCode != n.Type && ast.NodeMathBlockContent != n.Type && ast.NodeTaskListItemMarker != n.Type &&
				ast.NodeTableHead != n.Type && ast.NodeTableRow != n.Type && ast.NodeTableCell != n.Type {
				if ast.NodeBlockquoteMarker == n.Type {
					inlines = append(inlines, &ast.Node{Type: ast.NodeText, Tokens: []byte(">")})
					return ast.WalkSkipChildren
				}

				inlines = append(inlines, n)
				return ast.WalkSkipChildren
			}
		} else if ast.NodeHTMLBlock == n.Type {
			inlines = append(inlines, &ast.Node{Type: ast.NodeText, Tokens: n.Tokens})
			return ast.WalkSkipChildren
		}
		return ast.WalkContinue
	})

	var unlinks []*ast.Node
	for n := tree.Root.FirstChild; nil != n; n = n.Next {
		unlinks = append(unlinks, n)
	}
	for _, n := range unlinks {
		n.Unlink()
	}

	for _, n := range inlines {
		tree.Root.AppendChild(n)
	}

	renderer := render.NewProtyleRenderer(tree, lute.RenderOptions)
	output := renderer.Render()
	vIHTML = util.BytesToStr(output)
	vIHTML = strings.TrimSpace(vIHTML)
	return
}

func (lute *Lute) Md2BlockDOM(markdown string, reserveEmptyParagraph bool) (vHTML string) {
	vHTML, _ = lute.Md2BlockDOMTree(markdown, reserveEmptyParagraph)
	return
}

func (lute *Lute) Md2BlockDOMTree(markdown string, reserveEmptyParagraph bool) (vHTML string, tree *parse.Tree) {
	tree = parse.Parse("", []byte(markdown), lute.ParseOptions)

	parse.TextMarks2Inlines(tree) // 先将 TextMark 转换为 Inlines https://github.com/siyuan-note/siyuan/issues/13056
	parse.NestedInlines2FlattedSpansHybrid(tree, false)
	if reserveEmptyParagraph {
		ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
			if !entering {
				return ast.WalkContinue
			}

			if n.IsEmptyBlockIAL() {
				p := &ast.Node{Type: ast.NodeParagraph}
				p.KramdownIAL = parse.Tokens2IAL(n.Tokens)
				p.ID = p.IALAttr("id")
				n.InsertBefore(p)
				return ast.WalkContinue
			}
			return ast.WalkContinue
		})
	}

	renderer := render.NewProtyleRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	return
}

func (lute *Lute) InlineMd2BlockDOM(markdown string) (vHTML string) {
	tree := parse.Inline("", []byte(markdown), lute.ParseOptions)
	parse.NestedInlines2FlattedSpansHybrid(tree, false)
	renderer := render.NewProtyleRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	return
}

func (lute *Lute) BlockDOM2Md(htmlStr string) (kramdown string) {
	kramdown = lute.blockDOM2Md(htmlStr)
	kramdown = strings.ReplaceAll(kramdown, editor.Zwsp, "")
	return
}

func (lute *Lute) BlockDOM2StdMd(htmlStr string) (markdown string) {
	htmlStr = strings.ReplaceAll(htmlStr, editor.Zwsp, "")

	// DOM 转 AST
	tree := lute.BlockDOM2Tree(htmlStr)

	// 将 kramdown IAL 节点内容置空
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		if ast.NodeKramdownBlockIAL == n.Type || ast.NodeKramdownSpanIAL == n.Type {
			n.Tokens = nil
		}
		return ast.WalkContinue
	})

	options := render.NewOptions()
	options.AutoSpace = false
	options.FixTermTypo = false
	options.KramdownBlockIAL = true
	options.KramdownSpanIAL = true
	options.KeepParagraphBeginningSpace = true
	renderer := render.NewProtyleExportMdRenderer(tree, options)
	formatted := renderer.Render()
	markdown = util.BytesToStr(formatted)
	return
}

func (lute *Lute) BlockDOM2Text(htmlStr string) (text string) {
	tree := lute.BlockDOM2Tree(htmlStr)
	return tree.Root.Text()
}

func (lute *Lute) BlockDOM2TextLen(htmlStr string) int {
	tree := lute.BlockDOM2Tree(htmlStr)
	return tree.Root.TextLen()
}

func (lute *Lute) BlockDOM2Content(htmlStr string) (text string) {
	tree := lute.BlockDOM2Tree(htmlStr)
	return tree.Root.Content()
}

func (lute *Lute) BlockDOM2EscapeMarkerContent(htmlStr string) (text string) {
	tree := lute.BlockDOM2Tree(htmlStr)
	return tree.Root.EscapeMarkerContent()
}

func (lute *Lute) Tree2BlockDOM(tree *parse.Tree, options *render.Options) (vHTML string) {
	renderer := render.NewProtyleRenderer(tree, options)
	for nodeType, rendererFunc := range lute.Md2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	vHTML = strings.ReplaceAll(vHTML, editor.Caret, "<wbr>")
	return
}

func (lute *Lute) RenderNodeBlockDOM(node *ast.Node) string {
	root := &ast.Node{Type: ast.NodeDocument}
	tree := &parse.Tree{Root: root, Context: &parse.Context{ParseOption: lute.ParseOptions}}
	renderer := render.NewProtyleRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	renderer.Writer = &bytes.Buffer{}
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := renderer.RendererFuncs[n.Type]
		return rendererFunc(n, entering)
	})
	return renderer.Writer.String()
}

func (lute *Lute) BlockDOM2Tree(htmlStr string) (ret *parse.Tree) {
	htmlStr = strings.ReplaceAll(htmlStr, "\n<wbr>\n</strong>", "</strong>\n<wbr>\n")
	htmlStr = strings.ReplaceAll(htmlStr, "\n<wbr>\n</em>", "</em>\n<wbr>\n")
	htmlStr = strings.ReplaceAll(htmlStr, "\n<wbr>\n</s>", "</s>\n<wbr>\n")
	htmlStr = strings.ReplaceAll(htmlStr, "\n<wbr>\n</u>", "</u>\n<wbr>\n")
	htmlStr = strings.ReplaceAll(htmlStr, "\n<wbr>\n</span>", "</span>\n<wbr>\n")

	// Improve `inline code` markdown editing https://github.com/siyuan-note/siyuan/issues/9978
	// spinBlockDOMTests #212
	htmlStr = strings.ReplaceAll(htmlStr, "`<wbr></span>", "</span>`<wbr>")

	htmlStr = strings.ReplaceAll(htmlStr, "<wbr>", editor.Caret)

	var startSpaces, endSpaces int
	for _, c := range htmlStr {
		if ' ' == c {
			startSpaces++
		} else {
			break
		}
	}
	for i := len(htmlStr) - 1; i >= 0; i-- {
		if ' ' == htmlStr[i] {
			endSpaces++
		} else {
			break
		}
	}
	htmlStr = strings.TrimSpace(htmlStr)
	htmlStr = strings.Repeat("&nbsp;", startSpaces) + htmlStr + strings.Repeat("&nbsp;", endSpaces)

	// 替换结尾空白，否则 HTML 解析会产生冗余节点导致生成空的代码块
	for strings.HasSuffix(htmlStr, "\t\n") {
		htmlStr = strings.TrimSuffix(htmlStr, "\t\n") + "\n"
	}
	for strings.HasSuffix(htmlStr, "    \n") {
		htmlStr = strings.TrimSuffix(htmlStr, "    \n") + "\n"
	}

	// 将字符串解析为 DOM 树
	htmlRoot := lute.parseHTML(htmlStr)
	if nil == htmlRoot {
		return
	}

	// 调整 DOM 结构
	lute.adjustVditorDOM(htmlRoot)

	// 将 HTML 树转换为 Markdown AST
	ret = &parse.Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument}, Context: &parse.Context{ParseOption: lute.ParseOptions}}
	ret.Context.Tip = ret.Root
	for c := htmlRoot.FirstChild; nil != c; c = c.NextSibling {
		lute.genASTByBlockDOM(c, ret)
	}

	// 调整树结构
	ast.Walk(ret.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n.Type {
			case ast.NodeInlineHTML, ast.NodeHTMLBlock, ast.NodeCodeSpanContent, ast.NodeCodeBlockCode, ast.NodeInlineMathContent, ast.NodeMathBlockContent,
				ast.NodeCodeSpan, ast.NodeInlineMath:
				if nil != n.Next && ast.NodeCodeSpan == n.Next.Type && n.CodeMarkerLen == n.Next.CodeMarkerLen && nil != n.FirstChild && nil != n.FirstChild.Next {
					// 合并代码节点 https://github.com/Vanessa219/vditor/issues/167
					n.FirstChild.Next.Tokens = append(n.FirstChild.Next.Tokens, n.Next.FirstChild.Next.Tokens...)
					n.Next.Unlink()
				}
			case ast.NodeStrong, ast.NodeEmphasis, ast.NodeStrikethrough, ast.NodeUnderline:
				lute.MergeSameSpan(n)
			case ast.NodeTextMark:
				lute.MergeSameTextMark(n)
			case ast.NodeText:
				n.Tokens = bytes.ReplaceAll(n.Tokens, []byte("\u00a0"), []byte(" "))
			}
		}
		return ast.WalkContinue
	})
	return
}

func (lute *Lute) MergeSameTextMark(n *ast.Node) {
	if nil == n.Previous {
		return
	}

	mergeWithIAL := false
	mergeWithZwsp := false
	if ast.NodeKramdownSpanIAL == n.Previous.Type {
		if nil == n.Next || ast.NodeKramdownSpanIAL != n.Next.Type || nil == n.Previous.Previous {
			return
		}

		if !bytes.Equal(n.Previous.Tokens, n.Next.Tokens) {
			return
		}

		if !n.IsSameTextMarkType(n.Previous.Previous) {
			return
		}

		mergeWithIAL = true
	} else {
		previewText := n.Previous.TokensStr()
		if ast.NodeText == n.Previous.Type &&
			!strings.Contains(previewText, "　") && !strings.Contains(previewText, " ") && !strings.Contains(previewText, "\n") && !strings.Contains(previewText, "\t") &&
			"" == strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(previewText, editor.Zwsp, ""), editor.Caret, "")) &&
			nil != n.Previous.Previous && n.IsSameTextMarkType(n.Previous.Previous) {
			mergeWithZwsp = true
		} else {
			if n.Type != n.Previous.Type || !n.IsSameTextMarkType(n.Previous) {
				return
			}
		}
	}

	types := strings.Split(n.TextMarkType, " ")
	m := map[string]bool{}
	for _, t := range types {
		m[t] = true
	}
	var allowMerge []string
	for k, _ := range m {
		switch k {
		case "code", "em", "strong", "s", "mark", "u", "sub", "sup", "kbd", "text", "tag", "block-ref", "a":
			allowMerge = append(allowMerge, k)
		}
	}
	for _, k := range allowMerge {
		delete(m, k)
	}
	if 0 < len(m) {
		return
	}

	if mergeWithIAL || mergeWithZwsp {
		content := n.TextMarkTextContent
		n.TextMarkTextContent = n.Previous.Previous.TextMarkTextContent
		if strings.Contains(n.Previous.TokensStr(), editor.Caret) {
			n.TextMarkTextContent += editor.Caret
		}
		n.TextMarkTextContent += content
		n.Previous.Previous.Unlink()
	} else {
		n.TextMarkTextContent = n.Previous.TextMarkTextContent + n.TextMarkTextContent
	}
	n.Previous.Unlink()
	n.SortTextMarkDataTypes()
}

func (lute *Lute) MergeSameSpan(n *ast.Node) {
	if nil == n.Next || n.Type != n.Next.Type {
		return
	}
	if nil != n.Next.Next && ast.NodeKramdownSpanIAL == n.Next.Next.Type {
		return
	}

	var spanChildren []*ast.Node
	n.Next.FirstChild.Unlink() // open marker
	n.Next.LastChild.Unlink()  // close marker
	for c := n.Next.FirstChild; nil != c; c = c.Next {
		spanChildren = append(spanChildren, c)
	}
	for _, c := range spanChildren {
		n.LastChild.InsertBefore(c)
	}
	n.Next.Unlink()
}

func (lute *Lute) CancelSuperBlock(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	if ast.NodeSuperBlock != tree.Root.FirstChild.Type {
		return ivHTML
	}

	sb := tree.Root.FirstChild

	var blocks []*ast.Node
	for b := sb.FirstChild; nil != b; b = b.Next {
		blocks = append(blocks, b)
	}
	for _, b := range blocks {
		tree.Root.AppendChild(b)
	}
	sb.Unlink()

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) CancelList(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	if ast.NodeList != tree.Root.FirstChild.Type {
		return ivHTML
	}

	list := tree.Root.FirstChild

	var appends, unlinks []*ast.Node
	for li := list.FirstChild; nil != li; li = li.Next {
		for c := li.FirstChild; nil != c; c = c.Next {
			if ast.NodeTaskListItemMarker != c.Type {
				appends = append(appends, c)
			}
		}
		unlinks = append(unlinks, li)
	}
	for _, c := range appends {
		tree.Root.AppendChild(c)
	}
	for _, n := range unlinks {
		n.Unlink()
	}
	list.Unlink()

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) CancelBlockquote(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	if ast.NodeBlockquote != tree.Root.FirstChild.Type {
		return ivHTML
	}

	bq := tree.Root.FirstChild

	var appends, unlinks []*ast.Node
	for sub := bq.FirstChild; nil != sub; sub = sub.Next {
		if ast.NodeBlockquoteMarker != sub.Type {
			appends = append(appends, sub)
		}
		unlinks = append(unlinks, sub)
	}
	for _, c := range appends {
		tree.Root.AppendChild(c)
	}
	bq.Unlink()

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) Blocks2Ps(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	node := tree.Root.FirstChild

	var appends, unlinks []*ast.Node
	for n := node; nil != n; n = n.Next {
		switch n.Type {
		case ast.NodeHeading:
			n.Type = ast.NodeParagraph
		case ast.NodeBlockquote:
			for c := n.FirstChild; nil != c; c = c.Next {
				if ast.NodeBlockquoteMarker == c.Type {
					unlinks = append(unlinks, c)
					continue
				}
				appends = append(appends, c)
			}
			unlinks = append(unlinks, n)
		case ast.NodeList:
			for li := n.FirstChild; nil != li; li = li.Next {
				for c := li.FirstChild; nil != c; c = c.Next {
					if ast.NodeTaskListItemMarker != c.Type {
						appends = append(appends, c)
					}
				}
				unlinks = append(unlinks, li)
			}
			unlinks = append(unlinks, n)
		}
	}
	for _, n := range unlinks {
		n.Unlink()
	}
	for _, c := range appends {
		tree.Root.AppendChild(c)
	}
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) Blocks2Hs(ivHTML, level string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	node := tree.Root.FirstChild

	for p := node; nil != p; p = p.Next {
		if ast.NodeParagraph == p.Type || ast.NodeHeading == p.Type {
			p.Type = ast.NodeHeading
			if nil != p.FirstChild {
				p.FirstChild.Tokens = bytes.ReplaceAll(p.FirstChild.Tokens, []byte("\n"), nil)
				p.FirstChild.Tokens = bytes.TrimLeft(p.FirstChild.Tokens, " \t\n")
			}
			p.HeadingLevel, _ = strconv.Atoi(level)
		}
	}
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) OL2TL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)

	tree.Root.FirstChild.ListData.Typ = 3
	for li := tree.Root.FirstChild.FirstChild; nil != li; li = li.Next {
		if ast.NodeListItem == li.Type {
			li.ListData.Typ = 3
			li.PrependChild(&ast.Node{Type: ast.NodeTaskListItemMarker})
		}
	}
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) UL2TL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)

	tree.Root.FirstChild.ListData.Typ = 3
	for li := tree.Root.FirstChild.FirstChild; nil != li; li = li.Next {
		if ast.NodeListItem == li.Type {
			li.ListData.Typ = 3
			li.PrependChild(&ast.Node{Type: ast.NodeTaskListItemMarker})
		}
	}
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) TL2OL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	list := tree.Root.FirstChild
	if ast.NodeList != list.Type || 3 != list.ListData.Typ {
		return ivHTML
	}

	num := 1
	list.ListData.Typ = 1
	var unlinks []*ast.Node
	for li := list.FirstChild; nil != li; li = li.Next {
		if ast.NodeKramdownBlockIAL == li.Type {
			continue
		}
		unlinks = append(unlinks, li.FirstChild) // task marker
		li.ListData.Typ = 1
		li.ListData.Num = num
		num++
	}
	for _, n := range unlinks {
		n.Unlink()
	}

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) TL2UL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	list := tree.Root.FirstChild
	if ast.NodeList != list.Type || 3 != list.ListData.Typ {
		return ivHTML
	}

	list.ListData.Typ = 0
	var unlinks []*ast.Node
	for li := list.FirstChild; nil != li; li = li.Next {
		if ast.NodeKramdownBlockIAL == li.Type {
			continue
		}
		unlinks = append(unlinks, li.FirstChild) // task marker
		li.ListData.Typ = 0
	}
	for _, n := range unlinks {
		n.Unlink()
	}

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) OL2UL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	list := tree.Root.FirstChild
	if ast.NodeList != list.Type {
		return ivHTML
	}

	list.ListData.Typ = 0
	for li := list.FirstChild; nil != li; li = li.Next {
		if ast.NodeKramdownBlockIAL == li.Type {
			continue
		}
		li.ListData.Typ = 0
	}

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) UL2OL(ivHTML string) (ovHTML string) {
	tree := lute.BlockDOM2Tree(ivHTML)
	list := tree.Root.FirstChild
	if ast.NodeList != list.Type {
		return ivHTML
	}

	num := 1
	list.ListData.Typ = 1
	for li := list.FirstChild; nil != li; li = li.Next {
		if ast.NodeKramdownBlockIAL == li.Type {
			continue
		}
		li.ListData.Typ = 1
		li.ListData.Num = num
		num++
	}

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) blockDOM2Md(htmlStr string) (markdown string) {
	tree := lute.BlockDOM2Tree(htmlStr)

	// 将 AST 进行 Markdown 格式化渲染
	options := render.NewOptions()
	options.AutoSpace = false
	options.FixTermTypo = false
	options.KramdownBlockIAL = true
	options.KramdownSpanIAL = true
	options.KeepParagraphBeginningSpace = true
	options.ProtyleWYSIWYG = true
	options.SuperBlock = true
	renderer := render.NewFormatRenderer(tree, options)
	formatted := renderer.Render()
	markdown = string(formatted)
	return
}

func (lute *Lute) genASTByBlockDOM(n *html.Node, tree *parse.Tree) {
	class := util.DomAttrValue(n, "class")

	// Custom dom, which will be omitted when build tree https://github.com/88250/lute/issues/206
	if strings.Contains(class, "protyle-custom") {
		return
	}

	if "protyle-attr" == class ||
		strings.Contains(class, "__copy") ||
		strings.Contains(class, "protyle-linenumber__rows") ||
		strings.Contains(class, "hljs") {
		return
	}

	if "1" == util.DomAttrValue(n, "spin") {
		return
	}

	if strings.Contains(class, "protyle-action") {
		if ast.NodeCodeBlock == tree.Context.Tip.Type {
			languageNode := n.FirstChild
			language := ""
			if nil != languageNode.FirstChild {
				language = languageNode.FirstChild.Data
			}
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker, CodeBlockInfo: util.StrToBytes(language)})
			code := util.DomText(n.NextSibling.LastChild)
			if strings.HasSuffix(code, "\n\n"+editor.Caret) {
				code = strings.TrimSuffix(code, "\n\n"+editor.Caret)
				code += "\n" + editor.Caret + "\n"
			}
			lines := strings.Split(code, "\n")
			buf := bytes.Buffer{}
			for i, line := range lines {
				if strings.Contains(line, "```") {
					line = strings.ReplaceAll(line, editor.Zwj+"```", "```")
					line = strings.ReplaceAll(line, "```", editor.Zwj+"```")
				}
				buf.WriteString(line)
				if i < len(lines)-1 {
					buf.WriteByte('\n')
				}
			}
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeCodeBlockCode, Tokens: buf.Bytes()})
		} else if ast.NodeListItem == tree.Context.Tip.Type {
			if 3 == tree.Context.Tip.ListData.Typ { // 任务列表
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeTaskListItemMarker, TaskListItemChecked: strings.Contains(util.DomAttrValue(n.Parent, "class"), "protyle-task--done")})
			}
		}
		return
	}

	if "true" == util.DomAttrValue(n, "contenteditable") {
		lute.genASTContenteditable(n, tree)
		return
	}

	dataType := ast.Str2NodeType(util.DomAttrValue(n, "data-type"))

	nodeID := util.DomAttrValue(n, "data-node-id")
	node := &ast.Node{ID: nodeID}
	if "" != node.ID && !lute.parentIs(n, atom.Table) {
		node.KramdownIAL = [][]string{{"id", node.ID}}
		ialTokens := lute.setBlockIAL(n, node)
		ial := &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: ialTokens}
		defer tree.Context.TipAppendChild(ial)
	}

	switch dataType {
	case ast.NodeBlockQueryEmbed:
		node.Type = ast.NodeBlockQueryEmbed
		node.AppendChild(&ast.Node{Type: ast.NodeOpenBrace})
		node.AppendChild(&ast.Node{Type: ast.NodeOpenBrace})
		content := util.DomAttrValue(n, "data-content")
		// 嵌入块中存在换行 SQL 语句时会被转换为段落文本 https://github.com/siyuan-note/siyuan/issues/5728
		content = strings.ReplaceAll(content, "\n", editor.IALValEscNewLine)
		node.AppendChild(&ast.Node{Type: ast.NodeBlockQueryEmbedScript, Tokens: util.StrToBytes(content)})
		node.AppendChild(&ast.Node{Type: ast.NodeCloseBrace})
		node.AppendChild(&ast.Node{Type: ast.NodeCloseBrace})
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeTable:
		node.Type = ast.NodeTable
		var tableAligns []int
		if nil == n.FirstChild {
			node.Type = ast.NodeParagraph
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			tree.Context.ParentTip()
			return
		}

		if lute.parentIs(n, atom.Table) {
			text := util.DomText(n)
			node.Tokens = []byte(strings.TrimSpace(text))
			tree.Context.Tip.AppendChild(node)
			return
		}

		tableDiv := n.FirstChild
		table := lute.domChild(tableDiv, atom.Table)
		if nil == table {
			node.Type = ast.NodeParagraph
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			tree.Context.ParentTip()
			return
		}

		thead := lute.domChild(table, atom.Thead)
		if nil == thead || nil == thead.FirstChild || nil == thead.FirstChild.FirstChild {
			node.Type = ast.NodeParagraph
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			tree.Context.ParentTip()
			return
		}
		for th := thead.FirstChild.FirstChild; nil != th; th = th.NextSibling {
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
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()

		lute.genASTContenteditable(table, tree)
		return
	case ast.NodeParagraph:
		node.Type = ast.NodeParagraph
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeHeading:
		text := util.DomText(n)
		if lute.parentIs(n, atom.Table) {
			node.Tokens = []byte(strings.TrimSpace(text))
			for bytes.HasPrefix(node.Tokens, []byte("#")) {
				node.Tokens = bytes.TrimPrefix(node.Tokens, []byte("#"))
			}
			tree.Context.Tip.AppendChild(node)
			return
		}

		level := util.DomAttrValue(n, "data-subtype")[1:]
		tmp := strings.TrimPrefix(text, " ")
		if strings.HasPrefix(tmp, "#") {
			// Allow changing headings with `#` https://github.com/siyuan-note/siyuan/issues/7924
			if idx := strings.Index(tmp, " "+editor.Caret); 0 < idx {
				tmp = tmp[:idx]
				if nil != n.FirstChild && nil != n.FirstChild.FirstChild {
					headingContent := strings.TrimPrefix(strings.TrimPrefix(n.FirstChild.FirstChild.Data, tmp), " ")
					n.FirstChild.FirstChild.Data = headingContent
				}
				level = fmt.Sprintf("%d", strings.Count(tmp, "#"))
			}
		}

		node.Type = ast.NodeHeading
		node.HeadingLevel, _ = strconv.Atoi(level)
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeBlockquote:
		content := strings.TrimSpace(util.DomText(n))
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
	case ast.NodeList:
		node.Type = ast.NodeList
		marker := util.DomAttrValue(n, "data-marker")
		node.ListData = &ast.ListData{}
		subType := util.DomAttrValue(n, "data-subtype")
		if "u" == subType {
			node.ListData.Typ = 0
		} else if "o" == subType {
			node.ListData.Typ = 1
		} else if "t" == subType {
			node.ListData.Typ = 3
		}
		node.ListData.Marker = []byte(marker)
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeListItem:
		marker := util.DomAttrValue(n, "data-marker")
		if ast.NodeList != tree.Context.Tip.Type {
			parent := &ast.Node{}
			parent.Type = ast.NodeList
			parent.ListData = &ast.ListData{}
			subType := util.DomAttrValue(n, "data-subtype")
			if "u" == subType {
				parent.ListData.Typ = 0
				parent.ListData.BulletChar = '*'
			} else if "o" == subType {
				parent.ListData.Typ = 1
				parent.ListData.Num, _ = strconv.Atoi(marker[:len(marker)-1])
				parent.ListData.Delimiter = '.'
			} else if "t" == subType {
				parent.ListData.Typ = 3
				parent.ListData.BulletChar = '*'
			}
			tree.Context.Tip.AppendChild(parent)
			tree.Context.Tip = parent
		}

		node.Type = ast.NodeListItem
		node.ListData = &ast.ListData{}
		subType := util.DomAttrValue(n, "data-subtype")
		if "u" == subType {
			node.ListData.Typ = 0
			node.ListData.BulletChar = '*'
		} else if "o" == subType {
			node.ListData.Typ = 1
			node.ListData.Num, _ = strconv.Atoi(marker[:len(marker)-1])
			node.ListData.Delimiter = '.'
		} else if "t" == subType {
			node.ListData.Typ = 3
			node.ListData.BulletChar = '*'
		}
		node.ListData.Marker = []byte(marker)
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeGitConflict:
		node.Type = ast.NodeGitConflict
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeSuperBlock:
		node.Type = ast.NodeSuperBlock
		tree.Context.Tip.AppendChild(node)
		node.AppendChild(&ast.Node{Type: ast.NodeSuperBlockOpenMarker})
		layout := util.DomAttrValue(n, "data-sb-layout")
		node.AppendChild(&ast.Node{Type: ast.NodeSuperBlockLayoutMarker, Tokens: []byte(layout)})
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeMathBlock:
		node.Type = ast.NodeMathBlock
		node.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker})
		content := util.DomAttrValue(n, "data-content")
		content = html.UnescapeHTMLStr(content)
		node.AppendChild(&ast.Node{Type: ast.NodeMathBlockContent, Tokens: util.StrToBytes(content)})
		node.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker})
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeCodeBlock:
		node.Type = ast.NodeCodeBlock
		node.IsFencedCodeBlock = true
		node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: util.StrToBytes("```")})
		if language := util.DomAttrValue(n, "data-subtype"); "" != language {
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker, CodeBlockInfo: util.StrToBytes(language)})
			content := util.DomAttrValue(n, "data-content")
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockCode, Tokens: util.StrToBytes(content)})
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: util.StrToBytes("```")})
			tree.Context.Tip.AppendChild(node)
			return
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeHTMLBlock:
		node.Type = ast.NodeHTMLBlock
		content := util.DomAttrValue(n.FirstChild.NextSibling.FirstChild, "data-content")
		content = html.UnescapeHTMLStr(content)
		node.Tokens = util.StrToBytes(content)
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeYamlFrontMatter:
		node.Type = ast.NodeYamlFrontMatter
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeThematicBreak:
		node.Type = ast.NodeThematicBreak
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeIFrame:
		node.Type = ast.NodeIFrame
		n = lute.domChild(n.FirstChild, atom.Iframe)
		node.Tokens = util.DomHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeWidget:
		node.Type = ast.NodeWidget
		n = lute.domChild(n.FirstChild, atom.Iframe)
		node.Tokens = util.DomHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeVideo:
		node.Type = ast.NodeVideo
		n = lute.domChild(n.FirstChild, atom.Video)
		node.Tokens = util.DomHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeAudio:
		node.Type = ast.NodeAudio
		n = lute.domChild(n.FirstChild, atom.Audio)
		node.Tokens = util.DomHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeAttributeView:
		node.Type = ast.NodeAttributeView
		node.AttributeViewID = util.DomAttrValue(n, "data-av-id")
		if "" == node.AttributeViewID {
			node.AttributeViewID = ast.NewNodeID()
		}
		node.AttributeViewType = util.DomAttrValue(n, "data-av-type")
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeCustomBlock:
		node.Type = ast.NodeCustomBlock
		node.CustomBlockInfo = util.DomAttrValue(n, "data-info")
		node.Tokens = []byte(html.UnescapeHTMLStr(util.DomAttrValue(n, "data-content")))
		tree.Context.Tip.AppendChild(node)
		return
	default:
		switch n.DataAtom {
		case 0:
			node.Type = ast.NodeText
			node.Tokens = util.StrToBytes(n.Data)
			if ast.NodeDocument == tree.Context.Tip.Type {
				p := &ast.Node{Type: ast.NodeParagraph}
				tree.Context.Tip.AppendChild(p)
				tree.Context.Tip = p
			}
			lute.genASTContenteditable(n, tree)
			return
		case atom.U, atom.Code, atom.Strong, atom.Em, atom.Kbd, atom.Mark, atom.S, atom.Sub, atom.Sup, atom.Span:
			if ast.NodeDocument == tree.Context.Tip.Type {
				p := &ast.Node{Type: ast.NodeParagraph}
				tree.Context.Tip.AppendChild(p)
				tree.Context.Tip = p
			}
			lute.genASTContenteditable(n, tree)
			return
		}

		if ast.NodeListItem == tree.Context.Tip.Type && atom.Input == n.DataAtom {
			node.Type = ast.NodeTaskListItemMarker
			node.TaskListItemChecked = lute.hasAttr(n, "checked")
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeInlineHTML
		node.Tokens = util.DomHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTByBlockDOM(c, tree)
	}

	switch dataType {
	case ast.NodeSuperBlock:
		node.AppendChild(&ast.Node{Type: ast.NodeSuperBlockCloseMarker})
	case ast.NodeCodeBlock:
		node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: util.StrToBytes("```")})
	}
}

func (lute *Lute) genASTContenteditable(n *html.Node, tree *parse.Tree) {
	if ast.NodeCodeBlock == tree.Context.Tip.Type || ast.NodeCustomBlock == tree.Context.Tip.Type {
		return
	}

	if atom.Colgroup == n.DataAtom {
		return
	}

	class := util.DomAttrValue(n, "class")
	if "svg" == class {
		return
	}

	content := n.Data
	node := &ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(content)}
	switch n.DataAtom {
	case 0:
		if "" == content {
			return
		}

		if html.ElementNode == n.Type {
			node.Tokens = util.StrToBytes("<" + content + ">")
		}

		if ast.NodeLink == tree.Context.Tip.Type {
			node.Type = ast.NodeLinkText
		} else if ast.NodeHeading == tree.Context.Tip.Type {
			content = strings.ReplaceAll(content, "\n", "")
			node.Tokens = util.StrToBytes(content)
		} else if ast.NodeStrong == tree.Context.Tip.Type {
			content = strings.ReplaceAll(content, "**", "")
			content = strings.ReplaceAll(content, "*"+editor.Caret, editor.Caret)
			content = strings.ReplaceAll(content, editor.Caret+"*", editor.Caret)
			node.Tokens = util.StrToBytes(content)
		}

		if lute.parentIs(n, atom.Table) {
			content = strings.TrimSuffix(content, "\n")
			if (nil == n.NextSibling && !strings.Contains(content, "\n")) /* 外部内容粘贴到表格中后编辑导致换行丢失  https://github.com/siyuan-note/siyuan/issues/7501 */ ||
				(nil != n.NextSibling && atom.Br == n.NextSibling.DataAtom && strings.HasPrefix(content, "\n")) /* 表格内存在行级公式时编辑会产生换行 https://github.com/siyuan-note/siyuan/issues/2279 */ {
				content = strings.ReplaceAll(content, "\n", "")
			}

			if strings.Contains(content, "\\") {
				// After entering `\` in the table, the next column is merged incorrectly https://github.com/siyuan-note/siyuan/issues/7817
				tmp := strings.ReplaceAll(content, "\\", "")
				tmp = strings.TrimSpace(tmp)
				if "" == tmp {
					// 仅包含转义字符时转义自身 \
					content = strings.ReplaceAll(content, "\\", "\\\\")
				}
			}

			node.Tokens = util.StrToBytes(strings.ReplaceAll(content, "\n", "<br />"))
			array := lex.SplitWithoutBackslashEscape(node.Tokens, '|')
			node.Tokens = nil
			for i, tokens := range array {
				node.Tokens = append(node.Tokens, tokens...)
				if i < len(array)-1 {
					node.Tokens = append(node.Tokens, []byte("\\|")...)
				}
			}
		}
		if ast.NodeCodeSpan == tree.Context.Tip.Type || ast.NodeInlineMath == tree.Context.Tip.Type {
			if nil != tree.Context.Tip.Previous && tree.Context.Tip.Type == tree.Context.Tip.Previous.Type { // 合并相邻的代码
				tree.Context.Tip.FirstChild.Next.Tokens = util.StrToBytes(content)
			} else { // 叠加代码
				if nil != tree.Context.Tip.FirstChild.Next.Next && ast.NodeBackslash == tree.Context.Tip.FirstChild.Next.Next.Type {
					// 表格单元格中使用代码和 `|` 的问题 https://github.com/siyuan-note/siyuan/issues/4717
					content = util.BytesToStr(tree.Context.Tip.FirstChild.Next.Next.FirstChild.Tokens) + content
					tree.Context.Tip.FirstChild.Next.Next.Unlink()
				}
				tree.Context.Tip.FirstChild.Next.Tokens = append(tree.Context.Tip.FirstChild.Next.Tokens, util.StrToBytes(content)...)
			}
			return
		}
		if ast.NodeTextMark == tree.Context.Tip.Type {
			if "code" == tree.Context.Tip.TokensStr() {
				if nil != tree.Context.Tip.FirstChild && nil != tree.Context.Tip.FirstChild.Next && nil != tree.Context.Tip.FirstChild.Next.Next && ast.NodeBackslash == tree.Context.Tip.FirstChild.Next.Next.Type {
					// 表格单元格中使用代码和 `|` 的问题 https://github.com/siyuan-note/siyuan/issues/4717
					content = util.BytesToStr(tree.Context.Tip.FirstChild.Next.Next.FirstChild.Tokens) + content
					tree.Context.Tip.FirstChild.Next.Next.Unlink()
					tree.Context.Tip.FirstChild.Next.Tokens = append(tree.Context.Tip.FirstChild.Next.Tokens, util.StrToBytes(content)...)
					return
				}
			}
		}

		if ast.NodeKbd == tree.Context.Tip.Type {
			// `<kbd>` 中反斜杠转义问题 https://github.com/siyuan-note/siyuan/issues/2242
			node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("\\\\"), []byte("\\"))
			node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("\\"), []byte("\\\\"))

			if bytes.Equal(node.Tokens, editor.CaretTokens) {
				// `<kbd>` 无法删除 https://github.com/siyuan-note/siyuan/issues/4162
				parent := tree.Context.Tip.Parent
				tree.Context.Tip.Unlink()
				tree.Context.Tip = parent
			}
		}
		tree.Context.Tip.AppendChild(node)
	case atom.Thead:
		if lute.parentIs(n.Parent.Parent, atom.Table) {
			text := util.DomText(n.Parent.Parent)
			text = strings.ReplaceAll(text, editor.Caret, "")
			node.Tokens = []byte(strings.TrimSpace(text))
			tree.Context.Tip.AppendChild(node)
			return
		}

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
		tree.Context.Tip.AppendChild(node)
		parse.SetSpanIAL(node, n)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Code:
		isCaret, isEmpty := lute.isCaret(n)
		if isCaret {
			node.Type = ast.NodeText
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		if isEmpty {
			return
		}

		if lute.ParseOptions.TextMark {
			tree.Context.Tip.AppendChild(node)
			parse.SetTextMarkNode(node, n, lute.ParseOptions)
			return
		}

		node.Type = ast.NodeCodeSpan
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanOpenMarker})
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanContent})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Span:
		dataType := util.DomAttrValue(n, "data-type")
		if "" == dataType {
			dataType = "text"
		}
		if strings.Contains(dataType, "span") {
			// 某些情况下复制过来的 DOM 是该情况，这里按纯文本解析
			node.Type = ast.NodeText
			node.Tokens = util.StrToBytes(util.DomText(n))
			tree.Context.Tip.AppendChild(node)
			return
		}

		if strings.Contains(dataType, "img") {
			// 给文字和图片同时设置字体格式后图片丢失 https://github.com/siyuan-note/siyuan/issues/6297
			dataType = "img"
		}

		if nil != tree.Context.Tip && nil != tree.Context.Tip.LastChild {
			// 行级元素前输入转义符 `\` 导致异常 https://github.com/siyuan-note/siyuan/issues/6237
			previousEndText := tree.Context.Tip.LastChild.Text()
			backslashCaret := strings.HasSuffix(previousEndText, "\\"+editor.Caret)
			if backslashCaret {
				previousEndText = strings.TrimSuffix(previousEndText, editor.Caret)
			}
			if strings.HasSuffix(previousEndText, "\\") {
				backslashCount := 0
				for i := len(previousEndText) - 1; i >= 0; i-- {
					if '\\' == previousEndText[i] {
						backslashCount++
					} else {
						break
					}
				}
				if 0 != backslashCount%2 {
					if backslashCaret {
						tree.Context.Tip.LastChild.Tokens = bytes.TrimSuffix(tree.Context.Tip.LastChild.Tokens, []byte(editor.Caret))
						tree.Context.Tip.LastChild.Tokens = append(tree.Context.Tip.LastChild.Tokens, []byte("\\")...)
						tree.Context.Tip.LastChild.Tokens = append(tree.Context.Tip.LastChild.Tokens, []byte(editor.Caret)...)
					} else {
						tree.Context.Tip.LastChild.Tokens = append(tree.Context.Tip.LastChild.Tokens, []byte("\\")...)
					}
				}
			}
		}

		if "tag" == dataType {
			isCaret, isEmpty := lute.isCaret(n)
			if isCaret {
				node.Type = ast.NodeText
				node.Tokens = editor.CaretTokens
				tree.Context.Tip.AppendChild(node)
				return
			}
			if isEmpty {
				return
			}

			if lute.ParseOptions.TextMark {
				tree.Context.Tip.AppendChild(node)
				parse.SetTextMarkNode(node, n, lute.ParseOptions)
				return
			}

			n.FirstChild.Data = strings.ReplaceAll(n.FirstChild.Data, editor.Zwsp, "")
			node.Type = ast.NodeTag
			node.AppendChild(&ast.Node{Type: ast.NodeTagOpenMarker})
			// 开头结尾空格后会形成 * foo * 导致强调、加粗删除线标记失效，这里将空格移到右标记符前后 _*foo*_
			processSpanMarkerSpace(n, node)
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		} else if "inline-math" == dataType {
			inlineMathContent := util.GetTextMarkInlineMathData(n)
			if "" == inlineMathContent {
				return
			}

			if lute.ParseOptions.TextMark {
				tree.Context.Tip.AppendChild(node)
				parse.SetTextMarkNode(node, n, lute.ParseOptions)
				return
			}

			node.Type = ast.NodeInlineMath
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathOpenMarker})
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathContent, Tokens: util.StrToBytes(inlineMathContent)})
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathCloseMarker})
			tree.Context.Tip.AppendChild(node)
			return
		} else if "inline-memo" == dataType {
			isCaret, isEmpty := lute.isCaret(n)
			if isCaret {
				node.Type = ast.NodeText
				node.Tokens = editor.CaretTokens
				tree.Context.Tip.AppendChild(node)
				return
			}
			if isEmpty {
				return
			}

			if lute.ParseOptions.TextMark {
				tree.Context.Tip.AppendChild(node)
				parse.SetTextMarkNode(node, n, lute.ParseOptions)
				return
			}

			node.Type = ast.NodeText
			node.Tokens = util.StrToBytes(util.DomText(n))
			tree.Context.Tip.AppendChild(node)
			return
		} else if "a" == dataType {
			if nil == n.FirstChild {
				// 丢弃没有锚文本的链接
				return
			}

			if ast.NodeLink == tree.Context.Tip.Type {
				break
			}

			if lute.ParseOptions.TextMark {
				tree.Context.Tip.AppendChild(node)
				parse.SetTextMarkNode(node, n, lute.ParseOptions)
				return
			}

			node.Type = ast.NodeLink
			node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		} else if "block-ref" == dataType {
			refText := util.DomText(n)
			refText = strings.TrimSpace(refText)
			if "" == refText {
				return
			}
			if refText == editor.Caret {
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: editor.CaretTokens})
				return
			}

			if lute.ParseOptions.TextMark {
				tree.Context.Tip.AppendChild(node)
				parse.SetTextMarkNode(node, n, lute.ParseOptions)
				return
			}

			node.Type = ast.NodeBlockRef
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			id := util.DomAttrValue(n, "data-id")
			node.AppendChild(&ast.Node{Type: ast.NodeBlockRefID, Tokens: util.StrToBytes(id)})

			node.AppendChild(&ast.Node{Type: ast.NodeBlockRefSpace})
			var refTextNode *ast.Node
			subtype := util.DomAttrValue(n, "data-subtype")
			if "s" == subtype || "" == subtype {
				refTextNode = &ast.Node{Type: ast.NodeBlockRefText, Tokens: util.StrToBytes(refText)}
			} else {
				refTextNode = &ast.Node{Type: ast.NodeBlockRefDynamicText, Tokens: util.StrToBytes(refText)}
			}
			if lute.parentIs(n, atom.Table) {
				refTextNode.Tokens = bytes.ReplaceAll(refTextNode.Tokens, []byte("|"), []byte("&#124;"))
			}
			node.AppendChild(refTextNode)
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			tree.Context.Tip.AppendChild(node)
			return
		} else if "file-annotation-ref" == dataType {
			refText := util.DomText(n)
			refText = strings.TrimSpace(refText)
			if "" == refText {
				return
			}
			if refText == editor.Caret {
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: editor.CaretTokens})
				return
			}

			if lute.ParseOptions.TextMark {
				tree.Context.Tip.AppendChild(node)
				parse.SetTextMarkNode(node, n, lute.ParseOptions)
				return
			}

			node.Type = ast.NodeFileAnnotationRef
			node.AppendChild(&ast.Node{Type: ast.NodeLess})
			node.AppendChild(&ast.Node{Type: ast.NodeLess})
			id := util.DomAttrValue(n, "data-id")
			node.AppendChild(&ast.Node{Type: ast.NodeFileAnnotationRefID, Tokens: util.StrToBytes(id)})
			node.AppendChild(&ast.Node{Type: ast.NodeFileAnnotationRefSpace})
			refTextNode := &ast.Node{Type: ast.NodeFileAnnotationRefText, Tokens: util.StrToBytes(refText)}
			node.AppendChild(refTextNode)
			node.AppendChild(&ast.Node{Type: ast.NodeGreater})
			node.AppendChild(&ast.Node{Type: ast.NodeGreater})
			tree.Context.Tip.AppendChild(node)
			return
		} else if "img" == dataType {
			img := lute.domChild(n, atom.Img) //n.FirstChild.NextSibling.FirstChild.NextSibling
			if nil == img {
				return
			}

			node.Type = ast.NodeImage
			node.AppendChild(&ast.Node{Type: ast.NodeBang})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
			alt := util.DomAttrValue(img, "alt")
			node.AppendChild(&ast.Node{Type: ast.NodeLinkText, Tokens: util.StrToBytes(alt)})
			node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			src := strings.TrimSpace(util.DomAttrValue(img, "data-src"))
			node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: util.StrToBytes(src)})
			if title := util.DomAttrValue(img, "title"); "" != title {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
				node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: util.StrToBytes(title)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			tree.Context.Tip.AppendChild(node)
			parse.SetSpanIAL(tree.Context.Tip.LastChild, img)
			return
		} else if "backslash" == dataType {
			node.Type = ast.NodeBackslash
			if nil == n.FirstChild {
				return
			}
			if n.FirstChild == n.LastChild && nil != n.FirstChild.FirstChild {
				// 转义字符加行级样式后继续输入会出现标记符 https://github.com/siyuan-note/siyuan/issues/6134
				return
			}
			if nil == n.FirstChild.NextSibling && html.TextNode == n.FirstChild.Type {
				node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(n.FirstChild.Data)})
				tree.Context.Tip.AppendChild(node)
				return
			}
			if nil != n.FirstChild.NextSibling {
				data := n.FirstChild.NextSibling.Data
				data = strings.ReplaceAll(data, "\\\\", "\\")
				node.AppendChild(&ast.Node{Type: ast.NodeBackslashContent, Tokens: util.StrToBytes(data)})
			}
			tree.Context.Tip.AppendChild(node)
			return
		} else {
			// TextMark 节点
			isCaret, isEmpty := lute.isCaret(n)
			if isCaret {
				node.Type = ast.NodeText
				node.Tokens = editor.CaretTokens
				tree.Context.Tip.AppendChild(node)
				return
			}
			if isEmpty {
				return
			}

			dataType = lute.removeTempMark(dataType)
			tmpDataType := strings.ReplaceAll(dataType, "backslash", "")
			tmpDataType = strings.TrimSpace(tmpDataType)
			tree.Context.Tip.AppendChild(node)
			if "" == tmpDataType {
				node.Type = ast.NodeText
				node.Tokens = []byte(util.DomText(n))
				return
			}
			lute.setDOMAttrValue(n, "data-type", dataType)
			parse.SetTextMarkNode(node, n, lute.ParseOptions)
			return
		}
	case atom.Sub:
		isCaret, isEmpty := lute.isCaret(n)
		if isCaret {
			node.Type = ast.NodeText
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		if isEmpty {
			return
		}

		if lute.ParseOptions.TextMark {
			tree.Context.Tip.AppendChild(node)
			parse.SetTextMarkNode(node, n, lute.ParseOptions)
			return
		}

		node.Type = ast.NodeSub
		node.AppendChild(&ast.Node{Type: ast.NodeSubOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Sup:
		isCaret, isEmpty := lute.isCaret(n)
		if isCaret {
			node.Type = ast.NodeText
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		if isEmpty {
			return
		}

		if lute.ParseOptions.TextMark {
			tree.Context.Tip.AppendChild(node)
			parse.SetTextMarkNode(node, n, lute.ParseOptions)
			return
		}

		node.Type = ast.NodeSup
		node.AppendChild(&ast.Node{Type: ast.NodeSupOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.U:
		isCaret, isEmpty := lute.isCaret(n)
		if isCaret {
			node.Type = ast.NodeText
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		if isEmpty {
			return
		}

		if lute.ParseOptions.TextMark {
			tree.Context.Tip.AppendChild(node)
			parse.SetTextMarkNode(node, n, lute.ParseOptions)
			return
		}

		node.Type = ast.NodeUnderline
		node.AppendChild(&ast.Node{Type: ast.NodeUnderlineOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Kbd:
		isCaret, isEmpty := lute.isCaret(n)
		if isCaret {
			node.Type = ast.NodeText
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		if isEmpty {
			return
		}

		if lute.ParseOptions.TextMark {
			tree.Context.Tip.AppendChild(node)
			parse.SetTextMarkNode(node, n, lute.ParseOptions)
			return
		}

		node.Type = ast.NodeKbd
		node.AppendChild(&ast.Node{Type: ast.NodeKbdOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Br:
		if ast.NodeHeading == tree.Context.Tip.Type {
			return
		}
		if nil != n.PrevSibling && "\n" == n.PrevSibling.Data && lute.parentIs(n, atom.Table) {
			return
		}

		node.Type = ast.NodeBr
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Em, atom.I:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, editor.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(editor.Zwsp + "\n")})
		}
		isCaret, isEmpty := lute.isCaret(n)
		if isCaret {
			node.Type = ast.NodeText
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		if isEmpty {
			return
		}

		if lute.ParseOptions.TextMark {
			tree.Context.Tip.AppendChild(node)
			parse.SetTextMarkNode(node, n, lute.ParseOptions)
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
		processSpanMarkerSpace(n, node)
		lute.removeInnerMarker(n, "__")
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Strong, atom.B:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if nil != tree.Context.Tip.LastChild {
			// 转义符导致的行级元素样式属性暴露 https://github.com/siyuan-note/siyuan/issues/2969
			if bytes.HasSuffix(tree.Context.Tip.LastChild.Tokens, []byte("\\"+editor.Caret)) {
				// foo\‸**bar**
				tree.Context.Tip.LastChild.Tokens = bytes.ReplaceAll(tree.Context.Tip.LastChild.Tokens, []byte("\\"+editor.Caret), []byte("\\\\"+editor.Caret))
			}
			if bytes.HasSuffix(tree.Context.Tip.LastChild.Tokens, []byte("\\")) {
				// foo\**bar**
				tree.Context.Tip.LastChild.Tokens = bytes.ReplaceAll(tree.Context.Tip.LastChild.Tokens, []byte("\\"), []byte("\\\\"))
			}
		}

		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, editor.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(editor.Zwsp + "\n")})
		}
		isCaret, isEmpty := lute.isCaret(n)
		if isCaret {
			node.Type = ast.NodeText
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		if isEmpty {
			return
		}

		if lute.ParseOptions.TextMark {
			tree.Context.Tip.AppendChild(node)
			parse.SetTextMarkNode(node, n, lute.ParseOptions)
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

		processSpanMarkerSpace(n, node)
		lute.removeInnerMarker(n, "**")
		parse.SetSpanIAL(node, n)
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
		isCaret, isEmpty := lute.isCaret(n)
		if isCaret {
			node.Type = ast.NodeText
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		if isEmpty {
			return
		}

		if lute.ParseOptions.TextMark {
			tree.Context.Tip.AppendChild(node)
			parse.SetTextMarkNode(node, n, lute.ParseOptions)
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

		processSpanMarkerSpace(n, node)
		lute.removeInnerMarker(n, "~~")
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
		isCaret, isEmpty := lute.isCaret(n)
		if isCaret {
			node.Type = ast.NodeText
			node.Tokens = editor.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}
		if isEmpty {
			return
		}

		if lute.ParseOptions.TextMark {
			tree.Context.Tip.AppendChild(node)
			parse.SetTextMarkNode(node, n, lute.ParseOptions)
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

		processSpanMarkerSpace(n, node)
		lute.removeInnerMarker(n, "==")
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Img:
		if "emoji" == class {
			alt := util.DomAttrValue(n, "alt")
			node.Type = ast.NodeEmoji
			emojiImg := &ast.Node{Type: ast.NodeEmojiImg, Tokens: tree.EmojiImgTokens(alt, util.DomAttrValue(n, "src"))}
			emojiImg.AppendChild(&ast.Node{Type: ast.NodeEmojiAlias, Tokens: []byte(":" + alt + ":")})
			node.AppendChild(emojiImg)
			tree.Context.Tip.AppendChild(node)
			return
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTContenteditable(c, tree)
	}

	if lute.ParseOptions.TextMark {
		return
	}

	switch n.DataAtom {
	case atom.Code:
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanCloseMarker})
	case atom.Span:
		dataType := util.DomAttrValue(n, "data-type")
		if "tag" == dataType {
			node.AppendChild(&ast.Node{Type: ast.NodeTagCloseMarker})
		} else if "a" == dataType {
			node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			href := util.DomAttrValue(n, "data-href")
			if "" != lute.RenderOptions.LinkBase {
				href = strings.ReplaceAll(href, lute.RenderOptions.LinkBase, "")
			}
			if "" != lute.RenderOptions.LinkPrefix {
				href = strings.ReplaceAll(href, lute.RenderOptions.LinkPrefix, "")
			}
			node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: []byte(href)})
			linkTitle := util.DomAttrValue(n, "data-title")
			if "" != linkTitle {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
				node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: []byte(linkTitle)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
		}
	case atom.Sub:
		node.AppendChild(&ast.Node{Type: ast.NodeSubCloseMarker})
	case atom.Sup:
		node.AppendChild(&ast.Node{Type: ast.NodeSupCloseMarker})
	case atom.U:
		node.AppendChild(&ast.Node{Type: ast.NodeUnderlineCloseMarker})
	case atom.Kbd:
		node.AppendChild(&ast.Node{Type: ast.NodeKbdCloseMarker})
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
	}
}

func (lute *Lute) setBlockIAL(n *html.Node, node *ast.Node) (ialTokens []byte) {
	node.SetIALAttr("id", node.ID)

	if icon := util.DomAttrValue(n, "icon"); "" != icon {
		node.SetIALAttr("icon", icon)
		ialTokens = append(ialTokens, []byte(" icon=\""+icon+"\"")...)
	}

	if refcount := util.DomAttrValue(n, "refcount"); "" != refcount {
		node.SetIALAttr("refcount", refcount)
		ialTokens = append(ialTokens, []byte(" refcount=\""+refcount+"\"")...)
	}

	if avNames := util.DomAttrValue(n, "av-names"); "" != avNames {
		node.SetIALAttr("av-names", avNames)
		ialTokens = append(ialTokens, []byte(" av-names=\""+avNames+"\"")...)
	}

	if bookmark := util.DomAttrValue(n, "bookmark"); "" != bookmark {
		bookmark = html.UnescapeHTMLStr(bookmark)
		node.SetIALAttr("bookmark", bookmark)
		ialTokens = append(ialTokens, []byte(" bookmark=\""+bookmark+"\"")...)
	}

	if style := util.DomAttrValue(n, "style"); "" != style {
		style = html.UnescapeHTMLStr(style)
		style = parse.StyleValue(style)
		node.SetIALAttr("style", style)
		ialTokens = append(ialTokens, []byte(" style=\""+style+"\"")...)
	}

	if name := util.DomAttrValue(n, "name"); "" != name {
		name = html.UnescapeHTMLStr(name)
		node.SetIALAttr("name", name)
		ialTokens = append(ialTokens, []byte(" name=\""+name+"\"")...)
	}

	if memo := util.DomAttrValue(n, "memo"); "" != memo {
		memo = html.UnescapeHTMLStr(memo)
		node.SetIALAttr("memo", memo)
		ialTokens = append(ialTokens, []byte(" memo=\""+memo+"\"")...)
	}

	if alias := util.DomAttrValue(n, "alias"); "" != alias {
		alias = html.UnescapeHTMLStr(alias)
		node.SetIALAttr("alias", alias)
		ialTokens = append(ialTokens, []byte(" alias=\""+alias+"\"")...)
	}

	if fold := util.DomAttrValue(n, "fold"); "" != fold {
		node.SetIALAttr("fold", fold)
		ialTokens = append(ialTokens, []byte(" fold=\""+fold+"\"")...)
	}

	if headingFold := util.DomAttrValue(n, "heading-fold"); "" != headingFold {
		node.SetIALAttr("heading-fold", headingFold)
		ialTokens = append(ialTokens, []byte(" heading-fold=\""+headingFold+"\"")...)
	}

	if parentFold := util.DomAttrValue(n, "parent-fold"); "" != parentFold {
		node.SetIALAttr("parent-fold", parentFold)
		ialTokens = append(ialTokens, []byte(" parent-fold=\""+parentFold+"\"")...)
	}

	if updated := util.DomAttrValue(n, "updated"); "" != updated {
		node.SetIALAttr("updated", updated)
		ialTokens = append(ialTokens, []byte(" updated=\""+updated+"\"")...)
	}

	if linewrap := util.DomAttrValue(n, "linewrap"); "" != linewrap {
		node.SetIALAttr("linewrap", linewrap)
		ialTokens = append(ialTokens, []byte(" linewrap=\""+linewrap+"\"")...)
	}

	if ligatures := util.DomAttrValue(n, "ligatures"); "" != ligatures {
		node.SetIALAttr("ligatures", ligatures)
		ialTokens = append(ialTokens, []byte(" ligatures=\""+ligatures+"\"")...)
	}

	if linenumber := util.DomAttrValue(n, "linenumber"); "" != linenumber {
		node.SetIALAttr("linenumber", linenumber)
		ialTokens = append(ialTokens, []byte(" linenumber=\""+linenumber+"\"")...)
	}

	if breadcrumb := util.DomAttrValue(n, "breadcrumb"); "" != breadcrumb {
		node.SetIALAttr("breadcrumb", breadcrumb)
		ialTokens = append(ialTokens, []byte(" breadcrumb=\""+breadcrumb+"\"")...)
	}

	if dataExportMd := util.DomAttrValue(n, "data-export-md"); "" != dataExportMd {
		dataExportMd = html.UnescapeHTMLStr(dataExportMd)
		node.SetIALAttr("data-export-md", dataExportMd)
		ialTokens = append(ialTokens, []byte(" data-export-md=\""+dataExportMd+"\"")...)
	}

	if dataExportHtml := util.DomAttrValue(n, "data-export-html"); "" != dataExportHtml {
		dataExportHtml = html.UnescapeHTMLStr(dataExportHtml)
		node.SetIALAttr("data-export-html", dataExportHtml)
		ialTokens = append(ialTokens, []byte(" data-export-html=\""+dataExportHtml+"\"")...)
	}

	if customAttrs := util.DomCustomAttrs(n); nil != customAttrs {
		for k, v := range customAttrs {
			v = html.UnescapeHTMLStr(v)
			node.SetIALAttr(k, v)
			ialTokens = append(ialTokens, []byte(" "+k+"=\""+v+"\"")...)
		}
	}

	if "NodeTable" == util.DomAttrValue(n, "data-type") {
		colgroup := lute.domChild(n, atom.Colgroup)
		var colgroupAttrVal string
		if nil != colgroup {
			for col := colgroup.FirstChild; nil != col; col = col.NextSibling {
				colStyle := util.DomAttrValue(col, "style")
				colgroupAttrVal += colStyle
				if nil != col.NextSibling {
					colgroupAttrVal += "|"
				}
			}
			node.SetIALAttr("colgroup", colgroupAttrVal)
			ialTokens = append(ialTokens, []byte(" colgroup=\""+colgroupAttrVal+"\"")...)
		}
	}

	ialTokens = parse.IAL2Tokens(node.KramdownIAL)
	return ialTokens
}

func processSpanMarkerSpace(n *html.Node, node *ast.Node) {
	if strings.HasPrefix(n.FirstChild.Data, " ") && nil == n.FirstChild.PrevSibling {
		n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, " ")
		node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(" ")})
	}
	if strings.HasSuffix(n.FirstChild.Data, " ") && nil == n.FirstChild.NextSibling {
		n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " ")
		n.InsertAfter(&html.Node{Type: html.TextNode, Data: " "})
	}
	if strings.HasSuffix(n.FirstChild.Data, " "+editor.Caret) && nil == n.FirstChild.NextSibling {
		n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, " "+editor.Caret)
		n.InsertAfter(&html.Node{Type: html.TextNode, Data: " " + editor.Caret})
	}
	if strings.HasSuffix(n.FirstChild.Data, "\n") && nil == n.FirstChild.NextSibling {
		n.FirstChild.Data = strings.TrimRight(n.FirstChild.Data, "\n")
		n.InsertAfter(&html.Node{Type: html.TextNode, Data: "\n"})
	}
}

func (lute *Lute) removeInnerMarker(n *html.Node, marker string) {
	if html.TextNode == n.Type {
		n.Data = strings.ReplaceAll(n.Data, marker, "")
	}
	for child := n.FirstChild; nil != child; child = child.NextSibling {
		lute.removeInnerMarker0(child, marker)
	}
}

func (lute *Lute) removeInnerMarker0(n *html.Node, marker string) {
	if nil == n {
		return
	}
	if dataRender := util.DomAttrValue(n, "data-render"); "1" == dataRender || "2" == dataRender {
		return
	}

	if "svg" == n.Namespace {
		return
	}

	if 0 == n.DataAtom && html.ElementNode == n.Type { // 自定义标签
		return
	}

	switch n.DataAtom {
	case 0:
		n.Data = strings.ReplaceAll(n.Data, marker, "")
	}

	for child := n.FirstChild; nil != child; child = child.NextSibling {
		lute.removeInnerMarker0(child, marker)
	}
}

func (lute *Lute) removeTempMark(dataType string) (ret string) {
	ret = strings.ReplaceAll(dataType, "search-mark", "")
	ret = strings.ReplaceAll(ret, "virtual-block-ref", "")
	ret = strings.TrimSpace(ret)
	return
}
