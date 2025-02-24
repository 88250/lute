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
	"path"
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

	if strings.Contains(class, "mw-editsection") {
		// 忽略 Wikipedia [编辑] Do not clip the `Edit` element next to Wikipedia headings https://github.com/siyuan-note/siyuan/issues/11600
		return
	}

	if strings.Contains(class, "citation-comment") {
		// 忽略 Wikipedia 引用中的注释 https://github.com/siyuan-note/siyuan/issues/11640
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
		class := util.DomAttrValue(n.PrevSibling, "class")
		if "fn__space5" == class {
			// 链滴剪藏图片时多了长宽显示 https://github.com/siyuan-note/siyuan/issues/10987
			return
		}

		if nil != n.Parent && atom.A == n.Parent.DataAtom {
			node.Type = ast.NodeLinkText
		}

		// 将 \n空格空格* 转换为\n
		for strings.Contains(string(node.Tokens), "\n  ") {
			node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("\n  "), []byte("\n "))
		}
		node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("\n "), []byte("\n"))
		node.Tokens = bytes.Trim(node.Tokens, "\t\n")

		if lute.parentIs(n, atom.Table) {
			if "\n" == n.Data {
				if nil == tree.Context.Tip.FirstChild || nil == n.NextSibling {
					break
				}

				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeBr})
				break
			} else {
				if "" == strings.TrimSpace(n.Data) {
					node.Tokens = []byte(" ")
					tree.Context.Tip.AppendChild(node)
					break
				}
			}

			node.Tokens = bytes.TrimSpace(node.Tokens)
			node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("\n"), []byte(" "))
		}
		node.Tokens = bytes.ReplaceAll(node.Tokens, []byte{194, 160}, []byte{' '}) // 将 &nbsp; 转换为空格

		node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("\n"), []byte{' '}) // 将 \n 转换为空格 https://github.com/siyuan-note/siyuan/issues/6052
		if ast.NodeStrong == tree.Context.Tip.Type ||
			ast.NodeEmphasis == tree.Context.Tip.Type ||
			ast.NodeStrikethrough == tree.Context.Tip.Type ||
			ast.NodeMark == tree.Context.Tip.Type ||
			ast.NodeSup == tree.Context.Tip.Type ||
			ast.NodeSub == tree.Context.Tip.Type {
			if bytes.HasPrefix(node.Tokens, []byte(" ")) || 1 > len(bytes.TrimSpace(node.Tokens)) {
				if nil != tree.Context.Tip.LastChild && tree.Context.Tip.LastChild.IsMarker() {
					node.Tokens = append([]byte(editor.Zwsp), node.Tokens...)
				}
			}
			if bytes.HasSuffix(node.Tokens, []byte(" ")) && nil == n.NextSibling {
				node.Tokens = append(node.Tokens, []byte(editor.Zwsp)...)
			}
		}

		if nil != n.Parent && atom.Span == n.Parent.DataAtom && 0 == len(n.Parent.Attr) {
			// 按原文解析，不处理转义
		} else {
			if lute.ParseOptions.ProtyleWYSIWYG {
				node.Tokens = lex.EscapeProtyleMarkers(node.Tokens)
			} else {
				node.Tokens = lex.EscapeCommonMarkers(node.Tokens)
				if lute.parentIs(n, atom.Table) {
					node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("\\|"), []byte("|"))
					node.Tokens = bytes.ReplaceAll(node.Tokens, []byte("|"), []byte("\\|"))
				}
			}
		}

		if 1 > len(node.Tokens) {
			return
		}

		if ast.NodeListItem == tree.Context.Tip.Type {
			p := &ast.Node{Type: ast.NodeParagraph}
			p.AppendChild(node)
			tree.Context.Tip.AppendChild(p)
			tree.Context.Tip = p
		} else {
			tree.Context.Tip.AppendChild(node)
		}
	case atom.P, atom.Div, atom.Section, atom.Dd:
		if ast.NodeLink == tree.Context.Tip.Type {
			break
		}

		if lute.parentIs(n, atom.Table) {
			if nil != n.PrevSibling && strings.Contains(n.PrevSibling.Data, "\n") {
				break
			}

			if nil != n.NextSibling && strings.Contains(n.NextSibling.Data, "\n") {
				break
			}

			if nil == tree.Context.Tip.FirstChild {
				break
			}

			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeBr})
			break
		}

		if ast.NodeHeading == tree.Context.Tip.Type {
			// h 下存在 div/p/section 则忽略分块
			break
		}

		class := util.DomAttrValue(n, "class")
		if atom.Div == n.DataAtom || atom.Section == n.DataAtom {
			// 解析 GitHub 语法高亮代码块
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
				tokens := buf.Bytes()
				tokens = bytes.ReplaceAll(tokens, []byte("\u00A0"), []byte(" "))
				tokens = bytes.TrimSuffix(tokens, []byte("\n"+editor.Zwsp))
				content := &ast.Node{Type: ast.NodeCodeBlockCode, Tokens: tokens}
				node.AppendChild(content)
				node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: util.StrToBytes("```"), CodeBlockFenceLen: 3})
				tree.Context.Tip.AppendChild(node)
				return
			}

			// The browser extension supports CSDN formula https://github.com/siyuan-note/siyuan/issues/5624
			if strings.Contains(class, "MathJax") && nil != n.NextSibling && atom.Script == n.NextSibling.DataAtom && strings.Contains(util.DomAttrValue(n.NextSibling, "type"), "math/tex") {
				tex := util.DomText(n.NextSibling)
				appendMathBlock(tree, tex)
				n.NextSibling.Unlink()
				return
			}

			// The browser extension supports Wikipedia formula clipping https://github.com/siyuan-note/siyuan/issues/11583
			if tex := strings.TrimSpace(util.DomAttrValue(n, "data-tex")); "" != tex {
				appendMathBlock(tree, tex)
				return
			}
		}

		if strings.Contains(strings.ToLower(class), "mathjax") {
			return
		}

		if "" == strings.TrimSpace(util.DomText(n)) {
			for { // 这里用 for 是为了简化实现
				if util.DomExistChildByType(n, atom.Img, atom.Picture, atom.Annotation, atom.Iframe, atom.Video, atom.Audio, atom.Source, atom.Canvas, atom.Svg, atom.Math) {
					break
				}

				// span 可能是 TextMark 元素，也可能是公式，其他情况则忽略
				spans := util.DomChildrenByType(n, atom.Span)
				if 0 < len(spans) {
					span := spans[0]
					if "" != util.DomAttrValue(span, "data-type") || "" != util.DomAttrValue(span, "data-tex") {
						break
					}
				}
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
		firstc := n.FirstChild
		if nil == firstc {
			return
		}

		codes := util.DomChildrenByType(n, atom.Code)
		if 0 < len(codes) {
			// 删除第一个 code 之前的标签
			unlinks := []*html.Node{}
			for prev := codes[0].PrevSibling; nil != prev; prev = prev.PrevSibling {
				unlinks = append(unlinks, prev)
			}
			for _, unlink := range unlinks {
				unlink.Unlink()
			}
			firstc = n.FirstChild
			if nil == firstc {
				return
			}
		}

		if atom.Em == firstc.DataAtom && nil != firstc.NextSibling && atom.Em == firstc.NextSibling.DataAtom {
			// pre.em,em,code 的情况，这两个 em 是“复制代码”和“隐藏代码” https://github.com/siyuan-note/siyuan/issues/13026

			if 0 < len(codes) {
				firstc = codes[0]
				if nil != firstc {
					unlinks := []*html.Node{}
					for prev := firstc.PrevSibling; nil != prev; prev = prev.PrevSibling {
						unlinks = append(unlinks, prev)
					}
					for _, unlink := range unlinks {
						unlink.Unlink()
					}
				}
			}
		}

		if atom.Div == firstc.DataAtom {
			if nil == firstc.NextSibling {
				if 1 == len(codes) {
					code := codes[0]
					// pre 下只有一个 div，且 div 下只有一个 code，那么将 pre.div 替换为 pre.code https://github.com/siyuan-note/siyuan/issues/11131
					code.Unlink()
					n.AppendChild(code)
					firstc.Unlink()
					firstc = n.FirstChild
				}
			} else {
				// pre 下全是 div，每个 div 为一行代码 https://github.com/siyuan-note/siyuan/issues/14195
				// 将其转换为 pre.code， code, ... code，每个 div 为一行代码，然后交由后续处理
				var unlinks, codes []*html.Node
				for div := firstc; nil != div; div = div.NextSibling {
					code := &html.Node{Data: "code", DataAtom: atom.Code, Type: html.ElementNode}
					for child := div.FirstChild; nil != child; child = div.FirstChild {
						child.Unlink()
						code.AppendChild(child)
					}
					codes = append(codes, code)
					unlinks = append(unlinks, div)
				}
				for _, unlink := range unlinks {
					unlink.Unlink()
				}
				for _, code := range codes {
					n.AppendChild(code)
				}
				firstc = n.FirstChild
			}
		}

		// 改进两种 pre.ol.li 的代码块解析 https://github.com/siyuan-note/siyuan/issues/11296
		// 第一种：将 pre.ol.li.p.span, span, ... span 转换为 pre.ol.li.p.code, code, ... code，然后交由第二种处理
		span2Code := false
		if atom.Ol == firstc.DataAtom && nil == firstc.NextSibling && nil != firstc.FirstChild && atom.Li == firstc.FirstChild.DataAtom &&
			nil != firstc.FirstChild.FirstChild && atom.P == firstc.FirstChild.FirstChild.DataAtom &&
			nil != firstc.FirstChild.FirstChild.FirstChild && (atom.Span == firstc.FirstChild.FirstChild.FirstChild.DataAtom || html.TextNode == firstc.FirstChild.FirstChild.FirstChild.Type) {
			for li := firstc.FirstChild; nil != li; li = li.NextSibling {
				code := &html.Node{Data: "code", DataAtom: atom.Code, Type: html.ElementNode}

				var spans []*html.Node
				if nil == li.FirstChild {
					continue
				}

				for span := li.FirstChild.FirstChild; nil != span; span = span.NextSibling {
					spans = append(spans, span)
				}

				for _, span := range spans {
					span.Unlink()
					code.AppendChild(span)
				}

				li.FirstChild.AppendChild(code)
				span2Code = true
			}
		}
		// 第二种：将 pre.ol.li.p.code, code, ... code 转换为 pre.code, code, ... code，然后交由后续处理
		if atom.Ol == firstc.DataAtom && nil == firstc.NextSibling && nil != firstc.FirstChild && atom.Li == firstc.FirstChild.DataAtom &&
			nil != firstc.FirstChild.FirstChild && atom.P == firstc.FirstChild.FirstChild.DataAtom &&
			nil != firstc.FirstChild.FirstChild.FirstChild && atom.Code == firstc.FirstChild.FirstChild.FirstChild.DataAtom {
			var lis, codes []*html.Node
			for li := firstc.FirstChild; nil != li; li = li.NextSibling {
				lis = append(lis, li)
				if nil == li.FirstChild {
					continue
				}
				codes = append(codes, li.FirstChild.FirstChild)
			}
			for _, li := range lis {
				li.Unlink()
			}
			for _, code := range codes {
				code.Unlink()
				n.AppendChild(code)
			}
			firstc.Unlink()
			firstc = n.FirstChild
		}

		if html.TextNode == firstc.Type || atom.Span == firstc.DataAtom || atom.Code == firstc.DataAtom || atom.Section == firstc.DataAtom || atom.Pre == firstc.DataAtom || atom.A == firstc.DataAtom || atom.P == firstc.DataAtom {
			node.Type = ast.NodeCodeBlock
			node.IsFencedCodeBlock = true
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: util.StrToBytes("```"), CodeBlockFenceLen: 3})
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker})
			if atom.Code == firstc.DataAtom || atom.Span == firstc.DataAtom || atom.A == firstc.DataAtom {
				class := util.DomAttrValue(firstc, "class")
				if nil != n.Parent && nil != n.Parent.Parent {
					parentClass := util.DomAttrValue(n.Parent.Parent, "class")
					if strings.Contains(parentClass, "language-") {
						class = parentClass
					}
				}

				if !strings.Contains(class, "language-") {
					class = util.DomAttrValue(n, "class")
				}
				if strings.Contains(class, "language-") {
					language := class[strings.Index(class, "language-")+len("language-"):]
					language = strings.Split(language, " ")[0]
					if "fallback" != language && "chroma" != language {
						node.LastChild.CodeBlockInfo = []byte(language)
					}
				} else {
					if atom.Code == firstc.DataAtom && !span2Code {
						class := util.DomAttrValue(firstc, "class")
						if !strings.Contains(class, " ") {
							node.LastChild.CodeBlockInfo = []byte(class)
						}
					}
				}

				if 1 > len(node.LastChild.CodeBlockInfo) {
					class := util.DomAttrValue(n, "class")
					if !strings.Contains(class, " ") && "fallback" != class && "chroma" != class {
						node.LastChild.CodeBlockInfo = []byte(class)
					}
				}

				if 1 > len(node.LastChild.CodeBlockInfo) {
					lang := util.DomAttrValue(n, "data-language")
					if !strings.Contains(lang, " ") {
						node.LastChild.CodeBlockInfo = []byte(lang)
					}
				}

				if bytes.ContainsAny(node.LastChild.CodeBlockInfo, "-_ ") {
					node.LastChild.CodeBlockInfo = nil
				}
			}

			if atom.Code == firstc.DataAtom {
				if nil != firstc.NextSibling && atom.Code == firstc.NextSibling.DataAtom {
					// pre.code code 每个 code 为一行的结构，需要在 code 中间插入换行
					for c := firstc.NextSibling; nil != c; c = c.NextSibling {
						c.InsertBefore(&html.Node{DataAtom: atom.Br})
					}
				}
				if nil != firstc.FirstChild && atom.Ol == firstc.FirstChild.DataAtom {
					// CSDN 代码块：pre.code.ol.li
					for li := firstc.FirstChild.FirstChild; nil != li; li = li.NextSibling {
						if li != firstc.FirstChild.FirstChild {
							li.InsertBefore(&html.Node{DataAtom: atom.Br})
						}
					}
				}
				if nil != n.LastChild && atom.Ul == n.LastChild.DataAtom {
					// CSDN 代码块：pre.code,ul
					n.LastChild.Unlink() // 去掉最后一个代码行号子块 https://github.com/siyuan-note/siyuan/issues/5564
				}
			}

			if atom.Pre == firstc.DataAtom && nil != firstc.FirstChild {
				// pre.code code 每个 code 为一行的结构，需要在 code 中间插入换行
				for c := firstc.FirstChild.NextSibling; nil != c; c = c.NextSibling {
					c.InsertBefore(&html.Node{DataAtom: atom.Br})
				}
			}

			if atom.P == firstc.DataAtom {
				// 避免下面 util.DomText 把 p 转换为两个换行
				firstc.DataAtom = atom.Div
			}

			buf := &bytes.Buffer{}
			buf.WriteString(util.DomText(n))
			tokens := buf.Bytes()
			tokens = bytes.ReplaceAll(tokens, []byte("\u00A0"), []byte(" "))
			tokens = bytes.TrimSuffix(tokens, []byte("\n"+editor.Zwsp))
			content := &ast.Node{Type: ast.NodeCodeBlockCode, Tokens: tokens}
			node.AppendChild(content)
			node.AppendChild(&ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: util.StrToBytes("```"), CodeBlockFenceLen: 3})

			if tree.Context.Tip.ParentIs(ast.NodeTable) {
				// 如果表格中只有一行一列，那么丢弃表格直接使用代码块
				// Improve HTML parsing code blocks https://github.com/siyuan-note/siyuan/issues/11068
				for table := tree.Context.Tip.Parent; nil != table; table = table.Parent {
					if ast.NodeTable == table.Type {
						if nil != table.FirstChild && table.FirstChild == table.LastChild && ast.NodeTableHead == table.FirstChild.Type &&
							table.FirstChild.FirstChild == table.FirstChild.LastChild &&
							nil != table.FirstChild.FirstChild.FirstChild && ast.NodeTableCell == table.FirstChild.FirstChild.FirstChild.Type {
							table.InsertBefore(node)
							table.Unlink()
							tree.Context.Tip = node

							parent := n.Parent
							for i := 0; i < 32; i++ {
								if nil == parent {
									break
								}

								class := util.DomAttrValue(parent, "class")
								if strings.Contains(class, "language-") {
									node.ChildByType(ast.NodeCodeBlockFenceInfoMarker).CodeBlockInfo = []byte(class[strings.Index(class, "language-")+len("language-"):])
									break
								} else if strings.Contains(class, "highlight ") {
									node.ChildByType(ast.NodeCodeBlockFenceInfoMarker).CodeBlockInfo = []byte(class[strings.Index(class, "highlight ")+len("highlight "):])
									break
								}
								parent = parent.Parent
							}

							return
						}
					}
				}

				// 表格中不支持添加块级元素，所以这里只能将其转换为多个行级代码元素
				lines := bytes.Split(content.Tokens, []byte("\n"))
				for i, line := range lines {
					if 0 < len(line) {
						code := &ast.Node{Type: ast.NodeCodeSpan}
						if bytes.Contains(line, []byte("`")) {
							node.CodeMarkerLen = 2
						}
						code.AppendChild(&ast.Node{Type: ast.NodeCodeSpanOpenMarker, Tokens: []byte("`")})
						code.AppendChild(&ast.Node{Type: ast.NodeCodeSpanContent, Tokens: line})
						code.AppendChild(&ast.Node{Type: ast.NodeCodeSpanCloseMarker, Tokens: []byte("`")})
						tree.Context.Tip.AppendChild(code)
						if i < len(lines)-1 {
							if tree.Context.ParseOption.ProtyleWYSIWYG {
								tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeBr})
							} else {
								tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeHardBreak, Tokens: []byte("\n")})
							}
						}
					}
				}
			} else {
				tree.Context.Tip.AppendChild(node)
			}
		} else {
			node.Type = ast.NodeHTMLBlock
			node.Tokens = util.DomHTML(n)
			tree.Context.Tip.AppendChild(node)
		}
		return
	case atom.Em, atom.I:
		text := util.DomText(n)
		if "" == strings.TrimSpace(text) {
			break
		}

		if ast.NodeEmphasis == tree.Context.Tip.Type || tree.Context.Tip.ParentIs(ast.NodeEmphasis) {
			break
		}

		if nil != tree.Context.Tip.LastChild && (ast.NodeStrong == tree.Context.Tip.LastChild.Type || ast.NodeEmphasis == tree.Context.Tip.LastChild.Type) {
			// 在两个相邻的加粗或者斜体之间插入零宽空格，避免标记符重复
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(editor.Zwsp)})
		}

		if !lute.ParseOptions.InlineAsterisk || !lute.ParseOptions.InlineUnderscore {
			node.Type = ast.NodeHTMLTag
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagOpen, Tokens: util.StrToBytes("<em>")})
		} else {
			node.Type = ast.NodeEmphasis
			node.AppendChild(&ast.Node{Type: ast.NodeEmA6kOpenMarker, Tokens: util.StrToBytes("*")})
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Strong, atom.B:
		text := util.DomText(n)
		if "" == strings.TrimSpace(text) {
			break
		}

		if ast.NodeStrong == tree.Context.Tip.Type || tree.Context.Tip.ParentIs(ast.NodeStrong) {
			break
		}

		if nil != tree.Context.Tip.LastChild && (ast.NodeStrong == tree.Context.Tip.LastChild.Type || ast.NodeEmphasis == tree.Context.Tip.LastChild.Type) {
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(editor.Zwsp)})
		}

		if !lute.ParseOptions.InlineAsterisk || !lute.ParseOptions.InlineUnderscore {
			node.Type = ast.NodeHTMLTag
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagOpen, Tokens: util.StrToBytes("<strong>")})
		} else {
			node.Type = ast.NodeStrong
			node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kOpenMarker, Tokens: util.StrToBytes("**")})
		}
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

		allSpan := true
		for c := n.FirstChild; nil != c; c = c.NextSibling {
			if html.TextNode == c.Type {
				continue
			}
			if atom.Em == c.DataAtom || atom.Strong == c.DataAtom {
				// https://github.com/siyuan-note/siyuan/issues/11682
				continue
			}
			if atom.Span != c.DataAtom {
				allSpan = false
				break
			}
		}
		if allSpan {
			// 如果全部都是 span 子节点，那么直接使用 span 的内容 https://github.com/siyuan-note/siyuan/issues/11281
			code = []byte(util.DomText(n))
			code = bytes.ReplaceAll(code, []byte("\u00A0"), []byte(" "))
		}

		content := &ast.Node{Type: ast.NodeCodeSpanContent, Tokens: code}
		node.Type = ast.NodeCodeSpan
		if bytes.Contains(code, []byte("```")) {
			node.CodeMarkerLen = 2
		} else if bytes.Contains(code, []byte("``")) {
			node.CodeMarkerLen = 1
		} else if bytes.Contains(code, []byte("`")) {
			node.CodeMarkerLen = 2
		}
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanOpenMarker, Tokens: []byte("`")})
		if bytes.Contains(code, []byte("``")) {
			content.Tokens = append([]byte(" "), content.Tokens...)
		}
		node.AppendChild(content)
		if bytes.Contains(code, []byte("``")) {
			content.Tokens = append(content.Tokens, ' ')
		}
		node.AppendChild(&ast.Node{Type: ast.NodeCodeSpanCloseMarker, Tokens: []byte("`")})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
		return
	case atom.Br:
		if ast.NodeLink == tree.Context.Tip.Type || ast.NodeHeading == tree.Context.Tip.Type {
			break
		}

		if nil == n.NextSibling {
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
		text := strings.TrimSpace(util.DomText(n))
		if "" == text && nil != n.Parent && lute.parentIs(n, atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6, atom.Div, atom.Section) && nil == util.DomChildrenByType(n, atom.Img) {
			// 丢弃标题中文本为空的链接，这样的链接是没有锚文本的锚点
			// https://github.com/Vanessa219/vditor/issues/359
			// https://github.com/siyuan-note/siyuan/issues/11445
			return
		}
		if "" == text && nil == n.FirstChild {
			// 剪藏时过滤空的超链接 https://github.com/siyuan-note/siyuan/issues/5686
			return
		}

		if nil != n.FirstChild && atom.Img == n.FirstChild.DataAtom && strings.Contains(util.DomAttrValue(n.FirstChild, "src"), "wikimedia.org") {
			// Wikipedia 链接嵌套图片的情况只保留图片
			break
		}

		node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Img:
		imgClass := util.DomAttrValue(n, "class")
		imgAlt := util.DomAttrValue(n, "alt")
		if "emoji" == imgClass {
			if e := parse.EmojiUnicodeAlias[imgAlt]; "" != e {
				// 直接使用 alt 值（即 emoji 字符）https://github.com/siyuan-note/siyuan/issues/13342
				node.Type = ast.NodeText
				node.Tokens = []byte(imgAlt)
			} else {
				node.Type = ast.NodeEmoji
				emojiImg := &ast.Node{Type: ast.NodeEmojiImg, Tokens: tree.EmojiImgTokens(imgAlt, util.DomAttrValue(n, "src"))}
				emojiImg.AppendChild(&ast.Node{Type: ast.NodeEmojiAlias, Tokens: util.StrToBytes(":" + imgAlt + ":")})
				node.AppendChild(emojiImg)
			}
		} else {
			node.Type = ast.NodeImage
			node.AppendChild(&ast.Node{Type: ast.NodeBang})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
			if "" != imgAlt {
				imgAlt = strings.TrimSpace(imgAlt)
				imgAlt = strings.ReplaceAll(imgAlt, "动图封面", "动图")
				node.AppendChild(&ast.Node{Type: ast.NodeLinkText, Tokens: util.StrToBytes(imgAlt)})
			}
			node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
			node.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			src := util.DomAttrValue(n, "src")
			if strings.Contains(class, "ztext-gif") && strings.Contains(src, "zhimg.com") {
				// 处理知乎动图
				src = strings.Replace(src, ".jpg", ".webp", 1)
			}

			if strings.HasPrefix(src, "data:image") {
				// 处理可能存在的预加载情况
				if dataSrc := util.DomAttrValue(n, "data-src"); "" != dataSrc {
					src = dataSrc
				}
			}

			// 处理使用 data-original 属性的情况 https://github.com/siyuan-note/siyuan/issues/11826
			dataOriginal := util.DomAttrValue(n, "data-original")
			if "" != dataOriginal {
				if "" == src || !strings.HasSuffix(src, ".gif") {
					src = dataOriginal
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

			// Wikipedia 使用图片原图 https://github.com/siyuan-note/siyuan/issues/11640
			if strings.Contains(src, "wikipedia/commons/thumb/") {
				ext := path.Ext(src)
				if strings.Contains(src, ".svg.png") {
					ext = ".svg"
				}
				if idx := strings.Index(src, ext+"/"); 0 < idx {
					src = src[:idx+len(ext)]
					src = strings.ReplaceAll(src, "/commons/thumb/", "/commons/")
				}
			}

			node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: util.StrToBytes(src)})
			var linkTitle string

			var figcaption *html.Node
			for next := n.NextSibling; nil != next; next = next.NextSibling {
				if html.TextNode == next.Type && "" == strings.TrimSpace(next.Data) {
					continue
				}

				if atom.Figcaption == next.DataAtom {
					figcaption = next
					break
				}
			}

			if nil != figcaption {
				linkTitle = util.DomText(figcaption)
				figcaption.Unlink()
			}
			if "" == linkTitle {
				linkTitle = util.DomAttrValue(n, "title")
			}
			linkTitle = strings.TrimSpace(linkTitle)
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
		if !lute.ParseOptions.GFMStrikethrough {
			node.Type = ast.NodeHTMLTag
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagOpen, Tokens: util.StrToBytes("<s>")})
		} else {
			node.Type = ast.NodeStrikethrough
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2OpenMarker, Tokens: util.StrToBytes("~~")})
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.U:
		node.Type = ast.NodeHTMLTag
		node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagOpen, Tokens: util.StrToBytes("<u>")})
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Mark:
		if !lute.ParseOptions.Mark {
			node.Type = ast.NodeHTMLTag
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagOpen, Tokens: util.StrToBytes("<mark>")})
		} else {
			node.Type = ast.NodeMark
			node.AppendChild(&ast.Node{Type: ast.NodeMark2OpenMarker, Tokens: util.StrToBytes("==")})
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Sup:
		if !lute.ParseOptions.Sup {
			node.Type = ast.NodeHTMLTag
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagOpen, Tokens: util.StrToBytes("<sup>")})
		} else {
			node.Type = ast.NodeSup
			node.AppendChild(&ast.Node{Type: ast.NodeSupOpenMarker})
		}
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
	case atom.Sub:
		if !lute.ParseOptions.Sub {
			node.Type = ast.NodeHTMLTag
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagOpen, Tokens: util.StrToBytes("<sub>")})
		} else {
			node.Type = ast.NodeSub
			node.AppendChild(&ast.Node{Type: ast.NodeSubOpenMarker})
		}
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
		if nil == n.FirstChild {
			break
		}

		tbodys := util.DomChildrenByType(n.Parent, atom.Tbody)
		if 0 < len(tbodys) {
			tbody := tbodys[0]
			// 找到最多的 td 数
			var tdCount int
			for tr := tbody.FirstChild; nil != tr; tr = tr.NextSibling {
				if atom.Tr != tr.DataAtom {
					continue
				}

				var count int
				for td := tr.FirstChild; nil != td; td = td.NextSibling {
					if atom.Td == td.DataAtom {
						count++
					}
				}

				if count > tdCount {
					tdCount = count
				}
			}

			// 补全 thead 中 tr 的 th
			for tr := n.FirstChild; nil != tr; tr = tr.NextSibling {
				if atom.Tr != tr.DataAtom {
					continue
				}

				var count int
				for td := tr.FirstChild; nil != td; td = td.NextSibling {
					if atom.Th == td.DataAtom {
						count++
					}
				}

				if count < tdCount {
					for i := count; i < tdCount; i++ {
						th := &html.Node{Data: "th", DataAtom: atom.Th, Type: html.ElementNode}
						tr.AppendChild(th)
					}
				}
			}
		}

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

		if nil == tree.Context.Tip.ChildByType(ast.NodeTableHead) && 1 > len(util.DomChildrenByType(table, atom.Thead)) {
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
		class := util.DomAttrValue(n, "class")
		if "fn__space5" == class {
			return
		}

		if "tip" == class {
			// 转换为行级备注 https://github.com/siyuan-note/siyuan/issues/13998
			if nil != tree.Context.Tip.LastChild && ast.NodeText == tree.Context.Tip.LastChild.Type {
				tree.Context.Tip.LastChild.Type = ast.NodeTextMark
				tree.Context.Tip.LastChild.TextMarkType = "inline-memo"
				tree.Context.Tip.LastChild.TextMarkTextContent = tree.Context.Tip.LastChild.TokensStr()
				tree.Context.Tip.LastChild.TextMarkInlineMemoContent = util.DomText(n)
				if nil != tree.Context.Tip.LastChild.Previous && ast.NodeText == tree.Context.Tip.LastChild.Previous.Type {
					tree.Context.Tip.LastChild.Previous.Tokens = bytes.TrimSpace(tree.Context.Tip.LastChild.Previous.Tokens)
					if 0 == len(tree.Context.Tip.LastChild.Previous.Tokens) {
						tree.Context.Tip.LastChild.Previous.Unlink()
					}
				}
				return
			}
		}

		if title := strings.TrimSpace(util.DomAttrValue(n, "title")); "" != title {
			// 转换为行级备注 https://github.com/siyuan-note/siyuan/issues/13998
			node.Type = ast.NodeTextMark
			node.TextMarkType = "inline-memo"
			node.TextMarkTextContent = util.DomText(n)
			node.TextMarkInlineMemoContent = title
			tree.Context.Tip.AppendChild(node)
			if nil != tree.Context.Tip.LastChild.Previous && ast.NodeText == tree.Context.Tip.LastChild.Previous.Type {
				tree.Context.Tip.LastChild.Previous.Tokens = bytes.TrimSpace(tree.Context.Tip.LastChild.Previous.Tokens)
				if 0 == len(tree.Context.Tip.LastChild.Previous.Tokens) {
					tree.Context.Tip.LastChild.Previous.Unlink()
				}
			}
			return
		}

		// Improve inline elements pasting https://github.com/siyuan-note/siyuan/issues/11740
		dataType := util.DomAttrValue(n, "data-type")
		dataType = strings.Split(dataType, " ")[0] // 简化为只处理第一个类型
		switch dataType {
		case "inline-math":
			mathContent := util.DomAttrValue(n, "data-content")
			appendInlineMath(tree, mathContent)
			return
		case "code":
			if nil != tree.Context.Tip.LastChild && ast.NodeCodeSpan == tree.Context.Tip.LastChild.Type {
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(editor.Zwsp)})
			}

			code := &ast.Node{Type: ast.NodeCodeSpan}
			content := util.StrToBytes(util.DomText(n))
			if bytes.Contains(content, []byte("`")) {
				node.CodeMarkerLen = 2
			}
			code.AppendChild(&ast.Node{Type: ast.NodeCodeSpanOpenMarker, Tokens: []byte("`")})
			code.AppendChild(&ast.Node{Type: ast.NodeCodeSpanContent, Tokens: content})
			code.AppendChild(&ast.Node{Type: ast.NodeCodeSpanCloseMarker, Tokens: []byte("`")})
			tree.Context.Tip.AppendChild(code)
			return
		case "tag":
			tag := &ast.Node{Type: ast.NodeTag}
			tag.AppendChild(&ast.Node{Type: ast.NodeTagOpenMarker, Tokens: util.StrToBytes("#")})
			tag.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
			tag.AppendChild(&ast.Node{Type: ast.NodeTagCloseMarker, Tokens: util.StrToBytes("#")})
			tree.Context.Tip.AppendChild(tag)
			return
		case "kbd":
			if nil != tree.Context.Tip.LastChild && ast.NodeKbd == tree.Context.Tip.LastChild.Type {
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(editor.Zwsp)})
			}

			kbd := &ast.Node{Type: ast.NodeKbd}
			kbd.AppendChild(&ast.Node{Type: ast.NodeKbdOpenMarker})
			kbd.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
			kbd.AppendChild(&ast.Node{Type: ast.NodeKbdCloseMarker})
			tree.Context.Tip.AppendChild(kbd)
			return
		case "sub":
			sub := &ast.Node{Type: ast.NodeSub}
			sub.AppendChild(&ast.Node{Type: ast.NodeSubOpenMarker})
			sub.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
			sub.AppendChild(&ast.Node{Type: ast.NodeSubCloseMarker})
			tree.Context.Tip.AppendChild(sub)
			return
		case "sup":
			sup := &ast.Node{Type: ast.NodeSup}
			sup.AppendChild(&ast.Node{Type: ast.NodeSupOpenMarker})
			sup.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
			sup.AppendChild(&ast.Node{Type: ast.NodeSupCloseMarker})
			tree.Context.Tip.AppendChild(sup)
			return
		case "mark":
			mark := &ast.Node{Type: ast.NodeMark}
			mark.AppendChild(&ast.Node{Type: ast.NodeMark2OpenMarker, Tokens: util.StrToBytes("==")})
			mark.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
			mark.AppendChild(&ast.Node{Type: ast.NodeMark2CloseMarker, Tokens: util.StrToBytes("==")})
			tree.Context.Tip.AppendChild(mark)
			return
		case "s":
			s := &ast.Node{Type: ast.NodeStrikethrough}
			s.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2OpenMarker, Tokens: util.StrToBytes("~~")})
			s.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
			s.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2CloseMarker, Tokens: util.StrToBytes("~~")})
			tree.Context.Tip.AppendChild(s)
			return
		case "u":
			u := &ast.Node{Type: ast.NodeUnderline}
			u.AppendChild(&ast.Node{Type: ast.NodeUnderlineOpenMarker})
			u.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
			u.AppendChild(&ast.Node{Type: ast.NodeUnderlineCloseMarker})
			tree.Context.Tip.AppendChild(u)
			return
		case "em":
			em := &ast.Node{Type: ast.NodeEmphasis}
			em.AppendChild(&ast.Node{Type: ast.NodeEmA6kOpenMarker, Tokens: util.StrToBytes("*")})
			em.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
			em.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker, Tokens: util.StrToBytes("*")})
			tree.Context.Tip.AppendChild(em)
			return
		case "strong":
			strong := &ast.Node{Type: ast.NodeStrong}
			strong.AppendChild(&ast.Node{Type: ast.NodeStrongA6kOpenMarker, Tokens: util.StrToBytes("**")})
			strong.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
			strong.AppendChild(&ast.Node{Type: ast.NodeStrongA6kCloseMarker, Tokens: util.StrToBytes("**")})
			tree.Context.Tip.AppendChild(strong)
			return
		case "block-ref":
			ref := &ast.Node{Type: ast.NodeBlockRef}
			ref.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			ref.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
			ref.AppendChild(&ast.Node{Type: ast.NodeBlockRefID, Tokens: util.StrToBytes(util.DomAttrValue(n, "data-id"))})
			ref.AppendChild(&ast.Node{Type: ast.NodeBlockRefSpace})
			if "s" == util.DomAttrValue(n, "data-subtype") {
				ref.AppendChild(&ast.Node{Type: ast.NodeBlockRefText, Tokens: util.StrToBytes(util.DomText(n))})
			} else {
				ref.AppendChild(&ast.Node{Type: ast.NodeBlockRefDynamicText, Tokens: util.StrToBytes(util.DomText(n))})
			}
			ref.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
			ref.AppendChild(&ast.Node{Type: ast.NodeCloseParen})

			tree.Context.Tip.AppendChild(ref)
			return
		}

		// The browser extension supports Zhihu formula https://github.com/siyuan-note/siyuan/issues/5599
		if tex := strings.TrimSpace(util.DomAttrValue(n, "data-tex")); "" != tex {
			if nil != n.Parent && strings.Contains(util.DomAttrValue(n.Parent, "class"), "math-inline") {
				appendInlineMath(tree, tex)
				return
			}

			parentInline := nil != n.Parent && (lute.parentIs(n, atom.Strong, atom.Em, atom.I, atom.B, atom.Span, atom.P, atom.Td, atom.Th) || strings.Contains(util.DomAttrValue(n.Parent, "class"), "inline"))
			if parentInline && atom.Span == n.DataAtom && nil == n.PrevSibling && (nil == n.NextSibling || (html.TextNode == n.NextSibling.Type && "" == strings.TrimSpace(util.DomText(n.NextSibling)))) &&
				nil == n.Parent.PrevSibling && (nil == n.Parent.NextSibling || (html.TextNode == n.Parent.NextSibling.Type && "" == strings.TrimSpace(util.DomText(n.Parent.NextSibling)))) {
				appendMathBlock(tree, tex)
				return
			}

			if !parentInline && nil == n.PrevSibling && (nil == n.NextSibling || (html.TextNode == n.NextSibling.Type && "" == strings.TrimSpace(util.DomText(n.NextSibling)))) {
				appendMathBlock(tree, tex)
				return
			}

			if atom.Span == n.DataAtom && "katex-display" == util.DomAttrValue(n, "class") ||
				nil != util.DomChildByTypeAndClass(n, atom.Span, "MathJax_SVG_Display") {
				appendMathBlock(tree, tex)
				return
			}

			if strings.HasSuffix(strings.TrimSpace(tex), "\\\\") || strings.Contains(tex, "\n") || strings.Contains(tex, "\\tag{") {
				appendMathBlock(tree, tex)
			} else {
				appendInlineMath(tree, tex)
			}
			return
		}

		// The browser extension supports CSDN formula https://github.com/siyuan-note/siyuan/issues/5624
		if strings.Contains(strings.ToLower(strings.TrimSpace(util.DomAttrValue(n, "class"))), "katex") {
			if span := util.DomChildByTypeAndClass(n, atom.Span, "katex-mathml"); nil != span {
				if tex := util.DomText(span.FirstChild); "" != tex {
					tex = strings.TrimSpace(tex)
					for strings.Contains(tex, "\n ") {
						tex = strings.ReplaceAll(tex, "\n ", "\n")
					}
					// 根据最后 4 个换行符分隔公式内容
					if idx := strings.LastIndex(tex, "\n\n\n\n"); 0 < idx {
						tex = tex[idx+4:]
						tex = strings.TrimSpace(tex)
						appendInlineMath(tree, tex)
						return
					}
				}
			}
		}
		if strings.Contains(strings.ToLower(strings.TrimSpace(util.DomAttrValue(n, "class"))), "mathjax") {
			scripts := util.DomChildrenByType(n, atom.Script)
			if 0 < len(scripts) {
				script := scripts[0]
				if tex := util.DomText(script.FirstChild); "" != tex {
					appendInlineMath(tree, tex)
					return
				}
			}
			return
		}
	case atom.Kbd:
		if nil != tree.Context.Tip.LastChild && ast.NodeKbd == tree.Context.Tip.LastChild.Type {
			tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(editor.Zwsp)})
		}

		kbd := &ast.Node{Type: ast.NodeKbd}
		kbd.AppendChild(&ast.Node{Type: ast.NodeKbdOpenMarker})
		kbd.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(util.DomText(n))})
		kbd.AppendChild(&ast.Node{Type: ast.NodeKbdCloseMarker})
		tree.Context.Tip.AppendChild(kbd)
		return
	case atom.Font:
		node.Type = ast.NodeText
		tokens := []byte(util.DomText(n))
		for strings.Contains(string(tokens), "\n\n") {
			tokens = bytes.ReplaceAll(tokens, []byte("\n\n"), []byte("\n"))
		}

		for strings.Contains(string(tokens), "\n  ") {
			tokens = bytes.ReplaceAll(tokens, []byte("\n  "), []byte("\n "))
		}
		tokens = bytes.ReplaceAll(tokens, []byte("\n "), []byte("\n"))

		tokens = bytes.ReplaceAll(tokens, []byte("\n"), []byte(" "))
		node.Tokens = tokens
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Details:
		node.Type = ast.NodeList
		node.ListData = &ast.ListData{}
		node.ListData.Tight = true
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
		defer tree.Context.ParentTip()
		defer tree.Context.ParentTip()
	case atom.Summary:
		if ast.NodeList != tree.Context.Tip.Type || nil == n.Parent || atom.Details != n.Parent.DataAtom {
			return
		}

		li := &ast.Node{Type: ast.NodeListItem}
		li.ListData = &ast.ListData{Marker: []byte("*"), BulletChar: '*'}
		node.Type = ast.NodeParagraph
		li.AppendChild(node)
		tree.Context.Tip.AppendChild(li)
		tree.Context.Tip = node
	case atom.Iframe, atom.Audio, atom.Video:
		node.Type = ast.NodeHTMLBlock
		node.Tokens = util.DomHTML(n)
		tree.Context.Tip.AppendChild(node)
		return
	case atom.Noscript:
		return
	case atom.Script:
		if tex := util.DomText(n.FirstChild); "" != tex {
			if tree.Context.Tip.IsContainerBlock() ||
				(nil != n.Parent && strings.Contains(util.DomAttrValue(n.Parent, "class"), "math display") && n.Parent.LastChild == n) ||
				strings.Contains(util.DomAttrValue(n, "type"), "mode=display") {
				appendMathBlock(tree, tex)
			} else {
				appendInlineMath(tree, tex)
			}
			return
		}
	case atom.Figcaption:
		if tree.Context.Tip.IsBlock() {
			if ast.NodeDocument != tree.Context.Tip.Type {
				tree.Context.Tip.AppendChild(&ast.Node{Type: ast.NodeHardBreak})
				break
			}
		}

		node.Type = ast.NodeParagraph
		tree.Context.Tip.AppendChild(node)
		tree.Context.Tip = node
	case atom.Figure, atom.Picture:
		if !tree.Context.Tip.IsBlock() {
			break
		}

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
		if !lute.ParseOptions.InlineAsterisk || !lute.ParseOptions.InlineUnderscore {
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagClose, Tokens: util.StrToBytes("</em>")})
		} else {
			node.Type = ast.NodeEmphasis
			node.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker, Tokens: util.StrToBytes("*")})
		}
		appendSpace(n, tree, lute)
	case atom.Strong, atom.B:
		if !lute.ParseOptions.InlineAsterisk || !lute.ParseOptions.InlineUnderscore {
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagClose, Tokens: util.StrToBytes("</strong>")})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrongA6kCloseMarker, Tokens: util.StrToBytes("**")})
		}
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
		if !lute.ParseOptions.GFMStrikethrough {
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagClose, Tokens: util.StrToBytes("</s>")})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2CloseMarker, Tokens: util.StrToBytes("~~")})
		}
		appendSpace(n, tree, lute)
	case atom.U:
		node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagClose, Tokens: util.StrToBytes("</u>")})
	case atom.Mark:
		if !lute.ParseOptions.Mark {
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagClose, Tokens: util.StrToBytes("</mark>")})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeMark2CloseMarker, Tokens: util.StrToBytes("==")})
		}
		appendSpace(n, tree, lute)
	case atom.Sup:
		if !lute.ParseOptions.Sup {
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagClose, Tokens: util.StrToBytes("</sup>")})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeSupCloseMarker})
		}
		appendSpace(n, tree, lute)
	case atom.Sub:
		if !lute.ParseOptions.Sub {
			node.AppendChild(&ast.Node{Type: ast.NodeHTMLTagClose, Tokens: util.StrToBytes("</sub>")})
		} else {
			node.AppendChild(&ast.Node{Type: ast.NodeSubCloseMarker})
		}
		appendSpace(n, tree, lute)
	case atom.Details:
		tree.Context.ParentTip()
	case atom.Summary:
		tree.Context.ParentTip()
	}
}

