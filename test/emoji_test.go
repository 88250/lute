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

var emojiTests = []parseTest{

	{"15", "1 :+1::1::+1:1\n", "<p>1 👍:1:👍1</p>\n"},
	{"14", "1 :1: 1\n", "<p>1 :1: 1</p>\n"},
	{"13", ":1:\n", "<p>:1:</p>\n"}, // 冒号解析错误 https://github.com/b3log/lute/issues/12
	{"12", ":smile::smile:\n", "<p>😄😄</p>\n"},
	{"11", "::\n", "<p>::</p>\n"},
	{"10", "smile: :heart :smile:\n", "<p>smile: :heart 😄</p>\n"},
	{"9", ":smile: :heart :smile:\n", "<p>😄 :heart 😄</p>\n"},
	{"8", ":heart\n", "<p>:heart</p>\n"},
	{"7", ":heart 不是表情\n", "<p>:heart 不是表情</p>\n"},
	{"6", ":heart:开头表情\n", "<p>❤️开头表情</p>\n"},
	{"5", "结尾表情:heart:\n", "<p>结尾表情❤️</p>\n"},
	{"4", "没有表情\n", "<p>没有表情</p>\n"},
	{"3", "0 :b3log: 1 :heart: 2\n", "<p>0 <img alt=\"b3log\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/b3log.png\" title=\"b3log\" /> 1 ❤️ 2</p>\n"},
	{"2", ":smile: :heart:\n", "<p>😄 ❤️</p>\n"},
	{"1", ":b3log:\n", "<p><img alt=\"b3log\" class=\"emoji\" src=\"https://cdn.jsdelivr.net/npm/vditor/dist/images/emoji/b3log.png\" title=\"b3log\" /></p>\n"},
	{"0", "爱心:heart:一个\n", "<p>爱心❤️一个</p>\n"},
}

func TestEmoji(t *testing.T) {
	luteEngine := lute.New() // 默认已经开启 Emoji 处理

	for _, test := range emojiTests {
		html, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
