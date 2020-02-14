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

var autoLinkTests = []parseTest{

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

	for _, test := range autoLinkTests {
		result, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("test case [%s] unexpected: %s", test.name, err)
		}

		if test.to != result {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal html\n\t%q", test.name, test.to, result, test.from)
		}
	}
}
