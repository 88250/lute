// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

// Lute 是一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript。
package lute

import (
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/88250/lute/util"
	"github.com/gopherjs/gopherjs/js"
)

const Version = "1.6.6"

// Lute 描述了 Lute 引擎的顶层使用入口。
type Lute struct {
	*parse.Options // 解析和渲染选项配置

	HTML2MdRendererFuncs               map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2Md 渲染器函数
	HTML2VditorDOMRendererFuncs        map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2VditorDOM 渲染器函数
	HTML2VditorIRDOMRendererFuncs      map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2VditorIRDOM 渲染器函数
	HTML2VditorIRBlockDOMRendererFuncs map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2VditorIRBlockDOM 渲染器函数
	HTML2VditorSVDOMRendererFuncs      map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2VditorSVDOM 渲染器函数
	Md2HTMLRendererFuncs               map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2HTML 渲染器函数
	Md2VditorDOMRendererFuncs          map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2VditorDOM 渲染器函数
	Md2VditorIRDOMRendererFuncs        map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2VditorIRDOM 渲染器函数
	Md2VditorIRBlockDOMRendererFuncs   map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2VditorIRBlockDOM 渲染器函数
	Md2VditorSVDOMRendererFuncs        map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 Md2VditorSVDOM 渲染器函数
}

// New 创建一个新的 Lute 引擎，默认启用：
//  * GFM 支持
//  * 代码块语法高亮
//  * 软换行转硬换行
//  * 脚注
//  * 标题自定义 ID
//  * 中西文间插入空格
//  * 修正术语拼写
//  * 替换中文标点
//  * Emoji 别名替换，比如 :heart: 替换为 ❤️
//  * YAML Front Matter
func New(opts ...Option) (ret *Lute) {
	ret = &Lute{Options: NewOptions()}
	for _, opt := range opts {
		opt(ret)
	}

	ret.HTML2MdRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2VditorDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2VditorIRDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2VditorIRBlockDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2VditorSVDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2HTMLRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2VditorDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2VditorIRDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2VditorIRBlockDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.Md2VditorSVDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	return ret
}

func NewOptions() *parse.Options {
	emojis, emoji := parse.NewEmojis()
	return &parse.Options{
		GFMTable:                       true,
		GFMTaskListItem:                true,
		GFMTaskListItemClass:           "vditor-task",
		GFMStrikethrough:               true,
		GFMAutoLink:                    true,
		SoftBreak2HardBreak:            true,
		CodeSyntaxHighlight:            true,
		CodeSyntaxHighlightInlineStyle: false,
		CodeSyntaxHighlightLineNum:     false,
		CodeSyntaxHighlightStyleName:   "github",
		Footnotes:                      true,
		ToC:                            false,
		HeadingID:                      true,
		AutoSpace:                      true,
		FixTermTypo:                    true,
		ChinesePunct:                   true,
		Emoji:                          true,
		AliasEmoji:                     emojis,
		EmojiAlias:                     emoji,
		Terms:                          render.NewTerms(),
		EmojiSite:                      "https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji",
		LinkBase:                       "",
		LinkPrefix:                     "",
		VditorCodeBlockPreview:         true,
		VditorMathBlockPreview:         true,
		RenderListStyle:                false,
		ChineseParagraphBeginningSpace: false,
		YamlFrontMatter:                true,
		BlockRef:                       false,
		Mark:                           false,
		KramdownIAL:                    false,
		KramdownIALIDRenderName:        "id",
	}
}

// Markdown 将 markdown 文本字节数组处理为相应的 html 字节数组。name 参数仅用于标识文本，比如可传入 id 或者标题，也可以传入 ""。
func (lute *Lute) Markdown(name string, markdown []byte) (html []byte) {
	tree := parse.Parse(name, markdown, lute.Options)
	renderer := render.NewHtmlRenderer(tree)
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
	tree := parse.Parse(name, markdown, lute.Options)
	renderer := render.NewFormatRenderer(tree)
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
	tree := parse.Parse(name, markdown, lute.Options)
	renderer := render.NewTextBundleRenderer(tree, linkPrefixes)
	textbundle, originalLinks = renderer.Render()
	return
}

