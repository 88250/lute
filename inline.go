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

import (
	"html"
	"strings"
)

func (t *Tree) parseInlines() {
	t.context.delimiters = nil
	t.context.brackets = nil
	t.parseBlockInlines(t.Root.Children())
}

func (t *Tree) parseBlockInlines(blocks []Node) {
	for _, block := range blocks {
		cType := block.Type()
		switch cType {
		case NodeCodeBlock, NodeThematicBreak, NodeHTML:
			continue
		}

		cs := block.Children()
		if 0 < len(cs) {
			t.parseBlockInlines(cs)
			continue
		}

		tokens := block.Tokens()
		if nil == tokens {
			continue
		}

		t.context.pos = 0
		for {
			token := tokens[t.context.pos]
			var n Node
			switch token {
			case itemBackslash:
				n = t.parseBackslash(tokens)
			case itemBacktick:
				n = t.parseCodeSpan(tokens)
			case itemAsterisk, itemUnderscore:
				t.handleDelim(block, tokens)
			case itemNewline:
				n = t.parseNewline(block, tokens)
			case itemLess:
				n = t.parseAutolink(tokens)
				if nil == n {
					n = t.parseAutoEmailLink(tokens)
					if nil == n {
						n = t.parseInlineHTML(tokens)
					}
				}
			case itemOpenBracket:
				n = t.parseOpenBracket(tokens)
			case itemCloseBracket:
				n = t.parseCloseBracket(tokens)
			case itemAmpersand:
				n = t.parseEntity(tokens)
			case itemBang:
				n = t.parseBang(tokens)
			default:
				n = t.parseText(tokens)
			}

			if nil != n {
				block.AppendChild(block, n)
			}

			length := len(tokens)
			if 1 > length || t.context.pos >= length || itemEOF == tokens[t.context.pos] {
				break
			}
		}

		t.processEmphasis(nil)
	}
}

func (t *Tree) parseEntity(tokens items) (ret Node) {
	length := len(tokens)
	if 2 > length {
		t.context.pos++
		return &Text{&BaseNode{typ: NodeText, value: "&"}}
	}

	start := t.context.pos
	numeric := itemCrosshatch == tokens[start+1]
	i := t.context.pos
	var token item
	var endWithSemicolon bool
	for ; i < length; i++ {
		token = tokens[i]
		if token.isWhitespace() {
			break
		}
		if itemSemicolon == token {
			i++
			endWithSemicolon = true
			break
		}
	}

	entityName := tokens[start:i].rawText()
	if entityValue, ok := htmlEntities[entityName]; ok { // 通过查表优化
		t.context.pos += i - start
		return &Text{&BaseNode{typ: NodeText, value: entityValue}}
	}

	if !endWithSemicolon {
		t.context.pos++
		return &Text{&BaseNode{typ: NodeText, value: "&"}}
	}

	if numeric {
		entityNameLen := len(entityName)
		if 10 < entityNameLen || 4 > entityNameLen {
			t.context.pos++
			return &Text{&BaseNode{typ: NodeText, value: "&"}}
		}

		hex := 'x' == entityName[2] || 'X' == entityName[2]
		if hex {
			if 5 > entityNameLen {
				t.context.pos++
				return &Text{&BaseNode{typ: NodeText, value: "&"}}
			}
		}
	}

	v := html.UnescapeString(entityName)
	if v == entityName {
		t.context.pos++
		return &Text{&BaseNode{typ: NodeText, value: "&"}}
	}
	t.context.pos += i - start
	return &Text{&BaseNode{typ: NodeText, value: v}}
}

