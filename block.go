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
	"strings"
)

func (t *Tree) parseBlocks() {
	t.context.tip = t.Root
	t.context.oldtip = t.Root
	t.context.linkRefDef = map[string]*Link{}

	for line := t.nextLine(); !line.isEOF(); line = t.nextLine() {
		t.incorporateLine(line)
	}
	for nil != t.context.tip {
		t.context.finalize(t.context.tip)
	}
}

// Analyze a line of text and update the document appropriately.
// We parse markdown text by calling this on each line of input,
// then finalizing the document.
func (t *Tree) incorporateLine(line items) {
	t.context.currentLine = line

	allMatched := true
	var container Node
	container = t.Root
	for lastChild := container.LastChild(); nil != lastChild && lastChild.IsOpen(); lastChild = lastChild.LastChild() {
		container = lastChild

		switch container.Continue(t.context) {
		case 0: // we've matched, keep going
			break
		case 1: // we've failed to match a block
			allMatched = false
			break
		case 2: // we've hit end of line for fenced code close and can return
			return
		}

		if !allMatched {
			container = container.Parent() // back up to last matching block
			break
		}
	}

	t.context.allClosed = container == t.context.oldtip
	t.context.lastMatchedContainer = container

	matchedLeaf := container.Type() != NodeParagraph && container.AcceptLines()
	var startsLen = len(blockStarts)
	// Unless last matched container is a code block, try new container starts,
	// adding children to the last matched container:
	for !matchedLeaf {
		t.context.findNextNonspace()

		// TODO this is a little performance optimization:
		//if !t.context.indented &&
		//	!reMaybeSpecial.test(ln.slice(t.context.nextNonspace)) {
		//	t.context.advanceNextNonspace()
		//	break
		//}

		var i = 0
		for i < startsLen {
			var res = blockStarts[i](t, container)
			if res == 1 {
				container = t.context.tip
				break
			} else if res == 2 {
				container = t.context.tip
				matchedLeaf = true
				break
			} else {
				i++
			}
		}

		if i == startsLen { // nothing matched
			t.context.advanceNextNonspace()
			break
		}
	}

	// What remains at the offset is a text line.  Add the text to the
	// appropriate container.

	// First check for a lazy paragraph continuation:
	if !t.context.allClosed && !t.context.blank && t.context.tip.Type() == NodeParagraph {
		// lazy paragraph continuation
		t.addLine()
	} else { // not a lazy continuation

		// finalize any blocks not matched
		t.context.closeUnmatchedBlocks()
		if t.context.blank && nil != container.LastChild() {
			container.LastChild().SetLastLineBlank(true)
		}

		typ := container.Type()
		isFenced := false
		if NodeCodeBlock == typ {
			isFenced = container.(*CodeBlock).IsFenced
		}

		// Block quote lines are never blank as they start with >
		// and we don't count blanks in fenced code for purposes of tight/loose
		// lists or breaking out of lists.  We also don't set _lastLineBlank
		// on an empty list item, or if we just closed a fenced block.
		var lastLineBlank = t.context.blank &&
			!(typ == NodeBlockquote ||
				(typ == NodeCodeBlock && isFenced) ||
				(typ == NodeListItem &&
					nil != container.FirstChild() /*&&container.sourcepos[0][0] == = this.lineNumber*/))

		// propagate lastLineBlank up through parents:
		var cont = container
		for nil != cont {
			cont.SetLastLineBlank(lastLineBlank)
			cont = cont.Parent()
		}

		if container.AcceptLines() {
			t.addLine()
			// if HtmlBlock, check for end condition
			if typ == NodeHTML {
				html := container.(*HTML)
				if html.Typ >= 1 && html.Typ <= 5 {
					//if reHtmlBlockClose[container._htmlBlockType].test(this.currentLine.slice(this.offset))
					{
						t.context.finalize(container)
					}
				}
			}
		} else if t.context.offset < len(t.context.currentLine) && !t.context.blank {
			// create paragraph container for line
			t.context.addChild(&Paragraph{BaseNode: &BaseNode{typ: NodeParagraph}})
			t.context.advanceNextNonspace()
			t.addLine()
		}
	}
}

type startFunc func(t *Tree, container Node) int

