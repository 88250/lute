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

func isUnicodeWhitespace(token byte) bool {
	if unicode.Is(unicode.Zs, rune(token)) {
		return true
	}

	return itemTab == token || '\u000D' == token || itemNewline == token || '\u000C' == token
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

// isPunct 判断 token 是否是一个标点符号。
func isPunct(token byte) bool {
	return isASCIIPunct(token) || unicode.IsPunct(rune(token))
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
)

// items 定义了字节数组，每个字节是一个 token。
type items []byte

// splitWithoutBackslashEscape 使用 separator 作为分隔符将 tokens 切分为多个子串，被反斜杠 \ 转义的字符不会计入切分。
func (tokens items) splitWithoutBackslashEscape(separator byte) (ret []items) {
	length := len(tokens)
	var i int
	var token byte
	var line items
	for ; i < length; i++ {
		token = tokens[i]
		if separator != token || tokens.isBackslashEscapePunct(i) {
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

// split 使用 separator 作为分隔符将 tokens 切分为两个子串（仅分隔一次）。
func (tokens items) split(separator byte) (ret []items) {
	ret = append(ret, items{})
	var i int
	var token byte
	for i, token = range tokens {
		if separator == token {
			ret = append(ret, items{})
			ret[1] = append(ret[1], tokens[i+1:]...)
			return
		}
		ret[0] = append(ret[0], token)
	}
	return
}

// lines 会使用 \n 对 tokens 进行分隔转行。
func (tokens items) lines() (ret []items) {
	length := len(tokens)
	var i int
	var token byte
	var line items
	for ; i < length; i++ {
		token = tokens[i]
		if itemNewline != token {
			line = append(line, token)
		} else {
			ret = append(ret, line)
			line = items{}
		}
	}
	if 0 < len(line) {
		ret = append(ret, line)
	}
	return
}

// removeFirst 会将 tokens 中和 token 相等的字符删掉（只删出现相等情况的第一个字符），返回处理后的新串。
func (tokens items) removeFirst(token byte) items {
	length := len(tokens)
	var i int
	for ; i < length; i++ {
		if token == tokens[i] {
			return append(tokens[:i], tokens[i+1:]...)
		}
	}
	return tokens
}

// replaceAll 会将 tokens 中的所有 old 使用 new 替换。
func (tokens items) replaceAll(old, new byte) {
	length := len(tokens)
	var i int
	for ; i < length; i++ {
		if old == tokens[i] {
			tokens[i] = new
		}
	}
}

// replaceNewlineSpace 会将 tokens 中的所有 "\n " 替换为 "\n"。
func (tokens items) replaceNewlineSpace() items {
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

func (tokens items) equal(anotherTokens items) bool {
	if len(tokens) != len(anotherTokens) {
		return false
	}

	for i, token := range tokens {
		if token != anotherTokens[i] {
			return false
		}
	}
	return true
}

func (tokens items) string() string {
	return fromItems(tokens)
}

func (tokens items) trimLeftSpace() (spaces int, remains items) {
	length := len(tokens)
	if 1 > length {
		return 0, tokens
	}

	i := 0
	for ; i < length; i++ {
		if itemSpace == tokens[i] {
			spaces++
		} else if itemTab == tokens[i] {
			spaces += 4
		} else {
			break
		}
	}

	remains = tokens[i:]
	return
}

func (tokens items) trim() (ret items) {
	_, ret = tokens.trimLeft()
	ret = ret.trimRight()

	return
}

func (tokens items) trimLeft() (whitespaces, remains items) {
	length := len(tokens)
	if 1 > length {
		return nil, tokens
	}

	i := 0
	for ; i < length; i++ {
		if !isWhitespace(tokens[i]) {
			break
		} else {
			whitespaces = append(whitespaces, tokens[i])
		}
	}
	return whitespaces, tokens[i:]
}

func (tokens items) trimRight() items {
	length := len(tokens)
	if 1 > length {
		return tokens
	}

	i := length - 1
	for ; 0 <= i; i-- {
		if !isWhitespace(tokens[i]) {
			break
		}
	}
	return tokens[:i+1]
}

func (tokens items) firstNonSpace() (index int, token byte) {
	for index, token = range tokens {
		if itemSpace != token {
			return
		}
	}
	return
}

func (tokens items) accept(token byte) (pos int) {
	length := len(tokens)
	for ; pos < length; pos++ {
		if token != tokens[pos] {
			break
		}
	}
	return
}

func (tokens items) acceptTokenss(someTokenss []items) (pos int) {
	length := len(tokens)
	length2 := len(someTokenss)
	for i := 0; i < length; i++ {
		remains := tokens[i:]
		for j := 0; j < length2; j++ {
			someTokens := someTokenss[j]
			if pos = remains.acceptTokens(someTokens); 0 <= pos {
				return
			}
		}
	}
	return -1
}

func (tokens items) acceptTokens(someTokens items) (pos int) {
	length := len(someTokens)
	for ; pos < length; pos++ {
		if someTokens[pos] != tokens[pos] {
			return -1
		}
	}
	return
}

func (tokens items) containOne(someTokens ...byte) bool {
	for _, t := range tokens {
		for _, it := range someTokens {
			if t == it {
				return true
			}
		}
	}
	return false
}

func (tokens items) contain(anotherToken byte) bool {
	for _, t := range tokens {
		if t == anotherToken {
			return true
		}
	}
	return false
}

func (tokens items) containWhitespace() bool {
	for _, token := range tokens {
		if isWhitespace(token) {
			return true
		}
	}
	return false
}

func (tokens items) isBlankLine() bool {
	for _, token := range tokens {
		if itemSpace != token && itemTab != token && itemNewline != token {
			return false
		}
	}
	return true
}

func (tokens items) splitWhitespace() (ret []items) {
	i := 0
	ret = append(ret, items{})
	lastIsWhitespace := false
	for _, token := range tokens {
		if isWhitespace(token) {
			if !lastIsWhitespace {
				i++
				ret = append(ret, items{})
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
	return token == tokens[0]
}

func (tokens items) endWith(token byte) bool {
	length := len(tokens)
	if 1 > length {
		return false
	}
	return token == tokens[length-1]
}

// isBackslashEscapePunct 判断 tokens 中 pos 所指的值是否是由反斜杠 \ 转义的 ASCII 标点符号。
func (tokens items) isBackslashEscapePunct(pos int) bool {
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

func (tokens items) statWhitespace() (newlines, spaces, tabs int) {
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

func (tokens items) spnl() (ret bool, passed, remains items) {
	passed, remains = tokens.trimLeft()
	newlines, _, _ := passed.statWhitespace()
	if 1 < newlines {
		return false, nil, tokens
	}
	ret = true
	return
}

func (tokens items) peek(pos int) byte {
	if pos < len(tokens) {
		return tokens[pos]
	}
	return itemEnd
}
