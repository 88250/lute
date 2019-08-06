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
	"unicode"
)

func isWhitespace(token byte) bool {
	return itemSpace == token || itemNewline /* '\u000A' */ == token || itemTab == token || '\u000B' == token || '\u000C' == token || '\u000D' == token
}

func isUnicodeWhitespace(token byte) bool {
	if unicode.Is(unicode.Zs, rune(token)) {
		return true
	}

	return itemTab == token || '\u000D' == token || itemNewline == token || '\u000C' == token
}

func isDigit(token byte) bool {
	return unicode.IsDigit(rune(token))
}

func isPunct(token byte) bool {
	return isASCIIPunct(token) || unicode.IsPunct(rune(token))
}

func isASCIIPunct(token byte) bool {
	return (0x21 <= token && 0x2F >= token) || (0x3A <= token && 0x40 >= token) || (0x5B <= token && 0x60 >= token) || (0x7B <= token && 0x7E >= token)
}

func isASCIILetter(token byte) bool {
	return ('A' <= token && 'Z' >= token) || ('a' <= token && 'z' >= token)
}

func isASCIILetterNumHyphen(token byte) bool {
	return ('A' <= token && 'Z' >= token) || ('a' <= token && 'z' >= token) || ('0' <= token && '9' >= token) || '-' == token
}

func isControl(token byte) bool {
	return unicode.IsControl(rune(token))
}

const (
	itemEnd          = byte(0)
	itemBacktick     = byte('`')
	itemTilde        = byte('~')
	itemBang         = byte('!')
	itemCrosshatch   = byte('#')
	itemAsterisk     = byte('*')
	itemOpenParen    = byte('(')
	itemCloseParen   = byte(')')
	itemHyphen       = byte('-')
	itemUnderscore   = byte('_')
	itemPlus         = byte('+')
	itemEqual        = byte('=')
	itemTab          = byte('\t')
	itemOpenBracket  = byte('[')
	itemCloseBracket = byte(']')
	itemDoublequote  = byte('"')
	itemSinglequote  = byte('\'')
	itemLess         = byte('<')
	itemGreater      = byte('>')
	itemSpace        = byte(' ')
	itemNewline      = byte('\n')
	itemBackslash    = byte('\\')
	itemSlash        = byte('/')
	itemDot          = byte('.')
	itemColon        = byte(':')
	itemQuestion     = byte('?')
	itemAmpersand    = byte('&')
	itemSemicolon    = byte(';')
)

// items 定义了字节数组，每个字节是一个 token。
type items []byte

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
	size := len(tokens)
	if 1 > size {
		return 0, tokens
	}

	i := 0
	for ; i < size; i++ {
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
	size := len(tokens)
	if 1 > size {
		return nil, tokens
	}

	i := 0
	for ; i < size; i++ {
		if !isWhitespace(tokens[i]) {
			break
		} else {
			whitespaces = append(whitespaces, tokens[i])
		}
	}

	return whitespaces, tokens[i:]
}

func (tokens items) trimRight() items {
	size := len(tokens)
	if 1 > size {
		return tokens
	}

	i := size - 1
	for ; 0 <= i; i-- {
		if !isWhitespace(tokens[i]) && itemEnd != tokens[i] {
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

func (tokens items) contain(someTokens ...byte) bool {
	for _, t := range tokens {
		for _, it := range someTokens {
			if t == it {
				return true
			}
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
	ret = []items{}
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

func (tokens items) split(token byte) (ret []items) {
	ret = []items{}
	i := 0
	ret = append(ret, items{})
	for j, t := range tokens {
		if token == t {
			ret = append(ret, items{})
			ret[i+1] = append(ret[i+1], tokens[j+1:]...)
			return
		} else {
			ret[i] = append(ret[i], t)
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

func (tokens items) isBackslashEscape(pos int) bool {
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
