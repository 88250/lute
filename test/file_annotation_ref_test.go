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

var fileAnnotationRefTests = []parseTest{

	{"8", "<<assets/foo bar-20211115212742-80mbhnk.pdf/20211115213157-zifdhvu \"The phrase \\\"regular expression\\\" is often \">>", "<p>\"The phrase &quot;regular expression&quot; is often \"</p>\n"},
	{"7", "<<assets/foo bar.pdf/20210911230820-lhiaysx \"注解<锚文本>\">>", "<p>&lt;&lt;assets/foo bar.pdf/20210911230820-lhiaysx &quot;注解&lt;锚文本&gt;&quot;&gt;&gt;</p>\n"},
	{"6", "<<assets/foo bar.pdf/20210911230820-lhiaysx \"注解锚文本\">>", "<p>&lt;&lt;assets/foo bar.pdf/20210911230820-lhiaysx &quot;注解锚文本&quot;&gt;&gt;</p>\n"},
	{"5", "<<foo bar-20210911230735-pzlpdt.txt/20210911230820-lhiaysx \"注解锚文本\">>", "<p>&lt;&lt;foo bar-20210911230735-pzlpdt.txt/20210911230820-lhiaysx &quot;注解锚文本&quot;&gt;&gt;</p>\n"},
	{"4", "<<assets/foo bar-20210911230735-pzlpdt.txt/20210911230820-lhiaysx \"注解锚文本\">>", "<p>&lt;&lt;assets/foo bar-20210911230735-pzlpdt.txt/20210911230820-lhiaysx &quot;注解锚文本&quot;&gt;&gt;</p>\n"},
	{"3", "<<assets/foo bar-20210911230735-pzlpdt.pdf/20210911230820-lhiaysx \"注解锚文本\">>", "<p>&lt;&lt;assets/foo bar-20210911230735-pzlpdt.pdf/20210911230820-lhiaysx &quot;注解锚文本&quot;&gt;&gt;</p>\n"},
	{"2", "foo<<<bar>>>bazbazbazbazbazbazbazbazbazbazbazbazbazbazbazbaz", "<p>foo&lt;&lt;<bar>&gt;&gt;bazbazbazbazbazbazbazbazbazbazbazbazbazbazbazbaz</p>\n"},
	{"1", "<<assets/foo bar-20210911230735-pzlpdtf.pdf/20210911230820-lhiaysx \"注解锚文本\">>", "<p>\"注解锚文本\"</p>\n"},
	{"0", "<<assets/文件名-20210911230735-pzlpdtf.pdf/20210911230820-lhiaysx \"注解锚文本\">>", "<p>\"注解锚文本\"</p>\n"},
}

func TestFileAnnotationRef(t *testing.T) {
	luteEngine := lute.New()
	luteEngine.SetFileAnnotationRef(true)
	for _, test := range fileAnnotationRefTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
