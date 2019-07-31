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
	if itemSlash == tokens[startPos + 1].typ { // a closing tag
		tags = append(tags, tokens[startPos + 1])
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
		var attr items
		for {
			remains, attr = t.parseTagAttr(tokens)
			tokens = remains
			tags = append(tags, attr...)
			if 1 > len(attr) {
				break
			}
		}
	}

	if itemGreater != tokens[0].typ {
		t.context.pos = startPos + 1
		return
	}

	tags = append(tags, tokens[0])
	t.context.pos += len(tags)
	ret = &InlineHTML{&BaseNode{typ: NodeInlineHTML, tokens: tags, value: tags.rawText()}}

	return
}

func (t *Tree) parseTagAttr(tokens items) (remains, attr items) {
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
	tokens, valSpec = t.parseAttrValSpec(tokens)
	if 1 > len(valSpec) {
		return
	}

	remains = tokens
	attr = append(attr, whitespaces...)
	attr = append(attr, attrName...)
	attr = append(attr, valSpec...)

	return
}

func (t *Tree) parseAttrValSpec(tokens items) (remains, valSpec items) {
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
	if itemDoublequote == token.typ { // A double-quoted attribute value consists of ", zero or more characters not including ", and a final ".
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if itemDoublequote == token.typ {
				break
			}
		}
	} else if itemSinglequote == token.typ { // A single-quoted attribute value consists of ', zero or more characters not including ', and a final '.
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if itemSinglequote == token.typ {
				break
			}
		}
	} else { // An unquoted attribute value is a nonempty string of characters not including whitespace, ", ', =, <, >, or `.
		for i, token = range tokens {
			valSpec = append(valSpec, token)
			if token.isWhitespace() || itemDoublequote == token.typ || itemSinglequote == token.typ || itemEqual == token.typ || itemLess == token.typ || itemGreater == token.typ || itemBacktick == token.typ {
				break
			}
		}
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
	var i int
	var token *item
	for i, token = range tokens {
		if !token.isASCIILetterNumHyphen() {
			break
		}

		tagName = append(tagName, token)
	}

	remains = tokens[i:]

	return
}
