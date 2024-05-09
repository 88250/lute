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
	"strconv"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

// parseInline 解析并生成块节点 block 的行级子节点。
func (t *Tree) parseInline(block *ast.Node, ctx *InlineContext) {
	for ctx.pos < ctx.tokensLen {
		token := ctx.tokens[ctx.pos]
		var n *ast.Node
		switch token {
		case lex.ItemBackslash:
			n = t.parseBackslash(block, ctx)
		case lex.ItemBacktick:
			n = t.parseCodeSpan(block, ctx)
		case lex.ItemAsterisk, lex.ItemUnderscore, lex.ItemTilde, lex.ItemEqual, lex.ItemCrosshatch:
			t.handleDelim(block, ctx)
		case lex.ItemCaret:
			if t.Context.ParseOption.Sup {
				t.handleDelim(block, ctx)
			} else {
				n = t.parseText(ctx)
			}
		case lex.ItemNewline:
			n = t.parseNewline(block, ctx)
		case lex.ItemLess:
			if n = t.parseAutolink(ctx); nil == n {
				if n = t.parseAutoEmailLink(ctx); nil == n {
					if n = t.parseFileAnnotationRef(ctx); nil == n {
						n = t.parseInlineHTML(ctx)
						if t.Context.ParseOption.ProtyleWYSIWYG && nil != n && ast.NodeInlineHTML == n.Type {
							// Protyle 中不存在内联 HTML，使用文本
							n.Type = ast.NodeText
						}
					}
				}
			}
		case lex.ItemOpenBracket:
			n = t.parseOpenBracket(ctx)
		case lex.ItemCloseBracket:
			n = t.parseCloseBracket(ctx)
		case lex.ItemAmpersand:
			n = t.parseEntity(ctx)
		case lex.ItemBang:
			n = t.parseBang(ctx)
		case lex.ItemDollar:
			n = t.parseInlineMath(ctx)
		case lex.ItemOpenBrace:
			n = t.parseHeadingID(block, ctx)
		case lex.ItemOpenParen:
			n = t.parseBlockRef(ctx)
		default:
			n = t.parseText(ctx)
		}

		if nil != n {
			block.AppendChild(n)
		}
	}
	block.Tokens = nil
}

func (t *Tree) parseEntity(ctx *InlineContext) (ret *ast.Node) {
	and := []byte{ctx.tokens[ctx.pos]}
	if 2 > ctx.tokensLen || ctx.tokensLen <= ctx.pos+1 {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: and}
	}

	start := ctx.pos
	numeric := false
	if 3 < ctx.tokensLen {
		numeric = lex.ItemCrosshatch == ctx.tokens[start+1]
	}
	i := ctx.pos
	var token byte
	var endWithSemicolon bool
	for ; i < ctx.tokensLen; i++ {
		token = ctx.tokens[i]
		if lex.IsWhitespace(token) {
			break
		}
		if lex.ItemSemicolon == token {
			i++
			endWithSemicolon = true
			break
		}
	}

	if !endWithSemicolon {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: and}
	}

	entityName := util.BytesToStr(ctx.tokens[start:i])
	if numeric {
		entityNameLen := len(entityName)
		if 10 < entityNameLen || 4 > entityNameLen {
			ctx.pos++
			return &ast.Node{Type: ast.NodeText, Tokens: and}
		}

		if ('x' == entityName[2] || 'X' == entityName[2]) && 5 > entityNameLen {
			ctx.pos++
			return &ast.Node{Type: ast.NodeText, Tokens: and}
		}
	}

	v := html.HtmlUnescapeString(entityName)
	if v == entityName {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: and}
	}
	ctx.pos += i - start
	return &ast.Node{Type: ast.NodeHTMLEntity, Tokens: util.StrToBytes(v), HtmlEntityTokens: util.StrToBytes(entityName)}
}

