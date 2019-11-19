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

func (t *Tree) parseInlineHTML(ctx *InlineContext) (ret *Node) {
	tokens := ctx.tokens
	startPos := ctx.pos
	ret = &Node{typ: NodeText, tokens: []byte{tokens[ctx.pos]}}
	if 3 > ctx.tokensLen || ctx.tokensLen <= startPos+1 {
		ctx.pos++
		return
	}

	var tags []byte
	tags = append(tags, tokens[startPos])
	if itemSlash == tokens[startPos+1] { // a closing tag
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
			valid, remains, attr := t.parseTagAttr(tokens)
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
		ret = &Node{typ: NodeInlineHTML, tokens: tags}
		return
	} else if valid, remains, ins := t.parseProcessingInstruction(tokens[ctx.pos+1:]); valid {
		tags = append(tags, ins...)
		tokens = remains
		ctx.pos += len(tags)
		ret = &Node{typ: NodeInlineHTML, tokens: tags}
		return
	} else if valid, remains, decl := t.parseDeclaration(tokens[ctx.pos+1:]); valid {
		tags = append(tags, decl...)
		tokens = remains
		ctx.pos += len(tags)
		ret = &Node{typ: NodeInlineHTML, tokens: tags}
		return
	} else if valid, remains, cdata := t.parseCDATA(tokens[ctx.pos+1:]); valid {
		tags = append(tags, cdata...)
		tokens = remains
		ctx.pos += len(tags)
		ret = &Node{typ: NodeInlineHTML, tokens: tags}
		return
	} else {
		ctx.pos++
		return
	}

	length := len(tokens)
	if 1 > length {
		ctx.pos = startPos + 1
		return
	}

	whitespaces, tokens := trimLeft(tokens)
	if (itemGreater == tokens[0]) ||
		(1 < ctx.tokensLen && itemSlash == tokens[0] && itemGreater == tokens[1]) {
		tags = append(tags, whitespaces...)
		tags = append(tags, tokens[0])
		if itemSlash == tokens[0] {
			tags = append(tags, tokens[1])
		}
		ctx.pos += len(tags)
		ret = &Node{typ: NodeInlineHTML, tokens: tags}
		return
	}

	ctx.pos = startPos + 1
	return
}

func (t *Tree) parseCDATA(tokens []byte) (valid bool, remains, content []byte) {
	remains = tokens
	if itemBang != tokens[0] {
		return
	}
	if itemOpenBracket != tokens[1] {
		return
	}

	if 'C' != tokens[2] || 'D' != tokens[3] || 'A' != tokens[4] || 'T' != tokens[5] || 'A' != tokens[6] {
		return
	}
	if itemOpenBracket != tokens[7] {
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
		if i <= length-3 && itemCloseBracket == token && itemCloseBracket == tokens[i+1] && itemGreater == tokens[i+2] {
			break
		}
	}
	tokens = tokens[i:]
	if itemCloseBracket != tokens[0] || itemCloseBracket != tokens[1] || itemGreater != tokens[2] {
		return
	}
	content = append(content, tokens[1], tokens[2])
	valid = true
	remains = tokens[3:]

	return
}

func (t *Tree) parseDeclaration(tokens []byte) (valid bool, remains, content []byte) {
	remains = tokens
	if itemBang != tokens[0] {
		return
	}

	var token byte
	var i int
	for _, token = range tokens[1:] {
		if isWhitespace(token) {
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
		if itemGreater == token {
			break
		}
	}
	tokens = tokens[i:]
	if itemGreater != tokens[0] {
		return
	}
	valid = true
	remains = tokens[1:]

	return
}

func (t *Tree) parseProcessingInstruction(tokens []byte) (valid bool, remains, content []byte) {
	remains = tokens
	if itemQuestion != tokens[0] {
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
		if i <= length-2 && itemQuestion == token && itemGreater == tokens[i+1] {
			break
		}
	}
	tokens = tokens[i:]
	if 1 > len(tokens) {
		return
	}

	if itemQuestion != tokens[0] || itemGreater != tokens[1] {
		return
	}
	content = append(content, tokens[1])
	valid = true
	remains = tokens[2:]

	return
}

func (t *Tree) parseHTMLComment(tokens []byte) (valid bool, remains, comment []byte) {
	remains = tokens
	if itemBang != tokens[0] || itemHyphen != tokens[1] || itemHyphen != tokens[2] {
		return
	}

	comment = append(comment, tokens[0], tokens[1], tokens[2])
	tokens = tokens[3:]
	if itemGreater == tokens[0] {
		return
	}
	if itemHyphen == tokens[0] && itemGreater == tokens[1] {
		return
	}
	var token byte
	var i int
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		comment = append(comment, token)
		if i <= length-2 && itemHyphen == token && itemHyphen == tokens[i+1] {
			break
		}
		if i <= length-3 && itemHyphen == token && itemHyphen == tokens[i+1] && itemGreater == tokens[i+2] {
			break
		}
	}
	tokens = tokens[i:]
	if itemHyphen != tokens[0] || itemHyphen != tokens[1] || itemGreater != tokens[2] {
		return
	}
	comment = append(comment, tokens[1], tokens[2])
	valid = true
	remains = tokens[3:]

	return
}

