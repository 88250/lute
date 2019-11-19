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
	"unicode/utf8"
)

func (t *Tree) parseGFMAutoEmailLink(node *Node) {
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

func (t *Tree) parseGFMAutoLink(node *Node) {
	for child := node.firstChild; nil != child; {
		next := child.next
		if NodeText == child.typ {
			t.parseGFMAutoLink0(child)
		} else {
			t.parseGFMAutoLink(child) // 递归处理子节点
		}
		child = next
	}
}

var mailto = strToBytes("mailto:")

func (t *Tree) parseGFMAutoEmailLink0(node *Node) {
	tokens := node.tokens
	if 0 >= bytes.IndexByte(tokens, '@') {
		return
	}

	var i, j, k, atIndex int
	var token byte
	length := len(tokens)

	// 按空白分隔成多组并进行处理
loopPart:
	for i < length {
		var group []byte
		atIndex = 0
		j = i

		// 积攒组直到遇到空白符
		for ; j < length; j++ {
			token = tokens[j]
			if !isWhitespace(token) {
				group = append(group, tokens[j])
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
			text := &Node{typ: NodeText, tokens: []byte{tokens[j]}}
			node.InsertBefore(text)
			i++
			continue
		}

		// 移动主循环下标
		i = j

		if 0 >= atIndex {
			text := &Node{typ: NodeText, tokens: group}
			node.InsertBefore(text)
			continue
		}

		// 至此说明这一组中包含了 @，可尝试进行邮件地址解析

		k = 0
		for ; k < atIndex; k++ {
			token = group[k]
			if !t.isValidEmailSegment1(token) {
				text := &Node{typ: NodeText, tokens: group}
				node.InsertBefore(text)
				continue loopPart
			}
		}

		k++ // 跳过 @ 检查后面的部分
		var item byte
		for ; k < len(group); k++ {
			item = group[k]
			token = group[k]
			if !t.isValidEmailSegment2(token) {
				text := &Node{typ: NodeText, tokens: group}
				node.InsertBefore(text)
				continue loopPart
			}
		}

		if itemDot == token {
			// 如果以 . 结尾则剔除该 .
			lastIndex := len(group) - 1
			group = group[:lastIndex]
			link := t.newLink(NodeLink, group, append(mailto, group...), nil, 2)
			node.InsertBefore(link)
			// . 作为文本节点插入
			node.InsertBefore(&Node{typ: NodeText, tokens: []byte{item}})
		} else if itemHyphen == token || itemUnderscore == token {
			// 如果以 - 或者 _ 结尾则整个串都不能算作邮件链接
			text := &Node{typ: NodeText, tokens: group}
			node.InsertBefore(text)
			continue loopPart
		} else {
			// 以字母或者数字结尾
			link := &Node{typ: NodeLink, linkType: 2}
			text := &Node{typ: NodeLinkText, tokens: group}
			link.AppendChild(text)
			dest := &Node{typ: NodeLinkDest, tokens: append(mailto, group...)}
			link.AppendChild(dest)
			node.InsertBefore(link)
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

var (
	httpProto = strToBytes("http://")

	// validAutoLinkDomainSuffix 作为 GFM 自动连接解析时校验域名后缀用。
	validAutoLinkDomainSuffix = [][]byte{strToBytes("top"), strToBytes("com"), strToBytes("net"), strToBytes("org"), strToBytes("edu"), strToBytes("gov"),
		strToBytes("cn"), strToBytes("io"), strToBytes("me"), strToBytes("biz"), strToBytes("co"), strToBytes("live"), strToBytes("pro"), strToBytes("xyz"),
		strToBytes("win"), strToBytes("club"), strToBytes("tv"), strToBytes("wiki"), strToBytes("site"), strToBytes("tech"), strToBytes("space"), strToBytes("cc"),
		strToBytes("name"), strToBytes("social"), strToBytes("band"), strToBytes("pub"), strToBytes("info")}
)

// AddAutoLinkDomainSuffix 添加自动链接解析域名后缀 suffix。
func (lute *Lute) AddAutoLinkDomainSuffix(suffix string) {
	validAutoLinkDomainSuffix = append(validAutoLinkDomainSuffix, strToBytes(suffix))
}

func (t *Tree) parseGFMAutoLink0(node *Node) {
	tokens := node.tokens
	var i, j, k int
	length := len(tokens)
	minLinkLen := 10 // 太短的情况肯定不可能有链接，最短的情况是 www.xxx.xx
	if minLinkLen > length {
		return
	}

	var token byte
	var consumed = make([]byte, 0, 256)
	var tmp = make([]byte, 0, 16)
	www := false
	for i < length {
		token = tokens[i]
		var protocol []byte
		// 检查前缀
		tmp = tokens[i:]
		tmpLen := len(tmp)
		if 10 <= tmpLen /* www.xxx.xx */ && 'w' == tmp[0] && 'w' == tmp[1] && 'w' == tmp[2] && '.' == tmp[3] {
			protocol = httpProto
			www = true
		} else if 13 <= tmpLen /* http://xxx.xx */ && 'h' == tmp[0] && 't' == tmp[1] && 't' == tmp[2] && 'p' == tmp[3] && ':' == tmp[4] && '/' == tmp[5] && '/' == tmp[6] {
			protocol = tmp[0:7]
			i += 7
		} else if 14 <= tmpLen /* https://xxx.xx */ && 'h' == tmp[0] && 't' == tmp[1] && 't' == tmp[2] && 'p' == tmp[3] && 's' == tmp[4] && ':' == tmp[5] && '/' == tmp[6] && '/' == tmp[7] {
			protocol = tmp[0:8]
			i += 8
		} else if 12 <= tmpLen /* ftp://xxx.xx */ && 'f' == tmp[0] && 't' == tmp[1] && 'p' == tmp[2] && ':' == tmp[3] && '/' == tmp[4] && '/' == tmp[5] {
			protocol = tmp[0:6]
			i += 6
		} else {
			consumed = append(consumed, token)
			if length-i < minLinkLen && 0 < length-i {
				// 剩余字符不足，已经不可能形成链接了
				consumed = append(consumed, tokens[i+1:]...)
				node.InsertBefore(&Node{typ: NodeText, tokens: consumed})
				node.Unlink()
				return
			}
			i++
			continue
		}

		if 0 < len(consumed) {
			text := &Node{typ: NodeText, tokens: consumed}
			node.InsertBefore(text)
			consumed = make([]byte, 0, 256)
		}

		var url []byte
		j = i
		for ; j < length; j++ {
			token = tokens[j]
			if (isWhitespace(token) || itemLess == token) || (!isASCIIPunct(token) && !isASCIILetterNum(token)) {
				break
			}
			url = append(url, token)
		}
		if i == j { // 第一个字符就断开了
			if utf8.RuneSelf <= token {
				if !www {
					url = append(url, protocol...)
				}
				for ; i < length; i++ {
					token = tokens[i]
					if utf8.RuneSelf > token {
						break
					}
					url = append(url, token)
				}
			} else {
				url = append(url, token)
				i++
			}
			text := &Node{typ: NodeText, tokens: url}
			node.InsertBefore(text)
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
			var part []byte
			if www {
				part = url
			} else {
				part = append(protocol, url...)
			}
			text := &Node{typ: NodeText, tokens: part}
			node.InsertBefore(text)
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
		var addr []byte
		if !www {
			addr = append(addr, protocol...)
		}
		addr = append(addr, domain...)
		addr = append(addr, path...)

		link := t.newLink(NodeLink, addr, encodeDestination(dest), nil, 2)
		node.InsertBefore(link)
	}

	if 0 < len(consumed) {
		text := &Node{typ: NodeText, tokens: consumed}
		node.InsertBefore(text)
	}

	// 处理完后传入的文本节点 node 已经被拆分为多个节点，所以可以移除自身
	node.Unlink()
	return
}

// isValidDomain 校验 GFM 规范自动链接规则中定义的合法域名。
// https://github.github.com/gfm/#valid-domain
func (t *Tree) isValidDomain(domain []byte) bool {
	segments := split(domain, '.')
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
			for j := 0; j < len(validAutoLinkDomainSuffix); j++ {
				if bytes.Equal(segment, validAutoLinkDomainSuffix[j]) {
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

var markers = strToBytes(".!#$%&'*+/=?^_`{|}~")

func (t *Tree) parseAutoEmailLink(ctx *InlineContext) (ret *Node) {
	tokens := ctx.tokens[1:]
	var dest []byte
	var token byte
	length := len(tokens)
	passed := 0
	i := 0
	at := false
	for ; i < length; i++ {
		token = tokens[i]
		dest = append(dest, tokens[i])
		passed++
		if '@' == token {
			at = true
			break
		}

		if !isASCIILetterNumHyphen(token) && !bytes.Contains(markers, []byte{token}) {
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
		dest = append(dest, domainPart[i])
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
	return t.newLink(NodeLink, dest, append(mailto, dest...), nil, 2)
}

func (t *Tree) newLink(typ nodeType, text, dest, title []byte, linkType int) (ret *Node) {
	ret = &Node{typ: typ, linkType: linkType}
	if NodeImage == typ {
		ret.AppendChild(&Node{typ: NodeBang})
	}
	ret.AppendChild(&Node{typ: NodeOpenBracket})
	ret.AppendChild(&Node{typ: NodeLinkText, tokens: text})
	ret.AppendChild(&Node{typ: NodeCloseBracket})
	ret.AppendChild(&Node{typ: NodeOpenParen})
	ret.AppendChild(&Node{typ: NodeLinkDest, tokens: dest})
	if nil != title {
		ret.AppendChild(&Node{typ: NodeLinkTitle, tokens: title})
	}
	ret.AppendChild(&Node{typ: NodeCloseParen})
	return
}

func (t *Tree) parseAutolink(ctx *InlineContext) (ret *Node) {
	schemed := false
	scheme := ""
	var dest []byte
	var token byte
	i := ctx.pos + 1
	for ; i < ctx.tokensLen && itemGreater != ctx.tokens[i]; i++ {
		token = ctx.tokens[i]
		if itemSpace == token {
			return nil
		}

		dest = append(dest, ctx.tokens[i])
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

	if itemGreater != ctx.tokens[i] {
		return nil
	}
	ctx.pos = 1 + i
	return t.newLink(NodeLink, dest, encodeDestination(dest), nil, 2)
}
