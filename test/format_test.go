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
