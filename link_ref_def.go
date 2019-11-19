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

func (context *Context) parseLinkRefDef(tokens []byte) []byte {
	_, tokens = trimLeft(tokens)
	if 1 > len(tokens) {
		return nil
	}

	n, remains, label := context.parseLinkLabel(tokens)
	if 2 > n || 1 > len(label) {
		return nil
	}

	length := len(remains)
	if 1 > length {
		return nil
	}

	if ':' != remains[0] {
		return nil
	}

	remains = remains[1:]
	whitespaces, remains := trimLeft(remains)
	newlines, _, _ := statWhitespace(whitespaces)
	if 1 < newlines {
		return nil
	}

	tokens = remains
	linkDest, remains, destination := context.parseLinkDest(tokens)
	if nil == linkDest {
		return nil
	}

	whitespaces, remains = trimLeft(remains)
	if nil == whitespaces && 0 < len(remains) {
		return nil
	}
	newlines, spaces1, tabs1 := statWhitespace(whitespaces)
	if 1 < newlines {
		return nil
	}

	_, tokens = trimLeft(remains)
	validTitle, _, remains, title := context.parseLinkTitle(tokens)
	if !validTitle && 1 > newlines {
		return nil
	}
	if 0 < spaces1+tabs1 && !isBlankLine(remains) && itemNewline != remains[0] {
		return nil
	}

	titleLine := tokens
	whitespaces, tokens = trimLeft(remains)
	_, spaces2, tabs2 := statWhitespace(whitespaces)
	if !isBlankLine(tokens) && 0 < spaces2+tabs2 {
		remains = titleLine
	} else {
		remains = tokens
	}

	link := context.tree.newLink(NodeLink, label, destination, title, 1)
	lowerCaseLabel := bytes.ToLower(label)
	if _, ok := context.linkRefDef[bytesToStr(lowerCaseLabel)]; !ok {
		context.linkRefDef[bytesToStr(lowerCaseLabel)] = link
	}
	return remains
}

func (context *Context) parseLinkTitle(tokens []byte) (validTitle bool, passed, remains, title []byte) {
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
		if !context.option.VditorWYSIWYG {
			title = unescapeString(title)
		}
	}

	return
}

func (context *Context) parseLinkTitleMatch(opener, closer byte, tokens []byte) (validTitle bool, passed, remains, title []byte) {
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
		title = append(title, strToBytes(string(r))...)
		if closer == token && !isBackslashEscapePunct(tokens, i) {
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

func (context *Context) parseLinkDest(tokens []byte) (ret, remains, destination []byte) {
	ret, remains, destination = context.parseLinkDest1(tokens) // <autolink>
	if nil == ret {
		ret, remains, destination = context.parseLinkDest2(tokens) // [label](/url)
	}
	if nil != ret {
		if !context.option.VditorWYSIWYG {
			destination = encodeDestination(unescapeString(destination))
		}
	}
	return
}

func (context *Context) parseLinkDest2(tokens []byte) (ret, remains, destination []byte) {
	remains = tokens
	length := len(tokens)
	if 1 > length {
		return
	}

	ret = make([]byte, 0, 256)
	destination = make([]byte, 0, 256)

	var openParens int
	i := 0
	size := 0
	var r rune
	for i < length {
		token := tokens[i]
		ret = append(ret, token)
		r, size = utf8.DecodeRune(tokens[i:])
		for j := 1; j < size; j++ {
			ret = append(ret, tokens[i+j])
		}
		destination = append(destination, strToBytes(string(r))...)
		if isWhitespace(token) || isControl(token) {
			destination = destination[:len(destination)-1]
			ret = ret[:len(ret)-1]
			break
		}

		if itemOpenParen == token && !isBackslashEscapePunct(tokens, i) {
			openParens++
		}
		if itemCloseParen == token && !isBackslashEscapePunct(tokens, i) {
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

func (context *Context) parseLinkDest1(tokens []byte) (ret, remains, destination []byte) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemLess != tokens[0] {
		return
	}

	ret = make([]byte, 0, 256)
	destination = make([]byte, 0, 256)

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
			destination = append(destination, strToBytes(string(r))...)
			if itemLess == token && !isBackslashEscapePunct(tokens, i) {
				ret = nil
				return
			}
		}

		if itemGreater == token && !isBackslashEscapePunct(tokens, i) {
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

func (context *Context) parseLinkLabel(tokens []byte) (n int, remains, label []byte) {
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemOpenBracket != tokens[0] {
		return
	}

	passed := make([]byte, 0, len(tokens))
	passed = append(passed, tokens[0])

	closed := false
	i := 1
	for i < length {
		token := tokens[i]
		passed = append(passed, token)
		r, size := utf8.DecodeRune(tokens[i:])
		for j := 1; j < size; j++ {
			passed = append(passed, tokens[i+j])
		}
		label = append(label, strToBytes(string(r))...)
		if itemCloseBracket == token && !isBackslashEscapePunct(tokens, i) {
			closed = true
			label = label[0 : len(label)-1]
			remains = tokens[i+1:]
			break
		}
		if itemOpenBracket == token && !isBackslashEscapePunct(tokens, i) {
			passed = nil
			return
		}
		i += size
	}

	if !closed || nil == trimWhitespace(label) || 999 < len(label) {
		passed = nil
		return
	}

	label = trimWhitespace(label)
	if !context.option.VditorWYSIWYG {
		label = replaceAll(label, itemNewline, itemSpace)
		length := len(label)
		var token byte
		for i := 0; i < length; i++ {
			token = label[i]
			if token == itemSpace && i < length-1 && label[i+1] == itemSpace {
				label = append(label[:i], label[i+1:]...)
				length--
			}
		}
	}
	n = len(passed)
	return
}
