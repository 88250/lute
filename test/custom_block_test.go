// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"strings"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	lutehtml "github.com/88250/lute/html"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestCustomBlockOptionAndFence(t *testing.T) {
	markdown := ";;;plugin:test\ncontent\n;;;\nafter"
	luteEngine := lute.New()
	tree := parse.Parse("", []byte(markdown), luteEngine.ParseOptions)
	if ast.NodeCustomBlock == tree.Root.FirstChild.Type {
		t.Fatal("custom block parsing should be disabled by default")
	}

	luteEngine.SetCustomBlock(true)
	tree = parse.Parse("", []byte(markdown), luteEngine.ParseOptions)
	customBlock := tree.Root.FirstChild
	if ast.NodeCustomBlock != customBlock.Type || "plugin:test" != customBlock.CustomBlockInfo || "content\n" != customBlock.TokensStr() {
		t.Fatalf("unexpected custom block: type=%s, info=%q, content=%q", customBlock.Type, customBlock.CustomBlockInfo, customBlock.TokensStr())
	}
	if nil == customBlock.Next || ast.NodeParagraph != customBlock.Next.Type {
		t.Fatal("an exact three-semicolon closing fence should close the custom block")
	}
	for name, input := range map[string]string{
		"two-semicolon opening fence":  ";;plugin:test\ncontent\n;;;",
		"four-semicolon opening fence": ";;;;plugin:test\ncontent\n;;;",
		"semicolon in info":            ";;;plugin;test\ncontent\n;;;",
	} {
		t.Run(name, func(t *testing.T) {
			tree := parse.Parse("", []byte(input), luteEngine.ParseOptions)
			if ast.NodeCustomBlock == tree.Root.FirstChild.Type {
				t.Fatalf("invalid opening fence parsed as a custom block: %q", input)
			}
		})
	}

	tree = parse.Parse("", []byte(";;;plugin:test\ncontent\n;;;;\nafter"), luteEngine.ParseOptions)
	customBlock = tree.Root.FirstChild
	if ast.NodeCustomBlock != customBlock.Type || nil != customBlock.Next || !strings.Contains(customBlock.TokensStr(), ";;;;\nafter") {
		t.Fatalf("a four-semicolon line should remain custom block content: type=%s, content=%q", customBlock.Type, customBlock.TokensStr())
	}
}

func TestCustomBlockSafeFallbackRenderers(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetCustomBlock(true)
	info := "plugin:test\"\tonmouseover=\"alert(1)"
	content := "<script>alert(\"x\")</script>&amp;\n"
	markdown := ";;;" + info + "\n" + content + ";;;"

	outputs := map[string]string{
		"Protyle":         luteEngine.Md2BlockDOM(markdown, false),
		"HTML":            luteEngine.MarkdownStr("", markdown),
		"Protyle preview": luteEngine.ProtylePreviewStr("", markdown),
	}
	tree := parse.Parse("", []byte(markdown), luteEngine.ParseOptions)
	outputs["Protyle export"] = string(render.NewProtyleExportRenderer(tree, luteEngine.RenderOptions, luteEngine.ParseOptions).Render())
	tree = parse.Parse("", []byte(markdown), luteEngine.ParseOptions)
	outputs["Docx export"] = string(render.NewProtyleExportDocxRenderer(tree, luteEngine.RenderOptions, luteEngine.ParseOptions).Render())

	escapedInfo := lutehtml.EscapeHTMLStr(info)
	escapedContent := lutehtml.EscapeHTMLStr(content)
	for name, output := range outputs {
		if strings.Contains(output, "onmouseover=\"") {
			t.Fatalf("%s contains an injected event attribute:\n%s", name, output)
		}
		if !strings.Contains(output, `data-info="`+escapedInfo+`"`) {
			t.Fatalf("%s did not escape custom block info:\n%s", name, output)
		}
		if strings.Contains(output, "<script>") || !strings.Contains(output, "<pre>"+escapedContent+"</pre>") {
			t.Fatalf("%s did not render a safe visible fallback:\n%s", name, output)
		}
	}

	protyle := outputs["Protyle"]
	rootTag := protyle[:strings.IndexByte(protyle, '>')]
	if 1 != strings.Count(rootTag, `data-type="NodeCustomBlock"`) {
		t.Fatalf("Protyle custom block should contain one data-type attribute:\n%s", protyle)
	}
	if !strings.Contains(rootTag, `contenteditable="false"`) || !strings.Contains(protyle, `class="custom-block__content"`) {
		t.Fatalf("Protyle custom block is missing its non-editable root or content mount:\n%s", protyle)
	}
}

