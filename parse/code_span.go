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

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
)

func (t *Tree) parseCodeSpan(block *ast.Node, ctx *InlineContext) (ret *ast.Node) {
	startPos := ctx.pos
	n := 0
	for ; startPos+n < ctx.tokensLen; n++ {
		if lex.ItemBacktick != ctx.tokens[startPos+n] {
			break
		}
	}

	backticks := ctx.tokens[startPos : startPos+n]
	if ctx.tokensLen <= startPos+n {
		ctx.pos += n
		ret = &ast.Node{Type: ast.NodeText, Tokens: backticks}
		return
	}
	openMarker := &ast.Node{Type: ast.NodeCodeSpanOpenMarker, Tokens: backticks}

	endPos := t.matchCodeSpanEnd(ctx.tokens[startPos+n:], n)
	if 1 > endPos {
		ctx.pos += n
		ret = &ast.Node{Type: ast.NodeText, Tokens: backticks}
		return
	}
	endPos = startPos + endPos + n
	closeMarker := &ast.Node{Type: ast.NodeCodeSpanCloseMarker, Tokens: ctx.tokens[endPos : endPos+n]}

	textTokens := ctx.tokens[startPos+n : endPos]
	textTokens = lex.ReplaceAll(textTokens, lex.ItemNewline, lex.ItemSpace)
	if 2 < len(textTokens) && lex.ItemSpace == textTokens[0] && lex.ItemSpace == textTokens[len(textTokens)-1] && !lex.IsBlankLine(textTokens) {
		// 如果首尾是空格并且整行不是空行时剔除首尾的一个空格
		openMarker.Tokens = append(openMarker.Tokens, textTokens[0])
		closeMarker.Tokens = ctx.tokens[endPos-1 : endPos+n]
		textTokens = textTokens[1 : len(textTokens)-1]
	}

	if t.Context.ParseOption.GFMTable {
		if ast.NodeTableCell == block.Type {
			// 表格中的代码中带有管道符的处理 https://github.com/88250/lute/issues/63
			textTokens = bytes.ReplaceAll(textTokens, []byte("\\|"), []byte("|"))
		}
	}

	ret = &ast.Node{Type: ast.NodeCodeSpan, CodeMarkerLen: n}
	ret.AppendChild(openMarker)

	if t.Context.ParseOption.ProtyleWYSIWYG {
		// Improve `inline code` markdown editing https://github.com/siyuan-note/siyuan/issues/9978

		if !bytes.HasPrefix(textTokens, []byte("<span data-type=\"code\">")) {
			// HTML 转换 Markdown 时需要转义 HTML 实体
			textTokens = bytes.ReplaceAll(textTokens, []byte("&"), []byte("&amp;"))
		}

		inlineTree := Inline("", textTokens, t.Context.ParseOption)
		if nil != inlineTree {
			content := bytes.Buffer{}
			ast.Walk(inlineTree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
				if !entering {
					content.WriteString(n.Marker(entering))
					if ast.NodeLinkTitle == n.Type {
						content.WriteByte(lex.ItemDoublequote)
					} else if ast.NodeTextMark == n.Type {
						if "kbd" == n.TextMarkType {
							content.WriteString("</kbd>")
						} else if "u" == n.TextMarkType {
							content.WriteString("</u>")
						} else if "sup" == n.TextMarkType {
							content.WriteString("</sup>")
						} else if "sub" == n.TextMarkType {
							content.WriteString("</sub>")
						}
					}
					return ast.WalkContinue
				}

				content.WriteString(n.Marker(entering))

				if ast.NodeTextMark == n.Type {
					if "kbd" == n.TextMarkType {
						content.WriteString("<kbd>")
					} else if "u" == n.TextMarkType {
						content.WriteString("<u>")
					} else if "sup" == n.TextMarkType {
						content.WriteString("<sup>")
					} else if "sub" == n.TextMarkType {
						content.WriteString("<sub>")
					}

					content.WriteString(n.TextMarkTextContent)
				} else if ast.NodeText == n.Type || ast.NodeLinkText == n.Type || ast.NodeLinkTitle == n.Type || ast.NodeLinkDest == n.Type {
					if entering {
						if ast.NodeLinkTitle == n.Type {
							content.WriteByte(lex.ItemDoublequote)
						}
					}

					if spanIdx1 := bytes.Index(n.Tokens, []byte("<span data-type=")); 0 <= spanIdx1 {
						if spanIdx2 := bytes.Index(n.Tokens[spanIdx1:], []byte(">")); 0 <= spanIdx2 {
							content.Write(n.Tokens[:spanIdx1])
							content.Write(n.Tokens[spanIdx1+spanIdx2+1:])
							if closeIdx := bytes.Index(ctx.tokens[endPos+1:], []byte("</span>{: ")); 0 <= closeIdx {
								if closeIdx2 := bytes.Index(ctx.tokens[endPos+1+closeIdx:], []byte("}")); 0 <= closeIdx2 {
									ctx.tokens = append(ctx.tokens[:endPos+1+closeIdx], ctx.tokens[endPos+1+closeIdx+closeIdx2+1:]...)
									ctx.tokensLen = len(ctx.tokens)
								} else {
									ctx.tokens = bytes.Replace(ctx.tokens, []byte("</span>"), nil, 1)
									ctx.tokensLen = len(ctx.tokens)
								}
							} else if closeIdx := bytes.Index(ctx.tokens[endPos+1:], []byte("</span>")); 0 <= closeIdx {
								ctx.tokens = bytes.Replace(ctx.tokens, []byte("</span>"), nil, 1)
								ctx.tokensLen -= 7
							}
						} else {
							content.Write(n.Tokens)
						}
					} else {
						content.Write(n.Tokens)
					}
				} else if ast.NodeLinkSpace == n.Type {
					content.WriteByte(lex.ItemSpace)
				} else if ast.NodeBackslashContent == n.Type {
					content.WriteString("\\")
					content.Write(n.Tokens)
				} else if ast.NodeHTMLEntity == n.Type {
					content.Write(html.EscapeHTML(n.Tokens))
				} else if ast.NodeInlineMathContent == n.Type {
					content.Write(n.Tokens)
				} else if ast.NodeCodeSpanContent == n.Type {
					content.WriteByte(lex.ItemBacktick)
					content.Write(n.Tokens)
					content.WriteByte(lex.ItemBacktick)
				}
				return ast.WalkContinue
			})
			textTokens = html.UnescapeHTML(content.Bytes())
		}
	}

	ret.AppendChild(&ast.Node{Type: ast.NodeCodeSpanContent, Tokens: textTokens})
	ret.AppendChild(closeMarker)
	ctx.pos = endPos + n
	return
}

func (t *Tree) matchCodeSpanEnd(tokens []byte, num int) (pos int) {
	length := len(tokens)
	for pos < length {
		l := lex.Accept(tokens[pos:], lex.ItemBacktick)
		if num == l {
			next := pos + l
			if length-1 > next && lex.ItemBacktick == tokens[next] {
				continue
			}
			return pos
		}
		if 0 < l {
			pos += l
		} else {
			pos++
		}
	}
	return -1
}
