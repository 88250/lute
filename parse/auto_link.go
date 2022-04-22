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
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

func (t *Tree) parseGFMAutoEmailLink(node *ast.Node) {
	for child := node.FirstChild; nil != child; {
		next := child.Next
		if ast.NodeText == child.Type && nil != child.Parent &&
			ast.NodeLink != child.Parent.Type /* 不处理链接 label */ {
			t.parseGFMAutoEmailLink0(child)
		} else {
			t.parseGFMAutoEmailLink(child) // 递归处理子节点
		}
		child = next
	}
}

func (t *Tree) parseGFMAutoLink(node *ast.Node) {
	for child := node.FirstChild; nil != child; {
		next := child.Next
		if ast.NodeText == child.Type {
			t.parseGFMAutoLink0(child)
		} else {
			t.parseGFMAutoLink(child) // 递归处理子节点
		}
		child = next
	}
}

var mailto = util.StrToBytes("mailto:")

func (t *Tree) parseGFMAutoEmailLink0(node *ast.Node) {
	tokens := node.Tokens
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
			if !lex.IsWhitespace(token) {
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
			t.addPreviousText(node, []byte{tokens[j]})
			i++
			continue
		}

		// 移动主循环下标
		i = j

		if 0 >= atIndex {
			t.addPreviousText(node, group)
			continue
		}

		// 至此说明这一组中包含了 @，可尝试进行邮件地址解析

		k = 0
		for ; k < atIndex; k++ {
			token = group[k]
			if !t.isValidEmailSegment1(token) {
				t.addPreviousText(node, group)
				continue loopPart
			}
		}

		k++ // 跳过 @ 检查后面的部分
		var item byte
		for ; k < len(group); k++ {
			item = group[k]
			token = group[k]
			if !t.isValidEmailSegment2(token) {
				t.addPreviousText(node, group)
				continue loopPart
			}
		}

		if lex.ItemDot == token {
			// 如果以 . 结尾则剔除该 .
			lastIndex := len(group) - 1
			group = group[:lastIndex]
			link := t.newLink(ast.NodeLink, group, append(mailto, group...), nil, 2)
			node.InsertBefore(link)
			// . 作为文本节点插入
			t.addPreviousText(node, []byte{item})
		} else if lex.ItemHyphen == token || lex.ItemUnderscore == token {
			// 如果以 - 或者 _ 结尾则整个串都不能算作邮件链接
			t.addPreviousText(node, group)
			continue loopPart
		} else {
			// 以字母或者数字结尾
			link := &ast.Node{Type: ast.NodeLink, LinkType: 2}
			link.AppendChild(&ast.Node{Type: ast.NodeLinkText, Tokens: group})
			link.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: append(mailto, group...)})
			node.InsertBefore(link)
		}
	}

	// 处理完后传入的文本节点 node 已经被拆分为多个节点，所以可以移除自身
	node.Unlink()
	return
}

func (t *Tree) isValidEmailSegment1(token byte) bool {
	return lex.IsASCIILetterNumHyphen(token) || lex.ItemDot == token || lex.ItemPlus == token || lex.ItemUnderscore == token
}

func (t *Tree) isValidEmailSegment2(token byte) bool {
	return lex.IsASCIILetterNumHyphen(token) || lex.ItemDot == token || lex.ItemUnderscore == token
}

