// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func newIALHandlerTestLute() *lute.Lute {
	luteEngine := lute.New()
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetKramdownIAL(true)
	luteEngine.SetSanitize(true)
	return luteEngine
}

func TestProtylePreviewDropsHandlerIAL(t *testing.T) {
	luteEngine := newIALHandlerTestLute()
	tests := []struct {
		name     string
		markdown string
		element  string
	}{
		{name: "block IAL", markdown: "para\n{: onmouseover=\"window.xssFired=true\"}\n", element: "<p"},
		{name: "inline IAL", markdown: "**bold**{: onmouseover=\"window.xssFired=true\"}\n", element: "<strong"},
		{name: "upper case IAL", markdown: "para\n{: OnClick=\"window.xssFired=true\"}\n", element: "<p"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			html := luteEngine.ProtylePreviewStr("", test.markdown)
			if strings.Contains(strings.ToLower(html), "onclick") || strings.Contains(strings.ToLower(html), "onmouseover") ||
				strings.Contains(html, "xssFired") {
				t.Fatalf("preview contains a handler attribute: %s", html)
			}
			if !strings.Contains(html, test.element) {
				t.Fatalf("preview lost element %q: %s", test.element, html)
			}
		})
	}
}

func TestProtylePreviewKeepsSafeIAL(t *testing.T) {
	luteEngine := newIALHandlerTestLute()
	html := luteEngine.ProtylePreviewStr("", "para\n{: custom-foo=\"bar\"}\n")
	if !strings.Contains(html, `custom-foo="bar"`) {
		t.Fatalf("preview lost a safe attribute: %s", html)
	}
}

func TestProtyleExportDropsHandlerIAL(t *testing.T) {
	luteEngine := newIALHandlerTestLute()
	tree := parse.Parse("", []byte("text **bold**{: onmouseover=\"window.xssFired=true\"}\n"), luteEngine.ParseOptions)
	html := string(render.NewProtyleExportRenderer(tree, luteEngine.RenderOptions, luteEngine.ParseOptions).Render())
	if strings.Contains(strings.ToLower(html), "onmouseover") || strings.Contains(html, "xssFired") {
		t.Fatalf("export HTML contains a handler attribute: %s", html)
	}
	if !strings.Contains(html, "<strong>bold</strong>") {
		t.Fatalf("export HTML lost the inline markup: %s", html)
	}
}

func TestMarkdownDropsHandlerIAL(t *testing.T) {
	luteEngine := newIALHandlerTestLute()
	html := luteEngine.MarkdownStr("", "para\n{: onmouseover=\"window.xssFired=true\"}\n")
	if strings.Contains(strings.ToLower(html), "onmouseover") || strings.Contains(html, "xssFired") {
		t.Fatalf("HTML contains a handler attribute: %s", html)
	}
	if !strings.Contains(html, "para") {
		t.Fatalf("HTML lost the paragraph: %s", html)
	}
}
