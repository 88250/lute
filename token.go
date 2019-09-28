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
	"unicode"
)

// item 描述了词法分析的一个 token。
type item struct {
	term byte // 源码
	ln   int  // 源码行号
	col  int  // 源码列号
}

func isWhitespace(token byte) bool {
	return itemSpace == token || itemNewline == token || itemTab == token || '\u000B' == token || '\u000C' == token || '\u000D' == token
}

func isUnicodeWhitespace(r rune) bool {
	return unicode.IsSpace(r) || unicode.Is(unicode.Zs, r)
}

func isDigit(token byte) bool {
	return '0' <= token && '9' >= token
}

func isHexDigit(token byte) bool {
	return isDigit(token) || token >= 'a' && token <= 'f' || token >= 'A' && token <= 'F'
}

func tokenToUpper(token byte) byte {
	if token >= 'a' && token <= 'z' {
		return token - 'a' + 'A'
	}
	return token
}

// isASCIIPunct 判断 token 是否是一个 ASCII 标点符号。
func isASCIIPunct(token byte) bool {
	return (0x21 <= token && 0x2F >= token) || (0x3A <= token && 0x40 >= token) || (0x5B <= token && 0x60 >= token) || (0x7B <= token && 0x7E >= token)
}

// isASCIILetter 判断 token 是否是一个 ASCII 字母。
func isASCIILetter(token byte) bool {
	return ('A' <= token && 'Z' >= token) || ('a' <= token && 'z' >= token)
}

// isASCIILetterNum 判断 token 是否是一个 ASCII 字母或数字。
func isASCIILetterNum(token byte) bool {
	return ('A' <= token && 'Z' >= token) || ('a' <= token && 'z' >= token) || ('0' <= token && '9' >= token)
}

// isASCIILetterNumHyphen 判断 token 是否是一个 ASCII 字母、数字或者横线 -。
func isASCIILetterNumHyphen(token byte) bool {
	return ('A' <= token && 'Z' >= token) || ('a' <= token && 'z' >= token) || ('0' <= token && '9' >= token) || '-' == token
}

// isControl 判断 token 是否是一个控制字符。
func isControl(token byte) bool {
	return unicode.IsControl(rune(token))
}

// isBlank 判断 tokens 是否都为空格。
func isBlank(tokens []byte) bool {
	for _, token := range tokens {
		if itemSpace != token {
			return false
		}
	}
	return true
}

const (
	itemEnd            = byte(0)
	itemBacktick       = byte('`')
	itemTilde          = byte('~')
	itemBang           = byte('!')
	itemCrosshatch     = byte('#')
	itemAsterisk       = byte('*')
	itemOpenParen      = byte('(')
	itemCloseParen     = byte(')')
	itemHyphen         = byte('-')
	itemUnderscore     = byte('_')
	itemPlus           = byte('+')
	itemEqual          = byte('=')
	itemTab            = byte('\t')
	itemOpenBracket    = byte('[')
	itemCloseBracket   = byte(']')
	itemDoublequote    = byte('"')
	itemSinglequote    = byte('\'')
	itemLess           = byte('<')
	itemGreater        = byte('>')
	itemSpace          = byte(' ')
	itemNewline        = byte('\n')
	itemCarriageReturn = byte('\r')
	itemBackslash      = byte('\\')
	itemSlash          = byte('/')
	itemDot            = byte('.')
	itemColon          = byte(':')
	itemQuestion       = byte('?')
	itemAmpersand      = byte('&')
	itemSemicolon      = byte(';')
	itemPipe           = byte('|')
	itemDollar         = byte('$')
)

// items 定义了字节数组，每个字节是一个 token。
type items []*item

// strToItems 将 str 转为 items。
func strToItems(str string) (ret items) {
	ret = make(items, 0, len(str))
	length := len(str)
	for i := 0; i < length; i++ {
		ret = append(ret, &item{term: str[i]})
	}
	return
}

// itemsToStr 将 items 转为 string。
func itemsToStr(items items) string {
	return string(itemsToBytes(items))
}

// itemsToBytes 将 items 转为 []byte。
func itemsToBytes(items items) (ret []byte) {
	length := len(items)
	for i := 0; i < length; i++ {
		ret = append(ret, items[i].term)
	}
	return
}

// bytesToItems 将 bytes 转为 items。
func bytesToItems(bytes []byte) (ret items) {
	ret = make(items, 0, len(bytes))
	length := len(bytes)
	for i := 0; i < length; i++ {
		ret = append(ret, &item{term: bytes[i]})
	}
	return
}

func split(tokens items, separator byte) (ret []items) {
	length := len(tokens)
	var i int
	var token *item
	var line items
	for ; i < length; i++ {
		token = tokens[i]
		if separator != token.term {
			line = append(line, token)
			continue
		}

		ret = append(ret, line)
		line = items{}
	}
	if 0 < len(line) {
		ret = append(ret, line)
	}
	return
}

// splitWithoutBackslashEscape 使用 separator 作为分隔符将 tokens 切分为多个子串，被反斜杠 \ 转义的字符不会计入切分。
func (tokens items) splitWithoutBackslashEscape(separator byte) (ret []items) {
	length := len(tokens)
	var i int
	var token *item
	var line items
	for ; i < length; i++ {
		token = tokens[i]
		if separator != token.term || tokens.isBackslashEscapePunct(i) {
			line = append(line, token)
			continue
		}

		ret = append(ret, line)
		line = items{}
	}
	if 0 < len(line) {
		ret = append(ret, line)
	}
	return
}

