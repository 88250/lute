// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lex

import (
	"bytes"
	"unicode"
)

const (
	ItemBacktick       = byte('`')
	ItemTilde          = byte('~')
	ItemBang           = byte('!')
	ItemCrosshatch     = byte('#')
	ItemAsterisk       = byte('*')
	ItemOpenParen      = byte('(')
	ItemCloseParen     = byte(')')
	ItemHyphen         = byte('-')
	ItemUnderscore     = byte('_')
	ItemPlus           = byte('+')
	ItemEqual          = byte('=')
	ItemTab            = byte('\t')
	ItemOpenBracket    = byte('[')
	ItemCloseBracket   = byte(']')
	ItemDoublequote    = byte('"')
	ItemSinglequote    = byte('\'')
	ItemLess           = byte('<')
	ItemGreater        = byte('>')
	ItemSpace          = byte(' ')
	ItemNewline        = byte('\n')
	ItemCarriageReturn = byte('\r')
	ItemBackslash      = byte('\\')
	ItemSlash          = byte('/')
	ItemDot            = byte('.')
	ItemColon          = byte(':')
	ItemQuestion       = byte('?')
	ItemAmpersand      = byte('&')
	ItemSemicolon      = byte(';')
	ItemPipe           = byte('|')
	ItemDollar         = byte('$')
	ItemCaret          = byte('^')
	ItemOpenBrace      = byte('{')
	ItemCloseBrace     = byte('}')
)

// IsWhitespace 判断 token 是否是空白。
func IsWhitespace(token byte) bool {
	return ItemSpace == token || ItemNewline == token || ItemTab == token || '\u000B' == token || '\u000C' == token || '\u000D' == token
}

// IsUnicodeWhitespace 判断 token 是否是 Unicode 空白。
func IsUnicodeWhitespace(r rune) bool {
	return unicode.IsSpace(r) || unicode.Is(unicode.Zs, r)
}

// IsDigit 判断 token 是否为数字 0-9。
func IsDigit(token byte) bool {
	return '0' <= token && '9' >= token
}

// IsHexDigit 判断 token 是否是十六进制数字。
func IsHexDigit(token byte) bool {
	return IsDigit(token) || token >= 'a' && token <= 'f' || token >= 'A' && token <= 'F'
}

// TokenToUpper 将 token 转为大写。
func TokenToUpper(token byte) byte {
	if token >= 'a' && token <= 'z' {
		return token - 'a' + 'A'
	}
	return token
}

// IsASCIIPunct 判断 token 是否是一个 ASCII 标点符号。
func IsASCIIPunct(token byte) bool {
	return (0x21 <= token && 0x2F >= token) || (0x3A <= token && 0x40 >= token) || (0x5B <= token && 0x60 >= token) || (0x7B <= token && 0x7E >= token)
}

// IsASCIILetter 判断 token 是否是一个 ASCII 字母。
func IsASCIILetter(token byte) bool {
	return ('A' <= token && 'Z' >= token) || ('a' <= token && 'z' >= token)
}

// IsASCIILetterNum 判断 token 是否是一个 ASCII 字母或数字。
func IsASCIILetterNum(token byte) bool {
	return ('A' <= token && 'Z' >= token) || ('a' <= token && 'z' >= token) || ('0' <= token && '9' >= token)
}

// IsASCIILetterNums 判断 tokens 是否是 ASCII 字母或数字组成。
func IsASCIILetterNums(tokens []byte) bool {
	for _, token := range tokens {
		if !IsASCIILetterNum(token) {
			return false
		}
	}
	return true
}

// IsASCIILetterNumHyphen 判断 token 是否是一个 ASCII 字母、数字或者横线 -。
func IsASCIILetterNumHyphen(token byte) bool {
	return ('A' <= token && 'Z' >= token) || ('a' <= token && 'z' >= token) || ('0' <= token && '9' >= token) || '-' == token
}

// IsControl 判断 token 是否是一个控制字符。
func IsControl(token byte) bool {
	return unicode.IsControl(rune(token))
}

// IsBlank 判断 Tokens 是否都为空格。
func IsBlank(tokens []byte) bool {
	for _, token := range tokens {
		if ItemSpace != token {
			return false
		}
	}
	return true
}

