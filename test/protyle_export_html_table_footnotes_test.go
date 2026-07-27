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
	"bytes"
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/render"
)

func TestProtyleExportHTMLTableInFootnotes(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetFootnotes(true)
	luteEngine.SetKramdownIAL(true)
	luteEngine.SetProtyleWYSIWYG(true)

	tree := parse.Parse("", []byte("foo[^1]\n\n[^1]: bar\n\n[^2]: baz\n"), luteEngine.ParseOptions)
	defBlock := tree.Root.ChildByType(ast.NodeFootnotesDefBlock)
	if nil == defBlock || nil == defBlock.FirstChild || nil == defBlock.FirstChild.Next {
		t.Fatal("expected two footnotes definitions")
	}
	firstDef, secondDef := defBlock.FirstChild, defBlock.FirstChild.Next

	table := &ast.Node{Type: ast.NodeTable}
	table.SetIALAttr("id", "20260724130222-z0ruwey")
	head := &ast.Node{Type: ast.NodeTableHead}
	row := &ast.Node{Type: ast.NodeTableRow}
	cell := &ast.Node{Type: ast.NodeTableCell}
	cell.SetIALAttr("colspan", "2")
	cell.AppendChild(&ast.Node{Type: ast.NodeText, Tokens: []byte("merged ")})
	cell.AppendChild(&ast.Node{
		Type:              ast.NodeFootnotesRef,
		Tokens:            []byte("^2"),
		FootnotesRefId:    "2",
		FootnotesRefLabel: []byte("^2"),
	})
	row.AppendChild(cell)
	head.AppendChild(row)
	table.AppendChild(head)
	firstDef.AppendChild(table)

	renderer := render.NewProtyleExportRenderer(tree, luteEngine.RenderOptions, luteEngine.ParseOptions)
	output := renderer.Render()
	expectedRef := []byte("<sup class=\"footnotes-ref\" id=\"footnotes-ref-2\"><a href=\"#footnotes-def-2\">2</a></sup>")
	if !bytes.Contains(output, expectedRef) {
		t.Fatalf("expected rendered footnotes reference %q in %q", expectedRef, output)
	}
	if defBlock != tree.Root.ChildByType(ast.NodeFootnotesDefBlock) || tree.Root != defBlock.Parent {
		t.Fatal("footnotes definition block should remain on the document root")
	}
	if firstDef != defBlock.FirstChild || secondDef != firstDef.Next || secondDef != defBlock.LastChild || nil != secondDef.Next {
		t.Fatal("footnotes definitions should retain their order")
	}
	if firstDef != table.Parent || nil != table.ChildByType(ast.NodeFootnotesDefBlock) {
		t.Fatal("table should not retain a temporary footnotes definition block")
	}
}
