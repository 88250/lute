// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"bytes"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
)

// NestedInlines2FlattedSpansHybrid 将嵌套的行级节点转换为平铺的文本标记节点。
// 该函数不会移除转义节点。
func NestedInlines2FlattedSpansHybrid(tree *Tree, isExportMd bool) {
	var unlinks []*ast.Node
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		if ast.NodeLink == n.Type {
			var unlinkBackslashes []*ast.Node
			ast.Walk(n, func(c *ast.Node, entering bool) ast.WalkStatus {
				if !entering {
					return ast.WalkContinue
				}

				if ast.NodeBackslash == c.Type {
					cont := c.ChildByType(ast.NodeBackslashContent)
					if nil != cont {
						linkText := &ast.Node{Type: ast.NodeLinkText, Tokens: cont.Tokens}
						c.InsertBefore(linkText)
					}

					unlinkBackslashes = append(unlinkBackslashes, c)
				}
				return ast.WalkContinue
			})
			for _, backslash := range unlinkBackslashes {
				backslash.Unlink()
			}

			// 超链接嵌套图片情况下，图片子节点移到超链接节点前面
			img := n.ChildByType(ast.NodeImage)
			if nil == img {
				return ast.WalkContinue
			}
			n.InsertBefore(img)
			linkText := n.ChildByType(ast.NodeLinkText)
			if nil == linkText || 1 > len(bytes.TrimSpace(linkText.Tokens)) {
				if openBracket := n.ChildByType(ast.NodeOpenBracket); nil != openBracket {
					if dest := n.ChildByType(ast.NodeLinkDest); nil != dest {
						openBracket.InsertAfter(&ast.Node{Type: ast.NodeLinkText, Tokens: dest.Tokens})
					}
				}
			}
		}
		return ast.WalkContinue
	})
	for _, n := range unlinks {
		n.Unlink()
	}
	unlinks = nil

	var tags []string
	var span *ast.Node
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		switch n.Type {
		// case ast.NodeBackslash: Spin 过程中存在转义节点，转义节点本身就是嵌套的，所以需要在这里排除处理
		case ast.NodeCodeSpan:
			processNestedNode(n, "code", &tags, &unlinks, entering)
		case ast.NodeTag:
			processNestedNode(n, "tag", &tags, &unlinks, entering)
		case ast.NodeInlineMath:
			processNestedNode(n, "inline-math", &tags, &unlinks, entering)
		case ast.NodeEmphasis:
			processNestedNode(n, "em", &tags, &unlinks, entering)
		case ast.NodeStrong:
			processNestedNode(n, "strong", &tags, &unlinks, entering)
		case ast.NodeStrikethrough:
			processNestedNode(n, "s", &tags, &unlinks, entering)
		case ast.NodeMark:
			processNestedNode(n, "mark", &tags, &unlinks, entering)
		case ast.NodeUnderline:
			processNestedNode(n, "u", &tags, &unlinks, entering)
		case ast.NodeSub:
			processNestedNode(n, "sub", &tags, &unlinks, entering)
		case ast.NodeSup:
			processNestedNode(n, "sup", &tags, &unlinks, entering)
		case ast.NodeKbd:
			processNestedNode(n, "kbd", &tags, &unlinks, entering)
		case ast.NodeLink:
			processNestedNode(n, "a", &tags, &unlinks, entering)
		case ast.NodeBlockRef:
			processNestedNode(n, "block-ref", &tags, &unlinks, entering)
		case ast.NodeText, ast.NodeCodeSpanContent, ast.NodeInlineMathContent, ast.NodeLinkText, ast.NodeBlockRefID, ast.NodeHTMLEntity, ast.NodeBackslash:
			if 1 > len(tags) {
				return ast.WalkContinue
			}

			if entering {
				span = &ast.Node{Type: ast.NodeTextMark, TextMarkType: strings.Join(tags, " "), TextMarkTextContent: string(html.EscapeHTML(n.Tokens))}
				if ast.NodeInlineMathContent == n.Type {
					span.TextMarkTextContent = ""
					span.TextMarkInlineMathContent = string(html.EscapeHTML(n.Tokens))
					if n.ParentIs(ast.NodeTableCell) && !isExportMd {
						// Improve the handling of inline-math containing `|` in the table https://github.com/siyuan-note/siyuan/issues/9227
						span.TextMarkInlineMathContent = strings.ReplaceAll(span.TextMarkInlineMathContent, "\\|", "|")
					}
				} else if ast.NodeBackslash == n.Type {
					if c := n.ChildByType(ast.NodeBackslashContent); nil != c {
						span.TextMarkTextContent = string(html.EscapeHTML(c.Tokens))
					}
				} else if ast.NodeBlockRefID == n.Type {
					span.TextMarkBlockRefSubtype = "s"
					span.TextMarkTextContent = n.TokensStr()

					refText := n.Parent.ChildByType(ast.NodeBlockRefText)
					if nil == refText {
						refText = n.Parent.ChildByType(ast.NodeBlockRefDynamicText)
						span.TextMarkBlockRefSubtype = "d"
					}
					if nil != refText {
						span.TextMarkTextContent = refText.TokensStr()
					}

					span.TextMarkBlockRefID = n.Parent.ChildByType(ast.NodeBlockRefID).TokensStr()
				} else if n.ParentIs(ast.NodeLink) && !n.ParentIs(ast.NodeImage) {
					if next := n.Next; nil != next && ast.NodeLinkText == next.Type {
						// 合并相邻的链接文本节点
						n.Next.PrependTokens(n.Tokens)
						return ast.WalkContinue
					}

					var link *ast.Node
					for p := n.Parent; nil != p; p = p.Parent {
						if ast.NodeLink == p.Type {
							link = p
							break
						}
					}
					if nil != link {
						dest := link.ChildByType(ast.NodeLinkDest)
						if nil != dest {
							span.TextMarkAHref = string(dest.Tokens)
						}
						title := link.ChildByType(ast.NodeLinkTitle)
						if nil != title {
							span.TextMarkATitle = string(title.Tokens)
						}
					}
				}
			} else {
				if next := n.Next; nil != next && ast.NodeLinkText == next.Type {
					return ast.WalkContinue
				}

				span.KramdownIAL = n.Parent.KramdownIAL
				if n.IsMarker() {
					n.Parent.InsertBefore(span)
				} else {
					n.InsertBefore(span)
					if ast.NodeText == n.Type {
						unlinks = append(unlinks, n)
					}
				}
			}
		case ast.NodeTextMark:
			if 1 > len(tags) {
				return ast.WalkContinue
			}

			if entering {
				contain := false
				for _, tag := range tags {
					if n.IsTextMarkType(tag) {
						contain = true
						break
					}
				}
				if !contain {
					tags = append(tags, n.TextMarkType)
					n.TextMarkType = strings.Join(tags, " ")
				}
			} else {
				if nil == n.Next || n.Next.IsCloseMarker() {
					tags = tags[:len(tags)-1]
				}
			}
			return ast.WalkContinue
		}
		return ast.WalkContinue
	})

	for _, n := range unlinks {
		n.Unlink()
	}
}

