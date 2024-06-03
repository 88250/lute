// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"bytes"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

// HtmlBlockStart 判断 HTML 块（<）是否开始。
func HtmlBlockStart(t *Tree, container *ast.Node) int {
	if t.Context.indented {
		return 0
	}

	if lex.ItemLess != lex.Peek(t.Context.currentLine, t.Context.nextNonspace) {
		return 0
	}

	if t.Context.ParseOption.VditorWYSIWYG {
		if bytes.Contains(t.Context.currentLine, []byte("vditor-comment")) {
			return 0
		}
	}
	if t.Context.ParseOption.ProtyleWYSIWYG {
		if bytes.Contains(t.Context.currentLine, []byte("<span ")) {
			return 0
		}
	}

	tokens := t.Context.currentLine[t.Context.nextNonspace:]
	if htmlType := t.parseHTML(tokens); 0 != htmlType {
		t.Context.closeUnmatchedBlocks()

		if t.Context.ParseOption.ProtyleWYSIWYG {
			tokens = bytes.TrimSpace(tokens)
			if bytes.HasPrefix(tokens, []byte("<iframe")) && bytes.HasSuffix(tokens, []byte(">")) {
				if bytes.Contains(tokens, []byte("data-subtype=\"widget\"")) {
					t.Context.addChild(ast.NodeWidget)
				} else {
					t.Context.addChild(ast.NodeIFrame)
				}
				return 2
			} else if bytes.HasPrefix(tokens, []byte("<video")) && bytes.HasSuffix(tokens, []byte(">")) {
				t.Context.addChild(ast.NodeVideo)
				return 2
			} else if bytes.HasPrefix(tokens, []byte("<audio")) && bytes.HasSuffix(tokens, []byte(">")) {
				t.Context.addChild(ast.NodeAudio)
				return 2
			} else if bytes.HasPrefix(tokens, []byte("<div")) &&
				bytes.Contains(tokens, []byte("data-type=\"NodeAttributeView\"")) &&
				bytes.Contains(tokens, []byte("data-av-type=\"")) &&
				bytes.HasSuffix(tokens, []byte("</div>")) {
				av := t.Context.addChild(ast.NodeAttributeView)
				avTypeIdx := bytes.Index(tokens, []byte("data-av-type=\"")) + len("data-av-type=\"")
				avTypeEndIdx := avTypeIdx + bytes.Index(tokens[avTypeIdx:], []byte("\""))
				av.AttributeViewType = string(tokens[avTypeIdx:avTypeEndIdx])
				if avIdIdx := bytes.Index(tokens, []byte("data-av-id=\"")); 0 < avIdIdx {
					avIdIdx = avIdIdx + len("data-av-id=\"")
					avIdEndIdx := avIdIdx + bytes.Index(tokens[avIdIdx:], []byte("\""))
					av.AttributeViewID = string(tokens[avIdIdx:avIdEndIdx])
				} else {
					av.AttributeViewID = ast.NewNodeID()

				}
				return 2
			}
		}

		if t.Context.ParseOption.ProtyleWYSIWYG {
			// Protyle WYSIWYG 模式下，只有 <div 开头的块级元素才能被解析为 HTML 块
			// Only HTML code wrapped in `<div>` is supported to be parsed into HTML blocks https://github.com/siyuan-note/siyuan/issues/9758
			_, start := lex.TrimLeft(t.Context.currentLine)
			if !bytes.HasPrefix(start, []byte("<div")) {
				return 0
			}
		}

		block := t.Context.addChild(ast.NodeHTMLBlock)
		block.HtmlBlockType = htmlType
		return 2
	}
	return 0
}

func HtmlBlockContinue(html *ast.Node, context *Context) int {
	tokens := context.currentLine
	if context.ParseOption.KramdownBlockIAL && simpleCheckIsBlockIAL(tokens) {
		// 判断 IAL 打断
		if context.Tip.ParentIs(ast.NodeListItem) {
			_, tokens = lex.TrimLeft(tokens)
		}
		if ial := context.parseKramdownBlockIAL(tokens); 0 < len(ial) {
			context.Tip.ID = IAL2Map(ial)["id"]
			context.Tip.KramdownIAL = ial
			return 1
		}
	}

	if context.blank && (html.HtmlBlockType == 6 || html.HtmlBlockType == 7) {
		return 1
	}
	return 0
}

