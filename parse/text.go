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
	"github.com/88250/lute/util"
)

func (t *Tree) parseText(ctx *InlineContext) *ast.Node {
	start := ctx.pos
	for ; ctx.pos < ctx.tokensLen; ctx.pos++ {
		if t.isMarker(ctx.tokens[ctx.pos]) {
			// 遇到潜在的标记符时需要跳出该文本节点，回到行级解析主循环
			break
		}
	}
	return &ast.Node{Type: ast.NodeText, Tokens: ctx.tokens[start:ctx.pos]}
}

// isMarker 判断 token 是否是潜在的 Markdown 标记符。
func (t *Tree) isMarker(token byte) bool {
	if lex.IsMarker(token) {
		return true
	}

	if t.Context.ParseOption.Sup && lex.ItemCaret == token {
		return true
	}
	return false
}

var backslash = util.StrToBytes("\\")

func (t *Tree) parseBackslash(block *ast.Node, ctx *InlineContext) *ast.Node {
	if ctx.pos == ctx.tokensLen-1 {
		ctx.pos++
		return &ast.Node{Type: ast.NodeText, Tokens: backslash}
	}

	ctx.pos++
	token := ctx.tokens[ctx.pos]
	if lex.ItemNewline == token {
		ctx.pos++
		return &ast.Node{Type: ast.NodeHardBreak, Tokens: []byte{token}}
	}
	if lex.IsASCIIPunct(token) {
		if '<' == token && nil != t.Context.oldtip && ast.NodeTable == t.Context.oldtip.Type {
			// 表格单元格内存在多行时末尾输入转义符 `\` 导致 `<br />` 暴露 https://github.com/siyuan-note/siyuan/issues/7725
			isBr := ctx.tokens[ctx.pos:]
			if bytes.HasPrefix(isBr, []byte("<br />")) || bytes.HasPrefix(isBr, []byte("<br/>")) || bytes.HasPrefix(isBr, []byte("<br>")) {
				return &ast.Node{Type: ast.NodeText, Tokens: backslash}
			}
		}

		ctx.pos++
		n := &ast.Node{Type: ast.NodeBackslash}
		block.AppendChild(n)
		n.AppendChild(&ast.Node{Type: ast.NodeBackslashContent, Tokens: []byte{token}})
		return nil
	}
	if t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.ProtyleWYSIWYG {
		// 处理 \‸x 情况，插入符后的字符才是待转义的
		tokens := ctx.tokens[ctx.pos:]
		caret := editor.CaretTokens
		if len(caret) < len(tokens) && bytes.HasPrefix(tokens, caret) {
			token = ctx.tokens[ctx.pos+len(caret)]
			if lex.IsASCIIPunct(token) {
				if '<' == token && nil != t.Context.oldtip && ast.NodeTable == t.Context.oldtip.Type {
					// 表格单元格内存在多行时末尾输入转义符 `\` 导致 `<br />` 暴露 https://github.com/siyuan-note/siyuan/issues/7725
					isBr := ctx.tokens[ctx.pos+len(caret):]
					if bytes.HasPrefix(isBr, []byte("<br />")) || bytes.HasPrefix(isBr, []byte("<br/>")) || bytes.HasPrefix(isBr, []byte("<br>")) {
						return &ast.Node{Type: ast.NodeText, Tokens: backslash}
					}
				}

				ctx.pos += len(caret)
				ctx.pos++
				n := &ast.Node{Type: ast.NodeBackslash}
				block.AppendChild(n)
				n.AppendChild(&ast.Node{Type: ast.NodeBackslashContent, Tokens: []byte{token}})
				if t.Context.ParseOption.ProtyleWYSIWYG {
					// Protyle WYSIWYG 模式下插入符移到转义符节点前面
					n.InsertBefore(&ast.Node{Type: ast.NodeText, Tokens: caret})
				} else {
					block.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: caret})
				}
				return nil
			}
		}
	}
	return &ast.Node{Type: ast.NodeText, Tokens: backslash}
}

func (t *Tree) parseNewline(block *ast.Node, ctx *InlineContext) (ret *ast.Node) {
	pos := ctx.pos
	ctx.pos++

	isHardBreak := false
	// 检查前一个节点的结尾空格，如果大于等于两个则说明是硬换行
	if lastc := block.LastChild; nil != lastc && ast.NodeText == lastc.Type {
		tokens := lastc.Tokens
		if valueLen := len(tokens); lex.ItemSpace == tokens[valueLen-1] {
			_, lastc.Tokens = lex.TrimRight(tokens)
			if 1 < valueLen {
				isHardBreak = lex.ItemSpace == tokens[len(tokens)-2]
			}
		}
	}

	ret = &ast.Node{Type: ast.NodeSoftBreak, Tokens: []byte{ctx.tokens[pos]}}
	if isHardBreak {
		ret.Type = ast.NodeHardBreak
	}
	return
}

func (t *Tree) MergeText() {
	t.mergeText(t.Root)
}

// mergeText 合并 node 中所有（包括子节点）连续的文本节点。
// 合并后顺便进行中文排版优化以及 GFM 自动邮件链接识别。
func (t *Tree) mergeText(node *ast.Node) {
	for child := node.FirstChild; nil != child; {
		next := child.Next
		if ast.NodeText == child.Type {
			// 逐个合并后续兄弟节点
			for nil != next && ast.NodeText == next.Type {
				child.AppendTokens(next.Tokens)
				next.Unlink()
				next = child.Next
			}
		} else if ast.NodeLinkText == child.Type {
			for nil != next && ast.NodeLinkText == next.Type {
				child.AppendTokens(next.Tokens)
				next.Unlink()
				next = child.Next
			}
		} else {
			t.mergeText(child) // 递归处理子节点
		}
		child = next
	}
}
