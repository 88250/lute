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
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

func (t *Tree) parseInlineHTML(ctx *InlineContext) (ret *ast.Node) {
	tokens := ctx.tokens
	caretInTag := false
	caretLeftSpace := false
	if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG {
		caretIndex := bytes.Index(tokens, editor.CaretTokens)
		caretInTag = caretIndex > ctx.pos
		if caretInTag {
			caretLeftSpace = bytes.Contains(tokens, []byte(" "+editor.Caret))
			tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, []byte(editor.CaretReplacement))
			tokens = bytes.ReplaceAll(tokens, []byte("\""+editor.CaretReplacement), []byte("\" "+editor.CaretReplacement))
		}
	}

	startPos := ctx.pos
	ret = &ast.Node{Type: ast.NodeText, Tokens: []byte{tokens[ctx.pos]}}
	if 3 > ctx.tokensLen || ctx.tokensLen <= startPos+1 {
		ctx.pos++
		return
	}

	var tags []byte
	tags = append(tags, tokens[startPos])
	if lex.ItemSlash == tokens[startPos+1] && 1 < ctx.tokensLen-(startPos+1) { // a closing tag
		tags = append(tags, tokens[startPos+1])
		remains, tagName := t.parseTagName(tokens[ctx.pos+2:])
		if 1 > len(tagName) {
			ctx.pos++
			return
		}
		tags = append(tags, tagName...)
		tokens = remains
	} else if remains, tagName := t.parseTagName(tokens[ctx.pos+1:]); 0 < len(tagName) {
		tags = append(tags, tagName...)
		tokens = remains
		for {
			valid, remains, attr, _, _ := TagAttr(tokens)
			if !valid {
				ctx.pos++
				return
			}

			tokens = remains
			tags = append(tags, attr...)
			if 1 > len(attr) {
				break
			}
		}
	} else if valid, remains, comment := t.parseHTMLComment(tokens[ctx.pos+1:]); valid {
		tags = append(tags, comment...)
		tokens = remains
		ctx.pos += len(tags)
		ret = &ast.Node{Type: ast.NodeInlineHTML, Tokens: tags}
		return
	} else if valid, remains, ins := t.parseProcessingInstruction(tokens[ctx.pos+1:]); valid {
		tags = append(tags, ins...)
		tokens = remains
		ctx.pos += len(tags)
		ret = &ast.Node{Type: ast.NodeInlineHTML, Tokens: tags}
		return
	} else if valid, remains, decl := t.parseDeclaration(tokens[ctx.pos+1:]); valid {
		tags = append(tags, decl...)
		tokens = remains
		ctx.pos += len(tags)
		ret = &ast.Node{Type: ast.NodeInlineHTML, Tokens: tags}
		return
	} else if valid, remains, cdata := t.parseCDATA(tokens[ctx.pos+1:]); valid {
		tags = append(tags, cdata...)
		tokens = remains
		ctx.pos += len(tags)
		ret = &ast.Node{Type: ast.NodeInlineHTML, Tokens: tags}
		return
	} else {
		ctx.pos++
		return
	}

	whitespaces, tokens := lex.TrimLeft(tokens)
	length := len(tokens)
	if 1 > length {
		ctx.pos = startPos + 1
		return
	}

	if (lex.ItemGreater == tokens[0]) ||
		(1 < length && lex.ItemSlash == tokens[0] && lex.ItemGreater == tokens[1]) {
		tags = append(tags, whitespaces...)
		tags = append(tags, tokens[0])
		if lex.ItemSlash == tokens[0] {
			tags = append(tags, tokens[1])
		}
		if (t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV) && caretInTag || t.Context.ParseOption.ProtyleWYSIWYG {
			if !bytes.Contains(tags, []byte(editor.CaretReplacement+" ")) && !caretLeftSpace {
				tags = bytes.ReplaceAll(tags, []byte("\" "+editor.CaretReplacement), []byte("\""+editor.CaretReplacement))
			}
			tags = bytes.ReplaceAll(tags, []byte(editor.CaretReplacement), editor.CaretTokens)
		}
		ctx.pos += len(tags)

		if t.Context.ParseOption.ProtyleWYSIWYG {
			if bytes.Equal(tags, []byte("<br />")) || bytes.Equal(tags, []byte("<br/>")) {
				ret = &ast.Node{Type: ast.NodeBr}
				return
			} else if bytes.HasPrefix(tags, []byte("<span data-type=")) {
				ret = t.processSpanTag(startPos, tags, "<span data-type=", "</span>", ctx)
				return
			} else if bytes.Equal(tags, []byte("<kbd>")) {
				ret = t.processSpanTag(startPos, tags, "<kbd>", "</kbd>", ctx)
				return
			} else if bytes.Equal(tags, []byte("<u>")) {
				ret = t.processSpanTag(startPos, tags, "<u>", "</u>", ctx)
				return
			}
		}
		ret = &ast.Node{Type: ast.NodeInlineHTML, Tokens: tags}
		return
	}

	ctx.pos = startPos + 1
	return
}

