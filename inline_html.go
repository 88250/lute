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
	if itemSlash == tokens[startPos+1].term() { // a closing tag
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
	if (itemGreater == tokens[0].term()) ||
		(1 < ctx.tokensLen && itemSlash == tokens[0].term() && itemGreater == tokens[1].term()) {
		tags = append(tags, whitespaces...)
		tags = append(tags, tokens[0])
		if itemSlash == tokens[0].term() {
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
	if itemBang != tokens[0].term() {
		return
	}
	if itemOpenBracket != tokens[1].term() {
		return
	}

	if 'C' != tokens[2].term() || 'D' != tokens[3].term() || 'A' != tokens[4].term() || 'T' != tokens[5].term() || 'A' != tokens[6].term() {
		return
	}
	if itemOpenBracket != tokens[7].term() {
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
		if i <= length-3 && itemCloseBracket == token.term() && itemCloseBracket == tokens[i+1].term() && itemGreater == tokens[i+2].term() {
			break
		}
	}
	tokens = tokens[i:]
	if itemCloseBracket != tokens[0].term() || itemCloseBracket != tokens[1].term() || itemGreater != tokens[2].term() {
		return
	}
	content = append(content, tokens[1], tokens[2])
	valid = true
	remains = tokens[3:]

	return
}

func (t *Tree) parseDeclaration(tokens items) (valid bool, remains, content items) {
	remains = tokens
	if itemBang != tokens[0].term() {
		return
	}

	var token item
	var i int
	for _, token = range tokens[1:] {
		if isWhitespace(token.term()) {
			break
		}
		if !('A' <= token.term() && 'Z' >= token.term()) {
			return
		}
	}

	content = append(content, tokens[0], tokens[1])
	tokens = tokens[2:]
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		content = append(content, token)
		if itemGreater == token.term() {
			break
		}
	}
	tokens = tokens[i:]
	if itemGreater != tokens[0].term() {
		return
	}
	valid = true
	remains = tokens[1:]

	return
}

func (t *Tree) parseProcessingInstruction(tokens items) (valid bool, remains, content items) {
	remains = tokens
	if itemQuestion != tokens[0].term() {
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
		if i <= length-2 && itemQuestion == token.term() && itemGreater == tokens[i+1].term() {
			break
		}
	}
	tokens = tokens[i:]
	if 1 > len(tokens) {
		return
	}

	if itemQuestion != tokens[0].term() || itemGreater != tokens[1].term() {
		return
	}
	content = append(content, tokens[1])
	valid = true
	remains = tokens[2:]

	return
}

func (t *Tree) parseHTMLComment(tokens items) (valid bool, remains, comment items) {
	remains = tokens
	if itemBang != tokens[0].term() || itemHyphen != tokens[1].term() || itemHyphen != tokens[2].term() {
		return
	}

	comment = append(comment, tokens[0], tokens[1], tokens[2])
	tokens = tokens[3:]
	if itemGreater == tokens[0].term() {
		return
	}
	if itemHyphen == tokens[0].term() && itemGreater == tokens[1].term() {
		return
	}
	var token item
	var i int
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		comment = append(comment, token)
		if i <= length-2 && itemHyphen == token.term() && itemHyphen == tokens[i+1].term() {
			break
		}
		if i <= length-3 && itemHyphen == token.term() && itemHyphen == tokens[i+1].term() && itemGreater == tokens[i+2].term() {
			break
		}
	}
	tokens = tokens[i:]
	if itemHyphen != tokens[0].term() || itemHyphen != tokens[1].term() || itemGreater != tokens[2].term() {
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
		if !isWhitespace(token.term()) {
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
		if !isWhitespace(token.term()) {
			break
		}
		valSpec = append(valSpec, token)
	}
	if itemEqual != token.term() {
		valSpec = nil
		return
	}
	valSpec = append(valSpec, token)
	tokens = tokens[i+1:]
	for i, token = range tokens {
		if !isWhitespace(token.term()) {
			break
		}
		valSpec = append(valSpec, token)
	}
	token = tokens[i]
	valSpec = append(valSpec, token)
	tokens = tokens[i+1:]
	closed := false
	if itemDoublequote == token.term() { // A double-quoted attribute value consists of ", zero or more characters not including ", and a final ".
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if itemDoublequote == token.term() {
				closed = true
				break
			}
		}
	} else if itemSinglequote == token.term() { // A single-quoted attribute value consists of ', zero or more characters not including ', and a final '.
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if itemSinglequote == token.term() {
				closed = true
				break
			}
		}
	} else { // An unquoted attribute value is a nonempty string of characters not including whitespace, ", ', =, <, >, or `.
		for i, token = range tokens {
			if itemGreater == token.term() {
				i-- // 大于字符 > 不计入 valSpec
				break
			}
			valSpec = append(valSpec, token)
			if isWhitespace(token.term()) {
				// 属性使用空白分隔
				break
			}
			if itemDoublequote == token.term() || itemSinglequote == token.term() || itemEqual == token.term() || itemLess == token.term() || itemGreater == token.term() || itemBacktick == token.term() {
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
	if !isASCIILetter(tokens[0].term()) && itemUnderscore != tokens[0].term() && itemColon != tokens[0].term() {
		return
	}
	attrName = append(attrName, tokens[0])
	tokens = tokens[1:]
	var i int
	var token item
	for i, token = range tokens {
		if !isASCIILetterNumHyphen(token.term()) && itemUnderscore != token.term() && itemDot != token.term() && itemColon != token.term() {
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
	if !isASCIILetter(token.term()) {
		return tokens, nil
	}
	tagName = append(tagName, token)
	for i = 1; i < len(tokens); i++ {
		token = tokens[i]
		if !isASCIILetterNumHyphen(token.term()) {
			break
		}
		tagName = append(tagName, token)
	}
	remains = tokens[i:]

	return
}
