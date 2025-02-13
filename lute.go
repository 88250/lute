// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

// Package lute 提供了一款结构化的 Markdown 引擎，支持 Go 和 JavaScript。
package lute

import (
	"bytes"
	"errors"
	"strings"
	"sync"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/lex"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
	"github.com/gopherjs/gopherjs/js"
)

const Version = "1.7.6"

// Lute 描述了 Lute 引擎的顶层使用入口。
type Lute struct {
	ParseOptions  *parse.Options  // 解析选项
	RenderOptions *render.Options // 渲染选项

	HTML2MdRendererFuncs          map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2Md 渲染器函数
	HTML2VditorDOMRendererFuncs   map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2VditorDOM 渲染器函数
	HTML2VditorIRDOMRendererFuncs map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2VditorIRDOM 渲染器函数
	HTML2BlockDOMRendererFuncs    map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2BlockDOM 渲染器函数
	HTML2VditorSVDOMRendererFuncs map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2VditorSVDOM 渲染器函数
	Md2HTMLRendererFuncs          map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2HTML 渲染器函数
	Md2VditorDOMRendererFuncs     map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2VditorDOM 渲染器函数
	Md2VditorIRDOMRendererFuncs   map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2VditorIRDOM 渲染器函数
	Md2BlockDOMRendererFuncs      map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2BlockDOM 渲染器函数
	Md2VditorSVDOMRendererFuncs   map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2VditorSVDOM 渲染器函数
}

// New 创建一个新的 Lute 引擎。
//
// 默认启用的解析选项：
//   - GFM 支持
//   - 脚注
//   - 标题自定义 ID
//   - Emoji 别名替换，比如 :heart: 替换为 ❤️
//   - YAML Front Matter
//
// 默认启用的渲染选项：
//   - 软换行转硬换行
//   - 代码块语法高亮
//   - 中西文间插入空格
//   - 修正术语拼写
//   - 标题自定义 ID
func New(opts ...ParseOption) (ret *Lute) {
	ret = &Lute{ParseOptions: parse.NewOptions(), RenderOptions: render.NewOptions()}
	for _, opt := range opts {
		opt(ret)
	}

	ret.HTML2MdRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2VditorDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2VditorIRDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2BlockDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2VditorSVDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2HTMLRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2VditorDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2VditorIRDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2BlockDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2VditorSVDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	return ret
}

// Markdown 将 markdown 文本字节数组处理为相应的 html 字节数组。name 参数仅用于标识文本，比如可传入 id 或者标题，也可以传入 ""。
func (lute *Lute) Markdown(name string, markdown []byte) (html []byte) {
	tree := parse.Parse(name, markdown, lute.ParseOptions)
	renderer := render.NewHtmlRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2HTMLRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	html = renderer.Render()
	return
}

// MarkdownStr 接受 string 类型的 markdown 后直接调用 Markdown 进行处理。
func (lute *Lute) MarkdownStr(name, markdown string) (html string) {
	htmlBytes := lute.Markdown(name, []byte(markdown))
	html = util.BytesToStr(htmlBytes)
	return
}

// Format 将 markdown 文本字节数组进行格式化。
func (lute *Lute) Format(name string, markdown []byte) (formatted []byte) {
	tree := parse.Parse(name, markdown, lute.ParseOptions)
	renderer := render.NewFormatRenderer(tree, lute.RenderOptions)
	formatted = renderer.Render()
	return
}

// FormatStr 接受 string 类型的 markdown 后直接调用 Format 进行处理。
func (lute *Lute) FormatStr(name, markdown string) (formatted string) {
	formattedBytes := lute.Format(name, []byte(markdown))
	formatted = util.BytesToStr(formattedBytes)
	return
}

// TextBundle 将 markdown 文本字节数组进行 TextBundle 处理。
func (lute *Lute) TextBundle(name string, markdown []byte, linkPrefixes []string) (textbundle []byte, originalLinks []string) {
	tree := parse.Parse(name, markdown, lute.ParseOptions)
	renderer := render.NewTextBundleRenderer(tree, linkPrefixes, lute.RenderOptions)
	textbundle, originalLinks = renderer.Render()
	return
}

