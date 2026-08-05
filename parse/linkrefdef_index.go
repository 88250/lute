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

// linkRefDef 描述链接引用定义索引项，folded 为定义 label 的全折叠（full case fold）结果。
type linkRefDef struct {
	tokens []byte
	folded []byte
	link   *ast.Node
}

// indexLinkRefDefs 惰性构建链接引用定义索引，避免每次查找时都遍历整棵语法树。
func (t *Tree) indexLinkRefDefs() {
	if t.linkRefDefIndexed {
		return
	}
	t.linkRefDefIndexed = true

	ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering && ast.NodeLinkRefDef == n.Type {
			t.linkRefDefs = append(t.linkRefDefs, &linkRefDef{tokens: n.Tokens, folded: foldBytes(n.Tokens), link: n.FirstChild})
		}
		return ast.WalkContinue
	})
}

// indexFootnotesDefs 惰性构建脚注定义索引（保持文档顺序），避免每次查找时都遍历整棵语法树。
func (t *Tree) indexFootnotesDefs() {
	if t.footnotesDefsIndex {
		return
	}
	t.footnotesDefsIndex = true

	ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering && ast.NodeFootnotesDef == n.Type {
			t.footnotesDefs = append(t.footnotesDefs, n)
		}
		return ast.WalkContinue
	})
}

func (t *Tree) FindLinkRefDefLink(label []byte) (link *ast.Node) {
	if !t.Context.ParseOption.LinkRef {
		return
	}

	if t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.ProtyleWYSIWYG {
		label = bytes.ReplaceAll(label, editor.CaretTokens, nil)
	}

	t.indexLinkRefDefs()
	// 全折叠结果只计算一次，避免每次查找都重复构建折叠状态
	foldedLabel := foldBytes(label)
	for _, def := range t.linkRefDefs {
		// 按定义文档顺序依次比较，保证第一个匹配的定义优先
		if bytes.EqualFold(def.tokens, label) ||
			bytes.EqualFold(foldedLabel, def.tokens) || bytes.EqualFold(def.folded, label) {
			return def.link
		}
	}
	return
}

func (t *Tree) FindFootnotesDef(label []byte) (pos int, def *ast.Node) {
	pos = 0
	if nil != t.Context && (t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.ProtyleWYSIWYG) {
		label = bytes.ReplaceAll(label, editor.CaretTokens, nil)
	}

	t.indexFootnotesDefs()
	for i, n := range t.footnotesDefs {
		if bytes.EqualFold(n.Tokens, label) {
			pos = i + 1
			def = n
			return
		}
	}
	return
}
