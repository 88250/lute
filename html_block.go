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

type HTML struct {
	*BaseNode
	hType int
}

func (html *HTML) CanContain(nodeType NodeType) bool {
	return false
}

func (html *HTML) Continue(context *Context) int {
	if context.blank && (html.hType == 6 || html.hType == 7) {
		return 1
	}
	return 0
}

func (html *HTML) Finalize(context *Context) {
	html.value = strings.TrimRight(html.value, "\n ")
}

func (html *HTML) AcceptLines() bool {
	return true
}

var HTMLBlockTags = []string{"address", "article", "aside", "base", "basefont", "blockquote", "body", "caption", "center", "col", "colgroup", "dd", "details", "dialog", "dir", "div", "dl", "dt", "fieldset", "figcaption", "figure", "footer", "form", "frame", "frameset", "h1", "h2", "h3", "h4", "h5", "h6", "head", "header", "hr", "html", "iframe", "legend", "li", "link", "main", "menu", "menuitem", "nav", "noframes", "ol", "optgroup", "option", "p", "param", "section", "source", "summary", "table", "tbody", "td", "tfoot", "th", "thead", "title", "tr", "track", "ul"}

func (t *Tree) isHTMLBlockClose(tokens items, htmlType int) bool {
	length := len(tokens)
	switch htmlType {
	case 1:
		for i := 0; i < length-3; i++ {
			if itemLess == tokens[i].typ && itemSlash == tokens[i+1].typ && t.equalAnyIgnoreCase(tokens[i+2].Value(), "script", "pre", "style") && itemGreater == tokens[i+3].typ {
				return true
			}
		}
	case 2:
		for i := 0; i < length-3; i++ {
			if itemHyphen == tokens[i].typ && itemHyphen == tokens[i+1].typ && itemGreater == tokens[i+2].typ {
				return true
			}
		}
	case 3:
		for i := 0; i < length-2; i++ {
			if itemQuestion == tokens[i].typ && itemGreater == tokens[i+1].typ {
				return true
			}
		}
	case 4:
		return tokens.contain(itemGreater)
	case 5:
		for i := 0; i < length-2; i++ {
			if itemCloseBracket == tokens[i].typ && itemCloseBracket == tokens[i+1].typ {
				return true
			}
		}
	}

	return false
}

func (t *Tree) parseHTML(tokens items) (ret *HTML) {
	_, tokens = tokens.trimLeft()
	length := len(tokens)
	if 3 > length { // at least <? and a newline
		return nil
	}

	if itemLess != tokens[0].typ {
		return nil
	}

	if t.equalAnyIgnoreCase(tokens[1].Value(), "script", "pre", "style") {
		l := tokens[2:]
		if 1 > len(l) {
			return nil
		}

		if l[0].isWhitespace() || itemGreater == l[0].typ || l[0].isEOF() {
			return &HTML{&BaseNode{typ: NodeHTML}, 1}
		}
	}

	slash := itemSlash == tokens[1].typ
	i := 1
	if slash {
		i = 2
	}
	rule6 := t.equalAnyIgnoreCase(tokens[i].Value(), HTMLBlockTags...)
	if rule6 {
		i++
		if tokens[i].isWhitespace() || itemGreater == tokens[i].typ {
			return &HTML{&BaseNode{typ: NodeHTML}, 6}
		}
		if i < length && itemSlash == tokens[i].typ && itemGreater == tokens[i+1].typ {
			return &HTML{&BaseNode{typ: NodeHTML}, 6}
		}
	}

	tag := tokens.trim()
	isOpenTag, _ := tag.isOpenTag()
	if isOpenTag && t.context.tip.Type() != NodeParagraph {
		return &HTML{&BaseNode{typ: NodeHTML}, 7}
	}
	isCloseTag := tag.isCloseTag()
	if isCloseTag && t.context.tip.Type() != NodeParagraph {
		return &HTML{&BaseNode{typ: NodeHTML}, 7}
	}

	rawText := tokens.rawText()
	if 0 == strings.Index(rawText, "<!--") {
		return &HTML{&BaseNode{typ: NodeHTML}, 2}
	}

	if 0 == strings.Index(rawText, "<?") {
		return &HTML{&BaseNode{typ: NodeHTML}, 3}
	}

	if 2 < len(rawText) && 0 == strings.Index(rawText, "<!") {
		following := rawText[2:]
		if 'A' <= following[0] && 'Z' >= following[0] {
			return &HTML{&BaseNode{typ: NodeHTML}, 4}
		}
		if 0 == strings.Index(following, "[CDATA[") {
			return &HTML{&BaseNode{typ: NodeHTML}, 5}
		}
	}

	return nil
}

