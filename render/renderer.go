// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"

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

// Options 描述了渲染选项。
type Options struct {
	// SoftBreak2HardBreak 设置是否将软换行（\n）渲染为硬换行（<br />）。
	SoftBreak2HardBreak bool
	// AutoSpace 设置是否对普通文本中的中西文间自动插入空格。
	// https://github.com/sparanoid/chinese-copywriting-guidelines
	AutoSpace bool
	// RenderListStyle 设置在渲染 OL、UL 时是否添加 data-style 属性 https://github.com/88250/lute/issues/48
	RenderListStyle bool
	// CodeSyntaxHighlight 设置是否对代码块进行语法高亮。
	CodeSyntaxHighlight bool
	// CodeSyntaxHighlightDetectLang bool
	CodeSyntaxHighlightDetectLang bool
	// CodeSyntaxHighlightInlineStyle 设置语法高亮是否为内联样式，默认不内联。
	CodeSyntaxHighlightInlineStyle bool
	// CodeSyntaxHightLineNum 设置语法高亮是否显示行号，默认不显示。
	CodeSyntaxHighlightLineNum bool
	// CodeSyntaxHighlightStyleName 指定语法高亮样式名，默认为 "github"。
	CodeSyntaxHighlightStyleName string
	// Vditor 所见即所得支持。
	VditorWYSIWYG bool
	// Vditor 即时渲染支持。
	VditorIR bool
	// Vditor 分屏预览支持。
	VditorSV bool
	// Protyle 所见即所得支持。
	ProtyleWYSIWYG bool
	// KramdownBlockIAL 设置是否打开 kramdown 块级内联属性列表支持。 https://kramdown.gettalong.org/syntax.html#inline-attribute-lists
	KramdownBlockIAL bool
	// KramdownSpanIAL 设置是否打开 kramdown 行级内联属性列表支持。
	KramdownSpanIAL bool
	// SuperBlock 设置是否支持超级块。 https://github.com/88250/lute/issues/111
	SuperBlock bool
	// ImageLazyLoading 设置图片懒加载时使用的图片路径，配置该字段后将启用图片懒加载。
	// 图片 src 的值会复制给新属性 data-src，然后使用该参数值作为 src 的值 https://github.com/88250/lute/issues/55
	ImageLazyLoading string
	// ChineseParagraphBeginningSpace 设置是否使用传统中文排版“段落开头空两格”。
	ChineseParagraphBeginningSpace bool
	// Sanitize 设置是否启用 XSS 安全过滤 https://github.com/88250/lute/issues/51
	// 注意：Lute 目前的实现存在一些漏洞，请不要依赖它来防御 XSS 攻击。
	Sanitize bool
	// FixTermTypo 设置是否对普通文本中出现的术语进行修正。
	// https://github.com/sparanoid/chinese-copywriting-guidelines
	// 注意：开启术语修正的话会默认在中西文之间插入空格。
	FixTermTypo bool
	// Terms 将传入的 terms 合并覆盖到已有的 Terms 字典。
	Terms map[string]string
	// ToC 设置是否打开“目录”支持。
	ToC bool
	// HeadingID 设置是否打开“自定义标题 ID”支持。
	HeadingID bool
	// KramdownIALIDRenderName 设置 kramdown 内联属性列表中出现 id 属性时渲染 id 属性用的 name(key) 名称，默认为 "id"。
	// 仅在 HTML 渲染器 HtmlRenderer 中支持。
	KramdownIALIDRenderName string
	// HeadingAnchor 设置是否对标题生成链接锚点。
	HeadingAnchor bool
	// GFMTaskListItemClass 作为 GFM 任务列表项类名，默认为 "vditor-task"。
	GFMTaskListItemClass string
	// VditorCodeBlockPreview 设置 Vditor 代码块是否需要渲染预览部分
	VditorCodeBlockPreview bool
	// VditorMathBlockPreview 设置 Vditor 数学公式块是否需要渲染预览部分
	VditorMathBlockPreview bool
	// VditorHTMLBlockPreview 设置 Vditor HTML 块是否需要渲染预览部分
	VditorHTMLBlockPreview bool
	// LinkBase 设置链接、图片、脚注的基础路径。如果用户在链接或者图片地址中使用相对路径（没有协议前缀且不以 / 开头）并且 LinkBase 不为空则会用该值作为前缀。
	// 比如 LinkBase 设置为 http://domain.com/，对于 ![foo](bar.png) 则渲染为 <img src="http://domain.com/bar.png" alt="foo" />
	LinkBase string
	// LinkPrefix 设置连接、图片的路径前缀。一旦设置该值，链接渲染将强制添加该值作为链接前缀，这有别于 LinkBase。
	// 比如 LinkPrefix 设置为 http://domain.com，对于使用绝对路径的 ![foo](/local/path/bar.png) 则渲染为 <img src="http://domain.com/local/path/bar.png" alt="foo" />；
	// 在 LinkBase 和 LinkPrefix 同时设置的情况下，会先处理 LinkBase 逻辑，最后再在 LinkBase 处理结果上加上 LinkPrefix。
	LinkPrefix string
	// NodeIndexStart 用于设置块级节点编号起始值。
	NodeIndexStart int
	// ProtyleContenteditable 设置 Protyle 渲染时标签中的 contenteditable 属性。
	ProtyleContenteditable bool
	// KeepParagraphBeginningSpace 设置是否保留段首空格
	KeepParagraphBeginningSpace bool
	// NetImgMarker 设置 Protyle 是否标记网络图片
	ProtyleMarkNetImg bool
	// Spellcheck 设置是否启用拼写检查
	Spellcheck bool
}