// Try to match close bracket against an opening in the delimiter stack. Add either a link or image, or a plain [ character,
// to block's children. If there is a matching delimiter, remove it from the delimiter stack.
func (t *Tree) parseCloseBracket(tokens items) Node {
	// get last [ or ![
	opener := t.context.brackets
	if nil == opener {
		t.context.pos++
		// no matched opener, just return a literal
		return &Text{&BaseNode{typ: NodeText, value: "]"}}
	}

	if !opener.active {
		t.context.pos++
		// no matched opener, just return a literal
		// take opener off brackets stack
		t.removeBracket()
		return &Text{&BaseNode{typ: NodeText, value: "]"}}
	}

	// If we got here, open is a potential opener
	isImage := opener.image

	var dest, title string
	// Check to see if we have a link/image

	startPos := t.context.pos
	savepos := t.context.pos
	matched := false
	// Inline link?
	if t.context.pos+1 < len(tokens) && itemOpenParen == tokens[t.context.pos+1] {
		t.context.pos++
		isLink := false
		var passed, remains items
		if isLink, passed, remains = tokens[t.context.pos:].spnl(); isLink {
			t.context.pos += len(passed)
			if passed, remains, dest = t.context.parseInlineLinkDest(remains); nil != passed {
				t.context.pos += len(passed)
				if 0 < len(remains) && remains[0].isWhitespace() { // 跟空格的话后续尝试按 title 解析
					t.context.pos++
					if isLink, passed, remains = remains.spnl(); isLink {
						t.context.pos += len(passed)
						validTitle := false
						if validTitle, passed, remains, title = t.context.parseLinkTitle(remains); validTitle {
							t.context.pos += len(passed)
							isLink, passed, remains = remains.spnl()
							t.context.pos += len(passed)
							matched = isLink && itemCloseParen == remains[0]
						}
					}
				} else { // 没有 title
					t.context.pos--
					matched = true
				}
			}
		}
		if !matched {
			t.context.pos = savepos
		}
	}

	var reflabel string
	if !matched {
		// Next, see if there's a link label
		var beforelabel = t.context.pos + 1
		passed, _, label := t.context.parseLinkLabel(tokens[beforelabel:])
		var n = len(passed)
		if n > 0 { // label 解析出来的话说明满足格式 [text][label]
			reflabel = label
			t.context.pos += n + 1
		} else if !opener.bracketAfter {
			// [text][] 或者 [text][] 格式，将第一个 text 视为 label 进行解析
			passed, reflabel = t.extractTokens(tokens, opener.index, startPos)
			if len(passed) > 0 && len(tokens) > beforelabel && itemOpenBracket == tokens[beforelabel] {
				// [text][] 格式，跳过 []
				t.context.pos += 2
			}
		}

		if "" != reflabel {
			// lookup rawlabel in refmap
			var link = t.context.linkRefDef[strings.ToLower(reflabel)]
			if nil != link {
				dest = link.Destination
				title = link.Title
				matched = true
			}
		}
	}

	if matched {
		var node Node
		if isImage {
			node = &Image{&BaseNode{typ: NodeImage}, dest, title}
		} else {
			node = &Link{&BaseNode{typ: NodeLink}, dest, title}
		}

		var tmp, next Node
		tmp = opener.node.Next()
		for nil != tmp {
			next = tmp.Next()
			tmp.Unlink()
			node.AppendChild(node, tmp)
			tmp = next
		}

		t.processEmphasis(opener.previousDelimiter)
		t.removeBracket()
		opener.node.Unlink()

		// We remove this bracket and processEmphasis will remove later delimiters.
		// Now, for a link, we also deactivate earlier link openers.
		// (no links in links)
		if !isImage {
			opener = t.context.brackets
			for nil != opener {
				if !opener.image {
					opener.active = false // deactivate this opener
				}
				opener = opener.previous
			}
		}

		t.context.pos++
		return node
	} else { // no match
		t.removeBracket() // remove this opener from stack
		t.context.pos = startPos
		t.context.pos++
		return &Text{&BaseNode{typ: NodeText, value: "]"}}
	}
}

func (t *Tree) parseOpenBracket(tokens items) (ret Node) {
	t.context.pos++
	ret = &Text{&BaseNode{typ: NodeText, value: "["}}
	// Add entry to stack for this opener
	t.addBracket(ret, t.context.pos, false)

	return
}

func (t *Tree) addBracket(node Node, index int, image bool) {
	if nil != t.context.brackets {
		t.context.brackets.bracketAfter = true
	}

	t.context.brackets = &delimiter{
		node:              node,
		previous:          t.context.brackets,
		previousDelimiter: t.context.delimiters,
		index:             index,
		image:             image,
		active:            true,
	}
}

func (t *Tree) removeBracket() {
	t.context.brackets = t.context.brackets.previous
}

func (t *Tree) parseBackslash(tokens items) (ret Node) {
	if len(tokens)-1 > t.context.pos {
		t.context.pos++
	}
	token := tokens[t.context.pos]
	if itemNewline == token {
		ret = &HardBreak{&BaseNode{typ: NodeHardBreak}}
		t.context.pos++
	} else if token.isASCIIPunct() {
		ret = &Text{&BaseNode{typ: NodeText, value: string(token)}}
		t.context.pos++
	} else {
		ret = &Text{&BaseNode{typ: NodeText, value: "\\"}}
	}

	return
}

func (t *Tree) extractTokens(tokens items, startPos, endPos int) (subTokens items, text string) {
	b := &strings.Builder{}
	for i := startPos; i < endPos; i++ {
		b.WriteString(string(tokens[i]))
	}
	text = b.String()
	subTokens = tokens[startPos:endPos]

	return
}

func (t *Tree) parseText(tokens items) (ret Node) {
	length := len(tokens)
	var token item
	b := &strings.Builder{}
	for ; t.context.pos < length; t.context.pos++ {
		token = tokens[t.context.pos]
		if itemBackslash == token || itemBacktick == token || itemAsterisk == token || itemUnderscore == token ||
			itemNewline == token || itemLess == token || itemOpenBracket == token || itemCloseBracket == token ||
			itemAmpersand == token || itemBang == token {
			break
		}
		b.WriteString(string(token))
	}
	ret = &Text{&BaseNode{typ: NodeText, value: b.String()}}

	return
}

func (t *Tree) parseNewline(block Node, tokens items) (ret Node) {
	t.context.pos++
	// check previous node for trailing spaces
	lastc := block.LastChild()
	hardbreak := false
	if nil != lastc && NodeText == lastc.Type() {
		value := lastc.Value()
		if valueLen := len(value); ' ' == value[valueLen-1] {
			lastc.SetValue(strings.TrimRight(value, " "))
			if 1 < valueLen {
				hardbreak = ' ' == value[len(value)-2]
			}
		}
	}

	if hardbreak {
		ret = &HardBreak{&BaseNode{typ: NodeHardBreak}}
	} else {
		ret = &SoftBreak{&BaseNode{typ: NodeSoftBreak}}
	}

	return
}
