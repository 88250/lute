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

func (html *Node) htmlBlockContinue(context *Context) int {
	if context.blank && (html.htmlBlockType == 6 || html.htmlBlockType == 7) {
		return 1
	}
	return 0
}

func (html *Node) htmlBlockFinalize(context *Context) {
	_, html.tokens = trimRight(replaceNewlineSpace(html.tokens))
}

var (
	htmlBlockTags1      = []items{strToItems("<script"), strToItems("<pre"), strToItems("<style")}
	htmlBlockCloseTags1 = []items{strToItems("</script>"), strToItems("</pre>"), strToItems("</style>")}
	htmlBlockTags6      = []items{
		strToItems("<address"), strToItems("<article"), strToItems("<aside"), strToItems("<base"), strToItems("<basefont"), strToItems("<blockquote"), strToItems("<body"), strToItems("<caption"), strToItems("<center"), strToItems("<col"), strToItems("<colgroup"), strToItems("<dd"), strToItems("<details"), strToItems("<dialog"), strToItems("<dir"), strToItems("<div"), strToItems("<dl"), strToItems("<dt"), strToItems("<fieldset"), strToItems("<figcaption"), strToItems("<figure"), strToItems("<footer"), strToItems("<form"), strToItems("<frame"), strToItems("<frameset"), strToItems("<h1"), strToItems("<h2"), strToItems("<h3"), strToItems("<h4"), strToItems("<h5"), strToItems("<h6"), strToItems("<head"), strToItems("<header"), strToItems("<hr"), strToItems("<html"), strToItems("<iframe"), strToItems("<legend"), strToItems("<li"), strToItems("<link"), strToItems("<main"), strToItems("<menu"), strToItems("<menuitem"), strToItems("<nav"), strToItems("<noframes"), strToItems("<ol"), strToItems("<optgroup"), strToItems("<option"), strToItems("<p"), strToItems("<param"), strToItems("<section"), strToItems("<source"), strToItems("<summary"), strToItems("<table"), strToItems("<tbody"), strToItems("<td"), strToItems("<tfoot"), strToItems("<th"), strToItems("<thead"), strToItems("<title"), strToItems("<tr"), strToItems("<track"), strToItems("<ul"),
		strToItems("</address"), strToItems("</article"), strToItems("</aside"), strToItems("</base"), strToItems("</basefont"), strToItems("</blockquote"), strToItems("</body"), strToItems("</caption"), strToItems("</center"), strToItems("</col"), strToItems("</colgroup"), strToItems("</dd"), strToItems("</details"), strToItems("</dialog"), strToItems("</dir"), strToItems("</div"), strToItems("</dl"), strToItems("</dt"), strToItems("</fieldset"), strToItems("</figcaption"), strToItems("</figure"), strToItems("</footer"), strToItems("</form"), strToItems("</frame"), strToItems("</frameset"), strToItems("</h1"), strToItems("</h2"), strToItems("</h3"), strToItems("</h4"), strToItems("</h5"), strToItems("</h6"), strToItems("</head"), strToItems("</header"), strToItems("</hr"), strToItems("</html"), strToItems("</iframe"), strToItems("</legend"), strToItems("</li"), strToItems("</link"), strToItems("</main"), strToItems("</menu"), strToItems("</menuitem"), strToItems("</nav"), strToItems("</noframes"), strToItems("</ol"), strToItems("</optgroup"), strToItems("</option"), strToItems("</p"), strToItems("</param"), strToItems("</section"), strToItems("</source"), strToItems("</summary"), strToItems("</table"), strToItems("</tbody"), strToItems("</td"), strToItems("</tfoot"), strToItems("</th"), strToItems("</thead"), strToItems("</title"), strToItems("</tr"), strToItems("</track"), strToItems("</ul"),
	}
	htmlBlockSinglequote = strToItems("'")
	htmlBlockDoublequote = strToItems("\"")
	htmlBlockGreater     = strToItems(">")
)

func (t *Tree) isHTMLBlockClose(tokens items, htmlType int) bool {
	length := len(tokens)
	switch htmlType {
	case 1:
		if pos := acceptTokenss(tokens, htmlBlockCloseTags1); 0 <= pos {
			return true
		}
		return false
	case 2:
		for i := 0; i < length-3; i++ {
			if itemHyphen == tokens[i].term() && itemHyphen == tokens[i+1].term() && itemGreater == tokens[i+2].term() {
				return true
			}
		}
	case 3:
		for i := 0; i < length-2; i++ {
			if itemQuestion == tokens[i].term() && itemGreater == tokens[i+1].term() {
				return true
			}
		}
	case 4:
		return contains(tokens, htmlBlockGreater)
	case 5:
		for i := 0; i < length-2; i++ {
			if itemCloseBracket == tokens[i].term() && itemCloseBracket == tokens[i+1].term() {
				return true
			}
		}
	}

	return false
}

