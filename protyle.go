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

func (lute *Lute) SpinBlockDOM(ivHTML string) (ovHTML string) {
	markdown := lute.blockDOM2Md(ivHTML)
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	ovHTML = strings.ReplaceAll(ovHTML, parse.Zwsp, "")
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
	renderer := render.NewBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.HTML2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	return
}

func (lute *Lute) BlockDOM2HTML(vhtml string) (sHTML string) {
	markdown := lute.blockDOM2Md(vhtml)
	sHTML = lute.Md2HTML(markdown)
	return
}

func (lute *Lute) Md2BlockDOM(markdown string) (vHTML string) {
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	return
}

func (lute *Lute) InlineMd2BlockDOM(markdown string) (vHTML string) {
	tree := parse.Inline("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2BlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	return
}

func (lute *Lute) BlockDOM2Md(htmlStr string) (markdown string) {
	//fmt.Println(htmlStr)
	markdown = lute.blockDOM2Md(htmlStr)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	return
}

func (lute *Lute) BlockDOM2StdMd(htmlStr string) (markdown string) {
	htmlStr = strings.ReplaceAll(htmlStr, parse.Zwsp, "")

	// DOM 转 AST
	tree, err := lute.BlockDOM2Tree(htmlStr)
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
	markdown = util.BytesToStr(formatted)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	return
}

func (lute *Lute) BlockDOM2Text(htmlStr string) (text string) {
	tree, err := lute.BlockDOM2Tree(htmlStr)
	if nil != err {
		return ""
	}
	return tree.Root.Text()
}

func (lute *Lute) BlockDOM2TextLen(htmlStr string) int {
	tree, err := lute.BlockDOM2Tree(htmlStr)
	if nil != err {
		return 0
	}
	return tree.Root.TextLen()
}

func (lute *Lute) Tree2BlockDOM(tree *parse.Tree, options *render.Options) (vHTML string) {
	renderer := render.NewBlockRenderer(tree, options)
	output := renderer.Render()
	vHTML = util.BytesToStr(output)
	vHTML = strings.ReplaceAll(vHTML, util.Caret, "<wbr>")
	return
}

func RenderNodeBlockDOM(node *ast.Node, parseOptions *parse.Options, renderOptions *render.Options) string {
	root := &ast.Node{Type: ast.NodeDocument}
	tree := &parse.Tree{Root: root, Context: &parse.Context{ParseOption: parseOptions}}
	renderer := render.NewBlockRenderer(tree, renderOptions)
	renderer.Writer = &bytes.Buffer{}
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := renderer.RendererFuncs[n.Type]
		return rendererFunc(n, entering)
	})
	return renderer.Writer.String()
}

