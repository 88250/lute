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
	html.value = html.tokens.replaceNewlineSpace().trimRight().rawText()
	html.tokens = nil
}

func (html *HTML) AcceptLines() bool {
	return true
}

var HTMLBlockTags1, HTMLBlockCloseTags1, HTMLBlockTags6 []items

func init() {
	var tags = []string{"<script", "<pre", "<style"}
	for _, str := range tags {
		HTMLBlockTags1 = append(HTMLBlockTags1, tokenize(str))
	}
	tags = []string{"</script>", "</pre>", "</style>"}
	for _, str := range tags {
		HTMLBlockCloseTags1 = append(HTMLBlockTags1, tokenize(str))
	}
	tags = []string{"address", "article", "aside", "base", "basefont", "blockquote", "body", "caption", "center", "col", "colgroup", "dd", "details", "dialog", "dir", "div", "dl", "dt", "fieldset", "figcaption", "figure", "footer", "form", "frame", "frameset", "h1", "h2", "h3", "h4", "h5", "h6", "head", "header", "hr", "html", "iframe", "legend", "li", "link", "main", "menu", "menuitem", "nav", "noframes", "ol", "optgroup", "option", "p", "param", "section", "source", "summary", "table", "tbody", "td", "tfoot", "th", "thead", "title", "tr", "track", "ul"}
	for _, str := range tags {
		HTMLBlockTags6 = append(HTMLBlockTags6, tokenize("<"+str))
		HTMLBlockTags6 = append(HTMLBlockTags6, tokenize("</"+str))
	}
}

func (t *Tree) isHTMLBlockClose(tokens items, htmlType int) bool {
	length := len(tokens)
	switch htmlType {
	case 1:
		if pos := tokens.acceptTokenss(HTMLBlockCloseTags1); 0 <= pos {
			return true
		}
		return false
	case 2:
		for i := 0; i < length-3; i++ {
			if itemHyphen == tokens[i] && itemHyphen == tokens[i+1] && itemGreater == tokens[i+2] {
				return true
			}
		}
	case 3:
		for i := 0; i < length-2; i++ {
			if itemQuestion == tokens[i] && itemGreater == tokens[i+1] {
				return true
			}
		}
	case 4:
		return tokens.contain(itemGreater)
	case 5:
		for i := 0; i < length-2; i++ {
			if itemCloseBracket == tokens[i] && itemCloseBracket == tokens[i+1] {
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

	if itemLess != tokens[0] {
		return nil
	}

	if pos := tokens.acceptTokenss(HTMLBlockTags1); 0 <= pos {
		if tokens[pos].isWhitespace() || itemGreater == tokens[pos] {
			return &HTML{&BaseNode{typ: NodeHTML}, 1}
		}
	}

	if pos := tokens.acceptTokenss(HTMLBlockTags6); 0 <= pos {
		if tokens[pos].isWhitespace() || itemGreater == tokens[pos] {
			return &HTML{&BaseNode{typ: NodeHTML}, 6}
		}
		if itemSlash == tokens[pos] && itemGreater == tokens[pos+1] {
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

func tokenize(str string) (ret items) {
	for _, r := range str {
		ret = append(ret, item(r))
	}

	return
}

func (tokens items) isOpenTag() (isOpenTag, withAttr bool) {
	length := len(tokens)
	if 3 > length {
		return
	}

	if itemLess != tokens[0] {
		return
	}
	if itemGreater != tokens[length-1] {
		return
	}
	if itemSlash == tokens[length-2] {
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
		if !name[0].isASCIILetter() && itemUnderscore != name[0] && itemColon != name[0] {
			return
		}

		if 1 < len(name) {
			name = name[1:]
			for _, n := range name {
				if !n.isASCIILetter() && !n.isDigit() && itemUnderscore != n && itemDot != n && itemColon != n && itemHyphen != n {
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

	if itemLess != tokens[0] || itemSlash != tokens[1] {
		return false
	}
	if itemGreater != tokens[length-1] {
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
