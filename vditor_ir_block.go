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

// SpinVditorIRBlockDOM 自旋 Vditor Instant-Rendering Block DOM，用于即时渲染块模式下的编辑。
func (lute *Lute) SpinVditorIRBlockDOM(ivHTML string) (ovHTML string) {
	// 替换插入符
	ivHTML = strings.ReplaceAll(ivHTML, "<wbr>", util.Caret)

	markdown := lute.vditorIRBlockDOM2Md(ivHTML)
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)

	ovHTML = lute.Tree2VditorIRBlockDOM(tree, lute.RenderOptions)
	// 替换插入符
	ovHTML = strings.ReplaceAll(ovHTML, util.Caret, "<wbr>")
	return
}

// HTML2VditorIRBlockDOM 将 HTML 转换为 Vditor Instant-Rendering Block DOM，用于即时渲染块模式下粘贴。
func (lute *Lute) HTML2VditorIRBlockDOM(sHTML string) (vHTML string) {
	markdown, err := lute.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorIRBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.HTML2VditorIRBlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

// VditorIRBlockDOM2HTML 将 Vditor Instant-Rendering Block DOM 转换为 HTML，用于 Vditor.getHTML() 接口。
func (lute *Lute) VditorIRBlockDOM2HTML(vhtml string) (sHTML string) {
	markdown := lute.vditorIRBlockDOM2Md(vhtml)
	sHTML = lute.Md2HTML(markdown)
	return
}

// Md2VditorIRBlockDOM 将 markdown 转换为 Vditor Instant-Rendering Block DOM，用于从源码模式切换至即时渲染块模式。
func (lute *Lute) Md2VditorIRBlockDOM(markdown string) (vHTML string) {
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorIRBlockRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2VditorIRBlockDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

// VditorIRBlockDOM2Md 将 Vditor Instant-Rendering DOM 转换为 markdown，用于从即时渲染块模式切换至源码模式。
func (lute *Lute) VditorIRBlockDOM2Md(htmlStr string) (markdown string) {
	htmlStr = strings.ReplaceAll(htmlStr, parse.Zwsp, "")
	markdown = lute.vditorIRBlockDOM2Md(htmlStr)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	return
}

// VditorIRBlockDOM2StdMd 将 Vditor Instant-Rendering DOM 转换为标准 markdown，用于复制剪切。
func (lute *Lute) VditorIRBlockDOM2StdMd(htmlStr string) (markdown string) {
	htmlStr = strings.ReplaceAll(htmlStr, parse.Zwsp, "")

	// DOM 转 AST
	tree, err := lute.VditorIRBlockDOM2Tree(htmlStr)
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
	options.KramdownIAL = true
	renderer := render.NewFormatRenderer(tree, options)
	formatted := renderer.Render()
	markdown = string(formatted)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	return
}

func (lute *Lute) VditorIRBlockDOM2Text(htmlStr string) (text string) {
	tree, err := lute.VditorIRBlockDOM2Tree(htmlStr)
	if nil != err {
		return ""
	}
	return tree.Root.Text()
}

func (lute *Lute) VditorIRBlockDOM2TextLen(htmlStr string) int {
	tree, err := lute.VditorIRBlockDOM2Tree(htmlStr)
	if nil != err {
		return 0
	}
	return tree.Root.TextLen()
}

func (lute *Lute) Tree2VditorIRBlockDOM(tree *parse.Tree, options *render.Options) (vHTML string) {
	renderer := render.NewVditorIRBlockRenderer(tree, options)
	output := renderer.Render()
	vHTML = string(output)
	return
}

func RenderNodeVditorIRBlockDOM(node *ast.Node, parseOptions *parse.Options, renderOptions *render.Options) string {
	root := &ast.Node{Type: ast.NodeDocument}
	tree := &parse.Tree{Root: root, Context: &parse.Context{ParseOption: parseOptions}}
	renderer := render.NewVditorIRBlockRenderer(tree, renderOptions)
	renderer.Writer = &bytes.Buffer{}
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := renderer.RendererFuncs[n.Type]
		return rendererFunc(n, entering)
	})
	return renderer.Writer.String()
}

