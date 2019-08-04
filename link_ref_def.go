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
)

func (context *Context) parseLinkRefDef(tokens items) items {
	_, tokens = tokens.trimLeft()
	if 1 > len(tokens) {
		return nil
	}

	linkLabel, remains, label := context.parseLinkLabel(tokens)
	if nil == linkLabel {
		return nil
	}

	if 1 > len(remains) || itemColon != remains[0].typ {
		return nil
	}

	remains = remains[1:]
	whitespaces, remains := remains.trimLeft()
	newlines, _, _ := whitespaces.statWhitespace()
	if 1 < newlines {
		return nil
	}

	tokens = remains
	linkDest, remains, destination := context.parseLinkDest(tokens)
	if nil == linkDest {
		return nil
	}

	whitespaces, remains = remains.trimLeft()
	if nil == whitespaces && 0 < len(remains) {
		return nil
	}
	newlines, spaces1, tabs1 := whitespaces.statWhitespace()
	if 1 < newlines {
		return nil
	}

	_, tokens = remains.trimLeft()
	validTitle, _, remains, title := context.parseLinkTitle(tokens)
	if !validTitle && 1 > newlines {
		return nil
	}
	if 0 < spaces1+tabs1 && !remains.isBlankLine() && itemNewline != remains[0].typ {
		return nil
	}

	titleLine := tokens
	whitespaces, tokens = remains.trimLeft()
	_, spaces2, tabs2 := whitespaces.statWhitespace()
	if !tokens.isBlankLine() && 0 < spaces2+tabs2 {
		title = ""
		remains = titleLine
	} else {
		remains = tokens
	}

	link := &Link{&BaseNode{typ: NodeLink}, destination, ""}
	lowerCaseLabel := strings.ToLower(label)
	link.Title = title
	if _, ok := context.linkRefDef[lowerCaseLabel]; !ok {
		context.linkRefDef[lowerCaseLabel] = link
	}

	return remains
}

func (context *Context) parseLinkTitle(tokens items) (validTitle bool, passed, remains items, title string) {
	if 1 > len(tokens) {
		return true, nil, tokens, ""
	}
	if itemOpenBracket == tokens[0].typ {
		return true, nil, tokens, ""
	}

	validTitle, passed, remains, title = context.parseLinkTitleMatch(itemDoublequote, itemDoublequote, tokens)
	if !validTitle {
		validTitle, passed, remains, title = context.parseLinkTitleMatch(itemSinglequote, itemSinglequote, tokens)
		if !validTitle {
			validTitle, passed, remains, title = context.parseLinkTitleMatch(itemOpenParen, itemCloseParen, tokens)
		}
	}
	if "" != title {
		title = unescapeString(title)
	}

	return
}

func (context *Context) parseLinkTitleMatch(opener, closer itemType, tokens items) (validTitle bool, passed, remains items, title string) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if opener != tokens[0].typ {
		return
	}

	line := tokens
	closed := false
	i := 1
	for i < len(line) {
		token := line[i]
		title += token.Value()
		passed = append(passed, token)
		if closer == token.typ && !tokens.isBackslashEscape(i) {
			closed = true
			title = title[:len(title)-1]
			break
		}
		i++
	}

	if !closed {
		title = ""
		passed = nil
		return
	}

	validTitle = true
	remains = tokens[i+1:]

	return
}

func (context *Context) parseLinkDest(tokens items) (ret, remains items, destination string) {
	ret, remains, destination = context.parseLinkDest1(tokens) // <autolink>
	if nil == ret {
		ret, remains, destination = context.parseLinkDest2(tokens) // [label](/url)
	}
	if nil != ret {
		destination = encodeDestination(unescapeString(destination))
	}

	return
}

func (context *Context) parseLinkDest2(tokens items) (ret, remains items, destination string) {
	remains = tokens
	length := len(tokens)
	if 1 > length {
		return
	}

	var openParens int
	i := 0
	for ; i < length; i++ {
		token := tokens[i]
		ret = append(ret, token)
		destination += token.Value()
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
				i++
				break
			}
		}
	}

	remains = tokens[i:]
	if length > i && !tokens[i].isWhitespace() {
		ret = nil
		destination = ""
		return
	}

	return
}

func (context *Context) parseLinkDest1(tokens items) (ret, remains items, destination string) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemLess != tokens[0].typ {
		return
	}

	closed := false
	i := 0
	for ; i < length; i++ {
		token := tokens[i]
		ret = append(ret, token)
		if 0 < i {
			destination += token.Value()
			if itemLess == token.typ && !tokens.isBackslashEscape(i) {
				ret = nil
				destination = ""
				return
			}
		}

		if itemGreater == token.typ && !tokens.isBackslashEscape(i) {
			closed = true
			destination = destination[0 : len(destination)-1]
			break
		}
	}

	if !closed {
		ret = nil
		destination = ""

		return
	}

	remains = tokens[i+1:]

	return
}

func (context *Context) parseLinkLabel(tokens items) (passed, remains items, label string) {
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemOpenBracket != tokens[0].typ {
		return
	}

	line := tokens
	closed := false
	i := 1
	for {
		token := line[i]
		passed = append(passed, token)
		label += token.Value()
		if itemCloseBracket == token.typ && !tokens.isBackslashEscape(i) {
			closed = true
			label = label[0 : len(label)-1]
			remains = line[i+1:]
			break
		}
		if itemOpenBracket == token.typ && !tokens.isBackslashEscape(i) {
			passed = nil
			label = ""
			return
		}
		i++
	}

	if !closed || "" == strings.TrimSpace(label) || 999 < len(label) {
		passed = nil
	}

	label = strings.TrimSpace(label)
	label = strings.ReplaceAll(label, "\n", " ")
	for 0 <= strings.Index(label, "  ") {
		label = strings.ReplaceAll(label, "  ", " ")
	}

	return
}
