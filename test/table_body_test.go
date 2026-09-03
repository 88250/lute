// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of the License at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/render"
)

func TestTableRenderersOmitEmptyBody(t *testing.T) {
	luteEngine := lute.New()
	blockDOM, tree := luteEngine.Md2BlockDOMTree("| foo |\n| --- |\n", false)
	outputs := map[string]string{
		"block DOM": blockDOM,
		"preview":   luteEngine.ProtylePreview(tree, luteEngine.RenderOptions, luteEngine.ParseOptions),
		"export": string(render.NewProtyleExportRenderer(
			tree, luteEngine.RenderOptions, luteEngine.ParseOptions).Render()),
	}
	for name, output := range outputs {
		if strings.Contains(output, "<tbody") || strings.Contains(output, "</tbody>") {
			t.Fatalf("%s contains an empty table body: %q", name, output)
		}
	}
}

func TestTableRenderersKeepNonemptyBody(t *testing.T) {
	luteEngine := lute.New()
	blockDOM, tree := luteEngine.Md2BlockDOMTree("| foo |\n| --- |\n| bar |\n", false)
	outputs := map[string]string{
		"block DOM": blockDOM,
		"preview":   luteEngine.ProtylePreview(tree, luteEngine.RenderOptions, luteEngine.ParseOptions),
		"export": string(render.NewProtyleExportRenderer(
			tree, luteEngine.RenderOptions, luteEngine.ParseOptions).Render()),
	}
	for name, output := range outputs {
		if !strings.Contains(output, "<tbody") || !strings.Contains(output, "</tbody>") {
			t.Fatalf("%s omits a nonempty table body: %q", name, output)
		}
	}
}
