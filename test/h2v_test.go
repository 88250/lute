// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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

var html2VditorDOMTests = []parseTest{

	{"2", `<!--StartFragment--><span>Use<span>&nbsp;</span></span><code class="language-text">new Date()</code><span><span>&nbsp;</span>and</span><!--EndFragment-->`, "<p data-block=\"0\">Use <code>\u200bnew Date()</code>\u200b and\n</p>"},
	{"1", `<pre><code class="language-text">&gt;</code></pre>`, "<div class=\"vditor-wysiwyg__block\" data-type=\"code-block\" data-block=\"0\" data-marker=\"```\"><pre><code class=\"language-text\">&gt;\n</code></pre></div>"},
	{"0", `<code class="language-text">&gt;</code>`, "<p data-block=\"0\">\u200b<code>\u200b&gt;</code>\u200b\n</p>"},
}

func TestHTML2VditorDOM(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range html2VditorDOMTests {
		result := luteEngine.HTML2VditorDOM(test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}