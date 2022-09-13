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
)

// BlockQueryEmbedStart 判断内容块查询嵌入（{{ SELECT * FROM blocks WHERE content LIKE '%待办%' }}）是否开始。
func BlockQueryEmbedStart(t *Tree, container *ast.Node) int {
	if !t.Context.ParseOption.BlockRef || t.Context.indented {
		return 0
	}

	node := t.parseBlockQueryEmbed()
	if nil == node {
		return 0
	}

	t.Context.closeUnmatchedBlocks()

	for !t.Context.Tip.CanContain(ast.NodeBlockQueryEmbed) {
		t.Context.finalize(t.Context.Tip) // 注意调用 finalize 会向父节点方向进行迭代
	}
	t.Context.Tip.AppendChild(node)
	t.Context.Tip = node
	return 2
}

func (t *Tree) parseBlockQueryEmbed() (ret *ast.Node) {
	tokens := t.Context.currentLine[t.Context.nextNonspace:]
	startCaret := bytes.HasPrefix(tokens, []byte(editor.Caret+"{{"))
	if !bytes.HasPrefix(tokens, []byte("{{")) && !startCaret {
		return
	}
	if startCaret {
		tokens = bytes.Replace(tokens, []byte(editor.Caret+"{{"), []byte("{{"), 1)
	}

	tokens = tokens[2:]
	tokens = bytes.TrimSpace(tokens)
	if t.Context.ParseOption.ProtyleWYSIWYG {
		tokens = bytes.ReplaceAll(tokens, editor.CaretTokens, nil)
	}
	if !bytes.HasSuffix(tokens, []byte("}}")) {
		return
	}

	tokens = tokens[:len(tokens)-2] // 去掉结尾 }}
	script := bytes.TrimSpace(tokens)
	tokens = bytes.TrimSuffix(tokens, editor.CaretTokens)

	ret = &ast.Node{Type: ast.NodeBlockQueryEmbed}
	ret.AppendChild(&ast.Node{Type: ast.NodeOpenBrace})
	ret.AppendChild(&ast.Node{Type: ast.NodeOpenBrace})
	ret.AppendChild(&ast.Node{Type: ast.NodeBlockQueryEmbedScript, Tokens: script})
	ret.AppendChild(&ast.Node{Type: ast.NodeCloseBrace})
	ret.AppendChild(&ast.Node{Type: ast.NodeCloseBrace})
	return
}
