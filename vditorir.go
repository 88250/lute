// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lute

import (
	"strings"

	"github.com/88250/lute/html"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

// SpinVditorIRDOM 自旋 Vditor Instant-Rendering DOM，用于即时渲染模式下的编辑。
func (lute *Lute) SpinVditorIRDOM(ivHTML string) (ovHTML string) {
	lute.VditorIR = true
	lute.VditorWYSIWYG = true

	// 替换插入符
	ivHTML = strings.ReplaceAll(ivHTML, "<wbr>", parse.Caret)
	markdown := lute.vditorIRDOM2Md(ivHTML)
	tree := parse.Parse("", []byte(markdown), lute.Options)
	renderer := render.NewVditorIRRenderer(tree)
	output := renderer.Render()
	if renderer.Option.Footnotes && 0 < len(renderer.Tree.Context.FootnotesDefs) {
		output = renderer.RenderFootnotesDefs(renderer.Tree.Context)
	}
	// 替换插入符
	ovHTML = strings.ReplaceAll(string(output), parse.Caret, "<wbr>")
	return
}

// HTML2VditorIRDOM 将 HTML 转换为 Vditor Instant-Rendering DOM，用于即时渲染模式下粘贴。
func (lute *Lute) HTML2VditorIRDOM(sHTML string) (vHTML string) {
	lute.VditorIR = true
	lute.VditorWYSIWYG = true

	markdown, err := lute.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), lute.Options)
	renderer := render.NewVditorIRRenderer(tree)
	for nodeType, rendererFunc := range lute.HTML2VditorDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	if renderer.Option.Footnotes && 0 < len(renderer.Tree.Context.FootnotesDefs) {
		output = renderer.RenderFootnotesDefs(renderer.Tree.Context)
	}
	vHTML = string(output)
	return
}

// VditorIRDOM2HTML 将 Vditor Instant-Rendering DOM 转换为 HTML，用于 Vditor.getHTML() 接口。
func (lute *Lute) VditorIRDOM2HTML(vhtml string) (sHTML string) {
	lute.VditorIR = true
	lute.VditorWYSIWYG = true

	markdown := lute.vditorIRDOM2Md(vhtml)
	sHTML = lute.Md2HTML(markdown)
	return
}

// Md2VditorIRDOM 将 markdown 转换为 Vditor Instant-Rendering DOM，用于从源码模式切换至所见即所得模式。
func (lute *Lute) Md2VditorIRDOM(markdown string) (vHTML string) {
	lute.VditorIR = true
	lute.VditorWYSIWYG = true

	tree := parse.Parse("", []byte(markdown), lute.Options)
	renderer := render.NewVditorIRRenderer(tree)
	output := renderer.Render()
	if renderer.Option.Footnotes && 0 < len(renderer.Tree.Context.FootnotesDefs) {
		output = renderer.RenderFootnotesDefs(renderer.Tree.Context)
	}
	vHTML = string(output)
	return
}

// VditorIRDOM2Md 将 Vditor Instant-Rendering DOM 转换为 markdown，用于从所见即所得模式切换至源码模式。
func (lute *Lute) VditorIRDOM2Md(htmlStr string) (markdown string) {
	lute.VditorIR = true
	lute.VditorWYSIWYG = true

	htmlStr = strings.ReplaceAll(htmlStr, parse.Zwsp, "")
	markdown = lute.vditorIRDOM2Md(htmlStr)
	markdown = strings.ReplaceAll(markdown, parse.Zwsp, "")
	return
}

func (lute *Lute) vditorIRDOM2Md(htmlStr string) (markdown string) {
	// 删掉插入符
	htmlStr = strings.ReplaceAll(htmlStr, "<wbr>", "")

	// 替换结尾空白，否则 HTML 解析会产生冗余节点导致生成空的代码块
	htmlStr = strings.ReplaceAll(htmlStr, "\t\n", "\n")
	htmlStr = strings.ReplaceAll(htmlStr, "    \n", "  \n")

	// 将字符串解析为 DOM 树

	reader := strings.NewReader(htmlStr)
	htmlRoot := &html.Node{Type: html.ElementNode}
	htmlNodes, err := html.ParseFragment(reader, htmlRoot)
	if nil != err {
		markdown = err.Error()
		return
	}

	// 将 HTML 树转换为 Markdown AST

	var md string
	for _, htmlNode := range htmlNodes {
		md += lute.domText(htmlNode)
	}
	tree := parse.Parse("", []byte(md), lute.Options)

	// 将 AST 进行 Markdown 格式化渲染

	renderer := render.NewFormatRenderer(tree)
	formatted := renderer.Render()
	markdown = string(formatted)
	return
}
