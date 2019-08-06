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

// InlineHTML 描述了内联 HTML 节点结构。
type InlineHTML struct {
	*BaseNode
}

func (t *Tree) parseInlineHTML(tokens items) (ret Node) {
	startPos := t.context.pos
	ret = &Text{&BaseNode{typ: NodeText, value: "<"}}

	var tags items
	tags = append(tags, tokens[startPos])
	if itemSlash == tokens[startPos+1] { // a closing tag
		tags = append(tags, tokens[startPos+1])
		remains, tagName := t.parseTagName(tokens[t.context.pos+2:])
		if 1 > len(tagName) {
			t.context.pos++
			return
		}
		tags = append(tags, tagName...)
		tokens = remains
	} else if remains, tagName := t.parseTagName(tokens[t.context.pos+1:]); 0 < len(tagName) {
		tags = append(tags, tagName...)
		tokens = remains
		for {
			valid, remains, attr := t.parseTagAttr(tokens)
			if !valid {
				t.context.pos++
				return
			}

			tokens = remains
			tags = append(tags, attr...)
			if 1 > len(attr) {
				break
			}
		}
	} else if valid, remains, comment := t.parseHTMLComment(tokens[t.context.pos+1:]); valid {
		tags = append(tags, comment...)
		tokens = remains
		t.context.pos += len(tags)
		ret = &InlineHTML{&BaseNode{typ: NodeInlineHTML, tokens: tags, value: tags.rawText()}}
		return
	} else if valid, remains, ins := t.parseProcessingInstruction(tokens[t.context.pos+1:]); valid {
		tags = append(tags, ins...)
		tokens = remains
		t.context.pos += len(tags)
		ret = &InlineHTML{&BaseNode{typ: NodeInlineHTML, tokens: tags, value: tags.rawText()}}
		return
	} else if valid, remains, decl := t.parseDeclaration(tokens[t.context.pos+1:]); valid {
		tags = append(tags, decl...)
		tokens = remains
		t.context.pos += len(tags)
		ret = &InlineHTML{&BaseNode{typ: NodeInlineHTML, tokens: tags, value: tags.rawText()}}
		return
	} else if valid, remains, cdata := t.parseCDATA(tokens[t.context.pos+1:]); valid {
		tags = append(tags, cdata...)
		tokens = remains
		t.context.pos += len(tags)
		ret = &InlineHTML{&BaseNode{typ: NodeInlineHTML, tokens: tags, value: tags.rawText()}}
		return
	} else {
		t.context.pos++
		return
	}

	length := len(tokens)
	if 1 > length {
		t.context.pos = startPos + 1
		return
	}

	whitespaces, tokens := tokens.trimLeft()
	if (itemGreater == tokens[0]) ||
		(1 < length && itemSlash == tokens[0] && itemGreater == tokens[1]) {
		tags = append(tags, whitespaces...)
		tags = append(tags, tokens[0])
		if itemSlash == tokens[0] {
			tags = append(tags, tokens[1])
		}
		t.context.pos += len(tags)
		ret = &InlineHTML{&BaseNode{typ: NodeInlineHTML, tokens: tags, value: tags.rawText()}}
		return
	}

	t.context.pos = startPos + 1
	return
}

func (t *Tree) parseCDATA(tokens items) (valid bool, remains, content items) {
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
	var token item
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

func (t *Tree) parseDeclaration(tokens items) (valid bool, remains, content items) {
	remains = tokens
	if itemBang != tokens[0] {
		return
	}

	var token item
	var i int
	for _, token = range tokens[1:] {
		if token.isWhitespace() {
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

func (t *Tree) parseProcessingInstruction(tokens items) (valid bool, remains, content items) {
	remains = tokens
	if itemQuestion != tokens[0] {
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
		if i <= length-2 && itemQuestion == token && itemGreater == tokens[i+1] {
			break
		}
	}
	tokens = tokens[i:]
	if itemQuestion != tokens[0] || itemGreater != tokens[1] {
		return
	}
	content = append(content, tokens[1])
	valid = true
	remains = tokens[2:]

	return
}

func (t *Tree) parseHTMLComment(tokens items) (valid bool, remains, comment items) {
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
	var token item
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

func (t *Tree) parseTagAttr(tokens items) (valid bool, remains, attr items) {
	valid = true
	remains = tokens
	var whitespaces items
	var i int
	var token item
	for i, token = range tokens {
		if !token.isWhitespace() {
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
		if !token.isWhitespace() {
			break
		}
		valSpec = append(valSpec, token)
	}
	token = tokens[i]
	if itemEqual != token {
		valSpec = nil
		return
	}
	valSpec = append(valSpec, token)
	tokens = tokens[i+1:]
	for i, token = range tokens {
		if !token.isWhitespace() {
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
			valSpec = append(valSpec, token)
			if token.isWhitespace() {
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

func (t *Tree) parseAttrName(tokens items) (remains, attrName items) {
	remains = tokens
	if !tokens[0].isASCIILetter() && itemUnderscore != tokens[0] && itemColon != tokens[0] {
		return
	}
	attrName = append(attrName, tokens[0])
	tokens = tokens[1:]
	var i int
	var token item
	for i, token = range tokens {
		if !token.isASCIILetterNumHyphen() && itemUnderscore != token && itemDot != token && itemColon != token {
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
	if !token.isASCIILetter() {
		return tokens, nil
	}
	tagName = append(tagName, token)
	for i = 1; i < len(tokens); i++ {
		token = tokens[i]
		if !token.isASCIILetterNumHyphen() {
			break
		}
		tagName = append(tagName, token)
	}
	remains = tokens[i:]

	return
}