func (t *Tree) processSpanTag(startPos int, tags []byte, startTag, endTag string, ctx *InlineContext) (ret *ast.Node) {
	remains := ctx.tokens[ctx.pos:]
	end := bytes.Index(remains, []byte(endTag))
	closerLen := len(endTag)
	tokens := append(tags, remains[:end+closerLen]...)
	nodes, _ := html.ParseFragment(bytes.NewReader(tokens), &html.Node{Type: html.ElementNode})
	if 1 != len(nodes) {
		ctx.pos = startPos + 1
		return
	}
	node := nodes[0]

	var typ string
	startTagLen := len(startTag)
	if "<kbd>" == startTag || "<u>" == startTag {
		if !t.Context.ParseOption.HTMLTag2TextMark {
			ret = &ast.Node{Type: ast.NodeInlineHTML, Tokens: tags}
			return
		}
		typ = node.Data
	} else { // <span data-type="a">
		typ = string(tags[startTagLen+1:])
		typ = typ[:strings.Index(typ, "\"")]
	}
	ret = &ast.Node{Type: ast.NodeTextMark, TextMarkType: typ}
	SetTextMarkNode(ret, node, t.Context.ParseOption)
	ctx.pos += end + closerLen
	return
}

func (t *Tree) parseCDATA(tokens []byte) (valid bool, remains, content []byte) {
	remains = tokens
	if 8 > len(tokens) {
		return
	}
	if lex.ItemBang != tokens[0] {
		return
	}
	if lex.ItemOpenBracket != tokens[1] {
		return
	}

	if 'C' != tokens[2] || 'D' != tokens[3] || 'A' != tokens[4] || 'T' != tokens[5] || 'A' != tokens[6] {
		return
	}
	if lex.ItemOpenBracket != tokens[7] {
		return
	}

	content = append(content, tokens[:7]...)
	tokens = tokens[7:]
	var token byte
	var i int
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		content = append(content, token)
		if i <= length-3 && lex.ItemCloseBracket == token && lex.ItemCloseBracket == tokens[i+1] && lex.ItemGreater == tokens[i+2] {
			break
		}
	}
	tokens = tokens[i:]
	if 3 > len(tokens) {
		return
	}
	if lex.ItemCloseBracket != tokens[0] || lex.ItemCloseBracket != tokens[1] || lex.ItemGreater != tokens[2] {
		return
	}
	content = append(content, tokens[1], tokens[2])
	valid = true
	remains = tokens[3:]
	return
}

