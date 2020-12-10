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

var chineseParagraphBeginningSpaceTests = []parseTest{

	{"4", "* 段落开头空两格\n\n  第二段", "<ul>\n<li>\n<p>段落开头空两格</p>\n<p>第二段</p>\n</li>\n</ul>\n"},
	{"3", "* 段落开头空两格\n", "<ul>\n<li>段落开头空两格</li>\n</ul>\n"},
	{"2", "> 段落开头空两格\n", "<blockquote>\n<p>段落开头空两格</p>\n</blockquote>\n"},
	{"1", "段落开头空两格\n换行\n\n 第二段", "<p class=\"indent--2\">段落开头空两格<br />\n换行</p>\n<p class=\"indent--2\">第二段</p>\n"},
	{"0", "段落开头空两格\n", "<p class=\"indent--2\">段落开头空两格</p>\n"},
}

func TestChineseParagraphBeginningSpace(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.ChineseParagraphBeginningSpace = true

	for _, test := range chineseParagraphBeginningSpaceTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
