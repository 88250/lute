// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
type Walker func(n Node, entering bool) (WalkStatus, error)

// Walk 使用深度优先算法遍历指定的树节点 n。
func Walk(n Node, walker Walker) (err error) {
	var status WalkStatus

	// 进入节点
	status, err = walker(n, true)
	if err != nil || status == WalkStop {
		return err
	}

	if status != WalkSkipChildren {
		// 递归遍历子节点
		for c := n.FirstChild(); c != nil; c = c.Next() {
			if err := Walk(c, walker); err != nil {
				return err
			}
		}
	}

	// 离开节点
	status, err = walker(n, false)
	if err != nil || status == WalkStop {
		return err
	}
	return
}
