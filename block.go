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
	for line := t.nextLine(); !line.isEOF(); line = t.nextLine() {
		t.processLine(line)
	}
}

func (t *Tree) processLine(line items) {
	t.context.line = line

	allMatched := true
	var container Node
	container = t.Root
	for lastChild := container.LastChild(); nil != lastChild && lastChild.IsOpen(); container = container.LastChild() {
		container = lastChild

		switch container.Continuation(t.context.line) {
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



	t.appendBlock(openBlock)
}

func (t *Tree) parseBlock(tokens items) (node Node) {
	node = t.parseBlankLine(tokens)
	if nil == node {
		node = t.parseIndentCode(tokens)
	}
	if nil == node {
		node = t.parseATXHeading(tokens)
	}
	if nil == node {
		node = t.parseList(tokens)
	}
	if nil == node {
		node = t.parseBlockquote(tokens)
	}
	if nil == node {
		node = t.parseParagraph(tokens)
	}

	return
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
	t.context.tip = parent
}

// Add block of type tag as a child of the tip.  If the tip can't
// accept children, close and finalize it and try its parent,
// and so on til we find a block that can accept children.
func (t *Tree) addChild(newBlock Node) {
	for !t.context.tip.CanContain(newBlock) {
		t.finalize(t.context.tip)
	}

	t.context.tip.AppendChild(t.context.tip, newBlock)
	t.context.tip = newBlock
}
