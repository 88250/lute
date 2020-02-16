// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package render

import (
	"bytes"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

// RendererFunc 描述了渲染器函数签名。
type RendererFunc func(n *ast.Node, entering bool) ast.WalkStatus

// ExtRendererFunc 描述了用户自定义的渲染器函数签名。
type ExtRendererFunc func(n *ast.Node, entering bool) (string, ast.WalkStatus)

// Renderer 描述了渲染器接口。
type Renderer interface {
	// Render 渲染输出。
	Render() (output []byte, err error)
}

// BaseRenderer 描述了渲染器结构。
type BaseRenderer struct {
	writer              *bytes.Buffer                 // 输出缓冲
	lastOut             byte                          // 最新输出的一个字节
	rendererFuncs       map[ast.NodeType]RendererFunc // 渲染器
	defaultRendererFunc RendererFunc                  // 默认渲染器，在 rendererFuncs 中找不到节点渲染器时会使用该默认渲染器进行渲染
	disableTags         int                           // 标签嵌套计数器，用于判断不可能出现标签嵌套的情况，比如语法树允许图片节点包含链接节点，但是 HTML <img> 不能包含 <a>。
	option              *parse.Options                // 解析渲染选项
	tree                *parse.Tree                   // 待渲染的树

	ExtRendererFuncs map[ast.NodeType]ExtRendererFunc // 用户自定义的渲染器
}

// newBaseRenderer 构造一个 BaseRenderer。
func newBaseRenderer(tree *parse.Tree) *BaseRenderer {
	ret := &BaseRenderer{rendererFuncs: map[ast.NodeType]RendererFunc{}, ExtRendererFuncs: map[ast.NodeType]ExtRendererFunc{}, option: tree.Context.Option, tree: tree}
	ret.writer = &bytes.Buffer{}
	ret.writer.Grow(4096)
	return ret
}

// Render 从指定的根节点 root 开始遍历并渲染。
func (r *BaseRenderer) Render() (output []byte, err error) {
	defer util.RecoverPanic(&err)

	r.lastOut = lex.ItemNewline
	r.writer = &bytes.Buffer{}
	r.writer.Grow(4096)

	ast.Walk(r.tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		extRender := r.ExtRendererFuncs[n.Type]
		if nil != extRender {
			output, status := extRender(n, entering)
			r.writeString(output)
			return status
		}

		render := r.rendererFuncs[n.Type]
		if nil == render {
			render = r.rendererFuncs[n.Type]
			if nil == render {
				if nil != r.defaultRendererFunc {
					return r.defaultRendererFunc(n, entering)
				} else {
					return r.renderDefault(n, entering)
				}
			}
		}
		return render(n, entering)
	})
	if nil != err {
		return
	}

	output = r.writer.Bytes()
	return
}

func (r *BaseRenderer) RendererFuncs(nodeType ast.NodeType) RendererFunc {
	return r.rendererFuncs[nodeType]
}

func (r *BaseRenderer) renderDefault(n *ast.Node, entering bool) ast.WalkStatus {
	r.writeString("not found render function for node [type=" + n.Type.String() + ", Tokens=" + util.BytesToStr(n.Tokens) + "]")
	return ast.WalkContinue
}

// writeByte 输出一个字节 c。
func (r *BaseRenderer) writeByte(c byte) {
	r.writer.WriteByte(c)
	r.lastOut = c
}

// writeBytes 输出字节数组 bytes。
func (r *BaseRenderer) writeBytes(bytes []byte) {
	if length := len(bytes); 0 < length {
		r.writer.Write(bytes)
		r.lastOut = bytes[length-1]
	}
}

// write 输出指定的 Tokens 数组 content。
func (r *BaseRenderer) write(content []byte) {
	if length := len(content); 0 < length {
		r.writer.Write(content)
		r.lastOut = content[length-1]
	}
}

// writeString 输出指定的字符串 content。
func (r *BaseRenderer) writeString(content string) {
	if length := len(content); 0 < length {
		r.writer.WriteString(content)
		r.lastOut = content[length-1]
	}
}

// newline 会在最新内容不是换行符 \n 时输出一个换行符。
func (r *BaseRenderer) newline() {
	if lex.ItemNewline != r.lastOut {
		r.writer.WriteByte(lex.ItemNewline)
		r.lastOut = lex.ItemNewline
	}
}