func TestCustomBlockEntityRoundTripAndRenderedSubtree(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetCustomBlock(true)
	content := "literal &amp; &quot; &#123; & <tag> \"quoted\"\n"
	markdown := ";;;plugin:entity\n" + content + ";;;"
	blockDOM := luteEngine.Md2BlockDOM(markdown, false)
	if !strings.Contains(blockDOM, `data-content="literal &amp;amp; &amp;quot; &amp;#123; &amp; &lt;tag&gt; &quot;quoted&quot;`+"\n\"") {
		t.Fatalf("custom block data-content was not escaped exactly once:\n%s", blockDOM)
	}

	assertCustomBlockContent(t, luteEngine, luteEngine.BlockDOM2Md(blockDOM), content)
	assertCustomBlockDOMContent(t, luteEngine, luteEngine.SpinBlockDOM(blockDOM), content)

	contentStart := strings.Index(blockDOM, `<div class="custom-block__content">`)
	attrStart := strings.Index(blockDOM, `<div class="protyle-attr"`)
	if 0 > contentStart || contentStart >= attrStart {
		t.Fatalf("unexpected custom block DOM:\n%s", blockDOM)
	}
	pluginRenderTree := `<div class="custom-block__content"><div data-node-id="rendered-child" data-type="NodeParagraph" class="p"><div contenteditable="true">PLUGIN_RENDER_TREE</div></div></div>`
	blockDOM = blockDOM[:contentStart] + pluginRenderTree + blockDOM[attrStart:]
	markdown = luteEngine.BlockDOM2Md(blockDOM)
	if strings.Contains(markdown, "PLUGIN_RENDER_TREE") {
		t.Fatalf("plugin render subtree leaked into Markdown:\n%s", markdown)
	}
	assertCustomBlockContent(t, luteEngine, markdown, content)
	spun := luteEngine.SpinBlockDOM(blockDOM)
	if strings.Contains(spun, "PLUGIN_RENDER_TREE") {
		t.Fatalf("plugin render subtree leaked through SpinBlockDOM:\n%s", spun)
	}
	assertCustomBlockDOMContent(t, luteEngine, spun, content)
}

func assertCustomBlockContent(t *testing.T, luteEngine *lute.Lute, markdown, expected string) {
	t.Helper()
	tree := parse.Parse("", []byte(markdown), luteEngine.ParseOptions)
	node := tree.Root.FirstChild
	if nil == node {
		t.Fatalf("custom block missing after round trip: markdown=%q", markdown)
	}
	if ast.NodeCustomBlock != node.Type || expected != node.TokensStr() {
		t.Fatalf("unexpected custom block round trip: type=%v, content=%q, markdown=%q", node.Type, node.TokensStr(), markdown)
	}
}

func assertCustomBlockDOMContent(t *testing.T, luteEngine *lute.Lute, blockDOM, expected string) {
	t.Helper()
	tree := luteEngine.BlockDOM2Tree(blockDOM)
	node := tree.Root.FirstChild
	if nil == node {
		t.Fatalf("custom block missing after DOM round trip: DOM=%q", blockDOM)
	}
	if ast.NodeCustomBlock != node.Type || expected != node.TokensStr() {
		t.Fatalf("unexpected custom block DOM round trip: type=%v, content=%q, DOM=%q", node.Type, node.TokensStr(), blockDOM)
	}
}