func (context *Context) htmlBlockFinalize(html *ast.Node) {
	_, html.Tokens = lex.TrimRight(lex.ReplaceNewlineSpace(html.Tokens))
}

var (
	htmlBlockTags1      = [][]byte{util.StrToBytes("<script"), util.StrToBytes("<pre"), util.StrToBytes("<style"), util.StrToBytes("<textarea")}
	htmlBlockCloseTags1 = [][]byte{util.StrToBytes("</script>"), util.StrToBytes("</pre>"), util.StrToBytes("</style>"), util.StrToBytes("</textarea>")}
	htmlBlockTags6      = [][]byte{
		util.StrToBytes("<address"), util.StrToBytes("<article"), util.StrToBytes("<aside"), util.StrToBytes("<base"), util.StrToBytes("<basefont"), util.StrToBytes("<blockquote"), util.StrToBytes("<body"), util.StrToBytes("<caption"), util.StrToBytes("<center"), util.StrToBytes("<col"), util.StrToBytes("<colgroup"), util.StrToBytes("<dd"), util.StrToBytes("<details"), util.StrToBytes("<dialog"), util.StrToBytes("<dir"), util.StrToBytes("<div"), util.StrToBytes("<dl"), util.StrToBytes("<dt"), util.StrToBytes("<fieldset"), util.StrToBytes("<figcaption"), util.StrToBytes("<figure"), util.StrToBytes("<footer"), util.StrToBytes("<form"), util.StrToBytes("<frame"), util.StrToBytes("<frameset"), util.StrToBytes("<h1"), util.StrToBytes("<h2"), util.StrToBytes("<h3"), util.StrToBytes("<h4"), util.StrToBytes("<h5"), util.StrToBytes("<h6"), util.StrToBytes("<head"), util.StrToBytes("<header"), util.StrToBytes("<hr"), util.StrToBytes("<html"), util.StrToBytes("<iframe"), util.StrToBytes("<legend"), util.StrToBytes("<li"), util.StrToBytes("<link"), util.StrToBytes("<main"), util.StrToBytes("<menu"), util.StrToBytes("<menuitem"), util.StrToBytes("<nav"), util.StrToBytes("<noframes"), util.StrToBytes("<ol"), util.StrToBytes("<optgroup"), util.StrToBytes("<option"), util.StrToBytes("<p"), util.StrToBytes("<param"), util.StrToBytes("<section"), util.StrToBytes("<source"), util.StrToBytes("<summary"), util.StrToBytes("<table"), util.StrToBytes("<tbody"), util.StrToBytes("<td"), util.StrToBytes("<tfoot"), util.StrToBytes("<th"), util.StrToBytes("<thead"), util.StrToBytes("<title"), util.StrToBytes("<tr"), util.StrToBytes("<track"), util.StrToBytes("<ul"), util.StrToBytes("<video"), util.StrToBytes("<audio"),
		util.StrToBytes("</address"), util.StrToBytes("</article"), util.StrToBytes("</aside"), util.StrToBytes("</base"), util.StrToBytes("</basefont"), util.StrToBytes("</blockquote"), util.StrToBytes("</body"), util.StrToBytes("</caption"), util.StrToBytes("</center"), util.StrToBytes("</col"), util.StrToBytes("</colgroup"), util.StrToBytes("</dd"), util.StrToBytes("</details"), util.StrToBytes("</dialog"), util.StrToBytes("</dir"), util.StrToBytes("</div"), util.StrToBytes("</dl"), util.StrToBytes("</dt"), util.StrToBytes("</fieldset"), util.StrToBytes("</figcaption"), util.StrToBytes("</figure"), util.StrToBytes("</footer"), util.StrToBytes("</form"), util.StrToBytes("</frame"), util.StrToBytes("</frameset"), util.StrToBytes("</h1"), util.StrToBytes("</h2"), util.StrToBytes("</h3"), util.StrToBytes("</h4"), util.StrToBytes("</h5"), util.StrToBytes("</h6"), util.StrToBytes("</head"), util.StrToBytes("</header"), util.StrToBytes("</hr"), util.StrToBytes("</html"), util.StrToBytes("</iframe"), util.StrToBytes("</legend"), util.StrToBytes("</li"), util.StrToBytes("</link"), util.StrToBytes("</main"), util.StrToBytes("</menu"), util.StrToBytes("</menuitem"), util.StrToBytes("</nav"), util.StrToBytes("</noframes"), util.StrToBytes("</ol"), util.StrToBytes("</optgroup"), util.StrToBytes("</option"), util.StrToBytes("</p"), util.StrToBytes("</param"), util.StrToBytes("</section"), util.StrToBytes("</source"), util.StrToBytes("</summary"), util.StrToBytes("</table"), util.StrToBytes("</tbody"), util.StrToBytes("</td"), util.StrToBytes("</tfoot"), util.StrToBytes("</th"), util.StrToBytes("</thead"), util.StrToBytes("</title"), util.StrToBytes("</tr"), util.StrToBytes("</track"), util.StrToBytes("</ul"), util.StrToBytes("</video"), util.StrToBytes("</audio"),
	}
	htmlBlockSinglequote = util.StrToBytes("'")
	htmlBlockDoublequote = util.StrToBytes("\"")
	htmlBlockGreater     = util.StrToBytes(">")
)

