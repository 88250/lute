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
	"bytes"
	"strings"
)

func (t *Tree) parseGFMAutoEmailLink(node *BaseNode) {
	for child := node.firstChild; nil != child; {
		next := child.next
		if NodeText == child.typ && nil != child.parent &&
			NodeLink != child.parent.typ /* 不处理链接 label */ {
			t.parseGFMAutoEmailLink0(child)
		} else {
			t.parseGFMAutoEmailLink(child) // 递归处理子节点
		}
		child = next
	}
}

func (t *Tree) parseGFMAutoLink(node *BaseNode) {
	for child := node.firstChild; nil != child; {
		next := child.Next()
		if NodeText == child.typ && nil != child.parent &&
			NodeLink != child.parent.typ /* 不处理链接 label */ {
			t.parseGFMAutoLink0(child)
		} else {
			t.parseGFMAutoLink(child) // 递归处理子节点
		}
		child = next
	}
}

func (t *Tree) parseGFMAutoEmailLink0(node *BaseNode) {
	tokens := node.Tokens()
	if 0 >= bytes.Index(tokens, []byte("@")) {
		return
	}

	var i, j, k, atIndex int
	var token byte
	length := len(tokens)

	// 按空白分隔成多组并进行处理
loopPart:
	for i < length {
		var group items
		atIndex = 0
		j = i

		// 积攒组直到遇到空白符
		for ; j < length; j++ {
			token = tokens[j]
			if !isWhitespace(token) {
				group = append(group, token)
				if '@' == token {
					// 记录 @ 符号在组中的绝对位置，后面会用到
					atIndex = j - i
				}
				continue
			}
			break
		}
		if i == j {
			// 说明积攒组时第一个字符就是空白符，那就把这个空白符作为一个文本节点插到前面
			text := &BaseNode{typ: NodeText, tokens: items{token}}
			node.InsertBefore(node, text)
			i++
			continue
		}

		// 移动主循环下标
		i = j

		if 0 >= atIndex {
			text := &BaseNode{typ: NodeText, tokens: group}
			node.InsertBefore(node, text)
			continue
		}

		// 至此说明这一组中包含了 @，可尝试进行邮件地址解析

		k = 0
		for ; k < atIndex; k++ {
			token = group[k]
			if !t.isValidEmailSegment1(token) {
				text := &BaseNode{typ: NodeText, tokens: group}
				node.InsertBefore(node, text)
				continue loopPart
			}
		}

		k++ // 跳过 @ 检查后面的部分
		for ; k < len(group); k++ {
			token = group[k]
			if !t.isValidEmailSegment2(token) {
				text := &BaseNode{typ: NodeText, tokens: group}
				node.InsertBefore(node, text)
				continue loopPart
			}
		}

		if itemDot == token {
			// 如果以 . 结尾则剔除该 .
			lastIndex := len(group) - 1
			group = group[:lastIndex]
			link := &BaseNode{typ: NodeLink, Destination: append(items("mailto:"), group...)}
			text := &BaseNode{typ: NodeText, tokens: group}
			link.AppendChild(link, text)
			node.InsertBefore(node, link)
			// . 作为文本节点插入
			text = &BaseNode{typ: NodeText, tokens: items{itemDot}}
			node.InsertBefore(node, text)
		} else if itemHyphen == token || itemUnderscore == token {
			// 如果以 - 或者 _ 结尾则整个串都不能算作邮件链接
			text := &BaseNode{typ: NodeText, tokens: group}
			node.InsertBefore(node, text)
			continue loopPart
		} else {
			// 以字母或者数字结尾
			link := &BaseNode{typ: NodeLink, Destination: append(items("mailto:"), group...)}
			text := &BaseNode{typ: NodeText, tokens: group}
			link.AppendChild(link, text)
			node.InsertBefore(node, link)
		}
	}

	// 处理完后传入的文本节点 node 已经被拆分为多个节点，所以可以移除自身
	node.Unlink()
	return
}

func (t *Tree) isValidEmailSegment1(token byte) bool {
	return isASCIILetterNumHyphen(token) || itemDot == token || itemPlus == token || itemUnderscore == token
}

func (t *Tree) isValidEmailSegment2(token byte) bool {
	return isASCIILetterNumHyphen(token) || itemDot == token || itemUnderscore == token
}

