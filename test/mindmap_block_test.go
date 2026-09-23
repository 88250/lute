package test

import (
	"strings"
	"testing"

	"github.com/88250/lute/ast"
)

func TestMindmapBlockDOMRoundTrip(t *testing.T) {
	l := tabsEngine()
	const dom = `<div data-node-id="20260923120000-map0001" data-type="NodeMindmap" data-subtype="u" class="mindmap"><div data-node-id="20260923120000-item001" data-type="NodeMindmapItem" data-subtype="u" data-marker="*" class="mindmap-item"><div class="protyle-action"></div><div data-node-id="20260923120000-para001" data-type="NodeParagraph" class="p"><div contenteditable="true">Root</div></div><div data-node-id="20260923120000-map0002" data-type="NodeMindmap" data-subtype="u" class="mindmap"><div data-node-id="20260923120000-item002" data-type="NodeMindmapItem" data-subtype="u" data-marker="*" class="mindmap-item"><div class="protyle-action"></div><div data-node-id="20260923120000-para002" data-type="NodeParagraph" class="p"><div contenteditable="true">Child</div></div></div></div></div></div>`
	source := dom
	for round := 0; round < 3; round++ {
		tree := l.BlockDOM2Tree(source)
		mapNode := tree.Root.FirstChild
		if mapNode.Type != ast.NodeMindmap || mapNode.FirstChild.Type != ast.NodeMindmapItem ||
			mapNode.FirstChild.ChildByType(ast.NodeMindmap).FirstChild.Type != ast.NodeMindmapItem {
			t.Fatalf("round %d: mind map types were lost", round)
		}
		rendered := l.RenderNodeBlockDOM(mapNode)
		if !strings.Contains(rendered, `data-type="NodeMindmap"`) ||
			!strings.Contains(rendered, `data-type="NodeMindmapItem"`) || !strings.Contains(rendered, `class="mindmap"`) {
			t.Fatalf("round %d: %s", round, rendered)
		}
		if mapNode.ID != "20260923120000-map0001" || mapNode.FirstChild.ID != "20260923120000-item001" {
			t.Fatalf("round %d: IDs changed", round)
		}
		source = rendered
	}
	spun := l.SpinBlockDOM(source)
	if strings.Count(spun, `data-type="NodeMindmap"`) != 2 ||
		strings.Count(spun, `data-type="NodeMindmapItem"`) != 2 {
		t.Fatalf("spin lost mind map types: markdown=%q DOM=%s", l.BlockDOM2Md(source), spun)
	}
	standard := l.BlockDOM2StdMd(source)
	if !strings.Contains(standard, "Root") || !strings.Contains(standard, "Child") {
		t.Fatalf("standard Markdown lost mind map content: %s", standard)
	}
	html := l.BlockDOM2HTML(source)
	if !strings.Contains(html, "Root") || !strings.Contains(html, "Child") {
		t.Fatalf("HTML lost mind map content: %s", html)
	}
}

func TestMindmapContainment(t *testing.T) {
	container := &ast.Node{Type: ast.NodeMindmap}
	item := &ast.Node{Type: ast.NodeMindmapItem}
	if !container.CanContain(ast.NodeMindmapItem) || container.CanContain(ast.NodeListItem) ||
		!item.CanContain(ast.NodeMindmap) || !item.CanContain(ast.NodeParagraph) ||
		item.CanContain(ast.NodeList) || item.CanContain(ast.NodeListItem) || item.CanContain(ast.NodeMindmapItem) {
		t.Fatal("invalid mind map content boundaries")
	}
}
