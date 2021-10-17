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

func (t *Tree) parseInlineHTML(ctx *InlineContext) (ret *ast.Node) {
	tokens := ctx.tokens
	caretInTag := false
	caretLeftSpace := false
	if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG {
		caretIndex := bytes.Index(tokens, util.CaretTokens)
		caretInTag = caretIndex > ctx.pos
		if caretInTag {
			caretLeftSpace = bytes.Contains(tokens, []byte(" "+util.Caret))
			tokens = bytes.ReplaceAll(tokens, util.CaretTokens, []byte(util.CaretReplacement))
			tokens = bytes.ReplaceAll(tokens, []byte("\""+util.CaretReplacement), []byte("\" "+util.CaretReplacement))
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
			if !bytes.Contains(tags, []byte(util.CaretReplacement+" ")) && !caretLeftSpace {
				tags = bytes.ReplaceAll(tags, []byte("\" "+util.CaretReplacement), []byte("\""+util.CaretReplacement))
			}
			tags = bytes.ReplaceAll(tags, []byte(util.CaretReplacement), util.CaretTokens)
		}
		ctx.pos += len(tags)

		if t.Context.ParseOption.ProtyleWYSIWYG {
			if bytes.Equal(tags, []byte("<kbd>")) {
				ret = &ast.Node{Type: ast.NodeKbd}
				ret.AppendChild(&ast.Node{Type: ast.NodeKbdOpenMarker})
				return
			} else if bytes.Equal(tags, []byte("</kbd>")) {
				ret = &ast.Node{Type: ast.NodeKbdCloseMarker}
				return
			} else if bytes.Equal(tags, []byte("<u>")) {
				ret = &ast.Node{Type: ast.NodeUnderline}
				ret.AppendChild(&ast.Node{Type: ast.NodeUnderlineOpenMarker})
				return
			} else if bytes.Equal(tags, []byte("</u>")) {
				ret = &ast.Node{Type: ast.NodeUnderlineCloseMarker}
				return
			} else if bytes.Equal(tags, []byte("<br />")) || bytes.Equal(tags, []byte("<br/>")) {
				ret = &ast.Node{Type: ast.NodeBr}
				return
			} else if bytes.HasPrefix(tags, []byte("<span data-type=")) {
				ret = &ast.Node{Type: ast.NodeTextMark}
				typ := tags[len("<span data-type=")+1:]
				typ = typ[:bytes.Index(typ, []byte("\""))]
				ret.AppendChild(&ast.Node{Type: ast.NodeTextMarkOpenMarker, Tokens: typ})
				return
			} else if bytes.Equal(tags, []byte("</span>")) {
				ret = &ast.Node{Type: ast.NodeTextMarkCloseMarker}
				return
			}
		}
		ret = &ast.Node{Type: ast.NodeInlineHTML, Tokens: tags}
		return
	}

	ctx.pos = startPos + 1
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
