// Lute - 一款结构化的 Markdown 引擎，支持 Go 和 JavaScript
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
//         http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package main

import (
	"os"
	"testing"

	"github.com/88250/lute"
)

const spec = "commonmark-spec"

func BenchmarkLute(b *testing.B) {
	buf, err := os.ReadFile(spec + ".md")
	if nil != err {
		b.Fatalf("read spec text failed: " + err.Error())
	}

	luteEngine := lute.New()
	luteEngine.SetGFMTaskListItem(true)
	luteEngine.SetGFMTable(true)
	luteEngine.SetGFMAutoLink(true)
	luteEngine.SetGFMStrikethrough(true)
	luteEngine.SetSoftBreak2HardBreak(false)
	luteEngine.SetCodeSyntaxHighlight(false)
	luteEngine.SetFootnotes(false)
	luteEngine.SetToC(false)
	luteEngine.SetHeadingID(false)
	luteEngine.SetAutoSpace(false)
	luteEngine.SetFixTermTypo(false)
	luteEngine.SetEmoji(false)
	luteEngine.SetYamlFrontMatter(false)
	output := luteEngine.Markdown("spec text", buf)
	if nil != err {
		b.Fatalf("unexpected: %s", err)
	}
	if err := os.WriteFile(spec+".html", output, 0644); nil != err {
		b.Fatalf("write spec html failed: %s", err)
	}

	b.SetParallelism(12)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			luteEngine.Markdown("spec text", buf)
		}
	})
}
