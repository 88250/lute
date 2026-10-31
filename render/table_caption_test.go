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
	"testing"
)

func TestRenderTableCaption(t *testing.T) {
	tests := []struct {
		name     string
		caption  string
		editable bool
		want     string
	}{
		{
			name:     "editor keeps contenteditable",
			caption:  `&lt;caption contenteditable=&quot;false&quot;&gt;Title&lt;/caption&gt;`,
			editable: true,
			want:     `<caption contenteditable="false">Title</caption>`,
		},
		{
			name:     "preview omits contenteditable",
			caption:  `&lt;caption contenteditable=&quot;false&quot;&gt;Title&lt;/caption&gt;`,
			editable: false,
			want:     `<caption>Title</caption>`,
		},
		{
			name:     "bottom caption",
			caption:  `&lt;caption contenteditable=&quot;false&quot; style=&quot;caption-side: bottom;&quot;&gt;A &amp;amp; B&lt;/caption&gt;`,
			editable: true,
			want:     `<caption contenteditable="false" style="caption-side: bottom;">A &amp; B</caption>`,
		},
		{
			name:     "image payload is removed",
			caption:  `&lt;caption contenteditable=&quot;false&quot;&gt;&lt;img src=x onerror=window.xssFired=true&gt;&lt;/caption&gt;`,
			editable: true,
			want:     `<caption contenteditable="false"></caption>`,
		},
		{
			name:     "text around payload is preserved",
			caption:  `&lt;caption&gt;hello&lt;img src=x onerror=alert(1)&gt;world&lt;/caption&gt;`,
			editable: false,
			want:     `<caption>helloworld</caption>`,
		},
		{
			name:     "script content becomes text",
			caption:  `&lt;caption&gt;&lt;script&gt;alert(1)&lt;/script&gt;text&lt;/caption&gt;`,
			editable: false,
			want:     `<caption>alert(1)text</caption>`,
		},
		{
			name:     "event attribute is dropped",
			caption:  `&lt;caption onmouseover=&quot;alert(1)&quot;&gt;x&lt;/caption&gt;`,
			editable: true,
			want:     `<caption>x</caption>`,
		},
		{
			name:     "unknown style is dropped",
			caption:  `&lt;caption style=&quot;background: url(javascript:alert(1))&quot;&gt;x&lt;/caption&gt;`,
			editable: false,
			want:     `<caption>x</caption>`,
		},
		{
			name:     "plain text is escaped",
			caption:  `Title &amp;amp; more`,
			editable: false,
			want:     `<caption>Title &amp; more</caption>`,
		},
		{
			name:     "single encoded payload is decoded then removed",
			caption:  `&lt;caption&gt;&lt;img src=x onerror=alert(1)&gt;&lt;/caption&gt;`,
			editable: true,
			want:     `<caption></caption>`,
		},
		{
			name:     "double encoded payload stays text",
			caption:  `&amp;lt;caption&amp;gt;&amp;lt;img src=x onerror=alert(1)&amp;gt;&amp;lt;/caption&amp;gt;`,
			editable: true,
			want:     `<caption>&lt;caption&gt;&lt;img src=x onerror=alert(1)&gt;&lt;/caption&gt;</caption>`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := renderTableCaption(test.caption, test.editable)
			if test.want != got {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
			if strings.Contains(got, "<img") || strings.Contains(got, "<script") ||
				strings.Contains(got, `onerror="`) {
				t.Fatalf("unsafe caption output: %q", got)
			}
		})
	}
}

func TestSafeIALAttrs(t *testing.T) {
	attrs := safeIALAttrs([][]string{
		{"style", `width: 100px;" ><img src=x onerror=alert(1)>`},
		{"onmouseover", "alert(1)"},
		{"bad name", "x"},
		{"custom-ok", "value"},
	})
	if 2 != len(attrs) {
		t.Fatalf("expected 2 attrs, got %d: %v", len(attrs), attrs)
	}
	if "custom-ok" != attrs[1][0] || "value" != attrs[1][1] {
		t.Fatalf("unexpected custom attr: %v", attrs[1])
	}
	if strings.Contains(attrs[0][1], `"`) {
		t.Fatalf("quote was not escaped: %q", attrs[0][1])
	}
}
