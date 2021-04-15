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
	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

func (lute *Lute) OL2UL(ivHTML string) (ovHTML string) {
	tree, err := lute.VditorBlockDOM2Tree(ivHTML)
	if nil != err {
		return err.Error()
	}

	if ast.NodeList != tree.Root.FirstChild.Type {
		return ivHTML
	}

	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering || !n.IsBlock() || (ast.NodeList != n.Type && ast.NodeListItem != n.Type) {
			return ast.WalkContinue
		}

		n.ListData.Typ = 0
		return ast.WalkContinue
	})

	ovHTML = lute.Tree2VditorBlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) UL2OL(ivHTML string) (ovHTML string) {
	tree, err := lute.VditorBlockDOM2Tree(ivHTML)
	if nil != err {
		return err.Error()
	}

	if ast.NodeList != tree.Root.FirstChild.Type {
		return ivHTML
	}

	var num int
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering || !n.IsBlock() || (ast.NodeList != n.Type && ast.NodeListItem != n.Type) {
			return ast.WalkContinue
		}

		if ast.NodeList == n.Type {
			num = 0
		} else {
			num++
		}

		n.ListData.Typ = 1
		n.ListData.Num = num
		return ast.WalkContinue
	})

	ovHTML = lute.Tree2VditorBlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) SpinVditorBlockDOM(ivHTML string) (ovHTML string) {
	// 替换插入符
	ivHTML = strings.ReplaceAll(ivHTML, "<wbr>", util.Caret)

	markdown := lute.vditorBlockDOM2Md(ivHTML)
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)

	ovHTML = lute.Tree2VditorBlockDOM(tree, lute.RenderOptions)
	// 替换插入符
	ovHTML = strings.ReplaceAll(ovHTML, util.Caret, "<wbr>")
	return
}

func (lute *Lute) HTML2VditorBlockDOM(sHTML string) (vHTML string) {
	//fmt.Println(sHTML)
	markdown, err := lute.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.HTML2VditorBlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

func (lute *Lute) VditorBlockDOM2HTML(vhtml string) (sHTML string) {
	markdown := lute.vditorBlockDOM2Md(vhtml)
	sHTML = lute.Md2HTML(markdown)
	return
}

func (lute *Lute) Md2VditorBlockDOM(markdown string) (vHTML string) {
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2VditorBlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

func (lute *Lute) InlineMd2VditorBlockDOM(markdown string) (vHTML string) {
	tree := parse.Inline("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2VditorBlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

func (lute *Lute) VditorBlockDOM2Md(htmlStr string) (markdown string) {
	//fmt.Println(htmlStr)
	htmlStr = strings.ReplaceAll(htmlStr, parse.Zwsp, "")
	markdown = lute.vditorBlockDOM2Md(htmlStr)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	return
}

func (lute *Lute) VditorBlockDOM2StdMd(htmlStr string) (markdown string) {
	htmlStr = strings.ReplaceAll(htmlStr, parse.Zwsp, "")

	// DOM 转 AST
	tree, err := lute.VditorBlockDOM2Tree(htmlStr)
	if nil != err {
		return err.Error()
	}

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

	// 将 AST 进行 Markdown 格式化渲染
	options := render.NewOptions()
	options.AutoSpace = false
	options.FixTermTypo = false
	options.KramdownBlockIAL = true
	options.KramdownSpanIAL = true
	renderer := render.NewFormatRenderer(tree, options)
	formatted := renderer.Render()
	markdown = string(formatted)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	return
}

func (lute *Lute) VditorBlockDOM2Text(htmlStr string) (text string) {
	tree, err := lute.VditorBlockDOM2Tree(htmlStr)
	if nil != err {
		return ""
	}
	return tree.Root.Text()
}

func (lute *Lute) VditorBlockDOM2TextLen(htmlStr string) int {
	tree, err := lute.VditorBlockDOM2Tree(htmlStr)
	if nil != err {
		return 0
	}
	return tree.Root.TextLen()
}

func (lute *Lute) Tree2VditorBlockDOM(tree *parse.Tree, options *render.Options) (vHTML string) {
	renderer := render.NewVditorBlockRenderer(tree, options)
	output := renderer.Render()
	vHTML = string(output)
	return
}

func RenderNodeVditorBlockDOM(node *ast.Node, parseOptions *parse.Options, renderOptions *render.Options) string {
	root := &ast.Node{Type: ast.NodeDocument}
	tree := &parse.Tree{Root: root, Context: &parse.Context{ParseOption: parseOptions}}
	renderer := render.NewVditorBlockRenderer(tree, renderOptions)
	renderer.Writer = &bytes.Buffer{}
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := renderer.RendererFuncs[n.Type]
		return rendererFunc(n, entering)
	})
	return renderer.Writer.String()
}

func (lute *Lute) VditorBlockDOM2Tree(htmlStr string) (ret *parse.Tree, err error) {
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
	ret = &parse.Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument}, Context: &parse.Context{ParseOption: lute.ParseOptions}}
	ret.Context.Tip = ret.Root
	for c := htmlRoot.FirstChild; nil != c; c = c.NextSibling {
		lute.genASTByVditorBlockDOM(c, ret)
	}

	// 调整树结构
	ast.Walk(ret.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n.Type {
			case ast.NodeInlineHTML, ast.NodeCodeSpan, ast.NodeInlineMath, ast.NodeHTMLBlock, ast.NodeCodeBlockCode, ast.NodeMathBlockContent:
				n.Tokens = html.UnescapeHTML(n.Tokens)
				if nil != n.Next && ast.NodeCodeSpan == n.Next.Type && n.CodeMarkerLen == n.Next.CodeMarkerLen && nil != n.FirstChild && nil != n.FirstChild.Next {
					// 合并代码节点 https://github.com/Vanessa219/vditor/issues/167
					n.FirstChild.Next.Tokens = append(n.FirstChild.Next.Tokens, n.Next.FirstChild.Next.Tokens...)
					n.Next.Unlink()
				}
			}
		}
		return ast.WalkContinue
	})
	return
}