func hasSuffix(tokens, suffix items) bool {
	return len(tokens) >= len(suffix) && equal(tokens[len(tokens)-len(suffix):], suffix)
}

func equal(a, b items) bool {
	length := len(a)
	if length != len(b) {
		return false
	}
	for i := 0; i < length; i++ {
		if a[i].term != b[i].term {
			return false
		}
	}
	return true
}

func index(tokens, sep items) (pos int) {
	a := itemsToBytes(tokens)
	b := itemsToBytes(sep)
	return bytes.Index(a, b)
}

func contains(tokens, sub items) bool {
	return 0 <= index(tokens, sub)
}

// replaceAll 会将 tokens 中的所有 old 使用 new 替换。
func replaceAll(tokens, old, new items) items {
	t := itemsToBytes(tokens)
	o := itemsToBytes(old)
	n := itemsToBytes(new)
	return bytesToItems(bytes.ReplaceAll(t, o, n))
}

// replaceNewlineSpace 会将 tokens 中的所有 "\n " 替换为 "\n"。
func replaceNewlineSpace(tokens items) items {
	length := len(tokens)
	var token byte
	for i := length - 1; 0 <= i; i-- {
		token = tokens[i].term
		if itemNewline != token && itemSpace != token {
			break
		}
		if itemNewline == tokens[i-1].term && (itemSpace == token || itemNewline == token) {
			tokens = tokens[:i]
		}
	}
	return tokens
}

func trimWhitespace(tokens items) items {
	length := len(tokens)
	if 0 == length {
		return tokens
	}
	start, end := 0, length-1
	for ; start < length; start++ {
		if !isWhitespace(tokens[start].term) {
			break
		}
	}
	if start == length {
		start--
	}
	for ; 0 <= end; end-- {
		if !isWhitespace(tokens[end].term) {
			break
		}
	}
	if end < start {
		end = start - 1
	}
	return tokens[start : end+1]
}

func trimRight(tokens items) (whitespaces, remains items) {
	length := len(tokens)
	if 1 > length {
		return nil, tokens
	}

	i := length - 1
	for ; 0 <= i; i-- {
		if !isWhitespace(tokens[i].term) {
			break
		}
		whitespaces = append(whitespaces, tokens[i])
	}
	return whitespaces, tokens[:i+1]
}

func trimLeft(tokens items) (whitespaces, remains items) {
	length := len(tokens)
	if 1 > length {
		return nil, tokens
	}

	i := 0
	for ; i < length; i++ {
		if !isWhitespace(tokens[i].term) {
			break
		}
		whitespaces = append(whitespaces, tokens[i])
	}
	return whitespaces, tokens[i:]
}

func (tokens items) accept(token byte) (pos int) {
	length := len(tokens)
	for ; pos < length; pos++ {
		if token != tokens[pos].term {
			break
		}
	}
	return
}

func acceptTokenss(tokens []byte, someTokenss [][]byte) (pos int) {
	length := len(tokens)
	length2 := len(someTokenss)
	for i := 0; i < length; i++ {
		remains := tokens[i:]
		for j := 0; j < length2; j++ {
			someTokens := someTokenss[j]
			if pos = acceptTokens(remains, someTokens); 0 <= pos {
				return
			}
		}
	}
	return -1
}

func acceptTokens(remains, someTokens []byte) (pos int) {
	length := len(someTokens)
	for ; pos < length; pos++ {
		if someTokens[pos] != remains[pos] {
			return -1
		}
	}
	return
}

func (tokens items) isBlankLine() bool {
	for _, token := range tokens {
		if itemSpace != token.term && itemTab != token.term && itemNewline != token.term {
			return false
		}
	}
	return true
}

func splitWhitespace(tokens []byte) (ret [][]byte) {
	i := 0
	ret = append(ret, []byte{})
	lastIsWhitespace := false
	for _, token := range tokens {
		if isWhitespace(token) {
			if !lastIsWhitespace {
				i++
				ret = append(ret, []byte{})
			}
			lastIsWhitespace = true
		} else {
			ret[i] = append(ret[i], token)
			lastIsWhitespace = false
		}
	}
	return
}

func (tokens items) startWith(token byte) bool {
	if 1 > len(tokens) {
		return false
	}
	return token == tokens[0].term
}

func (tokens items) endWith(token byte) bool {
	length := len(tokens)
	if 1 > length {
		return false
	}
	return token == tokens[length-1].term
}

// isBackslashEscapePunct 判断 tokens 中 pos 所指的值是否是由反斜杠 \ 转义的 ASCII 标点符号。
func (tokens items) isBackslashEscapePunct(pos int) bool {
	if !isASCIIPunct(tokens[pos].term) {
		return false
	}

	backslashes := 0
	for i := pos - 1; 0 <= i; i-- {
		if itemBackslash != tokens[i].term {
			break
		}
		backslashes++
	}
	return 0 != backslashes%2
}

func (tokens items) statWhitespace() (newlines, spaces, tabs int) {
	for _, token := range tokens {
		if itemNewline == token.term {
			newlines++
		} else if itemSpace == token.term {
			spaces++
		} else if itemTab == token.term {
			tabs++
		}
	}
	return
}

func (tokens items) spnl() (ret bool, passed, remains items) {
	passed, remains = trimLeft(tokens)
	newlines, _, _ := passed.statWhitespace()
	if 1 < newlines {
		return false, nil, tokens
	}
	ret = true
	return
}

func (tokens items) peek(pos int) *item {
	if pos < len(tokens) {
		return tokens[pos]
	}
	return nil
}
