// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/html"
	"github.com/88250/lute/render"
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	js.Global.Set("Lute", map[string]interface{}{
		"Version":           lute.Version,
		"New":               New,
		"WalkStop":          ast.WalkStop,
		"WalkSkipChildren":  ast.WalkSkipChildren,
		"WalkContinue":      ast.WalkContinue,
		"GetHeadingID":      render.HeadingID,
		"Caret":             editor.Caret,
		"NewNodeID":         ast.NewNodeID,
		"EscapeHTMLStr":     html.EscapeHTMLStr,
		"UnEscapeHTMLStr":   html.UnescapeHTMLStr,
		"EChartsMindmapStr": render.EChartsMindmapStr,
		"Sanitize":          render.Sanitize,
		"BlockDOM2Content":  BlockDOM2Content,
	})
}

func New(options map[string]map[string]*js.Object) *js.Object {
	engine := lute.New()
	engine.SetJSRenderers(options)
	return js.MakeWrapper(engine)
}

func BlockDOM2Content(dom string) string {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetBlockRef(true)
	luteEngine.SetFileAnnotationRef(true)
	luteEngine.SetKramdownIAL(true)
	luteEngine.SetTag(true)
	luteEngine.SetSuperBlock(true)
	luteEngine.SetImgPathAllowSpace(true)
	luteEngine.SetGitConflict(true)
	luteEngine.SetMark(true)
	luteEngine.SetSup(true)
	luteEngine.SetSub(true)
	luteEngine.SetInlineMathAllowDigitAfterOpenMarker(true)
	luteEngine.SetFootnotes(false)
	luteEngine.SetToC(false)
	luteEngine.SetIndentCodeBlock(false)
	luteEngine.SetParagraphBeginningSpace(true)
	luteEngine.SetAutoSpace(false)
	luteEngine.SetHeadingID(false)
	luteEngine.SetSetext(false)
	luteEngine.SetYamlFrontMatter(false)
	luteEngine.SetLinkRef(false)
	luteEngine.SetCodeSyntaxHighlight(false)
	luteEngine.SetSanitize(true)
	return luteEngine.BlockDOM2Content(dom)
}
