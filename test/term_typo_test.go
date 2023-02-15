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

var termTypoTests = []parseTest{

	{"7", "customtest\n", "<p>CUSTOMtest</p>\n"},
	{"6", "gorm orm\n", "<p>GORM ORM</p>\n"},
	{"5", "test.html\n", "<p>test.html</p>\n"},
	{"4", "cookie ie ieo\n", "<p>cookie IE ieo</p>\n"},
	{"3", "github.com\n", "<p>github.com</p>\n"},
	{"2", "https://github.com\n", "<p><a href=\"https://github.com\">https://github.com</a></p>\n"},
	{"1", "特别是简历中千万不要出现这样的情况：熟练使用JAVA、Javascript、GIT，对android、ios开发有一定了解，熟练使用Mysql、postgresql数据库。\n", "<p>特别是简历中千万不要出现这样的情况：熟练使用 Java、JavaScript、Git，对 Android、iOS 开发有一定了解，熟练使用 MySQL、PostgreSQL 数据库。</p>\n"},
	{"0", "在github上做开源项目是一件很开心的事情，请不要把Github拼写成`github`哦！\n", "<p>在 GitHub 上做开源项目是一件很开心的事情，请不要把 GitHub 拼写成 <code>github</code> 哦！</p>\n"},
}

func TestTermTypo(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetAutoSpace(true)
	luteEngine.SetFixTermTypo(true)
	luteEngine.PutTerms(map[string]string{"customtest": "CUSTOMtest"})
	for _, test := range termTypoTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

var termTypoFormatTests = []parseTest{

	{"0", "这是一个 github 测试", "这是一个 GitHub 测试\n"},
}

func TestTermTypoFormat(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetAutoSpace(true)
	luteEngine.SetFixTermTypo(true)
	// luteEngine.PutTerms(map[string]string{"customtest": "CUSTOMtest"})
	for _, test := range termTypoFormatTests {
		html := luteEngine.FormatStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
