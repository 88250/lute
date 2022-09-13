// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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

	"github.com/88250/lute/editor"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

// SpinVditorSVDOM 自旋 Vditor Split-View DOM，用于分屏预览模式下的编辑。
func (lute *Lute) SpinVditorSVDOM(markdown string) (ovHTML string) {
	// 为空的特殊情况处理
	if editor.Caret == strings.TrimSpace(markdown) {
		return "<span data-type=\"text\"><wbr></span>" + string(render.NewlineSV)
	}

	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)

	renderer := render.NewVditorSVRenderer(tree, lute.RenderOptions)
	output := renderer.Render()
	// 替换插入符
	ovHTML = strings.ReplaceAll(string(output), editor.Caret, "<wbr>")
	return
}

// HTML2VditorSVDOM 将 HTML 转换为 Vditor Split-View DOM，用于分屏预览模式下粘贴。
func (lute *Lute) HTML2VditorSVDOM(sHTML string) (vHTML string) {
	markdown, err := lute.HTML2Markdown(sHTML)
	if nil != err {
		vHTML = err.Error()
		return
	}

	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorSVRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.HTML2VditorSVDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	vHTML = string(output)
	return
}

// Md2VditorSVDOM 将 markdown 转换为 Vditor Split-View DOM，用于从源码模式切换至分屏预览模式。
func (lute *Lute) Md2VditorSVDOM(markdown string) (vHTML string) {
	tree := parse.Parse("", []byte(markdown), lute.ParseOptions)
	renderer := render.NewVditorSVRenderer(tree, lute.RenderOptions)
	for nodeType, rendererFunc := range lute.Md2VditorSVDOMRendererFuncs {
		renderer.ExtRendererFuncs[nodeType] = rendererFunc
	}
	output := renderer.Render()
	// 替换插入符
	vHTML = strings.ReplaceAll(string(output), editor.Caret, "<wbr>")
	return
}
