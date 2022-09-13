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

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/lex"
)

// SuperBlockStart 判断超级块（{{{ blocks }}}）是否开始。
func SuperBlockStart(t *Tree, container *ast.Node) int {
	if !t.Context.ParseOption.SuperBlock || t.Context.indented {
		return 0
	}

	if ok, layout := t.parseSuperBlock(); ok {
		t.Context.closeUnmatchedBlocks()
		t.Context.addChild(ast.NodeSuperBlock)
		t.Context.addChildMarker(ast.NodeSuperBlockOpenMarker, nil)
		t.Context.addChildMarker(ast.NodeSuperBlockLayoutMarker, layout)
		t.Context.offset = t.Context.currentLineLen - 1 // 整行过
		return 1
	}
	return 0
}

func SuperBlockContinue(superBlock *ast.Node, context *Context) int {
	if nil != context.Tip.LastChild && ast.NodeSuperBlockCloseMarker == context.Tip.LastChild.Type && context.Tip.LastChild.Close {
		return 1
	}

	if context.isSuperBlockClose(context.currentLine[context.nextNonspace:]) {
		for p := context.Tip; nil != p; p = p.Parent {
			if ast.NodeSuperBlock == p.Type {
				return 3 // 闭合
			}
		}
	}
	return 0
}

func (context *Context) superBlockFinalize(superBlock *ast.Node) {
	// 最终化所有子块
	for child := superBlock.FirstChild; nil != child; child = child.Next {
		if child.Close {
			continue
		}
		context.finalize(child)
	}
}

func (t *Tree) parseSuperBlock() (ok bool, layout []byte) {
	marker := t.Context.currentLine[t.Context.nextNonspace]
	if lex.ItemOpenBrace != marker {
		return
	}

	fenceChar := marker
	var fenceLen int
	for i := t.Context.nextNonspace; i < t.Context.currentLineLen && fenceChar == t.Context.currentLine[i]; i++ {
		fenceLen++
	}

	if 3 != fenceLen {
		return
	}

	layout = t.Context.currentLine[t.Context.nextNonspace+fenceLen:]
	layout = lex.TrimWhitespace(layout)
	if !bytes.EqualFold(layout, nil) && !bytes.EqualFold(layout, []byte("row")) && !bytes.EqualFold(layout, []byte("col")) {
		return
	}
	return true, layout
}

func (context *Context) isSuperBlockClose(tokens []byte) (ok bool) {
	tokens = lex.TrimWhitespace(tokens)
	if bytes.Equal(tokens, []byte(editor.Caret+"}}}")) {
		p := &ast.Node{Type: ast.NodeParagraph, Tokens: editor.CaretTokens}
		context.TipAppendChild(p)
	}
	endCaret := bytes.HasSuffix(tokens, editor.CaretTokens)
	tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
	if !bytes.Equal([]byte("}}}"), tokens) {
		return
	}
	if endCaret {
		paras := context.Tip.ChildrenByType(ast.NodeParagraph)
		if length := len(paras); 0 < length {
			lastP := paras[length-1]
			lastP.Tokens = append(lastP.Tokens, editor.CaretTokens...)
		}
	}
	return true
}