func Split(tokens []byte, separator byte) (ret [][]byte) {
	length := len(tokens)
	var i int
	var token byte
	var line []byte
	for ; i < length; i++ {
		token = tokens[i]
		if separator != token {
			line = append(line, token)
			continue
		}

		ret = append(ret, line)
		line = []byte{}
	}
	if 0 < len(line) {
		ret = append(ret, line)
	}
	return
}

// SplitWithoutBackslashEscape 使用 separator 作为分隔符将 Tokens 切分为多个子串，被反斜杠 \ 转义的字符不会计入切分。
func SplitWithoutBackslashEscape(tokens []byte, separator byte) (ret [][]byte) {
	length := len(tokens)
	var token byte
	var line []byte
	for i := 0; i < length; i++ {
		token = tokens[i]
		if separator != token || IsBackslashEscapePunct(tokens, i) {
			line = append(line, token)
			continue
		}

		//if ItemPipe == token && inInlineMath(tokens, i) {
		//	line = append(line, token)
		//	continue
		//}

		ret = append(ret, line)
		line = []byte{}
	}
	if 0 < len(line) {
		ret = append(ret, line)
	}
	return
}

// ReplaceAll 会将 Tokens 中的所有 old 使用 new 替换。
func ReplaceAll(tokens []byte, old, new byte) []byte {
	for i, token := range tokens {
		if token == old {
			tokens[i] = new
		}
	}
	return tokens
}

// ReplaceNewlineSpace 会将 Tokens 中的所有 "\n " 替换为 "\n"。
func ReplaceNewlineSpace(tokens []byte) []byte {
	length := len(tokens)
	var token byte
	for i := length - 1; 0 <= i; i-- {
		token = tokens[i]
		if ItemNewline != token && ItemSpace != token {
			break
		}
		if ItemNewline == tokens[i-1] && (ItemSpace == token || ItemNewline == token) {
			tokens = tokens[:i]
		}
	}
	return tokens
}

func TrimWhitespace(tokens []byte) (ret []byte) {
	_, _, ret = Trim(tokens)
	return
}

func Trim(tokens []byte) (leftWhitespaces, rightWhitespaces, remains []byte) {
	length := len(tokens)
	if 0 == length {
		return nil, nil, tokens
	}
	start, end := 0, length-1
	for ; start < length; start++ {
		if !IsWhitespace(tokens[start]) {
			break
		}
	}
	leftWhitespaces = tokens[0:start]
	if start == length {
		start--
	}
	for ; 0 <= end; end-- {
		if !IsWhitespace(tokens[end]) {
			break
		}
	}
	if end < start {
		end = start - 1
	}
	if 0 < end {
		rightWhitespaces = tokens[end+1 : length]
	}
	remains = tokens[start : end+1]
	return
}

func TrimRight(tokens []byte) (whitespaces, remains []byte) {
	length := len(tokens)
	if 1 > length {
		return nil, tokens
	}

	i := length - 1
	for ; 0 <= i; i-- {
		if !IsWhitespace(tokens[i]) {
			break
		}
		whitespaces = append(whitespaces, tokens[i])
	}
	return whitespaces, tokens[:i+1]
}

func TrimLeft(tokens []byte) (whitespaces, remains []byte) {
	length := len(tokens)
	if 1 > length {
		return nil, tokens
	}

	i := 0
	for ; i < length; i++ {
		if !IsWhitespace(tokens[i]) {
			break
		}
		whitespaces = append(whitespaces, tokens[i])
	}
	return whitespaces, tokens[i:]
}

func Accept(tokens []byte, token byte) (pos int) {
	length := len(tokens)
	for ; pos < length; pos++ {
		if token != tokens[pos] {
			break
		}
	}
	return
}

func AcceptTokenss(tokens []byte, someTokenss [][]byte) (pos int) {
	length := len(tokens)
	length2 := len(someTokenss)
	for i := 0; i < length; i++ {
		remains := tokens[i:]
		for j := 0; j < length2; j++ {
			someTokens := someTokenss[j]
			if pos = AcceptTokens(remains, someTokens); 0 <= pos {
				return
			}
		}
	}
	return -1
}

func AcceptTokens(remains, someTokens []byte) (pos int) {
	length := len(someTokens)
	for ; pos < length; pos++ {
		if someTokens[pos] != remains[pos] {
			return -1
		}
	}
	return
}