// TextBundle 接受 string 类型的 markdown 后直接调用 TextBundle 进行处理。
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

// Space 用于在 text 中的中西文之间插入空格。
func (lute *Lute) Space(text string) string {
	return render.Space0(text)
}

// IsValidLinkDest 判断 str 是否为合法的链接地址。
func (lute *Lute) IsValidLinkDest(str string) bool {
	luteEngine := New()
	luteEngine.GFMAutoLink = true
	tree := parse.Parse("", []byte(str),luteEngine.Options)
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

// GetEmojis 返回 Emoji 别名和对应 Unicode 字符的字典列表。
func (lute *Lute) GetEmojis() (ret map[string]string) {
	ret = make(map[string]string, len(lute.AliasEmoji))
	placeholder := util.BytesToStr(parse.EmojiSitePlaceholder)
	for k, v := range lute.AliasEmoji {
		if strings.Contains(v, placeholder) {
			v = strings.ReplaceAll(v, placeholder, lute.EmojiSite)
		}
		ret[k] = v
	}
	return
}

// PutEmojis 将指定的 emojiMap 合并覆盖已有的 Emoji 字典。
func (lute *Lute) PutEmojis(emojiMap map[string]string) {
	for k, v := range emojiMap {
		lute.AliasEmoji[k] = v
		lute.EmojiAlias[v] = k
	}
}

// GetTerms 返回术语字典。
func (lute *Lute) GetTerms() map[string]string {
	return lute.Terms
}

// PutTerms 将制定的 termMap 合并覆盖已有的术语字典。
func (lute *Lute) PutTerms(termMap map[string]string) {
	for k, v := range termMap {
		lute.Terms[k] = v
	}
}

// Option 描述了解析渲染选项设置函数签名。
type Option func(lute *Lute)

// 以下 Setters 主要是给 JavaScript 端导出方法用。

func (lute *Lute) SetGFMTable(b bool) {
	lute.GFMTable = b
}

func (lute *Lute) SetGFMTaskListItem(b bool) {
	lute.GFMTaskListItem = b
}

func (lute *Lute) SetGFMTaskListItemClass(class string) {
	lute.GFMTaskListItemClass = class
}

func (lute *Lute) SetGFMStrikethrough(b bool) {
	lute.GFMStrikethrough = b
}

func (lute *Lute) SetGFMAutoLink(b bool) {
	lute.GFMAutoLink = b
}

func (lute *Lute) SetSoftBreak2HardBreak(b bool) {
	lute.SoftBreak2HardBreak = b
}

func (lute *Lute) SetCodeSyntaxHighlight(b bool) {
	lute.CodeSyntaxHighlight = b
}

func (lute *Lute) SetCodeSyntaxHighlightDetectLang(b bool) {
	lute.CodeSyntaxHighlightDetectLang = b
}

func (lute *Lute) SetCodeSyntaxHighlightInlineStyle(b bool) {
	lute.CodeSyntaxHighlightInlineStyle = b
}

func (lute *Lute) SetCodeSyntaxHighlightLineNum(b bool) {
	lute.CodeSyntaxHighlightLineNum = b
}

func (lute *Lute) SetCodeSyntaxHighlightStyleName(name string) {
	lute.CodeSyntaxHighlightStyleName = name
}

func (lute *Lute) SetFootnotes(b bool) {
	lute.Footnotes = b
}

func (lute *Lute) SetToC(b bool) {
	lute.ToC = b
}

func (lute *Lute) SetHeadingID(b bool) {
	lute.HeadingID = b
}

func (lute *Lute) SetAutoSpace(b bool) {
	lute.AutoSpace = b
}

func (lute *Lute) SetFixTermTypo(b bool) {
	lute.FixTermTypo = b
}

func (lute *Lute) SetChinesePunct(b bool) {
	lute.ChinesePunct = b
}

func (lute *Lute) SetEmoji(b bool) {
	lute.Emoji = b
}

func (lute *Lute) SetEmojis(emojis map[string]string) {
	lute.AliasEmoji = emojis
}

func (lute *Lute) SetEmojiSite(emojiSite string) {
	lute.EmojiSite = emojiSite
}

func (lute *Lute) SetHeadingAnchor(b bool) {
	lute.HeadingAnchor = b
}

func (lute *Lute) SetTerms(terms map[string]string) {
	lute.Terms = terms
}

func (lute *Lute) SetVditorWYSIWYG(b bool) {
	lute.VditorWYSIWYG = b
}

func (lute *Lute) SetInlineMathAllowDigitAfterOpenMarker(b bool) {
	lute.InlineMathAllowDigitAfterOpenMarker = b
}

func (lute *Lute) SetLinkPrefix(linkPrefix string) {
	lute.LinkPrefix = linkPrefix
}

func (lute *Lute) SetLinkBase(linkBase string) {
	lute.LinkBase = linkBase
}

func (lute *Lute) GetLinkBase() string {
	return lute.LinkBase
}

func (lute *Lute) SetVditorCodeBlockPreview(b bool) {
	lute.VditorCodeBlockPreview = b
}

func (lute *Lute) SetVditorMathBlockPreview(b bool) {
	lute.VditorMathBlockPreview = b
}

func (lute *Lute) SetRenderListStyle(b bool) {
	lute.RenderListStyle = b
}

func (lute *Lute) SetSanitize(b bool) {
	lute.Sanitize = b
}

func (lute *Lute) SetImageLazyLoading(dataSrc string) {
	lute.ImageLazyLoading = dataSrc
}

func (lute *Lute) SetChineseParagraphBeginningSpace(b bool) {
	lute.ChineseParagraphBeginningSpace = b
}

func (lute *Lute) SetYamlFrontMatter(b bool) {
	lute.YamlFrontMatter = b
}

func (lute *Lute) SetBlockRef(b bool) {
	lute.BlockRef = b
}

func (lute *Lute) SetMark(b bool) {
	lute.Mark = b
}

func (lute *Lute) SetKramdownIAL(b bool) {
	lute.KramdownIAL = b
}

func (lute *Lute) SetKramdownIALIDRenderName(name string) {
	lute.KramdownIALIDRenderName = name
}

func (lute *Lute) SetTag(b bool) {
	lute.Tag = b
}

func (lute *Lute) SetImgPathAllowSpace(b bool) {
	lute.ImgPathAllowSpace = b
}

func (lute *Lute) SetSuperBlock(b bool) {
	lute.SuperBlock = b
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
		} else if "HTML2VditorIRBlockDOM" == rendererType {
			rendererFuncs = lute.HTML2VditorIRBlockDOMRendererFuncs
		} else if "HTML2VditorSVDOM" == rendererType {
			rendererFuncs = lute.HTML2VditorSVDOMRendererFuncs
		} else if "Md2HTML" == rendererType {
			rendererFuncs = lute.Md2HTMLRendererFuncs
		} else if "Md2VditorDOM" == rendererType {
			rendererFuncs = lute.Md2VditorDOMRendererFuncs
		} else if "Md2VditorIRDOM" == rendererType {
			rendererFuncs = lute.Md2VditorIRDOMRendererFuncs
		} else if "Md2VditorIRBlockDOM" == rendererType {
			rendererFuncs = lute.Md2VditorIRBlockDOMRendererFuncs
		} else if "Md2VditorSVDOM" == rendererType {
			rendererFuncs = lute.Md2VditorSVDOMRendererFuncs
		} else {
			panic("unknown ext renderer func [" + rendererType + "]")
		}

		renderFuncs := extRenderer.Interface().(map[string]interface{})
		for funcName := range renderFuncs {
			nodeType := "Node" + funcName[len("render"):]
			rendererFuncs[ast.Str2NodeType(nodeType)] = func(node *ast.Node, entering bool) (string, ast.WalkStatus) {
				nodeType := node.Type.String()
				funcName = "render" + nodeType[len("Node"):]
				ret := extRenderer.Call(funcName, js.MakeWrapper(node), entering).Interface().([]interface{})
				return ret[0].(string), ast.WalkStatus(ret[1].(float64))
			}
		}
	}
}
