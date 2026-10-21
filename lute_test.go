// Lute - A structured Markdown engine supporting Go and JavaScript.
// Copyright (c) 2019-present, b3log.org
// Licensed under Mulan PSL v2. See LICENSE for details.

package lute

import (
	"runtime"
	"testing"
	"weak"

	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
)

func exportListNode(typ ast.NodeType) (weak.Pointer[ast.Node], error) {
	engine := New()
	tree := parse.Parse("", []byte("- first\n- second\n"), engine.ParseOptions)
	root := weak.Make(tree.Root)
	node := tree.Root
	for node.Type != typ {
		node = node.FirstChild
	}
	_, err := ProtyleExportMdNodeSync(node, engine.ParseOptions, engine.RenderOptions)
	return root, err
}

func TestProtyleExportMdNodeSyncReleasesSourceTree(t *testing.T) {
	for _, typ := range []ast.NodeType{ast.NodeDocument, ast.NodeList, ast.NodeListItem, ast.NodeParagraph} {
		t.Run(typ.String(), func(t *testing.T) {
			root, err := exportListNode(typ)
			if nil != err {
				t.Fatal(err)
			}
			for i := 0; i < 3; i++ {
				runtime.GC()
			}
			if nil != root.Value() {
				t.Fatal("node export retained its source document after garbage collection")
			}
		})
	}
}