// TextBundleStr 接受 string 类型的 markdown 后直接调用 TextBundle 进行处理。
func (lute *Lute) TextBundleStr(name, markdown string, linkPrefixes []string) (textbundle string, originalLinks []string) {
	textbundleBytes, originalLinks := lute.TextBundle(name, []byte(markdown), linkPrefixes)
	textbundle = util.BytesToStr(textbundleBytes)
	return
}

// HTML2Text 将指定的 HTMl dom 转换为文本。
func (lute *Lute) HTML2Text(dom string) string {
	tree := lute.HTML2Tree(dom)
	if nil == tree {
		return ""
	}
	return tree.Root.Text()
}

// RenderJSON 用于渲染 JSON 格式数据。
func (lute *Lute) RenderJSON(markdown string) (json string) {
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewJSONRenderer(tree, lute.RenderOptions)
	output := renderer.Render()
	json = util.BytesToStr(output)
	return
}

// Space 用于在 text 中的中西文之间插入空格。
func (lute *Lute) Space(text string) string {
	return render.Space0(text)
}

// IsValidLinkDest 判断 str 是否为合法的链接地址。
func (lute *Lute) IsValidLinkDest(str string) bool {
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, "[") {
		return false
	}

	luteEngine := New()
	luteEngine.ParseOptions.GFMAutoLink = true
	tree := parse.Parse("", []byte(str), luteEngine.ParseOptions)
	if nil == tree.Root.FirstChild || nil == tree.Root.FirstChild.FirstChild {
		return false
	}
	if tree.Root.LastChild != tree.Root.FirstChild {
		return false
	}
	if ast.NodeLink != tree.Root.FirstChild.FirstChild.Type {
		return false
	}
	return true
}

func (lute *Lute) GetLinkDest(str string) string {
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, "file://") {
		return str
	}

	luteEngine := New()
	luteEngine.ParseOptions.GFMAutoLink = true
	tree := parse.Parse("", []byte(str), luteEngine.ParseOptions)
	if nil == tree.Root.FirstChild || nil == tree.Root.FirstChild.FirstChild {
		return ""
	}
	if tree.Root.LastChild != tree.Root.FirstChild {
		return ""
	}
	if ast.NodeLink != tree.Root.FirstChild.FirstChild.Type {
		return ""
	}
	return tree.Root.FirstChild.FirstChild.ChildByType(ast.NodeLinkDest).TokensStr()
}

// GetEmojis 返回 Emoji 别名和对应 Unicode 字符的字典列表。
func (lute *Lute) GetEmojis() (ret map[string]string) {
	parse.EmojiLock.Lock()
	defer parse.EmojiLock.Unlock()

	ret = make(map[string]string, len(lute.ParseOptions.AliasEmoji))
	placeholder := util.BytesToStr(parse.EmojiSitePlaceholder)
	for k, v := range lute.ParseOptions.AliasEmoji {
		if strings.Contains(v, placeholder) {
			v = strings.ReplaceAll(v, placeholder, lute.ParseOptions.EmojiSite)
		}
		ret[k] = v
	}
	return
}

// PutEmojis 将指定的 emojiMap 合并覆盖已有的 Emoji 字典。
func (lute *Lute) PutEmojis(emojiMap map[string]string) {
	parse.EmojiLock.Lock()
	defer parse.EmojiLock.Unlock()

	for k, v := range emojiMap {
		lute.ParseOptions.AliasEmoji[k] = v
		lute.ParseOptions.EmojiAlias[v] = k
	}
}

// RemoveEmoji 用于删除 str 中的 Emoji Unicode。
func (lute *Lute) RemoveEmoji(str string) string {
	parse.EmojiLock.Lock()
	defer parse.EmojiLock.Unlock()

	for u := range lute.ParseOptions.EmojiAlias {
		str = strings.ReplaceAll(str, u, "")
	}
	return strings.TrimSpace(str)
}