func (t *Tree) parseGFMAutoLink0(node *BaseNode) {
	tokens := node.Tokens()
	var i, j, k int
	length := len(tokens)
	if 8 > length { // 太短的情况肯定不可能有链接
		return
	}

	var token byte
	var consumed = make(items, 0, 256)
	var tmp = make(items, 0, 256)
	www := false
	for i < length {
		token = tokens[i]
		var protocol items

		// 检查前缀
		tmp = tokens[i:]
		tmpLen := len(tmp)
		if 8 <= tmpLen /* www.x.xx */ && 'w' == tmp[0] && 'w' == tmp[1] && 'w' == tmp[2] && '.' == tmp[3] {
			protocol = items("http://")
			www = true
		} else if 11 <= tmpLen /* http://x.xx */ && 'h' == tmp[0] && 't' == tmp[1] && 't' == tmp[2] && 'p' == tmp[3] && ':' == tmp[4] && '/' == tmp[5] && '/' == tmp[6] {
			protocol = items("http://")
			i += 7
		} else if 12 <= tmpLen /* https://x.xx */ && 'h' == tmp[0] && 't' == tmp[1] && 't' == tmp[2] && 'p' == tmp[3] && 's' == tmp[4] && ':' == tmp[5] && '/' == tmp[6] && '/' == tmp[7] {
			protocol = items("https://")
			i += 8
		} else if 10 <= tmpLen /* ftp://x.xx */ && 'f' == tmp[0] && 't' == tmp[1] && 'p' == tmp[2] && ':' == tmp[3] && '/' == tmp[4] && '/' == tmp[5] {
			protocol = items("ftp://")
			i += 6
		} else {
			consumed = append(consumed, token)
			i++
			continue
		}

		if 0 < len(consumed) {
			text := &BaseNode{typ: NodeText, tokens: consumed}
			node.InsertBefore(node, text)
			consumed = make(items, 0, 256)
		}

		var url items
		j = i
		for ; j < length; j++ {
			token = tokens[j]
			if (isWhitespace(token) || itemLess == token) || (!isASCIIPunct(token) && !isASCIILetterNum(token)) {
				break
			}
			url = append(url, token)
		}
		if i == j { // 第一个字符就断开了
			url = append(url, token)
			text := &BaseNode{typ: NodeText, tokens: url}
			node.InsertBefore(node, text)
			i++
			continue
		}

		// 移动主循环下标
		i = j

		k = 0
		for ; k < len(url); k++ {
			token = url[k]
			if itemSlash == token {
				break
			}
		}
		domain := url[:k]
		if !t.isValidDomain(domain) {
			text := &BaseNode{typ: NodeText, tokens: url}
			node.InsertBefore(node, text)
			continue
		}

		var openParens, closeParens int
		// 最后一个字符如果是标点符号则剔掉
		path := url[k:]
		length := len(path)
		if 0 < length {
			var l int
			// 统计圆括号个数
			for l = 0; l < length; l++ {
				token = path[l]
				if itemOpenParen == token {
					openParens++
				} else if itemCloseParen == token {
					closeParens++
				}
			}

			trimmed := false
			lastToken := path[length-1]
			if itemCloseParen == lastToken {
				// 以 ) 结尾的话需要计算圆括号匹配
				unmatches := closeParens - openParens
				if 0 < unmatches {
					// 向前移动
					for l = length - 1; 0 < unmatches; l-- {
						token = path[l]
						if itemCloseParen != token {
							break
						}
						unmatches--
						i--
					}
					path = path[:l+1]
					trimmed = true
				} else { // 右圆括号 ) 数目小于等于左圆括号 ( 数目
					// 算作全匹配上了，不需要再处理结尾标点符号
					trimmed = true
				}
			} else if itemSemicolon == lastToken {
				// 检查 HTML 实体
				foundAmp := false
				// 向前检查 & 是否存在
				for l = length - 1; 0 <= l; l-- {
					token = path[l]
					if itemAmpersand == token {
						foundAmp = true
						break
					}
				}
				if foundAmp { // 如果 & 存在
					entity := path[l:length]
					if 3 <= len(entity) {
						// 检查截取的子串是否满足实体特征（&;中间需要是字母或数字）
						isEntity := true
						for j = 1; j < len(entity)-1; j++ {
							if !isASCIILetterNum(entity[j]) {
								isEntity = false
								break
							}
						}
						if isEntity {
							path = path[:l]
							trimmed = true
							i -= length - l
						}
					}
				}
			}

			// 如果之前的 ) 或者 ; 没有命中处理，则进行结尾的标点符号规则处理，即标点不计入链接，需要剔掉
			if !trimmed && isASCIIPunct(lastToken) && itemSlash != lastToken {
				path = path[:length-1]
				i--
			}
		} else {
			length = len(domain)
			lastToken := domain[length-1]
			if isASCIIPunct(lastToken) {
				domain = domain[:length-1]
				i--
			}
		}

		dest := protocol
		dest = append(dest, domain...)
		dest = append(dest, path...)
		var addr items
		if !www {
			addr = append(addr, protocol...)
		}
		addr = append(addr, domain...)
		addr = append(addr, path...)

		link := &BaseNode{typ: NodeLink, Destination: encodeDestination(dest)}
		text := &BaseNode{typ: NodeText, tokens: addr}
		link.AppendChild(link, text)
		node.InsertBefore(node, link)
	}

	if 0 < len(consumed) {
		text := &BaseNode{typ: NodeText, tokens: consumed}
		node.InsertBefore(node, text)
	}

	// 处理完后传入的文本节点 node 已经被拆分为多个节点，所以可以移除自身
	node.Unlink()
	return
}

