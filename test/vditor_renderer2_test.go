// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

// +build javascript

package test

import (
	"testing"

	"github.com/b3log/lute"
)

type vditor2Test struct {
	*parseTest
	startOffset, endOffset int
}

var vditorRenderer2Tests = []*vditorTest{

	{&parseTest{"1", "<p><strong data-marker=\"**\">foo</strong></p>", "<p><strong data-marker=\"**\">foo</strong></p>"}, 2, 2},
	{&parseTest{"0", "<p>foo</p>", "<p>foo</p>"}, 2, 2},
}

func TestVditorRenderer2(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorRenderer2Tests {
		html, err := luteEngine.RenderVditorDOM2(test.from, test.startOffset, test.endOffset)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