// NestedInlines2FlattedSpans 将嵌套的行级节点转换为平铺的文本标记节点。
// 该函数会移除转义节点。
func NestedInlines2FlattedSpans(tree *Tree, isExportMd bool) {
	var unlinks []*ast.Node
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		if ast.NodeLink == n.Type {
			// 超链接嵌套图片情况下，图片子节点移到超链接节点前面
			img := n.ChildByType(ast.NodeImage)
			if nil == img {
				return ast.WalkContinue
			}
			n.InsertBefore(img)
			linkText := n.ChildByType(ast.NodeLinkText)
			if nil == linkText || 1 > len(bytes.TrimSpace(linkText.Tokens)) {
				if openBracket := n.ChildByType(ast.NodeOpenBracket); nil != openBracket {
					if dest := n.ChildByType(ast.NodeLinkDest); nil != dest {
						openBracket.InsertAfter(&ast.Node{Type: ast.NodeLinkText, Tokens: dest.Tokens})
					}
				}
			}
		} else if ast.NodeBackslash == n.Type {
			if n.ParentIs(ast.NodeTableCell) {
				return ast.WalkContinue
			}

			// 不再需要反斜杠转义节点
			if c := n.FirstChild; nil != c {
				tokens := html.UnescapeHTML(c.Tokens)
				processed := false
				if previous := n.Previous; nil != previous && (ast.NodeText == previous.Type || ast.NodeLinkText == previous.Type) {
					previous.AppendTokens(tokens)
					processed = true
				} else if next := n.Next; nil != next && (ast.NodeText == next.Type || ast.NodeLinkText == next.Type) {
					next.PrependTokens(tokens)
					processed = true
				}
				if !processed {
					text := &ast.Node{Type: ast.NodeText, Tokens: tokens}
					n.InsertBefore(text)
				}
				unlinks = append(unlinks, n)
			}
		}
		return ast.WalkContinue
	})
	for _, n := range unlinks {
		n.Unlink()
	}
	unlinks = nil

	var tags []string
	var span *ast.Node
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		switch n.Type {
		case ast.NodeBackslash:
			if !n.ParentIs(ast.NodeTableCell) {
				unlinks = append(unlinks, n)
			}
		case ast.NodeCodeSpan:
			processNestedNode(n, "code", &tags, &unlinks, entering)
		case ast.NodeTag:
			processNestedNode(n, "tag", &tags, &unlinks, entering)
		case ast.NodeInlineMath:
			processNestedNode(n, "inline-math", &tags, &unlinks, entering)
		case ast.NodeEmphasis:
			processNestedNode(n, "em", &tags, &unlinks, entering)
		case ast.NodeStrong:
			processNestedNode(n, "strong", &tags, &unlinks, entering)
		case ast.NodeStrikethrough:
			processNestedNode(n, "s", &tags, &unlinks, entering)
		case ast.NodeMark:
			processNestedNode(n, "mark", &tags, &unlinks, entering)
		case ast.NodeUnderline:
			processNestedNode(n, "u", &tags, &unlinks, entering)
		case ast.NodeSub:
			processNestedNode(n, "sub", &tags, &unlinks, entering)
		case ast.NodeSup:
			processNestedNode(n, "sup", &tags, &unlinks, entering)
		case ast.NodeKbd:
			processNestedNode(n, "kbd", &tags, &unlinks, entering)
		case ast.NodeLink:
			processNestedNode(n, "a", &tags, &unlinks, entering)
		case ast.NodeBlockRef:
			processNestedNode(n, "block-ref", &tags, &unlinks, entering)
		case ast.NodeText, ast.NodeCodeSpanContent, ast.NodeInlineMathContent, ast.NodeLinkText, ast.NodeBlockRefID, ast.NodeHTMLEntity:
			if 1 > len(tags) {
				return ast.WalkContinue
			}

			if entering {
				span = &ast.Node{Type: ast.NodeTextMark, TextMarkType: strings.Join(tags, " "), TextMarkTextContent: string(html.EscapeHTML(n.Tokens))}
				if ast.NodeInlineMathContent == n.Type {
					span.TextMarkTextContent = ""
					span.TextMarkInlineMathContent = string(html.EscapeHTML(n.Tokens))
					if n.ParentIs(ast.NodeTableCell) && !isExportMd {
						// Improve the handling of inline-math containing `|` in the table https://github.com/siyuan-note/siyuan/issues/9227
						span.TextMarkInlineMathContent = strings.ReplaceAll(span.TextMarkInlineMathContent, "\\|", "|")
					}
				} else if ast.NodeBackslash == n.Type {
					if c := n.ChildByType(ast.NodeBackslashContent); nil != c {
						span.TextMarkTextContent = string(html.EscapeHTML(c.Tokens))
					}
				} else if ast.NodeBlockRefID == n.Type {
					span.TextMarkBlockRefSubtype = "s"
					span.TextMarkTextContent = n.TokensStr()

					refText := n.Parent.ChildByType(ast.NodeBlockRefText)
					if nil == refText {
						refText = n.Parent.ChildByType(ast.NodeBlockRefDynamicText)
						span.TextMarkBlockRefSubtype = "d"
					}
					if nil != refText {
						span.TextMarkTextContent = refText.TokensStr()
					}

					span.TextMarkBlockRefID = n.Parent.ChildByType(ast.NodeBlockRefID).TokensStr()
				} else if n.ParentIs(ast.NodeLink) && !n.ParentIs(ast.NodeImage) {
					if next := n.Next; nil != next && ast.NodeLinkText == next.Type {
						// 合并相邻的链接文本节点
						n.Next.PrependTokens(n.Tokens)
						return ast.WalkContinue
					}

					var link *ast.Node
					for p := n.Parent; nil != p; p = p.Parent {
						if ast.NodeLink == p.Type {
							link = p
							break
						}
					}
					if nil != link {
						dest := link.ChildByType(ast.NodeLinkDest)
						if nil != dest {
							span.TextMarkAHref = string(dest.Tokens)
						}
						title := link.ChildByType(ast.NodeLinkTitle)
						if nil != title {
							span.TextMarkATitle = string(title.Tokens)
						}
					}
				}
			} else {
				if next := n.Next; nil != next && ast.NodeLinkText == next.Type {
					return ast.WalkContinue
				}

				span.KramdownIAL = n.Parent.KramdownIAL
				if n.IsMarker() {
					n.Parent.InsertBefore(span)
				} else {
					n.InsertBefore(span)
					if ast.NodeText == n.Type {
						unlinks = append(unlinks, n)
					}
				}
			}
		case ast.NodeTextMark:
			if 1 > len(tags) {
				return ast.WalkContinue
			}

			if entering {
				contain := false
				for _, tag := range tags {
					if n.IsTextMarkType(tag) {
						contain = true
						break
					}
				}
				if !contain {
					tags = append(tags, n.TextMarkType)
					n.TextMarkType = strings.Join(tags, " ")
				}
			} else {
				if nil == n.Next || n.IsCloseMarker() {
					tags = tags[:len(tags)-1]
				}
			}
			return ast.WalkContinue
		}
		return ast.WalkContinue
	})

	for _, n := range unlinks {
		n.Unlink()
	}
}