func (t *Tree) isHTMLBlockClose(tokens []byte, htmlType int) bool {
	if t.Context.ParseOption.KramdownBlockIAL && simpleCheckIsBlockIAL(tokens) {
		// 判断 IAL 打断
		if ial := t.Context.parseKramdownBlockIAL(tokens); 0 < len(ial) {
			t.Context.Tip.ID = IAL2Map(ial)["id"]
			t.Context.Tip.KramdownIAL = ial
			t.Context.Tip.InsertAfter(&ast.Node{Type: ast.NodeKramdownBlockIAL, Tokens: tokens})
			return true
		}
	}

	length := len(tokens)
	switch htmlType {
	case 1:
		if pos := lex.AcceptTokenss(tokens, htmlBlockCloseTags1); 0 <= pos {
			return true
		}
		return false
	case 2:
		for i := 0; i < length-3; i++ {
			if lex.ItemHyphen == tokens[i] && lex.ItemHyphen == tokens[i+1] && lex.ItemGreater == tokens[i+2] {
				return true
			}
		}
	case 3:
		for i := 0; i < length-2; i++ {
			if lex.ItemQuestion == tokens[i] && lex.ItemGreater == tokens[i+1] {
				return true
			}
		}
	case 4:
		return bytes.Contains(tokens, htmlBlockGreater)
	case 5:
		for i := 0; i < length-2; i++ {
			if lex.ItemCloseBracket == tokens[i] && lex.ItemCloseBracket == tokens[i+1] {
				return true
			}
		}
	}
	return false
}

