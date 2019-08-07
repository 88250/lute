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

// HTMLBlock 描述了 HTML 块节点结构。
type HTMLBlock struct {
	*BaseNode
	hType int // 规范中定义的 HTML 块类型（1-7）
}

func (html *HTMLBlock) CanContain(nodeType int) bool {
	return false
}

func (html *HTMLBlock) Continue(context *Context) int {
	if context.blank && (html.hType == 6 || html.hType == 7) {
		return 1
	}
	return 0
}

func (html *HTMLBlock) Finalize(context *Context) {
	html.tokens = html.tokens.replaceNewlineSpace().trimRight()
}

func (html *HTMLBlock) AcceptLines() bool {
	return true
}

var htmlBlockTags1, htmlBlockCloseTags1, htmlBlockTags6 []items

func init() {
	var tags = []string{"<script", "<pre", "<style"}
	for _, str := range tags {
		htmlBlockTags1 = append(htmlBlockTags1, tokenize(str))
	}
	tags = []string{"</script>", "</pre>", "</style>"}
	for _, str := range tags {
		htmlBlockCloseTags1 = append(htmlBlockCloseTags1, tokenize(str))
	}
	tags = []string{"address", "article", "aside", "base", "basefont", "blockquote", "body", "caption", "center", "col", "colgroup", "dd", "details", "dialog", "dir", "div", "dl", "dt", "fieldset", "figcaption", "figure", "footer", "form", "frame", "frameset", "h1", "h2", "h3", "h4", "h5", "h6", "head", "header", "hr", "html", "iframe", "legend", "li", "link", "main", "menu", "menuitem", "nav", "noframes", "ol", "optgroup", "option", "p", "param", "section", "source", "summary", "table", "tbody", "td", "tfoot", "th", "thead", "title", "tr", "track", "ul"}
	for _, str := range tags {
		htmlBlockTags6 = append(htmlBlockTags6, tokenize("<"+str))
		htmlBlockTags6 = append(htmlBlockTags6, tokenize("</"+str))
	}
}

func (t *Tree) isHTMLBlockClose(tokens items, htmlType int) bool {
	length := len(tokens)
	switch htmlType {
	case 1:
		if pos := tokens.acceptTokenss(htmlBlockCloseTags1); 0 <= pos {
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

func (t *Tree) parseHTML(tokens items) (ret *HTMLBlock) {
	_, tokens = tokens.trimLeft()
	length := len(tokens)
	if 3 > length { // at least <? and a newline
		return nil
	}

	if itemLess != tokens[0] {
		return nil
	}

	ret = &HTMLBlock{&BaseNode{typ: NodeHTMLBlock, tokens: make(items, 0, 256)}, 1}

	if pos := tokens.acceptTokenss(htmlBlockTags1); 0 <= pos {
		if isWhitespace(tokens[pos]) || itemGreater == tokens[pos] {
			return
		}
	}

	if pos := tokens.acceptTokenss(htmlBlockTags6); 0 <= pos {
		if isWhitespace(tokens[pos]) || itemGreater == tokens[pos] {
			ret.hType = 6
			return
		}
		if itemSlash == tokens[pos] && itemGreater == tokens[pos+1] {
			ret.hType = 6
			return
		}
	}

	tag := tokens.trim()
	isOpenTag, _ := tag.isOpenTag()
	if isOpenTag && t.context.tip.Type() != NodeParagraph {
		ret.hType = 7
		return
	}
	isCloseTag := tag.isCloseTag()
	if isCloseTag && t.context.tip.Type() != NodeParagraph {
		ret.hType = 7
		return
	}

	rawText := tokens.string()
	if 0 == strings.Index(rawText, "<!--") {
		ret.hType = 2
		return
	}

	if 0 == strings.Index(rawText, "<?") {
		ret.hType = 3
		return
	}

	if 2 < len(rawText) && 0 == strings.Index(rawText, "<!") {
		following := rawText[2:]
		if 'A' <= following[0] && 'Z' >= following[0] {
			ret.hType = 4
			return
		}
		if 0 == strings.Index(following, "[CDATA[") {
			ret.hType = 5
			return
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

// tokenize 在 init 函数中调用，可以认为是静态分配，所以使用拷贝字符不会有性能问题。
// 另外，这里也必须要拷贝，因为调用点的 str 是局部变量，地址上的值会被覆盖。
func tokenize(str string) (ret items) {
	for _, r := range str {
		ret = append(ret, byte(r))
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

	if isWhitespace(tokens[0]) { // < 后面不能跟空白
		return
	}

	nameAndAttrs := tokens.splitWhitespace()
	name := nameAndAttrs[0]
	if !isASCIILetter(name[0]) {
		return
	}
	if 1 < len(name) {
		name = name[1:]
		for _, n := range name {
			if !isASCIILetterNumHyphen(n) {
				return
			}
		}
	}

	withAttr = true
	nameAndAttrs = nameAndAttrs[1:]
	for _, nameAndAttr := range nameAndAttrs {
		nameAndValue := nameAndAttr.split(itemEqual)
		name := nameAndValue[0]
		if !isASCIILetter(name[0]) && itemUnderscore != name[0] && itemColon != name[0] {
			return
		}

		if 1 < len(name) {
			name = name[1:]
			for _, n := range name {
				if !isASCIILetter(n) && !isDigit(n) && itemUnderscore != n && itemDot != n && itemColon != n && itemHyphen != n {
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
	if !isASCIILetter(name[0]) {
		return false
	}
	if 1 < len(name) {
		name = name[1:]
		for _, n := range name {
			if !isASCIILetterNumHyphen(n) {
				return false
			}
		}
	}

	return true
}
