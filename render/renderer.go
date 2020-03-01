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
	"unicode/utf8"

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
	Option              *parse.Options                   // 解析渲染选项
	RendererFuncs       map[ast.NodeType]RendererFunc    // 渲染器
	DefaultRendererFunc RendererFunc                     // 默认渲染器，在 RendererFuncs 中找不到节点渲染器时会使用该默认渲染器进行渲染
	ExtRendererFuncs    map[ast.NodeType]ExtRendererFunc // 用户自定义的渲染器
	Writer              *bytes.Buffer                    // 输出缓冲
	LastOut             byte                             // 最新输出的一个字节
	Tree                *parse.Tree                      // 待渲染的树
	DisableTags         int                              // 标签嵌套计数器，用于判断不可能出现标签嵌套的情况，比如语法树允许图片节点包含链接节点，但是 HTML <img> 不能包含 <a>。
}

// NewBaseRenderer 构造一个 BaseRenderer。
func NewBaseRenderer(tree *parse.Tree) *BaseRenderer {
	ret := &BaseRenderer{RendererFuncs: map[ast.NodeType]RendererFunc{}, ExtRendererFuncs: map[ast.NodeType]ExtRendererFunc{}, Option: tree.Context.Option, Tree: tree}
	ret.Writer = &bytes.Buffer{}
	ret.Writer.Grow(4096)
	return ret
}

// Render 从指定的根节点 root 开始遍历并渲染。
func (r *BaseRenderer) Render() (output []byte, err error) {
	defer util.RecoverPanic(&err)

	r.LastOut = lex.ItemNewline
	r.Writer = &bytes.Buffer{}
	r.Writer.Grow(4096)

	ast.Walk(r.Tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		extRender := r.ExtRendererFuncs[n.Type]
		if nil != extRender {
			output, status := extRender(n, entering)
			r.WriteString(output)
			return status
		}

		render := r.RendererFuncs[n.Type]
		if nil == render {
			render = r.RendererFuncs[n.Type]
			if nil == render {
				if nil != r.DefaultRendererFunc {
					return r.DefaultRendererFunc(n, entering)
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

	output = r.Writer.Bytes()
	return
}

func (r *BaseRenderer) renderDefault(n *ast.Node, entering bool) ast.WalkStatus {
	r.WriteString("not found render function for node [type=" + n.Type.String() + ", Tokens=" + util.BytesToStr(n.Tokens) + "]")
	return ast.WalkContinue
}

// WriteByte 输出一个字节 c。
func (r *BaseRenderer) WriteByte(c byte) {
	r.Writer.WriteByte(c)
	r.LastOut = c
}

// Write 输出指定的字节数组 content。
func (r *BaseRenderer) Write(content []byte) {
	if length := len(content); 0 < length {
		r.Writer.Write(content)
		r.LastOut = content[length-1]
	}
}

// WriteString 输出指定的字符串 content。
func (r *BaseRenderer) WriteString(content string) {
	if length := len(content); 0 < length {
		r.Writer.WriteString(content)
		r.LastOut = content[length-1]
	}
}

// Newline 会在最新内容不是换行符 \n 时输出一个换行符。
func (r *BaseRenderer) Newline() {
	if lex.ItemNewline != r.LastOut {
		r.Writer.WriteByte(lex.ItemNewline)
		r.LastOut = lex.ItemNewline
	}
}

func (r *BaseRenderer) TextAutoSpacePrevious(node *ast.Node) {
	if r.Option.AutoSpace {
		if text := node.ChildByType(ast.NodeText); nil != text && nil != text.Tokens {
			if previous := node.Previous; nil != previous && ast.NodeText == previous.Type {
				prevLast, _ := utf8.DecodeLastRune(previous.Tokens)
				first, _ := utf8.DecodeRune(text.Tokens)
				if allowSpace(prevLast, first) {
					r.Writer.WriteByte(lex.ItemSpace)
				}
			}
		}
	}
}

func (r *BaseRenderer) TextAutoSpaceNext(node *ast.Node) {
	if r.Option.AutoSpace {
		if text := node.ChildByType(ast.NodeText); nil != text && nil != text.Tokens {
			if next := node.Next; nil != next && ast.NodeText == next.Type {
				nextFirst, _ := utf8.DecodeRune(next.Tokens)
				last, _ := utf8.DecodeLastRune(text.Tokens)
				if allowSpace(last, nextFirst) {
					r.Writer.WriteByte(lex.ItemSpace)
				}
			}
		}
	}
}

func (r *BaseRenderer) LinkTextAutoSpacePrevious(node *ast.Node) {
	if r.Option.AutoSpace {
		if text := node.ChildByType(ast.NodeLinkText); nil != text && nil != text.Tokens {
			if previous := node.Previous; nil != previous && ast.NodeText == previous.Type {
				prevLast, _ := utf8.DecodeLastRune(previous.Tokens)
				first, _ := utf8.DecodeRune(text.Tokens)
				if allowSpace(prevLast, first) {
					r.Writer.WriteByte(lex.ItemSpace)
				}
			}
		}
	}
}

func (r *BaseRenderer) LinkTextAutoSpaceNext(node *ast.Node) {
	if r.Option.AutoSpace {
		if text := node.ChildByType(ast.NodeLinkText); nil != text && nil != text.Tokens {
			if next := node.Next; nil != next && ast.NodeText == next.Type {
				nextFirst, _ := utf8.DecodeRune(next.Tokens)
				last, _ := utf8.DecodeLastRune(text.Tokens)
				if allowSpace(last, nextFirst) {
					r.Writer.WriteByte(lex.ItemSpace)
				}
			}
		}
	}
}
