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
		if NodeText == child.typ && nil != child.parent &&
			NodeLink != child.parent.typ /* 不处理链接 label */ {
			t.parseGFMAutoLink0(child)
		} else {
			t.parseGFMAutoLink(child) // 递归处理子节点
		}
		child = next
	}
}

var mailto = strToItems("mailto:")
var at = strToItems("@")

func (t *Tree) parseGFMAutoEmailLink0(node *Node) {
	tokens := node.tokens
	if 0 >= index(tokens, at) {
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
			token = tokens[j].term
			if !isWhitespace(token) {
				group = append(group, &item{term: token})
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
			text := &Node{typ: NodeText, tokens: items{{term: token}}}
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
			token = group[k].term
			if !t.isValidEmailSegment1(token) {
				text := &Node{typ: NodeText, tokens: group}
				node.InsertBefore(text)
				continue loopPart
			}
		}

		k++ // 跳过 @ 检查后面的部分
		for ; k < len(group); k++ {
			token = group[k].term
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
			link := &Node{typ: NodeLink, destination: append(mailto, group...)}
			text := &Node{typ: NodeText, tokens: group}
			link.AppendChild(text)
			node.InsertBefore(link)
			// . 作为文本节点插入
			text = &Node{typ: NodeText, tokens: items{{term: itemDot}}}
			node.InsertBefore(text)
		} else if itemHyphen == token || itemUnderscore == token {
			// 如果以 - 或者 _ 结尾则整个串都不能算作邮件链接
			text := &Node{typ: NodeText, tokens: group}
			node.InsertBefore(text)
			continue loopPart
		} else {
			// 以字母或者数字结尾
			link := &Node{typ: NodeLink, destination: append(mailto, group...)}
			text := &Node{typ: NodeText, tokens: group}
			link.AppendChild(text)
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
	httpProto  = strToItems("http://")
	httpsProto = strToItems("https://")
	ftpProto   = strToItems("ftp://")

	// validAutoLinkDomainSuffix 作为 GFM 自动连接解析时校验域名后缀用。
	validAutoLinkDomainSuffix = [][]byte{[]byte("top"), []byte("com"), []byte("net"), []byte("org"), []byte("edu"), []byte("gov"),
		[]byte("cn"), []byte("io"), []byte("me"), []byte("biz"), []byte("co"), []byte("live"), []byte("pro"), []byte("xyz"),
		[]byte("win"), []byte("club"), []byte("tv"), []byte("wiki"), []byte("site"), []byte("tech"), []byte("space"), []byte("cc"),
		[]byte("name"), []byte("social"), []byte("band"), []byte("pub"), []byte("info")}
)

func (t *Tree) parseGFMAutoLink0(node *Node) {
	tokens := node.tokens
	var i, j, k int
	length := len(tokens)
	if 8 > length { // 太短的情况肯定不可能有链接
		return
	}

	var token *item
	var consumed = make(items, 0, 256)
	var tmp = make(items, 0, 16)
	www := false
	for i < length {
		token = tokens[i]
		var protocol items

		// 检查前缀
		tmp = tokens[i:]
		tmpLen := len(tmp)
		if 8 <= tmpLen /* www.x.xx */ && 'w' == tmp[0].term && 'w' == tmp[1].term && 'w' == tmp[2].term && '.' == tmp[3].term {
			protocol = httpProto
			www = true
		} else if 11 <= tmpLen /* http://x.xx */ && 'h' == tmp[0].term && 't' == tmp[1].term && 't' == tmp[2].term && 'p' == tmp[3].term && ':' == tmp[4].term && '/' == tmp[5].term && '/' == tmp[6].term {
			protocol = httpProto
			i += 7
		} else if 12 <= tmpLen /* https://x.xx */ && 'h' == tmp[0].term && 't' == tmp[1].term && 't' == tmp[2].term && 'p' == tmp[3].term && 's' == tmp[4].term && ':' == tmp[5].term && '/' == tmp[6].term && '/' == tmp[7].term {
			protocol = httpsProto
			i += 8
		} else if 10 <= tmpLen /* ftp://x.xx */ && 'f' == tmp[0].term && 't' == tmp[1].term && 'p' == tmp[2].term && ':' == tmp[3].term && '/' == tmp[4].term && '/' == tmp[5].term {
			protocol = ftpProto
			i += 6
		} else {
			consumed = append(consumed, token)
			i++
			continue
		}

		if 0 < len(consumed) {
			text := &Node{typ: NodeText, tokens: consumed}
			node.InsertBefore(text)
			consumed = make(items, 0, 256)
		}

		var url items
		j = i
		for ; j < length; j++ {
			token = tokens[j]
			if (isWhitespace(token.term) || itemLess == token.term) || (!isASCIIPunct(token.term) && !isASCIILetterNum(token.term)) {
				break
			}
			url = append(url, token)
		}
		if i == j { // 第一个字符就断开了
			url = append(url, token)
			text := &Node{typ: NodeText, tokens: url}
			node.InsertBefore(text)
			i++
			continue
		}

		// 移动主循环下标
		i = j

		k = 0
		for ; k < len(url); k++ {
			token = url[k]
			if itemSlash == token.term {
				break
			}
		}
		domain := url[:k]
		domainBytes := itemsToBytes(domain) // TODO 性能优化
		if !t.isValidDomain(domainBytes) {
			text := &Node{typ: NodeText, tokens: append(protocol, url...)}
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
				if itemOpenParen == token.term {
					openParens++
				} else if itemCloseParen == token.term {
					closeParens++
				}
			}

			trimmed := false
			lastToken := path[length-1]
			if itemCloseParen == lastToken.term {
				// 以 ) 结尾的话需要计算圆括号匹配
				unmatches := closeParens - openParens
				if 0 < unmatches {
					// 向前移动
					for l = length - 1; 0 < unmatches; l-- {
						token = path[l]
						if itemCloseParen != token.term {
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
			} else if itemSemicolon == lastToken.term {
				// 检查 HTML 实体
				foundAmp := false
				// 向前检查 & 是否存在
				for l = length - 1; 0 <= l; l-- {
					token = path[l]
					if itemAmpersand == token.term {
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
							if !isASCIILetterNum(entity[j].term) {
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
			if !trimmed && isASCIIPunct(lastToken.term) && itemSlash != lastToken.term {
				path = path[:length-1]
				i--
			}
		} else {
			length = len(domain)
			lastToken := domain[length-1]
			if isASCIIPunct(lastToken.term) {
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

		link := &Node{typ: NodeLink, destination: strToItems(encodeDestination(dest))}
		text := &Node{typ: NodeText, tokens: addr}
		link.AppendChild(text)
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

// invalidAutoLinkDomain 指定了 GFM 自动链接解析时跳过的域名。
var invalidAutoLinkDomain [][]byte

// isValidDomain 校验 GFM 规范自动链接规则中定义的合法域名。
// https://github.github.com/gfm/#valid-domain
func (t *Tree) isValidDomain(domain []byte) bool {
	segments := bytes.Split(domain, []byte("."))
	length := len(segments)
	if 2 > length { // 域名至少被 . 分隔为两部分，小于两部分的话不合法
		return false
	}

	for i := 0; i < len(invalidAutoLinkDomain); i++ {
		if bytes.EqualFold(domain, invalidAutoLinkDomain[i]) {
			return false
		}
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
				if bytes.EqualFold(segment, validAutoLinkDomainSuffix[j]) {
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

var markers = toBytes(".!#$%&'*+/=?^_`{|}~")

func (t *Tree) parseAutoEmailLink(ctx *InlineContext) (ret *Node) {
	tokens := ctx.tokens[1:]
	var dest items
	var token byte
	length := len(tokens)
	passed := 0
	i := 0
	at := false
	for ; i < length; i++ {
		token = tokens[i].term
		dest = append(dest, &item{term: token})
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
		token = domainPart[i].term
		passed++
		if itemGreater == token {
			closed = true
			break
		}
		dest = append(dest, &item{term: token})
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
	link := &Node{typ: NodeLink, destination: append(mailto, dest...)}
	text := &Node{typ: NodeText, tokens: dest}
	link.AppendChild(text)
	return link
}

func (t *Tree) parseAutolink(ctx *InlineContext) (ret *Node) {
	schemed := false
	scheme := ""
	var dest items
	var token byte
	i := ctx.pos + 1
	for ; i < ctx.tokensLen && itemGreater != ctx.tokens[i].term; i++ {
		token = ctx.tokens[i].term
		if itemSpace == token {
			return nil
		}

		dest = append(dest, &item{term: token})
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

	link := &Node{typ: NodeLink, destination: strToItems(encodeDestination(dest))}
	if itemGreater != ctx.tokens[i].term {
		return nil
	}

	ctx.pos = 1 + i
	text := &Node{typ: NodeText, tokens: dest}
	link.AppendChild(text)
	return link
}
