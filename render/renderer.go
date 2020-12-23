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
	"github.com/88250/lute/html"
	"strconv"
	"strings"
	"unicode"
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
	Render() (output []byte)
}

// BaseRenderer 描述了渲染器结构。
type BaseRenderer struct {
	Option              *parse.ParseOptions              // 解析渲染选项
	RendererFuncs       map[ast.NodeType]RendererFunc    // 渲染器
	DefaultRendererFunc RendererFunc                     // 默认渲染器，在 RendererFuncs 中找不到节点渲染器时会使用该默认渲染器进行渲染
	ExtRendererFuncs    map[ast.NodeType]ExtRendererFunc // 用户自定义的渲染器
	Writer              *bytes.Buffer                    // 输出缓冲
	LastOut             byte                             // 最新输出的一个字节
	Tree                *parse.Tree                      // 待渲染的树
	DisableTags         int                              // 标签嵌套计数器，用于判断不可能出现标签嵌套的情况，比如语法树允许图片节点包含链接节点，但是 HTML <img> 不能包含 <a>
	FootnotesDefs       []*ast.Node                      // 脚注定义集
	RenderingFootnotes  bool                             // 是否正在渲染脚注定义
}

// NewBaseRenderer 构造一个 BaseRenderer。
func NewBaseRenderer(tree *parse.Tree) *BaseRenderer {
	ret := &BaseRenderer{RendererFuncs: map[ast.NodeType]RendererFunc{}, ExtRendererFuncs: map[ast.NodeType]ExtRendererFunc{}, Option: tree.Context.ParseOption, Tree: tree}
	ret.Writer = &bytes.Buffer{}
	ret.Writer.Grow(4096)
	return ret
}

// Render 从根节点开始遍历并渲染。
func (r *BaseRenderer) Render() (output []byte) {
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
			if nil != r.DefaultRendererFunc {
				return r.DefaultRendererFunc(n, entering)
			}
			return r.renderDefault(n, entering)
		}
		return render(n, entering)
	})

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
	if !r.Option.AutoSpace {
		return
	}

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

func (r *BaseRenderer) TextAutoSpaceNext(node *ast.Node) {
	if !r.Option.AutoSpace {
		return
	}

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

func (r *BaseRenderer) LinkTextAutoSpacePrevious(node *ast.Node) {
	if !r.Option.AutoSpace {
		return
	}

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

func (r *BaseRenderer) LinkTextAutoSpaceNext(node *ast.Node) {
	if !r.Option.AutoSpace {
		return
	}

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

func SubStr(str string, length int) (ret string) {
	var count int
	for i := 0; i < len(str); {
		r, size := utf8.DecodeRuneInString(str[i:])
		i += size
		ret += string(r)
		count++
		if length <= count {
			break
		}
	}
	return
}

func HeadingID(heading *ast.Node) (ret string) {
	if 0 == len(util.StrToBytes(heading.HeadingNormalizedID)) {
		headingID0(heading)
	}
	return heading.HeadingNormalizedID
}

func headingID0(heading *ast.Node) {
	var root *ast.Node
	for root = heading.Parent; ast.NodeDocument != root.Type; root = root.Parent {
	}

	idOccurs := map[string]int{}
	ast.Walk(root, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			if ast.NodeHeading == n.Type {
				id := normalizeHeadingID(n)
				for ; 0 < idOccurs[id]; id += "-" {
				}
				n.HeadingNormalizedID = id
				idOccurs[id] = 1
			}
		}
		return ast.WalkContinue
	})
}

func normalizeHeadingID(heading *ast.Node) (ret string) {
	headingID := heading.ChildByType(ast.NodeHeadingID)
	var id string
	if nil != headingID {
		id = util.BytesToStr(headingID.Tokens)
	}
	if "" == id {
		id = heading.Text()
	}

	id = strings.TrimLeft(id, "#")
	id = strings.ReplaceAll(id, util.Caret, "")
	for _, r := range id {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			ret += string(r)
		} else {
			ret += "-"
		}
	}
	return
}

type Heading struct {
	URL      string     `json:"url"`
	Path     string     `json:"path"`
	ID       string     `json:"id"`
	Content  string     `json:"content"`
	Level    int        `json:"level"`
	Children []*Heading `json:"children"`
	parent   *Heading
}