func (t *Tree) parseDeclaration(tokens []byte) (valid bool, remains, content []byte) {
	remains = tokens
	if 2 > len(tokens) {
		return
	}

	if lex.ItemBang != tokens[0] {
		return
	}

	var token byte
	var i int
	for _, token = range tokens[1:] {
		if lex.IsWhitespace(token) {
			break
		}
		if !('A' <= token && 'Z' >= token) {
			return
		}
	}

	content = append(content, tokens[0], tokens[1])
	tokens = tokens[2:]
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		content = append(content, token)
		if lex.ItemGreater == token {
			break
		}
	}
	tokens = tokens[i:]
	if 1 > len(tokens) || lex.ItemGreater != tokens[0] {
		return
	}
	valid = true
	remains = tokens[1:]
	return
}

func (t *Tree) parseProcessingInstruction(tokens []byte) (valid bool, remains, content []byte) {
	remains = tokens
	if lex.ItemQuestion != tokens[0] {
		return
	}

	content = append(content, tokens[0])
	tokens = tokens[1:]
	var token byte
	var i int
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		content = append(content, token)
		if i <= length-2 && lex.ItemQuestion == token && lex.ItemGreater == tokens[i+1] {
			break
		}
	}
	tokens = tokens[i:]
	if 1 > len(tokens) {
		return
	}

	if lex.ItemQuestion != tokens[0] || lex.ItemGreater != tokens[1] {
		return
	}
	content = append(content, tokens[1])
	valid = true
	remains = tokens[2:]
	return
}

func (t *Tree) parseHTMLComment(tokens []byte) (valid bool, remains, comment []byte) {
	remains = tokens
	if 3 > len(tokens) {
		return
	}

	if lex.ItemBang != tokens[0] || lex.ItemHyphen != tokens[1] || lex.ItemHyphen != tokens[2] {
		return
	}

	comment = append(comment, tokens[0], tokens[1], tokens[2])
	tokens = tokens[3:]
	length := len(tokens)
	if 2 > length {
		return
	}
	if lex.ItemGreater == tokens[0] {
		return
	}
	if lex.ItemHyphen == tokens[0] && lex.ItemGreater == tokens[1] {
		return
	}
	var token byte
	var i int
	for ; i < length; i++ {
		token = tokens[i]
		comment = append(comment, token)
		if i <= length-2 && lex.ItemHyphen == token && lex.ItemHyphen == tokens[i+1] {
			break
		}
		if i <= length-3 && lex.ItemHyphen == token && lex.ItemHyphen == tokens[i+1] && lex.ItemGreater == tokens[i+2] {
			break
		}
	}
	tokens = tokens[i:]
	if 3 > len(tokens) || lex.ItemHyphen != tokens[0] || lex.ItemHyphen != tokens[1] || lex.ItemGreater != tokens[2] {
		return
	}
	comment = append(comment, tokens[1], tokens[2])
	valid = true
	remains = tokens[3:]
	return
}

func TagAttr(tokens []byte) (valid bool, remains, attr, name, val []byte) {
	valid = true
	remains = tokens
	var whitespaces []byte
	var i int
	var token byte
	for i, token = range tokens {
		if !lex.IsWhitespace(token) {
			break
		}
		whitespaces = append(whitespaces, token)
	}
	if 1 > len(whitespaces) {
		return
	}
	tokens = tokens[i:]

	var attrName []byte
	tokens, attrName = parseAttrName(tokens)
	if 1 > len(attrName) {
		return
	}

	var valSpec []byte
	valid, tokens, valSpec = parseAttrValSpec(tokens)
	if !valid {
		return
	}

	remains = tokens
	attr = append(attr, whitespaces...)
	attr = append(attr, attrName...)
	attr = append(attr, valSpec...)
	if nil != valSpec {
		name = attrName
		val = valSpec[2 : len(valSpec)-1]
	}
	return
}