// GetTerms 返回术语字典。
func (lute *Lute) GetTerms() map[string]string {
	return lute.RenderOptions.Terms
}

// PutTerms 将制定的 termMap 合并覆盖已有的术语字典。
func (lute *Lute) PutTerms(termMap map[string]string) {
	for k, v := range termMap {
		lute.RenderOptions.Terms[k] = v
	}
}

var (
	formatRendererSync = render.NewFormatRenderer(nil, nil)
	formatRendererLock = sync.Mutex{}
)

func FormatNodeSync(node *ast.Node, parseOptions *parse.Options, renderOptions *render.Options) (ret string, err error) {
	formatRendererLock.Lock()
	defer formatRendererLock.Unlock()
	defer util.RecoverPanic(&err)

	root := &ast.Node{Type: ast.NodeDocument}
	tree := &parse.Tree{Root: root, Context: &parse.Context{ParseOption: parseOptions}}
	formatRendererSync.Tree = tree
	formatRendererSync.Options = renderOptions
	formatRendererSync.LastOut = lex.ItemNewline
	formatRendererSync.NodeWriterStack = []*bytes.Buffer{formatRendererSync.Writer}

	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := formatRendererSync.RendererFuncs[n.Type]
		if nil == rendererFunc {
			err = errors.New("not found renderer for node [type=" + n.Type.String() + "]")
			return ast.WalkStop
		}
		return rendererFunc(n, entering)
	})

	ret = strings.TrimSpace(formatRendererSync.Writer.String())
	formatRendererSync.Tree = nil
	formatRendererSync.Options = nil
	formatRendererSync.Writer.Reset()
	formatRendererSync.NodeWriterStack = nil
	return
}

var (
	protyleExportMdRendererSync = render.NewProtyleExportMdRenderer(nil, nil)
	protyleExportMdRendererLock = sync.Mutex{}
)

func ProtyleExportMdNodeSync(node *ast.Node, parseOptions *parse.Options, renderOptions *render.Options) (ret string, err error) {
	protyleExportMdRendererLock.Lock()
	defer protyleExportMdRendererLock.Unlock()
	defer util.RecoverPanic(&err)

	root := &ast.Node{Type: ast.NodeDocument}
	tree := &parse.Tree{Root: root, Context: &parse.Context{ParseOption: parseOptions}}
	protyleExportMdRendererSync.Tree = tree
	protyleExportMdRendererSync.Options = renderOptions
	protyleExportMdRendererSync.LastOut = lex.ItemNewline
	protyleExportMdRendererSync.NodeWriterStack = []*bytes.Buffer{protyleExportMdRendererSync.Writer}

	ast.Walk(node, func(n *ast.Node, entering bool) ast.WalkStatus {
		rendererFunc := protyleExportMdRendererSync.RendererFuncs[n.Type]
		if nil == rendererFunc {
			err = errors.New("not found renderer for node [type=" + n.Type.String() + "]")
			return ast.WalkStop
		}
		return rendererFunc(n, entering)
	})

	ret = strings.TrimSpace(protyleExportMdRendererSync.Writer.String())
	protyleExportMdRendererSync.Tree = nil
	protyleExportMdRendererSync.Options = nil
	protyleExportMdRendererSync.Writer.Reset()
	protyleExportMdRendererSync.NodeWriterStack = nil
	return
}

// ProtylePreview 使用指定的 options 渲染 tree 为 Protyle 预览 HTML。
func (lute *Lute) ProtylePreview(tree *parse.Tree, options *render.Options) string {
	renderer := render.NewProtylePreviewRenderer(tree, options)
	output := renderer.Render()
	return util.BytesToStr(output)
}

// Tree2HTML 使用指定的 options 渲染 tree 为标准 HTML。
func (lute *Lute) Tree2HTML(tree *parse.Tree, options *render.Options) string {
	renderer := render.NewHtmlRenderer(tree, options)
	output := renderer.Render()
	return util.BytesToStr(output)
}

