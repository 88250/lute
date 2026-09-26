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
	"unicode"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
)

// linkRefDef 保存定义的文档顺序，用于跨大小写索引选择首个匹配。
type linkRefDef struct {
	order int
	link  *ast.Node
}

// linkRefFoldKey 用简单大小写折叠环中最小的字符建立与 bytes.EqualFold 一致的键。
func linkRefFoldKey(tokens []byte) string {
	var key strings.Builder
	key.Grow(len(tokens))
	for _, r := range string(tokens) {
		if 'a' <= r && r <= 'z' {
			r -= 'a' - 'A'
		} else if unicode.MaxASCII < r {
			minimum := r
			for folded := unicode.SimpleFold(r); folded != r; folded = unicode.SimpleFold(folded) {
				if folded < minimum {
					minimum = folded
				}
			}
			r = minimum
		}
		key.WriteRune(r)
	}
	return key.String()
}

// indexLinkRefDefs 惰性构建链接引用定义索引，避免每次查找时都遍历整棵语法树。
func (t *Tree) indexLinkRefDefs() {
	if t.linkRefDefIndexed {
		return
	}
	t.linkRefDefIndexed = true
	t.linkRefDefs = map[string]linkRefDef{}
	t.linkRefDefsFolded = map[string]linkRefDef{}
	order := 0

	ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering && ast.NodeLinkRefDef == n.Type {
			def := linkRefDef{order: order, link: n.FirstChild}
			key := linkRefFoldKey(n.Tokens)
			if _, exists := t.linkRefDefs[key]; !exists {
				t.linkRefDefs[key] = def
			}
			key = linkRefFoldKey(foldBytes(n.Tokens))
			if _, exists := t.linkRefDefsFolded[key]; !exists {
				t.linkRefDefsFolded[key] = def
			}
			order++
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
	if 0 == len(t.linkRefDefs) {
		return
	}
	key := linkRefFoldKey(label)
	first, found := t.linkRefDefs[key]
	// 分别查询原文、仅标签全折叠、仅定义全折叠，保持三种匹配及文档顺序。
	if def, ok := t.linkRefDefs[linkRefFoldKey(foldBytes(label))]; ok && (!found || def.order < first.order) {
		first, found = def, true
	}
	if def, ok := t.linkRefDefsFolded[key]; ok && (!found || def.order < first.order) {
		first, found = def, true
	}
	if found {
		return first.link
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
