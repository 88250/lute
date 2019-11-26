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
	"strings"
)

// Lute 描述了 Lute 引擎的顶层使用入口。
type Lute struct {
	*options
}

// New 创建一个新的 Lute 引擎，默认启用：
//  * GFM 支持
//  * 代码块语法高亮
//  * 软换行转硬换行
//  * 中西文间插入空格
//  * 修正术语拼写
//  * 替换中文标点
//  * Emoji 别名替换，比如 :heart: 替换为 ❤️
//  * 并行解析
func New(opts ...option) (ret *Lute) {
	ret = &Lute{options: &options{}}
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
	ret.AutoSpace = true
	ret.FixTermTypo = true
	ret.ChinesePunct = true
	ret.Emoji = true
	ret.AliasEmoji, ret.EmojiAlias = newEmojis()
	ret.Terms = newTerms()
	ret.EmojiSite = "https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji"
	ret.ParallelParsing = true
	for _, opt := range opts {
		opt(ret)
	}
	return ret
}

// Markdown 将 markdown 文本字符数组处理为相应的 html 字符数组。name 参数仅用于标识文本，比如可传入 id 或者标题，也可以传入 ""。
func (lute *Lute) Markdown(name string, markdown []byte) (html []byte, err error) {
	var tree *Tree
	tree, err = lute.parse(name, markdown)
	if nil != err {
		return
	}

	renderer := lute.newHTMLRenderer(tree)
	html, err = renderer.Render()
	return
}

// MarkdownStr 接受 string 类型的 markdown 后直接调用 Markdown 进行处理。
func (lute *Lute) MarkdownStr(name, markdown string) (html string, err error) {
	var htmlBytes []byte
	htmlBytes, err = lute.Markdown(name, []byte(markdown))
	if nil != err {
		return
	}

	html = bytesToStr(htmlBytes)
	return
}

// Format 将 markdown 文本字符数组进行格式化。
func (lute *Lute) Format(name string, markdown []byte) (formatted []byte, err error) {
	var tree *Tree
	tree, err = lute.parse(name, markdown)
	if nil != err {
		return
	}

	renderer := lute.newFormatRenderer(tree)
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

	formatted = bytesToStr(formattedBytes)
	return
}

// Space 用于在 text 中的中西文之间插入空格。
func (lute *Lute) Space(text string) string {
	return space0(text)
}

// GetEmojis 返回 Emoji 别名和对应 Unicode 字符的字典列表。
func (lute *Lute) GetEmojis() (ret map[string]string) {
	ret = make(map[string]string, len(lute.AliasEmoji))
	placeholder := bytesToStr(emojiSitePlaceholder)
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

// options 描述了一些列解析和渲染选项。
type options struct {
	// GFMTable 设置是否打开“GFM 表”支持。
	GFMTable bool
	// GFMTaskListItem 设置是否打开“GFM 任务列表项”支持。
	GFMTaskListItem bool
	// GFMTaskListItemClass 作为 GFM 任务列表项类名，默认为 "vditor-task"。
	GFMTaskListItemClass string
	// GFMStrikethrough 设置是否打开“GFM 删除线”支持。
	GFMStrikethrough bool
	// GFMAutoLink 设置是否打开“GFM 自动链接”支持。
	GFMAutoLink bool
	// SoftBreak2HardBreak 设置是否将软换行（\n）渲染为硬换行（<br />）。
	SoftBreak2HardBreak bool
	// CodeSyntaxHighlight 设置是否对代码块进行语法高亮。
	CodeSyntaxHighlight bool
	// CodeSyntaxHighlightInlineStyle 设置语法高亮是否为内联样式，默认不内联。
	CodeSyntaxHighlightInlineStyle bool
	// CodeSyntaxHightLineNum 设置语法高亮是否显示行号，默认不显示。
	CodeSyntaxHighlightLineNum bool
	// CodeSyntaxHighlightStyleName 指定语法高亮样式名，默认为 "github"。
	CodeSyntaxHighlightStyleName string
	// AutoSpace 设置是否对普通文本中的中西文间自动插入空格。
	// https://github.com/sparanoid/chinese-copywriting-guidelines
	AutoSpace bool
	// FixTermTypo 设置是否对普通文本中出现的术语进行修正。
	// https://github.com/sparanoid/chinese-copywriting-guidelines
	// 注意：开启术语修正的话会默认在中西文之间插入空格。
	FixTermTypo bool
	// ChinesePunct 设置是否对普通文本中出现中文后跟英文逗号句号等标点替换为中文对应标点。
	ChinesePunct bool
	// Emoji 设置是否对 Emoji 别名替换为原生 Unicode 字符。
	Emoji bool
	// AliasEmoji 存储 ASCII 别名到表情 Unicode 映射。
	AliasEmoji map[string]string
	// EmojiAlias 存储表情 Unicode 到 ASCII 别名映射。
	EmojiAlias map[string]string
	// EmojiSite 设置图片 Emoji URL 的路径前缀。
	EmojiSite string
	// HeadingAnchor 设置是否对标题生成链接锚点。
	HeadingAnchor bool
	// Terms 将传入的 terms 合并覆盖到已有的 Terms 字典。
	Terms map[string]string
	// Vditor 所见即所得支持
	VditorWYSIWYG bool
	// ParallelParsing 设置是否启用并行解析。
	ParallelParsing bool
	// InlineMathAllowDigitAfterOpenMarker 设置内联数学公式是否允许起始 $ 后紧跟数字 https://github.com/b3log/lute/issues/38
	InlineMathAllowDigitAfterOpenMarker bool
}

// option 描述了解析渲染选项设置函数签名。
type option func(lute *Lute)

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

func (lute *Lute) SetCodeSyntaxHighlightInlineStyle(b bool) {
	lute.CodeSyntaxHighlightInlineStyle = b
}

func (lute *Lute) SetCodeSyntaxHighlightLineNum(b bool) {
	lute.CodeSyntaxHighlightLineNum = b
}

func (lute *Lute) SetCodeSyntaxHighlightStyleName(name string) {
	lute.CodeSyntaxHighlightStyleName = name
}

func (lute *Lute) SetAutoSpace(b bool) {
	lute.AutoSpace = b
}

func (lute *Lute) SetFixTermTypo(b bool) {
	lute.FixTermTypo = b
}

func (lute *Lute) SetPunct(b bool) {
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

func (lute *Lute) SetParallelParsing(b bool) {
	lute.ParallelParsing = b
}

func (lute *Lute) SetInlineMathAllowDigitAfterOpenMarker(b bool) {
	lute.InlineMathAllowDigitAfterOpenMarker = b
}
