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

package benchmark

import (
	"github.com/b3log/lute"
	"io/ioutil"
	"testing"
)

func BenchmarkLute(b *testing.B) {
	spec := "../test/commonmark-spec"
	bytes, err := ioutil.ReadFile(spec + ".md")
	if nil != err {
		b.Fatalf("read spec text failed: " + err.Error())
	}

	luteEngine := lute.New(lute.GFM(false),
		lute.CodeSyntaxHighlight(false),
		lute.SoftBreak2HardBreak(false),
		lute.AutoSpace(false),
		lute.FixTermTypo(false),
	)
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
