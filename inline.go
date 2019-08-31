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

// +build !wasm

package lute

import (
	"bytes"
	"html"
	"sync"

	//"html"
	"strings"
)

// parseInlines 解析并生成行级节点。
func (t *Tree) parseInlines() {
	t.walkParseInline(t.Root, nil)
}

// walkParseInline 解析生成节点 node 的行级子节点。
func (t *Tree) walkParseInline(node *BaseNode, wg *sync.WaitGroup) {
	defer recoverPanic(nil)
	if nil != wg {
		defer wg.Done()
	}
	if nil == node {
		return
	}

	// 只有如下几种类型的块节点需要生成行级子节点
	if typ := node.typ; NodeParagraph == typ || NodeHeading == typ || NodeTableCell == typ {
		tokens := node.Tokens()
		if NodeParagraph == typ && nil == tokens {
			// 解析 GFM 表节点后段落内容 tokens 可能会被置换为空，具体可参看函数 Paragraph.Finalize()
			// 在这里从语法树上移除空段落节点
			next := node.next
			node.Unlink()
			// Unlink 会将后一个兄弟节点置空，此处是在在遍历过程中修改树结构，所以需要保持继续迭代后面的兄弟节点
			node.next = (next)
			return
		}

		length := len(tokens)
		if 1 > length {
			return
		}

		ctx := &InlineContext{
			tokens:    tokens,
			tokensLen: length,
		}

		// 生成该块节点的行级子节点
		t.parseInline(node, ctx)

		// 处理该块节点中的强调、加粗和删除线
		t.processEmphasis(nil, ctx)

		// 将连续的文本节点进行合并。
		// 规范只是定义了从输入的 Markdown 文本到输出的 HTML 的解析渲染规则，并未定义中间语法树的规则。
		// 也就是说语法树的节点结构没有标准，可以自行发挥。这里进行文本节点合并主要有两个目的：
		// 1. 减少节点数量，提升后续处理性能
		// 2. 方便后续功能方面的处理，比如 GFM 自动链接解析
		t.mergeText(node)

		if t.context.option.GFMAutoLink {
			t.parseGFMAutoEmailLink(node)
			t.parseGFMAutoLink(node)
		}

		if t.context.option.AutoSpace {
			t.space(node)
		}

		if t.context.option.FixTermTypo {
			t.fixTermTypo(node)
		}
		return
	}

	// 遍历处理子节点，通过并行处理提升性能
	cwg := &sync.WaitGroup{}
	for child := node.firstChild; nil != child; child = child.next {
		cwg.Add(1)
		go t.walkParseInline(child, cwg)
	}
	cwg.Wait()
}

// parseInline 解析并生成块节点 block 的行级子节点。
func (t *Tree) parseInline(block *BaseNode, ctx *InlineContext) {
	for {
		token := ctx.tokens[ctx.pos]
		var n *BaseNode
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
		default:
			n = t.parseText(ctx)
		}

		if nil != n {
			block.AppendChild(block, n)
		}

		if 1 > ctx.tokensLen || ctx.pos >= ctx.tokensLen || itemEnd == ctx.tokens[ctx.pos] {
			return
		}
	}
}

