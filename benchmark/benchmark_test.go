// Lute - 一款对中文语境优化的 Markdown 引擎，支持 Go 和 JavaScript
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

package main

import (
	"io/ioutil"
	"testing"

	"github.com/88250/lute"
)

const spec = "commonmark-spec"

func BenchmarkLute(b *testing.B) {
	bytes, err := ioutil.ReadFile(spec + ".md")
	if nil != err {
		b.Fatalf("read spec text failed: " + err.Error())
	}

	luteEngine := lute.New()
	luteEngine.GFMTaskListItem = true
	luteEngine.GFMTable = true
	luteEngine.GFMAutoLink = true
	luteEngine.GFMStrikethrough = true
	luteEngine.SoftBreak2HardBreak = false
	luteEngine.CodeSyntaxHighlight = false
	luteEngine.AutoSpace = false
	luteEngine.FixTermTypo = false
	luteEngine.ChinesePunct = false
	luteEngine.Emoji = false
	html, err := luteEngine.Markdown("spec text", bytes)
	if nil != err {
		b.Fatalf("unexpected: %s", err)
	}
	if err := ioutil.WriteFile(spec+".html", html, 0644); nil != err {
		b.Fatalf("write spec html failed: %s", err)
	}

	for i := 0; i < b.N; i++ {
		luteEngine.Markdown("spec text", bytes)
	}
}