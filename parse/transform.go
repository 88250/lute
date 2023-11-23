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
			if nil == n.ChildByType(ast.NodeLinkText) {
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
			if nil == n.ChildByType(ast.NodeLinkText) {
				if openBracket := n.ChildByType(ast.NodeOpenBracket); nil != openBracket {
					if dest := n.ChildByType(ast.NodeLinkDest); nil != dest {
						openBracket.InsertAfter(&ast.Node{Type: ast.NodeLinkText, Tokens: dest.Tokens})
					}
				}
			}
		} else if ast.NodeBackslash == n.Type {
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
			unlinks = append(unlinks, n)
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
				}
			}
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
		*tags = (*tags)[:len(*tags)-1]
		*unlinks = append(*unlinks, n)
		for c := n.FirstChild; nil != c; {
			next := c.Next
			if ast.NodeTextMark == c.Type {
				n.InsertBefore(c)
			}
			c = next
		}
	}
}
