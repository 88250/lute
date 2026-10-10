package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
)

func TestImageTitleDataAttribute(t *testing.T) {
	for _, fixture := range []struct {
		name  string
		attrs string
		title string
	}{
		{"legacy", `title="caption &amp; &lt;literal&gt;"`, "caption & <literal>"},
		{"data", `data-title="caption &amp; &lt;literal&gt;"`, "caption & <literal>"},
		{"preferred", `title="legacy" data-title="updated"`, "updated"},
		{"cleared", `title="legacy" data-title=""`, ""},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			engine := lute.New()
			engine.SetProtyleWYSIWYG(true)
			dom := `<div data-node-id="20261010130000-abcdefg" data-type="NodeParagraph"><div contenteditable="true"><span data-type="img" class="img"><span> </span><span><img src="assets/image.png" data-src="assets/image.png" alt="description" ` + fixture.attrs + `><span class="protyle-action__title"><span>caption</span></span></span><span> </span></span></div></div>`
			checkTitle := func(input string) {
				t.Helper()
				tree := engine.BlockDOM2Tree(input)
				image := tree.Root.FirstChild.ChildByType(ast.NodeImage)
				if nil == image {
					t.Fatal("image is missing")
				}
				var title string
				if node := image.ChildByType(ast.NodeLinkTitle); nil != node {
					title = string(node.Tokens)
				}
				if title != fixture.title {
					t.Fatalf("title: got %q, want %q", title, fixture.title)
				}
			}
			checkTitle(dom)
			checkTitle(engine.SpinBlockDOM(dom))
			checkTitle(engine.Md2BlockDOM(engine.BlockDOM2Md(dom), false))
			checkTitle(engine.Md2BlockDOM(engine.BlockDOM2StdMd(dom), false))
			if fixture.title != "" && !strings.Contains(engine.BlockDOM2HTML(dom), "title=\"") {
				t.Fatal("HTML export lost the image title")
			}
		})
	}
}
