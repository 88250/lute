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

func (t *Tree) parseLinkRefDef(line items) bool {
	_, line = line.trimLeft()
	if 1 > len(line) {
		return false
	}

	linkLabel, remains, label := t.parseLinkLabel(line)
	if nil == linkLabel {
		return false
	}

	if nil != t.context.LinkRefDef[label] {
		return false
	}

	if itemColon != remains[0].typ {
		return false
	}

	remains = remains[1:]
	whitespaces, remains := remains.trimLeft()
	newlines, _, _ := whitespaces.statWhitespace()
	if 1 < newlines {
		return false
	}

	tokens := remains
	linkDest, remains, url := t.parseLinkDest1(tokens)
	if nil == linkDest {
		linkDest, remains, url = t.parseLinkDest2(tokens)
	}
	if nil == linkDest {
		return false
	}

	whitespaces, remains = remains.trimLeft()
	newlines, _, _ = whitespaces.statWhitespace()
	if 1 < newlines {
		return false
	}

	tokens = remains
	linkTitle, remains, title := t.parseLinkTitle1(tokens)
	if nil == linkTitle {
		return false
	}
	_ = remains
	_ = url
	_ = title

	link := &Link{&BaseNode{typ: NodeLink}, url, title}
	t.context.LinkRefDef[label] = link

	return true
}

func (t *Tree) parseLinkTitle1(tokens items) (ret, remains items, title string) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemDoublequote != tokens[0].typ {
		return
	}

	close := false
	i := 0
	for ; i < length; i++ {
		token := tokens[i]
		ret = append(ret, token)
		if 0 < i {
			title += token.val
			if itemDoublequote == token.typ && !tokens.isBackslashEscape(i) {
				close = true
				title = title[:len(title)-1]
				break
			}
		}
	}

	if !close {
		ret = nil
		title = ""

		return
	}

	remains = tokens[i+1:]

	return
}

func (t *Tree) parseLinkDest2(tokens items) (ret, remains items, url string) {
	remains = tokens
	var leftParens, rightParens int
	i := 0
	for ; i < len(tokens); i++ {
		token := tokens[i]
		ret = append(ret, token)
		url += token.val
		if itemSpace == token.typ || token.isControl() {
			url = url[0 : len(url)-1]
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
		url = ""
		return
	}

	remains = tokens[i:]

	return
}

func (t *Tree) parseLinkDest1(tokens items) (ret, remains items, url string) {
	remains = tokens
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemLess != tokens[0].typ {
		return
	}

	close := false
	i := 0
	for ; i < length; i++ {
		token := tokens[i]
		ret = append(ret, token)
		if 0 < i {
			url += token.val
			if itemLess == token.typ && !tokens.isBackslashEscape(i) {
				ret = nil
				url = ""
				return
			}
		}

		if itemGreater == token.typ && !tokens.isBackslashEscape(i) {
			close = true
			url = url[0 : len(url)-1]
			break
		}
	}

	if !close {
		ret = nil
		url = ""

		return
	}

	remains = tokens[i+1:]

	return
}

func (t *Tree) parseLinkLabel(tokens items) (ret, remains items, label string) {
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemOpenBracket != tokens[0].typ {
		return
	}

	close := false
	i := 0
	for ; i < length; i++ {
		token := tokens[i]
		ret = append(ret, token)
		if 0 < i {
			label += token.val
		}

		if itemCloseBracket == token.typ && !tokens.isBackslashEscape(i) {
			close = true
			label = label[0 : len(label)-1]
			break
		}
	}

	if !close || "" == strings.TrimSpace(label) || 999 < len(label) {
		ret = nil
	}

	remains = tokens[i+1:]

	return
}
