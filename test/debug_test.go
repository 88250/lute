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

var debugTests = []parseTest{

	// 出现在 invalidAutoLinkDomain 中的不要解析自动链接
	{"28", "https://v.qq.com/x/cover/zf2z0xpqcculhcz/y0016tj0qvh.html", "<p>https://v.qq.com/x/cover/zf2z0xpqcculhcz/y0016tj0qvh.html</p>\n"},

	{"27", "[]( https://b3log.org )", "<p><a href=\"https://b3log.org\"></a></p>\n"},
	{"26", "[](https://b3log.org)", "<p><a href=\"https://b3log.org\"></a></p>\n"},
	{"25", "[]( https://b3log.org", "<p>[]( <a href=\"https://b3log.org\">https://b3log.org</a></p>\n"},

	// GFM 任务列表 li 加 class="vditor-task" https://github.com/b3log/lute/issues/10
	{"24", "- [x]\n", "<ul>\n<li class=\"vditor-task\"><input checked=\"\" disabled=\"\" type=\"checkbox\" /></li>\n</ul>\n"},

	// Empty list following GFM Table makes table broken https://github.com/b3log/lute/issues/9
	{"23", "0\n-:\n1\n-\n", "<table>\n<thead>\n<tr>\n<th align=\"right\">0</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td align=\"right\">1</td>\n</tr>\n</tbody>\n</table>\n<ul>\n<li></li>\n</ul>\n"},
	{"22", "0\n-:\n-\n", "<table>\n<thead>\n<tr>\n<th align=\"right\">0</th>\n</tr>\n</thead>\n</table>\n<ul>\n<li></li>\n</ul>\n"},

	// GFM Table rendered as h2 https://github.com/b3log/lute/issues/3
	{"21", "0\n-:\n", "<table>\n<thead>\n<tr>\n<th align=\"right\">0</th>\n</tr>\n</thead>\n</table>\n"},

	// HTMl 块解析，等号前面空格情况
	{"20", "<a href =\"https://github.com\">GitHub</a>\n", "<a href =\"https://github.com\">GitHub</a>\n"},

	// 链接结尾 / 处理
	{"19", "https://hacpai.com/ https://hacpai.com", "<p><a href=\"https://hacpai.com/\">https://hacpai.com/</a> <a href=\"https://hacpai.com\">https://hacpai.com</a></p>\n"},

	// 转义
	{"18", "`<a href=\"`\">`\n", "<p><code>&lt;a href=&quot;</code>&quot;&gt;`</p>\n"},

	// 原文不以 \n 结尾的话需要自动补上
	{"17", "- - ", "<ul>\n<li>\n<ul>\n<li></li>\n</ul>\n</li>\n</ul>\n"},

	// 强调优先级高于删除线
	{"16", "~~*~~Hi*\n", "<p>~~<em>~~Hi</em></p>\n"},

	{"15", "a*\"foo\"*\n", "<p>a*&quot;foo&quot;*</p>\n"},
	{"14", "5*6*78\n", "<p>5<em>6</em>78</p>\n"},
	{"13", "**莠**\n", "<p><strong>莠</strong></p>\n"},
	{"12", "**章**\n", "<p><strong>章</strong></p>\n"},
	{"11", "1>tag<\n", "<p>1&gt;tag&lt;</p>\n"},
	{"10", "<http:\n", "<p>&lt;http:</p>\n"},
	{"9", "<\n", "<p>&lt;</p>\n"},
	{"8", "~~~ \n", "<pre><code class=\"language-fallback\"></code></pre>\n"},
	{"7", "|||\n|||\n", "<p>|||<br />\n|||</p>\n"},
	{"6", "[https://github.com/b3log/lute](https://github.com/b3log/lute)\n", "<p><a href=\"https://github.com/b3log/lute\">https://github.com/b3log/lute</a></p>\n"},
	{"5", "[1\n--\n", "<h2>[1</h2>\n"},
	{"4", "[1 \n", "<p>[1</p>\n"},
	{"3", "- -\r\n", "<ul>\n<li>\n<ul>\n<li></li>\n</ul>\n</li>\n</ul>\n"},
	{"2", "foo@bar.baz\n", "<p><a href=\"mailto:foo@bar.baz\">foo@bar.baz</a></p>\n"},
	{"1", "B3log https://b3log.org Lute\n", "<p>B3log <a href=\"https://b3log.org\">https://b3log.org</a> Lute</p>\n"},
	{"0", "[https://b3log.org](https://b3log.org)\n", "<p><a href=\"https://b3log.org\">https://b3log.org</a></p>\n"},
}

func TestDebug(t *testing.T) {
	luteEngine := lute.New()
	for _, test := range debugTests {
		html, err := luteEngine.MarkdownStr(test.name, test.markdown)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.html != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.html, html, test.markdown)
		}
	}
}
