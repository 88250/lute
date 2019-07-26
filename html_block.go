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

func (html *HTML) Finalize() {
	// TODO html.rawText = html.rawText._string_content.replace(/(\n *)+$/, '');

}

func (html *HTML) AcceptLines() bool {
	return true
}

func (t *Tree) parseHTML(tokens items, typ int) (ret *HTML) {
	ret = &HTML{BaseNode: &BaseNode{typ: NodeHTML}}
	openTagName := tokens.split(itemGreater)[0][1].val
	for {
		matchEnd := false
		for i, token := range tokens {
			if 1 == typ {
				if !matchEnd && itemLess == token.typ && i < len(tokens)-3 && itemSlash == tokens[i+1].typ {
					if openTagName == tokens[i+2].val && itemGreater == tokens[i+3].typ {
						matchEnd = true
					}
				}
			} else if 2 == typ {
				if !matchEnd && itemHyphen == token.typ && i < len(tokens)-2 && itemHyphen == tokens[i+1].typ && itemGreater == tokens[i+2].typ {
					matchEnd = true
				}
			} else if 3 == typ {
				if !matchEnd && itemQuestion == token.typ && i < len(tokens)-1 && itemGreater == tokens[i+1].typ {
					matchEnd = true
				}
			} else if 5 == typ {
				if !matchEnd && itemCloseBracket == token.typ && i < len(tokens)-2 && itemCloseBracket == tokens[i+1].typ && itemGreater == tokens[i+2].typ {
					matchEnd = true
				}
			}

			ret.AppendValue(token.val)
		}

		if matchEnd && 6 != typ {
			break
		}

		tokens = t.nextLine()
		if tokens.isEOF() {
			break
		}

		if 1 == typ || 2 == typ || 4 == typ || 3 == typ || 5 == typ {
			continue
		}

		if 6 == typ || 7 == typ {
			if tokens.isBlankLine() {
				break
			}
		} else {
			break
		}
	}

	ret.value = strings.TrimRight(ret.value, "\n")

	return
}

var HTMLBlockTags = []string{"address", "article", "aside", "base", "basefont", "blockquote", "body", "caption", "center", "col", "colgroup", "dd", "details", "dialog", "dir", "div", "dl", "dt", "fieldset", "figcaption", "figure", "footer", "form", "frame", "frameset", "h1", "h2", "h3", "h4", "h5", "h6", "head", "header", "hr", "html", "iframe", "legend", "li", "link", "main", "menu", "menuitem", "nav", "noframes", "ol", "optgroup", "option", "p", "param", "section", "source", "summary", "table", "tbody", "td", "tfoot", "th", "thead", "title", "tr", "track", "ul"}

func (t *Tree) isHTML(tokens items) (htmlType int) {
	_, tokens = tokens.trimLeft()
	length := len(tokens)
	if 3 > length { // at least <? and a newline
		return -1
	}

	if itemLess != tokens[0].typ {
		return -1
	}

	if t.equalAnyIgnoreCase(tokens[1].val, "script", "pre", "style") {
		l := tokens[2:]
		if 1 > len(l) {
			return -1
		}

		if l[0].isWhitespace() || itemGreater == l[0].typ || l[0].isEOF() {
			return 1
		}
	}

	slash := itemSlash == tokens[1].typ
	i := 1
	if slash {
		i = 2
	}
	rule6 := t.equalAnyIgnoreCase(tokens[i].val, HTMLBlockTags...)
	if rule6 {
		i++
		if tokens[i].isWhitespace() || itemGreater == tokens[i].typ {
			return 6
		}
		if i < length && itemSlash == tokens[i].typ && itemGreater == tokens[i+1].typ {
			return 6
		}
	}

	tag := tokens.trim()
	isOpenTag, _ := tag.isOpenTag()
	if isOpenTag {
		return 7
	}
	isCloseTag := tag.isCloseTag()
	if isCloseTag {
		return 7
	}

	rawText := tokens.rawText()
	if 0 == strings.Index(rawText, "<!--") {
		return 2
	}

	if 0 == strings.Index(rawText, "<?") {
		return 3
	}

	if 2 < len(rawText) && 0 == strings.Index(rawText, "<!") {
		following := rawText[2:]
		if 'A' <= following[0] && 'Z' >= following[0] {
			return 4
		}
		if 0 == strings.Index(following, "[CDATA[") {
			return 5
		}
	}

	return -1
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
	tokens = tokens.trim()
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
