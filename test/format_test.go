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
	"io/ioutil"
	"strings"
	"testing"

	"github.com/88250/lute"
)

type formatTest struct {
	name      string
	original  string // 原始的 Markdown 文本
	formatted string // 格式化过的 Markdown 文本
}

var formatTests = []formatTest{

	// 链接引用格式化改进 https://github.com/88250/lute/issues/36
	{"50", "[text][foo]\n\n[foo]: bar\n", "[text][foo]\n\n[foo]: bar\n"},
	{"49", "[foo]\n\n[foo]: bar\n", "[foo]\n\n[foo]: bar\n"},

	// 格式化支持 Setext 标题 https://github.com/88250/lute/issues/29
	{"48", "Setext 标题\n==", "Setext 标题\n===========\n"},
	{"47", "Setext 标题\n------", "Setext 标题\n-----------\n"},

	// 链接前后自动空格改进 https://github.com/88250/lute/issues/24
	{"46", "中文 [黑客派](https://hacpai.com) 不需要", "中文 [黑客派](https://hacpai.com) 不需要\n"},
	{"45", "数字[1黑客派2](https://hacpai.com)需要", "数字 [1 黑客派 2](https://hacpai.com) 需要\n"},
	{"44", "数字1[HacPai](https://hacpai.com)2需要", "数字 1[HacPai](https://hacpai.com)2 需要\n"},
	{"43", "英文[HacPai](https://hacpai.com)需要", "英文 [HacPai](https://hacpai.com) 需要\n"},
	{"42", "中文[黑客派](https://hacpai.com)不需要", "中文[黑客派](https://hacpai.com)不需要\n"},
	{"41", "[黑客派HacPai](https://hacpai.com)需要", "[黑客派 HacPai](https://hacpai.com) 需要\n"},

	{"40", "foo[^bar]\nfoo[^baz]\nfoo[^bar]: bar\n[^baz]: baz\n", "foo[^bar]\nfoo[^baz]\nfoo[^bar]: bar\n\n[^baz]: baz\n"},
	{"39", "foo[^bar]\n[^bar]: bar\n", "foo[^bar]\n\n[^bar]: bar\n"},
	{"38", "[^bar]: bar\n", "[^bar]: bar\n"},
	{"37", "``foo``、`bar`\n", "``foo``、`bar`\n"},
	{"36", "`foo`、`bar`\n", "`foo`、`bar`\n"},
	{"35", "foo`bar`\n", "foo `bar`\n"},
	{"34", "`bar`\n", "`bar`\n"},
	{"33", "foo`bar`baz\n", "foo `bar` baz\n"},

	// 格式化后丢失反斜杠 https://github.com/88250/lute/issues/16
	{"32", "|foo|\n|-|\n|`\\|bar`|\n", "| foo     |\n| ------- |\n| `\\|bar` |\n"},
	{"31", "|foo|\n|-|\n|\\|bar|\n", "| foo   |\n| ----- |\n| \\|bar |\n"},
	{"30", "\\<foo>\n", "\\<foo>\n"},

	{"29", "1. [X] foo\n", "1. [X] foo\n"},
	{"28", "|f|\n|:-:|\nfoo|\n", "|  f  |\n| :-: |\n| foo |\n"},

	// 子列表格式化后缩进不对 https://github.com/b3log/lute/issues/22
	{"27", "* first\n   * sub first\n* second\n  *  sub second\n", "* first\n  * sub first\n* second\n  * sub second\n"},
	{"26", "* first\n  * sub first\n* second\n  * sub second\n", "* first\n  * sub first\n* second\n  * sub second\n"},

	{"25", "`` `Lute` ``\n", "`` `Lute` ``\n"},

	// 图片 Emoji 依然使用别名 https://github.com/b3log/lute/issues/14
	{"24", ":heart: :hacpai:\n", ":heart: :hacpai:\n"},

	// 原先是 HTML 实体 &amp; 格式化后反转义
	{"23", "&&amp;\n", "&&\n"},
	{"22", "&amp;123&emsp;456\n", "&123\u2003456\n"},

	{"21", "\u2003emsp\n", "\u2003emsp\n"},
	{"20", "~删除线~\n", "~删除线~\n"},
	{"19", "我们**需要Markdown Format**\n", "我们**需要 Markdown Format**\n"},
	{"18", "试下中西文间1自动插入lute空格\n", "试下中西文间 1 自动插入 Lute 空格\n"},
	{"17", "* [ ] 项一\n* [X] 项二\n", "* [ ] 项一\n* [X] 项二\n"},
	{"16", "| abc | defghi |\n:-: | -----------:\nbar | baz\n", "| abc | defghi |\n| :-: | -----: |\n| bar |    baz |\n"},
	{"15", "| abc | def |\n| --- | --- |\n", "| abc | def |\n| --- | --- |\n"},
	{"14", "~~删除线~~\n", "~~删除线~~\n"},
	{"13", "![B3log 开源](https://b3log.org \"B3log 开源\")\n", "![B3log 开源](https://b3log.org \"B3log 开源\")\n"},
	{"12", "[B3log 开源](https://b3log.org \"B3log 开源\")\n", "[B3log 开源](https://b3log.org \"B3log 开源\")\n"},
	{"11", "硬换行  \n第二行\n", "硬换行\n第二行\n"}, // 因为启用了软转硬
	{"10", "硬换行\\\n第二行\n", "硬换行\n第二行\n"}, // 因为启用了软转硬
	{"9", "分隔线\n\n---\n", "分隔线\n\n---\n"},
	{"8", "```go\nvar lute\n```\n", "```go\nvar lute\n```\n"},
	{"7", "`代码`\n", "`代码`\n"},
	{"6", ">块引用\n", "> 块引用\n"},
	{"5", "**加粗**格式化\n", "**加粗**格式化\n"},
	{"4", "_强调_ 格式化\n", "_强调_ 格式化\n"},
	{"3", "*强调*格式化\n", "*强调*格式化\n"},
	{"2", "1.  列表项\n    * 子列表项\n", "1. 列表项\n   * 子列表项\n"},
	{"1", "*  列表项\n", "* 列表项\n"},
	{"0", "# 标题\n\n段落用一个空行分隔就够了。\n\n\n这是第二段。", "# 标题\n\n段落用一个空行分隔就够了。\n\n这是第二段。\n"},
}

