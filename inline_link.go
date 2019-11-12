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

func (context *Context) parseInlineLinkDest(tokens items) (passed, remains, destination items) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	passed = make(items, 0, 256)
	destination = make(items, 0, 256)

	isPointyBrackets := itemLess == tokens[1].term()
	if isPointyBrackets {
		matchEnd := false
		passed = append(passed, tokens[0], tokens[1])
		i := 2
		size := 1
		var r rune
		var dest, runes items
		for ; i < length; i += size {
			size = 1
			token := tokens[i]
			if itemNewline == token.term() {
				passed = nil
				return
			}

			if token.term() < utf8.RuneSelf {
				passed = append(passed, token)
				dest = items{token}
			} else {
				dest = items{}
				r, size = utf8.DecodeRune(itemsToBytes(tokens[i:]))
				runes = strToItems(string(r))
				passed = append(passed, runes...)
				dest = append(dest, runes...)
			}
			destination = append(destination, dest...)
			if itemGreater == token.term() && !tokens.isBackslashEscapePunct(i) {
				destination = destination[:len(destination)-1]
				matchEnd = true
				break
			}
		}

		if !matchEnd || (length > i && itemCloseParen != tokens[i+1].term()) {
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
		var dest, runes items
		destStarted := false
		for ; i < length; i += size {
			size = 1
			token := tokens[i]
			if token.term() < utf8.RuneSelf {
				passed = append(passed, token)
				dest = items{token}
			} else {
				dest = items{}
				r, size = utf8.DecodeRune(itemsToBytes(tokens[i:]))
				runes = strToItems(string(r))
				passed = append(passed, runes...)
				dest = append(dest, runes...)
			}
			destination = append(destination, dest...)
			if !destStarted && !isWhitespace(token.term()) && 0 < i {
				destStarted = true
				destination = destination[size:]
				destination = trimWhitespace(destination)
			}
			if destStarted && (isWhitespace(token.term()) || isControl(token.term())) {
				destination = destination[:len(destination)-size]
				passed = passed[:len(passed)-1]
				openParens--
				break
			}
			if itemOpenParen == token.term() && !tokens.isBackslashEscapePunct(i) {
				openParens++
			}
			if itemCloseParen == token.term() && !tokens.isBackslashEscapePunct(i) {
				openParens--
				if 1 > openParens {
					if itemOpenParen == destination[0].term() {
						// TODO: 需要重写边界判断
						destination = destination[1:]
					}
					destination = destination[:len(destination)-1]
					break
				}
			}
		}

		remains = tokens[i:]
		if length > i && (itemCloseParen != tokens[i].term() && itemSpace != tokens[i].term() && itemNewline != tokens[i].term()) {
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