func (lute *Lute) VditorIRBlockDOM2Tree(htmlStr string) (ret *parse.Tree, err error) {
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
		return
	}

	// 调整 DOM 结构
	lute.adjustVditorDOM(htmlNodes)

	// 将 HTML 树转换为 Markdown AST

	ret = &parse.Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument}, Context: &parse.Context{ParseOption: lute.ParseOptions}}
	ret.Context.Tip = ret.Root
	for _, htmlNode := range htmlNodes {
		lute.genASTByVditorIRBlockDOM(htmlNode, ret)
	}

	// 调整树结构
	ast.Walk(ret.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
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
						if ast.NodeKramdownBlockIAL == previousLi.Type {
							previousLi = previousLi.Previous
						}
						listIAL := n.Next
						previousLi.AppendChild(n)
						previousLi.AppendChild(listIAL)
					}
				}
			}
		}
		return ast.WalkContinue
	})
	return
}

func (lute *Lute) vditorIRBlockDOM2Md(htmlStr string) (markdown string) {
	tree, err := lute.VditorIRBlockDOM2Tree(htmlStr)
	if nil != err {
		return err.Error()
	}

	// 将 AST 进行 Markdown 格式化渲染
	options := render.NewOptions()
	options.AutoSpace = false
	options.FixTermTypo = false
	options.KramdownIAL = true
	renderer := render.NewFormatRenderer(tree, options)
	formatted := renderer.Render()
	markdown = string(formatted)
	return
}

