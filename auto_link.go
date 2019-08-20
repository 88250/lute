// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lute

import (
	"bytes"
	"strings"
)

// parseGfmAutoEmailLink 解析 node 文本节点中的 tokens，如果有邮件地址则生成链接节点并插入到 node 之前或者之后。
func (t *Tree) parseGfmAutoEmailLink(node Node) {
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
		group := items{}
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
			text := &Text{tokens: items{token}}
			node.InsertBefore(node, text)
			i++ // 继续下一个字符
			continue
		}

		// 移动主循环下标
		i = j

		if 0 >= atIndex {
			text := &Text{tokens: group}
			node.InsertBefore(node, text)
			continue
		}

		// 至此说明这一组中包含了 @，可尝试进行邮件地址解析

		k = 0
		for ; k < atIndex; k++ {
			token = group[k]
			if !t.isValidEmailSegment1(token) {
				text := &Text{tokens: group}
				node.InsertBefore(node, text)
				continue loopPart
			}
		}

		k++ // 跳过 @ 检查后面的部分
		length := len(group)
		for ; k < length; k++ {
			token = group[k]
			if !t.isValidEmailSegment2(token) {
				text := &Text{tokens: group}
				node.InsertBefore(node, text)
				continue loopPart
			}
		}

		if itemDot == token {
			// 如果以 . 结尾则剔除该 .
			lastIndex := len(group) - 1
			group = group[:lastIndex]
			link := &Link{&BaseNode{typ: NodeLink}, "mailto:" + fromItems(group), ""}
			link.AppendChild(link, &Text{tokens: group})
			node.InsertBefore(node, link)
			// . 作为文本节点插入
			text := &Text{tokens: items{itemDot}}
			node.InsertBefore(node, text)
		} else if itemHyphen == token || itemUnderscore == token {
			// 如果以 - 或者 _ 结尾则整个串都不能算作邮件链接
			text := &Text{tokens: group}
			node.InsertBefore(node, text)
			continue loopPart
		} else {
			// 以字母或者数字结尾
			link := &Link{&BaseNode{typ: NodeLink}, "mailto:" + fromItems(group), ""}
			link.AppendChild(link, &Text{tokens: group})
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

func (t *Tree) parseGfmAutoLink(tokens items, protocol string) (ret Node) {
	tokens = tokens[t.context.pos:]
	index := bytes.Index(tokens, []byte(protocol))
	if 0 > index {
		return nil
	}

	var i int
	if 0 < index {
		// 检查是否有潜在的标记，有的话处理标记优先
		for ; i < index; i++ {
			if t.isMarker(tokens[i]) {
				break
			}
		}
		// 将标记出现之前的部分构造为文本节点
		t.context.pos += i
		return &Text{tokens: tokens[:i]}
	}

	length := len(tokens)
	var token byte
	for ; i < length; i++ {
		token = tokens[i]
		// 链接以空白或者 < 截断
		if isWhitespace(token) || itemLess == token {
			break
		}
	}

	www := "www." == protocol

	url := tokens[:i]
	length = len(url)
	var j int
	if !www {
		j = len(protocol)
	}
	for ; j < length; j++ {
		token = url[j]
		if itemSlash == token {
			break
		}
	}
	domain := url[:j]
	if !www {
		domain = domain[len(protocol):]
	}

	if !t.isValidDomain(domain) {
		t.context.pos += i
		return &Text{tokens: url}
	}

	var openParens, closeParens int
	// 最后一个字符如果是标点符号则剔掉
	path := url[j:]
	length = len(path)
	if 0 < length {
		var k int
		// 统计圆括号个数
		for k = 0; k < length; k++ {
			token = path[k]
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
				for k = length - 1; 0 < unmatches; k-- {
					token = path[k]
					if itemCloseParen != token {
						break
					}
					unmatches--
					i--
				}
				path = path[:k+1]
				trimmed = true
			} else { // 右圆括号 ) 数目小于等于左圆括号 ( 数目
				// 算作全匹配上了，不需要再处理结尾标点符号
				trimmed = true
			}
		} else if itemSemicolon == lastToken {
			// 检查 HTML 实体
			foundAmp := false
			// 向前检查 & 是否存在
			for k = length - 1; 0 <= k; k-- {
				token = path[k]
				if itemAmpersand == token {
					foundAmp = true
					break
				}
			}
			if foundAmp { // 如果 & 存在
				entity := path[k:length]
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
						path = path[:k]
						trimmed = true
						i -= length - k
					}
				}
			}
		}

		// 如果之前的 ) 或者 ; 没有命中处理，则进行结尾的标点符号规则处理，即标点不计入链接，需要剔掉
		if !trimmed && isASCIIPunct(lastToken) {
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

	var dest, domainPath items
	if www {
		dest = items("http://")
		domainPath = append(domain, path...)
		dest = append(dest, domainPath...)
	} else {
		dest = items(protocol)
		domain = append(dest, domain...)
		domainPath = append(domain, path...)
		dest = domainPath
	}

	ret = &Link{&BaseNode{typ: NodeLink}, encodeDestination(fromItems(dest)), ""}
	ret.AppendChild(ret, &Text{tokens: domainPath})
	t.context.pos += i
	return
}

var (
	// validDomainSuffix 用于列出所有认为合法的域名后缀，不够的话往里加就行。
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

func (t *Tree) parseAutoEmailLink(tokens items) (ret Node) {
	tokens = tokens[1:]
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

	t.context.pos += passed + 1
	ret = &Link{&BaseNode{typ: NodeLink}, "mailto:" + dest, ""}
	ret.AppendChild(ret, &Text{tokens: toItems(dest)})

	return
}

func (t *Tree) parseAutolink(tokens items) (ret Node) {
	schemed := false
	scheme := ""
	dest := ""
	var token byte
	i := t.context.pos + 1
	for ; i < len(tokens) && itemGreater != tokens[i]; i++ {
		token = tokens[i]
		if itemSpace == token {
			return nil
		}

		dest += string(token)
		if !schemed {
			if itemColon != token {
				scheme += string(token)
			} else {
				schemed = true
			}
		}
	}
	if !schemed || 3 > len(scheme) {
		return nil
	}

	ret = &Link{&BaseNode{typ: NodeLink}, encodeDestination(dest), ""}
	if itemGreater != tokens[i] {
		return nil
	}

	t.context.pos = 1 + i
	ret.AppendChild(ret, &Text{tokens: toItems(dest)})

	return
}
