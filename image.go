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

// Image 描述了图片节点结构。
type Image struct {
	*BaseNode
	Destination string // 图片链接地址
	Title       string // 图片标题
}

// parseBang 解析 !，可能是图片标记开始 ![ 也可能是普通文本 !。
func (t *Tree) parseBang(tokens items) (ret Node) {
	var startPos = t.context.pos
	t.context.pos++
	if t.context.pos < len(tokens) && itemOpenBracket == tokens[t.context.pos] {
		t.context.pos++
		ret = &Text{tokens: toItems("![")}
		// 将图片开始标记入栈
		t.addBracket(ret, startPos+2, true)
		return
	}

	ret = &Text{tokens: toItems("!")}
	return
}
