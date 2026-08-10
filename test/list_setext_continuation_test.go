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
	"github.com/88250/lute/util"
)

func TestFormatListSetextContinuation(t *testing.T) {
	luteEngine := lute.New()
	tests := []struct {
		name      string
		markdown  string
		formatted string
	}{
		{"unordered-equal", "* a\n* b\n=", "* a\n* b\n  \\=\n"},
		{"unordered-hyphens", "* a\n* b\n--", "* a\n* b\n  \\--\n"},
		{"ordered-equals", "1. a\n2. b\n===", "1. a\n2. b\n   \\===\n"},
		{"nested-equal", "* a\n  * b\n=", "* a\n  * b\n    \\=\n"},
		{"leading-space", "* a\n* b\n =", "* a\n* b\n  \\=\n"},
	}

	for _, test := range tests {
		formatted := luteEngine.FormatStr(test.name, test.markdown)
		if test.formatted != formatted {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q", test.name, test.formatted, formatted)
		}
		if formattedAgain := luteEngine.FormatStr(test.name, formatted); formatted != formattedAgain {
			t.Fatalf("test case [%s] is not idempotent\nfirst\n\t%q\nsecond\n\t%q", test.name, formatted, formattedAgain)
		}
		before := luteEngine.MarkdownStr(test.name, test.markdown)
		after := luteEngine.MarkdownStr(test.name, formatted)
		if before != after {
			t.Fatalf("test case [%s] changed semantics\nbefore\n\t%q\nafter\n\t%q", test.name, before, after)
		}
	}

	heading := "* a\n* b\n  ="
	if formatted := luteEngine.FormatStr("heading", heading); "* a\n* b\n  =\n" != formatted {
		t.Fatalf("list Setext heading should remain a heading, got %q", formatted)
	}
}

func TestVditorDOMListSetextContinuation(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetVditorWYSIWYG(true)

	dom := `<ul data-tight="true" data-marker="*" data-block="0"><li data-marker="*">a</li><li data-marker="*">b
-<wbr></li></ul>`
	markdown := luteEngine.VditorDOM2Md(dom)
	if expected := "* a\n* b\n  \\-\n"; expected != markdown {
		t.Fatalf("expected %q, got %q", expected, markdown)
	}
	spun := luteEngine.SpinVditorDOM(dom)
	if strings.Contains(spun, "<h2") || !strings.Contains(spun, `data-type="backslash"`) {
		t.Fatalf("list continuation should remain literal, got %q", spun)
	}
}

func TestSpinVditorSVListSetextContinuation(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetVditorSV(true)

	for _, marker := range []string{"-", "--", "=", "==="} {
		markdown := "* a\n* b\n" + marker + "‸"
		spun := luteEngine.SpinVditorSVDOM(markdown)
		if strings.Contains(spun, "heading-marker") || !strings.Contains(spun, "\\"+marker) {
			t.Fatalf("marker %q should remain literal, got %q", marker, spun)
		}

		root := util.ParseHTML(spun)
		spunAgain := luteEngine.SpinVditorSVDOM(util.DomText(root))
		if strings.Contains(spunAgain, "heading-marker") {
			t.Fatalf("marker %q became a heading after a second spin: %q", marker, spunAgain)
		}
	}
}

func TestSpinVditorSVThematicBreakAtDocumentStart(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetVditorSV(true)

	spun := luteEngine.SpinVditorSVDOM("***\n‸")
	if strings.Contains(spun, "yaml-front-matter") || !strings.Contains(spun, ">***</span>") {
		t.Fatalf("thematic break should not become YAML front matter, got %q", spun)
	}

	root := util.ParseHTML(spun)
	spunAgain := luteEngine.SpinVditorSVDOM(util.DomText(root))
	if strings.Contains(spunAgain, "yaml-front-matter") || !strings.Contains(spunAgain, ">***</span>") {
		t.Fatalf("thematic break became YAML front matter after a second spin: %q", spunAgain)
	}
}
