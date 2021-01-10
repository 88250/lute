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
	"github.com/88250/lute/ast"
	"testing"

	"github.com/88250/lute"
)

var kitymindJSONRendererTests = []parseTest{

	{"0", "# foo\n{: id=\"20210110005758-m303ovi\"}\n\nbar\n{: id=\"20210110115402-21ltd5v\"}\n\nbaz **bazz**\n{: id=\"20210110115405-17ng22v\"}\n\n\n{: id=\"20201228004131-bys3g5x\" type=\"doc\"}\n", "{\"root\":{\"data\":{\"text\":\"文档名 TODO\",\"id\":\"20201228004131-bys3g5x\"},\"children\":[{\"data\":{\"text\":\"# foo\",\"id\":\"20210110005758-m303ovi\"},\"children\":[{\"data\":{\"text\":\"baz **bazz**\",\"id\":\"20210110115405-17ng22v\"},\"children\":[]}]},{\"data\":{\"text\":\"bar\",\"id\":\"20210110115402-21ltd5v\"},\"children\":[]}]}}"},
}

func TestKityMinderJSONRenderer(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.KramdownIAL = true

	ast.Testing = true
	for _, test := range kitymindJSONRendererTests {
		jsonStr := luteEngine.RenderKityMinderJSON(test.from)
		t.Log(jsonStr)
		if test.to != jsonStr {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, jsonStr, test.from)
		}
	}
}
