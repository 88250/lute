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
	"strings"
	"unicode"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
)

// HTML2Markdown 将 HTML 转换为 Markdown。
func (lute *Lute) HTML2Markdown(htmlStr string) (markdown string, err error) {
	//fmt.Println(htmlStr)
	// 将字符串解析为 DOM 树
	tree := lute.HTML2Tree(htmlStr)

	// 将 AST 进行 Markdown 格式化渲染
	var formatted []byte
	renderer := render.NewFormatRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.HTML2MdRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	formatted = renderer.Render()
	markdown = util.BytesToStr(formatted)
	return
}

// HTML2Tree 将 HTML 转换为 AST。
func (lute *Lute) HTML2Tree(dom string) (ret *parse.Tree) {
	htmlRoot := lute.parseHTML(dom)
	if nil == htmlRoot {
		return
	}

	// 调整 DOM 结构
	lute.adjustVditorDOM(htmlRoot)

	// 将 HTML 树转换为 Markdown AST
	ret = &parse.Tree{Name: "", Root: &ast.Node{Type: ast.NodeDocument}, Context: &parse.Context{ParseOption: lute.ParseOptions}}
	ret.Context.Tip = ret.Root
	for c := htmlRoot.FirstChild; nil != c; c = c.NextSibling {
		lute.genASTByDOM(c, ret)
	}

	// 调整树结构
	ast.Walk(ret.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			if ast.NodeList == n.Type {
				// ul.ul => ul.li.ul
				if nil != n.Parent && ast.NodeList == n.Parent.Type {
					previousLi := n.Previous
					if nil != previousLi {
						n.Unlink()
						previousLi.AppendChild(n)
					}
				}
			}
		}
		return ast.WalkContinue
	})
	return
}

