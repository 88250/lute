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

import "strings"

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

		tokens := block.Tokens().trim()
		if nil == tokens {
			continue
		}

		t.context.pos = 0
		for {
			token := tokens[t.context.pos]
			var n Node
			switch token.typ {
			case itemBackslash:
				n = t.parseBackslash(tokens)
			case itemBacktick:
				n = t.parseInlineCode(tokens)
			case itemAsterisk, itemUnderscore:
				t.handleDelim(block, tokens)
			case itemNewline:
				n = t.parseNewline(block, tokens)
			case itemLess:
				n = t.parseInlineHTML(tokens)
			case itemOpenBracket:
				n = t.parseOpenBracket(tokens)
			case itemCloseBracket:
				t.parseCloseBracket(block, tokens)
			default:
				n = t.parseText(tokens)
			}

			if nil != n {
				block.AppendChild(block, n)
			}

			len := len(tokens)
			if 1 > len || t.context.pos >= len || tokens[t.context.pos].isEOF() {
				break
			}
		}

		t.processEmphasis(nil)
	}
}

func (t *Tree) parseCloseBracket(block Node, tokens items) {
	startPos := t.context.pos
	matched := false
	var dest, title, reflabel string

	// get last [ or ![
	opener := t.context.brackets

	if nil == opener {
		t.context.pos++
		return
	}

	if !opener.active {
		// take opener off brackets stack
		t.removeBracket()
		t.context.pos++
		return
	}

	// If we got here, open is a potential opener
	isImage := opener.image

	// Check to see if we have a link/image

	savepos := t.context.pos

	// Inline link?
	if itemOpenParen == tokens[t.context.pos].typ {
		t.context.pos++

		tmp := tokens[t.context.pos:]
		isLink, tmp := tmp.spnl()

		if isLink {
			_, tmp, dest := t.parseLinkDest(tmp)
			if "" != dest {
				isLink, tmp = tmp.spnl()
				if isLink {
					if tmp[0].isWhitespace() { // make sure there's a space before the title
						_, tmp, title := t.parseLinkTitle(tmp)
						if "" != title {
							isLink, tmp = tmp.spnl()
							if isLink && itemCloseParen == tmp[0].typ {
								t.context.pos++
								matched = true
							}
						}
					}
				}
			}
		}
	}

	if !matched {
		t.context.pos = savepos
	}

	if !matched {
		// Next, see if there's a link label
		var beforelabel = t.context.pos
		_, _, label := t.parseLinkLabel(tokens[t.context.pos:])
		var n = len(label)
		if n > 2 {
			reflabel = tokens[beforelabel:beforelabel+n].rawText()
		} else if !opener.bracketAfter {
			// Empty or missing second label means to use the first label as the reference.
			// The reference must not contain a bracket. If we know there's a bracket, we don't even bother checking it.
			_, reflabel = t.extractTokens(tokens, opener.index, startPos)
		}
		if n == 0 {
			// If shortcut reference link, rewind before spaces we skipped.
			t.context.pos = savepos
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

		block.AppendChild(block, node)
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

		return
	} else { // no match
		t.removeBracket() // remove this opener from stack
		t.context.pos = startPos
		block.AppendChild(block, &Text{&BaseNode{typ: NodeText, value: "]"}})

		return
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
	if token.isNewline() {
		ret = &HardBreak{&BaseNode{typ: NodeHardBreak}}
		t.context.pos++
	} else if token.isASCIIPunct() {
		ret = &Text{&BaseNode{typ: NodeText, value: token.val}}
		t.context.pos++
	} else {
		ret = &Text{&BaseNode{typ: NodeText, value: "\\"}}
	}

	return
}

func (t *Tree) extractTokens(tokens items, startPos, endPos int) (subTokens items, text string) {
	for i := startPos; i < endPos; i++ {
		text += tokens[i].val
		subTokens = append(subTokens, tokens[i])
	}

	return
}

func (t *Tree) parseInlineCode(tokens items) (ret Node) {
	startPos := t.context.pos
	marker := tokens[startPos]
	n := tokens[startPos:].accept(marker.typ)
	endPos := t.matchEnd(tokens[startPos+n:], marker, n)
	if 1 > endPos {
		marker.typ = itemStr
		t.context.pos++
		ret = &Text{&BaseNode{typ: NodeText, rawText: marker.val, value: marker.val}}
		return
	}
	endPos = startPos + endPos + n

	var textTokens = items{}
	for i := startPos + n; i < len(tokens) && i < endPos; i++ {
		token := tokens[i]
		if token.isNewline() {
			textTokens = append(textTokens, tSpace)
		} else {
			textTokens = append(textTokens, token)
		}
	}

	if 2 < len(textTokens) && textTokens[0].isSpace() && textTokens[len(textTokens)-1].isSpace() && !textTokens.isBlankLine() {
		textTokens = textTokens[1:]
		textTokens = textTokens[:len(textTokens)-1]
	}

	baseNode := &BaseNode{typ: NodeInlineCode, tokens: textTokens, value: textTokens.rawText()}
	ret = &InlineCode{baseNode}
	t.context.pos = endPos + n

	return
}

func (t *Tree) parseText(tokens items) (ret Node) {
	token := tokens[t.context.pos]
	t.context.pos++
	ret = &Text{&BaseNode{typ: NodeText, rawText: token.val, value: token.val}}

	return
}

func (t *Tree) parseInlineHTML(tokens items) (ret Node) {
	tag := tokens[t.context.pos:]
	tag = tag[:tag.index(itemGreater)+1]
	if 1 > len(tag) {
		token := tokens[t.context.pos]
		ret = &Text{&BaseNode{typ: NodeText, rawText: token.val, value: token.val}}
		t.context.pos++

		return
	}

	codeTokens := items{}
	for _, token := range tag {
		codeTokens = append(codeTokens, token)
	}
	t.context.pos += len(codeTokens)

	baseNode := &BaseNode{typ: NodeInlineHTML, tokens: codeTokens, value: codeTokens.rawText()}
	ret = &InlineHTML{baseNode}

	return
}

func (t *Tree) parseNewline(block Node, tokens items) (ret Node) {
	t.context.pos++
	// check previous node for trailing spaces
	var lastc = block.LastChild()
	if nil != lastc && lastc.Type() == NodeText && lastc.RawText() == " " {
		previous := lastc.Previous()
		hardbreak := nil != previous && previous.RawText() == " "
		if hardbreak {
			ret = &HardBreak{&BaseNode{typ: NodeHardBreak}}
			for nil != lastc && lastc.RawText() == " " {
				tmp := lastc.Previous()
				lastc.Unlink()
				lastc = tmp
			}
		} else {
			ret = &SoftBreak{&BaseNode{typ: NodeSoftBreak}}
		}
	} else {
		ret = &SoftBreak{&BaseNode{typ: NodeSoftBreak}}
	}

	return
}

func (t *Tree) matchEnd(tokens items, openMarker *item, num int) (pos int) {
	for ; pos < len(tokens); pos++ {
		len := tokens[pos:].accept(openMarker.typ)
		if num <= len {
			return pos
		}
	}

	return -1
}