func processNestedNode(n *ast.Node, tag string, tags *[]string, unlinks *[]*ast.Node, entering bool) {
	if entering {
		*tags = append(*tags, tag)
	} else {
		if 0 < len(*tags) {
			*tags = (*tags)[:len(*tags)-1]
		}

		*unlinks = append(*unlinks, n)
		for c := n.FirstChild; nil != c; {
			next := c.Next
			if ast.NodeTextMark == c.Type || ast.NodeText == c.Type {
				n.InsertBefore(c)
			} else if ast.NodeLinkDest == c.Type {
				if nil != n.Previous && ast.NodeTextMark == n.Previous.Type {
					n.Previous.TextMarkAHref = string(c.Tokens)
				}
			}
			c = next
		}
	}
}

func TextMarks2Inlines(tree *Tree) {
	inlines := func(content string) (ret []*ast.Node) {
		subTree := Inline("", []byte(content), tree.Context.ParseOption)
		for c := subTree.Root.FirstChild.FirstChild; nil != c; c = c.Next {
			ret = append(ret, c)
		}
		return
	}

	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		if ast.NodeTextMark == n.Type {
			switch n.TextMarkType {
			case "sup":
				n.Type = ast.NodeSup
				n.PrependChild(&ast.Node{Type: ast.NodeSupOpenMarker})
				nodes := inlines(n.TextMarkTextContent)
				for _, node := range nodes {
					n.AppendChild(node)
				}
				n.AppendChild(&ast.Node{Type: ast.NodeSupCloseMarker})
			case "sub":
				n.Type = ast.NodeSub
				n.PrependChild(&ast.Node{Type: ast.NodeSubOpenMarker})
				nodes := inlines(n.TextMarkTextContent)
				for _, node := range nodes {
					n.AppendChild(node)
				}
				n.AppendChild(&ast.Node{Type: ast.NodeSubCloseMarker})
			case "em":
				n.Type = ast.NodeEmphasis
				n.PrependChild(&ast.Node{Type: ast.NodeEmA6kOpenMarker})
				nodes := inlines(n.TextMarkTextContent)
				for _, node := range nodes {
					n.AppendChild(node)
				}
				n.AppendChild(&ast.Node{Type: ast.NodeEmA6kCloseMarker})
			case "strong":
				n.Type = ast.NodeStrong
				n.PrependChild(&ast.Node{Type: ast.NodeStrongA6kOpenMarker})
				nodes := inlines(n.TextMarkTextContent)
				for _, node := range nodes {
					n.AppendChild(node)
				}
				n.AppendChild(&ast.Node{Type: ast.NodeStrongA6kCloseMarker})
			case "mark":
				n.Type = ast.NodeMark
				n.PrependChild(&ast.Node{Type: ast.NodeMark2OpenMarker})
				nodes := inlines(n.TextMarkTextContent)
				for _, node := range nodes {
					n.AppendChild(node)
				}
				n.AppendChild(&ast.Node{Type: ast.NodeMark2CloseMarker})
			case "s":
				n.Type = ast.NodeStrikethrough
				n.PrependChild(&ast.Node{Type: ast.NodeStrikethrough2OpenMarker})
				nodes := inlines(n.TextMarkTextContent)
				for _, node := range nodes {
					n.AppendChild(node)
				}
				n.AppendChild(&ast.Node{Type: ast.NodeStrikethrough2CloseMarker})
			}
		}
		return ast.WalkContinue
	})
}