// genASTByVditorIRBlockDOM 根据指定的 Vditor IR DOM 节点 n 进行深度优先遍历并逐步生成 Markdown 语法树 tree。
func (lute *Lute) genASTByVditorIRBlockDOM(n *html.Node, tree *parse.Tree) {
	dataRender := lute.domAttrValue(n, "data-render")
	if "1" == dataRender || "2" == dataRender { // 1：浮动工具栏，2：preview 代码块、数学公式块或者不解析的节点
		return
	}

	dataType := lute.domAttrValue(n, "data-type")
	if "ref-text-tpl-render-result" == dataType { // 剔除渲染好的锚文本
		return
	}

	class := lute.domAttrValue(n, "class")
	content := strings.ReplaceAll(n.Data, parse.Zwsp, "")
	nodeID := lute.domAttrValue(n, "data-node-id")
	node := &ast.Node{ID: nodeID, Type: ast.NodeText, Tokens: []byte(content)}
	if "" == nodeID {
		if "p" == dataType || "ul" == dataType || "ol" == dataType || "blockquote" == dataType ||
			"math-block" == dataType || "code-block" == dataType || "table" == dataType || "h" == dataType ||
			"link-ref-defs-block" == dataType || "footnotes-block" == dataType || "super-block" == dataType {
			nodeID = ast.NewNodeID()
			node.ID = nodeID
		}
	}
	if "" != node.ID {
		node.KramdownIAL = [][]string{{"id", node.ID}}
		ialTokens := lute.setIAL(n, node)
		ial := &ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: ialTokens}
		defer tree.Context.TipAppendChild(ial)
	}

	if atom.Div == n.DataAtom {
		if "link-ref-defs-block" == dataType {
			text := lute.domText(n)
			if !strings.HasPrefix(text, "[") {
				subTree := parse.Parse("", []byte(text), lute.ParseOptions)
				if nil != subTree.Root.FirstChild {
					tree.Context.Tip.AppendChild(subTree.Root.FirstChild)
				}
				return
			}

			defBlock := &ast.Node{Type: ast.NodeLinkRefDefBlock}
			tree.Context.Tip.AppendChild(defBlock)
			for def := n.FirstChild; nil != def; def = def.NextSibling {
				text = lute.domText(def)
				subTree := parse.Parse("", []byte(text), lute.ParseOptions)
				child := subTree.Root.FirstChild.FirstChild
				if ast.NodeLinkRefDef == child.Type {
					defBlock.AppendChild(subTree.Root.FirstChild.FirstChild)
				} else {
					tree.Context.Tip.AppendChild(subTree.Root.FirstChild)
				}
			}
			return
		} else if "footnotes-def" == dataType {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if nil == c.FirstChild {
					continue
				}
				if c == n.FirstChild && !strings.HasPrefix(c.FirstChild.Data, "[^") {
					tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(lute.domText(c))})
					continue
				}

				if strings.HasPrefix(c.FirstChild.Data, "[^") && strings.Contains(c.FirstChild.Data, "]: ") {
					label := c.FirstChild.Data[1:strings.Index(c.FirstChild.Data, "]: ")]
					tree.Context.Tip.Tokens = []byte(label)
					c.FirstChild.Data = c.FirstChild.Data[strings.Index(c.FirstChild.Data, "]: ")+3:]
				}
				lute.genASTByVditorIRBlockDOM(c, tree)
			}
			return
		} else if "footnotes-block" == dataType {
			footnotesBlock := &ast.Node{Type: ast.NodeFootnotesDefBlock}
			tree.Context.Tip.AppendChild(footnotesBlock)
			for def := n.FirstChild; nil != def; def = def.NextSibling {
				defNode := &ast.Node{Type: ast.NodeFootnotesDef}
				originalHTML := &bytes.Buffer{}
				if err := html.Render(originalHTML, def); nil == err {
					subTree, _ := lute.VditorIRBlockDOM2Tree(originalHTML.String())
					if nil != subTree.Root.Tokens {
						var children []*ast.Node
						for c := subTree.Root.FirstChild; nil != c; c = c.Next {
							children = append(children, c)
						}
						for _, c := range children {
							defNode.AppendChild(c)
						}
						defNode.Tokens = subTree.Root.Tokens
						footnotesBlock.AppendChild(defNode)
					} else {
						footnotesBlock.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(subTree.Root.Text())})
					}
				}
			}
			return
		} else if "toc-block" == dataType {
			node := &ast.Node{Type: ast.NodeToC}
			tree.Context.Tip.AppendChild(node)
			return
		} else if "block-query-embed" == dataType {
			text := lute.domText(n)
			t := parse.Parse("", []byte(text), lute.ParseOptions)
			t.Root.LastChild.Unlink() // 移除 doc IAL
			if blockQueryEmbed := t.Root.FirstChild; nil != blockQueryEmbed && ast.NodeBlockQueryEmbed == blockQueryEmbed.Type {
				ial, id := node.KramdownIAL, node.ID
				node = blockQueryEmbed
				node.KramdownIAL, node.ID = ial, id
				next := blockQueryEmbed.Next
				tree.Context.Tip.AppendChild(node)
				appendNextToTip(next, tree)
				return
			}

			node := &ast.Node{Type: ast.NodeText, Tokens: []byte(text)}
			tree.Context.Tip.AppendChild(node)
			return
		}
	}

	switch n.DataAtom {
	case 0:
		if "" == content {
			return
		}

		if html.ElementNode == n.Type {
			break
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

		if ast.NodeCodeBlock == tree.Context.Tip.Type {
			// 开始代码块标记符后退格的情况
			tree.Context.Tip.Type = ast.NodeParagraph
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(content)})
			if nil != n.NextSibling && nil != n.NextSibling.NextSibling && nil != n.NextSibling.NextSibling.NextSibling &&
				nil != n.NextSibling.NextSibling.NextSibling.NextSibling {
				n.Parent.RemoveChild(n.NextSibling.NextSibling.NextSibling.NextSibling)
			}
			break
		}

		tokens := make([]byte, len(node.Tokens))
		copy(tokens, node.Tokens)

		// 尝试块级解析，处理列表代码块
		subTree := parse.Parse("", tokens, tree.Context.ParseOption)
		if nil != subTree.Root.FirstChild && ast.NodeCodeBlock == subTree.Root.FirstChild.Type {
			node.Tokens = bytes.TrimPrefix(node.Tokens, []byte("\n"))
			tree.Context.Tip.AppendChild(node)
		} else {
			// 尝试行级解析，处理段落图片文本节点转换为图片节点
			subTree = parse.Inline("", tokens, tree.Context.ParseOption)
			if ast.NodeSoftBreak == subTree.Root.FirstChild.FirstChild.Type || // 软换行
				(ast.NodeParagraph == subTree.Root.FirstChild.Type &&
					(ast.NodeImage == subTree.Root.FirstChild.FirstChild.Type ||
						(ast.NodeSoftBreak == subTree.Root.FirstChild.FirstChild.Type && nil != subTree.Root.FirstChild.FirstChild.Next &&
							(ast.NodeText == subTree.Root.FirstChild.FirstChild.Next.Type ||
								ast.NodeEmphasis == subTree.Root.FirstChild.FirstChild.Next.Type ||
								ast.NodeStrong == subTree.Root.FirstChild.FirstChild.Next.Type ||
								ast.NodeStrikethrough == subTree.Root.FirstChild.FirstChild.Next.Type ||
								ast.NodeCodeSpan == subTree.Root.FirstChild.FirstChild.Next.Type ||
								ast.NodeMark == subTree.Root.FirstChild.FirstChild.Next.Type))) && // 软换行后跟普通文本
					nil == subTree.Root.FirstChild.Next) {
				appendNextToTip(subTree.Root.FirstChild.FirstChild, tree)
			} else {
				node.Tokens = bytes.TrimPrefix(node.Tokens, []byte("\n"))
				tree.Context.Tip.AppendChild(node)
			}
		}
	case atom.P:
		node.Type = ast.NodeParagraph
		text := lute.domText(n)
		if "\n" == text && ast.NodeBlockquote == tree.Context.Tip.Type && nil == tree.Context.Tip.FirstChild.Next {
			// 不允许在 bq 第一个节点前换行
			return
		} else {
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		}
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		text := lute.domText(n)
		if "" == strings.TrimSpace(text) {
			return
		}
		node.Type = ast.NodeHeading
		marker := lute.domAttrValue(n, "data-marker")
		node.HeadingSetext = "=" == marker || "-" == marker
		if !node.HeadingSetext {
			marker := lute.domText(n.FirstChild)
			node.HeadingLevel = bytes.Count([]byte(marker), []byte("#"))
		} else {
			if n.FirstChild == n.LastChild || "" == strings.TrimSpace(strings.ReplaceAll(lute.domText(n.LastChild), util.Caret, "")) {
				node.Type = ast.NodeText
				node.Tokens = []byte(text)
				tree.Context.Tip.AppendChild(node)
				tree.Context.Tip = node
				return
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
	case atom.Ol, atom.Ul:
		if nil == n.FirstChild {
			return
		}

		node.Type = ast.NodeList
		node.ListData = &ast.ListData{}
		marker := lute.domAttrValue(n, "data-marker")
		if "" == marker {
			marker = "*"
		}
		if atom.Ol == n.DataAtom {
			node.ListData.Typ = 1
			start := lute.domAttrValue(n, "start")
			if "" == start {
				start = "1"
			}
			node.ListData.Start, _ = strconv.Atoi(start)
		} else {
			node.ListData.BulletChar = marker[0]
		}
		node.ListData.Marker = []byte(marker)
		tight := lute.domAttrValue(n, "data-tight")
		if "true" == tight || "" == tight {
			node.Tight = true
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Li:
		if ast.NodeList != tree.Context.Tip.Type {
			parent := &ast.Node{}
			parent.Type = ast.NodeList
			parent.ListData = &ast.ListData{Tight: true}
			marker := lute.domAttrValue(n, "data-marker")
			if "" == marker {
				marker = "*"
			}
			tree.Context.Tip.AppendChild(parent)
			tree.Context.Tip = parent
		}

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
			node.ListData.Num, _ = strconv.Atoi(marker[:len(marker)-1])
			node.ListData.Delimiter = marker[len(marker)-1]
		}
		if "vditor-task" == lute.domAttrValue(n, "class") {
			node.ListData.Typ = 3
			tree.Context.Tip.ListData.Typ = 3
		}

		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Pre:
		if atom.Code == n.FirstChild.DataAtom {
			var codeTokens []byte
			if nil != n.FirstChild.FirstChild {
				codeTokens = []byte(n.FirstChild.FirstChild.Data)
				for next := n.FirstChild.FirstChild.NextSibling; nil != next; next = next.NextSibling {
					// YAML Front Matter 中删除问题 https://github.com/siyuan-note/siyuan/issues/109
					codeTokens = append(codeTokens, []byte(lute.domText(next))...)
				}
			}

			divDataType := lute.domAttrValue(n.Parent, "data-type")
			switch divDataType {
			case "math-block":
				node.Type = ast.NodeMathBlockContent
				node.Tokens = codeTokens
				tree.Context.Tip.AppendChild(node)
			case "html-block":
				tree.Context.Tip.Tokens = codeTokens
			case "yaml-front-matter":
				node.Type = ast.NodeYamlFrontMatterContent
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
		if util.Caret == text {
			node.Tokens = util.CaretTokens
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
		if lute.starstWithNewline(n.FirstChild) {
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
		tree.Context.Tip.AppendChild(node)
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
		if util.Caret == text {
			node.Tokens = util.CaretTokens
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
		if lute.starstWithNewline(n.FirstChild) {
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
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Code:
		if nil == n.FirstChild {
			return
		}
		firstChildDataType := lute.domAttrValue(n.FirstChild, "data-type")
		if firstChildDataType == "html-inline" {
			break
		}

		contentStr := strings.ReplaceAll(n.FirstChild.Data, parse.Zwsp, "")
		if util.Caret == contentStr {
			node.Tokens = util.CaretTokens
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
				if (nil == n.PrevSibling || util.Caret == n.PrevSibling.Data) && (nil == n.NextSibling || util.Caret == n.NextSibling.Data) {
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
		if nil == n.Parent || (atom.P != n.Parent.DataAtom && atom.Li != n.Parent.DataAtom) {
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
		if nil != node.Parent.Parent.Parent && nil != node.Parent.Parent.Parent.ListData { // ul.li.p.input
			node.Parent.Parent.Parent.ListData.Typ = 3
			node.Parent.Parent.ListData.Typ = 3
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
		node.Tokens = nil
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

		text := lute.domText(n)
		if "footnotes-ref" == dataType {
			if strings.HasPrefix(text, "[^") && strings.HasSuffix(text, "]") {
				node.Type = ast.NodeFootnotesRef
				node.Tokens = []byte(text[1 : len(text)-1])
				tree.Context.Tip.AppendChild(node)
			} else {
				node.Type = ast.NodeText
				node.Tokens = []byte(text)
				tree.Context.Tip.AppendChild(node)
			}
		} else {
			node.Type = ast.NodeInlineHTML
			text := lute.domText(n)
			node.Tokens = []byte(text)
			tree.Context.Tip.AppendChild(node)
		}
		return
	case atom.Span:
		switch dataType {
		case "block-ref":
			text := lute.domText(n)
			if "" == strings.TrimSpace(text) {
				return
			}

			t := parse.Parse("", []byte(text), lute.ParseOptions)
			if textNode := t.Root.FirstChild.FirstChild; nil != textNode && ast.NodeText == textNode.Type &&
				nil != textNode.Next && ast.NodeBlockRef == textNode.Next.Type {
				content := textNode.Text()
				if ("！"+util.Caret == content) || ("!"+util.Caret == content) {
					text = strings.Replace(text, content, "!", 1) + util.Caret
					t = parse.Parse("", []byte(text), lute.ParseOptions)
				}
			}
			t.Root.LastChild.Unlink() // 移除 doc IAL
			if blockRef := t.Root.FirstChild.FirstChild; nil != blockRef && ast.NodeBlockRef == blockRef.Type {
				node = blockRef
				next := blockRef.Next
				tree.Context.Tip.AppendChild(node)
				appendNextToTip(next, tree)
				return
			}
			if blockEmbed := t.Root.FirstChild; nil != blockEmbed && ast.NodeBlockEmbed == blockEmbed.Type {
				node = blockEmbed
				next := blockEmbed.Next
				tree.Context.Tip.AppendChild(node)
				appendNextToTip(next, tree)
				return
			}

			node.Type = ast.NodeText
			node.Tokens = []byte(text)
			tree.Context.Tip.AppendChild(node)
			return
		case "heading-id":
			node.Type = ast.NodeHeadingID
			text := lute.domText(n)
			if "" == strings.TrimSpace(text) {
				return
			}

			if !strings.HasSuffix(text, "}") {
				node.Type = ast.NodeText
				node.Tokens = []byte(text)
				tree.Context.Tip.AppendChild(node)
				return
			}
			text = strings.TrimSpace(text)
			node.Tokens = []byte(text[1 : len(text)-1])
			tree.Context.Tip.AppendChild(node)
			return
		case "em", "strong", "s", "mark", "code", "inline-math", "tag", "sup", "sub", "span-ial":
			text := lute.domText(n)
			if "" == strings.TrimSpace(text) {
				return
			}

			t := parse.Parse("", []byte(text), lute.ParseOptions)
			t.Root.LastChild.Unlink() // 移除 doc IAL
			if inlineNode := t.Root.FirstChild.FirstChild; nil != inlineNode && (ast.NodeEmphasis == inlineNode.Type ||
				ast.NodeStrong == inlineNode.Type || ast.NodeStrikethrough == inlineNode.Type || ast.NodeMark == inlineNode.Type ||
				ast.NodeCodeSpan == inlineNode.Type || ast.NodeInlineMath == inlineNode.Type || ast.NodeTag == inlineNode.Type ||
				ast.NodeSup == inlineNode.Type || ast.NodeSub == inlineNode.Type) {
				node = inlineNode
				next := inlineNode.Next
				tree.Context.Tip.AppendChild(node)
				appendNextToTip(next, tree)
				style := lute.domAttrValue(n.FirstChild.NextSibling, "style")
				if "" != style {
					node.SetIALAttr("style", style)
					node.KramdownIAL = [][]string{{"style", style}}
					ialTokens := []byte("{: style=\"" + style + "\"}")
					ial := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
					node.InsertAfter(ial)
				}
				return
			} else if ast.NodeKramdownBlockIAL == t.Root.FirstChild.Type { // Span IAL 单独存在时会被解析为 Block IAL
				t.Root.FirstChild.Type = ast.NodeKramdownSpanIAL
				node = t.Root.FirstChild
				node.Tokens = bytes.TrimSpace(node.Tokens)
				next := t.Root.FirstChild.Next
				tree.Context.Tip.AppendChild(node)
				appendNextToTip(next, tree)
				return
			}
			node.Type = ast.NodeText
			node.Tokens = []byte(text)
			tree.Context.Tip.AppendChild(node)
			return
		case "a", "link-ref", "img":
			text := lute.domText(n)
			if "" == strings.TrimSpace(text) {
				return
			}

			t := parse.Parse("", []byte(text), lute.ParseOptions)
			t.Root.LastChild.Unlink() // 移除 doc IAL
			if inlineNode := t.Root.FirstChild.FirstChild; nil != inlineNode && (ast.NodeLink == inlineNode.Type || ast.NodeImage == inlineNode.Type) {
				node = inlineNode
				next := inlineNode.Next
				tree.Context.Tip.AppendChild(node)
				appendNextToTip(next, tree)
				if nil != t.Root.FirstChild.Next {
					tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeHardBreak})
					nextBlock := t.Root.FirstChild.Next
					tree.Context.Tip.InsertAfter(nextBlock)
					tree.Context.Tip = nextBlock
				}
				img := lute.domChild(n, atom.Img)
				style := lute.domAttrValue(img, "style")
				if "" != style {
					node.SetIALAttr("style", style)
					node.KramdownIAL = [][]string{{"style", style}}
					ialTokens := []byte("{: style=\"" + style + "\"}")
					ial := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
					node.InsertAfter(ial)
				}
				return
			}
			node.Type = ast.NodeText
			node.Tokens = []byte(text)
			tree.Context.Tip.AppendChild(node)
			return
		case "html-inline":
			text := lute.domText(n)
			if "" == strings.TrimSpace(text) {
				return
			}

			node.Type = ast.NodeInlineHTML
			node.Tokens = []byte(text)
			tree.Context.Tip.AppendChild(node)
			return
		case "html-entity":
			text := lute.domText(n)
			if "" == strings.TrimSpace(text) {
				return
			}

			t := parse.Parse("", []byte(text), lute.ParseOptions)
			t.Root.LastChild.Unlink() // 移除 doc IAL
			if inlineNode := t.Root.FirstChild.FirstChild; nil != inlineNode && (ast.NodeHTMLEntity == inlineNode.Type) {
				node = inlineNode
				next := inlineNode.Next
				tree.Context.Tip.AppendChild(node)
				appendNextToTip(next, tree)
				return
			}
			node.Type = ast.NodeText
			node.Tokens = []byte(text)
			tree.Context.Tip.AppendChild(node)
			return
		case "emoji":
			text := lute.domText(n.FirstChild.NextSibling)
			if "" == strings.TrimSpace(text) {
				return
			}

			t := parse.Parse("", []byte(text), lute.ParseOptions)
			t.Root.LastChild.Unlink() // 移除 doc IAL
			if inlineNode := t.Root.FirstChild.FirstChild; nil != inlineNode && (ast.NodeEmoji == inlineNode.Type) {
				node = inlineNode
				next := inlineNode.Next
				tree.Context.Tip.AppendChild(node)
				appendNextToTip(next, tree)
				return
			}
			node.Type = ast.NodeText
			node.Tokens = []byte(text)
			tree.Context.Tip.AppendChild(node)
			return
		case "backslash":
			text := lute.domText(n)
			if "" == strings.TrimSpace(text) {
				return
			}

			t := parse.Parse("", []byte(text), lute.ParseOptions)
			t.Root.LastChild.Unlink() // 移除 doc IAL
			if inlineNode := t.Root.FirstChild.FirstChild; nil != inlineNode && (ast.NodeBackslash == inlineNode.Type) {
				node = inlineNode
				next := inlineNode.Next
				tree.Context.Tip.AppendChild(node)
				appendNextToTip(next, tree)
				return
			}
			node.Type = ast.NodeText
			node.Tokens = []byte(text)
			tree.Context.Tip.AppendChild(node)
			return
		case "inline-node":
			text := lute.domText(n)
			node.Type = ast.NodeText
			if "</font>" == text {
				node.Type = ast.NodeInlineHTML
			}
			node.Tokens = []byte(text)
			tree.Context.Tip.AppendChild(node)
			return
		case "math-block-close-marker":
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker, Tokens: parse.MathBlockMarker})
			return
		case "math-block-open-marker":
			node.Type = ast.NodeMathBlockOpenMarker
			node.Tokens = parse.MathBlockMarker
			tree.Context.Tip.AppendChild(node)
			return
		case "yaml-front-matter-close-marker":
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeYamlFrontMatterCloseMarker, Tokens: parse.YamlFrontMatterMarker})
			return
		case "yaml-front-matter-open-marker":
			node.Type = ast.NodeYamlFrontMatterOpenMarker
			node.Tokens = parse.MathBlockMarker
			tree.Context.Tip.AppendChild(node)
			return
		case "code-block-open-marker":
			if atom.Pre == n.NextSibling.DataAtom { // DOM 后缺少 info span 节点
				n.InsertAfter(&html.Node{DataAtom: atom.Span, Attr: []*html.Attribute{{Key: "data-type", Val: "code-block-info"}}})
			}
			marker := []byte(lute.domText(n))
			lastBacktick := bytes.LastIndex(marker, []byte("`")) + 1
			if 0 < lastBacktick {
				// 把 ` 后面的字符调整到 info 节点
				n.NextSibling.AppendChild(&html.Node{Data: string(marker[lastBacktick:])})
				marker = marker[:lastBacktick]
			}
			tree.Context.Tip.IsFencedCodeBlock = true
			node.Type = ast.NodeCodeBlockFenceOpenMarker
			node.Tokens = marker
			node.CodeBlockFenceLen = len(marker)
			tree.Context.Tip.AppendChild(node)
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
			return
		case "super-block-open-marker":
			text := lute.domText(n)
			if "{{{" != strings.ReplaceAll(text, util.Caret, "") {
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(text)})
			} else {
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeSuperBlockOpenMarker})
			}
			return
		case "super-block-layout":
			layout := lute.domText(n)
			layout = strings.ReplaceAll(layout, parse.Zwsp, "")
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeSuperBlockLayoutMarker, Tokens: []byte(layout)})
			return
		case "super-block-close-marker":
			text := lute.domText(n)
			if "}}}" != strings.ReplaceAll(text, util.Caret, "") {
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte(text)})
			} else {
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeSuperBlockCloseMarker})
				if strings.Contains(text, util.Caret) {
					// 将插入符移到最后一个文本节点末尾
					paras := tree.Context.Tip.ChildrenByType(ast.NodeText)
					if length := len(paras); 0 < length {
						lastP := paras[length-1]
						lastP.Tokens = append(lastP.Tokens, util.CaretTokens...)
					}
				}
			}
			return
		case "heading-marker":
			text := lute.domText(n)
			if caretInMarker := strings.Contains(text, util.Caret); caretInMarker {
				caret := &html.Node{Type: html.TextNode, Data: util.Caret}
				n.InsertAfter(caret)
				text = strings.ReplaceAll(text, "#", "")
				text = strings.ReplaceAll(text, util.Caret, "")
				text = strings.TrimSpace(text)
				if 0 < len(text) {
					caret.Data = text + caret.Data
				}
			}
			return
		default:
			if nil != n.FirstChild && (atom.Audio == n.FirstChild.DataAtom || atom.Video == n.FirstChild.DataAtom) {
				node.Type = ast.NodeHTMLBlock
				node.Tokens = lute.domHTML(n.FirstChild)
				tree.Context.Tip.AppendChild(node)
				return
			}
		}

		text := lute.domText(n)
		node.Type = ast.NodeText
		node.Tokens = []byte(text)
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Div:
		switch dataType {
		case "super-block":
			node.Type = ast.NodeSuperBlock
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		case "math-block":
			node.Type = ast.NodeMathBlock
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		case "code-block":
			node.Type = ast.NodeCodeBlock
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		case "html-block":
			node.Type = ast.NodeHTMLBlock
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		case "yaml-front-matter":
			node.Type = ast.NodeYamlFrontMatter
			tree.Context.Tip.AppendChild(node)
			tree.Context.Tip = node
			defer tree.Context.ParentTip()
		case "block-ref-embed":
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
		}
	case atom.Font:
		text := lute.domText(n)
		inlineTree := parse.Inline("", []byte(text), tree.Context.ParseOption)
		appendNextToTip(inlineTree.Root.FirstChild.FirstChild, tree)
		return
	case atom.Details:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = lute.domHTML(n)
		node.Tokens = bytes.SplitAfter(node.Tokens, []byte("</summary>"))[0]
		tree.Context.Tip.AppendChild(node)
	case atom.Kbd:
		// kbd 标签由 code 标签构成节点
	case atom.Summary:
		return
	case atom.Audio, atom.Video:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = lute.domHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	default:
		if html.ElementNode == n.Type {
			break
		}

		node.Type = ast.NodeHTMLBlock
		node.Tokens = lute.domHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTByVditorIRBlockDOM(c, tree)
	}

	switch n.DataAtom {
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
	case atom.Details:
		tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeHTMLBlock, Tokens: []byte("</details>")})
	}
}

func (lute *Lute) setIAL(n *html.Node, node *ast.Node) (ialTokens []byte) {
	ialTokens = []byte("{: id=\"" + node.ID + "\"")
	if bookmark := lute.domAttrValue(n, "bookmark"); "" != bookmark {
		node.SetIALAttr("bookmark", bookmark)
		ialTokens = append(ialTokens, []byte(" bookmark=\""+bookmark+"\"")...)
	}
	if style := lute.domAttrValue(n, "style"); "" != style {
		node.SetIALAttr("style", style)
		ialTokens = append(ialTokens, []byte(" style=\""+style+"\"")...)
	}
	if name := lute.domAttrValue(n, "name"); "" != name {
		node.SetIALAttr("name", name)
		ialTokens = append(ialTokens, []byte(" name=\""+name+"\"")...)
	}
	if memo := lute.domAttrValue(n, "memo"); "" != memo {
		node.SetIALAttr("memo", memo)
		ialTokens = append(ialTokens, []byte(" memo=\""+memo+"\"")...)
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
