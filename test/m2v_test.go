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

	"github.com/88250/lute"
)

var md2VditorTests = []parseTest{

	{"2", "<p align=\"center\">\nfoo\n</p>\n\nbar", "<div class=\"vditor-wysiwyg__block\" data-type=\"html-block\" data-block=\"0\"><pre><code>&lt;p align=&quot;center&quot;&gt;\nfoo\n&lt;/p&gt;</code></pre></div><p data-block=\"0\">bar\n</p>"},
	{"1", `foo\<aa>bar`, "<p data-block=\"0\">foo\\&lt;aa&gt;bar\n</p>"},
	{"0", `<details>
<summary>foo</summary>

* bar

</details>`, "<details>\n<summary>foo</summary><ul data-tight=\"true\" data-block=\"0\"><li data-marker=\"*\">bar</li></ul></details>"},
}

func TestMd2Vditor(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range md2VditorTests {
		md := luteEngine.Md2VditorDOM(test.from)
		if test.to != md {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, md, test.from)
		}
	}
}