func parseAttrValSpec(tokens []byte) (valid bool, remains, valSpec []byte) {
	valid = true
	remains = tokens
	var i int
	var token byte
	for i, token = range tokens {
		if !lex.IsWhitespace(token) {
			break
		}
		valSpec = append(valSpec, token)
	}
	if lex.ItemEqual != token {
		valSpec = nil
		return
	}
	valSpec = append(valSpec, token)
	tokens = tokens[i+1:]
	if 1 > len(tokens) {
		valid = false
		return
	}

	for i, token = range tokens {
		if !lex.IsWhitespace(token) {
			break
		}
		valSpec = append(valSpec, token)
	}
	token = tokens[i]
	valSpec = append(valSpec, token)
	tokens = tokens[i+1:]
	closed := false
	if lex.ItemDoublequote == token { // A double-quoted attribute value consists of ", zero or more characters not including ", and a final ".
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if lex.ItemDoublequote == token {
				closed = true
				break
			}
		}
	} else if lex.ItemSinglequote == token { // A single-quoted attribute value consists of ', zero or more characters not including ', and a final '.
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if lex.ItemSinglequote == token {
				closed = true
				break
			}
		}
	} else { // An unquoted attribute value is a nonempty string of characters not including whitespace, ", ', =, <, >, or `.
		for i, token = range tokens {
			if lex.ItemGreater == token {
				i-- // 大于字符 > 不计入 valSpec
				break
			}
			valSpec = append(valSpec, token)
			if lex.IsWhitespace(token) {
				// 属性使用空白分隔
				break
			}
			if lex.ItemDoublequote == token || lex.ItemSinglequote == token || lex.ItemEqual == token || lex.ItemLess == token || lex.ItemGreater == token || lex.ItemBacktick == token {
				closed = false
				break
			}
			closed = true
		}
	}

	if !closed {
		valid = false
		valSpec = nil
		return
	}

	remains = tokens[i+1:]
	return
}

func parseAttrName(tokens []byte) (remains, attrName []byte) {
	remains = tokens
	if !lex.IsASCIILetter(tokens[0]) && lex.ItemUnderscore != tokens[0] && lex.ItemColon != tokens[0] {
		return
	}
	attrName = append(attrName, tokens[0])
	tokens = tokens[1:]
	var i int
	var token byte
	for i, token = range tokens {
		if !lex.IsASCIILetterNumHyphen(token) && lex.ItemUnderscore != token && lex.ItemDot != token && lex.ItemColon != token {
			break
		}
		attrName = append(attrName, token)
	}
	if 1 > len(attrName) {
		return
	}

	remains = tokens[i:]
	return
}

func (t *Tree) parseTagName(tokens []byte) (remains, tagName []byte) {
	if 1 > len(tokens) {
		return
	}

	i := 0
	token := tokens[i]
	if !lex.IsASCIILetter(token) {
		return tokens, nil
	}
	tagName = append(tagName, token)
	for i = 1; i < len(tokens); i++ {
		token = tokens[i]
		if !lex.IsASCIILetterNumHyphen(token) {
			break
		}
		tagName = append(tagName, token)
	}
	remains = tokens[i:]
	return
}

