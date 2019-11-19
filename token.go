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
	"unicode"
)

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

func split(tokens []byte, separator byte) (ret [][]byte) {
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

// splitWithoutBackslashEscape 使用 separator 作为分隔符将 tokens 切分为多个子串，被反斜杠 \ 转义的字符不会计入切分。
func splitWithoutBackslashEscape(tokens []byte, separator byte) (ret [][]byte) {
	length := len(tokens)
	var i int
	var token byte
	var line []byte
	for ; i < length; i++ {
		token = tokens[i]
		if separator != token || isBackslashEscapePunct(tokens, i) {
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

// replaceAll 会将 tokens 中的所有 old 使用 new 替换。
func replaceAll(tokens []byte, old, new byte) []byte {
	for i, token := range tokens {
		if token == old {
			tokens[i] = new
		}
	}
	return tokens
}

// replaceNewlineSpace 会将 tokens 中的所有 "\n " 替换为 "\n"。
func replaceNewlineSpace(tokens []byte) []byte {
	length := len(tokens)
	var token byte
	for i := length - 1; 0 <= i; i-- {
		token = tokens[i]
		if itemNewline != token && itemSpace != token {
			break
		}
		if itemNewline == tokens[i-1] && (itemSpace == token || itemNewline == token) {
			tokens = tokens[:i]
		}
	}
	return tokens
}

func trimWhitespace(tokens []byte) (ret []byte) {
	_, _, ret = trim(tokens)
	return
}

func trim(tokens []byte) (leftWhitespaces, rightWhitespaces, remains []byte) {
	length := len(tokens)
	if 0 == length {
		return nil, nil, tokens
	}
	start, end := 0, length-1
	for ; start < length; start++ {
		if !isWhitespace(tokens[start]) {
			break
		}
	}
	leftWhitespaces = tokens[0:start]
	if start == length {
		start--
	}
	for ; 0 <= end; end-- {
		if !isWhitespace(tokens[end]) {
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

func trimRight(tokens []byte) (whitespaces, remains []byte) {
	length := len(tokens)
	if 1 > length {
		return nil, tokens
	}

	i := length - 1
	for ; 0 <= i; i-- {
		if !isWhitespace(tokens[i]) {
			break
		}
		whitespaces = append(whitespaces, tokens[i])
	}
	return whitespaces, tokens[:i+1]
}

func trimLeft(tokens []byte) (whitespaces, remains []byte) {
	length := len(tokens)
	if 1 > length {
		return nil, tokens
	}

	i := 0
	for ; i < length; i++ {
		if !isWhitespace(tokens[i]) {
			break
		}
		whitespaces = append(whitespaces, tokens[i])
	}
	return whitespaces, tokens[i:]
}

func accept(tokens []byte, token byte) (pos int) {
	length := len(tokens)
	for ; pos < length; pos++ {
		if token != tokens[pos] {
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

func isBlankLine(tokens []byte) bool {
	for _, token := range tokens {
		if itemSpace != token && itemTab != token && itemNewline != token {
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

// isBackslashEscapePunct 判断 tokens 中 pos 所指的值是否是由反斜杠 \ 转义的 ASCII 标点符号。
func isBackslashEscapePunct(tokens []byte, pos int) bool {
	if !isASCIIPunct(tokens[pos]) {
		return false
	}

	backslashes := 0
	for i := pos - 1; 0 <= i; i-- {
		if itemBackslash != tokens[i] {
			break
		}
		backslashes++
	}
	return 0 != backslashes%2
}

func statWhitespace(tokens []byte) (newlines, spaces, tabs int) {
	for _, token := range tokens {
		if itemNewline == token {
			newlines++
		} else if itemSpace == token {
			spaces++
		} else if itemTab == token {
			tabs++
		}
	}
	return
}

func spnl(tokens []byte) (ret bool, passed, remains []byte) {
	passed, remains = trimLeft(tokens)
	newlines, _, _ := statWhitespace(passed)
	if 1 < newlines {
		return false, nil, tokens
	}
	ret = true
	return
}

func peek(tokens []byte, pos int) byte {
	if pos < len(tokens) {
		return tokens[pos]
	}
	return 0
}
