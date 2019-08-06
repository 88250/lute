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
	"strings"
	"unicode"
)

type item rune

func (token item) isNewline() bool {
	return itemNewline == token || '\r' == token /*|| '\u2424' == token || '\u2028' == token*/ || '\u0085' == token || '\u0000' == token
}

func (token item) isWhitespace() bool {
	return itemSpace == token || itemTab == token || itemNewline == token || '\u000A' == token || '\u000C' == token || '\u000D' == token
}

func (token item) isUnicodeWhitespace() bool {
	return unicode.Is(unicode.Zs, rune(token)) || itemTab == token || '\u000D' == token || itemNewline == token || '\u000C' == token
}

func (token item) isDigit() bool {
	return unicode.IsDigit(rune(token))
}

func (token item) isPunct() bool {
	return token.isASCIIPunct() || unicode.IsPunct(rune(token))
}

func (token item) isASCIIPunct() bool {
	return (0x21 <= token && 0x2F >= token) || (0x3A <= token && 0x40 >= token) || (0x5B <= token && 0x60 >= token) || (0x7B <= token && 0x7E >= token)
}

func (token item) isLetter() bool {
	return unicode.IsLetter(rune(token))
}

func (token item) isASCIILetter() bool {
	return ('A' <= token && 'Z' >= token) || ('a' <= token && 'z' >= token)
}

func (token item) isASCIILetterNumHyphen() bool {
	return ('A' <= token && 'Z' >= token) || ('a' <= token && 'z' >= token) || ('0' <= token && '9' >= token) || '-' == token
}

func (token item) isNumber() bool {
	return unicode.IsNumber(rune(token))
}

func (token item) isMark() bool {
	return unicode.IsMark(rune(token))
}

func (token item) isControl() bool {
	return unicode.IsControl(rune(token))
}

func (token item) isSpace() bool {
	return unicode.IsSpace(rune(token))
}

func (token item) isSymbol() bool {
	return unicode.IsSymbol(rune(token))
}

const (
	itemEnd          = item(0)
	itemBacktick     = item('`')
	itemTilde        = item('~')
	itemBang         = item('!')
	itemCrosshatch   = item('#')
	itemAsterisk     = item('*')
	itemOpenParen    = item('(')
	itemCloseParen   = item(')')
	itemHyphen       = item('-')
	itemUnderscore   = item('_')
	itemPlus         = item('+')
	itemEqual        = item('=')
	itemTab          = item('\t')
	itemOpenBracket  = item('[')
	itemCloseBracket = item(']')
	itemDoublequote  = item('"')
	itemSinglequote  = item('\'')
	itemLess         = item('<')
	itemGreater      = item('>')
	itemSpace        = item(' ')
	itemNewline      = item('\n')
	itemBackslash    = item('\\')
	itemSlash        = item('/')
	itemDot          = item('.')
	itemColon        = item(':')
	itemQuestion     = item('?')
	itemAmpersand    = item('&')
	itemSemicolon    = item(';')

	// 以下定义主要是为了方便调试时查看

	item_a = item('a')
	item_b = item('b')
	item_c = item('c')
	item_d = item('d')
	item_e = item('e')
	item_f = item('f')
	item_g = item('g')
	item_h = item('h')
	item_i = item('i')
	item_j = item('j')
	item_k = item('k')
	item_l = item('l')
	item_m = item('m')
	item_n = item('n')
	item_o = item('o')
	item_p = item('p')
	item_q = item('q')
	item_r = item('r')
	item_s = item('s')
	item_t = item('t')
	item_u = item('u')
	item_v = item('v')
	item_w = item('w')
	item_x = item('x')
	item_y = item('y')
	item_z = item('z')

	item_A = item('A')
	item_B = item('B')
	item_C = item('C')
	item_D = item('D')
	item_E = item('E')
	item_F = item('F')
	item_G = item('G')
	item_H = item('H')
	item_I = item('I')
	item_J = item('J')
	item_K = item('K')
	item_L = item('L')
	item_M = item('M')
	item_N = item('N')
	item_O = item('O')
	item_P = item('P')
	item_Q = item('Q')
	item_R = item('R')
	item_S = item('S')
	item_T = item('T')
	item_U = item('U')
	item_V = item('V')
	item_W = item('W')
	item_X = item('X')
	item_Y = item('Y')
	item_Z = item('Z')

	item_0 = item('0')
	item_1 = item('1')
	item_2 = item('2')
	item_3 = item('3')
	item_4 = item('4')
	item_5 = item('5')
	item_6 = item('6')
	item_7 = item('7')
	item_8 = item('8')
	item_9 = item('9')
)

type items []item

// replaceNewlineSpace 会将 tokens 中的所有 "\n " 替换为 "\n"。
func (tokens items) replaceNewlineSpace() items {
	length := len(tokens)
	var token item
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

func (tokens items) rawText() (ret string) {
	b := &strings.Builder{}
	for i := 0; i < len(tokens); i++ {
		b.WriteString(string(tokens[i]))
	}
	ret = b.String()

	return
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
		if !tokens[i].isWhitespace() {
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
		if !tokens[i].isWhitespace() && itemEnd != tokens[i] {
			break
		}
	}

	return tokens[:i+1]
}

func (tokens items) firstNonSpace() (index int, token item) {
	for index, token = range tokens {
		if itemSpace != token {
			return
		}
	}

	return
}

func (tokens items) accept(token item) (pos int) {
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

func (tokens items) contain(someTokens ...item) bool {
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
		if token.isWhitespace() {
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
		if token.isWhitespace() {
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

func (tokens items) split(token item) (ret []items) {
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

func (tokens items) startWith(token item) bool {
	if 1 > len(tokens) {
		return false
	}

	return token == tokens[0]
}

func (tokens items) endWith(token item) bool {
	length := len(tokens)
	if 1 > length {
		return false
	}

	return token == tokens[length-1]
}

func (tokens items) isBackslashEscape(pos int) bool {
	if !tokens[pos].isASCIIPunct() {
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

func (tokens items) peek(pos int) item {
	if pos < len(tokens) {
		return tokens[pos]
	}

	return itemEnd
}