var (
	httpProto = util.StrToBytes("http://")

	// validAutoLinkDomainSuffix 作为 GFM 自动连接解析时校验域名后缀用。
	validAutoLinkDomainSuffix = [][]byte{util.StrToBytes("top"), util.StrToBytes("com"), util.StrToBytes("net"), util.StrToBytes("org"), util.StrToBytes("edu"), util.StrToBytes("gov"),
		util.StrToBytes("cn"), util.StrToBytes("io"), util.StrToBytes("me"), util.StrToBytes("biz"), util.StrToBytes("co"), util.StrToBytes("live"), util.StrToBytes("pro"), util.StrToBytes("xyz"),
		util.StrToBytes("win"), util.StrToBytes("club"), util.StrToBytes("tv"), util.StrToBytes("wiki"), util.StrToBytes("site"), util.StrToBytes("tech"), util.StrToBytes("space"), util.StrToBytes("cc"),
		util.StrToBytes("name"), util.StrToBytes("social"), util.StrToBytes("band"), util.StrToBytes("pub"), util.StrToBytes("info"), util.StrToBytes("app"), util.StrToBytes("md"), util.StrToBytes("edu"),
		util.StrToBytes("hk"), util.StrToBytes("so"), util.StrToBytes("vip"), util.StrToBytes("ai"), util.StrToBytes("ink"), util.StrToBytes("mobi"), util.StrToBytes("pro")}
)

// AddAutoLinkDomainSuffix 添加自动链接解析域名后缀 suffix。
func AddAutoLinkDomainSuffix(suffix string) {
	validAutoLinkDomainSuffix = append(validAutoLinkDomainSuffix, util.StrToBytes(suffix))
}

