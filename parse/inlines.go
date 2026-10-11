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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
)

// parseInlines 解析并生成行级节点。
func (t *Tree) parseInlines() {
	t.walkParseInline(t.Root)

	if t.Context.ParseOption.KramdownSpanIAL {
		t.parseKramdownSpanIAL()
	}
	if t.Context.ParseOption.ProtyleWYSIWYG {
		NormalizeInlineHTMLTextStyles(t)
		// 未转换的标签仍作为字面文本呈现，避免未闭合或不支持的 HTML 进入编辑器。
		ast.Walk(t.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
			if entering && n.Type == ast.NodeInlineHTML {
				n.Type = ast.NodeText
			}
			if !entering && (n.Type == ast.NodeParagraph || n.Type == ast.NodeHeading || n.Type == ast.NodeTableCell) {
				t.mergeText(n)
			}
			return ast.WalkContinue
		})
	}
}

// walkParseInline 解析生成节点 node 的行级子节点。
func (t *Tree) walkParseInline(node *ast.Node) {
	if nil == node {
		return
	}

	typ := node.Type

	if ast.NodeSuperBlock == typ {
		if nil != node.LastChild && ast.NodeSuperBlockLayoutMarker == node.LastChild.Type {
			node.Type = ast.NodeParagraph
			node.Tokens = append([]byte("{{{"), node.LastChild.Tokens...)
			node.FirstChild.Unlink()
			node.LastChild.Unlink()
			typ = ast.NodeParagraph
		}
	}

	// 只有如下几种类型的块节点需要生成行级子节点
	if ast.NodeParagraph == typ || ast.NodeHeading == typ || ast.NodeTableCell == typ {
		tokens := node.Tokens
		if ast.NodeParagraph == typ {
			if nil == tokens && nil == node.FirstChild {
				if ast.NodeTabItem != node.Parent.Type && (ast.NodeListItem != node.Parent.Type ||
					t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV) {
					// 解析 GFM 表节点后段落内容 Tokens 可能会被置换为空，具体可参看函数 Paragraph.Finalize()
					// 在这里从语法树上移除空段落节点
					next := node.Next
					node.Unlink()
					// Unlink 会将后一个兄弟节点置空，此处是在在遍历过程中修改树结构，所以需要保持继续迭代后面的兄弟节点
					node.Next = next
				}
				return
			} else if ial := t.Context.parseKramdownIALInListItem(tokens); 0 < len(ial) {
				if nil != node.Previous {
					// 解析 kramdown 列表或者引述块时可能出现标记符后为空（* \n{id:foo} 或者 >\n{: custom-b="info" }），此时 IAL 应该进行合并
					for _, kv := range ial {
						node.Previous.SetIALAttr(kv[0], html.UnescapeAttrVal(kv[1]))
					}
					next := node.Next
					node.Unlink()
					node.Next = next
					if nil != node.Next && ast.NodeKramdownBlockIAL == node.Next.Type {
						node.Next.Tokens = IAL2Tokens(mergeIALPreservingOrder(Tokens2IAL(node.Next.Tokens), ial))
					}
					return
				}
			}
		}

		length := len(tokens)
		if 1 > length {
			return
		}

		ctx := &InlineContext{tokens: tokens, tokensLen: length}
		if typ == ast.NodeTableCell && t.Context.ParseOption.PreserveTableCellSource {
			if t.Context.TableCellSources == nil {
				t.Context.TableCellSources = map[*ast.Node]string{}
			}
			t.Context.TableCellSources[node] = string(tokens)
		}

		// 生成该块节点的行级子节点
		t.parseInline(node, ctx)

		// 处理该块节点中的强调、加粗和删除线
		t.processEmphasis(nil, ctx)

		// 将连续的文本节点进行合并。
		// 规范只是定义了从输入的 Markdown 文本到输出的 HTML 的解析渲染规则，并未定义中间语法树的规则。
		// 也就是说语法树的节点结构没有标准，可以自行发挥。这里进行文本节点合并主要有两个目的：
		// 1. 减少节点数量，提升后续处理性能
		// 2. 方便后续功能方面的处理，比如 GFM 自动链接解析
		t.mergeText(node)

		editorMode := t.Context.ParseOption.VditorWYSIWYG || t.Context.ParseOption.VditorIR || t.Context.ParseOption.VditorSV || t.Context.ParseOption.ProtyleWYSIWYG
		protyleAutoLink := t.Context.ParseOption.ProtyleWYSIWYG && t.Context.ParseOption.ProtyleWYSIWYGAutoLink
		if (t.Context.ParseOption.GFMAutoLink && !editorMode) || protyleAutoLink {
			t.parseGFMAutoEmailLink(node)
			t.parseGFMAutoLink(node)
		}

		if t.Context.ParseOption.Emoji {
			t.emoji(node)
		}
		return
	} else if ast.NodeCodeBlock == typ {
		if node.IsFencedCodeBlock {
			// 细化围栏代码块子节点
			openMarker := &ast.Node{Type: ast.NodeCodeBlockFenceOpenMarker, Tokens: node.CodeBlockOpenFence, CodeBlockFenceLen: node.CodeBlockFenceLen}
			node.PrependChild(openMarker)
			info := &ast.Node{Type: ast.NodeCodeBlockFenceInfoMarker, CodeBlockInfo: node.CodeBlockInfo}
			node.AppendChild(info)
			code := &ast.Node{Type: ast.NodeCodeBlockCode, Tokens: node.Tokens}
			node.AppendChild(code)
			if nil == node.CodeBlockCloseFence {
				node.CodeBlockCloseFence = node.CodeBlockOpenFence
			}
			closeMarker := &ast.Node{Type: ast.NodeCodeBlockFenceCloseMarker, Tokens: node.CodeBlockCloseFence, CodeBlockFenceLen: node.CodeBlockFenceLen}
			node.AppendChild(closeMarker)
		} else {
			// 细化缩进代码块子节点
			code := &ast.Node{Type: ast.NodeCodeBlockCode, Tokens: node.Tokens}
			node.AppendChild(code)
		}
		node.Tokens = nil
	}

	// 遍历处理子节点
	for child := node.FirstChild; nil != child; child = child.Next {
		t.walkParseInline(child)
	}
}
