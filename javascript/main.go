// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"bytes"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
	"github.com/gopherjs/gopherjs/js"
)

func New(options map[string]map[string]*js.Object) *js.Object {
	engine := lute.New()
	engine.SetJSRenderers(options)
	return js.MakeWrapper(engine)
}

func main() {
	js.Global.Set("Lute", map[string]interface{}{
		"Version":          lute.Version,
		"New":              New,
		"WalkStop":         ast.WalkStop,
		"WalkSkipChildren": ast.WalkSkipChildren,
		"WalkContinue":     ast.WalkContinue,
		"GetHeadingID":     render.HeadingID,
		"RenderMindmap":    RenderMindmap,
	})
}

func RenderMindmap(listContent string) string {
	tree := parse.Parse("", []byte(listContent), lute.NewOptions())
	buf := &bytes.Buffer{}
	ast.Walk(tree.Root, func(n *ast.Node, entering bool) ast.WalkStatus {
		switch n.Type {
		case ast.NodeDocument, ast.NodeList:
			return ast.WalkContinue
		case ast.NodeListItem:
			children := nil != n.ChildByType(ast.NodeList)
			if entering {
				buf.WriteString("{\"name\": \"" + n.Text())
				if children {
					buf.WriteString("\", \"children\": [")
				}
			} else {
				if children {
					buf.WriteString("]")
				}
				buf.WriteString("}")
				if nil != n.Next {
					buf.WriteString(", ")
				}
			}
		default:
			return ast.WalkStop
		}
		return ast.WalkContinue
	})
	return buf.String()
}