func (t *Tree) parseGFMAutoLink0(node *ast.Node) {
	tokens := node.Tokens
	length := len(tokens)
	minLinkLen := 10 // 太短的情况肯定不可能有链接，最短的情况是 www.xxx.xx
	if minLinkLen > length {
		return
	}

	var i, j, k int
	var textStart, textEnd int
	var token byte
	www := false
	needUnlink := false
	for i < length {
		token = tokens[i]
		var protocol []byte
		// 检查前缀
		tmpLen := length - i
		if 10 <= tmpLen /* www.xxx.xx */ && 'w' == tokens[i] && 'w' == tokens[i+1] && 'w' == tokens[i+2] && '.' == tokens[i+3] {
			protocol = httpProto
			www = true
		} else if 13 <= tmpLen /* http://xxx.xx */ && 'h' == tokens[i] && 't' == tokens[i+1] && 't' == tokens[i+2] && 'p' == tokens[i+3] && ':' == tokens[i+4] && '/' == tokens[i+5] && '/' == tokens[i+6] {
			protocol = tokens[i : i+7]
			i += 7
		} else if 14 <= tmpLen /* https://xxx.xx */ && 'h' == tokens[i] && 't' == tokens[i+1] && 't' == tokens[i+2] && 'p' == tokens[i+3] && 's' == tokens[i+4] && ':' == tokens[i+5] && '/' == tokens[i+6] && '/' == tokens[i+7] {
			protocol = tokens[i : i+8]
			i += 8
		} else if 12 <= tmpLen /* ftp://xxx.xx */ && 'f' == tokens[i] && 't' == tokens[i+1] && 'p' == tokens[i+2] && ':' == tokens[i+3] && '/' == tokens[i+4] && '/' == tokens[i+5] {
			protocol = tokens[i : i+6]
			i += 6
		} else {
			textEnd++
			if length-i < minLinkLen { // 剩余字符不足，已经不可能形成链接了
				if needUnlink {
					if textStart < textEnd {
						t.addPreviousText(node, tokens[textStart:])
					} else {
						t.addPreviousText(node, tokens[textEnd:])
					}
					node.Unlink()
				}
				return
			}
			i++
			continue
		}

		if textStart < textEnd {
			t.addPreviousText(node, tokens[textStart:textEnd])
			needUnlink = true
			textStart = textEnd
		}

		var url []byte
		j = i
		for ; j < length; j++ {
			token = tokens[j]
			if (lex.IsWhitespace(token) || lex.ItemLess == token) || (!lex.IsASCIIPunct(token) && !lex.IsASCIILetterNum(token)) {
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

			if nil != node.Previous {
				node.Previous.Tokens = append(node.Previous.Tokens, url...)
			}

			textStart = i
			textEnd = i
			continue
		}

		// 移动主循环下标
		i = j

		k = 0
		for ; k < len(url); k++ {
			token = url[k]
			if lex.ItemSlash == token {
				break
			}
		}
		domain := url[:k]
		var port []byte
		if idx := bytes.Index(domain, []byte(":")); 0 < idx {
			port = domain[idx:]
			domain = domain[:idx]
		}

		if !t.isValidDomain(domain) {
			t.addPreviousText(node, tokens[textStart:i])
			needUnlink = true
			textStart = i
			textEnd = i
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
				if lex.ItemOpenParen == token {
					openParens++
				} else if lex.ItemCloseParen == token {
					closeParens++
				}
			}

			trimmed := false
			lastToken := path[length-1]
			if lex.ItemCloseParen == lastToken {
				// 以 ) 结尾的话需要计算圆括号匹配
				unmatches := closeParens - openParens
				if 0 < unmatches {
					// 向前移动
					for l = length - 1; 0 < unmatches; l-- {
						token = path[l]
						if lex.ItemCloseParen != token {
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
			} else if lex.ItemSemicolon == lastToken {
				// 检查 HTML 实体
				foundAmp := false
				// 向前检查 & 是否存在
				for l = length - 1; 0 <= l; l-- {
					token = path[l]
					if lex.ItemAmpersand == token {
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
							if !lex.IsASCIILetterNum(entity[j]) {
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
			if !trimmed && lex.IsASCIIPunct(lastToken) && lex.ItemSlash != lastToken &&
				'}' != lastToken && '{' != lastToken /* 自动链接解析结尾 } 问题 https://github.com/88250/lute/issues/4 */ {
				path = path[:length-1]
				i--
			}
		} else {
			length = len(domain)
			lastToken := domain[length-1]
			if lex.IsASCIIPunct(lastToken) {
				domain = domain[:length-1]
				i--
			}
		}

		dest := append(protocol, domain...)
		dest = append(dest, port...)
		dest = append(dest, path...)
		var addr []byte
		if !www {
			addr = append(addr, protocol...)
		}
		addr = append(addr, domain...)
		addr = append(addr, path...)
		linkText := addr
		if bytes.HasPrefix(linkText, []byte("https://github.com/")) && bytes.Contains(linkText, []byte("/issues/")) {
			// 优化 GitHub Issues 自动链接文本 https://github.com/88250/lute/issues/161
			repo := linkText[len("https://github.com/"):]
			repo = repo[:bytes.Index(repo, []byte("/issues/"))]
			num := bytes.Split(linkText, []byte("/issues/"))[1]
			num = bytes.Split(num, []byte("?"))[0]
			if 0 < len(num) {
				isDigit := true
				for _, d := range num {
					if !lex.IsDigit(d) {
						isDigit = false
						break
					}
				}
				if isDigit {
					linkText = []byte("Issue #" + string(num) + " · " + string(repo))
				}
			}

		}

		link := t.newLink(ast.NodeLink, linkText, html.EncodeDestination(dest), nil, 2)
		node.InsertBefore(link)
		needUnlink = true

		textStart = i
		textEnd = i
	}

	if textStart < textEnd {
		t.addPreviousText(node, tokens[textStart:textEnd])
		needUnlink = true
	}
	if needUnlink {
		node.Unlink()
	}
	return
}

// isValidDomain 校验 GFM 规范自动链接规则中定义的合法域名。
// https://github.github.com/gfm/#valid-domain
func (t *Tree) isValidDomain(domain []byte) bool {
	segments := lex.Split(domain, '.')
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
			if !lex.IsASCIILetterNumHyphen(token) {
				return false
			}
			if 2 < i && (i == length-2 || i == length-1) {
				// 最后两个部分不能包含 _
				if lex.ItemUnderscore == token {
					return false
				}
			}
		}

		if i == length-1 {
			validSuffix := false
			suffixIsDigit := true // 校验后缀是否全为数字
			for _, b := range segment {
				if !lex.IsDigit(b) {
					suffixIsDigit = false
					break
				}
			}
			if !suffixIsDigit { // 如果后缀不是数字的话检查是否在后缀可用名单中
				for j := 0; j < len(validAutoLinkDomainSuffix); j++ {
					if bytes.Equal(segment, validAutoLinkDomainSuffix[j]) {
						validSuffix = true
						break
					}
				}
			} else { // 后缀全为数字的话可能是 IPv4 地址
				validSuffix = true
			}
			if !validSuffix {
				return false
			}
		}
	}
	return true
}

var markers = util.StrToBytes(".!#$%&'*+/=?^_`{|}~")

func (t *Tree) parseAutoEmailLink(ctx *InlineContext) (ret *ast.Node) {
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

		if !lex.IsASCIILetterNumHyphen(token) && !bytes.Contains(markers, []byte{token}) {
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
		if lex.ItemGreater == token {
			closed = true
			break
		}
		dest = append(dest, domainPart[i])
		if !lex.IsASCIILetterNumHyphen(token) && lex.ItemDot != token {
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
	return t.newLink(ast.NodeLink, dest, append(mailto, dest...), nil, 2)
}

func (t *Tree) newLink(typ ast.NodeType, text, dest, title []byte, linkType int) (ret *ast.Node) {
	appendCaret := t.Context.ParseOption.ProtyleWYSIWYG && bytes.HasSuffix(text, util.CaretTokens) && bytes.HasSuffix(dest, []byte("%E2%80%B8"))
	if appendCaret {
		text = bytes.ReplaceAll(text, util.CaretTokens, nil)
		dest = bytes.ReplaceAll(dest, []byte("%E2%80%B8"), nil)
	}

	ret = &ast.Node{Type: typ, LinkType: linkType}
	if ast.NodeImage == typ {
		ret.AppendChild(&ast.Node{Type: ast.NodeBang})
	}
	ret.AppendChild(&ast.Node{Type: ast.NodeOpenBracket})
	ret.AppendChild(&ast.Node{Type: ast.NodeLinkText, Tokens: text})
	ret.AppendChild(&ast.Node{Type: ast.NodeCloseBracket})
	ret.AppendChild(&ast.Node{Type: ast.NodeOpenParen})
	ret.AppendChild(&ast.Node{Type: ast.NodeLinkDest, Tokens: dest})
	if nil != title {
		ret.AppendChild(&ast.Node{Type: ast.NodeLinkTitle, Tokens: title})
	}
	ret.AppendChild(&ast.Node{Type: ast.NodeCloseParen})
	if appendCaret {
		ret.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: util.CaretTokens})
	}
	if 1 == linkType {
		ret.LinkRefLabel = text
	}
	return
}

func (t *Tree) parseAutolink(ctx *InlineContext) (ret *ast.Node) {
	schemed := false
	scheme := ""
	var dest []byte
	var token byte
	i := ctx.pos + 1
	for ; i < ctx.tokensLen && lex.ItemGreater != ctx.tokens[i]; i++ {
		token = ctx.tokens[i]
		if lex.ItemSpace == token {
			return nil
		}

		dest = append(dest, ctx.tokens[i])
		if !schemed {
			if lex.ItemColon != token {
				scheme += string(token)
			} else {
				schemed = true
			}
		}
	}
	if !schemed || 3 > len(scheme) || i == ctx.tokensLen {
		return nil
	}

	if lex.ItemGreater != ctx.tokens[i] {
		return nil
	}
	ctx.pos = 1 + i
	return t.newLink(ast.NodeLink, dest, html.EncodeDestination(dest), nil, 2)
}

func (t *Tree) addPreviousText(node *ast.Node, tokens []byte) {
	if nil == node.Previous || ast.NodeText != node.Previous.Type {
		node.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: tokens})
		return
	}
	node.Previous.AppendTokens(tokens)
}
