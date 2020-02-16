// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"strings"

	"github.com/88250/lute/ast"
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
			n = t.parseCodeSpan(ctx)
		case lex.ItemAsterisk, lex.ItemUnderscore, lex.ItemTilde:
			t.handleDelim(block, ctx)
		case lex.ItemNewline:
			n = t.parseNewline(block, ctx)
		case lex.ItemLess:
			n = t.parseAutolink(ctx)
			if nil == n {
				n = t.parseAutoEmailLink(ctx)
				if nil == n {
					n = t.parseInlineHTML(ctx)
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

	entityName := util.BytesToStr(ctx.tokens[start:i])
	if entityValue, ok := html.Entities[entityName]; ok {
		ctx.pos += i - start
		return &ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(entityValue)}
	}

	if !endWithSemicolon {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: and}
	}

	if numeric {
		entityNameLen := len(entityName)
		if 10 < entityNameLen || 4 > entityNameLen {
			ctx.pos++
			return &ast.Node{Type: ast.NodeText, Tokens: and}
		}

		hex := 'x' == entityName[2] || 'X' == entityName[2]
		if hex {
			if 5 > entityNameLen {
				ctx.pos++
				return &ast.Node{Type: ast.NodeText, Tokens: and}
			}
		}
	}

	v := util.HtmlUnescapeString(entityName)
	if v == entityName {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: and}
	}
	ctx.pos += i - start
	return &ast.Node{Type: ast.NodeText, Tokens: util.StrToBytes(v)}
}

// Try to match close bracket against an opening in the delimiter stack. Add either a link or image, or a plain [ character,
// to block's children. If there is a matching delimiter, remove it from the delimiter stack.
func (t *Tree) parseCloseBracket(ctx *InlineContext) *ast.Node {
	closeBracket := []byte{ctx.tokens[ctx.pos]}
	ctx.pos++
	startPos := ctx.pos

	// get last [ or ![
	opener := ctx.brackets
	if nil == opener {
		return &ast.Node{Type: ast.NodeText, Tokens: closeBracket}
	}

	if !opener.active {
		// no matched opener, just return a literal
		// take opener off brackets stack
		t.removeBracket(ctx)
		return &ast.Node{Type: ast.NodeText, Tokens: closeBracket}
	}

	// If we got here, open is a potential opener
	isImage := opener.image

	// Check to see if we have a link/image

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
			if t.Context.Option.VditorWYSIWYG {
				if !isImage && nil == opener.node.Next {
					break
				}
			}
			ctx.pos += len(passed)
			openParen = passed[0:1]
			closeParen = passed[len(passed)-1:]
			matched = lex.ItemCloseParen == passed[len(passed)-1]
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
			matched = isLink && 0 < len(remains) && lex.ItemCloseParen == remains[0]
			closeParen = remains[0:]
			break
		}
		if !matched {
			ctx.pos = savepos
		}
	}

	var reflabel []byte
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
				// TODO: 链接引用定义 key 还是包括方括号好些 [xxx]
				start++
			}
			reflabel = ctx.tokens[start : startPos-1]
			ctx.pos += 2
		}
		if 0 == n {
			ctx.pos = startPos
		}
		if nil != reflabel {
			if t.Context.Option.Footnotes {
				// 查找脚注
				if idx, footnotesDef := t.Context.FindFootnotesDef(reflabel); nil != footnotesDef {
					t.removeBracket(ctx)
					opener.node.Next.Unlink() // ^label
					opener.node.Unlink()      // [

					refId := strconv.Itoa(idx)
					refsLen := len(footnotesDef.FootnotesRefs)
					if 0 < refsLen {
						refId += ":" + strconv.Itoa(refsLen+1)
					}
					ref := &ast.Node{Type: ast.NodeFootnotesRef, Tokens: bytes.ToLower(reflabel), FootnotesRefId: refId}
					footnotesDef.FootnotesRefs = append(footnotesDef.FootnotesRefs, ref)
					return ref
				}
			}

			// 查找链接引用
			if link := t.Context.linkRefDefs[strings.ToLower(util.BytesToStr(reflabel))]; nil != link {
				dest = link.ChildByType(ast.NodeLinkDest).Tokens
				titleNode := link.ChildByType(ast.NodeLinkTitle)
				if nil != titleNode {
					title = titleNode.Tokens
				}
				matched = true
			}
		}
	}

	if matched {
		node := &ast.Node{Type: ast.NodeLink, LinkType: 0}
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
	} else { // no match
		t.removeBracket(ctx) // remove this opener from stack
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

var backslash = util.StrToBytes("\\")

func (t *Tree) parseBackslash(block *ast.Node, ctx *InlineContext) *ast.Node {
	if ctx.pos == ctx.tokensLen-1 {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: backslash}
	}

	ctx.pos++
	token := ctx.tokens[ctx.pos]
	if lex.ItemNewline == token {
		ctx.pos++
		return &ast.Node{Type: ast.NodeHardBreak, Tokens: []byte{token}}
	}
	if lex.IsASCIIPunct(token) {
		ctx.pos++
		n := &ast.Node{Type: ast.NodeBackslash}
		block.AppendChild(n)
		n.AppendChild(&ast.Node{Type: ast.NodeBackslashContent, Tokens: []byte{token}})
		return nil
	}
	return &ast.Node{Type: ast.NodeText, Tokens: backslash}
}

func (t *Tree) parseText(ctx *InlineContext) *ast.Node {
	start := ctx.pos
	for ; ctx.pos < ctx.tokensLen; ctx.pos++ {
		if t.isMarker(ctx.tokens[ctx.pos]) {
			// 遇到潜在的标记符时需要跳出该文本节点，回到行级解析主循环
			break
		}
	}
	return &ast.Node{Type: ast.NodeText, Tokens: ctx.tokens[start:ctx.pos]}
}

// IsMarker 判断 token 是否是潜在的 Markdown 标记符。
func (t *Tree) isMarker(token byte) bool {
	switch token {
	case lex.ItemAsterisk, lex.ItemUnderscore, lex.ItemOpenBracket, lex.ItemBang, lex.ItemNewline, lex.ItemBackslash, lex.ItemBacktick, lex.ItemLess,
		lex.ItemCloseBracket, lex.ItemAmpersand, lex.ItemTilde, lex.ItemDollar:
		return true
	default:
		return false
	}
}

func (t *Tree) parseNewline(block *ast.Node, ctx *InlineContext) (ret *ast.Node) {
	pos := ctx.pos
	ctx.pos++

	hardbreak := false
	// 检查前一个节点的结尾空格，如果大于等于两个则说明是硬换行
	if lastc := block.LastChild; nil != lastc {
		if ast.NodeText == lastc.Type {
			tokens := lastc.Tokens
			if valueLen := len(tokens); lex.ItemSpace == tokens[valueLen-1] {
				_, lastc.Tokens = lex.TrimRight(tokens)
				if 1 < valueLen {
					hardbreak = lex.ItemSpace == tokens[len(tokens)-2]
				}
			}
		}
	}

	ret = &ast.Node{Type: ast.NodeSoftBreak, Tokens: []byte{ctx.tokens[pos]}}
	if hardbreak {
		ret.Type = ast.NodeHardBreak
	}
	return
}
