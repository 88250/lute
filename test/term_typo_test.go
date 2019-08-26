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
"fmt"
"github.com/b3log/lute"
"testing"
)

var termTypoTests = []parseTest{

	{"0", "在github上做开源项目是一件很开心的事情，请不要把Github拼写成`github`哦！\n", "<p>在 GitHub 上做开源项目是一件很开心的事情，请不要把 GitHub 拼写成<code>github</code>哦！</p>\n"},
}

func TestTermTypo(t *testing.T) {
	luteEngine := lute.New() // 默认已经开启了术语修正

	for _, test := range termTypoTests {
		fmt.Println("Test [" + test.name + "]")
		html, err := luteEngine.MarkdownStr(test.name, test.markdown)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.html != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.html, html, test.markdown)
		}
		fmt.Println(html)
		fmt.Println(luteEngine.FormatStr(test.name, test.markdown))
	}
}
