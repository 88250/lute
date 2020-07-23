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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
	"strings"
)

// TextBundleRenderer 描述了 TextBundle 渲染器。https://github.com/88250/lute/issues/77
type TextBundleRenderer struct {
	*FormatRenderer

	linkPrefixes []string // 链接前缀列表
	originalLink []string // 原始链接列表
}

// NewTextBundleRenderer 创建一个 TextBundle 渲染器。
func NewTextBundleRenderer(tree *parse.Tree, linkPrefixes []string) *TextBundleRenderer {
	ret := &TextBundleRenderer{FormatRenderer: NewFormatRenderer(tree), linkPrefixes: linkPrefixes}
	ret.RendererFuncs[ast.NodeLinkDest] = ret.renderLinkDest
	return ret
}

func (r *TextBundleRenderer) Render() (output []byte, originalLink []string) {
	output = r.FormatRenderer.Render()
	originalLink = r.originalLink
	return
}

func (r *TextBundleRenderer) renderLinkDest(node *ast.Node, entering bool) ast.WalkStatus {
	dest := util.BytesToStr(node.Tokens)
	for _, linkPrefix := range r.linkPrefixes {
		if "" != linkPrefix && strings.HasPrefix(dest, linkPrefix) {
			r.originalLink = append(r.originalLink, dest)
			dest = "assets" + dest[len(linkPrefix):]
		}
	}
	r.WriteString(dest)
	return ast.WalkStop
}