func NewOptions() *Options {
	return &Options{
		SoftBreak2HardBreak:            true,
		AutoSpace:                      false,
		RenderListStyle:                false,
		CodeSyntaxHighlight:            true,
		CodeSyntaxHighlightInlineStyle: false,
		CodeSyntaxHighlightLineNum:     false,
		CodeSyntaxHighlightStyleName:   "github",
		VditorWYSIWYG:                  false,
		VditorIR:                       false,
		VditorSV:                       false,
		ProtyleWYSIWYG:                 false,
		KramdownBlockIAL:               false,
		ChineseParagraphBeginningSpace: false,
		FixTermTypo:                    false,
		ToC:                            false,
		HeadingID:                      false,
		KramdownIALIDRenderName:        "id",
		GFMTaskListItemClass:           "vditor-task",
		VditorCodeBlockPreview:         true,
		VditorMathBlockPreview:         true,
		VditorHTMLBlockPreview:         true,
		LinkBase:                       "",
		LinkPrefix:                     "",
		NodeIndexStart:                 1,
		ProtyleContenteditable:         true,
		ProtyleMarkNetImg:              true,
		Spellcheck:                     false,
		Terms:                          NewTerms(),
	}
}

// BaseRenderer 描述了渲染器结构。
type BaseRenderer struct {
	Options             *Options                         // 渲染选项
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
func NewBaseRenderer(tree *parse.Tree, options *Options) *BaseRenderer {
	ret := &BaseRenderer{RendererFuncs: make(map[ast.NodeType]RendererFunc, 192), ExtRendererFuncs: map[ast.NodeType]ExtRendererFunc{}, Options: options, Tree: tree}
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
	if !r.Options.AutoSpace {
		return
	}

	text := node.ChildByType(ast.NodeText)
	var tokens []byte
	if nil != text {
		tokens = text.Tokens
	}
	if ast.NodeTextMark == node.Type {
		tokens = []byte(node.TextMarkTextContent)
	}
	if 1 > len(tokens) {
		return
	}

	if previous := node.Previous; nil != previous && ast.NodeText == previous.Type {
		prevLast, _ := utf8.DecodeLastRune(previous.Tokens)
		first, _ := utf8.DecodeRune(tokens)
		if allowSpace(prevLast, first) {
			r.Writer.WriteByte(lex.ItemSpace)
		}
	}
}

func (r *BaseRenderer) TextAutoSpaceNext(node *ast.Node) {
	if !r.Options.AutoSpace {
		return
	}

	text := node.ChildByType(ast.NodeText)
	var tokens []byte
	if nil != text {
		tokens = text.Tokens
	}
	if ast.NodeTextMark == node.Type {
		tokens = []byte(node.TextMarkTextContent)
	}
	if 1 > len(tokens) {
		return
	}

	if next := node.Next; nil != next {
		if ast.NodeText == next.Type {
			nextFirst, _ := utf8.DecodeRune(next.Tokens)
			last, _ := utf8.DecodeLastRune(tokens)
			if allowSpace(last, nextFirst) {
				r.Writer.WriteByte(lex.ItemSpace)
			}
		} else if ast.NodeKramdownSpanIAL == next.Type {
			// 优化排版未处理样式文本 https://github.com/siyuan-note/siyuan/issues/6305
			next = next.Next
			if nil != next && ast.NodeText == next.Type {
				nextFirst, _ := utf8.DecodeRune(next.Tokens)
				last, _ := utf8.DecodeLastRune(tokens)
				if allowSpace(last, nextFirst) {
					next.Tokens = append([]byte{lex.ItemSpace}, next.Tokens...)
				}
			}
		}
	}
}

func (r *BaseRenderer) LinkTextAutoSpacePrevious(node *ast.Node) {
	if !r.Options.AutoSpace {
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
	if !r.Options.AutoSpace {
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
	id = strings.ReplaceAll(id, editor.Caret, "")
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
	ID       string     `json:"id"`
	Box      string     `json:"box"`
	Path     string     `json:"path"`
	HPath    string     `json:"hPath"`
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
		if r.Options.VditorWYSIWYG {
			id = "wysiwyg-" + id
		} else if r.Options.VditorIR {
			id = "ir-" + id
		}

		if r.Options.KramdownBlockIAL {
			for _, kv := range heading.KramdownIAL {
				if "id" == kv[0] {
					id = kv[1]
					break
				}
			}
		}

		h := &Heading{
			ID:      id,
			Box:     r.Tree.Box,
			Path:    r.Tree.Path,
			HPath:   r.Tree.HPath,
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
		case ast.NodeLinkText, ast.NodeBlockRefText, ast.NodeBlockRefDynamicText, ast.NodeFileAnnotationRefText:
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
						buf.Write(n.Tokens)
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

func (r *BaseRenderer) setextHeadingLen(node *ast.Node) (ret int) {
	buf := &bytes.Buffer{}
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		if (ast.NodeText == n.Type || ast.NodeLinkText == n.Type || ast.NodeSoftBreak == n.Type) && entering {
			buf.Write(n.Tokens)
		}
		return ast.WalkContinue
	})
	content := buf.String()
	content = strings.ReplaceAll(content, editor.Caret, "")
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
	if r.Options.RenderListStyle {
		switch node.ListData.Typ {
		case 0:
			*attrs = append(*attrs, []string{"data-style", string(node.ListData.Marker)})
		case 1:
			*attrs = append(*attrs, []string{"data-style", strconv.Itoa(node.ListData.Num) + string(node.ListData.Delimiter)})
		case 3:
			if 0 == node.ListData.BulletChar {
				*attrs = append(*attrs, []string{"data-style", strconv.Itoa(node.ListData.Num) + string(node.ListData.Delimiter)})
			} else {
				*attrs = append(*attrs, []string{"data-style", string(node.ListData.Marker)})
			}
		}
	}
}

func (r *BaseRenderer) tagSrc(tokens []byte) []byte {
	if srcIndex := bytes.Index(tokens, []byte("src=\"")); 0 > srcIndex {
		return nil
	} else {
		src := tokens[srcIndex+len("src=\""):]
		src = src[:bytes.Index(src, []byte("\""))]
		return src
	}
}

func (r *BaseRenderer) tagSrcPath(tokens []byte) []byte {
	if srcIndex := bytes.Index(tokens, []byte("src=\"")); 0 < srcIndex {
		src := tokens[srcIndex+len("src=\""):]
		if 1 > len(bytes.ReplaceAll(src, editor.CaretTokens, nil)) {
			return tokens
		}
		targetSrc := r.LinkPath(src)
		originSrc := string(targetSrc)
		if bytes.HasPrefix(targetSrc, []byte("//")) {
			originSrc = "https:" + originSrc
		}
		tokens = bytes.ReplaceAll(tokens, src, []byte(originSrc))
	}
	return tokens
}

func (r *BaseRenderer) isLastNode(treeRoot, node *ast.Node) bool {
	if treeRoot == node || nil == node || nil == node.Parent {
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
		if nil == n || nil == n.Parent {
			return true
		}
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

// languagesNoHighlight 中定义的语言不要进行代码语法高亮。这些代码块会在前端进行渲染，比如各种图表。
var languagesNoHighlight = []string{"mermaid", "echarts", "abc", "graphviz", "mindmap", "flowchart", "plantuml"}

func NoHighlight(language string) bool {
	if "" == language {
		return false
	}

	for _, langNoHighlight := range languagesNoHighlight {
		if language == langNoHighlight {
			return true
		}
	}
	return false
}

func (r *BaseRenderer) Text(node *ast.Node) (ret string) {
	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n.Type {
			case ast.NodeText, ast.NodeLinkText, ast.NodeLinkDest, ast.NodeLinkSpace, ast.NodeLinkTitle, ast.NodeCodeBlockCode,
				ast.NodeCodeSpanContent, ast.NodeInlineMathContent, ast.NodeMathBlockContent, ast.NodeYamlFrontMatterContent,
				ast.NodeHTMLBlock, ast.NodeInlineHTML, ast.NodeEmojiAlias, ast.NodeFileAnnotationRefText, ast.NodeFileAnnotationRefSpace,
				ast.NodeBlockRefText, ast.NodeBlockRefDynamicText, ast.NodeBlockRefSpace,
				ast.NodeKramdownSpanIAL:
				ret += string(n.Tokens)
			case ast.NodeCodeBlockFenceInfoMarker:
				ret += string(n.CodeBlockInfo)
			case ast.NodeLink:
				if 3 == n.LinkType {
					ret += string(n.LinkRefLabel)
				}
			}
		}
		return ast.WalkContinue
	})
	return
}

func (r *BaseRenderer) ParagraphContainImgOnly(paragraph *ast.Node) (ret bool) {
	ret = true
	containImg := false
	ast.Walk(paragraph, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		if ast.NodeText == n.Type {
			if !util.IsEmptyStr(string(n.Tokens)) {
				ret = false
				return ast.WalkStop
			}
		} else if ast.NodeTextMark == n.Type {
			ret = false
			return ast.WalkStop
		} else if ast.NodeImage == n.Type {
			containImg = true
		}
		return ast.WalkContinue
	})

	ret = containImg && ret
	return
}

func RenderHeadingText(n *ast.Node) (ret string) {
	buf := &bytes.Buffer{}
	ast.Walk(n, func(n *ast.Node, entering bool) ast.WalkStatus {
		if !entering {
			return ast.WalkContinue
		}

		switch n.Type {
		case ast.NodeLinkText, ast.NodeBlockRefText, ast.NodeBlockRefDynamicText, ast.NodeFileAnnotationRefText:
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
					if !bytes.HasPrefix(n.Previous.Tokens, []byte("</")) {
						buf.Write(n.Previous.Tokens)
						buf.Write(html.EscapeHTML(n.Tokens))
					} else {
						buf.Write(n.Previous.Tokens)
						buf.Write(html.EscapeHTML(n.Tokens))
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
