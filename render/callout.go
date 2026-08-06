// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"strings"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/html"
)

func calloutTitle(node *ast.Node) string {
	if ast.IsBuiltInCalloutType(node.CalloutType) &&
		node.CalloutTitle == ast.GetCalloutTitle(node.CalloutType) &&
		node.CalloutIcon == ast.GetCalloutIcon(node.CalloutType) {
		return ""
	}

	icon := node.CalloutIcon
	if 1 == node.CalloutIconType {
		if ast.IsValidCalloutImageSrc(icon) {
			icon = "![" + ast.CalloutIconImageAlt + "](<" + string(html.EncodeDestination([]byte(icon))) + ">)"
		} else {
			icon = ""
		}
	}
	return strings.TrimSpace(icon + " " + node.CalloutTitle)
}