func (t *Tree) parseHTML(tokens items) (typ int) {
	_, tokens = trimLeft(tokens)
	length := len(tokens)
	if 3 > length { // at least <? and a newline
		return
	}

	if itemLess != tokens[0].term() {
		return
	}

	typ = 1

	if pos := acceptTokenss(tokens, htmlBlockTags1); 0 <= pos {
		if isWhitespace(tokens[pos].term()) || itemGreater == tokens[pos].term() {
			return
		}
	}

	if pos := acceptTokenss(tokens, htmlBlockTags6); 0 <= pos {
		if isWhitespace(tokens[pos].term()) || itemGreater == tokens[pos].term() {
			typ = 6
			return
		}
		if itemSlash == tokens[pos].term() && itemGreater == tokens[pos+1].term() {
			typ = 6
			return
		}
	}

	tag := trimWhitespace(tokens)
	isOpenTag := t.isOpenTag(tag)
	if isOpenTag && t.context.tip.typ != NodeParagraph {
		typ = 7
		return
	}
	isCloseTag := t.isCloseTag(tag)
	if isCloseTag && t.context.tip.typ != NodeParagraph {
		typ = 7
		return
	}

	if 0 == index(tokens, strToItems("<!--")) {
		typ = 2
		return
	}

	if 0 == index(tokens, strToItems("<?")) {
		typ = 3
		return
	}

	if 2 < len(tokens) && 0 == index(tokens, strToItems("<!")) {
		following := tokens[2:]
		if 'A' <= following[0].term() && 'Z' >= following[0].term() {
			typ = 4
			return
		}
		if 0 == index(following, strToItems("[CDATA[")) {
			typ = 5
			return
		}
	}
	return 0
}

func (t *Tree) isOpenTag(tokens items) (isOpenTag bool) {
	length := len(tokens)
	if 3 > length {
		return
	}

	if itemLess != tokens[0].term() {
		return
	}
	if itemGreater != tokens[length-1].term() {
		return
	}
	if itemSlash == tokens[length-2].term() {
		tokens = tokens[1 : length-2]
	} else {
		tokens = tokens[1 : length-1]
	}

	length = len(tokens)
	if 0 == length {
		return
	}

	if isWhitespace(tokens[0].term()) { // < 后面不能跟空白
		return
	}

	nameAndAttrs := splitWhitespace(tokens)
	name := nameAndAttrs[0]
	if !isASCIILetter(name[0].term()) {
		return
	}
	if 1 < len(name) {
		name = name[1:]
		for _, n := range name {
			if !isASCIILetterNumHyphen(n.term()) {
				return
			}
		}
	}

	attrs := nameAndAttrs[1:]
	for _, attr := range attrs {
		if 1 >= len(attr) {
			continue
		}

		nameAndValue := split(attr, itemEqual)
		name := nameAndValue[0]
		if 1 > len(name) { // 等号前面空格的情况
			continue
		}
		if !isASCIILetter(name[0].term()) && itemUnderscore != name[0].term() && itemColon != name[0].term() {
			return
		}

		if 1 < len(name) {
			name = name[1:]
			for _, n := range name {
				if !isASCIILetter(n.term()) && !isDigit(n.term()) && itemUnderscore != n.term() && itemDot != n.term() && itemColon != n.term() && itemHyphen != n.term() {
					return
				}
			}
		}

		if 1 < len(nameAndValue) {
			value := nameAndValue[1]
			if hasPrefix(value, htmlBlockSinglequote) && hasSuffix(value, htmlBlockSinglequote) {
				value = value[1:]
				value = value[:len(value)-1]
				return !contains(value, htmlBlockSinglequote)
			}
			if hasPrefix(value, htmlBlockDoublequote) && hasSuffix(value, htmlBlockDoublequote) {
				value = value[1:]
				value = value[:len(value)-1]
				return !contains(value, htmlBlockDoublequote)
			}
			return !containsAny(value, " \t\n") && !containsAny(value, "\"'=<>`")
		}
	}
	return true
}

func (t *Tree) isCloseTag(tokens items) bool {
	tokens = trimWhitespace(tokens)
	length := len(tokens)
	if 4 > length {
		return false
	}

	if itemLess != tokens[0].term() || itemSlash != tokens[1].term() {
		return false
	}
	if itemGreater != tokens[length-1].term() {
		return false
	}

	tokens = tokens[2 : length-1]
	length = len(tokens)
	if 0 == length {
		return false
	}

	name := tokens[0:]
	if !isASCIILetter(name[0].term()) {
		return false
	}
	if 1 < len(name) {
		name = name[1:]
		for _, n := range name {
			if !isASCIILetterNumHyphen(n.term()) {
				return false
			}
		}
	}

	return true
}
