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

import "strings"

// RenderVditorDOM 用于渲染 Vditor DOM。
func (lute *Lute) RenderVditorDOM(htmlStr string) (html string, err error) {
	lute.VditorWYSIWYG = true
	lute.endNewline(&htmlStr)

	var md string
	md, err = lute.Html2Md(htmlStr)
	if nil != err {
		return
	}

	md = strings.ReplaceAll(md, "<wbr>", "\u2038")

	var tree *Tree
	tree, err = lute.parse("", []byte(md))
	if nil != err {
		return
	}

	renderer := lute.newVditorRenderer(tree)
	var output []byte
	output, err = renderer.Render()
	html = string(output)
	html = strings.ReplaceAll(html, "\u2038", "<wbr>")
	return
}
