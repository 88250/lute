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

type Paragraph struct {
	*BaseNode

	OpenTag, CloseTag string
}

func (p *Paragraph) CanContain(nodeType NodeType) bool {
	return false
}

func (p *Paragraph) Continue(context *Context) int {
	if context.blank {
		return 1
	}
	return 0
}

func (p *Paragraph) Finalize(context *Context) {
	p.value = strings.TrimSpace(p.value)
	p.tokens = p.tokens.trim()

	// try parsing the beginning as link reference definitions:
	hasReferenceDefs := false
	for tokens := p.tokens; 0 < len(tokens) && itemOpenBracket == tokens[0].typ; tokens = p.tokens {
		if tokens = context.parseLinkRefDef(tokens); nil != tokens {
			p.tokens = tokens
			hasReferenceDefs = true
			continue
		}
		break
	}

	if hasReferenceDefs && p.tokens.isBlankLine() {
		p.Unlink()
	}
}

func (p *Paragraph) AcceptLines() bool {
	return true
}
