package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/editor"
	"github.com/88250/lute/render"
)

func TestProtyleMindmapIsOrdinaryCode(t *testing.T) {
	engine := lute.New()
	engine.SetProtyleWYSIWYG(true)
	engine.SetKramdownIAL(true)
	for _, source := range []string{"", "1\n", "- Root\n  - Child\n", "\n  &lt;script&gt; &amp; `code`  \n\n"} {
		markdown := "```mindmap\n" + source + "```\n"
		dom, tree := engine.Md2BlockDOMTree(markdown, false)
		for _, output := range []string{dom, engine.SpinBlockDOM(dom), engine.RenderNodeBlockDOM(tree.Root.FirstChild)} {
			if !strings.Contains(output, `class="code-block"`) || !strings.Contains(output, `>mindmap</span>`) ||
				strings.Contains(output, `class="render-node"`) || strings.Contains(output, `data-subtype="mindmap"`) {
				t.Fatalf("mindmap did not follow ordinary code rendering: %s", output)
			}
			parsed := engine.BlockDOM2Tree(output).Root.FirstChild
			if parsed.Type != ast.NodeCodeBlock ||
				string(parsed.ChildByType(ast.NodeCodeBlockCode).Tokens) != source ||
				string(parsed.ChildByType(ast.NodeCodeBlockFenceInfoMarker).CodeBlockInfo) != "mindmap" {
				t.Fatalf("code content or language changed for %q", source)
			}
		}
		for name, output := range map[string]string{
			"preview": string(render.NewProtylePreviewRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()),
			"export":  string(render.NewProtyleExportRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()),
			"docx":    string(render.NewProtyleExportDocxRenderer(tree, engine.RenderOptions, engine.ParseOptions).Render()),
		} {
			if strings.Contains(output, `data-subtype="mindmap"`) || !strings.Contains(output, `class="code-block`) {
				t.Fatalf("%s still renders mindmap as a diagram: %s", name, output)
			}
		}
		md := engine.BlockDOM2StdMd(dom)
		if !strings.Contains(md, source) || !strings.Contains(md, "```mindmap") {
			t.Fatalf("Markdown export lost literal code: %q", md)
		}
		for _, output := range []string{engine.BlockDOM2HTML(dom), engine.BlockDOM2RichHTML(dom)} {
			if strings.Contains(output, "data-code=") || !strings.Contains(output, "<pre") {
				t.Fatalf("HTML export must contain ordinary code: %s", output)
			}
		}
	}
	caret := engine.Md2BlockDOM("```mindmap\n1"+editor.Caret+"2\n```", false)
	if output := engine.SpinBlockDOM(caret); !strings.Contains(output, "<wbr>") || strings.Contains(output, `class="render-node"`) {
		t.Fatalf("editing must retain the caret inside a code block: %s", output)
	}
}