func (lute *Lute) vditorBlockDOM2Md(htmlStr string) (markdown string) {
	tree, err := lute.VditorBlockDOM2Tree(htmlStr)
	if nil != err {
		return err.Error()
	}

	// 将 AST 进行 Markdown 格式化渲染
	options := render.NewOptions()
	options.AutoSpace = false
	options.FixTermTypo = false
	options.KramdownBlockIAL = true
	options.KramdownSpanIAL = true
	renderer := render.NewFormatRenderer(tree, options)
	formatted := renderer.Render()
	markdown = string(formatted)
	return
}

func (lute *Lute) genASTByVditorBlockDOM(n *html.Node, tree *parse.Tree) {
	if class := lute.domAttrValue(n, "class"); strings.Contains(class, "vditor-bullet") || "vditor-attr" == class {
		return
	}

	if "true" == lute.domAttrValue(n, "contenteditable") {
		lute.genASTContenteditable(n, tree)
		return
	}

	dataType := ast.Str2NodeType(lute.domAttrValue(n, "data-type"))

	nodeID := lute.domAttrValue(n, "data-node-id")
	node := &ast.Node{ID: nodeID}
	if "" != node.ID && !lute.parentIs(n, atom.Table) {
		node.KramdownIAL = [][]string{{"id", node.ID}}
		ialTokens := lute.setBlockIAL(n, node)
		ial := &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: ialTokens}
		defer tree.Context.TipAppendChild(ial)
	}

	switch dataType {
	case ast.NodeParagraph:
		node.Type = ast.NodeParagraph
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeHeading:
		text := lute.domText(n)
		if "" == strings.TrimSpace(text) {
			return
		}
		if lute.parentIs(n, atom.Table) {
			node.Tokens = []byte(strings.TrimSpace(text))
			for bytes.HasPrefix(node.Tokens, []byte("#")) {
				node.Tokens = bytes.TrimPrefix(node.Tokens, []byte("#"))
			}
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeHeading
		level := lute.domAttrValue(n, "data-subtype")[1:]
		node.HeadingLevel, _ = strconv.Atoi(level)
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeBlockquote:
		content := strings.TrimSpace(lute.domText(n))
		if "" == content || "&gt;" == content {
			return
		}
		if util.Caret == content {
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
		marker := lute.domAttrValue(n, "data-marker")
		node.ListData = &ast.ListData{}
		subType := lute.domAttrValue(n, "data-subtype")
		if "u" == subType {
			node.ListData.BulletChar = '*'
			node.ListData.Typ = 0
		} else if "o" == subType {
			node.ListData.Typ = 1
		} else if "t" == subType {
			node.ListData.BulletChar = '*'
			node.ListData.Typ = 3
		}
		node.ListData.Marker = []byte(marker)
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeListItem:
		marker := lute.domAttrValue(n, "data-marker")
		if ast.NodeList != tree.Context.Tip.Type {
			parent := &ast.Node{}
			parent.Type = ast.NodeList
			parent.ListData = &ast.ListData{}
			subType := lute.domAttrValue(n, "data-subtype")
			if "u" == subType {
				node.ListData.BulletChar = '*'
				node.ListData.Typ = 0
			} else if "o" == subType {
				node.ListData.Typ = 1
				node.ListData.Num, _ = strconv.Atoi(marker[:len(marker)-1])
			} else if "t" == subType {
				node.ListData.BulletChar = '*'
				node.ListData.Typ = 3
			}
			tree.Context.Tip.AppendChild(parent)
			tree.Context.Tip = parent
		}

		node.Type = ast.NodeListItem
		node.ListData = &ast.ListData{}
		subType := lute.domAttrValue(n, "data-subtype")
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
	case ast.NodeGitConflict:
		node.Type = ast.NodeGitConflict
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeSuperBlock:
		node.Type = ast.NodeSuperBlock
		tree.Context.Tip.AppendChild(node)
		node.AppendChild(&ast.Node{Type: ast.NodeSuperBlockOpenMarker})
		layout := lute.domAttrValue(n, "data-sb-layout")
		node.AppendChild(&ast.Node{Type: ast.NodeSuperBlockLayoutMarker, Tokens: []byte(layout)})
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeMathBlock:
		node.Type = ast.NodeMathBlock
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeCodeBlock:
		node.Type = ast.NodeCodeBlock
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeHTMLBlock:
		node.Type = ast.NodeHTMLBlock
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeYamlFrontMatter:
		node.Type = ast.NodeYamlFrontMatter
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeBlockEmbed:
		text := lute.domText(n)
		if "" == text {
			return
		}

		t := parse.Parse("", []byte(text), lute.ParseOptions)
		t.Root.LastChild.Unlink() // 移除 doc IAL
		if blockEmbed := t.Root.FirstChild; nil != blockEmbed && ast.NodeBlockEmbed == blockEmbed.Type {
			ial, id := node.KramdownIAL, node.ID
			node = blockEmbed
			node.KramdownIAL, node.ID = ial, id
			next := blockEmbed.Next
			tree.Context.Tip.AppendChild(node)
			appendNextToTip(next, tree)
			return
		}
		node.Type = ast.NodeText
		node.Tokens = []byte(text)
		tree.Context.Tip.AppendChild(node)
		defer tree.Context.ParentTip()
		return
	default:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = lute.domHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTByVditorBlockDOM(c, tree)
	}

	switch dataType {
	case ast.NodeSuperBlock:
		node.AppendChild(&ast.Node{Type: ast.NodeSuperBlockCloseMarker})
	}
}

func (lute *Lute) genASTContenteditable(n *html.Node, tree *parse.Tree) {
	content := n.Data
	node := &ast.Node{Type: ast.NodeText, Tokens: []byte(content)}
	switch n.DataAtom {
	case 0:
		if "" == content {
			return
		}

		checkIndentCodeBlock := strings.ReplaceAll(content, util.Caret, "")
		checkIndentCodeBlock = strings.ReplaceAll(checkIndentCodeBlock, "\t", "    ")
		if (!lute.isInline(n.PrevSibling)) && strings.HasPrefix(checkIndentCodeBlock, "    ") {
			node.Type = ast.NodeCodeBlock
			node.IsFencedCodeBlock = true
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: []byte("```"), CodeBlockFenceLen: 3})
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker})
			startCaret := strings.HasPrefix(content, util.Caret)
			if startCaret {
				content = strings.ReplaceAll(content, util.Caret, "")
			}
			content = strings.TrimSpace(content)
			if startCaret {
				content = util.Caret + content
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
	case atom.Em, atom.I:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if util.Caret == text {
			node.Tokens = util.CaretTokens
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

		if nil != n.FirstChild && util.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.CaretTokens})
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
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if util.Caret == text {
			node.Tokens = util.CaretTokens
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

		if nil != n.FirstChild && util.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.CaretTokens})
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
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if util.Caret == text {
			node.Tokens = util.CaretTokens
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

		if nil != n.FirstChild && util.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.CaretTokens})
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
	case atom.Mark:
		if nil == n.FirstChild || atom.Br == n.FirstChild.DataAtom {
			return
		}
		if lute.startsWithNewline(n.FirstChild) {
			n.FirstChild.Data = strings.TrimLeft(n.FirstChild.Data, parse.Zwsp+"\n")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(parse.Zwsp + "\n")})
		}
		text := strings.TrimSpace(lute.domText(n))
		if lute.isEmptyText(n) {
			return
		}
		if util.Caret == text {
			node.Tokens = util.CaretTokens
			tree.Context.Tip.AppendChild(node)
			return
		}

		node.Type = ast.NodeMark
		marker := lute.domAttrValue(n, "data-marker")
		if "=" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeMark1OpenMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeMark2OpenMarker, Tokens: []byte(marker)})
		}
		tree.Context.Tip.AppendChild(node)

		if nil != n.FirstChild && util.Caret == n.FirstChild.Data && nil != n.LastChild && "br" == n.LastChild.Data {
			// 处理结尾换行
			node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.CaretTokens})
			if "=" == marker {
				node.AppendChild(&ast.Node{Type: ast.NodeMark1CloseMarker, Tokens: []byte(marker)})
			} else {
				node.AppendChild(&ast.Node{Type: ast.NodeMark2CloseMarker, Tokens: []byte(marker)})
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
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTContenteditable(c, tree)
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
		if "" != lute.RenderOptions.LinkBase {
			href = strings.ReplaceAll(href, lute.RenderOptions.LinkBase, "")
		}
		if "" != lute.RenderOptions.LinkPrefix {
			href = strings.ReplaceAll(href, lute.RenderOptions.LinkPrefix, "")
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
	case atom.Mark:
		marker := lute.domAttrValue(n, "data-marker")
		if "=" == marker {
			node.AppendChild(&ast.Node{Type: ast.NodeMark1CloseMarker, Tokens: []byte(marker)})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeMark2CloseMarker, Tokens: []byte(marker)})
		}
	}
}

