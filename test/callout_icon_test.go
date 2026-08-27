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
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
)

func TestCalloutImageIconParse(t *testing.T) {
	tests := []struct {
		name  string
		icon  string
		title string
	}{
		{"local emoji", "/emojis/1/b3log.png", "Local"},
		{"dynamic icon", "api/icon/getDynamicIcon?type=1&color=%23d23f31&date=20260807", "Dynamic"},
		{"http icon", "http://example.com/avatar", "HTTP"},
		{"https icon", "https://example.com/a(b)?foo=bar&baz=qux#fragment", "HTTPS"},
		{"encoded space", "https://example.com/icon%20name.png", "Encoded"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			markdown := "> [!NOTE] ![callout-icon](<" + test.icon + ">) " + test.title + "\n> Content"
			callout := parseCallout(t, markdown)
			if 1 != callout.CalloutIconType || test.icon != callout.CalloutIcon || test.title != callout.CalloutTitle {
				t.Fatalf("unexpected callout icon state: type=%d, icon=%q, title=%q", callout.CalloutIconType,
					callout.CalloutIcon, callout.CalloutTitle)
			}
		})
	}

	callout := parseCallout(t, "> [!NOTE] ![callout-icon](https://example.com/icon.png) Plain\n> Content")
	if 1 != callout.CalloutIconType || "https://example.com/icon.png" != callout.CalloutIcon ||
		"Plain" != callout.CalloutTitle {
		t.Fatalf("unexpected plain callout icon state: type=%d, icon=%q, title=%q", callout.CalloutIconType,
			callout.CalloutIcon, callout.CalloutTitle)
	}
}

func TestCalloutImageIconDoesNotConsumeTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
	}{
		{"plain URL", "https://example.com/icon.png Note"},
		{"regular image", "![cover](<https://example.com/icon.png>) Note"},
		{"unsupported scheme", "![callout-icon](<javascript:alert(1)>) Note"},
		{"data image", "![callout-icon](<data:image/png;base64,AA>) Note"},
		{"protocol relative URL", "![callout-icon](<//example.com/icon.png>) Note"},
		{"missing host", "![callout-icon](<https://>) Note"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			callout := parseCallout(t, "> [!NOTE] "+test.title+"\n> Content")
			if 0 != callout.CalloutIconType || ast.GetCalloutIcon(ast.CalloutTypeNote) != callout.CalloutIcon ||
				test.title != callout.CalloutTitle {
				t.Fatalf("title was consumed as an icon: type=%d, icon=%q, title=%q", callout.CalloutIconType,
					callout.CalloutIcon, callout.CalloutTitle)
			}
		})
	}
}

func TestCalloutImageIconRoundTrip(t *testing.T) {
	icon := "https://b3logfile.com/avatar/1734703705652.png?size=64&theme=light"
	dom := "<div data-subtype=\"NOTE\" data-node-id=\"20260807120000-abcdefg\" " +
		"data-node-index=\"1\" data-type=\"NodeCallout\" class=\"callout\"><div class=\"callout-info\" contenteditable=\"false\">" +
		"<span class=\"callout-icon\"><img class=\"callout-img\" src=\"" + icon + "\"></span>" +
		"<span contenteditable=\"true\" spellcheck=\"false\" class=\"callout-title\">Note</span></div><div class=\"callout-content\">" +
		"<div data-node-id=\"20260807120001-abcdefg\" data-type=\"NodeParagraph\" class=\"p\">" +
		"<div contenteditable=\"true\" spellcheck=\"false\">Content</div>" +
		"<div class=\"protyle-attr\" contenteditable=\"false\">&ZeroWidthSpace;</div></div></div>" +
		"<div class=\"protyle-attr\" contenteditable=\"false\">&ZeroWidthSpace;</div></div>"

	luteEngine := newCalloutLute()
	markdown := luteEngine.BlockDOM2StdMd(dom)
	expected := "> [!NOTE] ![callout-icon](<" + icon + ">) Note\n> Content\n"
	if expected != markdown {
		t.Fatalf("unexpected markdown\nexpected\n\t%q\ngot\n\t%q", expected, markdown)
	}

	formatted := luteEngine.FormatStr("callout-image-icon", markdown)
	formattedAgain := luteEngine.FormatStr("callout-image-icon", formatted)
	if formatted != formattedAgain {
		t.Fatalf("formatting is not idempotent\nfirst\n\t%q\nsecond\n\t%q", formatted, formattedAgain)
	}
	callout := parseCalloutWithLute(t, luteEngine, formattedAgain)
	if 1 != callout.CalloutIconType || icon != callout.CalloutIcon || "Note" != callout.CalloutTitle {
		t.Fatalf("round trip lost icon state: type=%d, icon=%q, title=%q", callout.CalloutIconType,
			callout.CalloutIcon, callout.CalloutTitle)
	}

	blockDOM := luteEngine.Md2BlockDOM(formattedAgain, true)
	if !strings.Contains(blockDOM, "src=\"https://b3logfile.com/avatar/1734703705652.png?size=64&amp;theme=light\"") ||
		!strings.Contains(blockDOM, "<span contenteditable=\"true\" spellcheck=\"false\" class=\"callout-title\">Note</span>") {
		t.Fatalf("unexpected BlockDOM: %s", blockDOM)
	}
}