// ParseOption 描述了解析选项设置函数签名。
type ParseOption func(lute *Lute)

// 以下 Setters 主要是给 JavaScript 端导出方法用。

func (lute *Lute) SetGFMTable(b bool) {
	lute.ParseOptions.GFMTable = b
}

func (lute *Lute) SetGFMTaskListItem(b bool) {
	lute.ParseOptions.GFMTaskListItem = b
}

func (lute *Lute) SetGFMTaskListItemClass(class string) {
	lute.RenderOptions.GFMTaskListItemClass = class
}

func (lute *Lute) SetGFMStrikethrough(b bool) {
	lute.ParseOptions.GFMStrikethrough = b
}

func (lute *Lute) SetGFMStrikethrough1(b bool) {
	lute.ParseOptions.GFMStrikethrough1 = b
}

func (lute *Lute) SetGFMAutoLink(b bool) {
	lute.ParseOptions.GFMAutoLink = b
}

func (lute *Lute) SetSoftBreak2HardBreak(b bool) {
	lute.RenderOptions.SoftBreak2HardBreak = b
}

func (lute *Lute) SetCodeSyntaxHighlight(b bool) {
	lute.RenderOptions.CodeSyntaxHighlight = b
}

func (lute *Lute) SetCodeSyntaxHighlightDetectLang(b bool) {
	lute.RenderOptions.CodeSyntaxHighlightDetectLang = b
}

func (lute *Lute) SetCodeSyntaxHighlightInlineStyle(b bool) {
	lute.RenderOptions.CodeSyntaxHighlightInlineStyle = b
}

func (lute *Lute) SetCodeSyntaxHighlightLineNum(b bool) {
	lute.RenderOptions.CodeSyntaxHighlightLineNum = b
}

func (lute *Lute) SetCodeSyntaxHighlightStyleName(name string) {
	lute.RenderOptions.CodeSyntaxHighlightStyleName = name
}

func (lute *Lute) SetFootnotes(b bool) {
	lute.ParseOptions.Footnotes = b
}

func (lute *Lute) SetToC(b bool) {
	lute.ParseOptions.ToC = b
	lute.RenderOptions.ToC = b
}

func (lute *Lute) SetHeadingID(b bool) {
	lute.ParseOptions.HeadingID = b
	lute.RenderOptions.HeadingID = b
}

func (lute *Lute) SetAutoSpace(b bool) {
	lute.RenderOptions.AutoSpace = b
}

func (lute *Lute) SetFixTermTypo(b bool) {
	lute.RenderOptions.FixTermTypo = b
}

func (lute *Lute) SetEmoji(b bool) {
	lute.ParseOptions.Emoji = b
}

func (lute *Lute) SetEmojis(emojis map[string]string) {
	lute.ParseOptions.AliasEmoji = emojis
}

func (lute *Lute) SetEmojiSite(emojiSite string) {
	lute.ParseOptions.EmojiSite = emojiSite
}

func (lute *Lute) SetHeadingAnchor(b bool) {
	lute.RenderOptions.HeadingAnchor = b
}

func (lute *Lute) SetTerms(terms map[string]string) {
	lute.RenderOptions.Terms = terms
}

func (lute *Lute) SetVditorWYSIWYG(b bool) {
	lute.ParseOptions.VditorWYSIWYG = b
	lute.RenderOptions.VditorWYSIWYG = b
}

func (lute *Lute) SetProtyleWYSIWYG(b bool) {
	lute.ParseOptions.ProtyleWYSIWYG = b
	lute.RenderOptions.ProtyleWYSIWYG = b
}

func (lute *Lute) SetVditorIR(b bool) {
	lute.ParseOptions.VditorIR = b
	lute.RenderOptions.VditorIR = b
}

func (lute *Lute) SetVditorSV(b bool) {
	lute.ParseOptions.VditorSV = b
	lute.RenderOptions.VditorSV = b
}

