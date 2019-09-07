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

import (
	"bytes"
	"errors"
	"strconv"
)

// RendererFunc 描述了渲染器函数签名。
type RendererFunc func(n *Node, entering bool) (WalkStatus, error)

// Renderer 描述了渲染器结构。
type Renderer struct {
	writer        bytes.Buffer         // 输出缓冲
	lastOut       byte                 // 最新输出的一个字节
	rendererFuncs map[int]RendererFunc // 渲染器
	disableTags   int                  // 标签嵌套计数器，用于判断不可能出现标签嵌套的情况，比如语法树允许图片节点包含链接节点，但是 HTML <img> 不能包含 <a>。
	option        *options             // 解析渲染选项

	listLevel int // 列表级别，用于记录嵌套列表深度
}

// render 从指定的根节点 root 开始遍历并渲染。
func (r *Renderer) render(root *Node) error {
	r.lastOut = itemNewline
	r.writer.Grow(4096)

	return Walk(root, func(n *Node, entering bool) (WalkStatus, error) {
		f := r.rendererFuncs[n.typ]
		if nil == f {
			return WalkStop, errors.New("not found render function for node [type=" + strconv.Itoa(n.typ) + ", tokens=" + string(n.tokens) + "]")
		}

		return f(n, entering)
	})
}

// writeByte 输出一个字节 c。
func (r *Renderer) writeByte(c byte) {
	r.writer.WriteByte(c)
	r.lastOut = c
}

// write 输出指定的 tokens 数组 content。
func (r *Renderer) write(content items) {
	if length := len(content); 0 < length {
		r.writer.Write(content)
		r.lastOut = content[length-1]
	}
}

// writeString 输出指定的字符串 content。
func (r *Renderer) writeString(content string) {
	if length := len(content); 0 < length {
		r.writer.WriteString(content)
		r.lastOut = content[length-1]
	}
}

// newline 会在最新内容不是换行符 \n 时输出一个换行符。
func (r *Renderer) newline() {
	if itemNewline != r.lastOut {
		r.writer.WriteByte(itemNewline)
		r.lastOut = itemNewline
	}
}
