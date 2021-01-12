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

	{"0", "{{{\nfoo\n{: id=\"20210110174236-xe4fzwr\"}\n\nbar\n{: id=\"20210110174236-xe4fzwr\"}\n\n}}}\n{: id=\"20210110174239-afmuohb\"}\n\nbaz\n{: id=\"20210110174252-qte4wlu\"}\n\n\n{: id=\"20210110005118-gqrx3wm\" type=\"doc\"}\n", "{\"root\":{\"data\":{\"text\":\"\",\"id\":\"20210110005118-gqrx3wm\",\"type\":\"NodeDocument\"\"isContainer\":true},\"children\":[{\"data\":{\"text\":\"foobar\",\"id\":\"20210110174239-afmuohb\",\"type\":\"NodeSuperBlock\"\"isContainer\":true},\"children\":[{\"data\":{\"text\":\"foo\",\"id\":\"20210110174236-xe4fzwr\",\"type\":\"NodeParagraph\"\"isContainer\":false},\"children\":[]},{\"data\":{\"text\":\"bar\",\"id\":\"20210110174236-xe4fzwr\",\"type\":\"NodeParagraph\"\"isContainer\":false},\"children\":[]}]},{\"data\":{\"text\":\"baz\",\"id\":\"20210110174252-qte4wlu\",\"type\":\"NodeParagraph\"\"isContainer\":false},\"children\":[]}]}}"},
}

func TestKityMinderJSONRenderer(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetKramdownIAL(true)
	luteEngine.SetSuperBlock(true)

	ast.Testing = true
	for _, test := range kitymindJSONRendererTests {
		jsonStr := luteEngine.RenderKityMinderJSON(test.from)
		t.Log(jsonStr)
		if test.to != jsonStr {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, jsonStr, test.from)
		}
	}
}
