// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package ast

// WalkStatus 描述了遍历状态。
type WalkStatus int

const (
	// WalkStop 意味着不需要继续遍历。
	WalkStop = iota
	// WalkSkipChildren 意味着不要遍历子节点。
	WalkSkipChildren
	// WalkContinue 意味着继续遍历。
	WalkContinue
)

// Walker 函数定义了遍历节点 n 时需要执行的操作，进入节点设置 entering 为 true，离开节点设置为 false。
// 如果返回 WalkStop 或者 error 则结束遍历。
type Walker func(n *Node, entering bool) WalkStatus

// Walk 使用深度优先算法遍历指定的树节点 n。
func Walk(n *Node, walker Walker) {
	walk(n, walker)
}

func walk(n *Node, walker Walker) (ret WalkStatus) {
	// 进入节点
	ret = walker(n, true)
	if ret == WalkStop {
		return
	}

	if ret != WalkSkipChildren {
		// 递归遍历子节点
		for c := n.FirstChild; nil != c; c = c.Next {
			if ret = walk(c, walker); WalkStop == ret {
				return WalkStop
			}
		}
	}

	// 离开节点
	ret = walker(n, false)
	return
}
