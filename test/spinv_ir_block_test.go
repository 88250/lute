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

	{"0", "<p data-block=\"0\" data-node-id=\"fooid\">foo</p>\n<p data-block=\"0\" data-node-id=\"barid\">bar</p>", "<p data-block=\"0\" data-node-id=\"fooid\">foo</p><span data-type=\"kramdown-ial\">{: id=\"fooid\"}\n</span><p data-block=\"0\" data-node-id=\"barid\">bar</p><span data-type=\"kramdown-ial\">{: id=\"barid\"}\n</span>"},
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