func (t *Tree) parseHTML(tokens []byte) (typ int) {
	_, tokens = lex.TrimLeft(tokens)
	length := len(tokens)
	if 3 > length { // at least <? and a newline
		return
	}

	if lex.ItemLess != tokens[0] {
		return
	}

	typ = 1

	if pos := lex.AcceptTokenss(tokens, htmlBlockTags1); 0 <= pos {
		if lex.IsWhitespace(tokens[pos]) || lex.ItemGreater == tokens[pos] {
			return
		}
	}

	if pos := lex.AcceptTokenss(tokens, htmlBlockTags6); 0 <= pos {
		if lex.IsWhitespace(tokens[pos]) || lex.ItemGreater == tokens[pos] {
			typ = 6
			return
		}
		if lex.ItemSlash == tokens[pos] && lex.ItemGreater == tokens[pos+1] {
			typ = 6
			return
		}
	}

	tag := lex.TrimWhitespace(tokens)
	isOpenTag := t.isOpenTag(tag)
	if isOpenTag && t.Context.Tip.Type != ast.NodeParagraph {
		typ = 7
		return
	}
	isCloseTag := t.isCloseTag(tag)
	if isCloseTag && t.Context.Tip.Type != ast.NodeParagraph {
		typ = 7
		return
	}

	if 0 == bytes.Index(tokens, util.StrToBytes("<!--")) {
		typ = 2
		return
	}

	if 0 == bytes.Index(tokens, util.StrToBytes("<?")) {
		typ = 3
		return
	}

	if 2 < len(tokens) && 0 == bytes.Index(tokens, util.StrToBytes("<!")) {
		following := tokens[2:]
		if 'A' <= following[0] && 'Z' >= following[0] {
			typ = 4
			return
		}
		if 0 == bytes.Index(following, util.StrToBytes("[CDATA[")) {
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

	if lex.ItemLess != tokens[0] {
		return
	}
	if lex.ItemGreater != tokens[length-1] {
		return
	}
	if lex.ItemSlash == tokens[length-2] {
		tokens = tokens[1 : length-2]
	} else {
		tokens = tokens[1 : length-1]
	}

	length = len(tokens)
	if 0 == length {
		return
	}

	if lex.IsWhitespace(tokens[0]) { // < 后面不能跟空白
		return
	}

	nameAndAttrs := lex.SplitWhitespace(tokens)
	name := nameAndAttrs[0]
	if !lex.IsASCIILetter(name[0]) {
		return
	}
	if 1 < len(name) {
		name = name[1:]
		for _, n := range name {
			if !lex.IsASCIILetterNumHyphen(n) {
				return
			}
		}
	}

	attrs := nameAndAttrs[1:]
	for _, attr := range attrs {
		if 1 >= len(attr) {
			continue
		}

		nameAndValue := bytes.SplitN(attr, []byte("="), 2)
		name := nameAndValue[0]
		if 1 > len(name) { // 等号前面空格的情况
			continue
		}
		if !lex.IsASCIILetter(name[0]) && lex.ItemUnderscore != name[0] && lex.ItemColon != name[0] {
			return
		}

		if 1 < len(name) {
			name = name[1:]
			for _, n := range name {
				if !lex.IsASCIILetter(n) && !lex.IsDigit(n) && lex.ItemUnderscore != n && lex.ItemDot != n && lex.ItemColon != n && lex.ItemHyphen != n {
					return
				}
			}
		}

		if 1 < len(nameAndValue) {
			value := nameAndValue[1]
			if bytes.HasPrefix(value, htmlBlockSinglequote) && bytes.HasSuffix(value, htmlBlockSinglequote) {
				if value = value[1:]; 1 > len(value) {
					return
				}
				value = value[:len(value)-1]
				return !bytes.Contains(value, htmlBlockSinglequote)
			}
			if bytes.HasPrefix(value, htmlBlockDoublequote) && bytes.HasSuffix(value, htmlBlockDoublequote) {
				if value = value[1:]; 1 > len(value) {
					return
				}
				value = value[:len(value)-1]
				return !bytes.Contains(value, htmlBlockDoublequote)
			}
			return !bytes.ContainsAny(value, " \t\n") && !bytes.ContainsAny(value, "\"'=<>`")
		}
	}
	return true
}

func (t *Tree) isCloseTag(tokens []byte) bool {
	tokens = lex.TrimWhitespace(tokens)
	length := len(tokens)
	if 4 > length {
		return false
	}

	if lex.ItemLess != tokens[0] || lex.ItemSlash != tokens[1] {
		return false
	}
	if lex.ItemGreater != tokens[length-1] {
		return false
	}

	tokens = tokens[2 : length-1]
	length = len(tokens)
	if 0 == length {
		return false
	}

	name := tokens[0:]
	if !lex.IsASCIILetter(name[0]) {
		return false
	}
	if 1 < len(name) {
		name = name[1:]
		for _, n := range name {
			if !lex.IsASCIILetterNumHyphen(n) {
				return false
			}
		}
	}
	return true
}