func (lute *Lute) SetInlineMath(b bool) {
	lute.ParseOptions.InlineMath = b
}

func (lute *Lute) SetInlineMathAllowDigitAfterOpenMarker(b bool) {
	lute.ParseOptions.InlineMathAllowDigitAfterOpenMarker = b
}

func (lute *Lute) SetLinkPrefix(linkPrefix string) {
	lute.RenderOptions.LinkPrefix = linkPrefix
}

func (lute *Lute) SetLinkBase(linkBase string) {
	lute.RenderOptions.LinkBase = linkBase
}

func (lute *Lute) GetLinkBase() string {
	return lute.RenderOptions.LinkBase
}

func (lute *Lute) SetVditorCodeBlockPreview(b bool) {
	lute.RenderOptions.VditorCodeBlockPreview = b
}

func (lute *Lute) SetVditorMathBlockPreview(b bool) {
	lute.RenderOptions.VditorMathBlockPreview = b
}

func (lute *Lute) SetVditorHTMLBlockPreview(b bool) {
	lute.RenderOptions.VditorHTMLBlockPreview = b
}

func (lute *Lute) SetRenderListStyle(b bool) {
	lute.RenderOptions.RenderListStyle = b
}

// SetSanitize 设置为 true 时表示对输出进行 XSS 过滤。
// 注意：Lute 目前的实现存在一些漏洞，请不要依赖它来防御 XSS 攻击。
func (lute *Lute) SetSanitize(b bool) {
	lute.RenderOptions.Sanitize = b
}

func (lute *Lute) SetImageLazyLoading(dataSrc string) {
	lute.RenderOptions.ImageLazyLoading = dataSrc
}

func (lute *Lute) SetChineseParagraphBeginningSpace(b bool) {
	lute.RenderOptions.ChineseParagraphBeginningSpace = b
}

func (lute *Lute) SetYamlFrontMatter(b bool) {
	lute.ParseOptions.YamlFrontMatter = b
}

func (lute *Lute) SetSetext(b bool) {
	lute.ParseOptions.Setext = b
}

func (lute *Lute) SetBlockRef(b bool) {
	lute.ParseOptions.BlockRef = b
}

func (lute *Lute) SetFileAnnotationRef(b bool) {
	lute.ParseOptions.FileAnnotationRef = b
}

func (lute *Lute) SetMark(b bool) {
	lute.ParseOptions.Mark = b
}

func (lute *Lute) SetKramdownIAL(b bool) {
	lute.ParseOptions.KramdownBlockIAL = b
	lute.ParseOptions.KramdownSpanIAL = b
	lute.RenderOptions.KramdownBlockIAL = b
	lute.RenderOptions.KramdownSpanIAL = b
}

func (lute *Lute) SetKramdownBlockIAL(b bool) {
	lute.ParseOptions.KramdownBlockIAL = b
	lute.RenderOptions.KramdownBlockIAL = b
}

func (lute *Lute) SetKramdownSpanIAL(b bool) {
	lute.ParseOptions.KramdownSpanIAL = b
	lute.RenderOptions.KramdownSpanIAL = b
}

func (lute *Lute) SetKramdownIALIDRenderName(name string) {
	lute.RenderOptions.KramdownIALIDRenderName = name
}

func (lute *Lute) SetTag(b bool) {
	lute.ParseOptions.Tag = b
}

func (lute *Lute) SetImgPathAllowSpace(b bool) {
	lute.ParseOptions.ImgPathAllowSpace = b
}

func (lute *Lute) SetSuperBlock(b bool) {
	lute.ParseOptions.SuperBlock = b
	lute.RenderOptions.SuperBlock = b
}

func (lute *Lute) SetSup(b bool) {
	lute.ParseOptions.Sup = b
}

func (lute *Lute) SetSub(b bool) {
	lute.ParseOptions.Sub = b
}

