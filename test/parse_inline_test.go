// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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
	"testing"

	"github.com/88250/lute"
	"github.com/88250/lute/ast"
	"github.com/88250/lute/parse"
	"github.com/88250/lute/util"
)

var inlineTests = []parseTest{

	{"2", "<form enctype=", "<form enctype="},
	{"1", "**foo** [foo](bar)\n", "foo foo"},
	{"0", "1. foo\n", "1. foo"},
}

func TestInline(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range inlineTests {
		tree := parse.Inline("", util.StrToBytes(test.from), luteEngine.ParseOptions)
		if ast.NodeParagraph != tree.Root.FirstChild.Type {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, ast.NodeParagraph, tree.Root.FirstChild.Type, test.from)
		}

		if output := tree.Root.Text(); test.to != output {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, output, test.from)
		}
	}
}