func (lute *Lute) setBlockIAL(n *html.Node, node *ast.Node) (ialTokens []byte) {
	ialTokens = []byte("{: id=\"" + node.ID + "\"")
	if bookmark := lute.domAttrValue(n, "bookmark"); "" != bookmark {
		bookmark = html.UnescapeString(bookmark)
		bookmark = html.EscapeString(bookmark)
		node.SetIALAttr("bookmark", bookmark)
		ialTokens = append(ialTokens, []byte(" bookmark=\""+bookmark+"\"")...)
	}
	if style := lute.domAttrValue(n, "style"); "" != style {
		node.SetIALAttr("style", style)
		ialTokens = append(ialTokens, []byte(" style=\""+style+"\"")...)
	}
	if name := lute.domAttrValue(n, "name"); "" != name {
		name = html.UnescapeString(name)
		name = html.EscapeString(name)
		node.SetIALAttr("name", name)
		ialTokens = append(ialTokens, []byte(" name=\""+name+"\"")...)
	}
	if memo := lute.domAttrValue(n, "memo"); "" != memo {
		memo = html.UnescapeString(memo)
		memo = html.EscapeString(memo)
		node.SetIALAttr("memo", memo)
		ialTokens = append(ialTokens, []byte(" memo=\""+memo+"\"")...)
	}
	if alias := lute.domAttrValue(n, "alias"); "" != alias {
		alias = html.UnescapeString(alias)
		alias = html.EscapeString(alias)
		node.SetIALAttr("alias", alias)
		ialTokens = append(ialTokens, []byte(" alias=\""+alias+"\"")...)
	}
	if val := lute.domAttrValue(n, "fold"); "" != val {
		val = html.UnescapeString(val)
		val = html.EscapeString(val)
		node.SetIALAttr("fold", val)
		ialTokens = append(ialTokens, []byte(" fold=\""+val+"\"")...)
	}
	if val := lute.domAttrValue(n, "parent-fold"); "" != val {
		val = html.UnescapeString(val)
		val = html.EscapeString(val)
		node.SetIALAttr("parent-fold", val)
		ialTokens = append(ialTokens, []byte(" parent-fold=\""+val+"\"")...)
	}
	if val := lute.domAttrValue(n, "updated"); "" != val {
		val = html.UnescapeString(val)
		val = html.EscapeString(val)
		node.SetIALAttr("updated", val)
		ialTokens = append(ialTokens, []byte(" updated=\""+val+"\"")...)
	}
	ialTokens = append(ialTokens, '}')
	return ialTokens
}