func (lute *Lute) BlockDOM2Tree(htmlStr string) (ret *parse.Tree, err error) {
	htmlStr = strings.ReplaceAll(htmlStr, "<wbr>", util.Caret)

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
		lute.genASTByBlockDOM(c, ret)
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

func (lute *Lute) H2P(ivHTML string) (ovHTML string) {
	tree, err := lute.BlockDOM2Tree(ivHTML)
	if nil != err {
		return err.Error()
	}

	node := tree.Root.FirstChild
	if ast.NodeHeading != node.Type {
		return ivHTML
	}

	node.Type = ast.NodeParagraph
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) P2H(ivHTML, level string) (ovHTML string) {
	tree, err := lute.BlockDOM2Tree(ivHTML)
	if nil != err {
		return err.Error()
	}

	node := tree.Root.FirstChild
	if ast.NodeParagraph != node.Type {
		return ivHTML
	}

	node.Type = ast.NodeHeading
	node.HeadingLevel, _ = strconv.Atoi(level)
	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) OL2UL(ivHTML string) (ovHTML string) {
	tree, err := lute.BlockDOM2Tree(ivHTML)
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

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) UL2OL(ivHTML string) (ovHTML string) {
	tree, err := lute.BlockDOM2Tree(ivHTML)
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

	ovHTML = lute.Tree2BlockDOM(tree, lute.RenderOptions)
	return
}

func (lute *Lute) blockDOM2Md(htmlStr string) (markdown string) {
	tree, err := lute.BlockDOM2Tree(htmlStr)
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

func (lute *Lute) genASTByBlockDOM(n *html.Node, tree *parse.Tree) {
	class := lute.domAttrValue(n, "class")
	if "protyle-attr" == class ||
		strings.Contains(class, "__copy") ||
		strings.Contains(class, "protyle-linenumber__rows") {
		return
	}

	if "1" == lute.domAttrValue(n, "spin") {
		return
	}

	if "protyle-action" == class {
		if ast.NodeCodeBlock == tree.Context.Tip.Type {
			languageNode := n.FirstChild
			language := ""
			if nil != languageNode.FirstChild {
				language = languageNode.FirstChild.Data
			}
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker, Tokens: util.StrToBytes(language)})
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeCodeBlockCode, Tokens: util.StrToBytes(lute.domText(n.NextSibling))})
		}
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
	case ast.NodeTable:
		node.Type = ast.NodeTable
		var tableAligns []int
		if nil == n.FirstChild {
			return
		}

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
		node.Tokens = nil
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()

		lute.genASTContenteditable(n.FirstChild.NextSibling.FirstChild, tree)
		return
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
		marker := lute.domAttrValue(n, "data-marker")
		if ast.NodeList != tree.Context.Tip.Type {
			parent := &ast.Node{}
			parent.Type = ast.NodeList
			parent.ListData = &ast.ListData{}
			subType := lute.domAttrValue(n, "data-subtype")
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
			tree.Context.Tip.AppendChild(parent)
			tree.Context.Tip = parent
		}

		node.Type = ast.NodeListItem
		node.ListData = &ast.ListData{}
		subType := lute.domAttrValue(n, "data-subtype")
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
		layout := lute.domAttrValue(n, "data-sb-layout")
		node.AppendChild(&ast.Node{Type: ast.NodeSuperBlockLayoutMarker, Tokens: []byte(layout)})
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeMathBlock:
		node.Type = ast.NodeMathBlock
		node.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker})
		content := lute.domAttrValue(n, "data-content")
		node.AppendChild(&ast.Node{Type: ast.NodeMathBlockContent, Tokens: util.StrToBytes(content)})
		node.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker})
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeCodeBlock:
		node.Type = ast.NodeCodeBlock
		node.IsFencedCodeBlock = true
		node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: util.StrToBytes("```")})
		if language := lute.domAttrValue(n, "data-subtype"); "" != language {
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker, CodeBlockInfo: util.StrToBytes(language)})
			content := lute.domAttrValue(n, "data-content")
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
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeYamlFrontMatter:
		node.Type = ast.NodeYamlFrontMatter
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case ast.NodeThematicBreak:
		node.Type = ast.NodeThematicBreak
		tree.Context.Tip.AppendChild(node)
		return
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
	case ast.NodeIFrame:
		node.Type = ast.NodeIFrame
		node.Tokens = lute.domHTML(n.FirstChild.NextSibling.FirstChild)
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeVideo:
		node.Type = ast.NodeVideo
		node.Tokens = lute.domHTML(n.FirstChild.NextSibling.FirstChild)
		tree.Context.Tip.AppendChild(node)
		return
	case ast.NodeAudio:
		node.Type = ast.NodeAudio
		node.Tokens = lute.domHTML(n.FirstChild.NextSibling.FirstChild)
		tree.Context.Tip.AppendChild(node)
		return
	default:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = lute.domHTML(n)
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
	if ast.NodeCodeBlock == tree.Context.Tip.Type {
		return
	}

	class := lute.domAttrValue(n, "class")
	if "svg" == class {
		return
	}

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
		if ast.NodeLink == tree.Context.Tip.Type {
			node.Type = ast.NodeLinkText
		}
		tree.Context.Tip.AppendChild(node)
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
	case atom.Code:
		if lute.isEmptyText(n) {
			return
		}

		node.Type = ast.NodeCodeSpan
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanOpenMarker})
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanContent})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Span:
		dataType := lute.domAttrValue(n, "data-type")
		if "tag" == dataType {
			node.Type = ast.NodeTag
			node.AppendChild(&ast.Node{Type: ast.NodeTagOpenMarker})
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		} else if "inline-math" == dataType {
			node.Type = ast.NodeInlineMath
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathOpenMarker})
			content = lute.domAttrValue(n, "data-content")
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathContent, Tokens: util.StrToBytes(content)})
			node.AppendChild(&ast.Node{Type: ast.NodeInlineMathCloseMarker})
			tree.Context.Tip.AppendChild(node)
			return
		} else if "a" == dataType {
			node.Type = ast.NodeLink
			node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		} else if "block-ref" == dataType {
			node.Type = ast.NodeBlockRef
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			id := lute.domAttrValue(n, "data-id")
			node.AppendChild(&ast.Node{Type: ast.NodeBlockRefID, Tokens: util.StrToBytes(id)})
			node.AppendChild(&ast.Node{Type: ast.NodeBlockRefSpace})
			refText := lute.domAttrValue(n, "data-anchor")
			refTextNode := &ast.Node{Type: ast.NodeBlockRefText, Tokens: util.StrToBytes(refText)}
			node.AppendChild(refTextNode)
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			tree.Context.Tip.AppendChild(node)
			return
		} else if "img" == dataType {
			img := n.FirstChild.NextSibling
			node.Type = ast.NodeImage
			node.AppendChild(&ast.Node{Type: ast.NodeBang})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
			alt := lute.domAttrValue(img, "alt")
			node.AppendChild(&ast.Node{Type: ast.NodeLinkText, Tokens: util.StrToBytes(alt)})
			node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			src := lute.domAttrValue(img, "data-src")
			node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: util.StrToBytes(src)})
			if title := lute.domAttrValue(img, "title"); "" != title {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
				node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: util.StrToBytes(title)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			tree.Context.Tip.AppendChild(node)
			return
		}
	case atom.Sub:
		if lute.isEmptyText(n) {
			return
		}

		node.Type = ast.NodeSub
		node.AppendChild(&ast.Node{Type: ast.NodeSubOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Sup:
		if lute.isEmptyText(n) {
			return
		}

		node.Type = ast.NodeSup
		node.AppendChild(&ast.Node{Type: ast.NodeSupOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.U:
		if lute.isEmptyText(n) {
			return
		}
		node.Type = ast.NodeInlineHTML
		node.Tokens = []byte("<u>")
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
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

		lute.setSpanIAL(n, node)

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
	case atom.Code:
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanCloseMarker})
	case atom.Span:
		dataType := lute.domAttrValue(n, "data-type")
		if "tag" == dataType {
			node.AppendChild(&ast.Node{Type: ast.NodeTagCloseMarker})
		} else if "a" == dataType {
			node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			href := lute.domAttrValue(n, "data-href")
			if "" != lute.RenderOptions.LinkBase {
				href = strings.ReplaceAll(href, lute.RenderOptions.LinkBase, "")
			}
			if "" != lute.RenderOptions.LinkPrefix {
				href = strings.ReplaceAll(href, lute.RenderOptions.LinkPrefix, "")
			}
			node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: []byte(href)})
			linkTitle := lute.domAttrValue(n, "data-title")
			if "" != linkTitle {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
				node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: []byte(linkTitle)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
		} else if "block-ref" == dataType {
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
		}
	case atom.Sub:
		node.AppendChild(&ast.Node{Type: ast.NodeSubCloseMarker})
	case atom.Sup:
		node.AppendChild(&ast.Node{Type: ast.NodeSupCloseMarker})
	case atom.U:
		tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeInlineHTML, Tokens: []byte("</u>")})
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

func (lute *Lute) setSpanIAL(n *html.Node, node *ast.Node) {
	style := lute.domAttrValue(n, "style")
	if "" != style {
		node.SetIALAttr("style", style)
		node.KramdownIAL = [][]string{{"style", style}}
		ialTokens := parse.IAL2Tokens(node.KramdownIAL)
		ial := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
		node.InsertAfter(ial)
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

func appendNextToTip(next *ast.Node, tree *parse.Tree) {
	var nodes []*ast.Node
	for n := next; nil != n; n = n.Next {
		nodes = append(nodes, n)
	}
	for _, n := range nodes {
		tree.Context.Tip.AppendChild(n)
	}
}