func (t *Tree) parseTagAttr(tokens []byte) (valid bool, remains, attr []byte) {
	valid = true
	remains = tokens
	var whitespaces []byte
	var i int
	var token byte
	for i, token = range tokens {
		if !isWhitespace(token) {
			break
		}
		whitespaces = append(whitespaces, token)
	}
	if 1 > len(whitespaces) {
		return
	}
	tokens = tokens[i:]

	var attrName []byte
	tokens, attrName = t.parseAttrName(tokens)
	if 1 > len(attrName) {
		return
	}

	var valSpec []byte
	valid, tokens, valSpec = t.parseAttrValSpec(tokens)
	if !valid {
		return
	}

	remains = tokens
	attr = append(attr, whitespaces...)
	attr = append(attr, attrName...)
	attr = append(attr, valSpec...)

	return
}

func (t *Tree) parseAttrValSpec(tokens []byte) (valid bool, remains, valSpec []byte) {
	valid = true
	remains = tokens
	var i int
	var token byte
	for i, token = range tokens {
		if !isWhitespace(token) {
			break
		}
		valSpec = append(valSpec, token)
	}
	if itemEqual != token {
		valSpec = nil
		return
	}
	valSpec = append(valSpec, token)
	tokens = tokens[i+1:]
	for i, token = range tokens {
		if !isWhitespace(token) {
			break
		}
		valSpec = append(valSpec, token)
	}
	token = tokens[i]
	valSpec = append(valSpec, token)
	tokens = tokens[i+1:]
	closed := false
	if itemDoublequote == token { // A double-quoted attribute value consists of ", zero or more characters not including ", and a final ".
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if itemDoublequote == token {
				closed = true
				break
			}
		}
	} else if itemSinglequote == token { // A single-quoted attribute value consists of ', zero or more characters not including ', and a final '.
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if itemSinglequote == token {
				closed = true
				break
			}
		}
	} else { // An unquoted attribute value is a nonempty string of characters not including whitespace, ", ', =, <, >, or `.
		for i, token = range tokens {
			if itemGreater == token {
				i-- // 大于字符 > 不计入 valSpec
				break
			}
			valSpec = append(valSpec, token)
			if isWhitespace(token) {
				// 属性使用空白分隔
				break
			}
			if itemDoublequote == token || itemSinglequote == token || itemEqual == token || itemLess == token || itemGreater == token || itemBacktick == token {
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

func (t *Tree) parseAttrName(tokens []byte) (remains, attrName []byte) {
	remains = tokens
	if !isASCIILetter(tokens[0]) && itemUnderscore != tokens[0] && itemColon != tokens[0] {
		return
	}
	attrName = append(attrName, tokens[0])
	tokens = tokens[1:]
	var i int
	var token byte
	for i, token = range tokens {
		if !isASCIILetterNumHyphen(token) && itemUnderscore != token && itemDot != token && itemColon != token {
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
	i := 0
	token := tokens[i]
	if !isASCIILetter(token) {
		return tokens, nil
	}
	tagName = append(tagName, token)
	for i = 1; i < len(tokens); i++ {
		token = tokens[i]
		if !isASCIILetterNumHyphen(token) {
			break
		}
		tagName = append(tagName, token)
	}
	remains = tokens[i:]

	return
}