func SetSpanIAL(node *ast.Node, n *html.Node) {
	if nil == node || nil == n {
		return
	}

	insertedIAL := false
	if style := util.DomAttrValue(n, "style"); "" != style { // 比如设置表格列宽
		style = StyleValue(style)
		node.SetIALAttr("style", style)
		ialTokens := IAL2Tokens(node.KramdownIAL)
		ial := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
		if ast.NodeTableCell == node.Type {
			node.PrependChild(ial)
		} else {
			node.InsertAfter(ial)
		}
		insertedIAL = true
	}

	if atom.Th == n.DataAtom || atom.Td == n.DataAtom {
		// 设置表格合并单元格
		colspan := util.DomAttrValue(n, "colspan")
		if "" != colspan {
			node.SetIALAttr("colspan", colspan)
		}
		rowspan := util.DomAttrValue(n, "rowspan")
		if "" != rowspan {
			node.SetIALAttr("rowspan", rowspan)
		}
		class := util.DomAttrValue(n, "class")
		if "" != class {
			node.SetIALAttr("class", class)
		}
		if "" != colspan || "" != rowspan || "" != class {
			ialTokens := IAL2Tokens(node.KramdownIAL)
			if !insertedIAL {
				ial := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
				node.PrependChild(ial)
				insertedIAL = true
			} else {
				// 合并这两个 IAL
				node.FirstChild.Tokens = IAL2Tokens(node.KramdownIAL)
			}
		}
	}

	if nil != n.Parent && nil != n.Parent.Parent {
		if parentStyle := util.DomAttrValue(n.Parent.Parent, "style"); "" != parentStyle {
			if insertedIAL {
				m := Tokens2IAL(node.Next.Tokens)
				m = append(m, []string{"parent-style", parentStyle})
				node.Next.Tokens = IAL2Tokens(m)
				node.SetIALAttr("parent-style", parentStyle)
				node.KramdownIAL = m
			} else {
				node.SetIALAttr("parent-style", parentStyle)
				ialTokens := IAL2Tokens(node.KramdownIAL)
				ial := &ast.Node{Type: ast.NodeKramdownSpanIAL, Tokens: ialTokens}
				node.InsertAfter(ial)
			}
		}
	}
}

func ContainTextMark(node *ast.Node, dataTypes ...string) bool {
	parts := strings.Split(node.TextMarkType, " ")
	for _, typ := range parts {
		for _, dataType := range dataTypes {
			if typ == dataType {
				return true
			}
		}
	}
	return false
}

func SetTextMarkNode(node *ast.Node, n *html.Node, options *Options) {
	node.Type = ast.NodeTextMark
	dataType := util.DomAttrValue(n, "data-type")
	if "" == dataType {
		if n.DataAtom == atom.Span {
			dataType = "text"
		} else {
			dataType = n.DataAtom.String()
		}
	}
	node.TextMarkType = dataType
	node.Tokens = nil
	types := strings.Split(dataType, " ")
	isInlineMath := false
	for _, typ := range types {
		switch typ {
		case "a":
			node.TextMarkAHref, node.TextMarkATitle = util.GetTextMarkAData(n)
			node.TextMarkTextContent = util.GetTextMarkTextData(n)
		case "inline-math":
			node.TextMarkInlineMathContent = util.GetTextMarkInlineMathData(n)
			isInlineMath = true
		case "block-ref":
			node.TextMarkBlockRefID, node.TextMarkBlockRefSubtype = util.GetTextMarkBlockRefData(n)
			node.TextMarkTextContent = util.GetTextMarkTextData(n)
		case "file-annotation-ref":
			node.TextMarkFileAnnotationRefID = util.GetTextMarkFileAnnotationRefData(n)
			node.TextMarkTextContent = util.GetTextMarkTextData(n)
		case "inline-memo":
			node.TextMarkTextContent = util.GetTextMarkTextData(n)
			node.TextMarkInlineMemoContent = util.GetTextMarkInlineMemoData(n)
			inlineTree := Inline("", []byte(node.TextMarkInlineMemoContent), options)
			if nil != inlineTree {
				node.TextMarkInlineMemoContent = strings.ReplaceAll(inlineTree.Root.Content(), "\n", "")
				node.TextMarkInlineMemoContent = strings.ReplaceAll(node.TextMarkInlineMemoContent, "\"", "&quot;")
			}
		default:
			if !isInlineMath { // 带有字体样式的公式复制之后内容不正确 https://github.com/siyuan-note/siyuan/issues/6799
				node.TextMarkTextContent = util.GetTextMarkTextDataWithoutEscapeSingleQuote(n)
			}
		}
	}

	SetSpanIAL(node, n)
}

func StyleValue(style string) (ret string) {
	ret = strings.TrimSpace(style)
	ret = strings.ReplaceAll(ret, "\n", "")
	ret = strings.Join(strings.Fields(ret), " ")
	return
}