func IsBlankLine(tokens []byte) bool {
	for _, token := range tokens {
		if ItemSpace != token && ItemTab != token && ItemNewline != token {
			return false
		}
	}
	return true
}

func SplitWhitespace(tokens []byte) (ret [][]byte) {
	i := 0
	ret = append(ret, []byte{})
	lastIsWhitespace := false
	for _, token := range tokens {
		if IsWhitespace(token) {
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

// IsBackslashEscapePunct 判断 Tokens 中 pos 所指的值是否是由反斜杠 \ 转义的 ASCII 标点符号。
func IsBackslashEscapePunct(tokens []byte, pos int) bool {
	if !IsASCIIPunct(tokens[pos]) {
		return false
	}

	backslashes := 0
	for i := pos - 1; 0 <= i; i-- {
		if ItemBackslash != tokens[i] {
			break
		}
		backslashes++
	}
	return 0 != backslashes%2
}

func StatWhitespace(tokens []byte) (newlines, spaces, tabs int) {
	for _, token := range tokens {
		if ItemNewline == token {
			newlines++
		} else if ItemSpace == token {
			spaces++
		} else if ItemTab == token {
			tabs++
		}
	}
	return
}

func Spnl(tokens []byte) (ret bool, passed, remains []byte) {
	passed, remains = TrimLeft(tokens)
	newlines, _, _ := StatWhitespace(passed)
	if 1 < newlines {
		return false, nil, tokens
	}
	ret = true
	return
}

func Peek(tokens []byte, pos int) byte {
	if pos < len(tokens) {
		return tokens[pos]
	}
	return 0
}

// BytesShowLength 获取字节数组展示为 UTF8 字符串时的长度。
func BytesShowLength(bytes []byte) int {
	length := 0
	for i := 0; i < len(bytes); i++ {
		// 按位与 11000000 为 10000000 则表示为 UTF8 字节首位
		if (bytes[i] & 0xc0) != 0x80 {
			if bytes[i] < 0x7f {
				length++
			} else {
				length += 2
			}
		}
	}
	return length
}

func RepeatBackslashBeforePipe(content string) string {
	buf := bytes.Buffer{}
	var last byte
	backslashCnt := 0
	for i := 0; i < len(content); i++ {
		b := content[i]
		if ItemPipe == b {
			if ItemBackslash != last {
				buf.WriteByte(ItemBackslash)
			}
			if 1 <= backslashCnt {
				buf.WriteByte(ItemBackslash)
			}
		}
		last = b
		if ItemBackslash == last {
			backslashCnt++
		} else {
			backslashCnt = 0
		}
		buf.WriteByte(b)
	}
	return buf.String()
}

func EscapeCommonMarkers(tokens []byte) []byte {
	for i := 0; i < len(tokens); i++ {
		if IsCommonInlineMarker(tokens[i]) {
			remains := append([]byte{ItemBackslash}, tokens[i:]...)
			tokens = tokens[:i]
			tokens = append(tokens, remains...)
			i++
		}
	}
	return tokens
}

func EscapeProtyleMarkers(tokens []byte) []byte {
	for i := 0; i < len(tokens); i++ {
		if IsProtyleInlineMarker(tokens[i]) {
			remains := append([]byte{ItemBackslash}, tokens[i:]...)
			tokens = tokens[:i]
			tokens = append(tokens, remains...)
			i++
		}
	}
	return tokens
}

func IsCommonInlineMarker(token byte) bool {
	switch token {
	case ItemAsterisk, ItemUnderscore, ItemBackslash, ItemBacktick, ItemTilde, ItemDollar:
		return true
	default:
		return false
	}
}

func IsProtyleInlineMarker(token byte) bool {
	switch token {
	case ItemAsterisk, ItemUnderscore, ItemBackslash, ItemBacktick, ItemTilde, ItemDollar, ItemEqual, ItemCaret, ItemLess, ItemGreater:
		return true
	default:
		return false
	}
}

func IsMarker(token byte) bool {
	switch token {
	case ItemAsterisk, ItemUnderscore, ItemOpenBracket, ItemBang, ItemNewline, ItemBackslash, ItemBacktick, ItemLess,
		ItemCloseBracket, ItemAmpersand, ItemTilde, ItemDollar, ItemOpenBrace, ItemOpenParen, ItemEqual, ItemCrosshatch:
		return true
	default:
		return false
	}
}
