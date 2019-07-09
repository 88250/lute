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

func (t *Tree) parseHTML(line items, typ int) (ret Node) {
	baseNode := &BaseNode{typ: NodeHTML}
	html := &HTML{baseNode, ""}
	ret = html

Loop:
	for {
		html.Value += line.rawText()
		line = t.nextLine()
		switch typ {
		case 6:
			if line.isBlankLine() {
				break Loop
			}
		}
	}

	return
}

var HTMLBlockTags = []string{"address", "article", "aside", "base", "basefont", "blockquote", "body", "caption", "center", "col", "colgroup", "dd", "details", "dialog", "dir", "div", "dl", "dt", "fieldset", "figcaption", "figure", "footer", "form", "frame", "frameset", "h1", "h2", "h3", "h4", "h5", "h6", "head", "header", "hr", "html", "iframe", "legend", "li", "link", "main", "menu", "menuitem", "nav", "noframes", "ol", "optgroup", "option", "p", "param", "section", "source", "summary", "table", "tbody", "td", "tfoot", "th", "thead", "title", "tr", "track", "ul"}

func (t *Tree) isHTML(line items, htmlType *int) bool {
	length := len(line)
	if 3 > length { // at least <? and a newline
		return false
	}

	if itemLess != line[0].typ {
		return false
	}

	if (t.equalIgnoreCase("script", line[1].val) || t.equalIgnoreCase("pre", line[1].val) || t.equalIgnoreCase("style", line[1].val)) && 0 < line[1:].whitespaceCountLeft() {
		l := line[1:].trimLeft()
		if 1 > len(l) {
			return false
		}

		if itemGreater == l[0].typ || l[0].isNewline() || l[0].isEOF() {
			*htmlType = 1
			return true
		}
	}

	slash := itemSlash == line[1].typ
	i := 1
	if slash {
		i = 2
	}
	rule6 := t.equalAnyIgnoreCase(line[i].val, HTMLBlockTags)
	if rule6 {
		i++
		if line[i].isWhitespace() || itemGreater == line[i].typ {
			*htmlType = 6
			return true
		}
		if i < length && itemSlash == line[i].typ && itemGreater == line[i+1].typ {
			*htmlType = 6
			return true
		}
	}

	i = 1
	if slash {
		i = 2
	}
	rule7 := line[i].isASCIILetterNumHyphen()
	if rule7 {
		i++
		if line[i].isWhitespace() || itemGreater == line[i].typ {
			*htmlType = 7
			return true
		}
		if i < length && itemSlash == line[i].typ && itemGreater == line[i+1].typ {
			*htmlType = 7
			return true
		}
	}

	rawText := line.rawText()
	if 0 == strings.Index(rawText, "<!--") {
		*htmlType = 2
		return true
	}

	if 0 == strings.Index(rawText, "<?") {
		*htmlType = 3
		return true
	}

	if 2 < len(rawText) && 0 == strings.Index(rawText, "<!") {
		following := rawText[2:]
		if 'A' <= following[0] && 'Z' >= following[0] {
			*htmlType = 4
			return true
		}
		if 0 == strings.Index(following, "[CDATA[") {
			*htmlType = 5
			return true
		}
	}

	return false
}

func (t *Tree) startWithAnyIgnoreCase(s1 string, strs []string) (pos int) {
	for _, s := range strs {
		s1 = strings.ToLower(s1)
		s = strings.ToLower(s)
		if 0 == strings.Index(s1, s) {
			return len(s)
		}
	}

	return -1
}

func (t *Tree) equalAnyIgnoreCase(s1 string, strs []string) bool {
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