// Try to match close bracket against an opening in the delimiter stack. Add either a link or image, or a plain [ character,
// to block's children. If there is a matching delimiter, remove it from the delimiter stack.
func (t *Tree) parseCloseBracket(ctx *InlineContext) *ast.Node {
	closeBracket := []byte{ctx.tokens[ctx.pos]}
	ctx.pos++
	startPos := ctx.pos

	// 获取最新一个 [ 或者 ![
	opener := ctx.brackets
	if nil == opener {
		return &ast.Node{Type: ast.NodeText, Tokens: closeBracket}
	}

	if !opener.active {
		t.removeBracket(ctx)
		return &ast.Node{Type: ast.NodeText, Tokens: closeBracket}
	}

	isImage := opener.image

	// 检查是否满足链接或者图片规则

	var openParen, dest, space, title, closeParen []byte
	savepos := ctx.pos
	matched := false
	// 尝试解析内联链接 [text](url "tile")
	if ctx.pos+1 < ctx.tokensLen && lex.ItemOpenParen == ctx.tokens[ctx.pos] {
		ctx.pos++
		isLink := false
		var passed, remains []byte

		for { // 这里使用 for 是为了简化逻辑，不是为了循环
			if isLink, passed, remains = lex.Spnl(ctx.tokens[ctx.pos-1:]); !isLink {
				break
			}
			ctx.pos += len(passed)
			if passed, remains, dest = t.Context.parseInlineLinkDest(remains); nil == passed {
				break
			}
			if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG {
				if !isImage && nil == opener.node.Next {
					break
				}
			}
			ctx.pos += len(passed)
			openParen = passed[0:1]
			closeParen = passed[len(passed)-1:]
			matched = lex.ItemCloseParen == passed[len(passed)-1]
			if matched && 1 < len(remains) {
				// 如果 passed 是 ) 结尾，则继续判断 remains 是否以 空格" 开头
				// 解决 [foo](bar.com(baz) "bar.com(baz)") 这种情况，测试用例 debug_test.go #75
				matched = !lex.IsWhitespace(remains[0]) && lex.ItemDoublequote != remains[1]
			}

			if matched {
				ctx.pos--
				break
			}
			if 1 > len(remains) || !lex.IsWhitespace(remains[0]) {
				break
			}
			// 跟空格的话后续尝试 title 解析
			if isLink, passed, remains = lex.Spnl(remains); !isLink {
				break
			}
			space = passed
			ctx.pos += len(passed)
			matched = lex.ItemCloseParen == remains[0]
			closeParen = remains[0:1]
			if matched {
				break
			}
			ctx.pos++
			validTitle := false
			if validTitle, passed, remains, title = t.Context.parseLinkTitle(remains); !validTitle {
				break
			}
			ctx.pos += len(passed)
			isLink, passed, remains = lex.Spnl(remains)
			ctx.pos += len(passed)
			matched = isLink && 0 < len(remains)
			if matched {
				if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG {
					if bytes.HasPrefix(remains, []byte(editor.Caret+")")) {
						if 0 < len(title) {
							// 将 ‸) 换位为 )‸
							remains = remains[len([]byte(editor.Caret+")")):]
							remains = append([]byte(")"+editor.Caret), remains...)
							copy(ctx.tokens[ctx.pos-1:], remains) // 同时也将 tokens 换位，后续解析从插入符位置开始
						} else {
							// 将 ""‸ 换位为 "‸"
							title = editor.CaretTokens
							remains = remains[len(editor.CaretTokens):]
							ctx.pos += 3
						}
					} else if bytes.HasPrefix(remains, []byte(")"+editor.Caret)) {
						if 0 == len(title) {
							// 将 "")‸ 换位为 "‸")
							title = editor.CaretTokens
							remains = bytes.ReplaceAll(remains, editor.CaretTokens, nil)
							ctx.pos += 3
						}
					}
				}
				matched = lex.ItemCloseParen == remains[0]
			}
			closeParen = remains[0:]
			break
		}
		if !matched {
			ctx.pos = savepos
		}
	}

	var reflabel []byte
	var linkType int
	if !matched {
		// 尝试解析链接 label
		var beforelabel = ctx.pos
		n, _, label := t.Context.parseLinkLabel(ctx.tokens[beforelabel:])
		if 2 < n { // label 解析出来的话说明满足格式 [text][label]
			reflabel = label
			ctx.pos += n
		} else if !opener.bracketAfter {
			// [text][] 格式，将 text 视为 label 进行解析
			start := opener.index
			if lex.ItemOpenBracket == ctx.tokens[start] {
				start++
			}
			reflabel = ctx.tokens[start : startPos-1]
			ctx.pos += 2
		}
		if 0 == n {
			ctx.pos = startPos
		}
		if nil != reflabel {
			if t.Context.ParseOption.Footnotes {
				// 查找脚注
				if idx, footnotesDef := t.FindFootnotesDef(reflabel); nil != footnotesDef {
					t.removeBracket(ctx)

					if t.Context.ParseOption.Sup && nil != opener.node.Next.Next {
						opener.node.Next.Next.Unlink() // label
						opener.node.Next.Unlink()      // ^
					} else {
						opener.node.Next.Unlink() // ^label
					}
					opener.node.Unlink() // [

					refId := strconv.Itoa(idx)
					refsLen := len(footnotesDef.FootnotesRefs)
					if 0 < refsLen {
						refId += ":" + strconv.Itoa(refsLen+1)
					}
					ref := &ast.Node{Type: ast.NodeFootnotesRef, Tokens: reflabel, FootnotesRefId: refId, FootnotesRefLabel: bytes.ReplaceAll(reflabel, editor.CaretTokens, nil)}
					footnotesDef.FootnotesRefs = append(footnotesDef.FootnotesRefs, ref)
					return ref
				}
			}

			// 查找链接引用定义
			if link := t.FindLinkRefDefLink(reflabel); nil != link {
				dest = link.ChildByType(ast.NodeLinkDest).Tokens
				titleNode := link.ChildByType(ast.NodeLinkTitle)
				if nil != titleNode {
					title = titleNode.Tokens
				}
				matched = true
				linkType = 3
			}
		}
	}

	if matched {
		node := &ast.Node{Type: ast.NodeLink, LinkType: linkType, LinkRefLabel: reflabel}
		if isImage {
			node.Type = ast.NodeImage
			node.AppendChild(&ast.Node{Type: ast.NodeBang, Tokens: opener.node.Tokens[:1]})
			opener.node.Tokens = opener.node.Tokens[1:]
		}
		node.AppendChild(&ast.Node{Type: ast.NodeOpenBracket, Tokens: opener.node.Tokens})

		var tmp, next *ast.Node
		tmp = opener.node.Next
		for nil != tmp {
			next = tmp.Next
			tmp.Unlink()
			if ast.NodeText == tmp.Type {
				tmp.Type = ast.NodeLinkText
			}
			node.AppendChild(tmp)
			tmp = next
		}
		node.AppendChild(&ast.Node{Type: ast.NodeCloseBracket, Tokens: closeBracket})
		node.AppendChild(&ast.Node{Type: ast.NodeOpenParen, Tokens: openParen})
		node.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: dest})
		if nil != space {
			node.AppendChild(&ast.Node{Type: ast.NodeLinkSpace, Tokens: space})
		}
		if 0 < len(title) {
			node.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: title})
		}
		node.AppendChild(&ast.Node{Type: ast.NodeCloseParen, Tokens: closeParen})
		t.processEmphasis(opener.previousDelimiter, ctx)
		t.removeBracket(ctx)
		opener.node.Unlink()

		// We remove this bracket and processEmphasis will remove later delimiters.
		// Now, for a link, we also deactivate earlier link openers.
		// (no links in links)
		if !isImage {
			opener = ctx.brackets
			for nil != opener {
				if !opener.image {
					opener.active = false // deactivate this opener
				}
				opener = opener.previous
			}
		}

		return node
	} else { // 没有匹配到
		t.removeBracket(ctx)
		ctx.pos = startPos
		return &ast.Node{Type: ast.NodeText, Tokens: closeBracket}
	}
}

func (t *Tree) parseOpenBracket(ctx *InlineContext) (ret *ast.Node) {
	startPos := ctx.pos
	ctx.pos++
	ret = &ast.Node{Type: ast.NodeText, Tokens: ctx.tokens[startPos:ctx.pos]}
	// 将 [ 入栈
	t.addBracket(ret, ctx.pos-1, false, ctx)
	return
}

func (t *Tree) addBracket(node *ast.Node, index int, image bool, ctx *InlineContext) {
	if nil != ctx.brackets {
		ctx.brackets.bracketAfter = true
	}

	ctx.brackets = &delimiter{
		node:              node,
		previous:          ctx.brackets,
		previousDelimiter: ctx.delimiters,
		index:             index,
		image:             image,
		active:            true,
	}
}

func (t *Tree) removeBracket(ctx *InlineContext) {
	ctx.brackets = ctx.brackets.previous
}
