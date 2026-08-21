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
	"bytes"
	"strings"
	"testing"

	"github.com/alecthomas/chroma/v2/quick"
)

func TestChroma(t *testing.T) {
	java := `
	@RequestProcessing("/")
	public void index(final RequestContext context) {
		context.setRenderer(new SimpleFMRenderer("index.ftl"));
		final Map<String, Object> dataModel = context.getRenderer().getRenderDataModel();
		dataModel.put("greeting", "Hello, Latke!");
	}`

	writer := bytes.Buffer{}
	err := quick.Highlight(&writer, java, "java", "html", "github")
	if nil != err {
		t.Fatalf("%s", err.Error())
	}
	output := writer.String()
	for _, expected := range []string{
		`<pre class="chroma"><code>`,
		`<span class="nd">@RequestProcessing</span>`,
		`<span class="kd">public</span>`,
		`<span class="s">&#34;Hello, Latke!&#34;</span>`,
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q\ngot\n%s", expected, output)
		}
	}
	if strings.Contains(output, "tabindex=") {
		t.Fatalf("unexpected tabindex in output:\n%s", output)
	}
}
