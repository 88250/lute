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

var autoLinkTests = []parseTest{

	{"31", "https://github.com/88250/lute/pull/207", "<p><a href=\"https://github.com/88250/lute/pull/207\">Pull Request #207 · 88250/lute</a></p>\n"},
	{"30", "https://foo.com/bar:baz", "<p><a href=\"https://foo.com/bar:baz\">https://foo.com/bar:baz</a></p>\n"},
	{"29", "www.test.com:8080/bar/baz", "<p><a href=\"http://www.test.com:8080/bar/baz\">www.test.com:8080/bar/baz</a></p>\n"},
	{"28", "www.test.com:8080(测试地址)", "<p><a href=\"http://www.test.com:8080\">www.test.com:8080</a>(测试地址)</p>\n"},
	{"27", "http://127.0.0.1:8080(测试地址)", "<p><a href=\"http://127.0.0.1:8080\">http://127.0.0.1:8080</a>(测试地址)</p>\n"},
	{"26", "https://www.baidu.help", "<p><a href=\"https://www.baidu.help\">https://www.baidu.help</a></p>\n"},
	{"25", "https://www.baidu.wang", "<p><a href=\"https://www.baidu.wang\">https://www.baidu.wang</a></p>\n"},
	{"24", "https://www.google.com.np/", "<p><a href=\"https://www.google.com.np/\">https://www.google.com.np/</a></p>\n"},
	{"23", "https://www.google.co.kr/", "<p><a href=\"https://www.google.co.kr/\">https://www.google.co.kr/</a></p>\n"},
	{"22", "https://aka.ms", "<p><a href=\"https://aka.ms\">https://aka.ms</a></p>\n"},
	{"21", "https://scoop.sh/", "<p><a href=\"https://scoop.sh/\">https://scoop.sh/</a></p>\n"},
	{"20", "https://www.electron.build/", "<p><a href=\"https://www.electron.build/\">https://www.electron.build/</a></p>\n"},
	{"19", "https://rime.im", "<p><a href=\"https://rime.im\">https://rime.im</a></p>\n"},
	{"18", "https://bbs.125.la", "<p><a href=\"https://bbs.125.la\">https://bbs.125.la</a></p>\n"},
	{"17", "https://www.ghisler.ch", "<p><a href=\"https://www.ghisler.ch\">https://www.ghisler.ch</a></p>\n"},
	{"16", "abc://xyz", "<p><a href=\"abc://xyz\">abc://xyz</a></p>\n"},
	{"15", "中https://notaurl文\n", "<p>中 https://notaurl 文</p>\n"},
	{"14", "abc://xyz测试foo", "<p><a href=\"abc://xyz%E6%B5%8B%E8%AF%95foo\">abc://xyz 测试 foo</a></p>\n"},
	{"13", "siyuan://blocks/20220817180757-c57m8qi测试foo", "<p><a href=\"siyuan://blocks/20220817180757-c57m8qi%E6%B5%8B%E8%AF%95foo\">siyuan://blocks/20220817180757-c57m8qi 测试 foo</a></p>\n"},
	{"12", "https://github.com/siyuan-note/siyuan/issues/?page=1&q=is%3Aissue+is%3Aopen", "<p><a href=\"https://github.com/siyuan-note/siyuan/issues/?page=1&amp;q=is%3Aissue+is%3Aopen\">https://github.com/siyuan-note/siyuan/issues/?page=1&amp;q=is%3Aissue+is%3Aopen</a></p>\n"},
	{"11", "https://github.com/88250/lute/issues/101", "<p><a href=\"https://github.com/88250/lute/issues/101\">Issue #101 · 88250/lute</a></p>\n"},
	{"10", "https://github.com/pages#标题\nhttps://www.google.com.hk/search?q=博客\nhttps://例子.网站/pages#home\n", "<p><a href=\"https://github.com/pages#%E6%A0%87%E9%A2%98\">https://github.com/pages#标题</a><br />\n<a href=\"https://www.google.com.hk/search?q=%E5%8D%9A%E5%AE%A2\">https://www.google.com.hk/search?q=博客</a><br />\nhttps://例子.网站/pages#home</p>\n"},
	{"9", "中http://notaurl文\n", "<p>中 http://notaurl 文</p>\n"},
	{"8", "1 www.noturl 2\n", "<p>1 www.noturl 2</p>\n"},
	{"7", "www.我的网址/console\n", "<p>www.我的网址/console</p>\n"},
	{"6", "http://我的网址/console\n", "<p>http://我的网址/console</p>\n"},
	{"5", "http://mydomain/console\n", "<p>http://mydomain/console</p>\n"},
	{"4", "http://foo.com/bar\n", "<p><a href=\"http://foo.com/bar\">http://foo.com/bar</a></p>\n"},
	{"3", "http://mydomain/console\n", "<p>http://mydomain/console</p>\n"},
	{"2", "www.非链接\n", "<p>www.非链接</p>\n"},
	{"1", "foo bar baz\n", "<p>foo bar baz</p>\n"},
	{"0", "foo http://bar.com baz\nfoo http://bar.com baz\n", "<p>foo <a href=\"http://bar.com\">http://bar.com</a> baz<br />\nfoo <a href=\"http://bar.com\">http://bar.com</a> baz</p>\n"},
}

func TestAutoLink(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetAutoSpace(true)
	for _, test := range autoLinkTests {
		result := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}
