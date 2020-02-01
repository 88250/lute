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
	"os"
	"runtime/pprof"

	"github.com/88250/lute"
)

func main() {
	spec := "test/commonmark-spec"
	bytes, err := ioutil.ReadFile(spec + ".md")
	if nil != err {
		panic(err)
	}

	luteEngine := lute.New()
	luteEngine.GFMTaskListItem = false
	luteEngine.GFMTable = false
	luteEngine.GFMAutoLink = false
	luteEngine.GFMStrikethrough = false
	luteEngine.SoftBreak2HardBreak = false
	luteEngine.CodeSyntaxHighlight = false
	luteEngine.Footnotes = false
	luteEngine.AutoSpace = false
	luteEngine.FixTermTypo = false
	luteEngine.ChinesePunct = false
	luteEngine.Emoji = false
	luteEngine.ParallelParsing = false

	cpuProfile, _ := os.Create("pprof/cpu_profile")
	pprof.StartCPUProfile(cpuProfile)
	for i := 0; i < 300; i++ {
		_, err := luteEngine.Markdown("pprof "+spec, bytes)
		if nil != err {
			panic(err)
		}
	}
	pprof.StopCPUProfile()
}
