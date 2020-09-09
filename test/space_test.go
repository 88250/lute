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

var spaceTests = []parseTest{

	// 井号 # 前后自动空格问题 https://github.com/88250/lute/issues/62
	{"36", "前#foo", "<p>前 #foo</p>\n"},
	{"35", "foo#后", "<p>foo# 后</p>\n"},
	{"34", "foo#bar", "<p>foo#bar</p>\n"},
	{"33", "前#后", "<p>前 # 后</p>\n"},

	// 加粗、强调和删除线自动空格改进 https://github.com/88250/lute/issues/25
	{"32", "数字~1链滴2~需要", "<p>数字 <del>1 链滴 2</del> 需要</p>\n"},
	{"32", "数字*1链滴2*需要", "<p>数字 <em>1 链滴 2</em> 需要</p>\n"},
	{"31", "中文 **链滴** 不需要", "<p>中文 <strong>链滴</strong> 不需要</p>\n"},
	{"30", "数字**1链滴2**需要", "<p>数字 <strong>1 链滴 2</strong> 需要</p>\n"},
	{"29", "英文1**HacPai**2需要", "<p>英文 1<strong>HacPai</strong>2 需要</p>\n"},
	{"28", "英文**HacPai**需要", "<p>英文 <strong>HacPai</strong> 需要</p>\n"},
	{"27", "中文**链滴**不需要", "<p>中文<strong>链滴</strong>不需要</p>\n"},
	{"26", "**链滴HacPai**需要", "<p><strong>链滴 HacPai</strong> 需要</p>\n"},

	// 链接前后自动空格改进 https://github.com/88250/lute/issues/24
	{"25", "中文 [链滴](https://ld246.com) 不需要", "<p>中文 <a href=\"https://ld246.com\">链滴</a> 不需要</p>\n"},
	{"24", "数字[1链滴2](https://ld246.com)需要", "<p>数字 <a href=\"https://ld246.com\">1 链滴 2</a> 需要</p>\n"},
	{"23", "英文1[HacPai](https://ld246.com)2需要", "<p>英文 1<a href=\"https://ld246.com\">HacPai</a>2 需要</p>\n"},
	{"22", "英文[HacPai](https://ld246.com)需要", "<p>英文 <a href=\"https://ld246.com\">HacPai</a> 需要</p>\n"},
	{"21", "中文[链滴](https://ld246.com)不需要", "<p>中文<a href=\"https://ld246.com\">链滴</a>不需要</p>\n"},
	{"20", "[链滴HacPai](https://ld246.com)需要", "<p><a href=\"https://ld246.com\">链滴 HacPai</a> 需要</p>\n"},

	{"19", "测试ping空格", "<p>测试 ping 空格</p>\n"},
	{"18", "foo❤️bar", "<p>foo❤️bar</p>\n"},

	// ing 前不需要空格，如 打码ing https://github.com/88250/lute/issues/9
	{"17", "打码ing开源", "<p>打码ing 开源</p>\n"},
	{"16", "打码in", "<p>打码 in</p>\n"},
	{"15", "打码ing", "<p>打码ing</p>\n"},

	{"14", "foo%bar", "<p>foo%bar</p>\n"},
	{"13", "**[链接foo文本](/bar)**\n", "<p><strong><a href=\"/bar\">链接 foo 文本</a></strong></p>\n"},
	{"12", "[链接foo文本](/bar)\n", "<p><a href=\"/bar\">链接 foo 文本</a></p>\n"},
	{"11", "[链接foo文本](/bar)\n", "<p><a href=\"/bar\">链接 foo 文本</a></p>\n"},
	{"10", "插入\u2038符\n", "<p>插入\u2038符</p>\n"},
	{"9", "非\u200c打印&zwnj;字符\n", "<p>非\u200c打印\u200c字符</p>\n"},
	{"8", "逗号，1后面\n", "<p>逗号，1 后面</p>\n"},
	{"7", "人民币符号后￥100不加空格\n", "<p>人民币符号后 ￥100 不加空格</p>\n"},
	{"6", "今日气温25℃晴\n", "<p>今日气温 25℃ 晴</p>\n"},

	// 自动空格 % 问题 https://github.com/b3log/lute/issues/28
	{"5", "百分号%前后需要空格\n", "<p>百分号 % 前后需要空格</p>\n"},
	{"4", "不错100%不错\n", "<p>不错 100% 不错</p>\n"},

	{"3", "是符号但不是$标^点|符号的要自动插入空格\n", "<p>是符号但不是 $ 标 ^ 点 | 符号的要自动插入空格</p>\n"},
	{"2", "(括[号{问号?等!标.点,符-号*要_排除掉\n", "<p>(括[号{问号?等!标.点,符-号*要_排除掉</p>\n"},
	{"1", "Lute解析200K的Markdown文本在我的电脑上只需要5ms。\n", "<p>Lute 解析 200K 的 Markdown 文本在我的电脑上只需要 5ms。</p>\n"},
	{"0", "Lute是一款结构化的Markdown引擎，完整实现了最新的GFM / CommonMark规范，对中文语境支持更好。\n", "<p>Lute 是一款结构化的 Markdown 引擎，完整实现了最新的 GFM / CommonMark 规范，对中文语境支持更好。</p>\n"},
}

func TestAutoSpace(t *testing.T) {
	luteEngine := lute.New() // 默认已经开启自动空格优化
	luteEngine.ChinesePunct = false

	for _, test := range spaceTests {
		html := luteEngine.MarkdownStr(test.name, test.from)
		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
