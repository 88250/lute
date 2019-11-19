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
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/b3log/lute/html"
)

func unescapeString(tokens []byte) (ret []byte) {
	if nil == tokens {
		return
	}

	tokens = strToBytes(htmlUnescapeString(bytesToStr(tokens)))
	length := len(tokens)
	ret = make([]byte, 0, length)
	for i := 0; i < length; i++ {
		if isBackslashEscapePunct(tokens, i) {
			ret = ret[:len(ret)-1]
		}
		ret = append(ret, tokens[i])
	}
	return
}

func htmlUnescapeString(s string) string {
	// 鸣谢 https://gitlab.com/golang-commonmark

	i := strings.IndexByte(s, itemAmpersand)
	if i < 0 {
		return s
	}

	anyChanges := false
	var entityStr string
	var entityLen int
	for i < len(s) {
		if s[i] == '&' {
			entityStr, entityLen = parseEntity(s[i:])
			if entityLen > 0 {
				anyChanges = true
				break
			}
		}
		i++
	}

	if !anyChanges {
		return s
	}

	buf := make([]byte, len(s)-entityLen+len(entityStr))
	copy(buf[:i], s)
	n := copy(buf[i:], entityStr)
	j := i + n
	i += entityLen
	for i < len(s) {
		b := s[i]
		if b == '&' {
			entityStr, entityLen = parseEntity(s[i:])
			if entityLen > 0 {
				n = copy(buf[j:], entityStr)
				j += n
				i += entityLen
				continue
			}
		}

		buf[j] = b
		j++
		i++
	}

	return string(buf[:j])
}

func parseEntity(s string) (string, int) {
	st := 0
	var n int

	for i := 1; i < len(s); i++ {
		b := s[i]

		switch st {
		case 0: // initial state
			switch {
			case b == '#':
				st = 1
			case isASCIILetter(b):
				n = 1
				st = 2
			default:
				return "", 0
			}

		case 1: // &#
			switch {
			case b == 'x' || b == 'X':
				st = 3
			case isDigit(b):
				n = 1
				st = 4
			default:
				return "", 0
			}

		case 2: // &q
			switch {
			case isASCIILetterNum(b):
				n++
				if n > 31 {
					return "", 0
				}
			case b == ';':
				if e, ok := html.Entities[s[i-n:i]]; ok {
					return e, i + 1
				}
				return "", 0
			default:
				return "", 0
			}

		case 3: // &#x
			switch {
			case isHexDigit(b):
				n = 1
				st = 5
			default:
				return "", 0
			}

		case 4: // &#0
			switch {
			case isDigit(b):
				n++
				if n > 8 {
					return "", 0
				}
			case b == ';':
				c, _ := strconv.ParseInt(s[i-n:i], 10, 32)
				if !isValidEntityCode(c) {
					return BadEntity, i + 1
				}
				return string(rune(c)), i + 1
			default:
				return "", 0
			}

		case 5: // &#x0
			switch {
			case isHexDigit(b):
				n++
				if n > 8 {
					return "", 0
				}
			case b == ';':
				c, err := strconv.ParseInt(s[i-n:i], 16, 32)
				if nil != err {
					return BadEntity, i + 1
				}
				if !isValidEntityCode(c) {
					return BadEntity, i + 1
				}
				return string(rune(c)), i + 1
			default:
				return "", 0
			}
		}
	}
	return "", 0
}

const BadEntity = string(utf8.RuneError)

func isValidEntityCode(c int64) bool {
	switch {
	case !utf8.ValidRune(rune(c)):
		return false

	// never used
	case c >= 0xfdd0 && c <= 0xfdef:
		return false
	case c&0xffff == 0xffff || c&0xffff == 0xfffe:
		return false
	// control codes
	case c >= 0x00 && c <= 0x08:
		return false
	case c == 0x0b:
		return false
	case c >= 0x0e && c <= 0x1f:
		return false
	case c >= 0x7f && c <= 0x9f:
		return false
	}

	return true
}
