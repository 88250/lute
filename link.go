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
	"unicode/utf8"
)

// Link 描述了链接节点结构。
type Link struct {
	*BaseNode
	Destination string
	Title       string
}

func (t *Tree) parseAutoEmailLink(tokens items) (ret Node) {
	tokens = tokens[1:]
	var dest string
	var token byte
	length := len(tokens)
	passed := 0
	i := 0
	at := false
	for ; i < length; i++ {
		token = tokens[i]
		dest += string(token)
		passed++
		if '@' == token {
			at = true
			break
		}

		if !isASCIILetterNumHyphen(token) && !strings.Contains(".!#$%&'*+/=?^_`{|}~", string(token)) {
			return nil
		}
	}

	if 1 > i || !at {
		return nil
	}

	domainPart := tokens[i+1:]
	length = len(domainPart)
	i = 0
	closed := false
	for ; i < length; i++ {
		token = domainPart[i]
		passed++
		if itemGreater == token {
			closed = true
			break
		}
		dest += string(token)
		if !isASCIILetterNumHyphen(token) && itemDot != token {
			return nil
		}
		if 63 < i {
			return nil
		}
	}

	if 1 > i || !closed {
		return nil
	}

	t.context.pos += passed + 1
	ret = &Link{&BaseNode{typ: NodeLink}, "mailto:" + dest, ""}
	ret.AppendChild(ret, &Text{typ: NodeText, value: items(dest)})

	return
}

func (t *Tree) parseAutolink(tokens items) (ret Node) {
	schemed := false
	scheme := ""
	dest := ""
	var token byte
	i := t.context.pos + 1
	for ; i < len(tokens) && itemGreater != tokens[i]; i++ {
		token = tokens[i]
		if itemSpace == token {
			return nil
		}

		dest += string(token)
		if !schemed {
			if itemColon != token {
				scheme += string(token)
			} else {
				schemed = true
			}
		}
	}
	if !schemed || 3 > len(scheme) {
		return nil
	}

	ret = &Link{&BaseNode{typ: NodeLink}, encodeDestination(dest), ""}
	if itemGreater != tokens[i] {
		return nil
	}

	t.context.pos = 1 + i
	ret.AppendChild(ret, &Text{typ: NodeText, value: items(dest)})

	return
}

func (context *Context) parseInlineLinkDest(tokens items) (passed, remains items, destination string) {
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
				destination = ""
				return
			}

			passed = append(passed, token)
			r, size = utf8.DecodeRune(tokens[i:])
			if 1 < size {
				for j := 1; j < size; j++ {
					passed = append(passed, tokens[i+j])
				}
			}
			destination += string(r)
			if itemGreater == token && !tokens.isBackslashEscape(i) {
				destination = destination[:len(destination)-1]
				matchEnd = true
				break
			}
		}

		if !matchEnd || (length > i && itemCloseParen != tokens[i+1]) {
			passed = nil
			destination = ""
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
			destination += string(r)
			for j := 1; j < size; j++ {
				passed = append(passed, tokens[i+j])
			}

			if !destStarted && !isWhitespace(token) && 0 < i {
				destStarted = true
				destination = destination[size:]
				destination = strings.TrimSpace(destination)
			}
			if destStarted && (isWhitespace(token) || isControl(token)) {
				destination = destination[:len(destination)-size]
				passed = passed[:len(passed)-1]
				break
			}
			if itemOpenParen == token && !tokens.isBackslashEscape(i) {
				openParens++
			}
			if itemCloseParen == token && !tokens.isBackslashEscape(i) {
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
			destination = ""
			return
		}
	}

	if nil != passed {
		destination = encodeDestination(unescapeString(destination))
	}

	return
}
