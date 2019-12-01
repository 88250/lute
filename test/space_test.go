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

var spaceTests = []parseTest{

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
		html, err := luteEngine.MarkdownStr(test.name, test.from)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}
