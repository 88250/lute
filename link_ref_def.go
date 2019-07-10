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

func (t *Tree) parseLinkRefDef(line items) bool {
	line = line.trimLeft()
	if 1 > len(line) {
		return false
	}

	linkLabelTokens, label := t.parseLinkLabel(line)
	if nil == linkLabelTokens {
		return false
	}

	if nil != t.context.LinkRefDef[label] {
		link := &Link{&BaseNode{typ: NodeLink}, "url", "title"}
		t.context.LinkRefDef[label] = link
	}

	return true
}

func (t *Tree) parseLinkLabel(tokens items) (ret items, label string) {
	length := len(tokens)
	if 2 > length {
		return
	}

	if itemOpenBracket != tokens[0].typ {
		return
	}

	close := false
	for i := 0; i < length; i++ {
		token := tokens[i]
		ret = append(ret, token)
		if 0 < i && !token.isWhitespace() {
			label += token.val
		}

		if itemCloseBracket == token.typ && !tokens.isBackslashEscape(i) {
			close = true
			label = label[0 : len(label)-1]
			break
		}
	}

	if !close || "" == label || 999 < len(label) {
		ret = nil
	}

	return
}
