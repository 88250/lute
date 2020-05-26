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

var chinesePunctTests = []parseTest{

	{"13", "文件后缀.md测试\n", "<p>文件后缀.md 测试</p>\n"},
	{"13", "文件后缀.textbundle测试\n", "<p>文件后缀.textbundle 测试</p>\n"},

	// 英文逗号标点在英文和中文之间需要渲染为中文逗号 https://github.com/88250/lute/issues/54
	{"12", "test,测试\n", "<p>test，测试</p>\n"},

	{"11", "英文标点!?\n", "<p>英文标点!?</p>\n"},
	{"10", "英文叹号!！\n", "<p>英文叹号！！</p>\n"},
	{"9", "英文叹号!!\n", "<p>英文叹号!!</p>\n"},
	// 连续英文句号出现在中文后不优化 https://github.com/88250/lute/issues/2
	{"8", "英文句号.。\n", "<p>英文句号。。</p>\n"},
	{"7", "英文句号..\n", "<p>英文句号..</p>\n"},

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
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