func (t *Tree) parseEntity(ctx *InlineContext) (ret *BaseNode) {
	if 2 > ctx.tokensLen || ctx.tokensLen <= ctx.pos+1 {
		ctx.pos++
		return &BaseNode{typ: NodeText, tokens: toItems("&")}
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

	entityName := fromItems(ctx.tokens[start:i])
	if entityValue, ok := htmlEntities[entityName]; ok { // 通过查表优化
		ctx.pos += i - start
		return &BaseNode{typ: NodeText, tokens: toItems(entityValue)}
	}

	if !endWithSemicolon {
		ctx.pos++
		return &BaseNode{typ: NodeText, tokens: toItems("&")}
	}

	if numeric {
		entityNameLen := len(entityName)
		if 10 < entityNameLen || 4 > entityNameLen {
			ctx.pos++
			return &BaseNode{typ: NodeText, tokens: toItems("&")}
		}

		hex := 'x' == entityName[2] || 'X' == entityName[2]
		if hex {
			if 5 > entityNameLen {
				ctx.pos++
				return &BaseNode{typ: NodeText, tokens: toItems("&")}
			}
		}
	}

	v := html.UnescapeString(entityName)
	if v == entityName {
		ctx.pos++
		return &BaseNode{typ: NodeText, tokens: toItems("&")}
	}
	ctx.pos += i - start
	return &BaseNode{typ: NodeText, tokens: toItems(v)}
}

// Try to match close bracket against an opening in the delimiter stack. Add either a link or image, or a plain [ character,
// to block's children. If there is a matching delimiter, remove it from the delimiter stack.
func (t *Tree) parseCloseBracket(ctx *InlineContext) *BaseNode {
	// get last [ or ![
	opener := ctx.brackets
	if nil == opener {
		ctx.pos++
		// no matched opener, just return a literal
		return &BaseNode{typ: NodeText, tokens: toItems("]")}
	}

	if !opener.active {
		ctx.pos++
		// no matched opener, just return a literal
		// take opener off brackets stack
		t.removeBracket(ctx)
		return &BaseNode{typ: NodeText, tokens: toItems("]")}
	}

	// If we got here, open is a potential opener
	isImage := opener.image

	var dest, title items
	// Check to see if we have a link/image

	startPos := ctx.pos
	savepos := ctx.pos
	matched := false
	// Inline link?
	if ctx.pos+1 < ctx.tokensLen && itemOpenParen == ctx.tokens[ctx.pos+1] {
		ctx.pos++
		isLink := false
		var passed, remains items
		if isLink, passed, remains = ctx.tokens[ctx.pos:].spnl(); isLink {
			ctx.pos += len(passed)
			if passed, remains, dest = t.context.parseInlineLink(remains); nil != passed {
				ctx.pos += len(passed)
				if 0 < len(remains) && isWhitespace(remains[0]) { // 跟空格的话后续尝试按 title 解析
					ctx.pos++
					if isLink, passed, remains = remains.spnl(); isLink {
						ctx.pos += len(passed)
						validTitle := false
						if validTitle, passed, remains, title = t.context.parseLinkTitle(remains); validTitle {
							ctx.pos += len(passed)
							isLink, passed, remains = remains.spnl()
							ctx.pos += len(passed)
							matched = isLink && itemCloseParen == remains[0]
						}
					}
				} else { // 没有 title
					ctx.pos--
					matched = true
				}
			}
		}
		if !matched {
			ctx.pos = savepos
		}
	}

	var reflabel string
	if !matched {
		// 尝试解析链接 label
		var beforelabel = ctx.pos + 1
		passed, _, label := t.context.parseLinkLabel(ctx.tokens[beforelabel:])
		var n = len(passed)
		if n > 0 { // label 解析出来的话说明满足格式 [text][label]
			reflabel = label
			ctx.pos += n + 1
		} else if !opener.bracketAfter {
			// [text][] 或者 [text][] 格式，将第一个 text 视为 label 进行解析
			passed = ctx.tokens[opener.index:startPos]
			reflabel = fromItems(passed)
			if len(passed) > 0 && ctx.tokensLen > beforelabel && itemOpenBracket == ctx.tokens[beforelabel] {
				// [text][] 格式，跳过 []
				ctx.pos += 2
			}
		}

		if "" != reflabel {
			// 查找链接引用
			var link = t.context.linkRefDef[strings.ToLower(reflabel)]
			if nil != link {
				dest = link.Destination
				title = link.Title
				matched = true
			}
		}
	}

	if matched {
		var node *BaseNode
		if isImage {
			node = &BaseNode{typ: NodeImage, Destination: dest, Title: title}
		} else {
			node = &BaseNode{typ: NodeLink, Destination: dest, Title: title}
		}

		var tmp, next *BaseNode
		tmp = opener.node.next
		for nil != tmp {
			next = tmp.next
			tmp.Unlink()
			node.AppendChild(node, tmp)
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

		ctx.pos++
		return node
	} else { // no match
		t.removeBracket(ctx) // remove this opener from stack
		ctx.pos = startPos
		ctx.pos++
		return &BaseNode{typ: NodeText, tokens: toItems("]")}
	}
}

func (t *Tree) parseOpenBracket(ctx *InlineContext) (ret *BaseNode) {
	ctx.pos++
	ret = &BaseNode{typ: NodeText, tokens: toItems("[")}
	// 将 [ 入栈
	t.addBracket(ret, ctx.pos, false, ctx)
	return
}

func (t *Tree) addBracket(node *BaseNode, index int, image bool, ctx *InlineContext) {
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

func (t *Tree) parseBackslash(ctx *InlineContext) (ret *BaseNode) {
	if ctx.tokensLen-1 > ctx.pos {
		ctx.pos++
	}
	token := ctx.tokens[ctx.pos]
	if itemNewline == token {
		ret = &BaseNode{typ: NodeHardBreak}
		ctx.pos++
	} else if isASCIIPunct(token) {
		ret = &BaseNode{typ: NodeText, tokens: items{token}}
		ctx.pos++
	} else {
		ret = &BaseNode{typ: NodeText, tokens: toItems("\\")}
	}

	return
}

func (t *Tree) parseText(ctx *InlineContext) (ret *BaseNode) {
	var token byte
	start := ctx.pos
	for ; ctx.pos < ctx.tokensLen; ctx.pos++ {
		token = ctx.tokens[ctx.pos]
		if t.isMarker(token) {
			// 遇到潜在的标记符时需要跳出 text，回到行级解析主循环
			break
		}
	}

	ret = &BaseNode{typ: NodeText, tokens: ctx.tokens[start:ctx.pos]}
	return
}

// isMarker 判断 token 是否是潜在的 Markdown 标记。
func (t *Tree) isMarker(token byte) bool {
	return itemAsterisk == token || itemUnderscore == token || itemOpenBracket == token || itemBang == token ||
		itemNewline == token || itemBackslash == token || itemBacktick == token ||
		itemLess == token || itemCloseBracket == token || itemAmpersand == token || itemTilde == token
}

func (t *Tree) parseNewline(block *BaseNode, ctx *InlineContext) (ret *BaseNode) {
	ctx.pos++

	hardbreak := false
	// 检查前一个节点的结尾空格，如果大于等于两个则说明是硬换行
	if lastc := block.lastChild; nil != lastc {
		if NodeText == lastc.typ {
			tokens := lastc.tokens
			if valueLen := len(tokens); itemSpace == tokens[valueLen-1] {
				lastc.SetTokens(bytes.TrimRight(tokens, " \t\n"))
				if 1 < valueLen {
					hardbreak = itemSpace == tokens[len(tokens)-2]
				}
			}
		}
	}

	if hardbreak {
		ret = &BaseNode{typ: NodeHardBreak}
	} else {
		ret = &BaseNode{typ: NodeSoftBreak}
	}
	return
}
