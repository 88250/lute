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
	"github.com/b3log/lute/html"
	"strings"
)

// parseInline 解析并生成块节点 block 的行级子节点。
func (t *Tree) parseInline(block *Node, ctx *InlineContext) {
	for ctx.pos < ctx.tokensLen {
		token := ctx.tokens[ctx.pos]
		var n *Node
		switch token.term {
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
}

var and = strToItems("&")

func (t *Tree) parseEntity(ctx *InlineContext) (ret *Node) {
	if 2 > ctx.tokensLen || ctx.tokensLen <= ctx.pos+1 {
		ctx.pos++
		return &Node{typ: NodeText, tokens: and}
	}

	start := ctx.pos
	numeric := false
	if 3 < ctx.tokensLen {
		numeric = itemCrosshatch == ctx.tokens[start+1].term
	}
	i := ctx.pos
	var token byte
	var endWithSemicolon bool
	for ; i < ctx.tokensLen; i++ {
		token = ctx.tokens[i].term
		if isWhitespace(token) {
			break
		}
		if itemSemicolon == token {
			i++
			endWithSemicolon = true
			break
		}
	}

	entityName := itemsToStr(ctx.tokens[start:i])
	if entityValue, ok := html.Entities[entityName]; ok {
		ctx.pos += i - start
		return &Node{typ: NodeText, tokens: strToItems(entityValue)}
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
	return &Node{typ: NodeText, tokens: strToItems(v)}
}

var closeBracket = strToItems("]")

// Try to match close bracket against an opening in the delimiter stack. Add either a link or image, or a plain [ character,
// to block's children. If there is a matching delimiter, remove it from the delimiter stack.
func (t *Tree) parseCloseBracket(ctx *InlineContext) *Node {
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

	var dest, title items
	savepos := ctx.pos
	matched := false
	// 尝试解析内联链接 [text](url "tile")
	if ctx.pos+1 < ctx.tokensLen && itemOpenParen == ctx.tokens[ctx.pos].term {
		ctx.pos++
		isLink := false
		var passed, remains items

		for { // 这里使用 for 是为了简化逻辑，不是为了循环
			if isLink, passed, remains = ctx.tokens[ctx.pos-1:].spnl(); !isLink {
				break
			}
			ctx.pos += len(passed)
			if passed, remains, dest = t.context.parseInlineLink(remains); nil == passed {
				break
			}
			ctx.pos += len(passed)
			matched = itemCloseParen == passed[len(passed)-1].term
			if matched {
				ctx.pos--
				break
			}
			if 1 > len(remains) || !isWhitespace(remains[0].term) {
				break
			}
			// 跟空格的话后续尝试 title 解析
			ctx.pos++
			if isLink, passed, remains = remains.spnl(); !isLink {
				break
			}
			ctx.pos += len(passed)
			matched = itemCloseParen == remains[0].term
			if matched {
				break
			}
			validTitle := false
			if validTitle, passed, remains, title = t.context.parseLinkTitle(remains); !validTitle {
				break
			}
			ctx.pos += len(passed)
			isLink, passed, remains = remains.spnl()
			ctx.pos += len(passed)
			matched = isLink && itemCloseParen == remains[0].term
			break
		}
		if !matched {
			ctx.pos = savepos
		}
	}

	var reflabel items
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
			if itemOpenBracket == ctx.tokens[start].term {
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
			var link = t.context.linkRefDef[strings.ToLower(itemsToStr(reflabel))]
			if nil != link {
				dest = link.destination
				title = link.title
				matched = true
			}
		}
	}

	if matched {
		var node *Node
		if isImage {
			node = &Node{typ: NodeImage, destination: dest, title: title}
		} else {
			node = &Node{typ: NodeLink, destination: dest, title: title}
		}

		var tmp, next *Node
		tmp = opener.node.next
		for nil != tmp {
			next = tmp.next
			tmp.Unlink()
			node.AppendChild(tmp)
			tmp = next
		}

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

var openBracket = strToItems("[")

func (t *Tree) parseOpenBracket(ctx *InlineContext) (ret *Node) {
	ctx.pos++
	ret = &Node{typ: NodeText, tokens: openBracket}
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

var backslash = strToItems("\\")

func (t *Tree) parseBackslash(ctx *InlineContext) *Node {
	if ctx.pos == ctx.tokensLen-1 {
		ctx.pos++
		return &Node{typ: NodeText, tokens: backslash}
	}

	ctx.pos++
	token := ctx.tokens[ctx.pos]
	if itemNewline == token.term {
		ctx.pos++
		return &Node{typ: NodeHardBreak, tokens: items{token}}
	}
	if isASCIIPunct(token.term) {
		ctx.pos++
		return &Node{typ: NodeText, tokens: items{token}}
	}
	return &Node{typ: NodeText, tokens: backslash}
}

func (t *Tree) parseText(ctx *InlineContext) (ret *Node) {
	var token byte
	start := ctx.pos
	for ; ctx.pos < ctx.tokensLen; ctx.pos++ {
		token = ctx.tokens[ctx.pos].term
		if t.isMarker(token) {
			// 遇到潜在的标记符时需要跳出该文本节点，回到行级解析主循环
			break
		}
	}

	ret = &Node{typ: NodeText, tokens: ctx.tokens[start:ctx.pos]}
	return
}

// isMarker 判断 token 是否是潜在的 Markdown 标记符。
func (t *Tree) isMarker(token byte) bool {
	return itemAsterisk == token || itemUnderscore == token || itemOpenBracket == token || itemBang == token ||
		itemNewline == token || itemBackslash == token || itemBacktick == token ||
		itemLess == token || itemCloseBracket == token || itemAmpersand == token || itemTilde == token ||
		itemDollar == token
}

func (t *Tree) parseNewline(block *Node, ctx *InlineContext) *Node {
	ctx.pos++

	hardbreak := false
	// 检查前一个节点的结尾空格，如果大于等于两个则说明是硬换行
	if lastc := block.lastChild; nil != lastc {
		if NodeText == lastc.typ {
			tokens := lastc.tokens
			if valueLen := len(tokens); itemSpace == tokens[valueLen-1].term {
				_, lastc.tokens = trimRight(tokens)
				if 1 < valueLen {
					hardbreak = itemSpace == tokens[len(tokens)-2].term
				}
			}
		}
	}

	if hardbreak {
		return &Node{typ: NodeHardBreak}
	}
	return &Node{typ: NodeSoftBreak}
}
