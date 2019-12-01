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

	"github.com/88250/lute"
)

var debugTests = []parseTest{

	{"41", "|abc|def|\n|---|---|\n", "<table>\n<thead>\n<tr>\n<th>abc</th>\n<th>def</th>\n</tr>\n</thead>\n</table>\n"},

	// 链接解析括号匹配问题 https://github.com/b3log/lute/issues/36
	{"40", "[link](/u(ri\n)\n", "<p>[link](/u(ri<br />\n)</p>\n"},
	{"39", "[link](/u(ri )\n", "<p>[link](/u(ri )</p>\n"},

	{"38", "www.我的网址/console\n", "<p>www.我的网址/console</p>\n"},
	{"37", "http://我的网址/console\n", "<p>http://我的网址/console</p>\n"},
	{"36", "http://mydomain/console\n", "<p>http://mydomain/console</p>\n"},
	{"35", "* [ ] foo [foo](/bar)\n", "<ul>\n<li class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" /> foo <a href=\"/bar\">foo</a></li>\n</ul>\n"},
	{"34", "[foo](/bar )1\n", "<p><a href=\"/bar\">foo</a>1</p>\n"},
	{"33", "[foo](/bar \"baz\"\n", "<p>[foo](/bar &quot;baz&quot;</p>\n"},
	{"32", "пристаням_стремятся_", "<p>пристаням_стремятся_</p>\n"},
	{"31", "**foo*<br>", "<p>*<em>foo</em><br></p>\n"},
	{"30", "https://t.mex .mex 后缀不自动链接", "<p>https://t.mex .mex 后缀不自动链接</p>\n"},
	{"29", "https://t.me .me 后缀自动链接", "<p><a href=\"https://t.me\">https://t.me</a> .me 后缀自动链接</p>\n"},

	// 链接解析问题 https://github.com/b3log/lute/issues/20
	{"28", "[新建.zip](/路)foo", "<p><a href=\"/%E8%B7%AF\">新建.zip</a>foo</p>\n"},
	{"27", "[新建.zip](http://bar.com/文件.zip)[新建.zip](http://bar.com/文件.zip)", "<p><a href=\"http://bar.com/%E6%96%87%E4%BB%B6.zip\">新建.zip</a><a href=\"http://bar.com/%E6%96%87%E4%BB%B6.zip\">新建.zip</a></p>\n"},

	{"26", "[]( https://b3log.org )", "<p><a href=\"https://b3log.org\"></a></p>\n"},
	{"25", "[](https://b3log.org)", "<p><a href=\"https://b3log.org\"></a></p>\n"},
	{"24", "[]( https://b3log.org", "<p>[]( <a href=\"https://b3log.org\">https://b3log.org</a></p>\n"},

	// GFM 任务列表 li 加 class="vditor-task" https://github.com/b3log/lute/issues/10
	{"23", "- [x] foo\n", "<ul>\n<li class=\"vditor-task\"><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo</li>\n</ul>\n"},

	// Empty list following GFM Table makes table broken https://github.com/b3log/lute/issues/9
	{"22", "0\n-:\n1\n-\n", "<table>\n<thead>\n<tr>\n<th align=\"right\">0</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td align=\"right\">1</td>\n</tr>\n</tbody>\n</table>\n<ul>\n<li></li>\n</ul>\n"},
	{"21", "0\n-:\n-\n", "<table>\n<thead>\n<tr>\n<th align=\"right\">0</th>\n</tr>\n</thead>\n</table>\n<ul>\n<li></li>\n</ul>\n"},

	// GFM Table rendered as h2 https://github.com/b3log/lute/issues/3
	{"20", "0\n-:\n", "<table>\n<thead>\n<tr>\n<th align=\"right\">0</th>\n</tr>\n</thead>\n</table>\n"},

	// HTMl 块解析，等号前面空格情况
	{"19", "<a href =\"https://github.com\">GitHub</a>\n", "<a href =\"https://github.com\">GitHub</a>\n"},

	// 链接结尾 / 处理
	{"18", "https://hacpai.com/ https://hacpai.com", "<p><a href=\"https://hacpai.com/\">https://hacpai.com/</a> <a href=\"https://hacpai.com\">https://hacpai.com</a></p>\n"},

	// 转义
	{"17", "`<a href=\"`\">`\n", "<p><code>&lt;a href=&quot;</code>&quot;&gt;`</p>\n"},

	// 原文不以 \n 结尾的话需要自动补上
	{"16", "- - ", "<ul>\n<li>\n<ul>\n<li></li>\n</ul>\n</li>\n</ul>\n"},

	// 强调优先级高于删除线
	{"15", "~~*~~Hi*\n", "<p>~~<em>~~Hi</em></p>\n"},

	{"14", "a*\"foo\"*\n", "<p>a*&quot;foo&quot;*</p>\n"},
	{"13", "5*6*78\n", "<p>5<em>6</em>78</p>\n"},
	{"12", "**莠**\n", "<p><strong>莠</strong></p>\n"},
	{"11", "**章**\n", "<p><strong>章</strong></p>\n"},
	{"10", "1>tag<\n", "<p>1&gt;tag&lt;</p>\n"},
	{"9", "<http:\n", "<p>&lt;http:</p>\n"},
	{"8", "<\n", "<p>&lt;</p>\n"},
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
		html, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("test case [%s] unexpected: %s", test.name, err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
