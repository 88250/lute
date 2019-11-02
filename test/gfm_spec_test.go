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
	"github.com/b3log/lute"
	"testing"
)

var gfmSpecTests = []parseTest{
	// gfm spec block-level cases

	{"gfm198", "| foo | bar |\n| --- | --- |\n| baz | bim |\n", "<table>\n<thead>\n<tr>\n<th>foo</th>\n<th>bar</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>baz</td>\n<td>bim</td>\n</tr>\n</tbody>\n</table>\n"},
	{"gfm199", "| abc | defghi |\n:-: | -----------:\nbar | baz\n", "<table>\n<thead>\n<tr>\n<th align=\"center\">abc</th>\n<th align=\"right\">defghi</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td align=\"center\">bar</td>\n<td align=\"right\">baz</td>\n</tr>\n</tbody>\n</table>\n"},
	{"gfm200", "| f\\|oo  |\n| ------ |\n| b `\\|` az |\n| b **\\|** im |\n", "<table>\n<thead>\n<tr>\n<th>f|oo</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>b <code>|</code> az</td>\n</tr>\n<tr>\n<td>b <strong>|</strong> im</td>\n</tr>\n</tbody>\n</table>\n"},
	{"gfm201", "| abc | def |\n| --- | --- |\n| bar | baz |\n> bar\n", "<table>\n<thead>\n<tr>\n<th>abc</th>\n<th>def</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>bar</td>\n<td>baz</td>\n</tr>\n</tbody>\n</table>\n<blockquote>\n<p>bar</p>\n</blockquote>\n"},
	{"gfm202", "| abc | def |\n| --- | --- |\n| bar | baz |\nbar\n\nbar\n", "<table>\n<thead>\n<tr>\n<th>abc</th>\n<th>def</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>bar</td>\n<td>baz</td>\n</tr>\n<tr>\n<td>bar</td>\n<td></td>\n</tr>\n</tbody>\n</table>\n<p>bar</p>\n"},
	{"gfm203", "| abc | def |\n| --- |\n| bar |\n", "<p>| abc | def |\n| --- |\n| bar |</p>\n"},
	{"gfm204", "| abc | def |\n| --- | --- |\n| bar |\n| bar | baz | boo |\n", "<table>\n<thead>\n<tr>\n<th>abc</th>\n<th>def</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td>bar</td>\n<td></td>\n</tr>\n<tr>\n<td>bar</td>\n<td>baz</td>\n</tr>\n</tbody>\n</table>\n"},
	{"gfm205", "| abc | def |\n| --- | --- |\n", "<table>\n<thead>\n<tr>\n<th>abc</th>\n<th>def</th>\n</tr>\n</thead>\n</table>\n"},
	{"gfm279", "- [ ] foo\n- [x] bar\n", "<ul>\n<li><input disabled=\"\" type=\"checkbox\" /> foo</li>\n<li><input checked=\"\" disabled=\"\" type=\"checkbox\" /> bar</li>\n</ul>\n"},
	{"gfm280", "- [x] foo\n  - [ ] bar\n  - [x] baz\n- [ ] bim\n", "<ul>\n<li><input checked=\"\" disabled=\"\" type=\"checkbox\" /> foo\n<ul>\n<li><input disabled=\"\" type=\"checkbox\" /> bar</li>\n<li><input checked=\"\" disabled=\"\" type=\"checkbox\" /> baz</li>\n</ul>\n</li>\n<li><input disabled=\"\" type=\"checkbox\" /> bim</li>\n</ul>\n"},
	{"task3", "- [ ]foo\n", "<ul>\n<li>[ ]foo</li>\n</ul>\n"},
	{"task2", "- [ ] foo\n", "<ul>\n<li><input disabled=\"\" type=\"checkbox\" /> foo</li>\n</ul>\n"},
	{"task1", "- [ ] \n", "<ul>\n<li>[ ]</li>\n</ul>\n"},
	{"task0", "- [ ]\n", "<ul>\n<li>[ ]</li>\n</ul>\n"},

	// gfm spec inline-level cases

	{"gfm491", "~~Hi~~ Hello, world!\n", "<p><del>Hi</del> Hello, world!</p>\n"},
	{"gfm492", "This ~~has a\n\nnew paragraph~~.\n", "<p>This ~~has a</p>\n<p>new paragraph~~.</p>\n"},
	{"strikethrough0", "**~~Hi~~** Hello, world!\n", "<p><strong><del>Hi</del></strong> Hello, world!</p>\n"},
	{"strikethrough1", "~~**Hi**~~ Hello, world!\n", "<p><del><strong>Hi</strong></del> Hello, world!</p>\n"},
	{"strikethrough2", "~~**Hi~~** Hello, world!\n", "<p><del>**Hi</del>** Hello, world!</p>\n"},
	{"strikethrough3", "**~~**Hi~~ Hello, world!\n", "<p>**<del>**Hi</del> Hello, world!</p>\n"},
	{"gfm621", "www.commonmark.org\n", "<p><a href=\"http://www.commonmark.org\">www.commonmark.org</a></p>\n"},
	{"gfm622", "Visit www.commonmark.org/help for more information.\n", "<p>Visit <a href=\"http://www.commonmark.org/help\">www.commonmark.org/help</a> for more information.</p>\n"},
	{"gfm623", "Visit www.commonmark.org.\n\nVisit www.commonmark.org/a.b.\n", "<p>Visit <a href=\"http://www.commonmark.org\">www.commonmark.org</a>.</p>\n<p>Visit <a href=\"http://www.commonmark.org/a.b\">www.commonmark.org/a.b</a>.</p>\n"},
	{"autolink0", "www.google.com/search?q=Markup+(business)\n", "<p><a href=\"http://www.google.com/search?q=Markup+(business)\">www.google.com/search?q=Markup+(business)</a></p>\n"},
	{"autolink1", "www.google.com/search?q=Markup+(business)))\n", "<p><a href=\"http://www.google.com/search?q=Markup+(business)\">www.google.com/search?q=Markup+(business)</a>))</p>\n"},
	{"autolink2", "(www.google.com/search?q=Markup+(business))\n", "<p>(<a href=\"http://www.google.com/search?q=Markup+(business)\">www.google.com/search?q=Markup+(business)</a>)</p>\n"},
	{"autolink3", "(www.google.com/search?q=Markup+(business)\n", "<p>(<a href=\"http://www.google.com/search?q=Markup+(business)\">www.google.com/search?q=Markup+(business)</a></p>\n"},
	{"gfm624", "www.google.com/search?q=Markup+(business)\n\nwww.google.com/search?q=Markup+(business)))\n\n(www.google.com/search?q=Markup+(business))\n\n(www.google.com/search?q=Markup+(business)\n", "<p><a href=\"http://www.google.com/search?q=Markup+(business)\">www.google.com/search?q=Markup+(business)</a></p>\n<p><a href=\"http://www.google.com/search?q=Markup+(business)\">www.google.com/search?q=Markup+(business)</a>))</p>\n<p>(<a href=\"http://www.google.com/search?q=Markup+(business)\">www.google.com/search?q=Markup+(business)</a>)</p>\n<p>(<a href=\"http://www.google.com/search?q=Markup+(business)\">www.google.com/search?q=Markup+(business)</a></p>\n"},
	{"gfm625", "www.google.com/search?q=(business))+ok\n", "<p><a href=\"http://www.google.com/search?q=(business))+ok\">www.google.com/search?q=(business))+ok</a></p>\n"},
	{"gfm626", "www.google.com/search?q=commonmark&hl=en\n\nwww.google.com/search?q=commonmark&hl;\n", "<p><a href=\"http://www.google.com/search?q=commonmark&amp;hl=en\">www.google.com/search?q=commonmark&amp;hl=en</a></p>\n<p><a href=\"http://www.google.com/search?q=commonmark\">www.google.com/search?q=commonmark</a>&amp;hl;</p>\n"},
	{"gfm627", "www.commonmark.org/he<lp\n", "<p><a href=\"http://www.commonmark.org/he\">www.commonmark.org/he</a>&lt;lp</p>\n"},
	{"gfm628", "http://commonmark.org\n\n(Visit https://encrypted.google.com/search?q=Markup+(business))\n\nAnonymous FTP is available at ftp://foo.bar.baz.\n", "<p><a href=\"http://commonmark.org\">http://commonmark.org</a></p>\n<p>(Visit <a href=\"https://encrypted.google.com/search?q=Markup+(business)\">https://encrypted.google.com/search?q=Markup+(business)</a>)</p>\n<p>Anonymous FTP is available at <a href=\"ftp://foo.bar.baz\">ftp://foo.bar.baz</a>.</p>\n"},
	{"gfm629", "foo@bar.baz\n", "<p><a href=\"mailto:foo@bar.baz\">foo@bar.baz</a></p>\n"},
	{"gfm630", "hello@mail+xyz.example isn't valid, but hello+xyz@mail.example is.\n", "<p>hello@mail+xyz.example isn't valid, but <a href=\"mailto:hello+xyz@mail.example\">hello+xyz@mail.example</a> is.</p>\n"},
	{"auto email link0", "a.b-c_d@a.b\n", "<p><a href=\"mailto:a.b-c_d@a.b\">a.b-c_d@a.b</a></p>\n"},
	{"auto email link1", "a.b-c_d@a.b.\n", "<p><a href=\"mailto:a.b-c_d@a.b\">a.b-c_d@a.b</a>.</p>\n"},
	{"auto email link2", "a.b-c_d@a.b-\n", "<p>a.b-c_d@a.b-</p>\n"},
	{"auto email link3", "a.b-c_d@a.b_\n", "<p>a.b-c_d@a.b_</p>\n"},
	{"gfm631", "a.b-c_d@a.b\n\na.b-c_d@a.b.\n\na.b-c_d@a.b-\n\na.b-c_d@a.b_\n", "<p><a href=\"mailto:a.b-c_d@a.b\">a.b-c_d@a.b</a></p>\n<p><a href=\"mailto:a.b-c_d@a.b\">a.b-c_d@a.b</a>.</p>\n<p>a.b-c_d@a.b-</p>\n<p>a.b-c_d@a.b_</p>\n"},
}

func TestGFMSpec(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SoftBreak2HardBreak = false
	luteEngine.AutoSpace = false
	luteEngine.GFMTaskListItemClass = "" // 关闭类名渲染
	luteEngine.AddAutoLinkDomainSuffix("baz")

	for _, test := range gfmSpecTests {
		html, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
