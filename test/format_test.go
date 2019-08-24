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

type formatTest struct {
	name      string
	original  string // 原始的 Markdown 文本
	formatted string // 格式化过的 Markdown 文本
}

var formatTests = []formatTest{
	{"15", "| abc | def |\n| --- | --- |\n", "| abc | def |\n| --- | --- |\n"},
	{"14", "~~B3log~~\n", "~~B3log~~\n"},
	{"13", "![B3log 开源](https://b3log.org \"B3log 开源\")\n", "![B3log 开源](https://b3log.org \"B3log 开源\")\n"},
	{"12", "[B3log 开源](https://b3log.org \"B3log 开源\")\n", "[B3log 开源](https://b3log.org \"B3log 开源\")\n"},
	{"11", "硬换行  \n第二行\n", "硬换行\\\n第二行\n"},
	{"10", "硬换行\\\n第二行\n", "硬换行\\\n第二行\n"},
	{"9", "分隔线\n\n---\n", "分隔线\n\n---\n"},
	{"8", "```go\nvar lute\n```\n", "```go\nvar lute\n```\n"},
	{"7", "`代码`\n", "`代码`\n"},
	{"6", ">块引用\n", ">块引用\n"},
	{"5", "**加粗**格式化\n", "**加粗**格式化\n"},
	{"4", "_强调_ 格式化\n", "*强调* 格式化\n"},
	{"3", "*强调*格式化\n", "*强调*格式化\n"},
	{"2", "1.  列表项\n   * 子列表项\n", "1. 列表项\n   * 子列表项\n"},
	{"1", "*  列表项\n   * 子列表项\n", "* 列表项\n  * 子列表项\n"},
	{"0", "# 标题\n\n段落用一个空行分隔就够了。\n\n\n这是第二段。", "# 标题\n段落用一个空行分隔就够了。\n这是第二段。\n"},
}

func TestFormat(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range formatTests {
		t.Log("Test [" + test.name + "]")
		formatted, err := luteEngine.FormatStr(test.name, test.original)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.formatted != formatted {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.formatted, formatted, test.original)
		}
	}
}
