package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
)

func TestSpinCrossTextMarkDelimiters(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetTextMark(true)
	engine.SetKramdownIAL(true)
	engine.SetSpin(true)
	engine.SetMark(true)
	wrap := func(content string) string {
		return `<div data-node-id="20260908150000-abcdefg" data-type="NodeParagraph" class="p"><div contenteditable="true">` + content + `</div><div class="protyle-attr" contenteditable="false"></div></div>`
	}
	for _, tc := range []struct{ name, input, want string }{
		{"close-inside", `*<span data-type="strong">11*<wbr></span>`, `<span data-type="em strong">11</span>`},
		{"open-inside", `<span data-type="strong">*11</span>*<wbr>`, `<span data-type="em strong">11</span>`},
		{"different-styles", `<span data-type="strong">*foo</span><span data-type="mark">bar*</span><wbr>`, `<span data-type="em strong">foo</span><span data-type="em mark">bar</span>`},
		{"mixed-text", `*foo<span data-type="strong">bar*</span><wbr>`, `<span data-type="em">foo</span><span data-type="em strong">bar</span>`},
		{"attributes", `*<span data-type="strong text" style="color: red;">11*</span><wbr>`, `style="color: red;"`},
		{"unmatched", `*<span data-type="strong">11<wbr></span>`, `*<span data-type="strong">11<wbr></span>`},
		{"code-boundary", `*<span data-type="strong code">11*</span><wbr>`, `11*</span>`},
		{"escaped", `\*<span data-type="strong">11*</span><wbr>`, `<span data-type="strong">11*</span>`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := engine.SpinBlockDOM(wrap(tc.input))
			if !strings.Contains(got, tc.want) {
				t.Fatalf("expected %q in %q", tc.want, got)
			}
			if !strings.Contains(got, "<wbr>") {
				t.Fatalf("caret lost: %q", got)
			}
			// 空样式节点中的光标可在下一轮规范化时移到节点外，正文格式必须保持稳定。
			again := engine.SpinBlockDOM(got)
			if !strings.Contains(again, tc.want) {
				t.Fatalf("format changed: %q / %q", got, again)
			}
			if stable := engine.SpinBlockDOM(again); stable != again {
				t.Fatalf("not stable: %q / %q", again, stable)
			}
		})
	}
	engine.SetInlineAsterisk(false)
	got := engine.SpinBlockDOM(wrap(`*<span data-type="strong">11*</span><wbr>`))
	if strings.Contains(got, `data-type="em`) {
		t.Fatalf("disabled emphasis parsed: %q", got)
	}
}
