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
)

func TestVditorTableCellBrRoundTrip(t *testing.T) {
	luteEngine := lute.New()
	tests := []struct {
		name            string
		markdown        string
		markdownContent string
		domBrCount      int
		markdownBrCount int
	}{
		{"middle", "| c |\n| - |\n| 1<br />2 |\n", "1<br />2", 1, 1},
		{"compact", "| c |\n| - |\n| 1<br>2<br/>3 |\n", "1<br />2<br />3", 2, 2},
		{"uppercase", "| c |\n| - |\n| 1<BR />2 |\n", "1<br />2", 1, 1},
		{"trailing", "| c |\n| - |\n| 1<br /> |\n", "1<br />", 2, 1},
		{"consecutive", "| c |\n| - |\n| 1<br /><br /> |\n", "1<br /><br />", 3, 2},
	}

	for _, test := range tests {
		dom := luteEngine.Md2VditorIRDOM(test.markdown)
		if got := strings.Count(dom, "<br />"); test.domBrCount != got {
			t.Fatalf("test case [%s] DOM failed\nexpected br count\n\t%d\ngot\n\t%d\nDOM\n\t%q",
				test.name, test.domBrCount, got, dom)
		}
		if strings.Contains(dom, "data-type=\"html-inline\"") {
			t.Fatalf("test case [%s] should render table br directly\nDOM\n\t%q", test.name, dom)
		}

		markdown := luteEngine.VditorIRDOM2Md(dom)
		if got := strings.Count(markdown, "<br />"); test.markdownBrCount != got {
			t.Fatalf("test case [%s] Markdown failed\nexpected br count\n\t%d\ngot\n\t%d\nDOM\n\t%q\nMarkdown\n\t%q",
				test.name, test.markdownBrCount, got, dom, markdown)
		}
		if !strings.Contains(markdown, test.markdownContent) {
			t.Fatalf("test case [%s] Markdown failed\nexpected content\n\t%q\ngot\n\t%q",
				test.name, test.markdownContent, markdown)
		}
	}

	dom := luteEngine.Md2VditorIRDOM("1<br />2")
	if !strings.Contains(dom, "data-type=\"html-inline\"") {
		t.Fatalf("br outside a table should remain inline HTML\nDOM\n\t%q", dom)
	}
}

func TestSpinVditorTableCellBr(t *testing.T) {
	luteEngine := lute.New()
	tests := []struct {
		name            string
		cell            string
		domContent      string
		domBrCount      int
		markdownBrCount int
	}{
		{"trailing", "1<br><wbr><br>", "1<br /><wbr><br />", 2, 1},
		{"consecutive", "1<br><br><wbr><br>", "1<br /><br /><wbr><br />", 3, 2},
		{"deleted", "1<wbr><br>", "<td>1<wbr></td>", 0, 0},
	}

	for _, test := range tests {
		input := "<table data-block=\"0\"><thead><tr><th>c</th></tr></thead><tbody><tr><td>" + test.cell +
			"</td></tr></tbody></table>"
		dom := luteEngine.SpinVditorIRDOM(input)
		if !strings.Contains(dom, test.domContent) {
			t.Fatalf("test case [%s] DOM failed\nexpected content\n\t%q\ngot\n\t%q", test.name, test.domContent, dom)
		}
		if got := strings.Count(dom, "<br />"); test.domBrCount != got {
			t.Fatalf("test case [%s] DOM failed\nexpected br count\n\t%d\ngot\n\t%d\nDOM\n\t%q",
				test.name, test.domBrCount, got, dom)
		}
		if strings.Contains(dom, "data-type=\"html-inline\"") {
			t.Fatalf("test case [%s] should render table br directly\nDOM\n\t%q", test.name, dom)
		}

		markdown := luteEngine.VditorIRDOM2Md(dom)
		if got := strings.Count(markdown, "<br />"); test.markdownBrCount != got {
			t.Fatalf("test case [%s] Markdown failed\nexpected br count\n\t%d\ngot\n\t%d\nMarkdown\n\t%q",
				test.name, test.markdownBrCount, got, markdown)
		}

		respun := luteEngine.SpinVditorIRDOM(dom)
		if dom != respun {
			t.Fatalf("test case [%s] should be stable after another spin\nexpected\n\t%q\ngot\n\t%q", test.name, dom, respun)
		}
	}
}