var (
	// validDomainSuffix 用于列出所有认为合法的域名后缀。
	// TODO: 考虑提供接口支持开发者添加
	validDomainSuffix = [][]byte{[]byte("top"), []byte("com"), []byte("net"), []byte("org"), []byte("edu"), []byte("gov"), []byte("cn"), []byte("io")}
)

// isValidDomain 校验 GFM 规范自动链接规则中定义的合法域名。
// https://github.github.com/gfm/#valid-domain
func (t *Tree) isValidDomain(domain items) bool {
	segments := bytes.Split(domain, []byte("."))
	length := len(segments)
	if 2 > length { // 域名至少被 . 分隔为两部分，小于两部分的话不合法
		return false
	}

	var token byte
	for i := 0; i < length; i++ {
		segment := segments[i]
		segLen := len(segment)
		if 1 > segLen {
			continue
		}

		for j := 0; j < segLen; j++ {
			token = segment[j]
			if !isASCIILetterNumHyphen(token) {
				return false
			}
			if 2 < i && (i == length-2 || i == length-1) {
				// 最后两个部分不能包含 _
				if itemUnderscore == token {
					return false
				}
			}
		}

		if i == length-1 {
			validSuffix := false
			for j := 0; j < len(validDomainSuffix); j++ {
				if bytes.EqualFold(segment, validDomainSuffix[j]) {
					validSuffix = true
					break
				}
			}
			if !validSuffix {
				return false
			}
		}
	}
	return true
}

func (t *Tree) parseAutoEmailLink(ctx *InlineContext) (ret *BaseNode) {
	tokens := ctx.tokens[1:]
	var dest string
	var token byte
	length := len(tokens)
	passed := 0
	i := 0
	at := false
	for ; i < length; i++ {
		token = tokens[i]
		dest += string(token)
		passed++
		if '@' == token {
			at = true
			break
		}

		if !isASCIILetterNumHyphen(token) && !strings.Contains(".!#$%&'*+/=?^_`{|}~", string(token)) {
			return nil
		}
	}

	if 1 > i || !at {
		return nil
	}

	domainPart := tokens[i+1:]
	length = len(domainPart)
	i = 0
	closed := false
	for ; i < length; i++ {
		token = domainPart[i]
		passed++
		if itemGreater == token {
			closed = true
			break
		}
		dest += string(token)
		if !isASCIILetterNumHyphen(token) && itemDot != token {
			return nil
		}
		if 63 < i {
			return nil
		}
	}

	if 1 > i || !closed {
		return nil
	}

	ctx.pos += passed + 1
	link := &BaseNode{typ: NodeLink, Destination: append(items("mailto:"), dest...)}
	text := &BaseNode{typ: NodeText, tokens: toItems(dest)}
	link.AppendChild(ret, text)
	return link
}

func (t *Tree) parseAutolink(ctx *InlineContext) (ret *BaseNode) {
	schemed := false
	scheme := ""
	var dest items
	var token byte
	i := ctx.pos + 1
	for ; i < ctx.tokensLen && itemGreater != ctx.tokens[i]; i++ {
		token = ctx.tokens[i]
		if itemSpace == token {
			return nil
		}

		dest = append(dest, token)
		if !schemed {
			if itemColon != token {
				scheme += string(token)
			} else {
				schemed = true
			}
		}
	}
	if !schemed || 3 > len(scheme) || i == ctx.tokensLen {
		return nil
	}

	link := &BaseNode{typ: NodeLink, Destination: encodeDestination(dest)}
	if itemGreater != ctx.tokens[i] {
		return nil
	}

	ctx.pos = 1 + i
	text := &BaseNode{typ: NodeText, tokens: dest}
	link.AppendChild(ret, text)
	return link
}
