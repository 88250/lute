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

type textbundleTest struct {
	name          string
	original      string   // 原始的 Markdown 文本
	textbundle    string   // TextBundle 过的 Markdown 文本
	originalLinks []string // 原始的链接地址
}

var originalLinksCases = [][]string{
	{"https://img.hacpai.com/dir1/bar.zip", "https://b3logfile.com/dir2/baz.png"},
}

var textbundleTests = []textbundleTest{

	{"0", "[foo](" + originalLinksCases[0][0] + ")\n\n![foo](" + originalLinksCases[0][1] + ")", "[foo](assets/dir1/bar.zip)\n\n![foo](assets/dir2/baz.png)\n", originalLinksCases[0]},
}

func TestTextBundle(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range textbundleTests {
		textbundle, originalLinks := luteEngine.TextBundleStr(test.name, test.original, []string{"https://img.hacpai.com", "https://b3logfile.com"})
		if test.textbundle != textbundle {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.textbundle, textbundle, test.original)
		}
		if !equalStrs(test.originalLinks, originalLinks) {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.originalLinks, originalLinks, test.original)
		}
	}
}

func equalStrs(strs1, strs2 []string) bool {
	if len(strs1) != len(strs2) {
		return false
	}
	for i, str1 := range strs1 {
		if strs2[i] != str1 {
			return false
		}
	}
	return true
}
