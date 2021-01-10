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

	{"2", "foo\nbar", "{\"root\":{\"data\":{\"text\":\"文档名 TODO\",\"id\":\"20060102150405-1a2b3c4\"},\"children\":[{\"data\":{\"text\":\"foobar\",\"id\":\"\"},\"children\":[]}]}}"},
	{"1", "# foo\n\n para1\n\npara2", "{\"root\":{\"data\":{\"text\":\"文档名 TODO\",\"id\":\"20060102150405-1a2b3c4\"},\"children\":[{\"data\":{\"text\":\"# foo\",\"id\":\"\"},\"children\":[{\"data\":{\"text\":\"para1\",\"id\":\"\"},\"children\":[]},{\"data\":{\"text\":\"para2\",\"id\":\"\"},\"children\":[]}]}]}}"},
	{"0", "foo **bar**\n", "{\"root\":{\"data\":{\"text\":\"文档名 TODO\",\"id\":\"20060102150405-1a2b3c4\"},\"children\":[{\"data\":{\"text\":\"foo **bar**\",\"id\":\"\"},\"children\":[]}]}}"},
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