func TestCalloutEmptyTitleRoundTrip(t *testing.T) {
	luteEngine := newCalloutLute()
	defaultMarkdown := "> [!NOTE]\n> Content"
	defaultCallout := parseCalloutWithLute(t, luteEngine, defaultMarkdown)
	if defaultCallout.CalloutTitleExplicit {
		t.Fatal("default callout title should not be explicit")
	}
	defaultBlockDOM := luteEngine.Md2BlockDOM(defaultMarkdown, true)
	if !strings.Contains(defaultBlockDOM, "class=\"callout-title\">Note</span>") {
		t.Fatalf("default callout title was not rendered: %s", defaultBlockDOM)
	}

	emptyTitleMarkdown := "> [!NOTE] ✏️\n> Content\n"
	emptyTitleCallout := parseCalloutWithLute(t, luteEngine, emptyTitleMarkdown)
	if !emptyTitleCallout.CalloutTitleExplicit || "" != emptyTitleCallout.CalloutTitle {
		t.Fatalf("unexpected empty title state: explicit=%v, title=%q",
			emptyTitleCallout.CalloutTitleExplicit, emptyTitleCallout.CalloutTitle)
	}
	emptyTitleBlockDOM := luteEngine.Md2BlockDOM(emptyTitleMarkdown, true)
	if !strings.Contains(emptyTitleBlockDOM, "class=\"callout-title\"></span>") {
		t.Fatalf("empty callout title was not rendered: %s", emptyTitleBlockDOM)
	}
	if markdown := luteEngine.BlockDOM2StdMd(emptyTitleBlockDOM); emptyTitleMarkdown != markdown {
		t.Fatalf("empty callout title did not round trip\nexpected\n\t%q\ngot\n\t%q", emptyTitleMarkdown, markdown)
	}
}

func parseCallout(t *testing.T, markdown string) *ast.Node {
	return parseCalloutWithLute(t, newCalloutLute(), markdown)
}

func parseCalloutWithLute(t *testing.T, luteEngine *lute.Lute, markdown string) *ast.Node {
	t.Helper()
	tree := parse.Parse("callout-image-icon", []byte(markdown), luteEngine.ParseOptions)
	if nil == tree.Root.FirstChild || ast.NodeCallout != tree.Root.FirstChild.Type {
		t.Fatalf("callout was not parsed: %q", markdown)
	}
	return tree.Root.FirstChild
}

func newCalloutLute() *lute.Lute {
	luteEngine := lute.New()
	luteEngine.SetCallout(true)
	luteEngine.SetProtyleWYSIWYG(true)
	luteEngine.SetAutoSpace(true)
	return luteEngine
}
