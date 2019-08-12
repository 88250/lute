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

import (
	"bytes"
	"errors"
	"fmt"
)

// RendererFunc 描述了渲染器函数签名。
type RendererFunc func(n Node, entering bool) (WalkStatus, error)

// Renderer 描述了渲染器结构。
type Renderer struct {
	writer        bytes.Buffer         // 输出缓冲
	lastOut       byte                 // 最新输出的一个字节
	rendererFuncs map[int]RendererFunc // 渲染器
	disableTags   int                  // 标签嵌套计数器，用于判断不可能出现标签嵌套的情况。比如语法树允许图片节点包含链接节点，但是 HTML <img> 不能包含 <a>。
}

// render 渲染指定的节点 n。
func (r *Renderer) Render(n Node) error {
	r.lastOut = itemNewline
	r.writer.Grow(4096)

	return Walk(n, func(n Node, entering bool) (WalkStatus, error) {
		f := r.rendererFuncs[n.Type()]
		if nil == f {
			return WalkStop, errors.New(fmt.Sprintf("not found render function for node [type=%d, text=%s]", n.Type(), n.RawText()))
		}

		return f(n, entering)
	})
}

// Write 输出指定的 tokens 数组 content。
func (r *Renderer) Write(content items) {
	if length := len(content); 0 < length {
		r.writer.Write(content)
		r.lastOut = content[length-1]
	}
}

// WriteString 输出指定的字符串 content。
func (r *Renderer) WriteString(content string) {
	if length := len(content); 0 < length {
		r.writer.WriteString(content)
		r.lastOut = content[length-1]
	}
}

// Newline 会在最新内容不是换行符 \n 时输出一个换行符。
func (r *Renderer) Newline() {
	if itemNewline != r.lastOut {
		r.writer.WriteByte(itemNewline)
		r.lastOut = itemNewline
	}
}
