// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of the License at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package test

import (
	"testing"

	"github.com/88250/lute"
)

func TestHTMLBlockDataContentRoundTrip(t *testing.T) {
	luteEngine := lute.New()
	const content = `<div title="&quot;">&lt;test&gt; &amp; &#x3C;</div>`
	blockDOM := `<div data-node-id="20260724120000-abcdefg" data-type="NodeHTMLBlock" class="render-node" data-subtype="block"><div class="protyle-icons"></div><div><protyle-html data-content="&lt;div title=&#34;&amp;quot;&#34;&gt;&amp;lt;test&amp;gt; &amp;amp; &amp;#x3C;&lt;/div&gt;"></protyle-html><span style="position: absolute">​</span></div><div class="protyle-attr" contenteditable="false">​</div></div>`

	for i := 0; i < 3; i++ {
		tree := luteEngine.BlockDOM2Tree(blockDOM)
		node := tree.Root.FirstChild
		if nil == node {
			t.Fatalf("round trip [%d] failed: HTML block node is missing", i)
		}
		if got := string(node.Tokens); content != got {
			t.Fatalf("round trip [%d] failed\nexpected\n\t%q\ngot\n\t%q", i, content, got)
		}
		blockDOM = luteEngine.RenderNodeBlockDOM(node)
	}
}
