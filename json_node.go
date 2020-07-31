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
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

type JSONNode struct {
	Id      string
	Content string

	Children []*JSONNode
}

// TODO: RenderMarkdown

// RenderJSON 用于渲染 JSON 格式数据。
func (lute *Lute) RenderJSON(markdown string) (ret string) {
	tree := parse.Parse("", []byte(markdown), lute.Options)
	renderer := render.NewJSONRenderer(tree)
	output := renderer.Render()
	ret = string(output)
	return
}