func (t *Tree) startWithAnyIgnoreCase(s1 string, strs ...string) (pos int) {
	for _, s := range strs {
		s1 = strings.ToLower(s1)
		s = strings.ToLower(s)
		if 0 == strings.Index(s1, s) {
			return len(s)
		}
	}

	return -1
}

func (t *Tree) equalAnyIgnoreCase(s1 string, strs ...string) bool {
	for _, s := range strs {
		if t.equalIgnoreCase(s1, s) {
			return true
		}
	}

	return false
}

func (t *Tree) equalIgnoreCase(s1, s2 string) bool {
	return strings.ToLower(s1) == strings.ToLower(s2)
}

func (tokens items) isOpenTag() (isOpenTag, withAttr bool) {
	length := len(tokens)
	if 3 > length {
		return
	}

	if itemLess != tokens[0].typ {
		return
	}
	if itemGreater != tokens[length-1].typ {
		return
	}
	if itemSlash == tokens[length-2].typ {
		tokens = tokens[1 : length-2]
	} else {
		tokens = tokens[1 : length-1]
	}

	length = len(tokens)
	if 0 == length {
		return
	}

	if tokens[0].isWhitespace() { // < 后面不能跟空白
		return
	}

	nameAndAttrs := tokens.splitWhitespace()
	name := nameAndAttrs[0]
	if !name[0].isASCIILetter() {
		return
	}
	if 1 < len(name) {
		name = name[1:]
		for _, n := range name {
			if !n.isASCIILetterNumHyphen() {
				return
			}
		}
	}

	withAttr = true
	nameAndAttrs = nameAndAttrs[1:]
	for _, nameAndAttr := range nameAndAttrs {
		nameAndValue := nameAndAttr.split(itemEqual)
		name := nameAndValue[0]
		if !name[0].isASCIILetter() && itemUnderscore != name[0].typ && itemColon != name[0].typ {
			return
		}

		if 1 < len(name) {
			name = name[1:]
			for _, n := range name {
				if !n.isASCIILetter() && !n.isNumInt() && itemUnderscore != n.typ && itemDot != n.typ && itemColon != n.typ && itemHyphen != n.typ {
					return
				}
			}
		}

		if 1 < len(nameAndValue) {
			value := nameAndValue[1]
			if value.startWith(itemSinglequote) && value.endWith(itemSinglequote) {
				value = value[1:]
				value = value[:len(value)-1]
				return !value.contain(itemSinglequote), withAttr
			}
			if value.startWith(itemDoublequote) && value.endWith(itemDoublequote) {
				value = value[1:]
				value = value[:len(value)-1]
				return !value.contain(itemDoublequote), withAttr
			}
			return !value.containWhitespace() && !value.contain(itemSinglequote, itemDoublequote, itemEqual, itemLess, itemGreater, itemBacktick), withAttr
		}
	}
	return true, withAttr
}

func (tokens items) isCloseTag() bool {
	tokens = tokens.trim()
	length := len(tokens)
	if 4 > length {
		return false
	}

	if itemLess != tokens[0].typ || itemSlash != tokens[1].typ {
		return false
	}
	if itemGreater != tokens[length-1].typ {
		return false
	}

	tokens = tokens[2 : length-1]
	length = len(tokens)
	if 0 == length {
		return false
	}

	name := tokens[0:]
	if !name[0].isASCIILetter() {
		return false
	}
	if 1 < len(name) {
		name = name[1:]
		for _, n := range name {
			if !n.isASCIILetterNumHyphen() {
				return false
			}
		}
	}

	return true
}
