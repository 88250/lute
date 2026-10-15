package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
)

func TestBlockDOMTableSpans(t *testing.T) {
	cases := []struct {
		name, head, body, expected, imported string
	}{
		{
			name: "middle rowspan",
			head: "<tr><th>A</th><th>B</th><th>C</th></tr>",
			body: `<tr><td>a</td><td rowspan="2">merged</td><td>b</td></tr><tr><td>c</td><td>d</td></tr>`,
			expected: `<tbody><tr><td>a</td><td rowspan="2">merged</td><td>b</td></tr>` +
				`<tr><td>c</td><td class="fn__none"></td><td>d</td></tr></tbody>`,
		},
		{
			name: "staggered rowspans",
			head: "<tr><th>A</th><th>B</th><th>C</th><th>D</th><th>E</th><th>F</th></tr>",
			body: `<tr><td>A1</td><td>B1</td><td>C1</td><td rowspan="4">D</td><td>E1</td><td rowspan="4">F</td></tr>` +
				`<tr><td>A2</td><td>B2</td><td>C2</td><td rowspan="3">E2</td></tr>` +
				`<tr><td>A3</td><td>B3</td><td>C3</td></tr><tr><td>A4</td><td>B4</td><td>C4</td></tr>`,
			expected: `<tbody><tr><td>A1</td><td>B1</td><td>C1</td><td rowspan="4">D</td><td>E1</td><td rowspan="4">F</td></tr>` +
				`<tr><td>A2</td><td>B2</td><td>C2</td><td class="fn__none"></td><td rowspan="3">E2</td><td class="fn__none"></td></tr>` +
				`<tr><td>A3</td><td>B3</td><td>C3</td><td class="fn__none"></td><td class="fn__none"></td><td class="fn__none"></td></tr>` +
				`<tr><td>A4</td><td>B4</td><td>C4</td><td class="fn__none"></td><td class="fn__none"></td><td class="fn__none"></td></tr></tbody>`,
		},
		{
			name: "combined spans with existing placeholders",
			head: "<tr><th>A</th><th>B</th><th>C</th></tr>",
			body: `<tr><td colspan="2" rowspan="2" style="color: red;">merged</td><td class="fn__none"></td><td>b</td></tr>` +
				`<tr><td>c</td></tr>`,
			expected: `<tbody><tr><td style="color: red;" colspan="2" rowspan="2">merged</td><td class="fn__none"></td><td>b</td></tr>` +
				`<tr><td class="fn__none"></td><td class="fn__none"></td><td>c</td></tr></tbody>`,
			imported: `<tbody><tr><td colspan="2" rowspan="2"><span data-type="text" style="color: red;">merged</span></td><td class="fn__none"></td><td>b</td></tr>` +
				`<tr><td class="fn__none"></td><td class="fn__none"></td><td>c</td></tr></tbody>`,
		},
		{
			name: "colspan in header and body",
			head: `<tr><th colspan="2">head</th><th>C</th></tr>`,
			body: `<tr><td>a</td><td colspan="2">wide</td></tr>`,
			expected: `<thead><tr><th colspan="2">head</th><th class="fn__none"></th><th>C</th></tr></thead>` +
				`<tbody><tr><td>a</td><td colspan="2">wide</td><td class="fn__none"></td></tr></tbody>`,
		},
		{
			name: "fully covered row",
			head: "<tr><th>A</th><th>B</th></tr>",
			body: `<tr><td colspan="2" rowspan="2">merged</td></tr><tr></tr>`,
			expected: `<tbody><tr><td colspan="2" rowspan="2">merged</td><td class="fn__none"></td></tr>` +
				`<tr><td class="fn__none"></td><td class="fn__none"></td></tr></tbody>`,
		},
		{
			name:     "separate row groups",
			head:     `<tr><th rowspan="2">head</th><th>B</th></tr>`,
			body:     `<tr><td>a</td><td>b</td></tr>`,
			expected: `<tbody><tr><td>a</td><td>b</td></tr></tbody>`,
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			engine := lute.New()
			engine.SetProtyleWYSIWYG(true)
			engine.SetKramdownIAL(true)
			source := `<div data-type="NodeTable"><div><table><thead>` + test.head + `</thead><tbody>` + test.body + `</tbody></table></div></div>`
			for round := 0; round < 3; round++ {
				source = engine.SpinBlockDOM(source)
				if !strings.Contains(source, test.expected) {
					t.Fatalf("round %d: expected %s in %s", round, test.expected, source)
				}
			}
			// 导出时去掉内部占位，重新导入后仍须保持逻辑列数和内容。
			if "separate row groups" != test.name {
				markdown := engine.BlockDOM2StdMd(source)
				if strings.Contains(markdown, "fn__none") {
					t.Fatalf("placeholder leaked into export: %s", markdown)
				}
				restored := engine.HTML2BlockDOM(markdown)
				expected := test.imported
				if "" == expected {
					expected = test.expected
				}
				if !strings.Contains(restored, expected) {
					t.Fatalf("export round trip changed cells: %s", restored)
				}
			}
		})
	}
}

func TestBlockDOMTableSpansPreserveRichSource(t *testing.T) {
	engine := lute.New()
	rich := &ast.TableCellRich{Spec: 1, Format: "kramdown", Content: "- first\n- second"}
	source := `<div data-type="NodeTable" custom-sy-table-rich="1"><div><table><thead><tr><th>A</th><th>B</th></tr></thead>` +
		`<tbody><tr><td rowspan="2" data-sy-table-cell-rich="` + rich.Encode() + `"></td><td>b</td></tr><tr><td>c</td></tr></tbody></table></div></div>`
	tree := engine.BlockDOM2Tree(source)
	found := false
	ast.Walk(tree.Root, func(node *ast.Node, entering bool) ast.WalkStatus {
		if entering && nil != node.TableCellRich {
			found = true
			if *node.TableCellRich != *rich {
				t.Fatalf("rich source changed: %#v", node.TableCellRich)
			}
		}
		return ast.WalkContinue
	})
	if !found {
		t.Fatal("rich source was lost")
	}
}
