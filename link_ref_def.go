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

func (context *Context) parseLinkRefDef(line items) items {
	_, line = line.trimLeft()
	if 1 > len(line) {
		return nil
	}

	linkLabel, remains, label := context.parseLinkLabel(line)
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

	tokens := remains
	linkDest, remains, destination := context.parseLinkDest(tokens)
	if nil == linkDest {
		return nil
	}

	link := &Link{&BaseNode{typ: NodeLink}, destination, ""}

	whitespaces, remains = remains.trimLeft()
	if nil == whitespaces {
		return nil
	}
	newlines, _, _ = whitespaces.statWhitespace()
	if 1 < newlines {
		return nil
	}

	lowerCaseLabel := strings.ToLower(label)

	tokens = remains
	validTitle, remains, title := context.parseLinkTitle(tokens)
	if !validTitle {
		return nil
	}

	link.Title = title
	if _, ok := context.linkRefDef[lowerCaseLabel]; !ok {
		context.linkRefDef[lowerCaseLabel] = link
	}

	return remains
}

func (context *Context) parseLinkText(tokens items) (ret, remains items, text string) {

	return
}

func (context *Context) parseLinkTitle(tokens items) (validTitle bool, remains items, title string) {
	validTitle, remains, title = context.parseLinkTitleMatch(itemDoublequote, itemDoublequote, tokens)
	if !validTitle {
		validTitle, remains, title = context.parseLinkTitleMatch(itemSinglequote, itemSinglequote, tokens)
		if !validTitle{
			validTitle, remains, title = context.parseLinkTitleMatch(itemOpenParen, itemCloseParen, tokens)
		}
	}
	if "" != title {
		title = unescapeString(title)
	}

	return
}

func (context *Context) parseLinkTitleMatch(opener, closer itemType, tokens items) (validTitle bool, remains items, title string) {
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
	for i < len(line){
		token := line[i]
		title += token.val
		if closer == token.typ && !tokens.isBackslashEscape(i) {
			closed = true
			title = title[:len(title)-1]
			break
		}
		i++
	}

	if !closed {
		title = ""
		return
	}

	validTitle = true
	remains = tokens[i+1:]

	return
}

func (context *Context) parseLinkDest(tokens items) (ret, remains items, destination string) {
	ret, remains, destination = context.parseLinkDest1(tokens)
	if nil == ret {
		ret, remains, destination = context.parseLinkDest2(tokens)
	}
	if nil != ret {
		destination = encodeDestination(destination)
	}

	return
}

func (context *Context) parseLinkDest2(tokens items) (ret, remains items, destination string) {
	remains = tokens
	var leftParens, rightParens int
	i := 0
	length := len(tokens)
	for ; i < length; i++ {
		token := tokens[i]
		ret = append(ret, token)
		destination += token.val
		if itemSpace == token.typ || token.isControl() {
			destination = destination[:len(destination)-1]
			ret = ret[:len(ret)-1]
			break
		}

		if itemOpenParen == token.typ && !tokens.isBackslashEscape(i) {
			leftParens++
		}
		if itemCloseParen == token.typ && !tokens.isBackslashEscape(i) {
			rightParens++
		}
	}

	if leftParens != rightParens {
		ret = nil
		destination = ""
		return
	}

	if length <= i {
		i = length - 1
		destination = destination[:len(destination)-1]
	}

	remains = tokens[i:]

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
			destination += token.val
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

func (context *Context) parseLinkLabel(tokens items) (ret, remains items, label string) {
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
		ret = append(ret, token)
		label += token.val
		if itemCloseBracket == token.typ && !tokens.isBackslashEscape(i) {
			closed = true
			label = label[0 : len(label)-1]
			remains = line[i+1:]
			break
		}
		i++
	}

	if !closed || "" == strings.TrimSpace(label) || 999 < len(label) {
		ret = nil
	}

	label = strings.TrimSpace(label)

	return
}
