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
	"testing"

	"github.com/88250/lute"
)


var blockDOM2StdMd = []parseTest{

	{"2", "<span data-type=\"tag\">foo</span> bar <em>foo</em> bar", "#foo# bar *foo* bar\n"},
	{"1", "foo <code>bar</code> baz", "foo `bar` baz\n"},
	{"0", "foo<u>bar</u>baz~~abc~~xyz", "foo<u>bar</u>baz~~abc~~xyz\n"},
}

func TestBlockDOM2StdMd(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range blockDOM2StdMd {
		md := luteEngine.BlockDOM2StdMd(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}
