// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"github.com/88250/lute/lex"
	"github.com/88250/lute/util"
)

func SuperBlockContinue(superBlock *ast.Node, context *Context) int {
	if context.isSuperBlockClose(context.currentLine[context.nextNonspace:]) {
		level := 0
		for p := context.Tip; nil != p; p = p.Parent {
			if ast.NodeSuperBlock == p.Type {
				level++
			}
		}
		if 1 < level {
			return 4 // 嵌套层闭合
		}
		return 3 // 顶层闭合
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
	endCaret := bytes.HasSuffix(tokens, util.CaretTokens)
	tokens = bytes.ReplaceAll(tokens, util.CaretTokens, nil)
	if !bytes.Equal([]byte("}}}"), tokens) {
		return
	}
	if endCaret {
		paras := context.Tip.ChildrenByType(ast.NodeParagraph)
		if length := len(paras); 0 < length {
			lastP := paras[length-1]
			lastP.Tokens = append(lastP.Tokens, util.CaretTokens...)
		}
	}
	return true
}
