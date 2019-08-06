// Lute - A structured markdown engine.
// Copyright (C) 2019-present, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/b3log/lute"
	"io/ioutil"
	"os"
	"runtime/pprof"
)

func main() {
	spec := "commonmark-0.29-spec"
	bytes, err := ioutil.ReadFile(spec + ".md")
	if nil != err {
		panic(err)
	}

	cpuProfile, _ := os.Create("cpu_profile")
	pprof.StartCPUProfile(cpuProfile)
	for i := 0; i < 40; i++ {
		tree, err := lute.Parse("spec text", bytes)
		if nil != err {
			panic(err)
		}

		renderer := lute.NewHTMLRenderer()
		if _, err := tree.Render(renderer); nil != err {
			panic(err)
		}
	}
	pprof.StopCPUProfile()
}
