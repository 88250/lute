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

package test

import (
	"testing"
)

var vditorSpinTests = []parseTest{

	{"1", "<span data-ntype=\"10\" data-mtype=\"2\">*</span><span class=\"node\" data-ntype=\"11\" data-mtype=\"2\"><span class=\"marker\">*</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></em><span class=\"marker\">*</span></span>",
		           "<span data-ntype=\"10\" data-mtype=\"2\">*</span><span class=\"node\" data-ntype=\"11\" data-mtype=\"2\"><span class=\"marker\">*</span><em data-ntype=\"11\" data-mtype=\"2\"><span data-ntype=\"10\" data-mtype=\"2\">foo</span></em><span class=\"marker\">*</span></span>"},
	{"0", "<span data-ntype=\"10\" data-mtype=\"2\">Lute</span>", "<span data-ntype=\"10\" data-mtype=\"2\">Lute</span>"},
}

func TestVditorSpin(t *testing.T) {
	//luteEngine := lute.New()
	//
	//for _, test := range vditorSpinTests {
	//	html, err := luteEngine.SpinVditorDOM(test.from)
	//	if nil != err {
	//		t.Fatalf("unexpected: %s", err)
	//	}
	//
	//	if test.to != html {
	//		t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
	//	}
	//}
}
