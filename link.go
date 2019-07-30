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

type Link struct {
	*BaseNode
	Destination string
	Title       string
}

func (context *Context) parseInlineLinkDest(tokens items) (ret, remains items, destination string) {
	remains = tokens
	length := len(tokens)
	if 3 > length {
		return
	}

	if itemOpenParen != tokens[0].typ {
		return
	}

	var openParens int
	i := 0
	for ; i < length; i++ {
		token := tokens[i]
		ret = append(ret, token)
		destination += token.val
		if token.isWhitespace() || token.isControl() {
			destination = destination[:len(destination)-1]
			ret = ret[:len(ret)-1]
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
	if length > i && (itemCloseParen != tokens[i].typ && itemSpace != tokens[i].typ) {
		ret = nil
		destination = ""
		return
	}

	destination = destination[1:]

	if nil != ret {
		destination = encodeDestination(unescapeString(destination))
	}

	return
}
