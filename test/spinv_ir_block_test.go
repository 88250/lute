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
	"strings"
	"testing"

	"github.com/88250/lute"
)

var spinVditorIRBlockDOMTests = []*parseTest{

	// TODO: 图片路径包含空格的问题 #104
	//{"5", "<p data-block=\"0\" data-node-id=\"20200916084414-0hk6l8v\" data-type=\"p\">![foo中文.png](assets/foo中文<wbr>.png)</p>", "<h1 data-block=\"0\" class=\"vditor-ir__node\" data-node-id=\"test\" bookmark=\"bookmark\" data-type=\"h\" id=\"ir-foo\" data-marker=\"#\"><span class=\"vditor-ir__marker vditor-ir__marker--heading\" data-type=\"heading-marker\"># </span>foo</h1>"},
	{"4", "<h1 data-block=\"0\" class=\"vditor-ir__node\" data-node-id=\"test\" bookmark=\"bookmark\" data-type=\"h\" id=\"ir-foo\" data-marker=\"#\"><span class=\"vditor-ir__marker vditor-ir__marker--heading\" data-type=\"heading-marker\"># </span>foo</h1>", "<h1 data-block=\"0\" class=\"vditor-ir__node\" data-node-id=\"test\" bookmark=\"bookmark\" data-type=\"h\" id=\"ir-foo\" data-marker=\"#\"><span class=\"vditor-ir__marker vditor-ir__marker--heading\" data-type=\"heading-marker\"># </span>foo</h1>"},
	{"3", "<div data-block=\"0\" data-node-id=\"test\" bookmark=\"bookmark\" data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">\u200b</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code>foo\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><code>foo\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div>", "<div data-block=\"0\" data-node-id=\"test\" bookmark=\"bookmark\" data-type=\"code-block\" class=\"vditor-ir__node\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">\u200b</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code>foo\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><code>foo\n</code></pre><span data-type=\"code-block-close-marker\">```</span></div>"},
	{"2", "<p data-block=\"0\" data-node-id=\"20200915173154-1wi2p2h\" data-type=\"p\">$$<wbr></p>", "<div data-block=\"0\" data-node-id=\"20200915173154-1wi2p2h\" data-type=\"math-block\" class=\"vditor-ir__node vditor-ir__node--expand\"><span data-type=\"math-block-open-marker\">$$</span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code data-type=\"math-block\" class=\"language-math\"><wbr>\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><code data-type=\"math-block\" class=\"language-math\"></code></pre><span data-type=\"math-block-close-marker\">$$</span></div>"},
	{"1", "<p data-block=\"0\" data-node-id=\"20200915172226-iexs3bo\" data-type=\"p\">```<wbr></p>", "<div data-block=\"0\" data-node-id=\"20200915172226-iexs3bo\" data-type=\"code-block\" class=\"vditor-ir__node vditor-ir__node--expand\"><span data-type=\"code-block-open-marker\">```</span><span class=\"vditor-ir__marker vditor-ir__marker--info\" data-type=\"code-block-info\">\u200b<wbr></span><pre class=\"vditor-ir__marker--pre vditor-ir__marker\"><code>\n</code></pre><pre class=\"vditor-ir__preview\" data-render=\"2\"><code></code></pre><span data-type=\"code-block-close-marker\">```</span></div>"},
	{"0", "<p data-block=\"0\" data-node-id=\"20200914181352-laa3jyd\" data-type=\"p\">foo<wbr></p>", "<p data-block=\"0\" data-node-id=\"20200914181352-laa3jyd\" data-type=\"p\">foo<wbr></p>"},
}

func TestSpinVditorIRBlockDOM(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.BlockRef = true
	luteEngine.KramdownIAL = true
	luteEngine.SetLinkBase(" http://127.0.0.1:6807/webdav/q0fk7yv/测试笔记/")

	for _, test := range spinVditorIRBlockDOMTests {
		html := luteEngine.SpinVditorIRBlockDOM(test.from)

		if "15" == test.name || "18" == test.name {
			// 去掉最后一个新生成的节点 ID，因为这个节点 ID 是随机生成，每次测试用例运行都不一样，比较没有意义，长度一致即可
			lastNodeIDIdx := strings.LastIndex(html, "data-node-id=")
			length := len("data-node-id=\"20200813190226-1234567\" ")
			html = html[:lastNodeIDIdx] + html[lastNodeIDIdx+length:]
			lastNodeIDIdx = strings.LastIndex(test.to, "data-node-id=")
			test.to = test.to[:lastNodeIDIdx] + test.to[lastNodeIDIdx+length:]
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
