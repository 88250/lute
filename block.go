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

func (t *Tree) parseBlocks() {
	t.context.linkRefDef = map[string]*Link{}

	for line := t.nextLine(); !line.isEOF(); line = t.nextLine() {
		t.processLine(line)
	}
}

func (t *Tree) processLine(line items) {
	t.context.currentLine = line

	allMatched := true
	var container Node
	container = t.Root
	for lastChild := container.LastChild(); nil != lastChild && lastChild.IsOpen(); container = container.LastChild() {
		container = lastChild

		switch container.Continuation(t.context.currentLine) {
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
	var starts = blockStarts
	var startsLen = len(starts)
	// Unless last matched container is a code block, try new container starts,
	// adding children to the last matched container:
	for !matchedLeaf {
		t.findNextNonspace()

		// this is a little performance optimization:
		//if !t.context.indented &&
		//	!reMaybeSpecial.test(ln.slice(t.context.nextNonspace)) {
		//	t.context.advanceNextNonspace()
		//	break
		//}

		var i = 0
		for i < startsLen {
			var res = starts[i](t, container)
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
			t.advanceNextNonspace()
			break
		}
	}

}

func (t *Tree) advanceOffset(count int, columns bool) {
	var currentLine = t.context.currentLine
	var charsToTab, charsToAdvance int
	var c *item
	for c = currentLine[t.context.offset]; count > 0 && nil != c; {
		if c.isTab() {
			charsToTab = 4 - (t.context.column % 4)
			if columns {
				t.context.partiallyConsumedTab = charsToTab > count
				if charsToTab > count {
					charsToAdvance = count
				} else {
					charsToAdvance = charsToTab
				}
				t.context.column += charsToAdvance
				if !t.context.partiallyConsumedTab {
					t.context.offset += 1
				}
				count -= charsToAdvance
			} else {
				t.context.partiallyConsumedTab = false
				t.context.column += charsToTab
				t.context.offset += 1
				count -= 1
			}
		} else {
			t.context.partiallyConsumedTab = false
			t.context.offset += 1
			t.context.column += 1 // assume ascii; block starts are ascii
			count -= 1
		}
	}
}

func (t *Tree) advanceNextNonspace() {
	t.context.offset = t.context.nextNonspace
	t.context.column = t.context.nextNonspaceColumn
	t.context.partiallyConsumedTab = false
}

func (t *Tree) findNextNonspace() {
	currentLine := t.context.currentLine
	i := t.context.offset
	cols := t.context.column

	var c *item
	for _, c = range currentLine {
		if "" == c.val {
			break
		}

		if c.isSpace() {
			i++
			cols++
		} else if c.isTab() {
			i++
			cols += 4 - (cols % 4)
		} else {
			break
		}
	}
	t.context.blank = c.val == "\n" || c.val == "\r" || "" == c.val
	t.context.nextNonspace = i
	t.context.nextNonspaceColumn = cols
	t.context.indent = t.context.nextNonspaceColumn - t.context.column
	t.context.indented = t.context.indent >= 4
}

// Finalize and close any unmatched blocks.
func (t *Tree) closeUnmatchedBlocks() {
	if !t.context.allClosed {
		// finalize any blocks not matched
		for t.context.oldtip != t.context.lastMatchedContainer {
			parent := t.context.oldtip.Parent()
			t.finalize(t.context.oldtip)
			t.context.oldtip = parent
		}
		t.context.allClosed = true
	}
}

// Finalize a block.  Close it and do any necessary postprocessing,
// e.g. creating string_content from strings, setting the 'tight'
// or 'loose' status of a list, and parsing the beginnings
// of paragraphs for reference definitions.  Reset the tip to the
// parent of the closed block.
func (t *Tree) finalize(block Node) {
	var parent = block.Parent()
	block.Close()
	block.Finalize()
	t.context.tip = parent
}

// Add block of type tag as a child of the tip.  If the tip can't
// accept children, close and finalize it and try its parent,
// and so on til we find a block that can accept children.
func (t *Tree) addChild(typ NodeType) {
	for !t.context.tip.CanContain(typ) {
		t.finalize(t.context.tip)
	}

	newBlock := &BaseNode{typ: typ}
	t.context.tip.AppendChild(t.context.tip, newBlock)
	t.context.tip = newBlock
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
			token := peek(t.context.currentLine, t.context.nextNonspace)
			if nil != token && itemGreater == token.typ {
				t.advanceNextNonspace()
				t.advanceOffset(1, false)
				// optional following space
				token = peek(t.context.currentLine, t.context.offset)
				if token.isSpaceOrTab() {
					t.advanceOffset(1, true)
				}

				t.closeUnmatchedBlocks()
				t.addChild(NodeBlockquote)
				return 1
			}
		}

		return 0
	},
}

func peek(ln items, pos int) *item {
	if pos < len(ln) {
		return ln[pos]
	}

	return nil
}
