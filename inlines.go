// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

import "sync"

// parseInlines 解析并生成行级节点。
func (t *Tree) parseInlines() {
	t.walkParseInline(t.Root, nil)
}

// walkParseInline 解析生成节点 node 的行级子节点。
func (t *Tree) walkParseInline(node *Node, wg *sync.WaitGroup) {
	defer recoverPanic(nil)
	if nil != wg {
		defer wg.Done()
	}
	if nil == node {
		return
	}

	// 只有如下几种类型的块节点需要生成行级子节点
	if typ := node.typ; NodeParagraph == typ || NodeHeading == typ || NodeTableCell == typ {
		tokens := node.tokens
		if NodeParagraph == typ && nil == tokens {
			// 解析 GFM 表节点后段落内容 tokens 可能会被置换为空，具体可参看函数 Paragraph.Finalize()
			// 在这里从语法树上移除空段落节点
			next := node.next
			node.Unlink()
			// Unlink 会将后一个兄弟节点置空，此处是在在遍历过程中修改树结构，所以需要保持继续迭代后面的兄弟节点
			node.next = next
			return
		}

		length := len(tokens)
		if 1 > length {
			return
		}

		ctx := &InlineContext{tokens: tokens, tokensLen: length}

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

		if t.context.option.GFMAutoLink && !t.context.option.VditorWYSIWYG {
			t.parseGFMAutoEmailLink(node)
			t.parseGFMAutoLink(node)
		}

		if t.context.option.Emoji {
			t.emoji(node)
		}
		return
	} else if NodeCodeBlock == typ {
		if node.isFencedCodeBlock { // 如果是围栏代码块需要细化其子节点
			openMarker := &Node{typ: NodeCodeBlockFenceOpenMarker, tokens: node.codeBlockOpenFence, codeBlockFenceLen: node.codeBlockFenceLen}
			node.PrependChild(openMarker)
			info := &Node{typ: NodeCodeBlockFenceInfoMarker, codeBlockInfo: node.codeBlockInfo}
			node.AppendChild(info)
			code := &Node{typ: NodeCodeBlockCode, tokens: node.tokens}
			node.AppendChild(code)
			node.tokens = nil
			closeMarker := &Node{typ: NodeCodeBlockFenceCloseMarker, tokens: node.codeBlockCloseFence, codeBlockFenceLen: node.codeBlockFenceLen}
			node.AppendChild(closeMarker)
		}
	}

	// 遍历处理子节点

	if t.context.option.ParallelParsing {
		// 通过并行处理提升性能
		cwg := &sync.WaitGroup{}
		for child := node.firstChild; nil != child; child = child.next {
			cwg.Add(1)
			go t.walkParseInline(child, cwg)
		}
		cwg.Wait()
	} else {
		for child := node.firstChild; nil != child; child = child.next {
			t.walkParseInline(child, nil)
		}
	}
}
