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
	"testing"

	"github.com/88250/lute"
)

// 开关开启：空列表项下创建子列表前补一个空段落（HTML 渲染会省略空段落，效果与默认一致）
var disableNestedListTests = []parseTest{
	{"空列表项-子列表", "-\n  - bar\n",
		"<ul>\n<li>\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>\n"},
	{"一行-连写", "- - bar\n",
		"<ul>\n<li>\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>\n"},
	{"一行*连写", "* * bar\n",
		"<ul>\n<li>\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>\n"},
	{"一行有序连写", "1. 1. bar\n",
		"<ol>\n<li>\n<ol>\n<li>bar</li>\n</ol>\n</li>\n</ol>\n"},
	{"有内容列表项-子列表-不受影响", "- foo\n  - bar\n",
		"<ul>\n<li>foo\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>\n"},
	{"同级列表-不受影响", "- foo\n- bar\n",
		"<ul>\n<li>foo</li>\n<li>bar</li>\n</ul>\n"},
}

func TestEnsureListItemParagraph(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetEnsureListItemParagraph(true)
	luteEngine.SetHeadingID(false)
	luteEngine.SetKramdownIAL(false)

	for _, test := range disableNestedListTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal\n\t%q",
				test.name, test.to, html, test.from)
		}
	}
}

// 开关关闭（默认）：空列表项下直接挂子列表
var defaultNestedListTests = []parseTest{
	{"空列表项-子列表", "-\n  - bar\n",
		"<ul>\n<li>\n<ul>\n<li>bar</li>\n</ul>\n</li>\n</ul>\n"},
}

func TestDefaultNestedList(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetHeadingID(false)
	luteEngine.SetKramdownIAL(false)

	for _, test := range defaultNestedListTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal\n\t%q",
				test.name, test.to, html, test.from)
		}
	}
}
