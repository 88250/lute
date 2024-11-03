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

var debugTests = []parseTest{

	{"76", ":\n0", "<p>:<br />\n0</p>\n"},
	{"75", "[foo](bar.com(baz) \"bar.com(baz)\")", "<p><a href=\"bar.com(baz)\" title=\"bar.com(baz)\">foo</a></p>\n"},
	{"74", "https://foo.com:443/bar/baz", "<p><a href=\"https://foo.com:443/bar/baz\">https://foo.com:443/bar/baz</a></p>\n"},
	{"73", "| $foo\\\\\\|bar$ |\n| - |", "<table>\n<thead>\n<tr>\n<th><span class=\"language-math\">foo\\\\|bar</span></th>\n</tr>\n</thead>\n</table>\n"},
	{"72", "| $foo\\\\|bar$ |\n| - |", "<table>\n<thead>\n<tr>\n<th><span class=\"language-math\">foo\\|bar</span></th>\n</tr>\n</thead>\n</table>\n"},
	{"71", "| $foo\\|bar$ |\n| - |", "<table>\n<thead>\n<tr>\n<th><span class=\"language-math\">foo|bar</span></th>\n</tr>\n</thead>\n</table>\n"},
	{"70", "| `foo\\\\\\|bar` |\n| - |", "<table>\n<thead>\n<tr>\n<th><code>foo\\\\|bar</code></th>\n</tr>\n</thead>\n</table>\n"},
	{"69", "| `foo\\\\|bar` |\n| - |", "<table>\n<thead>\n<tr>\n<th><code>foo\\|bar</code></th>\n</tr>\n</thead>\n</table>\n"},
	{"68", "| `foo\\|bar` |\n| - |", "<table>\n<thead>\n<tr>\n<th><code>foo|bar</code></th>\n</tr>\n</thead>\n</table>\n"},
	{"67", "foo\n:", "<p>foo<br />\n:</p>\n"},
	{"66", "[foo](<www.bar.com> \"baz\")", "<p><a href=\"www.bar.com\" title=\"baz\">foo</a></p>\n"},
	{"65", "foo：bar://baz", "<p>foo：<a href=\"bar://baz\">bar://baz</a></p>\n"},
	{"64", "foo bar://baz", "<p>foo <a href=\"bar://baz\">bar://baz</a></p>\n"},
	{"63", "1. foo\n\nbar\n\n2. baz", "<ol>\n<li>foo</li>\n</ol>\n<p>bar</p>\n<ol start=\"2\">\n<li>baz</li>\n</ol>\n"},
	{"62", "[请从这里开始](siyuan://blocks/20200812220555-lj3enxa)", "<p><a href=\"siyuan://blocks/20200812220555-lj3enxa\">请从这里开始</a></p>\n"},
	{"61", "![a](\"<img src=xss onerror=alert(1)>)", "<p>![a](&quot;&lt;img src=xss onerror=alert(1)&gt;)</p>\n"},
	{"60", "123\n456\n| a | b |\n| ---| --- |\nd | e", "<p>123<br />\n456</p>\n<table>\n<thead>\n<tr>\n<th>a</th>\n<th>b</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>d</td>\n<td>e</td>\n</tr>\n</tbody>\n</table>\n"},
	{"59", "<img src=' foo.png'/>\n", "<p><img src=' foo.png'/></p>\n"},
	{"58", "<img src=\" foo.png\"/>\n", "<p><img src=\" foo.png\"/></p>\n"},

	// CommonMark 0.30 https://spec.commonmark.org/0.30/changes.html#part-66
	{"57", "[Толпой][Толпой] is a Russian word.\n\n[ТОЛПОЙ]: /url\n", "<p><a href=\"/url\">Толпой</a> is a Russian word.</p>\n"},
	{"56", "[SS]\n\n[ẞ]: /url\n", "<p><a href=\"/url\">SS</a></p>\n"},
	{"55", "[ẞ]\n\n[SS]: /url\n", "<p><a href=\"/url\">ẞ</a></p>\n"},

	{"54", "- foo\n\n    ```\nbar\n```", "<ul>\n<li>\n<p>foo</p>\n<pre><code class=\"highlight-chroma\"></code></pre>\n</li>\n</ul>\n<p>bar</p>\n<pre><code class=\"highlight-chroma\"></code></pre>\n"},
	{"53", "- foo\n\n    $$\nbar\n$$\n", "<ul>\n<li>\n<p>foo</p>\n<div class=\"language-math\"></div>\n</li>\n</ul>\n<p>bar</p>\n<div class=\"language-math\"></div>\n"},

	// Auto link `.app` domain suffix https://github.com/Vanessa219/vditor/issues/936
	{"52", "https://netlify.app/", "<p><a href=\"https://netlify.app/\">https://netlify.app/</a></p>\n"},

	// 表格和 Setext 标题解析冲突问题 https://github.com/88250/lute/issues/110
	{"51", "|   foo   | \n| :-----: |\n|   bar   |\n=======\nbaz\n", "<table>\n<thead>\n<tr>\n<th align=\"center\">foo</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td align=\"center\">bar</td>\n</tr>\n</tbody>\n</table>\n<p>=======<br />\nbaz</p>\n"},

	// 表格解析异常 https://github.com/88250/lute/issues/52
	{"50", "foo\nname | age |\n---- | ---\n\nbar", "<p>foo</p>\n<table>\n<thead>\n<tr>\n<th>name</th>\n<th>age</th>\n</tr>\n</thead>\n</table>\n<p>bar</p>\n"},
	{"49", "foo\n| bar |\n| --- |", "<p>foo</p>\n<table>\n<thead>\n<tr>\n<th>bar</th>\n</tr>\n</thead>\n</table>\n"},

	// 自动链接渲染问题 https://github.com/88250/lute/issues/41
	{"48", "中 https://foo bar\n", "<p>中 https://foo bar</p>\n"},
	{"47", "https://中 bar\n", "<p>https://中 bar</p>\n"},
	{"46", "foo https://” bar\n", "<p>foo https://” bar</p>\n"},

	{"45", "*~foo~*bar\n", "<p><em><del>foo</del></em>bar</p>\n"},
	{"44", "~~foo~\n", "<p>~~foo~</p>\n"},
	{"43", "1. [x]\n2. [x] foo\n", "<ol>\n<li>[x]</li>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo</li>\n</ol>\n"},
	{"42", "* [x]\n* [x] foo\n", "<ul>\n<li>[x]</li>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo</li>\n</ul>\n"},
	{"41", "f</\n", "<p>f&lt;/</p>\n"},

	// 自动链接解析结尾 } 问题 https://github.com/88250/lute/issues/4
	{"40", "https://foo.com/bar}", "<p><a href=\"https://foo.com/bar%7D\">https://foo.com/bar}</a></p>\n"},

	{"39", "[label][] 是 label\n\n[label]: https://b3log.org\n", "<p><a href=\"https://b3log.org\">label</a> 是 label</p>\n"},
	{"38", "|abc|def|\n|---|---|\n", "<table>\n<thead>\n<tr>\n<th>abc</th>\n<th>def</th>\n</tr>\n</thead>\n</table>\n"},

	// 链接解析括号匹配问题 https://github.com/b3log/lute/issues/36
	{"37", "[link](/u(ri\n)\n", "<p>[link](/u(ri<br />\n)</p>\n"},
	{"36", "[link](/u(ri )\n", "<p>[link](/u(ri )</p>\n"},

	{"35", "* [ ] foo [foo](/bar)\n", "<ul>\n<li class=\"vditor-task\"><input disabled=\"\" type=\"checkbox\" /> foo <a href=\"/bar\">foo</a></li>\n</ul>\n"},
	{"34", "[foo](/bar )1\n", "<p><a href=\"/bar\">foo</a>1</p>\n"},
	{"33", "[foo](/bar \"baz\"\n", "<p>[foo](/bar &quot;baz&quot;</p>\n"},
	{"32", "пристаням_стремятся_", "<p>пристаням_стремятся_</p>\n"},
	{"31", "*foo*<br>", "<p><em>foo</em><br></p>\n"},
	{"30", "https://t.mex .mex 后缀不自动链接", "<p>https://t.mex .mex 后缀不自动链接</p>\n"},
	{"29", "https://t.me .me 后缀自动链接", "<p><a href=\"https://t.me\">https://t.me</a> .me 后缀自动链接</p>\n"},

	// 链接解析问题 https://github.com/b3log/lute/issues/20
	{"28", "[新建.zip](/路)foo", "<p><a href=\"/%E8%B7%AF\">新建.zip</a>foo</p>\n"},
	{"27", "[新建.zip](http://bar.com/文件.zip)[新建.zip](http://bar.com/文件.zip)", "<p><a href=\"http://bar.com/%E6%96%87%E4%BB%B6.zip\">新建.zip</a><a href=\"http://bar.com/%E6%96%87%E4%BB%B6.zip\">新建.zip</a></p>\n"},

	{"26", "[]( https://b3log.org )", "<p><a href=\"https://b3log.org\"></a></p>\n"},
	{"25", "[](https://b3log.org)", "<p><a href=\"https://b3log.org\"></a></p>\n"},
	{"24", "[]( https://b3log.org", "<p>[]( <a href=\"https://b3log.org\">https://b3log.org</a></p>\n"},

	// GFM 任务列表 li 加 class="vditor-task" https://github.com/b3log/lute/issues/10
	{"23", "- [x] foo\n", "<ul>\n<li class=\"vditor-task vditor-task--done\"><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo</li>\n</ul>\n"},

	// Empty list following GFM Table makes table broken https://github.com/b3log/lute/issues/9
	{"22", "0\n-:\n1\n-\n", "<table>\n<thead>\n<tr>\n<th align=\"right\">0</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td align=\"right\">1</td>\n</tr>\n</tbody>\n</table>\n<ul>\n<li></li>\n</ul>\n"},
	{"21", "0\n-:\n-\n", "<table>\n<thead>\n<tr>\n<th align=\"right\">0</th>\n</tr>\n</thead>\n</table>\n<ul>\n<li></li>\n</ul>\n"},

	// GFM Table rendered as h2 https://github.com/b3log/lute/issues/3
	{"20", "0\n-:\n", "<table>\n<thead>\n<tr>\n<th align=\"right\">0</th>\n</tr>\n</thead>\n</table>\n"},

	// HTMl 块解析，等号前面空格情况
	{"19", "<a href =\"https://github.com\">GitHub</a>\n", "<a href =\"https://github.com\">GitHub</a>\n"},

	// 链接结尾 / 处理
	{"18", "https://ld246.com/ https://ld246.com", "<p><a href=\"https://ld246.com/\">https://ld246.com/</a> <a href=\"https://ld246.com\">https://ld246.com</a></p>\n"},

	// 转义
	{"17", "`<a href=\"`\">`\n", "<p><code>&lt;a href=&quot;</code>&quot;&gt;`</p>\n"},

	// 原文不以 \n 结尾的话需要自动补上
	{"16", "- - ", "<ul>\n<li>\n<ul>\n<li></li>\n</ul>\n</li>\n</ul>\n"},

	{"15", "~~*~~Hi*\n", "<p><del>*</del>Hi*</p>\n"},

	{"14", "a*\"foo\"*\n", "<p>a*&quot;foo&quot;*</p>\n"},
	{"13", "5*6*78\n", "<p>5<em>6</em>78</p>\n"},
	{"12", "**莠**\n", "<p><strong>莠</strong></p>\n"},
	{"11", "**章**\n", "<p><strong>章</strong></p>\n"},
	{"10", "1>tag<\n", "<p>1&gt;tag&lt;</p>\n"},
	{"9", "<http:\n", "<p>&lt;http:</p>\n"},
	{"8", "<\n", "<p>&lt;</p>\n"},
	{"7", "|||\n|||\n", "<p>|||<br />\n|||</p>\n"},
	{"6", "[https://github.com/88250/lute](https://github.com/88250/lute)\n", "<p><a href=\"https://github.com/88250/lute\">https://github.com/88250/lute</a></p>\n"},
	{"5", "[1\n--\n", "<h2 id=\"-1\">[1</h2>\n"},
	{"4", "[1 \n", "<p>[1</p>\n"},
	{"3", "- -\r\n", "<ul>\n<li>\n<ul>\n<li></li>\n</ul>\n</li>\n</ul>\n"},
	{"2", "foo@bar.baz\n", "<p><a href=\"mailto:foo@bar.baz\">foo@bar.baz</a></p>\n"},
	{"1", "B3log https://b3log.org Lute\n", "<p>B3log <a href=\"https://b3log.org\">https://b3log.org</a> Lute</p>\n"},
	{"0", "[https://b3log.org](https://b3log.org)\n", "<p><a href=\"https://b3log.org\">https://b3log.org</a></p>\n"},
}

func TestDebug(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetHeadingID(true)
	for _, test := range debugTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
