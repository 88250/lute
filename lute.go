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

const Version = "1.2.0"

// Lute 描述了 Lute 引擎的顶层使用入口。
type Lute struct {
	*parse.Options // 解析和渲染选项配置

	HTML2MdRendererFuncs        map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2Md 渲染器函数
	HTML2VditorDOMRendererFuncs map[ast.NodeType]render.ExtRendererFunc // 用户自定义的 HTML2VditorDOM 渲染器函数
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
//  * 并行解析
func New(opts ...Option) (ret *Lute) {
	ret = &Lute{Options: &parse.Options{}}
	ret.GFMTable = true
	ret.GFMTaskListItem = true
	ret.GFMTaskListItemClass = "vditor-task"
	ret.GFMStrikethrough = true
	ret.GFMAutoLink = true
	ret.SoftBreak2HardBreak = true
	ret.CodeSyntaxHighlight = true
	ret.CodeSyntaxHighlightInlineStyle = false
	ret.CodeSyntaxHighlightLineNum = false
	ret.CodeSyntaxHighlightStyleName = "github"
	ret.Footnotes = true
	ret.ToC = false
	ret.HeadingID = true
	ret.AutoSpace = true
	ret.FixTermTypo = true
	ret.ChinesePunct = true
	ret.Emoji = true
	ret.AliasEmoji, ret.EmojiAlias = parse.NewEmojis()
	ret.Terms = render.NewTerms()
	ret.EmojiSite = "https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji"
	ret.LinkBase = ""
	for _, opt := range opts {
		opt(ret)
	}
	ret.HTML2MdRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	ret.HTML2VditorDOMRendererFuncs = map[ast.NodeType]render.ExtRendererFunc{}
	return ret
}

// Markdown 将 markdown 文本字节数组处理为相应的 html 字节数组。name 参数仅用于标识文本，比如可传入 id 或者标题，也可以传入 ""。
func (lute *Lute) Markdown(name string, markdown []byte) (html []byte, err error) {
	var tree *parse.Tree
	tree, err = parse.Parse(name, markdown, lute.Options)
	if nil != err {
		return
	}

	renderer := render.NewHtmlRenderer(tree)
	html, err = renderer.Render()

	if lute.Options.Footnotes && 0 < len(tree.Context.FootnotesDefs) {
		html = renderer.(*render.HtmlRenderer).RenderFootnotesDefs(tree.Context)
	}

	return
}

// MarkdownStr 接受 string 类型的 markdown 后直接调用 Markdown 进行处理。
func (lute *Lute) MarkdownStr(name, markdown string) (html string, err error) {
	var htmlBytes []byte
	htmlBytes, err = lute.Markdown(name, []byte(markdown))
	if nil != err {
		return
	}

	html = util.BytesToStr(htmlBytes)
	return
}

// Format 将 markdown 文本字节数组进行格式化。
func (lute *Lute) Format(name string, markdown []byte) (formatted []byte, err error) {
	var tree *parse.Tree
	tree, err = parse.Parse(name, markdown, lute.Options)
	if nil != err {
		return
	}

	renderer := render.NewFormatRenderer(tree)
	formatted, err = renderer.Render()
	return
}

// FormatStr 接受 string 类型的 markdown 后直接调用 Format 进行处理。
func (lute *Lute) FormatStr(name, markdown string) (formatted string, err error) {
	var formattedBytes []byte
	formattedBytes, err = lute.Format(name, []byte(markdown))
	if nil != err {
		return
	}

	formatted = util.BytesToStr(formattedBytes)
	return
}

// Space 用于在 text 中的中西文之间插入空格。
func (lute *Lute) Space(text string) string {
	return render.Space0(text)
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

func (lute *Lute) SetLinkBase(linkBase string) {
	lute.LinkBase = linkBase
}

func (lute *Lute) SetJSRenderers(options map[string]map[string]*js.Object) {
	for rendererType, extRenderer := range options["renderers"] {
		switch extRenderer.Interface().(type) {
		case map[string]interface{}:
			break
		default:
			continue
		}

		var rendererFuncs map[ast.NodeType]render.ExtRendererFunc
		if "HTML2Md" == rendererType {
			rendererFuncs = lute.HTML2MdRendererFuncs
		} else if "HTML2VditorDOM" == rendererType {
			rendererFuncs = lute.HTML2VditorDOMRendererFuncs
		} else {
			continue
		}

		renderFuncs := extRenderer.Interface().(map[string]interface{})
		for funcName, _ := range renderFuncs {
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
