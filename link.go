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

import "strings"

type Link struct {
	*BaseNode
	Destination string
	Title       string
}

func (context *Context) parseInlineLinkDest(tokens items) (passed, remains items, destination string) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemOpenParen != tokens[0].typ {
		return
	}

	isPointyBrackets := itemLess == tokens[1].typ
	if isPointyBrackets {
		matchEnd := false
		passed = append(passed, tokens[0], tokens[1])
		i := 2
		for ; i < length; i++ {
			token := tokens[i]
			if token.isNewline() {
				passed = nil
				destination = ""
				return
			}

			passed = append(passed, token)
			destination += token.val
			if itemGreater == token.typ && !tokens.isBackslashEscape(i) {
				destination = destination[:len(destination)-1]
				matchEnd = true
				break
			}
		}

		if !matchEnd || (length > i && itemCloseParen != tokens[i+1].typ) {
			passed = nil
			destination = ""
			return
		}

		passed = append(passed, tokens[i+1])

		remains = tokens[i+2:]
	} else {
		var openParens int
		i := 0
		destStarted := false
		for ; i < length; i++ {
			token := tokens[i]
			passed = append(passed, token)
			destination += token.val
			if !destStarted && !token.isWhitespace() && 0 < i {
				destStarted = true
				destination = destination[1:]
				destination = strings.TrimSpace(destination)
			}
			if destStarted && (token.isWhitespace() || token.isControl()) {
				destination = destination[:len(destination)-1]
				passed = passed[:len(passed)-1]
				break
			}
			if itemOpenParen == token.typ && !tokens.isBackslashEscape(i) {
				openParens++
			}
			if itemCloseParen == token.typ && !tokens.isBackslashEscape(i) {
				openParens--
				if 1 > openParens {
					destination = destination[:len(destination)-1]
					break
				}
			}
		}

		remains = tokens[i:]
		if length > i && (itemCloseParen != tokens[i].typ && itemSpace != tokens[i].typ && itemNewline != tokens[i].typ) {
			passed = nil
			destination = ""
			return
		}
	}

	if nil != passed {
		destination = encodeDestination(unescapeString(destination))
	}

	return
}
