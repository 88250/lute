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
	"strings"
	"unicode/utf8"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
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

// AddAutoLinkDomainSuffix 添加自动链接解析域名后缀 suffix。
func AddAutoLinkDomainSuffix(suffix string) {
	validAutoLinkDomainSuffix[suffix] = true
}

func (t *Tree) parseGFMAutoLink0(node *ast.Node) {
	tokens := node.Tokens
	length := len(tokens)
	minLinkLen := 5 // 太短的情况肯定不可能有链接，最短的情况是 a://b 或者 www.xxx.xx
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
		} else if parts := bytes.Split(tokens[i:], []byte("://")); 2 == len(parts) && 0 < len(parts[0]) && 0 < len(parts[1]) && !bytes.Contains(tokens[i:], httpProto) && !bytes.Contains(tokens[i:], httpsProto) && !bytes.Contains(tokens[i:], ftpProto) {
			if !lex.IsASCIILetterNums(parts[0]) {
				textEnd++
				i++
				continue
			}

			// 自定义协议均认为是有效的 https://github.com/siyuan-note/siyuan/issues/5865
			protocol = append(parts[0], []byte("://")...)
			i += len(parts[0]) + 3
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
			if lex.IsWhitespace(token) || lex.ItemLess == token {
				break
			}

			// 判断端口后部分是否为数字
			if tmp := bytes.ReplaceAll(url, []byte("://"), nil); bytes.Contains(tmp, []byte(":")) && !bytes.Contains(tmp, []byte("/")) {
				tmp = tmp[bytes.Index(tmp, []byte(":"))+1:]
				if !bytes.Contains(tmp, []byte("/")) && !lex.IsDigit(token) && lex.ItemSlash != token {
					break
				}
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

		if !t.isValidDomain(protocol, domain) {
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
			// Trailing punctuation (specifically, ?, !, ., ,, :, *, _, and ~) will not be considered part of the autolink, though they may be included in the interior of the link:
			// https://github.github.com/gfm/#example-624
			if !trimmed && ('?' == lastToken || '!' == lastToken || '.' == lastToken || ',' == lastToken || ':' == lastToken || '*' == lastToken || '_' == lastToken || '~' == lastToken) {
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
		addr = append(addr, port...)
		addr = append(addr, path...)
		linkText := addr
		if bytes.HasPrefix(linkText, []byte("https://github.com/")) {
			if bytes.Contains(linkText, []byte("/issues/")) {
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
			} else if bytes.Contains(linkText, []byte("/pull/")) {
				// 优化 GitHub Pull Requests 自动链接文本 https://github.com/88250/lute/issues/208
				repo := linkText[len("https://github.com/"):]
				repo = repo[:bytes.Index(repo, []byte("/pull/"))]
				num := bytes.Split(linkText, []byte("/pull/"))[1]
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
						linkText = []byte("Pull Request #" + string(num) + " · " + string(repo))
					}
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
func (t *Tree) isValidDomain(protocol, domain []byte) bool {
	if 0 < len(protocol) && !bytes.Contains(protocol, httpProto) && !bytes.Contains(protocol, httpsProto) && !bytes.Contains(protocol, ftpProto) {
		// 自定义协议均认为是有效的 https://github.com/siyuan-note/siyuan/issues/5865
		return true
	}

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
				validSuffix = validAutoLinkDomainSuffix[util.BytesToStr(segment)]
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
	appendCaret := t.Context.ParseOption.ProtyleWYSIWYG && bytes.HasSuffix(text, editor.CaretTokens) && bytes.HasSuffix(dest, []byte("%E2%80%B8"))
	if appendCaret {
		text = bytes.ReplaceAll(text, editor.CaretTokens, nil)
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
		ret.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: editor.CaretTokens})
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

var (
	httpProto  = util.StrToBytes("http://")
	httpsProto = util.StrToBytes("https://")
	ftpProto   = util.StrToBytes("ftp://")

	// validAutoLinkDomainSuffix 作为 GFM 自动连接解析时校验域名后缀用。
	validAutoLinkDomainSuffix = getValidDomainSuffixFromStr()
)

// getValidDomainSuffixFromStr 将字符串形式顶级域名转换为 map。
func getValidDomainSuffixFromStr() (ret map[string]bool) {
	ret = map[string]bool{}
	split := strings.Split(allTLDs, "\n")
	var l1, l2 string
	for _, line := range split {
		l1 = strings.Trim(line, " ")
		l2 = strings.ToLower(l1)
		if 1 > len(l2) || strings.HasPrefix(l2, "#") {
			continue
		}
		ret[l2] = true
	}
	return
}

// 所有顶级域名(Top-Level Domain),
// 一行对应一个,#号开头为注释;
// xn--开头的是 Punycode,对应非英语域名; 比如中文网址"中国移动.中国" 转换后为 xn--fiq02ib9d179b.xn--fiqs8s
// 来自: List of Top-Level Domains - ICANN
// https://www.icann.org/resources/pages/tlds-2012-02-25-en
var allTLDs = `
# Version 2023021700, Last Updated Fri Feb 17 07:07:01 2023 UTC
AAA
AARP
ABARTH
ABB
ABBOTT
ABBVIE
ABC
ABLE
ABOGADO
ABUDHABI
AC
ACADEMY
ACCENTURE
ACCOUNTANT
ACCOUNTANTS
ACO
ACTOR
AD
ADS
ADULT
AE
AEG
AERO
AETNA
AF
AFL
AFRICA
AG
AGAKHAN
AGENCY
AI
AIG
AIRBUS
AIRFORCE
AIRTEL
AKDN
AL
ALFAROMEO
ALIBABA
ALIPAY
ALLFINANZ
ALLSTATE
ALLY
ALSACE
ALSTOM
AM
AMAZON
AMERICANEXPRESS
AMERICANFAMILY
AMEX
AMFAM
AMICA
AMSTERDAM
ANALYTICS
ANDROID
ANQUAN
ANZ
AO
AOL
APARTMENTS
APP
APPLE
AQ
AQUARELLE
AR
ARAB
ARAMCO
ARCHI
ARMY
ARPA
ART
ARTE
AS
ASDA
ASIA
ASSOCIATES
AT
ATHLETA
ATTORNEY
AU
AUCTION
AUDI
AUDIBLE
AUDIO
AUSPOST
AUTHOR
AUTO
AUTOS
AVIANCA
AW
AWS
AX
AXA
AZ
AZURE
BA
BABY
BAIDU
BANAMEX
BANANAREPUBLIC
BAND
BANK
BAR
BARCELONA
BARCLAYCARD
BARCLAYS
BAREFOOT
BARGAINS
BASEBALL
BASKETBALL
BAUHAUS
BAYERN
BB
BBC
BBT
BBVA
BCG
BCN
BD
BE
BEATS
BEAUTY
BEER
BENTLEY
BERLIN
BEST
BESTBUY
BET
BF
BG
BH
BHARTI
BI
BIBLE
BID
BIKE
BING
BINGO
BIO
BIZ
BJ
BLACK
BLACKFRIDAY
BLOCKBUSTER
BLOG
BLOOMBERG
BLUE
BM
BMS
BMW
BN
BNPPARIBAS
BO
BOATS
BOEHRINGER
BOFA
BOM
BOND
BOO
BOOK
BOOKING
BOSCH
BOSTIK
BOSTON
BOT
BOUTIQUE
BOX
BR
BRADESCO
BRIDGESTONE
BROADWAY
BROKER
BROTHER
BRUSSELS
BS
BT
BUILD
BUILDERS
BUSINESS
BUY
BUZZ
BV
BW
BY
BZ
BZH
CA
CAB
CAFE
CAL
CALL
CALVINKLEIN
CAM
CAMERA
CAMP
CANON
CAPETOWN
CAPITAL
CAPITALONE
CAR
CARAVAN
CARDS
CARE
CAREER
CAREERS
CARS
CASA
CASE
CASH
CASINO
CAT
CATERING
CATHOLIC
CBA
CBN
CBRE
CBS
CC
CD
CENTER
CEO
CERN
CF
CFA
CFD
CG
CH
CHANEL
CHANNEL
CHARITY
CHASE
CHAT
CHEAP
CHINTAI
CHRISTMAS
CHROME
CHURCH
CI
CIPRIANI
CIRCLE
CISCO
CITADEL
CITI
CITIC
CITY
CITYEATS
CK
CL
CLAIMS
CLEANING
CLICK
CLINIC
CLINIQUE
CLOTHING
CLOUD
CLUB
CLUBMED
CM
CN
CO
COACH
CODES
COFFEE
COLLEGE
COLOGNE
COM
COMCAST
COMMBANK
COMMUNITY
COMPANY
COMPARE
COMPUTER
COMSEC
CONDOS
CONSTRUCTION
CONSULTING
CONTACT
CONTRACTORS
COOKING
COOKINGCHANNEL
COOL
COOP
CORSICA
COUNTRY
COUPON
COUPONS
COURSES
CPA
CR
CREDIT
CREDITCARD
CREDITUNION
CRICKET
CROWN
CRS
CRUISE
CRUISES
CU
CUISINELLA
CV
CW
CX
CY
CYMRU
CYOU
CZ
DABUR
DAD
DANCE
DATA
DATE
DATING
DATSUN
DAY
DCLK
DDS
DE
DEAL
DEALER
DEALS
DEGREE
DELIVERY
DELL
DELOITTE
DELTA
DEMOCRAT
DENTAL
DENTIST
DESI
DESIGN
DEV
DHL
DIAMONDS
DIET
DIGITAL
DIRECT
DIRECTORY
DISCOUNT
DISCOVER
DISH
DIY
DJ
DK
DM
DNP
DO
DOCS
DOCTOR
DOG
DOMAINS
DOT
DOWNLOAD
DRIVE
DTV
DUBAI
DUNLOP
DUPONT
DURBAN
DVAG
DVR
DZ
EARTH
EAT
EC
ECO
EDEKA
EDU
EDUCATION
EE
EG
EMAIL
EMERCK
ENERGY
ENGINEER
ENGINEERING
ENTERPRISES
EPSON
EQUIPMENT
ER
ERICSSON
ERNI
ES
ESQ
ESTATE
ET
ETISALAT
EU
EUROVISION
EUS
EVENTS
EXCHANGE
EXPERT
EXPOSED
EXPRESS
EXTRASPACE
FAGE
FAIL
FAIRWINDS
FAITH
FAMILY
FAN
FANS
FARM
FARMERS
FASHION
FAST
FEDEX
FEEDBACK
FERRARI
FERRERO
FI
FIAT
FIDELITY
FIDO
FILM
FINAL
FINANCE
FINANCIAL
FIRE
FIRESTONE
FIRMDALE
FISH
FISHING
FIT
FITNESS
FJ
FK
FLICKR
FLIGHTS
FLIR
FLORIST
FLOWERS
FLY
FM
FO
FOO
FOOD
FOODNETWORK
FOOTBALL
FORD
FOREX
FORSALE
FORUM
FOUNDATION
FOX
FR
FREE
FRESENIUS
FRL
FROGANS
FRONTDOOR
FRONTIER
FTR
FUJITSU
FUN
FUND
FURNITURE
FUTBOL
FYI
GA
GAL
GALLERY
GALLO
GALLUP
GAME
GAMES
GAP
GARDEN
GAY
GB
GBIZ
GD
GDN
GE
GEA
GENT
GENTING
GEORGE
GF
GG
GGEE
GH
GI
GIFT
GIFTS
GIVES
GIVING
GL
GLASS
GLE
GLOBAL
GLOBO
GM
GMAIL
GMBH
GMO
GMX
GN
GODADDY
GOLD
GOLDPOINT
GOLF
GOO
GOODYEAR
GOOG
GOOGLE
GOP
GOT
GOV
GP
GQ
GR
GRAINGER
GRAPHICS
GRATIS
GREEN
GRIPE
GROCERY
GROUP
GS
GT
GU
GUARDIAN
GUCCI
GUGE
GUIDE
GUITARS
GURU
GW
GY
HAIR
HAMBURG
HANGOUT
HAUS
HBO
HDFC
HDFCBANK
HEALTH
HEALTHCARE
HELP
HELSINKI
HERE
HERMES
HGTV
HIPHOP
HISAMITSU
HITACHI
HIV
HK
HKT
HM
HN
HOCKEY
HOLDINGS
HOLIDAY
HOMEDEPOT
HOMEGOODS
HOMES
HOMESENSE
HONDA
HORSE
HOSPITAL
HOST
HOSTING
HOT
HOTELES
HOTELS
HOTMAIL
HOUSE
HOW
HR
HSBC
HT
HU
HUGHES
HYATT
HYUNDAI
IBM
ICBC
ICE
ICU
ID
IE
IEEE
IFM
IKANO
IL
IM
IMAMAT
IMDB
IMMO
IMMOBILIEN
IN
INC
INDUSTRIES
INFINITI
INFO
ING
INK
INSTITUTE
INSURANCE
INSURE
INT
INTERNATIONAL
INTUIT
INVESTMENTS
IO
IPIRANGA
IQ
IR
IRISH
IS
ISMAILI
IST
ISTANBUL
IT
ITAU
ITV
JAGUAR
JAVA
JCB
JE
JEEP
JETZT
JEWELRY
JIO
JLL
JM
JMP
JNJ
JO
JOBS
JOBURG
JOT
JOY
JP
JPMORGAN
JPRS
JUEGOS
JUNIPER
KAUFEN
KDDI
KE
KERRYHOTELS
KERRYLOGISTICS
KERRYPROPERTIES
KFH
KG
KH
KI
KIA
KIDS
KIM
KINDER
KINDLE
KITCHEN
KIWI
KM
KN
KOELN
KOMATSU
KOSHER
KP
KPMG
KPN
KR
KRD
KRED
KUOKGROUP
KW
KY
KYOTO
KZ
LA
LACAIXA
LAMBORGHINI
LAMER
LANCASTER
LANCIA
LAND
LANDROVER
LANXESS
LASALLE
LAT
LATINO
LATROBE
LAW
LAWYER
LB
LC
LDS
LEASE
LECLERC
LEFRAK
LEGAL
LEGO
LEXUS
LGBT
LI
LIDL
LIFE
LIFEINSURANCE
LIFESTYLE
LIGHTING
LIKE
LILLY
LIMITED
LIMO
LINCOLN
LINDE
LINK
LIPSY
LIVE
LIVING
LK
LLC
LLP
LOAN
LOANS
LOCKER
LOCUS
LOL
LONDON
LOTTE
LOTTO
LOVE
LPL
LPLFINANCIAL
LR
LS
LT
LTD
LTDA
LU
LUNDBECK
LUXE
LUXURY
LV
LY
MA
MACYS
MADRID
MAIF
MAISON
MAKEUP
MAN
MANAGEMENT
MANGO
MAP
MARKET
MARKETING
MARKETS
MARRIOTT
MARSHALLS
MASERATI
MATTEL
MBA
MC
MCKINSEY
MD
ME
MED
MEDIA
MEET
MELBOURNE
MEME
MEMORIAL
MEN
MENU
MERCKMSD
MG
MH
MIAMI
MICROSOFT
MIL
MINI
MINT
MIT
MITSUBISHI
MK
ML
MLB
MLS
MM
MMA
MN
MO
MOBI
MOBILE
MODA
MOE
MOI
MOM
MONASH
MONEY
MONSTER
MORMON
MORTGAGE
MOSCOW
MOTO
MOTORCYCLES
MOV
MOVIE
MP
MQ
MR
MS
MSD
MT
MTN
MTR
MU
MUSEUM
MUSIC
MUTUAL
MV
MW
MX
MY
MZ
NA
NAB
NAGOYA
NAME
NATURA
NAVY
NBA
NC
NE
NEC
NET
NETBANK
NETFLIX
NETWORK
NEUSTAR
NEW
NEWS
NEXT
NEXTDIRECT
NEXUS
NF
NFL
NG
NGO
NHK
NI
NICO
NIKE
NIKON
NINJA
NISSAN
NISSAY
NL
NO
NOKIA
NORTHWESTERNMUTUAL
NORTON
NOW
NOWRUZ
NOWTV
NP
NR
NRA
NRW
NTT
NU
NYC
NZ
OBI
OBSERVER
OFFICE
OKINAWA
OLAYAN
OLAYANGROUP
OLDNAVY
OLLO
OM
OMEGA
ONE
ONG
ONL
ONLINE
OOO
OPEN
ORACLE
ORANGE
ORG
ORGANIC
ORIGINS
OSAKA
OTSUKA
OTT
OVH
PA
PAGE
PANASONIC
PARIS
PARS
PARTNERS
PARTS
PARTY
PASSAGENS
PAY
PCCW
PE
PET
PF
PFIZER
PG
PH
PHARMACY
PHD
PHILIPS
PHONE
PHOTO
PHOTOGRAPHY
PHOTOS
PHYSIO
PICS
PICTET
PICTURES
PID
PIN
PING
PINK
PIONEER
PIZZA
PK
PL
PLACE
PLAY
PLAYSTATION
PLUMBING
PLUS
PM
PN
PNC
POHL
POKER
POLITIE
PORN
POST
PR
PRAMERICA
PRAXI
PRESS
PRIME
PRO
PROD
PRODUCTIONS
PROF
PROGRESSIVE
PROMO
PROPERTIES
PROPERTY
PROTECTION
PRU
PRUDENTIAL
PS
PT
PUB
PW
PWC
PY
QA
QPON
QUEBEC
QUEST
RACING
RADIO
RE
READ
REALESTATE
REALTOR
REALTY
RECIPES
RED
REDSTONE
REDUMBRELLA
REHAB
REISE
REISEN
REIT
RELIANCE
REN
RENT
RENTALS
REPAIR
REPORT
REPUBLICAN
REST
RESTAURANT
REVIEW
REVIEWS
REXROTH
RICH
RICHARDLI
RICOH
RIL
RIO
RIP
RO
ROCHER
ROCKS
RODEO
ROGERS
ROOM
RS
RSVP
RU
RUGBY
RUHR
RUN
RW
RWE
RYUKYU
SA
SAARLAND
SAFE
SAFETY
SAKURA
SALE
SALON
SAMSCLUB
SAMSUNG
SANDVIK
SANDVIKCOROMANT
SANOFI
SAP
SARL
SAS
SAVE
SAXO
SB
SBI
SBS
SC
SCA
SCB
SCHAEFFLER
SCHMIDT
SCHOLARSHIPS
SCHOOL
SCHULE
SCHWARZ
SCIENCE
SCOT
SD
SE
SEARCH
SEAT
SECURE
SECURITY
SEEK
SELECT
SENER
SERVICES
SEVEN
SEW
SEX
SEXY
SFR
SG
SH
SHANGRILA
SHARP
SHAW
SHELL
SHIA
SHIKSHA
SHOES
SHOP
SHOPPING
SHOUJI
SHOW
SHOWTIME
SI
SILK
SINA
SINGLES
SITE
SJ
SK
SKI
SKIN
SKY
SKYPE
SL
SLING
SM
SMART
SMILE
SN
SNCF
SO
SOCCER
SOCIAL
SOFTBANK
SOFTWARE
SOHU
SOLAR
SOLUTIONS
SONG
SONY
SOY
SPA
SPACE
SPORT
SPOT
SR
SRL
SS
ST
STADA
STAPLES
STAR
STATEBANK
STATEFARM
STC
STCGROUP
STOCKHOLM
STORAGE
STORE
STREAM
STUDIO
STUDY
STYLE
SU
SUCKS
SUPPLIES
SUPPLY
SUPPORT
SURF
SURGERY
SUZUKI
SV
SWATCH
SWISS
SX
SY
SYDNEY
SYSTEMS
SZ
TAB
TAIPEI
TALK
TAOBAO
TARGET
TATAMOTORS
TATAR
TATTOO
TAX
TAXI
TC
TCI
TD
TDK
TEAM
TECH
TECHNOLOGY
TEL
TEMASEK
TENNIS
TEVA
TF
TG
TH
THD
THEATER
THEATRE
TIAA
TICKETS
TIENDA
TIFFANY
TIPS
TIRES
TIROL
TJ
TJMAXX
TJX
TK
TKMAXX
TL
TM
TMALL
TN
TO
TODAY
TOKYO
TOOLS
TOP
TORAY
TOSHIBA
TOTAL
TOURS
TOWN
TOYOTA
TOYS
TR
TRADE
TRADING
TRAINING
TRAVEL
TRAVELCHANNEL
TRAVELERS
TRAVELERSINSURANCE
TRUST
TRV
TT
TUBE
TUI
TUNES
TUSHU
TV
TVS
TW
TZ
UA
UBANK
UBS
UG
UK
UNICOM
UNIVERSITY
UNO
UOL
UPS
US
UY
UZ
VA
VACATIONS
VANA
VANGUARD
VC
VE
VEGAS
VENTURES
VERISIGN
VERSICHERUNG
VET
VG
VI
VIAJES
VIDEO
VIG
VIKING
VILLAS
VIN
VIP
VIRGIN
VISA
VISION
VIVA
VIVO
VLAANDEREN
VN
VODKA
VOLKSWAGEN
VOLVO
VOTE
VOTING
VOTO
VOYAGE
VU
VUELOS
WALES
WALMART
WALTER
WANG
WANGGOU
WATCH
WATCHES
WEATHER
WEATHERCHANNEL
WEBCAM
WEBER
WEBSITE
WED
WEDDING
WEIBO
WEIR
WF
WHOSWHO
WIEN
WIKI
WILLIAMHILL
WIN
WINDOWS
WINE
WINNERS
WME
WOLTERSKLUWER
WOODSIDE
WORK
WORKS
WORLD
WOW
WS
WTC
WTF
XBOX
XEROX
XFINITY
XIHUAN
XIN
XN--11B4C3D
XN--1CK2E1B
XN--1QQW23A
XN--2SCRJ9C
XN--30RR7Y
XN--3BST00M
XN--3DS443G
XN--3E0B707E
XN--3HCRJ9C
XN--3PXU8K
XN--42C2D9A
XN--45BR5CYL
XN--45BRJ9C
XN--45Q11C
XN--4DBRK0CE
XN--4GBRIM
XN--54B7FTA0CC
XN--55QW42G
XN--55QX5D
XN--5SU34J936BGSG
XN--5TZM5G
XN--6FRZ82G
XN--6QQ986B3XL
XN--80ADXHKS
XN--80AO21A
XN--80AQECDR1A
XN--80ASEHDB
XN--80ASWG
XN--8Y0A063A
XN--90A3AC
XN--90AE
XN--90AIS
XN--9DBQ2A
XN--9ET52U
XN--9KRT00A
XN--B4W605FERD
XN--BCK1B9A5DRE4C
XN--C1AVG
XN--C2BR7G
XN--CCK2B3B
XN--CCKWCXETD
XN--CG4BKI
XN--CLCHC0EA0B2G2A9GCD
XN--CZR694B
XN--CZRS0T
XN--CZRU2D
XN--D1ACJ3B
XN--D1ALF
XN--E1A4C
XN--ECKVDTC9D
XN--EFVY88H
XN--FCT429K
XN--FHBEI
XN--FIQ228C5HS
XN--FIQ64B
XN--FIQS8S
XN--FIQZ9S
XN--FJQ720A
XN--FLW351E
XN--FPCRJ9C3D
XN--FZC2C9E2C
XN--FZYS8D69UVGM
XN--G2XX48C
XN--GCKR3F0F
XN--GECRJ9C
XN--GK3AT1E
XN--H2BREG3EVE
XN--H2BRJ9C
XN--H2BRJ9C8C
XN--HXT814E
XN--I1B6B1A6A2E
XN--IMR513N
XN--IO0A7I
XN--J1AEF
XN--J1AMH
XN--J6W193G
XN--JLQ480N2RG
XN--JVR189M
XN--KCRX77D1X4A
XN--KPRW13D
XN--KPRY57D
XN--KPUT3I
XN--L1ACC
XN--LGBBAT1AD8J
XN--MGB9AWBF
XN--MGBA3A3EJT
XN--MGBA3A4F16A
XN--MGBA7C0BBN0A
XN--MGBAAKC7DVF
XN--MGBAAM7A8H
XN--MGBAB2BD
XN--MGBAH1A3HJKRD
XN--MGBAI9AZGQP6J
XN--MGBAYH7GPA
XN--MGBBH1A
XN--MGBBH1A71E
XN--MGBC0A9AZCG
XN--MGBCA7DZDO
XN--MGBCPQ6GPA1A
XN--MGBERP4A5D4AR
XN--MGBGU82A
XN--MGBI4ECEXP
XN--MGBPL2FH
XN--MGBT3DHD
XN--MGBTX2B
XN--MGBX4CD0AB
XN--MIX891F
XN--MK1BU44C
XN--MXTQ1M
XN--NGBC5AZD
XN--NGBE9E0A
XN--NGBRX
XN--NODE
XN--NQV7F
XN--NQV7FS00EMA
XN--NYQY26A
XN--O3CW4H
XN--OGBPF8FL
XN--OTU796D
XN--P1ACF
XN--P1AI
XN--PGBS0DH
XN--PSSY2U
XN--Q7CE6A
XN--Q9JYB4C
XN--QCKA1PMC
XN--QXA6A
XN--QXAM
XN--RHQV96G
XN--ROVU88B
XN--RVC1E0AM3E
XN--S9BRJ9C
XN--SES554G
XN--T60B56A
XN--TCKWE
XN--TIQ49XQYJ
XN--UNUP4Y
XN--VERMGENSBERATER-CTB
XN--VERMGENSBERATUNG-PWB
XN--VHQUV
XN--VUQ861B
XN--W4R85EL8FHU5DNRA
XN--W4RS40L
XN--WGBH1C
XN--WGBL6A
XN--XHQ521B
XN--XKC2AL3HYE2A
XN--XKC2DL3A5EE0H
XN--Y9A3AQ
XN--YFRO4I67O
XN--YGBI2AMMX
XN--ZFR164B
XXX
XYZ
YACHTS
YAHOO
YAMAXUN
YANDEX
YE
YODOBASHI
YOGA
YOKOHAMA
YOU
YOUTUBE
YT
YUN
ZA
ZAPPOS
ZARA
ZERO
ZIP
ZM
ZONE
ZUERICH
ZW
`
