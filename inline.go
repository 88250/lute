// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

import (
	"strings"

	"github.com/88250/lute/html"
)

// parseInline 解析并生成块节点 block 的行级子节点。
func (t *Tree) parseInline(block *Node, ctx *InlineContext) {
	for ctx.pos < ctx.tokensLen {
		token := ctx.tokens[ctx.pos]
		var n *Node
		switch token {
		case itemBackslash:
			n = t.parseBackslash(ctx)
		case itemBacktick:
			n = t.parseCodeSpan(ctx)
		case itemAsterisk, itemUnderscore, itemTilde:
			t.handleDelim(block, ctx)
		case itemNewline:
			n = t.parseNewline(block, ctx)
		case itemLess:
			n = t.parseAutolink(ctx)
			if nil == n {
				n = t.parseAutoEmailLink(ctx)
				if nil == n {
					n = t.parseInlineHTML(ctx)
				}
			}
		case itemOpenBracket:
			n = t.parseOpenBracket(ctx)
		case itemCloseBracket:
			n = t.parseCloseBracket(ctx)
		case itemAmpersand:
			n = t.parseEntity(ctx)
		case itemBang:
			n = t.parseBang(ctx)
		case itemDollar:
			n = t.parseInlineMath(ctx)
		default:
			n = t.parseText(ctx)
		}

		if nil != n {
			block.AppendChild(n)
		}
	}
	block.tokens = nil
}

func (t *Tree) parseEntity(ctx *InlineContext) (ret *Node) {
	and := []byte{ctx.tokens[ctx.pos]}
	if 2 > ctx.tokensLen || ctx.tokensLen <= ctx.pos+1 {
		ctx.pos++
		return &Node{typ: NodeText, tokens: and}
	}

	start := ctx.pos
	numeric := false
	if 3 < ctx.tokensLen {
		numeric = itemCrosshatch == ctx.tokens[start+1]
	}
	i := ctx.pos
	var token byte
	var endWithSemicolon bool
	for ; i < ctx.tokensLen; i++ {
		token = ctx.tokens[i]
		if isWhitespace(token) {
			break
		}
		if itemSemicolon == token {
			i++
			endWithSemicolon = true
			break
		}
	}

	entityName := bytesToStr(ctx.tokens[start:i])
	if entityValue, ok := html.Entities[entityName]; ok {
		ctx.pos += i - start
		return &Node{typ: NodeText, tokens: strToBytes(entityValue)}
	}

	if !endWithSemicolon {
		ctx.pos++
		return &Node{typ: NodeText, tokens: and}
	}

	if numeric {
		entityNameLen := len(entityName)
		if 10 < entityNameLen || 4 > entityNameLen {
			ctx.pos++
			return &Node{typ: NodeText, tokens: and}
		}

		hex := 'x' == entityName[2] || 'X' == entityName[2]
		if hex {
			if 5 > entityNameLen {
				ctx.pos++
				return &Node{typ: NodeText, tokens: and}
			}
		}
	}

	v := htmlUnescapeString(entityName)
	if v == entityName {
		ctx.pos++
		return &Node{typ: NodeText, tokens: and}
	}
	ctx.pos += i - start
	return &Node{typ: NodeText, tokens: strToBytes(v)}
}

