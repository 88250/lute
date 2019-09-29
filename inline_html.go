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
	ret = &Node{typ: NodeText, tokens: items{tokens[ctx.pos]}}
	if 3 > ctx.tokensLen || ctx.tokensLen <= startPos+1 {
		ctx.pos++
		return
	}

	var tags items
	tags = append(tags, tokens[startPos])
	if itemSlash == term(tokens[startPos+1]) { // a closing tag
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
	if (itemGreater == term(tokens[0])) ||
		(1 < ctx.tokensLen && itemSlash == term(tokens[0]) && itemGreater == term(tokens[1])) {
		tags = append(tags, whitespaces...)
		tags = append(tags, tokens[0])
		if itemSlash == term(tokens[0]) {
			tags = append(tags, tokens[1])
		}
		ctx.pos += len(tags)
		ret = &Node{typ: NodeInlineHTML, tokens: tags}
		return
	}

	ctx.pos = startPos + 1
	return
}

func (t *Tree) parseCDATA(tokens items) (valid bool, remains, content items) {
	remains = tokens
	if itemBang != term(tokens[0]) {
		return
	}
	if itemOpenBracket != term(tokens[1]) {
		return
	}

	if 'C' != term(tokens[2]) || 'D' != term(tokens[3]) || 'A' != term(tokens[4]) || 'T' != term(tokens[5]) || 'A' != term(tokens[6]) {
		return
	}
	if itemOpenBracket != term(tokens[7]) {
		return
	}

	content = append(content, tokens[:7]...)
	tokens = tokens[7:]
	var token item
	var i int
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		content = append(content, token)
		if i <= length-3 && itemCloseBracket == term(token) && itemCloseBracket == term(tokens[i+1]) && itemGreater == term(tokens[i+2]) {
			break
		}
	}
	tokens = tokens[i:]
	if itemCloseBracket != term(tokens[0]) || itemCloseBracket != term(tokens[1]) || itemGreater != term(tokens[2]) {
		return
	}
	content = append(content, tokens[1], tokens[2])
	valid = true
	remains = tokens[3:]

	return
}

func (t *Tree) parseDeclaration(tokens items) (valid bool, remains, content items) {
	remains = tokens
	if itemBang != term(tokens[0]) {
		return
	}

	var token item
	var i int
	for _, token = range tokens[1:] {
		if isWhitespace(term(token)) {
			break
		}
		if !('A' <= term(token) && 'Z' >= term(token)) {
			return
		}
	}

	content = append(content, tokens[0], tokens[1])
	tokens = tokens[2:]
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		content = append(content, token)
		if itemGreater == term(token) {
			break
		}
	}
	tokens = tokens[i:]
	if itemGreater != term(tokens[0]) {
		return
	}
	valid = true
	remains = tokens[1:]

	return
}

func (t *Tree) parseProcessingInstruction(tokens items) (valid bool, remains, content items) {
	remains = tokens
	if itemQuestion != term(tokens[0]) {
		return
	}

	content = append(content, tokens[0])
	tokens = tokens[1:]
	var token item
	var i int
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		content = append(content, token)
		if i <= length-2 && itemQuestion == term(token) && itemGreater == term(tokens[i+1]) {
			break
		}
	}
	tokens = tokens[i:]
	if 1 > len(tokens) {
		return
	}

	if itemQuestion != term(tokens[0]) || itemGreater != term(tokens[1]) {
		return
	}
	content = append(content, tokens[1])
	valid = true
	remains = tokens[2:]

	return
}