// genASTByDOM 根据指定的 DOM 节点 n 进行深度优先遍历并逐步生成 Markdown 语法树 tree。
func (lute *Lute) genASTByDOM(n *html.Node, tree *parse.Tree) {
	if html.CommentNode == n.Type || atom.Meta == n.DataAtom {
		return
	}

	if "svg" == n.Namespace {
		return
	}

	dataRender := util.DomAttrValue(n, "data-render")
	if "1" == dataRender {
		return
	}

	class := util.DomAttrValue(n, "class")
	if strings.HasPrefix(class, "line-number") &&
		!strings.HasPrefix(class, "line-numbers" /* 简书代码块 https://github.com/siyuan-note/siyuan/issues/4361 */) {
		return
	}

	if 0 == n.DataAtom && html.ElementNode == n.Type { // 自定义标签
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lute.genASTByDOM(c, tree)
		}
		return
	}

	node := &ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(n.Data)}
	switch n.DataAtom {
	case 0:
		if nil != n.Parent && atom.A == n.Parent.DataAtom {
			node.Type = ast.NodeLinkText
		}
		if lute.parentIs(n, atom.Table) {
			if "" == strings.TrimSpace(n.Data) {
				node.Tokens = []byte(" ")
				tree.Context.Tip.AppendChild(node)
				break
			}
			node.Tokens = bytes.TrimSpace(node.Tokens)
			node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("\n"), []byte(""))
		}
		node.Tokens = bytes.ReplaceAll(node.Tokens, []byte{194, 160}, []byte{' '}) // 将 &nbsp; 转换为空格
		node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("\n"), []byte{' '})     // 将 \n 转换为空格 https://github.com/siyuan-note/siyuan/issues/6052
		if lute.ParseOptions.ProtyleWYSIWYG {
			node.Tokens = lex.EscapeProtyleMarkers(node.Tokens)
		} else {
			node.Tokens = lex.EscapeMarkers(node.Tokens)
		}
		if nil != tree.Context.Tip && tree.Context.Tip.IsBlock() && nil != n.Parent && atom.Span != n.Parent.DataAtom && 1 > len(lex.TrimWhitespace(node.Tokens)) {
			// 块级节点下非 span 包裹需要忽略空白
			// 剪藏时列表下方块缩进不正确 https://github.com/siyuan-note/siyuan/issues/6650
			return
		}

		if nil != n.Parent && atom.Span == n.Parent.DataAtom && ("" != util.DomAttrValue(n.Parent, "class")) {
			if lastc := tree.Context.Tip.LastChild; nil == lastc || (ast.NodeText == lastc.Type && !bytes.HasSuffix(lastc.Tokens, []byte("**"))) {
				node.Tokens = []byte("**" + util.BytesToStr(node.Tokens) + "**")
			}
		}
		tree.Context.Tip.AppendChild(node)
	case atom.P, atom.Div, atom.Section:
		if lute.parentIs(n, atom.Table) || ast.NodeLink == tree.Context.Tip.Type {
			break
		}

		if atom.Div == n.DataAtom {
			// 解析 GitHub 语法高亮代码块
			class := util.DomAttrValue(n, "class")
			language := ""
			if strings.Contains(class, "-source-") {
				language = class[strings.LastIndex(class, "-source-")+len("-source-"):]
			} else if strings.Contains(class, "-text-html-basic") {
				language = "html"
			}
			if "" != language {
				node.Type = ast.NodeCodeBlock
				node.IsFencedCodeBlock = true
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: util.StrToBytes("```"), CodeBlockFenceLen: 3})
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker})
				buf := &bytes.Buffer{}
				node.LastChild.CodeBlockInfo = []byte(language)
				buf.WriteString(util.DomText(n))
				content := &ast.Node{Type: ast.NodeCodeBlockCode, Tokens: buf.Bytes()}
				node.AppendChild(content)
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: util.StrToBytes("```"), CodeBlockFenceLen: 3})
				tree.Context.Tip.AppendChild(node)
				return
			}
		}
		node.Type = ast.NodeParagraph
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		if ast.NodeLink == tree.Context.Tip.Type {
			break
		}

		node.Type = ast.NodeHeading
		node.HeadingLevel = int(node.Tokens[1] - byte('0'))
		node.AppendChild(&ast.Node{Type: ast.NodeHeadingC8hMarker, Tokens: util.StrToBytes(strings.Repeat("#", node.HeadingLevel))})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Hr:
		node.Type = ast.NodeThematicBreak
		tree.Context.Tip.AppendChild(node)
	case atom.Blockquote:
		node.Type = ast.NodeBlockquote
		node.AppendChild(&ast.Node{Type: ast.NodeBlockquoteMarker, Tokens: util.StrToBytes(">")})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Ol, atom.Ul:
		node.Type = ast.NodeList
		node.ListData = &ast.ListData{}
		if atom.Ol == n.DataAtom {
			node.ListData.Typ = 1
		}
		node.ListData.Tight = true
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Li:
		node.Type = ast.NodeListItem
		marker := util.DomAttrValue(n, "data-marker")
		var bullet byte
		if "" == marker {
			if nil != n.Parent && atom.Ol == n.Parent.DataAtom {
				start := util.DomAttrValue(n.Parent, "start")
				if "" == start {
					marker = "1."
				} else {
					marker = start + "."
				}
			} else {
				marker = "*"
				bullet = marker[0]
			}
		} else {
			if nil != n.Parent && "1." != marker && atom.Ol == n.Parent.DataAtom && nil != n.Parent.Parent && (atom.Ol == n.Parent.Parent.DataAtom || atom.Ul == n.Parent.Parent.DataAtom) {
				// 子有序列表必须从 1 开始
				marker = "1."
			}
		}
		node.ListData = &ast.ListData{Marker: []byte(marker), BulletChar: bullet}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Pre:
		if firstc := n.FirstChild; nil != firstc {
			if html.TextNode == firstc.Type || atom.Span == firstc.DataAtom || atom.Code == firstc.DataAtom {
				node.Type = ast.NodeCodeBlock
				node.IsFencedCodeBlock = true
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: util.StrToBytes("```"), CodeBlockFenceLen: 3})
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker})
				if atom.Code == firstc.DataAtom || atom.Span == firstc.DataAtom {
					class := util.DomAttrValue(firstc, "class")
					if !strings.Contains(class, "language-") {
						class = util.DomAttrValue(n, "class")
					}
					if strings.Contains(class, "language-") {
						language := class[strings.Index(class, "language-")+len("language-"):]
						language = strings.Split(language, " ")[0]
						node.LastChild.CodeBlockInfo = []byte(language)
					}
				}

				if atom.Code == n.FirstChild.DataAtom {
					if nil != n.FirstChild.NextSibling && atom.Code == n.FirstChild.NextSibling.DataAtom {
						// pre.code code 每个 code 为一行的结构，需要在 code 中间插入换行
						for c := n.FirstChild.NextSibling; nil != c; c = c.NextSibling {
							c.InsertBefore(&html.Node{DataAtom: atom.Br})
						}
					}
					if nil != n.FirstChild.FirstChild && atom.Ol == n.FirstChild.FirstChild.DataAtom {
						// CSDN 代码块：pre.code.ol.li
						for li := n.FirstChild.FirstChild.FirstChild; nil != li; li = li.NextSibling {
							if li != n.FirstChild.FirstChild.FirstChild {
								li.InsertBefore(&html.Node{DataAtom: atom.Br})
							}
						}
					}
					if nil != n.LastChild && atom.Ul == n.LastChild.DataAtom {
						// CSDN 代码块：pre.code,ul
						n.LastChild.Unlink() // 去掉最后一个代码行号子块 https://github.com/siyuan-note/siyuan/issues/5564
					}
				}

				buf := &bytes.Buffer{}
				buf.WriteString(util.DomText(n))
				content := &ast.Node{Type: ast.NodeCodeBlockCode, Tokens: buf.Bytes()}
				node.AppendChild(content)
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: util.StrToBytes("```"), CodeBlockFenceLen: 3})
				tree.Context.Tip.AppendChild(node)
			} else {
				node.Type = ast.NodeHTMLBlock
				node.Tokens = util.DomHTML(n)
				tree.Context.Tip.AppendChild(node)
			}
		}
		return
	case atom.Em, atom.I:
		text := util.DomText(n)
		if "" == strings.TrimSpace(text) {
			break
		}

		if nil != tree.Context.Tip.LastChild && (ast.NodeStrong == tree.Context.Tip.LastChild.Type || ast.NodeEmphasis == tree.Context.Tip.LastChild.Type) {
			// 在两个相邻的加粗或者斜体之间插入零宽空格，避免标记符重复
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(editor.Zwsp)})
		}

		node.Type = ast.NodeEmphasis
		marker := "*"
		node.AppendChild(&ast.Node{Type: ast.NodeEmA6kOpenMarker, Tokens: util.StrToBytes(marker)})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Strong, atom.B:
		text := util.DomText(n)
		if "" == strings.TrimSpace(text) {
			break
		}

		if nil != tree.Context.Tip.LastChild && (ast.NodeStrong == tree.Context.Tip.LastChild.Type || ast.NodeEmphasis == tree.Context.Tip.LastChild.Type) {
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(editor.Zwsp)})
		}

		node.Type = ast.NodeStrong
		marker := "**"
		node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kOpenMarker, Tokens: util.StrToBytes(marker)})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Code:
		if nil == n.FirstChild {
			return
		}

		if nil != tree.Context.Tip.LastChild && ast.NodeCodeSpan == tree.Context.Tip.LastChild.Type {
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(editor.Zwsp)})
		}

		code := util.DomHTML(n)
		if bytes.Contains(code, []byte(">")) {
			code = code[bytes.Index(code, []byte(">"))+1:]
		}
		code = bytes.TrimSuffix(code, []byte("</code>"))
		unescaped := html.UnescapeString(string(code))
		code = []byte(unescaped)
		content := &ast.Node{Type: ast.NodeCodeSpanContent, Tokens: code}
		node.Type = ast.NodeCodeSpan
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanOpenMarker, Tokens: []byte("`")})
		node.AppendChild(content)
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanCloseMarker, Tokens: []byte("`")})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
		return
	case atom.Br:
		if ast.NodeLink == tree.Context.Tip.Type {
			break
		}

		if tree.Context.ParseOption.ProtyleWYSIWYG && lute.parentIs(n, atom.Table) {
			node.Type = ast.NodeBr
		} else {
			node.Type = ast.NodeHardBreak
			node.Tokens = util.StrToBytes("\n")
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.A:
		node.Type = ast.NodeLink
		text := util.DomText(n)
		if "" == text && nil != n.Parent && (atom.H1 == n.Parent.DataAtom || atom.H2 == n.Parent.DataAtom || atom.H3 == n.Parent.DataAtom || atom.H4 == n.Parent.DataAtom || atom.H5 == n.Parent.DataAtom || atom.H6 == n.Parent.DataAtom) {
			// 丢弃标题中文本为空的链接，这样的链接可能是锚点 https://github.com/Vanessa219/vditor/issues/359
			return
		}
		if "" == text && nil == n.FirstChild {
			// 剪藏时过滤空的超链接 https://github.com/siyuan-note/siyuan/issues/5686
			return
		}

		node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Img:
		imgClass := util.DomAttrValue(n, "class")
		imgAlt := util.DomAttrValue(n, "alt")
		if "emoji" == imgClass {
			node.Type = ast.NodeEmoji
			emojiImg := &ast.Node{Type: ast.NodeEmojiImg, Tokens: tree.EmojiImgTokens(imgAlt, util.DomAttrValue(n, "src"))}
			emojiImg.AppendChild(&ast.Node{Type: ast.NodeEmojiAlias, Tokens: util.StrToBytes(":" + imgAlt + ":")})
			node.AppendChild(emojiImg)
		} else {
			node.Type = ast.NodeImage
			node.AppendChild(&ast.Node{Type: ast.NodeBang})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
			if "" != imgAlt {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkText, Tokens: util.StrToBytes(imgAlt)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			src := util.DomAttrValue(n, "src")
			if strings.HasPrefix(src, "data:image") {
				// 处理可能存在的预加载情况
				if dataSrc := util.DomAttrValue(n, "data-src"); "" != dataSrc {
					src = dataSrc
				}
			}
			if "" == src {
				// 处理使用 srcset 属性的情况
				if srcset := util.DomAttrValue(n, "srcset"); "" != srcset {
					if strings.Contains(srcset, ",") {
						src = strings.Split(srcset, ",")[len(strings.Split(srcset, ","))-1]
						src = strings.TrimSpace(src)
						if strings.Contains(src, " ") {
							src = strings.TrimSpace(strings.Split(src, " ")[0])
						}
					} else {
						src = strings.TrimSpace(src)
						if strings.Contains(src, " ") {
							src = strings.TrimSpace(strings.Split(srcset, " ")[0])
						}
					}
				}
			}
			node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: util.StrToBytes(src)})
			linkTitle := util.DomAttrValue(n, "title")
			if "" != linkTitle {
				node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
				node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: []byte(linkTitle)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
		}

		if ast.NodeDocument == tree.Context.Tip.Type {
			p := &ast.Node{Type: ast.NodeParagraph}
			tree.Context.Tip.AppendChild(p)
			tree.Context.Tip = p
			defer tree.Context.ParentTip()
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Input:
		node.Type = ast.NodeTaskListItemMarker
		node.TaskListItemChecked = lute.hasAttr(n, "checked")
		tree.Context.Tip.AppendChild(node)
		if nil != node.Parent.Parent {
			if nil == node.Parent.Parent.ListData {
				node.Parent.Parent.ListData = &ast.ListData{Typ: 3}
			} else {
				node.Parent.Parent.ListData.Typ = 3
			}
		}
	case atom.Del, atom.S, atom.Strike:
		node.Type = ast.NodeStrikethrough
		marker := "~"
		node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1OpenMarker, Tokens: util.StrToBytes(marker)})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Mark:
		node.Type = ast.NodeMark
		marker := "=="
		node.AppendChild(&ast.Node{Type: ast.NodeMark2OpenMarker, Tokens: util.StrToBytes(marker)})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Sup:
		node.Type = ast.NodeSup
		node.AppendChild(&ast.Node{Type: ast.NodeSupOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Sub:
		node.Type = ast.NodeSub
		node.AppendChild(&ast.Node{Type: ast.NodeSubOpenMarker})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Table:
		node.Type = ast.NodeTable
		var tableAligns []int
		if nil != n.FirstChild && nil != n.FirstChild.FirstChild && nil != n.FirstChild.FirstChild.FirstChild {
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
		if nil == n.FirstChild {
			break
		}
		table := n.Parent.Parent
		node.Type = ast.NodeTableRow
		if atom.Thead != table.FirstChild.DataAtom && n == n.Parent.FirstChild {
			// 补全 thread 节点
			thead := &ast.Node{Type: ast.NodeTableHead}
			tree.Context.Tip.AppendChild(thead)
			tree.Context.Tip = thead
			defer tree.Context.ParentTip()
		}
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
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Colgroup, atom.Col:
		return
	case atom.Span:
		if nil == n.FirstChild {
			return
		}
	case atom.Font:
		node.Type = ast.NodeText
		node.Tokens = []byte(util.DomText(n))
		node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("\n"), nil)
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Details:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = util.DomHTML(n)
		node.Tokens = bytes.SplitAfter(node.Tokens, []byte("</summary>"))[0]
		tree.Context.Tip.AppendChild(node)
	case atom.Summary:
		return
	case atom.Iframe, atom.Audio, atom.Video:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = util.DomHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Noscript:
		return
	case atom.Figcaption:
		node.Type = ast.NodeParagraph
		node.AppendChild(&ast.Node{Type: ast.NodeHardBreak})
		node.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Figure:
		node.Type = ast.NodeParagraph
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	default:
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		lute.genASTByDOM(c, tree)
	}

	switch n.DataAtom {
	case atom.Em, atom.I:
		marker := "*"
		node.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker, Tokens: util.StrToBytes(marker)})
		appendSpace(n, tree, lute)
	case atom.Strong, atom.B:
		marker := "**"
		node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kCloseMarker, Tokens: util.StrToBytes(marker)})
		appendSpace(n, tree, lute)
	case atom.A:
		node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
		node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
		node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: util.StrToBytes(util.DomAttrValue(n, "href"))})
		linkTitle := util.DomAttrValue(n, "title")
		if "" != linkTitle {
			node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace})
			node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: util.StrToBytes(linkTitle)})
		}
		node.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
	case atom.Del, atom.S, atom.Strike:
		marker := "~"
		node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough1CloseMarker, Tokens: util.StrToBytes(marker)})
		appendSpace(n, tree, lute)
	case atom.Mark:
		marker := "=="
		node.AppendChild(&ast.Node{Type: ast.NodeMark2CloseMarker, Tokens: util.StrToBytes(marker)})
		appendSpace(n, tree, lute)
	case atom.Sup:
		node.AppendChild(&ast.Node{Type: ast.NodeSupCloseMarker})
		appendSpace(n, tree, lute)
	case atom.Sub:
		node.AppendChild(&ast.Node{Type: ast.NodeSubCloseMarker})
		appendSpace(n, tree, lute)
	case atom.Details:
		tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeHTMLBlock, Tokens: []byte("</details>")})
	}
}

