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
	ret = &Text{&BaseNode{typ: NodeText, rawText: "<", value: "<"}}
	startPos := t.context.pos

	remains, tagName := t.parseTagName(tokens[startPos+1:])
	if "" == tagName {
		return
	}

	var tags items
	tokens = remains
	var attr items
	for {
		remains, attr = t.parseTagAttr(tokens)
		if 1 > len(attr) {
			break
		}
		tokens = remains
		tags = append(tags, attr...)
	}

	if 1 > len(tags) {
		t.context.pos = startPos
	} else {
		t.context.pos += len(tags)
	}
	
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
	tokens = tokens[i:]
	for i, token = range tokens {
		if !token.isWhitespace() {
			break
		}
		valSpec = append(valSpec, token)
	}
	token = tokens[i]
	valSpec = append(valSpec, token)
	tokens = tokens[i:]
	if itemDoublequote == token.typ { // A double-quoted attribute value consists of ", zero or more characters not including ", and a final ".
		for i, token = range tokens {
			if itemDoublequote == token.typ {
				break
			}
			valSpec = append(valSpec, token)
		}
	} else if itemSinglequote == token.typ { // A single-quoted attribute value consists of ', zero or more characters not including ', and a final '.
		for i, token = range tokens {
			if itemSinglequote == token.typ {
				break
			}
			valSpec = append(valSpec, token)
		}
	} else { // An unquoted attribute value is a nonempty string of characters not including whitespace, ", ', =, <, >, or `.
		for i, token = range tokens {
			if token.isWhitespace() || itemDoublequote == token.typ || itemSinglequote == token.typ || itemEqual == token.typ || itemLess == token.typ || itemGreater == token.typ || itemBacktick == token.typ {
				break
			}
			valSpec = append(valSpec, token)
		}
	}

	remains = tokens[i:]

	return
}

func (t *Tree) parseAttrName(tokens items) (remains, attrName items) {
	remains = tokens
	if !tokens[0].isASCIILetter() && itemUnderscore != tokens[0].typ && itemColon != tokens[0].typ {
		return
	}
	attrName = append(attrName, tokens[0])
	var retRemains items
	retRemains = append(retRemains, tokens[0])
	tokens = tokens[1:]
	var i int
	var token *item
	for i, token = range tokens {
		if !token.isASCIILetterNumHyphen() && itemUnderscore != token.typ && itemDot != token.typ && itemColon != token.typ {
			attrName = nil
			return
		}
		retRemains = append(retRemains, token)
		attrName = append(attrName, token)
	}
	remains = retRemains

	return
}

func (t *Tree) parseTagAttrSpec(tokens items) (remains items, spec string) {

}

func (t *Tree) parseTagName(tokens items) (remains items, tagName string) {
	var name string
	var i int
	var token *item
	for i, token = range tokens {
		if !token.isASCIILetterNumHyphen() {
			return
		}

		name += token.val
		t.context.pos++
	}

	tagName = name
	remains = tokens[i:]
}