func appendInlineMath(tree *parse.Tree, tex string) {
	tex = strings.TrimSpace(tex)
	if "" == tex {
		return
	}

	inlineMath := &ast.Node{Type: ast.NodeInlineMath}
	inlineMath.AppendChild(&ast.Node{Type: ast.NodeInlineMathOpenMarker, Tokens: []byte("$")})
	inlineMath.AppendChild(&ast.Node{Type: ast.NodeInlineMathContent, Tokens: util.StrToBytes(tex)})
	inlineMath.AppendChild(&ast.Node{Type: ast.NodeInlineMathCloseMarker, Tokens: []byte("$")})
	tree.Context.Tip.AppendChild(inlineMath)
	tree.Context.Tip = inlineMath
	defer tree.Context.ParentTip()
}

func appendMathBlock(tree *parse.Tree, tex string) {
	tex = strings.TrimSpace(tex)
	if "" == tex {
		return
	}

	mathBlock := &ast.Node{Type: ast.NodeMathBlock}
	mathBlock.AppendChild(&ast.Node{Type: ast.NodeMathBlockOpenMarker, Tokens: []byte("$$")})
	mathBlock.AppendChild(&ast.Node{Type: ast.NodeMathBlockContent, Tokens: util.StrToBytes(tex)})
	mathBlock.AppendChild(&ast.Node{Type: ast.NodeMathBlockCloseMarker, Tokens: []byte("$$")})

	if ast.NodeParagraph == tree.Context.Tip.Type {
		tree.Context.Tip.InsertAfter(mathBlock)
		if nil == tree.Context.Tip.FirstChild {
			tree.Context.Tip.Unlink()
		}
	} else {
		tree.Context.Tip.AppendChild(mathBlock)
	}
	tree.Context.Tip = mathBlock
	defer tree.Context.ParentTip()
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
