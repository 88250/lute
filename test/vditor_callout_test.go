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
	"github.com/88250/lute/html"
	"github.com/88250/lute/util"
)

func TestVditorCalloutRoundTrip(t *testing.T) {
	markdowns := []string{
		"> [!NOTE]\n> Content",
		"> [!TIP] Custom **title**\n> First\n>\n> Second",
		"> [!CUSTOM] ✂️ Custom title\n> * List\n>   > [!WARNING]\n>   > Nested",
		"> [!IMPORTANT] ![callout-icon](<https://example.com/icon.png>) Image\n> Content",
	}

	for _, markdown := range markdowns {
		luteEngine := newVditorCalloutLute()
		expected := luteEngine.Md2HTML(markdown)

		wysiwygDOM := luteEngine.Md2VditorDOM(markdown)
		if strings.Contains(wysiwygDOM, "not found render function") ||
			!strings.Contains(wysiwygDOM, "data-type=\"callout\"") {
			t.Fatalf("unexpected WYSIWYG DOM for %q:\n%s", markdown, wysiwygDOM)
		}
		wysiwygMarkdown := luteEngine.VditorDOM2Md(wysiwygDOM)
		if actual := luteEngine.Md2HTML(wysiwygMarkdown); expected != actual {
			t.Fatalf("WYSIWYG round trip failed for %q\nexpected:\n%s\nactual:\n%s\nmarkdown:\n%s",
				markdown, expected, actual, wysiwygMarkdown)
		}

		irDOM := luteEngine.Md2VditorIRDOM(markdown)
		if strings.Contains(irDOM, "not found render function") || !strings.Contains(irDOM, "data-type=\"callout\"") {
			t.Fatalf("unexpected IR DOM for %q:\n%s", markdown, irDOM)
		}
		irMarkdown := luteEngine.VditorIRDOM2Md(irDOM)
		if actual := luteEngine.Md2HTML(irMarkdown); expected != actual {
			t.Fatalf("IR round trip failed for %q\nexpected:\n%s\nactual:\n%s\nmarkdown:\n%s",
				markdown, expected, actual, irMarkdown)
		}

		svDOM := luteEngine.Md2VditorSVDOM(markdown)
		if strings.Contains(svDOM, "not found render function") || !strings.Contains(svDOM, "[!") {
			t.Fatalf("unexpected SV DOM for %q:\n%s", markdown, svDOM)
		}
		svMarkdown := browserTextContent(util.ParseHTML(svDOM))
		if actual := luteEngine.Md2HTML(svMarkdown); expected != actual {
			t.Fatalf("SV round trip failed for %q\nexpected:\n%s\nactual:\n%s\nDOM:\n%s\nmarkdown:\n%q",
				markdown, expected, actual, svDOM, svMarkdown)
		}
	}
}

func browserTextContent(node *html.Node) (ret string) {
	if nil == node {
		return
	}
	if html.TextNode == node.Type {
		return node.Data
	}
	for child := node.FirstChild; nil != child; child = child.NextSibling {
		ret += browserTextContent(child)
	}
	return
}

func TestCalloutTypeAttributeEscaping(t *testing.T) {
	luteEngine := newVditorCalloutLute()
	markdown := "> [!\" onclick=\"alert(1)] Title\n> Content"
	blockDOM := luteEngine.Md2BlockDOM(markdown, false)
	for name, output := range map[string]string{
		"HTML":            luteEngine.Md2HTML(markdown),
		"WYSIWYG":         luteEngine.Md2VditorDOM(markdown),
		"IR":              luteEngine.Md2VditorIRDOM(markdown),
		"BlockDOM":        blockDOM,
		"Protyle preview": luteEngine.ProtylePreviewStr("", markdown),
		"Protyle export":  luteEngine.BlockDOM2HTML(blockDOM),
	} {
		if strings.Contains(output, " onclick=\"") {
			t.Fatalf("%s contains an injected event attribute:\n%s", name, output)
		}
	}
}

func newVditorCalloutLute() *lute.Lute {
	luteEngine := lute.New()
	luteEngine.SetCallout(true)
	luteEngine.SetSanitize(true)
	return luteEngine
}