func TestFormat(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range formatTests {
		formatted, err := luteEngine.FormatStr(test.name, test.original)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.formatted != formatted {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.formatted, formatted, test.original)
		}
	}
}

func TestFormatCases(t *testing.T) {
	files, err := ioutil.ReadDir(".")
	if nil != err {
		t.Fatalf("read test dir failed: %s", err)
	}

	//skips := "format-case0.md,format-case1.md,format-case2.md,format-case3.md,format-case4.md,format-case5.md" // 用于跳过测试文件，例如 format-case0.md
	skips := ""

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "format-case") || strings.Contains(file.Name(), "formatted") {
			continue
		}
		if strings.Contains(skips, file.Name()) {
			continue
		}

		caseName := file.Name()[:len(file.Name())-3]
		bytes, err := ioutil.ReadFile(caseName + ".md")
		if nil != err {
			t.Fatalf("read case failed: %s", err)
		}

		luteEngine := lute.New()
		htmlBytes, err := luteEngine.Format(caseName+".md", bytes)
		if nil != err {
			t.Fatalf("test case [%s] markdown format failed: %s", caseName, err)
		}
		html := string(htmlBytes)

		bytes, err = ioutil.ReadFile(caseName + "-formatted.md")
		if nil != err {
			t.Fatalf("read case cailed: %s", err)
		}
		expected := string(bytes)
		if expected != html {
			t.Fatalf("test case [%s] failed\nexpected\n%s\ngot\n%s\n", caseName, expected, html)
		}
	}
}
