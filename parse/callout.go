// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"bytes"
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
)

// CalloutStart 判断提示块（> [!Note]）是否开始。
func CalloutStart(t *Tree, container *ast.Node) int {
	if !t.Context.ParseOption.Callout {
		return 0
	}

	if t.Context.indented {
		return 0
	}

	marker := lex.Peek(t.Context.currentLine, t.Context.nextNonspace)
	if lex.ItemGreater != marker {
		return 0
	}

	t.Context.advanceNextNonspace()
	t.Context.advanceOffset(1, false)

	content := string(bytes.TrimSpace(t.Context.currentLine[t.Context.offset:]))
	if !strings.HasPrefix(content, "[!") {
		return 0
	}

	idx := strings.Index(content, "]")
	if 0 > idx {
		return 0
	}

	t.Context.closeUnmatchedBlocks()
	t.Context.addChild(ast.NodeCallout)
	return 1
}

func CalloutContinue(callout *ast.Node, context *Context) int {
	ln := context.currentLine
	if !context.indented && lex.Peek(ln, context.nextNonspace) == lex.ItemGreater {
		context.advanceNextNonspace()
		context.advanceOffset(1, false)
		if token := lex.Peek(ln, context.offset); lex.ItemSpace == token || lex.ItemTab == token {
			context.advanceOffset(1, true)
		}
		return 0
	}
	return 1
}

func (context *Context) calloutFinalize(callout *ast.Node) {
	p := callout.FirstChild
	lines := bytes.Split(p.Tokens, []byte("\n"))
	content := bytes.TrimSpace(lines[0])
	typ := string(bytes.TrimSpace(content[2:bytes.IndexByte(content, ']')]))
	title := string(bytes.TrimSpace(content[bytes.IndexByte(content, ']')+1:]))
	callout.CalloutType = typ
	icon := strings.Split(title, " ")[0]
	if "" != icon {
		if "" != EmojiUnicodeAlias[icon] {
			callout.CalloutIcon = icon
			title = strings.TrimSpace(title[len(icon):])
		}
	}

	switch typ {
	case ast.CalloutTypeNote:
		callout.CalloutIcon = "✏️"
	case ast.CalloutTypeTip:
		callout.CalloutIcon = "💡"
	case ast.CalloutTypeImportant:
		callout.CalloutIcon = "❗"
	case ast.CalloutTypeWarning:
		callout.CalloutIcon = "⚠️"
	case ast.CalloutTypeCaution:
		callout.CalloutIcon = "🚨"
	}

	callout.CalloutTitle = title
	p.Tokens = bytes.Join(lines[1:], []byte("\n"))
	if 1 > len(p.Tokens) {
		p.Tokens = nil
	}
}