// Try to match close bracket against an opening in the delimiter stack. Add either a link or image, or a plain [ character,
// to block's children. If there is a matching delimiter, remove it from the delimiter stack.
func (t *Tree) parseCloseBracket(ctx *InlineContext) *Node {
	closeBracket := []byte{ctx.tokens[ctx.pos]}
	ctx.pos++
	startPos := ctx.pos

	// get last [ or ![
	opener := ctx.brackets
	if nil == opener {
		return &Node{typ: NodeText, tokens: closeBracket}
	}

	if !opener.active {
		// no matched opener, just return a literal
		// take opener off brackets stack
		t.removeBracket(ctx)
		return &Node{typ: NodeText, tokens: closeBracket}
	}

	// If we got here, open is a potential opener
	isImage := opener.image

	// Check to see if we have a link/image

	var openParen, dest, space, title, closeParen []byte
	savepos := ctx.pos
	matched := false
	// 尝试解析内联链接 [text](url "tile")
	if ctx.pos+1 < ctx.tokensLen && itemOpenParen == ctx.tokens[ctx.pos] {
		ctx.pos++
		isLink := false
		var passed, remains []byte

		for { // 这里使用 for 是为了简化逻辑，不是为了循环
			if isLink, passed, remains = spnl(ctx.tokens[ctx.pos-1:]); !isLink {
				break
			}
			ctx.pos += len(passed)
			if passed, remains, dest = t.context.parseInlineLinkDest(remains); nil == passed {
				break
			}
			if t.context.option.VditorWYSIWYG && (1 > len(dest) || (nil == opener.node.next && !isImage)) {
				break
			}
			ctx.pos += len(passed)
			openParen = passed[0:1]
			closeParen = passed[len(passed)-1:]
			matched = itemCloseParen == passed[len(passed)-1]
			if matched {
				ctx.pos--
				break
			}
			if 1 > len(remains) || !isWhitespace(remains[0]) {
				break
			}
			// 跟空格的话后续尝试 title 解析
			if isLink, passed, remains = spnl(remains); !isLink {
				break
			}
			space = passed
			ctx.pos += len(passed)
			matched = itemCloseParen == remains[0]
			closeParen = remains[0:1]
			if matched {
				break
			}
			ctx.pos++
			validTitle := false
			if validTitle, passed, remains, title = t.context.parseLinkTitle(remains); !validTitle {
				break
			}
			ctx.pos += len(passed)
			isLink, passed, remains = spnl(remains)
			ctx.pos += len(passed)
			matched = isLink && 0 < len(remains) && itemCloseParen == remains[0]
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
		n, _, label := t.context.parseLinkLabel(ctx.tokens[beforelabel:])
		if 2 < n { // label 解析出来的话说明满足格式 [text][label]
			reflabel = label
			ctx.pos += n
		} else if !opener.bracketAfter {
			// [text][] 格式，将 text 视为 label 进行解析
			start := opener.index
			if itemOpenBracket == ctx.tokens[start] {
				// TODO: 链接引用定义 key 还是包括方括号好些 [xxx]
				start++
			}
			reflabel = ctx.tokens[start : startPos-1]
			ctx.pos += len(reflabel)
		}
		if 0 == n {
			ctx.pos = startPos
		}
		if nil != reflabel {
			// 查找链接引用
			var link = t.context.linkRefDef[strings.ToLower(bytesToStr(reflabel))]
			if nil != link {
				dest = link.ChildByType(NodeLinkDest).tokens
				titleNode := link.ChildByType(NodeLinkTitle)
				if nil != titleNode {
					title = titleNode.tokens
				}
				matched = true
			}
		}
	}

	if matched {
		node := &Node{typ: NodeLink, linkType: 0}
		if isImage {
			node.typ = NodeImage
			node.AppendChild(&Node{typ: NodeBang, tokens: opener.node.tokens[:1]})
			opener.node.tokens = opener.node.tokens[1:]
		}
		node.AppendChild(&Node{typ: NodeOpenBracket, tokens: opener.node.tokens})

		var tmp, next *Node
		tmp = opener.node.next
		for nil != tmp {
			next = tmp.next
			tmp.Unlink()
			if NodeText == tmp.typ {
				tmp.typ = NodeLinkText
			}
			node.AppendChild(tmp)
			tmp = next
		}
		node.AppendChild(&Node{typ: NodeCloseBracket, tokens: closeBracket})
		node.AppendChild(&Node{typ: NodeOpenParen, tokens: openParen})
		node.AppendChild(&Node{typ: NodeLinkDest, tokens: dest})
		if nil != space {
			node.AppendChild(&Node{typ: NodeLinkSpace, tokens: space})
		}
		if 0 < len(title) {
			node.AppendChild(&Node{typ: NodeLinkTitle, tokens: title})
		}
		node.AppendChild(&Node{typ: NodeCloseParen, tokens: closeParen})

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
		return &Node{typ: NodeText, tokens: closeBracket}
	}
}

func (t *Tree) parseOpenBracket(ctx *InlineContext) (ret *Node) {
	startPos := ctx.pos
	ctx.pos++
	ret = &Node{typ: NodeText, tokens: ctx.tokens[startPos:ctx.pos]}
	// 将 [ 入栈
	t.addBracket(ret, ctx.pos-1, false, ctx)
	return
}

func (t *Tree) addBracket(node *Node, index int, image bool, ctx *InlineContext) {
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

var backslash = strToBytes("\\")

func (t *Tree) parseBackslash(ctx *InlineContext) *Node {
	if ctx.pos == ctx.tokensLen-1 {
		ctx.pos++
		return &Node{typ: NodeText, tokens: backslash}
	}

	ctx.pos++
	token := ctx.tokens[ctx.pos]
	if itemNewline == token {
		ctx.pos++
		return &Node{typ: NodeHardBreak, tokens: []byte{token}}
	}
	if isASCIIPunct(token) {
		ctx.pos++
		return &Node{typ: NodeText, tokens: []byte{token}}
	}
	return &Node{typ: NodeText, tokens: backslash}
}

func (t *Tree) parseText(ctx *InlineContext) (ret *Node) {
	start := ctx.pos
	for ; ctx.pos < ctx.tokensLen; ctx.pos++ {
		if t.isMarker(ctx.tokens[ctx.pos]) {
			// 遇到潜在的标记符时需要跳出该文本节点，回到行级解析主循环
			break
		}
	}

	ret = &Node{typ: NodeText, tokens: ctx.tokens[start:ctx.pos]}
	return
}

// isMarker 判断 token 是否是潜在的 Markdown 标记符。
func (t *Tree) isMarker(token byte) bool {
	switch token {
	case itemAsterisk, itemUnderscore, itemOpenBracket, itemBang, itemNewline, itemBackslash, itemBacktick, itemLess,
		itemCloseBracket, itemAmpersand, itemTilde, itemDollar:
		return true
	default:
		return false
	}
}

func (t *Tree) parseNewline(block *Node, ctx *InlineContext) (ret *Node) {
	pos := ctx.pos
	ctx.pos++

	hardbreak := false
	// 检查前一个节点的结尾空格，如果大于等于两个则说明是硬换行
	if lastc := block.lastChild; nil != lastc {
		if NodeText == lastc.typ {
			tokens := lastc.tokens
			if valueLen := len(tokens); itemSpace == tokens[valueLen-1] {
				_, lastc.tokens = trimRight(tokens)
				if 1 < valueLen {
					hardbreak = itemSpace == tokens[len(tokens)-2]
				}
			}
		}
	}

	ret = &Node{typ: NodeSoftBreak, tokens: []byte{ctx.tokens[pos]}}
	if hardbreak {
		ret.typ = NodeHardBreak
	}
	return
}