func appendSpace(n *html.Node, tree *parse.Tree, lute *Lute) {
	if nil != n.NextSibling {
		if nextText := util.DomText(n.NextSibling); "" != nextText {
			if runes := []rune(nextText); !unicode.IsSpace(runes[0]) {
				if unicode.IsPunct(runes[0]) || unicode.IsSymbol(runes[0]) {
					tree.Context.Tip.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: []byte(editor.Zwsp)})
					tree.Context.Tip.InsertAfter(&ast.Node{Type: ast.NodeText, Tokens: []byte(editor.Zwsp)})
					return
				}

				if curText := util.DomText(n); "" != curText {
					runes = []rune(curText)
					if lastC := runes[len(runes)-1]; unicode.IsPunct(lastC) || unicode.IsSymbol(lastC) {
						text := tree.Context.Tip.ChildByType(ast.NodeText)
						if nil != text {
							text.Tokens = append([]byte(editor.Zwsp), text.Tokens...)
							text.Tokens = append(text.Tokens, []byte(editor.Zwsp)...)
						}
						return
					}

					spaces := lute.prefixSpaces(curText)
					if "" != spaces {
						previous := tree.Context.Tip.Previous
						if nil != previous {
							if ast.NodeText == previous.Type {
								previous.Tokens = append(previous.Tokens, util.StrToBytes(spaces)...)
							} else {
								previous.InsertAfter(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(spaces)})
							}
						} else {
							tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(spaces)})
						}

						text := tree.Context.Tip.ChildByType(ast.NodeText)
						text.Tokens = bytes.TrimLeft(text.Tokens, " \u0160")
					}
					spaces = lute.suffixSpaces(curText)
					if "" != spaces {
						texts := tree.Context.Tip.ChildrenByType(ast.NodeText)
						if 0 < len(texts) {
							text := texts[len(texts)-1]
							text.Tokens = bytes.TrimRight(text.Tokens, " \u0160")
							if 1 > len(text.Tokens) {
								text.Unlink()
							}
						}
						if nil != n.NextSibling {
							if html.TextNode == n.NextSibling.Type {
								n.NextSibling.Data = spaces + n.NextSibling.Data
							} else {
								tree.Context.Tip.InsertAfter(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(spaces)})
							}
						} else {
							tree.Context.Tip.InsertAfter(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(spaces)})
						}
					}
				}
			}
		}
	}
}
