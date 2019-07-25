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

type Paragraph struct {
	*BaseNode

	OpenTag, CloseTag string
}

func (p *Paragraph) Continuation(tokens items) int {
	if tokens.isBlankLine() {
		return 1
	}
	return 0
}

func (p *Paragraph) Finalize() {
	// TODO try parsing the beginning as link reference definitions:
	//var pos int
	//hasReferenceDefs := false
	// for (peek(p.rawText, 0) == itemOpenBracket && (pos =
	//context.inlineParser.parseReference(block._string_content,
	//	parser.refmap))) {
	//block._string_content = block._string_content.slice(pos);
	//hasReferenceDefs = true;
	//}
	//if hasReferenceDefs && isBlank(block._string_content) {
	//	block.unlink()
	//}
}

func (p *Paragraph) AcceptLines() bool {
	return true
}
