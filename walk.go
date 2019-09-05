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
type Walker func(n *Node, entering bool) (WalkStatus, error)

// Walk 使用深度优先算法遍历指定的树节点 n。
func Walk(n *Node, walker Walker) (err error) {
	var status WalkStatus

	// 进入节点
	status, err = walker(n, true)
	if nil != err || status == WalkStop {
		return
	}

	if status != WalkSkipChildren {
		// 递归遍历子节点
		for c := n.firstChild; nil != c; c = c.next {
			if err := Walk(c, walker); nil != err {
				return err
			}
		}
	}

	// 离开节点
	status, err = walker(n, false)
	return
}
