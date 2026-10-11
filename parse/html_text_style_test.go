// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package parse

import (
	"strings"
	"testing"

	"github.com/88250/lute/html"
	"github.com/88250/lute/html/atom"
)

func TestResolveHTMLTextStyle(t *testing.T) {
	parent := HTMLTextStyle{Color: "red", Family: "Arial", Size: 24, SizeUnit: "px"}
	for _, tc := range []struct {
		name, style string
		tag         atom.Atom
		want        HTMLTextStyle
	}{
		{"inherit", "", atom.Span, parent},
		{"small", "", atom.Small, HTMLTextStyle{Color: "red", Family: "Arial", Size: 20, SizeUnit: "px"}},
		{"explicit small", "font-size:18px", atom.Small, HTMLTextStyle{Color: "red", Family: "Arial", Size: 18, SizeUnit: "px"}},
		{"percentage", "font-size:50%", atom.Span, HTMLTextStyle{Color: "red", Family: "Arial", Size: 12, SizeUnit: "px"}},
		{"em", "font-size:1.5em", atom.Span, HTMLTextStyle{Color: "red", Family: "Arial", Size: 36, SizeUnit: "px"}},
		{"pt", "font-size:12pt", atom.Span, HTMLTextStyle{Color: "red", Family: "Arial", Size: 16, SizeUnit: "px"}},
		{"priority", "COLOR: blue !important;color: green;background-color:currentColor", atom.Span, HTMLTextStyle{Color: "blue", Background: "blue", Family: "Arial", Size: 24, SizeUnit: "px"}},
		{"family", `font-family: 'Noto Sans', "Microsoft YaHei", sans-serif`, atom.Span, HTMLTextStyle{Color: "red", Family: `'Noto Sans', "Microsoft YaHei", sans-serif`, Size: 24, SizeUnit: "px"}},
		{"unsafe", `color:url(https://example.com);background-color:expression(alert(1));font-family:"x";position:fixed;font-size:calc(1px + 2px)`, atom.Span, HTMLTextStyle{Color: "red", Family: `"x"`, Size: 24, SizeUnit: "px"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := ResolveHTMLTextStyle(&html.Node{DataAtom: tc.tag, Attr: []*html.Attribute{{Key: "style", Val: tc.style}}}, parent)
			if got != tc.want {
				t.Errorf("got %#v, want %#v", got, tc.want)
			}
			if strings.Contains(got.CSS(), "position") {
				t.Fatal("layout style leaked")
			}
		})
	}
}

func TestHTMLStyleDeclarations(t *testing.T) {
	got := HTMLStyleDeclarations(`font-family:"A;B";color:rgb(1, 2, 3); COLOR:blue ! important;color:red;`)
	if got["font-family"] != `"A;B"` || got["color"] != "blue" {
		t.Fatal(got)
	}
}

func TestHTMLBackgroundShorthand(t *testing.T) {
	for _, tc := range []struct{ css, expected string }{
		{"background-color:red;background:blue", "blue"},
		{"background:blue;background-color:red", "red"},
		{"background:blue!important;background-color:red", "blue"},
		{"background-color:red!important;background:blue", "red"},
		{"background:rgb(255,165,61)", "rgb(255,165,61)"},
		{"background:var(--b3-font-background6)", "var(--b3-font-background6)"},
		{"background-color:red;background:var(--b3-font-background6)", "var(--b3-font-background6)"},
		{"background:url(https://example.com/image.png)", ""},
	} {
		got := ResolveHTMLTextStyle(&html.Node{Attr: []*html.Attribute{{Key: "style", Val: tc.css}}}, HTMLTextStyle{})
		if got.Background != tc.expected {
			t.Errorf("%s: got %q, want %q", tc.css, got.Background, tc.expected)
		}
	}
}

func TestHTMLTextThemeColors(t *testing.T) {
	for _, value := range []string{"var(--b3-font-color6)", "var(--b3-font-background13)", "VAR( --b3-font-color10 )"} {
		got := ResolveHTMLTextStyle(&html.Node{Attr: []*html.Attribute{{Key: "style", Val: "color:" + value + ";background-color:" + value}}}, HTMLTextStyle{})
		if got.Color != value || got.Background != value {
			t.Errorf("lost theme color %q: %#v", value, got)
		}
	}
	for _, value := range []string{"var(--other)", "var(--B3-font-color6)", "var(--b3-font-color0)", "var(--b3-font-color6, url(https://example.com))", "var(--b3-font-color6) url(https://example.com)", "var(--b3-font-color6"} {
		got := ResolveHTMLTextStyle(&html.Node{Attr: []*html.Attribute{{Key: "style", Val: "color:" + value}}}, HTMLTextStyle{})
		if got.Color != "" {
			t.Errorf("accepted unsupported color %q: %#v", value, got)
		}
	}
}
