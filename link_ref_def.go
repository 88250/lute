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
	"strings"
	"unicode/utf8"
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

	if 1 > len(remains) || itemColon != remains[0] {
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
	if 0 < spaces1+tabs1 && !remains.isBlankLine() && itemNewline != remains[0] {
		return nil
	}

	titleLine := tokens
	whitespaces, tokens = remains.trimLeft()
	_, spaces2, tabs2 := whitespaces.statWhitespace()
	if !tokens.isBlankLine() && 0 < spaces2+tabs2 {
		remains = titleLine
	} else {
		remains = tokens
	}

	link := &Link{&BaseNode{typ: NodeLink}, destination, nil}
	lowerCaseLabel := strings.ToLower(label)
	link.Title = title
	if _, ok := context.linkRefDef[lowerCaseLabel]; !ok {
		context.linkRefDef[lowerCaseLabel] = link
	}

	return remains
}

func (context *Context) parseLinkTitle(tokens items) (validTitle bool, passed, remains, title items) {
	if 1 > len(tokens) {
		return true, nil, tokens, nil
	}
	if itemOpenBracket == tokens[0] {
		return true, nil, tokens, nil
	}

	validTitle, passed, remains, title = context.parseLinkTitleMatch(itemDoublequote, itemDoublequote, tokens)
	if !validTitle {
		validTitle, passed, remains, title = context.parseLinkTitleMatch(itemSinglequote, itemSinglequote, tokens)
		if !validTitle {
			validTitle, passed, remains, title = context.parseLinkTitleMatch(itemOpenParen, itemCloseParen, tokens)
		}
	}
	if nil != title {
		title = unescapeString(title)
	}

	return
}

func (context *Context) parseLinkTitleMatch(opener, closer byte, tokens items) (validTitle bool, passed, remains, title items) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if opener != tokens[0] {
		return
	}

	line := tokens
	length = len(line)
	closed := false
	i := 1
	size := 0
	var r rune
	for ; i < length; i += size {
		token := line[i]
		passed = append(passed, token)
		r, size = utf8.DecodeRune(line[i:])
		for j := 1; j < size; j++ {
			passed = append(passed, tokens[i+j])
		}
		title = append(title, items(string(r))...)
		if closer == token && !tokens.isBackslashEscapePunct(i) {
			closed = true
			title = title[:len(title)-1]
			break
		}
	}

	if !closed {
		passed = nil
		return
	}

	validTitle = true
	remains = tokens[i+1:]
	return
}

func (context *Context) parseLinkDest(tokens items) (ret, remains, destination items) {
	ret, remains, destination = context.parseLinkDest1(tokens) // <autolink>
	if nil == ret {
		ret, remains, destination = context.parseLinkDest2(tokens) // [label](/url)
	}
	if nil != ret {
		destination = encodeDestination(unescapeString(destination))
	}

	return
}

func (context *Context) parseLinkDest2(tokens items) (ret, remains, destination items) {
	remains = tokens
	length := len(tokens)
	if 1 > length {
		return
	}

	var openParens int
	i := 0
	size := 0
	var r rune
	for ; i < length; {
		token := tokens[i]
		ret = append(ret, token)
		r, size = utf8.DecodeRune(tokens[i:])
		for j := 1; j < size; j++ {
			ret = append(ret, tokens[i+j])
		}
		destination = append(destination, items(string(r))...)
		if isWhitespace(token) || isControl(token) {
			destination = destination[:len(destination)-1]
			ret = ret[:len(ret)-1]
			break
		}

		if itemOpenParen == token && !tokens.isBackslashEscapePunct(i) {
			openParens++
		}
		if itemCloseParen == token && !tokens.isBackslashEscapePunct(i) {
			openParens--
			if 1 > openParens {
				i++
				break
			}
		}

		i += size
	}

	remains = tokens[i:]
	if length > i && !isWhitespace(tokens[i]) {
		ret = nil
		return
	}

	return
}

func (context *Context) parseLinkDest1(tokens items) (ret, remains, destination items) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemLess != tokens[0] {
		return
	}

	closed := false
	i := 0
	size := 0
	var r rune
	for ; i < length; i += size {
		token := tokens[i]
		ret = append(ret, token)
		size = 1
		if 0 < i {
			r, size = utf8.DecodeRune(tokens[i:])
			for j := 1; j < size; j++ {
				ret = append(ret, tokens[i+j])
			}
			destination = append(destination, items(string(r))...)
			if itemLess == token && !tokens.isBackslashEscapePunct(i) {
				ret = nil
				return
			}
		}

		if itemGreater == token && !tokens.isBackslashEscapePunct(i) {
			closed = true
			destination = destination[0 : len(destination)-1]
			break
		}
	}

	if !closed {
		ret = nil
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

	if itemOpenBracket != tokens[0] {
		return
	}

	passed = make(items, 0, len(tokens))

	closed := false
	i := 1
	for i < length {
		token := tokens[i]
		passed = append(passed, token)
		r, size := utf8.DecodeRune(tokens[i:])
		for j := 1; j < size; j++ {
			passed = append(passed, tokens[i+j])
		}
		label += string(r)
		if itemCloseBracket == token && !tokens.isBackslashEscapePunct(i) {
			closed = true
			label = label[0 : len(label)-1]
			remains = tokens[i+1:]
			break
		}
		if itemOpenBracket == token && !tokens.isBackslashEscapePunct(i) {
			passed = nil
			label = ""
			return
		}
		i += size
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
