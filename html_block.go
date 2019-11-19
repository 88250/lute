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

import "bytes"

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
	htmlBlockTags1      = [][]byte{strToBytes("<script"), strToBytes("<pre"), strToBytes("<style")}
	htmlBlockCloseTags1 = [][]byte{strToBytes("</script>"), strToBytes("</pre>"), strToBytes("</style>")}
	htmlBlockTags6      = [][]byte{
		strToBytes("<address"), strToBytes("<article"), strToBytes("<aside"), strToBytes("<base"), strToBytes("<basefont"), strToBytes("<blockquote"), strToBytes("<body"), strToBytes("<caption"), strToBytes("<center"), strToBytes("<col"), strToBytes("<colgroup"), strToBytes("<dd"), strToBytes("<details"), strToBytes("<dialog"), strToBytes("<dir"), strToBytes("<div"), strToBytes("<dl"), strToBytes("<dt"), strToBytes("<fieldset"), strToBytes("<figcaption"), strToBytes("<figure"), strToBytes("<footer"), strToBytes("<form"), strToBytes("<frame"), strToBytes("<frameset"), strToBytes("<h1"), strToBytes("<h2"), strToBytes("<h3"), strToBytes("<h4"), strToBytes("<h5"), strToBytes("<h6"), strToBytes("<head"), strToBytes("<header"), strToBytes("<hr"), strToBytes("<html"), strToBytes("<iframe"), strToBytes("<legend"), strToBytes("<li"), strToBytes("<link"), strToBytes("<main"), strToBytes("<menu"), strToBytes("<menuitem"), strToBytes("<nav"), strToBytes("<noframes"), strToBytes("<ol"), strToBytes("<optgroup"), strToBytes("<option"), strToBytes("<p"), strToBytes("<param"), strToBytes("<section"), strToBytes("<source"), strToBytes("<summary"), strToBytes("<table"), strToBytes("<tbody"), strToBytes("<td"), strToBytes("<tfoot"), strToBytes("<th"), strToBytes("<thead"), strToBytes("<title"), strToBytes("<tr"), strToBytes("<track"), strToBytes("<ul"),
		strToBytes("</address"), strToBytes("</article"), strToBytes("</aside"), strToBytes("</base"), strToBytes("</basefont"), strToBytes("</blockquote"), strToBytes("</body"), strToBytes("</caption"), strToBytes("</center"), strToBytes("</col"), strToBytes("</colgroup"), strToBytes("</dd"), strToBytes("</details"), strToBytes("</dialog"), strToBytes("</dir"), strToBytes("</div"), strToBytes("</dl"), strToBytes("</dt"), strToBytes("</fieldset"), strToBytes("</figcaption"), strToBytes("</figure"), strToBytes("</footer"), strToBytes("</form"), strToBytes("</frame"), strToBytes("</frameset"), strToBytes("</h1"), strToBytes("</h2"), strToBytes("</h3"), strToBytes("</h4"), strToBytes("</h5"), strToBytes("</h6"), strToBytes("</head"), strToBytes("</header"), strToBytes("</hr"), strToBytes("</html"), strToBytes("</iframe"), strToBytes("</legend"), strToBytes("</li"), strToBytes("</link"), strToBytes("</main"), strToBytes("</menu"), strToBytes("</menuitem"), strToBytes("</nav"), strToBytes("</noframes"), strToBytes("</ol"), strToBytes("</optgroup"), strToBytes("</option"), strToBytes("</p"), strToBytes("</param"), strToBytes("</section"), strToBytes("</source"), strToBytes("</summary"), strToBytes("</table"), strToBytes("</tbody"), strToBytes("</td"), strToBytes("</tfoot"), strToBytes("</th"), strToBytes("</thead"), strToBytes("</title"), strToBytes("</tr"), strToBytes("</track"), strToBytes("</ul"),
	}
	htmlBlockSinglequote = strToBytes("'")
	htmlBlockDoublequote = strToBytes("\"")
	htmlBlockGreater     = strToBytes(">")
)

func (t *Tree) isHTMLBlockClose(tokens []byte, htmlType int) bool {
	length := len(tokens)
	switch htmlType {
	case 1:
		if pos := acceptTokenss(tokens, htmlBlockCloseTags1); 0 <= pos {
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
		return bytes.Contains(tokens, htmlBlockGreater)
	case 5:
		for i := 0; i < length-2; i++ {
			if itemCloseBracket == tokens[i] && itemCloseBracket == tokens[i+1] {
				return true
			}
		}
	}

	return false
}

func (t *Tree) parseHTML(tokens []byte) (typ int) {
	_, tokens = trimLeft(tokens)
	length := len(tokens)
	if 3 > length { // at least <? and a newline
		return
	}

	if itemLess != tokens[0] {
		return
	}

	typ = 1

	if pos := acceptTokenss(tokens, htmlBlockTags1); 0 <= pos {
		if isWhitespace(tokens[pos]) || itemGreater == tokens[pos] {
			return
		}
	}

	if pos := acceptTokenss(tokens, htmlBlockTags6); 0 <= pos {
		if isWhitespace(tokens[pos]) || itemGreater == tokens[pos] {
			typ = 6
			return
		}
		if itemSlash == tokens[pos] && itemGreater == tokens[pos+1] {
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

	if 0 == bytes.Index(tokens, strToBytes("<!--")) {
		typ = 2
		return
	}

	if 0 == bytes.Index(tokens, strToBytes("<?")) {
		typ = 3
		return
	}

	if 2 < len(tokens) && 0 == bytes.Index(tokens, strToBytes("<!")) {
		following := tokens[2:]
		if 'A' <= following[0] && 'Z' >= following[0] {
			typ = 4
			return
		}
		if 0 == bytes.Index(following, strToBytes("[CDATA[")) {
			typ = 5
			return
		}
	}
	return 0
}

func (t *Tree) isOpenTag(tokens []byte) (isOpenTag bool) {
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

	nameAndAttrs := splitWhitespace(tokens)
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
			if bytes.HasPrefix(value, htmlBlockSinglequote) && bytes.HasSuffix(value, htmlBlockSinglequote) {
				value = value[1:]
				value = value[:len(value)-1]
				return !bytes.Contains(value, htmlBlockSinglequote)
			}
			if bytes.HasPrefix(value, htmlBlockDoublequote) && bytes.HasSuffix(value, htmlBlockDoublequote) {
				value = value[1:]
				value = value[:len(value)-1]
				return !bytes.Contains(value, htmlBlockDoublequote)
			}
			return !bytes.ContainsAny(value, " \t\n") && !bytes.ContainsAny(value, "\"'=<>`")
		}
	}
	return true
}

func (t *Tree) isCloseTag(tokens []byte) bool {
	tokens = trimWhitespace(tokens)
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