func (r *BaseRenderer) renderToC(node *ast.Node, entering bool) ast.WalkStatus {
	if entering {
		headings := r.headings()
		length := len(headings)
		r.WriteString("<div class=\"vditor-toc\" data-block=\"0\" data-type=\"toc-block\" contenteditable=\"false\">")
		if 0 < length {
			r.WriteString("<ul>")
			for _, child := range headings {
				r.renderToC0(child)
			}
			r.WriteString("</ul>")
		} else {
			r.WriteString("[toc]<br>")
		}
		r.WriteString("</div>")
	}
	return ast.WalkContinue
}

func (r *BaseRenderer) renderToC0(heading *Heading) {
	r.WriteString("<li>")
	r.Tag("span", [][]string{{"data-target-id", heading.ID}}, false)
	r.WriteString(heading.Content)
	r.Tag("/span", nil, false)
	if 0 < len(heading.Children) {
		r.WriteString("<ul>")
		for _, child := range heading.Children {
			r.renderToC0(child)
		}
		r.WriteString("</ul>")
	}
	r.WriteString("</li>")
}

func (r *BaseRenderer) Tag(name string, attrs [][]string, selfclosing bool) {
	if r.DisableTags > 0 {
		return
	}

	r.WriteString("<")
	r.WriteString(name)
	if 0 < len(attrs) {
		for _, attr := range attrs {
			r.WriteString(" " + attr[0] + "=\"" + attr[1] + "\"")
		}
	}
	if selfclosing {
		r.WriteString(" /")
	}
	r.WriteString(">")
}

func (r *BaseRenderer) headings() (ret []*Heading) {
	headings := r.Tree.Root.ChildrenByType(ast.NodeHeading)
	var tip *Heading
	for _, heading := range headings {
		if r.Tree.Root != heading.Parent {
			continue
		}

		id := HeadingID(heading)
		if r.Option.VditorWYSIWYG {
			id = "wysiwyg-" + id
		} else if r.Option.VditorIR {
			id = "ir-" + id
		}

		if r.Option.KramdownIAL {
			for _, kv := range heading.KramdownIAL {
				if "id" == kv[0] {
					id = kv[1]
					break
				}
			}
		}

		h := &Heading{
			URL:     r.Tree.URL,
			Path:    r.Tree.Path,
			ID:      id,
			Content: headingText(heading),
			Level:   heading.HeadingLevel,
		}

		if nil == tip {
			ret = append(ret, h)
		} else {
			if tip.Level < h.Level {
				tip.Children = append(tip.Children, h)
				h.parent = tip
			} else {
				if parent := parentTip(h, tip); nil == parent {
					ret = append(ret, h)
				} else {
					parent.Children = append(parent.Children, h)
					h.parent = tip.parent
				}
			}
		}
		tip = h
	}
	return
}

func parentTip(currentHeading, tip *Heading) *Heading {
	if nil == tip.parent {
		return nil
	}

	for parent := tip.parent; nil != parent; parent = parent.parent {
		if parent.Level < currentHeading.Level {
			return parent
		}
	}
	return nil
}

func headingText(n *ast.Node) (ret string) {
	buf := &bytes.Buffer{}
	ast.Walk(n, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		switch n.Type {
		case ast.NodeLinkText, ast.NodeBlockRefText, ast.NodeBlockEmbedText:
			buf.Write(n.Tokens)
		case ast.NodeInlineMathContent:
			buf.WriteString("<span class=\"language-math\">")
			buf.Write(n.Tokens)
			buf.WriteString("</span>")
		case ast.NodeCodeSpanContent:
			buf.WriteString("<code>")
			buf.Write(n.Tokens)
			buf.WriteString("</code>")
		case ast.NodeText:
			if n.ParentIs(ast.NodeStrong) {
				buf.WriteString("<strong>")
				buf.Write(n.Tokens)
				buf.WriteString("</strong>")
			} else if n.ParentIs(ast.NodeEmphasis) {
				buf.WriteString("<em>")
				buf.Write(n.Tokens)
				buf.WriteString("</em>")
			} else {
				if nil != n.Previous && ast.NodeInlineHTML == n.Previous.Type {
					if bytes.HasPrefix(n.Previous.Tokens, []byte("<font ")) {
						buf.Write(n.Previous.Tokens)
						buf.Write(n.Tokens)
					}
					if nil != n.Next && bytes.Equal(n.Next.Tokens, []byte("</font>")) {
						buf.Write(n.Next.Tokens)
					}
				} else {
					buf.Write(n.Tokens)
				}
			}
		}
		return ast.WalkContinue
	})
	return buf.String()
}

