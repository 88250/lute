package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestImageCustomAttributesRoundTrip(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetKramdownIAL(true)
	engine.SetSanitize(true)
	md := `![a](assets/photo-a.heic){: custom-sy-live-photo="1" custom-sy-live-photo-video="assets/video-a.mov"} ![b](assets/photo-b.heic){: custom-sy-live-photo="1" custom-sy-live-photo-video="assets/video-b.mov"}`
	want := map[string]string{"assets/photo-a.heic": "assets/video-a.mov", "assets/photo-b.heic": "assets/video-b.mov"}
	check := func(dom string) {
		t.Helper()
		count := 0
		ast.Walk(engine.BlockDOM2Tree(dom).Root, func(n *ast.Node, entering bool) ast.WalkStatus {
			if entering && n.Type == ast.NodeImage {
				count++
				dest := n.ChildByType(ast.NodeLinkDest).TokensStr()
				if n.IALAttr("custom-sy-live-photo-video") != want[dest] || n.IALAttr("custom-sy-live-photo") != "1" {
					t.Errorf("image %q lost its own attributes: %v", dest, n.KramdownIAL)
				}
			}
			return ast.WalkContinue
		})
		if count != 2 {
			t.Fatalf("got %d images", count)
		}
	}
	dom := engine.Md2BlockDOM(md, false)
	for i := 0; i < 3; i++ {
		check(dom)
		dom = engine.SpinBlockDOM(dom)
		check(engine.Md2BlockDOM(engine.BlockDOM2Md(dom), false))
	}
	if html := engine.BlockDOM2HTML(dom); !strings.Contains(html, `custom-sy-live-photo-video="assets/video-a.mov"`) {
		t.Fatalf("HTML lost custom image attributes: %s", html)
	}
	tree := parse.Parse("", []byte(md), engine.ParseOptions)
	preview := render.NewProtylePreviewRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()
	if !strings.Contains(string(preview), `custom-sy-live-photo-video="assets/video-b.mov"`) {
		t.Fatalf("preview lost custom image attributes: %s", preview)
	}
}

func TestImageCustomAttributeValuesRemainEscaped(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetKramdownIAL(true)
	engine.SetSanitize(true)
	const value = `A & B "quoted" <literal>`
	tree := parse.Parse("", []byte("![a](assets/photo.heic)"), engine.ParseOptions)
	image := tree.Root.FirstChild.ChildByType(ast.NodeImage)
	image.SetIALAttr("custom-caption", value)
	image.SetIALAttr(`custom-invalid" onerror`, "unexpected")
	for i := 0; i < 3; i++ {
		dom := engine.Tree2BlockDOM(tree, engine.RenderOptions, engine.ParseOptions)
		if strings.Contains(dom, " onerror=") || strings.Contains(dom, "custom-invalid") {
			t.Fatalf("unsafe attribute name was rendered: %s", dom)
		}
		tree = engine.BlockDOM2Tree(dom)
		image = tree.Root.FirstChild.ChildByType(ast.NodeImage)
		if image == nil || image.IALAttr("custom-caption") != value {
			t.Fatalf("round %d did not preserve an escaped attribute value", i)
		}
	}
}
