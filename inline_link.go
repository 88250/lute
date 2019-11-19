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
	"unicode/utf8"
)

func (context *Context) parseInlineLinkDest(tokens []byte) (passed, remains, destination []byte) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	passed = make([]byte, 0, 256)
	destination = make([]byte, 0, 256)

	isPointyBrackets := itemLess == tokens[1]
	if isPointyBrackets {
		matchEnd := false
		passed = append(passed, tokens[0], tokens[1])
		i := 2
		size := 1
		var r rune
		var dest, runes []byte
		for ; i < length; i += size {
			size = 1
			token := tokens[i]
			if itemNewline == token {
				passed = nil
				return
			}

			if token < utf8.RuneSelf {
				passed = append(passed, token)
				dest = []byte{token}
			} else {
				dest = []byte{}
				r, size = utf8.DecodeRune(tokens[i:])
				runes = strToBytes(string(r))
				passed = append(passed, runes...)
				dest = append(dest, runes...)
			}
			destination = append(destination, dest...)
			if itemGreater == token && !isBackslashEscapePunct(tokens, i) {
				destination = destination[:len(destination)-1]
				matchEnd = true
				break
			}
		}

		if !matchEnd || (length > i && itemCloseParen != tokens[i+1]) {
			passed = nil
			return
		}

		passed = append(passed, tokens[i+1])
		remains = tokens[i+2:]
	} else {
		var openParens int
		i := 0
		size := 1
		var r rune
		var dest, runes []byte
		destStarted := false
		for ; i < length; i += size {
			size = 1
			token := tokens[i]
			if token < utf8.RuneSelf {
				passed = append(passed, token)
				dest = []byte{token}
			} else {
				dest = []byte{}
				r, size = utf8.DecodeRune(tokens[i:])
				runes = strToBytes(string(r))
				passed = append(passed, runes...)
				dest = append(dest, runes...)
			}
			destination = append(destination, dest...)
			if !destStarted && !isWhitespace(token) && 0 < i {
				destStarted = true
				destination = destination[size:]
				destination = trimWhitespace(destination)
			}
			if destStarted && (isWhitespace(token) || isControl(token)) {
				destination = destination[:len(destination)-size]
				passed = passed[:len(passed)-1]
				openParens--
				break
			}
			if itemOpenParen == token && !isBackslashEscapePunct(tokens, i) {
				openParens++
			}
			if itemCloseParen == token && !isBackslashEscapePunct(tokens, i) {
				openParens--
				if 1 > openParens {
					if itemOpenParen == destination[0] {
						// TODO: 需要重写边界判断
						destination = destination[1:]
					}
					destination = destination[:len(destination)-1]
					break
				}
			}
		}

		remains = tokens[i:]
		if length > i && (itemCloseParen != tokens[i] && itemSpace != tokens[i] && itemNewline != tokens[i]) {
			passed = nil
			return
		}

		if 0 != openParens {
			passed = nil
			return
		}
	}

	if nil != passed {
		if !context.option.VditorWYSIWYG {
			destination = encodeDestination(unescapeString(destination))
		}
	}
	return
}