func (r *BaseRenderer) setextHeadingLen(node *ast.Node) (ret int) {
	buf := &bytes.Buffer{}
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		if (ast.NodeText == n.Type || ast.NodeLinkText == n.Type || ast.NodeSoftBreak == n.Type) && entering {
			buf.Write(n.Tokens)
		}
		return ast.WalkContinue
	})
	content := buf.String()
	content = strings.ReplaceAll(content, util.Caret, "")
	lines := strings.Split(content, "\n")
	lastLine := lines[len(lines)-1]
	for _, r := range lastLine {
		if utf8.RuneSelf <= r {
			ret += 2
		} else {
			ret++
		}
	}
	if 0 == ret {
		ret = 3
	}
	return
}

func (r *BaseRenderer) renderListStyle(node *ast.Node, attrs *[][]string) {
	if r.Option.RenderListStyle {
		switch node.ListData.Typ {
		case 0:
			*attrs = append(*attrs, []string{"data-style", string(node.Marker)})
		case 1:
			*attrs = append(*attrs, []string{"data-style", strconv.Itoa(node.Num) + string(node.ListData.Delimiter)})
		case 3:
			if 0 == node.ListData.BulletChar {
				*attrs = append(*attrs, []string{"data-style", strconv.Itoa(node.Num) + string(node.ListData.Delimiter)})
			} else {
				*attrs = append(*attrs, []string{"data-style", string(node.Marker)})
			}
		}
	}
}

func (r *BaseRenderer) isLastNode(treeRoot, node *ast.Node) bool {
	if treeRoot == node {
		return true
	}
	if nil != node.Next {
		return false
	}
	if ast.NodeDocument == node.Parent.Type {
		return treeRoot.LastChild == node
	}

	var n *ast.Node
	for n = node.Parent; ; n = n.Parent {
		if ast.NodeDocument == n.Parent.Type {
			break
		}
	}
	return treeRoot.LastChild == n
}

func (r *BaseRenderer) NodeID(node *ast.Node) (ret string) {
	for _, kv := range node.KramdownIAL {
		if "id" == kv[0] {
			return kv[1]
		}
	}
	if ast.NodeListItem == node.Type { // 列表项暂时不生成 ID，等确定是否需要列表项块类型后再打开
		return ""
	}
	return ast.NewNodeID()
}

func (r *BaseRenderer) NodeAttrs(node *ast.Node) (ret [][]string) {
	for _, kv := range node.KramdownIAL {
		if "id" == kv[0] {
			continue
		}
		ret = append(ret, kv)
	}
	return
}

func (r *BaseRenderer) NodeAttrsStr(node *ast.Node) (ret string) {
	for _, kv := range node.KramdownIAL {
		if "id" == kv[0] {
			continue
		}
		ret += kv[0] + "=\"" + kv[1] + "\" "
	}
	if "" != ret {
		ret = ret[:len(ret)-1]
	}
	return
}

func RenderHeadingText(n *ast.Node) (ret string) {
	buf := &bytes.Buffer{}
	ast.Walk(n, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		switch n.Type {
		case ast.NodeLinkText, ast.NodeBlockRefText, ast.NodeBlockEmbedText:
			buf.Write(n.Tokens)
		case ast.NodeInlineMathContent:
			buf.WriteString("<span class=\"language-math\">")
			buf.Write(html.EscapeHTML(n.Tokens))
			buf.WriteString("</span>")
		case ast.NodeCodeSpanContent:
			buf.WriteString("<code>")
			buf.Write(html.EscapeHTML(n.Tokens))
			buf.WriteString("</code>")
		case ast.NodeText:
			if n.ParentIs(ast.NodeStrong) {
				buf.WriteString("<strong>")
				buf.Write(html.EscapeHTML(n.Tokens))
				buf.WriteString("</strong>")
			} else if n.ParentIs(ast.NodeEmphasis) {
				buf.WriteString("<em>")
				buf.Write(html.EscapeHTML(n.Tokens))
				buf.WriteString("</em>")
			} else {
				if nil != n.Previous && ast.NodeInlineHTML == n.Previous.Type {
					if bytes.HasPrefix(n.Previous.Tokens, []byte("<font ")) {
						buf.Write(n.Previous.Tokens)
						buf.Write(html.EscapeHTML(n.Tokens))
					}
					if nil != n.Next && bytes.Equal(n.Next.Tokens, []byte("</font>")) {
						buf.Write(n.Next.Tokens)
					}
				} else {
					buf.Write(html.EscapeHTML(n.Tokens))
				}
			}
		}
		return ast.WalkContinue
	})
	return buf.String()
}
