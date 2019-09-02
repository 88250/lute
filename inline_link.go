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
	"unicode/utf8"
)

func (context *Context) parseInlineLink(tokens items) (passed, remains, destination items) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemOpenParen != tokens[0] {
		return
	}

	isPointyBrackets := itemLess == tokens[1]
	if isPointyBrackets {
		matchEnd := false
		passed = append(passed, tokens[0], tokens[1])
		i := 2
		size := 0
		var r rune
		for ; i < length; i += size {
			token := tokens[i]
			if itemNewline == token {
				passed = nil
				return
			}

			passed = append(passed, token)
			r, size = utf8.DecodeRune(tokens[i:])
			if 1 < size {
				for j := 1; j < size; j++ {
					passed = append(passed, tokens[i+j])
				}
			}
			destination = append(destination, toItems(string(r))...)
			if itemGreater == token && !tokens.isBackslashEscapePunct(i) {
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
		size := 0
		var r rune
		destStarted := false
		for ; i < length; i += size {
			token := tokens[i]
			passed = append(passed, token)
			r, size = utf8.DecodeRune(tokens[i:])
			destination = append(destination, []byte(string(r))...)
			for j := 1; j < size; j++ {
				passed = append(passed, tokens[i+j])
			}

			if !destStarted && !isWhitespace(token) && 0 < i {
				destStarted = true
				destination = destination[size:]
				destination = bytes.TrimSpace(destination)
			}
			if destStarted && (isWhitespace(token) || isControl(token)) {
				destination = destination[:len(destination)-size]
				passed = passed[:len(passed)-1]
				break
			}
			if itemOpenParen == token && !tokens.isBackslashEscapePunct(i) {
				openParens++
			}
			if itemCloseParen == token && !tokens.isBackslashEscapePunct(i) {
				openParens--
				if 1 > openParens {
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
	}

	if nil != passed {
		destination = encodeDestination(unescapeString(destination))
	}

	return
}
