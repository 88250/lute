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

	"github.com/b3log/lute/html"
	"github.com/b3log/lute/html/atom"
)

// Vditor DOM Parser

// parseVditorDOM 解析 Vditor DOM 生成 Markdown 文本。
func (lute *Lute) parseVditorDOM(htmlStr string) (tree *Tree, err error) {
	defer recoverPanic(&err)

	reader := strings.NewReader(htmlStr)
	doc, err := html.Parse(reader)
	if nil != err {
		return
	}

	// HTML Tree to Markdown AST
	tree = &Tree{Name: "", Root: &Node{typ: NodeDocument}}

	var walker func(*html.Node)
	walker = func(n *html.Node) {
		if html.ElementNode == n.Type && atom.Html == n.DataAtom {
			// Do something with n...
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}

	walker(doc)

	return
}