func (t *Tree) parseHTMLComment(tokens items) (valid bool, remains, comment items) {
	remains = tokens
	if itemBang != term(tokens[0]) || itemHyphen != term(tokens[1]) || itemHyphen != term(tokens[2]) {
		return
	}

	comment = append(comment, tokens[0], tokens[1], tokens[2])
	tokens = tokens[3:]
	if itemGreater == term(tokens[0]) {
		return
	}
	if itemHyphen == term(tokens[0]) && itemGreater == term(tokens[1]) {
		return
	}
	var token item
	var i int
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		comment = append(comment, token)
		if i <= length-2 && itemHyphen == term(token) && itemHyphen == term(tokens[i+1]) {
			break
		}
		if i <= length-3 && itemHyphen == term(token) && itemHyphen == term(tokens[i+1]) && itemGreater == term(tokens[i+2]) {
			break
		}
	}
	tokens = tokens[i:]
	if itemHyphen != term(tokens[0]) || itemHyphen != term(tokens[1]) || itemGreater != term(tokens[2]) {
		return
	}
	comment = append(comment, tokens[1], tokens[2])
	valid = true
	remains = tokens[3:]

	return
}

func (t *Tree) parseTagAttr(tokens items) (valid bool, remains, attr items) {
	valid = true
	remains = tokens
	var whitespaces items
	var i int
	var token item
	for i, token = range tokens {
		if !isWhitespace(term(token)) {
			break
		}
		whitespaces = append(whitespaces, token)
	}
	if 1 > len(whitespaces) {
		return
	}
	tokens = tokens[i:]

	var attrName items
	tokens, attrName = t.parseAttrName(tokens)
	if 1 > len(attrName) {
		return
	}

	var valSpec items
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

func (t *Tree) parseAttrValSpec(tokens items) (valid bool, remains, valSpec items) {
	valid = true
	remains = tokens
	var i int
	var token item
	for i, token = range tokens {
		if !isWhitespace(term(token)) {
			break
		}
		valSpec = append(valSpec, token)
	}
	if itemEqual != term(token) {
		valSpec = nil
		return
	}
	valSpec = append(valSpec, token)
	tokens = tokens[i+1:]
	for i, token = range tokens {
		if !isWhitespace(term(token)) {
			break
		}
		valSpec = append(valSpec, token)
	}
	token = tokens[i]
	valSpec = append(valSpec, token)
	tokens = tokens[i+1:]
	closed := false
	if itemDoublequote == term(token) { // A double-quoted attribute value consists of ", zero or more characters not including ", and a final ".
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if itemDoublequote == term(token) {
				closed = true
				break
			}
		}
	} else if itemSinglequote == term(token) { // A single-quoted attribute value consists of ', zero or more characters not including ', and a final '.
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if itemSinglequote == term(token) {
				closed = true
				break
			}
		}
	} else { // An unquoted attribute value is a nonempty string of characters not including whitespace, ", ', =, <, >, or `.
		for i, token = range tokens {
			if itemGreater == term(token) {
				i-- // 大于字符 > 不计入 valSpec
				break
			}
			valSpec = append(valSpec, token)
			if isWhitespace(term(token)) {
				// 属性使用空白分隔
				break
			}
			if itemDoublequote == term(token) || itemSinglequote == term(token) || itemEqual == term(token) || itemLess == term(token) || itemGreater == term(token) || itemBacktick == term(token) {
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

func (t *Tree) parseAttrName(tokens items) (remains, attrName items) {
	remains = tokens
	if !isASCIILetter(term(tokens[0])) && itemUnderscore != term(tokens[0]) && itemColon != term(tokens[0]) {
		return
	}
	attrName = append(attrName, tokens[0])
	tokens = tokens[1:]
	var i int
	var token item
	for i, token = range tokens {
		if !isASCIILetterNumHyphen(term(token)) && itemUnderscore != term(token) && itemDot != term(token) && itemColon != term(token) {
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

func (t *Tree) parseTagName(tokens items) (remains, tagName items) {
	i := 0
	token := tokens[i]
	if !isASCIILetter(term(token)) {
		return tokens, nil
	}
	tagName = append(tagName, token)
	for i = 1; i < len(tokens); i++ {
		token = tokens[i]
		if !isASCIILetterNumHyphen(term(token)) {
			break
		}
		tagName = append(tagName, token)
	}
	remains = tokens[i:]

	return
}
