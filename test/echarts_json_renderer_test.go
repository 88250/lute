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
)

var echartsJSONRendererTests = []parseTest{

	{"7", "| foo |\n| - |\n|   |\n{: id=\"20201122125005-kc4sl0l\"}\n\n{: id=\"20201122125319-t74imwc\" type=\"doc\"}", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Table\\ntable\",\"children\":[{\"name\":\"Table Head\\nthead\"},{\"name\":\"Table Row\\ntr\"}]},{\"name\":\"Block IAL\\n{: 20201122125005-kc4sl0l}\"},{\"name\":\"Block IAL\\n{: 20201122125319-t74imwc}\"}]}]"},
	{"6", "foo\n{: id=\"20201122125005-kc4sl0l\"}\n\n{: id=\"20201122125319-t74imwc\" type=\"doc\"}", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Paragraph\\np\",\"children\":[{\"name\":\"Text\\nfoo\"}]},{\"name\":\"Block IAL\\n{: 20201122125005-kc4sl0l}\"},{\"name\":\"Block IAL\\n{: 20201122125319-t74imwc}\"}]}]"},
	{"5", "foo\n{: id=\"fooid\"}\n\n{: id=\"20201122125319-t74imwc\" type=\"doc\"}", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Paragraph\\np\",\"children\":[{\"name\":\"Text\\nfoo\"}]},{\"name\":\"Block IAL\\n{: fooid}\"},{\"name\":\"Block IAL\\n{: 20201122125319-t74imwc}\"}]}]"},
	{"4", "&hearts;\n\n{: id=\"20201122125319-t74imwc\" type=\"doc\"}", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Paragraph\\np\",\"children\":[{\"name\":\"HTML Entity\\nspan\"}]},{\"name\":\"Block IAL\\n{: 20201122125319-t74imwc}\"}]}]"},
	{"3", ":smile:\n\n{: id=\"20201122125319-t74imwc\" type=\"doc\"}", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Paragraph\\np\",\"children\":[{\"name\":\"Emoji Unicode\\n\"}]},{\"name\":\"Block IAL\\n{: 20201122125319-t74imwc}\"}]}]"},
	{"2", "~foo~\n\n{: id=\"20201122125319-t74imwc\" type=\"doc\"}", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Paragraph\\np\",\"children\":[{\"name\":\"Strikethrough\\ndel\"}]},{\"name\":\"Block IAL\\n{: 20201122125319-t74imwc}\"}]}]"},
	{"1", "# foo\n*bar*\n\n{: id=\"20201122125319-t74imwc\" type=\"doc\"}", "[{\"name\":\"Document\",\"children\":[{\"name\":\"Heading\\nh1\",\"children\":[{\"name\":\"Text\\nfoo\"}]},{\"name\":\"Paragraph\\np\",\"children\":[{\"name\":\"Emphasis\\nem\",\"children\":[{\"name\":\"Text\\nbar\"}]}]},{\"name\":\"Block IAL\\n{: 20201122125319-t74imwc}\"}]}]"},
	{"0", "", "[{\"name\":\"Document\",\"children\":[]}]"},
}

func TestEChartsJSONRenderer(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ParseOptions.KramdownIAL = true

	for _, test := range echartsJSONRendererTests {
		jsonStr := luteEngine.RenderEChartsJSON(test.from)
		if test.to != jsonStr {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, jsonStr, test.from)
		}
	}
}
