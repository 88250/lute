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

// mergeText 合并 node 中所有（包括子节点）连续的文本节点。
// 合并后顺便进行中文排版优化以及 GFM 自动邮件链接识别。
func (t *Tree) mergeText(node *Node) {
	for child := node.firstChild; nil != child; {
		next := child.next
		if NodeText == child.typ {
			// 逐个合并后续兄弟节点
			for nil != next && NodeText == next.typ {
				child.AppendTokens(next.tokens)
				next.Unlink()
				next = child.next
			}
		} else {
			t.mergeText(child) // 递归处理子节点
		}
		child = next
	}
}