func (lute *Lute) SetInlineAsterisk(b bool) {
	lute.ParseOptions.InlineAsterisk = b
}

func (lute *Lute) SetInlineUnderscore(b bool) {
	lute.ParseOptions.InlineUnderscore = b
}

func (lute *Lute) SetGitConflict(b bool) {
	lute.ParseOptions.GitConflict = b
}

func (lute *Lute) SetLinkRef(b bool) {
	lute.ParseOptions.LinkRef = b
}

func (lute *Lute) SetIndentCodeBlock(b bool) {
	lute.ParseOptions.IndentCodeBlock = b
}

func (lute *Lute) SetDataImage(b bool) {
	lute.ParseOptions.DataImage = b
}

func (lute *Lute) SetTextMark(b bool) {
	lute.ParseOptions.TextMark = b
}

func (lute *Lute) SetSpin(b bool) {
	lute.ParseOptions.Spin = b
}

func (lute *Lute) SetHTMLTag2TextMark(b bool) {
	lute.ParseOptions.HTMLTag2TextMark = b
}

func (lute *Lute) SetParagraphBeginningSpace(b bool) {
	lute.ParseOptions.ParagraphBeginningSpace = b
	lute.RenderOptions.KeepParagraphBeginningSpace = b
}

func (lute *Lute) SetProtyleMarkNetImg(b bool) {
	lute.RenderOptions.ProtyleMarkNetImg = b
}

func (lute *Lute) SetSpellcheck(b bool) {
	lute.RenderOptions.Spellcheck = b
}

func (lute *Lute) SetJSRenderers(options map[string]map[string]*js.Object) {
	for rendererType, extRenderer := range options["renderers"] {
		switch extRenderer.Interface().(type) { // 稍微进行一点格式校验
		case map[string]interface{}:
			break
		default:
			panic("invalid type [" + rendererType + "]")
		}

		var rendererFuncs map[ast.NodeType]render.ExtRendererFunc
		if "HTML2Md" == rendererType {
			rendererFuncs = lute.HTML2MdRendererFuncs
		} else if "HTML2VditorDOM" == rendererType {
			rendererFuncs = lute.HTML2VditorDOMRendererFuncs
		} else if "HTML2VditorIRDOM" == rendererType {
			rendererFuncs = lute.HTML2VditorIRDOMRendererFuncs
		} else if "HTML2BlockDOM" == rendererType {
			rendererFuncs = lute.HTML2BlockDOMRendererFuncs
		} else if "HTML2VditorSVDOM" == rendererType {
			rendererFuncs = lute.HTML2VditorSVDOMRendererFuncs
		} else if "Md2HTML" == rendererType {
			rendererFuncs = lute.Md2HTMLRendererFuncs
		} else if "Md2VditorDOM" == rendererType {
			rendererFuncs = lute.Md2VditorDOMRendererFuncs
		} else if "Md2VditorIRDOM" == rendererType {
			rendererFuncs = lute.Md2VditorIRDOMRendererFuncs
		} else if "Md2BlockDOM" == rendererType {
			rendererFuncs = lute.Md2BlockDOMRendererFuncs
		} else if "Md2VditorSVDOM" == rendererType {
			rendererFuncs = lute.Md2VditorSVDOMRendererFuncs
		} else {
			panic("unknown ext renderer func [" + rendererType + "]")
		}

		extRenderer := extRenderer // https://go.dev/blog/loopvar-preview
		renderFuncs := extRenderer.Interface().(map[string]interface{})
		for funcName := range renderFuncs {
			nodeType := "Node" + funcName[len("render"):]
			rendererFuncs[ast.Str2NodeType(nodeType)] = func(node *ast.Node, entering bool) (string, ast.WalkStatus) {
				funcName = "render" + node.Type.String()[len("Node"):]
				ret := extRenderer.Call(funcName, js.MakeWrapper(node), entering).Interface().([]interface{})
				return ret[0].(string), ast.WalkStatus(ret[1].(float64))
			}
		}
	}
}