// block start functions.  Return values:
// 0 = no match
// 1 = matched container, keep going
// 2 = matched leaf, no more block starts
var blockStarts = []startFunc{
	// block quote
	func(t *Tree, container Node) int {
		if !t.context.indented {
			token := t.context.currentLine.peek(t.context.nextNonspace)
			if nil != token && itemGreater == token.typ {
				t.context.advanceNextNonspace()
				t.context.advanceOffset(1, false)
				// optional following space
				token = t.context.currentLine.peek(t.context.offset)
				if token.isSpaceOrTab() {
					t.context.advanceOffset(1, true)
				}

				t.context.closeUnmatchedBlocks()
				t.context.addChild(&Blockquote{BaseNode: &BaseNode{typ: NodeBlockquote}})
				return 1
			}
		}

		return 0
	},

	// ATX heading
	func(t *Tree, container Node) int {
		if !t.context.indented {
			if heading := t.parseATXHeading(); nil != heading {
				t.context.advanceNextNonspace()
				t.context.advanceOffset(len(heading.tokens), false)
				t.context.closeUnmatchedBlocks()

				t.context.addChild(heading)
				t.context.advanceOffset(len(t.context.currentLine)-t.context.offset, false)
				return 2
			}
		}

		return 0

	},

	// fenced code block
	//func(t *Tree, container Node) int {
	//	var match;
	//	if (!t.context.indented &&
	//		(match = t.context.currentLine.slice(parser.nextNonspace).match(reCodeFence))) {
	//	var fenceLength = match[0].length;
	//	parser.closeUnmatchedBlocks();
	//	var container = parser.addChild('code_block', parser.nextNonspace);
	//	container._isFenced = true;
	//	container._fenceLength = fenceLength;
	//	container._fenceChar = match[0][0];
	//	container._fenceOffset = parser.indent;
	//	parser.advanceNextNonspace();
	//	parser.advanceOffset(fenceLength, false);
	//	return 2;
	//	} else {
	//	return 0;
	//	}
	//},

	// list item
	func(t *Tree, container Node) int {
		if !t.context.indented || container.Type() == NodeList {
			data := t.parseListMarker(container)
			if nil == data {
				return 0
			}

			t.context.closeUnmatchedBlocks()

			listsMatch := false
			if container.Type() == NodeList {
				listsMatch = t.context.listsMatch(container.(*List).ListData, data)
			}

			// add the list if needed
			if t.context.tip.Type() != NodeList || !listsMatch {
				list := &List{&BaseNode{typ: NodeList}, data}
				t.context.addChild(list)
			}

			// add the list item
			li := &ListItem{&BaseNode{typ: NodeListItem}, data}
			t.context.addChild(li)

			return 1
		}
		return 0
	},

	// indented code block
	func(t *Tree, container Node) int {
		if t.context.indented && t.context.tip.Type() != NodeParagraph && !t.context.blank {
			t.context.advanceOffset(4, true)
			t.context.closeUnmatchedBlocks()
			t.context.addChild(&CodeBlock{BaseNode: &BaseNode{typ: NodeCodeBlock}})
			return 2
		} else {
			return 0
		}
	},
}

// Add a line to the block at the tip.  We assume the tip
// can accept lines -- that check should be done before calling this.
func (t *Tree) addLine() {
	if t.context.partiallyConsumedTab {
		t.context.offset += 1 // skip over tab
		// add space characters:
		var charsToTab = 4 - (t.context.column % 4)
		t.context.tip.AppendValue(strings.Repeat(" ", charsToTab))
	}
	t.context.tip.AddTokens(t.context.currentLine[t.context.offset:])
	t.context.tip.AppendValue(t.context.currentLine[t.context.offset:].rawText())
}

// Returns true if block ends with a blank line, descending if needed
// into lists and sublists.
func endsWithBlankLine(block Node) bool {
	for nil != block {
		if block.LastLineBlank() {
			return true
		}

		var t = block.Type()
		if !block.LastLineChecked() && (t == NodeList || t == NodeListItem) {
			block.SetLastLineBlank(true)
			block = block.LastChild()
		} else {
			block.SetLastLineChecked(true)
			break
		}
	}

	return false
}
