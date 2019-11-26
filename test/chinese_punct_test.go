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

	"github.com/b3log/lute"
)

var chinesePunctTests = []parseTest{

	// 标点转换需要排除文件后缀 https://github.com/b3log/lute/issues/41
	{"6", "名字@数字.exe\n", "<p>名字@数字.exe</p>\n"},
	{"5", "主页.1html 主页.html1\n", "<p>主页。1html 主页.html1</p>\n"},
	{"4", "程序.exe 程序.exee 程序.no\n", "<p>程序.exe 程序.exee 程序。no</p>\n"},

	{"3", "感叹号!问号?\n", "<p>感叹号！问号？</p>\n"},
	{"2", "中文,。冒号:bar.英文句号在前\n", "<p>中文，。冒号：bar.英文句号在前</p>\n"},
	{"1", "foo,bar.\n", "<p>foo,bar.</p>\n"},
	{"0", "中文,逗号句号.\n", "<p>中文，逗号句号。</p>\n"},
}

func TestChinesePunct(t *testing.T) {
	luteEngine := lute.New() // 默认已经开启标点替换

	for _, test := range chinesePunctTests {
		html, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
