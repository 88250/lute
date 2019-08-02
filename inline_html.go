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

type InlineHTML struct {
	*BaseNode
}

func (t *Tree) parseInlineHTML(tokens items) (ret Node) {
	startPos := t.context.pos
	ret = &Text{&BaseNode{typ: NodeText, rawText: "<", value: "<"}}

	var tags items
	tags = append(tags, tokens[startPos])
	if itemSlash == tokens[startPos+1].typ { // a closing tag
		tags = append(tags, tokens[startPos+1])
		remains, tagName := t.parseTagName(tokens[t.context.pos+2:])
		if 1 > len(tagName) {
			t.context.pos++
			return
		}
		tags = append(tags, tagName...)
		tokens = remains
	} else { // an open tag
		remains, tagName := t.parseTagName(tokens[t.context.pos+1:])
		if 1 > len(tagName) {
			t.context.pos++
			return
		}

		tags = append(tags, tagName...)
		tokens = remains
		for {
			validAttr, remains, attr := t.parseTagAttr(tokens)
			if !validAttr {
				t.context.pos++
				return
			}

			tokens = remains
			tags = append(tags, attr...)
			if 1 > len(attr) {
				break
			}
		}
	}

	length := len(tokens)
	if 1 > length {
		t.context.pos = startPos + 1
		return
	}

	whitespaces, tokens := tokens.trimLeft()
	if (itemGreater == tokens[0].typ) ||
		(1 < length && itemSlash == tokens[0].typ && itemGreater == tokens[1].typ) {
		tags = append(tags, whitespaces...)
		tags = append(tags, tokens[0])
		if itemSlash == tokens[0].typ {
			tags = append(tags, tokens[1])
		}
		t.context.pos += len(tags)
		ret = &InlineHTML{&BaseNode{typ: NodeInlineHTML, tokens: tags, value: tags.rawText()}}
		return
	}

	t.context.pos = startPos + 1
	return
}

func (t *Tree) parseTagAttr(tokens items) (validAttr bool, remains, attr items) {
	validAttr = true
	remains = tokens
	var whitespaces items
	var i int
	var token *item
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
	validValSpec, tokens, valSpec := t.parseAttrValSpec(tokens)
	if !validValSpec {
		validAttr = false
		return
	}

	remains = tokens
	attr = append(attr, whitespaces...)
	attr = append(attr, attrName...)
	attr = append(attr, valSpec...)

	return
}

func (t *Tree) parseAttrValSpec(tokens items) (validValSpec bool, remains, valSpec items) {
	validValSpec = true
	remains = tokens
	var i int
	var token *item
	for i, token = range tokens {
		if !token.isWhitespace() {
			break
		}
		valSpec = append(valSpec, token)
	}
	token = tokens[i]
	if itemEqual != token.typ {
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
	if itemDoublequote == token.typ { // A double-quoted attribute value consists of ", zero or more characters not including ", and a final ".
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if itemDoublequote == token.typ {
				closed = true
				break
			}
		}
	} else if itemSinglequote == token.typ { // A single-quoted attribute value consists of ', zero or more characters not including ', and a final '.
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if itemSinglequote == token.typ {
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
			if itemDoublequote == token.typ || itemSinglequote == token.typ || itemEqual == token.typ || itemLess == token.typ || itemGreater == token.typ || itemBacktick == token.typ {
				closed = false
				break
			}
			closed = true
		}

	}

	if !closed {
		validValSpec = false
		valSpec = nil
		return
	}

	remains = tokens[i+1:]

	return
}

func (t *Tree) parseAttrName(tokens items) (remains, attrName items) {
	remains = tokens
	if !tokens[0].isASCIILetter() && itemUnderscore != tokens[0].typ && itemColon != tokens[0].typ {
		return
	}
	attrName = append(attrName, tokens[0])
	tokens = tokens[1:]
	var i int
	var token *item
	for i, token = range tokens {
		if !token.isASCIILetterNumHyphen() && itemUnderscore != token.typ && itemDot != token.typ && itemColon != token.typ {
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
	c := tokens[0].val[0]
	if !('A' <= c && 'Z' >= c) && !('a' <= c && 'z' >= c) {
		return tokens, nil
	}
	for i := 0; i < len(tokens[0].val); i++ {
		c = tokens[0].val[i]
		if !('A' <= c && 'Z' >= c) && !('a' <= c && 'z' >= c) &&
			!('0' <= c && '9' >= c) && '-' != c {
			return tokens, nil
		}
	}
	tagName = append(tagName, tokens[0])

	var token *item
	i := 1
	length := len(tokens)
	for ; i < length; i++ {
		token = tokens[i]
		if !token.isASCIILetterNumHyphen() {
			break
		}
		tagName = append(tagName, token)
	}

	remains = tokens[i:]

	return
}
